package ws

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/suite"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

const (
	serverPort = 8887
	serverPath = "/ws/{id}"
	testPath   = "/ws/testws"
	// Default sub-protocol to send to peer upon connection.
	defaultSubProtocol = "ocpp1.6"
)

func newWebsocketServer(t *testing.T, onMessage func(data []byte) ([]byte, error)) *server {
	wsServer := NewServer()
	innerS, ok := wsServer.(*server)
	require.True(t, ok)
	innerS.SetMessageHandler(func(ws Channel, data []byte) error {
		assert.NotNil(t, ws)
		assert.NotNil(t, data)
		if onMessage != nil {
			response, err := onMessage(data)
			assert.Nil(t, err)
			if response != nil {
				err = innerS.Write(ws.ID(), data)
				assert.Nil(t, err)
			}
		}
		return nil
	})
	return innerS
}

func newWebsocketClient(t *testing.T, onMessage func(data []byte) ([]byte, error)) *client {
	wsClient := NewClient()
	innerC, ok := wsClient.(*client)
	require.True(t, ok)
	innerC.SetRequestedSubProtocol(defaultSubProtocol)
	innerC.SetMessageHandler(func(data []byte) error {
		assert.NotNil(t, data)
		if onMessage != nil {
			response, err := onMessage(data)
			assert.Nil(t, err)
			if response != nil {
				err = innerC.Write(data)
				assert.Nil(t, err)
			}
		}
		return nil
	})
	return innerC
}

type WebSocketSuite struct {
	suite.Suite
	client *client
	server *server
}

func (s *WebSocketSuite) SetupTest() {
	s.server = newWebsocketServer(s.T(), nil)
	s.client = newWebsocketClient(s.T(), nil)
}

func (s *WebSocketSuite) TearDownTest() {
	if s.client != nil {
		s.client.Stop()
	}
	if s.server != nil {
		s.server.Stop()
	}
}

func (s *WebSocketSuite) TestPingTicker() {
	defaultPeriod := 1 * time.Millisecond
	testTable := []struct {
		name              string
		active            bool
		pingPeriod        time.Duration
		expectedTickerNil bool
	}{
		{
			"real ticker",
			true,
			defaultPeriod,
			false,
		},
		{
			"dummy ticker",
			false,
			defaultPeriod,
			true,
		},
		{
			"dummy ticker due to invalid period",
			true,
			0,
			true,
		},
	}
	for _, tc := range testTable {
		var pc *PingConfig
		if tc.active {
			pc = &PingConfig{
				PingPeriod: tc.pingPeriod,
			}
		}
		t := newOptTicker(pc)
		if tc.expectedTickerNil {
			s.Nil(t.ticker, tc.name)
			s.NotNil(t.c, tc.name)
		} else {
			s.NotNil(t.ticker, tc.name)
			s.Nil(t.c, tc.name)
		}
		// Test retrieving channel
		c := t.T()
		s.NotNil(c)
		// Test waiting for tick
		select {
		case <-c:
			if tc.expectedTickerNil {
				s.Fail("unexpected tick from nil ticker", tc.name)
			}
		case <-time.After(2 * defaultPeriod):
			if !tc.expectedTickerNil {
				s.Fail("unexpected timeout from real ticker", tc.name)
			}
		}
		// Test waiting for tick after stop
		t.Stop()
		select {
		case <-c:
			s.Fail("unexpected tick from stopped ticker", tc.name)
		case <-time.After(2 * defaultPeriod):
			break
		}
	}
}

func (s *WebSocketSuite) TestWebsocketConnectionState() {
	s.False(s.client.IsConnected())
	closeC := make(chan struct{}, 1)
	s.client.SetMessageHandler(func(data []byte) error {
		s.Fail("unexpected message")
		return nil
	})
	s.client.SetDisconnectedHandler(func(err error) {
		closeC <- struct{}{}
	})
	// Simulate connection
	go s.server.Start(serverPort, serverPath)
	time.Sleep(50 * time.Millisecond)
	host := fmt.Sprintf("localhost:%v", serverPort)
	u := url.URL{Scheme: "ws", Host: host, Path: testPath}
	err := s.client.Start(u.String())
	s.NoError(err)
	// Check connection state on internal web socket
	ws := s.client.webSocket
	s.NotNil(ws)
	s.True(ws.IsConnected())
	// Close connection
	err = ws.Close(websocket.CloseError{Code: websocket.CloseNormalClosure, Text: ""})
	s.NoError(err)
	select {
	case <-closeC:
	case <-time.After(1 * time.Second):
		s.Fail("timeout waiting for connection to close")
	}
	s.False(ws.IsConnected())
}

func (s *WebSocketSuite) TestWebsocketGetReadTimeout() {
	// Create default timeout settings and handlers
	serverTimeoutConfig := NewServerTimeoutConfig()
	s.server.SetTimeoutConfig(serverTimeoutConfig)
	ctrlC := make(chan struct{}, 1)
	s.server.SetNewClientHandler(func(ws Channel) {
		ctrlC <- struct{}{}
	})
	// Simulate connection to initialize a websocket
	go s.server.Start(serverPort, serverPath)
	time.Sleep(50 * time.Millisecond)
	host := fmt.Sprintf("localhost:%v", serverPort)
	u := url.URL{Scheme: "ws", Host: host, Path: testPath}
	err := s.client.Start(u.String())
	s.NoError(err)
	// Wait for connection to be established
	select {
	case <-ctrlC:
	case <-time.After(1 * time.Second):
		s.Fail("timeout waiting for connection to establish")
	}
	// Test server timeout for default settings
	serverW, ok := s.server.connections["testws"]
	s.True(ok)
	now := time.Now()
	timeout := serverW.getReadTimeout()
	s.GreaterOrEqual(timeout.Unix(), now.Add(s.server.timeoutConfig.PingWait).Unix())
	// Test server timeout for zero setting
	cfg := serverW.cfg
	cfg.ReadWait = 0
	serverW.updateConfig(cfg)
	timeout = serverW.getReadTimeout()
	s.Equal(time.Time{}, timeout)
	// Test client timeout for default settings
	clientW := s.client.webSocket
	s.NotNil(clientW)
	now = time.Now()
	timeout = clientW.getReadTimeout()
	s.GreaterOrEqual(timeout.Unix(), now.Add(s.client.timeoutConfig.PongWait).Unix())
	// Test client timeout for zero setting
	cfg = clientW.cfg
	cfg.PingConfig.PongWait = 0
	cfg.ReadWait = 0
	clientW.updateConfig(cfg)
	timeout = clientW.getReadTimeout()
	s.Equal(time.Time{}, timeout)
}

