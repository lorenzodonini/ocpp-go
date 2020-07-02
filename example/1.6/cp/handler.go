package main

import (
	"fmt"
	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/firmware"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/remotetrigger"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/reservation"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/smartcharging"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
)

// ConnectorInfo contains some simple state about a single connector.
type ConnectorInfo struct {
	status             core.ChargePointStatus
	availability       core.AvailabilityType
	currentTransaction int
	currentReservation int
}

// ChargePointHandler contains some simple state that a charge point needs to keep.
// In production this will typically be replaced by database/API calls.
type ChargePointHandler struct {
	status               core.ChargePointStatus
	connectors           map[int]*ConnectorInfo
	errorCode            core.ChargePointErrorCode
	configuration        ConfigMap
	meterValue           int
	localAuthList        []localauth.AuthorizationData
	localAuthListVersion int
}

var asyncRequestChan chan func()

var chargePoint ocpp16.ChargePoint

func (handler *ChargePointHandler) isValidConnectorID(ID int) bool {
	_, ok := handler.connectors[ID]
	return ok || ID == 0
}

// ------------- Core profile callbacks -------------

func (handler *ChargePointHandler) OnChangeAvailability(request *core.ChangeAvailabilityRequest) (confirmation *core.ChangeAvailabilityConfirmation, err error) {
	handler.connectors[request.ConnectorId].availability = request.Type
	return core.NewChangeAvailabilityConfirmation(core.AvailabilityStatusAccepted), nil
}

func (handler *ChargePointHandler) OnChangeConfiguration(request *core.ChangeConfigurationRequest) (confirmation *core.ChangeConfigurationConfirmation, err error) {
	configKey, ok := handler.configuration[request.Key]
	if !ok {
		logDefault(request.GetFeatureName()).Infof("couldn't change configuration for unsupported parameter %v", configKey.Key)
		return core.NewChangeConfigurationConfirmation(core.ConfigurationStatusNotSupported), nil
	} else if configKey.Readonly {
		logDefault(request.GetFeatureName()).Infof("couldn't change configuration for readonly parameter %v", configKey.Key)
		return core.NewChangeConfigurationConfirmation(core.ConfigurationStatusRejected), nil
	}
	configKey.Value = request.Value
	handler.configuration[request.Key] = configKey
	logDefault(request.GetFeatureName()).Infof("changed configuration for parameter %v to %v", configKey.Key, configKey.Value)
	return core.NewChangeConfigurationConfirmation(core.ConfigurationStatusAccepted), nil
}

func (handler *ChargePointHandler) OnClearCache(request *core.ClearCacheRequest) (confirmation *core.ClearCacheConfirmation, err error) {
	return core.NewClearCacheConfirmation(core.ClearCacheStatusAccepted), nil
}

func (handler *ChargePointHandler) OnDataTransfer(request *core.DataTransferRequest) (confirmation *core.DataTransferConfirmation, err error) {
	logDefault(request.GetFeatureName()).Infof("data transfer [Vendor: %v Message: %v]: %v", request.VendorId, request.MessageId, request.Data)
	return core.NewDataTransferConfirmation(core.DataTransferStatusAccepted), nil
}

func (handler *ChargePointHandler) OnGetConfiguration(request *core.GetConfigurationRequest) (confirmation *core.GetConfigurationConfirmation, err error) {
	var resultKeys []core.ConfigurationKey
	var unknownKeys []string
	for _, key := range request.Key {
		configKey, ok := handler.configuration[key]
		if !ok {
			unknownKeys = append(unknownKeys, configKey.Value)
		} else {
			resultKeys = append(resultKeys, configKey)
		}
	}
	conf := core.NewGetConfigurationConfirmation(resultKeys)
	conf.UnknownKey = unknownKeys
	return conf, nil
}

func (handler *ChargePointHandler) OnRemoteStartTransaction(request *core.RemoteStartTransactionRequest) (confirmation *core.RemoteStartTransactionConfirmation, err error) {
	connector, ok := handler.connectors[request.ConnectorId]
	if !ok {
		return core.NewRemoteStartTransactionConfirmation(types.RemoteStartStopStatusRejected), nil
	} else if connector.availability != core.AvailabilityTypeOperative || connector.status != core.ChargePointStatusAvailable || connector.currentTransaction > 0 {
		return core.NewRemoteStartTransactionConfirmation(types.RemoteStartStopStatusRejected), nil
	}
	logDefault(request.GetFeatureName()).Infof("started transaction %v on connector %v", connector.currentTransaction, request.ConnectorId)
	connector.currentTransaction = request.ConnectorId
	return core.NewRemoteStartTransactionConfirmation(types.RemoteStartStopStatusAccepted), nil
}

