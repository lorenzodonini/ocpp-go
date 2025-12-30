package der

import (
	"github.com/xBlaz3kx/ocpp-go/ocpp2.1/types"
	"reflect"
)

// -------------------- GetDERControl (CSMS -> CS) --------------------

const GetDERControl = "GetDERControl"

// The field definition of the GetDERControlRequest request payload sent by the CSMS to the Charging Station.
type GetDERControlRequest struct {
	RequestId   int         `json:"requestId" validate:"required"`
	IsDefault   *bool       `json:"isDefault,omitempty" validate:"omitempty"`
	ControlType *DERControl `json:"controlType,omitempty" validate:"omitempty"`
	ControlId   string      `json:"controlId,omitempty" validate:"omitempty,max=36"` // Optional field, max length 36 characters
}

// This field definition of the GetDERControlResponse
type GetDERControlResponse struct {
	Status     DERControlStatus  `json:"status" validate:"required,derControlStatus"`
	StatusInfo *types.StatusInfo `json:"statusInfo,omitempty" validate:"omitempty,dive"`
}

type GetDERControlFeature struct{}

func (f GetDERControlFeature) GetFeatureName() string {
	return GetDERControl
}

func (f GetDERControlFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(GetDERControlRequest{})
}

func (f GetDERControlFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(GetDERControlResponse{})
}

func (r GetDERControlRequest) GetFeatureName() string {
	return GetDERControl
}

func (c GetDERControlResponse) GetFeatureName() string {
	return GetDERControl
}

// Creates a new GetDERControlRequest, containing all required fields. Optional fields may be set afterwards.
func NewGetDERControlResponseRequest(requestId int) *GetDERControlRequest {
	return &GetDERControlRequest{
		RequestId: requestId,
	}
}

// Creates a new GetDERControlResponse, containing all required fields. Optional fields may be set afterwards.
func NewGetDERControlResponseResponse(status DERControlStatus) *GetDERControlResponse {
	return &GetDERControlResponse{
		Status: status,
	}
}
