package main

import (
	"fmt"
	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

type ConnectorInfo struct {
	status             ocpp16.ChargePointStatus
	availability       ocpp16.AvailabilityType
	currentTransaction int
	currentReservation int
}

type ChargePointHandler struct {
	status               ocpp16.ChargePointStatus
	connectors           map[int]*ConnectorInfo
	errorCode            ocpp16.ChargePointErrorCode
	configuration        map[string]ocpp16.ConfigurationKey
	meterValue           int
	localAuthList        []ocpp16.AuthorizationData
	localAuthListVersion int
}

// Core profile callbacks
func (handler *ChargePointHandler) OnChangeAvailability(request *ocpp16.ChangeAvailabilityRequest) (confirmation *ocpp16.ChangeAvailabilityConfirmation, err error) {
	handler.connectors[request.ConnectorId].availability = request.Type
	return ocpp16.NewChangeAvailabilityConfirmation(ocpp16.AvailabilityStatusAccepted), nil
}

func (handler *ChargePointHandler) OnChangeConfiguration(request *ocpp16.ChangeConfigurationRequest) (confirmation *ocpp16.ChangeConfigurationConfirmation, err error) {
	configKey, ok := handler.configuration[request.Key]
	if !ok {
		return ocpp16.NewChangeConfigurationConfirmation(ocpp16.ConfigurationStatusNotSupported), nil
	} else if configKey.Readonly {
		return ocpp16.NewChangeConfigurationConfirmation(ocpp16.ConfigurationStatusRejected), nil
	}
	configKey.Value = request.Value
	handler.configuration[request.Key] = configKey
	return ocpp16.NewChangeConfigurationConfirmation(ocpp16.ConfigurationStatusAccepted), nil
}

func (handler *ChargePointHandler) OnClearCache(request *ocpp16.ClearCacheRequest) (confirmation *ocpp16.ClearCacheConfirmation, err error) {
	return ocpp16.NewClearCacheConfirmation(ocpp16.ClearCacheStatusAccepted), nil
}

func (handler *ChargePointHandler) OnDataTransfer(request *ocpp16.DataTransferRequest) (confirmation *ocpp16.DataTransferConfirmation, err error) {
	logDefault(request.GetFeatureName()).Infof("data transfer [Vendor: %v Message: %v]: %v", request.VendorId, request.MessageId, request.Data)
	return ocpp16.NewDataTransferConfirmation(ocpp16.DataTransferStatusAccepted), nil
}

func (handler *ChargePointHandler) OnGetConfiguration(request *ocpp16.GetConfigurationRequest) (confirmation *ocpp16.GetConfigurationConfirmation, err error) {
	var resultKeys []ocpp16.ConfigurationKey
	var unknownKeys []string
	for _, key := range request.Key {
		configKey, ok := handler.configuration[key]
		if !ok {
			unknownKeys = append(unknownKeys, configKey.Value)
		} else {
			resultKeys = append(resultKeys, configKey)
		}
	}
	conf := ocpp16.NewGetConfigurationConfirmation(resultKeys)
	conf.UnknownKey = unknownKeys
	return conf, nil
}

func (handler *ChargePointHandler) OnRemoteStartTransaction(request *ocpp16.RemoteStartTransactionRequest) (confirmation *ocpp16.RemoteStartTransactionConfirmation, err error) {
	connector, ok := handler.connectors[request.ConnectorId]
	if !ok {
		return ocpp16.NewRemoteStartTransactionConfirmation(ocpp16.RemoteStartStopStatusRejected), nil
	} else if connector.availability != ocpp16.AvailabilityTypeOperative || connector.status != ocpp16.ChargePointStatusAvailable || connector.currentTransaction > 0 {
		return ocpp16.NewRemoteStartTransactionConfirmation(ocpp16.RemoteStartStopStatusRejected), nil
	}
	logDefault(request.GetFeatureName()).Infof("started transaction %v on connector %v", connector.currentTransaction, request.ConnectorId)
	connector.currentTransaction = request.ConnectorId
	return ocpp16.NewRemoteStartTransactionConfirmation(ocpp16.RemoteStartStopStatusAccepted), nil
}

func (handler *ChargePointHandler) OnRemoteStopTransaction(request *ocpp16.RemoteStopTransactionRequest) (confirmation *ocpp16.RemoteStopTransactionConfirmation, err error) {
	for key, val := range handler.connectors {
		if val.currentTransaction == request.TransactionId {
			logDefault(request.GetFeatureName()).Infof("stopped transaction %v on connector %v", val.currentTransaction, key)
			val.currentTransaction = 0
			val.currentReservation = 0
			val.status = ocpp16.ChargePointStatusAvailable
			return ocpp16.NewRemoteStopTransactionConfirmation(ocpp16.RemoteStartStopStatusAccepted), nil
		}
	}
	return ocpp16.NewRemoteStopTransactionConfirmation(ocpp16.RemoteStartStopStatusRejected), nil
}

func (handler *ChargePointHandler) OnReset(request *ocpp16.ResetRequest) (confirmation *ocpp16.ResetConfirmation, err error) {
	//TODO: stop all ongoing transactions
	return ocpp16.NewResetConfirmation(ocpp16.ResetStatusAccepted), nil
}

func (handler *ChargePointHandler) OnUnlockConnector(request *ocpp16.UnlockConnectorRequest) (confirmation *ocpp16.UnlockConnectorConfirmation, err error) {
	_, ok := handler.connectors[request.ConnectorId]
	if !ok {
		return ocpp16.NewUnlockConnectorConfirmation(ocpp16.UnlockStatusNotSupported), nil
	}
	return ocpp16.NewUnlockConnectorConfirmation(ocpp16.UnlockStatusUnlocked), nil
}

// Local authorization list profile callbacks
func (handler *ChargePointHandler) OnGetLocalListVersion(request *ocpp16.GetLocalListVersionRequest) (confirmation *ocpp16.GetLocalListVersionConfirmation, err error) {
	return ocpp16.NewGetLocalListVersionConfirmation(handler.localAuthListVersion), nil
}

func (handler *ChargePointHandler) OnSendLocalList(request *ocpp16.SendLocalListRequest) (confirmation *ocpp16.SendLocalListConfirmation, err error) {
	if request.ListVersion <= handler.localAuthListVersion {
		return ocpp16.NewSendLocalListConfirmation(ocpp16.UpdateStatusVersionMismatch), nil
	}
	if request.UpdateType == ocpp16.UpdateTypeFull {
		handler.localAuthList = request.LocalAuthorizationList
		handler.localAuthListVersion = request.ListVersion
	} else if request.UpdateType == ocpp16.UpdateTypeDifferential {
		handler.localAuthList = append(handler.localAuthList, request.LocalAuthorizationList...)
		handler.localAuthListVersion = request.ListVersion
	}
	return ocpp16.NewSendLocalListConfirmation(ocpp16.UpdateStatusAccepted), nil
}

// Firmware management profile callbacks
func (handler *ChargePointHandler) OnGetDiagnostics(request *ocpp16.GetDiagnosticsRequest) (confirmation *ocpp16.GetDiagnosticsConfirmation, err error) {
	return ocpp16.NewGetDiagnosticsConfirmation(), nil
	//TODO: perform diagnostics upload out-of-band
}

func (handler *ChargePointHandler) OnUpdateFirmware(request *ocpp16.UpdateFirmwareRequest) (confirmation *ocpp16.UpdateFirmwareConfirmation, err error) {
	return ocpp16.NewUpdateFirmwareConfirmation(), nil
	//TODO: download new firmware out-of-band
}

// Remote trigger profile callbacks
func (handler *ChargePointHandler) OnTriggerMessage(request *ocpp16.TriggerMessageRequest) (confirmation *ocpp16.TriggerMessageConfirmation, err error) {
	switch request.RequestedMessage {
	case ocpp16.BootNotificationFeatureName:
		//TODO: schedule boot notification message
		break
	case ocpp16.DiagnosticsStatusNotificationFeatureName:
		//TODO: schedule diagnostics status notification message
		break
	case ocpp16.FirmwareStatusNotificationFeatureName:
		//TODO: schedule firmware status notification message
		break
	case ocpp16.HeartbeatFeatureName:
		//TODO: schedule heartbeat message
		break
	case ocpp16.MeterValuesFeatureName:
		//TODO: schedule meter values message
		break
		//TODO: schedule status notification message
	case ocpp16.StatusNotificationFeatureName:
		break
	default:
		return ocpp16.NewTriggerMessageConfirmation(ocpp16.TriggerMessageStatusNotImplemented), nil
	}
	return ocpp16.NewTriggerMessageConfirmation(ocpp16.TriggerMessageStatusAccepted), nil
}

// Reservation profile callbacks
func (handler *ChargePointHandler) OnReserveNow(request *ocpp16.ReserveNowRequest) (confirmation *ocpp16.ReserveNowConfirmation, err error) {
	connector := handler.connectors[request.ConnectorId]
	if connector == nil {
		return ocpp16.NewReserveNowConfirmation(ocpp16.ReservationStatusUnavailable), nil
	} else if connector.status != ocpp16.ChargePointStatusAvailable {
		return ocpp16.NewReserveNowConfirmation(ocpp16.ReservationStatusOccupied), nil
	}
	connector.status = ocpp16.ChargePointStatusReserved
	connector.currentReservation = request.ReservationId
	// TODO: notify status change
	// TODO: automatically remove reservation after expiryDate
	return ocpp16.NewReserveNowConfirmation(ocpp16.ReservationStatusAccepted), nil
}

func (handler *ChargePointHandler) OnCancelReservation(request *ocpp16.CancelReservationRequest) (confirmation *ocpp16.CancelReservationConfirmation, err error) {
	for _,v := range handler.connectors {
		if v.currentReservation == request.ReservationId {
			v.currentReservation = 0
			if v.status == ocpp16.ChargePointStatusReserved {
				// TODO: notify status change
				v.status = ocpp16.ChargePointStatusAvailable
			}
			return ocpp16.NewCancelReservationConfirmation(ocpp16.CancelReservationStatusAccepted), nil
		}
	}
	return ocpp16.NewCancelReservationConfirmation(ocpp16.CancelReservationStatusRejected), nil
}

// Smart charging profile callbacks
func (handler *ChargePointHandler) OnSetChargingProfile(request *ocpp16.SetChargingProfileRequest) (confirmation *ocpp16.SetChargingProfileConfirmation, err error) {
	//TODO: handle logic
	return ocpp16.NewSetChargingProfileConfirmation(ocpp16.ChargingProfileStatusNotImplemented), nil
}

func (handler *ChargePointHandler) OnClearChargingProfile(request *ocpp16.ClearChargingProfileRequest) (confirmation *ocpp16.ClearChargingProfileConfirmation, err error) {
	//TODO: handle logic
	return ocpp16.NewClearChargingProfileConfirmation(ocpp16.ClearChargingProfileStatusUnknown), nil
}

func (handler *ChargePointHandler) OnGetCompositeSchedule(request *ocpp16.GetCompositeScheduleRequest) (confirmation *ocpp16.GetCompositeScheduleConfirmation, err error) {
	//TODO: handle logic
	return ocpp16.NewGetCompositeScheduleConfirmation(ocpp16.GetCompositeScheduleStatusRejected), nil
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getExpiryDate(info *ocpp16.IdTagInfo) string {
	if info.ExpiryDate != nil {
		return fmt.Sprintf("authorized until %v", info.ExpiryDate.String())
	}
	return ""
}

func exampleRoutine(chargePoint ocpp16.ChargePoint, stateHandler *ChargePointHandler) {
	dummyClientIdTag := "12345"
	chargingConnector := 1
	// Boot
	bootConf, err := chargePoint.BootNotification("model1", "vendor1")
	checkError(err)
	logDefault(bootConf.GetFeatureName()).Infof("status: %v, interval: %v, current time: %v", bootConf.Status, bootConf.Interval, bootConf.CurrentTime.String())
	// Notify connector status
	statusConf, err := chargePoint.StatusNotification(0, ocpp16.NoError, ocpp16.ChargePointStatusAvailable)
	checkError(err)
	logDefault(statusConf.GetFeatureName()).Infof("status for all connectors updated to %v", ocpp16.ChargePointStatusAvailable)
	// Wait for some time ...
	time.Sleep(5 * time.Second)
	// Simulate charging for connector 1
	authConf, err := chargePoint.Authorize(dummyClientIdTag)
	checkError(err)
	logDefault(authConf.GetFeatureName()).Infof("status: %v %v", authConf.IdTagInfo.Status, getExpiryDate(authConf.IdTagInfo))
	stateHandler.connectors[chargingConnector].status = ocpp16.ChargePointStatusPreparing
	// Update connector status
	statusConf, err = chargePoint.StatusNotification(chargingConnector, ocpp16.NoError, stateHandler.connectors[chargingConnector].status)
	checkError(err)
	logDefault(statusConf.GetFeatureName()).Infof("status for connector %v updated to %v", chargingConnector, stateHandler.connectors[chargingConnector].status)
	startConf, err := chargePoint.StartTransaction(chargingConnector, dummyClientIdTag, stateHandler.meterValue, ocpp16.NewDateTime(time.Now()))
	checkError(err)
	logDefault(startConf.GetFeatureName()).Infof("status: %v, transaction %v %v", startConf.IdTagInfo.Status, startConf.TransactionId, getExpiryDate(startConf.IdTagInfo))
	stateHandler.connectors[chargingConnector].currentTransaction = startConf.TransactionId
	stateHandler.connectors[chargingConnector].status = ocpp16.ChargePointStatusCharging
	// Update connector status
	statusConf, err = chargePoint.StatusNotification(chargingConnector, ocpp16.NoError, stateHandler.connectors[chargingConnector].status)
	checkError(err)
	logDefault(statusConf.GetFeatureName()).Infof("status for connector %v updated to %v", chargingConnector, stateHandler.connectors[chargingConnector].status)
	// Periodically send meter values
	for i := 0; i < 5; i++ {
		time.Sleep(5 * time.Second)
		stateHandler.meterValue += 10
		sampledValue := ocpp16.SampledValue{Value: fmt.Sprintf("%v", stateHandler.meterValue), Unit: ocpp16.UnitOfMeasureWh, Format: ocpp16.ValueFormatRaw, Measurand: ocpp16.MeasurandEnergyActiveExportRegister, Context: ocpp16.ReadingContextSamplePeriodic, Location: ocpp16.LocationOutlet}
		meterValue := ocpp16.MeterValue{Timestamp: ocpp16.NewDateTime(time.Now()), SampledValue: []ocpp16.SampledValue{sampledValue}}
		meterConf, err := chargePoint.MeterValues(chargingConnector, []ocpp16.MeterValue{meterValue})
		checkError(err)
		logDefault(meterConf.GetFeatureName()).Infof("sent updated %v", sampledValue.Measurand)
	}
	stateHandler.meterValue += 2
	// Stop charging for connector 1
	stateHandler.connectors[chargingConnector].status = ocpp16.ChargePointStatusFinishing
	statusConf, err = chargePoint.StatusNotification(chargingConnector, ocpp16.NoError, stateHandler.connectors[chargingConnector].status)
	checkError(err)
	logDefault(statusConf.GetFeatureName()).Infof("status for connector %v updated to %v", chargingConnector, stateHandler.connectors[chargingConnector].status)
	stopConf, err := chargePoint.StopTransaction(stateHandler.meterValue, ocpp16.NewDateTime(time.Now()), startConf.TransactionId, func(request *ocpp16.StopTransactionRequest) {
		sampledValue := ocpp16.SampledValue{Value: fmt.Sprintf("%v", stateHandler.meterValue), Unit: ocpp16.UnitOfMeasureWh, Format: ocpp16.ValueFormatRaw, Measurand: ocpp16.MeasurandEnergyActiveExportRegister, Context: ocpp16.ReadingContextSamplePeriodic, Location: ocpp16.LocationOutlet}
		meterValue := ocpp16.MeterValue{Timestamp: ocpp16.NewDateTime(time.Now()), SampledValue: []ocpp16.SampledValue{sampledValue}}
		request.TransactionData = []ocpp16.MeterValue{meterValue}
		request.Reason = ocpp16.ReasonEVDisconnected
	})
	checkError(err)
	logDefault(stopConf.GetFeatureName()).Infof("transaction %v stopped", startConf.TransactionId)
	// Update connector status
	stateHandler.connectors[chargingConnector].status = ocpp16.ChargePointStatusAvailable
	statusConf, err = chargePoint.StatusNotification(chargingConnector, ocpp16.NoError, stateHandler.connectors[chargingConnector].status)
	checkError(err)
	logDefault(statusConf.GetFeatureName()).Infof("status for connector %v updated to %v", chargingConnector, stateHandler.connectors[chargingConnector].status)
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
	chargePoint := ocpp16.NewChargePoint(id, nil, nil)
	// Set a handler for all callback functions
	connectors := map[int]*ConnectorInfo{
		1: {status: ocpp16.ChargePointStatusAvailable, availability: ocpp16.AvailabilityTypeOperative, currentTransaction: 0},
	}
	handler := &ChargePointHandler{
		status:               ocpp16.ChargePointStatusAvailable,
		connectors:           connectors,
		configuration:        map[string]ocpp16.ConfigurationKey{},
		errorCode:            ocpp16.NoError,
		localAuthList:        []ocpp16.AuthorizationData{},
		localAuthListVersion: 0}
	chargePoint.SetChargePointCoreListener(handler)
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
