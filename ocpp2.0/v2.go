// The package contains an implementation of the OCPP 2.0 communication protocol between a Charging Station and an Charging Station Management System in an EV charging infrastructure.
package ocpp2

import (
	"github.com/gorilla/websocket"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocppj"
	"github.com/lorenzodonini/ocpp-go/ws"
	log "github.com/sirupsen/logrus"
)

// -------------------- v2.0 Charging Station --------------------

// A Charging Station represents the physical system where an EV can be charged.
// You can instantiate a default Charging Station struct by calling NewChargingStation.
//
// The logic for incoming messages needs to be implemented, and the message handler has to be registered with the charging station:
//	handler := &ChargingStationHandler{}
//	centralStation.SetMessageHandler(handler)
// Refer to the ChargingStationHandler interface for the implementation requirements.
//
// A charging station can be started and stopped using the Start and Stop functions.
// While running, messages can be sent to the CSMS by calling the Charging Station's functions, e.g.
//	bootConf, err := chargingStation.BootNotification(BootReasonPowerUp, "model1", "vendor1")
//
// All messages are synchronous blocking, and return either the response from the CSMS or an error.
// To send asynchronous messages and avoid blocking the calling thread, refer to SendRequestAsync.
type ChargingStation interface {
	// Sends a BootNotificationRequest to the CSMS, along with information about the charging station.
	BootNotification(reason BootReason, model string, chargePointVendor string, props ...func(request *BootNotificationRequest)) (*BootNotificationConfirmation, error)
	// Requests explicit authorization to the CSMS, provided a valid IdToken (typically the customer's). The CSMS may either authorize or reject the token.
	Authorize(idToken string, tokenType IdTokenType, props ...func(request *AuthorizeRequest)) (*AuthorizeConfirmation, error)
	// Notifies the CSMS, that a previously set charging limit was cleared.
	ClearedChargingLimit(chargingLimitSource ChargingLimitSourceType, props ...func(request *ClearedChargingLimitRequest)) (*ClearedChargingLimitConfirmation, error)
	// Performs a custom data transfer to the CSMS. The message payload is not pre-defined and must be supported by the CSMS.
	DataTransfer(vendorId string, props ...func(request *DataTransferRequest)) (*DataTransferConfirmation, error)
	// Notifies the CSMS of a status change during a firmware update procedure (download, installation).
	FirmwareStatusNotification(status FirmwareStatus, requestID int, props ...func(request *FirmwareStatusNotificationRequest)) (*FirmwareStatusNotificationConfirmation, error)
	// Requests a new certificate, required for an ISO 15118 EV, from the CSMS.
	Get15118EVCertificate(schemaVersion string, exiRequest string, props ...func(request *Get15118EVCertificateRequest)) (*Get15118EVCertificateConfirmation, error)
	// Requests the CSMS to provide OCSP certificate status for the charging station's 15118 certificates.
	GetCertificateStatus(ocspRequestData OCSPRequestDataType, props ...func(request *GetCertificateStatusRequest)) (*GetCertificateStatusConfirmation, error)
	//Heartbeat(props ...func(request *HeartbeatRequest)) (*HeartbeatConfirmation, error)
	//MeterValues(connectorId int, meterValues []MeterValue, props ...func(request *MeterValuesRequest)) (*MeterValuesConfirmation, error)
	//StartTransaction(connectorId int, idTag string, meterStart int, timestamp *DateTime, props ...func(request *StartTransactionRequest)) (*StartTransactionConfirmation, error)
	//StopTransaction(meterStop int, timestamp *DateTime, transactionId int, props ...func(request *StopTransactionRequest)) (*StopTransactionConfirmation, error)
	//StatusNotification(connectorId int, errorCode ChargePointErrorCode, status ChargePointStatus, props ...func(request *StatusNotificationRequest)) (*StatusNotificationConfirmation, error)
	//DiagnosticsStatusNotification(status DiagnosticsStatus, props ...func(request *DiagnosticsStatusNotificationRequest)) (*DiagnosticsStatusNotificationConfirmation, error)
	//FirmwareStatusNotification(status FirmwareStatus, props ...func(request *FirmwareStatusNotificationRequest)) (*FirmwareStatusNotificationConfirmation, error)

	// SetMessageHandler sets a handler for incoming messages from the CSMS.
	// Refer to ChargingStationHandler for info on how to handle the callbacks.
	SetMessageHandler(handler ChargingStationHandler)
	// Sends a request to the CSMS.
	// The CSMS will respond with a confirmation, or with an error if the request was invalid or could not be processed.
	// In case of network issues (i.e. the remote host couldn't be reached), the function also returns an error.
	//
	// The request is synchronous blocking.
	SendRequest(request ocpp.Request) (ocpp.Confirmation, error)
	// Sends an asynchronous request to the CSMS.
	// The CSMS will respond with a confirmation message, or with an error if the request was invalid or could not be processed.
	// This result is propagated via a callback, called asynchronously.
	//
	// In case of network issues (i.e. the remote host couldn't be reached), the function returns an error directly. In this case, the callback is never invoked.
	SendRequestAsync(request ocpp.Request, callback func(confirmation ocpp.Confirmation, protoError error)) error
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
}

