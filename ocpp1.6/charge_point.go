package ocpp16

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocppj"
	log "github.com/sirupsen/logrus"
)

type chargePoint struct {
	chargePoint           *ocppj.ChargePoint
	coreListener          ChargePointCoreHandler
	localAuthListListener ChargePointLocalAuthListHandler
	firmwareListener      ChargePointFirmwareManagementHandler
	reservationListener   ChargePointReservationHandler
	remoteTriggerListener ChargePointRemoteTriggerHandler
	smartChargingListener ChargePointSmartChargingHandler
	confirmationListener  chan ocpp.Confirmation
	errorListener         chan error
}

func (cp *chargePoint) BootNotification(chargePointModel string, chargePointVendor string, props ...func(request *BootNotificationRequest)) (*BootNotificationConfirmation, error) {
	request := NewBootNotificationRequest(chargePointModel, chargePointVendor)
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

func (cp *chargePoint) Authorize(idTag string, props ...func(request *AuthorizeRequest)) (*AuthorizeConfirmation, error) {
	request := NewAuthorizationRequest(idTag)
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

func (cp *chargePoint) DataTransfer(vendorId string, props ...func(request *DataTransferRequest)) (*DataTransferConfirmation, error) {
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

func (cp *chargePoint) Heartbeat(props ...func(request *HeartbeatRequest)) (*HeartbeatConfirmation, error) {
	request := NewHeartbeatRequest()
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cp.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*HeartbeatConfirmation), err
	}
}

func (cp *chargePoint) MeterValues(connectorId int, meterValues []MeterValue, props ...func(request *MeterValuesRequest)) (*MeterValuesConfirmation, error) {
	request := NewMeterValuesRequest(connectorId, meterValues)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cp.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*MeterValuesConfirmation), err
	}
}

func (cp *chargePoint) StartTransaction(connectorId int, idTag string, meterStart int, timestamp *DateTime, props ...func(request *StartTransactionRequest)) (*StartTransactionConfirmation, error) {
	request := NewStartTransactionRequest(connectorId, idTag, meterStart, timestamp)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cp.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*StartTransactionConfirmation), err
	}
}

func (cp *chargePoint) StopTransaction(meterStop int, timestamp *DateTime, transactionId int, props ...func(request *StopTransactionRequest)) (*StopTransactionConfirmation, error) {
	request := NewStopTransactionRequest(meterStop, timestamp, transactionId)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cp.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*StopTransactionConfirmation), err
	}
}

func (cp *chargePoint) StatusNotification(connectorId int, errorCode ChargePointErrorCode, status ChargePointStatus, props ...func(request *StatusNotificationRequest)) (*StatusNotificationConfirmation, error) {
	request := NewStatusNotificationRequest(connectorId, errorCode, status)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cp.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*StatusNotificationConfirmation), err
	}
}

func (cp *chargePoint) DiagnosticsStatusNotification(status DiagnosticsStatus, props ...func(request *DiagnosticsStatusNotificationRequest)) (*DiagnosticsStatusNotificationConfirmation, error) {
	request := NewDiagnosticsStatusNotificationRequest(status)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cp.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*DiagnosticsStatusNotificationConfirmation), err
	}
}

