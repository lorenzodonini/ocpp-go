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

// WsServer defines a websocket server, which passively listens for incoming connections on ws or wss protocol.
// The offered API are of asynchronous nature, and each incoming connection/message is handled using callbacks.
//
// To create a new ws server, use:
//
//	server := NewServer()
//
// If you need a TLS ws server instead, use:
//
//	server := NewTLSServer("cert.pem", "privateKey.pem")
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
type WsServer interface {
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
}

// Default implementation of a Websocket server.
//
// Use the NewServer or NewTLSServer functions to create a new server.
type Server struct {
	connections         map[string]*webSocket
	httpServer          *http.Server
	messageHandler      func(ws Channel, data []byte) error
	checkClientHandler  CheckClientHandler
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
	httpHandler         *mux.Router
}

// Creates a new simple websocket server (the websockets are not secured).
func NewServer() *Server {
	router := mux.NewRouter()
	return &Server{
		httpServer:    &http.Server{},
		timeoutConfig: NewServerTimeoutConfig(),
		upgrader:      websocket.Upgrader{Subprotocols: []string{}},
		httpHandler:   router,
	}
}

// NewTLSServer creates a new secure websocket server. All created websocket channels will use TLS.
//
// You need to pass a filepath to the server TLS certificate and key.
//
// It is recommended to pass a valid TLSConfig for the server to use.
// For example to require client certificate verification:
//
//	tlsConfig := &tls.Config{
//		ClientAuth: tls.RequireAndVerifyClientCert,
//		ClientCAs: clientCAs,
//	}
//
// If no tlsConfig parameter is passed, the server will by default
// not perform any client certificate verification.
func NewTLSServer(certificatePath string, certificateKey string, tlsConfig *tls.Config) *Server {
	router := mux.NewRouter()
	return &Server{
		tlsCertificatePath: certificatePath,
		tlsCertificateKey:  certificateKey,
		httpServer: &http.Server{
			TLSConfig: tlsConfig,
		},
		timeoutConfig: NewServerTimeoutConfig(),
		upgrader:      websocket.Upgrader{Subprotocols: []string{}},
		httpHandler:   router,
	}
}

func (server *Server) SetMessageHandler(handler MessageHandler) {
	server.messageHandler = handler
}

func (server *Server) SetCheckClientHandler(handler CheckClientHandler) {
	server.checkClientHandler = handler
}

