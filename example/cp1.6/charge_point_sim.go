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
	"os"
	"time"
)

type ConnectorInfo struct {
	status             core.ChargePointStatus
	availability       core.AvailabilityType
	currentTransaction int
	currentReservation int
}

type ChargePointHandler struct {
	status               core.ChargePointStatus
	connectors           map[int]*ConnectorInfo
	errorCode            core.ChargePointErrorCode
	configuration        map[string]core.ConfigurationKey
	meterValue           int
	localAuthList        []localauth.AuthorizationData
	localAuthListVersion int
}

var chargePoint ocpp16.ChargePoint

// Core profile callbacks
func (handler *ChargePointHandler) OnChangeAvailability(request *core.ChangeAvailabilityRequest) (confirmation *core.ChangeAvailabilityConfirmation, err error) {
	handler.connectors[request.ConnectorId].availability = request.Type
	return core.NewChangeAvailabilityConfirmation(core.AvailabilityStatusAccepted), nil
}

func (handler *ChargePointHandler) OnChangeConfiguration(request *core.ChangeConfigurationRequest) (confirmation *core.ChangeConfigurationConfirmation, err error) {
	configKey, ok := handler.configuration[request.Key]
	if !ok {
		return core.NewChangeConfigurationConfirmation(core.ConfigurationStatusNotSupported), nil
	} else if configKey.Readonly {
		return core.NewChangeConfigurationConfirmation(core.ConfigurationStatusRejected), nil
	}
	configKey.Value = request.Value
	handler.configuration[request.Key] = configKey
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

// Local authorization list profile callbacks
func (handler *ChargePointHandler) OnGetLocalListVersion(request *localauth.GetLocalListVersionRequest) (confirmation *localauth.GetLocalListVersionConfirmation, err error) {
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

// Firmware management profile callbacks
func (handler *ChargePointHandler) OnGetDiagnostics(request *firmware.GetDiagnosticsRequest) (confirmation *firmware.GetDiagnosticsConfirmation, err error) {
	return firmware.NewGetDiagnosticsConfirmation(), nil
	//TODO: perform diagnostics upload out-of-band
}

func (handler *ChargePointHandler) OnUpdateFirmware(request *firmware.UpdateFirmwareRequest) (confirmation *firmware.UpdateFirmwareConfirmation, err error) {
	return firmware.NewUpdateFirmwareConfirmation(), nil
	//TODO: download new firmware out-of-band
}

// Remote trigger profile callbacks
func (handler *ChargePointHandler) OnTriggerMessage(request *remotetrigger.TriggerMessageRequest) (confirmation *remotetrigger.TriggerMessageConfirmation, err error) {
	switch request.RequestedMessage {
	case core.BootNotificationFeatureName:
		//TODO: schedule boot notification message
		break
	case firmware.DiagnosticsStatusNotificationFeatureName:
		//TODO: schedule diagnostics status notification message
		break
	case firmware.FirmwareStatusNotificationFeatureName:
		//TODO: schedule firmware status notification message
		break
	case core.HeartbeatFeatureName:
		//TODO: schedule heartbeat message
		break
	case core.MeterValuesFeatureName:
		//TODO: schedule meter values message
		break
		//TODO: schedule status notification message
	case core.StatusNotificationFeatureName:
		break
	default:
		return remotetrigger.NewTriggerMessageConfirmation(remotetrigger.TriggerMessageStatusNotImplemented), nil
	}
	return remotetrigger.NewTriggerMessageConfirmation(remotetrigger.TriggerMessageStatusAccepted), nil
}

// Reservation profile callbacks
func (handler *ChargePointHandler) OnReserveNow(request *reservation.ReserveNowRequest) (confirmation *reservation.ReserveNowConfirmation, err error) {
	connector := handler.connectors[request.ConnectorId]
	if connector == nil {
		return reservation.NewReserveNowConfirmation(reservation.ReservationStatusUnavailable), nil
	} else if connector.status != core.ChargePointStatusAvailable {
		return reservation.NewReserveNowConfirmation(reservation.ReservationStatusOccupied), nil
	}
	connector.currentReservation = request.ReservationId
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
			return reservation.NewCancelReservationConfirmation(reservation.CancelReservationStatusAccepted), nil
		}
	}
	return reservation.NewCancelReservationConfirmation(reservation.CancelReservationStatusRejected), nil
}

// Smart charging profile callbacks
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

func updateStatus(stateHandler *ChargePointHandler, connector int, status core.ChargePointStatus) {
	if connector == 0 {
		stateHandler.status = status
	} else {
		stateHandler.connectors[connector].status = status
	}
	statusConfirmation, err := chargePoint.StatusNotification(connector, core.NoError, status)
	checkError(err)
	if connector == 0 {
		logDefault(statusConfirmation.GetFeatureName()).Infof("status for all connectors updated to %v", status)
	} else {
		logDefault(statusConfirmation.GetFeatureName()).Infof("status for connector %v updated to %v", connector, status)
	}
}

