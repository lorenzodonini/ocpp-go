package battey_swap

import (
	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
	"reflect"
)

// -------------------- RequestBatterySwap (CSMS -> CS) --------------------

const RequestBatterySwap = "RequestBatterySwap"

// The field definition of the RequestBatterySwapRequest request payload sent by the CSMS to the Charging Station.
type RequestBatterySwapRequest struct {
	RequestId int               `json:"requestId" validate:"required"`
	IdToken   types.IdTokenType `json:"idToken" validate:"required,dive"`
}

// This field definition of the RequestBatterySwapResponse
type RequestBatterySwapResponse struct {
	Status     types.GenericStatus `json:"status" validate:"required,genericStatus21"`
	StatusInfo *types.StatusInfo   `json:"statusInfo,omitempty" validate:"omitempty,dive"`
}

type RequestBatterySwapFeature struct{}

func (f RequestBatterySwapFeature) GetFeatureName() string {
	return RequestBatterySwap
}

func (f RequestBatterySwapFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(RequestBatterySwapRequest{})
}

func (f RequestBatterySwapFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(RequestBatterySwapResponse{})
}

func (r RequestBatterySwapRequest) GetFeatureName() string {
	return RequestBatterySwap
}

func (c RequestBatterySwapResponse) GetFeatureName() string {
	return RequestBatterySwap
}

// Creates a new RequestBatterySwapRequest, containing all required fields. Optional fields may be set afterwards.
func NewRequestBatterySwapRequest(requestId int, idToken types.IdTokenType) *RequestBatterySwapRequest {
	return &RequestBatterySwapRequest{
		RequestId: requestId,
		IdToken:   idToken,
	}
}

// Creates a new RequestBatterySwapResponse, containing all required fields. Optional fields may be set afterwards.
func NewRequestBatterySwapResponse(status types.GenericStatus, statusInfo *types.StatusInfo) *RequestBatterySwapResponse {
	return &RequestBatterySwapResponse{
		Status:     status,
		StatusInfo: statusInfo,
	}
}
