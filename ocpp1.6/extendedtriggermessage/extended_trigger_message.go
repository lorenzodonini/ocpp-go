package extendedtriggermessage

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"gopkg.in/go-playground/validator.v9"
)

const ExtendedTriggerMessageFeatureName = "ExtendedTriggerMessage"

type ExtendedTriggerMessageFeature struct{}

func (e ExtendedTriggerMessageFeature) GetFeatureName() string {
	return ExtendedTriggerMessageFeatureName
}

func (e ExtendedTriggerMessageFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(ExtendedTriggerMessageRequest{})
}

func (e ExtendedTriggerMessageFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(ExtendedTriggerMessageResponse{})
}

type ExtendedTriggerMessageType string

type ExtendedTriggerMessageStatus string

const (
	ExtendedTriggerMessageTypeBootNotification               ExtendedTriggerMessageType = "BootNotification"           // This contains the field definition of a diagnostics log file
	ExtendedTriggerMessageTypeLogStatusNotification          ExtendedTriggerMessageType = "LogStatusNotification"      // Sent by the CSMS to the Charging Station to request that the Charging Station uploads the security log
	ExtendedTriggerMessageTypeHeartbeat                      ExtendedTriggerMessageType = "Heartbeat"                  // Accepted this log upload. This does not mean the log file is uploaded is successfully, the Charging Station will now start the log file upload.
	ExtendedTriggerMessageTypeMeterValues                    ExtendedTriggerMessageType = "MeterValues"                // Log update request rejected.
	ExtendedTriggerMessageTypeSignChargingStationCertificate ExtendedTriggerMessageType = "SignChargePointCertificate" // Accepted this log upload, but in doing this has canceled an ongoing log file upload.
	ExtendedTriggerMessageTypeFirmwareStatusNotification     ExtendedTriggerMessageType = "FirmwareStatusNotification" // Accepted this log upload, but in doing this has canceled an ongoing log file upload.
	ExtendedTriggerMessageTypeStatusNotification             ExtendedTriggerMessageType = "StatusNotification"         // Accepted this log upload, but in doing this has canceled an ongoing log file upload.

	ExtendedTriggerMessageStatusAccepted       ExtendedTriggerMessageStatus = "Accepted"
	ExtendedTriggerMessageStatusRejected       ExtendedTriggerMessageStatus = "Rejected"
	ExtendedTriggerMessageStatusNotImplemented ExtendedTriggerMessageStatus = "NotImplemented"
)

func isValidExtendedTriggerMessageType(fl validator.FieldLevel) bool {
	status := ExtendedTriggerMessageType(fl.Field().String())
	switch status {
	case ExtendedTriggerMessageTypeBootNotification,
		ExtendedTriggerMessageTypeLogStatusNotification,
		ExtendedTriggerMessageTypeHeartbeat,
		ExtendedTriggerMessageTypeMeterValues,
		ExtendedTriggerMessageTypeSignChargingStationCertificate,
		ExtendedTriggerMessageTypeFirmwareStatusNotification,
		ExtendedTriggerMessageTypeStatusNotification:
		return true
	default:
		return false
	}
}

func isValidExtendedTriggerMessageStatus(fl validator.FieldLevel) bool {
	status := ExtendedTriggerMessageStatus(fl.Field().String())
	switch status {
	case ExtendedTriggerMessageStatusAccepted,
		ExtendedTriggerMessageStatusRejected,
		ExtendedTriggerMessageStatusNotImplemented:
		return true
	default:
		return false
	}
}

// The field definition of the LogStatusNotification request payload sent by a Charging Station to the CSMS.
type ExtendedTriggerMessageRequest struct {
	RequestedMessage ExtendedTriggerMessageType `json:"requestedMessage" validate:"required,extendedTriggerMessageType"`
	ConnectorId      *int                       `json:"connectorId" validate:"gt=0,omitempty"`
}

// This field definition of the LogStatusNotification response payload, sent by the CSMS to the Charging Station in response to a ExtendedTriggerMessageRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type ExtendedTriggerMessageResponse struct {
	Status ExtendedTriggerMessageStatus `json:"status" validate:"required,extendedTriggerMessageStatus"`
}

func (r ExtendedTriggerMessageRequest) GetFeatureName() string {
	return ExtendedTriggerMessageFeatureName
}

func (c ExtendedTriggerMessageResponse) GetFeatureName() string {
	return ExtendedTriggerMessageFeatureName
}

// Creates a new ExtendedTriggerMessageRequest, containing all required fields. There are no optional fields for this message.
func NewExtendedTriggerMessageRequest(requestedMessage ExtendedTriggerMessageType) *ExtendedTriggerMessageRequest {
	return &ExtendedTriggerMessageRequest{RequestedMessage: requestedMessage}
}

// Creates a new ExtendedTriggerMessageResponse, which doesn't contain any required or optional fields.
func NewExtendedTriggerMessageResponse(status ExtendedTriggerMessageStatus) *ExtendedTriggerMessageResponse {
	return &ExtendedTriggerMessageResponse{status}
}

func init() {
	_ = types.Validate.RegisterValidation("extendedTriggerMessageType", isValidExtendedTriggerMessageType)
	_ = types.Validate.RegisterValidation("extendedTriggerMessageStatus", isValidExtendedTriggerMessageStatus)
}
