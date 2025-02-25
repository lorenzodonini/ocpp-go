package ocpp2_test

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/authorization"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/availability"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/data"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/diagnostics"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/display"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/firmware"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/iso15118"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/meter"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/provisioning"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/remotecontrol"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/reservation"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/security"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/smartcharging"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/tariffcost"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/transactions"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/ocppj"
	"github.com/lorenzodonini/ocpp-go/ws"
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

func (websocketServer *MockWebsocketServer) SetCheckClientHandler(handler func(id string, r *http.Request) (string, bool)) {
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

type MockResponse struct {
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
	return reflect.TypeOf(MockResponse{})
}

func (r *MockRequest) GetFeatureName() string {
	return MockFeatureName
}

func (c *MockResponse) GetFeatureName() string {
	return MockFeatureName
}

func newMockRequest(value string) *MockRequest {
	return &MockRequest{MockValue: value}
}

func newMockConfirmation(value string) *MockResponse {
	return &MockResponse{MockValue: value}
}

// ---------------------- MOCK CSMS SECURITY HANDLER ----------------------

type MockCSMSSecurityHandler struct {
	mock.Mock
}

func (handler *MockCSMSSecurityHandler) OnSecurityEventNotification(chargingStationID string, request *security.SecurityEventNotificationRequest) (response *security.SecurityEventNotificationResponse, err error) {
	args := handler.MethodCalled("OnSecurityEventNotification", chargingStationID, request)
	response = args.Get(0).(*security.SecurityEventNotificationResponse)
	return response, args.Error(1)
}

func (handler *MockCSMSSecurityHandler) OnSignCertificate(chargingStationID string, request *security.SignCertificateRequest) (response *security.SignCertificateResponse, err error) {
	args := handler.MethodCalled("OnSignCertificate", chargingStationID, request)
	response = args.Get(0).(*security.SignCertificateResponse)
	return response, args.Error(1)
}

// ---------------------- MOCK CS SECURITY HANDLER ----------------------

type MockChargingStationSecurityHandler struct {
	mock.Mock
}

func (handler *MockChargingStationSecurityHandler) OnCertificateSigned(request *security.CertificateSignedRequest) (response *security.CertificateSignedResponse, err error) {
	args := handler.MethodCalled("OnCertificateSigned", request)
	conf := args.Get(0).(*security.CertificateSignedResponse)
	return conf, args.Error(1)
}

// ---------------------- MOCK CSMS PROVISIONING HANDLER ----------------------

type MockCSMSProvisioningHandler struct {
	mock.Mock
}

func (handler *MockCSMSProvisioningHandler) OnBootNotification(chargingStationId string, request *provisioning.BootNotificationRequest) (confirmation *provisioning.BootNotificationResponse, err error) {
	args := handler.MethodCalled("OnBootNotification", chargingStationId, request)
	conf := args.Get(0).(*provisioning.BootNotificationResponse)
	return conf, args.Error(1)
}

func (handler *MockCSMSProvisioningHandler) OnNotifyReport(chargingStationID string, request *provisioning.NotifyReportRequest) (confirmation *provisioning.NotifyReportResponse, err error) {
	args := handler.MethodCalled("OnNotifyReport", chargingStationID, request)
	conf := args.Get(0).(*provisioning.NotifyReportResponse)
	return conf, args.Error(1)
}

// ---------------------- MOCK CS PROVISIONING HANDLER ----------------------

type MockChargingStationProvisioningHandler struct {
	mock.Mock
}

func (handler *MockChargingStationProvisioningHandler) OnGetBaseReport(request *provisioning.GetBaseReportRequest) (confirmation *provisioning.GetBaseReportResponse, err error) {
	args := handler.MethodCalled("OnGetBaseReport", request)
	conf := args.Get(0).(*provisioning.GetBaseReportResponse)
	return conf, args.Error(1)
}

func (handler *MockChargingStationProvisioningHandler) OnGetReport(request *provisioning.GetReportRequest) (response *provisioning.GetReportResponse, err error) {
	args := handler.MethodCalled("OnGetReport", request)
	conf := args.Get(0).(*provisioning.GetReportResponse)
	return conf, args.Error(1)
}

func (handler *MockChargingStationProvisioningHandler) OnGetVariables(request *provisioning.GetVariablesRequest) (response *provisioning.GetVariablesResponse, err error) {
	args := handler.MethodCalled("OnGetVariables", request)
	conf := args.Get(0).(*provisioning.GetVariablesResponse)
	return conf, args.Error(1)
}

func (handler *MockChargingStationProvisioningHandler) OnReset(request *provisioning.ResetRequest) (response *provisioning.ResetResponse, err error) {
	args := handler.MethodCalled("OnReset", request)
	response = args.Get(0).(*provisioning.ResetResponse)
	return response, args.Error(1)
}

func (handler *MockChargingStationProvisioningHandler) OnSetNetworkProfile(request *provisioning.SetNetworkProfileRequest) (response *provisioning.SetNetworkProfileResponse, err error) {
	args := handler.MethodCalled("OnSetNetworkProfile", request)
	response = args.Get(0).(*provisioning.SetNetworkProfileResponse)
	return response, args.Error(1)
}

func (handler *MockChargingStationProvisioningHandler) OnSetVariables(request *provisioning.SetVariablesRequest) (response *provisioning.SetVariablesResponse, err error) {
	args := handler.MethodCalled("OnSetVariables", request)
	response = args.Get(0).(*provisioning.SetVariablesResponse)
	return response, args.Error(1)
}

// ---------------------- MOCK CSMS AUTHORIZATION HANDLER ----------------------

type MockCSMSAuthorizationHandler struct {
	mock.Mock
}

func (handler *MockCSMSAuthorizationHandler) OnAuthorize(chargingStationId string, request *authorization.AuthorizeRequest) (confirmation *authorization.AuthorizeResponse, err error) {
	args := handler.MethodCalled("OnAuthorize", chargingStationId, request)
	conf := args.Get(0).(*authorization.AuthorizeResponse)
	return conf, args.Error(1)
}

// ---------------------- MOCK CS AUTHORIZATION HANDLER ----------------------

type MockChargingStationAuthorizationHandler struct {
	mock.Mock
}

func (handler *MockChargingStationAuthorizationHandler) OnClearCache(request *authorization.ClearCacheRequest) (confirmation *authorization.ClearCacheResponse, err error) {
	args := handler.MethodCalled("OnClearCache", request)
	conf := args.Get(0).(*authorization.ClearCacheResponse)
	return conf, args.Error(1)
}

// ---------------------- MOCK CS RESERVATION HANDLER ----------------------

type MockChargingStationReservationHandler struct {
	mock.Mock
}

func (handler *MockChargingStationReservationHandler) OnCancelReservation(request *reservation.CancelReservationRequest) (confirmation *reservation.CancelReservationResponse, err error) {
	args := handler.MethodCalled("OnCancelReservation", request)
	conf := args.Get(0).(*reservation.CancelReservationResponse)
	return conf, args.Error(1)
}

func (handler *MockChargingStationReservationHandler) OnReserveNow(request *reservation.ReserveNowRequest) (resp *reservation.ReserveNowResponse, err error) {
	args := handler.MethodCalled("OnReserveNow", request)
	conf := args.Get(0).(*reservation.ReserveNowResponse)
	return conf, args.Error(1)
}

// ---------------------- MOCK CSMS RESERVATION HANDLER ----------------------

type MockCSMSReservationHandler struct {
	mock.Mock
}

func (handler *MockCSMSReservationHandler) OnReservationStatusUpdate(chargingStationID string, request *reservation.ReservationStatusUpdateRequest) (response *reservation.ReservationStatusUpdateResponse, err error) {
	args := handler.MethodCalled("OnReservationStatusUpdate", chargingStationID, request)
	resp := args.Get(0).(*reservation.ReservationStatusUpdateResponse)
	return resp, args.Error(1)
}

// ---------------------- MOCK CS AVAILABILITY HANDLER ----------------------

type MockChargingStationAvailabilityHandler struct {
	mock.Mock
}

func (handler *MockChargingStationAvailabilityHandler) OnChangeAvailability(request *availability.ChangeAvailabilityRequest) (confirmation *availability.ChangeAvailabilityResponse, err error) {
	args := handler.MethodCalled("OnChangeAvailability", request)
	conf := args.Get(0).(*availability.ChangeAvailabilityResponse)
	return conf, args.Error(1)
}

// ---------------------- MOCK CSMS AVAILABILITY HANDLER ----------------------

type MockCSMSAvailabilityHandler struct {
	mock.Mock
}

func (handler *MockCSMSAvailabilityHandler) OnHeartbeat(chargingStationID string, request *availability.HeartbeatRequest) (response *availability.HeartbeatResponse, err error) {
	args := handler.MethodCalled("OnHeartbeat", chargingStationID, request)
	resp := args.Get(0).(*availability.HeartbeatResponse)
	return resp, args.Error(1)
}

func (handler *MockCSMSAvailabilityHandler) OnStatusNotification(chargingStationID string, request *availability.StatusNotificationRequest) (response *availability.StatusNotificationResponse, err error) {
	args := handler.MethodCalled("OnStatusNotification", chargingStationID, request)
	response = args.Get(0).(*availability.StatusNotificationResponse)
	return response, args.Error(1)
}

// ---------------------- MOCK CS DATA HANDLER ----------------------

type MockChargingStationDataHandler struct {
	mock.Mock
}

func (handler *MockChargingStationDataHandler) OnDataTransfer(request *data.DataTransferRequest) (response *data.DataTransferResponse, err error) {
	args := handler.MethodCalled("OnDataTransfer", request)
	rawResp := args.Get(0)
	err = args.Error(1)
	if rawResp != nil {
		response = rawResp.(*data.DataTransferResponse)
	}
	return
}

// ---------------------- MOCK CSMS DATA HANDLER ----------------------

type MockCSMSDataHandler struct {
	mock.Mock
}

func (handler *MockCSMSDataHandler) OnDataTransfer(chargingStationID string, request *data.DataTransferRequest) (response *data.DataTransferResponse, err error) {
	args := handler.MethodCalled("OnDataTransfer", chargingStationID, request)
	rawResp := args.Get(0)
	err = args.Error(1)
	if rawResp != nil {
		response = rawResp.(*data.DataTransferResponse)
	}
	return
}

// ---------------------- MOCK CS DIAGNOSTICS HANDLER ----------------------

type MockChargingStationDiagnosticsHandler struct {
	mock.Mock
}

func (handler *MockChargingStationDiagnosticsHandler) OnClearVariableMonitoring(request *diagnostics.ClearVariableMonitoringRequest) (confirmation *diagnostics.ClearVariableMonitoringResponse, err error) {
	args := handler.MethodCalled("OnClearVariableMonitoring", request)
	conf := args.Get(0).(*diagnostics.ClearVariableMonitoringResponse)
	return conf, args.Error(1)
}

func (handler *MockChargingStationDiagnosticsHandler) OnCustomerInformation(request *diagnostics.CustomerInformationRequest) (confirmation *diagnostics.CustomerInformationResponse, err error) {
	args := handler.MethodCalled("OnCustomerInformation", request)
	conf := args.Get(0).(*diagnostics.CustomerInformationResponse)
	return conf, args.Error(1)
}

func (handler *MockChargingStationDiagnosticsHandler) OnGetLog(request *diagnostics.GetLogRequest) (confirmation *diagnostics.GetLogResponse, err error) {
	args := handler.MethodCalled("OnGetLog", request)
	conf := args.Get(0).(*diagnostics.GetLogResponse)
	return conf, args.Error(1)
}

func (handler *MockChargingStationDiagnosticsHandler) OnGetMonitoringReport(request *diagnostics.GetMonitoringReportRequest) (confirmation *diagnostics.GetMonitoringReportResponse, err error) {
	args := handler.MethodCalled("OnGetMonitoringReport", request)
	conf := args.Get(0).(*diagnostics.GetMonitoringReportResponse)
	return conf, args.Error(1)
}

func (handler *MockChargingStationDiagnosticsHandler) OnSetMonitoringBase(request *diagnostics.SetMonitoringBaseRequest) (response *diagnostics.SetMonitoringBaseResponse, err error) {
	args := handler.MethodCalled("OnSetMonitoringBase", request)
	response = args.Get(0).(*diagnostics.SetMonitoringBaseResponse)
	return response, args.Error(1)
}

func (handler *MockChargingStationDiagnosticsHandler) OnSetMonitoringLevel(request *diagnostics.SetMonitoringLevelRequest) (response *diagnostics.SetMonitoringLevelResponse, err error) {
	args := handler.MethodCalled("OnSetMonitoringLevel", request)
	response = args.Get(0).(*diagnostics.SetMonitoringLevelResponse)
	return response, args.Error(1)
}

func (handler *MockChargingStationDiagnosticsHandler) OnSetVariableMonitoring(request *diagnostics.SetVariableMonitoringRequest) (response *diagnostics.SetVariableMonitoringResponse, err error) {
	args := handler.MethodCalled("OnSetVariableMonitoring", request)
	response = args.Get(0).(*diagnostics.SetVariableMonitoringResponse)
	return response, args.Error(1)
}

// ---------------------- MOCK CSMS DIAGNOSTICS HANDLER ----------------------

type MockCSMSDiagnosticsHandler struct {
	mock.Mock
}

func (handler *MockCSMSDiagnosticsHandler) OnLogStatusNotification(chargingStationID string, request *diagnostics.LogStatusNotificationRequest) (response *diagnostics.LogStatusNotificationResponse, err error) {
	args := handler.MethodCalled("OnLogStatusNotification", chargingStationID, request)
	resp := args.Get(0).(*diagnostics.LogStatusNotificationResponse)
	return resp, args.Error(1)
}

func (handler *MockCSMSDiagnosticsHandler) OnNotifyCustomerInformation(chargingStationID string, request *diagnostics.NotifyCustomerInformationRequest) (response *diagnostics.NotifyCustomerInformationResponse, err error) {
	args := handler.MethodCalled("OnNotifyCustomerInformation", chargingStationID, request)
	resp := args.Get(0).(*diagnostics.NotifyCustomerInformationResponse)
	return resp, args.Error(1)
}

func (handler *MockCSMSDiagnosticsHandler) OnNotifyEvent(chargingStationID string, request *diagnostics.NotifyEventRequest) (response *diagnostics.NotifyEventResponse, err error) {
	args := handler.MethodCalled("OnNotifyEvent", chargingStationID, request)
	resp := args.Get(0).(*diagnostics.NotifyEventResponse)
	return resp, args.Error(1)
}

func (handler *MockCSMSDiagnosticsHandler) OnNotifyMonitoringReport(chargingStationID string, request *diagnostics.NotifyMonitoringReportRequest) (response *diagnostics.NotifyMonitoringReportResponse, err error) {
	args := handler.MethodCalled("OnNotifyMonitoringReport", chargingStationID, request)
	resp := args.Get(0).(*diagnostics.NotifyMonitoringReportResponse)
	return resp, args.Error(1)
}

// ---------------------- MOCK CS DISPLAY HANDLER ----------------------

type MockChargingStationDisplayHandler struct {
	mock.Mock
}

func (handler *MockChargingStationDisplayHandler) OnClearDisplay(request *display.ClearDisplayRequest) (confirmation *display.ClearDisplayResponse, err error) {
	args := handler.MethodCalled("OnClearDisplay", request)
	conf := args.Get(0).(*display.ClearDisplayResponse)
	return conf, args.Error(1)
}

func (handler *MockChargingStationDisplayHandler) OnGetDisplayMessages(request *display.GetDisplayMessagesRequest) (confirmation *display.GetDisplayMessagesResponse, err error) {
	args := handler.MethodCalled("OnGetDisplayMessages", request)
	conf := args.Get(0).(*display.GetDisplayMessagesResponse)
	return conf, args.Error(1)
}

func (handler *MockChargingStationDisplayHandler) OnSetDisplayMessage(request *display.SetDisplayMessageRequest) (response *display.SetDisplayMessageResponse, err error) {
	args := handler.MethodCalled("OnSetDisplayMessage", request)
	response = args.Get(0).(*display.SetDisplayMessageResponse)
	return response, args.Error(1)
}

// ---------------------- MOCK CSMS DISPLAY HANDLER ----------------------

type MockCSMSDisplayHandler struct {
	mock.Mock
}

func (handler *MockCSMSDisplayHandler) OnNotifyDisplayMessages(chargingStationID string, request *display.NotifyDisplayMessagesRequest) (response *display.NotifyDisplayMessagesResponse, err error) {
	args := handler.MethodCalled("OnNotifyDisplayMessages", chargingStationID, request)
	conf := args.Get(0).(*display.NotifyDisplayMessagesResponse)
	return conf, args.Error(1)
}

// ---------------------- MOCK CS FIRMWARE HANDLER ----------------------

type MockChargingStationFirmwareHandler struct {
	mock.Mock
}

func (handler *MockChargingStationFirmwareHandler) OnPublishFirmware(request *firmware.PublishFirmwareRequest) (response *firmware.PublishFirmwareResponse, err error) {
	args := handler.MethodCalled("OnPublishFirmware", request)
	resp := args.Get(0).(*firmware.PublishFirmwareResponse)
	return resp, args.Error(1)
}

func (handler *MockChargingStationFirmwareHandler) OnUnpublishFirmware(request *firmware.UnpublishFirmwareRequest) (response *firmware.UnpublishFirmwareResponse, err error) {
	args := handler.MethodCalled("OnUnpublishFirmware", request)
	response = args.Get(0).(*firmware.UnpublishFirmwareResponse)
	return response, args.Error(1)
}

func (handler *MockChargingStationFirmwareHandler) OnUpdateFirmware(request *firmware.UpdateFirmwareRequest) (response *firmware.UpdateFirmwareResponse, err error) {
	args := handler.MethodCalled("OnUpdateFirmware", request)
	response = args.Get(0).(*firmware.UpdateFirmwareResponse)
	return response, args.Error(1)
}

// ---------------------- MOCK CSMS FIRMWARE HANDLER ----------------------

type MockCSMSFirmwareHandler struct {
	mock.Mock
}

func (handler *MockCSMSFirmwareHandler) OnFirmwareStatusNotification(chargingStationID string, request *firmware.FirmwareStatusNotificationRequest) (response *firmware.FirmwareStatusNotificationResponse, err error) {
	args := handler.MethodCalled("OnFirmwareStatusNotification", chargingStationID, request)
	resp := args.Get(0).(*firmware.FirmwareStatusNotificationResponse)
	return resp, args.Error(1)
}

func (handler *MockCSMSFirmwareHandler) OnPublishFirmwareStatusNotification(chargingStationID string, request *firmware.PublishFirmwareStatusNotificationRequest) (response *firmware.PublishFirmwareStatusNotificationResponse, err error) {
	args := handler.MethodCalled("OnPublishFirmwareStatusNotification", chargingStationID, request)
	resp := args.Get(0).(*firmware.PublishFirmwareStatusNotificationResponse)
	return resp, args.Error(1)
}

// ---------------------- MOCK CS ISO15118 HANDLER ----------------------

type MockChargingStationIso15118Handler struct {
	mock.Mock
}

func (handler *MockChargingStationIso15118Handler) OnDeleteCertificate(request *iso15118.DeleteCertificateRequest) (response *iso15118.DeleteCertificateResponse, err error) {
	args := handler.MethodCalled("OnDeleteCertificate", request)
	resp := args.Get(0).(*iso15118.DeleteCertificateResponse)
	return resp, args.Error(1)
}

func (handler *MockChargingStationIso15118Handler) OnGetInstalledCertificateIds(request *iso15118.GetInstalledCertificateIdsRequest) (response *iso15118.GetInstalledCertificateIdsResponse, err error) {
	args := handler.MethodCalled("OnGetInstalledCertificateIds", request)
	resp := args.Get(0).(*iso15118.GetInstalledCertificateIdsResponse)
	return resp, args.Error(1)
}

func (handler *MockChargingStationIso15118Handler) OnInstallCertificate(request *iso15118.InstallCertificateRequest) (response *iso15118.InstallCertificateResponse, err error) {
	args := handler.MethodCalled("OnInstallCertificate", request)
	resp := args.Get(0).(*iso15118.InstallCertificateResponse)
	return resp, args.Error(1)
}

// ---------------------- MOCK CSMS ISO15118 HANDLER ----------------------

type MockCSMSIso15118Handler struct {
	mock.Mock
}

func (handler *MockCSMSIso15118Handler) OnGet15118EVCertificate(chargingStationID string, request *iso15118.Get15118EVCertificateRequest) (confirmation *iso15118.Get15118EVCertificateResponse, err error) {
	args := handler.MethodCalled("OnGet15118EVCertificate", chargingStationID, request)
	conf := args.Get(0).(*iso15118.Get15118EVCertificateResponse)
	return conf, args.Error(1)
}

func (handler *MockCSMSIso15118Handler) OnGetCertificateStatus(chargingStationID string, request *iso15118.GetCertificateStatusRequest) (confirmation *iso15118.GetCertificateStatusResponse, err error) {
	args := handler.MethodCalled("OnGetCertificateStatus", chargingStationID, request)
	conf := args.Get(0).(*iso15118.GetCertificateStatusResponse)
	return conf, args.Error(1)
}

// ---------------------- MOCK CS LOCAL AUTH HANDLER ----------------------

type MockChargingStationLocalAuthHandler struct {
	mock.Mock
}

func (handler *MockChargingStationLocalAuthHandler) OnGetLocalListVersion(request *localauth.GetLocalListVersionRequest) (confirmation *localauth.GetLocalListVersionResponse, err error) {
	args := handler.MethodCalled("OnGetLocalListVersion", request)
	conf := args.Get(0).(*localauth.GetLocalListVersionResponse)
	return conf, args.Error(1)
}

func (handler *MockChargingStationLocalAuthHandler) OnSendLocalList(request *localauth.SendLocalListRequest) (response *localauth.SendLocalListResponse, err error) {
	args := handler.MethodCalled("OnSendLocalList", request)
	response = args.Get(0).(*localauth.SendLocalListResponse)
	return response, args.Error(1)
}

// ---------------------- MOCK CSMS LOCAL AUTH HANDLER ----------------------

type MockCSMSLocalAuthHandler struct {
	mock.Mock
}

// ---------------------- MOCK CS METER HANDLER ----------------------

type MockChargingStationMeterHandler struct {
	mock.Mock
}

// ---------------------- MOCK CSMS METER HANDLER ----------------------

type MockCSMSMeterHandler struct {
	mock.Mock
}

func (handler *MockCSMSMeterHandler) OnMeterValues(chargingStationID string, request *meter.MeterValuesRequest) (response *meter.MeterValuesResponse, err error) {
	args := handler.MethodCalled("OnMeterValues", chargingStationID, request)
	r := args.Get(0).(*meter.MeterValuesResponse)
	return r, args.Error(1)
}

// ---------------------- MOCK CS REMOTE CONTROL HANDLER ----------------------

type MockChargingStationRemoteControlHandler struct {
	mock.Mock
}

func (handler *MockChargingStationRemoteControlHandler) OnRequestStartTransaction(request *remotecontrol.RequestStartTransactionRequest) (response *remotecontrol.RequestStartTransactionResponse, err error) {
	args := handler.MethodCalled("OnRequestStartTransaction", request)
	conf := args.Get(0).(*remotecontrol.RequestStartTransactionResponse)
	return conf, args.Error(1)
}

func (handler *MockChargingStationRemoteControlHandler) OnRequestStopTransaction(request *remotecontrol.RequestStopTransactionRequest) (response *remotecontrol.RequestStopTransactionResponse, err error) {
	args := handler.MethodCalled("OnRequestStopTransaction", request)
	conf := args.Get(0).(*remotecontrol.RequestStopTransactionResponse)
	return conf, args.Error(1)
}

func (handler *MockChargingStationRemoteControlHandler) OnTriggerMessage(request *remotecontrol.TriggerMessageRequest) (response *remotecontrol.TriggerMessageResponse, err error) {
	args := handler.MethodCalled("OnTriggerMessage", request)
	response = args.Get(0).(*remotecontrol.TriggerMessageResponse)
	return response, args.Error(1)
}

func (handler *MockChargingStationRemoteControlHandler) OnUnlockConnector(request *remotecontrol.UnlockConnectorRequest) (response *remotecontrol.UnlockConnectorResponse, err error) {
	args := handler.MethodCalled("OnUnlockConnector", request)
	response = args.Get(0).(*remotecontrol.UnlockConnectorResponse)
	return response, args.Error(1)
}

// ---------------------- MOCK CSMS REMOTE CONTROL HANDLER ----------------------

type MockCSMSRemoteControlHandler struct {
	mock.Mock
}

// ---------------------- MOCK CS SMART CHARGING HANDLER ----------------------

type MockChargingStationSmartChargingHandler struct {
	mock.Mock
}

func (handler *MockChargingStationSmartChargingHandler) OnClearChargingProfile(request *smartcharging.ClearChargingProfileRequest) (confirmation *smartcharging.ClearChargingProfileResponse, err error) {
	args := handler.MethodCalled("OnClearChargingProfile", request)
	conf := args.Get(0).(*smartcharging.ClearChargingProfileResponse)
	return conf, args.Error(1)
}

func (handler *MockChargingStationSmartChargingHandler) OnGetChargingProfiles(request *smartcharging.GetChargingProfilesRequest) (confirmation *smartcharging.GetChargingProfilesResponse, err error) {
	args := handler.MethodCalled("OnGetChargingProfiles", request)
	conf := args.Get(0).(*smartcharging.GetChargingProfilesResponse)
	return conf, args.Error(1)
}

func (handler *MockChargingStationSmartChargingHandler) OnGetCompositeSchedule(request *smartcharging.GetCompositeScheduleRequest) (confirmation *smartcharging.GetCompositeScheduleResponse, err error) {
	args := handler.MethodCalled("OnGetCompositeSchedule", request)
	conf := args.Get(0).(*smartcharging.GetCompositeScheduleResponse)
	return conf, args.Error(1)
}

func (handler *MockChargingStationSmartChargingHandler) OnSetChargingProfile(request *smartcharging.SetChargingProfileRequest) (response *smartcharging.SetChargingProfileResponse, err error) {
	args := handler.MethodCalled("OnSetChargingProfile", request)
	response = args.Get(0).(*smartcharging.SetChargingProfileResponse)
	return response, args.Error(1)
}

// ---------------------- MOCK CSMS SMART CHARGING HANDLER ----------------------

type MockCSMSSmartChargingHandler struct {
	mock.Mock
}

func (handler *MockCSMSSmartChargingHandler) OnClearedChargingLimit(chargingStationID string, request *smartcharging.ClearedChargingLimitRequest) (confirmation *smartcharging.ClearedChargingLimitResponse, err error) {
	args := handler.MethodCalled("OnClearedChargingLimit", chargingStationID, request)
	r := args.Get(0).(*smartcharging.ClearedChargingLimitResponse)
	return r, args.Error(1)
}

func (handler *MockCSMSSmartChargingHandler) OnNotifyChargingLimit(chargingStationID string, request *smartcharging.NotifyChargingLimitRequest) (response *smartcharging.NotifyChargingLimitResponse, err error) {
	args := handler.MethodCalled("OnNotifyChargingLimit", chargingStationID, request)
	r := args.Get(0).(*smartcharging.NotifyChargingLimitResponse)
	return r, args.Error(1)
}

func (handler *MockCSMSSmartChargingHandler) OnNotifyEVChargingNeeds(chargingStationID string, request *smartcharging.NotifyEVChargingNeedsRequest) (response *smartcharging.NotifyEVChargingNeedsResponse, err error) {
	args := handler.MethodCalled("OnNotifyEVChargingNeeds", chargingStationID, request)
	r := args.Get(0).(*smartcharging.NotifyEVChargingNeedsResponse)
	return r, args.Error(1)
}

func (handler *MockCSMSSmartChargingHandler) OnNotifyEVChargingSchedule(chargingStationID string, request *smartcharging.NotifyEVChargingScheduleRequest) (response *smartcharging.NotifyEVChargingScheduleResponse, err error) {
	args := handler.MethodCalled("OnNotifyEVChargingSchedule", chargingStationID, request)
	r := args.Get(0).(*smartcharging.NotifyEVChargingScheduleResponse)
	return r, args.Error(1)
}

func (handler *MockCSMSSmartChargingHandler) OnReportChargingProfiles(chargingStationID string, request *smartcharging.ReportChargingProfilesRequest) (reponse *smartcharging.ReportChargingProfilesResponse, err error) {
	args := handler.MethodCalled("OnReportChargingProfiles", chargingStationID, request)
	r := args.Get(0).(*smartcharging.ReportChargingProfilesResponse)
	return r, args.Error(1)
}

// ---------------------- MOCK CS TARIFF COST HANDLER ----------------------

type MockChargingStationTariffCostHandler struct {
	mock.Mock
}

func (handler *MockChargingStationTariffCostHandler) OnCostUpdated(request *tariffcost.CostUpdatedRequest) (confirmation *tariffcost.CostUpdatedResponse, err error) {
	args := handler.MethodCalled("OnCostUpdated", request)
	conf := args.Get(0).(*tariffcost.CostUpdatedResponse)
	return conf, args.Error(1)
}

// ---------------------- MOCK CSMS TARIFF COST HANDLER ----------------------

type MockCSMSTariffCostHandler struct {
	mock.Mock
}

// ---------------------- MOCK CS TRANSACTIONS HANDLER ----------------------

type MockChargingStationTransactionHandler struct {
	mock.Mock
}

func (handler *MockChargingStationTransactionHandler) OnGetTransactionStatus(request *transactions.GetTransactionStatusRequest) (response *transactions.GetTransactionStatusResponse, err error) {
	args := handler.MethodCalled("OnGetTransactionStatus", request)
	conf := args.Get(0).(*transactions.GetTransactionStatusResponse)
	return conf, args.Error(1)
}

// ---------------------- MOCK CSMS TRANSACTIONS HANDLER ----------------------

type MockCSMSTransactionsHandler struct {
	mock.Mock
}

func (handler *MockCSMSTransactionsHandler) OnTransactionEvent(chargingStationID string, request *transactions.TransactionEventRequest) (response *transactions.TransactionEventResponse, err error) {
	args := handler.MethodCalled("OnTransactionEvent", chargingStationID, request)
	response = args.Get(0).(*transactions.TransactionEventResponse)
	return response, args.Error(1)
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

type expectedCSMSOptions struct {
	clientId              string
	rawWrittenMessage     []byte
	startReturnArgument   interface{}
	writeReturnArgument   interface{}
	forwardWrittenMessage bool
}

type expectedChargingStationOptions struct {
	serverUrl             string
	clientId              string
	createChannelOnStart  bool
	channel               ws.Channel
	rawWrittenMessage     []byte
	startReturnArgument   interface{}
	writeReturnArgument   interface{}
	forwardWrittenMessage bool
}

func setupDefaultCSMSHandlers(suite *OcppV2TestSuite, options expectedCSMSOptions, handlers ...interface{}) {
	t := suite.T()
	for _, h := range handlers {
		switch h := h.(type) {
		case *MockCSMSAuthorizationHandler:
			suite.csms.SetAuthorizationHandler(h)
		case *MockCSMSAvailabilityHandler:
			suite.csms.SetAvailabilityHandler(h)
		case *MockCSMSDataHandler:
			suite.csms.SetDataHandler(h)
		case *MockCSMSDiagnosticsHandler:
			suite.csms.SetDiagnosticsHandler(h)
		case *MockCSMSDisplayHandler:
			suite.csms.SetDisplayHandler(h)
		case *MockCSMSFirmwareHandler:
			suite.csms.SetFirmwareHandler(h)
		case *MockCSMSIso15118Handler:
			suite.csms.SetISO15118Handler(h)
		case *MockCSMSLocalAuthHandler:
			suite.csms.SetLocalAuthListHandler(h)
		case *MockCSMSMeterHandler:
			suite.csms.SetMeterHandler(h)
		case *MockCSMSProvisioningHandler:
			suite.csms.SetProvisioningHandler(h)
		case *MockCSMSRemoteControlHandler:
			suite.csms.SetRemoteControlHandler(h)
		case *MockCSMSReservationHandler:
			suite.csms.SetReservationHandler(h)
		case *MockCSMSSecurityHandler:
			suite.csms.SetSecurityHandler(h)
		case *MockCSMSSmartChargingHandler:
			suite.csms.SetSmartChargingHandler(h)
		case *MockCSMSTariffCostHandler:
			suite.csms.SetTariffCostHandler(h)
		case *MockCSMSTransactionsHandler:
			suite.csms.SetTransactionsHandler(h)
		}
	}
	suite.csms.SetNewChargingStationHandler(func(chargingStation ocpp2.ChargingStationConnection) {
		assert.Equal(t, options.clientId, chargingStation.ID())
	})
	suite.mockWsServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(options.startReturnArgument)
	suite.mockWsServer.On("Stop").Return()
	suite.mockWsServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(options.writeReturnArgument).Run(func(args mock.Arguments) {
		clientId := args.String(0)
		data := args.Get(1)
		bytes := data.([]byte)
		assert.Equal(t, options.clientId, clientId)
		if options.rawWrittenMessage != nil {
			assert.NotNil(t, bytes)
			assert.Equal(t, string(options.rawWrittenMessage), string(bytes))
		}
		if options.forwardWrittenMessage {
			// Notify client of incoming response
			err := suite.mockWsClient.MessageHandler(bytes)
			assert.Nil(t, err)
		}
	})
}

func setupDefaultChargingStationHandlers(suite *OcppV2TestSuite, options expectedChargingStationOptions, handlers ...interface{}) {
	t := suite.T()
	for _, h := range handlers {
		switch h := h.(type) {
		case *MockChargingStationAuthorizationHandler:
			suite.chargingStation.SetAuthorizationHandler(h)
		case *MockChargingStationAvailabilityHandler:
			suite.chargingStation.SetAvailabilityHandler(h)
		case *MockChargingStationDataHandler:
			suite.chargingStation.SetDataHandler(h)
		case *MockChargingStationDiagnosticsHandler:
			suite.chargingStation.SetDiagnosticsHandler(h)
		case *MockChargingStationDisplayHandler:
			suite.chargingStation.SetDisplayHandler(h)
		case *MockChargingStationFirmwareHandler:
			suite.chargingStation.SetFirmwareHandler(h)
		case *MockChargingStationIso15118Handler:
			suite.chargingStation.SetISO15118Handler(h)
		case *MockChargingStationLocalAuthHandler:
			suite.chargingStation.SetLocalAuthListHandler(h)
		case *MockChargingStationMeterHandler:
			suite.chargingStation.SetMeterHandler(h)
		case *MockChargingStationProvisioningHandler:
			suite.chargingStation.SetProvisioningHandler(h)
		case *MockChargingStationRemoteControlHandler:
			suite.chargingStation.SetRemoteControlHandler(h)
		case *MockChargingStationReservationHandler:
			suite.chargingStation.SetReservationHandler(h)
		case *MockChargingStationSecurityHandler:
			suite.chargingStation.SetSecurityHandler(h)
		case *MockChargingStationSmartChargingHandler:
			suite.chargingStation.SetSmartChargingHandler(h)
		case *MockChargingStationTariffCostHandler:
			suite.chargingStation.SetTariffCostHandler(h)
		case *MockChargingStationTransactionHandler:
			suite.chargingStation.SetTransactionsHandler(h)
		}
	}
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
			assert.Equal(t, string(options.rawWrittenMessage), string(bytes))
		}
		// Notify server of incoming request
		if options.forwardWrittenMessage {
			err := suite.mockWsServer.MessageHandler(options.channel, bytes)
			assert.Nil(t, err)
		}
	})
}

