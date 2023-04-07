package ocppj_test

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"testing"

	ut "github.com/go-playground/universal-translator"

	"github.com/lorenzodonini/ocpp-go/logging"

	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocppj"
	"github.com/lorenzodonini/ocpp-go/ws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gopkg.in/go-playground/validator.v9"
)

// ---------------------- MOCK WEBSOCKET ----------------------
type MockWebSocket struct {
	id string
}

func (websocket MockWebSocket) ID() string {
	return websocket.id

}

func (websocket MockWebSocket) RemoteAddr() net.Addr {
	return &net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 80,
	}
}

func (websocket MockWebSocket) TLSConnectionState() *tls.ConnectionState {
	return nil
}

func NewMockWebSocket(id string) MockWebSocket {
	return MockWebSocket{id: id}
}

// ---------------------- MOCK WEBSOCKET SERVER ----------------------

type MockWebsocketServer struct {
	mock.Mock
	ws.WsServer
	MessageHandler            func(ws ws.Channel, data []byte) error
	NewClientHandler          func(ws ws.Channel)
	CheckClientHandler        ws.CheckClientHandler
	DisconnectedClientHandler func(ws ws.Channel)
	errC                      chan error
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
	websocketServer.MessageHandler = handler
}

func (websocketServer *MockWebsocketServer) SetNewClientHandler(handler func(ws ws.Channel)) {
	websocketServer.NewClientHandler = handler
}

func (websocketServer *MockWebsocketServer) SetDisconnectedClientHandler(handler func(ws ws.Channel)) {
	websocketServer.DisconnectedClientHandler = handler
}

func (websocketServer *MockWebsocketServer) AddSupportedSubprotocol(subProto string) {
}

func (websocketServer *MockWebsocketServer) Errors() <-chan error {
	if websocketServer.errC == nil {
		websocketServer.errC = make(chan error, 1)
	}
	return websocketServer.errC
}

func (websocketServer *MockWebsocketServer) ThrowError(err error) {
	if websocketServer.errC != nil {
		websocketServer.errC <- err
	}
}

func (websocketServer *MockWebsocketServer) NewClient(websocketId string, client interface{}) {
	websocketServer.MethodCalled("NewClient", websocketId, client)
}

func (websocketServer *MockWebsocketServer) SetCheckClientHandler(handler func(id string, r *http.Request) bool) {
	websocketServer.CheckClientHandler = handler
}

// ---------------------- MOCK WEBSOCKET CLIENT ----------------------

type MockWebsocketClient struct {
	mock.Mock
	ws.WsClient
	MessageHandler      func(data []byte) error
	ReconnectedHandler  func()
	DisconnectedHandler func(err error)
	errC                chan error
}

func (websocketClient *MockWebsocketClient) Start(url string) error {
	args := websocketClient.MethodCalled("Start", url)
	return args.Error(0)
}

func (websocketClient *MockWebsocketClient) Stop() {
	websocketClient.MethodCalled("Stop")
}

func (websocketClient *MockWebsocketClient) SetMessageHandler(handler func(data []byte) error) {
	websocketClient.MessageHandler = handler
}

func (websocketClient *MockWebsocketClient) SetReconnectedHandler(handler func()) {
	websocketClient.ReconnectedHandler = handler
}

func (websocketClient *MockWebsocketClient) SetDisconnectedHandler(handler func(err error)) {
	websocketClient.DisconnectedHandler = handler
}

func (websocketClient *MockWebsocketClient) ThrowError(err error) {
	if websocketClient.errC != nil {
		websocketClient.errC <- err
	}
}

func (websocketClient *MockWebsocketClient) Write(data []byte) error {
	args := websocketClient.MethodCalled("Write", data)
	return args.Error(0)
}

func (websocketClient *MockWebsocketClient) AddOption(option interface{}) {
}

func (websocketClient *MockWebsocketClient) SetRequestedSubProtocol(subProto string) {
}

func (websocketClient *MockWebsocketClient) SetBasicAuth(username string, password string) {
}

func (websocketClient *MockWebsocketClient) SetTimeoutConfig(config ws.ClientTimeoutConfig) {
}

func (websocketClient *MockWebsocketClient) Errors() <-chan error {
	if websocketClient.errC == nil {
		websocketClient.errC = make(chan error, 1)
	}
	return websocketClient.errC
}

