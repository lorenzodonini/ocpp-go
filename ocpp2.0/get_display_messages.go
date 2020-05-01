package ocpp2

import (
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Get Display Messages (CSMS -> CS) --------------------

// Priority with which a message should be displayed on a Charging Station.
// Used within a GetDisplayMessagesRequest.
type MessagePriority string

// State of the Charging Station during which a message SHALL be displayed.
// Used within a GetDisplayMessagesRequest.
type MessageState string

// MessageStatus represents the status of the request, used in a GetDisplayMessagesConfirmation.
type MessageStatus string

const (
	MessagePriorityAlwaysFront MessagePriority = "AlwaysFront"
	MessagePriorityInFront     MessagePriority = "InFront"
	MessagePriorityNormalCycle MessagePriority = "NormalCycle"
	MessageStateCharging       MessageState    = "Charging"
	MessageStateFaulted        MessageState    = "Faulted"
	MessageStateIdle           MessageState    = "Idle"
	MessageStateUnavailable    MessageState    = "Unavailable"
	MessageStatusAccepted      MessageStatus   = "Accepted"
	MessageStatusUnknown       MessageStatus   = "Unknown"
)

func isValidMessagePriority(fl validator.FieldLevel) bool {
	priority := MessagePriority(fl.Field().String())
	switch priority {
	case MessagePriorityAlwaysFront, MessagePriorityInFront, MessagePriorityNormalCycle:
		return true
	default:
		return false
	}
}

func isValidMessageState(fl validator.FieldLevel) bool {
	priority := MessageState(fl.Field().String())
	switch priority {
	case MessageStateCharging, MessageStateFaulted, MessageStateIdle, MessageStateUnavailable:
		return true
	default:
		return false
	}
}

func isValidMessageStatus(fl validator.FieldLevel) bool {
	priority := MessageStatus(fl.Field().String())
	switch priority {
	case MessageStatusAccepted, MessageStatusUnknown:
		return true
	default:
		return false
	}
}

// The field definition of the GetDisplayMessages request payload sent by the CSMS to the Charging Station.
type GetDisplayMessagesRequest struct {
	RequestID int             `json:"requestId" validate:"gte=0"`
	Priority  MessagePriority `json:"priority,omitempty" validate:"omitempty,messagePriority"`
	State     MessageState    `json:"state,omitempty" validate:"omitempty,messageState"`
	ID        []int           `json:"id,omitempty" validate:"omitempty,dive,gte=0"`
}

// This field definition of the GetDisplayMessages confirmation payload, sent by the Charging Station to the CSMS in response to a GetDisplayMessagesRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type GetDisplayMessagesConfirmation struct {
	Status MessageStatus `json:"status" validate:"required,messageStatus"`
}

// A Charging Station can remove messages when they are out-dated, or transactions have ended. It can be very useful for a CSO to be able to view to current list of messages, so the CSO knows which messages are (still) configured.
//
// A CSO MAY request all the installed DisplayMessages configured via OCPP in a Charging Station. For this the CSO asks the CSMS to retrieve all messages.
// The CSMS sends a GetDisplayMessagesRequest message to the Charging Station.
// The Charging Station responds with a GetDisplayMessagesResponse Accepted, indicating it has configured messages and will send them.
//
// The Charging Station asynchronously sends one or more NotifyDisplayMessagesRequest messages to the
// CSMS (depending on the amount of messages to be sent).
type GetDisplayMessagesFeature struct{}

func (f GetDisplayMessagesFeature) GetFeatureName() string {
	return GetDisplayMessagesFeatureName
}

func (f GetDisplayMessagesFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(GetDisplayMessagesRequest{})
}

func (f GetDisplayMessagesFeature) GetConfirmationType() reflect.Type {
	return reflect.TypeOf(GetDisplayMessagesConfirmation{})
}

func (r GetDisplayMessagesRequest) GetFeatureName() string {
	return GetDisplayMessagesFeatureName
}

func (c GetDisplayMessagesConfirmation) GetFeatureName() string {
	return GetDisplayMessagesFeatureName
}

// Creates a new GetDisplayMessagesRequest, containing all required fields. Optional fields may be set afterwards.
func NewGetDisplayMessagesRequest(requestId int) *GetDisplayMessagesRequest {
	return &GetDisplayMessagesRequest{RequestID: requestId}
}

// Creates a new GetDisplayMessagesConfirmation, containing all required fields. There are no optional fields for this message.
func NewGetDisplayMessagesConfirmation(status MessageStatus) *GetDisplayMessagesConfirmation {
	return &GetDisplayMessagesConfirmation{Status: status}
}

func init() {
	_ = Validate.RegisterValidation("messagePriority", isValidMessagePriority)
	_ = Validate.RegisterValidation("messageState", isValidMessageState)
	_ = Validate.RegisterValidation("messageStatus", isValidMessageStatus)
}