func (handler *ChargePointHandler) OnRemoteStopTransaction(request *core.RemoteStopTransactionRequest) (confirmation *core.RemoteStopTransactionConfirmation, err error) {
	for key, val := range handler.connectors {
		if val.currentTransaction == request.TransactionId {
			logDefault(request.GetFeatureName()).Infof("stopped transaction %v on connector %v", val.currentTransaction, key)
			val.currentTransaction = 0
			val.currentReservation = 0
			val.status = core.ChargePointStatusAvailable
			return core.NewRemoteStopTransactionConfirmation(types.RemoteStartStopStatusAccepted), nil
		}
	}
	return core.NewRemoteStopTransactionConfirmation(types.RemoteStartStopStatusRejected), nil
}

func (handler *ChargePointHandler) OnReset(request *core.ResetRequest) (confirmation *core.ResetConfirmation, err error) {
	//TODO: stop all ongoing transactions
	return core.NewResetConfirmation(core.ResetStatusAccepted), nil
}

func (handler *ChargePointHandler) OnUnlockConnector(request *core.UnlockConnectorRequest) (confirmation *core.UnlockConnectorConfirmation, err error) {
	_, ok := handler.connectors[request.ConnectorId]
	if !ok {
		return core.NewUnlockConnectorConfirmation(core.UnlockStatusNotSupported), nil
	}
	return core.NewUnlockConnectorConfirmation(core.UnlockStatusUnlocked), nil
}

// ------------- Local authorization list profile callbacks -------------

func (handler *ChargePointHandler) OnGetLocalListVersion(request *localauth.GetLocalListVersionRequest) (confirmation *localauth.GetLocalListVersionConfirmation, err error) {
	logDefault(request.GetFeatureName()).Infof("returning current local list version: %v", handler.localAuthListVersion)
	return localauth.NewGetLocalListVersionConfirmation(handler.localAuthListVersion), nil
}

func (handler *ChargePointHandler) OnSendLocalList(request *localauth.SendLocalListRequest) (confirmation *localauth.SendLocalListConfirmation, err error) {
	if request.ListVersion <= handler.localAuthListVersion {
		return localauth.NewSendLocalListConfirmation(localauth.UpdateStatusVersionMismatch), nil
	}
	if request.UpdateType == localauth.UpdateTypeFull {
		handler.localAuthList = request.LocalAuthorizationList
		handler.localAuthListVersion = request.ListVersion
	} else if request.UpdateType == localauth.UpdateTypeDifferential {
		handler.localAuthList = append(handler.localAuthList, request.LocalAuthorizationList...)
		handler.localAuthListVersion = request.ListVersion
	}
	return localauth.NewSendLocalListConfirmation(localauth.UpdateStatusAccepted), nil
}

// ------------- Firmware management profile callbacks -------------

func (handler *ChargePointHandler) OnGetDiagnostics(request *firmware.GetDiagnosticsRequest) (confirmation *firmware.GetDiagnosticsConfirmation, err error) {
	return firmware.NewGetDiagnosticsConfirmation(), nil
	//TODO: perform diagnostics upload out-of-band
}

func (handler *ChargePointHandler) OnUpdateFirmware(request *firmware.UpdateFirmwareRequest) (confirmation *firmware.UpdateFirmwareConfirmation, err error) {
	return firmware.NewUpdateFirmwareConfirmation(), nil
	//TODO: download new firmware out-of-band
}

// ------------- Remote trigger profile callbacks -------------

func (handler *ChargePointHandler) OnTriggerMessage(request *remotetrigger.TriggerMessageRequest) (confirmation *remotetrigger.TriggerMessageConfirmation, err error) {
	logDefault(request.GetFeatureName()).Infof("received trigger for %v", request.RequestedMessage)
	status := remotetrigger.TriggerMessageStatusRejected
	switch request.RequestedMessage {
	case core.BootNotificationFeatureName:
		//TODO: schedule boot notification message
		break
	case firmware.DiagnosticsStatusNotificationFeatureName:
		// Schedule diagnostics status notification request
		fn := func() {
			_, e := chargePoint.DiagnosticsStatusNotification(firmware.DiagnosticsStatusIdle)
			checkError(e)
			logDefault(firmware.DiagnosticsStatusNotificationFeatureName).Info("diagnostics status notified")
		}
		scheduleAsyncRequest(fn)
		status = remotetrigger.TriggerMessageStatusAccepted
	case firmware.FirmwareStatusNotificationFeatureName:
		//TODO: schedule firmware status notification message
		break
	case core.HeartbeatFeatureName:
		// Schedule heartbeat request
		fn := func() {
			conf, e := chargePoint.Heartbeat()
			checkError(e)
			logDefault(core.HeartbeatFeatureName).Infof("clock synchronized: %v", conf.CurrentTime.FormatTimestamp())
		}
		scheduleAsyncRequest(fn)
		status = remotetrigger.TriggerMessageStatusAccepted
	case core.MeterValuesFeatureName:
		//TODO: schedule meter values message
		break
	case core.StatusNotificationFeatureName:
		connectorID := request.ConnectorId
		// Check if requested connector is valid and status can be retrieved
		if !handler.isValidConnectorID(connectorID) {
			logDefault(request.GetFeatureName()).Errorf("cannot trigger %v: requested invalid connector %v", request.RequestedMessage, request.ConnectorId)
			return remotetrigger.NewTriggerMessageConfirmation(remotetrigger.TriggerMessageStatusRejected), nil
		}
		// Schedule status notification request
		fn := func() {
			status := handler.status
			if c, ok := handler.connectors[request.ConnectorId]; ok {
				status = c.status
			}
			statusConfirmation, err := chargePoint.StatusNotification(connectorID, handler.errorCode, status)
			checkError(err)
			logDefault(statusConfirmation.GetFeatureName()).Infof("status for connector %v sent: %v", connectorID, status)
		}
		scheduleAsyncRequest(fn)
		status = remotetrigger.TriggerMessageStatusAccepted
	default:
		return remotetrigger.NewTriggerMessageConfirmation(remotetrigger.TriggerMessageStatusNotImplemented), nil
	}
	return remotetrigger.NewTriggerMessageConfirmation(status), nil
}

