package remotetrigger

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/firmware"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Trigger Message (CS -> CP) --------------------

const TriggerMessageFeatureName = "TriggerMessage"

// Status reported in TriggerMessageConfirmation.
type TriggerMessageStatus string

// Type of request to be triggered in a TriggerMessageRequest
type MessageTrigger string

const (
	TriggerMessageStatusAccepted       TriggerMessageStatus = "Accepted"
	TriggerMessageStatusRejected       TriggerMessageStatus = "Rejected"
	TriggerMessageStatusNotImplemented TriggerMessageStatus = "NotImplemented"
)

func isValidTriggerMessageStatus(fl validator.FieldLevel) bool {
	status := TriggerMessageStatus(fl.Field().String())
	switch status {
	case TriggerMessageStatusAccepted, TriggerMessageStatusRejected, TriggerMessageStatusNotImplemented:
		return true
	default:
		return false
	}
}

func isValidMessageTrigger(fl validator.FieldLevel) bool {
	trigger := MessageTrigger(fl.Field().String())
	switch trigger {
	case core.BootNotificationFeatureName, firmware.DiagnosticsStatusNotificationFeatureName, firmware.FirmwareStatusNotificationFeatureName, core.HeartbeatFeatureName, core.MeterValuesFeatureName, core.StatusNotificationFeatureName:
		return true
	default:
		return false
	}
}

// The field definition of the TriggerMessage request payload sent by the Central System to the Charge Point.
type TriggerMessageRequest struct {
	RequestedMessage MessageTrigger `json:"requestedMessage" validate:"required,messageTrigger"`
	ConnectorId      *int           `json:"connectorId,omitempty" validate:"omitempty,gt=0"`
}

// This field definition of the TriggerMessage confirmation payload, sent by the Charge Point to the Central System in response to a TriggerMessageRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type TriggerMessageConfirmation struct {
	Status TriggerMessageStatus `json:"status" validate:"required,triggerMessageStatus"`
}

// During normal operation, the Charge Point informs the Central System of its state and any relevant occurrences.
// If there is nothing to report the Charge Point will send at least a heartBeat at a predefined interval.
// Under normal circumstances this is just fine, but what if the Central System has (whatever) reason to doubt the last known state?
// What can a Central System do if a firmware update is in progress and the last status notification it received about it was much longer ago than could reasonably be expected?
// The same can be asked for the progress of a diagnostics request. The problem in these situations is not that the information needed isn’t covered by existing messages, the problem is strictly a timing issue.
// The Charge Point has the information, but has no way of knowing that the Central System would like an update.
// The TriggerMessageRequest makes it possible for the Central System, to request the Charge Point, to send Charge Point-initiated messages.
// In the request the Central System indicates which message it wishes to receive.
// For every such requested message the Central System MAY optionally indicate to which connector this request applies.
// The requested message is leading: if the specified connectorId is not relevant to the message, it should be ignored. In such cases the requested message should still be sent.
// Inversely, if the connectorId is relevant but absent, this should be interpreted as “for all allowed connectorId values”.
// For example, a request for a statusNotification for connectorId 0 is a request for the status of the Charge Point.
// A request for a statusNotification without connectorId is a request for multiple statusNotifications: the notification for the Charge Point itself and a notification for each of its connectors.
// The Charge Point SHALL first send the TriggerMessage response, before sending the requested message.
// In the TriggerMessageConfirmation the Charge Point SHALL indicate whether it will send it or not, by returning ACCEPTED or REJECTED.
// It is up to the Charge Point if it accepts or rejects the request to send.
// If the requested message is unknown or not implemented the Charge Point SHALL return NOT_IMPLEMENTED.
// Messages that the Charge Point marks as accepted SHOULD be sent. The situation could occur that, between accepting the request and actually sending the requested message, that same message gets sent because of normal operations. In such cases the message just sent MAY be considered as complying with the request.
// The TriggerMessage mechanism is not intended to retrieve historic data. The messages it triggers should only give current information.
// A MeterValuesRequest message triggered in this way for instance SHALL return the most recent measurements for all measurands configured in configuration key MeterValuesSampledData.
// StartTransaction and StopTransaction have been left out of this mechanism because they are not state related, but by their nature describe a transition.
type TriggerMessageFeature struct{}

func (f TriggerMessageFeature) GetFeatureName() string {
	return TriggerMessageFeatureName
}

func (f TriggerMessageFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(TriggerMessageRequest{})
}

func (f TriggerMessageFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(TriggerMessageConfirmation{})
}

func (r TriggerMessageRequest) GetFeatureName() string {
	return TriggerMessageFeatureName
}

func (c TriggerMessageConfirmation) GetFeatureName() string {
	return TriggerMessageFeatureName
}

// Creates a new TriggerMessageRequest, containing all required fields. Optional fields may be set afterwards.
func NewTriggerMessageRequest(requestedMessage MessageTrigger) *TriggerMessageRequest {
	return &TriggerMessageRequest{RequestedMessage: requestedMessage}
}

// Creates a new TriggerMessageConfirmation, containing all required fields. There are no optional fields for this message.
func NewTriggerMessageConfirmation(status TriggerMessageStatus) *TriggerMessageConfirmation {
	return &TriggerMessageConfirmation{Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("triggerMessageStatus", isValidTriggerMessageStatus)
	_ = types.Validate.RegisterValidation("messageTrigger", isValidMessageTrigger)
}
