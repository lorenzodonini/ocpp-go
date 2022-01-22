package ws

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// ClientTimeoutConfig contains optional configuration parameters for a websocket client.
// Setting the parameter allows to define custom timeout intervals for websocket network operations.
//
// To set a custom configuration, refer to the client's SetTimeoutConfig method.
// If no configuration is passed, a default configuration is generated via the NewClientTimeoutConfig function.
type ClientTimeoutConfig struct {
	WriteWait           time.Duration
	HandshakeTimeout    time.Duration
	PongWait            time.Duration
	PingPeriod          time.Duration
	ReconnectBackoff    time.Duration
	ReconnectMaxBackoff time.Duration
}

// NewClientTimeoutConfig creates a default timeout configuration for a websocket endpoint.
//
// You may change fields arbitrarily and pass the struct to a SetTimeoutConfig method.
func NewClientTimeoutConfig() ClientTimeoutConfig {
	return ClientTimeoutConfig{
		WriteWait:           defaultWriteWait,
		HandshakeTimeout:    defaultHandshakeTimeout,
		PongWait:            defaultPongWait,
		PingPeriod:          defaultPingPeriod,
		ReconnectBackoff:    defaultReconnectBackoff,
		ReconnectMaxBackoff: defaultReconnectMaxBackoff,
	}
}

// WsClient defines a websocket client, needed to connect to a websocket server.
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
//		// Your option
//	})
//
// To add basic HTTP authentication, use:
//	client.SetBasicAuth("username","password")
//
// If you need to set a specific timeout configuration (e.g. ping), refer to the SetTimeoutConfig method.
//
// Using Start and Stop you can respectively open/close a websocket to a websocket server.
//
// To receive incoming messages, you will need to set your own handler using SetMessageHandler.
// To write data on the open socket, simply call the Write function.
type WsClient interface {
	// Start the client and attempts to connect to the server on a specified URL.
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
	// Stop closes the output of the websocket Channel, effectively closing the connection to
	// the server with a normal closure.
	Stop()
	// Errors returns a channel for error messages. If it doesn't exist it es created.
	// The channel is closed by the client when stopped.
	Errors() <-chan error
	// SetMessageHandler sets a callback function for all incoming messages.
	SetMessageHandler(handler func(data []byte) error)
	// SetTimeoutConfig set custom timeout configuration parameters. If not passed, a default
	// ClientTimeoutConfig struct will be used.
	//
	// This function must be called before connecting to the server, otherwise it may lead to unexpected behavior.
	SetTimeoutConfig(config ClientTimeoutConfig)
	// SetDisconnectedHandler sets a callback function for receiving notifications about an
	// unexpected disconnection from the server.
	// The callback is invoked even if the automatic reconnection mechanism is active.
	//
	// If the client was stopped using the Stop function, the callback will still be invoked.
	SetDisconnectedHandler(handler func(err error))
	// SetReconnectedHandler sets a callback function for receiving notifications whenever the connection to the
	// server is re-established.
	// Connections are re-established automatically thanks to the auto-reconnection mechanism.
	//
	// If set, the DisconnectedHandler will always be invoked before the Reconnected callback is invoked.
	SetReconnectedHandler(handler func())
	// IsConnected Returns information about the current connection status.
	// If the client is currently attempting to auto-reconnect to the server, the function returns false.
	IsConnected() bool
	// Write a message to the server over the websocket.
	//
	// The data is queued and will be sent asynchronously in the background.
	Write(data []byte) error
	// AddOption adds a websocket option to the client.
	// Currently, only websocket.Dialer options are supported.
	AddOption(option interface{})
	// SetBasicAuth adds basic authentication credentials, to use when connecting to the server.
	// The credentials are automatically encoded in base64.
	SetBasicAuth(username string, password string)
	// SetHeaderValue sets a value on the HTTP header sent when opening a websocket connection to the server.
	//
	// The function overwrites previous header fields with the same key.
	SetHeaderValue(key string, value string)
}

// Client is the default implementation of a Websocket client.
//
// Use the NewClient or NewTLSClient functions to create a new client.
type Client struct {
	webSocket      WebSocket
	url            url.URL
	messageHandler func(data []byte) error
	dialOptions    []func(*websocket.Dialer)
	header         http.Header
	timeoutConfig  ClientTimeoutConfig
	connected      bool
	onDisconnected func(err error)
	onReconnected  func()
	mutex          sync.Mutex
	errC           chan error
	reconnectC     chan struct{} // used for signaling, that a reconnection attempt should be interrupted
}

// NewClient creates a new simple websocket client (the channel is not secured).
//
// Additional options may be added using the AddOption function.
// Basic authentication can be set using the SetBasicAuth function.
func NewClient() *Client {
	return &Client{dialOptions: []func(*websocket.Dialer){}, timeoutConfig: NewClientTimeoutConfig(), header: http.Header{}}
}

// NewTLSClient creates a new secure websocket client. If supported by the server, the websocket channel will use TLS.
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

func (client *Client) getReadTimeout() time.Time {
	if client.timeoutConfig.PongWait == 0 {
		return time.Time{}
	}
	return time.Now().Add(client.timeoutConfig.PongWait)
}