// Creates a new OCPP 2.0 charging station client.
// The id parameter is required to uniquely identify the charge point.
//
// The dispatcher and client parameters may be omitted, in order to use a default configuration:
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
//	cp := NewChargingStation("someUniqueId", nil, ws.NewTLSClient(&tls.Config{
//		RootCAs: certPool,
//	})
//
// For more advanced options, or if a custom networking/occpj layer is required,
// please refer to ocppj.ChargingStation and ws.WsClient.
func NewChargingStation(id string, dispatcher *ocppj.ChargePoint, client ws.WsClient) ChargingStation {
	if client == nil {
		client = ws.NewClient()
	}
	client.AddOption(func(dialer *websocket.Dialer) {
		// Look for v2.0 subprotocol and add it, if not found
		alreadyExists := false
		for _, proto := range dialer.Subprotocols {
			if proto == V2Subprotocol {
				alreadyExists = true
				break
			}
		}
		if !alreadyExists {
			dialer.Subprotocols = append(dialer.Subprotocols, V2Subprotocol)
		}
	})
	if dispatcher == nil {
		dispatcher = ocppj.NewChargePoint(id, client, CoreProfile)
	}
	cp := chargingStation{client: dispatcher, confirmationListener: make(chan ocpp.Confirmation), errorListener: make(chan error)}
	cp.client.SetConfirmationHandler(func(confirmation ocpp.Confirmation, requestId string) {
		cp.confirmationListener <- confirmation
	})
	cp.client.SetErrorHandler(func(err *ocpp.Error, details interface{}) {
		cp.errorListener <- err
	})
	cp.client.SetRequestHandler(cp.handleIncomingRequest)
	return &cp
}

// -------------------- v2.0 CSMS --------------------

// A Charging Station represents the physical system where an EV can be charged.

