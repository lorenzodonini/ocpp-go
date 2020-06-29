package display

import (
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
	"gopkg.in/go-playground/validator.v9"
)

// Priority with which a message should be displayed on a Charging Station.
// Used within a GetDisplayMessagesRequest.
type MessagePriority string

// State of the Charging Station during which a message SHALL be displayed.
// Used within a GetDisplayMessagesRequest.
type MessageState string

// MessageStatus represents the status of the request, used in a GetDisplayMessagesResponse.
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

// Contains message details, for a message to be displayed on a Charging Station.
type MessageInfo struct {
	ID            int                  `json:"id" validate:"gte=0"`                                 // Master resource identifier, unique within an exchange context. It is defined within the OCPP context as a positive Integer value (greater or equal to zero).
	Priority      MessagePriority      `json:"priority" validate:"required,messagePriority"`        // With what priority should this message be shown
	State         MessageState         `json:"state,omitempty" validate:"omitempty,messageState"`   // During what state should this message be shown. When omitted this message should be shown in any state of the Charging Station.
	StartDateTime *types.DateTime      `json:"startDateTime,omitempty" validate:"omitempty"`        // From what date-time should this message be shown. If omitted: directly.
	EndDateTime   *types.DateTime      `json:"endDateTime,omitempty" validate:"omitempty"`          // Until what date-time should this message be shown, after this date/time this message SHALL be removed.
	TransactionID string               `json:"transactionId,omitempty" validate:"omitempty,max=36"` // During which transaction shall this message be shown. Message SHALL be removed by the Charging Station after transaction has ended.
	Message       types.MessageContent `json:"message" validate:"required"`                         // Contains message details for the message to be displayed on a Charging Station.
	Display       *types.Component     `json:"display,omitempty" validate:"omitempty"`              // When a Charging Station has multiple Displays, this field can be used to define to which Display this message belongs.
}

func init() {
	_ = types.Validate.RegisterValidation("messagePriority", isValidMessagePriority)
	_ = types.Validate.RegisterValidation("messageState", isValidMessageState)
	_ = types.Validate.RegisterValidation("messageStatus", isValidMessageStatus)
}
