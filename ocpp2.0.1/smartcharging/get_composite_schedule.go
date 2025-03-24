package smartcharging

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"gopkg.in/go-playground/validator.v9"
)

// -------------------- Get Composite Schedule (CSMS -> CS) --------------------

const GetCompositeScheduleFeatureName = "GetCompositeSchedule"

// Status reported in GetCompositeScheduleResponse.
type GetCompositeScheduleStatus string

const (
	GetCompositeScheduleStatusAccepted GetCompositeScheduleStatus = "Accepted"
	GetCompositeScheduleStatusRejected GetCompositeScheduleStatus = "Rejected"
)

func isValidGetCompositeScheduleStatus(fl validator.FieldLevel) bool {
	status := GetCompositeScheduleStatus(fl.Field().String())
	switch status {
	case GetCompositeScheduleStatusAccepted, GetCompositeScheduleStatusRejected:
		return true
	default:
		return false
	}
}

type CompositeSchedule struct {
	StartDateTime    *types.DateTime         `json:"startDateTime,omitempty" validate:"omitempty"`
	ChargingSchedule *types.ChargingSchedule `json:"chargingSchedule,omitempty" validate:"omitempty"`
}

// The field definition of the GetCompositeSchedule request payload sent by the CSMS to the Charging System.
type GetCompositeScheduleRequest struct {
	Duration         int                        `json:"duration" validate:"gte=0"`
	ChargingRateUnit types.ChargingRateUnitType `json:"chargingRateUnit,omitempty" validate:"omitempty,chargingRateUnit201"`
	EvseID           int                        `json:"evseId" validate:"gte=0"`
}

// This field definition of the GetCompositeSchedule response payload, sent by the Charging System to the CSMS in response to a GetCompositeScheduleRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type GetCompositeScheduleResponse struct {
	Status     GetCompositeScheduleStatus `json:"status" validate:"required,getCompositeScheduleStatus"`
	StatusInfo *types.StatusInfo          `json:"statusInfo,omitempty" validate:"omitempty"`
	Schedule   *CompositeSchedule         `json:"schedule,omitempty" validate:"omitempty"`
}

// The CSMS MAY request the Charging System to report the Composite Charging Schedule by sending a GetCompositeScheduleRequest.
// The Charging System SHALL calculate the Composite Charging Schedule intervals, from the moment the request payload is received: Time X, up to X + Duration, and send them in the GetCompositeScheduleResponse to the CSMS.
// The reported schedule, in the GetCompositeScheduleResponse payload, is the result of the calculation of all active schedules and possible local limits present in the Charging System.
// If the ConnectorId in the request is set to '0', the Charging System SHALL report the total expected power or current the Charging System expects to consume from the grid during the requested time period.
// If the Charging System is not able to report the requested schedule, for instance if the connectorId is unknown, it SHALL respond with a status Rejected.
type GetCompositeScheduleFeature struct{}

func (f GetCompositeScheduleFeature) GetFeatureName() string {
	return GetCompositeScheduleFeatureName
}

func (f GetCompositeScheduleFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(GetCompositeScheduleRequest{})
}

func (f GetCompositeScheduleFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(GetCompositeScheduleResponse{})
}

func (r GetCompositeScheduleRequest) GetFeatureName() string {
	return GetCompositeScheduleFeatureName
}

func (c GetCompositeScheduleResponse) GetFeatureName() string {
	return GetCompositeScheduleFeatureName
}

// Creates a new GetCompositeScheduleRequest, containing all required fields. Optional fields may be set afterwards.
func NewGetCompositeScheduleRequest(duration int, evseId int) *GetCompositeScheduleRequest {
	return &GetCompositeScheduleRequest{Duration: duration, EvseID: evseId}
}

// Creates a new GetCompositeScheduleResponse, containing all required fields. Optional fields may be set afterwards.
func NewGetCompositeScheduleResponse(status GetCompositeScheduleStatus) *GetCompositeScheduleResponse {
	return &GetCompositeScheduleResponse{Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("getCompositeScheduleStatus", isValidGetCompositeScheduleStatus)
}
