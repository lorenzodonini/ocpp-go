package smartcharging

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Get Composite Schedule (CS -> CP) --------------------

const GetCompositeScheduleFeatureName = "GetCompositeSchedule"

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
	ConnectorId      int                        `json:"connectorId" validate:"gte=0"`
	Duration         int                        `json:"duration" validate:"gte=0"`
	ChargingRateUnit types.ChargingRateUnitType `json:"chargingRateUnit,omitempty" validate:"omitempty,chargingRateUnit"`
}

// This field definition of the GetCompositeSchedule confirmation payload, sent by the Charge Point to the Central System in response to a GetCompositeScheduleRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type GetCompositeScheduleConfirmation struct {
	Status           GetCompositeScheduleStatus `json:"status" validate:"required,compositeScheduleStatus"`
	ConnectorId      *int                       `json:"connectorId,omitempty" validate:"omitempty,gt=0"`
	ScheduleStart    *types.DateTime            `json:"scheduleStart,omitempty"`
	ChargingSchedule *types.ChargingSchedule    `json:"chargingSchedule,omitempty" validate:"omitempty"`
}

// The CSMS requests the Charging Station to report the Composite Charging Schedule by sending a GetCompositeScheduleRequest.
// The Charging Station calculates the schedule, according to the parameters specified in the request.
// The composite schedule is the result of the calculation of all active schedules and possible local limits present in the Charging Station.
// The Charging Station responds with a GetCompositeScheduleResponse with the status and ChargingSchedule.
// If the Charging Station is not able to report the requested schedule, for instance if the evseID is unknown, it SHALL respond with a status Rejected.
type GetCompositeScheduleFeature struct{}

func (f GetCompositeScheduleFeature) GetFeatureName() string {
	return GetCompositeScheduleFeatureName
}

func (f GetCompositeScheduleFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(GetCompositeScheduleRequest{})
}

func (f GetCompositeScheduleFeature) GetResponseType() reflect.Type {
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
	_ = types.Validate.RegisterValidation("compositeScheduleStatus", isValidGetCompositeScheduleStatus)
}
