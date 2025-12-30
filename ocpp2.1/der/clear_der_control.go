package der

import (
	"github.com/xBlaz3kx/ocpp-go/ocpp2.1/types"
	"reflect"
)

// -------------------- ClearDERControl (CSMS -> CS) --------------------

const ClearDERControl = "ClearDERControl"

// The field definition of the ClearDERControlRequest request payload sent by the CSMS to the Charging Station.
type ClearDERControlRequest struct {
	IsDefault   bool        `json:"isDefault" validate:"required"`
	ControlType *DERControl `json:"controlType,omitempty" validate:"omitempty"`
	ControlId   string      `json:"controlId,omitempty" validate:"omitempty,max=36"`
}

// This field definition of the ClearDERControlResponse
type ClearDERControlResponse struct {
	Status     DERControlStatus  `json:"status" validate:"required,derControlStatus"`
	StatusInfo *types.StatusInfo `json:"statusInfo,omitempty" validate:"omitempty"`
}

type ClearDERControlFeature struct{}

func (f ClearDERControlFeature) GetFeatureName() string {
	return ClearDERControl
}

func (f ClearDERControlFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(ClearDERControlRequest{})
}

func (f ClearDERControlFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(ClearDERControlResponse{})
}

func (r ClearDERControlRequest) GetFeatureName() string {
	return ClearDERControl
}

func (c ClearDERControlResponse) GetFeatureName() string {
	return ClearDERControl
}

// Creates a new ClearDERControlRequest, containing all required fields. Optional fields may be set afterwards.
func NewClearDERControlRequest(isDefault bool) *ClearDERControlRequest {
	return &ClearDERControlRequest{
		IsDefault: isDefault,
	}
}

// Creates a new ClearDERControlResponse, containing all required fields. Optional fields may be set afterwards.
func NewClearDERControlResponse(status DERControlStatus) *ClearDERControlResponse {
	return &ClearDERControlResponse{
		Status: status,
	}
}
