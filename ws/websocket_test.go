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
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

const (
	serverPort = 8887
	serverPath = "/ws/{id}"
	testPath   = "/ws/testws"
)

func NewWebsocketServer(t *testing.T, onMessage func(data []byte) ([]byte, error)) *Server {
	wsServer := NewServer()
	wsServer.SetMessageHandler(func(ws Channel, data []byte) error {
		assert.NotNil(t, ws)
		assert.NotNil(t, data)
		if onMessage != nil {
			response, err := onMessage(data)
			assert.Nil(t, err)
			if response != nil {
				err = wsServer.Write(ws.GetID(), data)
				assert.Nil(t, err)
			}
		}
		return nil
	})
	return wsServer
}

func NewWebsocketClient(t *testing.T, onMessage func(data []byte) ([]byte, error)) *Client {
	wsClient := NewClient()
	wsClient.SetMessageHandler(func(data []byte) error {
		assert.NotNil(t, data)
		if onMessage != nil {
			response, err := onMessage(data)
			assert.Nil(t, err)
			if response != nil {
				err = wsClient.Write(data)
				assert.Nil(t, err)
			}
		}
		return nil
	})
	return wsClient
}

func TestWebsocketSetConnected(t *testing.T) {
	wsClient := NewWebsocketClient(t, func(data []byte) ([]byte, error) {
		return nil, nil
	})
	assert.False(t, wsClient.IsConnected())
	wsClient.setConnected(true)
	assert.True(t, wsClient.IsConnected())
	wsClient.setConnected(false)
	assert.False(t, wsClient.IsConnected())
}

func TestWebsocketEcho(t *testing.T) {
	message := []byte("Hello WebSocket!")
	var wsServer *Server
	wsServer = NewWebsocketServer(t, func(data []byte) ([]byte, error) {
		assert.True(t, bytes.Equal(message, data))
		return data, nil
	})
	go wsServer.Start(serverPort, serverPath)
	time.Sleep(1 * time.Second)

	// Test message
	wsClient := NewWebsocketClient(t, func(data []byte) ([]byte, error) {
		assert.True(t, bytes.Equal(message, data))
		return nil, nil
	})
	host := fmt.Sprintf("localhost:%v", serverPort)
	u := url.URL{Scheme: "ws", Host: host, Path: testPath}
	// Wait for connection to be established, then send a message
	go func() {
		timer := time.NewTimer(1 * time.Second)
		<-timer.C
		err := wsClient.Write(message)
		assert.Nil(t, err)
	}()
	done := make(chan bool)
	// Wait for messages to be exchanged, then close connection
	go func() {
		timer := time.NewTimer(3 * time.Second)
		<-timer.C
		wsClient.Stop()
		done <- true
	}()
	err := wsClient.Start(u.String())
	assert.Nil(t, err)
	assert.True(t, wsClient.IsConnected())
	result := <-done
	assert.True(t, result)
	// Cleanup
	wsServer.Stop()
}