func (websocketClient *MockWebsocketClient) IsConnected() bool {
	args := websocketClient.MethodCalled("IsConnected")
	return args.Bool(0)
}

// ---------------------- MOCK FEATURE ----------------------
const (
	MockFeatureName = "Mock"
)

type MockRequest struct {
	mock.Mock
	MockValue string      `json:"mockValue" validate:"required,max=10"`
	MockAny   interface{} `json:"mockAny"`
}

type MockConfirmation struct {
	mock.Mock
	MockValue string `json:"mockValue" validate:"required,min=5"`
}

type MockFeature struct {
	mock.Mock
}

func (f MockFeature) GetFeatureName() string {
	return MockFeatureName
}

func (f MockFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(MockRequest{})
}

func (f MockFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(MockConfirmation{})
}

func (r MockRequest) GetFeatureName() string {
	return MockFeatureName
}

func (c MockConfirmation) GetFeatureName() string {
	return MockFeatureName
}

func newMockRequest(value string) *MockRequest {
	return &MockRequest{MockValue: value}
}

func newMockConfirmation(value string) *MockConfirmation {
	return &MockConfirmation{MockValue: value}
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
				err = wsServer.Write(ws.ID(), data)
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
				err = wsClient.Write(data)
				assert.Nil(t, err)
			}
		}
		return nil
	})
	return &wsClient
}

func ParseCall(endpoint *ocppj.Endpoint, state ocppj.ClientState, json string, t *testing.T) *ocppj.Call {
	parsedData, err := ocppj.ParseJsonMessage(json)
	require.NoError(t, err)
	require.NotNil(t, parsedData)
	result, err := endpoint.ParseMessage(parsedData, state)
	require.NoError(t, err)
	require.NotNil(t, result)
	call, ok := result.(*ocppj.Call)
	assert.Equal(t, true, ok)
	assert.NotNil(t, call)
	return call
}

func CheckCall(call *ocppj.Call, t *testing.T, expectedAction string, expectedId string) {
	assert.Equal(t, ocppj.CALL, call.GetMessageTypeId())
	assert.Equal(t, expectedAction, call.Action)
	assert.Equal(t, expectedId, call.GetUniqueId())
	assert.NotNil(t, call.Payload)
	err := Validate.Struct(call)
	assert.Nil(t, err)
}

func ParseCallResult(endpoint *ocppj.Endpoint, state ocppj.ClientState, json string, t *testing.T) *ocppj.CallResult {
	parsedData, err := ocppj.ParseJsonMessage(json)
	require.NoError(t, err)
	require.NotNil(t, parsedData)
	result, ocppErr := endpoint.ParseMessage(parsedData, state)
	require.NoError(t, ocppErr)
	require.NotNil(t, result)
	callResult, ok := result.(*ocppj.CallResult)
	assert.Equal(t, true, ok)
	assert.NotNil(t, callResult)
	return callResult
}

func CheckCallResult(result *ocppj.CallResult, t *testing.T, expectedId string) {
	assert.Equal(t, ocppj.CALL_RESULT, result.GetMessageTypeId())
	assert.Equal(t, expectedId, result.GetUniqueId())
	assert.NotNil(t, result.Payload)
	err := Validate.Struct(result)
	assert.Nil(t, err)
}

func ParseCallError(endpoint *ocppj.Endpoint, state ocppj.ClientState, json string, t *testing.T) *ocppj.CallError {
	parsedData, err := ocppj.ParseJsonMessage(json)
	require.NoError(t, err)
	require.NotNil(t, parsedData)
	result, ocppErr := endpoint.ParseMessage(parsedData, state)
	require.NoError(t, ocppErr)
	require.NotNil(t, result)
	callError, ok := result.(*ocppj.CallError)
	assert.Equal(t, true, ok)
	assert.NotNil(t, callError)
	return callError
}

func CheckCallError(t *testing.T, callError *ocppj.CallError, expectedId string, expectedError ocpp.ErrorCode, expectedDescription string, expectedDetails interface{}) {
	assert.Equal(t, ocppj.CALL_ERROR, callError.GetMessageTypeId())
	assert.Equal(t, expectedId, callError.GetUniqueId())
	assert.Equal(t, expectedError, callError.ErrorCode)
	assert.Equal(t, expectedDescription, callError.ErrorDescription)
	assert.Equal(t, expectedDetails, callError.ErrorDetails)
	err := Validate.Struct(callError)
	assert.Nil(t, err)
}

