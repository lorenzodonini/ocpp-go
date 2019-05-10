package test

import (
	"github.com/lorenzodonini/go-ocpp/ocpp"
	"github.com/lorenzodonini/go-ocpp/ws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

// ---------------------- MOCK WEBSOCKET ----------------------
type MockWebSocket struct {
	id string
}

func (websocket MockWebSocket)GetId() string {
	return websocket.id
}

// ---------------------- MOCK WEBSOCKET SERVER ----------------------
type MockWebsocketServer struct {
	mock.Mock
	ws.WsServer
	messageHandler func(ws ws.Channel, data []byte) error
	newClientHandler func(ws ws.Channel)
}

func (websocketServer* MockWebsocketServer)Start(port int, listenPath string) {
	websocketServer.MethodCalled("Start", port, listenPath)
}

func (websocketServer* MockWebsocketServer)Stop() {
	websocketServer.MethodCalled("Stop")
}

func (websocketServer* MockWebsocketServer)Write(webSocketId string, data []byte) error {
	args := websocketServer.MethodCalled("Write", webSocketId, data)
	return args.Error(0)
}

func (websocketServer* MockWebsocketServer)SetMessageHandler(handler func(ws ws.Channel, data []byte) error) {
	websocketServer.messageHandler = handler
}

func (websocketServer* MockWebsocketServer)SetNewClientHandler(handler func(ws ws.Channel)) {
	websocketServer.newClientHandler = handler
}

func (websocketServer* MockWebsocketServer)NewClient(websocketId string, client interface{}) {
	websocketServer.MethodCalled("NewClient", websocketId, client)
}


// ---------------------- MOCK WEBSOCKET CLIENT ----------------------
type MockWebsocketClient struct {
	mock.Mock
	ws.WsClient
	messageHandler func(data []byte) error
}

func (websocketClient* MockWebsocketClient)Start(url string) {
	websocketClient.Called(url)
}

func (websocketClient* MockWebsocketClient)Stop() {
	websocketClient.MethodCalled("Stop")
}

func (websocketClient* MockWebsocketClient)SetMessageHandler(handler func(data []byte) error) {
	websocketClient.messageHandler = handler
}

func (websocketClient* MockWebsocketClient)Write(data []byte) {
	websocketClient.MethodCalled("Write", data)
}


// ---------------------- COMMON UTILITY METHODS ----------------------
func NewWebsocketServer(t *testing.T, onMessage func(data []byte) ([]byte, error)) *ws.Server {
	wsServer := ws.Server{}
	wsServer.SetMessageHandler(func(ws ws.Channel, data []byte) error {
		assert.NotNil(t, ws)
		assert.NotNil(t, data)
		if onMessage != nil {
			response, err := onMessage(data)
			assert.Nil(t, err)
			if response != nil {
				err = wsServer.Write(ws.GetId(), data)
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
	assert.Equal(t, ocpp.CALL, int(call.MessageTypeId))
	assert.Equal(t, expectedAction, call.Action)
	assert.Equal(t, expectedId, call.UniqueId)
	assert.NotNil(t, call.Payload)
}

func ParseCallResult(json string, t* testing.T) *ocpp.CallResult {
	parsedData := ocpp.ParseJsonMessage(json)
	err, result := ocpp.ParseMessage(parsedData)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	call, ok := result.(ocpp.CallResult)
	assert.Equal(t, true, ok)
	assert.NotNil(t, call)
	return &call
}

func CheckCallResult(result* ocpp.CallResult, t *testing.T, expectedId string) {
	assert.Equal(t, ocpp.CALL_RESULT, int(result.MessageTypeId))
	assert.Equal(t, expectedId, result.UniqueId)
	assert.NotNil(t, result.Payload)
}
