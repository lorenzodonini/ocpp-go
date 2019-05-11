package test

import (
	"bytes"
	"fmt"
	"github.com/lorenzodonini/go-ocpp/ws"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
	"time"
)

var (
	serverPort = 8887
	serverPath = "/ws/{id}"
	testPath = "/ws/testws"
)

func TestWebsocketEcho(t *testing.T) {
	message := []byte("Hello WebSocket!")
	var wsServer *ws.Server
	wsServer = NewWebsocketServer(t, func(data []byte) ([]byte, error) {
		assert.True(t, bytes.Equal(message, data))
		return data, nil
	})
	go wsServer.Start(serverPort, serverPath)

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
		wsClient.Write(message)
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
	result := <- done
	assert.True(t, result)
}