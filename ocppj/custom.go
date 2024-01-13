package ocppj

import (
	"github.com/lorenzodonini/ocpp-go/ocpp"
)

// A CustomRequest allows to parse/serialize a custom request payload to/from a valid OCPP request.
// Custom types are commonly used for non-spec compliant requests (especially on server side),
// to convert them into a valid OCPP type.
//
// The struct needs to implement the ocpp.Request interface as well,
// so that the payload may be handled correctly in the ocpp-j layer.
type CustomRequest interface {
	ocpp.Request
	// Parse converts a custom request payload to a valid OCPP request.
	// The method is invoked after the payload was unmarshalled from JSON into a custom struct.
	// The passed targetPayload is the target OCPP request type, which needs to be filled with the parsed data.
	//
	// Example:
	//  func (r *CustomGetConfigurationRequest) Parse(targetPayload ocpp.Request) error {
	//  	target, ok := targetPayload.(*core.GetConfigurationRequest)
	//  	if !ok {
	//  		return fmt.Errorf("invalid type %T", targetPayload)
	//  	}
	//  	target.Key = make([]string, len(r.ConfigurationKeys))
	//  	for i, key := range r.ConfigurationKeys {
	// 			Fill target.Key[i] with data from key with appropriate transformations
	//  	}
	//  	return nil
	//  }
	//
	// The implementation may be left empty, if the request will never be parsed on this endpoint.
	Parse(targetPayload ocpp.Request) error
	// Serialize converts a valid OCPP request to a custom request payload.
	// The method is invoked before the payload is marshalled into JSON.
	// The passed srcPayload is the source OCPP request type, which needs to be converted into a custom struct.
	//
	// Example:
	//  func (r *CustomGetConfigurationRequest) Serialize(srcPayload ocpp.Request) error {
	//  	src, ok := srcPayload.(*core.CustomGetConfigurationRequest)
	//  	if !ok {
	//  		return fmt.Errorf("invalid type %T", srcPayload)
	//  	}
	//  	r.ConfigurationKeys = make([]string, len(src.Key))
	//  	for i, key := range src.Key {
	//			// Fill r.ConfigurationKeys[i] with data from key with appropriate transformations
	//  	}
	//  	return nil
	//  }
	//
	// The implementation may be left empty, if the request will never be serialized on this endpoint.
	Serialize(srcPayload ocpp.Request) error
}

// A CustomResponse allows to parse/serialize a custom response payload to/from a valid OCPP response.
// Custom types are commonly used for non-spec compliant responses (especially on server side),
// to convert them into a valid OCPP type.
//
// The struct needs to implement the ocpp.Response interface as well,
// so that the payload may be handled correctly in the ocpp-j layer.
type CustomResponse interface {
	ocpp.Response
	// Parse converts a custom response payload to a valid OCPP response.
	// The method is invoked after the payload was unmarshalled from JSON into a custom struct.
	// The passed targetPayload is the target OCPP response type, which needs to be filled with the parsed data.
	//
	// Example:
	//  func (r *CustomGetConfigurationConfirmation) Parse(targetPayload ocpp.Response) error {
	//  	target, ok := targetPayload.(*core.GetConfigurationConfirmation)
	//  	if !ok {
	//  		return fmt.Errorf("invalid type %T", targetPayload)
	//  	}
	//  	target.UnknownKey = r.UnknownKey
	//  	target.ConfigurationKey = make([]core.ConfigurationKey, len(r.ConfigurationKey))
	//  	for i, key := range r.ConfigurationKey {
	//			// Fill target.ConfigurationKey[i] with data from key with appropriate transformations
	//  	}
	//  	return nil
	//  }
	//
	// The implementation may be left empty, if the response will never be parsed on this endpoint.
	Parse(targetPayload ocpp.Response) error
	// Serialize converts a valid OCPP response to a custom response payload.
	// The method is invoked before the payload is marshalled into JSON.
	// The passed srcPayload is the source OCPP response type, which needs to be converted into a custom struct.
	//
	// Example:
	//  func (r *CustomGetConfigurationConfirmation) Serialize(srcPayload ocpp.Response) error {
	//  	src, ok := srcPayload.(*core.GetConfigurationConfirmation)
	//  	if !ok {
	//  		return fmt.Errorf("invalid type %T", srcPayload)
	//  	}
	//  	r.UnknownKey = src.UnknownKey
	//  	r.ConfigurationKey = make([]SomeCustomStruct, len(src.ConfigurationKey))
	//  	for i, key := range src.ConfigurationKey {
	//			// Fill r.ConfigurationKey[i] with data from key with appropriate transformations
	//  	}
	//  	return nil
	//  }
	//
	// The implementation may be left empty, if the response will never be serialized on this endpoint.
	Serialize(srcPayload ocpp.Response) error
}

// CustomTypeMapper allows to retrieve custom type converters for parsing and serializing
// from/to OCPP types and custom types.
type CustomTypeMapper interface {
	// GetCustomRequest returns a custom request parser/serializer for the passed feature name.
	// The feature always points to an OCPP type conforming to the ocpp.Request interface.
	// The returned CustomRequest is used for custom type transformations during JSON parsing/serialization.
	//
	// The method returns false if no custom request parser/serializer is registered for the passed type.
	GetCustomRequest(featureName string) (CustomRequest, bool)
	// GetCustomResponse returns a custom response parser/serializer for the passed feature name.
	// The feature always points to an OCPP type conforming to the ocpp.Response interface.
	// The returned CustomResponse is used for custom type transformations during JSON parsing/serialization.
	//
	// The method returns false if no custom response parser/serializer is registered for the passed type.
	GetCustomResponse(featureName string) (CustomResponse, bool)
	// SetCustomRequest adds a custom request parser/serializer for a specific OCPP type.
	// The type is retrieved via the GetOcppType() method of the passed CustomRequest.
	SetCustomRequest(request CustomRequest)
	// SetCustomResponse adds a custom response parser/serializer for a specific OCPP type.
	// The type is retrieved via the GetOcppType() method of the passed CustomResponse.
	SetCustomResponse(response CustomResponse)
}

type customTypeMapper struct {
	requests  map[string]CustomRequest
	responses map[string]CustomResponse
}

// NewCustomTypeMapper creates a simple implementation of CustomTypeMapper.
// The custom types are stored in a map, which is not thread-safe.
// To safely use the mapper in a concurrent environment, it is recommended to register all custom
// types via the SetCustomRequest/SetCustomResponse methods before usage.
//
// The mapper is meant to be read-only and can be shared between multiple endpoints.
// Refer to the SetCustomTypeMapper method of the ClientState interface for more information.
func NewCustomTypeMapper() CustomTypeMapper {
	return &customTypeMapper{
		requests:  make(map[string]CustomRequest),
		responses: make(map[string]CustomResponse),
	}
}

func (m *customTypeMapper) GetCustomRequest(featureName string) (CustomRequest, bool) {
	request, ok := m.requests[featureName]
	return request, ok
}

func (m *customTypeMapper) GetCustomResponse(featureName string) (CustomResponse, bool) {
	response, ok := m.responses[featureName]
	return response, ok
}

func (m *customTypeMapper) SetCustomRequest(request CustomRequest) {
	key := request.GetFeatureName()
	m.requests[key] = request
}

func (m *customTypeMapper) SetCustomResponse(response CustomResponse) {
	key := response.GetFeatureName()
	m.responses[key] = response
}
