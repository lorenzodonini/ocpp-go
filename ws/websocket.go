// The package is a wrapper around gorilla websockets,
// aimed at simplifying the creation and usage of a websocket client/server.
//
// Check the Client and Server structure to get started.
package ws

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/lorenzodonini/ocpp-go/logging"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	defaultWriteWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	defaultPongWait = 60 * time.Second
	// Time allowed waiting for a ping on the server, before closing a connection due to inactivity.
	defaultPingWait = defaultPongWait
	// Send pings to peer with this period. Must be less than pongWait.
	defaultPingPeriod = (defaultPongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 512
	// Time allowed for the initial handshake to complete.
	defaultHandshakeTimeout = 30 * time.Second
	// Default sub-protocol to send to peer upon connection.
	defaultSubProtocol = "ocpp1.6"
	// The base delay to be used for automatic reconnection. Will double for every attempt up to maxReconnectionDelay.
	defaultReconnectBackoff = 5 * time.Second
	// Default maximum reconnection delay for websockets
	defaultReconnectMaxBackoff = 2 * time.Minute
)

// The internal verbose logger
var log logging.Logger

// SetLogger sets a custom Logger implementation, allowing the package to log events.
// By default, a VoidLogger is used, so no logs will be sent to any output.
//
// The function panics, if a nil logger is passed.
func SetLogger(logger logging.Logger) {
	if logger == nil {
		panic("cannot set a nil logger")
	}
	log = logger
}

// Channel represents a bi-directional communication channel, which provides at least a unique ID.
type Channel interface {
	// ID retrieves the unique Identifier of the channel (for a websocket this is typically the URL suffix).
	ID() string
	// RemoteAddr returns the address of the remote peer.
	RemoteAddr() net.Addr
	// TLSConnectionState returns the TLS connection state of the connection, if any.
	TLSConnectionState() *tls.ConnectionState
}

// WebSocket is a wrapper for a single websocket channel.
// The connection itself is provided by the gorilla websocket package.
//
// Don't use a websocket directly, but refer to WsServer and WsClient.
type WebSocket struct {
	connection         *websocket.Conn
	id                 string
	outQueue           chan []byte
	closeC             chan websocket.CloseError // used to gracefully close a websocket connection.
	forceCloseC        chan error                // used by the readPump to notify a forcefully closed connection to the writePump.
	pingMessage        chan []byte
	tlsConnectionState *tls.ConnectionState
}

func (websocket *WebSocket) ID() string {
	return websocket.id
}

func (websocket *WebSocket) RemoteAddr() net.Addr {
	return websocket.connection.RemoteAddr()
}

func (websocket *WebSocket) TLSConnectionState() *tls.ConnectionState {
	return websocket.tlsConnectionState
}

// HttpConnectionError wraps an http error, that may be raised when connecting to a websocket server.
type HttpConnectionError struct {
	Message    string
	HttpStatus string
	HttpCode   int
	Details    string
}

func (e HttpConnectionError) Error() string {
	return fmt.Sprintf("%v, http status: %v", e.Message, e.HttpStatus)
}

func init() {
	log = &logging.VoidLogger{}
}
