package ws

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Time allowed to wait for a ping on the server, before closing a connection due to inactivity.
	pingWait = pongWait

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512

	// Time allowed for the initial handshake to complete.
	handshakeTimeout = 30 * time.Second

	// Default sub-protocol to send to peer upon connection.
	defaultSubProtocol = "ocpp1.6"
)

var upgrader = websocket.Upgrader{Subprotocols: []string{}}

// Channel represents a bi-directional communication channel, which provides at least a unique ID.
type Channel interface {
	GetID() string
}

// WebSocket is a wrapper for a single websocket channel.
// The connection itself is provided by the gorilla websocket package.
// Don't use a websocket directly, but refer to WsServer and WsClient.
type WebSocket struct {
	connection  *websocket.Conn
	id          string
	outQueue    chan []byte
	closeSignal chan bool
	pingMessage chan []byte
}

// Retrieves the unique Identifier of the websocket (typically, the URL suffix).
func (websocket *WebSocket) GetID() string {
	return websocket.id
}

// ---------------------- SERVER ----------------------

// A Websocket server, which passively listens for incoming connections on ws or wss protocol.
// The offered API are of asynchronous nature, and each incoming connection/message is handled using callbacks.
type WsServer interface {
	Start(port int, listenPath string)
	Stop()
	SetMessageHandler(handler func(ws Channel, data []byte) error)
	SetNewClientHandler(handler func(ws Channel))
	SetDisconnectedClientHandler(handler func(ws Channel))
	Write(webSocketId string, data []byte) error
	AddSupportedSubprotocol(subProto string)
}

// Default implementation of a Websocket server.
//
// To create a new ws server, use:
//	server := NewServer()
//
// If you need a TLS ws server instead, use:
//	server := NewTLSServer("cert.pem", "privateKey.pem")
type Server struct {
	connections         map[string]*WebSocket
	httpServer          *http.Server
	messageHandler      func(ws Channel, data []byte) error
	newClientHandler    func(ws Channel)
	disconnectedHandler func(ws Channel)
	tlsCertificatePath  string
	tlsCertificateKey   string
}

// Creates a new simple websocket server (the websockets are not secured).
func NewServer() *Server {
	return &Server{}
}

// Creates a new secure websocket server. All created websocket channels will use TLS.
func NewTLSServer(certificatePath string, certificateKey string) *Server {
	return &Server{tlsCertificatePath: certificatePath, tlsCertificateKey: certificateKey}
}

// Sets a callback function for all incoming messages.
// The callbacks accepts a Channel and the received data.
// It is up to the callback receiver, to check the identifier of the channel, to determine the source of the message.
func (server *Server) SetMessageHandler(handler func(ws Channel, data []byte) error) {
	server.messageHandler = handler
}

// Sets a callback function for all new incoming client connections.
// It is recommended to store a reference to the Channel in the received entity, so that the Channel may be recognized later on.
func (server *Server) SetNewClientHandler(handler func(ws Channel)) {
	server.newClientHandler = handler
}

// Sets a callback function for all client disconnection events.
// Once a client is disconnected, it is not possible to read/write on the respective Channel any longer.
func (server *Server) SetDisconnectedClientHandler(handler func(ws Channel)) {
	server.disconnectedHandler = handler
}

// Adds support for a specified subprotocol.
// This is recommended in order to communicate the capabilities to the client during the handshake.
// If left empty, any subprotocol will be accepted.
// Duplicates will be removed automatically.
func (server *Server) AddSupportedSubprotocol(subProto string) {
	for _, sub := range upgrader.Subprotocols {
		if sub == subProto {
			// Don't add duplicates
			return
		}
	}
	upgrader.Subprotocols = append(upgrader.Subprotocols, subProto)
}