func (client *Client) writePump() {
	ticker := time.NewTicker(client.timeoutConfig.PingPeriod)
	conn := client.webSocket.connection
	// Closure function correctly closes the current connection
	closure := func(err error) {
		ticker.Stop()
		client.cleanup()
		// Invoke callback
		if client.onDisconnected != nil {
			client.onDisconnected(err)
		}
	}

	for {
		select {
		case data, _ := <-client.webSocket.outQueue:
			// Send data
			log.Debugf("sending data")
			_ = conn.SetWriteDeadline(time.Now().Add(client.timeoutConfig.WriteWait))
			err := conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				client.error(fmt.Errorf("write failed: %w", err))
				closure(err)
				client.handleReconnection()
				return
			}
			log.Debugf("written %d bytes", len(data))
		case <-ticker.C:
			// Send periodic ping
			_ = conn.SetWriteDeadline(time.Now().Add(client.timeoutConfig.WriteWait))
			if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				client.error(fmt.Errorf("failed to send ping message: %w", err))
				closure(err)
				client.handleReconnection()
				return
			}
			log.Debugf("ping sent")
		case closeErr, _ := <-client.webSocket.closeC:
			log.Debugf("closing connection")
			// Closing connection gracefully
			if err := conn.WriteControl(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(closeErr.Code, closeErr.Text),
				time.Now().Add(client.timeoutConfig.WriteWait),
			); err != nil {
				client.error(fmt.Errorf("failed to write close message: %w", err))
			}
			// Disconnected by user command. Not calling auto-reconnect.
			// Passing nil will also not call onDisconnected.
			closure(nil)
			return
		case closed, ok := <-client.webSocket.forceCloseC:
			log.Debugf("handling forced close signal")
			// Read pump sent a forceClose signal (reading failed -> aborting the connection)
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
	_ = conn.SetReadDeadline(client.getReadTimeout())
	conn.SetPongHandler(func(string) error {
		log.Debugf("pong received")
		return conn.SetReadDeadline(client.getReadTimeout())
	})
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure) {
				client.error(fmt.Errorf("read failed: %w", err))
			}
			// Notify writePump of error. Forced close will be handled there
			client.webSocket.forceCloseC <- err
			return
		}

		log.Debugf("received %v bytes", len(message))
		if client.messageHandler != nil {
			err = client.messageHandler(message)
			if err != nil {
				client.error(fmt.Errorf("handle failed: %w", err))
				continue
			}
		}
	}
}

// Frees internal resources after a websocket connection was signaled to be closed.
// From this moment onwards, no new messages may be sent.
func (client *Client) cleanup() {
	client.setConnected(false)
	ws := client.webSocket
	_ = ws.connection.Close()
	client.mutex.Lock()
	defer client.mutex.Unlock()
	close(ws.outQueue)
	close(ws.closeC)
}

func (client *Client) handleReconnection() {
	log.Info("started automatic reconnection handler")
	delay := client.timeoutConfig.ReconnectBackoff
	for {
		// Wait before reconnecting
		select {
		case <-time.After(delay):
		case <-client.reconnectC:
			return
		}

		err := client.Start(client.url.String())
		if err == nil {
			// Re-connection was successful
			log.Info("reconnected successfully to server")
			if client.onReconnected != nil {
				client.onReconnected()
			}
			return
		}
		client.error(fmt.Errorf("reconnection failed: %w", err))
		// Re-connection failed, double the delay
		delay *= 2
		if delay >= client.timeoutConfig.ReconnectMaxBackoff {
			delay = client.timeoutConfig.ReconnectMaxBackoff
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
		return fmt.Errorf("client is currently not connected, cannot send data")
	}
	log.Debugf("queuing data for server")
	client.webSocket.outQueue <- data
	return nil
}

func (client *Client) Start(urlStr string) error {
	u, err := url.Parse(urlStr)
	if err != nil {
		return err
	}

	dialer := websocket.Dialer{
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		HandshakeTimeout: client.timeoutConfig.HandshakeTimeout,
		Subprotocols:     []string{},
	}
	for _, option := range client.dialOptions {
		option(&dialer)
	}
	// Connect
	log.Info("connecting to server")
	ws, resp, err := dialer.Dial(urlStr, client.header)
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
		return err
	}

	// The id of the charge point is the final path element
	id := path.Base(u.Path)
	client.url = *u
	client.webSocket = WebSocket{
		connection:         ws,
		id:                 id,
		outQueue:           make(chan []byte, 1),
		closeC:             make(chan websocket.CloseError, 1),
		forceCloseC:        make(chan error, 1),
		tlsConnectionState: resp.TLS,
	}
	log.Infof("connected to server as %s", id)
	client.reconnectC = make(chan struct{})
	client.setConnected(true)
	//Start reader and write routine
	go client.writePump()
	go client.readPump()
	return nil
}

func (client *Client) Stop() {
	log.Infof("closing connection to server")
	client.mutex.Lock()
	if client.connected {
		client.connected = false
		// Send signal for gracefully shutting down the connection
		select {
		case client.webSocket.closeC <- websocket.CloseError{Code: websocket.CloseNormalClosure, Text: ""}:
		default:
		}
	}
	client.mutex.Unlock()
	// Notify reconnection goroutine to stop (if any)
	close(client.reconnectC)
	if client.errC != nil {
		close(client.errC)
		client.errC = nil
	}
	// Wait for connection to actually close

}

func (client *Client) error(err error) {
	log.Error(err)
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
