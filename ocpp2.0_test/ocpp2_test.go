package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	ocpp2 "github.com/lorenzodonini/ocpp-go/ocpp2.0"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/authorization"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/availability"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/data"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/diagnostics"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/display"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/firmware"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/iso15118"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/meter"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/provisioning"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/remotecontrol"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/reservation"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/security"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/smartcharging"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/tariffcost"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/transactions"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
	"github.com/lorenzodonini/ocpp-go/ocppj"
	"github.com/lorenzodonini/ocpp-go/ws"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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

// ---------------------- MOCK CSMS SECURITY HANDLER ----------------------
type MockCSMSSecurityHandler struct {
	mock.Mock
}

// ---------------------- MOCK CS SECURITY HANDLER ----------------------
type MockChargingStationSecurityHandler struct {
	mock.Mock
}

func (handler MockChargingStationSecurityHandler) OnCertificateSigned(request *security.CertificateSignedRequest) (response *security.CertificateSignedConfirmation, err error) {
	args := handler.MethodCalled("OnCertificateSigned", request)
	conf := args.Get(0).(*security.CertificateSignedConfirmation)
	return conf, args.Error(1)
}

// ---------------------- MOCK CSMS PROVISIONING HANDLER ----------------------
type MockCSMSProvisioningHandler struct {
	mock.Mock
}

func (handler MockCSMSProvisioningHandler) OnBootNotification(chargePointId string, request *provisioning.BootNotificationRequest) (confirmation *provisioning.BootNotificationConfirmation, err error) {
	args := handler.MethodCalled("OnBootNotification", chargePointId, request)
	conf := args.Get(0).(*provisioning.BootNotificationConfirmation)
	return conf, args.Error(1)
}

// ---------------------- MOCK CS PROVISIONING HANDLER ----------------------
type MockChargingStationProvisioningHandler struct {
	mock.Mock
}

func (handler MockChargingStationProvisioningHandler) OnGetBaseReport(request *provisioning.GetBaseReportRequest) (confirmation *provisioning.GetBaseReportConfirmation, err error) {
	args := handler.MethodCalled("OnGetBaseReport", request)
	conf := args.Get(0).(*provisioning.GetBaseReportConfirmation)
	return conf, args.Error(1)
}

// ---------------------- MOCK CSMS AUTHORIZATION HANDLER ----------------------
type MockCSMSAuthorizationHandler struct {
	mock.Mock
}

func (handler MockCSMSAuthorizationHandler) OnAuthorize(chargePointId string, request *authorization.AuthorizeRequest) (confirmation *authorization.AuthorizeConfirmation, err error) {
	args := handler.MethodCalled("OnAuthorize", chargePointId, request)
	conf := args.Get(0).(*authorization.AuthorizeConfirmation)
	return conf, args.Error(1)
}

// ---------------------- MOCK CS AUTHORIZATION HANDLER ----------------------
type MockChargingStationAuthorizationHandler struct {
	mock.Mock
}

func (handler MockChargingStationAuthorizationHandler) OnClearCache(request *authorization.ClearCacheRequest) (confirmation *authorization.ClearCacheConfirmation, err error) {
	args := handler.MethodCalled("OnClearCache", request)
	conf := args.Get(0).(*authorization.ClearCacheConfirmation)
	return conf, args.Error(1)
}

// ---------------------- MOCK CS RESERVATION HANDLER ----------------------
type MockChargingStationReservationHandler struct {
	mock.Mock
}

func (handler MockChargingStationReservationHandler) OnCancelReservation(request *reservation.CancelReservationRequest) (confirmation *reservation.CancelReservationConfirmation, err error) {
	args := handler.MethodCalled("OnCancelReservation", request)
	conf := args.Get(0).(*reservation.CancelReservationConfirmation)
	return conf, args.Error(1)
}

// ---------------------- MOCK CSMS RESERVATION HANDLER ----------------------
type MockCSMSReservationHandler struct {
	mock.Mock
}

// ---------------------- MOCK CS AVAILABILITY HANDLER ----------------------
type MockChargingStationAvailabilityHandler struct {
	mock.Mock
}

func (handler MockChargingStationAvailabilityHandler) OnChangeAvailability(request *availability.ChangeAvailabilityRequest) (confirmation *availability.ChangeAvailabilityConfirmation, err error) {
	args := handler.MethodCalled("OnChangeAvailability", request)
	conf := args.Get(0).(*availability.ChangeAvailabilityConfirmation)
	return conf, args.Error(1)
}

// ---------------------- MOCK CSMS AVAILABILITY HANDLER ----------------------
type MockCSMSAvailabilityHandler struct {
	mock.Mock
}

// ---------------------- MOCK CS DATA HANDLER ----------------------
type MockChargingStationDataHandler struct {
	mock.Mock
}