func assertDateTimeEquality(t *testing.T, expected *types.DateTime, actual *types.DateTime) {
	assert.Equal(t, expected.FormatTimestamp(), actual.FormatTimestamp())
}

func testUnsupportedRequestFromChargingStation(suite *OcppV2TestSuite, request ocpp.Request, requestJson string, messageId string, handlers ...interface{}) {
	t := suite.T()
	wsId := "test_id"
	wsUrl := "someUrl"
	expectedError := fmt.Sprintf("unsupported action %v on charging station, cannot send request", request.GetFeatureName())
	errorDescription := fmt.Sprintf("unsupported action %v on CSMS", request.GetFeatureName())
	errorJson := fmt.Sprintf(`[4,"%v","%v","%v",{}]`, messageId, ocppj.NotSupported, errorDescription)
	channel := NewMockWebSocket(wsId)

	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(errorJson), forwardWrittenMessage: false})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(errorJson), forwardWrittenMessage: true}, handlers...)
	resultChannel := make(chan bool, 1)
	suite.ocppjClient.SetErrorHandler(func(err *ocpp.Error, details interface{}) {
		assert.Equal(t, messageId, err.MessageId)
		assert.Equal(t, ocppj.NotSupported, err.Code)
		assert.Equal(t, errorDescription, err.Description)
		assert.Equal(t, map[string]interface{}{}, details)
		resultChannel <- true
	})
	// Start
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	// 1. Test sending an unsupported request, expecting an error
	err = suite.chargingStation.SendRequestAsync(request, func(confirmation ocpp.Response, err error) {
		t.Fail()
	})
	require.Error(t, err)
	assert.Equal(t, expectedError, err.Error())
	// 2. Test receiving an unsupported request on the other endpoint and receiving an error
	// Mark mocked request as pending, otherwise response will be ignored
	suite.ocppjClient.RequestState.AddPendingRequest(messageId, request)
	err = suite.mockWsServer.MessageHandler(channel, []byte(requestJson))
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
	// Stop the CSMS
	suite.csms.Stop()
}

