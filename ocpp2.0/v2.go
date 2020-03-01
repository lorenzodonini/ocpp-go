package ocpp2

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocppj"
	"github.com/lorenzodonini/ocpp-go/ws"
	log "github.com/sirupsen/logrus"
)

// -------------------- v1.6 Charge Point --------------------
type ChargePoint interface {
	// Messages
	BootNotification(reason BootReason, chargePointModel string, chargePointVendor string, props ...func(request *BootNotificationRequest)) (*BootNotificationConfirmation, error)
	Authorize(idToken string, tokenType IdTokenType, props ...func(request *AuthorizeRequest)) (*AuthorizeConfirmation, error)
	ClearChargingLimit(chargingLimitSource ChargingLimitSourceType, props ...func(request *ClearedChargingLimitRequest)) (*ClearedChargingLimitConfirmation, error)
	//DataTransfer(vendorId string, props ...func(request *DataTransferRequest)) (*DataTransferConfirmation, error)
	//Heartbeat(props ...func(request *HeartbeatRequest)) (*HeartbeatConfirmation, error)
	//MeterValues(connectorId int, meterValues []MeterValue, props ...func(request *MeterValuesRequest)) (*MeterValuesConfirmation, error)
	//StartTransaction(connectorId int, idTag string, meterStart int, timestamp *DateTime, props ...func(request *StartTransactionRequest)) (*StartTransactionConfirmation, error)
	//StopTransaction(meterStop int, timestamp *DateTime, transactionId int, props ...func(request *StopTransactionRequest)) (*StopTransactionConfirmation, error)
	//StatusNotification(connectorId int, errorCode ChargePointErrorCode, status ChargePointStatus, props ...func(request *StatusNotificationRequest)) (*StatusNotificationConfirmation, error)
	//DiagnosticsStatusNotification(status DiagnosticsStatus, props ...func(request *DiagnosticsStatusNotificationRequest)) (*DiagnosticsStatusNotificationConfirmation, error)
	//FirmwareStatusNotification(status FirmwareStatus, props ...func(request *FirmwareStatusNotificationRequest)) (*FirmwareStatusNotificationConfirmation, error)

	// Logic
	SetChargePointCoreListener(listener ChargePointCoreListener)
	//SetLocalAuthListListener(listener ChargePointLocalAuthListListener)
	//SetFirmwareManagementListener(listener ChargePointFirmwareManagementListener)
	//SetReservationListener(listener ChargePointReservationListener)
	//SetRemoteTriggerListener(listener ChargePointRemoteTriggerListener)
	//SetSmartChargingListener(listener ChargePointSmartChargingListener)
	SendRequest(request ocpp.Request) (ocpp.Confirmation, error)
	SendRequestAsync(request ocpp.Request, callback func(confirmation ocpp.Confirmation, protoError error)) error
	Start(centralSystemUrl string) error
	Stop()
}

type chargePoint struct {
	chargePoint           *ocppj.ChargePoint
	coreListener          ChargePointCoreListener
	//localAuthListListener ChargePointLocalAuthListListener
	//firmwareListener      ChargePointFirmwareManagementListener
	//reservationListener   ChargePointReservationListener
	//remoteTriggerListener ChargePointRemoteTriggerListener
	//smartChargingListener ChargePointSmartChargingListener
	confirmationListener  chan ocpp.Confirmation
	errorListener         chan error
}

// Sends a BootNotificationRequest to the central system, along with information about the charge point.
func (cp *chargePoint) BootNotification(reason BootReason, chargePointModel string, chargePointVendor string, props ...func(request *BootNotificationRequest)) (*BootNotificationConfirmation, error) {
	request := NewBootNotificationRequest(reason, chargePointModel, chargePointVendor)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cp.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*BootNotificationConfirmation), err
	}
}

// Requests explicit authorization to the central system, provided a valid IdTag (typically the client's). The central system may either authorize or reject the client.
func (cp *chargePoint) Authorize(idToken string, tokenType IdTokenType, props ...func(request *AuthorizeRequest)) (*AuthorizeConfirmation, error) {
	request := NewAuthorizationRequest(idToken, tokenType)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cp.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*AuthorizeConfirmation), err
	}
}

func (cp *chargePoint) ClearChargingLimit(chargingLimitSource ChargingLimitSourceType, props ...func(request *ClearedChargingLimitRequest)) (*ClearedChargingLimitConfirmation, error) {
	request := NewClearedChargingLimitRequest(chargingLimitSource)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cp.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*ClearedChargingLimitConfirmation), err
	}
}

