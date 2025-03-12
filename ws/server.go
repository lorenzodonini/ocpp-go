package ws

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"path"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// ---------------------- SERVER ----------------------

type CheckClientHandler func(id string, r *http.Request) bool

// Server defines a websocket server, which passively listens for incoming connections on ws or wss protocol.
// The offered API are of asynchronous nature, and each incoming connection/message is handled using callbacks.
//
// To create a new ws server, use:
//
//	server := NewServer()
//
// If you need a server with TLS support, pass the following option:
//
//	server := NewServer(WithServerTLSConfig("cert.pem", "privateKey.pem", nil))
//
// To support client basic authentication, use:
//
//	server.SetBasicAuthHandler(func (user, pass) bool {
//		ok := authenticate(user, pass) // ... check for user and pass correctness
//		return ok
//	})
//
// To specify supported sub-protocols, use:
//
//	server.AddSupportedSubprotocol("ocpp1.6")
//
// If you need to set a specific timeout configuration, refer to the SetTimeoutConfig method.
//
// Using Start and Stop you can respectively start and stop listening for incoming client websocket connections.
//
// To be notified of new and terminated connections,
// refer to SetNewClientHandler and SetDisconnectedClientHandler functions.
//
// To receive incoming messages, you will need to set your own handler using SetMessageHandler.
// To write data on the open socket, simply call the Write function.
type Server interface {
	// Starts and runs the websocket server on a specific port and URL.
	// After start, incoming connections and messages are handled automatically, so no explicit read operation is required.
	//
	// The functions blocks forever, hence it is suggested to invoke it in a goroutine, if the caller thread needs to perform other work, e.g.:
	//	go server.Start(8887, "/ws/{id}")
	//	doStuffOnMainThread()
	//	...
	//
	// To stop a running server, call the Stop function.
	Start(port int, listenPath string)
	// Shuts down a running websocket server.
	// All open channels will be forcefully closed, and the previously called Start function will return.
	Stop()
	// Closes a specific websocket connection.
	StopConnection(id string, closeError websocket.CloseError) error
	// Errors returns a channel for error messages. If it doesn't exist it es created.
	// The channel is closed by the server when stopped.
	Errors() <-chan error
	// Sets a callback function for all incoming messages.
	// The callbacks accept a Channel and the received data.
	// It is up to the callback receiver, to check the identifier of the channel, to determine the source of the message.
	SetMessageHandler(handler MessageHandler)
	// SetNewClientHandler sets a callback function for all new incoming client connections.
	// It is recommended to store a reference to the Channel in the received entity, so that the Channel may be recognized later on.
	//
	// The callback is invoked after a connection was established and upgraded successfully.
	// If custom checks need to be run beforehand, refer to SetCheckClientHandler.
	SetNewClientHandler(handler ConnectedHandler)
	// Sets a callback function for all client disconnection events.
	// Once a client is disconnected, it is not possible to read/write on the respective Channel any longer.
	SetDisconnectedClientHandler(handler func(ws Channel))
	// Set custom timeout configuration parameters. If not passed, a default ServerTimeoutConfig struct will be used.
	//
	// This function must be called before starting the server, otherwise it may lead to unexpected behavior.
	SetTimeoutConfig(config ServerTimeoutConfig)
	// Write sends a message on a specific Channel, identifier by the webSocketId parameter.
	// If the passed ID is invalid, an error is returned.
	//
	// The data is queued and will be sent asynchronously in the background.
	Write(webSocketId string, data []byte) error
	// AddSupportedSubprotocol adds support for a specified subprotocol.
	// This is recommended in order to communicate the capabilities to the client during the handshake.
	// If left empty, any subprotocol will be accepted.
	//
	// Duplicates will be removed automatically.
	AddSupportedSubprotocol(subProto string)
	// SetChargePointIdResolver sets the callback function to use for resolving the charge point ID of a charger connecting to
	// the websocket server. By default, this will just be the path in the URL used by the client.
	SetChargePointIdResolver(resolver func(r *http.Request) (string, error))
	// SetBasicAuthHandler enables HTTP Basic Authentication and requires clients to pass credentials.
	// The handler function is called whenever a new client attempts to connect, to check for credentials correctness.
	// The handler must return true if the credentials were correct, false otherwise.
	SetBasicAuthHandler(handler func(chargePointID string, username string, password string) bool)
	// SetCheckOriginHandler sets a handler for incoming websocket connections, allowing to perform
	// custom cross-origin checks.
	//
	// By default, if the Origin header is present in the request, and the Origin host is not equal
	// to the Host request header, the websocket handshake fails.
	SetCheckOriginHandler(handler func(r *http.Request) bool)
	// SetCheckClientHandler sets a handler for validate incoming websocket connections, allowing to perform
	// custom client connection checks.
	// The handler is executed before any connection upgrade and allows optionally returning a custom
	// configuration for the web socket that will be created.
	//
	// Changes to the http request at runtime may lead to undefined behavior.
	SetCheckClientHandler(handler CheckClientHandler)
	// Addr gives the address on which the server is listening, useful if, for
	// example, the port is system-defined (set to 0).
	Addr() *net.TCPAddr
	// GetChannel retrieves an active Channel connection by its unique identifier.
	// If a connection with the given ID exists, it returns the corresponding webSocket instance.
	// If no connection is found with the specified ID, it returns nil and a false flag.
	GetChannel(websocketId string) (Channel, bool)
}

