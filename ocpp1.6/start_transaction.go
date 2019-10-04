package ocpp16

import (
	"reflect"
)

// -------------------- Start Transaction (CP -> CS) --------------------

// This field definition of the StartTransactionRequest payload sent by the Charge Point to the Central System.
type StartTransactionRequest struct {
	ConnectorId   int       `json:"connectorId" validate:"gt=0"`
	IdTag         string    `json:"idTag" validate:"required,max=20"`
	MeterStart    int       `json:"meterStart" validate:"gte=0"`
	ReservationId int       `json:"reservationId,omitempty"`
	Timestamp     *DateTime `json:"timestamp" validate:"required"`
}

// This field definition of the StartTransactionConfirmation payload sent by the Central System to the Charge Point in response to a StartTransactionRequest.
type StartTransactionConfirmation struct {
	IdTagInfo     *IdTagInfo `json:"idTagInfo" validate:"required"`
	TransactionId int        `json:"transactionId" validate:"gte=0"`
}

// The Charge Point SHALL send a StartTransactionRequest to the Central System to inform about a transaction that has been started.
// If this transaction ends a reservation (see ReserveNow operation), then the StartTransaction MUST contain the reservationId.
// Upon receipt of a StartTransaction.req PDU, the Central System SHOULD respond with a StartTransactionConfirmation.
// This response MUST include a transaction id and an authorization status value.
// The Central System MUST verify validity of the identifier in the StartTransactionRequest, because the identifier might have been authorized locally by the Charge Point using outdated information.
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

// Creates a new StartTransaction request. All signature parameters are required fields.
// Optional fields may be set directly on the created request.
func NewStartTransactionRequest(connectorId int, idTag string, meterStart int, timestamp *DateTime) *StartTransactionRequest {
	return &StartTransactionRequest{ConnectorId: connectorId, IdTag: idTag, MeterStart: meterStart, Timestamp: timestamp}
}

// Creates a new StartTransaction confirmation. All signature parameters are required fields. There are no optional fields for this message.
func NewStartTransactionConfirmation(idTagInfo *IdTagInfo, transactionId int) *StartTransactionConfirmation {
	return &StartTransactionConfirmation{IdTagInfo: idTagInfo, TransactionId: transactionId}
}
