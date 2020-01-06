package main

import (
	"fmt"
	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	log "github.com/sirupsen/logrus"
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
	status            ocpp16.ChargePointStatus
	diagnosticsStatus ocpp16.DiagnosticsStatus
	firmwareStatus    ocpp16.FirmwareStatus
	connectors        map[int]*ConnectorInfo // No assumptions about the # of connectors
	transactions      map[int]*TransactionInfo
	errorCode         ocpp16.ChargePointErrorCode
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
	logDefault(chargePointId, request.GetFeatureName()).Infof("client authorized")
	return ocpp16.NewAuthorizationConfirmation(ocpp16.NewIdTagInfo(ocpp16.AuthorizationStatusAccepted)), nil
}

func (handler *CentralSystemHandler) OnBootNotification(chargePointId string, request *ocpp16.BootNotificationRequest) (confirmation *ocpp16.BootNotificationConfirmation, err error) {
	logDefault(chargePointId, request.GetFeatureName()).Infof("boot confirmed")
	return ocpp16.NewBootNotificationConfirmation(ocpp16.NewDateTime(time.Now()), defaultHeartbeatInterval, ocpp16.RegistrationStatusAccepted), nil
}

func (handler *CentralSystemHandler) OnDataTransfer(chargePointId string, request *ocpp16.DataTransferRequest) (confirmation *ocpp16.DataTransferConfirmation, err error) {
	logDefault(chargePointId, request.GetFeatureName()).Infof("received data %d", request.Data)
	return ocpp16.NewDataTransferConfirmation(ocpp16.DataTransferStatusAccepted), nil
}

func (handler *CentralSystemHandler) OnHeartbeat(chargePointId string, request *ocpp16.HeartbeatRequest) (confirmation *ocpp16.HeartbeatConfirmation, err error) {
	return ocpp16.NewHeartbeatConfirmation(ocpp16.NewDateTime(time.Now())), nil
}

func (handler *CentralSystemHandler) OnMeterValues(chargePointId string, request *ocpp16.MeterValuesRequest) (confirmation *ocpp16.MeterValuesConfirmation, err error) {
	logDefault(chargePointId, request.GetFeatureName()).Infof("received meter values for connector %v, transaction %v. Meter values:\n", request.ConnectorId, request.TransactionId)
	for _, mv := range request.MeterValue {
		logDefault(chargePointId, request.GetFeatureName()).Printf("\t %v", mv)
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
		logDefault(chargePointId, request.GetFeatureName()).Infof("connector %v updated status to %v", request.ConnectorId, request.Status)
	} else {
		info.status = request.Status
		logDefault(chargePointId, request.GetFeatureName()).Infof("all connectors updated status to %v", request.Status)
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
		return nil, fmt.Errorf("connector %v is currently busy with another transaction", request.ConnectorId)
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
	logDefault(chargePointId, request.GetFeatureName()).Infof("started transaction %v for connector %v", transaction.id, transaction.connectorId)
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
	logDefault(chargePointId, request.GetFeatureName()).Infof("stopped transaction %v - %v. Meter values:", request.TransactionId, request.Reason)
	for _, mv := range request.TransactionData {
		logDefault(chargePointId, request.GetFeatureName()).Printf("\t %v", mv)
	}
	return ocpp16.NewStopTransactionConfirmation(), nil
}

// Firmware management callbacks
func (handler *CentralSystemHandler) OnDiagnosticsStatusNotification(chargePointId string, request *ocpp16.DiagnosticsStatusNotificationRequest) (confirmation *ocpp16.DiagnosticsStatusNotificationConfirmation, err error) {
	info, ok := handler.chargePoints[chargePointId]
	if !ok {
		return nil, fmt.Errorf("unknown charge point %v", chargePointId)
	}
	info.diagnosticsStatus = request.Status
	logDefault(chargePointId, request.GetFeatureName()).Infof("updated diagnostics status to %v", request.Status)
	return ocpp16.NewDiagnosticsStatusNotificationConfirmation(), nil
}

func (handler *CentralSystemHandler) OnFirmwareStatusNotification(chargePointId string, request *ocpp16.FirmwareStatusNotificationRequest) (confirmation *ocpp16.FirmwareStatusNotificationConfirmation, err error) {
	info, ok := handler.chargePoints[chargePointId]
	if !ok {
		return nil, fmt.Errorf("unknown charge point %v", chargePointId)
	}
	info.firmwareStatus = request.Status
	logDefault(chargePointId, request.GetFeatureName()).Infof("updated firmware status to %v", request.Status)
	return &ocpp16.FirmwareStatusNotificationConfirmation{}, nil
}

// No callbacks for Local Auth management, Reservation, Remote trigger or Smart Charging profile on central system

// Start function
func main() {
	args := os.Args[1:]
	centralSystem := ocpp16.NewCentralSystem(nil, nil)
	handler := &CentralSystemHandler{chargePoints: map[string]*ChargePointState{}}
	centralSystem.SetNewChargePointHandler(func(chargePointId string) {
		handler.chargePoints[chargePointId] = &ChargePointState{connectors: map[int]*ConnectorInfo{}, transactions: map[int]*TransactionInfo{}}
		log.WithField("client", chargePointId).Info("new charge point connected")
	})
	centralSystem.SetChargePointDisconnectedHandler(func(chargePointId string) {
		log.WithField("client", chargePointId).Info("charge point disconnected")
		delete(handler.chargePoints, chargePointId)
	})
	centralSystem.SetCentralSystemCoreListener(handler)
	var listenPort = defaultListenPort
	if len(args) > 0 {
		port, err := strconv.Atoi(args[0])
		if err != nil {
			listenPort = port
		}
	}
	log.Infof("starting central system on port %v", listenPort)
	centralSystem.Start(listenPort, "/{ws}")
	log.Info("stopped central system")
}

// Utility functions
func logDefault(chargePointId string, feature string) *log.Entry {
	return log.WithFields(log.Fields{"client": chargePointId, "message": feature})
}

func init() {
	log.SetLevel(log.InfoLevel)
}
