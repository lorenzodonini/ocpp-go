package provisioning

import (
	"reflect"

	"gopkg.in/go-playground/validator.v9"

	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
)

// -------------------- Reset (CSMS -> CS) --------------------

const ResetFeatureName = "Reset"

// ResetType indicates the type of reset requested.
type ResetType string

const (
	ResetTypeImmediate          ResetType = "Immediate"
	ResetTypeOnIdle             ResetType = "OnIdle"
	ResetTypeImmediateAndResume ResetType = "ImmediateAndResume"
)

func isValidResetType(fl validator.FieldLevel) bool {
	t := ResetType(fl.Field().String())
	switch t {
	case ResetTypeImmediate, ResetTypeOnIdle, ResetTypeImmediateAndResume:
		return true
	default:
		return false
	}
}

// ResetStatus indicates the response status to a reset request.
type ResetStatus string

const (
	ResetStatusAccepted  ResetStatus = "Accepted"
	ResetStatusRejected  ResetStatus = "Rejected"
	ResetStatusScheduled ResetStatus = "Scheduled"
)

func isValidResetStatus(fl validator.FieldLevel) bool {
	s := ResetStatus(fl.Field().String())
	switch s {
	case ResetStatusAccepted, ResetStatusRejected, ResetStatusScheduled:
		return true
	default:
		return false
	}
}

// ResetRequest is the payload for a Reset request from CSMS to CS.
type ResetRequest struct {
	Type       ResetType             `json:"type" validate:"required,resetType21"`
	EvseID     *int                  `json:"evseId,omitempty" validate:"omitempty,gte=0"`
	CustomData *types.CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// ResetResponse is the payload for a Reset response from CS to CSMS.
type ResetResponse struct {
	Status     ResetStatus           `json:"status" validate:"required,resetStatus21"`
	StatusInfo *types.StatusInfo     `json:"statusInfo,omitempty" validate:"omitempty"`
	CustomData *types.CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// ResetFeature represents the Reset feature.
type ResetFeature struct{}

func (f ResetFeature) GetFeatureName() string {
	return ResetFeatureName
}

func (f ResetFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(ResetRequest{})
}

func (f ResetFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(ResetResponse{})
}

func (r ResetRequest) GetFeatureName() string {
	return ResetFeatureName
}

func (c ResetResponse) GetFeatureName() string {
	return ResetFeatureName
}

// NewResetRequest creates a new ResetRequest.
func NewResetRequest(resetType ResetType) *ResetRequest {
	return &ResetRequest{Type: resetType}
}

// NewResetResponse creates a new ResetResponse.
func NewResetResponse(status ResetStatus) *ResetResponse {
	return &ResetResponse{Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("resetType21", isValidResetType)
	_ = types.Validate.RegisterValidation("resetStatus21", isValidResetStatus)
}
