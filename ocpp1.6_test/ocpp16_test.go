package ocpp16_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	"github.com/lorenzodonini/ocpp-go/ocppj"
	"github.com/lorenzodonini/ocpp-go/ws"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"reflect"
	"testing"
)

// ---------------------- MOCK WEBSOCKET ----------------------
type MockWebSocket struct {
	id string
}

func (websocket MockWebSocket) GetID() string {
	return websocket.id
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
	DisconnectedClientHandler func(ws ws.Channel)
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

func (websocketServer *MockWebsocketServer) NewClient(websocketId string, client interface{}) {
	websocketServer.MethodCalled("NewClient", websocketId, client)
}

// ---------------------- MOCK WEBSOCKET CLIENT ----------------------
type MockWebsocketClient struct {
	mock.Mock
	ws.WsClient
	MessageHandler func(data []byte) error
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

func (websocketClient *MockWebsocketClient) Write(data []byte) error {
	args := websocketClient.MethodCalled("Write", data)
	return args.Error(0)
}

func (websocketClient *MockWebsocketClient) AddOption(option interface{}) {
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

func (coreListener MockCentralSystemCoreListener) OnDataTransfer(chargePointId string, request *ocpp16.DataTransferRequest) (confirmation *ocpp16.DataTransferConfirmation, err error) {
	args := coreListener.MethodCalled("OnDataTransfer", chargePointId, request)
	conf := args.Get(0).(*ocpp16.DataTransferConfirmation)
	return conf, args.Error(1)
}

func (coreListener MockCentralSystemCoreListener) OnHeartbeat(chargePointId string, request *ocpp16.HeartbeatRequest) (confirmation *ocpp16.HeartbeatConfirmation, err error) {
	args := coreListener.MethodCalled("OnHeartbeat", chargePointId, request)
	conf := args.Get(0).(*ocpp16.HeartbeatConfirmation)
	return conf, args.Error(1)
}

func (coreListener MockCentralSystemCoreListener) OnMeterValues(chargePointId string, request *ocpp16.MeterValuesRequest) (confirmation *ocpp16.MeterValuesConfirmation, err error) {
	args := coreListener.MethodCalled("OnMeterValues", chargePointId, request)
	conf := args.Get(0).(*ocpp16.MeterValuesConfirmation)
	return conf, args.Error(1)
}

func (coreListener MockCentralSystemCoreListener) OnStartTransaction(chargePointId string, request *ocpp16.StartTransactionRequest) (confirmation *ocpp16.StartTransactionConfirmation, err error) {
	args := coreListener.MethodCalled("OnStartTransaction", chargePointId, request)
	conf := args.Get(0).(*ocpp16.StartTransactionConfirmation)
	return conf, args.Error(1)
}

func (coreListener MockCentralSystemCoreListener) OnStatusNotification(chargePointId string, request *ocpp16.StatusNotificationRequest) (confirmation *ocpp16.StatusNotificationConfirmation, err error) {
	args := coreListener.MethodCalled("OnStatusNotification", chargePointId, request)
	conf := args.Get(0).(*ocpp16.StatusNotificationConfirmation)
	return conf, args.Error(1)
}

func (coreListener MockCentralSystemCoreListener) OnStopTransaction(chargePointId string, request *ocpp16.StopTransactionRequest) (confirmation *ocpp16.StopTransactionConfirmation, err error) {
	args := coreListener.MethodCalled("OnStopTransaction", chargePointId, request)
	conf := args.Get(0).(*ocpp16.StopTransactionConfirmation)
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

func (coreListener MockChargePointCoreListener) OnDataTransfer(request *ocpp16.DataTransferRequest) (confirmation *ocpp16.DataTransferConfirmation, err error) {
	args := coreListener.MethodCalled("OnDataTransfer", request)
	conf := args.Get(0).(*ocpp16.DataTransferConfirmation)
	return conf, args.Error(1)
}

func (coreListener MockChargePointCoreListener) OnChangeConfiguration(request *ocpp16.ChangeConfigurationRequest) (confirmation *ocpp16.ChangeConfigurationConfirmation, err error) {
	args := coreListener.MethodCalled("OnChangeConfiguration", request)
	conf := args.Get(0).(*ocpp16.ChangeConfigurationConfirmation)
	return conf, args.Error(1)
}

func (coreListener MockChargePointCoreListener) OnClearCache(request *ocpp16.ClearCacheRequest) (confirmation *ocpp16.ClearCacheConfirmation, err error) {
	args := coreListener.MethodCalled("OnClearCache", request)
	conf := args.Get(0).(*ocpp16.ClearCacheConfirmation)
	return conf, args.Error(1)
}

func (coreListener MockChargePointCoreListener) OnGetConfiguration(request *ocpp16.GetConfigurationRequest) (confirmation *ocpp16.GetConfigurationConfirmation, err error) {
	args := coreListener.MethodCalled("OnGetConfiguration", request)
	conf := args.Get(0).(*ocpp16.GetConfigurationConfirmation)
	return conf, args.Error(1)
}

func (coreListener MockChargePointCoreListener) OnReset(request *ocpp16.ResetRequest) (confirmation *ocpp16.ResetConfirmation, err error) {
	args := coreListener.MethodCalled("OnReset", request)
	conf := args.Get(0).(*ocpp16.ResetConfirmation)
	return conf, args.Error(1)
}

func (coreListener MockChargePointCoreListener) OnUnlockConnector(request *ocpp16.UnlockConnectorRequest) (confirmation *ocpp16.UnlockConnectorConfirmation, err error) {
	args := coreListener.MethodCalled("OnUnlockConnector", request)
	conf := args.Get(0).(*ocpp16.UnlockConnectorConfirmation)
	return conf, args.Error(1)
}

func (coreListener MockChargePointCoreListener) OnRemoteStartTransaction(request *ocpp16.RemoteStartTransactionRequest) (confirmation *ocpp16.RemoteStartTransactionConfirmation, err error) {
	args := coreListener.MethodCalled("OnRemoteStartTransaction", request)
	conf := args.Get(0).(*ocpp16.RemoteStartTransactionConfirmation)
	return conf, args.Error(1)
}

func (coreListener MockChargePointCoreListener) OnRemoteStopTransaction(request *ocpp16.RemoteStopTransactionRequest) (confirmation *ocpp16.RemoteStopTransactionConfirmation, err error) {
	args := coreListener.MethodCalled("OnRemoteStopTransaction", request)
	conf := args.Get(0).(*ocpp16.RemoteStopTransactionConfirmation)
	return conf, args.Error(1)
}

// ---------------------- MOCK CS LOCAL AUTH LIST LISTENER ----------------------
type MockCentralSystemLocalAuthListListener struct {
	mock.Mock
}

// ---------------------- MOCK CP LOCAL AUTH LIST LISTENER ----------------------
type MockChargePointLocalAuthListListener struct {
	mock.Mock
}

func (localAuthListListener MockChargePointLocalAuthListListener) OnGetLocalListVersion(request *ocpp16.GetLocalListVersionRequest) (confirmation *ocpp16.GetLocalListVersionConfirmation, err error) {
	args := localAuthListListener.MethodCalled("OnGetLocalListVersion", request)
	conf := args.Get(0).(*ocpp16.GetLocalListVersionConfirmation)
	return conf, args.Error(1)
}

func (localAuthListListener MockChargePointLocalAuthListListener) OnSendLocalList(request *ocpp16.SendLocalListRequest) (confirmation *ocpp16.SendLocalListConfirmation, err error) {
	args := localAuthListListener.MethodCalled("OnSendLocalList", request)
	conf := args.Get(0).(*ocpp16.SendLocalListConfirmation)
	return conf, args.Error(1)
}

// ---------------------- MOCK CS FIRMWARE MANAGEMENT LISTENER ----------------------
type MockCentralSystemFirmwareManagementListener struct {
	mock.Mock
}

func (firmwareListener MockCentralSystemFirmwareManagementListener) OnDiagnosticsStatusNotification(chargePointId string, request *ocpp16.DiagnosticsStatusNotificationRequest) (confirmation *ocpp16.DiagnosticsStatusNotificationConfirmation, err error) {
	args := firmwareListener.MethodCalled("OnDiagnosticsStatusNotification", chargePointId, request)
	conf := args.Get(0).(*ocpp16.DiagnosticsStatusNotificationConfirmation)
	return conf, args.Error(1)
}

func (firmwareListener MockCentralSystemFirmwareManagementListener) OnFirmwareStatusNotification(chargePointId string, request *ocpp16.FirmwareStatusNotificationRequest) (confirmation *ocpp16.FirmwareStatusNotificationConfirmation, err error) {
	args := firmwareListener.MethodCalled("OnFirmwareStatusNotification", chargePointId, request)
	conf := args.Get(0).(*ocpp16.FirmwareStatusNotificationConfirmation)
	return conf, args.Error(1)
}

// ---------------------- MOCK CP FIRMWARE MANAGEMENT LISTENER ----------------------
type MockChargePointFirmwareManagementListener struct {
	mock.Mock
}

func (firmwareListener MockChargePointFirmwareManagementListener) OnGetDiagnostics(request *ocpp16.GetDiagnosticsRequest) (confirmation *ocpp16.GetDiagnosticsConfirmation, err error) {
	args := firmwareListener.MethodCalled("OnGetDiagnostics", request)
	conf := args.Get(0).(*ocpp16.GetDiagnosticsConfirmation)
	return conf, args.Error(1)
}

func (firmwareListener MockChargePointFirmwareManagementListener) OnUpdateFirmware(request *ocpp16.UpdateFirmwareRequest) (confirmation *ocpp16.UpdateFirmwareConfirmation, err error) {
	args := firmwareListener.MethodCalled("OnUpdateFirmware", request)
	conf := args.Get(0).(*ocpp16.UpdateFirmwareConfirmation)
	return conf, args.Error(1)
}

// ---------------------- MOCK CS RESERVATION LISTENER ----------------------
type MockCentralSystemReservationListener struct {
	mock.Mock
}

// ---------------------- MOCK CP RESERVATION LISTENER ----------------------
type MockChargePointReservationListener struct {
	mock.Mock
}

func (reservationListener MockChargePointReservationListener) OnReserveNow(request *ocpp16.ReserveNowRequest) (confirmation *ocpp16.ReserveNowConfirmation, err error) {
	args := reservationListener.MethodCalled("OnReserveNow", request)
	conf := args.Get(0).(*ocpp16.ReserveNowConfirmation)
	return conf, args.Error(1)
}

func (reservationListener MockChargePointReservationListener) OnCancelReservation(request *ocpp16.CancelReservationRequest) (confirmation *ocpp16.CancelReservationConfirmation, err error) {
	args := reservationListener.MethodCalled("OnCancelReservation", request)
	conf := args.Get(0).(*ocpp16.CancelReservationConfirmation)
	return conf, args.Error(1)
}

// ---------------------- MOCK CS REMOTE TRIGGER LISTENER ----------------------
type MockCentralSystemRemoteTriggerListener struct {
	mock.Mock
}

// ---------------------- MOCK CP REMOTE TRIGGER LISTENER ----------------------
type MockChargePointRemoteTriggerListener struct {
	mock.Mock
}

func (remoteTriggerListener MockChargePointRemoteTriggerListener) OnTriggerMessage(request *ocpp16.TriggerMessageRequest) (confirmation *ocpp16.TriggerMessageConfirmation, err error) {
	args := remoteTriggerListener.MethodCalled("OnTriggerMessage", request)
	conf := args.Get(0).(*ocpp16.TriggerMessageConfirmation)
	return conf, args.Error(1)
}

// ---------------------- MOCK CS SMART CHARGING LISTENER ----------------------
type MockCentralSystemSmartChargingListener struct {
	mock.Mock
}

// ---------------------- MOCK CP SMART CHARGING LISTENER ----------------------
type MockChargePointSmartChargingListener struct {
	mock.Mock
}

func (smartChargingListener MockChargePointSmartChargingListener) OnSetChargingProfile(request *ocpp16.SetChargingProfileRequest) (confirmation *ocpp16.SetChargingProfileConfirmation, err error) {
	args := smartChargingListener.MethodCalled("OnSetChargingProfile", request)
	conf := args.Get(0).(*ocpp16.SetChargingProfileConfirmation)
	return conf, args.Error(1)
}

func (smartChargingListener MockChargePointSmartChargingListener) OnClearChargingProfile(request *ocpp16.ClearChargingProfileRequest) (confirmation *ocpp16.ClearChargingProfileConfirmation, err error) {
	args := smartChargingListener.MethodCalled("OnClearChargingProfile", request)
	conf := args.Get(0).(*ocpp16.ClearChargingProfileConfirmation)
	return conf, args.Error(1)
}

func (smartChargingListener MockChargePointSmartChargingListener) OnGetCompositeSchedule(request *ocpp16.GetCompositeScheduleRequest) (confirmation *ocpp16.GetCompositeScheduleConfirmation, err error) {
	args := smartChargingListener.MethodCalled("OnGetCompositeSchedule", request)
	conf := args.Get(0).(*ocpp16.GetCompositeScheduleConfirmation)
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
				err = wsServer.Write(ws.GetID(), data)
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

func testUnsupportedRequestFromChargePoint(suite *OcppV16TestSuite, request ocpp.Request, requestJson string, messageId string) {
	t := suite.T()
	wsId := "test_id"
	wsUrl := "someUrl"
	expectedError := fmt.Sprintf("unsupported action %v on charge point, cannot send request", request.GetFeatureName())
	errorDescription := fmt.Sprintf("unsupported action %v on central system", request.GetFeatureName())
	errorJson := fmt.Sprintf(`[4,"%v","%v","%v",null]`, messageId, ocppj.NotSupported, errorDescription)
	channel := NewMockWebSocket(wsId)

	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(errorJson), forwardWrittenMessage: false})
	coreListener := MockCentralSystemCoreListener{}
	setupDefaultCentralSystemHandlers(suite, coreListener, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(errorJson), forwardWrittenMessage: true})
	resultChannel := make(chan bool, 1)
	suite.ocppjChargePoint.SetErrorHandler(func(err *ocpp.Error, details interface{}) {
		assert.Equal(t, messageId, err.MessageId)
		assert.Equal(t, ocppj.NotSupported, err.Code)
		assert.Equal(t, errorDescription, err.Description)
		assert.Nil(t, details)
		resultChannel <- true
	})
	// Start
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	assert.Nil(t, err)
	// Run request test
	err = suite.chargePoint.SendRequestAsync(request, func(confirmation ocpp.Confirmation, err error) {
		t.Fail()
	})
	assert.Error(t, err)
	assert.Equal(t, expectedError, err.Error())
	// Run response test
	suite.ocppjChargePoint.AddPendingRequest(messageId, request)
	err = suite.mockWsServer.MessageHandler(channel, []byte(requestJson))
	assert.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func testUnsupportedRequestFromCentralSystem(suite *OcppV16TestSuite, request ocpp.Request, requestJson string, messageId string) {
	t := suite.T()
	wsId := "test_id"
	wsUrl := "someUrl"
	expectedError := fmt.Sprintf("unsupported action %v on central system, cannot send request", request.GetFeatureName())
	errorDescription := fmt.Sprintf("unsupported action %v on charge point", request.GetFeatureName())
	errorJson := fmt.Sprintf(`[4,"%v","%v","%v",null]`, messageId, ocppj.NotSupported, errorDescription)
	channel := NewMockWebSocket(wsId)

	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: false})
	coreListener := MockChargePointCoreListener{}
	setupDefaultChargePointHandlers(suite, coreListener, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(errorJson), forwardWrittenMessage: true})
	suite.ocppjCentralSystem.SetErrorHandler(func(chargePointId string, err *ocpp.Error, details interface{}) {
		assert.Equal(t, messageId, err.MessageId)
		assert.Equal(t, wsId, chargePointId)
		assert.Equal(t, ocppj.NotSupported, err.Code)
		assert.Equal(t, errorDescription, err.Description)
		assert.Nil(t, details)
	})
	// Start
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	assert.Nil(t, err)
	// Run request test
	err = suite.centralSystem.SendRequestAsync(wsId, request, func(confirmation ocpp.Confirmation, err error) {
		t.Fail()
	})
	assert.Error(t, err)
	assert.Equal(t, expectedError, err.Error())
	// Run response test
	suite.ocppjCentralSystem.AddPendingRequest(messageId, request)
	err = suite.mockWsClient.MessageHandler([]byte(requestJson))
	assert.Nil(t, err)
}

