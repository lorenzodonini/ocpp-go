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

	"github.com/caarlos0/env/v11"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/suite"

	toxiproxy "github.com/Shopify/toxiproxy/client"
)

type proxyConfig struct {
	ToxiProxyHost     string `env:"TOXIPROXY_HOST" envDefault:"localhost"`
	ToxiProxyPort     int    `env:"TOXIPROXY_PORT" envDefault:"8474"`
	ProxyOcppListener string `env:"PROXY_OCPP_LISTENER" envDefault:"localhost:8886"`
	ProxyOcppUpstream string `env:"PROXY_OCPP_UPSTREAM" envDefault:"localhost:8887"`
}

type NetworkTestSuite struct {
	suite.Suite
	proxy     *toxiproxy.Proxy
	proxyPort int
	server    *server
	client    *client
}

func (s *NetworkTestSuite) SetupSuite() {
	var cfg proxyConfig
	err := env.Parse(&cfg)
	s.Require().NoError(err)

	endpoint := fmt.Sprintf("%v:%v", cfg.ToxiProxyHost, cfg.ToxiProxyPort)
	client := toxiproxy.NewClient(endpoint)
	s.proxyPort = 8886
	// Proxy listens on 8886 and upstreams to 8887 (where ocpp server is actually listening)
	oldProxy, err := client.Proxy("ocpp")
	if oldProxy != nil {
		s.Require().NoError(oldProxy.Delete())
	}
	p, err := client.CreateProxy("ocpp", cfg.ProxyOcppListener, cfg.ProxyOcppUpstream)
	s.NoError(err)
	s.proxy = p
}

func (s *NetworkTestSuite) TearDownSuite() {
	s.Require().NoError(s.proxy.Delete())
}

func (s *NetworkTestSuite) SetupTest() {
	s.server = newWebsocketServer(s.T(), nil)
	s.client = newWebsocketClient(s.T(), nil)
}

func (s *NetworkTestSuite) TearDownTest() {
	if s.client != nil {
		s.client.Stop()
	}
	if s.server != nil {
		s.server.Stop()
	}
}

func (s *NetworkTestSuite) TestClientConnectionFailed() {
	s.server.SetNewClientHandler(func(ws Channel) {
		s.Fail("should not accept new clients")
	})
	go s.server.Start(serverPort, serverPath)
	time.Sleep(100 * time.Millisecond)

	// Test client
	host := s.proxy.Listen
	u := url.URL{Scheme: "ws", Host: host, Path: testPath}

	// Disable network
	_ = s.proxy.Disable()
	defer s.proxy.Enable()
	// Attempt connection
	err := s.client.Start(u.String())
	s.Error(err)
	var netError *net.OpError
	ok := s.ErrorAs(err, &netError)
	s.True(ok)
	s.Error(netError.Err)
	var sysError *os.SyscallError
	ok = s.ErrorAs(netError.Err, &sysError)
	s.True(ok)
	s.Equal("connect", sysError.Syscall)
	s.Equal(syscall.ECONNREFUSED, sysError.Err)
}

func (s *NetworkTestSuite) TestClientConnectionFailedTimeout() {
	// Set timeouts for test
	s.client.timeoutConfig.HandshakeTimeout = 2 * time.Second
	// Setup
	s.server.SetNewClientHandler(func(ws Channel) {
		s.Fail("should not accept new clients")
	})
	go s.server.Start(serverPort, serverPath)
	time.Sleep(100 * time.Millisecond)

	// Test client
	host := s.proxy.Listen
	u := url.URL{Scheme: "ws", Host: host, Path: testPath}

	// Add connection timeout
	_, err := s.proxy.AddToxic("connectTimeout", "timeout", "upstream", 1, toxiproxy.Attributes{
		"timeout": 3000, // 3 seconds
	})
	defer s.proxy.RemoveToxic("connectTimeout")
	s.NoError(err)
	// Attempt connection
	err = s.client.Start(u.String())
	s.Error(err)
	var netError *net.OpError
	ok := s.ErrorAs(err, &netError)
	s.True(ok)
	s.Error(netError.Err)
	s.True(strings.Contains(netError.Error(), "timeout"))
	s.True(netError.Timeout())
}