func (s *WebSocketSuite) TestWebsocketEcho() {
	msg := []byte("Hello webSocket!")
	triggerC := make(chan struct{}, 1)
	done := make(chan struct{}, 1)
	s.server = newWebsocketServer(s.T(), func(data []byte) ([]byte, error) {
		s.True(bytes.Equal(msg, data))
		// Echo reply received, notifying flow routine
		triggerC <- struct{}{}
		return data, nil
	})
	s.server.SetNewClientHandler(func(ws Channel) {
		tlsState := ws.TLSConnectionState()
		s.Nil(tlsState)
	})
	s.client = newWebsocketClient(s.T(), func(data []byte) ([]byte, error) {
		s.True(bytes.Equal(msg, data))
		// Echo response received, notifying flow routine
		done <- struct{}{}
		return nil, nil
	})
	// Start server
	go s.server.Start(serverPort, serverPath)
	// Start flow routine
	go func() {
		// Wait for messages to be exchanged in a dedicate routine.
		// Will reply to client.
		sig := <-triggerC
		s.NotNil(sig)
		err := s.server.Write(path.Base(testPath), msg)
		s.NoError(err)
		sig = <-triggerC
		s.NotNil(sig)
	}()
	time.Sleep(100 * time.Millisecond)
	// Test connection
	host := fmt.Sprintf("localhost:%v", serverPort)
	u := url.URL{Scheme: "ws", Host: host, Path: testPath}
	err := s.client.Start(u.String())
	s.NoError(err)
	s.True(s.client.IsConnected())
	// Test message
	err = s.client.Write(msg)
	s.NoError(err)
	// Wait for echo result
	select {
	case result := <-done:
		s.NotNil(result)
	case <-time.After(1 * time.Second):
		s.Fail("timeout waiting for echo result")
	}
}

func (s *WebSocketSuite) TestWebsocketChargePointIdResolver() {
	connected := make(chan string)
	s.server = newWebsocketServer(s.T(), func(data []byte) ([]byte, error) {
		s.Fail("no message should be received from client!")
		return nil, nil
	})
	s.server.SetChargePointIdResolver(func(*http.Request) (string, error) {
		return "my-custom-id", nil
	})
	s.server.SetNewClientHandler(func(ws Channel) {
		connected <- ws.ID()
	})
	go s.server.Start(serverPort, serverPath)
	time.Sleep(500 * time.Millisecond)

	// Test message
	s.client = newWebsocketClient(s.T(), func(data []byte) ([]byte, error) {
		s.Fail("no message should be received from server!")
		return nil, nil
	})

	host := fmt.Sprintf("localhost:%v", serverPort)
	u := url.URL{Scheme: "ws", Host: host, Path: testPath}
	// Attempt to connect and expect the custom resolved charge point id
	err := s.client.Start(u.String())
	s.NoError(err)
	result := <-connected
	s.Equal("my-custom-id", result)
}

func (s *WebSocketSuite) TestWebsocketChargePointIdResolverFailure() {
	s.server = newWebsocketServer(s.T(), func(data []byte) ([]byte, error) {
		s.Fail("no message should be received from client!")
		return nil, nil
	})
	s.server.SetChargePointIdResolver(func(*http.Request) (string, error) {
		return "", fmt.Errorf("test error")
	})
	go s.server.Start(serverPort, serverPath)
	time.Sleep(500 * time.Millisecond)

	// Test message
	s.client = newWebsocketClient(s.T(), func(data []byte) ([]byte, error) {
		s.Fail("no message should be received from server!")
		return nil, nil
	})

	host := fmt.Sprintf("localhost:%v", serverPort)
	u := url.URL{Scheme: "ws", Host: host, Path: testPath}
	// Attempt to connect and expect the custom resolved charge point id
	err := s.client.Start(u.String())
	s.Error(err)
	httpErr, ok := err.(HttpConnectionError)
	s.True(ok)
	s.Equal(http.StatusNotFound, httpErr.HttpCode)
	s.Equal("websocket: bad handshake", httpErr.Message)
}

func (s *WebSocketSuite) TestWebsocketBootRetries() {
	verifyConnection := func(client *client, connected bool) {
		maxAttempts := 20
		for i := 0; i <= maxAttempts; i++ {
			if client.IsConnected() != connected {
				time.Sleep(200 * time.Millisecond)
				continue
			}
		}
		s.Equal(connected, client.IsConnected())
	}
	s.server = newWebsocketServer(s.T(), func(data []byte) ([]byte, error) {
		return data, nil
	})
	s.client = newWebsocketClient(s.T(), func(data []byte) ([]byte, error) {
		return nil, nil
	})
	// Reduce timeout to make test faster
	s.client.timeoutConfig.RetryBackOffWaitMinimum = 1 * time.Second
	s.client.timeoutConfig.RetryBackOffRandomRange = 2

	go func() {
		// Start websocket client
		host := fmt.Sprintf("localhost:%v", serverPort)
		u := url.URL{Scheme: "ws", Host: host, Path: testPath}
		s.client.StartWithRetries(u.String())
	}()
	// Initial connection attempt fails, as server isn't listening yet
	s.False(s.client.IsConnected())

	time.Sleep(500 * time.Millisecond)

	go s.server.Start(serverPort, serverPath)
	verifyConnection(s.client, true)

	s.server.Stop()
	verifyConnection(s.client, false)
}

