package remotecontrol

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

// -------------------- Request Start Transaction (CSMS -> CS) --------------------

const RequestStopTransactionFeatureName = "RequestStopTransaction"

// The field definition of the RequestStopTransaction request payload sent by the CSMS to the Charging Station.
type RequestStopTransactionRequest struct {
	TransactionID string `json:"transactionId" validate:"required,max=36"`
}

// This field definition of the RequestStopTransaction response payload, sent by the Charging Station to the CSMS in response to a RequestStopTransactionRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type RequestStopTransactionResponse struct {
	Status     RequestStartStopStatus `json:"status" validate:"required,requestStartStopStatus"`
	StatusInfo *types.StatusInfo      `json:"statusInfo,omitempty"`
}

// The CSMS may remotely stop an ongoing transaction for a user.
// This functionality may be triggered by:
//   - a CSO, to help out a user, that is having trouble stopping a transaction
//   - a third-party event (e.g. mobile app)
//   - the ISO15118-1 use-case F2
//
// The CSMS sends a RequestStopTransactionRequest to the Charging Station.
// The Charging Stations will reply with a RequestStopTransactionResponse.
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

// Creates a new RequestStopTransactionRequest, containing all required fields. There are no optional fields for this message.
func NewRequestStopTransactionRequest(transactionID string) *RequestStopTransactionRequest {
	return &RequestStopTransactionRequest{TransactionID: transactionID}
}

// Creates a new RequestStopTransactionResponse, containing all required fields. Optional fields may be set afterwards.
func NewRequestStopTransactionResponse(status RequestStartStopStatus) *RequestStopTransactionResponse {
	return &RequestStopTransactionResponse{Status: status}
}
