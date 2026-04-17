package multiplex

import (
	"github.com/lorenzodonini/ocpp-go/logging"
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

// server implements Server.
type server struct {
	wsServer ws.Server
}

// ServerOption is a function that configures a Server.
type ServerOption func(*server)

// NewServer creates a new multiplex server with a default ws.Server.
//
// Example usage:
//
//	// Create multiplex server
//	mux := multiplex.NewServer()
//
//	// Create only the OCPP servers you need
//	cs := ocpp16.NewCentralSystem(nil, mux.WebSocketServer())
//	csms := ocpp2.NewCSMS(nil, mux.WebSocketServer())
//
//	// Register handlers
//	cs.SetCoreHandler(&myOCPP16Handler{})
//	csms.SetProvisioningHandler(&myOCPP201Handler{})
//
//	// Start servers (ws.Server.Start is idempotent)
//	go cs.Start(8080, "/ocpp/{id}")
//	csms.Start(8080, "/ocpp/{id}")
func NewServer(opts ...ServerOption) Server {
	wsServer := ws.NewServer()
	return NewServerWithWebSocket(wsServer, opts...)
}

// NewServerWithWebSocket creates a Server using a custom ws.Server.
// This allows for custom TLS configuration or other WebSocket options.
//
// Example with TLS:
//
//	wsServer := ws.NewServer(ws.WithServerTLSConfig("cert.pem", "key.pem", nil))
//	mux := multiplex.NewServerWithWebSocket(wsServer)
func NewServerWithWebSocket(wsServer ws.Server, opts ...ServerOption) Server {
	if wsServer == nil {
		wsServer = ws.NewServer()
	}

	s := &server{
		wsServer: wsServer,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *server) WebSocketServer() ws.Server {
	return s.wsServer
}

func (s *server) SetSubprotocolSelector(selector SubprotocolSelector) {
	s.wsServer.SetSubprotocolSelector(func(id string, requestedSubprotocols []string) string {
		return selector(id, requestedSubprotocols)
	})
}

func (s *server) SetBasicAuthHandler(handler func(username, password string) bool) {
	s.wsServer.SetBasicAuthHandler(handler)
}

func (s *server) SetCheckClientHandler(handler ws.CheckClientHandler) {
	s.wsServer.SetCheckClientHandler(handler)
}
