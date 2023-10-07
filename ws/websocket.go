// The package is a wrapper around gorilla websockets,
// aimed at simplifying the creation and usage of a websocket client/server.
//
// Check the Client and Server structure to get started.
package ws

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lorenzodonini/ocpp-go/logging"
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
	// Time allowed for the initial handshake to complete.
	defaultHandshakeTimeout = 30 * time.Second
	// When the Charging Station is reconnecting, after a connection loss, it will use this variable for the amount of time
	// it will double the previous back-off time. When the maximum number of increments is reached, the Charging
	// Station keeps connecting with the same back-off time.
	defaultRetryBackOffRepeatTimes = 5
	// When the Charging Station is reconnecting, after a connection loss, it will use this variable as the maximum value
	// for the random part of the back-off time. It will add a new random value to every increasing back-off time,
	// including the first connection attempt (with this maximum), for the amount of times it will double the previous
	// back-off time. When the maximum number of increments is reached, the Charging Station will keep connecting
	// with the same back-off time.
	defaultRetryBackOffRandomRange = 15 // seconds
	// When the Charging Station is reconnecting, after a connection loss, it will use this variable as the minimum backoff
	// time, the first time it tries to reconnect.
	defaultRetryBackOffWaitMinimum = 10 * time.Second
)

// The internal verbose logger
var log logging.Logger

// Sets a custom Logger implementation, allowing the package to log events.
// By default, a VoidLogger is used, so no logs will be sent to any output.
//
// The function panics, if a nil logger is passed.
func SetLogger(logger logging.Logger) {
	if logger == nil {
		panic("cannot set a nil logger")
	}
	log = logger
}

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
	WriteWait               time.Duration
	HandshakeTimeout        time.Duration
	PongWait                time.Duration
	PingPeriod              time.Duration
	RetryBackOffRepeatTimes int
	RetryBackOffRandomRange int
	RetryBackOffWaitMinimum time.Duration
}

// NewClientTimeoutConfig creates a default timeout configuration for a websocket endpoint.
//
// You may change fields arbitrarily and pass the struct to a SetTimeoutConfig method.
func NewClientTimeoutConfig() ClientTimeoutConfig {
	return ClientTimeoutConfig{
		WriteWait:               defaultWriteWait,
		HandshakeTimeout:        defaultHandshakeTimeout,
		PongWait:                defaultPongWait,
		PingPeriod:              defaultPingPeriod,
		RetryBackOffRepeatTimes: defaultRetryBackOffRepeatTimes,
		RetryBackOffRandomRange: defaultRetryBackOffRandomRange,
		RetryBackOffWaitMinimum: defaultRetryBackOffWaitMinimum,
	}
}

// Channel represents a bi-directional communication channel, which provides at least a unique ID.
type Channel interface {
	ID() string
	RemoteAddr() net.Addr
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

// Retrieves the unique Identifier of the websocket (typically, the URL suffix).
func (websocket *WebSocket) ID() string {
	return websocket.id
}

// Returns the address of the remote peer.
func (websocket *WebSocket) RemoteAddr() net.Addr {
	return websocket.connection.RemoteAddr()
}

// Returns the TLS connection state of the connection, if any.
func (websocket *WebSocket) TLSConnectionState() *tls.ConnectionState {
	return websocket.tlsConnectionState
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
	// SetCheckClientHandler sets a handler for validate incoming websocket connections, allowing to perform
	// custom client connection checks.
	SetCheckClientHandler(handler func(id string, r *http.Request) bool)
	// Addr gives the address on which the server is listening, useful if, for
	// example, the port is system-defined (set to 0).
	Addr() *net.TCPAddr
}

// ---------------------- CLIENT ----------------------

// WsClient defines a websocket client, needed to connect to a websocket server.
// The offered API are of asynchronous nature, and each incoming message is handled using callbacks.
//
// To create a new ws client, use:
//
//	client := NewClient()
//
// If you need a TLS ws client instead, use:
//
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
//
//	client.AddOption(func(*websocket.Dialer) {
//		// Your option ...
//	)}
//
// To add basic HTTP authentication, use:
//
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
	// SetRequestedSubProtocol will negotiate the specified sub-protocol during the websocket handshake.
	// Internally this creates a dialer option and invokes the AddOption method on the client.
	//
	// Duplicates generated by invoking this method multiple times will be ignored.
	SetRequestedSubProtocol(subProto string)
	// SetBasicAuth adds basic authentication credentials, to use when connecting to the server.
	// The credentials are automatically encoded in base64.
	SetBasicAuth(username string, password string)
	// SetHeaderValue sets a value on the HTTP header sent when opening a websocket connection to the server.
	//
	// The function overwrites previous header fields with the same key.
	SetHeaderValue(key string, value string)
}

func init() {
	log = &logging.VoidLogger{}
}
