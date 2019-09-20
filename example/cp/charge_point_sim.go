package main

import (
	ocpp16 "github.com/lorenzodonini/go-ocpp/ocpp1.6"
	"log"
	"os"
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

func exampleRoutine(chargePoint ocpp16.ChargePoint, stateHandler *ChargePointHandler) {

}

// Start function
func main() {
	// Parse arguments from cmd line
	args := os.Args[1:]
	if len(args) != 2 {
		log.Print("Usage:\n\tocppClientId\n\tocppServerUrl")
		return
	}
	id := args[0]
	csUrl := args[1]
	// Create a default OCPP 1.6 charge point
	chargePoint := ocpp16.NewChargePoint(id, nil, nil)
	// Set a handler for all callback functions
	handler := &ChargePointHandler{}
	chargePoint.SetChargePointCoreListener(handler)
	// Connects to central system
	err := chargePoint.Start(csUrl)
	if err != nil {
		log.Println(err)
	} else {
		exampleRoutine(chargePoint, handler)
	}
}
