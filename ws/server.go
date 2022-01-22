package ws

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"path"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// ServerTimeoutConfig contains optional configuration parameters for a websocket server.
// Setting the parameter allows defining custom timeout intervals for websocket network operations.
//
// To set a custom configuration, refer to the server's SetTimeoutConfig method.
// If no configuration is passed, a default configuration is generated via the NewServerTimeoutConfig function.
type ServerTimeoutConfig struct {
	WriteWait time.Duration
	PingWait  time.Duration
}

// NewServerTimeoutConfig creates a default timeout configuration for a websocket endpoint.
//
// You may change fields arbitrarily and pass the struct to a SetTimeoutConfig method.
func NewServerTimeoutConfig() ServerTimeoutConfig {
	return ServerTimeoutConfig{WriteWait: defaultWriteWait, PingWait: defaultPingWait}
}

// WsServer defines a websocket server, which passively listens for incoming connections on ws or wss protocol.
// The offered API are of asynchronous nature, and each incoming connection/message is handled using callbacks.
//
// To create a new ws server, use:
//	server := NewServer()
//
// If you need a TLS ws server instead, use:
//	server := NewTLSServer("cert.pem", "privateKey.pem")
//
// To support client basic authentication, use:
//	server.SetBasicAuthHandler(func (user, pass) bool {
//		ok := authenticate(user, pass) // ... check for user and pass correctness
//		return ok
//	})
//
// To specify supported sub-protocols, use:
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
type WsServer interface {
	// Start and run the websocket server on a specific port and URL.
	// After start, incoming connections and messages are handled automatically, so no explicit read operation is required.
	//
	// The function blocks forever, hence it is suggested to invoke it in a goroutine, if the caller thread needs to perform other work, e.g.:
	//	go server.Start(8887, "/ws/{id}")
	//	doStuffOnMainThread()
	//	...
	//
	// To stop a running server, call the Stop function.
	Start(port int, listenPath string)
	// Stop shuts down a running websocket server.
	// All open channels will be forcefully closed, and the previously called Start function will return.
	Stop()
	// StopConnection closes a specific websocket connection.
	StopConnection(id string, closeError websocket.CloseError) error
	// Errors returns a channel for error messages. If it doesn't exist it es created.
	// The channel is closed by the server when stopped.
	Errors() <-chan error
	// SetMessageHandler sets a callback function for all incoming messages.
	// The callbacks accept a Channel and the received data.
	// It is up to the callback receiver, to check the identifier of the channel, to determine the source of the message.
	SetMessageHandler(handler func(ws Channel, data []byte) error)
	// SetNewClientHandler sets a callback function for all new incoming client connections.
	// It is recommended to store a reference to the Channel in the received entity, so that the Channel may be recognized later on.
	SetNewClientHandler(handler func(ws Channel))
	// SetDisconnectedClientHandler sets a callback function for all client disconnection events.
	// Once a client is disconnected, it is not possible to read/write on the respective Channel any longer.
	SetDisconnectedClientHandler(handler func(ws Channel))
	// SetTimeoutConfig set custom timeout configuration parameters. If not passed, a default
	// ServerTimeoutConfig struct will be used.
	//
	// This function must be called before starting the server, otherwise it may lead to unexpected behavior.
	SetTimeoutConfig(config ServerTimeoutConfig)
	// Write a message on a specific Channel, identified by the webSocketID parameter.
	// If the passed ID is invalid, an error is returned.
	//
	// The data is queued and will be sent asynchronously in the background.
	Write(webSocketID string, data []byte) error
	// AddSupportedSubprotocol adds support for a specified subprotocol.
	// This is recommended in order to communicate the capabilities to the client during the handshake.
	// If left empty, any subprotocol will be accepted.
	//
	// Duplicates will be removed automatically.
	AddSupportedSubprotocol(subProto string)
	// SetBasicAuthHandler enables HTTP Basic Authentication and requires clients to pass credentials.
	// The handler function is called whenever a new client attempts to connect, to check for credentials correctness.
	// The handler must return true if the credentials were correct, false otherwise.
	SetBasicAuthHandler(handler func(username string, password string) bool)
	// SetCheckOriginHandler sets a handler for incoming websocket connections, allowing to perform
	// custom cross-origin checks.
	//
	// By default, if the Origin header is present in the request, and the Origin host is not equal
	// to the Host request header, the websocket handshake fails.
	SetCheckOriginHandler(handler func(r *http.Request) bool)
	// Addr gives the address on which the server is listening, useful if, for
	// example, the port is system-defined (set to 0).
	Addr() *net.TCPAddr
}

