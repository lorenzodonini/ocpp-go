package diagnostics

import (
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Notify Event (CS -> CSMS) --------------------

const NotifyEventFeatureName = "NotifyEvent"

// EventTrigger defines the type of monitor that triggered an event.
type EventTrigger string

const (
	EventTriggerAlerting EventTrigger = "Alerting" // Monitored variable has passed an Alert or Critical threshold.
	EventTriggerDelta    EventTrigger = "Delta"    // Delta Monitored Variable value has changed by more than specified amount.
	EventTriggerPeriodic EventTrigger = "Periodic" // Periodic Monitored Variable has been sampled for reporting at the specified interval.
)

func isValidEventTrigger(fl validator.FieldLevel) bool {
	status := EventTrigger(fl.Field().String())
	switch status {
	case EventTriggerAlerting, EventTriggerDelta, EventTriggerPeriodic:
		return true
	default:
		return false
	}
}

// EventNotification specifies the event notification type of the message.
type EventNotification string

const (
	EventHardWiredNotification EventNotification = "HardWiredNotification" // The software implemented by the manufacturer triggered a hardwired notification.
	EventHardWiredMonitor      EventNotification = "HardWiredMonitor"      // Triggered by a monitor, which is hardwired by the manufacturer.
	EventPreconfiguredMonitor  EventNotification = "PreconfiguredMonitor"  // Triggered by a monitor, which is preconfigured by the manufacturer.
	EventCustomMonitor         EventNotification = "CustomMonitor"         // Triggered by a monitor, which is set with the setvariablemonitoringrequest message by the Charging Station Operator.
)

func isValidEventNotification(fl validator.FieldLevel) bool {
	status := EventNotification(fl.Field().String())
	switch status {
	case EventHardWiredMonitor, EventHardWiredNotification, EventPreconfiguredMonitor, EventCustomMonitor:
		return true
	default:
		return false
	}
}

// An EventData element contains only the Component, Variable and VariableMonitoring data that caused an event.
type EventData struct {
	EventID               int               `json:"eventId" validate:"gte=0"`
	Timestamp             *types.DateTime   `json:"timestamp" validate:"required"`
	Trigger               EventTrigger      `json:"trigger" validate:"required,eventTrigger"`
	Cause                 *int              `json:"cause,omitempty" validate:"omitempty"`
	ActualValue           string            `json:"actualValue" validate:"required,max=2500"`
	TechCode              string            `json:"techCode,omitempty" validate:"omitempty,max=50"`
	TechInfo              string            `json:"techInfo,omitempty" validate:"omitempty,max=500"`
	Cleared               bool              `json:"cleared,omitempty"`
	TransactionID         string            `json:"transactionId,omitempty" validate:"omitempty,max=36"`
	VariableMonitoringID  *int              `json:"variableMonitoringId,omitempty" validate:"omitempty"`
	EventNotificationType EventNotification `json:"eventNotificationType" validate:"required,eventNotification"`
	Component             types.Component   `json:"component" validate:"required"`
	Variable              types.Variable    `json:"variable" validate:"required"`
}

// The field definition of the NotifyEvent request payload sent by a Charging Station to the CSMS.
type NotifyEventRequest struct {
	GeneratedAt *types.DateTime `json:"generatedAt" validate:"required"`          // Timestamp of the moment this message was generated at the Charging Station.
	SeqNo       int             `json:"seqNo" validate:"gte=0"`                   // Sequence number of this message. First message starts at 0.
	Tbc         bool            `json:"tbc,omitempty" validate:"omitempty"`       // “to be continued” indicator. Indicates whether another part of the monitoringData follows in an upcoming notifyMonitoringReportRequest message. Default value when omitted is false.
	EventData   []EventData     `json:"eventData" validate:"required,min=1,dive"` // The list of EventData will usually contain one eventData element, but the Charging Station may decide to group multiple events in one notification.
}

// This field definition of the NotifyEvent response payload, sent by the CSMS to the Charging Station in response to a NotifyEventRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type NotifyEventResponse struct {
}

// The NotifyEvent feature gives Charging Stations the ability to notify the CSMS (periodically) about monitoring events.
// If a threshold or a delta value has exceeded, the Charging Station sends a NotifyEventRequest to the CSMS.
// A request reports every Component/Variable for which a VariableMonitoring setting was triggered.
// Only the VariableMonitoring settings that are responsible for triggering an event are included.
// The monitoring setting(s) might have been configured explicitly via a SetVariableMonitoring message or
// it might be "hard-wired" in the Charging Station’s firmware.
//
// The CSMS responds to the request with a NotifyEventResponse.
type NotifyEventFeature struct{}

func (f NotifyEventFeature) GetFeatureName() string {
	return NotifyEventFeatureName
}

func (f NotifyEventFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(NotifyEventRequest{})
}

func (f NotifyEventFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(NotifyEventResponse{})
}

func (r NotifyEventRequest) GetFeatureName() string {
	return NotifyEventFeatureName
}

func (c NotifyEventResponse) GetFeatureName() string {
	return NotifyEventFeatureName
}

// Creates a new NotifyEventRequest, containing all required fields. Optional fields may be set afterwards.
func NewNotifyEventRequest(generatedAt *types.DateTime, seqNo int, eventData []EventData) *NotifyEventRequest {
	return &NotifyEventRequest{GeneratedAt: generatedAt, SeqNo: seqNo, EventData: eventData}
}

// Creates a new NotifyEventResponse, which doesn't contain any required or optional fields.
func NewNotifyEventResponse() *NotifyEventResponse {
	return &NotifyEventResponse{}
}

func init() {
	_ = types.Validate.RegisterValidation("eventTrigger", isValidEventTrigger)
	_ = types.Validate.RegisterValidation("eventNotification", isValidEventNotification)
}