func (s *WebSocketSuite) TestTLSWebsocketEcho() {
	msg := []byte("Hello Secure webSocket!")
	triggerC := make(chan struct{}, 1)
	done := make(chan struct{}, 1)
	// Use NewServer(WithServerTLSConfig(...)) when in different package
	s.server = newWebsocketServer(s.T(), func(data []byte) ([]byte, error) {
		s.True(bytes.Equal(msg, data))
		// Message received, notifying flow routine
		triggerC <- struct{}{}
		return data, nil
	})
	s.server.SetNewClientHandler(func(ws Channel) {
		tlsState := ws.TLSConnectionState()
		s.NotNil(tlsState)
	})
	s.server.SetDisconnectedClientHandler(func(ws Channel) {
		// Connection closed, completing test
		done <- struct{}{}
	})
	// Create self-signed TLS certificate
	// TODO: use FiloSottile's lib for this
	certFilename := "/tmp/cert.pem"
	keyFilename := "/tmp/key.pem"
	err := createTLSCertificate(certFilename, keyFilename, "localhost", nil, nil)
	s.NoError(err)
	defer os.Remove(certFilename)
	defer os.Remove(keyFilename)

	// Set self-signed TLS certificate
	s.server.tlsCertificatePath = certFilename
	s.server.tlsCertificateKey = keyFilename
	// Create TLS client
	s.client = newWebsocketClient(s.T(), func(data []byte) ([]byte, error) {
		s.True(bytes.Equal(msg, data))
		// Echo response received, notifying flow routine
		done <- struct{}{}
		return nil, nil
	})
	s.client.AddOption(func(dialer *websocket.Dialer) {
		certPool := x509.NewCertPool()
		data, err := os.ReadFile(certFilename)
		s.NoError(err)
		ok := certPool.AppendCertsFromPEM(data)
		s.True(ok)
		dialer.TLSClientConfig = &tls.Config{
			RootCAs: certPool,
		}
	})

	// Start server
	go s.server.Start(serverPort, serverPath)
	// Start flow routine
	go func() {
		// Wait for messages to be exchanged, then close connection
		sig := <-triggerC
		s.NotNil(sig)
		err = s.server.Write(path.Base(testPath), msg)
		s.NoError(err)
		sig = <-triggerC
		s.NotNil(sig)
	}()
	time.Sleep(100 * time.Millisecond)

	// Test connection
	host := fmt.Sprintf("localhost:%v", serverPort)
	u := url.URL{Scheme: "wss", Host: host, Path: testPath}
	err = s.client.Start(u.String())
	s.NoError(err)
	s.True(s.client.IsConnected())
	// Test message
	err = s.client.Write(msg)
	s.NoError(err)
	// Wait for echo result
	select {
	case result := <-done:
		s.NotNil(result)
	case <-time.After(1 * time.Second):
		s.Fail("timeout waiting for echo result")
	}
}

func (s *WebSocketSuite) TestServerStartErrors() {
	triggerC := make(chan struct{}, 1)
	s.server = newWebsocketServer(s.T(), nil)
	s.server.SetNewClientHandler(func(ws Channel) {
		triggerC <- struct{}{}
	})
	// Make sure http server is initialized on start
	s.server.httpServer = nil
	// Listen for errors
	go func() {
		err, ok := <-s.server.Errors()
		s.True(ok)
		s.Error(err)
		triggerC <- struct{}{}
	}()
	time.Sleep(100 * time.Millisecond)
	go s.server.Start(serverPort, serverPath)
	time.Sleep(100 * time.Millisecond)
	// Starting server again throws error
	s.server.Start(serverPort, serverPath)
	r := <-triggerC
	s.NotNil(r)
}

func (s *WebSocketSuite) TestClientDuplicateConnectionWithKeepCurrentConnectionBehavior() {
	s.server = newWebsocketServer(s.T(), nil)
	s.server.SetNewClientHandler(func(ws Channel) {
	})
	// Start server
	go s.server.Start(serverPort, serverPath)
	time.Sleep(100 * time.Millisecond)
	// Connect client 1
	s.client = newWebsocketClient(s.T(), func(data []byte) ([]byte, error) {
		return nil, nil
	})
	host := fmt.Sprintf("localhost:%v", serverPort)
	u := url.URL{Scheme: "ws", Host: host, Path: testPath}
	err := s.client.Start(u.String())
	s.NoError(err)
	// Try to connect client 2
	disconnectC := make(chan struct{})
	wsClient2 := newWebsocketClient(s.T(), func(data []byte) ([]byte, error) {
		return nil, nil
	})
	wsClient2.SetDisconnectedHandler(func(err error) {
		s.IsType(&websocket.CloseError{}, err)
		var wsErr *websocket.CloseError
		ok := errors.As(err, &wsErr)
		s.True(ok)
		s.Equal(websocket.ClosePolicyViolation, wsErr.Code)
		s.Equal("a connection with this ID already exists", wsErr.Text)
		wsClient2.SetDisconnectedHandler(nil)
		disconnectC <- struct{}{}
	})
	err = wsClient2.Start(u.String())
	s.NoError(err)
	// Expect new connection to be closed immediately
	select {
	case _, ok := <-disconnectC:
		s.True(ok, "expected new client to have been disconnected")
	case <-time.After(1 * time.Second):
		s.Fail("timeout waiting for new client to disconnect")
	}
}

func (s *WebSocketSuite) TestClientDuplicateConnectionWithKeepNewConnectionBehavior() {
	s.server = newWebsocketServer(s.T(), nil)
	s.server.SetNewClientHandler(func(ws Channel) {
	})
	s.server.SetDuplicateConnectionBehavior(DuplicateConnectionBehaviorKeepNew)
	// Start server
	go s.server.Start(serverPort, serverPath)
	time.Sleep(100 * time.Millisecond)
	// Connect client 1
	disconnectC := make(chan struct{})
	s.client = newWebsocketClient(s.T(), func(data []byte) ([]byte, error) {
		return nil, nil
	})
	s.client.SetDisconnectedHandler(func(err error) {
		s.IsType(&websocket.CloseError{}, err)
		var wsErr *websocket.CloseError
		ok := errors.As(err, &wsErr)
		s.True(ok)
		s.Equal(websocket.ClosePolicyViolation, wsErr.Code)
		s.Equal("a connection with this ID has reconnected", wsErr.Text)
		s.client.SetDisconnectedHandler(nil)
		disconnectC <- struct{}{}
	})
	host := fmt.Sprintf("localhost:%v", serverPort)
	u := url.URL{Scheme: "ws", Host: host, Path: testPath}
	err := s.client.Start(u.String())
	s.NoError(err)
	// Connect client 2
	wsClient2 := newWebsocketClient(s.T(), func(data []byte) ([]byte, error) {
		return nil, nil
	})
	err = wsClient2.Start(u.String())
	s.NoError(err)
	// Expect current connection to be closed immediately
	select {
	case _, ok := <-disconnectC:
		s.True(ok, "expected current client to have been disconnected")
	case <-time.After(1 * time.Second):
		s.Fail("timeout waiting for current client to disconnect")
	}
}

