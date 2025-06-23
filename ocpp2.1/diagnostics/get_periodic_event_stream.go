package diagnostics

import "reflect"

// --------------------  Get Periodic EventStream (CSMS -> CS) --------------------

const GetPeriodicEventStream = "GetPeriodicEventStream"

// The field definition of the GetPeriodicEventStreamRequest request payload sent by the CSMS to the Charging Station.
type GetPeriodicEventStreamRequest struct {
}

// This field definition of the GetPeriodicEventStream response payload, sent by the Charging Station to the CSMS in response to a GetPeriodicEventStreamRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type GetPeriodicEventStreamResponse struct {
	ConstantStreamData []ConstantStreamData `json:"constantStreamData,omitempty" validate:"omitempty,dive"`
}

type ConstantStreamData struct {
	Id                   int                       `json:"id" validate:"required,gte=0"`
	VariableMonitoringId int                       `json:"variableMonitoringId" validate:"required,gte=0"`
	Params               PeriodicEventStreamParams `json:"params" validate:"required,dive"`
}

type GetPeriodicEventStreamFeature struct{}

func (f GetPeriodicEventStreamFeature) GetFeatureName() string {
	return GetPeriodicEventStream
}

func (f GetPeriodicEventStreamFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(GetPeriodicEventStreamRequest{})
}

func (f GetPeriodicEventStreamFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(GetPeriodicEventStreamResponse{})
}

func (r GetPeriodicEventStreamRequest) GetFeatureName() string {
	return GetPeriodicEventStream
}

func (c GetPeriodicEventStreamResponse) GetFeatureName() string {
	return GetPeriodicEventStream
}

// Creates a new GetPeriodicEventStreamRequest, containing all required fields. There are no optional fields for this message.
func NewGetPeriodicEventStreamsRequest() *GetPeriodicEventStreamRequest {
	return &GetPeriodicEventStreamRequest{}
}

// Creates a new GetPeriodicEventStreamResponse, which doesn't contain any required or optional fields.
func NewGetPeriodicEventStreamResponse() *GetPeriodicEventStreamResponse {
	return &GetPeriodicEventStreamResponse{}
}
