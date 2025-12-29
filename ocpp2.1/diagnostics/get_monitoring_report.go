package diagnostics

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
	"gopkg.in/go-playground/validator.v9"
)

// -------------------- Get Monitoring Report (CSMS -> CS) --------------------

const GetMonitoringReportFeatureName = "GetMonitoringReport"

// Monitoring criteria contained in GetMonitoringReportRequest.
type MonitoringCriteriaType string

const (
	MonitoringCriteriaThresholdMonitoring MonitoringCriteriaType = "ThresholdMonitoring"
	MonitoringCriteriaDeltaMonitoring     MonitoringCriteriaType = "DeltaMonitoring"
	MonitoringCriteriaPeriodicMonitoring  MonitoringCriteriaType = "PeriodicMonitoring"
)

func isValidMonitoringCriteriaType(fl validator.FieldLevel) bool {
	status := MonitoringCriteriaType(fl.Field().String())
	switch status {
	case MonitoringCriteriaThresholdMonitoring, MonitoringCriteriaDeltaMonitoring, MonitoringCriteriaPeriodicMonitoring:
		return true
	default:
		return false
	}
}

// The field definition of the GetMonitoringReport request payload sent by the CSMS to the Charging Station.
type GetMonitoringReportRequest struct {
	RequestID          int                       `json:"requestId" validate:"gte=0"`                                                        // The Id of the request.
	MonitoringCriteria []MonitoringCriteriaType  `json:"monitoringCriteria,omitempty" validate:"omitempty,max=3,dive,monitoringCriteria21"` // This field contains criteria for components for which a monitoring report is requested.
	ComponentVariable  []types.ComponentVariable `json:"componentVariable,omitempty" validate:"omitempty,dive"`                             // This field specifies the components and variables for which a monitoring report is requested.
}

// This field definition of the GetMonitoringReport response payload, sent by the Charging Station to the CSMS in response to a GetMonitoringReportRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type GetMonitoringReportResponse struct {
	Status     types.GenericDeviceModelStatus `json:"status" validate:"required,genericDeviceModelStatus21"` // This field indicates whether the Charging Station was able to accept the request.
	StatusInfo *types.StatusInfo              `json:"statusInfo,omitempty" validate:"omitempty"`
}

// A CSMS can request the Charging Station to send a report about configured monitoring settings per component and variable.
// Optionally, this list can be filtered on monitoringCriteria and componentVariables.
// The CSMS sends a GetMonitoringReportRequest to the Charging Station.
// The Charging Station then responds with a GetMonitoringReportResponse.
// Asynchronously, the Charging Station will then send a NotifyMonitoringReportRequest to the CSMS for each report part.
type GetMonitoringReportFeature struct{}

func (f GetMonitoringReportFeature) GetFeatureName() string {
	return GetMonitoringReportFeatureName
}

func (f GetMonitoringReportFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(GetMonitoringReportRequest{})
}

func (f GetMonitoringReportFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(GetMonitoringReportResponse{})
}

func (r GetMonitoringReportRequest) GetFeatureName() string {
	return GetMonitoringReportFeatureName
}

func (c GetMonitoringReportResponse) GetFeatureName() string {
	return GetMonitoringReportFeatureName
}

// Creates a new GetMonitoringReportRequest. All fields are optional and may be set afterwards.
func NewGetMonitoringReportRequest() *GetMonitoringReportRequest {
	return &GetMonitoringReportRequest{}
}

// Creates a new GetMonitoringReportResponse, containing all required fields. There are no optional fields for this message.
func NewGetMonitoringReportResponse(status types.GenericDeviceModelStatus) *GetMonitoringReportResponse {
	return &GetMonitoringReportResponse{Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("monitoringCriteria21", isValidMonitoringCriteriaType)
}