// Default implementation of a Websocket server.
//
// Use the NewServer function to create a new server.
type server struct {
	connections           map[string]*webSocket
	httpServer            *http.Server
	messageHandler        func(ws Channel, data []byte) error
	chargePointIdResolver func(*http.Request) (string, error)
	checkClientHandler    CheckClientHandler
	newClientHandler      func(ws Channel)
	disconnectedHandler   func(ws Channel)
	basicAuthHandler      func(chargePointID string, username string, password string) bool
	tlsCertificatePath    string
	tlsCertificateKey     string
	timeoutConfig         ServerTimeoutConfig
	upgrader              websocket.Upgrader
	errC                  chan error
	connMutex             sync.RWMutex
	addr                  *net.TCPAddr
	httpHandler           *mux.Router
}

// ServerOpt is a function that can be used to set options on a server during creation.
type ServerOpt func(s *server)

// WithServerTLSConfig sets the TLS configuration for the server.
// If the passed tlsConfig is nil, the client will not use TLS.
func WithServerTLSConfig(certificatePath string, certificateKey string, tlsConfig *tls.Config) ServerOpt {
	return func(s *server) {
		s.tlsCertificatePath = certificatePath
		s.tlsCertificateKey = certificateKey
		if tlsConfig != nil {
			s.httpServer.TLSConfig = tlsConfig
		}
	}
}

// NewServer Creates a new websocket server.
//
// Additional options may be added using the AddOption function.
//
// By default, the websockets are not secure, and the server will not perform any client certificate verification.
//
// To add TLS support to the server, a valid server certificate path and key must be passed.
// To also add support for client certificate verification, a valid TLSConfig needs to be configured.
// For example:
//
//		tlsConfig := &tls.Config{
//			ClientAuth: tls.RequireAndVerifyClientCert,
//			ClientCAs: clientCAs,
//		}
//	 server := ws.NewServer(ws.WithServerTLSConfig("cert.pem", "privateKey.pem", tlsConfig))
//
// When TLS is correctly configured, the server will automatically use it for all created websocket channels.
func NewServer(opts ...ServerOpt) Server {
	router := mux.NewRouter()
	s := &server{
		httpServer:    &http.Server{},
		timeoutConfig: NewServerTimeoutConfig(),
		upgrader:      websocket.Upgrader{Subprotocols: []string{}},
		httpHandler:   router,
		chargePointIdResolver: func(r *http.Request) (string, error) {
			url := r.URL
			return path.Base(url.Path), nil
		},
	}
	for _, o := range opts {
		o(s)
	}
	return s
}

func (s *server) SetMessageHandler(handler MessageHandler) {
	s.messageHandler = handler
}

func (s *server) SetCheckClientHandler(handler CheckClientHandler) {
	s.checkClientHandler = handler
}

func (s *server) SetNewClientHandler(handler ConnectedHandler) {
	s.newClientHandler = handler
}

func (s *server) SetDisconnectedClientHandler(handler func(ws Channel)) {
	s.disconnectedHandler = handler
}

