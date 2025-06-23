package diagnostics

import (
	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
	"reflect"
)

// --------------------  Adjust Periodic EventStream (CSMS -> CS) --------------------

const AdjustPeriodicEventStream = "ClosePeriodicEventStream"

// The field definition of the AdjustPeriodicEventStreamRequest request payload sent by the CSMS to the Charging Station.
type AdjustPeriodicEventStreamRequest struct {
	Id     int                       `json:"id" validate:"required,gte=0"`
	Params PeriodicEventStreamParams `json:"params" validate:"required,dive"`
}

// This field definition of the AdjustPeriodicEventStream response payload, sent by the Charging Station to the CSMS in response to a AdjustPeriodicEventStreamRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type AdjustPeriodicEventStreamResponse struct {
	Status     types.GenericStatus `json:"status" validate:"required,genericStatus21"`
	StatusInfo *types.StatusInfo   `json:"statusInfo,omitempty" validate:"omitempty,dive"`
}

type AdjustPeriodicEventStreamFeature struct{}

func (f AdjustPeriodicEventStreamFeature) GetFeatureName() string {
	return AdjustPeriodicEventStream
}

func (f AdjustPeriodicEventStreamFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(AdjustPeriodicEventStreamRequest{})
}

func (f AdjustPeriodicEventStreamFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(AdjustPeriodicEventStreamResponse{})
}

func (r AdjustPeriodicEventStreamRequest) GetFeatureName() string {
	return AdjustPeriodicEventStream
}

func (c AdjustPeriodicEventStreamResponse) GetFeatureName() string {
	return AdjustPeriodicEventStream
}

// Creates a new AdjustPeriodicEventStreamRequest, containing all required fields. There are no optional fields for this message.
func NewAdjustPeriodicEventStreamsRequest(id int, params PeriodicEventStreamParams) *AdjustPeriodicEventStreamRequest {
	return &AdjustPeriodicEventStreamRequest{
		Id:     id,
		Params: params,
	}
}

// Creates a new AdjustPeriodicEventStreamResponse, which doesn't contain any required or optional fields.
func NewAdjustPeriodicEventStreamResponse(status types.GenericStatus) *AdjustPeriodicEventStreamResponse {
	return &AdjustPeriodicEventStreamResponse{
		Status: status,
	}
}
