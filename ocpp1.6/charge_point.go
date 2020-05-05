package ocpp16

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/auth"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/firmware"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/remotetrigger"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/reservation"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/smartcharging"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/lorenzodonini/ocpp-go/ocppj"
	log "github.com/sirupsen/logrus"
)

type chargePoint struct {
	chargePoint           *ocppj.ChargePoint
	coreListener          core.ChargePointCoreHandler
	localAuthListListener auth.ChargePointLocalAuthListHandler
	firmwareListener      firmware.ChargePointFirmwareManagementHandler
	reservationListener   reservation.ChargePointReservationHandler
	remoteTriggerListener remotetrigger.ChargePointRemoteTriggerHandler
	smartChargingListener smartcharging.ChargePointSmartChargingHandler
	confirmationListener  chan ocpp.Response
	errorListener         chan error
}

func (cp *chargePoint) BootNotification(chargePointModel string, chargePointVendor string, props ...func(request *core.BootNotificationRequest)) (*core.BootNotificationConfirmation, error) {
	request := core.NewBootNotificationRequest(chargePointModel, chargePointVendor)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cp.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*core.BootNotificationConfirmation), err
	}
}

func (cp *chargePoint) Authorize(idTag string, props ...func(request *core.AuthorizeRequest)) (*core.AuthorizeConfirmation, error) {
	request := core.NewAuthorizationRequest(idTag)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cp.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*core.AuthorizeConfirmation), err
	}
}

func (cp *chargePoint) DataTransfer(vendorId string, props ...func(request *core.DataTransferRequest)) (*core.DataTransferConfirmation, error) {
	request := core.NewDataTransferRequest(vendorId)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cp.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*core.DataTransferConfirmation), err
	}
}

func (cp *chargePoint) Heartbeat(props ...func(request *core.HeartbeatRequest)) (*core.HeartbeatConfirmation, error) {
	request := core.NewHeartbeatRequest()
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cp.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*core.HeartbeatConfirmation), err
	}
}

func (cp *chargePoint) MeterValues(connectorId int, meterValues []types.MeterValue, props ...func(request *core.MeterValuesRequest)) (*core.MeterValuesConfirmation, error) {
	request := core.NewMeterValuesRequest(connectorId, meterValues)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cp.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*core.MeterValuesConfirmation), err
	}
}

func (cp *chargePoint) StartTransaction(connectorId int, idTag string, meterStart int, timestamp *types.DateTime, props ...func(request *core.StartTransactionRequest)) (*core.StartTransactionConfirmation, error) {
	request := core.NewStartTransactionRequest(connectorId, idTag, meterStart, timestamp)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cp.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*core.StartTransactionConfirmation), err
	}
}

func (cp *chargePoint) StopTransaction(meterStop int, timestamp *types.DateTime, transactionId int, props ...func(request *core.StopTransactionRequest)) (*core.StopTransactionConfirmation, error) {
	request := core.NewStopTransactionRequest(meterStop, timestamp, transactionId)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cp.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*core.StopTransactionConfirmation), err
	}
}

func (cp *chargePoint) StatusNotification(connectorId int, errorCode core.ChargePointErrorCode, status core.ChargePointStatus, props ...func(request *core.StatusNotificationRequest)) (*core.StatusNotificationConfirmation, error) {
	request := core.NewStatusNotificationRequest(connectorId, errorCode, status)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cp.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*core.StatusNotificationConfirmation), err
	}
}

func (cp *chargePoint) DiagnosticsStatusNotification(status firmware.DiagnosticsStatus, props ...func(request *firmware.DiagnosticsStatusNotificationRequest)) (*firmware.DiagnosticsStatusNotificationConfirmation, error) {
	request := firmware.NewDiagnosticsStatusNotificationRequest(status)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cp.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*firmware.DiagnosticsStatusNotificationConfirmation), err
	}
}

func (cp *chargePoint) FirmwareStatusNotification(status firmware.FirmwareStatus, props ...func(request *firmware.FirmwareStatusNotificationRequest)) (*firmware.FirmwareStatusNotificationConfirmation, error) {
	request := firmware.NewFirmwareStatusNotificationRequest(status)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cp.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*firmware.FirmwareStatusNotificationConfirmation), err
	}
}

func (cp *chargePoint) SetChargePointCoreHandler(listener core.ChargePointCoreHandler) {
	cp.coreListener = listener
}

func (cp *chargePoint) SetLocalAuthListHandler(listener auth.ChargePointLocalAuthListHandler) {
	cp.localAuthListListener = listener
}

func (cp *chargePoint) SetFirmwareManagementHandler(listener firmware.ChargePointFirmwareManagementHandler) {
	cp.firmwareListener = listener
}

func (cp *chargePoint) SetReservationHandler(listener reservation.ChargePointReservationHandler) {
	cp.reservationListener = listener
}

func (cp *chargePoint) SetRemoteTriggerHandler(listener remotetrigger.ChargePointRemoteTriggerHandler) {
	cp.remoteTriggerListener = listener
}

func (cp *chargePoint) SetSmartChargingHandler(listener smartcharging.ChargePointSmartChargingHandler) {
	cp.smartChargingListener = listener
}

