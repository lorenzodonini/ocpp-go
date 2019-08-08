package ocpp16

import (
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Reset (CS -> CP) --------------------
type ResetType string
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

type ResetRequest struct {
	Type ResetType `json:"type" validate:"required,resetType"`
}

type ResetConfirmation struct {
	Status ResetStatus `json:"status" validate:"required,resetStatus"`
}

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

func NewResetRequest(resetType ResetType) *ResetRequest {
	return &ResetRequest{Type: resetType}
}

func NewResetConfirmation(status ResetStatus) *ResetConfirmation {
	return &ResetConfirmation{Status: status}
}

func init() {
	_ = Validate.RegisterValidation("resetType", isValidResetType)
	_ = Validate.RegisterValidation("resetStatus", isValidResetStatus)
}