func (server *Server) SetNewClientHandler(handler ConnectedHandler) {
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

func (server *Server) AddHttpHandler(listenPath string, handler func(w http.ResponseWriter, r *http.Request)) {
	server.httpHandler.HandleFunc(listenPath, handler)
}

func (server *Server) Start(port int, listenPath string) {
	server.connMutex.Lock()
	server.connections = make(map[string]*webSocket)
	server.connMutex.Unlock()

	if server.httpServer == nil {
		server.httpServer = &http.Server{}
	}

	addr := fmt.Sprintf(":%v", port)
	server.httpServer.Addr = addr

	server.AddHttpHandler(listenPath, func(w http.ResponseWriter, r *http.Request) {
		server.wsHandler(w, r)
	})
	server.httpServer.Handler = server.httpHandler

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

	if !errors.Is(err, http.ErrServerClosed) {
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
	w, ok := server.connections[id]
	server.connMutex.RUnlock()

	if !ok {
		return fmt.Errorf("couldn't stop websocket connection. No connection with id %s is open", id)
	}
	log.Debugf("sending stop signal for websocket %s", w.ID())
	return w.Close(closeError)
}

func (server *Server) stopConnections() {
	server.connMutex.RLock()
	defer server.connMutex.RUnlock()
	for _, conn := range server.connections {
		_ = conn.Close(websocket.CloseError{Code: websocket.CloseNormalClosure, Text: ""})
	}
}

func (server *Server) Write(webSocketId string, data []byte) error {
	server.connMutex.RLock()
	defer server.connMutex.RUnlock()
	w, ok := server.connections[webSocketId]
	if !ok {
		return fmt.Errorf("couldn't write to websocket. No socket with id %v is open", webSocketId)
	}
	log.Debugf("queuing data for websocket %s", webSocketId)
	return w.Write(data)
}

func (server *Server) wsHandler(w http.ResponseWriter, r *http.Request) {
	responseHeader := http.Header{}
	url := r.URL
	id := path.Base(url.Path)
	log.Debugf("handling new connection for %s from %s", id, r.RemoteAddr)
	// Negotiate sub-protocol
	clientSubProtocols := websocket.Subprotocols(r)
	negotiatedSubProtocol := ""
out:
	for _, requestedProto := range clientSubProtocols {
		if len(server.upgrader.Subprotocols) == 0 {
			// All subProtocols are accepted, pick first
			negotiatedSubProtocol = requestedProto
			break
		}
		// Check if requested suprotocol is supported by server
		for _, supportedProto := range server.upgrader.Subprotocols {
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
	if server.basicAuthHandler != nil {
		username, password, ok := r.BasicAuth()
		if ok {
			ok = server.basicAuthHandler(username, password)
		}
		if !ok {
			server.error(fmt.Errorf("basic auth failed: credentials invalid"))
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}
	// Custom client checks
	if server.checkClientHandler != nil {
		ok := server.checkClientHandler(id, r)
		if !ok {
			server.error(fmt.Errorf("client validation: invalid client"))
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

	log.Debugf("upgraded websocket connection for %s from %s", id, conn.RemoteAddr().String())
	// If unsupported sub-protocol, terminate the connection immediately
	if negotiatedSubProtocol == "" {
		server.error(fmt.Errorf("unsupported subprotocols %v for new client %v (%v)", clientSubProtocols, id, r.RemoteAddr))
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
	// Create web socket for client, state is automatically set to connected
	ws := newWebSocket(
		id,
		conn,
		r.TLS,
		NewDefaultWebSocketConfig(
			false,
			server.timeoutConfig.WriteWait,
			server.timeoutConfig.PingWait,
			0,
			0),
		server.handleMessage,
		server.handleDisconnect,
		func(_ Channel, err error) {
			server.error(err)
		},
	)
	// Add new client
	server.connections[ws.id] = ws
	server.connMutex.Unlock()
	// Start reader and write routine
	ws.run()
	if server.newClientHandler != nil {
		var channel Channel = ws
		server.newClientHandler(channel)
	}
}

//func (server *Server) getReadTimeout() time.Time {
//	if server.timeoutConfig.PingWait == 0 {
//		return time.Time{}
//	}
//	return time.Now().Add(server.timeoutConfig.PingWait)
//}

//func (server *Server) readPump(ws *webSocket) {
//	conn := ws.connection
//
//	conn.SetPingHandler(func(appData string) error {
//		log.Debugf("ping received from %s", ws.ID())
//		ws.pingMessage <- []byte(appData)
//		err := conn.SetReadDeadline(server.getReadTimeout())
//		return err
//	})
//	_ = conn.SetReadDeadline(server.getReadTimeout())
//
//	for {
//		_, message, err := conn.ReadMessage()
//		if err != nil {
//			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure) {
//				server.error(fmt.Errorf("read failed unexpectedly for %s: %w", ws.ID(), err))
//			}
//			log.Debugf("handling read error for %s: %v", ws.ID(), err.Error())
//			// Notify writePump of error. Force close will be handled there
//			ws.forceCloseC <- err
//			return
//		}
//
//		if server.messageHandler != nil {
//			var channel Channel = ws
//			err = server.messageHandler(channel, message)
//			if err != nil {
//				server.error(fmt.Errorf("handling failed for %s: %w", ws.ID(), err))
//				continue
//			}
//		}
//		_ = conn.SetReadDeadline(server.getReadTimeout())
//	}
//}
//
//func (server *Server) writePump(ws *webSocket) {
//	conn := ws.connection
//
//	for {
//		select {
//		case data, ok := <-ws.outQueue:
//			_ = conn.SetWriteDeadline(time.Now().Add(server.timeoutConfig.WriteWait))
//			if !ok {
//				// Unexpected closed queue, should never happen
//				server.error(fmt.Errorf("output queue for socket %v was closed, forcefully closing", ws.id))
//				// Don't invoke cleanup
//				return
//			}
//			// Send data
//			err := conn.WriteMessage(websocket.TextMessage, data)
//			if err != nil {
//				server.error(fmt.Errorf("write failed for %s: %w", ws.ID(), err))
//				// Invoking cleanup, as socket was forcefully closed
//				server.cleanupConnection(ws)
//				return
//			}
//			log.Debugf("written %d bytes to %s", len(data), ws.ID())
//		case ping := <-ws.pingMessage:
//			_ = conn.SetWriteDeadline(time.Now().Add(server.timeoutConfig.WriteWait))
//			err := conn.WriteMessage(websocket.PongMessage, ping)
//			if err != nil {
//				server.error(fmt.Errorf("write failed for %s: %w", ws.ID(), err))
//				// Invoking cleanup, as socket was forcefully closed
//				server.cleanupConnection(ws)
//				return
//			}
//			log.Debugf("pong sent to %s", ws.ID())
//		case closeErr := <-ws.closeC:
//			log.Debugf("closing connection to %s", ws.ID())
//			// Closing connection gracefully
//			if err := conn.WriteControl(
//				websocket.CloseMessage,
//				websocket.FormatCloseMessage(closeErr.Code, closeErr.Text),
//				time.Now().Add(server.timeoutConfig.WriteWait),
//			); err != nil {
//				server.error(fmt.Errorf("failed to write close message for connection %s: %w", ws.id, err))
//			}
//			// Invoking cleanup
//			server.cleanupConnection(ws)
//			return
//		case closed, ok := <-ws.forceCloseC:
//			if !ok || closed != nil {
//				// Connection was forcefully closed, invoke cleanup
//				log.Debugf("handling forced close signal for %s", ws.ID())
//				server.cleanupConnection(ws)
//			}
//			return
//		}
//	}
//}

// --------- Internal callbacks webSocket -> Server ---------
func (server *Server) handleMessage(w Channel, data []byte) error {
	if server.messageHandler != nil {
		return server.messageHandler(w, data)
	}
	return fmt.Errorf("no message handler set")
}

func (server *Server) handleDisconnect(w Channel, _ error) {
	// Server never attempts to auto-reconnect to client. Resources are simply freed up
	server.connMutex.Lock()
	delete(server.connections, w.ID())
	server.connMutex.Unlock()
	log.Infof("closed connection to %s", w.ID())
	if server.disconnectedHandler != nil {
		server.disconnectedHandler(w)
	}
}