func TestTLSWebsocketEcho(t *testing.T) {
	message := []byte("Hello Secure WebSocket!")
	var wsServer *Server
	// Use NewTLSServer() when in different package
	wsServer = NewWebsocketServer(t, func(data []byte) ([]byte, error) {
		assert.True(t, bytes.Equal(message, data))
		return data, nil
	})
	// Create self-signed TLS certificate
	certFilename := "/tmp/cert.pem"
	keyFilename := "/tmp/key.pem"
	err := createTLSCertificate(certFilename, keyFilename, "localhost", nil, nil)
	require.Nil(t, err)
	defer os.Remove(certFilename)
	defer os.Remove(keyFilename)

	// Set self-signed TLS certificate
	wsServer.tlsCertificatePath = certFilename
	wsServer.tlsCertificateKey = keyFilename
	go wsServer.Start(serverPort, serverPath)
	time.Sleep(1 * time.Second)

	// Create TLS client
	wsClient := NewWebsocketClient(t, func(data []byte) ([]byte, error) {
		assert.True(t, bytes.Equal(message, data))
		return nil, nil
	})
	wsClient.dialOptions = append(wsClient.dialOptions, func(dialer *websocket.Dialer) {
		certPool := x509.NewCertPool()
		data, err := ioutil.ReadFile(certFilename)
		assert.Nil(t, err)
		ok := certPool.AppendCertsFromPEM(data)
		assert.True(t, ok)
		dialer.TLSClientConfig = &tls.Config{
			RootCAs: certPool,
		}
	})
	// Test message
	host := fmt.Sprintf("localhost:%v", serverPort)
	u := url.URL{Scheme: "wss", Host: host, Path: testPath}
	// Wait for connection to be established, then send a message to server
	go func() {
		timer := time.NewTimer(1 * time.Second)
		<-timer.C
		err := wsClient.Write(message)
		assert.Nil(t, err)
	}()
	done := make(chan bool)
	// Wait for messages to be exchanged, then close connection
	go func() {
		timer := time.NewTimer(3 * time.Second)
		<-timer.C
		wsClient.Stop()
		done <- true
	}()
	err = wsClient.Start(u.String())
	assert.Nil(t, err)
	result := <-done
	assert.True(t, result)
	// Cleanup
	wsServer.Stop()
}

func TestWebsocketClientConnectionBreak(t *testing.T) {
	newClient := make(chan bool)
	disconnected := make(chan bool)
	var wsServer *Server
	wsServer = NewWebsocketServer(t, nil)
	wsServer.SetNewClientHandler(func(ws Channel) {
		newClient <- true
	})
	wsServer.SetDisconnectedClientHandler(func(ws Channel) {
		disconnected <- true
	})
	go wsServer.Start(serverPort, serverPath)
	time.Sleep(1 * time.Second)

	// Test
	wsClient := NewWebsocketClient(t, nil)
	host := fmt.Sprintf("localhost:%v", serverPort)
	u := url.URL{Scheme: "ws", Host: host, Path: testPath}
	// Wait for connection to be established, then break the connection
	go func() {
		timer := time.NewTimer(1 * time.Second)
		<-timer.C
		err := wsClient.webSocket.connection.Close()
		assert.Nil(t, err)
	}()
	err := wsClient.Start(u.String())
	assert.Nil(t, err)
	result := <-newClient
	assert.True(t, result)
	result = <-disconnected
	assert.True(t, result)
	// Cleanup
	wsServer.Stop()
}

func TestWebsocketServerConnectionBreak(t *testing.T) {
	var wsServer *Server
	disconnected := make(chan bool)
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
		disconnected <- true
	})
	go wsServer.Start(serverPort, serverPath)
	time.Sleep(1 * time.Second)

	// Test
	wsClient := NewWebsocketClient(t, nil)
	host := fmt.Sprintf("localhost:%v", serverPort)
	u := url.URL{Scheme: "ws", Host: host, Path: testPath}
	err := wsClient.Start(u.String())
	assert.Nil(t, err)
	result := <-disconnected
	assert.True(t, result)
	// Cleanup
	wsServer.Stop()
}

