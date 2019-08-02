package ocpp16_test

import (
	"fmt"
	"github.com/gorilla/websocket"
	ocpp16 "github.com/lorenzodonini/go-ocpp/ocpp1.6"
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

// ---------------------- MOCK CS CORE LISTENER ----------------------
type MockCentralSystemCoreListener struct {
	mock.Mock
}

func (coreListener MockCentralSystemCoreListener) OnAuthorize(chargePointId string, request *ocpp16.AuthorizeRequest) (confirmation *ocpp16.AuthorizeConfirmation, err error) {
	args := coreListener.MethodCalled("OnAuthorize", chargePointId, request)
	conf := args.Get(0).(*ocpp16.AuthorizeConfirmation)
	return conf, args.Error(1)
}

func (coreListener MockCentralSystemCoreListener) OnBootNotification(chargePointId string, request *ocpp16.BootNotificationRequest) (confirmation *ocpp16.BootNotificationConfirmation, err error) {
	args := coreListener.MethodCalled("OnBootNotification", chargePointId, request)
	conf := args.Get(0).(*ocpp16.BootNotificationConfirmation)
	return conf, args.Error(1)
}

// ---------------------- MOCK CP CORE LISTENER ----------------------
type MockChargePointCoreListener struct {
	mock.Mock
}

func (coreListener MockChargePointCoreListener) OnChangeAvailability(request *ocpp16.ChangeAvailabilityRequest) (confirmation *ocpp16.ChangeAvailabilityConfirmation, err error) {
	args := coreListener.MethodCalled("OnChangeAvailability", request)
	conf := args.Get(0).(*ocpp16.ChangeAvailabilityConfirmation)
	return conf, args.Error(1)
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

type expectedCentralSystemOptions struct {
	clientId              string
	rawWrittenMessage     []byte
	startReturnArgument   interface{}
	writeReturnArgument   interface{}
	forwardWrittenMessage bool
}

type expectedChargePointOptions struct {
	serverUrl             string
	clientId              string
	createChannelOnStart  bool
	channel               ws.Channel
	rawWrittenMessage     []byte
	startReturnArgument   interface{}
	writeReturnArgument   interface{}
	forwardWrittenMessage bool
}

func setupDefaultCentralSystemHandlers(suite *OcppV16TestSuite, coreListener ocpp16.CentralSystemCoreListener, options expectedCentralSystemOptions) {
	t := suite.T()
	suite.centralSystem.SetNewChargePointHandler(func(chargePointId string) {
		assert.Equal(t, options.clientId, chargePointId)
	})
	suite.centralSystem.SetCentralSystemCoreListener(coreListener)
	// TODO: parametrize return arguments
	suite.mockWsServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(options.startReturnArgument)
	suite.mockWsServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(options.writeReturnArgument).Run(func(args mock.Arguments) {
		clientId := args.String(0)
		data := args.Get(1)
		bytes := data.([]byte)
		assert.Equal(t, options.clientId, clientId)
		if options.rawWrittenMessage != nil {
			assert.NotNil(t, bytes)
			assert.Equal(t, options.rawWrittenMessage, bytes)
		}
		if options.forwardWrittenMessage {
			// Notify client of incoming response
			err := suite.mockWsClient.MessageHandler(bytes)
			assert.Nil(t, err)
		}
	})
}

func setupDefaultChargePointHandlers(suite *OcppV16TestSuite, coreListener ocpp16.ChargePointCoreListener, options expectedChargePointOptions) {
	t := suite.T()
	suite.chargePoint.SetChargePointCoreListener(coreListener)
	suite.mockWsClient.On("Start", mock.AnythingOfType("string")).Return(options.startReturnArgument).Run(func(args mock.Arguments) {
		u := args.String(0)
		assert.Equal(t, fmt.Sprintf("%s/%s", options.serverUrl, options.clientId), u)
		// Notify server of incoming connection
		if options.createChannelOnStart {
			suite.mockWsServer.NewClientHandler(options.channel)
		}
	})
	suite.mockWsClient.On("Write", mock.Anything).Return(options.writeReturnArgument).Run(func(args mock.Arguments) {
		data := args.Get(0)
		bytes := data.([]byte)
		if options.rawWrittenMessage != nil {
			assert.NotNil(t, bytes)
			assert.Equal(t, options.rawWrittenMessage, bytes)
		}
		// Notify server of incoming request
		if options.forwardWrittenMessage {
			err := suite.mockWsServer.MessageHandler(options.channel, bytes)
			assert.Nil(t, err)
		}
	})
}

func assertDateTimeEquality(t *testing.T, expected ocpp16.DateTime, actual ocpp16.DateTime) {
	assert.Equal(t, expected.Time.Format(ocpp16.ISO8601), actual.Time.Format(ocpp16.ISO8601))
}

type RequestTestEntry struct {
	Request       ocppj.Request
	ExpectedValid bool
}

type ConfirmationTestEntry struct {
	Confirmation  ocppj.Confirmation
	ExpectedValid bool
}

func ExecuteRequestTestTable(t *testing.T, testTable []RequestTestEntry) {
	for _, testCase := range testTable {
		err := Validate.Struct(testCase.Request)
		assert.Equal(t, testCase.ExpectedValid, err == nil)
	}
}

func ExecuteConfirmationTestTable(t *testing.T, testTable []ConfirmationTestEntry) {
	for _, testCase := range testTable {
		err := Validate.Struct(testCase.Confirmation)
		assert.Equal(t, testCase.ExpectedValid, err == nil)
	}
}

var Validate = validator.New()

// ---------------------- TESTS ----------------------
type OcppV16TestSuite struct {
	suite.Suite
	ocppjChargePoint   *ocppj.ChargePoint
	ocppjCentralSystem *ocppj.CentralSystem
	mockWsServer       *MockWebsocketServer
	mockWsClient       *MockWebsocketClient
	chargePoint        ocpp16.ChargePoint
	centralSystem      ocpp16.CentralSystem
	messageIdGenerator TestRandomIdGenerator
}

type TestRandomIdGenerator struct {
	generator func() string
}

func (testGenerator * TestRandomIdGenerator) generateId() string {
	return testGenerator.generator()
}

var defaultMessageId = "1234"

func (suite *OcppV16TestSuite) SetupTest() {
	// TODO: replace coreProfile with ocpp16.CoreProfile
	coreProfile := ocppj.NewProfile("core", ocpp16.BootNotificationFeature{}, ocpp16.AuthorizeFeature{}, ocpp16.ChangeAvailabilityFeature{})
	mockClient := MockWebsocketClient{}
	mockServer := MockWebsocketServer{}
	suite.mockWsClient = &mockClient
	suite.mockWsServer = &mockServer
	suite.ocppjChargePoint = ocppj.NewChargePoint("test_id", suite.mockWsClient, coreProfile)
	suite.ocppjCentralSystem = ocppj.NewCentralSystem(suite.mockWsServer, coreProfile)
	suite.chargePoint = ocpp16.NewChargePoint("test_id", suite.ocppjChargePoint, suite.mockWsClient)
	suite.centralSystem = ocpp16.NewCentralSystem(suite.ocppjCentralSystem, suite.mockWsServer)
	suite.messageIdGenerator = TestRandomIdGenerator{generator: func() string {
		return defaultMessageId
	}}
	ocppj.SetMessageIdGenerator(suite.messageIdGenerator.generateId)
}

//TODO: implement generic protocol tests

func TestOcpp16Protocol(t *testing.T) {
	suite.Run(t, new(OcppV16TestSuite))
}