func (s *server) SetTimeoutConfig(config ServerTimeoutConfig) {
	s.timeoutConfig = config
}

func (s *server) AddSupportedSubprotocol(subProto string) {
	for _, sub := range s.upgrader.Subprotocols {
		if sub == subProto {
			// Don't add duplicates
			return
		}
	}
	s.upgrader.Subprotocols = append(s.upgrader.Subprotocols, subProto)
}

func (s *server) SetChargePointIdResolver(resolver func(r *http.Request) (string, error)) {
	s.chargePointIdResolver = resolver
}

func (s *server) SetBasicAuthHandler(handler func(chargePointID string, username string, password string) bool) {
	s.basicAuthHandler = handler
}

func (s *server) SetCheckOriginHandler(handler func(r *http.Request) bool) {
	s.upgrader.CheckOrigin = handler
}

func (s *server) error(err error) {
	log.Error(err)
	if s.errC != nil {
		s.errC <- err
	}
}

func (s *server) Errors() <-chan error {
	if s.errC == nil {
		s.errC = make(chan error, 1)
	}
	return s.errC
}

func (s *server) Addr() *net.TCPAddr {
	return s.addr
}

func (s *server) AddHttpHandler(listenPath string, handler func(w http.ResponseWriter, r *http.Request)) {
	s.httpHandler.HandleFunc(listenPath, handler)
}

func (s *server) Start(port int, listenPath string) {
	s.connMutex.Lock()
	s.connections = make(map[string]*webSocket)
	s.connMutex.Unlock()

	if s.httpServer == nil {
		s.httpServer = &http.Server{}
	}

	addr := fmt.Sprintf(":%v", port)
	s.httpServer.Addr = addr

	s.AddHttpHandler(listenPath, func(w http.ResponseWriter, r *http.Request) {
		s.wsHandler(w, r)
	})
	s.httpServer.Handler = s.httpHandler

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		s.error(fmt.Errorf("failed to listen: %w", err))
		return
	}

	s.addr = ln.Addr().(*net.TCPAddr)

	defer ln.Close()

	log.Infof("listening on tcp network %v", addr)
	s.httpServer.RegisterOnShutdown(s.stopConnections)
	if s.tlsCertificatePath != "" && s.tlsCertificateKey != "" {
		err = s.httpServer.ServeTLS(ln, s.tlsCertificatePath, s.tlsCertificateKey)
	} else {
		err = s.httpServer.Serve(ln)
	}

	if !errors.Is(err, http.ErrServerClosed) {
		s.error(fmt.Errorf("failed to listen: %w", err))
	}
}

func (s *server) Stop() {
	log.Info("stopping websocket server")
	err := s.httpServer.Shutdown(context.TODO())
	if err != nil {
		s.error(fmt.Errorf("shutdown failed: %w", err))
	}

	if s.errC != nil {
		close(s.errC)
		s.errC = nil
	}
}

func (s *server) StopConnection(id string, closeError websocket.CloseError) error {
	s.connMutex.RLock()
	w, ok := s.connections[id]
	s.connMutex.RUnlock()

	if !ok {
		return fmt.Errorf("couldn't stop websocket connection. No connection with id %s is open", id)
	}
	log.Debugf("sending stop signal for websocket %s", w.ID())
	return w.Close(closeError)
}

func (s *server) GetChannel(websocketId string) (Channel, bool) {
	s.connMutex.RLock()
	defer s.connMutex.RUnlock()
	c, ok := s.connections[websocketId]
	return c, ok
}

func (s *server) stopConnections() {
	s.connMutex.RLock()
	defer s.connMutex.RUnlock()
	for _, conn := range s.connections {
		_ = conn.Close(websocket.CloseError{Code: websocket.CloseNormalClosure, Text: ""})
	}
}

func (s *server) Write(webSocketId string, data []byte) error {
	s.connMutex.RLock()
	defer s.connMutex.RUnlock()
	w, ok := s.connections[webSocketId]
	if !ok {
		return fmt.Errorf("couldn't write to websocket. No socket with id %v is open", webSocketId)
	}
	log.Debugf("queuing data for websocket %s", webSocketId)
	return w.Write(data)
}

