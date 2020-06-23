package smartcharging

import (
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
	"reflect"
)

// -------------------- Notify Charging Limit (CS -> CSMS) --------------------

const NotifyChargingLimitFeatureName = "NotifyChargingLimit"

// ChargingLimit contains the source of the charging limit and whether it is grid critical.
type ChargingLimit struct {
	ChargingLimitSource types.ChargingLimitSourceType `json:"chargingLimitSource" validate:"required,chargingLimitSource"` // Represents the source of the charging limit.
	IsGridCritical      *bool                         `json:"isGridCritical,omitempty" validate:"omitempty"`               // Indicates whether the charging limit is critical for the grid.
}

// The field definition of the NotifyChargingLimit request payload sent by the Charging Station to the CSMS.
type NotifyChargingLimitRequest struct {
	EvseID           *int                     `json:"evseId,omitempty" validate:"omitempty,gte=0"`
	ChargingLimit    ChargingLimit            `json:"chargingLimit" validate:"required"`
	ChargingSchedule []types.ChargingSchedule `json:"chargingSchedule,omitempty" validate:"omitempty,dive"`
}

// This field definition of the NotifyChargingLimit response payload, sent by the CSMS to the Charging Station in response to a NotifyChargingLimitRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type NotifyChargingLimitResponse struct {
}

// When an external control system sends a signal to release a previously imposed charging limit to a Charging Station,
// the Charging Station adjusts the charging speed of the ongoing transaction(s).
// If the charging limit changed by more than: LimitChangeSignificance, the Charging Station sends a NotifyChargingLimitRequest message to CSMS with optionally the set charging
// limit/schedule.
//
// The CSMS responds with NotifyChargingLimitResponse to the Charging Station.
//
// If the charging rate changes by more than: LimitChangeSignificance, the Charging Station sends a TransactionEventRequest message to inform the CSMS.
type NotifyChargingLimitFeature struct{}

func (f NotifyChargingLimitFeature) GetFeatureName() string {
	return NotifyChargingLimitFeatureName
}

func (f NotifyChargingLimitFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(NotifyChargingLimitRequest{})
}

func (f NotifyChargingLimitFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(NotifyChargingLimitResponse{})
}

func (r NotifyChargingLimitRequest) GetFeatureName() string {
	return NotifyChargingLimitFeatureName
}

func (c NotifyChargingLimitResponse) GetFeatureName() string {
	return NotifyChargingLimitFeatureName
}

// Creates a new NotifyChargingLimitRequest, containing all required fields. Optional fields may be set afterwards.
func NewNotifyChargingLimitRequest(chargingLimit ChargingLimit) *NotifyChargingLimitRequest {
	return &NotifyChargingLimitRequest{ChargingLimit: chargingLimit}
}

// Creates a new NotifyChargingLimitResponse, which doesn't contain any required or optional fields.
func NewNotifyChargingLimitResponse() *NotifyChargingLimitResponse {
	return &NotifyChargingLimitResponse{}
}
