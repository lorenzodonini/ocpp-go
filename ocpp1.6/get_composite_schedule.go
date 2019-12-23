package ocpp16

import (
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Get Composite Schedule (CS -> CP) --------------------

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

// The field definition of the GetCompositeSchedule request payload sent by the Central System to the Charge Point.
type GetCompositeScheduleRequest struct {
	ConnectorId      int                  `json:"connectorId" validate:"gte=0"`
	Duration         int                  `json:"duration" validate:"gte=0"`
	ChargingRateUnit ChargingRateUnitType `json:"chargingRateUnit,omitempty" validate:"omitempty,chargingRateUnit"`
}

// This field definition of the GetCompositeSchedule confirmation payload, sent by the Charge Point to the Central System in response to a GetCompositeScheduleRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type GetCompositeScheduleConfirmation struct {
	Status           GetCompositeScheduleStatus `json:"status" validate:"required,chargingProfileStatus"`
	ConnectorId      int                        `json:"connectorId,omitempty" validate:"omitempty,gt=0"`
	ScheduleStart    *DateTime                  `json:"scheduleStart,omitempty"`
	ChargingSchedule *ChargingSchedule          `json:"chargingSchedule,omitempty" validate:"omitempty,dive"`
}

// The Central System MAY request the Charge Point to report the Composite Charging Schedule by sending a GetCompositeScheduleRequest.
// The Charge Point SHALL calculate the Composite Charging Schedule intervals, from the moment the request payload is received: Time X, up to X + Duration, and send them in the GetCompositeScheduleConfirmation to the Central System.
// The reported schedule, in the GetCompositeScheduleConfirmation payload, is the result of the calculation of all active schedules and possible local limits present in the Charge Point.
// If the ConnectorId in the request is set to '0', the Charge Point SHALL report the total expected power or current the Charge Point expects to consume from the grid during the requested time period.
// If the Charge Point is not able to report the requested schedule, for instance if the connectorId is unknown, it SHALL respond with a status Rejected.
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
func NewGetCompositeScheduleRequest(connectorId int, duration int) *GetCompositeScheduleRequest {
	return &GetCompositeScheduleRequest{ConnectorId: connectorId, Duration: duration}
}

// Creates a new GetCompositeScheduleConfirmation, containing all required fields. Optional fields may be set afterwards.
func NewGetCompositeScheduleConfirmation(status GetCompositeScheduleStatus) *GetCompositeScheduleConfirmation {
	return &GetCompositeScheduleConfirmation{Status: status}
}

func init() {
	_ = Validate.RegisterValidation("chargingProfileStatus", isValidChargingProfileStatus)
}
