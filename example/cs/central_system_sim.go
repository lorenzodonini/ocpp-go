package main

import (
	"fmt"
	ocpp16 "github.com/lorenzodonini/go-ocpp/ocpp1.6"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	defaultListenPort        = 8887
	defaultHeartbeatInterval = 600
)

var (
	nextTransactionId = 0
)

//TODO: cache authorization

// Charge Point state
type TransactionInfo struct {
	id          int
	startTime   *ocpp16.DateTime
	endTime     *ocpp16.DateTime
	startMeter  int
	endMeter    int
	connectorId int
	idTag       string
}

func (ti *TransactionInfo) hasTransactionEnded() bool {
	return ti.endTime != nil && !ti.endTime.IsZero()
}

type ConnectorInfo struct {
	status             ocpp16.ChargePointStatus
	currentTransaction int
}

func (ci *ConnectorInfo) hasTransactionInProgress() bool {
	return ci.currentTransaction >= 0
}

type ChargePointState struct {
	status       ocpp16.ChargePointStatus
	connectors   map[int]*ConnectorInfo // No assumptions about the # of connectors
	transactions map[int]*TransactionInfo
	errorCode    ocpp16.ChargePointErrorCode
}

func (cps *ChargePointState) getConnector(id int) *ConnectorInfo {
	ci, ok := cps.connectors[id]
	if !ok {
		ci = &ConnectorInfo{currentTransaction: -1}
		cps.connectors[id] = ci
	}
	return ci
}

type CentralSystemHandler struct {
	chargePoints map[string]*ChargePointState
}

// Core profile callbacks
func (handler *CentralSystemHandler) OnAuthorize(chargePointId string, request *ocpp16.AuthorizeRequest) (confirmation *ocpp16.AuthorizeConfirmation, err error) {
	return ocpp16.NewAuthorizationConfirmation(ocpp16.NewIdTagInfo(ocpp16.AuthorizationStatusAccepted)), nil
}

func (handler *CentralSystemHandler) OnBootNotification(chargePointId string, request *ocpp16.BootNotificationRequest) (confirmation *ocpp16.BootNotificationConfirmation, err error) {
	return ocpp16.NewBootNotificationConfirmation(ocpp16.NewDateTime(time.Now()), defaultHeartbeatInterval, ocpp16.RegistrationStatusAccepted), nil
}

func (handler *CentralSystemHandler) OnDataTransfer(chargePointId string, request *ocpp16.DataTransferRequest) (confirmation *ocpp16.DataTransferConfirmation, err error) {
	log.Printf("[%v] Received data %d", chargePointId, request.Data)
	return ocpp16.NewDataTransferConfirmation(ocpp16.DataTransferStatusAccepted), nil
}

func (handler *CentralSystemHandler) OnHeartbeat(chargePointId string, request *ocpp16.HeartbeatRequest) (confirmation *ocpp16.HeartbeatConfirmation, err error) {
	return ocpp16.NewHeartbeatConfirmation(ocpp16.NewDateTime(time.Now())), nil
}

func (handler *CentralSystemHandler) OnMeterValues(chargePointId string, request *ocpp16.MeterValuesRequest) (confirmation *ocpp16.MeterValuesConfirmation, err error) {
	log.Printf("[%v] Received meter values for connector %v, transaction %v. Meter values:\n", chargePointId, request.ConnectorId, request.TransactionId)
	for _, mv := range request.MeterValue {
		log.Printf("\t %v", mv)
	}
	return ocpp16.NewMeterValuesConfirmation(), nil
}

func (handler *CentralSystemHandler) OnStatusNotification(chargePointId string, request *ocpp16.StatusNotificationRequest) (confirmation *ocpp16.StatusNotificationConfirmation, err error) {
	info, ok := handler.chargePoints[chargePointId]
	if !ok {
		return nil, fmt.Errorf("unknown charge point %v", chargePointId)
	}
	info.errorCode = request.ErrorCode
	if request.ConnectorId > 0 {
		connectorInfo := info.getConnector(request.ConnectorId)
		connectorInfo.status = request.Status
	} else {
		info.status = request.Status
	}
	return ocpp16.NewStatusNotificationConfirmation(), nil
}

func (handler *CentralSystemHandler) OnStartTransaction(chargePointId string, request *ocpp16.StartTransactionRequest) (confirmation *ocpp16.StartTransactionConfirmation, err error) {
	info, ok := handler.chargePoints[chargePointId]
	if !ok {
		return nil, fmt.Errorf("unknown charge point %v", chargePointId)
	}
	connector := info.getConnector(request.ConnectorId)
	if connector.currentTransaction >= 0 {
		return nil, fmt.Errorf("connector %v is currently busy with another transaction")
	}
	transaction := &TransactionInfo{}
	transaction.idTag = request.IdTag
	transaction.connectorId = request.ConnectorId
	transaction.startMeter = request.MeterStart
	transaction.startTime = request.Timestamp
	transaction.id = nextTransactionId
	nextTransactionId += 1
	connector.currentTransaction = transaction.id
	info.transactions[transaction.id] = transaction
	//TODO: check billable clients
	log.Printf("[%v] started transaction %v for connector %v", chargePointId, transaction.id, transaction.connectorId)
	return ocpp16.NewStartTransactionConfirmation(ocpp16.NewIdTagInfo(ocpp16.AuthorizationStatusAccepted), transaction.id), nil
}

func (handler *CentralSystemHandler) OnStopTransaction(chargePointId string, request *ocpp16.StopTransactionRequest) (confirmation *ocpp16.StopTransactionConfirmation, err error) {
	info, ok := handler.chargePoints[chargePointId]
	if !ok {
		return nil, fmt.Errorf("unknown charge point %v", chargePointId)
	}
	transaction, ok := info.transactions[request.TransactionId]
	if ok {
		connector := info.getConnector(transaction.connectorId)
		connector.currentTransaction = -1
		transaction.endTime = request.Timestamp
		transaction.endMeter = request.MeterStop
		//TODO: meter data
	}
	log.Printf("[%v] stopped transaction %v - %v. Meter values:", chargePointId, request.TransactionId, request.Reason)
	for _, mv := range request.TransactionData {
		log.Printf("\t %v", mv)
	}
	return ocpp16.NewStopTransactionConfirmation(), nil
}

// Start function
func main() {
	args := os.Args[1:]
	centralSystem := ocpp16.NewCentralSystem(nil, nil)
	handler := &CentralSystemHandler{chargePoints: map[string]*ChargePointState{}}
	centralSystem.SetNewChargePointHandler(func(chargePointId string) {
		handler.chargePoints[chargePointId] = &ChargePointState{connectors: map[int]*ConnectorInfo{}, transactions: map[int]*TransactionInfo{}}
		log.Printf("new charge point %v connected", chargePointId)
	})
	centralSystem.SetChargePointDisconnectedHandler(func(chargePointId string) {
		log.Printf("charge point %v disconnected", chargePointId)
		delete(handler.chargePoints, chargePointId)
	})
	centralSystem.SetCentralSystemCoreListener(handler)
	var listenPort = defaultListenPort
	if len(args) > 0 {
		port, err := strconv.Atoi(args[1])
		if err != nil {
			listenPort = port
		}
	}
	log.Printf("starting central system on port %v", listenPort)
	centralSystem.Start(listenPort, "/{ws}")
	log.Println("stopped central system")
}
