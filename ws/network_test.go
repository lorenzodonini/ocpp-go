package ws

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	toxiproxy "github.com/Shopify/toxiproxy/client"
)

type NetworkTestSuite struct {
	suite.Suite
	proxy     *toxiproxy.Proxy
	proxyPort int
	server    *Server
	client    *Client
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

func (s *NetworkTestSuite) SetupTest() {
	s.server = NewWebsocketServer(s.T(), nil)
	s.client = NewWebsocketClient(s.T(), nil)
}

func (s *NetworkTestSuite) TearDownTest() {
	s.server = nil
	s.client = nil
}

func (s *NetworkTestSuite) TestClientConnectionFailed() {
	t := s.T()
	s.server = NewWebsocketServer(t, nil)
	s.server.SetNewClientHandler(func(ws Channel) {
		assert.Fail(t, "should not accept new clients")
	})
	go s.server.Start(serverPort, serverPath)
	time.Sleep(500 * time.Millisecond)

	// Test client
	host := s.proxy.Listen
	u := url.URL{Scheme: "ws", Host: host, Path: testPath}

	// Disable network
	_ = s.proxy.Disable()
	defer s.proxy.Enable()
	// Attempt connection
	err := s.client.Start(u.String())
	require.Error(t, err)
	netError, ok := err.(*net.OpError)
	require.True(t, ok)
	require.NotNil(t, netError.Err)
	sysError, ok := netError.Err.(*os.SyscallError)
	require.True(t, ok)
	assert.Equal(t, "connect", sysError.Syscall)
	assert.Equal(t, syscall.ECONNREFUSED, sysError.Err)
	// Cleanup
	s.server.Stop()
}

func (s *NetworkTestSuite) TestClientConnectionFailedTimeout() {
	t := s.T()
	// Set timeouts for test
	s.client.timeoutConfig.HandshakeTimeout = 2 * time.Second
	// Setup
	s.server = NewWebsocketServer(t, nil)
	s.server.SetNewClientHandler(func(ws Channel) {
		assert.Fail(t, "should not accept new clients")
	})
	go s.server.Start(serverPort, serverPath)
	time.Sleep(500 * time.Millisecond)

	// Test client
	host := s.proxy.Listen
	u := url.URL{Scheme: "ws", Host: host, Path: testPath}

	// Add connection timeout
	_, err := s.proxy.AddToxic("connectTimeout", "timeout", "upstream", 1, toxiproxy.Attributes{
		"timeout": 3000, // 3 seconds
	})
	defer s.proxy.RemoveToxic("connectTimeout")
	require.NoError(t, err)
	// Attempt connection
	err = s.client.Start(u.String())
	require.Error(t, err)
	netError, ok := err.(*net.OpError)
	require.True(t, ok)
	require.NotNil(t, netError.Err)
	assert.True(t, strings.Contains(netError.Error(), "timeout"))
	assert.True(t, netError.Timeout())
	// Cleanup
	s.server.Stop()
}

func (s *NetworkTestSuite) TestClientAutoReconnect() {
	t := s.T()
	// Set timeouts for test
	s.client.timeoutConfig.ReconnectBackoff = 1 * time.Second
	// Setup
	serverOnDisconnected := make(chan bool, 1)
	clientOnDisconnected := make(chan bool, 1)
	reconnected := make(chan bool, 1)
	s.server = NewWebsocketServer(t, nil)
	s.server.SetNewClientHandler(func(ws Channel) {
		assert.NotNil(t, ws)
		conn := s.server.connections[ws.GetID()]
		require.NotNil(t, conn)
	})
	s.server.SetDisconnectedClientHandler(func(ws Channel) {
		serverOnDisconnected <- true
	})
	go s.server.Start(serverPort, serverPath)
	time.Sleep(500 * time.Millisecond)

	// Test bench
	s.client.SetDisconnectedHandler(func(err error) {
		assert.NotNil(t, err)
		closeError, ok := err.(*websocket.CloseError)
		require.True(t, ok)
		assert.Equal(t, websocket.CloseAbnormalClosure, closeError.Code)
		assert.False(t, s.client.IsConnected())
		clientOnDisconnected <- true
	})
	s.client.SetReconnectedHandler(func() {
		time.Sleep(50 * time.Millisecond) // Make sure we reconnected after backoff
		reconnected <- true
	})
	// Connect client
	host := s.proxy.Listen
	u := url.URL{Scheme: "ws", Host: host, Path: testPath}
	err := s.client.Start(u.String())
	require.Nil(t, err)
	// Close all connection from server side
	time.Sleep(500 * time.Millisecond)
	for _, s := range s.server.connections {
		err = s.connection.Close()
		require.Nil(t, err)
	}
	// Wait for disconnect to propagate
	result := <-serverOnDisconnected
	require.True(t, result)
	result = <-clientOnDisconnected
	require.True(t, result)
	start := time.Now()
	// Wait for reconnection
	result = <-reconnected
	elapsed := time.Since(start)
	assert.True(t, result)
	assert.True(t, s.client.IsConnected())
	assert.GreaterOrEqual(t, elapsed.Milliseconds(), s.client.timeoutConfig.ReconnectBackoff.Milliseconds())
	// Cleanup
	s.client.Stop()
	s.server.Stop()
}

func (s *NetworkTestSuite) TestClientPongTimeout() {
	t := s.T()
	// Set timeouts for test
	// Will attempt to send ping after 1 second, and server expects ping within 1.4 seconds
	// Server will close connection
	s.client.timeoutConfig.PongWait = 2 * time.Second
	s.client.timeoutConfig.PingPeriod = (s.client.timeoutConfig.PongWait * 5) / 10
	s.client.timeoutConfig.ReconnectBackoff = 1 * time.Second
	s.server.timeoutConfig.PingWait = (s.client.timeoutConfig.PongWait * 7) / 10
	// Setup
	serverOnDisconnected := make(chan bool, 1)
	clientOnDisconnected := make(chan bool, 1)
	reconnected := make(chan bool, 1)
	s.server.SetNewClientHandler(func(ws Channel) {
		assert.NotNil(t, ws)
	})
	s.server.SetDisconnectedClientHandler(func(ws Channel) {
		serverOnDisconnected <- true
	})
	s.server.SetMessageHandler(func(ws Channel, data []byte) error {
		assert.Fail(t, "unexpected message received")
		return fmt.Errorf("unexpected message received")
	})
	go s.server.Start(serverPort, serverPath)
	time.Sleep(500 * time.Millisecond)

	// Test client
	s.client.SetDisconnectedHandler(func(err error) {
		defer func() {
			clientOnDisconnected <- true
		}()
		require.Error(t, err)
		closeError, ok := err.(*websocket.CloseError)
		require.True(t, ok)
		assert.Equal(t, websocket.CloseAbnormalClosure, closeError.Code)
	})
	s.client.SetReconnectedHandler(func() {
		time.Sleep(50 * time.Millisecond) // Make sure we reconnected after backoff
		reconnected <- true
	})
	host := s.proxy.Listen
	u := url.URL{Scheme: "ws", Host: host, Path: testPath}

	// Attempt connection
	err := s.client.Start(u.String())
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
	assert.GreaterOrEqual(t, elapsed.Milliseconds(), s.client.timeoutConfig.ReconnectBackoff.Milliseconds())
	// Cleanup
	s.client.Stop()
	s.server.Stop()
}

func (s *NetworkTestSuite) TestClientReadTimeout() {
	t := s.T()
	// Set timeouts for test
	s.client.timeoutConfig.PongWait = 2 * time.Second
	s.client.timeoutConfig.PingPeriod = (s.client.timeoutConfig.PongWait * 7) / 10
	s.client.timeoutConfig.ReconnectBackoff = 1 * time.Second
	s.server.timeoutConfig.PingWait = s.client.timeoutConfig.PongWait
	// Setup
	serverOnDisconnected := make(chan bool, 1)
	clientOnDisconnected := make(chan bool, 1)
	reconnected := make(chan bool, 1)
	s.server.SetNewClientHandler(func(ws Channel) {
		assert.NotNil(t, ws)
	})
	s.server.SetDisconnectedClientHandler(func(ws Channel) {
		serverOnDisconnected <- true
	})
	s.server.SetMessageHandler(func(ws Channel, data []byte) error {
		assert.Fail(t, "unexpected message received")
		return fmt.Errorf("unexpected message received")
	})
	go s.server.Start(serverPort, serverPath)
	time.Sleep(500 * time.Millisecond)

	// Test client
	s.client.SetDisconnectedHandler(func(err error) {
		defer func() {
			clientOnDisconnected <- true
		}()
		require.Error(t, err)
		errMsg := err.Error()
		c := strings.Contains(errMsg, "timeout")
		if !c {
			fmt.Println(errMsg)
		}
		assert.True(t, c)
	})
	s.client.SetReconnectedHandler(func() {
		time.Sleep(50 * time.Millisecond) // Make sure we reconnected after backoff
		reconnected <- true
	})
	host := s.proxy.Listen
	u := url.URL{Scheme: "ws", Host: host, Path: testPath}

	// Attempt connection
	err := s.client.Start(u.String())
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
	assert.GreaterOrEqual(t, elapsed.Milliseconds(), s.client.timeoutConfig.ReconnectBackoff.Milliseconds())
	// Cleanup
	s.client.Stop()
	s.server.Stop()
}

//TODO: test error channel from websocket

func TestNetworkErrors(t *testing.T) {
	suite.Run(t, new(NetworkTestSuite))
}