func (handler MockChargingStationDataHandler) OnDataTransfer(request *data.DataTransferRequest) (confirmation *data.DataTransferConfirmation, err error) {
	args := handler.MethodCalled("OnDataTransfer", request)
	conf := args.Get(0).(*data.DataTransferConfirmation)
	return conf, args.Error(1)
}

// ---------------------- MOCK CSMS DATA HANDLER ----------------------
type MockCSMSDataHandler struct {
	mock.Mock
}

func (handler MockCSMSDataHandler) OnDataTransfer(chargingStationID string, request *data.DataTransferRequest) (confirmation *data.DataTransferConfirmation, err error) {
	args := handler.MethodCalled("OnDataTransfer", chargingStationID, request)
	conf := args.Get(0).(*data.DataTransferConfirmation)
	return conf, args.Error(1)
}

// ---------------------- MOCK CS DIAGNOSTICS HANDLER ----------------------
type MockChargingStationDiagnosticsHandler struct {
	mock.Mock
}

func (handler MockChargingStationDiagnosticsHandler) OnClearVariableMonitoring(request *diagnostics.ClearVariableMonitoringRequest) (confirmation *diagnostics.ClearVariableMonitoringConfirmation, err error) {
	args := handler.MethodCalled("OnClearVariableMonitoring", request)
	conf := args.Get(0).(*diagnostics.ClearVariableMonitoringConfirmation)
	return conf, args.Error(1)
}

func (handler MockChargingStationDiagnosticsHandler) OnCustomerInformation(request *diagnostics.CustomerInformationRequest) (confirmation *diagnostics.CustomerInformationConfirmation, err error) {
	args := handler.MethodCalled("OnCustomerInformation", request)
	conf := args.Get(0).(*diagnostics.CustomerInformationConfirmation)
	return conf, args.Error(1)
}

// ---------------------- MOCK CSMS DIAGNOSTICS HANDLER ----------------------
type MockCSMSDiagnosticsHandler struct {
	mock.Mock
}

// ---------------------- MOCK CS DISPLAY HANDLER ----------------------
type MockChargingStationDisplayHandler struct {
	mock.Mock
}

func (handler MockChargingStationDisplayHandler) OnClearDisplay(request *display.ClearDisplayRequest) (confirmation *display.ClearDisplayConfirmation, err error) {
	args := handler.MethodCalled("OnClearDisplay", request)
	conf := args.Get(0).(*display.ClearDisplayConfirmation)
	return conf, args.Error(1)
}

// ---------------------- MOCK CSMS DISPLAY HANDLER ----------------------
type MockCSMSDisplayHandler struct {
	mock.Mock
}

// ---------------------- MOCK CS FIRMWARE HANDLER ----------------------
type MockChargingStationFirmwareHandler struct {
	mock.Mock
}

// ---------------------- MOCK CSMS FIRMWARE HANDLER ----------------------
type MockCSMSFirmwareHandler struct {
	mock.Mock
}

func (handler MockCSMSFirmwareHandler) OnFirmwareStatusNotification(chargingStationID string, request *firmware.FirmwareStatusNotificationRequest) (confirmation *firmware.FirmwareStatusNotificationConfirmation, err error) {
	args := handler.MethodCalled("OnFirmwareStatusNotification", chargingStationID, request)
	conf := args.Get(0).(*firmware.FirmwareStatusNotificationConfirmation)
	return conf, args.Error(1)
}

// ---------------------- MOCK CS ISO15118 HANDLER ----------------------
type MockChargingStationIso15118Handler struct {
	mock.Mock
}

func (handler MockChargingStationIso15118Handler) OnDeleteCertificate(request *iso15118.DeleteCertificateRequest) (confirmation *iso15118.DeleteCertificateConfirmation, err error) {
	args := handler.MethodCalled("OnDeleteCertificate", request)
	conf := args.Get(0).(*iso15118.DeleteCertificateConfirmation)
	return conf, args.Error(1)
}

// ---------------------- MOCK CSMS ISO15118 HANDLER ----------------------
type MockCSMSIso15118Handler struct {
	mock.Mock
}