func (cp *chargePoint) SendRequest(request ocpp.Request) (ocpp.Response, error) {
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

func (cp *chargePoint) SendRequestAsync(request ocpp.Request, callback func(confirmation ocpp.Response, err error)) error {
	switch request.GetFeatureName() {
	case core.AuthorizeFeatureName, core.BootNotificationFeatureName, core.DataTransferFeatureName, core.HeartbeatFeatureName, core.MeterValuesFeatureName, core.StartTransactionFeatureName, core.StopTransactionFeatureName, core.StatusNotificationFeatureName,
		firmware.DiagnosticsStatusNotificationFeatureName, firmware.FirmwareStatusNotificationFeatureName:
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

func (cp *chargePoint) sendResponse(confirmation ocpp.Response, err error, requestId string) {
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
		case core.ProfileName:
			if cp.coreListener == nil {
				cp.notSupportedError(requestId, action)
				return
			}
		case auth.ProfileName:
			if cp.localAuthListListener == nil {
				cp.notSupportedError(requestId, action)
				return
			}
		case firmware.ProfileName:
			if cp.firmwareListener == nil {
				cp.notSupportedError(requestId, action)
				return
			}
		case reservation.ProfileName:
			if cp.reservationListener == nil {
				cp.notSupportedError(requestId, action)
				return
			}
		case remotetrigger.ProfileName:
			if cp.remoteTriggerListener == nil {
				cp.notSupportedError(requestId, action)
				return
			}
		case smartcharging.ProfileName:
			if cp.smartChargingListener == nil {
				cp.notSupportedError(requestId, action)
				return
			}
		}
	}
	// Process request
	var confirmation ocpp.Response = nil
	cp.chargePoint.GetProfileForFeature(action)
	var err error = nil
	switch action {
	case core.ChangeAvailabilityFeatureName:
		confirmation, err = cp.coreListener.OnChangeAvailability(request.(*core.ChangeAvailabilityRequest))
	case core.ChangeConfigurationFeatureName:
		confirmation, err = cp.coreListener.OnChangeConfiguration(request.(*core.ChangeConfigurationRequest))
	case core.ClearCacheFeatureName:
		confirmation, err = cp.coreListener.OnClearCache(request.(*core.ClearCacheRequest))
	case core.DataTransferFeatureName:
		confirmation, err = cp.coreListener.OnDataTransfer(request.(*core.DataTransferRequest))
	case core.GetConfigurationFeatureName:
		confirmation, err = cp.coreListener.OnGetConfiguration(request.(*core.GetConfigurationRequest))
	case core.RemoteStartTransactionFeatureName:
		confirmation, err = cp.coreListener.OnRemoteStartTransaction(request.(*core.RemoteStartTransactionRequest))
	case core.RemoteStopTransactionFeatureName:
		confirmation, err = cp.coreListener.OnRemoteStopTransaction(request.(*core.RemoteStopTransactionRequest))
	case core.ResetFeatureName:
		confirmation, err = cp.coreListener.OnReset(request.(*core.ResetRequest))
	case core.UnlockConnectorFeatureName:
		confirmation, err = cp.coreListener.OnUnlockConnector(request.(*core.UnlockConnectorRequest))
	case auth.GetLocalListVersionFeatureName:
		confirmation, err = cp.localAuthListListener.OnGetLocalListVersion(request.(*auth.GetLocalListVersionRequest))
	case auth.SendLocalListFeatureName:
		confirmation, err = cp.localAuthListListener.OnSendLocalList(request.(*auth.SendLocalListRequest))
	case firmware.GetDiagnosticsFeatureName:
		confirmation, err = cp.firmwareListener.OnGetDiagnostics(request.(*firmware.GetDiagnosticsRequest))
	case firmware.UpdateFirmwareFeatureName:
		confirmation, err = cp.firmwareListener.OnUpdateFirmware(request.(*firmware.UpdateFirmwareRequest))
	case reservation.ReserveNowFeatureName:
		confirmation, err = cp.reservationListener.OnReserveNow(request.(*reservation.ReserveNowRequest))
	case reservation.CancelReservationFeatureName:
		confirmation, err = cp.reservationListener.OnCancelReservation(request.(*reservation.CancelReservationRequest))
	case remotetrigger.TriggerMessageFeatureName:
		confirmation, err = cp.remoteTriggerListener.OnTriggerMessage(request.(*remotetrigger.TriggerMessageRequest))
	case smartcharging.SetChargingProfileFeatureName:
		confirmation, err = cp.smartChargingListener.OnSetChargingProfile(request.(*smartcharging.SetChargingProfileRequest))
	case smartcharging.ClearChargingProfileFeatureName:
		confirmation, err = cp.smartChargingListener.OnClearChargingProfile(request.(*smartcharging.ClearChargingProfileRequest))
	case smartcharging.GetCompositeScheduleFeatureName:
		confirmation, err = cp.smartChargingListener.OnGetCompositeSchedule(request.(*smartcharging.GetCompositeScheduleRequest))
	default:
		cp.notSupportedError(requestId, action)
		return
	}
	cp.sendResponse(confirmation, err, requestId)
}
