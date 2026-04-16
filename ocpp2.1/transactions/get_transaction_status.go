package transactions

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
)

// -------------------- Get Transaction Status (CSMS -> CS) --------------------

const GetTransactionStatusFeatureName = "GetTransactionStatus"

// GetTransactionStatusRequest is the payload for a GetTransactionStatus request from CSMS to CS.
type GetTransactionStatusRequest struct {
	TransactionID string                `json:"transactionId,omitempty" validate:"omitempty,max=36"`
	CustomData    *types.CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// GetTransactionStatusResponse is the payload for a GetTransactionStatus response from CS to CSMS.
type GetTransactionStatusResponse struct {
	OngoingIndicator *bool                 `json:"ongoingIndicator,omitempty"`
	MessagesInQueue  bool                  `json:"messagesInQueue"`
	CustomData       *types.CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// GetTransactionStatusFeature represents the GetTransactionStatus feature.
type GetTransactionStatusFeature struct{}

func (f GetTransactionStatusFeature) GetFeatureName() string {
	return GetTransactionStatusFeatureName
}

func (f GetTransactionStatusFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(GetTransactionStatusRequest{})
}

func (f GetTransactionStatusFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(GetTransactionStatusResponse{})
}

func (r GetTransactionStatusRequest) GetFeatureName() string {
	return GetTransactionStatusFeatureName
}

func (c GetTransactionStatusResponse) GetFeatureName() string {
	return GetTransactionStatusFeatureName
}

// NewGetTransactionStatusRequest creates a new GetTransactionStatusRequest.
func NewGetTransactionStatusRequest() *GetTransactionStatusRequest {
	return &GetTransactionStatusRequest{}
}

// NewGetTransactionStatusResponse creates a new GetTransactionStatusResponse.
func NewGetTransactionStatusResponse(messagesInQueue bool) *GetTransactionStatusResponse {
	return &GetTransactionStatusResponse{MessagesInQueue: messagesInQueue}
}
