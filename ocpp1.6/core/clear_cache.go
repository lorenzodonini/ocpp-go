package core

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Clear Cache (CS -> CP) --------------------

const ClearCacheFeatureName = "ClearCache"

// Status returned in response to ClearCacheRequest.
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

// The field definition of the ClearCache request payload sent by the Central System to the Charge Point.
type ClearCacheRequest struct {
}

// This field definition of the ClearCache confirmation payload, sent by the Charge Point to the Central System in response to a ClearCacheRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type ClearCacheConfirmation struct {
	Status ClearCacheStatus `json:"status" validate:"required,cacheStatus16"`
}

// Central System can request a Charge Point to clear its Authorization Cache.
// The Central System SHALL send a ClearCacheRequest PDU for clearing the Charge Point’s Authorization Cache.
// Upon receipt of a ClearCacheRequest, the Charge Point SHALL respond with a ClearCacheConfirmation PDU.
// The response PDU SHALL indicate whether the Charge Point was able to clear its Authorization Cache.
type ClearCacheFeature struct{}

func (f ClearCacheFeature) GetFeatureName() string {
	return ClearCacheFeatureName
}

func (f ClearCacheFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(ClearCacheRequest{})
}

func (f ClearCacheFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(ClearCacheConfirmation{})
}

func (r ClearCacheRequest) GetFeatureName() string {
	return ClearCacheFeatureName
}

func (c ClearCacheConfirmation) GetFeatureName() string {
	return ClearCacheFeatureName
}

// Creates a new ClearCacheRequest, which doesn't contain any required or optional fields.
func NewClearCacheRequest() *ClearCacheRequest {
	return &ClearCacheRequest{}
}

// Creates a new ClearCacheConfirmation, containing all required fields. There are no optional fields for this message.
func NewClearCacheConfirmation(status ClearCacheStatus) *ClearCacheConfirmation {
	return &ClearCacheConfirmation{Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("cacheStatus16", isValidClearCacheStatus)
}
