package ocpp2

import (
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Get Monitoring Report (CSMS -> CS) --------------------

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
	RequestID          *int                     `json:"requestId,omitempty" validate:"omitempty,gte=0"`               // The Id of the request.
	MonitoringCriteria []MonitoringCriteriaType `json:"monitoringCriteria,omitempty" validate:"omitempty,max=3,dive,monitoringCriteria"` // This field contains criteria for components for which a monitoring report is requested.
	ComponentVariable  []ComponentVariable      `json:"componentVariable,omitempty" validate:"omitempty,dive"`        // This field specifies the components and variables for which a monitoring report is requested.
}

// This field definition of the GetMonitoringReport confirmation payload, sent by the Charging Station to the CSMS in response to a GetMonitoringReportRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type GetMonitoringReportConfirmation struct {
	Status GenericDeviceModelStatus `json:"status" validate:"required,genericDeviceModelStatus"` // This field indicates whether the Charging Station was able to accept the request.
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

func (f GetMonitoringReportFeature) GetConfirmationType() reflect.Type {
	return reflect.TypeOf(GetMonitoringReportConfirmation{})
}

func (r GetMonitoringReportRequest) GetFeatureName() string {
	return GetMonitoringReportFeatureName
}

func (c GetMonitoringReportConfirmation) GetFeatureName() string {
	return GetMonitoringReportFeatureName
}

// Creates a new GetMonitoringReportRequest. All fields are optional and may be set afterwards.
func NewGetMonitoringReportRequest() *GetMonitoringReportRequest {
	return &GetMonitoringReportRequest{}
}

// Creates a new GetMonitoringReportConfirmation, containing all required fields. There are no optional fields for this message.
func NewGetMonitoringReportConfirmation(status GenericDeviceModelStatus) *GetMonitoringReportConfirmation {
	return &GetMonitoringReportConfirmation{Status: status}
}

func init() {
	_ = Validate.RegisterValidation("monitoringCriteria", isValidMonitoringCriteriaType)
}