func (s *WebSocketSuite) TestServerStopConnection() {
	triggerC := make(chan struct{}, 1)
	disconnectedClientC := make(chan struct{}, 1)
	disconnectedServerC := make(chan struct{}, 1)
	closeError := websocket.CloseError{
		Code: websocket.CloseGoingAway,
		Text: "CloseClientConnection",
	}
	wsID := "testws"
	s.server = newWebsocketServer(s.T(), nil)
	s.server.SetNewClientHandler(func(ws Channel) {
		triggerC <- struct{}{}
	})
	s.server.SetDisconnectedClientHandler(func(ws Channel) {
		disconnectedServerC <- struct{}{}
	})
	s.client = newWebsocketClient(s.T(), func(data []byte) ([]byte, error) {
		return nil, nil
	})
	s.client.SetDisconnectedHandler(func(err error) {
		s.IsType(&closeError, err)
		var closeErr *websocket.CloseError
		ok := errors.As(err, &closeErr)
		s.True(ok)
		s.Equal(closeError.Code, closeErr.Code)
		s.Equal(closeError.Text, closeErr.Text)
		disconnectedClientC <- struct{}{}
	})
	// Start server
	go s.server.Start(serverPort, serverPath)
	time.Sleep(100 * time.Millisecond)
	var c Channel
	var ok bool
	c, ok = s.server.GetChannel(wsID)
	s.False(ok)
	s.Nil(c)
	// Connect client
	host := fmt.Sprintf("localhost:%v", serverPort)
	u := url.URL{Scheme: "ws", Host: host, Path: testPath}
	err := s.client.Start(u.String())
	s.NoError(err)
	// Wait for client to connect
	_, ok = <-triggerC
	s.True(ok)
	// Verify channel
	c, ok = s.server.GetChannel(wsID)
	s.True(ok)
	s.NotNil(c)
	s.Equal(wsID, c.ID())
	s.True(c.IsConnected())
	// Close connection and wait for client to be closed
	err = s.server.StopConnection(path.Base(testPath), closeError)
	s.NoError(err)
	_, ok = <-disconnectedClientC
	s.True(ok)
	_, ok = <-disconnectedServerC
	s.True(ok)
	s.False(s.client.IsConnected())
	time.Sleep(100 * time.Millisecond)
	s.Empty(s.server.connections)
	// client will attempt to reconnect under the hood, but test finishes before this can happen
}

func (s *WebSocketSuite) TestWebsocketServerStopAllConnections() {
	triggerC := make(chan struct{}, 1)
	numClients := 5
	disconnectedServerC := make(chan struct{}, 1)
	s.server = newWebsocketServer(s.T(), nil)
	s.server.SetNewClientHandler(func(ws Channel) {
		triggerC <- struct{}{}
	})
	s.server.SetDisconnectedClientHandler(func(ws Channel) {
		disconnectedServerC <- struct{}{}
	})
	// Start server
	go s.server.Start(serverPort, serverPath)
	time.Sleep(100 * time.Millisecond)
	// Connect clients
	clients := []Client{}
	wg := sync.WaitGroup{}
	host := fmt.Sprintf("localhost:%v", serverPort)
	for i := 0; i < numClients; i++ {
		wsClient := newWebsocketClient(s.T(), func(data []byte) ([]byte, error) {
			return nil, nil
		})
		wsClient.SetDisconnectedHandler(func(err error) {
			s.IsType(&websocket.CloseError{}, err)
			var closeErr *websocket.CloseError
			ok := errors.As(err, &closeErr)
			s.True(ok)
			s.Equal(websocket.CloseNormalClosure, closeErr.Code)
			s.Equal("", closeErr.Text)
			wg.Done()
		})
		u := url.URL{Scheme: "ws", Host: host, Path: fmt.Sprintf("%v-%v", testPath, i)}
		err := wsClient.Start(u.String())
		s.NoError(err)
		clients = append(clients, wsClient)
		// Wait for client to connect
		_, ok := <-triggerC
		s.True(ok)
		wg.Add(1)
	}
	// Stop server and wait for clients to disconnect
	s.server.Stop()
	waitC := make(chan struct{}, 1)
	go func() {
		wg.Wait()
		waitC <- struct{}{}
	}()
	select {
	case <-waitC:
		break
	case <-time.After(1 * time.Second):
		s.Fail("timeout waiting for clients to disconnect")
	}
	// Double-check disconnection status
	for _, c := range clients {
		s.False(c.IsConnected())
		// client will attempt to reconnect under the hood, but test finishes before this can happen
		c.Stop()
	}
	time.Sleep(100 * time.Millisecond)
	s.Empty(s.server.connections)
}

func (s *WebSocketSuite) TestWebsocketClientConnectionBreak() {
	newClient := make(chan struct{})
	disconnected := make(chan struct{})
	s.server = newWebsocketServer(s.T(), nil)
	s.server.SetNewClientHandler(func(ws Channel) {
		newClient <- struct{}{}
	})
	s.server.SetDisconnectedClientHandler(func(ws Channel) {
		disconnected <- struct{}{}
	})
	go s.server.Start(serverPort, serverPath)
	time.Sleep(100 * time.Millisecond)

	// Test
	s.client = newWebsocketClient(s.T(), nil)
	host := fmt.Sprintf("localhost:%v", serverPort)
	u := url.URL{Scheme: "ws", Host: host, Path: testPath}
	// Wait for connection to be established, then break the connection asynchronously
	go func() {
		<-time.After(200 * time.Millisecond)
		err := s.client.webSocket.connection.Close()
		s.NoError(err)
	}()
	// Connect and wait
	err := s.client.Start(u.String())
	s.NoError(err)
	result := <-newClient
	s.NotNil(result)
	// Wait for internal disconnect
	select {
	case result = <-disconnected:
		s.NotNil(result)
	case <-time.After(1 * time.Second):
		s.Fail("timeout waiting for client disconnect")
	}
}

func (s *WebSocketSuite) TestWebsocketServerConnectionBreak() {
	disconnected := make(chan struct{}, 1)
	s.server = newWebsocketServer(s.T(), nil)
	s.server.SetNewClientHandler(func(ws Channel) {
		s.NotNil(ws)
		conn := s.server.connections[ws.ID()]
		s.NotNil(conn)
		// Simulate connection closed as soon client is connected
		err := conn.connection.Close()
		s.NoError(err)
	})
	s.server.SetDisconnectedClientHandler(func(ws Channel) {
		disconnected <- struct{}{}
	})
	go s.server.Start(serverPort, serverPath)
	time.Sleep(100 * time.Millisecond)

	// Test
	s.client = newWebsocketClient(s.T(), nil)
	host := fmt.Sprintf("localhost:%v", serverPort)
	u := url.URL{Scheme: "ws", Host: host, Path: testPath}
	err := s.client.Start(u.String())
	s.NoError(err)

	select {
	case result := <-disconnected:
		s.NotNil(result)
	case <-time.After(1 * time.Second):
		s.Fail("timeout waiting for server disconnect")
	}
}