func assertPanic(t *testing.T, f func(), recoveredAssertion func(interface{})) {
	defer func() {
		r := recover()
		require.NotNil(t, r)
		recoveredAssertion(r)
	}()
	f()
}

var Validate = validator.New()

func init() {
	_ = Validate.RegisterValidation("errorCode", ocppj.IsErrorCodeValid)
}

// ---------------------- TESTS ----------------------

type OcppJTestSuite struct {
	suite.Suite
	chargePoint        *ocppj.Client
	centralSystem      *ocppj.Server
	mockServer         *MockWebsocketServer
	mockClient         *MockWebsocketClient
	clientDispatcher   ocppj.ClientDispatcher
	serverDispatcher   ocppj.ServerDispatcher
	clientRequestQueue ocppj.RequestQueue
	serverRequestMap   ocppj.ServerQueueMap
}

func (suite *OcppJTestSuite) SetupTest() {
	mockProfile := ocpp.NewProfile("mock", MockFeature{})
	mockClient := MockWebsocketClient{}
	mockServer := MockWebsocketServer{}
	suite.mockClient = &mockClient
	suite.mockServer = &mockServer
	suite.clientRequestQueue = ocppj.NewFIFOClientQueue(queueCapacity)
	suite.clientDispatcher = ocppj.NewDefaultClientDispatcher(suite.clientRequestQueue)
	suite.chargePoint = ocppj.NewClient("mock_id", suite.mockClient, suite.clientDispatcher, nil, mockProfile)
	suite.serverRequestMap = ocppj.NewFIFOQueueMap(queueCapacity)
	suite.serverDispatcher = ocppj.NewDefaultServerDispatcher(suite.serverRequestMap)
	suite.centralSystem = ocppj.NewServer(suite.mockServer, suite.serverDispatcher, nil, mockProfile)
}

func (suite *OcppJTestSuite) TearDownTest() {
	if suite.clientDispatcher.IsRunning() {
		suite.clientDispatcher.Stop()
	}
	if suite.serverDispatcher.IsRunning() {
		suite.serverDispatcher.Stop()
	}
}

func (suite *OcppJTestSuite) TestGetProfile() {
	t := suite.T()
	profile, ok := suite.chargePoint.GetProfile("mock")
	assert.True(t, ok)
	assert.NotNil(t, profile)
	feature := profile.GetFeature(MockFeatureName)
	assert.NotNil(t, feature)
	assert.Equal(t, reflect.TypeOf(MockRequest{}), feature.GetRequestType())
	assert.Equal(t, reflect.TypeOf(MockConfirmation{}), feature.GetResponseType())
}

func (suite *OcppJTestSuite) TestGetProfileForFeature() {
	t := suite.T()
	profile, ok := suite.chargePoint.GetProfileForFeature(MockFeatureName)
	assert.True(t, ok)
	assert.NotNil(t, profile)
	assert.Equal(t, "mock", profile.Name)
}

//func (suite *OcppJTestSuite) TestAddFeature

func (suite *OcppJTestSuite) TestGetProfileForInvalidFeature() {
	t := suite.T()
	profile, ok := suite.chargePoint.GetProfileForFeature("test")
	assert.False(t, ok)
	assert.Nil(t, profile)
}

func (suite *OcppJTestSuite) TestCallMaxValidation() {
	t := suite.T()
	mockLongValue := "somelongvalue"
	request := newMockRequest(mockLongValue)
	// Test invalid call
	call, err := suite.chargePoint.CreateCall(request)
	assert.Nil(t, call)
	assert.NotNil(t, err)
	assert.IsType(t, validator.ValidationErrors{}, err)
	errors := err.(validator.ValidationErrors)
	assert.Equal(t, 1, len(errors))
	validationError := errors[0]
	assert.Equal(t, "max", validationError.Tag())
}

func (suite *OcppJTestSuite) TestCallRequiredValidation() {
	t := suite.T()
	mockLongValue := ""
	request := newMockRequest(mockLongValue)
	// Test invalid call
	call, err := suite.chargePoint.CreateCall(request)
	assert.Nil(t, call)
	assert.NotNil(t, err)
	assert.IsType(t, validator.ValidationErrors{}, err)
	errors := err.(validator.ValidationErrors)
	assert.Equal(t, 1, len(errors))
	validationError := errors[0]
	assert.Equal(t, "required", validationError.Tag())
}