func testUnsupportedRequestFromCentralSystem(suite *OcppV2TestSuite, request ocpp.Request, requestJson string, messageId string, handlers ...interface{}) {
	t := suite.T()
	wsId := "test_id"
	wsUrl := "someUrl"
	expectedError := fmt.Sprintf("unsupported action %v on CSMS, cannot send request", request.GetFeatureName())
	errorDescription := fmt.Sprintf("unsupported action %v on charging station", request.GetFeatureName())
	errorJson := fmt.Sprintf(`[4,"%v","%v","%v",{}]`, messageId, ocppj.NotSupported, errorDescription)
	channel := NewMockWebSocket(wsId)

	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: false})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(errorJson), forwardWrittenMessage: true}, handlers...)
	resultChannel := make(chan struct{}, 1)
	suite.ocppjServer.SetErrorHandler(func(channel ws.Channel, err *ocpp.Error, details interface{}) {
		assert.Equal(t, messageId, err.MessageId)
		assert.Equal(t, wsId, channel.ID())
		assert.Equal(t, ocppj.NotSupported, err.Code)
		assert.Equal(t, errorDescription, err.Description)
		assert.Equal(t, map[string]interface{}{}, details)
		resultChannel <- struct{}{}
	})
	// Start
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	// 1. Test sending an unsupported request, expecting an error
	err = suite.csms.SendRequestAsync(wsId, request, func(response ocpp.Response, err error) {
		t.Fail()
	})
	require.Error(t, err)
	assert.Equal(t, expectedError, err.Error())
	// 2. Test receiving an unsupported request on the other endpoint and receiving an error
	// Mark mocked request as pending, otherwise response will be ignored
	suite.ocppjServer.RequestState.AddPendingRequest(wsId, messageId, request)
	// Run response test
	err = suite.mockWsClient.MessageHandler([]byte(requestJson))
	assert.Nil(t, err)
	_, ok := <-resultChannel
	assert.True(t, ok)
	// Stop the CSMS
	suite.csms.Stop()
}