func (s *WebSocketSuite) TestValidBasicAuth() {
	var ok bool
	authUsername := "testUsername"
	authPassword := "testPassword"
	// Create self-signed TLS certificate
	// TODO: replace with FiloSottile's lib
	certFilename := "/tmp/cert.pem"
	keyFilename := "/tmp/key.pem"
	err := createTLSCertificate(certFilename, keyFilename, "localhost", nil, nil)
	s.NoError(err)
	defer os.Remove(certFilename)
	defer os.Remove(keyFilename)

	// Create TLS server with self-signed certificate
	tlsServer := NewServer(WithServerTLSConfig(certFilename, keyFilename, nil))
	s.server, ok = tlsServer.(*server)
	s.True(ok)
	// Add basic auth handler
	s.server.SetBasicAuthHandler(func(username string, password string) bool {
		s.Equal(authUsername, username)
		s.Equal(authPassword, password)
		return true
	})
	connected := make(chan struct{}, 1)
	s.server.SetNewClientHandler(func(ws Channel) {
		connected <- struct{}{}
	})
	// Run server
	go s.server.Start(serverPort, serverPath)
	time.Sleep(100 * time.Millisecond)

	// Create TLS client
	certPool := x509.NewCertPool()
	data, err := os.ReadFile(certFilename)
	s.NoError(err)
	ok = certPool.AppendCertsFromPEM(data)
	s.True(ok)
	tlsClient := NewClient(WithClientTLSConfig(&tls.Config{
		RootCAs: certPool,
	}))
	s.client, ok = tlsClient.(*client)
	s.True(ok)
	s.client.SetRequestedSubProtocol(defaultSubProtocol)
	// Add basic auth
	s.client.SetBasicAuth(authUsername, authPassword)
	// Test connection
	host := fmt.Sprintf("localhost:%v", serverPort)
	u := url.URL{Scheme: "wss", Host: host, Path: testPath}
	err = s.client.Start(u.String())
	s.NoError(err)
	result := <-connected
	s.NotNil(result)
	s.True(s.client.IsConnected())
}

func (s *WebSocketSuite) TestInvalidBasicAuth() {
	var ok bool
	authUsername := "testUsername"
	authPassword := "testPassword"
	// Create self-signed TLS certificate
	// TODO: replace with FiloSottile's lib
	certFilename := "/tmp/cert.pem"
	keyFilename := "/tmp/key.pem"
	err := createTLSCertificate(certFilename, keyFilename, "localhost", nil, nil)
	s.NoError(err)
	defer os.Remove(certFilename)
	defer os.Remove(keyFilename)

	// Create TLS server with self-signed certificate
	tlsServer := NewServer(WithServerTLSConfig(certFilename, keyFilename, nil))
	s.server, ok = tlsServer.(*server)
	s.True(ok)
	// Add basic auth handler
	s.server.SetBasicAuthHandler(func(username string, password string) bool {
		validCredentials := authUsername == username && authPassword == password
		s.False(validCredentials)
		return validCredentials
	})
	s.server.SetNewClientHandler(func(ws Channel) {
		// Should never reach this
		s.Fail("no new connection should be received from client!")
	})
	// Run server
	go s.server.Start(serverPort, serverPath)
	time.Sleep(100 * time.Millisecond)

	// Create TLS client
	certPool := x509.NewCertPool()
	data, err := os.ReadFile(certFilename)
	s.NoError(err)
	ok = certPool.AppendCertsFromPEM(data)
	s.True(ok)
	wsClient := NewClient(WithClientTLSConfig(&tls.Config{
		RootCAs: certPool,
	}))
	// Test connection without bssic auth -> error expected
	host := fmt.Sprintf("localhost:%v", serverPort)
	u := url.URL{Scheme: "wss", Host: host, Path: testPath}
	err = wsClient.Start(u.String())
	// Assert HTTP error
	s.Error(err)
	var httpErr HttpConnectionError
	ok = errors.As(err, &httpErr)
	s.True(ok)
	s.Equal(http.StatusUnauthorized, httpErr.HttpCode)
	s.Equal("401 Unauthorized", httpErr.HttpStatus)
	s.Equal("websocket: bad handshake", httpErr.Message)
	s.True(strings.Contains(err.Error(), "http status:"))
	// Add basic auth
	wsClient.SetBasicAuth(authUsername, "invalidPassword")
	// Test connection
	err = wsClient.Start(u.String())
	s.Error(err)
	var httpError HttpConnectionError
	ok = errors.As(err, &httpError)
	s.True(ok)
	s.Equal(http.StatusUnauthorized, httpError.HttpCode)
}

func (s *WebSocketSuite) TestInvalidOriginHeader() {
	s.server = newWebsocketServer(s.T(), func(data []byte) ([]byte, error) {
		s.Fail("no message should be received from client!")
		return nil, nil
	})
	s.server.SetNewClientHandler(func(ws Channel) {
		s.Fail("no new connection should be received from client!")
	})
	go s.server.Start(serverPort, serverPath)
	time.Sleep(100 * time.Millisecond)

	// Test message
	s.client = newWebsocketClient(s.T(), func(data []byte) ([]byte, error) {
		s.Fail("no message should be received from server!")
		return nil, nil
	})
	// Set invalid origin header
	s.client.SetHeaderValue("Origin", "example.org")
	host := fmt.Sprintf("localhost:%v", serverPort)
	u := url.URL{Scheme: "ws", Host: host, Path: testPath}
	// Attempt to connect and expect cross-origin error
	err := s.client.Start(u.String())
	s.Error(err)
	var httpErr HttpConnectionError
	ok := errors.As(err, &httpErr)
	s.True(ok)
	s.Equal(http.StatusForbidden, httpErr.HttpCode)
	s.Equal("websocket: bad handshake", httpErr.Message)
}