func TestValidBasicAuth(t *testing.T) {
	authUsername := "testUsername"
	authPassword := "testPassword"
	var wsServer *Server
	// Create self-signed TLS certificate
	certFilename := "/tmp/cert.pem"
	keyFilename := "/tmp/key.pem"
	err := createTLSCertificate(certFilename, keyFilename, "localhost", nil, nil)
	require.Nil(t, err)
	defer os.Remove(certFilename)
	defer os.Remove(keyFilename)

	// Create TLS server with self-signed certificate
	wsServer = NewTLSServer(certFilename, keyFilename, nil)
	// Add basic auth handler
	wsServer.SetBasicAuthHandler(func(username string, password string) bool {
		require.Equal(t, authUsername, username)
		require.Equal(t, authPassword, password)
		return true
	})
	connected := make(chan bool)
	wsServer.SetNewClientHandler(func(ws Channel) {
		connected <- true
	})
	// Run server
	go wsServer.Start(serverPort, serverPath)
	time.Sleep(1 * time.Second)

	// Create TLS client
	certPool := x509.NewCertPool()
	data, err := ioutil.ReadFile(certFilename)
	require.Nil(t, err)
	ok := certPool.AppendCertsFromPEM(data)
	require.True(t, ok)
	wsClient := NewTLSClient(&tls.Config{
		RootCAs: certPool,
	})
	// Add basic auth
	wsClient.SetBasicAuth(authUsername, authPassword)
	// Test connection
	host := fmt.Sprintf("localhost:%v", serverPort)
	u := url.URL{Scheme: "wss", Host: host, Path: testPath}
	err = wsClient.Start(u.String())
	require.Nil(t, err)
	result := <-connected
	assert.True(t, result)
	// Cleanup
	wsClient.Stop()
	wsServer.Stop()
}

func TestInvalidBasicAuth(t *testing.T) {
	authUsername := "testUsername"
	authPassword := "testPassword"
	var wsServer *Server
	// Create self-signed TLS certificate
	certFilename := "/tmp/cert.pem"
	keyFilename := "/tmp/key.pem"
	err := createTLSCertificate(certFilename, keyFilename, "localhost", nil, nil)
	require.Nil(t, err)
	defer os.Remove(certFilename)
	defer os.Remove(keyFilename)

	// Create TLS server with self-signed certificate
	wsServer = NewTLSServer(certFilename, keyFilename, nil)
	// Add basic auth handler
	wsServer.SetBasicAuthHandler(func(username string, password string) bool {
		validCredentials := authUsername == username && authPassword == password
		require.False(t, validCredentials)
		return validCredentials
	})
	wsServer.SetNewClientHandler(func(ws Channel) {
		// Should never reach this
		t.Fail()
	})
	// Run server
	go wsServer.Start(serverPort, serverPath)
	time.Sleep(1 * time.Second)

	// Create TLS client
	certPool := x509.NewCertPool()
	data, err := ioutil.ReadFile(certFilename)
	require.Nil(t, err)
	ok := certPool.AppendCertsFromPEM(data)
	require.True(t, ok)
	wsClient := NewTLSClient(&tls.Config{
		RootCAs: certPool,
	})
	// Add basic auth
	wsClient.SetBasicAuth(authUsername, "invalidPassword")
	// Test connection
	host := fmt.Sprintf("localhost:%v", serverPort)
	u := url.URL{Scheme: "wss", Host: host, Path: testPath}
	err = wsClient.Start(u.String())
	assert.NotNil(t, err)
	httpError, ok := err.(HttpConnectionError)
	require.True(t, ok)
	require.NotNil(t, httpError)
	assert.Equal(t, http.StatusUnauthorized, httpError.HttpCode)
	// Cleanup
	wsServer.Stop()
}

func TestInvalidOriginHeader(t *testing.T) {
	var wsServer *Server
	wsServer = NewWebsocketServer(t, func(data []byte) ([]byte, error) {
		assert.Fail(t, "no message should be received from client!")
		return nil, nil
	})
	wsServer.SetNewClientHandler(func(ws Channel) {
		assert.Fail(t, "no new connection should be received from client!")
	})
	go wsServer.Start(serverPort, serverPath)
	time.Sleep(500 * time.Millisecond)

	// Test message
	wsClient := NewWebsocketClient(t, func(data []byte) ([]byte, error) {
		assert.Fail(t, "no message should be received from server!")
		return nil, nil
	})
	// Set invalid origin header
	wsClient.SetHeaderValue("Origin", "example.org")
	host := fmt.Sprintf("localhost:%v", serverPort)
	u := url.URL{Scheme: "ws", Host: host, Path: testPath}
	// Attempt to connect and expect cross-origin error
	err := wsClient.Start(u.String())
	require.Error(t, err)
	httpErr, ok := err.(HttpConnectionError)
	require.True(t, ok)
	assert.Equal(t, http.StatusForbidden, httpErr.HttpCode)
	assert.Equal(t, http.StatusForbidden, httpErr.HttpCode)
	assert.Equal(t, "websocket: bad handshake", httpErr.Message)
	// Cleanup
	wsServer.Stop()
}

