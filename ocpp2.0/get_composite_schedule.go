package ocpp2

import (
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Get Composite Schedule (CSMS -> CS) --------------------

// Status reported in GetCompositeScheduleConfirmation.
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
	StartDateTime    *DateTime         `json:"startDateTime,omitempty" validate:"omitempty"`
	ChargingSchedule *ChargingSchedule `json:"chargingSchedule,omitempty" validate:"omitempty"`
}

// The field definition of the GetCompositeSchedule request payload sent by the CSMS to the Charging System.
type GetCompositeScheduleRequest struct {
	Duration         int                  `json:"duration" validate:"gte=0"`
	ChargingRateUnit ChargingRateUnitType `json:"chargingRateUnit,omitempty" validate:"omitempty,chargingRateUnit"`
	EvseID           int                  `json:"evseId" validate:"gte=0"`
}

// This field definition of the GetCompositeSchedule confirmation payload, sent by the Charging System to the CSMS in response to a GetCompositeScheduleRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type GetCompositeScheduleConfirmation struct {
	Status   GetCompositeScheduleStatus `json:"status" validate:"required,getCompositeScheduleStatus"`
	EvseID   int                        `json:"evseId" validate:"gte=0"`
	Schedule *CompositeSchedule         `json:"schedule,omitempty" validate:"omitempty"`
}

// The CSMS MAY request the Charging System to report the Composite Charging Schedule by sending a GetCompositeScheduleRequest.
// The Charging System SHALL calculate the Composite Charging Schedule intervals, from the moment the request payload is received: Time X, up to X + Duration, and send them in the GetCompositeScheduleConfirmation to the CSMS.
// The reported schedule, in the GetCompositeScheduleConfirmation payload, is the result of the calculation of all active schedules and possible local limits present in the Charging System.
// If the ConnectorId in the request is set to '0', the Charging System SHALL report the total expected power or current the Charging System expects to consume from the grid during the requested time period.
// If the Charging System is not able to report the requested schedule, for instance if the connectorId is unknown, it SHALL respond with a status Rejected.
type GetCompositeScheduleFeature struct{}

func (f GetCompositeScheduleFeature) GetFeatureName() string {
	return GetCompositeScheduleFeatureName
}

func (f GetCompositeScheduleFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(GetCompositeScheduleRequest{})
}

func (f GetCompositeScheduleFeature) GetConfirmationType() reflect.Type {
	return reflect.TypeOf(GetCompositeScheduleConfirmation{})
}

func (r GetCompositeScheduleRequest) GetFeatureName() string {
	return GetCompositeScheduleFeatureName
}

func (c GetCompositeScheduleConfirmation) GetFeatureName() string {
	return GetCompositeScheduleFeatureName
}

// Creates a new GetCompositeScheduleRequest, containing all required fields. Optional fields may be set afterwards.
func NewGetCompositeScheduleRequest(duration int, evseId int) *GetCompositeScheduleRequest {
	return &GetCompositeScheduleRequest{Duration: duration, EvseID: evseId}
}

// Creates a new GetCompositeScheduleConfirmation, containing all required fields. Optional fields may be set afterwards.
func NewGetCompositeScheduleConfirmation(status GetCompositeScheduleStatus, evseId int) *GetCompositeScheduleConfirmation {
	return &GetCompositeScheduleConfirmation{Status: status, EvseID: evseId}
}

func init() {
	_ = Validate.RegisterValidation("getCompositeScheduleStatus", isValidGetCompositeScheduleStatus)
}
