package remotecontrol

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
)

// -------------------- Request Start Transaction (CSMS -> CS) --------------------

const RequestStartTransactionFeatureName = "RequestStartTransaction"

// RequestStartTransactionRequest is the payload for a RequestStartTransaction request from CSMS to CS.
type RequestStartTransactionRequest struct {
	EvseID          *int                   `json:"evseId,omitempty" validate:"omitempty,gte=1"`
	RemoteStartID   int                    `json:"remoteStartId" validate:"gte=0"`
	IdToken         types.IdToken          `json:"idToken" validate:"required"`
	ChargingProfile *types.ChargingProfile `json:"chargingProfile,omitempty" validate:"omitempty"`
	GroupIdToken    *types.GroupIdToken    `json:"groupIdToken,omitempty" validate:"omitempty"`
	CustomData      *types.CustomDataType  `json:"customData,omitempty" validate:"omitempty"`
}

// RequestStartTransactionResponse is the payload for a RequestStartTransaction response from CS to CSMS.
type RequestStartTransactionResponse struct {
	Status        types.RemoteStartStopStatus `json:"status" validate:"required,remoteStartStopStatus21"`
	StatusInfo    *types.StatusInfo           `json:"statusInfo,omitempty" validate:"omitempty"`
	TransactionID string                      `json:"transactionId,omitempty" validate:"omitempty,max=36"`
	CustomData    *types.CustomDataType       `json:"customData,omitempty" validate:"omitempty"`
}

// RequestStartTransactionFeature represents the RequestStartTransaction feature.
type RequestStartTransactionFeature struct{}

func (f RequestStartTransactionFeature) GetFeatureName() string {
	return RequestStartTransactionFeatureName
}

func (f RequestStartTransactionFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(RequestStartTransactionRequest{})
}

func (f RequestStartTransactionFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(RequestStartTransactionResponse{})
}

func (r RequestStartTransactionRequest) GetFeatureName() string {
	return RequestStartTransactionFeatureName
}

func (c RequestStartTransactionResponse) GetFeatureName() string {
	return RequestStartTransactionFeatureName
}

// NewRequestStartTransactionRequest creates a new RequestStartTransactionRequest.
func NewRequestStartTransactionRequest(remoteStartID int, idToken types.IdToken) *RequestStartTransactionRequest {
	return &RequestStartTransactionRequest{RemoteStartID: remoteStartID, IdToken: idToken}
}

// NewRequestStartTransactionResponse creates a new RequestStartTransactionResponse.
func NewRequestStartTransactionResponse(status types.RemoteStartStopStatus) *RequestStartTransactionResponse {
	return &RequestStartTransactionResponse{Status: status}
}