func exampleRoutine(chargePoint ocpp16.ChargePoint, stateHandler *ChargePointHandler) {
	dummyClientIdTag := "12345"
	chargingConnector := 1
	// Boot
	bootConf, err := chargePoint.BootNotification("model1", "vendor1")
	checkError(err)
	logDefault(bootConf.GetFeatureName()).Infof("status: %v, interval: %v, current time: %v", bootConf.Status, bootConf.Interval, bootConf.CurrentTime.String())
	// Notify connector status
	updateStatus(stateHandler, 0, core.ChargePointStatusAvailable)
	// Wait for some time ...
	time.Sleep(5 * time.Second)
	// Simulate charging for connector 1
	authConf, err := chargePoint.Authorize(dummyClientIdTag)
	checkError(err)
	logDefault(authConf.GetFeatureName()).Infof("status: %v %v", authConf.IdTagInfo.Status, getExpiryDate(authConf.IdTagInfo))
	// Update connector status
	updateStatus(stateHandler, chargingConnector, core.ChargePointStatusPreparing)
	// Start transaction
	startConf, err := chargePoint.StartTransaction(chargingConnector, dummyClientIdTag, stateHandler.meterValue, types.NewDateTime(time.Now()))
	checkError(err)
	logDefault(startConf.GetFeatureName()).Infof("status: %v, transaction %v %v", startConf.IdTagInfo.Status, startConf.TransactionId, getExpiryDate(startConf.IdTagInfo))
	stateHandler.connectors[chargingConnector].currentTransaction = startConf.TransactionId
	// Update connector status
	updateStatus(stateHandler, chargingConnector, core.ChargePointStatusCharging)
	// Periodically send meter values
	for i := 0; i < 5; i++ {
		time.Sleep(5 * time.Second)
		stateHandler.meterValue += 10
		sampledValue := types.SampledValue{Value: fmt.Sprintf("%v", stateHandler.meterValue), Unit: types.UnitOfMeasureWh, Format: types.ValueFormatRaw, Measurand: types.MeasurandEnergyActiveExportRegister, Context: types.ReadingContextSamplePeriodic, Location: types.LocationOutlet}
		meterValue := types.MeterValue{Timestamp: types.NewDateTime(time.Now()), SampledValue: []types.SampledValue{sampledValue}}
		meterConf, err := chargePoint.MeterValues(chargingConnector, []types.MeterValue{meterValue})
		checkError(err)
		logDefault(meterConf.GetFeatureName()).Infof("sent updated %v", sampledValue.Measurand)
	}
	stateHandler.meterValue += 2
	// Stop charging for connector 1
	updateStatus(stateHandler, chargingConnector, core.ChargePointStatusFinishing)
	stopConf, err := chargePoint.StopTransaction(stateHandler.meterValue, types.NewDateTime(time.Now()), startConf.TransactionId, func(request *core.StopTransactionRequest) {
		sampledValue := types.SampledValue{Value: fmt.Sprintf("%v", stateHandler.meterValue), Unit: types.UnitOfMeasureWh, Format: types.ValueFormatRaw, Measurand: types.MeasurandEnergyActiveExportRegister, Context: types.ReadingContextSamplePeriodic, Location: types.LocationOutlet}
		meterValue := types.MeterValue{Timestamp: types.NewDateTime(time.Now()), SampledValue: []types.SampledValue{sampledValue}}
		request.TransactionData = []types.MeterValue{meterValue}
		request.Reason = core.ReasonEVDisconnected
	})
	checkError(err)
	logDefault(stopConf.GetFeatureName()).Infof("transaction %v stopped", startConf.TransactionId)
	// Update connector status
	updateStatus(stateHandler, chargingConnector, core.ChargePointStatusAvailable)
}

// Start function
func main() {
	// Parse arguments from env variables
	id, ok := os.LookupEnv("CLIENT_ID")
	if !ok {
		log.Print("Usage:\n\tocppClientId\n\tocppServerUrl")
		return
	}
	csUrl, ok := os.LookupEnv("CENTRAL_SYSTEM_URL")
	if !ok {
		log.Print("Usage:\n\tocppClientId\n\tocppServerUrl")
		return
	}
	// Create a default OCPP 1.6 charge point
	chargePoint = ocpp16.NewChargePoint(id, nil, nil)
	// Set a handler for all callback functions
	connectors := map[int]*ConnectorInfo{
		1: {status: core.ChargePointStatusAvailable, availability: core.AvailabilityTypeOperative, currentTransaction: 0},
	}
	handler := &ChargePointHandler{
		status:               core.ChargePointStatusAvailable,
		connectors:           connectors,
		configuration:        map[string]core.ConfigurationKey{},
		errorCode:            core.NoError,
		localAuthList:        []localauth.AuthorizationData{},
		localAuthListVersion: 0}
	chargePoint.SetChargePointCoreHandler(handler)
	// Connects to central system
	err := chargePoint.Start(csUrl)
	if err != nil {
		log.Println(err)
	} else {
		log.Infof("connected to central system at %v", csUrl)
		exampleRoutine(chargePoint, handler)
		// Disconnect
		chargePoint.Stop()
		log.Infof("disconnected from central system")
	}
}

func init() {
	log.SetLevel(log.InfoLevel)
}

// Utility functions
func logDefault(feature string) *log.Entry {
	return log.WithField("message", feature)
}