// Server is the default implementation of a Websocket server.
//
// Use the NewServer or NewTLSServer functions to create a new server.
type Server struct {
	connections         map[string]*WebSocket
	httpServer          *http.Server
	messageHandler      func(ws Channel, data []byte) error
	newClientHandler    func(ws Channel)
	disconnectedHandler func(ws Channel)
	basicAuthHandler    func(username string, password string) bool
	tlsCertificatePath  string
	tlsCertificateKey   string
	timeoutConfig       ServerTimeoutConfig
	upgrader            websocket.Upgrader
	errC                chan error
	connMutex           sync.RWMutex
	addr                *net.TCPAddr
}

// NewServer creates a new simple websocket server (the websockets are not secured).
func NewServer() *Server {
	return &Server{
		httpServer:    &http.Server{},
		timeoutConfig: NewServerTimeoutConfig(),
		upgrader:      websocket.Upgrader{Subprotocols: []string{}},
	}
}

// NewTLSServer creates a new secure websocket server. All created websocket channels will use TLS.
//
// You need to pass a filepath to the server TLS certificate and key.
//
// It is recommended to pass a valid TLSConfig for the server to use.
// For example to require client certificate verification:
//	tlsConfig := &tls.Config{
//		ClientAuth: tls.RequireAndVerifyClientCert,
//		ClientCAs: clientCAs,
//	}
//
// If no tlsConfig parameter is passed, the server will by default
// not perform any client certificate verification.
func NewTLSServer(certificatePath string, certificateKey string, tlsConfig *tls.Config) *Server {
	return &Server{
		tlsCertificatePath: certificatePath,
		tlsCertificateKey:  certificateKey,
		httpServer: &http.Server{
			TLSConfig: tlsConfig,
		},
		timeoutConfig: NewServerTimeoutConfig(),
		upgrader:      websocket.Upgrader{Subprotocols: []string{}},
	}
}

func (server *Server) SetMessageHandler(handler func(ws Channel, data []byte) error) {
	server.messageHandler = handler
}

func (server *Server) SetNewClientHandler(handler func(ws Channel)) {
	server.newClientHandler = handler
}

func (server *Server) SetDisconnectedClientHandler(handler func(ws Channel)) {
	server.disconnectedHandler = handler
}

func (server *Server) SetTimeoutConfig(config ServerTimeoutConfig) {
	server.timeoutConfig = config
}

func (server *Server) AddSupportedSubprotocol(subProto string) {
	for _, sub := range server.upgrader.Subprotocols {
		if sub == subProto {
			// Don't add duplicates
			return
		}
	}
	server.upgrader.Subprotocols = append(server.upgrader.Subprotocols, subProto)
}

func (server *Server) SetBasicAuthHandler(handler func(username string, password string) bool) {
	server.basicAuthHandler = handler
}

func (server *Server) SetCheckOriginHandler(handler func(r *http.Request) bool) {
	server.upgrader.CheckOrigin = handler
}

func (server *Server) error(err error) {
	log.Error(err)
	if server.errC != nil {
		server.errC <- err
	}
}

func (server *Server) Errors() <-chan error {
	if server.errC == nil {
		server.errC = make(chan error, 1)
	}
	return server.errC
}

func (server *Server) Addr() *net.TCPAddr {
	return server.addr
}

func (server *Server) Start(port int, listenPath string) {
	router := mux.NewRouter()
	router.HandleFunc(listenPath, func(w http.ResponseWriter, r *http.Request) {
		server.wsHandler(w, r)
	})
	server.connections = make(map[string]*WebSocket)
	if server.httpServer == nil {
		server.httpServer = &http.Server{}
	}

	addr := fmt.Sprintf(":%v", port)
	server.httpServer.Addr = addr
	server.httpServer.Handler = router

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		server.error(fmt.Errorf("failed to listen: %w", err))
		return
	}

	server.addr = ln.Addr().(*net.TCPAddr)

	defer ln.Close()

	log.Infof("listening on tcp network %v", addr)
	server.httpServer.RegisterOnShutdown(server.stopConnections)
	if server.tlsCertificatePath != "" && server.tlsCertificateKey != "" {
		err = server.httpServer.ServeTLS(ln, server.tlsCertificatePath, server.tlsCertificateKey)
	} else {
		err = server.httpServer.Serve(ln)
	}

	if err != http.ErrServerClosed {
		server.error(fmt.Errorf("failed to listen: %w", err))
	}
}