func TestCustomOriginHeaderHandler(t *testing.T) {
	var wsServer *Server
	origin := "example.org"
	connected := make(chan bool)
	wsServer = NewWebsocketServer(t, func(data []byte) ([]byte, error) {
		assert.Fail(t, "no message should be received from client!")
		return nil, nil
	})
	wsServer.SetNewClientHandler(func(ws Channel) {
		connected <- true
	})
	wsServer.SetCheckOriginHandler(func(r *http.Request) bool {
		return r.Header.Get("Origin") == origin
	})
	go wsServer.Start(serverPort, serverPath)
	time.Sleep(500 * time.Millisecond)

	// Test message
	wsClient := NewWebsocketClient(t, func(data []byte) ([]byte, error) {
		assert.Fail(t, "no message should be received from server!")
		return nil, nil
	})
	// Set invalid origin header (not example.org)
	wsClient.SetHeaderValue("Origin", "localhost")
	host := fmt.Sprintf("localhost:%v", serverPort)
	u := url.URL{Scheme: "ws", Host: host, Path: testPath}
	// Attempt to connect and expect cross-origin error
	err := wsClient.Start(u.String())
	require.Error(t, err)
	httpErr, ok := err.(HttpConnectionError)
	require.True(t, ok)
	assert.Equal(t, http.StatusForbidden, httpErr.HttpCode)
	assert.Equal(t, http.StatusForbidden, httpErr.HttpCode)
	assert.Equal(t, "websocket: bad handshake", httpErr.Message)

	// Re-attempt with correct header
	wsClient.SetHeaderValue("Origin", "example.org")
	err = wsClient.Start(u.String())
	require.NoError(t, err)
	result := <-connected
	assert.True(t, result)
	// Cleanup
	wsServer.Stop()
}

func TestValidClientTLSCertificate(t *testing.T) {
	var wsServer *Server
	// Create self-signed TLS certificate
	clientCertFilename := "/tmp/client.pem"
	clientKeyFilename := "/tmp/client_key.pem"
	err := createTLSCertificate(clientCertFilename, clientKeyFilename, "localhost", nil, nil)
	defer os.Remove(clientCertFilename)
	defer os.Remove(clientKeyFilename)
	require.Nil(t, err)
	serverCertFilename := "/tmp/cert.pem"
	serverKeyFilename := "/tmp/key.pem"
	err = createTLSCertificate(serverCertFilename, serverKeyFilename, "localhost", nil, nil)
	require.Nil(t, err)
	defer os.Remove(serverCertFilename)
	defer os.Remove(serverKeyFilename)

	// Create TLS server with self-signed certificate
	certPool := x509.NewCertPool()
	data, err := ioutil.ReadFile(clientCertFilename)
	require.Nil(t, err)
	ok := certPool.AppendCertsFromPEM(data)
	require.True(t, ok)
	wsServer = NewTLSServer(serverCertFilename, serverKeyFilename, &tls.Config{
		ClientCAs:  certPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	})
	// Add basic auth handler
	connected := make(chan bool)
	wsServer.SetNewClientHandler(func(ws Channel) {
		connected <- true
	})
	// Run server
	go wsServer.Start(serverPort, serverPath)
	time.Sleep(1 * time.Second)

	// Create TLS client
	certPool = x509.NewCertPool()
	data, err = ioutil.ReadFile(serverCertFilename)
	require.Nil(t, err)
	ok = certPool.AppendCertsFromPEM(data)
	require.True(t, ok)
	loadedCert, err := tls.LoadX509KeyPair(clientCertFilename, clientKeyFilename)
	require.Nil(t, err)
	wsClient := NewTLSClient(&tls.Config{
		RootCAs:      certPool,
		Certificates: []tls.Certificate{loadedCert},
	})
	// Test connection
	host := fmt.Sprintf("localhost:%v", serverPort)
	u := url.URL{Scheme: "wss", Host: host, Path: testPath}
	err = wsClient.Start(u.String())
	assert.Nil(t, err)
	result := <-connected
	assert.True(t, result)
	// Cleanup
	wsServer.Stop()
}