// ---------------------- MOCK CS LOCAL AUTH HANDLER ----------------------
type MockChargingStationLocalAuthHandler struct {
	mock.Mock
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

// ---------------------- MOCK CS REMOTE CONTROL HANDLER ----------------------
type MockChargingStationRemoteControlHandler struct {
	mock.Mock
}

// ---------------------- MOCK CSMS REMOTE CONTROL HANDLER ----------------------
type MockCSMSRemoteControlHandler struct {
	mock.Mock
}

// ---------------------- MOCK CS SMART CHARGING HANDLER ----------------------
type MockChargingStationSmartChargingHandler struct {
	mock.Mock
}

func (handler MockChargingStationSmartChargingHandler) OnClearChargingProfile(request *smartcharging.ClearChargingProfileRequest) (confirmation *smartcharging.ClearChargingProfileConfirmation, err error) {
	args := handler.MethodCalled("OnClearChargingProfile", request)
	conf := args.Get(0).(*smartcharging.ClearChargingProfileConfirmation)
	return conf, args.Error(1)
}

// ---------------------- MOCK CSMS SMART CHARGING HANDLER ----------------------
type MockCSMSSmartChargingHandler struct {
	mock.Mock
}

func (handler MockCSMSSmartChargingHandler) OnClearedChargingLimit(chargingStationID string, request *smartcharging.ClearedChargingLimitRequest) (confirmation *smartcharging.ClearedChargingLimitConfirmation, err error) {
	args := handler.MethodCalled("OnClearedChargingLimit", chargingStationID, request)
	conf := args.Get(0).(*smartcharging.ClearedChargingLimitConfirmation)
	return conf, args.Error(1)
}

// ---------------------- MOCK CS TARIFF COST HANDLER ----------------------
type MockChargingStationTariffCostHandler struct {
	mock.Mock
}

func (handler MockChargingStationTariffCostHandler) OnCostUpdated(request *tariffcost.CostUpdatedRequest) (confirmation *tariffcost.CostUpdatedConfirmation, err error) {
	args := handler.MethodCalled("OnCostUpdated", request)
	conf := args.Get(0).(*tariffcost.CostUpdatedConfirmation)
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

// ---------------------- MOCK CSMS TRANSACTIONS HANDLER ----------------------
type MockCSMSTransactionsHandler struct {
	mock.Mock
}


// ---------------------- MOCK CSMS CORE LISTENER ----------------------
type MockCSMSHandler struct {
	mock.Mock
}

func (coreListener MockCSMSHandler) OnClearedChargingLimit(chargePointId string, request *smartcharging.ClearedChargingLimitRequest) (confirmation *smartcharging.ClearedChargingLimitConfirmation, err error) {
	args := coreListener.MethodCalled("OnClearedChargingLimit", chargePointId, request)
	conf := args.Get(0).(*smartcharging.ClearedChargingLimitConfirmation)
	return conf, args.Error(1)
}

func (coreListener MockCSMSHandler) OnDataTransfer(chargePointId string, request *data.DataTransferRequest) (confirmation *data.DataTransferConfirmation, err error) {
	args := coreListener.MethodCalled("OnDataTransfer", chargePointId, request)
	conf := args.Get(0).(*data.DataTransferConfirmation)
	return conf, args.Error(1)
}

func (coreListener MockCSMSHandler) OnFirmwareStatusNotification(chargePointId string, request *firmware.FirmwareStatusNotificationRequest) (confirmation *firmware.FirmwareStatusNotificationConfirmation, err error) {
	args := coreListener.MethodCalled("OnFirmwareStatusNotification", chargePointId, request)
	conf := args.Get(0).(*firmware.FirmwareStatusNotificationConfirmation)
	return conf, args.Error(1)
}

func (coreListener MockCSMSHandler) OnGet15118EVCertificate(chargePointId string, request *ocpp2.Get15118EVCertificateRequest) (confirmation *ocpp2.Get15118EVCertificateConfirmation, err error) {
	args := coreListener.MethodCalled("OnGet15118EVCertificate", chargePointId, request)
	conf := args.Get(0).(*ocpp2.Get15118EVCertificateConfirmation)
	return conf, args.Error(1)
}

func (coreListener MockCSMSHandler) OnGetCertificateStatus(chargePointId string, request *ocpp2.GetCertificateStatusRequest) (confirmation *ocpp2.GetCertificateStatusConfirmation, err error) {
	args := coreListener.MethodCalled("OnGetCertificateStatus", chargePointId, request)
	conf := args.Get(0).(*ocpp2.GetCertificateStatusConfirmation)
	return conf, args.Error(1)
}

//func (coreListener MockCSMSHandler) OnDataTransfer(chargePointId string, request *ocpp2.DataTransferRequest) (confirmation *ocpp2.DataTransferConfirmation, err error) {
//	args := coreListener.MethodCalled("OnDataTransfer", chargePointId, request)
//	conf := args.Get(0).(*ocpp2.DataTransferConfirmation)
//	return conf, args.Error(1)
//}
//
//func (coreListener MockCSMSHandler) OnHeartbeat(chargePointId string, request *ocpp2.HeartbeatRequest) (confirmation *ocpp2.HeartbeatConfirmation, err error) {
//	args := coreListener.MethodCalled("OnHeartbeat", chargePointId, request)
//	conf := args.Get(0).(*ocpp2.HeartbeatConfirmation)
//	return conf, args.Error(1)
//}
//
//func (coreListener MockCSMSHandler) OnMeterValues(chargePointId string, request *ocpp2.MeterValuesRequest) (confirmation *ocpp2.MeterValuesConfirmation, err error) {
//	args := coreListener.MethodCalled("OnMeterValues", chargePointId, request)
//	conf := args.Get(0).(*ocpp2.MeterValuesConfirmation)
//	return conf, args.Error(1)
//}
//
//func (coreListener MockCSMSHandler) OnStartTransaction(chargePointId string, request *ocpp2.StartTransactionRequest) (confirmation *ocpp2.StartTransactionConfirmation, err error) {
//	args := coreListener.MethodCalled("OnStartTransaction", chargePointId, request)
//	conf := args.Get(0).(*ocpp2.StartTransactionConfirmation)
//	return conf, args.Error(1)
//}
//
//func (coreListener MockCSMSHandler) OnStatusNotification(chargePointId string, request *ocpp2.StatusNotificationRequest) (confirmation *ocpp2.StatusNotificationConfirmation, err error) {
//	args := coreListener.MethodCalled("OnStatusNotification", chargePointId, request)
//	conf := args.Get(0).(*ocpp2.StatusNotificationConfirmation)
//	return conf, args.Error(1)
//}
//
//func (coreListener MockCSMSHandler) OnStopTransaction(chargePointId string, request *ocpp2.StopTransactionRequest) (confirmation *ocpp2.StopTransactionConfirmation, err error) {
//	args := coreListener.MethodCalled("OnStopTransaction", chargePointId, request)
//	conf := args.Get(0).(*ocpp2.StopTransactionConfirmation)
//	return conf, args.Error(1)
//}

// ---------------------- MOCK CP CORE LISTENER ----------------------
type MockChargePointCoreListener struct {
	mock.Mock
}

func (coreListener MockChargePointCoreListener) OnCancelReservation(request *reservation.CancelReservationRequest) (confirmation *reservation.CancelReservationConfirmation, err error) {
	args := coreListener.MethodCalled("OnCancelReservation", request)
	conf := args.Get(0).(*reservation.CancelReservationConfirmation)
	return conf, args.Error(1)
}

func (coreListener MockChargePointCoreListener) OnChangeAvailability(request *availability.ChangeAvailabilityRequest) (confirmation *availability.ChangeAvailabilityConfirmation, err error) {
	args := coreListener.MethodCalled("OnChangeAvailability", request)
	conf := args.Get(0).(*availability.ChangeAvailabilityConfirmation)
	return conf, args.Error(1)
}

//
//func (coreListener MockChargePointCoreListener) OnDataTransfer(request *ocpp2.DataTransferRequest) (confirmation *ocpp2.DataTransferConfirmation, err error) {
//	args := coreListener.MethodCalled("OnDataTransfer", request)
//	conf := args.Get(0).(*ocpp2.DataTransferConfirmation)
//	return conf, args.Error(1)
//}
//
//func (coreListener MockChargePointCoreListener) OnChangeConfiguration(request *ocpp2.ChangeConfigurationRequest) (confirmation *ocpp2.ChangeConfigurationConfirmation, err error) {
//	args := coreListener.MethodCalled("OnChangeConfiguration", request)
//	conf := args.Get(0).(*ocpp2.ChangeConfigurationConfirmation)
//	return conf, args.Error(1)
//}

func (coreListener MockChargePointCoreListener) OnClearCache(request *authorization.ClearCacheRequest) (confirmation *authorization.ClearCacheConfirmation, err error) {
	args := coreListener.MethodCalled("OnClearCache", request)
	conf := args.Get(0).(*authorization.ClearCacheConfirmation)
	return conf, args.Error(1)
}

func (coreListener MockChargePointCoreListener) OnClearChargingProfile(request *smartcharging.ClearChargingProfileRequest) (confirmation *smartcharging.ClearChargingProfileConfirmation, err error) {
	args := coreListener.MethodCalled("OnClearChargingProfile", request)
	conf := args.Get(0).(*smartcharging.ClearChargingProfileConfirmation)
	return conf, args.Error(1)
}

func (coreListener MockChargePointCoreListener) OnClearDisplay(request *display.ClearDisplayRequest) (confirmation *display.ClearDisplayConfirmation, err error) {
	args := coreListener.MethodCalled("OnClearDisplay", request)
	conf := args.Get(0).(*display.ClearDisplayConfirmation)
	return conf, args.Error(1)
}

func (coreListener MockChargePointCoreListener) OnClearVariableMonitoring(request *diagnostics.ClearVariableMonitoringRequest) (confirmation *diagnostics.ClearVariableMonitoringConfirmation, err error) {
	args := coreListener.MethodCalled("OnClearVariableMonitoring", request)
	conf := args.Get(0).(*diagnostics.ClearVariableMonitoringConfirmation)
	return conf, args.Error(1)
}

func (coreListener MockChargePointCoreListener) OnCostUpdated(request *tariffcost.CostUpdatedRequest) (confirmation *tariffcost.CostUpdatedConfirmation, err error) {
	args := coreListener.MethodCalled("OnCostUpdated", request)
	conf := args.Get(0).(*tariffcost.CostUpdatedConfirmation)
	return conf, args.Error(1)
}

func (coreListener MockChargePointCoreListener) OnCustomerInformation(request *diagnostics.CustomerInformationRequest) (confirmation *diagnostics.CustomerInformationConfirmation, err error) {
	args := coreListener.MethodCalled("OnCustomerInformation", request)
	conf := args.Get(0).(*diagnostics.CustomerInformationConfirmation)
	return conf, args.Error(1)
}

func (coreListener MockChargePointCoreListener) OnDataTransfer(request *data.DataTransferRequest) (confirmation *data.DataTransferConfirmation, err error) {
	args := coreListener.MethodCalled("OnDataTransfer", request)
	conf := args.Get(0).(*data.DataTransferConfirmation)
	return conf, args.Error(1)
}

func (coreListener MockChargePointCoreListener) OnDeleteCertificate(request *iso15118.DeleteCertificateRequest) (confirmation *iso15118.DeleteCertificateConfirmation, err error) {
	args := coreListener.MethodCalled("OnDeleteCertificate", request)
	conf := args.Get(0).(*iso15118.DeleteCertificateConfirmation)
	return conf, args.Error(1)
}

func (coreListener MockChargePointCoreListener) OnGetChargingProfiles(request *ocpp2.GetChargingProfilesRequest) (confirmation *ocpp2.GetChargingProfilesConfirmation, err error) {
	args := coreListener.MethodCalled("OnGetChargingProfiles", request)
	conf := args.Get(0).(*ocpp2.GetChargingProfilesConfirmation)
	return conf, args.Error(1)
}

func (coreListener MockChargePointCoreListener) OnGetCompositeSchedule(request *ocpp2.GetCompositeScheduleRequest) (confirmation *ocpp2.GetCompositeScheduleConfirmation, err error) {
	args := coreListener.MethodCalled("OnGetCompositeSchedule", request)
	conf := args.Get(0).(*ocpp2.GetCompositeScheduleConfirmation)
	return conf, args.Error(1)
}

func (coreListener MockChargePointCoreListener) OnGetDisplayMessages(request *ocpp2.GetDisplayMessagesRequest) (confirmation *ocpp2.GetDisplayMessagesConfirmation, err error) {
	args := coreListener.MethodCalled("OnGetDisplayMessages", request)
	conf := args.Get(0).(*ocpp2.GetDisplayMessagesConfirmation)
	return conf, args.Error(1)
}

func (coreListener MockChargePointCoreListener) OnGetInstalledCertificateIds(request *ocpp2.GetInstalledCertificateIdsRequest) (confirmation *ocpp2.GetInstalledCertificateIdsConfirmation, err error) {
	args := coreListener.MethodCalled("OnGetInstalledCertificateIds", request)
	conf := args.Get(0).(*ocpp2.GetInstalledCertificateIdsConfirmation)
	return conf, args.Error(1)
}

func (coreListener MockChargePointCoreListener) OnGetLocalListVersion(request *ocpp2.GetLocalListVersionRequest) (confirmation *ocpp2.GetLocalListVersionConfirmation, err error) {
	args := coreListener.MethodCalled("OnGetLocalListVersion", request)
	conf := args.Get(0).(*ocpp2.GetLocalListVersionConfirmation)
	return conf, args.Error(1)
}

func (coreListener MockChargePointCoreListener) OnGetLog(request *ocpp2.GetLogRequest) (confirmation *ocpp2.GetLogConfirmation, err error) {
	args := coreListener.MethodCalled("OnGetLog", request)
	conf := args.Get(0).(*ocpp2.GetLogConfirmation)
	return conf, args.Error(1)
}

func (coreListener MockChargePointCoreListener) OnGetMonitoringReport(request *ocpp2.GetMonitoringReportRequest) (confirmation *ocpp2.GetMonitoringReportConfirmation, err error) {
	args := coreListener.MethodCalled("OnGetMonitoringReport", request)
	conf := args.Get(0).(*ocpp2.GetMonitoringReportConfirmation)
	return conf, args.Error(1)
}

//func (coreListener MockChargePointCoreListener) OnGetConfiguration(request *ocpp2.GetConfigurationRequest) (confirmation *ocpp2.GetConfigurationConfirmation, err error) {
//	args := coreListener.MethodCalled("OnGetConfiguration", request)
//	conf := args.Get(0).(*ocpp2.GetConfigurationConfirmation)
//	return conf, args.Error(1)
//}
//
//func (coreListener MockChargePointCoreListener) OnReset(request *ocpp2.ResetRequest) (confirmation *ocpp2.ResetConfirmation, err error) {
//	args := coreListener.MethodCalled("OnReset", request)
//	conf := args.Get(0).(*ocpp2.ResetConfirmation)
//	return conf, args.Error(1)
//}
//
//func (coreListener MockChargePointCoreListener) OnUnlockConnector(request *ocpp2.UnlockConnectorRequest) (confirmation *ocpp2.UnlockConnectorConfirmation, err error) {
//	args := coreListener.MethodCalled("OnUnlockConnector", request)
//	conf := args.Get(0).(*ocpp2.UnlockConnectorConfirmation)
//	return conf, args.Error(1)
//}
//
//func (coreListener MockChargePointCoreListener) OnRemoteStartTransaction(request *ocpp2.RemoteStartTransactionRequest) (confirmation *ocpp2.RemoteStartTransactionConfirmation, err error) {
//	args := coreListener.MethodCalled("OnRemoteStartTransaction", request)
//	conf := args.Get(0).(*ocpp2.RemoteStartTransactionConfirmation)
//	return conf, args.Error(1)
//}
//
//func (coreListener MockChargePointCoreListener) OnRemoteStopTransaction(request *ocpp2.RemoteStopTransactionRequest) (confirmation *ocpp2.RemoteStopTransactionConfirmation, err error) {
//	args := coreListener.MethodCalled("OnRemoteStopTransaction", request)
//	conf := args.Get(0).(*ocpp2.RemoteStopTransactionConfirmation)
//	return conf, args.Error(1)
//}

// ---------------------- MOCK CS LOCAL AUTH LIST LISTENER ----------------------
type MockCentralSystemLocalAuthListListener struct {
	mock.Mock
}

// ---------------------- MOCK CP LOCAL AUTH LIST LISTENER ----------------------
type MockChargePointLocalAuthListListener struct {
	mock.Mock
}

//func (localAuthListListener MockChargePointLocalAuthListListener) OnGetLocalListVersion(request *ocpp2.GetLocalListVersionRequest) (confirmation *ocpp2.GetLocalListVersionConfirmation, err error) {
//	args := localAuthListListener.MethodCalled("OnGetLocalListVersion", request)
//	conf := args.Get(0).(*ocpp2.GetLocalListVersionConfirmation)
//	return conf, args.Error(1)
//}
//
//func (localAuthListListener MockChargePointLocalAuthListListener) OnSendLocalList(request *ocpp2.SendLocalListRequest) (confirmation *ocpp2.SendLocalListConfirmation, err error) {
//	args := localAuthListListener.MethodCalled("OnSendLocalList", request)
//	conf := args.Get(0).(*ocpp2.SendLocalListConfirmation)
//	return conf, args.Error(1)
//}

// ---------------------- MOCK CS FIRMWARE MANAGEMENT LISTENER ----------------------
type MockCentralSystemFirmwareManagementListener struct {
	mock.Mock
}

//func (firmwareListener MockCentralSystemFirmwareManagementListener) OnDiagnosticsStatusNotification(chargePointId string, request *ocpp2.DiagnosticsStatusNotificationRequest) (confirmation *ocpp2.DiagnosticsStatusNotificationConfirmation, err error) {
//	args := firmwareListener.MethodCalled("OnDiagnosticsStatusNotification", chargePointId, request)
//	conf := args.Get(0).(*ocpp2.DiagnosticsStatusNotificationConfirmation)
//	return conf, args.Error(1)
//}
//
//func (firmwareListener MockCentralSystemFirmwareManagementListener) OnFirmwareStatusNotification(chargePointId string, request *ocpp2.FirmwareStatusNotificationRequest) (confirmation *ocpp2.FirmwareStatusNotificationConfirmation, err error) {
//	args := firmwareListener.MethodCalled("OnFirmwareStatusNotification", chargePointId, request)
//	conf := args.Get(0).(*ocpp2.FirmwareStatusNotificationConfirmation)
//	return conf, args.Error(1)
//}

// ---------------------- MOCK CP FIRMWARE MANAGEMENT LISTENER ----------------------
type MockChargePointFirmwareManagementListener struct {
	mock.Mock
}

//func (firmwareListener MockChargePointFirmwareManagementListener) OnGetDiagnostics(request *ocpp2.GetDiagnosticsRequest) (confirmation *ocpp2.GetDiagnosticsConfirmation, err error) {
//	args := firmwareListener.MethodCalled("OnGetDiagnostics", request)
//	conf := args.Get(0).(*ocpp2.GetDiagnosticsConfirmation)
//	return conf, args.Error(1)
//}
//
//func (firmwareListener MockChargePointFirmwareManagementListener) OnUpdateFirmware(request *ocpp2.UpdateFirmwareRequest) (confirmation *ocpp2.UpdateFirmwareConfirmation, err error) {
//	args := firmwareListener.MethodCalled("OnUpdateFirmware", request)
//	conf := args.Get(0).(*ocpp2.UpdateFirmwareConfirmation)
//	return conf, args.Error(1)
//}

// ---------------------- MOCK CS RESERVATION LISTENER ----------------------
type MockCentralSystemReservationListener struct {
	mock.Mock
}

// ---------------------- MOCK CP RESERVATION LISTENER ----------------------
type MockChargePointReservationListener struct {
	mock.Mock
}

//func (reservationListener MockChargePointReservationListener) OnReserveNow(request *ocpp2.ReserveNowRequest) (confirmation *ocpp2.ReserveNowConfirmation, err error) {
//	args := reservationListener.MethodCalled("OnReserveNow", request)
//	conf := args.Get(0).(*ocpp2.ReserveNowConfirmation)
//	return conf, args.Error(1)
//}
//

// ---------------------- MOCK CS REMOTE TRIGGER LISTENER ----------------------
type MockCentralSystemRemoteTriggerListener struct {
	mock.Mock
}

// ---------------------- MOCK CP REMOTE TRIGGER LISTENER ----------------------
type MockChargePointRemoteTriggerListener struct {
	mock.Mock
}

//func (remoteTriggerListener MockChargePointRemoteTriggerListener) OnTriggerMessage(request *ocpp2.TriggerMessageRequest) (confirmation *ocpp2.TriggerMessageConfirmation, err error) {
//	args := remoteTriggerListener.MethodCalled("OnTriggerMessage", request)
//	conf := args.Get(0).(*ocpp2.TriggerMessageConfirmation)
//	return conf, args.Error(1)
//}

// ---------------------- MOCK CS SMART CHARGING LISTENER ----------------------
type MockCentralSystemSmartChargingListener struct {
	mock.Mock
}

// ---------------------- MOCK CP SMART CHARGING LISTENER ----------------------
type MockChargePointSmartChargingListener struct {
	mock.Mock
}

//func (smartChargingListener MockChargePointSmartChargingListener) OnSetChargingProfile(request *ocpp2.SetChargingProfileRequest) (confirmation *ocpp2.SetChargingProfileConfirmation, err error) {
//	args := smartChargingListener.MethodCalled("OnSetChargingProfile", request)
//	conf := args.Get(0).(*ocpp2.SetChargingProfileConfirmation)
//	return conf, args.Error(1)
//}
//
//func (smartChargingListener MockChargePointSmartChargingListener) OnClearChargingProfile(request *ocpp2.ClearChargingProfileRequest) (confirmation *ocpp2.ClearChargingProfileConfirmation, err error) {
//	args := smartChargingListener.MethodCalled("OnClearChargingProfile", request)
//	conf := args.Get(0).(*ocpp2.ClearChargingProfileConfirmation)
//	return conf, args.Error(1)
//}
//
//func (smartChargingListener MockChargePointSmartChargingListener) OnGetCompositeSchedule(request *ocpp2.GetCompositeScheduleRequest) (confirmation *ocpp2.GetCompositeScheduleConfirmation, err error) {
//	args := smartChargingListener.MethodCalled("OnGetCompositeSchedule", request)
//	conf := args.Get(0).(*ocpp2.GetCompositeScheduleConfirmation)
//	return conf, args.Error(1)
//}

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

func setupDefaultCentralSystemHandlers(suite *OcppV2TestSuite, coreListener ocpp2.CSMSHandler, options expectedCentralSystemOptions) {
	t := suite.T()
	suite.csms.SetNewChargingStationHandler(func(chargePointId string) {
		assert.Equal(t, options.clientId, chargePointId)
	})
	suite.csms.SetMessageHandler(coreListener)
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

func setupDefaultChargePointHandlers(suite *OcppV2TestSuite, coreListener ocpp2.ChargingStationHandler, options expectedChargePointOptions) {
	t := suite.T()
	suite.chargingStation.SetMessageHandler(coreListener)
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

func assertDateTimeEquality(t *testing.T, expected *types.DateTime, actual *types.DateTime) {
	assert.Equal(t, expected.FormatTimestamp(), actual.FormatTimestamp())
}

func testUnsupportedRequestFromChargePoint(suite *OcppV2TestSuite, request ocpp.Request, requestJson string, messageId string) {
	t := suite.T()
	wsId := "test_id"
	wsUrl := "someUrl"
	expectedError := fmt.Sprintf("unsupported action %v on charge point, cannot send request", request.GetFeatureName())
	errorDescription := fmt.Sprintf("unsupported action %v on central system", request.GetFeatureName())
	errorJson := fmt.Sprintf(`[4,"%v","%v","%v",null]`, messageId, ocppj.NotSupported, errorDescription)
	channel := NewMockWebSocket(wsId)

	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(errorJson), forwardWrittenMessage: false})
	coreListener := MockCSMSHandler{}
	setupDefaultCentralSystemHandlers(suite, coreListener, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(errorJson), forwardWrittenMessage: true})
	resultChannel := make(chan bool, 1)
	suite.ocppjClient.SetErrorHandler(func(err *ocpp.Error, details interface{}) {
		assert.Equal(t, messageId, err.MessageId)
		assert.Equal(t, ocppj.NotSupported, err.Code)
		assert.Equal(t, errorDescription, err.Description)
		assert.Nil(t, details)
		resultChannel <- true
	})
	// Start
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	// Run request test
	err = suite.chargingStation.SendRequestAsync(request, func(confirmation ocpp.Response, err error) {
		t.Fail()
	})
	require.Error(t, err)
	assert.Equal(t, expectedError, err.Error())
	// Run response test
	suite.ocppjClient.AddPendingRequest(messageId, request)
	err = suite.mockWsServer.MessageHandler(channel, []byte(requestJson))
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func testUnsupportedRequestFromCentralSystem(suite *OcppV2TestSuite, request ocpp.Request, requestJson string, messageId string) {
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
	suite.ocppjServer.SetErrorHandler(func(chargePointId string, err *ocpp.Error, details interface{}) {
		assert.Equal(t, messageId, err.MessageId)
		assert.Equal(t, wsId, chargePointId)
		assert.Equal(t, ocppj.NotSupported, err.Code)
		assert.Equal(t, errorDescription, err.Description)
		assert.Nil(t, details)
	})
	// Start
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	// Run request test
	err = suite.csms.SendRequestAsync(wsId, request, func(confirmation ocpp.Response, err error) {
		t.Fail()
	})
	require.Error(t, err)
	assert.Equal(t, expectedError, err.Error())
	// Run response test
	suite.ocppjServer.AddPendingRequest(messageId, request)
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
type OcppV2TestSuite struct {
	suite.Suite
	ocppjClient        *ocppj.Client
	ocppjServer        *ocppj.Server
	mockWsServer       *MockWebsocketServer
	mockWsClient       *MockWebsocketClient
	chargingStation    ocpp2.ChargingStation
	csms               ocpp2.CSMS
	messageIdGenerator TestRandomIdGenerator
}

type TestRandomIdGenerator struct {
	generator func() string
}

func (testGenerator *TestRandomIdGenerator) generateId() string {
	return testGenerator.generator()
}

var defaultMessageId = "1234"

func (suite *OcppV2TestSuite) SetupTest() {
	coreProfile := ocpp2.CoreProfile
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
	// TODO: init additional profiles
	mockClient := MockWebsocketClient{}
	mockServer := MockWebsocketServer{}
	suite.mockWsClient = &mockClient
	suite.mockWsServer = &mockServer
	suite.ocppjClient = ocppj.NewClient("test_id", suite.mockWsClient, coreProfile, securityProfile, provisioningProfile, authProfile, availabilityProfile, reservationProfile, diagnosticsProfile, dataProfile, displayProfile, firmwareProfile, isoProfile, localAuthProfile, meterProfile, remoteProfile, smartChargingProfile, tariffProfile, transactionsProfile)
	suite.ocppjServer = ocppj.NewServer(suite.mockWsServer, coreProfile, securityProfile, provisioningProfile, authProfile, availabilityProfile, reservationProfile, diagnosticsProfile, dataProfile, displayProfile, firmwareProfile, isoProfile, localAuthProfile, meterProfile, remoteProfile, smartChargingProfile, tariffProfile, transactionsProfile)
	suite.chargingStation = ocpp2.NewChargingStation("test_id", suite.ocppjClient, suite.mockWsClient)
	suite.csms = ocpp2.NewCSMS(suite.ocppjServer, suite.mockWsServer)
	suite.messageIdGenerator = TestRandomIdGenerator{generator: func() string {
		return defaultMessageId
	}}
	ocppj.SetMessageIdGenerator(suite.messageIdGenerator.generateId)
}

//TODO: implement generic protocol tests

func TestOcpp2Protocol(t *testing.T) {
	logrus.SetLevel(logrus.PanicLevel)
	suite.Run(t, new(OcppV2TestSuite))
}
