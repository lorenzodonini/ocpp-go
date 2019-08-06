package ocpp

import (
	errors2 "github.com/pkg/errors"
	"reflect"
)

type Feature interface {
	GetFeatureName() string
	GetRequestType() reflect.Type
	GetConfirmationType() reflect.Type
}

type Request interface {
	GetFeatureName() string
}

type Confirmation interface {
	GetFeatureName() string
}

type ErrorCode string

type OcppError struct {
	Error error
	ErrorCode ErrorCode
	MessageId string
}

// -------------------- Profile --------------------
type Profile struct {
	Name     string
	Features map[string]Feature
}

func NewProfile(name string, features ...Feature) *Profile {
	profile := Profile{Name: name, Features: make(map[string]Feature)}
	for _, feature := range features {
		profile.AddFeature(feature)
	}
	return &profile
}

func (p *Profile) AddFeature(feature Feature) {
	p.Features[feature.GetFeatureName()] = feature
}

func (p *Profile) SupportsFeature(name string) bool {
	_, ok := p.Features[name]
	return ok
}

func (p *Profile) GetFeature(name string) Feature {
	return p.Features[name]
}

func (p *Profile) ParseRequest(featureName string, rawRequest interface{}, requestParser func(raw interface{}, requestType reflect.Type) (Request, error)) (Request, error) {
	feature, ok := p.Features[featureName]
	if !ok {
		return nil, errors2.Errorf("Feature %s not found", featureName)
	}
	requestType := feature.GetRequestType()
	return requestParser(rawRequest, requestType)
}

func (p *Profile) ParseConfirmation(featureName string, rawConfirmation interface{}, confirmationParser func(raw interface{}, confirmationType reflect.Type) (Confirmation, error)) (Confirmation, error) {
	feature, ok := p.Features[featureName]
	if !ok {
		return nil, errors2.Errorf("Feature %s not found", featureName)
	}
	confirmationType := feature.GetConfirmationType()
	return confirmationParser(rawConfirmation, confirmationType)
}
