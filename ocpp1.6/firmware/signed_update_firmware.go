package firmware

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

const (
	SignedUpdateFirmwareFeatureName = "SignedUpdateFirmware"
)

type SignedUpdateFirmwareConfirmationStatus string

const (
	SignedUpdateFirmwareConfirmationStatusAccepted           SignedUpdateFirmwareConfirmationStatus = "Accepted"
	SignedUpdateFirmwareConfirmationStatusRejected           SignedUpdateFirmwareConfirmationStatus = "Rejected"
	SignedUpdateFirmwareConfirmationStatusAcceptedCanceled   SignedUpdateFirmwareConfirmationStatus = "AcceptedCanceled"
	SignedUpdateFirmwareConfirmationStatusInvalidCertificate SignedUpdateFirmwareConfirmationStatus = "InvalidCertificate"
	SignedUpdateFirmwareConfirmationStatusRevokedCertificate SignedUpdateFirmwareConfirmationStatus = "RevokedCertificate"
)

// The field definition of the SignedUpdateFirmware request payload sent by the Central System to the Charge Point.
type SignedUpdateFirmwareRequest struct {
	Retries       *int           `json:"retries,omitempty" validate:"omitempty,gte=0"`
	RetryInterval *int           `json:"retryInterval,omitempty"`
	RequestId     int            `json:"requestId" validate:"required"`
	Firmware      SignedFirmware `json:"firmware" validate:"required"`
}

type SignedFirmware struct {
	Location           string          `json:"location" validate:"required,uri"`
	RetrieveDateTime   *types.DateTime `json:"retrieveDateTime" validate:"required"`
	InstallDateTime    *types.DateTime `json:"installDateTime,omitempty"`
	Signature          string          `json:"signature" validate:"required"`
	SigningCertificate string          `json:"signingCertificate" validate:"required"`
}

type SignedUpdateFirmwareFeature struct{}

func (r SignedUpdateFirmwareRequest) GetFeatureName() string {
	return SignedUpdateFirmwareFeatureName
}

func (f SignedUpdateFirmwareFeature) GetFeatureName() string {
	return SignedUpdateFirmwareFeatureName
}

func (f SignedUpdateFirmwareFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(SignedUpdateFirmwareRequest{})
}

func (f SignedUpdateFirmwareFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(SignedUpdateFirmwareConfirmation{})
}

// This field definition of the SignedUpdateFirmware confirmation payload, sent by the Charge Point to the Central System in response to a SignedUpdateFirmwareRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type SignedUpdateFirmwareConfirmation struct {
	Status SignedUpdateFirmwareConfirmationStatus `json:"status" validate:"required,signedUpdateFirmwareConfirmationStatus"`
}

func (c SignedUpdateFirmwareConfirmation) GetFeatureName() string {
	return SignedUpdateFirmwareFeatureName
}

// Creates a new SignedUpdateFirmwareRequest, which doesn't contain any required or optional fields.
func NewSignedUpdateFirmwareRequest(requestId int, location, signature, signingCertificate string, retrieveDate *types.DateTime) *SignedUpdateFirmwareRequest {
	return &SignedUpdateFirmwareRequest{
		RequestId: requestId,
		Firmware: SignedFirmware{
			Location:           location,
			Signature:          signature,
			SigningCertificate: signingCertificate,
			RetrieveDateTime:   retrieveDate,
		},
	}
}

func NewSignedUpdateFirmwareConfirmation(status SignedUpdateFirmwareConfirmationStatus) *SignedUpdateFirmwareConfirmation {
	return &SignedUpdateFirmwareConfirmation{
		Status: status,
	}
}

func isValidSignedFirmwareConfirmationStatus(fl validator.FieldLevel) bool {
	status := SignedUpdateFirmwareConfirmationStatus(fl.Field().String())
	switch status {
	case SignedUpdateFirmwareConfirmationStatusAccepted, SignedUpdateFirmwareConfirmationStatusRejected, SignedUpdateFirmwareConfirmationStatusAcceptedCanceled,
		SignedUpdateFirmwareConfirmationStatusInvalidCertificate, SignedUpdateFirmwareConfirmationStatusRevokedCertificate:
		return true
	default:
		return false
	}
}

func init() {
	_ = types.Validate.RegisterValidation("signedUpdateFirmwareConfirmationStatus", isValidSignedFirmwareConfirmationStatus)
}
