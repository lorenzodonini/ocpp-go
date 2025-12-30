package der

import (
	"github.com/xBlaz3kx/ocpp-go/ocpp2.0.1/types"
	"reflect"
)

// -------------------- NotifyDERStartStop (CSMS -> CS) --------------------

const NotifyDERStartStop = "NotifyDERStartStop"

// The field definition of the NotifyDERStartStopRequest request payload sent by the CSMS to the Charging Station.
type NotifyDERStartStopRequest struct {
	ControlId     string         `json:"controlId" validate:"required,max=36"`
	Started       bool           `json:"started" validate:"required"` // Indicates whether the DER is started or stopped.
	Timestamp     types.DateTime `json:"timestamp" validate:"required"`
	SupersededIds []string       `json:"supersededIds,omitempty" validate:"omitempty,max=24"`
}

// This field definition of the NotifyDERStartStopResponse
type NotifyDERStartStopResponse struct {
}

type NotifyDERStartStopFeature struct{}

func (f NotifyDERStartStopFeature) GetFeatureName() string {
	return NotifyDERStartStop
}

func (f NotifyDERStartStopFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(NotifyDERStartStopRequest{})
}

func (f NotifyDERStartStopFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(NotifyDERStartStopResponse{})
}

func (r NotifyDERStartStopRequest) GetFeatureName() string {
	return NotifyDERStartStop
}

func (c NotifyDERStartStopResponse) GetFeatureName() string {
	return NotifyDERStartStop
}

// Creates a new NotifyDERStartStopRequest, containing all required fields. Optional fields may be set afterwards.
func NewNotifyDERStartStopRequest(
	controlId string,
	started bool,
	timestamp types.DateTime,
) *NotifyDERStartStopRequest {
	return &NotifyDERStartStopRequest{
		ControlId: controlId,
		Started:   started,
		Timestamp: timestamp,
	}
}

// Creates a new NotifyDERStartStopResponse, containing all required fields. Optional fields may be set afterwards.
func NewNotifyDERStartStopResponse() *NotifyDERStartStopResponse {
	return &NotifyDERStartStopResponse{}
}