func (s *NetworkTestSuite) TestClientAutoReconnect() {
	// Set timeouts for test
	s.client.timeoutConfig.RetryBackOffWaitMinimum = 1 * time.Second
	s.client.timeoutConfig.RetryBackOffRandomRange = 1 // seconds
	// Setup
	serverOnDisconnected := make(chan struct{}, 1)
	clientOnDisconnected := make(chan struct{}, 1)
	reconnected := make(chan struct{}, 1)
	s.server.SetNewClientHandler(func(ws Channel) {
		s.NotNil(ws)
		conn := s.server.connections[ws.ID()]
		s.NotNil(conn)
	})
	s.server.SetDisconnectedClientHandler(func(ws Channel) {
		serverOnDisconnected <- struct{}{}
	})
	go s.server.Start(serverPort, serverPath)
	time.Sleep(100 * time.Millisecond)

	// Test bench
	s.client.SetDisconnectedHandler(func(err error) {
		s.Error(err)
		var closeError *websocket.CloseError
		ok := s.ErrorAs(err, &closeError)
		s.True(ok)
		s.Equal(websocket.CloseAbnormalClosure, closeError.Code)
		s.False(s.client.IsConnected())
		clientOnDisconnected <- struct{}{}
	})
	s.client.SetReconnectedHandler(func() {
		time.Sleep(time.Duration(s.client.timeoutConfig.RetryBackOffRandomRange)*time.Second + 50*time.Millisecond) // Make sure we reconnected after backoff
		reconnected <- struct{}{}
	})
	// Connect client
	host := s.proxy.Listen
	u := url.URL{Scheme: "ws", Host: host, Path: testPath}
	err := s.client.Start(u.String())
	s.NoError(err)
	// Close all connection from server side
	time.Sleep(100 * time.Millisecond)
	for _, c := range s.server.connections {
		err = c.connection.Close()
		s.NoError(err)
	}
	// Wait for disconnect to propagate
	result := <-serverOnDisconnected
	s.NotNil(result)
	result = <-clientOnDisconnected
	s.NotNil(result)
	start := time.Now()
	// Wait for reconnection
	result = <-reconnected
	elapsed := time.Since(start)
	s.NotNil(result)
	s.True(s.client.IsConnected())
	s.GreaterOrEqual(elapsed.Milliseconds(), s.client.timeoutConfig.RetryBackOffWaitMinimum.Milliseconds())
	// Cleanup
	s.client.SetDisconnectedHandler(func(err error) {
		s.NoError(err)
		clientOnDisconnected <- struct{}{}
	})
	s.client.Stop()
	result = <-clientOnDisconnected
	s.NotNil(result)
	s.server.Stop()
}

