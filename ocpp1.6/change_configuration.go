package ocpp16

import (
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Change Configuration (CS -> CP) --------------------

// Status in ChangeConfigurationConfirmation.
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

// The field definition of the ChangeConfiguration request payload sent by the Central System to the Charge Point.
type ChangeConfigurationRequest struct {
	Key   string `json:"key" validate:"required,max=50"`
	Value string `json:"value" validate:"required,max=500"`
}

// This field definition of the ChangeConfiguration confirmation payload, sent by the Charge Point to the Central System in response to a ChangeConfigurationRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type ChangeConfigurationConfirmation struct {
	Status ConfigurationStatus `json:"status" validate:"required,configurationStatus"`
}

// Central System can request a Charge Point to change configuration parameters.
// To achieve this, Central System SHALL send a ChangeConfigurationRequest.
// This request contains a key-value pair, where "key" is the name of the configuration setting to change and "value" contains the new setting for the configuration setting.
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

// Creates a new ChangeConfigurationRequest, containing all required fields. There are no optional fields for this message.
func NewChangeConfigurationRequest(key string, value string) *ChangeConfigurationRequest {
	return &ChangeConfigurationRequest{Key: key, Value: value}
}

// Creates a new ChangeConfigurationConfirmation, containing all required fields. There are no optional fields for this message.
func NewChangeConfigurationConfirmation(status ConfigurationStatus) *ChangeConfigurationConfirmation {
	return &ChangeConfigurationConfirmation{Status: status}
}

func init() {
	_ = Validate.RegisterValidation("configurationStatus", isValidConfigurationStatus)
}
