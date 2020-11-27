// The package is a wrapper around gorilla websockets,
// aimed at simplifying the creation and usage of a websocket client/server.
//
// Check the Client and Server structure to get started.
package ws

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	defaultWriteWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	defaultPongWait = 60 * time.Second
	// Time allowed to wait for a ping on the server, before closing a connection due to inactivity.
	defaultPingWait = defaultPongWait
	// Send pings to peer with this period. Must be less than pongWait.
	defaultPingPeriod = (defaultPongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 512
	// Time allowed for the initial handshake to complete.
	defaultHandshakeTimeout = 30 * time.Second
	// Default sub-protocol to send to peer upon connection.
	defaultSubProtocol = "ocpp1.6"
	// The base delay to be used for automatic reconnection. Will increase exponentially up to maxReconnectionDelay.
	defaultAutoReconnectDelay = 5 * time.Second
	// Default maximum reconnection delay for websockets
	defaultMaxReconnectionDelay = 2 * time.Minute
)

// Config contains optional configuration parameters for a websocket server.
// Setting the parameter allows to define custom timeout intervals for websocket network operations.
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

// Config contains optional configuration parameters for a websocket client.
// Setting the parameter allows to define custom timeout intervals for websocket network operations.
//
// To set a custom configuration, refer to the client's SetTimeoutConfig method.
// If no configuration is passed, a default configuration is generated via the NewClientTimeoutConfig function.
type ClientTimeoutConfig struct {
	WriteWait             time.Duration
	HandshakeTimeout      time.Duration
	PongWait              time.Duration
	PingPeriod            time.Duration
	AutoReconnectDelay    time.Duration
	MaxAutoReconnectDelay time.Duration
}

// NewClientTimeoutConfig creates a default timeout configuration for a websocket endpoint.
//
// You may change fields arbitrarily and pass the struct to a SetTimeoutConfig method.
func NewClientTimeoutConfig() ClientTimeoutConfig {
	return ClientTimeoutConfig{
		WriteWait:             defaultWriteWait,
		HandshakeTimeout:      defaultHandshakeTimeout,
		PongWait:              defaultPongWait,
		PingPeriod:            defaultPingPeriod,
		AutoReconnectDelay:    defaultAutoReconnectDelay,
		MaxAutoReconnectDelay: defaultMaxReconnectionDelay,
	}
}

// Channel represents a bi-directional communication channel, which provides at least a unique ID.
type Channel interface {
	GetID() string
}

// WebSocket is a wrapper for a single websocket channel.
// The connection itself is provided by the gorilla websocket package.
//
// Don't use a websocket directly, but refer to WsServer and WsClient.
type WebSocket struct {
	connection  *websocket.Conn
	id          string
	outQueue    chan []byte
	closeSignal chan error
	pingMessage chan []byte
}

// Retrieves the unique Identifier of the websocket (typically, the URL suffix).
func (websocket *WebSocket) GetID() string {
	return websocket.id
}

// ConnectionError is a websocket
type HttpConnectionError struct {
	Message    string
	HttpStatus string
	HttpCode   int
	Details    string
}

func (e HttpConnectionError) Error() string {
	return fmt.Sprintf("%v, http status: %v", e.Message, e.HttpStatus)
}

// ---------------------- SERVER ----------------------

// A Websocket server, which passively listens for incoming connections on ws or wss protocol.
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
	// Errors returns a channel for error messages. If it doesn't exist it es created.
	// The channel is closed by the server when stopped.
	Errors() <-chan error
	// Sets a callback function for all incoming messages.
	// The callbacks accept a Channel and the received data.
	// It is up to the callback receiver, to check the identifier of the channel, to determine the source of the message.
	SetMessageHandler(handler func(ws Channel, data []byte) error)
	// Sets a callback function for all new incoming client connections.
	// It is recommended to store a reference to the Channel in the received entity, so that the Channel may be recognized later on.
	SetNewClientHandler(handler func(ws Channel))
	// Sets a callback function for all client disconnection events.
	// Once a client is disconnected, it is not possible to read/write on the respective Channel any longer.
	SetDisconnectedClientHandler(handler func(ws Channel))
	// Set custom timeout configuration parameters. If not passed, a default ServerTimeoutConfig struct will be used.
	//
	// This function must be called before starting the server, otherwise it may lead to unexpected behavior.
	SetTimeoutConfig(config ServerTimeoutConfig)
	// Sends a message on a specific Channel, identifier by the webSocketId parameter.
	// If the passed ID is invalid, an error is returned.
	//
	// The data is queued and will be sent asynchronously in the background.
	Write(webSocketId string, data []byte) error
	// Adds support for a specified subprotocol.
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
}

