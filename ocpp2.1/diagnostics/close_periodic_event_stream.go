package diagnostics

import "reflect"

// --------------------  Close Periodic EventStream (CSMS -> CS) --------------------

const ClosePeriodicEventStream = "ClosePeriodicEventStream"

// The field definition of the ClosePeriodicEventStreamRequest request payload sent by the CSMS to the Charging Station.
type ClosePeriodicEventStreamRequest struct {
	Id int `json:"id" validate:"required,gte=0"`
}

// This field definition of the ClosePeriodicEventStream response payload, sent by the Charging Station to the CSMS in response to a ClosePeriodicEventStreamRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type ClosePeriodicEventStreamResponse struct {
}

type ClosePeriodicEventStreamFeature struct{}

func (f ClosePeriodicEventStreamFeature) GetFeatureName() string {
	return ClosePeriodicEventStream
}

func (f ClosePeriodicEventStreamFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(ClosePeriodicEventStreamRequest{})
}

func (f ClosePeriodicEventStreamFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(ClosePeriodicEventStreamResponse{})
}

func (r ClosePeriodicEventStreamRequest) GetFeatureName() string {
	return ClosePeriodicEventStream
}

func (c ClosePeriodicEventStreamResponse) GetFeatureName() string {
	return ClosePeriodicEventStream
}

// Creates a new ClosePeriodicEventStreamRequest, containing all required fields. There are no optional fields for this message.
func NewClosePeriodicEventStreamsRequest(id int) *ClosePeriodicEventStreamRequest {
	return &ClosePeriodicEventStreamRequest{
		Id: id,
	}
}

// Creates a new ClosePeriodicEventStreamResponse, which doesn't contain any required or optional fields.
func NewClosePeriodicEventStreamResponse() *ClosePeriodicEventStreamResponse {
	return &ClosePeriodicEventStreamResponse{}
}
