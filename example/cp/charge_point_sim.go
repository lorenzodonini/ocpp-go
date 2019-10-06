package main

import (
	"fmt"
	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	"log"
	"os"
	"time"
)

type ConnectorInfo struct {
	status ocpp16.ChargePointStatus
	availability ocpp16.AvailabilityType
	currentTransaction int
}

type ChargePointHandler struct {
	status ocpp16.ChargePointStatus
	connectors map[int]*ConnectorInfo
	errorCode ocpp16.ChargePointErrorCode
	configuration map[string]ocpp16.ConfigurationKey
	meterValue int
}

func (handler * ChargePointHandler) OnChangeAvailability(request *ocpp16.ChangeAvailabilityRequest) (confirmation *ocpp16.ChangeAvailabilityConfirmation, err error) {
	handler.connectors[request.ConnectorId].availability = request.Type
	return ocpp16.NewChangeAvailabilityConfirmation(ocpp16.AvailabilityStatusAccepted), nil
}

func (handler * ChargePointHandler) OnChangeConfiguration(request *ocpp16.ChangeConfigurationRequest) (confirmation *ocpp16.ChangeConfigurationConfirmation, err error) {
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

func (handler * ChargePointHandler) OnClearCache(request *ocpp16.ClearCacheRequest) (confirmation *ocpp16.ClearCacheConfirmation, err error) {
	return ocpp16.NewClearCacheConfirmation(ocpp16.ClearCacheStatusAccepted), nil
}

func (handler * ChargePointHandler) OnDataTransfer(request *ocpp16.DataTransferRequest) (confirmation *ocpp16.DataTransferConfirmation, err error) {
	log.Printf("data transfer [Vendor: %v Message: %v]: %v", request.VendorId, request.MessageId, request.Data)
	return ocpp16.NewDataTransferConfirmation(ocpp16.DataTransferStatusAccepted), nil
}

func (handler * ChargePointHandler) OnGetConfiguration(request *ocpp16.GetConfigurationRequest) (confirmation *ocpp16.GetConfigurationConfirmation, err error) {
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

func (handler * ChargePointHandler) OnRemoteStartTransaction(request *ocpp16.RemoteStartTransactionRequest) (confirmation *ocpp16.RemoteStartTransactionConfirmation, err error) {
	connector, ok := handler.connectors[request.ConnectorId]
	if !ok {
		return ocpp16.NewRemoteStartTransactionConfirmation(ocpp16.RemoteStartStopStatusRejected), nil
	} else if connector.availability != ocpp16.AvailabilityTypeOperative || connector.currentTransaction > 0 {
		return ocpp16.NewRemoteStartTransactionConfirmation(ocpp16.RemoteStartStopStatusRejected), nil
	}
	log.Printf("started transaction %v on connector %v", connector.currentTransaction, request.ConnectorId)
	connector.currentTransaction = request.ConnectorId
	return ocpp16.NewRemoteStartTransactionConfirmation(ocpp16.RemoteStartStopStatusAccepted), nil
}

func (handler * ChargePointHandler) OnRemoteStopTransaction(request *ocpp16.RemoteStopTransactionRequest) (confirmation *ocpp16.RemoteStopTransactionConfirmation, err error) {
	for key, val := range handler.connectors {
		if val.currentTransaction == request.TransactionId {
			log.Printf("stopped transaction %v on connector %v", val.currentTransaction, key)
			val.currentTransaction = 0
			val.status = ocpp16.ChargePointStatusAvailable
			return ocpp16.NewRemoteStopTransactionConfirmation(ocpp16.RemoteStartStopStatusAccepted), nil
		}
	}
	return ocpp16.NewRemoteStopTransactionConfirmation(ocpp16.RemoteStartStopStatusRejected), nil
}

func (handler * ChargePointHandler) OnReset(request *ocpp16.ResetRequest) (confirmation *ocpp16.ResetConfirmation, err error) {
	//TODO: stop all ongoing transactions
	return ocpp16.NewResetConfirmation(ocpp16.ResetStatusAccepted), nil
}

func (handler * ChargePointHandler) OnUnlockConnector(request *ocpp16.UnlockConnectorRequest) (confirmation *ocpp16.UnlockConnectorConfirmation, err error) {
	_, ok := handler.connectors[request.ConnectorId]
	if !ok {
		return ocpp16.NewUnlockConnectorConfirmation(ocpp16.UnlockStatusNotSupported), nil
	}
	return ocpp16.NewUnlockConnectorConfirmation(ocpp16.UnlockStatusUnlocked), nil
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
	log.Printf("[Boot confirmation] status: %v, interval: %v, current time: %v", bootConf.Status, bootConf.Interval, bootConf.CurrentTime.String())
	// Notify connector status
	_, err = chargePoint.StatusNotification(0, ocpp16.NoError, ocpp16.ChargePointStatusAvailable)
	checkError(err)
	// Wait for some time ...
	time.Sleep(5 * time.Second)
	// Simulate charging for connector 1
	authConf, err := chargePoint.Authorize(dummyClientIdTag)
	checkError(err)
	log.Printf("[Authorize confirmation] status: %v %v", authConf.IdTagInfo.Status, getExpiryDate(authConf.IdTagInfo))
	stateHandler.connectors[chargingConnector].status = ocpp16.ChargePointStatusPreparing
	_, err = chargePoint.StatusNotification(chargingConnector, ocpp16.NoError, stateHandler.connectors[chargingConnector].status)
	checkError(err)
	startConf, err := chargePoint.StartTransaction(chargingConnector, dummyClientIdTag, stateHandler.meterValue, ocpp16.NewDateTime(time.Now()))
	checkError(err)
	log.Printf("[Start transaction confirmation] status: %v, transaction %v %v", startConf.IdTagInfo.Status, startConf.TransactionId, getExpiryDate(startConf.IdTagInfo))
	stateHandler.connectors[chargingConnector].currentTransaction = startConf.TransactionId
	stateHandler.connectors[chargingConnector].status = ocpp16.ChargePointStatusCharging
	_, err = chargePoint.StatusNotification(chargingConnector, ocpp16.NoError, stateHandler.connectors[chargingConnector].status)
	checkError(err)
	// Periodically send meter values
	for i := 0; i<5; i++ {
		time.Sleep(5 * time.Second)
		stateHandler.meterValue += 10
		sampledValue := ocpp16.SampledValue{Value: fmt.Sprintf("%v", stateHandler.meterValue), Unit:ocpp16.UnitOfMeasureWh, Format: ocpp16.ValueFormatRaw, Measurand: ocpp16.MeasurandEnergyActiveExportRegister, Context: ocpp16.ReadingContextSamplePeriodic, Location: ocpp16.LocationOutlet }
		meterValue := ocpp16.MeterValue{Timestamp: ocpp16.NewDateTime(time.Now()), SampledValue: []ocpp16.SampledValue{sampledValue}}
		_, err := chargePoint.MeterValues(chargingConnector, []ocpp16.MeterValue{meterValue})
		log.Printf("[Meter values confirmation] sent updated %v", sampledValue.Measurand)
		checkError(err)
	}
	stateHandler.meterValue += 2
	// Stop charging for connector 1
	stateHandler.connectors[chargingConnector].status = ocpp16.ChargePointStatusFinishing
	_, err = chargePoint.StatusNotification(chargingConnector, ocpp16.NoError, stateHandler.connectors[chargingConnector].status)
	checkError(err)
	_, err = chargePoint.StopTransaction(stateHandler.meterValue, ocpp16.NewDateTime(time.Now()), startConf.TransactionId, func(request *ocpp16.StopTransactionRequest) {
		sampledValue := ocpp16.SampledValue{Value: fmt.Sprintf("%v", stateHandler.meterValue), Unit:ocpp16.UnitOfMeasureWh, Format: ocpp16.ValueFormatRaw, Measurand:ocpp16.MeasurandEnergyActiveExportRegister, Context: ocpp16.ReadingContextSamplePeriodic, Location: ocpp16.LocationOutlet }
		meterValue := ocpp16.MeterValue{Timestamp: ocpp16.NewDateTime(time.Now()), SampledValue: []ocpp16.SampledValue{sampledValue}}
		request.TransactionData = []ocpp16.MeterValue{meterValue}
		request.Reason = ocpp16.ReasonEVDisconnected
	})
	checkError(err)
	log.Printf("[Stop transaction confirmation] transaction %v stopped", startConf.TransactionId)
	stateHandler.connectors[chargingConnector].status = ocpp16.ChargePointStatusFinishing
	_, err = chargePoint.StatusNotification(chargingConnector, ocpp16.NoError, stateHandler.connectors[chargingConnector].status)
	checkError(err)
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
	handler := &ChargePointHandler{status: ocpp16.ChargePointStatusAvailable, connectors: connectors, configuration: map[string]ocpp16.ConfigurationKey{}, errorCode: ocpp16.NoError}
	chargePoint.SetChargePointCoreListener(handler)
	// Connects to central system
	err := chargePoint.Start(csUrl)
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("connected to central system at %v", csUrl)
		exampleRoutine(chargePoint, handler)
		// Disconnect
		chargePoint.Stop()
		log.Printf("disconnected from central system")
	}
}
