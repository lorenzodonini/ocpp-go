package multiplex

import (
	"fmt"
	"sync"

	"github.com/lorenzodonini/ocpp-go/logging"
	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	ocpp2 "github.com/lorenzodonini/ocpp-go/ocpp2.0.1"
	"github.com/lorenzodonini/ocpp-go/ws"
)

var log logging.Logger

func init() {
	log = &logging.VoidLogger{}
}

// SetLogger sets a custom logger for the multiplex package.
func SetLogger(logger logging.Logger) {
	if logger == nil {
		panic("cannot set nil logger")
	}
	log = logger
}

// multiProtocolServer implements MultiProtocolServer.
type multiProtocolServer struct {
	wsServer             ws.Server
	ocpp16Server         ocpp16.CentralSystem
	ocpp201Server        ocpp2.CSMS
	clientProtocols      map[string]ProtocolVersion
	clientProtocolsMutex sync.RWMutex
}

// ServerOption is a function that configures a MultiProtocolServer.
type ServerOption func(*multiProtocolServer)

// NewMultiProtocolServer creates a new server that handles both OCPP 1.6 and OCPP 2.0.1.
//
// The server registers both "ocpp1.6" and "ocpp2.0.1" as supported WebSocket subprotocols.
// When a client connects with a Sec-WebSocket-Protocol header (e.g., "ocpp2.0.1, ocpp1.6"),
// the server negotiates the first mutually-supported protocol and routes the connection
// to the appropriate handler.
//
// Example usage:
//
//	server := multiplex.NewMultiProtocolServer()
//
//	// Register OCPP 1.6 handlers
//	server.OCPP16Server().SetCoreHandler(&myOCPP16Handler{})
//
//	// Register OCPP 2.0.1 handlers
//	server.OCPP201Server().SetProvisioningHandler(&myOCPP201Handler{})
//
//	// Start on single port
//	server.Start(8080, "/ocpp/{id}")
func NewMultiProtocolServer(opts ...ServerOption) MultiProtocolServer {
	wsServer := ws.NewServer()
	return newMultiProtocolServerWithWebSocket(wsServer, opts...)
}

// NewMultiProtocolServerWithWebSocket creates a MultiProtocolServer using a custom ws.Server.
// This allows for custom TLS configuration or other WebSocket options.
//
// Example with TLS:
//
//	wsServer := ws.NewServer(ws.WithServerTLSConfig("cert.pem", "key.pem", nil))
//	server := multiplex.NewMultiProtocolServerWithWebSocket(wsServer)
func NewMultiProtocolServerWithWebSocket(wsServer ws.Server, opts ...ServerOption) MultiProtocolServer {
	if wsServer == nil {
		wsServer = ws.NewServer()
	}
	return newMultiProtocolServerWithWebSocket(wsServer, opts...)
}

func newMultiProtocolServerWithWebSocket(wsServer ws.Server, opts ...ServerOption) *multiProtocolServer {
	// Create the OCPP servers - they will register their subprotocols and handlers
	ocpp16Server := ocpp16.NewCentralSystem(nil, wsServer)
	ocpp201Server := ocpp2.NewCSMS(nil, wsServer)

	s := &multiProtocolServer{
		wsServer:        wsServer,
		ocpp16Server:    ocpp16Server,
		ocpp201Server:   ocpp201Server,
		clientProtocols: make(map[string]ProtocolVersion),
	}

	// Apply options
	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *multiProtocolServer) Start(port int, listenPath string) {
	log.Infof("starting multi-protocol OCPP server on port %d", port)

	// Start both OCPP servers - ws.Server.Start() is idempotent,
	// so only the first call actually starts the server
	go s.ocpp16Server.Start(port, listenPath)
	s.ocpp201Server.Start(port, listenPath)
}

func (s *multiProtocolServer) Stop() {
	log.Info("stopping multi-protocol OCPP server")
	// Both servers share the same ws.Server, so stopping either stops both
	s.ocpp16Server.Stop()
}

func (s *multiProtocolServer) OCPP16Server() ocpp16.CentralSystem {
	return s.ocpp16Server
}

func (s *multiProtocolServer) OCPP201Server() ocpp2.CSMS {
	return s.ocpp201Server
}

func (s *multiProtocolServer) SetSubprotocolSelector(selector SubprotocolSelector) {
	s.wsServer.SetSubprotocolSelector(func(id string, requestedSubprotocols []string) string {
		return selector(id, requestedSubprotocols)
	})
}