//
//// Starts a custom data transfer request. Every vendor may implement their own proprietary logic for this message.
//func (cp *chargePoint) DataTransfer(vendorId string, props ...func(request *DataTransferRequest)) (*DataTransferConfirmation, error) {
//	request := NewDataTransferRequest(vendorId)
//	for _, fn := range props {
//		fn(request)
//	}
//	confirmation, err := cp.SendRequest(request)
//	if err != nil {
//		return nil, err
//	} else {
//		return confirmation.(*DataTransferConfirmation), err
//	}
//}
//
//// Notifies the central system that the charge point is still online. The central system's response is used for time synchronization purposes. It is recommended to perform this operation once every 24 hours.
//func (cp *chargePoint) Heartbeat(props ...func(request *HeartbeatRequest)) (*HeartbeatConfirmation, error) {
//	request := NewHeartbeatRequest()
//	for _, fn := range props {
//		fn(request)
//	}
//	confirmation, err := cp.SendRequest(request)
//	if err != nil {
//		return nil, err
//	} else {
//		return confirmation.(*HeartbeatConfirmation), err
//	}
//}
//
//// Sends a batch of collected meter values to the central system, for billing and analysis. May be done periodically during ongoing transactions.
//func (cp *chargePoint) MeterValues(connectorId int, meterValues []MeterValue, props ...func(request *MeterValuesRequest)) (*MeterValuesConfirmation, error) {
//	request := NewMeterValuesRequest(connectorId, meterValues)
//	for _, fn := range props {
//		fn(request)
//	}
//	confirmation, err := cp.SendRequest(request)
//	if err != nil {
//		return nil, err
//	} else {
//		return confirmation.(*MeterValuesConfirmation), err
//	}
//}
//
//// Requests to start a transaction for a specific connector. The central system will verify the client's IdTag and either accept or reject the transaction.
//func (cp *chargePoint) StartTransaction(connectorId int, idTag string, meterStart int, timestamp *DateTime, props ...func(request *StartTransactionRequest)) (*StartTransactionConfirmation, error) {
//	request := NewStartTransactionRequest(connectorId, idTag, meterStart, timestamp)
//	for _, fn := range props {
//		fn(request)
//	}
//	confirmation, err := cp.SendRequest(request)
//	if err != nil {
//		return nil, err
//	} else {
//		return confirmation.(*StartTransactionConfirmation), err
//	}
//}
//
//// Stops an ongoing transaction. Typically a batch of meter values is passed along with this message.
//func (cp *chargePoint) StopTransaction(meterStop int, timestamp *DateTime, transactionId int, props ...func(request *StopTransactionRequest)) (*StopTransactionConfirmation, error) {
//	request := NewStopTransactionRequest(meterStop, timestamp, transactionId)
//	for _, fn := range props {
//		fn(request)
//	}
//	confirmation, err := cp.SendRequest(request)
//	if err != nil {
//		return nil, err
//	} else {
//		return confirmation.(*StopTransactionConfirmation), err
//	}
//}
//
//// Notifies the central system of a status update. This may apply to the entire charge point or to a single connector.
//func (cp *chargePoint) StatusNotification(connectorId int, errorCode ChargePointErrorCode, status ChargePointStatus, props ...func(request *StatusNotificationRequest)) (*StatusNotificationConfirmation, error) {
//	request := NewStatusNotificationRequest(connectorId, errorCode, status)
//	for _, fn := range props {
//		fn(request)
//	}
//	confirmation, err := cp.SendRequest(request)
//	if err != nil {
//		return nil, err
//	} else {
//		return confirmation.(*StatusNotificationConfirmation), err
//	}
//}
//
//// Notifies the central system of a status change in the upload of diagnostics data.
//func (cp *chargePoint) DiagnosticsStatusNotification(status DiagnosticsStatus, props ...func(request *DiagnosticsStatusNotificationRequest)) (*DiagnosticsStatusNotificationConfirmation, error) {
//	request := NewDiagnosticsStatusNotificationRequest(status)
//	for _, fn := range props {
//		fn(request)
//	}
//	confirmation, err := cp.SendRequest(request)
//	if err != nil {
//		return nil, err
//	} else {
//		return confirmation.(*DiagnosticsStatusNotificationConfirmation), err
//	}
//}
//
//// Notifies the central system of a status change during the download of a new firmware version.
//func (cp *chargePoint) FirmwareStatusNotification(status FirmwareStatus, props ...func(request *FirmwareStatusNotificationRequest)) (*FirmwareStatusNotificationConfirmation, error) {
//	request := NewFirmwareStatusNotificationRequest(status)
//	for _, fn := range props {
//		fn(request)
//	}
//	confirmation, err := cp.SendRequest(request)
//	if err != nil {
//		return nil, err
//	} else {
//		return confirmation.(*FirmwareStatusNotificationConfirmation), err
//	}
//}

// Registers a handler for incoming core profile messages
func (cp *chargePoint) SetChargePointCoreListener(listener ChargePointCoreListener) {
	cp.coreListener = listener
}

//// Registers a handler for incoming local authorization profile messages
//func (cp *chargePoint) SetLocalAuthListListener(listener ChargePointLocalAuthListListener) {
//	cp.localAuthListListener = listener
//}
//
//// Registers a handler for incoming firmware management profile messages
//func (cp *chargePoint) SetFirmwareManagementListener(listener ChargePointFirmwareManagementListener) {
//	cp.firmwareListener = listener
//}
//
//// Registers a handler for incoming reservation profile messages
//func (cp *chargePoint) SetReservationListener(listener ChargePointReservationListener) {
//	cp.reservationListener = listener
//}
//
//// Registers a handler for incoming remote trigger profile messages
//func (cp *chargePoint) SetRemoteTriggerListener(listener ChargePointRemoteTriggerListener) {
//	cp.remoteTriggerListener = listener
//}
//
//// Registers a handler for incoming smart charging profile messages
//func (cp *chargePoint) SetSmartChargingListener(listener ChargePointSmartChargingListener) {
//	cp.smartChargingListener = listener
//}

// Sends a request to the central system.
// The central system will respond with a confirmation, or with an error if the request was invalid or could not be processed.
// In case of network issues (i.e. the remote host couldn't be reached), the function also returns an error.
// The request is synchronous blocking.
func (cp *chargePoint) SendRequest(request ocpp.Request) (ocpp.Confirmation, error) {
	// TODO: check for supported feature
	err := cp.chargePoint.SendRequest(request)
	if err != nil {
		return nil, err
	}
	//TODO: timeouts
	select {
	case confirmation := <-cp.confirmationListener:
		return confirmation, nil
	case err = <-cp.errorListener:
		return nil, err
	}
}

// Sends an asynchronous request to the central system.
// The central system will respond with a confirmation messages, or with an error if the request was invalid or could not be processed.
// This result is propagated via a callback, called asynchronously.
// In case of network issues (i.e. the remote host couldn't be reached), the function returns an error directly. In this case, the callback is never called.
func (cp *chargePoint) SendRequestAsync(request ocpp.Request, callback func(confirmation ocpp.Confirmation, err error)) error {
	switch request.GetFeatureName() {
	case AuthorizeFeatureName, BootNotificationFeatureName, ClearedChargingLimitFeatureName:
		break
	default:
		return fmt.Errorf("unsupported action %v on charge point, cannot send request", request.GetFeatureName())
	}
	err := cp.chargePoint.SendRequest(request)
	if err == nil {
		// Retrieve result asynchronously
		go func() {
			select {
			case confirmation := <-cp.confirmationListener:
				callback(confirmation, nil)
			case protoError := <-cp.errorListener:
				callback(nil, protoError)
			}
		}()
	}
	return err
}