// Starts and runs the websocket server on a specific port and URL.
// After start, incoming connections and messages are handled automatically, so no explicit read operation is required.
//
// The functions blocks forever, hence it is suggested to invoke it in a goroutine, if the caller thread needs to perform other work, e.g.:
//	go server.Start(8887, "/ws/{id}")
//	doStuffOnMainThread()
//	...
//
// To stop a running server, call the Stop function.
func (server *Server) Start(port int, listenPath string) {
	router := mux.NewRouter()
	router.HandleFunc(listenPath, func(w http.ResponseWriter, r *http.Request) {
		server.wsHandler(w, r)
	})
	server.connections = make(map[string]*WebSocket)
	addr := fmt.Sprintf(":%v", port)
	server.httpServer = &http.Server{Addr: addr, Handler: router}
	if server.tlsCertificatePath != "" && server.tlsCertificateKey != "" {
		if err := server.httpServer.ListenAndServeTLS(server.tlsCertificatePath, server.tlsCertificateKey); err != http.ErrServerClosed {
			log.Errorf("websocket server error: %v", err)
		}
	} else {
		if err := server.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Errorf("websocket server error: %v", err)
		}
	}
}

// Shuts down a running websocket server.
// All open channels will be forcefully closed, and the previously called Start function will return.
func (server *Server) Stop() {
	err := server.httpServer.Shutdown(context.TODO())
	if err != nil {
		log.Errorf("error while shutting down server: %v", err)
	}
}

// Sends a message on a specific Channel, identifier by the webSocketId parameter.
// If the passed ID is invalid, an error is returned.
//
// The data is queued and will be sent asynchronously in the background.
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
	if len(upgrader.Subprotocols) == 0 {
		// All subProtocols are accepted
		supported = true
	}
	for _, supportedProto := range upgrader.Subprotocols {
		for _, requestedProto := range clientSubprotocols {
			if requestedProto == supportedProto {
				supported = true
				break
			}
		}
	}
	if !supported {
		log.Warnf("client on %v requested unsupported subprotocols %v, closing socket", url.String(), clientSubprotocols)
		http.Error(w, "unsupported subprotocol", http.StatusBadRequest)
		return
	}
	// Upgrade websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("new client on URL %v", url.String())
	ws := WebSocket{connection: conn, id: url.Path, outQueue: make(chan []byte), closeSignal: make(chan bool, 1), pingMessage: make(chan []byte, 1)}
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
		ws.closeSignal <- true
	}()

	conn.SetPingHandler(func(appData string) error {
		log.WithField("client", ws.GetID()).Debug("ping received")
		ws.pingMessage <- []byte(appData)
		err := conn.SetReadDeadline(time.Now().Add(pingWait))
		return err
	})
	_ = conn.SetReadDeadline(time.Now().Add(pingWait))

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure) {
				log.WithFields(log.Fields{"client": ws.GetID()}).Errorf("error while reading from ws: %v", err)
			}
			if server.disconnectedHandler != nil {
				server.disconnectedHandler(ws)
			}
			break
		}
		log.WithFields(log.Fields{"client": ws.GetID()}).Debug("received message")
		if server.messageHandler != nil {
			var channel Channel = ws
			err = server.messageHandler(channel, message)
			if err != nil {
				log.WithFields(log.Fields{"client": ws.GetID()}).Errorf("error while handling message: %v", err)
				continue
			}
		}
		_ = conn.SetReadDeadline(time.Now().Add(pingWait))
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
			_ = conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Closing connection
				err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					log.WithFields(log.Fields{"client": ws.GetID()}).Errorf("error while closing: %v", err)
				}
				return
			}
			err := conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.WithFields(log.Fields{"client": ws.GetID()}).Errorf("error writing to websocket: %v", err)
				return
			}
		case ping := <-ws.pingMessage:
			_ = conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := conn.WriteMessage(websocket.PongMessage, ping)
			if err != nil {
				log.WithFields(log.Fields{"client": ws.GetID()}).Errorf("error writing to websocket: %v", err)
				return
			}
		case closed, ok := <-ws.closeSignal:
			if !ok || closed {
				return
			}
		}
	}
}

// ---------------------- CLIENT ----------------------