// Default implementation of a Websocket server.
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
}

// Creates a new simple websocket server (the websockets are not secured).
func NewServer() *Server {
	return &Server{timeoutConfig: NewServerTimeoutConfig(), upgrader: websocket.Upgrader{Subprotocols: []string{}}}
}

// Creates a new secure websocket server. All created websocket channels will use TLS.
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
	if server.tlsCertificatePath != "" && server.tlsCertificateKey != "" {
		if err := server.httpServer.ListenAndServeTLS(server.tlsCertificatePath, server.tlsCertificateKey); err != http.ErrServerClosed {
			server.error(fmt.Errorf("failed to listen: %w", err))
		}
	} else {
		if err := server.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			server.error(fmt.Errorf("failed to listen: %w", err))
		}
	}
}

func (server *Server) Stop() {
	err := server.httpServer.Shutdown(context.TODO())
	if err != nil {
		server.error(fmt.Errorf("shutdown failed: %w", err))
	}

	if server.errC != nil {
		close(server.errC)
		server.errC = nil
	}
}

func (server *Server) Write(webSocketId string, data []byte) error {
	ws, ok := server.connections[webSocketId]
	if !ok {
		return errors.New(fmt.Sprintf("couldn't write to websocket. No socket with id %v is open", webSocketId))
	}
	ws.outQueue <- data
	return nil
}

func (server *Server) wsHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL
	// Check if requested subprotocol is supported
	clientSubprotocols := websocket.Subprotocols(r)
	supported := false
	if len(server.upgrader.Subprotocols) == 0 {
		// All subProtocols are accepted
		supported = true
	}
	for _, supportedProto := range server.upgrader.Subprotocols {
		for _, requestedProto := range clientSubprotocols {
			if requestedProto == supportedProto {
				supported = true
				break
			}
		}
	}
	if !supported {
		server.error(fmt.Errorf("unsupported subprotocol: %v", clientSubprotocols))
		http.Error(w, "unsupported subprotocol", http.StatusBadRequest)
		return
	}
	// Handle client authentication
	if server.basicAuthHandler != nil {
		username, password, ok := r.BasicAuth()
		if !ok {
			server.error(errors.New("basic auth failed: credentials not found"))
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ok = server.basicAuthHandler(username, password)
		if !ok {
			server.error(errors.New("basic auth failed: credentials invalid"))
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}
	// Upgrade websocket
	conn, err := server.upgrader.Upgrade(w, r, nil)
	if err != nil {
		server.error(fmt.Errorf("upgrade failed: %w", err))
		return
	}
	ws := WebSocket{connection: conn, id: url.Path, outQueue: make(chan []byte), closeSignal: make(chan error, 1), pingMessage: make(chan []byte, 1)}
	server.connections[url.Path] = &ws
	// Read and write routines are started in separate goroutines and function will return immediately
	go server.writePump(&ws)
	go server.readPump(&ws)
	if server.newClientHandler != nil {
		var channel Channel = &ws
		server.newClientHandler(channel)
	}
}

func (server *Server) readPump(ws *WebSocket) {
	conn := ws.connection
	defer func() {
		_ = conn.Close()
		//TODO: close signal
		//ws.closeSignal <- true
	}()

	conn.SetPingHandler(func(appData string) error {
		ws.pingMessage <- []byte(appData)
		err := conn.SetReadDeadline(time.Now().Add(server.timeoutConfig.PingWait))
		return err
	})
	_ = conn.SetReadDeadline(time.Now().Add(server.timeoutConfig.PingWait))

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure) {
				server.error(fmt.Errorf("read failed for %s: %w", ws.GetID(), err))
			}
			if server.disconnectedHandler != nil {
				server.disconnectedHandler(ws)
			}
			break
		}

		if server.messageHandler != nil {
			var channel Channel = ws
			err = server.messageHandler(channel, message)
			if err != nil {
				server.error(fmt.Errorf("handling failed for %s: %w", ws.GetID(), err))
				continue
			}
		}
		_ = conn.SetReadDeadline(time.Now().Add(server.timeoutConfig.PingWait))
	}
}