type GenericTestEntry struct {
	Element       interface{}
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

type OcppV2TestSuite struct {
	suite.Suite
	ocppjClient        *ocppj.Client
	ocppjServer        *ocppj.Server
	mockWsServer       *MockWebsocketServer
	mockWsClient       *MockWebsocketClient
	chargingStation    ocpp2.ChargingStation
	csms               ocpp2.CSMS
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

func (suite *OcppV2TestSuite) SetupTest() {
	securityProfile := security.Profile
	provisioningProfile := provisioning.Profile
	authProfile := authorization.Profile
	availabilityProfile := availability.Profile
	reservationProfile := reservation.Profile
	diagnosticsProfile := diagnostics.Profile
	dataProfile := data.Profile
	displayProfile := display.Profile
	firmwareProfile := firmware.Profile
	isoProfile := iso15118.Profile
	localAuthProfile := localauth.Profile
	meterProfile := meter.Profile
	remoteProfile := remotecontrol.Profile
	smartChargingProfile := smartcharging.Profile
	tariffProfile := tariffcost.Profile
	transactionsProfile := transactions.Profile
	mockClient := MockWebsocketClient{}
	mockServer := MockWebsocketServer{}
	suite.mockWsClient = &mockClient
	suite.mockWsServer = &mockServer
	suite.clientDispatcher = ocppj.NewDefaultClientDispatcher(ocppj.NewFIFOClientQueue(queueCapacity))
	suite.serverDispatcher = ocppj.NewDefaultServerDispatcher(ocppj.NewFIFOQueueMap(queueCapacity))
	suite.ocppjClient = ocppj.NewClient("test_id", suite.mockWsClient, suite.clientDispatcher, nil, securityProfile, provisioningProfile, authProfile, availabilityProfile, reservationProfile, diagnosticsProfile, dataProfile, displayProfile, firmwareProfile, isoProfile, localAuthProfile, meterProfile, remoteProfile, smartChargingProfile, tariffProfile, transactionsProfile)
	suite.ocppjServer = ocppj.NewServer(suite.mockWsServer, suite.serverDispatcher, nil, securityProfile, provisioningProfile, authProfile, availabilityProfile, reservationProfile, diagnosticsProfile, dataProfile, displayProfile, firmwareProfile, isoProfile, localAuthProfile, meterProfile, remoteProfile, smartChargingProfile, tariffProfile, transactionsProfile)
	suite.chargingStation = ocpp2.NewChargingStation("test_id", suite.ocppjClient, suite.mockWsClient)
	suite.csms = ocpp2.NewCSMS(suite.ocppjServer, suite.mockWsServer)
	suite.messageIdGenerator = TestRandomIdGenerator{generator: func() string {
		return defaultMessageId
	}}
	ocppj.SetMessageIdGenerator(suite.messageIdGenerator.generateId)
	types.DateTimeFormat = time.RFC3339
}

func (suite *OcppV2TestSuite) TestIsConnected() {
	t := suite.T()
	// Simulate ws connected
	mockCall := suite.mockWsClient.On("IsConnected").Return(true)
	assert.True(t, suite.chargingStation.IsConnected())
	// Simulate ws disconnected
	mockCall.Return(false)
	assert.False(t, suite.chargingStation.IsConnected())
}

//TODO: implement generic protocol tests

func TestOcpp2Protocol(t *testing.T) {
	logrus.SetLevel(logrus.PanicLevel)
	suite.Run(t, new(OcppV2TestSuite))
}