func (suite *OcppJTestSuite) TestCallResultMinValidation() {
	t := suite.T()
	mockShortValue := "val"
	mockUniqueId := "123456"
	confirmation := newMockConfirmation(mockShortValue)
	// Test invalid call
	callResult, err := suite.chargePoint.CreateCallResult(confirmation, mockUniqueId)
	assert.Nil(t, callResult)
	assert.NotNil(t, err)
	assert.IsType(t, validator.ValidationErrors{}, err)
	errors := err.(validator.ValidationErrors)
	assert.Equal(t, 1, len(errors))
	validationError := errors[0]
	assert.Equal(t, "min", validationError.Tag())
}

func (suite *OcppJTestSuite) TestCallResultRequiredValidation() {
	t := suite.T()
	mockShortValue := ""
	mockUniqueId := "123456"
	confirmation := newMockConfirmation(mockShortValue)
	// Test invalid call
	callResult, err := suite.chargePoint.CreateCallResult(confirmation, mockUniqueId)
	assert.Nil(t, callResult)
	assert.NotNil(t, err)
	assert.IsType(t, validator.ValidationErrors{}, err)
	errors := err.(validator.ValidationErrors)
	assert.Equal(t, 1, len(errors))
	validationError := errors[0]
	assert.Equal(t, "required", validationError.Tag())
}

func (suite *OcppJTestSuite) TestCreateCall() {
	t := suite.T()
	mockValue := "somevalue"
	request := newMockRequest(mockValue)
	call, err := suite.chargePoint.CreateCall(request)
	assert.Nil(t, err)
	CheckCall(call, t, MockFeatureName, call.UniqueId)
	message, ok := call.Payload.(*MockRequest)
	assert.True(t, ok)
	assert.NotNil(t, message)
	assert.Equal(t, mockValue, message.MockValue)
	// Check that request was not yet stored as pending request
	pendingRequest, exists := suite.chargePoint.RequestState.GetPendingRequest(call.UniqueId)
	assert.False(t, exists)
	assert.Nil(t, pendingRequest)
}

func (suite *OcppJTestSuite) TestCreateCallResult() {
	t := suite.T()
	mockValue := "someothervalue"
	mockUniqueId := "123456"
	confirmation := newMockConfirmation(mockValue)
	callResult, err := suite.chargePoint.CreateCallResult(confirmation, mockUniqueId)
	assert.Nil(t, err)
	CheckCallResult(callResult, t, mockUniqueId)
	message, ok := callResult.Payload.(*MockConfirmation)
	assert.True(t, ok)
	assert.NotNil(t, message)
	assert.Equal(t, mockValue, message.MockValue)
}

func (suite *OcppJTestSuite) TestCreateCallError() {
	t := suite.T()
	mockUniqueId := "123456"
	mockDescription := "somedescription"
	mockDetailString := "somedetailstring"
	type MockDetails struct {
		DetailString string
	}
	mockDetails := MockDetails{DetailString: mockDetailString}
	callError, err := suite.chargePoint.CreateCallError(mockUniqueId, ocppj.GenericError, mockDescription, mockDetails)
	assert.Nil(t, err)
	assert.NotNil(t, callError)
	CheckCallError(t, callError, mockUniqueId, ocppj.GenericError, mockDescription, mockDetails)
}

func (suite *OcppJTestSuite) TestParseMessageInvalidLength() {
	t := suite.T()
	mockMessage := make([]interface{}, 2)
	messageId := "12345"
	// Test invalid message length
	mockMessage[0] = ocppj.CALL // Message Type ID
	mockMessage[1] = messageId  // Unique ID
	message, err := suite.chargePoint.ParseMessage(mockMessage, suite.chargePoint.RequestState)
	require.Nil(t, message)
	require.Error(t, err)
	protoErr := err.(*ocpp.Error)
	require.NotNil(t, protoErr)
	assert.Equal(t, "", protoErr.MessageId)
	assert.Equal(t, ocppj.FormationViolation, protoErr.Code)
	assert.Equal(t, "Invalid message. Expected array length >= 3", protoErr.Description)
}