func (server *Server) writePump(ws *WebSocket) {
	conn := ws.connection
	defer func() {
		_ = conn.Close()
	}()

	for {
		select {
		case data, ok := <-ws.outQueue:
			_ = conn.SetWriteDeadline(time.Now().Add(server.timeoutConfig.WriteWait))
			if !ok {
				// Closing connection
				err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					server.error(fmt.Errorf("close failed: %w", err))
				}
				return
			}
			err := conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				server.error(fmt.Errorf("write failed for %s: %w", ws.GetID(), err))
				return
			}
		case ping := <-ws.pingMessage:
			_ = conn.SetWriteDeadline(time.Now().Add(server.timeoutConfig.WriteWait))
			err := conn.WriteMessage(websocket.PongMessage, ping)
			if err != nil {
				server.error(fmt.Errorf("write failed for %s: %w", ws.GetID(), err))
				return
			}
		case closed, ok := <-ws.closeSignal:
			if !ok || closed != nil {
				//TODO: handle signal
				return
			}
		}
	}
}

// ---------------------- CLIENT ----------------------

// A Websocket client, needed to connect to a websocket server.
// The offered API are of asynchronous nature, and each incoming message is handled using callbacks.
//
// To create a new ws client, use:
//	client := NewClient()
//
// If you need a TLS ws client instead, use:
//	certPool, err := x509.SystemCertPool()
//	if err != nil {
//		log.Fatal(err)
//	}
//	// You may add more trusted certificates to the pool before creating the TLSClientConfig
//	client := NewTLSClient(&tls.Config{
//		RootCAs: certPool,
//	})
//
// To add additional dial options, use:
//	client.AddOption(func(*websocket.Dialer) {
//		// Your option ...
//	)}
//
// To add basic HTTP authentication, use:
//	client.SetBasicAuth("username","password")
//
// If you need to set a specific timeout configuration, refer to the SetTimeoutConfig method.
//
// Using Start and Stop you can respectively open/close a websocket to a websocket server.
//
// To receive incoming messages, you will need to set your own handler using SetMessageHandler.
// To write data on the open socket, simply call the Write function.
type WsClient interface {
	// Starts the client and attempts to connect to the server on a specified URL.
	// If the connection fails, an error is returned.
	//
	// For example:
	//	err := client.Start("ws://localhost:8887/ws/1234")
	//
	// The function returns immediately, after the connection has been established.
	// Incoming messages are passed automatically to the callback function, so no explicit read operation is required.
	//
	// To stop a running client, call the Stop function.
	Start(url string) error
	// Closes the output of the websocket Channel, effectively closing the connection to the server with a normal closure.
	Stop()
	// Errors returns a channel for error messages. If it doesn't exist it es created.
	// The channel is closed by the client when stopped.
	Errors() <-chan error
	// Sets a callback function for all incoming messages.
	SetMessageHandler(handler func(data []byte) error)
	// Set custom timeout configuration parameters. If not passed, a default ClientTimeoutConfig struct will be used.
	//
	// This function must be called before connecting to the server, otherwise it may lead to unexpected behavior.
	SetTimeoutConfig(config ClientTimeoutConfig)
	// Sets a callback function for receiving notifications about an unexpected disconnection from the server.
	// The callback is invoked even if the automatic reconnection mechanism is active.
	//
	// If the client was stopped using the Stop function, the callback will NOT be invoked.
	SetDisconnectedHandler(handler func(err error))
	// Sets a callback function for receiving notifications whenever the connection to the server is re-established.
	// Connections are re-established automatically thanks to the auto-reconnection mechanism.
	//
	// If set, the DisconnectedHandler will always be invoked before the Reconnected callback is invoked.
	SetReconnectedHandler(handler func())
	// IsConnected Returns information about the current connection status.
	// If the client is currently attempting to auto-reconnect to the server, the function returns false.
	IsConnected() bool
	// Sends a message to the server over the websocket.
	//
	// The data is queued and will be sent asynchronously in the background.
	Write(data []byte) error
	// Adds a websocket option to the client.
	AddOption(option interface{})
	// SetBasicAuth adds basic authentication credentials, to use when connecting to the server.
	// The credentials are automatically encoded in base64.
	SetBasicAuth(username string, password string)
	// SetHeaderValue sets a value on the HTTP header sent when opening a websocket connection to the server.
	//
	// The function overwrites previous header fields with the same key.
	SetHeaderValue(key string, value string)
}

