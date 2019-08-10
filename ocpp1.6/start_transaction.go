package ocpp16

import (
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Heartbeat (CP -> CS) --------------------
type StartTransactionRequest struct {
	ConnectorId   int      `json:"connectorId" validate:"gt=0"`
	IdTag         string   `json:"idTag" validate:"required,max=20"`
	MeterStart    int      `json:"meterStart" validate:"required"`
	ReservationId int      `json:"reservationId,omitempty"`
	Timestamp     DateTime `json:"timestamp"`
}

type StartTransactionConfirmation struct {
	IdTagInfo     IdTagInfo `json:"idTagInfo" validate:"required"`
	TransactionId int       `json:"transactionId" validate:"required"`
}

type StartTransactionFeature struct{}

func (f StartTransactionFeature) GetFeatureName() string {
	return StartTransactionFeatureName
}

func (f StartTransactionFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(StartTransactionRequest{})
}

func (f StartTransactionFeature) GetConfirmationType() reflect.Type {
	return reflect.TypeOf(StartTransactionConfirmation{})
}

func (r StartTransactionRequest) GetFeatureName() string {
	return StartTransactionFeatureName
}

func (c StartTransactionConfirmation) GetFeatureName() string {
	return StartTransactionFeatureName
}

func NewStartTransactionRequest(connectorId int, idTag string, meterStart int, timestamp DateTime) *StartTransactionRequest {
	return &StartTransactionRequest{ConnectorId: connectorId, IdTag: idTag, MeterStart: meterStart, Timestamp: timestamp}
}

func NewStartTransactionConfirmation(idTagInfo IdTagInfo, transactionId int) *StartTransactionConfirmation {
	return &StartTransactionConfirmation{IdTagInfo: idTagInfo, TransactionId: transactionId}
}

func validateStartTransactionRequest(sl validator.StructLevel) {
	confirmation := sl.Current().Interface().(StartTransactionRequest)
	if dateTimeIsNull(confirmation.Timestamp) {
		sl.ReportError(confirmation.Timestamp, "Timestamp", "timestamp", "required", "")
	}
	//if !validateDateTimeNow(confirmation.CurrentTime) {
	//	sl.ReportError(confirmation.CurrentTime, "CurrentTime", "currentTime", "eq", "")
	//}
}

func init() {
	Validate.RegisterStructValidation(validateStartTransactionRequest, StartTransactionRequest{})
}