// A Charging Station Management System (CSMS) manages Charging Stations and has the information for authorizing Management Users for using its Charging Stations.
// You can instantiate a default CSMS struct by calling the NewCSMS function.
//
// The logic for handling incoming messages needs to be implemented, and the message handler has to be registered with the CSMS:
//	handler := &CSMSHandler{}
//	csms.SetMessageHandler(handler)
// Refer to the CSMSHandler interface for the implementation requirements.
//
// A CSMS station can be started by using the Start function.
// To be notified of incoming (dis)connections from charging stations refer to the SetNewChargingStationHandler and SetChargingStationDisconnectedHandler functions.
//
// While running, messages can be sent to a Charging Station by calling the CSMS's, e.g.:
//	callback := func(conf *ClearDisplayConfirmation, err error) {
//		// handle the response...
//	}
//	clearDisplayConf, err := csms.ClearDisplay("cs0001", callback, 10)
// All messages are sent asynchronously and do not block the caller.
type CSMS interface {
	// Cancel a pending reservation, provided the reservationId, on a charging station.
	CancelReservation(clientId string, callback func(*CancelReservationConfirmation, error), reservationId int, props ...func(*CancelReservationRequest)) error
	// The CSMS installs a new certificate, signed by the CA, on the charging station. This typically follows a SignCertificate message, initiated by the charging station.
	CertificateSigned(clientId string, callback func(*CertificateSignedConfirmation, error), certificate []string, props ...func(*CertificateSignedRequest)) error
	// Instructs a charging station to change its availability to the desired operational status.
	ChangeAvailability(clientId string, callback func(*ChangeAvailabilityConfirmation, error), evseID int, operationalStatus OperationalStatus, props ...func(*ChangeAvailabilityRequest)) error
	// Instructs a charging station to clear its current authorization cache. All authorization saved locally will be invalidated.
	ClearCache(clientId string, callback func(*ClearCacheConfirmation, error), props ...func(*ClearCacheRequest)) error
	// Instructs a charging station to clear some or all charging profiles, previously sent to the charging station.
	ClearChargingProfile(clientId string, callback func(*ClearChargingProfileConfirmation, error), props ...func(request *ClearChargingProfileRequest)) error
	// Removes a specific display message, currently configured in a charging station.
	ClearDisplay(clientId string, callback func(*ClearDisplayConfirmation, error), id int, props ...func(*ClearDisplayRequest)) error
	// Removes one or more monitoring settings from a charging station for the given variable IDs.
	ClearVariableMonitoring(clientId string, callback func(*ClearVariableMonitoringConfirmation, error), id []int, props ...func(*ClearVariableMonitoringRequest)) error
	// Instructs a charging station to display the updated current total cost of an ongoing transaction.
	CostUpdated(clientId string, callback func(*CostUpdatedConfirmation, error), totalCost float64, transactionId string, props ...func(*CostUpdatedRequest)) error
	// Instructs a charging station to send one or more reports, containing raw customer information.
	CustomerInformation(clientId string, callback func(*CustomerInformationConfirmation, error), requestId int, report bool, clear bool, props ...func(*CustomerInformationRequest)) error
	// Performs a custom data transfer to a charging station. The message payload is not pre-defined and must be supported by the charging station.
	DataTransfer(clientId string, callback func(*DataTransferConfirmation, error), vendorId string, props ...func(*DataTransferRequest)) error
	// Deletes a previously installed certificate on a charging station.
	DeleteCertificate(clientId string, callback func(*DeleteCertificateConfirmation, error), data CertificateHashData, props ...func(*DeleteCertificateRequest)) error
	// Requests a report from a charging station. The charging station will asynchronously send the report in chunks using NotifyReportRequest messages.
	GetBaseReport(clientId string, callback func(*GetBaseReportConfirmation, error), requestId int, reportBase ReportBaseType, props ...func(*GetBaseReportRequest)) error
	// Request a charging station to report some or all installed charging profiles. The charging station will report these asynchronously using ReportChargingProfiles messages.
	GetChargingProfiles(clientId string, callback func(*GetChargingProfilesConfirmation, error), chargingProfile ChargingProfileCriterion, props ...func(*GetChargingProfilesRequest)) error
	// Requests a charging station to report the composite charging schedule for the indicated duration and evseID.
	GetCompositeSchedule(clientId string, callback func(*GetCompositeScheduleConfirmation, error), duration int, evseId int, props ...func(*GetCompositeScheduleRequest)) error
	// Retrieves all messages currently configured on a charging station.
	GetDisplayMessages(clientId string, callback func(*GetDisplayMessagesConfirmation, error), requestId int, props ...func(*GetDisplayMessagesRequest)) error
	// Retrieves all installed certificates on a charging station.
	GetInstalledCertificateIds(clientId string, callback func(*GetInstalledCertificateIdsConfirmation, error), typeOfCertificate CertificateUse, props ...func(*GetInstalledCertificateIdsRequest)) error
	// Queries a charging station for version number of the Local Authorization List.
	GetLocalListVersion(clientId string, callback func(*GetLocalListVersionConfirmation, error), props ...func(*GetLocalListVersionRequest)) error
	// Instructs a charging station to upload a diagnostics or security logfile to the CSMS.
	GetLog(clientId string, callback func(*GetLogConfirmation, error), logType LogType, requestID int, logParameters LogParameters, props ...func(*GetLogRequest)) error
	// Requests a report about configured monitoring settings per component and variable from a charging station. The reports will be uploaded asynchronously using NotifyMonitoringReport messages.
	GetMonitoringReport(clientId string, callback func(*GetMonitoringReportConfirmation, error), props ...func(*GetMonitoringReportRequest)) error
	//GetConfiguration(clientId string, callback func(*GetConfigurationConfirmation, error), keys []string, props ...func(*GetConfigurationRequest)) error
	//RemoteStartTransaction(clientId string, callback func(*RemoteStartTransactionConfirmation, error), idTag string, props ...func(*RemoteStartTransactionRequest)) error
	//RemoteStopTransaction(clientId string, callback func(*RemoteStopTransactionConfirmation, error), transactionId int, props ...func(request *RemoteStopTransactionRequest)) error
	//Reset(clientId string, callback func(*ResetConfirmation, error), resetType ResetType, props ...func(*ResetRequest)) error
	//UnlockConnector(clientId string, callback func(*UnlockConnectorConfirmation, error), connectorId int, props ...func(*UnlockConnectorRequest)) error
	//GetLocalListVersion(clientId string, callback func(*GetLocalListVersionConfirmation, error), props ...func(request *GetLocalListVersionRequest)) error
	//SendLocalList(clientId string, callback func(*SendLocalListConfirmation, error), version int, updateType UpdateType, props ...func(request *SendLocalListRequest)) error
	//GetDiagnostics(clientId string, callback func(*GetDiagnosticsConfirmation, error), location string, props ...func(request *GetDiagnosticsRequest)) error
	//UpdateFirmware(clientId string, callback func(*UpdateFirmwareConfirmation, error), location string, retrieveDate *DateTime, props ...func(request *UpdateFirmwareRequest)) error
	//ReserveNow(clientId string, callback func(*ReserveNowConfirmation, error), connectorId int, expiryDate *DateTime, idTag string, reservationId int, props ...func(request *ReserveNowRequest)) error
	//CancelReservation(clientId string, callback func(*CancelReservationConfirmation, error), reservationId int, props ...func(request *CancelReservationRequest)) error
	//TriggerMessage(clientId string, callback func(*TriggerMessageConfirmation, error), requestedMessage MessageTrigger, props ...func(request *TriggerMessageRequest)) error
	//SetChargingProfile(clientId string, callback func(*SetChargingProfileConfirmation, error), connectorId int, chargingProfile *ChargingProfile, props ...func(request *SetChargingProfileRequest)) error
	//GetCompositeSchedule(clientId string, callback func(*GetCompositeScheduleConfirmation, error), connectorId int, duration int, props ...func(request *GetCompositeScheduleRequest)) error

	// SetMessageHandler sets a handler for incoming messages from the Charging station.
	// Refer to CSMSHandler for info on how to handle the callbacks.
	SetMessageHandler(handler CSMSHandler)
	// Registers a handler for new incoming Charging station connections.
	SetNewChargingStationHandler(handler func(chargePointId string))
	// Registers a handler for Charging station disconnections.
	SetChargingStationDisconnectedHandler(handler func(chargePointId string))
	// Sends an asynchronous request to a Charging Station, identified by the clientId.
	// The charging station will respond with a confirmation message, or with an error if the request was invalid or could not be processed.
	// This result is propagated via a callback, called asynchronously.
	// In case of network issues (i.e. the remote host couldn't be reached), the function returns an error directly. In this case, the callback is never invoked.
	SendRequestAsync(clientId string, request ocpp.Request, callback func(ocpp.Confirmation, error)) error
	// Starts running the CSMS on the specified port and URL.
	// The central system runs as a daemon and handles incoming charge point connections and messages.

	// The function blocks forever, so it is suggested to wrap it in a goroutine, in case other functionality needs to be executed on the main program thread.
	Start(listenPort int, listenPath string)
}

