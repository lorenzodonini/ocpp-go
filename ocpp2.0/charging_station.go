package ocpp2

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocppj"
	log "github.com/sirupsen/logrus"
)

type chargingStation struct {
	client       *ocppj.ChargePoint
	coreListener ChargePointCoreListener
	//localAuthListListener ChargePointLocalAuthListListener
	//firmwareListener      ChargePointFirmwareManagementListener
	//reservationListener   ChargePointReservationListener
	//remoteTriggerListener ChargePointRemoteTriggerListener
	//smartChargingListener ChargePointSmartChargingListener
	confirmationListener chan ocpp.Confirmation
	errorListener        chan error
}

// Sends a BootNotificationRequest to the central system, along with information about the charge point.
func (cp *chargingStation) BootNotification(reason BootReason, chargePointModel string, chargePointVendor string, props ...func(request *BootNotificationRequest)) (*BootNotificationConfirmation, error) {
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
func (cp *chargingStation) Authorize(idToken string, tokenType IdTokenType, props ...func(request *AuthorizeRequest)) (*AuthorizeConfirmation, error) {
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

func (cp *chargingStation) ClearChargingLimit(chargingLimitSource ChargingLimitSourceType, props ...func(request *ClearedChargingLimitRequest)) (*ClearedChargingLimitConfirmation, error) {
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

// Starts a custom data transfer request. Every vendor may implement their own proprietary logic for this message.
func (cp *chargingStation) DataTransfer(vendorId string, props ...func(request *DataTransferRequest)) (*DataTransferConfirmation, error) {
	request := NewDataTransferRequest(vendorId)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cp.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*DataTransferConfirmation), err
	}
}

func (cp *chargingStation) FirmwareStatusNotification(status FirmwareStatus, requestID int, props ...func(request *FirmwareStatusNotificationRequest)) (*FirmwareStatusNotificationConfirmation, error) {
	request := NewFirmwareStatusNotificationRequest(status, requestID)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cp.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*FirmwareStatusNotificationConfirmation), err
	}
}

func (cp *chargingStation) Get15118EVCertificate(schemaVersion string, exiRequest string, props ...func(request *Get15118EVCertificateRequest)) (*Get15118EVCertificateConfirmation, error) {
	request := NewGet15118EVCertificateRequest(schemaVersion, exiRequest)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cp.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*Get15118EVCertificateConfirmation), err
	}
}

func (cp *chargingStation) GetCertificateStatus(ocspRequestData OCSPRequestDataType, props ...func(request *GetCertificateStatusRequest)) (*GetCertificateStatusConfirmation, error) {
	request := NewGetCertificateStatusRequest(ocspRequestData)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cp.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*GetCertificateStatusConfirmation), err
	}
}

//
//// Notifies the central system that the charge point is still online. The central system's response is used for time synchronization purposes. It is recommended to perform this operation once every 24 hours.
//func (cp *chargingStation) Heartbeat(props ...func(request *HeartbeatRequest)) (*HeartbeatConfirmation, error) {
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
//func (cp *chargingStation) MeterValues(connectorId int, meterValues []MeterValue, props ...func(request *MeterValuesRequest)) (*MeterValuesConfirmation, error) {
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
//func (cp *chargingStation) StartTransaction(connectorId int, idTag string, meterStart int, timestamp *DateTime, props ...func(request *StartTransactionRequest)) (*StartTransactionConfirmation, error) {
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
//func (cp *chargingStation) StopTransaction(meterStop int, timestamp *DateTime, transactionId int, props ...func(request *StopTransactionRequest)) (*StopTransactionConfirmation, error) {
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
//func (cp *chargingStation) StatusNotification(connectorId int, errorCode ChargePointErrorCode, status ChargePointStatus, props ...func(request *StatusNotificationRequest)) (*StatusNotificationConfirmation, error) {
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
//func (cp *chargingStation) DiagnosticsStatusNotification(status DiagnosticsStatus, props ...func(request *DiagnosticsStatusNotificationRequest)) (*DiagnosticsStatusNotificationConfirmation, error) {
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
//func (cp *chargingStation) FirmwareStatusNotification(status FirmwareStatus, props ...func(request *FirmwareStatusNotificationRequest)) (*FirmwareStatusNotificationConfirmation, error) {
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
func (cp *chargingStation) SetChargePointCoreListener(listener ChargePointCoreListener) {
	cp.coreListener = listener
}

//// Registers a handler for incoming local authorization profile messages
//func (cp *chargingStation) SetLocalAuthListListener(listener ChargePointLocalAuthListListener) {
//	cp.localAuthListListener = listener
//}
//
//// Registers a handler for incoming firmware management profile messages
//func (cp *chargingStation) SetFirmwareManagementListener(listener ChargePointFirmwareManagementListener) {
//	cp.firmwareListener = listener
//}
//
//// Registers a handler for incoming reservation profile messages
//func (cp *chargingStation) SetReservationListener(listener ChargePointReservationListener) {
//	cp.reservationListener = listener
//}
//
//// Registers a handler for incoming remote trigger profile messages
//func (cp *chargingStation) SetRemoteTriggerListener(listener ChargePointRemoteTriggerListener) {
//	cp.remoteTriggerListener = listener
//}
//
//// Registers a handler for incoming smart charging profile messages
//func (cp *chargingStation) SetSmartChargingListener(listener ChargePointSmartChargingListener) {
//	cp.smartChargingListener = listener
//}