func (s *WebSocketSuite) TestCustomOriginHeaderHandler() {
	origin := "example.org"
	connected := make(chan struct{})
	s.server = newWebsocketServer(s.T(), func(data []byte) ([]byte, error) {
		s.Fail("no message should be received from client!")
		return nil, nil
	})
	s.server.SetNewClientHandler(func(ws Channel) {
		connected <- struct{}{}
	})
	s.server.SetCheckOriginHandler(func(r *http.Request) bool {
		return r.Header.Get("Origin") == origin
	})
	go s.server.Start(serverPort, serverPath)
	time.Sleep(100 * time.Millisecond)

	// Test message
	s.client = newWebsocketClient(s.T(), func(data []byte) ([]byte, error) {
		s.Fail("no message should be received from server!")
		return nil, nil
	})
	// Set invalid origin header (not example.org)
	s.client.SetHeaderValue("Origin", "localhost")
	host := fmt.Sprintf("localhost:%v", serverPort)
	u := url.URL{Scheme: "ws", Host: host, Path: testPath}
	// Attempt to connect and expect cross-origin error
	err := s.client.Start(u.String())
	s.Error(err)
	var httpErr HttpConnectionError
	ok := errors.As(err, &httpErr)
	s.True(ok)
	s.Equal(http.StatusForbidden, httpErr.HttpCode)
	s.Equal("websocket: bad handshake", httpErr.Message)

	// Re-attempt with correct header
	s.client.SetHeaderValue("Origin", "example.org")
	err = s.client.Start(u.String())
	s.NoError(err)
	result := <-connected
	s.NotNil(result)
}

func (s *WebSocketSuite) TestCustomCheckClientHandler() {
	invalidTestPath := "/ws/invalid-testws"
	id := path.Base(testPath)
	connected := make(chan struct{})
	s.server = newWebsocketServer(s.T(), func(data []byte) ([]byte, error) {
		s.Fail("no message should be received from client!")
		return nil, nil
	})
	s.server.SetNewClientHandler(func(ws Channel) {
		connected <- struct{}{}
	})
	s.server.SetCheckClientHandler(func(clientId string, r *http.Request) bool {
		return id == clientId
	})
	go s.server.Start(serverPort, serverPath)
	time.Sleep(100 * time.Millisecond)

	// Test message
	s.client = newWebsocketClient(s.T(), func(data []byte) ([]byte, error) {
		s.Fail("no message should be received from server!")
		return nil, nil
	})

	host := fmt.Sprintf("localhost:%v", serverPort)
	// Set invalid client (not /ws/testws)
	u := url.URL{Scheme: "ws", Host: host, Path: invalidTestPath}
	// Attempt to connect and expect invalid client id error
	err := s.client.Start(u.String())
	s.Error(err)
	var httpErr HttpConnectionError
	ok := errors.As(err, &httpErr)
	s.True(ok)
	s.Equal(http.StatusUnauthorized, httpErr.HttpCode)
	s.Equal("websocket: bad handshake", httpErr.Message)

	// Re-attempt with correct client id
	u = url.URL{Scheme: "ws", Host: host, Path: testPath}
	err = s.client.Start(u.String())
	s.NoError(err)
	select {
	case result := <-connected:
		s.NotNil(result)
	case <-time.After(1 * time.Second):
		s.Fail("timeout waiting for client to connect")
	}
}

func (s *WebSocketSuite) TestValidClientTLSCertificate() {
	var ok bool
	// Create self-signed TLS certificate
	clientCertFilename := "/tmp/client.pem"
	clientKeyFilename := "/tmp/client_key.pem"
	// TODO: replace with FiloSottile's lib
	err := createTLSCertificate(clientCertFilename, clientKeyFilename, "localhost", nil, nil)
	s.NoError(err)
	defer os.Remove(clientCertFilename)
	defer os.Remove(clientKeyFilename)
	serverCertFilename := "/tmp/cert.pem"
	serverKeyFilename := "/tmp/key.pem"
	err = createTLSCertificate(serverCertFilename, serverKeyFilename, "localhost", nil, nil)
	s.NoError(err)
	defer os.Remove(serverCertFilename)
	defer os.Remove(serverKeyFilename)

	// Create TLS server with self-signed certificate
	certPool := x509.NewCertPool()
	data, err := os.ReadFile(clientCertFilename)
	s.NoError(err)
	ok = certPool.AppendCertsFromPEM(data)
	s.True(ok)
	tlsServer := NewServer(WithServerTLSConfig(serverCertFilename, serverKeyFilename, &tls.Config{
		ClientCAs:  certPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}))
	s.server, ok = tlsServer.(*server)
	s.True(ok)
	// Add basic auth handler
	connected := make(chan struct{})
	s.server.SetNewClientHandler(func(ws Channel) {
		connected <- struct{}{}
	})
	// Run server
	go s.server.Start(serverPort, serverPath)
	time.Sleep(100 * time.Millisecond)

	// Create TLS client
	certPool = x509.NewCertPool()
	data, err = os.ReadFile(serverCertFilename)
	s.NoError(err)
	ok = certPool.AppendCertsFromPEM(data)
	s.True(ok)
	loadedCert, err := tls.LoadX509KeyPair(clientCertFilename, clientKeyFilename)
	s.NoError(err)
	tlsClient := NewClient(WithClientTLSConfig(&tls.Config{
		RootCAs:      certPool,
		Certificates: []tls.Certificate{loadedCert},
	}))
	s.client, ok = tlsClient.(*client)
	s.True(ok)
	s.client.SetRequestedSubProtocol(defaultSubProtocol)
	// Test connection
	host := fmt.Sprintf("localhost:%v", serverPort)
	u := url.URL{Scheme: "wss", Host: host, Path: testPath}
	err = s.client.Start(u.String())
	s.NoError(err)
	select {
	case result := <-connected:
		s.NotNil(result)
	case <-time.After(1 * time.Second):
		s.Fail("timeout waiting for client to connect")
	}
}

