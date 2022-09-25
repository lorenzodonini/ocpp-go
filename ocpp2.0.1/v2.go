// The package contains an implementation of the OCPP 2.0.1 communication protocol between a Charging Station and an Charging Station Management System in an EV charging infrastructure.
package ocpp2

import (
	"crypto/tls"
	"net"

	"github.com/lorenzodonini/ocpp-go/internal/callbackqueue"
	"github.com/lorenzodonini/ocpp-go/ocpp"
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

type ChargingStationConnection interface {
	ID() string
	RemoteAddr() net.Addr
	TLSConnectionState() *tls.ConnectionState
}

type ChargingStationConnectionHandler func(chargePoint ChargingStationConnection)

// -------------------- v2.0 Charging Station --------------------

// A Charging Station represents the physical system where an EV can be charged.
// You can instantiate a default Charging Station struct by calling NewChargingStation.
//
// The logic for incoming messages needs to be implemented, and message handlers need to be registered with the charging station:
//	handler := &ChargingStationHandler{} 				// Custom struct
//  chargingStation.SetAuthorizationHandler(handler)
//  chargingStation.SetProvisioningHandler(handler)
//  // set more handlers...
// Refer to the ChargingStationHandler interface of each profile for the implementation requirements.
//
// If a handler for a profile is not set, the OCPP library will reply to incoming messages for that profile with a NotImplemented error.
//
// A charging station can be started and stopped using the Start and Stop functions.
// While running, messages can be sent to the CSMS by calling the Charging Station's functions, e.g.
//	bootConf, err := chargingStation.BootNotification(BootReasonPowerUp, "model1", "vendor1")
//
// All messages are synchronous blocking, and return either the response from the CSMS or an error.
// To send asynchronous messages and avoid blocking the calling thread, refer to SendRequestAsync.
type ChargingStation interface {
	// Sends a BootNotificationRequest to the CSMS, along with information about the charging station.
	BootNotification(reason provisioning.BootReason, model string, chargePointVendor string, props ...func(request *provisioning.BootNotificationRequest)) (*provisioning.BootNotificationResponse, error)
	// Requests explicit authorization to the CSMS, provided a valid IdToken (typically the customer's). The CSMS may either authorize or reject the token.
	Authorize(idToken string, tokenType types.IdTokenType, props ...func(request *authorization.AuthorizeRequest)) (*authorization.AuthorizeResponse, error)
	// Notifies the CSMS, that a previously set charging limit was cleared.
	ClearedChargingLimit(chargingLimitSource types.ChargingLimitSourceType, props ...func(request *smartcharging.ClearedChargingLimitRequest)) (*smartcharging.ClearedChargingLimitResponse, error)
	// Performs a custom data transfer to the CSMS. The message payload is not pre-defined and must be supported by the CSMS. Every vendor may implement their own proprietary logic for this message.
	DataTransfer(vendorId string, props ...func(request *data.DataTransferRequest)) (*data.DataTransferResponse, error)
	// Notifies the CSMS of a status change during a firmware update procedure (download, installation).
	FirmwareStatusNotification(status firmware.FirmwareStatus, props ...func(request *firmware.FirmwareStatusNotificationRequest)) (*firmware.FirmwareStatusNotificationResponse, error)
	// Requests a new certificate, required for an ISO 15118 EV, from the CSMS.
	Get15118EVCertificate(schemaVersion string, action iso15118.CertificateAction, exiRequest string, props ...func(request *iso15118.Get15118EVCertificateRequest)) (*iso15118.Get15118EVCertificateResponse, error)
	// Requests the CSMS to provide OCSP certificate status for the charging station's 15118 certificates.
	GetCertificateStatus(ocspRequestData types.OCSPRequestDataType, props ...func(request *iso15118.GetCertificateStatusRequest)) (*iso15118.GetCertificateStatusResponse, error)
	// Notifies the CSMS that the Charging Station is still alive. The response is used for time synchronization purposes.
	Heartbeat(props ...func(request *availability.HeartbeatRequest)) (*availability.HeartbeatResponse, error)
	// Updates the CSMS with the current log upload status.
	LogStatusNotification(status diagnostics.UploadLogStatus, requestID int, props ...func(request *diagnostics.LogStatusNotificationRequest)) (*diagnostics.LogStatusNotificationResponse, error)
	// Sends electrical meter values, not related to a transaction, to the CSMS. This message is deprecated and will be replaced by Device Management Monitoring events.
	MeterValues(evseID int, meterValues []types.MeterValue, props ...func(request *meter.MeterValuesRequest)) (*meter.MeterValuesResponse, error)
	// Informs the CSMS of a charging schedule or charging limit imposed by an External Control System on the Charging Station with ongoing transaction(s).
	NotifyChargingLimit(chargingLimit smartcharging.ChargingLimit, props ...func(request *smartcharging.NotifyChargingLimitRequest)) (*smartcharging.NotifyChargingLimitResponse, error)
	// Notifies the CSMS with raw customer data, previously requested by the CSMS (see CustomerInformationFeature).
	NotifyCustomerInformation(data string, seqNo int, generatedAt types.DateTime, requestID int, props ...func(request *diagnostics.NotifyCustomerInformationRequest)) (*diagnostics.NotifyCustomerInformationResponse, error)
	// Notifies the CSMS of the display messages currently configured on the Charging Station.
	NotifyDisplayMessages(requestID int, props ...func(request *display.NotifyDisplayMessagesRequest)) (*display.NotifyDisplayMessagesResponse, error)
	// Forwards the charging needs of an EV to the CSMS.
	NotifyEVChargingNeeds(evseID int, chargingNeeds smartcharging.ChargingNeeds, props ...func(request *smartcharging.NotifyEVChargingNeedsRequest)) (*smartcharging.NotifyEVChargingNeedsResponse, error)
	// Communicates the charging schedule as calculated by the EV to the CSMS.
	NotifyEVChargingSchedule(timeBase *types.DateTime, evseID int, schedule types.ChargingSchedule, props ...func(request *smartcharging.NotifyEVChargingScheduleRequest)) (*smartcharging.NotifyEVChargingScheduleResponse, error)
	// Notifies the CSMS about monitoring events.
	NotifyEvent(generatedAt *types.DateTime, seqNo int, eventData []diagnostics.EventData, props ...func(request *diagnostics.NotifyEventRequest)) (*diagnostics.NotifyEventResponse, error)
	// Sends a monitoring report to the CSMS, according to parameters specified in the GetMonitoringReport request, previously sent by the CSMS.
	NotifyMonitoringReport(requestID int, seqNo int, generatedAt *types.DateTime, monitorData []diagnostics.MonitoringData, props ...func(request *diagnostics.NotifyMonitoringReportRequest)) (*diagnostics.NotifyMonitoringReportResponse, error)
	// Sends a base report to the CSMS, according to parameters specified in the GetBaseReport request, previously sent by the CSMS.
	NotifyReport(requestID int, generatedAt *types.DateTime, seqNo int, props ...func(request *provisioning.NotifyReportRequest)) (*provisioning.NotifyReportResponse, error)
	// Notifies the CSMS about the current progress of a PublishFirmware operation.
	PublishFirmwareStatusNotification(status firmware.PublishFirmwareStatus, props ...func(request *firmware.PublishFirmwareStatusNotificationRequest)) (*firmware.PublishFirmwareStatusNotificationResponse, error)
	// Reports charging profiles installed in the Charging Station, as requested previously by the CSMS.
	ReportChargingProfiles(requestID int, chargingLimitSource types.ChargingLimitSourceType, evseID int, chargingProfile []types.ChargingProfile, props ...func(request *smartcharging.ReportChargingProfilesRequest)) (*smartcharging.ReportChargingProfilesResponse, error)
	// Notifies the CSMS about a reservation status having changed (i.e. the reservation has expired)
	ReservationStatusUpdate(reservationID int, status reservation.ReservationUpdateStatus, props ...func(request *reservation.ReservationStatusUpdateRequest)) (*reservation.ReservationStatusUpdateResponse, error)
	// Informs the CSMS about critical security events.
	SecurityEventNotification(typ string, timestamp *types.DateTime, props ...func(request *security.SecurityEventNotificationRequest)) (*security.SecurityEventNotificationResponse, error)
	// Requests the CSMS to issue a new certificate.
	SignCertificate(csr string, props ...func(request *security.SignCertificateRequest)) (*security.SignCertificateResponse, error)
	// Informs the CSMS about a connector status change.
	StatusNotification(timestamp *types.DateTime, status availability.ConnectorStatus, evseID int, connectorID int, props ...func(request *availability.StatusNotificationRequest)) (*availability.StatusNotificationResponse, error)
	// Sends information to the CSMS about a transaction, used for billing purposes.
	TransactionEvent(t transactions.TransactionEvent, timestamp *types.DateTime, reason transactions.TriggerReason, seqNo int, info transactions.Transaction, props ...func(request *transactions.TransactionEventRequest)) (*transactions.TransactionEventResponse, error)
	// Registers a handler for incoming security profile messages
	SetSecurityHandler(handler security.ChargingStationHandler)
	// Registers a handler for incoming provisioning profile messages
	SetProvisioningHandler(handler provisioning.ChargingStationHandler)
	// Registers a handler for incoming authorization profile messages
	SetAuthorizationHandler(handler authorization.ChargingStationHandler)
	// Registers a handler for incoming local authorization list profile messages
	SetLocalAuthListHandler(handler localauth.ChargingStationHandler)
	// Registers a handler for incoming transactions profile messages
	SetTransactionsHandler(handler transactions.ChargingStationHandler)
	// Registers a handler for incoming remote control profile messages
	SetRemoteControlHandler(handler remotecontrol.ChargingStationHandler)
	// Registers a handler for incoming availability profile messages
	SetAvailabilityHandler(handler availability.ChargingStationHandler)
	// Registers a handler for incoming reservation profile messages
	SetReservationHandler(handler reservation.ChargingStationHandler)
	// Registers a handler for incoming tariff and cost profile messages
	SetTariffCostHandler(handler tariffcost.ChargingStationHandler)
	// Registers a handler for incoming meter profile messages
	SetMeterHandler(handler meter.ChargingStationHandler)
	// Registers a handler for incoming smart charging messages
	SetSmartChargingHandler(handler smartcharging.ChargingStationHandler)
	// Registers a handler for incoming firmware management messages
	SetFirmwareHandler(handler firmware.ChargingStationHandler)
	// Registers a handler for incoming ISO15118 management messages
	SetISO15118Handler(handler iso15118.ChargingStationHandler)
	// Registers a handler for incoming diagnostics messages
	SetDiagnosticsHandler(handler diagnostics.ChargingStationHandler)
	// Registers a handler for incoming display messages
	SetDisplayHandler(handler display.ChargingStationHandler)
	// Registers a handler for incoming data transfer messages
	SetDataHandler(handler data.ChargingStationHandler)
	// Sends a request to the CSMS.
	// The CSMS will respond with a confirmation, or with an error if the request was invalid or could not be processed.
	// In case of network issues (i.e. the remote host couldn't be reached), the function also returns an error.
	//
	// The request is synchronous blocking.
	SendRequest(request ocpp.Request) (ocpp.Response, error)
	// Sends an asynchronous request to the CSMS.
	// The CSMS will respond with a confirmation message, or with an error if the request was invalid or could not be processed.
	// This result is propagated via a callback, called asynchronously.
	//
	// In case of network issues (i.e. the remote host couldn't be reached), the function returns an error directly. In this case, the callback is never invoked.
	SendRequestAsync(request ocpp.Request, callback func(confirmation ocpp.Response, protoError error)) error
	// Connects to the CSMS and starts the charging station routine.
	// The function doesn't block and returns right away, after having attempted to open a connection to the CSMS.
	// If the connection couldn't be opened, an error is returned.
	//
	// Optional client options must be set before calling this function. Refer to NewChargingStation.
	//
	// No auto-reconnect logic is implemented as of now, but is planned for the future.
	Start(csmsUrl string) error
	// Stops the charging station routine, disconnecting it from the CSMS.
	// Any pending requests are discarded.
	Stop()
	// Errors returns a channel for error messages. If it doesn't exist it es created.
	// The channel is closed by the charging station when stopped.
	Errors() <-chan error
}

// Creates a new OCPP 2.0 charging station client.
// The id parameter is required to uniquely identify the charge point.
//
// The endpoint and client parameters may be omitted, in order to use a default configuration:
//   chargingStation := NewChargingStation("someUniqueId", nil, nil)
//
// Additional networking parameters (e.g. TLS or proxy configuration) may be passed, by creating a custom client.
// Here is an example for a client using TLS configuration with a self-signed certificate:
//	certPool := x509.NewCertPool()
//	data, err := ioutil.ReadFile("serverSelfSignedCertFilename")
//	if err != nil {
//		log.Fatal(err)
//	}
//	ok = certPool.AppendCertsFromPEM(data)
//	if !ok {
//		log.Fatal("couldn't parse PEM certificate")
//	}
//	cs := NewChargingStation("someUniqueId", nil, ws.NewTLSClient(&tls.Config{
//		RootCAs: certPool,
//	})
//
// For more advanced options, or if a custom networking/occpj layer is required,
// please refer to ocppj.Client and ws.WsClient.
func NewChargingStation(id string, endpoint *ocppj.Client, client ws.WsClient) ChargingStation {
	if client == nil {
		client = ws.NewClient()
	}
	client.SetRequestedSubProtocol(types.V201Subprotocol)
	cs := chargingStation{responseHandler: make(chan ocpp.Response, 1), errorHandler: make(chan error, 1), callbacks: callbackqueue.New()}

	if endpoint == nil {
		dispatcher := ocppj.NewDefaultClientDispatcher(ocppj.NewFIFOClientQueue(0))
		endpoint = ocppj.NewClient(id, client, dispatcher, nil, authorization.Profile, availability.Profile, data.Profile, diagnostics.Profile, display.Profile, firmware.Profile, iso15118.Profile, localauth.Profile, meter.Profile, provisioning.Profile, remotecontrol.Profile, reservation.Profile, security.Profile, smartcharging.Profile, tariffcost.Profile, transactions.Profile)
	}

	// Callback invoked by dispatcher, whenever a queued request is canceled, due to timeout.
	endpoint.SetOnRequestCanceled(cs.onRequestTimeout)
	cs.client = endpoint

	cs.client.SetResponseHandler(func(confirmation ocpp.Response, requestId string) {
		cs.responseHandler <- confirmation
	})
	cs.client.SetErrorHandler(func(err *ocpp.Error, details interface{}) {
		cs.errorHandler <- err
	})
	cs.client.SetRequestHandler(cs.handleIncomingRequest)
	return &cs
}

// -------------------- v2.0 CSMS --------------------

// A Charging Station Management System (CSMS) manages Charging Stations and has the information for authorizing Management Users for using its Charging Stations.
// You can instantiate a default CSMS struct by calling the NewCSMS function.
//
// The logic for handling incoming messages needs to be implemented, and message handlers need to be registered with the CSMS:
//	handler := &CSMSHandler{} 				// Custom struct
//  csms.SetAuthorizationHandler(handler)
//  csms.SetProvisioningHandler(handler)
//  // set more handlers...
// Refer to the CSMSHandler interface of each profile for the implementation requirements.
//
// If a handler for a profile is not set, the OCPP library will reply to incoming messages for that profile with a NotImplemented error.
//
// A CSMS can be started by using the Start function.
// To be notified of incoming (dis)connections from charging stations refer to the SetNewChargingStationHandler and SetChargingStationDisconnectedHandler functions.
//
// While running, messages can be sent to a Charging Station by calling the CSMS's functions, e.g.:
//	callback := func(conf *ClearDisplayResponse, err error) {
//		// handle the response...
//	}
//	clearDisplayConf, err := csms.ClearDisplay("cs0001", callback, 10)
// All messages are sent asynchronously and do not block the caller.
type CSMS interface {
	// Cancel a pending reservation, provided the reservationId, on a charging station.
	CancelReservation(clientId string, callback func(*reservation.CancelReservationResponse, error), reservationId int, props ...func(*reservation.CancelReservationRequest)) error
	// Installs a new certificate (chain), signed by the CA, on the charging station. This typically follows a SignCertificate message, initiated by the charging station.
	CertificateSigned(clientId string, callback func(*security.CertificateSignedResponse, error), CertificateSigned string, props ...func(*security.CertificateSignedRequest)) error
	// Instructs a charging station to change its availability to the desired operational status.
	ChangeAvailability(clientId string, callback func(*availability.ChangeAvailabilityResponse, error), operationalStatus availability.OperationalStatus, props ...func(*availability.ChangeAvailabilityRequest)) error
	// Instructs a charging station to clear its current authorization cache. All authorization saved locally will be invalidated.
	ClearCache(clientId string, callback func(*authorization.ClearCacheResponse, error), props ...func(*authorization.ClearCacheRequest)) error
	// Instructs a charging station to clear some or all charging profiles, previously sent to the charging station.
	ClearChargingProfile(clientId string, callback func(*smartcharging.ClearChargingProfileResponse, error), props ...func(request *smartcharging.ClearChargingProfileRequest)) error
	// Removes a specific display message, currently configured in a charging station.
	ClearDisplay(clientId string, callback func(*display.ClearDisplayResponse, error), id int, props ...func(*display.ClearDisplayRequest)) error
	// Removes one or more monitoring settings from a charging station for the given variable IDs.
	ClearVariableMonitoring(clientId string, callback func(*diagnostics.ClearVariableMonitoringResponse, error), id []int, props ...func(*diagnostics.ClearVariableMonitoringRequest)) error
	// Instructs a charging station to display the updated current total cost of an ongoing transaction.
	CostUpdated(clientId string, callback func(*tariffcost.CostUpdatedResponse, error), totalCost float64, transactionId string, props ...func(*tariffcost.CostUpdatedRequest)) error
	// Instructs a charging station to send one or more reports, containing raw customer information.
	CustomerInformation(clientId string, callback func(*diagnostics.CustomerInformationResponse, error), requestId int, report bool, clear bool, props ...func(*diagnostics.CustomerInformationRequest)) error
	// Performs a custom data transfer to a charging station. The message payload is not pre-defined and must be supported by the charging station. Every vendor may implement their own proprietary logic for this message.
	DataTransfer(clientId string, callback func(*data.DataTransferResponse, error), vendorId string, props ...func(*data.DataTransferRequest)) error
	// Deletes a previously installed certificate on a charging station.
	DeleteCertificate(clientId string, callback func(*iso15118.DeleteCertificateResponse, error), data types.CertificateHashData, props ...func(*iso15118.DeleteCertificateRequest)) error
	// Requests a report from a charging station. The charging station will asynchronously send the report in chunks using NotifyReportRequest messages.
	GetBaseReport(clientId string, callback func(*provisioning.GetBaseReportResponse, error), requestId int, reportBase provisioning.ReportBaseType, props ...func(*provisioning.GetBaseReportRequest)) error
	// Request a charging station to report some or all installed charging profiles. The charging station will report these asynchronously using ReportChargingProfiles messages.
	GetChargingProfiles(clientId string, callback func(*smartcharging.GetChargingProfilesResponse, error), chargingProfile smartcharging.ChargingProfileCriterion, props ...func(*smartcharging.GetChargingProfilesRequest)) error
	// Requests a charging station to report the composite charging schedule for the indicated duration and evseID.
	GetCompositeSchedule(clientId string, callback func(*smartcharging.GetCompositeScheduleResponse, error), duration int, evseId int, props ...func(*smartcharging.GetCompositeScheduleRequest)) error
	// Retrieves all messages currently configured on a charging station.
	GetDisplayMessages(clientId string, callback func(*display.GetDisplayMessagesResponse, error), requestId int, props ...func(*display.GetDisplayMessagesRequest)) error
	// Retrieves all installed certificates on a charging station.
	GetInstalledCertificateIds(clientId string, callback func(*iso15118.GetInstalledCertificateIdsResponse, error), typeOfCertificate types.CertificateUse, props ...func(*iso15118.GetInstalledCertificateIdsRequest)) error
	// Queries a charging station for version number of the Local Authorization List.
	GetLocalListVersion(clientId string, callback func(*localauth.GetLocalListVersionResponse, error), props ...func(*localauth.GetLocalListVersionRequest)) error
	// Instructs a charging station to upload a diagnostics or security logfile to the CSMS.
	GetLog(clientId string, callback func(*diagnostics.GetLogResponse, error), logType diagnostics.LogType, requestID int, logParameters diagnostics.LogParameters, props ...func(*diagnostics.GetLogRequest)) error
	// Requests a report about configured monitoring settings per component and variable from a charging station. The reports will be uploaded asynchronously using NotifyMonitoringReport messages.
	GetMonitoringReport(clientId string, callback func(*diagnostics.GetMonitoringReportResponse, error), props ...func(*diagnostics.GetMonitoringReportRequest)) error
	// Requests a custom report about configured monitoring settings per criteria, component and variable from a charging station. The reports will be uploaded asynchronously using NotifyMonitoringReport messages.
	GetReport(clientId string, callback func(*provisioning.GetReportResponse, error), props ...func(*provisioning.GetReportRequest)) error
	// Asks a Charging Station whether it has transaction-related messages waiting to be delivered to the CSMS. When a transactionId is provided, only messages for a specific transaction are asked for.
	GetTransactionStatus(clientId string, callback func(*transactions.GetTransactionStatusResponse, error), props ...func(*transactions.GetTransactionStatusRequest)) error
	// Retrieves from a Charging Station the value of an attribute for one or more Variable of one or more Components.
	GetVariables(clientId string, callback func(*provisioning.GetVariablesResponse, error), variableData []provisioning.GetVariableData, props ...func(*provisioning.GetVariablesRequest)) error
	// Installs a new CA certificate on a Charging station.
	InstallCertificate(clientId string, callback func(*iso15118.InstallCertificateResponse, error), certificateType types.CertificateUse, certificate string, props ...func(*iso15118.InstallCertificateRequest)) error
	// Publishes a firmware to a local controller, allowing charging stations to download the same firmware from the local controller directly.
	PublishFirmware(clientId string, callback func(*firmware.PublishFirmwareResponse, error), location string, checksum string, requestID int, props ...func(request *firmware.PublishFirmwareRequest)) error
	// Remotely triggers a transaction to be started on a charging station.
	RequestStartTransaction(clientId string, callback func(*remotecontrol.RequestStartTransactionResponse, error), remoteStartID int, IdToken types.IdTokenType, props ...func(request *remotecontrol.RequestStartTransactionRequest)) error
	// Remotely triggers an ongoing transaction to be stopped on a charging station.
	RequestStopTransaction(clientId string, callback func(*remotecontrol.RequestStopTransactionResponse, error), transactionID string, props ...func(request *remotecontrol.RequestStopTransactionRequest)) error
	// Attempts to reserve a connector for an EV, on a specific charging station.
	ReserveNow(clientId string, callback func(*reservation.ReserveNowResponse, error), id int, expiryDateTime *types.DateTime, idToken types.IdTokenType, props ...func(request *reservation.ReserveNowRequest)) error
	// Instructs the Charging Station to reset itself.
	Reset(clientId string, callback func(*provisioning.ResetResponse, error), t provisioning.ResetType, props ...func(request *provisioning.ResetRequest)) error
	// Sends a local authorization list to a charging station, which can be used for the authorization of idTokens.
	SendLocalList(clientId string, callback func(*localauth.SendLocalListResponse, error), version int, updateType localauth.UpdateType, props ...func(request *localauth.SendLocalListRequest)) error
	// Sends a charging profile to a charging station, to influence the power/current drawn by EVs.
	SetChargingProfile(clientId string, callback func(*smartcharging.SetChargingProfileResponse, error), evseID int, chargingProfile *types.ChargingProfile, props ...func(request *smartcharging.SetChargingProfileRequest)) error
	// Asks a charging station to configure a new display message, that should be displayed (in the future).
	SetDisplayMessage(clientId string, callback func(*display.SetDisplayMessageResponse, error), message display.MessageInfo, props ...func(request *display.SetDisplayMessageRequest)) error
	// Requests a charging station to activate a set of preconfigured monitoring settings, as denoted by the value of MonitoringBase.
	SetMonitoringBase(clientId string, callback func(*diagnostics.SetMonitoringBaseResponse, error), monitoringBase diagnostics.MonitoringBase, props ...func(request *diagnostics.SetMonitoringBaseRequest)) error
	// Restricts a Charging Station to reporting only monitoring events with a severity number lower than or equal to a certain severity.
	SetMonitoringLevel(clientId string, callback func(*diagnostics.SetMonitoringLevelResponse, error), severity int, props ...func(request *diagnostics.SetMonitoringLevelRequest)) error
	// Updates the connection details on a Charging Station.
	SetNetworkProfile(clientId string, callback func(*provisioning.SetNetworkProfileResponse, error), configurationSlot int, connectionData provisioning.NetworkConnectionProfile, props ...func(request *provisioning.SetNetworkProfileRequest)) error
	// Requests a Charging Station to set monitoring triggers on variables.
	SetVariableMonitoring(clientId string, callback func(*diagnostics.SetVariableMonitoringResponse, error), data []diagnostics.SetMonitoringData, props ...func(request *diagnostics.SetVariableMonitoringRequest)) error
	// Configures/changes the values of a set of variables on a charging station.
	SetVariables(clientId string, callback func(*provisioning.SetVariablesResponse, error), data []provisioning.SetVariableData, props ...func(request *provisioning.SetVariablesRequest)) error
	// Requests a Charging Station to send a charging station-initiated message.
	TriggerMessage(clientId string, callback func(*remotecontrol.TriggerMessageResponse, error), requestedMessage remotecontrol.MessageTrigger, props ...func(request *remotecontrol.TriggerMessageRequest)) error
	// Instructs the Charging Station to unlock a connector, to help out an EV-driver.
	UnlockConnector(clientId string, callback func(*remotecontrol.UnlockConnectorResponse, error), evseID int, connectorID int, props ...func(request *remotecontrol.UnlockConnectorRequest)) error
	// Instructs a Local Controller to stops serving a firmware update to connected Charging Stations.
	UnpublishFirmware(clientId string, callback func(*firmware.UnpublishFirmwareResponse, error), checksum string, props ...func(request *firmware.UnpublishFirmwareRequest)) error
	// Instructs a Charging Station to download and install a firmware update.
	UpdateFirmware(clientId string, callback func(*firmware.UpdateFirmwareResponse, error), requestID int, firmware firmware.Firmware, props ...func(request *firmware.UpdateFirmwareRequest)) error

	// Registers a handler for incoming security profile messages.
	SetSecurityHandler(handler security.CSMSHandler)
	// Registers a handler for incoming provisioning profile messages.
	SetProvisioningHandler(handler provisioning.CSMSHandler)
	// Registers a handler for incoming authorization profile messages.
	SetAuthorizationHandler(handler authorization.CSMSHandler)
	// Registers a handler for incoming local authorization list profile messages.
	SetLocalAuthListHandler(handler localauth.CSMSHandler)
	// Registers a handler for incoming transactions profile messages
	SetTransactionsHandler(handler transactions.CSMSHandler)
	// Registers a handler for incoming remote control profile messages
	SetRemoteControlHandler(handler remotecontrol.CSMSHandler)
	// Registers a handler for incoming availability profile messages
	SetAvailabilityHandler(handler availability.CSMSHandler)
	// Registers a handler for incoming reservation profile messages
	SetReservationHandler(handler reservation.CSMSHandler)
	// Registers a handler for incoming tariff and cost profile messages
	SetTariffCostHandler(handler tariffcost.CSMSHandler)
	// Registers a handler for incoming meter profile messages
	SetMeterHandler(handler meter.CSMSHandler)
	// Registers a handler for incoming smart charging messages
	SetSmartChargingHandler(handler smartcharging.CSMSHandler)
	// Registers a handler for incoming firmware management messages
	SetFirmwareHandler(handler firmware.CSMSHandler)
	// Registers a handler for incoming ISO15118 management messages
	SetISO15118Handler(handler iso15118.CSMSHandler)
	// Registers a handler for incoming diagnostics messages
	SetDiagnosticsHandler(handler diagnostics.CSMSHandler)
	// Registers a handler for incoming display messages
	SetDisplayHandler(handler display.CSMSHandler)
	// Registers a handler for incoming data transfer messages
	SetDataHandler(handler data.CSMSHandler)
	// Registers a handler for new incoming Charging station connections.
	SetNewChargingStationHandler(handler ChargingStationConnectionHandler)
	// Registers a handler for Charging station disconnections.
	SetChargingStationDisconnectedHandler(handler ChargingStationConnectionHandler)
	// Sends an asynchronous request to a Charging Station, identified by the clientId.
	// The charging station will respond with a confirmation message, or with an error if the request was invalid or could not be processed.
	// This result is propagated via a callback, called asynchronously.
	// In case of network issues (i.e. the remote host couldn't be reached), the function returns an error directly. In this case, the callback is never invoked.
	SendRequestAsync(clientId string, request ocpp.Request, callback func(ocpp.Response, error)) error
	// Starts running the CSMS on the specified port and URL.
	// The central system runs as a daemon and handles incoming charge point connections and messages.

	// The function blocks forever, so it is suggested to wrap it in a goroutine, in case other functionality needs to be executed on the main program thread.
	Start(listenPort int, listenPath string)
	// Errors returns a channel for error messages. If it doesn't exist it es created.
	Errors() <-chan error
}

// Creates a new OCPP 2.0 CSMS.
//
// The endpoint and client parameters may be omitted, in order to use a default configuration:
//   csms := NewCSMS(nil, nil)
//
// It is recommended to use the default configuration, unless a custom networking / ocppj layer is required.
// The default dispatcher supports all implemented OCPP 2.0 features out-of-the-box.
//
// If you need a TLS server, you may use the following:
//	csms := NewCSMS(nil, ws.NewTLSServer("certificatePath", "privateKeyPath"))
func NewCSMS(endpoint *ocppj.Server, server ws.WsServer) CSMS {
	if server == nil {
		server = ws.NewServer()
	}
	server.AddSupportedSubprotocol(types.V201Subprotocol)
	if endpoint == nil {
		dispatcher := ocppj.NewDefaultServerDispatcher(ocppj.NewFIFOQueueMap(0))
		endpoint = ocppj.NewServer(server, dispatcher, nil, authorization.Profile, availability.Profile, data.Profile, diagnostics.Profile, display.Profile, firmware.Profile, iso15118.Profile, localauth.Profile, meter.Profile, provisioning.Profile, remotecontrol.Profile, reservation.Profile, security.Profile, smartcharging.Profile, tariffcost.Profile, transactions.Profile)
	}
	cs := newCSMS(endpoint)
	cs.server.SetRequestHandler(func(client ws.Channel, request ocpp.Request, requestId string, action string) {
		cs.handleIncomingRequest(client, request, requestId, action)
	})
	cs.server.SetResponseHandler(func(client ws.Channel, response ocpp.Response, requestId string) {
		cs.handleIncomingResponse(client, response, requestId)
	})
	cs.server.SetErrorHandler(func(client ws.Channel, err *ocpp.Error, details interface{}) {
		cs.handleIncomingError(client, err, details)
	})
	cs.server.SetCanceledRequestHandler(func(clientID string, requestID string, request ocpp.Request, err *ocpp.Error) {
		cs.handleCanceledRequest(clientID, request, err)
	})
	return &cs
}
