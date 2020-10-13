package smartcharging

import (
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
	"reflect"
)

// -------------------- Notify EV Charging Schedule (CS -> CSMS) --------------------

const NotifyEVChargingScheduleFeatureName = "NotifyEVChargingSchedule"

// The field definition of the NotifyEVChargingSchedule request payload sent by the Charging Station to the CSMS.
type NotifyEVChargingScheduleRequest struct {
	TimeBase         *types.DateTime        `json:"timeBase" validate:"required"`
	EvseID           int                    `json:"evseId" validate:"gt=0"`
	ChargingSchedule types.ChargingSchedule `json:"chargingSchedule" validate:"required,dive"`
}

// This field definition of the NotifyEVChargingSchedule response payload, sent by the CSMS to the Charging Station in response to a NotifyEVChargingScheduleRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type NotifyEVChargingScheduleResponse struct {
	Status     types.GenericStatus `json:"status" validate:"required,genericStatus"`
	StatusInfo *types.StatusInfo   `json:"statusInfo,omitempty" validate:"omitempty,dive"` // Detailed status information.
}

// A power renegotiation, either initiated by the EV or by the CSMS, may involve the EV providing a power profile.
// If a charging profile was provided, after receiving a PowerDeliveryResponse from the CSMS,
// the Charging Station will send a NotifyEVChargingScheduleRequest to the CSMS.
//
// The CSMS replies to the Charging Station with a NotifyEVChargingScheduleResponse.
type NotifyEVChargingScheduleFeature struct{}

func (f NotifyEVChargingScheduleFeature) GetFeatureName() string {
	return NotifyEVChargingScheduleFeatureName
}

func (f NotifyEVChargingScheduleFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(NotifyEVChargingScheduleRequest{})
}

func (f NotifyEVChargingScheduleFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(NotifyEVChargingScheduleResponse{})
}

func (r NotifyEVChargingScheduleRequest) GetFeatureName() string {
	return NotifyEVChargingScheduleFeatureName
}

func (c NotifyEVChargingScheduleResponse) GetFeatureName() string {
	return NotifyEVChargingScheduleFeatureName
}

// Creates a new NotifyEVChargingScheduleRequest, containing all required fields. Optional fields may be set afterwards.
func NewNotifyEVChargingScheduleRequest(timeBase *types.DateTime, evseID int, chargingSchedule types.ChargingSchedule) *NotifyEVChargingScheduleRequest {
	return &NotifyEVChargingScheduleRequest{TimeBase: timeBase, EvseID: evseID, ChargingSchedule: chargingSchedule}
}

// Creates a new NotifyEVChargingScheduleResponse, containing all required fields. Optional fields may be set afterwards.
func NewNotifyEVChargingScheduleResponse(status types.GenericStatus) *NotifyEVChargingScheduleResponse {
	return &NotifyEVChargingScheduleResponse{Status: status}
}
