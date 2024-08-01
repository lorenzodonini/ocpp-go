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

// Indicates whether the Charging Station was able to accept the request.
type UpdateFirmwareStatus string

const (
	UpdateFirmwareStatusAccepted           UpdateFirmwareStatus = "Accepted"
	UpdateFirmwareStatusRejected           UpdateFirmwareStatus = "Rejected"
	UpdateFirmwareStatusAcceptedCanceled   UpdateFirmwareStatus = "AcceptedCanceled"
	UpdateFirmwareStatusInvalidCertificate UpdateFirmwareStatus = "InvalidCertificate"
	UpdateFirmwareStatusRevokedCertificate UpdateFirmwareStatus = "RevokedCertificate"
)

func isValidUpdateFirmwareStatus(fl validator.FieldLevel) bool {
	status := UpdateFirmwareStatus(fl.Field().String())
	switch status {
	case UpdateFirmwareStatusAccepted,
		UpdateFirmwareStatusRejected,
		UpdateFirmwareStatusAcceptedCanceled,
		UpdateFirmwareStatusInvalidCertificate,
		UpdateFirmwareStatusRevokedCertificate:
		return true
	default:
		return false
	}
}

// The field definition of the SignedUpdateFirmwareRequest request payload sent by a Charging Station to the CSMS.
type SignedUpdateFirmwareRequest struct {
	Retries       *int     `json:"retries,omitempty" validate:"omitempty,gte=0"`       // This specifies how many times Charging Station must try to download the firmware before giving up. If this field is not present, it is left to Charging Station to decide how many times it wants to retry.
	RetryInterval *int     `json:"retryInterval,omitempty" validate:"omitempty,gte=0"` // The interval in seconds after which a retry may be attempted. If this field is not present, it is left to Charging Station to decide how long to wait between attempts.
	RequestID     int      `json:"requestId" validate:"gte=0"`                         // The Id of the request.
	Firmware      Firmware `json:"firmware" validate:"required"`                       // Specifies the firmware to be updated on the Charging Station.
}

// Represents a copy of the firmware that can be loaded/updated on the Charging Station.
type Firmware struct {
	Location           string          `json:"location" validate:"required,max=512,uri"`         // URI defining the origin of the firmware.
	RetrieveDateTime   *types.DateTime `json:"retrieveDateTime" validate:"required"`             // Date and time at which the firmware shall be retrieved.
	InstallDateTime    *types.DateTime `json:"installDateTime,omitempty" validate:"omitempty"`   // Date and time at which the firmware shall be installed.
	SigningCertificate string          `json:"signingCertificate,omitempty" validate:"max=5500"` // Certificate with which the firmware was signed. PEM encoded X.509 certificate.
	Signature          string          `json:"signature,omitempty" validate:"max=800"`           // Base64 encoded firmware signature.
}

// This field definition of the LogStatusNotification response payload, sent by the CSMS to the Charging Station in response to a SignedUpdateFirmwareRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type SignedUpdateFirmwareResponse struct {
	Status UpdateFirmwareStatus `json:"status" validate:"required,signedUpdateFirmwareStatus"`
}

func (r SignedUpdateFirmwareRequest) GetFeatureName() string {
	return SignedUpdateFirmwareFeatureName
}

func (c SignedUpdateFirmwareResponse) GetFeatureName() string {
	return SignedUpdateFirmwareFeatureName
}

// Creates a new SignedUpdateFirmwareRequest, containing all required fields. There are no optional fields for this message.
func NewSignedUpdateFirmwareRequest(requestId int, firmware Firmware) *SignedUpdateFirmwareRequest {
	return &SignedUpdateFirmwareRequest{RequestID: requestId, Firmware: firmware}
}

// Creates a new SignedUpdateFirmwareResponse, which doesn't contain any required or optional fields.
func NewSignedUpdateFirmwareResponse(status UpdateFirmwareStatus) *SignedUpdateFirmwareResponse {
	return &SignedUpdateFirmwareResponse{Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("signedUpdateFirmwareStatus", isValidUpdateFirmwareStatus)
}
