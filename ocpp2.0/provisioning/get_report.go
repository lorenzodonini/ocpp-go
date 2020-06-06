package provisioning

import (
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Get Base Report (CSMS -> CS) --------------------

const GetReportFeatureName = "GetReport"

// ComponentCriterion indicates the criterion for components requested in GetReportRequest.
type ComponentCriterion string

const (
	ComponentCriterionActive    ComponentCriterion = "Active"
	ComponentCriterionAvailable ComponentCriterion = "Available"
	ComponentCriterionEnabled   ComponentCriterion = "Enabled"
	ComponentCriterionProblem   ComponentCriterion = "Problem"
)

func isValidComponentCriterion(fl validator.FieldLevel) bool {
	status := ComponentCriterion(fl.Field().String())
	switch status {
	case ComponentCriterionActive, ComponentCriterionAvailable, ComponentCriterionEnabled, ComponentCriterionProblem:
		return true
	default:
		return false
	}
}

// The field definition of the GetReport request payload sent by the CSMS to the Charging Station.
type GetReportRequest struct {
	RequestID         *int                      `json:"requestId,omitempty" validate:"omitempty,gte=0"`
	ComponentCriteria []ComponentCriterion      `json:"componentCriteria,omitempty" validate:"omitempty,max=4,dive,componentCriterion"`
	ComponentVariable []types.ComponentVariable `json:"componentVariable,omitempty" validate:"omitempty,dive"`
}

// This field definition of the GetReport response payload, sent by the Charging Station to the CSMS in response to a GetReportRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type GetReportResponse struct {
	Status types.GenericDeviceModelStatus `json:"status" validate:"required,genericDeviceModelStatus"`
}

// The CSO may trigger the CSMS to request a report from a Charging Station.
// The CSMS shall then request a Charging Station to send a report of all Components and Variables limited to those that match ComponentCriteria and/or the list of ComponentVariables.
// The Charging Station responds with GetReportResponse.
// The result will be returned asynchronously in one or more NotifyReportRequest messages (one for each report part).
type GetReportFeature struct{}

func (f GetReportFeature) GetFeatureName() string {
	return GetReportFeatureName
}

func (f GetReportFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(GetReportRequest{})
}

func (f GetReportFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(GetReportResponse{})
}

func (r GetReportRequest) GetFeatureName() string {
	return GetReportFeatureName
}

func (c GetReportResponse) GetFeatureName() string {
	return GetReportFeatureName
}

// Creates a new GetReportRequest, containing all required fields. Optional fields may be set afterwards.
func NewGetReportRequest() *GetReportRequest {
	return &GetReportRequest{}
}

// Creates a new GetReportResponse, containing all required fields. There are no optional fields for this message.
func NewGetReportResponse(status types.GenericDeviceModelStatus) *GetReportResponse {
	return &GetReportResponse{Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("componentCriterion", isValidComponentCriterion)
}
