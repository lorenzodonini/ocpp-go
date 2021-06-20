package remotecontrol

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"gopkg.in/go-playground/validator.v9"
)

// -------------------- Trigger Message (CSMS -> CS) --------------------

const TriggerMessageFeatureName = "TriggerMessage"

// Type of request to be triggered by trigger messages.
type MessageTrigger string

// Status in TriggerMessageResponse.
type TriggerMessageStatus string

const (
	MessageTriggerBootNotification                  MessageTrigger = "BootNotification"
	MessageTriggerLogStatusNotification             MessageTrigger = "LogStatusNotification"
	MessageTriggerFirmwareStatusNotification        MessageTrigger = "FirmwareStatusNotification"
	MessageTriggerHeartbeat                         MessageTrigger = "Heartbeat"
	MessageTriggerMeterValues                       MessageTrigger = "MeterValues"
	MessageTriggerSignChargingStationCertificate    MessageTrigger = "SignChargingStationCertificate"
	MessageTriggerSignV2GCertificate                MessageTrigger = "SignV2GCertificate"
	MessageTriggerStatusNotification                MessageTrigger = "StatusNotification"
	MessageTriggerTransactionEvent                  MessageTrigger = "TransactionEvent"
	MessageTriggerSignCombinedCertificate           MessageTrigger = "SignCombinedCertificate"
	MessageTriggerPublishFirmwareStatusNotification MessageTrigger = "PublishFirmwareStatusNotification"

	TriggerMessageStatusAccepted       TriggerMessageStatus = "Accepted"
	TriggerMessageStatusRejected       TriggerMessageStatus = "Rejected"
	TriggerMessageStatusNotImplemented TriggerMessageStatus = "NotImplemented"
)

func isValidMessageTrigger(fl validator.FieldLevel) bool {
	status := MessageTrigger(fl.Field().String())
	switch status {
	case MessageTriggerBootNotification, MessageTriggerLogStatusNotification, MessageTriggerFirmwareStatusNotification,
		MessageTriggerHeartbeat, MessageTriggerMeterValues, MessageTriggerSignChargingStationCertificate,
		MessageTriggerSignV2GCertificate, MessageTriggerStatusNotification, MessageTriggerTransactionEvent,
		MessageTriggerSignCombinedCertificate, MessageTriggerPublishFirmwareStatusNotification:
		return true
	default:
		return false
	}
}

func isValidTriggerMessageStatus(fl validator.FieldLevel) bool {
	status := TriggerMessageStatus(fl.Field().String())
	switch status {
	case TriggerMessageStatusAccepted, TriggerMessageStatusRejected, TriggerMessageStatusNotImplemented:
		return true
	default:
		return false
	}
}

// The field definition of the TriggerMessage request payload sent by the CSMS to the Charging Station.
type TriggerMessageRequest struct {
	RequestedMessage MessageTrigger `json:"requestedMessage" validate:"required,messageTrigger"`
	Evse             *types.EVSE    `json:"evse,omitempty" validate:"omitempty"`
}

// This field definition of the TriggerMessage response payload, sent by the Charging Station to the CSMS in response to a TriggerMessageRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type TriggerMessageResponse struct {
	Status     TriggerMessageStatus `json:"status" validate:"required,triggerMessageStatus"`
	StatusInfo *types.StatusInfo    `json:"statusInfo,omitempty"`
}

// The CSMS may request a Charging Station to send a Charging Station-initiated message.
// This is achieved sending a TriggerMessageRequest to a charging station, indicating which message should be received.
// The charging station responds with a TriggerMessageResponse, indicating whether it will send a message or not.
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

// Creates a new TriggerMessageRequest, containing all required fields. Optional fields may be set afterwards.
func NewTriggerMessageRequest(requestedMessage MessageTrigger) *TriggerMessageRequest {
	return &TriggerMessageRequest{RequestedMessage: requestedMessage}
}

// Creates a new TriggerMessageResponse, containing all required fields. Optional fields may be set afterwards.
func NewTriggerMessageResponse(status TriggerMessageStatus) *TriggerMessageResponse {
	return &TriggerMessageResponse{Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("messageTrigger", isValidMessageTrigger)
	_ = types.Validate.RegisterValidation("triggerMessageStatus", isValidTriggerMessageStatus)
}
