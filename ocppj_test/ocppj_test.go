package ocppj_test

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/lorenzodonini/go-ocpp/ocppj"
	"github.com/lorenzodonini/go-ocpp/ws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
	"testing"
)

// ---------------------- MOCK WEBSOCKET ----------------------
type MockWebSocket struct {
	id string
}

func (websocket MockWebSocket) GetId() string {
	return websocket.id
}

func NewMockWebSocket(id string) MockWebSocket {
	return MockWebSocket{id: id}
}

// ---------------------- MOCK WEBSOCKET SERVER ----------------------
type MockWebsocketServer struct {
	mock.Mock
	ws.WsServer
	MessageHandler   func(ws ws.Channel, data []byte) error
	NewClientHandler func(ws ws.Channel)
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

func (websocketServer *MockWebsocketServer) NewClient(websocketId string, client interface{}) {
	websocketServer.MethodCalled("NewClient", websocketId, client)
}

// ---------------------- MOCK WEBSOCKET CLIENT ----------------------
type MockWebsocketClient struct {
	mock.Mock
	ws.WsClient
	MessageHandler func(data []byte) error
}

func (websocketClient *MockWebsocketClient) Start(url string, dialOptions ...func(websocket.Dialer)) error {
	args := websocketClient.MethodCalled("Start", url)
	return args.Error(0)
}

func (websocketClient *MockWebsocketClient) Stop() {
	websocketClient.MethodCalled("Stop")
}

func (websocketClient *MockWebsocketClient) SetMessageHandler(handler func(data []byte) error) {
	websocketClient.MessageHandler = handler
}

//TODO: Write should return error, same as for server
func (websocketClient *MockWebsocketClient) Write(data []byte) error {
	args := websocketClient.MethodCalled("Write", data)
	return args.Error(0)
}

// ---------------------- MOCK FEATURE ----------------------
const (
	MockFeatureName = "Mock"
)

type MockRequest struct {
	mock.Mock
	MockValue string `json:"mockValue" validate:"required,max=10"`
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

func (f MockFeature) GetConfirmationType() reflect.Type {
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
				err = wsClient.Write(data)
				assert.Nil(t, err)
			}
		}
		return nil
	})
	return &wsClient
}