func (cp *chargePoint) sendResponse(confirmation ocpp.Confirmation, err error, requestId string) {
	if confirmation != nil {
		err := cp.chargePoint.SendConfirmation(requestId, confirmation)
		if err != nil {
			log.WithField("request", requestId).Errorf("unknown error %v while replying to message with CallError", err)
			//TODO: handle error somehow
		}
	} else {
		err = cp.chargePoint.SendError(requestId, ocppj.ProtocolError, err.Error(), nil)
		if err != nil {
			log.WithField("request", requestId).Errorf("unknown error %v while replying to message with CallError", err)
		}
	}
}

// Connects to the central system and starts the charge point routine.
// The function doesn't block and returns right away, after having attempted to open a connection to the central system.
// If the connection couldn't be opened, an error is returned.
//
// Optional client options must be set before calling this function. Refer to NewChargePoint.
//
// No auto-reconnect logic is implemented as of now, but is planned for the future.
func (cp *chargePoint) Start(centralSystemUrl string) error {
	// TODO: implement auto-reconnect logic
	return cp.chargePoint.Start(centralSystemUrl)
}

// Stops the charge point routine, disconnecting it from the central system.
// Any pending requests are discarded.
func (cp *chargePoint) Stop() {
	cp.chargePoint.Stop()
}

func (cp *chargePoint) notImplementedError(requestId string, action string) {
	log.WithField("request", requestId).Errorf("cannot handle call from central system. Sending CallError instead")
	err := cp.chargePoint.SendError(requestId, ocppj.NotImplemented, fmt.Sprintf("no handler for action %v implemented", action), nil)
	if err != nil {
		log.WithField("request", requestId).Errorf("unknown error %v while replying to message with CallError", err)
	}
}

func (cp *chargePoint) notSupportedError(requestId string, action string) {
	log.WithField("request", requestId).Errorf("cannot handle call from central system. Sending CallError instead")
	err := cp.chargePoint.SendError(requestId, ocppj.NotSupported, fmt.Sprintf("unsupported action %v on charge point", action), nil)
	if err != nil {
		log.WithField("request", requestId).Errorf("unknown error %v while replying to message with CallError", err)
	}
}

func (cp *chargePoint) handleIncomingRequest(request ocpp.Request, requestId string, action string) {
	profile, found := cp.chargePoint.GetProfileForFeature(action)
	// Check whether action is supported and a listener for it exists
	if !found {
		cp.notImplementedError(requestId, action)
		return
	} else {
		switch profile.Name {
		case CoreProfileName:
			if cp.coreListener == nil {
				cp.notSupportedError(requestId, action)
				return
			}
		//case LocalAuthListProfileName:
		//	if cp.localAuthListListener == nil {
		//		cp.notSupportedError(requestId, action)
		//		return
		//	}
		//case FirmwareManagementProfileName:
		//	if cp.firmwareListener == nil {
		//		cp.notSupportedError(requestId, action)
		//		return
		//	}
		//case ReservationProfileName:
		//	if cp.reservationListener == nil {
		//		cp.notSupportedError(requestId, action)
		//		return
		//	}
		//case RemoteTriggerProfileName:
		//	if cp.remoteTriggerListener == nil {
		//		cp.notSupportedError(requestId, action)
		//		return
		//	}
		//case SmartChargingProfileName:
		//	if cp.smartChargingListener == nil {
		//		cp.notSupportedError(requestId, action)
		//		return
		//	}
		}
	}
	// Process request
	var confirmation ocpp.Confirmation = nil
	cp.chargePoint.GetProfileForFeature(action)
	var err error = nil
	switch action {
	case CancelReservationFeatureName:
		confirmation, err = cp.coreListener.OnCancelReservation(request.(*CancelReservationRequest))
	case CertificateSignedFeatureName:
		confirmation, err = cp.coreListener.OnCertificateSigned(request.(*CertificateSignedRequest))
	case ChangeAvailabilityFeatureName:
		confirmation, err = cp.coreListener.OnChangeAvailability(request.(*ChangeAvailabilityRequest))
	//case ChangeConfigurationFeatureName:
	//	confirmation, err = cp.coreListener.OnChangeConfiguration(request.(*ChangeConfigurationRequest))
	case ClearCacheFeatureName:
		confirmation, err = cp.coreListener.OnClearCache(request.(*ClearCacheRequest))
	case ClearChargingProfileFeatureName:
		confirmation, err = cp.coreListener.OnClearChargingProfile(request.(*ClearChargingProfileRequest))
	case ClearDisplayFeatureName:
		confirmation, err = cp.coreListener.OnClearDisplay(request.(*ClearDisplayRequest))
	case ClearVariableMonitoringFeatureName:
		confirmation, err = cp.coreListener.OnClearVariableMonitoring(request.(*ClearVariableMonitoringRequest))
	//case DataTransferFeatureName:
	//	confirmation, err = cp.coreListener.OnDataTransfer(request.(*DataTransferRequest))
	//case GetConfigurationFeatureName:
	//	confirmation, err = cp.coreListener.OnGetConfiguration(request.(*GetConfigurationRequest))
	//case RemoteStartTransactionFeatureName:
	//	confirmation, err = cp.coreListener.OnRemoteStartTransaction(request.(*RemoteStartTransactionRequest))
	//case RemoteStopTransactionFeatureName:
	//	confirmation, err = cp.coreListener.OnRemoteStopTransaction(request.(*RemoteStopTransactionRequest))
	//case ResetFeatureName:
	//	confirmation, err = cp.coreListener.OnReset(request.(*ResetRequest))
	//case UnlockConnectorFeatureName:
	//	confirmation, err = cp.coreListener.OnUnlockConnector(request.(*UnlockConnectorRequest))
	//case GetLocalListVersionFeatureName:
	//	confirmation, err = cp.localAuthListListener.OnGetLocalListVersion(request.(*GetLocalListVersionRequest))
	//case SendLocalListFeatureName:
	//	confirmation, err = cp.localAuthListListener.OnSendLocalList(request.(*SendLocalListRequest))
	//case GetDiagnosticsFeatureName:
	//	confirmation, err = cp.firmwareListener.OnGetDiagnostics(request.(*GetDiagnosticsRequest))
	//case UpdateFirmwareFeatureName:
	//	confirmation, err = cp.firmwareListener.OnUpdateFirmware(request.(*UpdateFirmwareRequest))
	//case ReserveNowFeatureName:
	//	confirmation, err = cp.reservationListener.OnReserveNow(request.(*ReserveNowRequest))
	//case CancelReservationFeatureName:
	//	confirmation, err = cp.reservationListener.OnCancelReservation(request.(*CancelReservationRequest))
	//case TriggerMessageFeatureName:
	//	confirmation, err = cp.remoteTriggerListener.OnTriggerMessage(request.(*TriggerMessageRequest))
	//case SetChargingProfileFeatureName:
	//	confirmation, err = cp.smartChargingListener.OnSetChargingProfile(request.(*SetChargingProfileRequest))
	//case ClearChargingProfileFeatureName:
	//	confirmation, err = cp.smartChargingListener.OnClearChargingProfile(request.(*ClearChargingProfileRequest))
	//case GetCompositeScheduleFeatureName:
	//	confirmation, err = cp.smartChargingListener.OnGetCompositeSchedule(request.(*GetCompositeScheduleRequest))
	default:
		cp.notSupportedError(requestId, action)
		return
	}
	cp.sendResponse(confirmation, err, requestId)
}

