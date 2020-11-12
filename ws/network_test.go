package ws

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"net"
	"net/url"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/Shopify/toxiproxy/client"
)

type NetworkTestSuite struct {
	suite.Suite
	proxy     *toxiproxy.Proxy
	proxyPort int
}

func (s *NetworkTestSuite) SetupSuite() {
	client := toxiproxy.NewClient("localhost:8474")
	s.proxyPort = 8886
	// Proxy listens on 8886 and upstreams to 8887 (where ocpp server is actually listening)
	oldProxy, err := client.Proxy("ocpp")
	if oldProxy != nil {
		oldProxy.Delete()
	}
	p, err := client.CreateProxy("ocpp", "localhost:8886", fmt.Sprintf("localhost:%v", serverPort))
	require.NoError(s.T(), err)
	s.proxy = p
}

func (s *NetworkTestSuite) TearDownSuite() {
	s.proxy.Delete()
}

func (s *NetworkTestSuite) TearDownTest() {
	// Reset websocket timeouts
	pongWait = defaultPongWait
	pingWait = defaultPingWait
	pingPeriod = defaultPingPeriod
	autoReconnectDelay = defaultAutoReconnectDelay
	maxReconnectionDelay = defaultMaxReconnectionDelay
}

func (s *NetworkTestSuite) TestClientConnectionFailed() {
	t := s.T()
	var wsServer *Server
	wsServer = NewWebsocketServer(t, nil)
	wsServer.SetNewClientHandler(func(ws Channel) {
		assert.Fail(t, "should not accept new clients")
	})
	go wsServer.Start(serverPort, serverPath)
	time.Sleep(1 * time.Second)

	// Test client
	wsClient := NewWebsocketClient(t, nil)
	host := s.proxy.Listen
	u := url.URL{Scheme: "ws", Host: host, Path: testPath}

	// Disable network
	_ = s.proxy.Disable()
	defer s.proxy.Enable()
	// Attempt connection
	err := wsClient.Start(u.String())
	require.Error(t, err)
	netError, ok := err.(*net.OpError)
	require.True(t, ok)
	require.NotNil(t, netError.Err)
	sysError, ok := netError.Err.(*os.SyscallError)
	require.True(t, ok)
	assert.Equal(t, "connect", sysError.Syscall)
	assert.Equal(t, syscall.ECONNREFUSED, sysError.Err)
	// Cleanup
	wsServer.Stop()
}

func (s *NetworkTestSuite) TestClientConnectionFailedTimeout() {
	t := s.T()
	// Set timeouts for test
	handshakeTimeout = 2 * time.Second
	// Setup
	var wsServer *Server
	wsServer = NewWebsocketServer(t, nil)
	wsServer.SetNewClientHandler(func(ws Channel) {
		assert.Fail(t, "should not accept new clients")
	})
	go wsServer.Start(serverPort, serverPath)
	time.Sleep(1 * time.Second)

	// Test client
	wsClient := NewWebsocketClient(t, nil)
	host := s.proxy.Listen
	u := url.URL{Scheme: "ws", Host: host, Path: testPath}

	// Add connection timeout
	_, err := s.proxy.AddToxic("connectTimeout", "timeout", "upstream", 1, toxiproxy.Attributes{
		"timeout": 3000, // 3 seconds
	})
	defer s.proxy.RemoveToxic("connectTimeout")
	require.NoError(t, err)
	// Attempt connection
	err = wsClient.Start(u.String())
	require.Error(t, err)
	netError, ok := err.(*net.OpError)
	require.True(t, ok)
	require.NotNil(t, netError.Err)
	assert.True(t, strings.Contains(netError.Error(), "timeout"))
	assert.True(t, netError.Timeout())
	// Cleanup
	wsServer.Stop()
}

func (s *NetworkTestSuite) TestClientAutoReconnect() {
	t := s.T()
	// Set timeouts for test
	autoReconnectDelay = 1 * time.Second
	// Setup
	var wsServer *Server
	serverOnDisconnected := make(chan bool, 1)
	clientOnDisconnected := make(chan bool, 1)
	reconnected := make(chan bool, 1)
	wsServer = NewWebsocketServer(t, nil)
	wsServer.SetNewClientHandler(func(ws Channel) {
		assert.NotNil(t, ws)
		conn := wsServer.connections[ws.GetID()]
		assert.NotNil(t, conn)
		// Simulate connection closed as soon client is connected
		err := conn.connection.Close()
		assert.Nil(t, err)
	})
	wsServer.SetDisconnectedClientHandler(func(ws Channel) {
		serverOnDisconnected <- true
	})
	go wsServer.Start(serverPort, serverPath)
	time.Sleep(1 * time.Second)

	// Test bench
	wsClient := NewWebsocketClient(t, nil)
	wsClient.SetDisconnectedHandler(func(err error) {
		assert.NotNil(t, err)
		closeError, ok := err.(*websocket.CloseError)
		require.True(t, ok)
		assert.Equal(t, websocket.CloseAbnormalClosure, closeError.Code)
		assert.False(t, wsClient.IsConnected())
		clientOnDisconnected <- true
	})
	wsClient.SetReconnectedHandler(func() {
		reconnected <- true
	})
	// Connect client
	host := s.proxy.Listen
	u := url.URL{Scheme: "ws", Host: host, Path: testPath}
	err := wsClient.Start(u.String())
	require.Nil(t, err)
	result := <-serverOnDisconnected
	require.True(t, result)
	result = <-clientOnDisconnected
	require.True(t, result)
	start := time.Now()
	// Wait for reconnection
	result = <-reconnected
	elapsed := time.Since(start)
	assert.True(t, result)
	assert.True(t, wsClient.IsConnected())
	assert.True(t, elapsed >= autoReconnectDelay)
	// Cleanup
	wsClient.Stop()
	wsServer.Stop()
}

