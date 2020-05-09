package core

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Change Configuration (CS -> CP) --------------------

const ChangeConfigurationFeatureName = "ChangeConfiguration"

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

// Central System can request a Charge Point to change configuration parameters, by sending a ChangeConfigurationRequest.
// This request contains a key-value pair, where "key" is the name of the configuration setting to change and "value" contains the new setting for the configuration setting.
// A Charge Point SHALL reply with a ChangeConfigurationConfirmation indicating whether it was able to apply the change to its configuration.
// The Charge Point SHALL set the status field in the ChangeConfiguration.conf according to the following rules:
// - If the change was applied successfully, and the change if effective immediately, the Charge Point SHALL respond with a status 'Accepted'.
// - If the change was applied successfully, but a reboot is needed to make it effective, the Charge Point SHALL respond with status 'RebootRequired'.
// - If "key" does not correspond to a configuration setting supported by Charge Point, it SHALL respond with a status 'NotSupported'.
// - If the Charge Point did not set the configuration, and none of the previous statuses applies, the Charge Point SHALL respond with status 'Rejected'.
//
// If a key value is defined as a CSL, it MAY be accompanied with a [KeyName]MaxLength key, indicating the max length of the CSL in items. If this key is not set, a safe value of 1 (one) item SHOULD be assumed.
type ChangeConfigurationFeature struct{}

func (f ChangeConfigurationFeature) GetFeatureName() string {
	return ChangeConfigurationFeatureName
}

func (f ChangeConfigurationFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(ChangeConfigurationRequest{})
}

func (f ChangeConfigurationFeature) GetResponseType() reflect.Type {
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
	_ = types.Validate.RegisterValidation("configurationStatus", isValidConfigurationStatus)
}
