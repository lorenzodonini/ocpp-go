package availability

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
)

// -------------------- Heartbeat (CS -> CSMS) --------------------

const HeartbeatFeatureName = "Heartbeat"

// HeartbeatRequest is the payload for a Heartbeat request from CS to CSMS.
type HeartbeatRequest struct {
	CustomData *types.CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// HeartbeatResponse is the payload for a Heartbeat response from CSMS to CS.
type HeartbeatResponse struct {
	CurrentTime types.DateTime        `json:"currentTime" validate:"required"`
	CustomData  *types.CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// HeartbeatFeature represents the Heartbeat feature.
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

// NewHeartbeatRequest creates a new HeartbeatRequest.
func NewHeartbeatRequest() *HeartbeatRequest {
	return &HeartbeatRequest{}
}

// NewHeartbeatResponse creates a new HeartbeatResponse.
func NewHeartbeatResponse(currentTime types.DateTime) *HeartbeatResponse {
	return &HeartbeatResponse{CurrentTime: currentTime}
}