func (server *Server) Stop() {
	log.Info("stopping websocket server")
	err := server.httpServer.Shutdown(context.TODO())
	if err != nil {
		server.error(fmt.Errorf("shutdown failed: %w", err))
	}

	if server.errC != nil {
		close(server.errC)
		server.errC = nil
	}
}

func (server *Server) StopConnection(id string, closeError websocket.CloseError) error {
	server.connMutex.RLock()
	ws, ok := server.connections[id]
	server.connMutex.RUnlock()

	if !ok {
		return fmt.Errorf("couldn't stop websocket connection. No connection with id %s is open", id)
	}
	log.Debugf("sending stop signal for websocket %s", ws.ID())
	ws.closeC <- closeError
	return nil
}

func (server *Server) stopConnections() {
	server.connMutex.Lock()
	defer server.connMutex.Unlock()
	for _, conn := range server.connections {
		conn.closeC <- websocket.CloseError{Code: websocket.CloseNormalClosure, Text: ""}
	}
}

func (server *Server) Write(webSocketId string, data []byte) error {
	server.connMutex.Lock()
	defer server.connMutex.Unlock()
	ws, ok := server.connections[webSocketId]
	if !ok {
		return fmt.Errorf("couldn't write to websocket. No socket with id %v is open", webSocketId)
	}
	log.Debugf("queuing data for websocket %s", webSocketId)
	ws.outQueue <- data
	return nil
}

func (server *Server) wsHandler(w http.ResponseWriter, r *http.Request) {
	responseHeader := http.Header{}
	url := r.URL
	id := path.Base(url.Path)
	log.Debugf("handling new connection for %s from %s", id, r.RemoteAddr)
	// Negotiate sub-protocol
	clientSubprotocols := websocket.Subprotocols(r)
	negotiatedSuprotocol := ""
out:
	for _, requestedProto := range clientSubprotocols {
		if len(server.upgrader.Subprotocols) == 0 {
			// All subProtocols are accepted, pick first
			negotiatedSuprotocol = requestedProto
			break
		}
		// Check if requested suprotocol is supported by server
		for _, supportedProto := range server.upgrader.Subprotocols {
			if requestedProto == supportedProto {
				negotiatedSuprotocol = requestedProto
				break out
			}
		}
	}
	if negotiatedSuprotocol != "" {
		responseHeader.Add("Sec-WebSocket-Protocol", negotiatedSuprotocol)
	}
	// Handle client authentication
	if server.basicAuthHandler != nil {
		username, password, ok := r.BasicAuth()
		ok = server.basicAuthHandler(username, password)
		if !ok {
			server.error(fmt.Errorf("basic auth failed: credentials invalid"))
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}
	// Upgrade websocket
	conn, err := server.upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		server.error(fmt.Errorf("upgrade failed: %w", err))
		return
	}

	// The id of the charge point is the final path element
	ws := WebSocket{
		connection:         conn,
		id:                 id,
		outQueue:           make(chan []byte, 1),
		closeC:             make(chan websocket.CloseError, 1),
		forceCloseC:        make(chan error, 1),
		pingMessage:        make(chan []byte, 1),
		tlsConnectionState: r.TLS,
	}
	log.Debugf("upgraded websocket connection for %s from %s", id, conn.RemoteAddr().String())
	// If unsupported subprotocol, terminate the connection immediately
	if negotiatedSuprotocol == "" {
		server.error(fmt.Errorf("unsupported subprotocols %v for new client %v (%v)", clientSubprotocols, id, r.RemoteAddr))
		_ = conn.WriteControl(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseProtocolError, "invalid or unsupported subprotocol"),
			time.Now().Add(server.timeoutConfig.WriteWait))
		_ = conn.Close()
		return
	}
	// Check whether client exists
	server.connMutex.Lock()
	// There is already a connection with the same ID. Close the new one immediately with a PolicyViolation.
	if _, exists := server.connections[id]; exists {
		server.connMutex.Unlock()
		server.error(fmt.Errorf("client %s already exists, closing duplicate client", id))
		_ = conn.WriteControl(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "a connection with this ID already exists"),
			time.Now().Add(server.timeoutConfig.WriteWait))
		_ = conn.Close()
		return
	}
	// Add new client
	server.connections[ws.id] = &ws
	server.connMutex.Unlock()
	// Read and write routines are started in separate goroutines and function will return immediately
	go server.writePump(&ws)
	go server.readPump(&ws)
	if server.newClientHandler != nil {
		var channel Channel = &ws
		server.newClientHandler(channel)
	}
}

