package ocpp16

import (
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Reset (CS -> CP) --------------------

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
	Type ResetType `json:"type" validate:"required,resetType"`
}

// This field definition of the Reset confirmation payload, sent by the Charge Point to the Central System in response to a ResetRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type ResetConfirmation struct {
	Status ResetStatus `json:"status" validate:"required,resetStatus"`
}

// The Central System SHALL send a ResetRequest for requesting a Charge Point to reset itself.
// The Central System can request a hard or a soft reset. Upon receipt of a ResetRequest, the Charge Point SHALL respond with a ResetConfirmation message.
// The response SHALL include whether the Charge Point will attempt to reset itself.
// After receipt of a ResetRequest, The Charge Point SHALL send a StopTransactionRequest for any ongoing transaction before performing the reset.
type ResetFeature struct{}

func (f ResetFeature) GetFeatureName() string {
	return ResetFeatureName
}

func (f ResetFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(ResetRequest{})
}

func (f ResetFeature) GetConfirmationType() reflect.Type {
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
	_ = Validate.RegisterValidation("resetType", isValidResetType)
	_ = Validate.RegisterValidation("resetStatus", isValidResetStatus)
}
