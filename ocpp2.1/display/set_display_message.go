package display

import (
	"reflect"

	"github.com/xBlaz3kx/ocpp-go/ocpp2.1/types"
	"gopkg.in/go-playground/validator.v9"
)

// -------------------- Clear Display (CSMS -> CS) --------------------

const SetDisplayMessageFeatureName = "SetDisplayMessage"

// Status returned in response to SetDisplayMessageRequest.
type DisplayMessageStatus string

const (
	DisplayMessageStatusAccepted                  DisplayMessageStatus = "Accepted"
	DisplayMessageStatusNotSupportedMessageFormat DisplayMessageStatus = "NotSupportedMessageFormat"
	DisplayMessageStatusRejected                  DisplayMessageStatus = "Rejected"
	DisplayMessageStatusNotSupportedPriority      DisplayMessageStatus = "NotSupportedPriority"
	DisplayMessageStatusNotSupportedState         DisplayMessageStatus = "NotSupportedState"
	DisplayMessageStatusUnknownTransaction        DisplayMessageStatus = "UnknownTransaction"
)

func isValidDisplayMessageStatus(fl validator.FieldLevel) bool {
	status := DisplayMessageStatus(fl.Field().String())
	switch status {
	case DisplayMessageStatusAccepted,
		DisplayMessageStatusNotSupportedMessageFormat,
		DisplayMessageStatusRejected,
		DisplayMessageStatusNotSupportedPriority,
		DisplayMessageStatusNotSupportedState,
		DisplayMessageStatusUnknownTransaction:
		return true
	default:
		return false
	}
}

// The field definition of the SetDisplayMessage request payload sent by the CSMS to the Charging Station.
type SetDisplayMessageRequest struct {
	Message MessageInfo `json:"message" validate:"required"`
}

// This field definition of the SetDisplayMessage response payload, sent by the Charging Station to the CSMS in response to a SetDisplayMessageRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type SetDisplayMessageResponse struct {
	Status     DisplayMessageStatus `json:"status" validate:"required,displayMessageStatus21"`
	StatusInfo *types.StatusInfo    `json:"statusInfo,omitempty" validate:"omitempty"`
}

// The CSMS may send a SetDisplayMessageRequest message to a Charging Station, instructing it to display a new message,
// which is not part of its firmware.
// The Charging Station accepts the request by replying with a SetDisplayMessageResponse.
//
// Depending on different parameters, the message may be displayed in different ways and/or at a configured time.
type SetDisplayMessageFeature struct{}

func (f SetDisplayMessageFeature) GetFeatureName() string {
	return SetDisplayMessageFeatureName
}

func (f SetDisplayMessageFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(SetDisplayMessageRequest{})
}

func (f SetDisplayMessageFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(SetDisplayMessageResponse{})
}

func (r SetDisplayMessageRequest) GetFeatureName() string {
	return SetDisplayMessageFeatureName
}

func (c SetDisplayMessageResponse) GetFeatureName() string {
	return SetDisplayMessageFeatureName
}

// Creates a new SetDisplayMessageRequest, containing all required fields. There are no optional fields for this message.
func NewSetDisplayMessageRequest(message MessageInfo) *SetDisplayMessageRequest {
	return &SetDisplayMessageRequest{Message: message}
}

// Creates a new SetDisplayMessageResponse, containing all required fields. Optional fields may be set afterwards.
func NewSetDisplayMessageResponse(status DisplayMessageStatus) *SetDisplayMessageResponse {
	return &SetDisplayMessageResponse{Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("displayMessageStatus21", isValidDisplayMessageStatus)
}
