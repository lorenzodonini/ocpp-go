package remotecontrol

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
)

// -------------------- Unlock Connector (CSMS -> CS) --------------------

const UnlockConnectorFeatureName = "UnlockConnector"

// UnlockConnectorRequest is the payload for an UnlockConnector request from CSMS to CS.
type UnlockConnectorRequest struct {
	EvseID      int                   `json:"evseId" validate:"gte=0"`
	ConnectorID int                   `json:"connectorId" validate:"gte=0"`
	CustomData  *types.CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// UnlockConnectorResponse is the payload for an UnlockConnector response from CS to CSMS.
type UnlockConnectorResponse struct {
	Status     types.UnlockStatus    `json:"status" validate:"required,unlockStatus21"`
	StatusInfo *types.StatusInfo     `json:"statusInfo,omitempty" validate:"omitempty"`
	CustomData *types.CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// UnlockConnectorFeature represents the UnlockConnector feature.
type UnlockConnectorFeature struct{}

func (f UnlockConnectorFeature) GetFeatureName() string {
	return UnlockConnectorFeatureName
}

func (f UnlockConnectorFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(UnlockConnectorRequest{})
}

func (f UnlockConnectorFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(UnlockConnectorResponse{})
}

func (r UnlockConnectorRequest) GetFeatureName() string {
	return UnlockConnectorFeatureName
}

func (c UnlockConnectorResponse) GetFeatureName() string {
	return UnlockConnectorFeatureName
}

// NewUnlockConnectorRequest creates a new UnlockConnectorRequest.
func NewUnlockConnectorRequest(evseID int, connectorID int) *UnlockConnectorRequest {
	return &UnlockConnectorRequest{EvseID: evseID, ConnectorID: connectorID}
}

// NewUnlockConnectorResponse creates a new UnlockConnectorResponse.
func NewUnlockConnectorResponse(status types.UnlockStatus) *UnlockConnectorResponse {
	return &UnlockConnectorResponse{Status: status}
}