func (s *NetworkTestSuite) TestClientPongTimeout() {
	// Set timeouts for test
	// Will attempt to send ping after 1 second, and server expects ping within 1.4 seconds
	// server will close connection
	s.client.timeoutConfig.PongWait = 2 * time.Second
	s.client.timeoutConfig.PingPeriod = (s.client.timeoutConfig.PongWait * 5) / 10
	s.client.timeoutConfig.RetryBackOffWaitMinimum = 1 * time.Second
	s.client.timeoutConfig.RetryBackOffWaitMinimum = 0 // remove randomness
	s.server.timeoutConfig.PingWait = (s.client.timeoutConfig.PongWait * 7) / 10
	// Setup
	serverOnDisconnected := make(chan struct{}, 1)
	clientOnDisconnected := make(chan struct{}, 1)
	reconnected := make(chan struct{}, 1)
	s.server.SetNewClientHandler(func(ws Channel) {
		s.NotNil(ws)
	})
	s.server.SetDisconnectedClientHandler(func(ws Channel) {
		serverOnDisconnected <- struct{}{}
	})
	s.server.SetMessageHandler(func(ws Channel, data []byte) error {
		s.Fail("unexpected message received")
		return fmt.Errorf("unexpected message received")
	})
	go s.server.Start(serverPort, serverPath)
	time.Sleep(100 * time.Millisecond)

	// Test client
	s.client.SetDisconnectedHandler(func(err error) {
		defer func() {
			clientOnDisconnected <- struct{}{}
		}()
		s.Error(err)
		var closeError *websocket.CloseError
		ok := s.ErrorAs(err, &closeError)
		s.True(ok)
		s.Equal(websocket.CloseAbnormalClosure, closeError.Code)
	})
	s.client.SetReconnectedHandler(func() {
		time.Sleep(50 * time.Millisecond) // Make sure we reconnected after backoff
		reconnected <- struct{}{}
	})
	host := s.proxy.Listen
	u := url.URL{Scheme: "ws", Host: host, Path: testPath}

	// Attempt connection
	err := s.client.Start(u.String())
	s.NoError(err)
	// Slow upstream network -> ping won't get through and server-side close will be triggered
	_, err = s.proxy.AddToxic("readTimeout", "timeout", "upstream", 1, toxiproxy.Attributes{
		"timeout": 5000, // 5 seconds
	})
	s.NoError(err)
	// Attempt to send message
	result := <-clientOnDisconnected
	s.NotNil(result)
	result = <-serverOnDisconnected
	s.NotNil(result)
	// Reconnect time starts
	_ = s.proxy.RemoveToxic("readTimeout")
	startTimeout := time.Now()
	result = <-reconnected
	s.NotNil(result)
	elapsed := time.Since(startTimeout)
	s.GreaterOrEqual(elapsed.Milliseconds(), s.client.timeoutConfig.RetryBackOffWaitMinimum.Milliseconds())
	// Cleanup
	s.client.SetDisconnectedHandler(nil)
	s.client.Stop()
	s.server.Stop()
}

func (s *NetworkTestSuite) TestClientReadTimeout() {
	// Set timeouts for test
	s.client.timeoutConfig.PongWait = 2 * time.Second
	s.client.timeoutConfig.PingPeriod = (s.client.timeoutConfig.PongWait * 7) / 10
	s.client.timeoutConfig.RetryBackOffWaitMinimum = 1 * time.Second
	s.client.timeoutConfig.RetryBackOffRandomRange = 0 // remove randomness
	s.server.timeoutConfig.PingWait = s.client.timeoutConfig.PongWait
	// Setup
	serverOnDisconnected := make(chan struct{}, 1)
	clientOnDisconnected := make(chan struct{}, 1)
	reconnected := make(chan struct{}, 1)
	s.server.SetNewClientHandler(func(ws Channel) {
		s.NotNil(ws)
	})
	s.server.SetDisconnectedClientHandler(func(ws Channel) {
		serverOnDisconnected <- struct{}{}
	})
	s.server.SetMessageHandler(func(ws Channel, data []byte) error {
		s.Fail("unexpected message received")
		return fmt.Errorf("unexpected message received")
	})
	go s.server.Start(serverPort, serverPath)
	time.Sleep(100 * time.Millisecond)

	// Test client
	s.client.SetDisconnectedHandler(func(err error) {
		defer func() {
			clientOnDisconnected <- struct{}{}
		}()
		s.Error(err)
		errMsg := err.Error()
		s.Contains(errMsg, "timeout")
	})
	s.client.SetReconnectedHandler(func() {
		time.Sleep(50 * time.Millisecond) // Make sure we reconnected after backoff
		reconnected <- struct{}{}
	})
	host := s.proxy.Listen
	u := url.URL{Scheme: "ws", Host: host, Path: testPath}

	// Attempt connection
	err := s.client.Start(u.String())
	s.NoError(err)
	// Slow down network. Ping will be received but pong won't go through
	_, err = s.proxy.AddToxic("writeTimeout", "timeout", "downstream", 1, toxiproxy.Attributes{
		"timeout": 5000, // 5 seconds
	})
	s.NoError(err)
	// Attempt to send message
	result := <-serverOnDisconnected
	s.NotNil(result)
	result = <-clientOnDisconnected
	s.NotNil(result)
	// Reconnect time starts
	s.proxy.RemoveToxic("writeTimeout")
	startTimeout := time.Now()
	result = <-reconnected
	s.NotNil(result)
	elapsed := time.Since(startTimeout)
	s.GreaterOrEqual(elapsed.Milliseconds(), s.client.timeoutConfig.RetryBackOffWaitMinimum.Milliseconds())
	// Cleanup
	s.client.SetDisconnectedHandler(nil)
	s.client.Stop()
	s.server.Stop()
}