func (suite *OcppJTestSuite) TestParseMessageInvalidTypeId() {
	t := suite.T()
	mockMessage := make([]interface{}, 3)
	invalidTypeId := "2"
	messageId := "12345"
	// Test invalid message length
	mockMessage[0] = invalidTypeId // Message Type ID
	mockMessage[1] = messageId     // Unique ID
	message, err := suite.chargePoint.ParseMessage(mockMessage, suite.chargePoint.RequestState)
	require.Nil(t, message)
	require.Error(t, err)
	protoErr := err.(*ocpp.Error)
	require.NotNil(t, protoErr)
	assert.Equal(t, "", protoErr.MessageId)
	assert.Equal(t, ocppj.FormationViolation, protoErr.Code)
	assert.Equal(t, fmt.Sprintf("Invalid element %v at 0, expected message type (int)", invalidTypeId), protoErr.Description)
}

func (suite *OcppJTestSuite) TestParseMessageInvalidMessageId() {
	t := suite.T()
	mockMessage := make([]interface{}, 3)
	invalidMessageId := 12345
	// Test invalid message length
	mockMessage[0] = float64(ocppj.CALL) // Message Type ID
	mockMessage[1] = invalidMessageId    // Unique ID
	message, err := suite.chargePoint.ParseMessage(mockMessage, suite.chargePoint.RequestState)
	require.Nil(t, message)
	require.Error(t, err)
	protoErr := err.(*ocpp.Error)
	require.NotNil(t, protoErr)
	assert.Equal(t, "", protoErr.MessageId)
	assert.Equal(t, ocppj.FormationViolation, protoErr.Code)
	assert.Equal(t, fmt.Sprintf("Invalid element %v at 1, expected unique ID (string)", invalidMessageId), protoErr.Description)
}

func (suite *OcppJTestSuite) TestParseMessageUnknownTypeId() {
	t := suite.T()
	mockMessage := make([]interface{}, 3)
	messageId := "12345"
	invalidTypeId := 1
	// Test invalid message length
	mockMessage[0] = float64(invalidTypeId) // Message Type ID
	mockMessage[1] = messageId              // Unique ID
	message, err := suite.chargePoint.ParseMessage(mockMessage, suite.chargePoint.RequestState)
	require.Nil(t, message)
	require.Error(t, err)
	protoErr := err.(*ocpp.Error)
	require.NotNil(t, protoErr)
	assert.Equal(t, messageId, protoErr.MessageId)
	assert.Equal(t, ocppj.MessageTypeNotSupported, protoErr.Code)
	assert.Equal(t, fmt.Sprintf("Invalid message type ID %v", invalidTypeId), protoErr.Description)
}

func (suite *OcppJTestSuite) TestParseMessageUnsupported() {
	t := suite.T()
	mockMessage := make([]interface{}, 4)
	messageId := "12345"
	invalidAction := "SomeAction"
	// Test invalid message length
	mockMessage[0] = float64(ocppj.CALL) // Message Type ID
	mockMessage[1] = messageId           // Unique ID
	mockMessage[2] = invalidAction       // Action
	message, err := suite.chargePoint.ParseMessage(mockMessage, suite.chargePoint.RequestState)
	require.Nil(t, message)
	require.Error(t, err)
	protoErr := err.(*ocpp.Error)
	require.NotNil(t, protoErr)
	assert.Equal(t, messageId, protoErr.MessageId)
	assert.Equal(t, ocppj.NotSupported, protoErr.Code)
	assert.Equal(t, fmt.Sprintf("Unsupported feature %v", invalidAction), protoErr.Description)
}

func (suite *OcppJTestSuite) TestParseMessageInvalidCall() {
	t := suite.T()
	mockMessage := make([]interface{}, 3)
	messageId := "12345"
	// Test invalid message length
	mockMessage[0] = float64(ocppj.CALL) // Message Type ID
	mockMessage[1] = messageId           // Unique ID
	mockMessage[2] = MockFeatureName
	message, err := suite.chargePoint.ParseMessage(mockMessage, suite.chargePoint.RequestState)
	require.Nil(t, message)
	require.Error(t, err)
	protoErr := err.(*ocpp.Error)
	require.NotNil(t, protoErr)
	assert.Equal(t, messageId, protoErr.MessageId)
	assert.Equal(t, ocppj.FormationViolation, protoErr.Code)
	assert.Equal(t, "Invalid Call message. Expected array length 4", protoErr.Description)
}