func TestInvalidClientTLSCertificate(t *testing.T) {
	var wsServer *Server
	// Create self-signed TLS certificate
	clientCertFilename := "/tmp/client.pem"
	clientKeyFilename := "/tmp/client_key.pem"
	err := createTLSCertificate(clientCertFilename, clientKeyFilename, "localhost", nil, nil)
	defer os.Remove(clientCertFilename)
	defer os.Remove(clientKeyFilename)
	require.Nil(t, err)
	serverCertFilename := "/tmp/cert.pem"
	serverKeyFilename := "/tmp/key.pem"
	err = createTLSCertificate(serverCertFilename, serverKeyFilename, "localhost", nil, nil)
	require.Nil(t, err)
	defer os.Remove(serverCertFilename)
	defer os.Remove(serverKeyFilename)

	// Create TLS server with self-signed certificate
	certPool := x509.NewCertPool()
	data, err := ioutil.ReadFile(serverCertFilename)
	require.Nil(t, err)
	ok := certPool.AppendCertsFromPEM(data)
	require.True(t, ok)
	wsServer = NewTLSServer(serverCertFilename, serverKeyFilename, &tls.Config{
		ClientCAs:  certPool,                       // Contains server certificate as allowed client CA
		ClientAuth: tls.RequireAndVerifyClientCert, // Requires client certificate signed by allowed CA (server)
	})
	// Add basic auth handler
	wsServer.SetNewClientHandler(func(ws Channel) {
		// Should never reach this
		t.Fail()
	})
	// Run server
	go wsServer.Start(serverPort, serverPath)
	time.Sleep(1 * time.Second)

	// Create TLS client
	certPool = x509.NewCertPool()
	data, err = ioutil.ReadFile(serverCertFilename)
	require.Nil(t, err)
	ok = certPool.AppendCertsFromPEM(data)
	require.True(t, ok)
	loadedCert, err := tls.LoadX509KeyPair(clientCertFilename, clientKeyFilename)
	require.Nil(t, err)
	wsClient := NewTLSClient(&tls.Config{
		RootCAs:      certPool,                      // Contains server certificate as allowed server CA
		Certificates: []tls.Certificate{loadedCert}, // Contains self-signed client certificate. Will be rejected by server
	})
	// Test connection
	host := fmt.Sprintf("localhost:%v", serverPort)
	u := url.URL{Scheme: "wss", Host: host, Path: testPath}
	err = wsClient.Start(u.String())
	assert.NotNil(t, err)
	netError, ok := err.(net.Error)
	require.True(t, ok)
	assert.Equal(t, "remote error: tls: bad certificate", netError.Error()) // tls.alertBadCertificate = 42
	// Cleanup
	wsServer.Stop()
}

func TestUnsupportedSubprotocol(t *testing.T) {
	var wsServer *Server
	disconnected := make(chan bool)
	wsServer = NewWebsocketServer(t, nil)
	wsServer.SetNewClientHandler(func(ws Channel) {
		assert.Fail(t, "invalid subprotocol expected, but hit client handler instead")
		t.Fail()
	})
	wsServer.SetDisconnectedClientHandler(func(ws Channel) {
		disconnected <- true
	})
	wsServer.AddSupportedSubprotocol(defaultSubProtocol)
	go wsServer.Start(serverPort, serverPath)
	time.Sleep(1 * time.Second)

	wsClient := NewWebsocketClient(t, nil)
	// Set invalid subprotocol
	wsClient.dialOptions = append(wsClient.dialOptions, func(dialer *websocket.Dialer) {
		dialer.Subprotocols = []string{"unsupportedSubProto"}
	})
	// Test
	host := fmt.Sprintf("localhost:%v", serverPort)
	u := url.URL{Scheme: "ws", Host: host, Path: testPath}
	err := wsClient.Start(u.String())
	assert.NotNil(t, err)
	// Cleanup
	wsServer.Stop()
}