// Sends a request to the central system.
// The central system will respond with a confirmation, or with an error if the request was invalid or could not be processed.
// In case of network issues (i.e. the remote host couldn't be reached), the function also returns an error.
// The request is synchronous blocking.
func (cp *chargingStation) SendRequest(request ocpp.Request) (ocpp.Confirmation, error) {
	// TODO: check for supported feature
	err := cp.client.SendRequest(request)
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
func (cp *chargingStation) SendRequestAsync(request ocpp.Request, callback func(confirmation ocpp.Confirmation, err error)) error {
	switch request.GetFeatureName() {
	case AuthorizeFeatureName, BootNotificationFeatureName, ClearedChargingLimitFeatureName, DataTransferFeatureName, FirmwareStatusNotificationFeatureName, Get15118EVCertificateFeatureName, GetCertificateStatusFeatureName:
		break
	default:
		return fmt.Errorf("unsupported action %v on charge point, cannot send request", request.GetFeatureName())
	}
	err := cp.client.SendRequest(request)
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

func (cp *chargingStation) sendResponse(confirmation ocpp.Confirmation, err error, requestId string) {
	if confirmation != nil {
		err := cp.client.SendConfirmation(requestId, confirmation)
		if err != nil {
			log.WithField("request", requestId).Errorf("unknown error %v while replying to message with CallError", err)
			//TODO: handle error somehow
		}
	} else {
		err = cp.client.SendError(requestId, ocppj.ProtocolError, err.Error(), nil)
		if err != nil {
			log.WithField("request", requestId).Errorf("unknown error %v while replying to message with CallError", err)
		}
	}
}

// Connects to the central system and starts the charge point routine.
// The function doesn't block and returns right away, after having attempted to open a connection to the central system.
// If the connection couldn't be opened, an error is returned.
//
// Optional client options must be set before calling this function. Refer to NewChargingStation.
//
// No auto-reconnect logic is implemented as of now, but is planned for the future.
func (cp *chargingStation) Start(centralSystemUrl string) error {
	// TODO: implement auto-reconnect logic
	return cp.client.Start(centralSystemUrl)
}

// Stops the charge point routine, disconnecting it from the central system.
// Any pending requests are discarded.
func (cp *chargingStation) Stop() {
	cp.client.Stop()
}

func (cp *chargingStation) notImplementedError(requestId string, action string) {
	log.WithField("request", requestId).Errorf("cannot handle call from central system. Sending CallError instead")
	err := cp.client.SendError(requestId, ocppj.NotImplemented, fmt.Sprintf("no handler for action %v implemented", action), nil)
	if err != nil {
		log.WithField("request", requestId).Errorf("unknown error %v while replying to message with CallError", err)
	}
}

func (cp *chargingStation) notSupportedError(requestId string, action string) {
	log.WithField("request", requestId).Errorf("cannot handle call from central system. Sending CallError instead")
	err := cp.client.SendError(requestId, ocppj.NotSupported, fmt.Sprintf("unsupported action %v on charge point", action), nil)
	if err != nil {
		log.WithField("request", requestId).Errorf("unknown error %v while replying to message with CallError", err)
	}
}

func (cp *chargingStation) handleIncomingRequest(request ocpp.Request, requestId string, action string) {
	profile, found := cp.client.GetProfileForFeature(action)
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
	cp.client.GetProfileForFeature(action)
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
	case CostUpdatedFeatureName:
		confirmation, err = cp.coreListener.OnCostUpdated(request.(*CostUpdatedRequest))
	case CustomerInformationFeatureName:
		confirmation, err = cp.coreListener.OnCustomerInformation(request.(*CustomerInformationRequest))
	case DataTransferFeatureName:
		confirmation, err = cp.coreListener.OnDataTransfer(request.(*DataTransferRequest))
	case DeleteCertificateFeatureName:
		confirmation, err = cp.coreListener.OnDeleteCertificate(request.(*DeleteCertificateRequest))
	case GetBaseReportFeatureName:
		confirmation, err = cp.coreListener.OnGetBaseReport(request.(*GetBaseReportRequest))
	case GetChargingProfilesFeatureName:
		confirmation, err = cp.coreListener.OnGetChargingProfiles(request.(*GetChargingProfilesRequest))
	case GetCompositeScheduleFeatureName:
		confirmation, err = cp.coreListener.OnGetCompositeSchedule(request.(*GetCompositeScheduleRequest))
	case GetDisplayMessagesFeatureName:
		confirmation, err = cp.coreListener.OnGetDisplayMessages(request.(*GetDisplayMessagesRequest))
	case GetInstalledCertificateIdsFeatureName:
		confirmation, err = cp.coreListener.OnGetInstalledCertificateIds(request.(*GetInstalledCertificateIdsRequest))
	case GetLocalListVersionFeatureName:
		confirmation, err = cp.coreListener.OnGetLocalListVersion(request.(*GetLocalListVersionRequest))
	case GetLogFeatureName:
		confirmation, err = cp.coreListener.OnGetLog(request.(*GetLogRequest))
	case GetMonitoringReportFeatureName:
		confirmation, err = cp.coreListener.OnGetMonitoringReport(request.(*GetMonitoringReportRequest))
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