func ParseCall(endpoint *ocppj.Endpoint, json string, t *testing.T) *ocppj.Call {
	parsedData := ocppj.ParseJsonMessage(json)
	result, err := endpoint.ParseMessage(parsedData)
	assert.Nil(t, err)
	assert.NotNil(t, result)
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

func ParseCallResult(endpoint *ocppj.Endpoint, json string, t *testing.T) *ocppj.CallResult {
	parsedData := ocppj.ParseJsonMessage(json)
	result, err := endpoint.ParseMessage(parsedData)
	assert.Nil(t, err)
	assert.NotNil(t, result)
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

func ParseCallError(endpoint *ocppj.Endpoint, json string, t *testing.T) *ocppj.CallError {
	parsedData := ocppj.ParseJsonMessage(json)
	result, err := endpoint.ParseMessage(parsedData)
	assert.Nil(t, err)
	callError, ok := result.(*ocppj.CallError)
	assert.Equal(t, true, ok)
	assert.NotNil(t, callError)
	return callError
}

func CheckCallError(t *testing.T, callError *ocppj.CallError, expectedId string, expectedError ocppj.ErrorCode, expectedDescription string, expectedDetails interface{}) {
	assert.Equal(t, ocppj.CALL_ERROR, callError.GetMessageTypeId())
	assert.Equal(t, expectedId, callError.GetUniqueId())
	assert.Equal(t, expectedError, callError.ErrorCode)
	assert.Equal(t, expectedDescription, callError.ErrorDescription)
	assert.Equal(t, expectedDetails, callError.ErrorDetails)
	err := Validate.Struct(callError)
	assert.Nil(t, err)
}

var Validate = validator.New()

// ---------------------- TESTS ----------------------
type OcppJTestSuite struct {
	suite.Suite
	chargePoint   *ocppj.ChargePoint
	centralSystem *ocppj.CentralSystem
	mockServer    *MockWebsocketServer
	mockClient    *MockWebsocketClient
}

func TestOcppJ(t *testing.T) {
	suite.Run(t, new(OcppJTestSuite))
}

func (suite *OcppJTestSuite) SetupTest() {
	mockProfile := ocppj.NewProfile("mock", MockFeature{})
	mockClient := MockWebsocketClient{}
	mockServer := MockWebsocketServer{}
	suite.mockClient = &mockClient
	suite.mockServer = &mockServer
	suite.chargePoint = ocppj.NewChargePoint("mock_id", suite.mockClient, mockProfile)
	suite.centralSystem = ocppj.NewCentralSystem(suite.mockServer, mockProfile)
}

// Protocol functions test
func (suite *OcppJTestSuite) TestGetProfile() {
	t := suite.T()
	profile, ok := suite.chargePoint.GetProfile("mock")
	assert.True(t, ok)
	assert.NotNil(t, profile)
	feature := profile.GetFeature(MockFeatureName)
	assert.NotNil(t, feature)
	assert.Equal(t, reflect.TypeOf(MockRequest{}), feature.GetRequestType())
	assert.Equal(t, reflect.TypeOf(MockConfirmation{}), feature.GetConfirmationType())
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
	// Check that request was stored as pending request
	pendingRequest, exists := suite.chargePoint.GetPendingRequest(call.UniqueId)
	assert.True(t, exists)
	assert.NotNil(t, pendingRequest)
	suite.chargePoint.DeletePendingRequest(call.UniqueId)
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
	callError := suite.chargePoint.CreateCallError(mockUniqueId, ocppj.GenericError, mockDescription, mockDetails)
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
	message, protoErr := suite.chargePoint.ParseMessage(mockMessage)
	assert.Nil(t, message)
	assert.NotNil(t, protoErr)
	assert.Equal(t, "", protoErr.MessageId)
	assert.Equal(t, ocppj.FormationViolation, protoErr.ErrorCode)
	assert.Equal(t, "Invalid message. Expected array length >= 3", protoErr.Error.Error())
}

func (suite *OcppJTestSuite) TestParseMessageInvalidTypeId() {
	t := suite.T()
	mockMessage := make([]interface{}, 3)
	invalidTypeId := "2"
	messageId := "12345"
	// Test invalid message length
	mockMessage[0] = invalidTypeId 	// Message Type ID
	mockMessage[1] = messageId 		// Unique ID
	message, protoErr := suite.chargePoint.ParseMessage(mockMessage)
	assert.Nil(t, message)
	assert.NotNil(t, protoErr)
	assert.Equal(t, "", protoErr.MessageId)
	assert.Equal(t, ocppj.FormationViolation, protoErr.ErrorCode)
	assert.Equal(t, fmt.Sprintf("Invalid element %v at 0, expected message type (int)", invalidTypeId), protoErr.Error.Error())
}

func (suite *OcppJTestSuite) TestParseMessageInvalidMessageId() {
	t := suite.T()
	mockMessage := make([]interface{}, 3)
	invalidMessageId := 12345
	// Test invalid message length
	mockMessage[0] = float64(ocppj.CALL) // Message Type ID
	mockMessage[1] = invalidMessageId    // Unique ID
	message, protoErr := suite.chargePoint.ParseMessage(mockMessage)
	assert.Nil(t, message)
	assert.NotNil(t, protoErr)
	assert.Equal(t, "", protoErr.MessageId)
	assert.Equal(t, ocppj.FormationViolation, protoErr.ErrorCode)
	assert.Equal(t, fmt.Sprintf("Invalid element %v at 1, expected unique ID (string)", invalidMessageId), protoErr.Error.Error())
}

func (suite *OcppJTestSuite) TestParseMessageUnknownTypeId() {
	t := suite.T()
	mockMessage := make([]interface{}, 3)
	messageId := "12345"
	invalidTypeId := 1
	// Test invalid message length
	mockMessage[0] = float64(invalidTypeId)		// Message Type ID
	mockMessage[1] = messageId 					// Unique ID
	message, protoErr := suite.chargePoint.ParseMessage(mockMessage)
	assert.Nil(t, message)
	assert.NotNil(t, protoErr)
	assert.Equal(t, messageId, protoErr.MessageId)
	assert.Equal(t, ocppj.FormationViolation, protoErr.ErrorCode)
	assert.Equal(t, fmt.Sprintf("Invalid message type ID %v", invalidTypeId), protoErr.Error.Error())
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
	message, protoErr := suite.chargePoint.ParseMessage(mockMessage)
	assert.Nil(t, message)
	assert.NotNil(t, protoErr)
	assert.Equal(t, messageId, protoErr.MessageId)
	assert.Equal(t, ocppj.NotSupported, protoErr.ErrorCode)
	assert.Equal(t, fmt.Sprintf("Unsupported feature %v", invalidAction), protoErr.Error.Error())
}

func (suite *OcppJTestSuite) TestParseMessageInvalidCall() {
	t := suite.T()
	mockMessage := make([]interface{}, 3)
	messageId := "12345"
	// Test invalid message length
	mockMessage[0] = float64(ocppj.CALL) // Message Type ID
	mockMessage[1] = messageId           // Unique ID
	mockMessage[2] = MockFeatureName
	message, protoErr := suite.chargePoint.ParseMessage(mockMessage)
	assert.Nil(t, message)
	assert.NotNil(t, protoErr)
	assert.Equal(t, messageId, protoErr.MessageId)
	assert.Equal(t, ocppj.FormationViolation, protoErr.ErrorCode)
	assert.Equal(t, "Invalid Call message. Expected array length 4", protoErr.Error.Error())
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
	message, protoErr := suite.chargePoint.ParseMessage(mockMessage)
	// Both message and error should be nil
	assert.Nil(t, message)
	assert.Nil(t, protoErr)
}

func (suite *OcppJTestSuite) TestParseMessageInvalidCallError() {
	t := suite.T()
	mockMessage := make([]interface{}, 3)
	messageId := "12345"
	// Test invalid message length
	mockMessage[0] = float64(ocppj.CALL_ERROR) // Message Type ID
	mockMessage[1] = messageId                 // Unique ID
	mockMessage[2] = ocppj.GenericError
	message, protoErr := suite.chargePoint.ParseMessage(mockMessage)
	assert.Nil(t, message)
	assert.NotNil(t, protoErr)
	assert.Equal(t, messageId, protoErr.MessageId)
	assert.Equal(t, ocppj.FormationViolation, protoErr.ErrorCode)
	assert.Equal(t, "Invalid Call Error message. Expected array length >= 4", protoErr.Error.Error())
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
	message, protoErr := suite.chargePoint.ParseMessage(mockMessage)
	assert.Nil(t, message)
	assert.NotNil(t, protoErr)
	assert.Equal(t, messageId, protoErr.MessageId)
	assert.Equal(t, ocppj.OccurrenceConstraintViolation, protoErr.ErrorCode)
	// Test invalid request -> max constraint wrong
	mockRequest.MockValue = "somelongvalue"
	message, protoErr = suite.chargePoint.ParseMessage(mockMessage)
	assert.Nil(t, message)
	assert.NotNil(t, protoErr)
	assert.Equal(t, messageId, protoErr.MessageId)
	assert.Equal(t, ocppj.PropertyConstraintViolation, protoErr.ErrorCode)
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
	suite.chargePoint.AddPendingRequest(messageId, pendingRequest)
	message, protoErr := suite.chargePoint.ParseMessage(mockMessage)
	assert.Nil(t, message)
	assert.NotNil(t, protoErr)
	assert.Equal(t, messageId, protoErr.MessageId)
	assert.Equal(t, ocppj.OccurrenceConstraintViolation, protoErr.ErrorCode)
	// Test invalid request -> max constraint wrong
	mockConfirmation.MockValue = "min"
	suite.chargePoint.AddPendingRequest(messageId, pendingRequest)
	message, protoErr = suite.chargePoint.ParseMessage(mockMessage)
	assert.Nil(t, message)
	assert.NotNil(t, protoErr)
	assert.Equal(t, messageId, protoErr.MessageId)
	assert.Equal(t, ocppj.PropertyConstraintViolation, protoErr.ErrorCode)
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
	message, protoErr := suite.chargePoint.ParseMessage(mockMessage)
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
