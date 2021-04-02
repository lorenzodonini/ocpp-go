package diagnostics

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
)

// -------------------- Set Monitoring Level (CSMS -> CS) --------------------

const SetMonitoringLevelFeatureName = "SetMonitoringLevel"

// The field definition of the SetMonitoringLevel request payload sent by the CSMS to the Charging Station.
type SetMonitoringLevelRequest struct {
	// Severity levels have the following meaning:
	//
	//  - 0 Danger:
	// Indicates lives are potentially in danger. Urgent attention
	// is needed and action should be taken immediately.
	//  - 1 Hardware Failure:
	// Indicates that the Charging Station is unable to continue regular operations due to Hardware issues.
	//  - 2 System Failure:
	// Indicates that the Charging Station is unable to continue regular operations due to software or minor hardware
	// issues.
	//  - 3 Critical:
	// Indicates a critical error.
	//  - 4 Error:
	// Indicates a non-urgent error.
	//  - 5 Alert:
	// Indicates an alert event. Default severity for any type of monitoring event.
	//  - 6 Warning:
	// Indicates a warning event.
	//  - 7 Notice:
	// Indicates an unusual event.
	//  - 8 Informational:
	// Indicates a regular operational event. May be used for reporting, measuring throughput, etc.
	//  - 9 Debug:
	// Indicates information useful to developers for debugging, not useful during operations.
	Severity int `json:"severity" validate:"min=0,max=9"`
}

// This field definition of the SetMonitoringLevel response payload, sent by the Charging Station to the CSMS in response to a SetMonitoringLevelRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type SetMonitoringLevelResponse struct {
	Status     types.GenericDeviceModelStatus `json:"status" validate:"required,genericDeviceModelStatus"` // Indicates whether the Charging Station was able to accept the request.
	StatusInfo *types.StatusInfo              `json:"statusInfo,omitempty" validate:"omitempty"`           // Detailed status information.
}

// It may be desirable to restrict the reporting of monitoring events, to only those monitors with a
// severity number lower than or equal to a certain severity. For example when the data-traffic between
// Charging Station and CSMS needs to limited for some reason.
//
// The CSMS can control which events it will to be notified of by the Charging Station with the
// SetMonitoringLevelRequest message. The charging station responds with a SetMonitoringLevelResponse.
// Monitoring events, reported later on via NotifyEventRequest messages,
// will be restricted according to the set monitoring level.
type SetMonitoringLevelFeature struct{}

func (f SetMonitoringLevelFeature) GetFeatureName() string {
	return SetMonitoringLevelFeatureName
}

func (f SetMonitoringLevelFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(SetMonitoringLevelRequest{})
}

func (f SetMonitoringLevelFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(SetMonitoringLevelResponse{})
}

func (r SetMonitoringLevelRequest) GetFeatureName() string {
	return SetMonitoringLevelFeatureName
}

func (c SetMonitoringLevelResponse) GetFeatureName() string {
	return SetMonitoringLevelFeatureName
}

// Creates a new SetMonitoringLevelRequest, containing all required fields. There are no optional fields for this message.
func NewSetMonitoringLevelRequest(severity int) *SetMonitoringLevelRequest {
	return &SetMonitoringLevelRequest{Severity: severity}
}

// Creates a new SetMonitoringLevelResponse, containing all required fields. Optional fields may be set afterwards.
func NewSetMonitoringLevelResponse(status types.GenericDeviceModelStatus) *SetMonitoringLevelResponse {
	return &SetMonitoringLevelResponse{Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("monitoringBase", isValidMonitoringBase)
}
