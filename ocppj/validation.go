package ocppj

import (
	"sync/atomic"

	"gopkg.in/go-playground/validator.v9"
)

// The validator, used for validating incoming/outgoing OCPP messages.
var Validate = validator.New()

// The internal validation settings. Enabled by default.
// Safe for concurrent use.
var validationEnabled atomic.Bool

func init() {
	_ = Validate.RegisterValidation("errorCode", IsErrorCodeValid)
	validationEnabled.Store(true)
}

// Allows to enable/disable automatic validation for OCPP messages
// (this includes the field constraints defined for every request/response).
// The feature may be useful when working with OCPP implementations that don't fully comply to the specs.
//
// Validation is enabled by default.
//
// ⚠️ Use at your own risk! When disabled, outgoing and incoming OCPP messages will not be validated anymore,
// potentially leading to errors.
func SetMessageValidation(enabled bool) {
	validationEnabled.Store(enabled)
}
