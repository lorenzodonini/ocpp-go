package ocpp16

import (
	"reflect"
)

// -------------------- Heartbeat (CP -> CS) --------------------
type MeterValuesRequest struct {
	ConnectorId   int          `json:"connectorId" validate:"gte=0"`
	TransactionId int          `json:"reservationId,omitempty"`
	MeterValue    []MeterValue `json:"meterValue" validate:"required,min=1,dive"`
}

type MeterValuesConfirmation struct {
}

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

func NewMeterValuesRequest(connectorId int, meterValues []MeterValue) *MeterValuesRequest {
	return &MeterValuesRequest{ConnectorId: connectorId, MeterValue: meterValues}
}

func NewMeterValuesConfirmation() *MeterValuesConfirmation {
	return &MeterValuesConfirmation{}
}
