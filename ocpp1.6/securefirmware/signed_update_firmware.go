package securefirmware

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"gopkg.in/go-playground/validator.v9"
)

const SignedUpdateFirmwareFeatureName = "SignedUpdateFirmware"

type SignedUpdateFirmwareFeature struct{}

func (e SignedUpdateFirmwareFeature) GetFeatureName() string {
	return SignedUpdateFirmwareFeatureName
}

func (e SignedUpdateFirmwareFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(SignedUpdateFirmwareRequest{})
}

func (e SignedUpdateFirmwareFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(SignedUpdateFirmwareResponse{})
}

type SignedUpdateFirmwareType string

type SignedUpdateFirmwareStatus string

const (
	SignedUpdateFirmwareTypeBootNotification               SignedUpdateFirmwareType = "BootNotification"               // This contains the field definition of a diagnostics log file
	SignedUpdateFirmwareTypeLogStatusNotification          SignedUpdateFirmwareType = "LogStatusNotification"          // Sent by the CSMS to the Charging Station to request that the Charging Station uploads the security log
	SignedUpdateFirmwareTypeHeartbeat                      SignedUpdateFirmwareType = "Heartbeat"                      // Accepted this log upload. This does not mean the log file is uploaded is successfully, the Charging Station will now start the log file upload.
	SignedUpdateFirmwareTypeMeterValues                    SignedUpdateFirmwareType = "MeterValues"                    // Log update request rejected.
	SignedUpdateFirmwareTypeSignChargingStationCertificate SignedUpdateFirmwareType = "SignChargingStationCertificate" // Accepted this log upload, but in doing this has canceled an ongoing log file upload.
	SignedUpdateFirmwareTypeFirmwareStatusNotification     SignedUpdateFirmwareType = "FirmwareStatusNotification"     // Accepted this log upload, but in doing this has canceled an ongoing log file upload.
	SignedUpdateFirmwareTypeStatusNotification             SignedUpdateFirmwareType = "StatusNotification"             // Accepted this log upload, but in doing this has canceled an ongoing log file upload.

	SignedUpdateFirmwareStatusAccepted       SignedUpdateFirmwareStatus = "Accepted"
	SignedUpdateFirmwareStatusRejected       SignedUpdateFirmwareStatus = "Rejected"
	SignedUpdateFirmwareStatusNotImplemented SignedUpdateFirmwareStatus = "NotImplemented"
)

func isValidSignedUpdateFirmwareType(fl validator.FieldLevel) bool {
	status := SignedUpdateFirmwareType(fl.Field().String())
	switch status {
	case SignedUpdateFirmwareTypeBootNotification,
		SignedUpdateFirmwareTypeLogStatusNotification,
		SignedUpdateFirmwareTypeHeartbeat,
		SignedUpdateFirmwareTypeMeterValues,
		SignedUpdateFirmwareTypeSignChargingStationCertificate,
		SignedUpdateFirmwareTypeFirmwareStatusNotification,
		SignedUpdateFirmwareTypeStatusNotification:
		return true
	default:
		return false
	}
}

func isValidSignedUpdateFirmwareStatus(fl validator.FieldLevel) bool {
	status := SignedUpdateFirmwareStatus(fl.Field().String())
	switch status {
	case SignedUpdateFirmwareStatusAccepted,
		SignedUpdateFirmwareStatusRejected,
		SignedUpdateFirmwareStatusNotImplemented:
		return true
	default:
		return false
	}
}

// The field definition of the LogStatusNotification request payload sent by a Charging Station to the CSMS.
type SignedUpdateFirmwareRequest struct {
	RequestedMessage SignedUpdateFirmwareType `json:"requestedMessage" validate:"required,extendedTriggerMessageType"`
	ConnectorId      *int                     `json:"connectorId" validate:"gt=0,omitempty"`
}

// This field definition of the LogStatusNotification response payload, sent by the CSMS to the Charging Station in response to a SignedUpdateFirmwareRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type SignedUpdateFirmwareResponse struct {
	Status SignedUpdateFirmwareStatus `json:"status" validate:"required,extendedTriggerMessageStatus"`
}

func (r SignedUpdateFirmwareRequest) GetFeatureName() string {
	return SignedUpdateFirmwareFeatureName
}

func (c SignedUpdateFirmwareResponse) GetFeatureName() string {
	return SignedUpdateFirmwareFeatureName
}

// Creates a new SignedUpdateFirmwareRequest, containing all required fields. There are no optional fields for this message.
func NewSignedUpdateFirmwareRequest(requestedMessage SignedUpdateFirmwareType) *SignedUpdateFirmwareRequest {
	return &SignedUpdateFirmwareRequest{RequestedMessage: requestedMessage}
}

// Creates a new SignedUpdateFirmwareResponse, which doesn't contain any required or optional fields.
func NewSignedUpdateFirmwareResponse() *SignedUpdateFirmwareResponse {
	return &SignedUpdateFirmwareResponse{}
}

func init() {
	_ = types.Validate.RegisterValidation("extendedTriggerMessageType", isValidSignedUpdateFirmwareType)
	_ = types.Validate.RegisterValidation("extendedTriggerMessageStatus", isValidSignedUpdateFirmwareStatus)
}
