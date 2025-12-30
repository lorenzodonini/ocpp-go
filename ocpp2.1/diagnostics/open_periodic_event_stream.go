package diagnostics

import (
	"reflect"

	"github.com/xBlaz3kx/ocpp-go/ocpp2.1/types"
)

// --------------------  Open Periodic EventStream (CSMS -> CS) --------------------

const OpenPeriodicEventStream = "OpenPeriodicEventStream"

// The field definition of the OpenPeriodicEventStreamRequest request payload sent by the CSMS to the Charging Station.
type OpenPeriodicEventStreamRequest struct {
	ConstantStreamData ConstantStreamData `json:"constantStreamData" validate:"required,dive"`
}

// This field definition of the OpenPeriodicEventStream response payload, sent by the Charging Station to the CSMS in response to a OpenPeriodicEventStreamRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type OpenPeriodicEventStreamResponse struct {
	Status     types.GenericStatus `json:"status" validate:"required,genericStatus21"`
	StatusInfo *types.StatusInfo   `json:"statusInfo,omitempty" validate:"omitempty,dive"`
}

type OpenPeriodicEventStreamFeature struct{}

func (f OpenPeriodicEventStreamFeature) GetFeatureName() string {
	return OpenPeriodicEventStream
}

func (f OpenPeriodicEventStreamFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(OpenPeriodicEventStreamRequest{})
}

func (f OpenPeriodicEventStreamFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(OpenPeriodicEventStreamResponse{})
}

func (r OpenPeriodicEventStreamRequest) GetFeatureName() string {
	return OpenPeriodicEventStream
}

func (c OpenPeriodicEventStreamResponse) GetFeatureName() string {
	return OpenPeriodicEventStream
}

// Creates a new OpenPeriodicEventStreamRequest, containing all required fields. There are no optional fields for this message.
func NewOpenPeriodicEventStreamsRequest(data ConstantStreamData) *OpenPeriodicEventStreamRequest {
	return &OpenPeriodicEventStreamRequest{
		ConstantStreamData: data,
	}
}

// Creates a new OpenPeriodicEventStreamResponse, which doesn't contain any required or optional fields.
func NewOpenPeriodicEventStreamResponse(status types.GenericStatus) *OpenPeriodicEventStreamResponse {
	return &OpenPeriodicEventStreamResponse{
		Status: status,
	}
}
