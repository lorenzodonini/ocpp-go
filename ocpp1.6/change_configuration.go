package ocpp16

import (
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Change Configuration (CS -> CP) --------------------
type ConfigurationStatus string

const (
	ConfigurationStatusAccepted       ConfigurationStatus = "Accepted"
	ConfigurationStatusRejected       ConfigurationStatus = "Rejected"
	ConfigurationStatusRebootRequired ConfigurationStatus = "RebootRequired"
	ConfigurationStatusNotSupported   ConfigurationStatus = "NotSupported"
)

func isValidConfigurationStatus(fl validator.FieldLevel) bool {
	status := ConfigurationStatus(fl.Field().String())
	switch status {
	case ConfigurationStatusAccepted, ConfigurationStatusRejected, ConfigurationStatusRebootRequired, ConfigurationStatusNotSupported:
		return true
	default:
		return false
	}
}

type ChangeConfigurationRequest struct {
	Key           string `json:"key" validate:"required,max=50"`
	Value         string `json:"value" validate:"required,max=500"`
}

type ChangeConfigurationConfirmation struct {
	Status             ConfigurationStatus `json:"status" validate:"required,configurationStatus"`
}

type ChangeConfigurationFeature struct{}

func (f ChangeConfigurationFeature) GetFeatureName() string {
	return ChangeConfigurationFeatureName
}

func (f ChangeConfigurationFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(ChangeConfigurationRequest{})
}

func (f ChangeConfigurationFeature) GetConfirmationType() reflect.Type {
	return reflect.TypeOf(ChangeConfigurationConfirmation{})
}

func (r ChangeConfigurationRequest) GetFeatureName() string {
	return ChangeConfigurationFeatureName
}

func (c ChangeConfigurationConfirmation) GetFeatureName() string {
	return ChangeConfigurationFeatureName
}

func NewChangeConfigurationRequest(key string, value string) *ChangeConfigurationRequest {
	return &ChangeConfigurationRequest{Key: key, Value: value}
}

func NewChangeConfigurationConfirmation(status ConfigurationStatus) *ChangeConfigurationConfirmation {
	return &ChangeConfigurationConfirmation{Status: status}
}

func init() {
	_ = Validate.RegisterValidation("configurationStatus", isValidConfigurationStatus)
}
