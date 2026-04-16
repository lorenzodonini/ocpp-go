package remotecontrol

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
)

// -------------------- Trigger Message (CSMS -> CS) --------------------

const TriggerMessageFeatureName = "TriggerMessage"

// TriggerMessageRequest is the payload for a TriggerMessage request from CSMS to CS.
type TriggerMessageRequest struct {
	RequestedMessage types.MessageTrigger  `json:"requestedMessage" validate:"required,messageTrigger21"`
	Evse             *types.EVSE           `json:"evse,omitempty" validate:"omitempty"`
	CustomTrigger    string                `json:"customTrigger,omitempty" validate:"omitempty,max=50"`
	CustomData       *types.CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// TriggerMessageResponse is the payload for a TriggerMessage response from CS to CSMS.
type TriggerMessageResponse struct {
	Status     types.TriggerMessageStatus `json:"status" validate:"required,triggerMessageStatus21"`
	StatusInfo *types.StatusInfo          `json:"statusInfo,omitempty" validate:"omitempty"`
	CustomData *types.CustomDataType      `json:"customData,omitempty" validate:"omitempty"`
}

// TriggerMessageFeature represents the TriggerMessage feature.
type TriggerMessageFeature struct{}

func (f TriggerMessageFeature) GetFeatureName() string {
	return TriggerMessageFeatureName
}

func (f TriggerMessageFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(TriggerMessageRequest{})
}

func (f TriggerMessageFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(TriggerMessageResponse{})
}

func (r TriggerMessageRequest) GetFeatureName() string {
	return TriggerMessageFeatureName
}

func (c TriggerMessageResponse) GetFeatureName() string {
	return TriggerMessageFeatureName
}

// NewTriggerMessageRequest creates a new TriggerMessageRequest.
func NewTriggerMessageRequest(requestedMessage types.MessageTrigger) *TriggerMessageRequest {
	return &TriggerMessageRequest{RequestedMessage: requestedMessage}
}

// NewTriggerMessageResponse creates a new TriggerMessageResponse.
func NewTriggerMessageResponse(status types.TriggerMessageStatus) *TriggerMessageResponse {
	return &TriggerMessageResponse{Status: status}
}
