package display

import (
	"reflect"
)

// -------------------- Notify Display Messages (CS -> CSMS) --------------------

const NotifyDisplayMessagesFeatureName = "NotifyDisplayMessages"

// The field definition of the NotifyDisplayMessages request payload sent by the CSMS to the Charging Station.
type NotifyDisplayMessagesRequest struct {
	RequestID   int           `json:"requestId" validate:"gte=0"`                      // The id of the GetDisplayMessagesRequest that requested this message.
	Tbc         bool          `json:"tbc,omitempty" validate:"omitempty"`              // "to be continued" indicator. Indicates whether another part of the report follows in an upcoming NotifyDisplayMessagesRequest message. Default value when omitted is false.
	MessageInfo []MessageInfo `json:"messageInfo,omitempty" validate:"omitempty,dive"` // The requested display message as configured in the Charging Station.
}

// This field definition of the NotifyDisplayMessages response payload, sent by the Charging Station to the CSMS in response to a NotifyDisplayMessagesRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type NotifyDisplayMessagesResponse struct {
}

// A CSO MAY request all the installed DisplayMessages configured via OCPP in a Charging Station. For this the CSO asks the CSMS to retrieve all messages (see NotifyDisplayMessages).
// If the Charging Station responded with a NotifyDisplayMessagesResponse Accepted, it will then send these messages asynchronously to the CSMS.
//
// The Charging Station sends one or more NotifyDisplayMessagesRequest message to the CSMS (depending on the amount of messages to be send).
// The CSMS responds to every notification with a NotifyDisplayMessagesResponse message.
type NotifyDisplayMessagesFeature struct{}

func (f NotifyDisplayMessagesFeature) GetFeatureName() string {
	return NotifyDisplayMessagesFeatureName
}

func (f NotifyDisplayMessagesFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(NotifyDisplayMessagesRequest{})
}

func (f NotifyDisplayMessagesFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(NotifyDisplayMessagesResponse{})
}

func (r NotifyDisplayMessagesRequest) GetFeatureName() string {
	return NotifyDisplayMessagesFeatureName
}

func (c NotifyDisplayMessagesResponse) GetFeatureName() string {
	return NotifyDisplayMessagesFeatureName
}

// Creates a new NotifyDisplayMessagesRequest, containing all required fields. Optional fields may be set afterwards.
func NewNotifyDisplayMessagesRequest(requestID int) *NotifyDisplayMessagesRequest {
	return &NotifyDisplayMessagesRequest{RequestID: requestID}
}

// Creates a new NotifyDisplayMessagesResponse, which doesn't contain any required or optional fields.
func NewNotifyDisplayMessagesResponse() *NotifyDisplayMessagesResponse {
	return &NotifyDisplayMessagesResponse{}
}