// Creates a new OCPP 2.0 CSMS.
//
// The dispatcher and client parameters may be omitted, in order to use a default configuration:
//   chargingStation := NewCSMS(nil, nil)
//
// It is recommended to use the default configuration, unless a custom networking / ocppj layer is required.
// The default dispatcher supports all OCPP 2.0 features out-of-the-box.
//
// If you need a TLS server, you may use the following:
//	cs := NewCSMS(nil, ws.NewTLSServer("certificatePath", "privateKeyPath"))
func NewCSMS(dispatcher *ocppj.CentralSystem, server ws.WsServer) CSMS {
	if server == nil {
		server = ws.NewServer()
	}
	server.AddSupportedSubprotocol(V2Subprotocol)
	if dispatcher == nil {
		dispatcher = ocppj.NewCentralSystem(server, CoreProfile)
	}
	cs := csms{
		server:    dispatcher,
		callbacks: map[string]func(confirmation ocpp.Confirmation, err error){}}
	cs.server.SetRequestHandler(cs.handleIncomingRequest)
	cs.server.SetConfirmationHandler(cs.handleIncomingConfirmation)
	cs.server.SetErrorHandler(cs.handleIncomingError)
	return &cs
}

func init() {
	log.New()
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetLevel(log.InfoLevel)
}
