package ocpp16_test

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp"
	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/certificates"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/extendedtriggermessage"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/firmware"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/remotetrigger"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/reservation"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/securefirmware"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/security"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/smartcharging"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/lorenzodonini/ocpp-go/ocppj"
	"github.com/lorenzodonini/ocpp-go/ws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
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

// Default queue capacity
const queueCapacity = 10

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

func (f *MockFeature) GetFeatureName() string {
	return MockFeatureName
}

func (f *MockFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(MockRequest{})
}

func (f *MockFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(MockConfirmation{})
}

func (r *MockRequest) GetFeatureName() string {
	return MockFeatureName
}

func (c *MockConfirmation) GetFeatureName() string {
	return MockFeatureName
}

// ---------------------- MOCK CS CORE LISTENER ----------------------
type MockCentralSystemCoreListener struct {
	mock.Mock
}

func (coreListener *MockCentralSystemCoreListener) OnAuthorize(chargePointId string, request *core.AuthorizeRequest) (confirmation *core.AuthorizeConfirmation, err error) {
	args := coreListener.MethodCalled("OnAuthorize", chargePointId, request)
	conf := args.Get(0).(*core.AuthorizeConfirmation)
	return conf, args.Error(1)
}

func (coreListener *MockCentralSystemCoreListener) OnBootNotification(chargePointId string, request *core.BootNotificationRequest) (confirmation *core.BootNotificationConfirmation, err error) {
	args := coreListener.MethodCalled("OnBootNotification", chargePointId, request)
	conf := args.Get(0).(*core.BootNotificationConfirmation)
	return conf, args.Error(1)
}

func (coreListener *MockCentralSystemCoreListener) OnDataTransfer(chargePointId string, request *core.DataTransferRequest) (confirmation *core.DataTransferConfirmation, err error) {
	args := coreListener.MethodCalled("OnDataTransfer", chargePointId, request)
	rawConf := args.Get(0)
	err = args.Error(1)
	if rawConf != nil {
		confirmation = rawConf.(*core.DataTransferConfirmation)
	}
	return
}

func (coreListener *MockCentralSystemCoreListener) OnHeartbeat(chargePointId string, request *core.HeartbeatRequest) (confirmation *core.HeartbeatConfirmation, err error) {
	args := coreListener.MethodCalled("OnHeartbeat", chargePointId, request)
	conf := args.Get(0).(*core.HeartbeatConfirmation)
	return conf, args.Error(1)
}

func (coreListener *MockCentralSystemCoreListener) OnMeterValues(chargePointId string, request *core.MeterValuesRequest) (confirmation *core.MeterValuesConfirmation, err error) {
	args := coreListener.MethodCalled("OnMeterValues", chargePointId, request)
	conf := args.Get(0).(*core.MeterValuesConfirmation)
	return conf, args.Error(1)
}

func (coreListener *MockCentralSystemCoreListener) OnStartTransaction(chargePointId string, request *core.StartTransactionRequest) (confirmation *core.StartTransactionConfirmation, err error) {
	args := coreListener.MethodCalled("OnStartTransaction", chargePointId, request)
	conf := args.Get(0).(*core.StartTransactionConfirmation)
	return conf, args.Error(1)
}

func (coreListener *MockCentralSystemCoreListener) OnStatusNotification(chargePointId string, request *core.StatusNotificationRequest) (confirmation *core.StatusNotificationConfirmation, err error) {
	args := coreListener.MethodCalled("OnStatusNotification", chargePointId, request)
	conf := args.Get(0).(*core.StatusNotificationConfirmation)
	return conf, args.Error(1)
}

func (coreListener *MockCentralSystemCoreListener) OnStopTransaction(chargePointId string, request *core.StopTransactionRequest) (confirmation *core.StopTransactionConfirmation, err error) {
	args := coreListener.MethodCalled("OnStopTransaction", chargePointId, request)
	conf := args.Get(0).(*core.StopTransactionConfirmation)
	return conf, args.Error(1)
}

// ---------------------- MOCK CP CORE LISTENER ----------------------
type MockChargePointCoreListener struct {
	mock.Mock
}

func (coreListener *MockChargePointCoreListener) OnChangeAvailability(request *core.ChangeAvailabilityRequest) (confirmation *core.ChangeAvailabilityConfirmation, err error) {
	args := coreListener.MethodCalled("OnChangeAvailability", request)
	conf := args.Get(0).(*core.ChangeAvailabilityConfirmation)
	return conf, args.Error(1)
}

