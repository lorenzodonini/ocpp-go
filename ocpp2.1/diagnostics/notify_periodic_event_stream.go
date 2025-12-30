package diagnostics

import (
	"reflect"

	"github.com/xBlaz3kx/ocpp-go/ocpp2.1/types"
)

const (
	NotifyPeriodicEventStreamFeat = "NotifyPeriodicEventStream"
)

// The field definition of the NotifyPeriodicEventStream request payload sent by the Charging station to the CSMS.
type NotifyPeriodicEventStream struct {
	ID       int                 `json:"id" validate:"required,gte=0"`
	Pending  int                 `json:"pending" validate:"required,gte=0"`
	BaseTime types.DateTime      `json:"baseTime" validate:"required"`
	Data     []StreamDataElement `json:"data" validate:"required,dive"` // A list of StreamDataElements, each containing a stream of data.
}

type StreamDataElement struct {
	T float64 `json:"t" validate:"required"`
	V string  `json:"v" validate:"required,max=2500"`
}

// Note: This feature does not have a response. This needs to be reflected in the websocket layer.
type NotifyPeriodicEventStreamFeature struct{}

func (f NotifyPeriodicEventStreamFeature) GetFeatureName() string {
	return NotifyPeriodicEventStreamFeat
}

func (f NotifyPeriodicEventStreamFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(NotifyPeriodicEventStream{})
}

func (f NotifyPeriodicEventStreamFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(NotifyPeriodicEventStream{})
}

func (r NotifyPeriodicEventStream) GetFeatureName() string {
	return NotifyPeriodicEventStreamFeat
}

// Creates a new NotifyPeriodicEventStream, containing all required fields. Additional optional fields may be set afterwards.
func NewNotifyPeriodicEventStream(id, pending int, baseTime types.DateTime, data []StreamDataElement) *NotifyPeriodicEventStream {
	return &NotifyPeriodicEventStream{
		ID:       id,
		Pending:  pending,
		BaseTime: baseTime,
		Data:     data,
	}
}
