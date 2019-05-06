package test

import (
	"github.com/lorenzodonini/go-ocpp/ocpp"
	"github.com/lorenzodonini/go-ocpp/websocket"
	"github.com/stretchr/testify/assert"
	"testing"
)

func NewWebsocketServer(t *testing.T, onMessage func(data []byte) ([]byte, error)) *websocket.Server {
	wsServer := websocket.Server{}
	wsServer.SetMessageHandler(func(ws *websocket.WebSocket, data []byte) error {
		assert.NotNil(t, ws)
		assert.NotNil(t, data)
		if onMessage != nil {
			response, err := onMessage(data)
			assert.Nil(t, err)
			if response != nil {
				err = wsServer.Write(ws.Id, data)
				assert.Nil(t, err)
			}
		}
		return nil
	})
	return &wsServer
}

func NewWebsocketClient(t *testing.T, onMessage func(data []byte) ([]byte, error)) *websocket.Client {
	wsClient := websocket.Client{}
	wsClient.SetMessageHandler(func(data []byte) error {
		assert.NotNil(t, data)
		if onMessage != nil {
			response, err := onMessage(data)
			assert.Nil(t, err)
			if response != nil {
				wsClient.Write(data)
			}
		}
		return nil
	})
	return &wsClient
}

func ParseCall(json string, t* testing.T) *ocpp.Call {
	parsedData := ocpp.ParseJsonMessage(json)
	err, result := ocpp.ParseMessage(parsedData)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	call, ok := result.(ocpp.Call)
	assert.Equal(t, true, ok)
	assert.NotNil(t, call)
	return &call
}

func CheckCall(call* ocpp.Call, t *testing.T, expectedAction string, expectedId string) {
	assert.Equal(t, ocpp.CALL, call.MessageTypeId)
	assert.Equal(t, expectedAction, call.Action)
	assert.Equal(t, expectedId, call.UniqueId)
	assert.NotNil(t, call.Payload)
}
