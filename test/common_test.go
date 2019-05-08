package test

import (
	"github.com/lorenzodonini/go-ocpp/ocpp"
	"github.com/lorenzodonini/go-ocpp/ws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

// ---------------------- MOCK WEBSOCKET SERVER ----------------------
type MockWebSocket struct {

}

type MockWebsocketServer struct {
	mock.Mock
	ws.WsServer
}

func (websocketServer MockWebsocketServer)Start(port int, listenPath string) {
	websocketServer.MethodCalled("Start", port, listenPath)
}

func (websocketServer MockWebsocketServer)Stop() {
	websocketServer.MethodCalled("Stop")
}

func (websocketServer MockWebsocketServer)Write(webSocketId string, data []byte) error {
	args := websocketServer.MethodCalled("Write", webSocketId, data)
	return args.Error(0)
}

func (websocketServer MockWebsocketServer)SetMessageHandler(handler func(ws *ws.WebSocket, data []byte) error) {

}

func (websocketServer MockWebsocketServer)SetNewClientHandler(handler func(ws *ws.WebSocket)) {

}

func (websocketServer MockWebsocketServer)Receive(webSocketId string, data []byte) {

}


// ---------------------- MOCK WEBSOCKET CLIENT ----------------------
type MockWebsocketClient struct {
	mock.Mock
	ws.WsClient
}

func (websocketClient MockWebsocketClient)Start(url string) {
	websocketClient.MethodCalled("Start", url)
}

func (websocketClient MockWebsocketClient)Stop() {
	websocketClient.MethodCalled("Stop")
}

func (websocketClient MockWebsocketClient)SetMessageHandler(handler func(data []byte) error) {
}

func (websocketClient MockWebsocketClient)Write(data []byte) {
	websocketClient.MethodCalled("Write", data)
}


// ---------------------- COMMON UTILITY METHODS ----------------------
func NewWebsocketServer(t *testing.T, onMessage func(data []byte) ([]byte, error)) *ws.Server {
	wsServer := ws.Server{}
	wsServer.SetMessageHandler(func(ws *ws.WebSocket, data []byte) error {
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

func NewWebsocketClient(t *testing.T, onMessage func(data []byte) ([]byte, error)) *ws.Client {
	wsClient := ws.Client{}
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