// A Websocket client, needed to connect to a websocket server.
// The offered API are of asynchronous nature, and each incoming message is handled using callbacks.
type WsClient interface {
	Start(url string) error
	Stop()
	SetMessageHandler(handler func(data []byte) error)
	Write(data []byte) error
	AddOption(option interface{})
}

// Default implementation of a Websocket client.
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
type Client struct {
	webSocket      WebSocket
	messageHandler func(data []byte) error
	dialOptions    []func(*websocket.Dialer)
}

// Creates a new simple websocket client (the channel is not secured).
// Additional options may be added using the AddOption function.
func NewClient() *Client {
	return &Client{dialOptions: []func(*websocket.Dialer){}}
}

// Creates a new secure websocket client. If supported by the server, the websocket channel will use TLS.
// Additional options may be passed.
func NewTLSClient(tlsConfig *tls.Config) *Client {
	cli := &Client{dialOptions: []func(*websocket.Dialer){}}
	cli.dialOptions = append(cli.dialOptions, func(dialer *websocket.Dialer) {
		dialer.TLSClientConfig = tlsConfig
	})
	return cli
}

// Sets a callback function for all incoming messages.
func (client *Client) SetMessageHandler(handler func(data []byte) error) {
	client.messageHandler = handler
}

// Adds a websocket option to the
func (client *Client) AddOption(option interface{}) {
	dialOption, ok := option.(func(*websocket.Dialer))
	if ok {
		client.dialOptions = append(client.dialOptions, dialOption)
	}
}

func (client *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	conn := client.webSocket.connection
	defer func() {
		ticker.Stop()
		_ = conn.Close()
	}()

	for {
		select {
		case data, ok := <-client.webSocket.outQueue:
			_ = conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Closing connection normally
				err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					log.Errorf("error while closing: %v", err)
				}
				return
			}
			err := conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Errorf("error writing to websocket: %v", err)
				return
			}
		case <-ticker.C:
			log.Debug("will send ping to server")
			_ = conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Errorf("couldn't send ping message: %v", err)
				return
			}
		case closed, ok := <-client.webSocket.closeSignal:
			if !ok || closed {
				return
			}
		}
	}
}

func (client *Client) readPump() {
	conn := client.webSocket.connection
	defer func() {
		_ = conn.Close()
		client.webSocket.closeSignal <- true
	}()
	_ = conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		log.Debug("pong received")
		return conn.SetReadDeadline(time.Now().Add(pongWait))
	})
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure) {
				log.Errorf("error while reading from websocket: %v", err)
			}
			return
		}
		log.Debugf("received message from server: %v", string(message))
		if client.messageHandler != nil {
			err = client.messageHandler(message)
			if err != nil {
				log.Errorf("error while handling message: %v", err)
				continue
			}
		}
	}
}

// Sends a message to the server over the websocket.
//
// The data is queued and will be sent asynchronously in the background.
func (client *Client) Write(data []byte) error {
	client.webSocket.outQueue <- data
	return nil
}

// Starts the client and attempts to connect to the server on a sepcified URL.
// If the connection fails, an error is returned.
//
// For example:
//	err := client.Start("ws://localhost:8887/ws/1234")
//
// The function returns immediately, after the connection has been established.
// Incoming messages are passed automatically to the callback function, so no explicit read operation is required.
//
// To stop a running client, call the Stop function.
func (client *Client) Start(url string) error {
	dialer := websocket.Dialer{
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		HandshakeTimeout: handshakeTimeout,
		Subprotocols:     []string{},
	}
	for _, option := range client.dialOptions {
		option(&dialer)
	}
	ws, _, err := dialer.Dial(url, nil)
	if err != nil {
		log.Errorf("couldn't connect to server: %v", err)
		return err
	}
	client.webSocket = WebSocket{connection: ws, id: url, outQueue: make(chan []byte), closeSignal: make(chan bool, 1)}
	//Start reader and write routine
	go client.writePump()
	go client.readPump()
	return nil
}

// Closes the output of the websocket Channel, effectively closing the connection to the server with a normal closure.
func (client *Client) Stop() {
	close(client.webSocket.outQueue)
}
