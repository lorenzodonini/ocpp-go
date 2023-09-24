// Open Charge Point Protocol (OCPP) is a standard open protocol for communication between Charge Points and a Central System and is designed to accommodate any type of charging technique.
// This package contains the base interfaces used for the OCPP 1.6 and OCPP 2.0 implementation.
package ocpp

import (
	"fmt"
	"reflect"
)

// Feature represents a single functionality, associated to a unique name.
// Every feature describes a Request and a Response message, specified in the OCPP protocol.
type Feature interface {
	// Returns the unique name of the feature.
	GetFeatureName() string
	// Returns the type of the request message.
	GetRequestType() reflect.Type
	// Returns the type of the response message.
	GetResponseType() reflect.Type
}

// Request message
type Request interface {
	// Returns the unique name of the feature, to which this request belongs to.
	GetFeatureName() string
}

// Response message
type Response interface {
	// Returns the unique name of the feature, to which this request belongs to.
	GetFeatureName() string
}

// ErrorCode defines a common code name for an error.
type ErrorCode string

// Error wraps an OCPP error, containing an ErrorCode, a Description and the ID of the message.
type Error struct {
	Code        ErrorCode
	Description string
	MessageId   string
}

// Creates a new OCPP Error.
func NewError(errorCode ErrorCode, description string, messageId string) *Error {
	return &Error{Code: errorCode, Description: description, MessageId: messageId}
}

// Creates a new OCPP Error without messageId, which is added by the handlers parent.
func NewHandlerError(errorCode ErrorCode, description string) *Error {
	return &Error{Code: errorCode, Description: description, MessageId: ""}
}

func (err *Error) Error() string {
	return fmt.Sprintf("ocpp message (%s): %v - %v", err.MessageId, err.Code, err.Description)
}

// -------------------- Profile --------------------

// Profile defines a specific set of features, grouped by functionality.
//
// Some vendor may want to keep the protocol as slim as possible, and only support some feature profiles.
// This can easily be achieved by only registering certain profiles, while remaining compliant with the specifications.
type Profile struct {
	Name     string
	Features map[string]Feature
}

// Creates a new profile, identified by a name and a set of features.
func NewProfile(name string, features ...Feature) *Profile {
	profile := Profile{Name: name, Features: make(map[string]Feature)}
	for _, feature := range features {
		profile.AddFeature(feature)
	}
	return &profile
}

// Adds a feature to the profile.
func (p *Profile) AddFeature(feature Feature) {
	p.Features[feature.GetFeatureName()] = feature
}

// SupportsFeature returns true if a feature matching the the passed name is registered with this profile, false otherwise.
func (p *Profile) SupportsFeature(name string) bool {
	_, ok := p.Features[name]
	return ok
}

// Retrieves a feature, identified by a unique name.
// Returns nil in case the feature is not registered with this profile.
func (p *Profile) GetFeature(name string) Feature {
	return p.Features[name]
}

// ParseRequest checks whether a feature is supported and passes the rawRequest message to the requestParser function.
// The type of the request message is passed to the requestParser function, which has to perform type assertion.
func (p *Profile) ParseRequest(featureName string, rawRequest interface{}, requestParser func(raw interface{}, requestType reflect.Type) (Request, error)) (Request, error) {
	feature, ok := p.Features[featureName]
	if !ok {
		return nil, fmt.Errorf("Feature %s not found", featureName)
	}
	requestType := feature.GetRequestType()
	return requestParser(rawRequest, requestType)
}

// ParseRequest checks whether a feature is supported and passes the rawResponse message to the responseParser function.
// The type of the response message is passed to the responseParser function, which has to perform type assertion.
func (p *Profile) ParseResponse(featureName string, rawResponse interface{}, responseParser func(raw interface{}, responseType reflect.Type) (Response, error)) (Response, error) {
	feature, ok := p.Features[featureName]
	if !ok {
		return nil, fmt.Errorf("Feature %s not found", featureName)
	}
	responseType := feature.GetResponseType()
	return responseParser(rawResponse, responseType)
}
