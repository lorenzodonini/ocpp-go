package smartcharging

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
)

// -------------------- PullDynamicScheduleUpdate (CS -> CSMS) --------------------

const PullDynamicScheduleUpdateFeatureName = "PullDynamicScheduleUpdate"

// ChargingScheduleUpdateType contains the fields that can be updated in a dynamic charging schedule.
type ChargingScheduleUpdateType struct {
	Limit              *float64 `json:"limit,omitempty" validate:"omitempty"`
	LimitL2            *float64 `json:"limit_L2,omitempty" validate:"omitempty"`
	LimitL3            *float64 `json:"limit_L3,omitempty" validate:"omitempty"`
	DischargeLimit     *float64 `json:"dischargeLimit,omitempty" validate:"omitempty,lte=0"`
	DischargeLimitL2   *float64 `json:"dischargeLimit_L2,omitempty" validate:"omitempty,lte=0"`
	DischargeLimitL3   *float64 `json:"dischargeLimit_L3,omitempty" validate:"omitempty,lte=0"`
	SetPoint           *float64 `json:"setpoint,omitempty" validate:"omitempty"`
	SetPointL2         *float64 `json:"setpoint_L2,omitempty" validate:"omitempty"`
	SetPointL3         *float64 `json:"setpoint_L3,omitempty" validate:"omitempty"`
	SetpointReactive   *float64 `json:"setpointReactive,omitempty" validate:"omitempty"`
	SetpointReactiveL2 *float64 `json:"setpointReactive_L2,omitempty" validate:"omitempty"`
	SetpointReactiveL3 *float64 `json:"setpointReactive_L3,omitempty" validate:"omitempty"`
}

// The field definition of the PullDynamicScheduleUpdateRequest request payload sent by the Charging Station to the CSMS.
type PullDynamicScheduleUpdateRequest struct {
	ChargingProfileId int `json:"chargingProfileId" validate:"required"`
}

// This field definition of the PullDynamicScheduleUpdateResponse response payload, sent by the CSMS to the Charging Station.
type PullDynamicScheduleUpdateResponse struct {
	Status         ChargingProfileStatus       `json:"status" validate:"required,chargingProfileStatus21"`
	StatusInfo     *types.StatusInfo           `json:"statusInfo,omitempty" validate:"omitempty,dive"`
	ScheduleUpdate *ChargingScheduleUpdateType `json:"scheduleUpdate,omitempty" validate:"omitempty,dive"`
}

type PullDynamicScheduleUpdateFeature struct{}

func (f PullDynamicScheduleUpdateFeature) GetFeatureName() string {
	return PullDynamicScheduleUpdateFeatureName
}

func (f PullDynamicScheduleUpdateFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(PullDynamicScheduleUpdateRequest{})
}

func (f PullDynamicScheduleUpdateFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(PullDynamicScheduleUpdateResponse{})
}

func (r PullDynamicScheduleUpdateRequest) GetFeatureName() string {
	return PullDynamicScheduleUpdateFeatureName
}

func (c PullDynamicScheduleUpdateResponse) GetFeatureName() string {
	return PullDynamicScheduleUpdateFeatureName
}

// Creates a new PullDynamicScheduleUpdateRequest, containing all required fields.
func NewPullDynamicScheduleUpdateRequest(chargingProfileId int) *PullDynamicScheduleUpdateRequest {
	return &PullDynamicScheduleUpdateRequest{ChargingProfileId: chargingProfileId}
}

// Creates a new PullDynamicScheduleUpdateResponse, containing all required fields. Optional fields may be set afterwards.
func NewPullDynamicScheduleUpdateResponse(status ChargingProfileStatus) *PullDynamicScheduleUpdateResponse {
	return &PullDynamicScheduleUpdateResponse{Status: status}
}