// Client is the the default implementation of a Websocket client.
//
// Use the NewClient or NewTLSClient functions to create a new client.
type Client struct {
	webSocket      WebSocket
	messageHandler func(data []byte) error
	dialOptions    []func(*websocket.Dialer)
	header         http.Header
	timeoutConfig  ClientTimeoutConfig
	connected      bool
	onDisconnected func(err error)
	onReconnected  func()
	mutex          sync.Mutex
	errC           chan error
}

// Creates a new simple websocket client (the channel is not secured).
//
// Additional options may be added using the AddOption function.
// Basic authentication can be set using the SetBasicAuth function.
func NewClient() *Client {
	return &Client{dialOptions: []func(*websocket.Dialer){}, timeoutConfig: NewClientTimeoutConfig(), header: http.Header{}}
}

// Creates a new secure websocket client. If supported by the server, the websocket channel will use TLS.
//
// Additional options may be added using the AddOption function.
// Basic authentication can be set using the SetBasicAuth function.
//
// To set a client certificate, you may do:
//	certificate, _ := tls.LoadX509KeyPair(clientCertPath, clientKeyPath)
//	clientCertificates := []tls.Certificate{certificate}
//	client := ws.NewTLSClient(&tls.Config{
//		RootCAs:      certPool,
//		Certificates: clientCertificates,
//	})
//
// You can set any other TLS option within the same constructor as well.
// For example, if you wish to test connecting to a server having a
// self-signed certificate (do not use in production!), pass:
//	InsecureSkipVerify: true
func NewTLSClient(tlsConfig *tls.Config) *Client {
	client := &Client{dialOptions: []func(*websocket.Dialer){}, timeoutConfig: NewClientTimeoutConfig(), header: http.Header{}}
	client.dialOptions = append(client.dialOptions, func(dialer *websocket.Dialer) {
		dialer.TLSClientConfig = tlsConfig
	})
	return client
}

func (client *Client) SetMessageHandler(handler func(data []byte) error) {
	client.messageHandler = handler
}

func (client *Client) SetTimeoutConfig(config ClientTimeoutConfig) {
	client.timeoutConfig = config
}

func (client *Client) SetDisconnectedHandler(handler func(err error)) {
	client.onDisconnected = handler
}

func (client *Client) SetReconnectedHandler(handler func()) {
	client.onReconnected = handler
}

func (client *Client) AddOption(option interface{}) {
	dialOption, ok := option.(func(*websocket.Dialer))
	if ok {
		client.dialOptions = append(client.dialOptions, dialOption)
	}
}

func (client *Client) SetBasicAuth(username string, password string) {
	client.header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(username+":"+password)))
}

func (client *Client) SetHeaderValue(key string, value string) {
	client.header.Set(key, value)
}

