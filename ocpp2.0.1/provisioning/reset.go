package provisioning

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"gopkg.in/go-playground/validator.v9"
)

// -------------------- Reset (CSMS -> CS) --------------------

const ResetFeatureName = "Reset"

// ResetType indicates the type of reset that the charging station or EVSE should perform,
// as requested by the CSMS in a ResetRequest.
type ResetType string

const (
	ResetTypeImmediate ResetType = "Immediate"
	ResetTypeOnIdle    ResetType = "OnIdle"
)

func isValidResetType(fl validator.FieldLevel) bool {
	status := ResetType(fl.Field().String())
	switch status {
	case ResetTypeImmediate, ResetTypeOnIdle:
		return true
	default:
		return false
	}
}

// Result of a ResetRequest.
// This indicates whether the Charging Station is able to perform the reset.
type ResetStatus string

const (
	ResetStatusAccepted  ResetStatus = "Accepted"
	ResetStatusRejected  ResetStatus = "Rejected"
	ResetStatusScheduled ResetStatus = "Scheduled"
)

func isValidResetStatus(fl validator.FieldLevel) bool {
	status := ResetStatus(fl.Field().String())
	switch status {
	case ResetStatusAccepted, ResetStatusRejected, ResetStatusScheduled:
		return true
	default:
		return false
	}
}

// The field definition of the Reset request payload sent by the CSMS to the Charging Station.
type ResetRequest struct {
	Type   ResetType `json:"type" validate:"resetType201"`
	EvseID *int      `json:"evseId,omitempty" validate:"omitempty,gte=0"`
}

// This field definition of the Reset response payload, sent by the Charging Station to the CSMS in response to a ResetRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type ResetResponse struct {
	Status     ResetStatus       `json:"status" validate:"required,resetStatus201"`
	StatusInfo *types.StatusInfo `json:"statusInfo" validate:"omitempty"`
}

// The CSO may trigger the CSMS to request a Charging Station to reset itself or an EVSE.
// This can be used when a Charging Station is not functioning correctly, or when the configuration
// (e.g. network, security profiles, etc.) on the Charging Station changed, that needs an explicit reset.
//
// The CSMS sends a ResetRequest to the Charging Station.
// The Charging Station replies with a ResetResponse, then proceeds to reset itself,
// either immediately or whenever possible.
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

// Creates a new ResetRequest, containing all required fields. Optional fields may be set afterwards.
func NewResetRequest(t ResetType) *ResetRequest {
	return &ResetRequest{Type: t}
}

// Creates a new ResetResponse, containing all required fields. Optional fields may be set afterwards.
func NewResetResponse(status ResetStatus) *ResetResponse {
	return &ResetResponse{Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("resetType201", isValidResetType)
	_ = types.Validate.RegisterValidation("resetStatus201", isValidResetStatus)
}
