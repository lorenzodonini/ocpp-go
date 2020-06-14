// The package contains an implementation of the OCPP 2.0 communication protocol between a Charging Station and an Charging Station Management System in an EV charging infrastructure.
package ocpp2

import (
	"crypto/tls"

	"github.com/gorilla/websocket"

	"github.com/lorenzodonini/ocpp-go/internal/callbackqueue"
	"github.com/lorenzodonini/ocpp-go/ocpp"
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
)

type ChargingStationConnection interface {
	ID() string
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
	FirmwareStatusNotification(status firmware.FirmwareStatus, requestID int, props ...func(request *firmware.FirmwareStatusNotificationRequest)) (*firmware.FirmwareStatusNotificationResponse, error)
	// Requests a new certificate, required for an ISO 15118 EV, from the CSMS.
	Get15118EVCertificate(schemaVersion string, exiRequest string, props ...func(request *iso15118.Get15118EVCertificateRequest)) (*iso15118.Get15118EVCertificateResponse, error)
	// Requests the CSMS to provide OCSP certificate status for the charging station's 15118 certificates.
	GetCertificateStatus(ocspRequestData types.OCSPRequestDataType, props ...func(request *iso15118.GetCertificateStatusRequest)) (*iso15118.GetCertificateStatusResponse, error)
	// Notifies the CSMS that the Charging Station is still alive. The response is used for time synchronization purposes.
	Heartbeat(props ...func(request *availability.HeartbeatRequest)) (*availability.HeartbeatResponse, error)

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
	client.AddOption(func(dialer *websocket.Dialer) {
		// Look for v2.0 subprotocol and add it, if not found
		alreadyExists := false
		for _, proto := range dialer.Subprotocols {
			if proto == types.V2Subprotocol {
				alreadyExists = true
				break
			}
		}
		if !alreadyExists {
			dialer.Subprotocols = append(dialer.Subprotocols, types.V2Subprotocol)
		}
	})
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
	// The CSMS installs a new certificate (chain), signed by the CA, on the charging station. This typically follows a SignCertificate message, initiated by the charging station.
	CertificateSigned(clientId string, callback func(*security.CertificateSignedResponse, error), certificate []string, props ...func(*security.CertificateSignedRequest)) error
	// Instructs a charging station to change its availability to the desired operational status.
	ChangeAvailability(clientId string, callback func(*availability.ChangeAvailabilityResponse, error), evseID int, operationalStatus availability.OperationalStatus, props ...func(*availability.ChangeAvailabilityRequest)) error
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
	GetVariables(clientId string, callback func(*provisioning.GetVariablesResponse, error), variableData []provisioning.VariableData, props ...func(*provisioning.GetVariablesRequest)) error
	// Installs a new CA certificate on a Charging station.
	InstallCertificate(clientId string, callback func(*iso15118.InstallCertificateResponse, error), certificateType types.CertificateUse, certificate string, props ...func(*iso15118.InstallCertificateRequest)) error
	//GetConfiguration(clientId string, callback func(*GetConfigurationConfirmation, error), keys []string, props ...func(*GetConfigurationRequest)) error
	//RemoteStartTransaction(clientId string, callback func(*RemoteStartTransactionConfirmation, error), idTag string, props ...func(*RemoteStartTransactionRequest)) error
	//RemoteStopTransaction(clientId string, callback func(*RemoteStopTransactionConfirmation, error), transactionId int, props ...func(request *RemoteStopTransactionRequest)) error
	//Reset(clientId string, callback func(*ResetConfirmation, error), resetType ResetType, props ...func(*ResetRequest)) error
	//UnlockConnector(clientId string, callback func(*UnlockConnectorConfirmation, error), connectorId int, props ...func(*UnlockConnectorRequest)) error
	//GetLocalListVersion(clientId string, callback func(*GetLocalListVersionResponse, error), props ...func(request *GetLocalListVersionRequest)) error
	//SendLocalList(clientId string, callback func(*SendLocalListConfirmation, error), version int, updateType UpdateType, props ...func(request *SendLocalListRequest)) error
	//GetDiagnostics(clientId string, callback func(*GetDiagnosticsConfirmation, error), location string, props ...func(request *GetDiagnosticsRequest)) error
	//UpdateFirmware(clientId string, callback func(*UpdateFirmwareConfirmation, error), location string, retrieveDate *DateTime, props ...func(request *UpdateFirmwareRequest)) error
	//ReserveNow(clientId string, callback func(*ReserveNowConfirmation, error), connectorId int, expiryDate *DateTime, idTag string, reservationId int, props ...func(request *ReserveNowRequest)) error
	//CancelReservation(clientId string, callback func(*CancelReservationResponse, error), reservationId int, props ...func(request *CancelReservationRequest)) error
	//TriggerMessage(clientId string, callback func(*TriggerMessageConfirmation, error), requestedMessage MessageTrigger, props ...func(request *TriggerMessageRequest)) error
	//SetChargingProfile(clientId string, callback func(*SetChargingProfileConfirmation, error), connectorId int, chargingProfile *ChargingProfile, props ...func(request *SetChargingProfileRequest)) error
	//GetCompositeSchedule(clientId string, callback func(*GetCompositeScheduleResponse, error), connectorId int, duration int, props ...func(request *GetCompositeScheduleRequest)) error

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
	server.AddSupportedSubprotocol(types.V2Subprotocol)
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
	return &cs
}
