package diagnostics

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
	"gopkg.in/go-playground/validator.v9"
)

// -------------------- Set Variable Monitoring (CSMS -> CS) --------------------

const SetVariableMonitoringFeatureName = "SetVariableMonitoring"

// Status contained inside a SetMonitoringResult struct.
type SetMonitoringStatus string

const (
	SetMonitoringStatusAccepted               SetMonitoringStatus = "Accepted"
	SetMonitoringStatusUnknownComponent       SetMonitoringStatus = "UnknownComponent"
	SetMonitoringStatusUnknownVariable        SetMonitoringStatus = "UnknownVariable"
	SetMonitoringStatusUnsupportedMonitorType SetMonitoringStatus = "UnsupportedMonitorType"
	SetMonitoringStatusRejected               SetMonitoringStatus = "Rejected"
	SetMonitoringStatusDuplicate              SetMonitoringStatus = "Duplicate"
)

func isValidSetMonitoringStatus(fl validator.FieldLevel) bool {
	status := SetMonitoringStatus(fl.Field().String())
	switch status {
	case SetMonitoringStatusAccepted, SetMonitoringStatusUnknownComponent, SetMonitoringStatusUnknownVariable, SetMonitoringStatusUnsupportedMonitorType, SetMonitoringStatusRejected, SetMonitoringStatusDuplicate:
		return true
	default:
		return false
	}
}

// Hold parameters of a SetVariableMonitoring request.
type SetMonitoringData struct {
	ID                  *int                       `json:"id,omitempty" validate:"omitempty"`    // An id SHALL only be given to replace an existing monitor. The Charging Station handles the generation of id’s for new monitors.
	Transaction         bool                       `json:"transaction,omitempty"`                // Monitor only active when a transaction is ongoing on a component relevant to this transaction.
	Value               float64                    `json:"value"`                                // Value for threshold or delta monitoring. For Periodic or PeriodicClockAligned this is the interval in seconds.
	Type                MonitorType                `json:"type" validate:"required,monitorType"` // The type of this monitor, e.g. a threshold, delta or periodic monitor.
	Severity            int                        `json:"severity" validate:"min=0,max=9"`      // The severity that will be assigned to an event that is triggered by this monitor. The severity range is 0-9, with 0 as the highest and 9 as the lowest severity level.
	Component           types.Component            `json:"component" validate:"required"`        // Component for which monitor is set.
	Variable            types.Variable             `json:"variable" validate:"required"`         // Variable for which monitor is set.
	PeriodicEventStream *PeriodicEventStreamParams `json:"periodicEventStream,omitempty" validate:"omitempty,dive"`
}

type PeriodicEventStreamParams struct {
	Interval *int `json:"interval,omitempty" validate:"omitempty,gte=0"` // Interval in seconds for periodic monitoring.
	Values   *int `json:"Values,omitempty" validate:"omitempty,gte=0"`
}

// Holds the result of SetVariableMonitoring request.
type SetMonitoringResult struct {
	ID         *int                `json:"id,omitempty" validate:"omitempty"`                // Id given to the VariableMonitor by the Charging Station. The Id is only returned when status is accepted.
	Status     SetMonitoringStatus `json:"status" validate:"required,setMonitoringStatus21"` // Status is OK if a value could be returned. Otherwise this will indicate the reason why a value could not be returned.
	Type       MonitorType         `json:"type" validate:"required,monitorType"`             // The type of this monitor, e.g. a threshold, delta or periodic monitor.
	Severity   int                 `json:"severity" validate:"min=0,max=9"`                  // The severity that will be assigned to an event that is triggered by this monitor. The severity range is 0-9, with 0 as the highest and 9 as the lowest severity level.
	Component  types.Component     `json:"component" validate:"required"`                    // Component for which status is returned.
	Variable   types.Variable      `json:"variable" validate:"required"`                     // Variable for which status is returned.
	StatusInfo *types.StatusInfo   `json:"statusInfo,omitempty" validate:"omitempty"`        // Detailed status information.
}

// The field definition of the SetVariableMonitoring request payload sent by the CSMS to the Charging Station.
type SetVariableMonitoringRequest struct {
	MonitoringData []SetMonitoringData `json:"setMonitoringData" validate:"required,min=1,dive"` // List of MonitoringData containing monitoring settings.
}

// This field definition of the SetVariableMonitoring response payload, sent by the Charging Station to the CSMS in response to a SetVariableMonitoringRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type SetVariableMonitoringResponse struct {
	MonitoringResult []SetMonitoringResult `json:"setMonitoringResult" validate:"required,min=1,dive"` //  List of result statuses per monitor.
}

// The CSMS may request the Charging Station to set monitoring triggers on Variables. Multiple triggers can be
// set for upper or lower thresholds, delta changes or periodic reporting.
//
// To achieve this, the CSMS sends a SetVariableMonitoringRequest to the Charging Station.
// The Charging Station responds with a SetVariableMonitoringResponse.
type SetVariableMonitoringFeature struct{}

func (f SetVariableMonitoringFeature) GetFeatureName() string {
	return SetVariableMonitoringFeatureName
}

func (f SetVariableMonitoringFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(SetVariableMonitoringRequest{})
}

func (f SetVariableMonitoringFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(SetVariableMonitoringResponse{})
}

func (r SetVariableMonitoringRequest) GetFeatureName() string {
	return SetVariableMonitoringFeatureName
}

func (c SetVariableMonitoringResponse) GetFeatureName() string {
	return SetVariableMonitoringFeatureName
}

// Creates a new SetVariableMonitoringRequest, containing all required fields. There are no optional fields for this message.
func NewSetVariableMonitoringRequest(data []SetMonitoringData) *SetVariableMonitoringRequest {
	return &SetVariableMonitoringRequest{MonitoringData: data}
}

// Creates a new SetVariableMonitoringResponse, containing all required fields. There are no optional fields for this message.
func NewSetVariableMonitoringResponse(result []SetMonitoringResult) *SetVariableMonitoringResponse {
	return &SetVariableMonitoringResponse{MonitoringResult: result}
}

func init() {
	_ = types.Validate.RegisterValidation("setMonitoringStatus21", isValidSetMonitoringStatus)
}
