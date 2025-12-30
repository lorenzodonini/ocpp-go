package smartcharging

import (
	"reflect"

	"github.com/xBlaz3kx/ocpp-go/ocpp2.1/types"
)

// -------------------- Cleared Charging Limit (CS -> CSMS) --------------------

const ClearedChargingLimitFeatureName = "ClearedChargingLimit"

// The field definition of the ClearedChargingLimit request payload sent by the Charging Station to the CSMS.
type ClearedChargingLimitRequest struct {
	ChargingLimitSource types.ChargingLimitSourceType `json:"chargingLimitSource" validate:"required,chargingLimitSource21"`
	EvseID              *int                          `json:"evseId,omitempty" validate:"omitempty,gte=0"`
}

// This field definition of the ClearedChargingLimit response payload, sent by the CSMS to the Charging Station in response to a ClearedChargingLimitRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type ClearedChargingLimitResponse struct {
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

func (f ClearedChargingLimitFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(ClearedChargingLimitResponse{})
}

func (r ClearedChargingLimitRequest) GetFeatureName() string {
	return ClearedChargingLimitFeatureName
}

func (c ClearedChargingLimitResponse) GetFeatureName() string {
	return ClearedChargingLimitFeatureName
}

// Creates a new ClearedChargingLimitRequest, containing all required fields. Optional fields may be set afterwards.
func NewClearedChargingLimitRequest(chargingLimitSource types.ChargingLimitSourceType) *ClearedChargingLimitRequest {
	return &ClearedChargingLimitRequest{ChargingLimitSource: chargingLimitSource}
}

// Creates a new ClearedChargingLimitResponse, which doesn't contain any required or optional fields.
func NewClearedChargingLimitResponse() *ClearedChargingLimitResponse {
	return &ClearedChargingLimitResponse{}
}