func (suite *OcppJTestSuite) TestParseMessageInvalidCallResult() {
	t := suite.T()
	mockMessage := make([]interface{}, 3)
	messageId := "12345"
	mockConfirmation := newMockConfirmation("testValue")
	// Test invalid message length
	mockMessage[0] = float64(ocppj.CALL_RESULT) // Message Type ID
	mockMessage[1] = messageId                  // Unique ID
	mockMessage[2] = mockConfirmation
	message, err := suite.chargePoint.ParseMessage(mockMessage, suite.chargePoint.RequestState)
	// Both message and error should be nil
	require.Nil(t, message)
	require.NoError(t, err)
}

func (suite *OcppJTestSuite) TestParseMessageInvalidCallError() {
	t := suite.T()
	mockMessage := make([]interface{}, 3)
	messageId := "12345"
	pendingRequest := newMockRequest("request")
	// Test invalid message length
	mockMessage[0] = float64(ocppj.CALL_ERROR) // Message Type ID
	mockMessage[1] = messageId                 // Unique ID
	mockMessage[2] = ocppj.GenericError
	suite.chargePoint.RequestState.AddPendingRequest(messageId, pendingRequest) // Manually add a pending request, so that response is not rejected
	message, err := suite.chargePoint.ParseMessage(mockMessage, suite.chargePoint.RequestState)
	require.Nil(t, message)
	require.Error(t, err)
	protoErr := err.(*ocpp.Error)
	require.NotNil(t, protoErr)
	assert.Equal(t, messageId, protoErr.MessageId)
	assert.Equal(t, ocppj.FormationViolation, protoErr.Code)
	assert.Equal(t, "Invalid Call Error message. Expected array length >= 4", protoErr.Description)
}

func (suite *OcppJTestSuite) TestParseMessageInvalidRequest() {
	t := suite.T()
	mockMessage := make([]interface{}, 4)
	messageId := "12345"
	// Test invalid request -> required field missing
	mockRequest := newMockRequest("")
	mockMessage[0] = float64(ocppj.CALL) // Message Type ID
	mockMessage[1] = messageId           // Unique ID
	mockMessage[2] = MockFeatureName
	mockMessage[3] = mockRequest
	message, err := suite.chargePoint.ParseMessage(mockMessage, suite.chargePoint.RequestState)
	require.Nil(t, message)
	require.Error(t, err)
	protoErr := err.(*ocpp.Error)
	require.NotNil(t, protoErr)
	assert.Equal(t, messageId, protoErr.MessageId)
	assert.Equal(t, ocppj.OccurrenceConstraintViolation, protoErr.Code)
	// Test invalid request -> max constraint wrong
	mockRequest.MockValue = "somelongvalue"
	message, err = suite.chargePoint.ParseMessage(mockMessage, suite.chargePoint.RequestState)
	require.Nil(t, message)
	require.Error(t, err)
	protoErr = err.(*ocpp.Error)
	require.NotNil(t, protoErr)
	assert.Equal(t, messageId, protoErr.MessageId)
	assert.Equal(t, ocppj.PropertyConstraintViolation, protoErr.Code)
}

func (suite *OcppJTestSuite) TestParseMessageInvalidConfirmation() {
	t := suite.T()
	mockMessage := make([]interface{}, 3)
	messageId := "12345"
	// Test invalid confirmation -> required field missing
	pendingRequest := newMockRequest("request")
	mockConfirmation := newMockConfirmation("")
	mockMessage[0] = float64(ocppj.CALL_RESULT) // Message Type ID
	mockMessage[1] = messageId                  // Unique ID
	mockMessage[2] = mockConfirmation
	suite.chargePoint.RequestState.AddPendingRequest(messageId, pendingRequest) // Manually add a pending request, so that response is not rejected
	message, err := suite.chargePoint.ParseMessage(mockMessage, suite.chargePoint.RequestState)
	require.Nil(t, message)
	require.Error(t, err)
	protoErr := err.(*ocpp.Error)
	require.NotNil(t, protoErr)
	assert.Equal(t, messageId, protoErr.MessageId)
	assert.Equal(t, ocppj.OccurrenceConstraintViolation, protoErr.Code)
	// Test invalid request -> max constraint wrong
	mockConfirmation.MockValue = "min"
	suite.chargePoint.RequestState.AddPendingRequest(messageId, pendingRequest) // Manually add a pending request, so that responses are not rejected
	message, err = suite.chargePoint.ParseMessage(mockMessage, suite.chargePoint.RequestState)
	require.Nil(t, message)
	require.Error(t, err)
	protoErr = err.(*ocpp.Error)
	require.NotNil(t, protoErr)
	assert.Equal(t, messageId, protoErr.MessageId)
	assert.Equal(t, ocppj.PropertyConstraintViolation, protoErr.Code)
}