func (client *Client) writePump() {
	ticker := time.NewTicker(client.timeoutConfig.PingPeriod)
	conn := client.webSocket.connection
	// Closure function shuts down the current connection
	closure := func(err error) {
		ticker.Stop()
		_ = conn.Close()
		client.setConnected(false)
		if client.onDisconnected != nil && err != nil {
			client.onDisconnected(err)
		}
	}

	for {
		select {
		case data, ok := <-client.webSocket.outQueue:
			_ = conn.SetWriteDeadline(time.Now().Add(client.timeoutConfig.WriteWait))
			if !ok {
				// Closing connection normally
				err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					client.error(fmt.Errorf("close failed: %w", err))
				}
				// Disconnected by user command. Not calling auto-reconnect.
				// Passing nil will also not call onDisconnected
				closure(nil)
				return
			}
			err := conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				client.error(fmt.Errorf("write failed: %w", err))
				closure(err)
				client.handleReconnection()
				return
			}
		case <-ticker.C:
			// Send periodic ping
			_ = conn.SetWriteDeadline(time.Now().Add(client.timeoutConfig.WriteWait))
			if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				client.error(fmt.Errorf("write failed: %w", err))
				closure(err)
				client.handleReconnection()
				return
			}
		case closed, ok := <-client.webSocket.closeSignal:
			// Read pump sent a closeSignal (i.e. a message couldn't be read in that moment)
			if !ok || closed != nil {
				closure(closed)
				client.handleReconnection()
				return
			}
		}
	}
}

func (client *Client) readPump() {
	conn := client.webSocket.connection
	_ = conn.SetReadDeadline(time.Now().Add(client.timeoutConfig.PongWait))
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(client.timeoutConfig.PongWait))
	})
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure) {
				client.error(fmt.Errorf("read failed: %w", err))
			}
			// Notify writePump of error. Disconnection will be handled there
			client.webSocket.closeSignal <- err
			return
		}

		if client.messageHandler != nil {
			err = client.messageHandler(message)
			if err != nil {
				// TODO: Handle?
				client.error(fmt.Errorf("handle failed: %w", err))
				continue
			}
		}
	}
}

func (client *Client) handleReconnection() {
	delay := defaultAutoReconnectDelay
	for {
		// Wait before reconnecting
		time.Sleep(delay)
		err := client.Start(client.webSocket.id)
		if err == nil {
			// Re-connection was successful
			if client.onReconnected != nil {
				client.onReconnected()
			}
			return
		}
		// Re-connection failed, increase delay exponentially
		delay *= 2
		if delay >= defaultMaxReconnectionDelay {
			delay = defaultMaxReconnectionDelay
		}
	}
}

func (client *Client) setConnected(connected bool) {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	client.connected = connected
}

func (client *Client) IsConnected() bool {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	return client.connected
}

func (client *Client) Write(data []byte) error {
	if !client.IsConnected() {
		return errors.New("client is currently not connected, cannot send data")
	}
	client.webSocket.outQueue <- data
	return nil
}

func (client *Client) Start(url string) error {
	dialer := websocket.Dialer{
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		HandshakeTimeout: client.timeoutConfig.HandshakeTimeout,
		Subprotocols:     []string{},
	}
	for _, option := range client.dialOptions {
		option(&dialer)
	}
	ws, resp, err := dialer.Dial(url, client.header)
	if err != nil {
		if resp != nil {
			httpError := HttpConnectionError{Message: err.Error(), HttpStatus: resp.Status, HttpCode: resp.StatusCode}
			// Parse http response details
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			if body != nil {
				httpError.Details = string(body)
			}
			err = httpError
		}
		client.error(fmt.Errorf("connect failed: %w", err))
		return err
	}

	client.webSocket = WebSocket{connection: ws, id: url, outQueue: make(chan []byte), closeSignal: make(chan error, 1)}
	client.setConnected(true)
	//Start reader and write routine
	go client.writePump()
	go client.readPump()
	return nil
}

func (client *Client) Stop() {
	close(client.webSocket.outQueue)

	if client.errC != nil {
		close(client.errC)
		client.errC = nil
	}
}

func (client *Client) error(err error) {
	if client.errC != nil {
		client.errC <- err
	}
}

func (client *Client) Errors() <-chan error {
	if client.errC == nil {
		client.errC = make(chan error, 1)
	}
	return client.errC
}