func (s *server) wsHandler(w http.ResponseWriter, r *http.Request) {
	responseHeader := http.Header{}
	id, err := s.chargePointIdResolver(r)
	if err != nil {
		s.error(fmt.Errorf("failed to resolve charge point id"))
		http.Error(w, "NotFound", http.StatusNotFound)
		return
	}
	log.Debugf("handling new connection for %s from %s", id, r.RemoteAddr)
	// Negotiate sub-protocol
	clientSubProtocols := websocket.Subprotocols(r)
	negotiatedSubProtocol := ""
out:
	for _, requestedProto := range clientSubProtocols {
		if len(s.upgrader.Subprotocols) == 0 {
			// All subProtocols are accepted, pick first
			negotiatedSubProtocol = requestedProto
			break
		}
		// Check if requested suprotocol is supported by server
		for _, supportedProto := range s.upgrader.Subprotocols {
			if requestedProto == supportedProto {
				negotiatedSubProtocol = requestedProto
				break out
			}
		}
	}
	if negotiatedSubProtocol != "" {
		responseHeader.Add("Sec-WebSocket-Protocol", negotiatedSubProtocol)
	}
	// Handle client authentication
	if s.basicAuthHandler != nil {
		username, password, ok := r.BasicAuth()
		if ok {
			ok = s.basicAuthHandler(id, username, password)
		}
		if !ok {
			s.error(fmt.Errorf("basic auth failed: credentials invalid"))
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}
	// Custom client checks
	if s.checkClientHandler != nil {
		ok := s.checkClientHandler(id, r)
		if !ok {
			s.error(fmt.Errorf("client validation: invalid client"))
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	// Upgrade websocket
	conn, err := s.upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		s.error(fmt.Errorf("upgrade failed: %w", err))
		return
	}

	log.Debugf("upgraded websocket connection for %s from %s", id, conn.RemoteAddr().String())
	// If unsupported sub-protocol, terminate the connection immediately
	if negotiatedSubProtocol == "" {
		s.error(fmt.Errorf("unsupported subprotocols %v for new client %v (%v)", clientSubProtocols, id, r.RemoteAddr))
		_ = conn.WriteControl(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseProtocolError, "invalid or unsupported subprotocol"),
			time.Now().Add(s.timeoutConfig.WriteWait))
		_ = conn.Close()
		return
	}
	// Check whether client exists
	s.connMutex.Lock()
	// There is already a connection with the same ID. Close the new one immediately with a PolicyViolation.
	if _, exists := s.connections[id]; exists {
		s.connMutex.Unlock()
		s.error(fmt.Errorf("client %s already exists, closing duplicate client", id))
		_ = conn.WriteControl(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "a connection with this ID already exists"),
			time.Now().Add(s.timeoutConfig.WriteWait))
		_ = conn.Close()
		return
	}
	// Create web socket for client, state is automatically set to connected
	ws := newWebSocket(
		id,
		conn,
		r.TLS,
		NewDefaultWebSocketConfig(
			s.timeoutConfig.WriteWait,
			s.timeoutConfig.PingWait,
			s.timeoutConfig.PingPeriod,
			s.timeoutConfig.PongWait),
		s.handleMessage,
		s.handleDisconnect,
		func(_ Channel, err error) {
			s.error(err)
		},
	)
	// Add new client
	s.connections[ws.id] = ws
	s.connMutex.Unlock()
	// Start reader and write routine
	ws.run()
	if s.newClientHandler != nil {
		var channel Channel = ws
		s.newClientHandler(channel)
	}
}

// --------- Internal callbacks webSocket -> server ---------
func (s *server) handleMessage(w Channel, data []byte) error {
	if s.messageHandler != nil {
		return s.messageHandler(w, data)
	}
	return fmt.Errorf("no message handler set")
}

func (s *server) handleDisconnect(w Channel, _ error) {
	// server never attempts to auto-reconnect to client. Resources are simply freed up
	s.connMutex.Lock()
	delete(s.connections, w.ID())
	s.connMutex.Unlock()
	log.Infof("closed connection to %s", w.ID())
	if s.disconnectedHandler != nil {
		s.disconnectedHandler(w)
	}
}