func TestSetServerTimeoutConfig(t *testing.T) {
	var wsServer *Server
	disconnected := make(chan bool)
	wsServer = NewWebsocketServer(t, nil)
	wsServer.SetNewClientHandler(func(ws Channel) {
	})
	wsServer.SetDisconnectedClientHandler(func(ws Channel) {
		// TODO: check for error with upcoming API
		disconnected <- true
	})
	// Setting server timeout
	config := NewServerTimeoutConfig()
	pingWait := 2 * time.Second
	writeWait := 2 * time.Second
	config.PingWait = pingWait
	config.WriteWait = writeWait
	wsServer.SetTimeoutConfig(config)
	// Start server
	go wsServer.Start(serverPort, serverPath)
	time.Sleep(500 * time.Millisecond)
	assert.Equal(t, wsServer.timeoutConfig.PingWait, pingWait)
	assert.Equal(t, wsServer.timeoutConfig.WriteWait, writeWait)
	// Run test
	wsClient := NewWebsocketClient(t, nil)
	host := fmt.Sprintf("localhost:%v", serverPort)
	u := url.URL{Scheme: "ws", Host: host, Path: testPath}
	err := wsClient.Start(u.String())
	assert.NoError(t, err)
	result := <-disconnected
	assert.True(t, result)
	// Cleanup
	wsClient.Stop()
	wsServer.Stop()
}

func TestSetClientTimeoutConfig(t *testing.T) {
	var wsServer *Server
	disconnected := make(chan bool)
	wsServer = NewWebsocketServer(t, nil)
	wsServer.SetNewClientHandler(func(ws Channel) {
	})
	wsServer.SetDisconnectedClientHandler(func(ws Channel) {
		// TODO: check for error with upcoming API
		disconnected <- true
	})
	// Start server
	go wsServer.Start(serverPort, serverPath)
	time.Sleep(500 * time.Millisecond)
	// Run test
	wsClient := NewWebsocketClient(t, nil)
	host := fmt.Sprintf("localhost:%v", serverPort)
	u := url.URL{Scheme: "ws", Host: host, Path: testPath}
	// Set client timeout
	config := NewClientTimeoutConfig()
	handshakeTimeout := 1 * time.Nanosecond // Very low timeout, handshake will fail
	pongWait := 2 * time.Second
	writeWait := 2 * time.Second
	pingPeriod := 5 * time.Second
	config.PongWait = pongWait
	config.HandshakeTimeout = handshakeTimeout
	config.WriteWait = writeWait
	config.PingPeriod = pingPeriod
	wsClient.SetTimeoutConfig(config)
	// Start client and expect handshake error
	err := wsClient.Start(u.String())
	opError, ok := err.(*net.OpError)
	require.True(t, ok)
	assert.Equal(t, "dial", opError.Op)
	assert.True(t, opError.Timeout())
	assert.Error(t, opError.Err, "i/o timeout")
	config.HandshakeTimeout = defaultHandshakeTimeout
	wsClient.SetTimeoutConfig(config)
	// Start client
	err = wsClient.Start(u.String())
	require.NoError(t, err)
	assert.Equal(t, wsClient.timeoutConfig.PongWait, pongWait)
	assert.Equal(t, wsClient.timeoutConfig.WriteWait, writeWait)
	assert.Equal(t, wsClient.timeoutConfig.PingPeriod, pingPeriod)
	result := <-disconnected
	assert.True(t, result)
	// Cleanup
	wsClient.Stop()
	wsServer.Stop()
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