func (coreListener *MockChargePointCoreListener) OnDataTransfer(request *core.DataTransferRequest) (confirmation *core.DataTransferConfirmation, err error) {
	args := coreListener.MethodCalled("OnDataTransfer", request)
	rawConf := args.Get(0)
	err = args.Error(1)
	if rawConf != nil {
		confirmation = rawConf.(*core.DataTransferConfirmation)
	}
	return
}

func (coreListener *MockChargePointCoreListener) OnChangeConfiguration(request *core.ChangeConfigurationRequest) (confirmation *core.ChangeConfigurationConfirmation, err error) {
	args := coreListener.MethodCalled("OnChangeConfiguration", request)
	conf := args.Get(0).(*core.ChangeConfigurationConfirmation)
	return conf, args.Error(1)
}

func (coreListener *MockChargePointCoreListener) OnClearCache(request *core.ClearCacheRequest) (confirmation *core.ClearCacheConfirmation, err error) {
	args := coreListener.MethodCalled("OnClearCache", request)
	conf := args.Get(0).(*core.ClearCacheConfirmation)
	return conf, args.Error(1)
}

func (coreListener *MockChargePointCoreListener) OnGetConfiguration(request *core.GetConfigurationRequest) (confirmation *core.GetConfigurationConfirmation, err error) {
	args := coreListener.MethodCalled("OnGetConfiguration", request)
	conf := args.Get(0).(*core.GetConfigurationConfirmation)
	return conf, args.Error(1)
}

func (coreListener *MockChargePointCoreListener) OnReset(request *core.ResetRequest) (confirmation *core.ResetConfirmation, err error) {
	args := coreListener.MethodCalled("OnReset", request)
	conf := args.Get(0).(*core.ResetConfirmation)
	return conf, args.Error(1)
}

func (coreListener *MockChargePointCoreListener) OnUnlockConnector(request *core.UnlockConnectorRequest) (confirmation *core.UnlockConnectorConfirmation, err error) {
	args := coreListener.MethodCalled("OnUnlockConnector", request)
	conf := args.Get(0).(*core.UnlockConnectorConfirmation)
	return conf, args.Error(1)
}

func (coreListener *MockChargePointCoreListener) OnRemoteStartTransaction(request *core.RemoteStartTransactionRequest) (confirmation *core.RemoteStartTransactionConfirmation, err error) {
	args := coreListener.MethodCalled("OnRemoteStartTransaction", request)
	conf := args.Get(0).(*core.RemoteStartTransactionConfirmation)
	return conf, args.Error(1)
}