func (s *NetworkTestSuite) TestServerReadTimeout() {
	// Set timeouts for test
	s.client.timeoutConfig.PongWait = 2 * time.Second
	s.client.timeoutConfig.PingPeriod = 3 * time.Second
	s.client.timeoutConfig.RetryBackOffWaitMinimum = 1 * time.Second
	s.client.timeoutConfig.RetryBackOffRandomRange = 0 // remove randomness
	s.server.timeoutConfig.PingWait = s.client.timeoutConfig.PongWait
	// Setup
	serverOnDisconnected := make(chan struct{}, 1)
	clientOnDisconnected := make(chan struct{}, 1)
	reconnected := make(chan struct{}, 1)
	s.server.SetNewClientHandler(func(ws Channel) {
		s.NotNil(ws)
	})
	s.server.SetDisconnectedClientHandler(func(ws Channel) {
		serverOnDisconnected <- struct{}{}
	})
	s.server.SetMessageHandler(func(ws Channel, data []byte) error {
		s.Fail("unexpected message received")
		return fmt.Errorf("unexpected message received")
	})
	go s.server.Start(serverPort, serverPath)
	time.Sleep(100 * time.Millisecond)

	// Test client
	s.client.SetDisconnectedHandler(func(err error) {
		defer func() {
			clientOnDisconnected <- struct{}{}
		}()
		s.Error(err)
		errMsg := err.Error()
		s.Contains(errMsg, "timeout")
	})
	s.client.SetReconnectedHandler(func() {
		time.Sleep(50 * time.Millisecond) // Make sure we reconnected after backoff
		reconnected <- struct{}{}
	})
	host := s.proxy.Listen
	u := url.URL{Scheme: "ws", Host: host, Path: testPath}

	// Attempt connection
	err := s.client.Start(u.String())
	s.NoError(err)
	// Send me
	// Slow down network. Ping will be received but pong won't go through
	_, err = s.proxy.AddToxic("writeTimeout", "timeout", "downstream", 1, toxiproxy.Attributes{
		"timeout": 5000, // 5 seconds
	})
	s.NoError(err)
	// Attempt to send message
	result := <-serverOnDisconnected
	s.NotNil(result)
	result = <-clientOnDisconnected
	s.NotNil(result)
	// Reconnect time starts
	s.proxy.RemoveToxic("writeTimeout")
	startTimeout := time.Now()
	result = <-reconnected
	s.NotNil(result)
	elapsed := time.Since(startTimeout)
	s.GreaterOrEqual(elapsed.Milliseconds(), s.client.timeoutConfig.RetryBackOffWaitMinimum.Milliseconds())
	// Cleanup
	s.client.SetDisconnectedHandler(nil)
	s.client.Stop()
	s.server.Stop()
}

//TODO: test error channel from websocket

func TestNetworkErrors(t *testing.T) {
	suite.Run(t, new(NetworkTestSuite))
}