func (s *NetworkTestSuite) TestClientPongTimeout() {
	t := s.T()
	// Set timeouts for test
	// Will attempt to send ping after 1 second, and server expects ping within 1.4 seconds
	// Server will close connection
	pongWait = 2 * time.Second
	pingWait = (pongWait * 7) / 10
	pingPeriod = (pongWait * 5) / 10
	autoReconnectDelay = 1 * time.Second
	// Setup
	var wsServer *Server
	serverOnDisconnected := make(chan bool, 1)
	clientOnDisconnected := make(chan bool, 1)
	reconnected := make(chan bool, 1)
	wsServer = NewWebsocketServer(t, nil)
	wsServer.SetNewClientHandler(func(ws Channel) {
		assert.NotNil(t, ws)
	})
	wsServer.SetDisconnectedClientHandler(func(ws Channel) {
		serverOnDisconnected <- true
	})
	wsServer.SetMessageHandler(func(ws Channel, data []byte) error {
		assert.Fail(t, "unexpected message received")
		return errors.New("unexpected message received")
	})
	go wsServer.Start(serverPort, serverPath)
	time.Sleep(1 * time.Second)

	// Test client
	wsClient := NewWebsocketClient(t, nil)
	wsClient.SetDisconnectedHandler(func(err error) {
		defer func() {
			clientOnDisconnected <- true
		}()
		require.Error(t, err)
		closeError, ok := err.(*websocket.CloseError)
		require.True(t, ok)
		assert.Equal(t, websocket.CloseAbnormalClosure, closeError.Code)
	})
	wsClient.SetReconnectedHandler(func() {
		reconnected <- true
	})
	host := s.proxy.Listen
	u := url.URL{Scheme: "ws", Host: host, Path: testPath}

	// Attempt connection
	err := wsClient.Start(u.String())
	require.NoError(t, err)
	// Slow upstream network -> ping won't get through and server-side close will be triggered
	_, err = s.proxy.AddToxic("readTimeout", "timeout", "upstream", 1, toxiproxy.Attributes{
		"timeout": 5000, // 5 seconds
	})
	require.NoError(t, err)
	// Attempt to send message
	require.NoError(t, err)
	result := <-clientOnDisconnected
	require.True(t, result)
	result = <-serverOnDisconnected
	require.True(t, result)
	// Reconnect time starts
	s.proxy.RemoveToxic("readTimeout")
	startTimeout := time.Now()
	result = <-reconnected
	require.True(t, result)
	elapsed := time.Since(startTimeout)
	assert.GreaterOrEqual(t, elapsed.Milliseconds(), autoReconnectDelay.Milliseconds())
	// Cleanup
	wsClient.Stop()
	wsServer.Stop()
}

func (s *NetworkTestSuite) TestClientReadTimeout() {
	t := s.T()
	// Set timeouts for test
	pongWait = 2 * time.Second
	pingWait = pongWait
	pingPeriod = (pongWait * 8) / 10
	autoReconnectDelay = 1 * time.Second
	// Setup
	var wsServer *Server
	serverOnDisconnected := make(chan bool, 1)
	clientOnDisconnected := make(chan bool, 1)
	reconnected := make(chan bool, 1)
	wsServer = NewWebsocketServer(t, nil)
	wsServer.SetNewClientHandler(func(ws Channel) {
		assert.NotNil(t, ws)
	})
	wsServer.SetDisconnectedClientHandler(func(ws Channel) {
		serverOnDisconnected <- true
	})
	wsServer.SetMessageHandler(func(ws Channel, data []byte) error {
		assert.Fail(t, "unexpected message received")
		return errors.New("unexpected message received")
	})
	go wsServer.Start(serverPort, serverPath)
	time.Sleep(1 * time.Second)

	// Test client
	wsClient := NewWebsocketClient(t, nil)
	wsClient.SetDisconnectedHandler(func(err error) {
		defer func() {
			clientOnDisconnected <- true
		}()
		require.Error(t, err)
		errMsg := err.Error()
		c := strings.Contains(errMsg, "timeout")
		if !c {
			fmt.Println(errMsg)
		}
		//TODO: not deterministic. Sometimes abnormal closure, sometimes timeout
		assert.True(t, c)
	})
	wsClient.SetReconnectedHandler(func() {
		reconnected <- true
	})
	host := s.proxy.Listen
	u := url.URL{Scheme: "ws", Host: host, Path: testPath}

	// Attempt connection
	err := wsClient.Start(u.String())
	require.NoError(t, err)
	// Slow down network. Ping will be received but pong won't go through
	_, err = s.proxy.AddToxic("writeTimeout", "timeout", "downstream", 1, toxiproxy.Attributes{
		"timeout": 5000, // 5 seconds
	})
	require.NoError(t, err)
	// Attempt to send message
	require.NoError(t, err)
	result := <-serverOnDisconnected
	require.True(t, result)
	result = <-clientOnDisconnected
	require.True(t, result)
	// Reconnect time starts
	s.proxy.RemoveToxic("writeTimeout")
	startTimeout := time.Now()
	result = <-reconnected
	require.True(t, result)
	elapsed := time.Since(startTimeout)
	assert.GreaterOrEqual(t, elapsed.Milliseconds(), autoReconnectDelay.Milliseconds())
	// Cleanup
	wsClient.Stop()
	wsServer.Stop()
}

func TestNetworkErrors(t *testing.T) {
	suite.Run(t, new(NetworkTestSuite))
}
