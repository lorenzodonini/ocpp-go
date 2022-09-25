package validate

import "gopkg.in/go-playground/validator.v9"

// The validator, used for validating incoming/outgoing OCPP messages.
var Validator = validator.New()

func MustRegisterValidation(tag string, fn validator.Func, callValidationEvenIfNull ...bool) {
	if err := Validator.RegisterValidation(tag, fn, callValidationEvenIfNull...); err != nil {
		panic(err)
	}
}

func MustRegisterStructValidation(fn validator.StructLevelFunc, types ...interface{}) {
	Validator.RegisterStructValidation(fn, types...)
}
