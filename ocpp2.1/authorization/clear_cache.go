package authorization

import (
	"reflect"

	"gopkg.in/go-playground/validator.v9"

	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
)

// -------------------- Clear Cache (CSMS -> CS) --------------------

const ClearCacheFeatureName = "ClearCache"

// ClearCacheStatus represents the status of a clear cache request.
type ClearCacheStatus string

const (
	ClearCacheStatusAccepted ClearCacheStatus = "Accepted"
	ClearCacheStatusRejected ClearCacheStatus = "Rejected"
)

func isValidClearCacheStatus(fl validator.FieldLevel) bool {
	s := ClearCacheStatus(fl.Field().String())
	switch s {
	case ClearCacheStatusAccepted, ClearCacheStatusRejected:
		return true
	default:
		return false
	}
}

// ClearCacheRequest is the payload for a ClearCache request from CSMS to CS.
type ClearCacheRequest struct {
	CustomData *types.CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// ClearCacheResponse is the payload for a ClearCache response from CS to CSMS.
type ClearCacheResponse struct {
	Status     ClearCacheStatus      `json:"status" validate:"required,clearCacheStatus21"`
	StatusInfo *types.StatusInfo     `json:"statusInfo,omitempty" validate:"omitempty"`
	CustomData *types.CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// ClearCacheFeature represents the ClearCache feature.
type ClearCacheFeature struct{}

func (f ClearCacheFeature) GetFeatureName() string {
	return ClearCacheFeatureName
}

func (f ClearCacheFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(ClearCacheRequest{})
}

func (f ClearCacheFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(ClearCacheResponse{})
}

func (r ClearCacheRequest) GetFeatureName() string {
	return ClearCacheFeatureName
}

func (c ClearCacheResponse) GetFeatureName() string {
	return ClearCacheFeatureName
}

// NewClearCacheRequest creates a new ClearCacheRequest.
func NewClearCacheRequest() *ClearCacheRequest {
	return &ClearCacheRequest{}
}

// NewClearCacheResponse creates a new ClearCacheResponse.
func NewClearCacheResponse(status ClearCacheStatus) *ClearCacheResponse {
	return &ClearCacheResponse{Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("clearCacheStatus21", isValidClearCacheStatus)
}
