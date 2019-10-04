package ocpp16

import (
	"reflect"
)

// -------------------- Heartbeat (CP -> CS) --------------------

// The field definition of the MeterValues request payload sent by the Charge Point to the Central System.
type MeterValuesRequest struct {
	ConnectorId   int          `json:"connectorId" validate:"gte=0"`
	TransactionId int          `json:"reservationId,omitempty"`
	MeterValue    []MeterValue `json:"meterValue" validate:"required,min=1,dive"`
}

// This field definition of the Authorize confirmation payload, sent by the Charge Point to the Central System in response to an AuthorizeRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type MeterValuesConfirmation struct {
}

// A Charge Point MAY sample the electrical meter or other sensor/transducer hardware to provide extra information about its meter values.
// It is up to the Charge Point to decide when it will send meter values.
// This can be configured using the ChangeConfiguration message to specify data acquisition intervals and specify data to be acquired & reported.
type MeterValuesFeature struct{}

func (f MeterValuesFeature) GetFeatureName() string {
	return MeterValuesFeatureName
}

func (f MeterValuesFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(MeterValuesRequest{})
}

func (f MeterValuesFeature) GetConfirmationType() reflect.Type {
	return reflect.TypeOf(MeterValuesConfirmation{})
}

func (r MeterValuesRequest) GetFeatureName() string {
	return MeterValuesFeatureName
}

func (c MeterValuesConfirmation) GetFeatureName() string {
	return MeterValuesFeatureName
}

// Creates a new MeterValuesRequest, containing all required fields. Optional fields may be set afterwards.
func NewMeterValuesRequest(connectorId int, meterValues []MeterValue) *MeterValuesRequest {
	return &MeterValuesRequest{ConnectorId: connectorId, MeterValue: meterValues}
}

// Creates a new MeterValuesConfirmation, which doesn't contain any required or optional fields.
func NewMeterValuesConfirmation() *MeterValuesConfirmation {
	return &MeterValuesConfirmation{}
}
