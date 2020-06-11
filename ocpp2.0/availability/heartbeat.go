package availability

import (
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Heartbeat (CS -> CSMS) --------------------

const HeartbeatFeatureName = "Heartbeat"

// The field definition of the Heartbeat request payload sent by the Charging Station to the CSMS.
type HeartbeatRequest struct {
}

// This field definition of the Heartbeat response payload, sent by the CSMS to the Charging Station in response to a HeartbeatRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type HeartbeatResponse struct {
	CurrentTime types.DateTime `json:"currentTime" validate:"required"`
}

// A Charging Station may send a heartbeat to let the CSMS know the Charging Station is still connected, after a configurable time interval.
//
// Upon receipt of HeartbeatRequest, the CSMS responds with HeartbeatResponse.
// The response message contains the current time of the CSMS, which the Charging Station MAY use to synchronize its internal clock.
type HeartbeatFeature struct{}

func (f HeartbeatFeature) GetFeatureName() string {
	return HeartbeatFeatureName
}

func (f HeartbeatFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(HeartbeatRequest{})
}

func (f HeartbeatFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(HeartbeatResponse{})
}

func (r HeartbeatRequest) GetFeatureName() string {
	return HeartbeatFeatureName
}

func (c HeartbeatResponse) GetFeatureName() string {
	return HeartbeatFeatureName
}

// Creates a new HeartbeatRequest, which doesn't contain any required or optional fields.
func NewHeartbeatRequest() *HeartbeatRequest {
	return &HeartbeatRequest{}
}

// Creates a new HeartbeatResponse, containing all required fields. There are no optional fields for this message.
func NewHeartbeatResponse(currentTime types.DateTime) *HeartbeatResponse {
	return &HeartbeatResponse{CurrentTime: currentTime}
}

func validateHeartbeatResponse(sl validator.StructLevel) {
	response := sl.Current().Interface().(HeartbeatResponse)
	if types.DateTimeIsNull(&response.CurrentTime) {
		sl.ReportError(response.CurrentTime, "CurrentTime", "currentTime", "required", "")
	}
}

func init() {
	types.Validate.RegisterStructValidation(validateHeartbeatResponse, HeartbeatResponse{})
}
