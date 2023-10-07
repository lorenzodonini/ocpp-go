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

// Default implementation of a Websocket server.
//
// Use the NewServer or NewTLSServer functions to create a new server.
type Server struct {
	connections         map[string]*WebSocket
	httpServer          *http.Server
	messageHandler      func(ws Channel, data []byte) error
	checkClientHandler  func(id string, r *http.Request) bool
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

func (server *Server) SetMessageHandler(handler func(ws Channel, data []byte) error) {
	server.messageHandler = handler
}

func (server *Server) SetCheckClientHandler(handler func(id string, r *http.Request) bool) {
	server.checkClientHandler = handler
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

func (server *Server) AddHttpHandler(listenPath string, handler func(w http.ResponseWriter, r *http.Request)) {
	server.httpHandler.HandleFunc(listenPath, handler)
}

func (server *Server) Start(ln net.Listener, listenPath string) {
	server.connections = make(map[string]*WebSocket)
	if server.httpServer == nil {
		server.httpServer = &http.Server{}
	}

	server.AddHttpHandler(listenPath, func(w http.ResponseWriter, r *http.Request) {
		server.wsHandler(w, r)
	})
	server.httpServer.Handler = server.httpHandler

	server.addr = ln.Addr().(*net.TCPAddr)
	server.httpServer.Addr = fmt.Sprintf(":%d", server.addr.Port)

	log.Infof("listening on tcp network %v", server.httpServer.Addr)
	server.httpServer.RegisterOnShutdown(server.stopConnections)

	var err error
	if server.tlsCertificatePath != "" && server.tlsCertificateKey != "" {
		err = server.httpServer.ServeTLS(ln, server.tlsCertificatePath, server.tlsCertificateKey)
	} else {
		err = server.httpServer.Serve(ln)
	}

	if err != http.ErrServerClosed {
		server.error(fmt.Errorf("server failed: %w", err))
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
		case closeErr := <-ws.closeC:
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
