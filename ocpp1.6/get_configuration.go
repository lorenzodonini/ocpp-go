package ocpp16

import (
	"reflect"
)

// -------------------- Get Configuration (CS -> CP) --------------------
type ConfigurationKey struct {
	Key      string `json:"key" validate:"required,max=50"`
	Readonly bool   `json:"readonly"`
	Value    string `json:"value,omitempty" validate:"max=500"`
}

type GetConfigurationRequest struct {
	Key           []string `json:"key" validate:"required,min=1,unique,dive,max=50"`
}

// TODO: validation of cardinalities for the two fields should be handled somewhere (#configurationKey + #unknownKey > 0)
// TODO: add uniqueness of configurationKey in slice, once PR is merged (https://github.com/go-playground/validator/pull/496)
type GetConfigurationConfirmation struct {
	ConfigurationKey   []ConfigurationKey `json:"configurationKey,omitempty" validate:"dive"`
	UnknownKey         []string           `json:"unknownKey,omitempty" validate:"dive,max=50"`
}

type GetConfigurationFeature struct{}

func (f GetConfigurationFeature) GetFeatureName() string {
	return GetConfigurationFeatureName
}

func (f GetConfigurationFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(GetConfigurationRequest{})
}

func (f GetConfigurationFeature) GetConfirmationType() reflect.Type {
	return reflect.TypeOf(GetConfigurationConfirmation{})
}

func (r GetConfigurationRequest) GetFeatureName() string {
	return GetConfigurationFeatureName
}

func (c GetConfigurationConfirmation) GetFeatureName() string {
	return GetConfigurationFeatureName
}

func NewGetConfigurationRequest(keys []string) *GetConfigurationRequest {
	return &GetConfigurationRequest{Key: keys}
}

func NewGetConfigurationConfirmation(configurationKey []ConfigurationKey) *GetConfigurationConfirmation {
	return &GetConfigurationConfirmation{ConfigurationKey: configurationKey}
}