func (cp *chargePoint) FirmwareStatusNotification(status FirmwareStatus, props ...func(request *FirmwareStatusNotificationRequest)) (*FirmwareStatusNotificationConfirmation, error) {
	request := NewFirmwareStatusNotificationRequest(status)
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

func (cp *chargePoint) SetChargePointCoreHandler(listener ChargePointCoreHandler) {
	cp.coreListener = listener
}

func (cp *chargePoint) SetLocalAuthListHandler(listener ChargePointLocalAuthListHandler) {
	cp.localAuthListListener = listener
}

func (cp *chargePoint) SetFirmwareManagementHandler(listener ChargePointFirmwareManagementHandler) {
	cp.firmwareListener = listener
}

func (cp *chargePoint) SetReservationHandler(listener ChargePointReservationHandler) {
	cp.reservationListener = listener
}

func (cp *chargePoint) SetRemoteTriggerHandler(listener ChargePointRemoteTriggerHandler) {
	cp.remoteTriggerListener = listener
}

func (cp *chargePoint) SetSmartChargingHandler(listener ChargePointSmartChargingHandler) {
	cp.smartChargingListener = listener
}

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

func (cp *chargePoint) SendRequestAsync(request ocpp.Request, callback func(confirmation ocpp.Confirmation, err error)) error {
	switch request.GetFeatureName() {
	case AuthorizeFeatureName, BootNotificationFeatureName, DataTransferFeatureName, HeartbeatFeatureName, MeterValuesFeatureName, StartTransactionFeatureName, StopTransactionFeatureName, StatusNotificationFeatureName,
		DiagnosticsStatusNotificationFeatureName, FirmwareStatusNotificationFeatureName:
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

func (cp *chargePoint) Start(centralSystemUrl string) error {
	// TODO: implement auto-reconnect logic
	return cp.chargePoint.Start(centralSystemUrl)
}

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
		case LocalAuthListProfileName:
			if cp.localAuthListListener == nil {
				cp.notSupportedError(requestId, action)
				return
			}
		case FirmwareManagementProfileName:
			if cp.firmwareListener == nil {
				cp.notSupportedError(requestId, action)
				return
			}
		case ReservationProfileName:
			if cp.reservationListener == nil {
				cp.notSupportedError(requestId, action)
				return
			}
		case RemoteTriggerProfileName:
			if cp.remoteTriggerListener == nil {
				cp.notSupportedError(requestId, action)
				return
			}
		case SmartChargingProfileName:
			if cp.smartChargingListener == nil {
				cp.notSupportedError(requestId, action)
				return
			}
		}
	}
	// Process request
	var confirmation ocpp.Confirmation = nil
	cp.chargePoint.GetProfileForFeature(action)
	var err error = nil
	switch action {
	case ChangeAvailabilityFeatureName:
		confirmation, err = cp.coreListener.OnChangeAvailability(request.(*ChangeAvailabilityRequest))
	case ChangeConfigurationFeatureName:
		confirmation, err = cp.coreListener.OnChangeConfiguration(request.(*ChangeConfigurationRequest))
	case ClearCacheFeatureName:
		confirmation, err = cp.coreListener.OnClearCache(request.(*ClearCacheRequest))
	case DataTransferFeatureName:
		confirmation, err = cp.coreListener.OnDataTransfer(request.(*DataTransferRequest))
	case GetConfigurationFeatureName:
		confirmation, err = cp.coreListener.OnGetConfiguration(request.(*GetConfigurationRequest))
	case RemoteStartTransactionFeatureName:
		confirmation, err = cp.coreListener.OnRemoteStartTransaction(request.(*RemoteStartTransactionRequest))
	case RemoteStopTransactionFeatureName:
		confirmation, err = cp.coreListener.OnRemoteStopTransaction(request.(*RemoteStopTransactionRequest))
	case ResetFeatureName:
		confirmation, err = cp.coreListener.OnReset(request.(*ResetRequest))
	case UnlockConnectorFeatureName:
		confirmation, err = cp.coreListener.OnUnlockConnector(request.(*UnlockConnectorRequest))
	case GetLocalListVersionFeatureName:
		confirmation, err = cp.localAuthListListener.OnGetLocalListVersion(request.(*GetLocalListVersionRequest))
	case SendLocalListFeatureName:
		confirmation, err = cp.localAuthListListener.OnSendLocalList(request.(*SendLocalListRequest))
	case GetDiagnosticsFeatureName:
		confirmation, err = cp.firmwareListener.OnGetDiagnostics(request.(*GetDiagnosticsRequest))
	case UpdateFirmwareFeatureName:
		confirmation, err = cp.firmwareListener.OnUpdateFirmware(request.(*UpdateFirmwareRequest))
	case ReserveNowFeatureName:
		confirmation, err = cp.reservationListener.OnReserveNow(request.(*ReserveNowRequest))
	case CancelReservationFeatureName:
		confirmation, err = cp.reservationListener.OnCancelReservation(request.(*CancelReservationRequest))
	case TriggerMessageFeatureName:
		confirmation, err = cp.remoteTriggerListener.OnTriggerMessage(request.(*TriggerMessageRequest))
	case SetChargingProfileFeatureName:
		confirmation, err = cp.smartChargingListener.OnSetChargingProfile(request.(*SetChargingProfileRequest))
	case ClearChargingProfileFeatureName:
		confirmation, err = cp.smartChargingListener.OnClearChargingProfile(request.(*ClearChargingProfileRequest))
	case GetCompositeScheduleFeatureName:
		confirmation, err = cp.smartChargingListener.OnGetCompositeSchedule(request.(*GetCompositeScheduleRequest))
	default:
		cp.notSupportedError(requestId, action)
		return
	}
	cp.sendResponse(confirmation, err, requestId)
}
