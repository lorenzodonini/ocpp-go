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

type WebSocket struct {
	connection *websocket.Conn
	Id string
	outQueue chan []byte
}

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

// ---------------------- SERVER ----------------------
type WsServer interface {
	Start(port int, listenPath string)
	Stop()
	SetMessageHandler(handler func(ws *WebSocket, data []byte) error)
	SetNewClientHandler(handler func(ws *WebSocket))
	Write(webSocketId string, data []byte) error
}

type Server struct {
	connections map[string]*WebSocket
	httpServer *http.Server
	messageHandler func(ws *WebSocket, data []byte) error
	newClientHandler func(ws *WebSocket)
}

func (server *Server)SetMessageHandler(handler func(ws *WebSocket, data []byte) error) {
	server.messageHandler = handler
}

func (server *Server)SetNewClientHandler(handler func(ws *WebSocket)) {
	server.newClientHandler = handler
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

func (server *Server)Write(webSocketId string, data []byte) error {
	ws, ok := server.connections[webSocketId]
	if !ok {
		return errors.New(fmt.Sprintf("Couldn't write to websocket. No socket with id %v is open", webSocketId))
	}
	ws.outQueue <- data
	return nil
}

func (server *Server)wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	url := r.URL
	log.Printf("New client on URL %v", url.String())
	ws := WebSocket{connection:conn, Id:url.Path, outQueue: make(chan []byte)}
	server.connections[url.Path] = &ws
	// Read and write routines are started in separate goroutines and function will return immediately
	go server.writePump(&ws)
	go server.readPump(&ws)
	if server.newClientHandler != nil {
		server.newClientHandler(&ws)
	}
}

func (server *Server)readPump(ws *WebSocket) {
	conn := ws.connection
	defer func() {
		_ = conn.Close()
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
		log.Printf("Received message from client %v", ws.Id)
		if server.messageHandler != nil {
			err = server.messageHandler(ws, message)
			if err != nil {
				log.Printf("Error while handling message: %v", err)
				//TODO: handle error
			}
		}
	}
}

func (server *Server)writePump(ws *WebSocket) {
	conn := ws.connection

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
				//TODO: handle error
				log.Printf("Error writing to websocket %v", err)
			}
		}
	}
}

// ---------------------- CLIENT ----------------------
type WsClient interface {
	Start(url string)
	Stop()
	SetMessageHandler(handler func(data []byte) error)
	Write(data []byte)
}

type Client struct {
	webSocket WebSocket
	messageHandler func(data []byte) error
}

func (client *Client)SetMessageHandler(handler func(data []byte) error) {
	client.messageHandler = handler
}

func (client *Client)writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	conn := client.webSocket.connection

	for {
		select {
		case data, ok := <-client.webSocket.outQueue:
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
				//TODO: handle error
				log.Printf("Error writing to websocket %v", err)
			}
		case <- ticker.C:
			log.Println("Ping triggered")
			_ = conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Couldn't send ping message -> %v", err)
				return
			}
		}
	}
}

func (client *Client)readPump() {
	conn := client.webSocket.connection
	defer func() {
		_ = conn.Close()
	}()
	_ = conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
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
		log.Printf("Received message from server")
		if client.messageHandler != nil {
			err = client.messageHandler(message)
			if err != nil {
				log.Printf("Error while handling message: %v", err)
				//TODO: handle error
			}
		}
	}
}

func (client *Client)Write(data []byte) {
	client.webSocket.outQueue <- data
}

func (client* Client) Start(url string) {
	dialer := websocket.Dialer{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
		HandshakeTimeout: 30 * time.Second,
		Subprotocols: []string{"ocpp1.6"},
	}
	ws, _, err := dialer.Dial(url, nil)
	if err != nil {
		log.Printf("Error %v", err)
	}
	client.webSocket = WebSocket{connection: ws, Id: url, outQueue: make(chan []byte)}
	//Start reader and write routine
	go client.writePump()
	client.readPump()
}

func (client* Client) Stop() {
	close(client.webSocket.outQueue)
}