func (s *multiProtocolServer) SetNewClientHandler(handler func(ws.Channel)) {
	s.ocpp16Server.SetNewChargePointHandler(func(cp ocpp16.ChargePointConnection) {
		if ch, ok := s.wsServer.GetChannel(cp.ID()); ok {
			s.trackNewClient(ch)
			if handler != nil {
				handler(ch)
			}
		}
	})
	s.ocpp201Server.SetNewChargingStationHandler(func(cs ocpp2.ChargingStationConnection) {
		if ch, ok := s.wsServer.GetChannel(cs.ID()); ok {
			s.trackNewClient(ch)
			if handler != nil {
				handler(ch)
			}
		}
	})
}

func (s *multiProtocolServer) SetDisconnectedClientHandler(handler func(ws.Channel)) {
	s.ocpp16Server.SetChargePointDisconnectedHandler(func(cp ocpp16.ChargePointConnection) {
		if ch, ok := s.wsServer.GetChannel(cp.ID()); ok {
			s.cleanupClient(ch)
			if handler != nil {
				handler(ch)
			}
		}
	})
	s.ocpp201Server.SetChargingStationDisconnectedHandler(func(cs ocpp2.ChargingStationConnection) {
		if ch, ok := s.wsServer.GetChannel(cs.ID()); ok {
			s.cleanupClient(ch)
			if handler != nil {
				handler(ch)
			}
		}
	})
}

func (s *multiProtocolServer) SetBasicAuthHandler(handler func(username, password string) bool) {
	s.wsServer.SetBasicAuthHandler(handler)
}

func (s *multiProtocolServer) SetCheckClientHandler(handler ws.CheckClientHandler) {
	s.wsServer.SetCheckClientHandler(handler)
}

// trackNewClient records the protocol version for a newly connected client.
func (s *multiProtocolServer) trackNewClient(channel ws.Channel) {
	clientID := channel.ID()
	subprotocol := channel.Subprotocol()

	var protocolVersion ProtocolVersion
	switch subprotocol {
	case V16Subprotocol:
		protocolVersion = V16
	case V201Subprotocol:
		protocolVersion = V201
	default:
		log.Errorf("unknown subprotocol %s for client %s", subprotocol, clientID)
		return
	}

	s.clientProtocolsMutex.Lock()
	s.clientProtocols[clientID] = protocolVersion
	s.clientProtocolsMutex.Unlock()

	log.Infof("client %s connected with protocol %s", clientID, protocolVersion)
}

// cleanupClient removes protocol tracking for a disconnected client.
func (s *multiProtocolServer) cleanupClient(channel ws.Channel) {
	clientID := channel.ID()

	s.clientProtocolsMutex.Lock()
	delete(s.clientProtocols, clientID)
	s.clientProtocolsMutex.Unlock()

	log.Infof("client %s disconnected", clientID)
}

// GetClientProtocol returns the protocol version for a connected client.
func (s *multiProtocolServer) GetClientProtocol(clientID string) (ProtocolVersion, bool) {
	s.clientProtocolsMutex.RLock()
	defer s.clientProtocolsMutex.RUnlock()
	v, ok := s.clientProtocols[clientID]
	return v, ok
}

// GetChannel retrieves the WebSocket channel for a client by ID.
func (s *multiProtocolServer) GetChannel(clientID string) (ws.Channel, bool) {
	return s.wsServer.GetChannel(clientID)
}

// Write sends data to a specific client.
func (s *multiProtocolServer) Write(clientID string, data []byte) error {
	return s.wsServer.Write(clientID, data)
}

// Errors returns a channel for error messages from the WebSocket server.
func (s *multiProtocolServer) Errors() <-chan error {
	return s.wsServer.Errors()
}

// IsClientConnected checks if a client with the given ID is currently connected.
func (s *multiProtocolServer) IsClientConnected(clientID string) bool {
	_, ok := s.wsServer.GetChannel(clientID)
	return ok
}

// GetClientSubprotocol returns the negotiated subprotocol for a client.
func (s *multiProtocolServer) GetClientSubprotocol(clientID string) (string, error) {
	channel, ok := s.wsServer.GetChannel(clientID)
	if !ok {
		return "", fmt.Errorf("client %s not found", clientID)
	}
	return channel.Subprotocol(), nil
}
