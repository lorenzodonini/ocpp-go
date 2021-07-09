package core

import (
	"reflect"
)

// -------------------- Get Configuration (CS -> CP) --------------------

const GetConfigurationFeatureName = "GetConfiguration"

// Contains information about a specific configuration key. It is returned in GetConfigurationConfirmation
type ConfigurationKey struct {
	Key      string  `json:"key" validate:"required,max=50"`
	Readonly bool    `json:"readonly"`
	Value    *string `json:"value,omitempty" validate:"max=500"`
}

// The field definition of the GetConfiguration request payload sent by the Central System to the Charge Point.
type GetConfigurationRequest struct {
	Key []string `json:"key,omitempty" validate:"omitempty,unique,dive,max=50"`
}

// TODO: validation of cardinalities for the two fields should be handled somewhere (#configurationKey + #unknownKey > 0)
// TODO: add uniqueness of configurationKey in slice, once PR is merged (https://github.com/go-playground/validator/pull/496)
// This field definition of the GetConfiguration confirmation payload, sent by the Charge Point to the Central System in response to a GetConfigurationRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type GetConfigurationConfirmation struct {
	ConfigurationKey []ConfigurationKey `json:"configurationKey,omitempty" validate:"omitempty,dive"`
	UnknownKey       []string           `json:"unknownKey,omitempty" validate:"omitempty,dive,max=50"`
}

// To retrieve the value of configuration settings, the Central System SHALL send a GetConfigurationRequest to the Charge Point.
// If the list of keys in the request is empty or missing (it is optional), the Charge Point SHALL return a list of all configuration settings in GetConfigurationConfirmation.
// Otherwise Charge Point SHALL return a list of recognized keys and their corresponding values and read-only state.
// Unrecognized keys SHALL be placed in the response payload as part of the optional unknown key list element of GetConfigurationConfirmation.
// The number of configuration keys requested in a single payload MAY be limited by the Charge Point.
// This maximum can be retrieved by reading the configuration key GetConfigurationMaxKeys.
type GetConfigurationFeature struct{}

func (f GetConfigurationFeature) GetFeatureName() string {
	return GetConfigurationFeatureName
}

func (f GetConfigurationFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(GetConfigurationRequest{})
}

func (f GetConfigurationFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(GetConfigurationConfirmation{})
}

func (r GetConfigurationRequest) GetFeatureName() string {
	return GetConfigurationFeatureName
}

func (c GetConfigurationConfirmation) GetFeatureName() string {
	return GetConfigurationFeatureName
}

// Creates a new GetConfigurationRequest, containing all required fields. There are no optional fields for this message.
func NewGetConfigurationRequest(keys []string) *GetConfigurationRequest {
	return &GetConfigurationRequest{Key: keys}
}

// Creates a new GetConfigurationConfirmation, containing all required fields. Optional fields may be set afterwards.
func NewGetConfigurationConfirmation(configurationKey []ConfigurationKey) *GetConfigurationConfirmation {
	return &GetConfigurationConfirmation{ConfigurationKey: configurationKey}
}
