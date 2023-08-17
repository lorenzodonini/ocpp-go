package core

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"gopkg.in/go-playground/validator.v9"
)

// -------------------- Reset (CS -> CP) --------------------

const ResetFeatureName = "Reset"

// Type of reset requested by ResetRequest.
type ResetType string

// Result of ResetRequest.
type ResetStatus string

const (
	ResetTypeHard       ResetType   = "Hard"
	ResetTypeSoft       ResetType   = "Soft"
	ResetStatusAccepted ResetStatus = "Accepted"
	ResetStatusRejected ResetStatus = "Rejected"
)

func isValidResetType(fl validator.FieldLevel) bool {
	status := ResetType(fl.Field().String())
	switch status {
	case ResetTypeHard, ResetTypeSoft:
		return true
	default:
		return false
	}
}

func isValidResetStatus(fl validator.FieldLevel) bool {
	status := ResetStatus(fl.Field().String())
	switch status {
	case ResetStatusAccepted, ResetStatusRejected:
		return true
	default:
		return false
	}
}

// The field definition of the Reset request payload sent by the Central System to the Charge Point.
type ResetRequest struct {
	Type ResetType `json:"type" validate:"required,resetType16"`
}

// This field definition of the Reset confirmation payload, sent by the Charge Point to the Central System in response to a ResetRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type ResetConfirmation struct {
	Status ResetStatus `json:"status" validate:"required,resetStatus16"`
}

// The Central System SHALL send a ResetRequest for requesting a Charge Point to reset itself.
// The Central System can request a hard or a soft reset. Upon receipt of a ResetRequest, the Charge Point SHALL respond with a ResetConfirmation message.
// The response SHALL include whether the Charge Point will attempt to reset itself.
// After receipt of a ResetRequest, The Charge Point SHALL send a StopTransactionRequest for any ongoing transaction before performing the reset.
// If the Charge Point fails to receive a StopTransactionConfirmation form the Central System, it shall queue the StopTransactionRequest.
// At receipt of a soft reset, the Charge Point SHALL stop ongoing transactions gracefully and send StopTransactionRequest for every ongoing transaction.
// It should then restart the application software (if possible, otherwise restart the processor/controller).
// At receipt of a hard reset the Charge Point SHALL restart (all) the hardware, it is not required to gracefully stop ongoing transaction.
// If possible the Charge Point sends a StopTransactionRequest for previously ongoing transactions after having restarted and having been accepted by the Central System via a BootNotificationConfirmation.
// This is a last resort solution for a not correctly functioning Charge Points, by sending a "hard" reset, (queued) information might get lost.
type ResetFeature struct{}

func (f ResetFeature) GetFeatureName() string {
	return ResetFeatureName
}

func (f ResetFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(ResetRequest{})
}

func (f ResetFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(ResetConfirmation{})
}

func (r ResetRequest) GetFeatureName() string {
	return ResetFeatureName
}

func (c ResetConfirmation) GetFeatureName() string {
	return ResetFeatureName
}

// Creates a new ResetRequest, containing all required fields. There are no optional fields for this message.
func NewResetRequest(resetType ResetType) *ResetRequest {
	return &ResetRequest{Type: resetType}
}

// Creates a new ResetConfirmation, containing all required fields. There are no optional fields for this message.
func NewResetConfirmation(status ResetStatus) *ResetConfirmation {
	return &ResetConfirmation{Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("resetType16", isValidResetType)
	_ = types.Validate.RegisterValidation("resetStatus16", isValidResetStatus)
}
