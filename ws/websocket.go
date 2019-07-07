package ws

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{}

type Channel interface {
	GetId() string
}

type WebSocket struct {
	connection  *websocket.Conn
	id          string
	outQueue    chan []byte
	closeSignal chan bool
}

func (websocket *WebSocket) GetId() string {
	return websocket.id
}

// ---------------------- SERVER ----------------------
type WsServer interface {
	Start(port int, listenPath string)
	Stop()
	SetMessageHandler(handler func(ws Channel, data []byte) error)
	SetNewClientHandler(handler func(ws Channel))
	SetDisconnectedHandler(handler func(ws Channel))
	Write(webSocketId string, data []byte) error
}

type Server struct {
	connections         map[string]*WebSocket
	httpServer          *http.Server
	messageHandler      func(ws Channel, data []byte) error
	newClientHandler    func(ws Channel)
	disconnectedHandler func(ws Channel)
}

func NewServer() *Server {
	return &Server{}
}

func (server *Server) SetMessageHandler(handler func(ws Channel, data []byte) error) {
	server.messageHandler = handler
}

func (server *Server) SetNewClientHandler(handler func(ws Channel)) {
	server.newClientHandler = handler
}

func (server *Server) SetDisconnectedHandler(handler func(ws Channel)) {
	server.disconnectedHandler = handler
}

func (server *Server) Start(port int, listenPath string) {
	router := mux.NewRouter()
	router.HandleFunc(listenPath, func(w http.ResponseWriter, r *http.Request) {
		server.wsHandler(w, r)
	})
	server.connections = make(map[string]*WebSocket)
	addr := fmt.Sprintf(":%v", port)
	server.httpServer = &http.Server{Addr: addr, Handler: router}
	if err := server.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Listen and server error: %v", err)
	}
}

func (server *Server) Stop() {
	err := server.httpServer.Shutdown(context.TODO())
	if err != nil {
		log.Printf("Error while shutting down server %v", err)
	}
}

func (server *Server) Write(webSocketId string, data []byte) error {
	ws, ok := server.connections[webSocketId]
	if !ok {
		return errors.New(fmt.Sprintf("Couldn't write to websocket. No socket with id %v is open", webSocketId))
	}
	ws.outQueue <- data
	return nil
}

func (server *Server) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	url := r.URL
	log.Printf("New client on URL %v", url.String())
	ws := WebSocket{connection: conn, id: url.Path, outQueue: make(chan []byte)}
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
		if server.disconnectedHandler != nil {
			server.disconnectedHandler(ws)
		}
	}()

	_ = conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		log.Printf("Ping received")
		_ = conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure) {
				log.Printf("error while reading from ws: %v", err)
			}
			break
		}
		log.Printf("Received message from client %v", ws.GetId())
		if server.messageHandler != nil {
			var channel Channel = ws
			err = server.messageHandler(channel, message)
			if err != nil {
				log.Printf("Error while handling message: %v", err)
				continue
			}
		}
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
					log.Printf("Error while closing client -> %v", err)
				}
				return
			}
			err := conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Printf("Error writing to websocket %v", err)
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
type WsClient interface {
	Start(url string) error
	Stop()
	SetMessageHandler(handler func(data []byte) error)
	SetDisconnectHandler(handler func(err error))
	Write(data []byte) error
}

type Client struct {
	webSocket           WebSocket
	messageHandler      func(data []byte) error
	disconnectedHandler func(err error)
}

func NewClient() *Client {
	return &Client{}
}

func (client *Client) SetMessageHandler(handler func(data []byte) error) {
	client.messageHandler = handler
}

func (client *Client) SetDisconnectHandler(handler func(err error)) {
	client.disconnectedHandler = handler
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
					log.Printf("Error while closing client -> %v", err)
				}
				return
			}
			err := conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Printf("Error writing to websocket %v", err)
				return
			}
		case <-ticker.C:
			log.Println("Ping triggered")
			_ = conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Couldn't send ping message -> %v", err)
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
		return conn.SetReadDeadline(time.Now().Add(pongWait))
	})
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure) {
				log.Printf("error while reading from ws: %v", err)
			}
			return
		}
		log.Printf("Received message from server")
		if client.messageHandler != nil {
			err = client.messageHandler(message)
			if err != nil {
				log.Printf("Error while handling message: %v", err)
				continue
			}
		}
	}
}

func (client *Client) Write(data []byte) error {
	client.webSocket.outQueue <- data
	return nil
}

func (client *Client) Start(url string) error {
	dialer := websocket.Dialer{
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		HandshakeTimeout: 30 * time.Second,
		Subprotocols:     []string{"ocpp1.6"}, //TODO: move out of websocket file
	}
	ws, _, err := dialer.Dial(url, nil)
	if err != nil {
		log.Printf("Error %v", err)
		return err
	}
	client.webSocket = WebSocket{connection: ws, id: url, outQueue: make(chan []byte)}
	//Start reader and write routine
	go client.writePump()
	client.readPump()
	return nil
}

func (client *Client) Stop() {
	close(client.webSocket.outQueue)
}
