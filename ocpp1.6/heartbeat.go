package ocpp16

import (
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Heartbeat (CP -> CS) --------------------
type HeartbeatRequest struct {
}

type HeartbeatConfirmation struct {
	CurrentTime DateTime `json:"currentTime" validate:"required"`
}

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

func NewHeartbeatRequest() *HeartbeatRequest {
	return &HeartbeatRequest{}
}

func NewHeartbeatConfirmation(currentTime DateTime) *HeartbeatConfirmation {
	return &HeartbeatConfirmation{CurrentTime: currentTime}
}

func validateHeartbeatConfirmation(sl validator.StructLevel) {
	confirmation := sl.Current().Interface().(HeartbeatConfirmation)
	if dateTimeIsNull(confirmation.CurrentTime) {
		sl.ReportError(confirmation.CurrentTime, "CurrentTime", "currentTime", "required", "")
	}
	//if !validateDateTimeNow(confirmation.CurrentTime) {
	//	sl.ReportError(confirmation.CurrentTime, "CurrentTime", "currentTime", "eq", "")
	//}
}

func init() {
	Validate.RegisterStructValidation(validateHeartbeatConfirmation, HeartbeatConfirmation{})
}