func (suite *OcppJTestSuite) TestParseCall() {
	t := suite.T()
	mockMessage := make([]interface{}, 4)
	messageId := "12345"
	mockValue := "somevalue"
	mockRequest := newMockRequest(mockValue)
	// Test invalid message length
	mockMessage[0] = float64(ocppj.CALL) // Message Type ID
	mockMessage[1] = messageId           // Unique ID
	mockMessage[2] = MockFeatureName
	mockMessage[3] = mockRequest
	message, protoErr := suite.chargePoint.ParseMessage(mockMessage, suite.chargePoint.RequestState)
	assert.Nil(t, protoErr)
	assert.NotNil(t, message)
	assert.Equal(t, ocppj.CALL, message.GetMessageTypeId())
	assert.Equal(t, messageId, message.GetUniqueId())
	assert.IsType(t, new(ocppj.Call), message)
	call := message.(*ocppj.Call)
	assert.Equal(t, MockFeatureName, call.Action)
	assert.IsType(t, new(MockRequest), call.Payload)
	mockRequest = call.Payload.(*MockRequest)
	assert.Equal(t, mockValue, mockRequest.MockValue)
}

//TODO: implement further ocpp-j protocol tests

type testLogger struct {
	c chan string
}

func (l *testLogger) Debug(args ...interface{}) {
	l.c <- "debug"
}
func (l *testLogger) Debugf(format string, args ...interface{}) {
	l.c <- "debugf"
}
func (l *testLogger) Info(args ...interface{}) {
	l.c <- "info"
}
func (l *testLogger) Infof(format string, args ...interface{}) {
	l.c <- "infof"
}
func (l *testLogger) Error(args ...interface{}) {
	l.c <- "error"
}
func (l *testLogger) Errorf(format string, args ...interface{}) {
	l.c <- "errorf"
}

func (suite *OcppJTestSuite) TestLogger() {
	t := suite.T()
	logger := testLogger{c: make(chan string, 1)}
	// Test with custom logger
	ocppj.SetLogger(&logger)
	defer ocppj.SetLogger(&logging.VoidLogger{})
	// Expect an error
	arr, err := ocppj.ParseRawJsonMessage([]byte("[3,\"1234\",{}]"))
	require.NoError(t, err)
	_, _ = suite.chargePoint.ParseMessage(arr, suite.chargePoint.RequestState)
	s := <-logger.c
	assert.Equal(t, "infof", s)
	// Nil logger must cause a panic
	assertPanic(t, func() {
		ocppj.SetLogger(nil)
	}, func(r interface{}) {
		assert.Equal(t, "cannot set a nil logger", r.(string))
	})
}

type MockValidationError struct {
	tag       string
	namespace string
	param     string
	value     string
	typ       reflect.Type
}

func (m MockValidationError) ActualTag() string                 { return m.tag }
func (m MockValidationError) Tag() string                       { return m.tag }
func (m MockValidationError) Namespace() string                 { return m.namespace }
func (m MockValidationError) StructNamespace() string           { return m.namespace }
func (m MockValidationError) Field() string                     { return m.namespace }
func (m MockValidationError) StructField() string               { return m.namespace }
func (m MockValidationError) Value() interface{}                { return m.value }
func (m MockValidationError) Param() string                     { return m.param }
func (m MockValidationError) Kind() reflect.Kind                { return m.typ.Kind() }
func (m MockValidationError) Type() reflect.Type                { return m.typ }
func (m MockValidationError) Translate(ut ut.Translator) string { return "" }
func (m MockValidationError) Error() string                     { return fmt.Sprintf("some error for value %s", m.value) }

func TestMockOcppJ(t *testing.T) {
	suite.Run(t, new(ClientQueueTestSuite))
	suite.Run(t, new(ServerQueueMapTestSuite))
	suite.Run(t, new(ClientStateTestSuite))
	suite.Run(t, new(ServerStateTestSuite))
	suite.Run(t, new(ClientDispatcherTestSuite))
	suite.Run(t, new(ServerDispatcherTestSuite))
	suite.Run(t, new(OcppJTestSuite))
}