// Creates a new OCPP 2.0 charge point client.
// The id parameter is required to uniquely identify the charge point.
//
// The dispatcher and client parameters may be omitted, in order to use a default configuration:
//   chargePoint := NewChargePoint("someUniqueId", nil, nil)
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
//	cp := NewChargePoint("someUniqueId", nil, ws.NewTLSClient(&tls.Config{
//		RootCAs: certPool,
//	})
//
// For more advanced options, or if a customer networking/occpj layer is required,
// please refer to ocppj.ChargePoint and ws.WsClient.
func NewChargePoint(id string, dispatcher *ocppj.ChargePoint, client ws.WsClient) ChargePoint {
	if client == nil {
		client = ws.NewClient()
	}
	client.AddOption(func (dialer *websocket.Dialer) {
		// Look for v1.6 subprotocol and add it, if not found
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
	cp := chargePoint{chargePoint: dispatcher, confirmationListener: make(chan ocpp.Confirmation), errorListener: make(chan error)}
	cp.chargePoint.SetConfirmationHandler(func(confirmation ocpp.Confirmation, requestId string) {
		cp.confirmationListener <- confirmation
	})
	cp.chargePoint.SetErrorHandler(func(err *ocpp.Error, details interface{}) {
		cp.errorListener <- err
	})
	cp.chargePoint.SetRequestHandler(cp.handleIncomingRequest)
	return &cp
}

// -------------------- v2.0 CSMS --------------------
type CSMS interface {
	// Messages
	CancelReservation(clientId string, callback func(*CancelReservationConfirmation, error), reservationId int, props ...func(*CancelReservationRequest)) error
	CertificateSigned(clientId string, callback func(*CertificateSignedConfirmation, error), certificate []string, props ...func(*CertificateSignedRequest)) error
	ChangeAvailability(clientId string, callback func(*ChangeAvailabilityConfirmation, error), evseID int, operationalStatus OperationalStatus, props ...func(*ChangeAvailabilityRequest)) error
	//ChangeConfiguration(clientId string, callback func(*ChangeConfigurationConfirmation, error), key string, value string, props ...func(*ChangeConfigurationRequest)) error
	ClearCache(clientId string, callback func(*ClearCacheConfirmation, error), props ...func(*ClearCacheRequest)) error
	ClearChargingProfile(clientId string, callback func(*ClearChargingProfileConfirmation, error), props ...func(request *ClearChargingProfileRequest)) error
	ClearDisplay(clientId string, callback func(*ClearDisplayConfirmation, error), id int, props ...func(*ClearDisplayRequest)) error
	ClearVariableMonitoring(clientId string, callback func(*ClearVariableMonitoringConfirmation, error), id []int, props ...func(*ClearVariableMonitoringRequest)) error
	//DataTransfer(clientId string, callback func(*DataTransferConfirmation, error), vendorId string, props ...func(*DataTransferRequest)) error
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

	// Logic
	SetCentralSystemCoreListener(listener CentralSystemCoreListener)
	//SetLocalAuthListListener(listener CentralSystemLocalAuthListListener)
	//SetFirmwareManagementListener(listener CentralSystemFirmwareManagementListener)
	//SetReservationListener(listener CentralSystemReservationListener)
	//SetRemoteTriggerListener(listener CentralSystemRemoteTriggerListener)
	//SetSmartChargingListener(listener CentralSystemSmartChargingListener)
	SetNewChargePointHandler(handler func(chargePointId string))
	SetChargePointDisconnectedHandler(handler func(chargePointId string))
	SendRequestAsync(clientId string, request ocpp.Request, callback func(ocpp.Confirmation, error)) error
	Start(listenPort int, listenPath string)
}

type csms struct {
	centralSystem         *ocppj.CentralSystem
	coreListener          CentralSystemCoreListener
	//localAuthListListener CentralSystemLocalAuthListListener
	//firmwareListener      CentralSystemFirmwareManagementListener
	//reservationListener   CentralSystemReservationListener
	//remoteTriggerListener CentralSystemRemoteTriggerListener
	//smartChargingListener CentralSystemSmartChargingListener
	callbacks             map[string]func(confirmation ocpp.Confirmation, err error)
}

// Cancels a previously reserved charge point or connector, given the reservation Id.
func (cs *csms) CancelReservation(clientId string, callback func(*CancelReservationConfirmation, error), reservationId int, props ...func(request *CancelReservationRequest)) error {
	request := NewCancelReservationRequest(reservationId)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*CancelReservationConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

// Sends a new certificate (chain) to the charging station.
func (cs *csms) CertificateSigned(clientId string, callback func(*CertificateSignedConfirmation, error), certificate []string, props ...func(*CertificateSignedRequest)) error {
	request := NewCertificateSignedRequest(certificate)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*CertificateSignedConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

// Instructs a charge point to change its availability. The target availability can be set for a single evse of for the whole charging station.
func (cs *csms) ChangeAvailability(clientId string, callback func(confirmation *ChangeAvailabilityConfirmation, err error), evseID int, operationalStatus OperationalStatus, props ...func(request *ChangeAvailabilityRequest)) error {
	request := NewChangeAvailabilityRequest(evseID, operationalStatus)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*ChangeAvailabilityConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}
//
//// Changes the configuration of a charge point, by setting a specific key-value pair.
//// The configuration key must be supported by the target charge point, in order for the configuration to be accepted.
//func (cs *centralSystem) ChangeConfiguration(clientId string, callback func(confirmation *ChangeConfigurationConfirmation, err error), key string, value string, props ...func(request *ChangeConfigurationRequest)) error {
//	request := NewChangeConfigurationRequest(key, value)
//	for _, fn := range props {
//		fn(request)
//	}
//	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
//		if confirmation != nil {
//			callback(confirmation.(*ChangeConfigurationConfirmation), protoError)
//		} else {
//			callback(nil, protoError)
//		}
//	}
//	return cs.SendRequestAsync(clientId, request, genericCallback)
//}

// Instructs the charge point to clear its current authorization cache. All authorization saved locally will be invalidated.
func (cs *csms) ClearCache(clientId string, callback func(confirmation *ClearCacheConfirmation, err error), props ...func(*ClearCacheRequest)) error {
	request := NewClearCacheRequest()
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*ClearCacheConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

// Removes one or more charging profiles from a charging station.
func (cs *csms) ClearChargingProfile(clientId string, callback func(*ClearChargingProfileConfirmation, error), props ...func(request *ClearChargingProfileRequest)) error {
	request := NewClearChargingProfileRequest()
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*ClearChargingProfileConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) ClearDisplay(clientId string, callback func(*ClearDisplayConfirmation, error), id int, props ...func(*ClearDisplayRequest)) error {
	request := NewClearDisplayRequest(id)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*ClearDisplayConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) ClearVariableMonitoring(clientId string, callback func(*ClearVariableMonitoringConfirmation, error), id []int, props ...func(*ClearVariableMonitoringRequest)) error {
	request := NewClearVariableMonitoringRequest(id)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*ClearVariableMonitoringConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

//// Starts a custom data transfer request. Every vendor may implement their own proprietary logic for this message.
//func (cs *centralSystem) DataTransfer(clientId string, callback func(confirmation *DataTransferConfirmation, err error), vendorId string, props ...func(request *DataTransferRequest)) error {
//	request := NewDataTransferRequest(vendorId)
//	for _, fn := range props {
//		fn(request)
//	}
//	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
//		if confirmation != nil {
//			callback(confirmation.(*DataTransferConfirmation), protoError)
//		} else {
//			callback(nil, protoError)
//		}
//	}
//	return cs.SendRequestAsync(clientId, request, genericCallback)
//}
//
//// Retrieves the configuration values for the provided configuration keys.
//func (cs *centralSystem) GetConfiguration(clientId string, callback func(confirmation *GetConfigurationConfirmation, err error), keys []string, props ...func(request *GetConfigurationRequest)) error {
//	request := NewGetConfigurationRequest(keys)
//	for _, fn := range props {
//		fn(request)
//	}
//	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
//		if confirmation != nil {
//			callback(confirmation.(*GetConfigurationConfirmation), protoError)
//		} else {
//			callback(nil, protoError)
//		}
//	}
//	return cs.SendRequestAsync(clientId, request, genericCallback)
//}
//
//// Instructs a charge point to start a transaction for a specified client on a provided connector.
//// Depending on the configuration, an explicit authorization message may still be required, before the transaction can start.
//func (cs *centralSystem) RemoteStartTransaction(clientId string, callback func(*RemoteStartTransactionConfirmation, error), idTag string, props ...func(*RemoteStartTransactionRequest)) error {
//	request := NewRemoteStartTransactionRequest(idTag)
//	for _, fn := range props {
//		fn(request)
//	}
//	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
//		if confirmation != nil {
//			callback(confirmation.(*RemoteStartTransactionConfirmation), protoError)
//		} else {
//			callback(nil, protoError)
//		}
//	}
//	return cs.SendRequestAsync(clientId, request, genericCallback)
//}
//
//// Instructs a charge point to stop an ongoing transaction, given the transaction's ID.
//func (cs *centralSystem) RemoteStopTransaction(clientId string, callback func(*RemoteStopTransactionConfirmation, error), transactionId int, props ...func(request *RemoteStopTransactionRequest)) error {
//	request := NewRemoteStopTransactionRequest(transactionId)
//	for _, fn := range props {
//		fn(request)
//	}
//	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
//		if confirmation != nil {
//			callback(confirmation.(*RemoteStopTransactionConfirmation), protoError)
//		} else {
//			callback(nil, protoError)
//		}
//	}
//	return cs.SendRequestAsync(clientId, request, genericCallback)
//}
//
//// Forces a charge point to perform an internal hard or soft reset. In both cases, all ongoing transactions are stopped.
//func (cs *centralSystem) Reset(clientId string, callback func(*ResetConfirmation, error), resetType ResetType, props ...func(request *ResetRequest)) error {
//	request := NewResetRequest(resetType)
//	for _, fn := range props {
//		fn(request)
//	}
//	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
//		if confirmation != nil {
//			callback(confirmation.(*ResetConfirmation), protoError)
//		} else {
//			callback(nil, protoError)
//		}
//	}
//	return cs.SendRequestAsync(clientId, request, genericCallback)
//}
//
//// Attempts to unlock a specific connector on a charge point. Used for remote support purposes.
//func (cs *centralSystem) UnlockConnector(clientId string, callback func(*UnlockConnectorConfirmation, error), connectorId int, props ...func(*UnlockConnectorRequest)) error {
//	request := NewUnlockConnectorRequest(connectorId)
//	for _, fn := range props {
//		fn(request)
//	}
//	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
//		if confirmation != nil {
//			callback(confirmation.(*UnlockConnectorConfirmation), protoError)
//		} else {
//			callback(nil, protoError)
//		}
//	}
//	return cs.SendRequestAsync(clientId, request, genericCallback)
//}
//
//// Queries the current version of the local authorization list from a charge point.
//func (cs *centralSystem) GetLocalListVersion(clientId string, callback func(*GetLocalListVersionConfirmation, error), props ...func(request *GetLocalListVersionRequest)) error {
//	request := NewGetLocalListVersionRequest()
//	for _, fn := range props {
//		fn(request)
//	}
//	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
//		if confirmation != nil {
//			callback(confirmation.(*GetLocalListVersionConfirmation), protoError)
//		} else {
//			callback(nil, protoError)
//		}
//	}
//	return cs.SendRequestAsync(clientId, request, genericCallback)
//}
//
//// Sends or updates a local authorization list on a charge point. Versioning rules must be followed.
//func (cs *centralSystem) SendLocalList(clientId string, callback func(*SendLocalListConfirmation, error), version int, updateType UpdateType, props ...func(request *SendLocalListRequest)) error {
//	request := NewSendLocalListRequest(version, updateType)
//	for _, fn := range props {
//		fn(request)
//	}
//	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
//		if confirmation != nil {
//			callback(confirmation.(*SendLocalListConfirmation), protoError)
//		} else {
//			callback(nil, protoError)
//		}
//	}
//	return cs.SendRequestAsync(clientId, request, genericCallback)
//}
//
//// Requests diagnostics data from a charge point. The data will be uploaded out-of-band to the provided URL location.
//func (cs *centralSystem) GetDiagnostics(clientId string, callback func(*GetDiagnosticsConfirmation, error), location string, props ...func(request *GetDiagnosticsRequest)) error {
//	request := NewGetDiagnosticsRequest(location)
//	for _, fn := range props {
//		fn(request)
//	}
//	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
//		if confirmation != nil {
//			callback(confirmation.(*GetDiagnosticsConfirmation), protoError)
//		} else {
//			callback(nil, protoError)
//		}
//	}
//	return cs.SendRequestAsync(clientId, request, genericCallback)
//}
//
//// Instructs the charge point to download and install a new firmware version. The firmware binary will be downloaded out-of-band from the provided URL location.
//func (cs *centralSystem) UpdateFirmware(clientId string, callback func(*UpdateFirmwareConfirmation, error), location string, retrieveDate *DateTime, props ...func(request *UpdateFirmwareRequest)) error {
//	request := NewUpdateFirmwareRequest(location, retrieveDate)
//	for _, fn := range props {
//		fn(request)
//	}
//	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
//		if confirmation != nil {
//			callback(confirmation.(*UpdateFirmwareConfirmation), protoError)
//		} else {
//			callback(nil, protoError)
//		}
//	}
//	return cs.SendRequestAsync(clientId, request, genericCallback)
//}
//
//// Instructs the charge point to reserve a connector for a specific IdTag (client). The connector, or the entire charge point, will be reserved until the provided expiration time.
//func (cs *centralSystem) ReserveNow(clientId string, callback func(*ReserveNowConfirmation, error), connectorId int, expiryDate *DateTime, idTag string, reservationId int, props ...func(request *ReserveNowRequest)) error {
//	request := NewReserveNowRequest(connectorId, expiryDate, idTag, reservationId)
//	for _, fn := range props {
//		fn(request)
//	}
//	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
//		if confirmation != nil {
//			callback(confirmation.(*ReserveNowConfirmation), protoError)
//		} else {
//			callback(nil, protoError)
//		}
//	}
//	return cs.SendRequestAsync(clientId, request, genericCallback)
//}
//
//// Instructs a charge point to send a specific message to the central system. This is used for forcefully triggering status updates, when the last known state is either too old or not clear to the central system.
//func (cs *centralSystem) TriggerMessage(clientId string, callback func(*TriggerMessageConfirmation, error), requestedMessage MessageTrigger, props ...func(request *TriggerMessageRequest)) error {
//	request := NewTriggerMessageRequest(requestedMessage)
//	for _, fn := range props {
//		fn(request)
//	}
//	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
//		if confirmation != nil {
//			callback(confirmation.(*TriggerMessageConfirmation), protoError)
//		} else {
//			callback(nil, protoError)
//		}
//	}
//	return cs.SendRequestAsync(clientId, request, genericCallback)
//}
//
//// Sends a smart charging profile to a charge point. Refer to the smart charging documentation for more information.
//func (cs *centralSystem) SetChargingProfile(clientId string, callback func(*SetChargingProfileConfirmation, error), connectorId int, chargingProfile *ChargingProfile, props ...func(request *SetChargingProfileRequest)) error {
//	request := NewSetChargingProfileRequest(connectorId, chargingProfile)
//	for _, fn := range props {
//		fn(request)
//	}
//	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
//		if confirmation != nil {
//			callback(confirmation.(*SetChargingProfileConfirmation), protoError)
//		} else {
//			callback(nil, protoError)
//		}
//	}
//	return cs.SendRequestAsync(clientId, request, genericCallback)
//}
//
//// Queries a charge point to the composite smart charging schedules and rules for a specified time interval.
//func (cs *centralSystem) GetCompositeSchedule(clientId string, callback func(*GetCompositeScheduleConfirmation, error), connectorId int, duration int, props ...func(request *GetCompositeScheduleRequest)) error {
//	request := NewGetCompositeScheduleRequest(connectorId, duration)
//	for _, fn := range props {
//		fn(request)
//	}
//	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
//		if confirmation != nil {
//			callback(confirmation.(*GetCompositeScheduleConfirmation), protoError)
//		} else {
//			callback(nil, protoError)
//		}
//	}
//	return cs.SendRequestAsync(clientId, request, genericCallback)
//}

// Registers a handler for incoming core profile messages.
func (cs *csms) SetCentralSystemCoreListener(listener CentralSystemCoreListener) {
	cs.coreListener = listener
}

// Registers a handler for incoming local authorization profile messages.
//func (cs *centralSystem) SetLocalAuthListListener(listener CentralSystemLocalAuthListListener) {
//	cs.localAuthListListener = listener
//}
//
//// Registers a handler for incoming firmware management profile messages.
//func (cs *centralSystem) SetFirmwareManagementListener(listener CentralSystemFirmwareManagementListener) {
//	cs.firmwareListener = listener
//}
//
//// Registers a handler for incoming reservation profile messages.
//func (cs *centralSystem) SetReservationListener(listener CentralSystemReservationListener) {
//	cs.reservationListener = listener
//}
//
//// Registers a handler for incoming remote trigger profile messages.
//func (cs *centralSystem) SetRemoteTriggerListener(listener CentralSystemRemoteTriggerListener) {
//	cs.remoteTriggerListener = listener
//}
//
//// Registers a handler for incoming smart charging profile messages.
//func (cs *centralSystem) SetSmartChargingListener(listener CentralSystemSmartChargingListener) {
//	cs.smartChargingListener = listener
//}

// Registers a handler for new incoming charge point connections.
func (cs *csms) SetNewChargePointHandler(handler func(chargePointId string)) {
	cs.centralSystem.SetNewChargePointHandler(handler)
}

// Registers a handler for charge point disconnections.
func (cs *csms) SetChargePointDisconnectedHandler(handler func(chargePointId string)) {
	cs.centralSystem.SetDisconnectedChargePointHandler(handler)
}

// Sends an asynchronous request to the charge point.
// The charge point will respond with a confirmation message, or with an error if the request was invalid or could not be processed.
// This result is propagated via a callback, called asynchronously.
// In case of network issues (i.e. the remote host couldn't be reached), the function returns an error directly. In this case, the callback is never called.
func (cs *csms) SendRequestAsync(clientId string, request ocpp.Request, callback func(confirmation ocpp.Confirmation, err error)) error {
	switch request.GetFeatureName() {
	case CancelReservationFeatureName, CertificateSignedFeatureName, ChangeAvailabilityFeatureName, ClearCacheFeatureName, ClearChargingProfileFeatureName, ClearDisplayFeatureName, ClearVariableMonitoringFeatureName:
		break
	//case ChangeConfigurationFeatureName, DataTransferFeatureName, GetConfigurationFeatureName, RemoteStartTransactionFeatureName, RemoteStopTransactionFeatureName, ResetFeatureName, UnlockConnectorFeatureName,
	//	GetLocalListVersionFeatureName, SendLocalListFeatureName,
	//	GetDiagnosticsFeatureName, UpdateFirmwareFeatureName,
	//	ReserveNowFeatureName,
	//	TriggerMessageFeatureName,
	//	SetChargingProfileFeatureName, ClearChargingProfileFeatureName, GetCompositeScheduleFeatureName:
	default:
		return fmt.Errorf("unsupported action %v on central system, cannot send request", request.GetFeatureName())
	}
	cs.callbacks[clientId] = callback
	err := cs.centralSystem.SendRequest(clientId, request)
	if err != nil {
		delete(cs.callbacks, clientId)
		return err
	}
	return nil
}

// Starts running the central system on the specified port and URL.
// The central system runs as a daemon and handles incoming charge point connections and messages.

// The function blocks forever, so it is suggested to wrap it in a goroutine, in case other functionality needs to be executed on the main program thread.
func (cs *csms) Start(listenPort int, listenPath string) {
	cs.centralSystem.Start(listenPort, listenPath)
}

func (cs *csms) sendResponse(chargePointId string, confirmation ocpp.Confirmation, err error, requestId string) {
	if confirmation != nil {
		err := cs.centralSystem.SendConfirmation(chargePointId, requestId, confirmation)
		if err != nil {
			//TODO: handle error somehow
			log.Print(err)
		}
	} else {
		err := cs.centralSystem.SendError(chargePointId, requestId, ocppj.ProtocolError, "Couldn't generate valid confirmation", nil)
		if err != nil {
			log.WithFields(log.Fields{
				"client":  chargePointId,
				"request": requestId,
			}).Errorf("unknown error %v while replying to message with CallError", err)
		}
	}
}

func (cs *csms) notImplementedError(chargePointId string, requestId string, action string) {
	log.Warnf("Cannot handle call %v from charge point %v. Sending CallError instead", requestId, chargePointId)
	err := cs.centralSystem.SendError(chargePointId, requestId, ocppj.NotImplemented, fmt.Sprintf("no handler for action %v implemented", action), nil)
	if err != nil {
		log.WithFields(log.Fields{
			"client":  chargePointId,
			"request": requestId,
		}).Errorf("unknown error %v while replying to message with CallError", err)
	}
}

func (cs *csms) notSupportedError(chargePointId string, requestId string, action string) {
	log.Warnf("Cannot handle call %v from charge point %v. Sending CallError instead", requestId, chargePointId)
	err := cs.centralSystem.SendError(chargePointId, requestId, ocppj.NotSupported, fmt.Sprintf("unsupported action %v on central system", action), nil)
	if err != nil {
		log.WithFields(log.Fields{
			"client":  chargePointId,
			"request": requestId,
		}).Errorf("unknown error %v while replying to message with CallError", err)
	}
}

func (cs *csms) handleIncomingRequest(chargePointId string, request ocpp.Request, requestId string, action string) {
	profile, found := cs.centralSystem.GetProfileForFeature(action)
	// Check whether action is supported and a listener for it exists
	if !found {
		cs.notImplementedError(chargePointId, requestId, action)
		return
	} else {
		switch profile.Name {
		case CoreProfileName:
			if cs.coreListener == nil {
				cs.notSupportedError(chargePointId, requestId, action)
				return
			}
		//case LocalAuthListProfileName:
		//	if cs.localAuthListListener == nil {
		//		cs.notSupportedError(chargePointId, requestId, action)
		//		return
		//	}
		//case FirmwareManagementProfileName:
		//	if cs.firmwareListener == nil {
		//		cs.notSupportedError(chargePointId, requestId, action)
		//		return
		//	}
		//case ReservationProfileName:
		//	if cs.reservationListener == nil {
		//		cs.notSupportedError(chargePointId, requestId, action)
		//		return
		//	}
		//case RemoteTriggerProfileName:
		//	if cs.remoteTriggerListener == nil {
		//		cs.notSupportedError(chargePointId, requestId, action)
		//		return
		//	}
		//case SmartChargingProfileName:
		//	if cs.smartChargingListener == nil {
		//		cs.notSupportedError(chargePointId, requestId, action)
		//		return
		//	}
		}
	}
	var confirmation ocpp.Confirmation = nil
	var err error = nil
	// Execute in separate goroutine, so the caller goroutine is available
	go func() {
		switch action {
		case BootNotificationFeatureName:
			confirmation, err = cs.coreListener.OnBootNotification(chargePointId, request.(*BootNotificationRequest))
		case AuthorizeFeatureName:
			confirmation, err = cs.coreListener.OnAuthorize(chargePointId, request.(*AuthorizeRequest))
		case ClearedChargingLimitFeatureName:
			confirmation, err = cs.coreListener.OnClearedChargingLimit(chargePointId, request.(*ClearedChargingLimitRequest))
		//case DataTransferFeatureName:
		//	confirmation, err = cs.coreListener.OnDataTransfer(chargePointId, request.(*DataTransferRequest))
		//case HeartbeatFeatureName:
		//	confirmation, err = cs.coreListener.OnHeartbeat(chargePointId, request.(*HeartbeatRequest))
		//case MeterValuesFeatureName:
		//	confirmation, err = cs.coreListener.OnMeterValues(chargePointId, request.(*MeterValuesRequest))
		//case StartTransactionFeatureName:
		//	confirmation, err = cs.coreListener.OnStartTransaction(chargePointId, request.(*StartTransactionRequest))
		//case StopTransactionFeatureName:
		//	confirmation, err = cs.coreListener.OnStopTransaction(chargePointId, request.(*StopTransactionRequest))
		//case StatusNotificationFeatureName:
		//	confirmation, err = cs.coreListener.OnStatusNotification(chargePointId, request.(*StatusNotificationRequest))
		//case DiagnosticsStatusNotificationFeatureName:
		//	confirmation, err = cs.firmwareListener.OnDiagnosticsStatusNotification(chargePointId, request.(*DiagnosticsStatusNotificationRequest))
		//case FirmwareStatusNotificationFeatureName:
		//	confirmation, err = cs.firmwareListener.OnFirmwareStatusNotification(chargePointId, request.(*FirmwareStatusNotificationRequest))
		default:
			cs.notSupportedError(chargePointId, requestId, action)
			return
		}
		cs.sendResponse(chargePointId, confirmation, err, requestId)
	}()
}

func (cs *csms) handleIncomingConfirmation(chargePointId string, confirmation ocpp.Confirmation, requestId string) {
	if callback, ok := cs.callbacks[chargePointId]; ok {
		delete(cs.callbacks, chargePointId)
		callback(confirmation, nil)
	} else {
		log.WithFields(log.Fields{
			"client":  chargePointId,
			"request": requestId,
		}).Errorf("no handler available for Call Result of type %v", confirmation.GetFeatureName())
	}
}

func (cs *csms) handleIncomingError(chargePointId string, err *ocpp.Error, details interface{}) {
	if callback, ok := cs.callbacks[chargePointId]; ok {
		delete(cs.callbacks, chargePointId)
		callback(nil, err)
	} else {
		log.WithFields(log.Fields{
			"client":  chargePointId,
			"request": err.MessageId,
		}).Errorf("no handler available for Call Error %v", err.Code)
	}
}

// Creates a new OCPP 2.0 CSMS.
//
// The dispatcher and client parameters may be omitted, in order to use a default configuration:
//   chargePoint := NewCentralSystem(nil, nil)
//
// It is recommended to use the default configuration, unless a custom networking / ocppj layer is required.
// The default dispatcher supports all OCPP 1.6 profiles out-of-the-box.
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
		centralSystem: dispatcher,
		callbacks:     map[string]func(confirmation ocpp.Confirmation, err error){}}
	cs.centralSystem.SetRequestHandler(cs.handleIncomingRequest)
	cs.centralSystem.SetConfirmationHandler(cs.handleIncomingConfirmation)
	cs.centralSystem.SetErrorHandler(cs.handleIncomingError)
	return &cs
}

func init() {
	log.New()
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetLevel(log.InfoLevel)
}
