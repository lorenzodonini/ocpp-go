package test

import (
	"bytes"
	"fmt"
	"github.com/lorenzodonini/go-ocpp/websocket"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
	"time"
)

var (
	serverPort = 8887
	serverPath = "/ws/{id}"
	testPath = "/ws/test1"
)

func TestWebsocketEcho(t *testing.T) {
	message := []byte("Hello WebSocket!")
	wsServer := websocket.Server{}
	wsServer.SetMessageHandler(func(ws *websocket.WebSocket, data []byte) error {
		assert.NotNil(t, ws)
		assert.NotNil(t, data)
		assert.True(t, bytes.Equal(message, data))
		err := wsServer.Write(ws.Id, data)
		assert.Nil(t, err)
		return nil
	})
	go wsServer.Start(serverPort, serverPath)

	// Test message
	wsClient := websocket.Client{}
	wsClient.SetMessageHandler(func(data []byte) error {
		assert.NotNil(t, data)
		assert.True(t, bytes.Equal(message, data))
		return nil
	})
	host := fmt.Sprintf("localhost:%v", serverPort)
	u := url.URL{Scheme: "ws", Host: host, Path: "/ws/testws"}
	go func() {
		// Wait for connection to be established, then send a message
		timer := time.NewTimer(1 * time.Second)
		<-timer.C
		wsClient.Write(message)
	}()
	done := make(chan bool)
	go func() {
		timer := time.NewTimer(3 * time.Second)
		<-timer.C
		wsClient.Stop()
		done <- true
	}()
	wsClient.Start(u.String())
	result := <- done
	assert.True(t, result)
}