type GenericTestEntry struct {
	Element       interface{}
	ExpectedValid bool
}

type RequestTestEntry struct {
	Request       ocpp.Request
	ExpectedValid bool
}

type ConfirmationTestEntry struct {
	Confirmation  ocpp.Confirmation
	ExpectedValid bool
}

// TODO: pass expected error value for improved validation and error message
func ExecuteGenericTestTable(t *testing.T, testTable []GenericTestEntry) {
	for _, testCase := range testTable {
		err := ocpp16.Validate.Struct(testCase.Element)
		if err != nil {
			assert.Equal(t, testCase.ExpectedValid, false, err.Error())
		} else {
			assert.Equal(t, testCase.ExpectedValid, true, "%v is valid", testCase.Element)
		}
	}
}

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

func (testGenerator *TestRandomIdGenerator) generateId() string {
	return testGenerator.generator()
}

var defaultMessageId = "1234"

func (suite *OcppV16TestSuite) SetupTest() {
	coreProfile := ocpp16.CoreProfile
	localAuthListProfile := ocpp16.LocalAuthListProfile
	firmwareProfile := ocpp16.FirmwareManagementProfile
	reservationProfile := ocpp16.ReservationProfile
	remoteTriggerProfile := ocpp16.RemoteTriggerProfile
	smartChargingProfile := ocpp16.SmartChargingProfile
	mockClient := MockWebsocketClient{}
	mockServer := MockWebsocketServer{}
	suite.mockWsClient = &mockClient
	suite.mockWsServer = &mockServer
	suite.ocppjChargePoint = ocppj.NewChargePoint("test_id", suite.mockWsClient, coreProfile, localAuthListProfile, firmwareProfile, reservationProfile, remoteTriggerProfile, smartChargingProfile)
	suite.ocppjCentralSystem = ocppj.NewCentralSystem(suite.mockWsServer, coreProfile, localAuthListProfile, firmwareProfile, reservationProfile, remoteTriggerProfile, smartChargingProfile)
	suite.chargePoint = ocpp16.NewChargePoint("test_id", suite.ocppjChargePoint, suite.mockWsClient)
	suite.centralSystem = ocpp16.NewCentralSystem(suite.ocppjCentralSystem, suite.mockWsServer)
	suite.messageIdGenerator = TestRandomIdGenerator{generator: func() string {
		return defaultMessageId
	}}
	ocppj.SetMessageIdGenerator(suite.messageIdGenerator.generateId)
}

//TODO: implement generic protocol tests

func TestOcpp16Protocol(t *testing.T) {
	logrus.SetLevel(logrus.PanicLevel)
	suite.Run(t, new(OcppV16TestSuite))
}
