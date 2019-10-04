package ocpp16

import (
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Heartbeat (CP -> CS) --------------------

// The field definition of the Heartbeat request payload sent by the Charge Point to the Central System.
type HeartbeatRequest struct {
}

// This field definition of the Heartbeat confirmation payload, sent by the Central System to the Charge Point in response to a HeartbeatRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type HeartbeatConfirmation struct {
	CurrentTime *DateTime `json:"currentTime" validate:"required"`
}

// To let the Central System know that a Charge Point is still connected, a Charge Point sends a heartbeat after a configurable time interval.
type HeartbeatFeature struct{}

func (f HeartbeatFeature) GetFeatureName() string {
	return HeartbeatFeatureName
}

func (f HeartbeatFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(HeartbeatRequest{})
}

func (f HeartbeatFeature) GetConfirmationType() reflect.Type {
	return reflect.TypeOf(HeartbeatConfirmation{})
}

func (r HeartbeatRequest) GetFeatureName() string {
	return HeartbeatFeatureName
}

func (c HeartbeatConfirmation) GetFeatureName() string {
	return HeartbeatFeatureName
}

// Creates a new HeartbeatRequest, which doesn't contain any required or optional fields.
func NewHeartbeatRequest() *HeartbeatRequest {
	return &HeartbeatRequest{}
}

// Creates a new HeartbeatConfirmation, containing all required fields.
func NewHeartbeatConfirmation(currentTime *DateTime) *HeartbeatConfirmation {
	return &HeartbeatConfirmation{CurrentTime: currentTime}
}

func validateHeartbeatConfirmation(sl validator.StructLevel) {
	confirmation := sl.Current().Interface().(HeartbeatConfirmation)
	if dateTimeIsNull(confirmation.CurrentTime) {
		sl.ReportError(confirmation.CurrentTime, "CurrentTime", "currentTime", "required", "")
	}
}

func init() {
	Validate.RegisterStructValidation(validateHeartbeatConfirmation, HeartbeatConfirmation{})
}