// ------------- Reservation profile callbacks -------------

func (handler *ChargePointHandler) OnReserveNow(request *reservation.ReserveNowRequest) (confirmation *reservation.ReserveNowConfirmation, err error) {
	connector := handler.connectors[request.ConnectorId]
	if connector == nil {
		return reservation.NewReserveNowConfirmation(reservation.ReservationStatusUnavailable), nil
	} else if connector.status != core.ChargePointStatusAvailable {
		return reservation.NewReserveNowConfirmation(reservation.ReservationStatusOccupied), nil
	}
	connector.currentReservation = request.ReservationId
	logDefault(request.GetFeatureName()).Infof("reservation %v for connector %v accepted", request.ReservationId, request.ConnectorId)
	go updateStatus(handler, request.ConnectorId, core.ChargePointStatusReserved)
	// TODO: automatically remove reservation after expiryDate
	return reservation.NewReserveNowConfirmation(reservation.ReservationStatusAccepted), nil
}

func (handler *ChargePointHandler) OnCancelReservation(request *reservation.CancelReservationRequest) (confirmation *reservation.CancelReservationConfirmation, err error) {
	for k, v := range handler.connectors {
		if v.currentReservation == request.ReservationId {
			v.currentReservation = 0
			if v.status == core.ChargePointStatusReserved {
				go updateStatus(handler, k, core.ChargePointStatusAvailable)
			}
			logDefault(request.GetFeatureName()).Infof("reservation %v for connector %v canceled", request.ReservationId, k)
			return reservation.NewCancelReservationConfirmation(reservation.CancelReservationStatusAccepted), nil
		}
	}
	logDefault(request.GetFeatureName()).Infof("couldn't cancel reservation %v: reservation not found!", request.ReservationId)
	return reservation.NewCancelReservationConfirmation(reservation.CancelReservationStatusRejected), nil
}

// ------------- Smart charging profile callbacks -------------

func (handler *ChargePointHandler) OnSetChargingProfile(request *smartcharging.SetChargingProfileRequest) (confirmation *smartcharging.SetChargingProfileConfirmation, err error) {
	//TODO: handle logic
	return smartcharging.NewSetChargingProfileConfirmation(smartcharging.ChargingProfileStatusNotImplemented), nil
}

func (handler *ChargePointHandler) OnClearChargingProfile(request *smartcharging.ClearChargingProfileRequest) (confirmation *smartcharging.ClearChargingProfileConfirmation, err error) {
	//TODO: handle logic
	return smartcharging.NewClearChargingProfileConfirmation(smartcharging.ClearChargingProfileStatusUnknown), nil
}

func (handler *ChargePointHandler) OnGetCompositeSchedule(request *smartcharging.GetCompositeScheduleRequest) (confirmation *smartcharging.GetCompositeScheduleConfirmation, err error) {
	//TODO: handle logic
	return smartcharging.NewGetCompositeScheduleConfirmation(smartcharging.GetCompositeScheduleStatusRejected), nil
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getExpiryDate(info *types.IdTagInfo) string {
	if info.ExpiryDate != nil {
		return fmt.Sprintf("authorized until %v", info.ExpiryDate.String())
	}
	return ""
}

func updateStatus(stateHandler *ChargePointHandler, connector int, status core.ChargePointStatus, props ...func(request *core.StatusNotificationRequest)) {
	if connector == 0 {
		stateHandler.status = status
	} else {
		stateHandler.connectors[connector].status = status
	}
	statusConfirmation, err := chargePoint.StatusNotification(connector, stateHandler.errorCode, status, props...)
	checkError(err)
	if connector == 0 {
		logDefault(statusConfirmation.GetFeatureName()).Infof("status for all connectors updated to %v", status)
	} else {
		logDefault(statusConfirmation.GetFeatureName()).Infof("status for connector %v updated to %v", connector, status)
	}
}