func (s *WebSocketSuite) TestInvalidClientTLSCertificate() {
	var ok bool
	// Create self-signed TLS certificate
	clientCertFilename := "/tmp/client.pem"
	clientKeyFilename := "/tmp/client_key.pem"
	// TODO: replace with FiloSottile's lib
	err := createTLSCertificate(clientCertFilename, clientKeyFilename, "localhost", nil, nil)
	s.NoError(err)
	defer os.Remove(clientCertFilename)
	defer os.Remove(clientKeyFilename)
	serverCertFilename := "/tmp/cert.pem"
	serverKeyFilename := "/tmp/key.pem"
	err = createTLSCertificate(serverCertFilename, serverKeyFilename, "localhost", nil, nil)
	s.NoError(err)
	defer os.Remove(serverCertFilename)
	defer os.Remove(serverKeyFilename)

	// Create TLS server with self-signed certificate
	certPool := x509.NewCertPool()
	data, err := os.ReadFile(serverCertFilename)
	s.NoError(err)
	ok = certPool.AppendCertsFromPEM(data)
	s.True(ok)
	tlsServer := NewServer(WithServerTLSConfig(serverCertFilename, serverKeyFilename, &tls.Config{
		ClientCAs:  certPool,                       // Contains server certificate as allowed client CA
		ClientAuth: tls.RequireAndVerifyClientCert, // Requires client certificate signed by allowed CA (server)
	}))
	s.server, ok = tlsServer.(*server)
	s.True(ok)
	// Add basic auth handler
	s.server.SetNewClientHandler(func(ws Channel) {
		// Should never reach this
		s.Fail("no new connection should be received from client!")
	})
	// Run server
	go s.server.Start(serverPort, serverPath)
	time.Sleep(100 * time.Millisecond)

	// Create TLS client
	certPool = x509.NewCertPool()
	data, err = os.ReadFile(serverCertFilename)
	s.NoError(err)
	ok = certPool.AppendCertsFromPEM(data)
	s.True(ok)
	loadedCert, err := tls.LoadX509KeyPair(clientCertFilename, clientKeyFilename)
	s.NoError(err)
	tlsClient := NewClient(WithClientTLSConfig(&tls.Config{
		RootCAs:      certPool,                      // Contains server certificate as allowed server CA
		Certificates: []tls.Certificate{loadedCert}, // Contains self-signed client certificate. Will be rejected by server
	}))
	s.client, ok = tlsClient.(*client)
	s.True(ok)
	s.client.SetRequestedSubProtocol(defaultSubProtocol)
	// Test connection
	host := fmt.Sprintf("localhost:%v", serverPort)
	u := url.URL{Scheme: "wss", Host: host, Path: testPath}
	err = s.client.Start(u.String())
	s.Error(err)
	var netError net.Error
	ok = errors.As(err, &netError)
	s.True(ok)
	s.Equal("remote error: tls: unknown certificate authority", netError.Error()) // tls.alertUnknownCA = 48
}

func (s *WebSocketSuite) TestUnsupportedSubProtocol() {
	s.server.SetNewClientHandler(func(ws Channel) {
	})
	s.server.SetDisconnectedClientHandler(func(ws Channel) {
	})
	s.server.AddSupportedSubprotocol(defaultSubProtocol)
	s.Len(s.server.upgrader.Subprotocols, 1)
	// Test duplicate sub-protocol
	s.server.AddSupportedSubprotocol(defaultSubProtocol)
	s.Len(s.server.upgrader.Subprotocols, 1)
	// Start server
	go s.server.Start(serverPort, serverPath)
	time.Sleep(100 * time.Millisecond)

	// Setup client
	disconnectC := make(chan struct{})
	s.client.SetDisconnectedHandler(func(err error) {
		var wsErr *websocket.CloseError
		ok := s.ErrorAs(err, &wsErr)
		s.True(ok)
		s.Equal(websocket.CloseProtocolError, wsErr.Code)
		s.Equal("invalid or unsupported subprotocol", wsErr.Text)
		s.client.SetDisconnectedHandler(nil)
		close(disconnectC)
	})
	// Set invalid sub-protocol
	s.client.AddOption(func(dialer *websocket.Dialer) {
		dialer.Subprotocols = []string{"unsupportedSubProto"}
	})
	// Test
	host := fmt.Sprintf("localhost:%v", serverPort)
	u := url.URL{Scheme: "ws", Host: host, Path: testPath}
	err := s.client.Start(u.String())
	s.NoError(err)
	// Expect connection to be closed directly after start
	select {
	case _, ok := <-disconnectC:
		s.False(ok)
	case <-time.After(1 * time.Second):
		s.Fail("timeout waiting for client disconnect")
	}
}

func (s *WebSocketSuite) TestSetServerTimeoutConfig() {
	disconnected := make(chan struct{})
	s.server.SetNewClientHandler(func(ws Channel) {
	})
	s.server.SetDisconnectedClientHandler(func(ws Channel) {
		// TODO: check for error with upcoming API
		close(disconnected)
	})
	// Setting server timeout
	config := NewServerTimeoutConfig()
	pingWait := 400 * time.Millisecond
	writeWait := 500 * time.Millisecond
	config.PingWait = pingWait
	config.WriteWait = writeWait
	s.server.SetTimeoutConfig(config)
	// Start server
	go s.server.Start(serverPort, serverPath)
	time.Sleep(100 * time.Millisecond)
	s.Equal(s.server.timeoutConfig.PingWait, pingWait)
	s.Equal(s.server.timeoutConfig.WriteWait, writeWait)
	// Run test
	host := fmt.Sprintf("localhost:%v", serverPort)
	u := url.URL{Scheme: "ws", Host: host, Path: testPath}
	err := s.client.Start(u.String())
	s.NoError(err)
	select {
	case _, ok := <-disconnected:
		s.False(ok)
	case <-time.After(1 * time.Second):
		s.Fail("timeout waiting for client disconnect")
	}
}

func (s *WebSocketSuite) TestSetClientTimeoutConfig() {
	disconnected := make(chan struct{})
	s.server.SetNewClientHandler(func(ws Channel) {
	})
	s.server.SetDisconnectedClientHandler(func(ws Channel) {
		// TODO: check for error with upcoming API
		close(disconnected)
	})
	// Start server
	go s.server.Start(serverPort, serverPath)
	time.Sleep(100 * time.Millisecond)
	// Run test
	host := fmt.Sprintf("localhost:%v", serverPort)
	u := url.URL{Scheme: "ws", Host: host, Path: testPath}
	// Set client timeout
	config := NewClientTimeoutConfig()
	handshakeTimeout := 1 * time.Nanosecond // Very low timeout, handshake will fail
	writeWait := 1 * time.Second
	// Ping period > pong wait, this is a nonsensical config that will trigger a pong timeout
	pingPeriod := 3 * time.Second
	pongWait := 500 * time.Millisecond
	config.PongWait = pongWait
	config.HandshakeTimeout = handshakeTimeout
	config.WriteWait = writeWait
	config.PingPeriod = pingPeriod
	s.client.SetTimeoutConfig(config)
	// Start client and expect handshake error
	err := s.client.Start(u.String())
	var opError *net.OpError
	ok := s.ErrorAs(err, &opError)
	s.True(ok)
	s.Equal("dial", opError.Op)
	s.True(opError.Timeout())
	s.Error(opError.Err, "i/o timeout")
	// Reset handshake to reasonable value
	config.HandshakeTimeout = defaultHandshakeTimeout
	s.client.SetTimeoutConfig(config)
	// Start client
	err = s.client.Start(u.String())
	s.NoError(err)
	s.Equal(s.client.timeoutConfig.PongWait, pongWait)
	s.Equal(s.client.timeoutConfig.WriteWait, writeWait)
	s.Equal(s.client.timeoutConfig.PingPeriod, pingPeriod)
	select {
	case _, closed := <-disconnected:
		s.False(closed)
	case <-time.After(1 * time.Second):
		s.Fail("timeout waiting for client disconnect")
	}
}

