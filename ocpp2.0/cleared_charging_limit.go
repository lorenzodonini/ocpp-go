package ocpp2

import (
	"reflect"
)

// -------------------- Cleared Charging Limit (CS -> CSMS) --------------------

// The field definition of the ClearedChargingLimit request payload sent by the Charging Station to the CSMS.
type ClearedChargingLimitRequest struct {
	ChargingLimitSource ChargingLimitSourceType `json:"chargingLimitSource" validate:"required,chargingLimitSource"`
	EvseID              *int                    `json:"evseId,omitempty" validate:"omitempty,gte=0"`
}

// This field definition of the ClearedChargingLimit confirmation payload, sent by the CSMS to the Charging Station in response to a ClearedChargingLimitRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type ClearedChargingLimitConfirmation struct {
}

// When an external control system sends a signal to release a previously imposed charging limit to a Charging Station,
// the Charging Station sends a ClearedChargingLimitRequest to notify the CSMS about this.
// The CSMS acknowledges with a ClearedChargingLimitResponse to the Charging Station.
// When the change has impact on an ongoing charging transaction and is more than: LimitChangeSignificance,
// the Charging Station needs to send a TransactionEventRequest to notify the CSMS.
type ClearedChargingLimitFeature struct{}

func (f ClearedChargingLimitFeature) GetFeatureName() string {
	return ClearedChargingLimitFeatureName
}

func (f ClearedChargingLimitFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(ClearedChargingLimitRequest{})
}

func (f ClearedChargingLimitFeature) GetConfirmationType() reflect.Type {
	return reflect.TypeOf(ClearedChargingLimitConfirmation{})
}

func (r ClearedChargingLimitRequest) GetFeatureName() string {
	return ClearedChargingLimitFeatureName
}

func (c ClearedChargingLimitConfirmation) GetFeatureName() string {
	return ClearedChargingLimitFeatureName
}

// Creates a new ClearedChargingLimitRequest, containing all required fields. Optional fields may be set afterwards.
func NewClearedChargingLimitRequest(chargingLimitSource ChargingLimitSourceType) *ClearedChargingLimitRequest {
	return &ClearedChargingLimitRequest{ChargingLimitSource: chargingLimitSource}
}

// Creates a new ClearedChargingLimitConfirmation, which doesn't contain any required or optional fields.
func NewClearedChargingLimitConfirmation() *ClearedChargingLimitConfirmation {
	return &ClearedChargingLimitConfirmation{}
}
