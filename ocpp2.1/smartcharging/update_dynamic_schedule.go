package smartcharging

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
)

// -------------------- UpdateDynamicSchedule (CSMS -> CS) --------------------

const UpdateDynamicScheduleFeatureName = "UpdateDynamicSchedule"

// The field definition of the UpdateDynamicScheduleRequest request payload sent by the CSMS to the Charging Station.
type UpdateDynamicScheduleRequest struct {
	ChargingProfileId int                        `json:"chargingProfileId" validate:"required"`
	ScheduleUpdate    ChargingScheduleUpdateType `json:"scheduleUpdate" validate:"required,dive"`
}

// This field definition of the UpdateDynamicScheduleResponse response payload, sent by the Charging Station to the CSMS.
type UpdateDynamicScheduleResponse struct {
	Status     ChargingProfileStatus `json:"status" validate:"required,chargingProfileStatus21"`
	StatusInfo *types.StatusInfo     `json:"statusInfo,omitempty" validate:"omitempty,dive"`
}

type UpdateDynamicScheduleFeature struct{}

func (f UpdateDynamicScheduleFeature) GetFeatureName() string {
	return UpdateDynamicScheduleFeatureName
}

func (f UpdateDynamicScheduleFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(UpdateDynamicScheduleRequest{})
}

func (f UpdateDynamicScheduleFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(UpdateDynamicScheduleResponse{})
}

func (r UpdateDynamicScheduleRequest) GetFeatureName() string {
	return UpdateDynamicScheduleFeatureName
}

func (c UpdateDynamicScheduleResponse) GetFeatureName() string {
	return UpdateDynamicScheduleFeatureName
}

// Creates a new UpdateDynamicScheduleRequest, containing all required fields.
func NewUpdateDynamicScheduleRequest(chargingProfileId int, scheduleUpdate ChargingScheduleUpdateType) *UpdateDynamicScheduleRequest {
	return &UpdateDynamicScheduleRequest{ChargingProfileId: chargingProfileId, ScheduleUpdate: scheduleUpdate}
}

// Creates a new UpdateDynamicScheduleResponse, containing all required fields. Optional fields may be set afterwards.
func NewUpdateDynamicScheduleResponse(status ChargingProfileStatus) *UpdateDynamicScheduleResponse {
	return &UpdateDynamicScheduleResponse{Status: status}
}