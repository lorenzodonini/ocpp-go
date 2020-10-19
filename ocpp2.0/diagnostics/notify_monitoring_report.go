package diagnostics

import (
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
	"reflect"
)

// -------------------- Notify Monitoring Report (CS -> CSMS) --------------------

const NotifyMonitoringReportFeatureName = "NotifyMonitoringReport"

// VariableMonitoring describes a monitoring setting for a variable.
type VariableMonitoring struct {
	ID          int         `json:"id" validate:"gte=0"`                  // Identifies the monitor.
	Transaction bool        `json:"transaction"`                          // Monitor only active when a transaction is ongoing on a component relevant to this transaction.
	Value       float64     `json:"value"`                                // Value for threshold or delta monitoring. For Periodic or PeriodicClockAligned this is the interval in seconds.
	Type        MonitorType `json:"type" validate:"required,monitorType"` // The type of this monitor, e.g. a threshold, delta or periodic monitor.
	Severity    int         `json:"severity" validate:"min=0,max=9"`      // The severity that will be assigned to an event that is triggered by this monitor. The severity range is 0-9, with 0 as the highest and 9 as the lowest severity level.
}

// NewVariableMonitoring is a utility function for creating a VariableMonitoring struct.
func NewVariableMonitoring(id int, transaction bool, value float64, t MonitorType, severity int) VariableMonitoring {
	return VariableMonitoring{ID: id, Transaction: transaction, Value: value, Type: t, Severity: severity}
}

// MonitoringData holds parameters of SetVariableMonitoring request.
type MonitoringData struct {
	Component          types.Component      `json:"component" validate:"required"`
	Variable           types.Variable       `json:"variable" validate:"required"`
	VariableMonitoring []VariableMonitoring `json:"variableMonitoring" validate:"required,min=1,dive"`
}

// The field definition of the NotifyMonitoringReport request payload sent by a Charging Station to the CSMS.
type NotifyMonitoringReportRequest struct {
	RequestID   int              `json:"requestId" validate:"gte=0"`                  // The id of the GetMonitoringRequest that requested this report.
	Tbc         bool             `json:"tbc,omitempty" validate:"omitempty"`          // “to be continued” indicator. Indicates whether another part of the monitoringData follows in an upcoming notifyMonitoringReportRequest message. Default value when omitted is false.
	SeqNo       int              `json:"seqNo" validate:"gte=0"`                      // Sequence number of this message. First message starts at 0.
	GeneratedAt *types.DateTime  `json:"generatedAt" validate:"required"`             // Timestamp of the moment this message was generated at the Charging Station.
	Monitor     []MonitoringData `json:"monitor,omitempty" validate:"omitempty,dive"` // List of MonitoringData containing monitoring settings.
}

// This field definition of the NotifyMonitoringReport response payload, sent by the CSMS to the Charging Station in response to a NotifyMonitoringReportRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type NotifyMonitoringReportResponse struct {
}

// The NotifyMonitoringReport feature is used by a Charging Station to send a report to the CSMS about configured
// monitoring settings per component and variable.
// Optionally, this list can be filtered on monitoringCriteria and componentVariables.
// After responding to a GetMonitoringReportRequest, a Charging Station will send one or more
// NotifyMonitoringReportRequest asynchronously to the CSMS, until all data of the monitoring report has been sent.
//
// The CSMS responds with a NotifyMonitoringReportResponse for every received received request.
type NotifyMonitoringReportFeature struct{}

func (f NotifyMonitoringReportFeature) GetFeatureName() string {
	return NotifyMonitoringReportFeatureName
}

func (f NotifyMonitoringReportFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(NotifyMonitoringReportRequest{})
}

func (f NotifyMonitoringReportFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(NotifyMonitoringReportResponse{})
}

func (r NotifyMonitoringReportRequest) GetFeatureName() string {
	return NotifyMonitoringReportFeatureName
}

func (c NotifyMonitoringReportResponse) GetFeatureName() string {
	return NotifyMonitoringReportFeatureName
}

// Creates a new NotifyMonitoringReportRequest, containing all required fields. Optional fields may be set afterwards.
func NewNotifyMonitoringReportRequest(requestID int, seqNo int, generatedAt *types.DateTime, monitorData []MonitoringData) *NotifyMonitoringReportRequest {
	return &NotifyMonitoringReportRequest{RequestID: requestID, SeqNo: seqNo, GeneratedAt: generatedAt, Monitor: monitorData}
}

// Creates a new NotifyMonitoringReportResponse, which doesn't contain any required or optional fields.
func NewNotifyMonitoringReportResponse() *NotifyMonitoringReportResponse {
	return &NotifyMonitoringReportResponse{}
}
