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

func (websocket MockWebSocket) GetId() string {
	return websocket.id
}

// ---------------------- MOCK WEBSOCKET SERVER ----------------------
type MockWebsocketServer struct {
	mock.Mock
	ws.WsServer
	messageHandler   func(ws ws.Channel, data []byte) error
	newClientHandler func(ws ws.Channel)
}

func (websocketServer *MockWebsocketServer) Start(port int, listenPath string) {
	websocketServer.MethodCalled("Start", port, listenPath)
}

func (websocketServer *MockWebsocketServer) Stop() {
	websocketServer.MethodCalled("Stop")
}

func (websocketServer *MockWebsocketServer) Write(webSocketId string, data []byte) error {
	args := websocketServer.MethodCalled("Write", webSocketId, data)
	return args.Error(0)
}

func (websocketServer *MockWebsocketServer) SetMessageHandler(handler func(ws ws.Channel, data []byte) error) {
	websocketServer.messageHandler = handler
}

func (websocketServer *MockWebsocketServer) SetNewClientHandler(handler func(ws ws.Channel)) {
	websocketServer.newClientHandler = handler
}

func (websocketServer *MockWebsocketServer) NewClient(websocketId string, client interface{}) {
	websocketServer.MethodCalled("NewClient", websocketId, client)
}

// ---------------------- MOCK WEBSOCKET CLIENT ----------------------
type MockWebsocketClient struct {
	mock.Mock
	ws.WsClient
	messageHandler func(data []byte) error
}

func (websocketClient *MockWebsocketClient) Start(url string) error {
	websocketClient.Called(url)
	return nil
}

func (websocketClient *MockWebsocketClient) Stop() {
	websocketClient.MethodCalled("Stop")
}

func (websocketClient *MockWebsocketClient) SetMessageHandler(handler func(data []byte) error) {
	websocketClient.messageHandler = handler
}

func (websocketClient *MockWebsocketClient) Write(data []byte) {
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

func ParseCall(endpoint *ocpp.Endpoint, json string, t *testing.T) *ocpp.Call {
	parsedData := ocpp.ParseJsonMessage(json)
	result, err := endpoint.ParseMessage(parsedData)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	call, ok := result.(*ocpp.Call)
	assert.Equal(t, true, ok)
	assert.NotNil(t, call)
	return call
}

func CheckCall(call *ocpp.Call, t *testing.T, expectedAction string, expectedId string) {
	assert.Equal(t, ocpp.CALL, int(call.GetMessageTypeId()))
	assert.Equal(t, expectedAction, call.Action)
	assert.Equal(t, expectedId, call.GetUniqueId())
	assert.NotNil(t, call.Payload)
	err := validate.Struct(call)
	assert.Nil(t, err)
}

func ParseCallResult(endpoint *ocpp.Endpoint, json string, t *testing.T) *ocpp.CallResult {
	parsedData := ocpp.ParseJsonMessage(json)
	result, err := endpoint.ParseMessage(parsedData)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	callResult, ok := result.(*ocpp.CallResult)
	assert.Equal(t, true, ok)
	assert.NotNil(t, callResult)
	return callResult
}

func CheckCallResult(result *ocpp.CallResult, t *testing.T, expectedId string) {
	assert.Equal(t, ocpp.CALL_RESULT, int(result.GetMessageTypeId()))
	assert.Equal(t, expectedId, result.GetUniqueId())
	assert.NotNil(t, result.Payload)
	err := validate.Struct(result)
	assert.Nil(t, err)
}

func ParseCallError(endpoint *ocpp.Endpoint, json string, t *testing.T) *ocpp.CallError {
	parsedData := ocpp.ParseJsonMessage(json)
	result, err := endpoint.ParseMessage(parsedData)
	assert.Nil(t, err)
	callError, ok := result.(*ocpp.CallError)
	assert.Equal(t, true, ok)
	assert.NotNil(t, callError)
	return callError
}

func CheckCallError(t *testing.T, callError *ocpp.CallError, expectedId string, expectedError ocpp.CallError, expectedDescription string, expectedDetails interface{}) {
	assert.Equal(t, ocpp.CALL_ERROR, int(callError.GetMessageTypeId()))
	assert.Equal(t, expectedId, callError.GetUniqueId())
	assert.Equal(t, expectedError, callError.ErrorCode)
	assert.Equal(t, expectedDescription, callError.ErrorDescription)
	assert.Equal(t, expectedDetails, callError.ErrorDetails)
	err := validate.Struct(callError)
	assert.Nil(t, err)
}

type RequestTestEntry struct {
	request       ocpp.Request
	expectedValid bool
}

type ConfirmationTestEntry struct {
	confirmation  ocpp.Confirmation
	expectedValid bool
}

func executeRequestTestTable(t *testing.T, testTable []RequestTestEntry) {
	for _, testCase := range testTable {
		err := validate.Struct(testCase.request)
		assert.Equal(t, testCase.expectedValid, err == nil)
	}
}

func executeConfirmationTestTable(t *testing.T, testTable []ConfirmationTestEntry) {
	for _, testCase := range testTable {
		err := validate.Struct(testCase.confirmation)
		assert.Equal(t, testCase.expectedValid, err == nil)
	}
}