func (s *WebSocketSuite) TestServerErrors() {
	triggerC := make(chan struct{}, 1)
	finishC := make(chan struct{}, 1)
	defer close(finishC)
	s.server.SetNewClientHandler(func(ws Channel) {
		triggerC <- struct{}{}
	})
	// Intercept errors asynchronously
	s.Nil(s.server.errC)
	go func() {
		for {
			select {
			case err, ok := <-s.server.Errors():
				triggerC <- struct{}{}
				if ok {
					s.Error(err)
				}
			case <-finishC:
				return
			}
		}
	}()
	s.server.SetMessageHandler(func(ws Channel, data []byte) error {
		return fmt.Errorf("this is a dummy error")
	})
	// Will trigger an out-of-bound error
	time.Sleep(50 * time.Millisecond)
	s.server.Stop()
	r := <-triggerC
	s.NotNil(r)
	// Start server for real
	s.server.httpServer = &http.Server{}
	go s.server.Start(serverPort, serverPath)
	time.Sleep(100 * time.Millisecond)
	// Connect client
	host := fmt.Sprintf("localhost:%v", serverPort)
	u := url.URL{Scheme: "ws", Host: host, Path: testPath}
	err := s.client.Start(u.String())
	s.NoError(err)
	// Wait for new client callback
	r = <-triggerC
	s.NotNil(r)
	// Send a dummy message and expect error on server side
	err = s.client.Write([]byte("dummy message"))
	s.NoError(err)
	r = <-triggerC
	s.NotNil(r)
	// Send message to non-existing client
	err = s.server.Write("fakeId", []byte("dummy response"))
	s.Error(err)
	// Send unexpected close message and wait for error to be thrown
	err = s.client.webSocket.connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseUnsupportedData, ""))
	s.NoError(err)
	r = <-triggerC
	s.NotNil(r)
	// Stop and wait for errors channel cleanup
	s.server.Stop()
	r = <-triggerC
	s.NotNil(r)
}

func (s *WebSocketSuite) TestClientErrors() {
	triggerC := make(chan struct{}, 1)
	finishC := make(chan struct{}, 1)
	defer close(finishC)
	s.server.SetNewClientHandler(func(ws Channel) {
		triggerC <- struct{}{}
	})
	s.client.SetMessageHandler(func(data []byte) error {
		return fmt.Errorf("this is a dummy error")
	})
	// Intercept errors asynchronously
	s.Nil(s.client.errC)
	go func() {
		for {
			select {
			case err, ok := <-s.client.Errors():
				triggerC <- struct{}{}
				if ok {
					s.Error(err)
				}
			case <-finishC:
				return
			}
		}
	}()
	go s.server.Start(serverPort, serverPath)
	time.Sleep(100 * time.Millisecond)
	// Attempt to write a message without being connected
	err := s.client.Write([]byte("dummy message"))
	s.Error(err)
	// Connect client
	host := fmt.Sprintf("localhost:%v", serverPort)
	u := url.URL{Scheme: "ws", Host: host, Path: testPath}
	err = s.client.Start(u.String())
	s.NoError(err)
	// Wait for new client callback
	r := <-triggerC
	s.NotNil(r)
	// Send a dummy message and expect error on client side
	err = s.server.Write(path.Base(testPath), []byte("dummy message"))
	s.NoError(err)
	r = <-triggerC
	s.NotNil(r)
	// Send unexpected close message and wait for error to be thrown
	conn := s.server.connections[path.Base(testPath)]
	s.NotNil(conn)
	err = conn.WriteManual(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseUnsupportedData, ""))
	// err = conn.connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseUnsupportedData, ""))
	s.NoError(err)
	r = <-triggerC
	s.NotNil(r)
	// Stop server and client and wait for errors channel cleanup
	s.server.Stop()
	s.client.Stop()
	r = <-triggerC
	s.NotNil(r)
}

// Utility functions

func createCACertificate(certificateFilename string, keyFilename string) (*x509.Certificate, *ecdsa.PrivateKey, error) {
	// Generate ed25519 key-pair
	privateKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	// Create CA
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, nil, err
	}
	notBefore := time.Now()
	notAfter := notBefore.Add(time.Hour * 24)
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"ocpp-go"},
			CommonName:   "ocpp-go-CA",
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost"},
	}
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, nil, err
	}
	// Save certificate to disk
	certOut, err := os.Create(certificateFilename)
	if err != nil {
		return nil, nil, err
	}
	defer certOut.Close()
	err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if err != nil {
		return nil, nil, err
	}
	// Save key to disk
	keyOut, err := os.Create(keyFilename)
	if err != nil {
		return nil, nil, err
	}
	defer keyOut.Close()
	privateBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, nil, err
	}
	err = pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privateBytes})
	if err != nil {
		return nil, nil, err
	}
	return &template, privateKey, nil
}

func createTLSCertificate(certificateFilename string, keyFilename string, cn string, ca *x509.Certificate, caKey *ecdsa.PrivateKey) error {
	// Generate ed25519 key-pair
	privateKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return err
	}
	// Create self-signed certificate
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return err
	}
	notBefore := time.Now()
	notAfter := notBefore.Add(time.Hour * 24)
	keyUsage := x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature
	extKeyUsage := []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth}
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"ocpp-go"},
			CommonName:   cn,
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              keyUsage,
		ExtKeyUsage:           extKeyUsage,
		BasicConstraintsValid: true,
		DNSNames:              []string{cn},
	}
	var derBytes []byte
	if ca != nil && caKey != nil {
		// Certificate signed by CA
		derBytes, err = x509.CreateCertificate(rand.Reader, &template, ca, &privateKey.PublicKey, caKey)
	} else {
		// Self-signed certificate
		derBytes, err = x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	}
	if err != nil {
		return err
	}
	// Save certificate to disk
	certOut, err := os.Create(certificateFilename)
	if err != nil {
		return err
	}
	defer certOut.Close()
	err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if err != nil {
		return err
	}
	// Save key to disk
	keyOut, err := os.Create(keyFilename)
	if err != nil {
		return err
	}
	defer keyOut.Close()
	privateBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return err
	}
	err = pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privateBytes})
	if err != nil {
		return err
	}
	return nil
}

func TestWebSockets(t *testing.T) {
	suite.Run(t, new(WebSocketSuite))
}
