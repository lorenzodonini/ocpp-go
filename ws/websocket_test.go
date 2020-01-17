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
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

const (
	serverPort = 8887
	serverPath = "/ws/{id}"
	testPath   = "/ws/testws"
)

func NewWebsocketServer(t *testing.T, onMessage func(data []byte) ([]byte, error)) *Server {
	wsServer := Server{}
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
	return &wsServer
}

func NewWebsocketClient(t *testing.T, onMessage func(data []byte) ([]byte, error)) *Client {
	wsClient := Client{}
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
	return &wsClient
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
	err := createTLSCertificate(certFilename, keyFilename)
	assert.Nil(t, err)
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

// Utility function
func createTLSCertificate(certificateFilename string, keyFilename string) error {
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
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"ocpp-go"},
			CommonName:   "localhost",
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
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
