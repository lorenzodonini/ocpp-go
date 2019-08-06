package ocpp16

import (
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Clear Cache (CS -> CP) --------------------
type ClearCacheStatus string

const (
	ClearCacheStatusAccepted ClearCacheStatus = "Accepted"
	ClearCacheStatusRejected ClearCacheStatus = "Rejected"
)

func isValidClearCacheStatus(fl validator.FieldLevel) bool {
	status := ClearCacheStatus(fl.Field().String())
	switch status {
	case ClearCacheStatusAccepted, ClearCacheStatusRejected:
		return true
	default:
		return false
	}
}

type ClearCacheRequest struct {
}

type ClearCacheConfirmation struct {
	Status             ClearCacheStatus `json:"status" validate:"required,cacheStatus"`
}

type ClearCacheFeature struct{}

func (f ClearCacheFeature) GetFeatureName() string {
	return ClearCacheFeatureName
}

func (f ClearCacheFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(ClearCacheRequest{})
}

func (f ClearCacheFeature) GetConfirmationType() reflect.Type {
	return reflect.TypeOf(ClearCacheConfirmation{})
}

func (r ClearCacheRequest) GetFeatureName() string {
	return ClearCacheFeatureName
}

func (c ClearCacheConfirmation) GetFeatureName() string {
	return ClearCacheFeatureName
}

func NewClearCacheRequest() *ClearCacheRequest {
	return &ClearCacheRequest{}
}

func NewClearCacheConfirmation(status ClearCacheStatus) *ClearCacheConfirmation {
	return &ClearCacheConfirmation{Status: status}
}

func init() {
	_ = Validate.RegisterValidation("cacheStatus", isValidClearCacheStatus)
}