func (coreListener *MockChargePointCoreListener) OnRemoteStopTransaction(request *core.RemoteStopTransactionRequest) (confirmation *core.RemoteStopTransactionConfirmation, err error) {
	args := coreListener.MethodCalled("OnRemoteStopTransaction", request)
	conf := args.Get(0).(*core.RemoteStopTransactionConfirmation)
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

func (localAuthListListener *MockChargePointLocalAuthListListener) OnGetLocalListVersion(request *localauth.GetLocalListVersionRequest) (confirmation *localauth.GetLocalListVersionConfirmation, err error) {
	args := localAuthListListener.MethodCalled("OnGetLocalListVersion", request)
	conf := args.Get(0).(*localauth.GetLocalListVersionConfirmation)
	return conf, args.Error(1)
}

func (localAuthListListener *MockChargePointLocalAuthListListener) OnSendLocalList(request *localauth.SendLocalListRequest) (confirmation *localauth.SendLocalListConfirmation, err error) {
	args := localAuthListListener.MethodCalled("OnSendLocalList", request)
	conf := args.Get(0).(*localauth.SendLocalListConfirmation)
	return conf, args.Error(1)
}

// ---------------------- MOCK CS FIRMWARE MANAGEMENT LISTENER ----------------------
type MockCentralSystemFirmwareManagementListener struct {
	mock.Mock
}

func (firmwareListener *MockCentralSystemFirmwareManagementListener) OnDiagnosticsStatusNotification(chargePointId string, request *firmware.DiagnosticsStatusNotificationRequest) (confirmation *firmware.DiagnosticsStatusNotificationConfirmation, err error) {
	args := firmwareListener.MethodCalled("OnDiagnosticsStatusNotification", chargePointId, request)
	conf := args.Get(0).(*firmware.DiagnosticsStatusNotificationConfirmation)
	return conf, args.Error(1)
}

func (firmwareListener *MockCentralSystemFirmwareManagementListener) OnFirmwareStatusNotification(chargePointId string, request *firmware.FirmwareStatusNotificationRequest) (confirmation *firmware.FirmwareStatusNotificationConfirmation, err error) {
	args := firmwareListener.MethodCalled("OnFirmwareStatusNotification", chargePointId, request)
	conf := args.Get(0).(*firmware.FirmwareStatusNotificationConfirmation)
	return conf, args.Error(1)
}

// ---------------------- MOCK CP FIRMWARE MANAGEMENT LISTENER ----------------------
type MockChargePointFirmwareManagementListener struct {
	mock.Mock
}

func (firmwareListener *MockChargePointFirmwareManagementListener) OnGetDiagnostics(request *firmware.GetDiagnosticsRequest) (confirmation *firmware.GetDiagnosticsConfirmation, err error) {
	args := firmwareListener.MethodCalled("OnGetDiagnostics", request)
	conf := args.Get(0).(*firmware.GetDiagnosticsConfirmation)
	return conf, args.Error(1)
}

func (firmwareListener *MockChargePointFirmwareManagementListener) OnUpdateFirmware(request *firmware.UpdateFirmwareRequest) (confirmation *firmware.UpdateFirmwareConfirmation, err error) {
	args := firmwareListener.MethodCalled("OnUpdateFirmware", request)
	conf := args.Get(0).(*firmware.UpdateFirmwareConfirmation)
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

func (reservationListener *MockChargePointReservationListener) OnReserveNow(request *reservation.ReserveNowRequest) (confirmation *reservation.ReserveNowConfirmation, err error) {
	args := reservationListener.MethodCalled("OnReserveNow", request)
	conf := args.Get(0).(*reservation.ReserveNowConfirmation)
	return conf, args.Error(1)
}

func (reservationListener *MockChargePointReservationListener) OnCancelReservation(request *reservation.CancelReservationRequest) (confirmation *reservation.CancelReservationConfirmation, err error) {
	args := reservationListener.MethodCalled("OnCancelReservation", request)
	conf := args.Get(0).(*reservation.CancelReservationConfirmation)
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

func (remoteTriggerListener *MockChargePointRemoteTriggerListener) OnTriggerMessage(request *remotetrigger.TriggerMessageRequest) (confirmation *remotetrigger.TriggerMessageConfirmation, err error) {
	args := remoteTriggerListener.MethodCalled("OnTriggerMessage", request)
	conf := args.Get(0).(*remotetrigger.TriggerMessageConfirmation)
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

func (smartChargingListener *MockChargePointSmartChargingListener) OnSetChargingProfile(request *smartcharging.SetChargingProfileRequest) (confirmation *smartcharging.SetChargingProfileConfirmation, err error) {
	args := smartChargingListener.MethodCalled("OnSetChargingProfile", request)
	conf := args.Get(0).(*smartcharging.SetChargingProfileConfirmation)
	return conf, args.Error(1)
}

func (smartChargingListener *MockChargePointSmartChargingListener) OnClearChargingProfile(request *smartcharging.ClearChargingProfileRequest) (confirmation *smartcharging.ClearChargingProfileConfirmation, err error) {
	args := smartChargingListener.MethodCalled("OnClearChargingProfile", request)
	conf := args.Get(0).(*smartcharging.ClearChargingProfileConfirmation)
	return conf, args.Error(1)
}

func (smartChargingListener *MockChargePointSmartChargingListener) OnGetCompositeSchedule(request *smartcharging.GetCompositeScheduleRequest) (confirmation *smartcharging.GetCompositeScheduleConfirmation, err error) {
	args := smartChargingListener.MethodCalled("OnGetCompositeSchedule", request)
	conf := args.Get(0).(*smartcharging.GetCompositeScheduleConfirmation)
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

func setupDefaultCentralSystemHandlers(suite *OcppV16TestSuite, coreListener core.CentralSystemHandler, options expectedCentralSystemOptions) {
	t := suite.T()
	suite.centralSystem.SetNewChargePointHandler(func(chargePoint ocpp16.ChargePointConnection) {
		assert.Equal(t, options.clientId, chargePoint.ID())
	})
	suite.centralSystem.SetCoreHandler(coreListener)
	suite.mockWsServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(options.startReturnArgument)
	suite.mockWsServer.On("Stop").Return()
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

func setupDefaultChargePointHandlers(suite *OcppV16TestSuite, coreListener core.ChargePointHandler, options expectedChargePointOptions) {
	t := suite.T()
	suite.chargePoint.SetCoreHandler(coreListener)
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

func assertDateTimeEquality(t *testing.T, expected types.DateTime, actual types.DateTime) {
	assert.Equal(t, expected.FormatTimestamp(), actual.FormatTimestamp())
}

func testUnsupportedRequestFromChargePoint(suite *OcppV16TestSuite, request ocpp.Request, requestJson string, messageId string) {
	t := suite.T()
	wsId := "test_id"
	wsUrl := "someUrl"
	expectedError := fmt.Sprintf("unsupported action %v on charge point, cannot send request", request.GetFeatureName())
	errorDescription := fmt.Sprintf("unsupported action %v on central system", request.GetFeatureName())
	errorJson := fmt.Sprintf(`[4,"%v","%v","%v",{}]`, messageId, ocppj.NotSupported, errorDescription)
	channel := NewMockWebSocket(wsId)

	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(errorJson), forwardWrittenMessage: false})
	coreListener := &MockCentralSystemCoreListener{}
	setupDefaultCentralSystemHandlers(suite, coreListener, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(errorJson), forwardWrittenMessage: true})
	resultChannel := make(chan struct{}, 1)
	suite.ocppjChargePoint.SetErrorHandler(func(err *ocpp.Error, details interface{}) {
		assert.Equal(t, messageId, err.MessageId)
		assert.Equal(t, ocppj.NotSupported, err.Code)
		assert.Equal(t, errorDescription, err.Description)
		assert.Equal(t, map[string]interface{}{}, details)
		resultChannel <- struct{}{}
	})
	// Start
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	require.Nil(t, err)
	// 1. Test sending an unsupported request, expecting an error
	err = suite.chargePoint.SendRequestAsync(request, func(confirmation ocpp.Response, err error) {
		t.Fail()
	})
	require.Error(t, err)
	assert.Equal(t, expectedError, err.Error())
	// 2. Test receiving an unsupported request on the other endpoint and receiving an error
	// Mark mocked request as pending, otherwise response will be ignored
	suite.ocppjChargePoint.RequestState.AddPendingRequest(messageId, request)
	err = suite.mockWsServer.MessageHandler(channel, []byte(requestJson))
	assert.Nil(t, err)
	_, ok := <-resultChannel
	assert.True(t, ok)
	// Stop the central system
	suite.centralSystem.Stop()
}

func testUnsupportedRequestFromCentralSystem(suite *OcppV16TestSuite, request ocpp.Request, requestJson string, messageId string) {
	t := suite.T()
	wsId := "test_id"
	wsUrl := "someUrl"
	expectedError := fmt.Sprintf("unsupported action %v on central system, cannot send request", request.GetFeatureName())
	errorDescription := fmt.Sprintf("unsupported action %v on charge point", request.GetFeatureName())
	errorJson := fmt.Sprintf(`[4,"%v","%v","%v",{}]`, messageId, ocppj.NotSupported, errorDescription)
	channel := NewMockWebSocket(wsId)

	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: false})
	coreListener := &MockChargePointCoreListener{}
	setupDefaultChargePointHandlers(suite, coreListener, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(errorJson), forwardWrittenMessage: true})
	resultChannel := make(chan struct{}, 1)
	suite.ocppjCentralSystem.SetErrorHandler(func(chargePoint ws.Channel, err *ocpp.Error, details interface{}) {
		assert.Equal(t, messageId, err.MessageId)
		assert.Equal(t, wsId, chargePoint.ID())
		assert.Equal(t, ocppj.NotSupported, err.Code)
		assert.Equal(t, errorDescription, err.Description)
		assert.Equal(t, map[string]interface{}{}, details)
		resultChannel <- struct{}{}
	})
	// Start
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	require.Nil(t, err)
	// 1. Test sending an unsupported request, expecting an error
	err = suite.centralSystem.SendRequestAsync(wsId, request, func(confirmation ocpp.Response, err error) {
		t.Fail()
	})
	require.Error(t, err)
	assert.Equal(t, expectedError, err.Error())
	// 2. Test receiving an unsupported request on the other endpoint and receiving an error
	// Mark mocked request as pending, otherwise response will be ignored
	suite.ocppjCentralSystem.RequestState.AddPendingRequest(wsId, messageId, request)
	err = suite.mockWsClient.MessageHandler([]byte(requestJson))
	assert.Nil(t, err)
	_, ok := <-resultChannel
	assert.True(t, ok)
	// Stop the central system
	suite.centralSystem.Stop()
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
	Confirmation  ocpp.Response
	ExpectedValid bool
}

// TODO: pass expected error value for improved validation and error message
func ExecuteGenericTestTable(t *testing.T, testTable []GenericTestEntry) {
	for _, testCase := range testTable {
		err := types.Validate.Struct(testCase.Element)
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
	ocppjChargePoint   *ocppj.Client
	ocppjCentralSystem *ocppj.Server
	mockWsServer       *MockWebsocketServer
	mockWsClient       *MockWebsocketClient
	chargePoint        ocpp16.ChargePoint
	centralSystem      ocpp16.CentralSystem
	messageIdGenerator TestRandomIdGenerator
	clientDispatcher   ocppj.ClientDispatcher
	serverDispatcher   ocppj.ServerDispatcher
}

type TestRandomIdGenerator struct {
	generator func() string
}

func (testGenerator *TestRandomIdGenerator) generateId() string {
	return testGenerator.generator()
}

var defaultMessageId = "1234"

func (suite *OcppV16TestSuite) SetupTest() {
	coreProfile := core.Profile
	localAuthListProfile := localauth.Profile
	firmwareProfile := firmware.Profile
	reservationProfile := reservation.Profile
	remoteTriggerProfile := remotetrigger.Profile
	smartChargingProfile := smartcharging.Profile
	certificatesProfile := certificates.Profile
	secureFirmwareUpdateProfile := securefirmware.Profile
	extendedTriggerMessageProfile := extendedtriggermessage.Profile
	securityProfile := security.Profile
	mockClient := MockWebsocketClient{}
	mockServer := MockWebsocketServer{}
	suite.mockWsClient = &mockClient
	suite.mockWsServer = &mockServer
	suite.clientDispatcher = ocppj.NewDefaultClientDispatcher(ocppj.NewFIFOClientQueue(queueCapacity))
	suite.serverDispatcher = ocppj.NewDefaultServerDispatcher(ocppj.NewFIFOQueueMap(queueCapacity))
	suite.ocppjChargePoint = ocppj.NewClient(
		"test_id",
		suite.mockWsClient,
		suite.clientDispatcher,
		nil,
		coreProfile,
		localAuthListProfile,
		firmwareProfile,
		reservationProfile,
		remoteTriggerProfile,
		smartChargingProfile,
		certificatesProfile,
		extendedTriggerMessageProfile,
		securityProfile,
		secureFirmwareUpdateProfile,
	)
	suite.ocppjCentralSystem = ocppj.NewServer(
		suite.mockWsServer,
		suite.serverDispatcher,
		nil,
		coreProfile,
		localAuthListProfile,
		firmwareProfile,
		reservationProfile,
		remoteTriggerProfile,
		smartChargingProfile,
		certificatesProfile,
		extendedTriggerMessageProfile,
		securityProfile,
		secureFirmwareUpdateProfile,
	)
	suite.chargePoint = ocpp16.NewChargePoint("test_id", suite.ocppjChargePoint, suite.mockWsClient)
	suite.centralSystem = ocpp16.NewCentralSystem(suite.ocppjCentralSystem, suite.mockWsServer)
	suite.messageIdGenerator = TestRandomIdGenerator{generator: func() string {
		return defaultMessageId
	}}
	ocppj.SetMessageIdGenerator(suite.messageIdGenerator.generateId)
	types.DateTimeFormat = time.RFC3339
}

func (suite *OcppV16TestSuite) TestIsConnected() {
	t := suite.T()
	// Simulate ws connected
	mockCall := suite.mockWsClient.On("IsConnected").Return(true)
	assert.True(t, suite.chargePoint.IsConnected())
	// Simulate ws disconnected
	mockCall.Return(false)
	assert.False(t, suite.chargePoint.IsConnected())
}

// TODO: implement generic protocol tests
func TestOcpp16Protocol(t *testing.T) {
	suite.Run(t, new(OcppV16TestSuite))
}
