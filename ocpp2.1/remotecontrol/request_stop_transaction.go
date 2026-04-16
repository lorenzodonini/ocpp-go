package remotecontrol

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
)

// -------------------- Request Stop Transaction (CSMS -> CS) --------------------

const RequestStopTransactionFeatureName = "RequestStopTransaction"

// RequestStopTransactionRequest is the payload for a RequestStopTransaction request from CSMS to CS.
type RequestStopTransactionRequest struct {
	TransactionID string                `json:"transactionId" validate:"required,max=36"`
	CustomData    *types.CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// RequestStopTransactionResponse is the payload for a RequestStopTransaction response from CS to CSMS.
type RequestStopTransactionResponse struct {
	Status     types.RemoteStartStopStatus `json:"status" validate:"required,remoteStartStopStatus21"`
	StatusInfo *types.StatusInfo           `json:"statusInfo,omitempty" validate:"omitempty"`
	CustomData *types.CustomDataType       `json:"customData,omitempty" validate:"omitempty"`
}

// RequestStopTransactionFeature represents the RequestStopTransaction feature.
type RequestStopTransactionFeature struct{}

func (f RequestStopTransactionFeature) GetFeatureName() string {
	return RequestStopTransactionFeatureName
}

func (f RequestStopTransactionFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(RequestStopTransactionRequest{})
}

func (f RequestStopTransactionFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(RequestStopTransactionResponse{})
}

func (r RequestStopTransactionRequest) GetFeatureName() string {
	return RequestStopTransactionFeatureName
}

func (c RequestStopTransactionResponse) GetFeatureName() string {
	return RequestStopTransactionFeatureName
}

// NewRequestStopTransactionRequest creates a new RequestStopTransactionRequest.
func NewRequestStopTransactionRequest(transactionID string) *RequestStopTransactionRequest {
	return &RequestStopTransactionRequest{TransactionID: transactionID}
}

// NewRequestStopTransactionResponse creates a new RequestStopTransactionResponse.
func NewRequestStopTransactionResponse(status types.RemoteStartStopStatus) *RequestStopTransactionResponse {
	return &RequestStopTransactionResponse{Status: status}
}