func (server *Server) getReadTimeout() time.Time {
	if server.timeoutConfig.PingWait == 0 {
		return time.Time{}
	}
	return time.Now().Add(server.timeoutConfig.PingWait)
}

func (server *Server) readPump(ws *WebSocket) {
	conn := ws.connection

	conn.SetPingHandler(func(appData string) error {
		log.Debugf("ping received from %s", ws.ID())
		ws.pingMessage <- []byte(appData)
		err := conn.SetReadDeadline(server.getReadTimeout())
		return err
	})
	_ = conn.SetReadDeadline(server.getReadTimeout())

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure) {
				server.error(fmt.Errorf("read failed unexpectedly for %s: %w", ws.ID(), err))
			}
			log.Debugf("handling read error for %s: %v", ws.ID(), err.Error())
			// Notify writePump of error. Force close will be handled there
			ws.forceCloseC <- err
			return
		}

		if server.messageHandler != nil {
			var channel Channel = ws
			err = server.messageHandler(channel, message)
			if err != nil {
				server.error(fmt.Errorf("handling failed for %s: %w", ws.ID(), err))
				continue
			}
		}
		_ = conn.SetReadDeadline(server.getReadTimeout())
	}
}

func (server *Server) writePump(ws *WebSocket) {
	conn := ws.connection

	for {
		select {
		case data, ok := <-ws.outQueue:
			_ = conn.SetWriteDeadline(time.Now().Add(server.timeoutConfig.WriteWait))
			if !ok {
				// Unexpected closed queue, should never happen
				server.error(fmt.Errorf("output queue for socket %v was closed, forcefully closing", ws.id))
				// Don't invoke cleanup
				return
			}
			// Send data
			err := conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				server.error(fmt.Errorf("write failed for %s: %w", ws.ID(), err))
				// Invoking cleanup, as socket was forcefully closed
				server.cleanupConnection(ws)
				return
			}
			log.Debugf("written %d bytes to %s", len(data), ws.ID())
		case ping := <-ws.pingMessage:
			_ = conn.SetWriteDeadline(time.Now().Add(server.timeoutConfig.WriteWait))
			err := conn.WriteMessage(websocket.PongMessage, ping)
			if err != nil {
				server.error(fmt.Errorf("write failed for %s: %w", ws.ID(), err))
				// Invoking cleanup, as socket was forcefully closed
				server.cleanupConnection(ws)
				return
			}
			log.Debugf("pong sent to %s", ws.ID())
		case closeErr, _ := <-ws.closeC:
			log.Debugf("closing connection to %s", ws.ID())
			// Closing connection gracefully
			if err := conn.WriteControl(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(closeErr.Code, closeErr.Text),
				time.Now().Add(server.timeoutConfig.WriteWait),
			); err != nil {
				server.error(fmt.Errorf("failed to write close message for connection %s: %w", ws.id, err))
			}
			// Invoking cleanup
			server.cleanupConnection(ws)
			return
		case closed, ok := <-ws.forceCloseC:
			if !ok || closed != nil {
				// Connection was forcefully closed, invoke cleanup
				log.Debugf("handling forced close signal for %s", ws.ID())
				server.cleanupConnection(ws)
			}
			return
		}
	}
}

// Frees internal resources after a websocket connection was signaled to be closed.
// From this moment onwards, no new messages may be sent.
func (server *Server) cleanupConnection(ws *WebSocket) {
	_ = ws.connection.Close()
	server.connMutex.Lock()
	close(ws.outQueue)
	close(ws.closeC)
	delete(server.connections, ws.id)
	server.connMutex.Unlock()
	log.Infof("closed connection to %s", ws.ID())
	if server.disconnectedHandler != nil {
		server.disconnectedHandler(ws)
	}
}
