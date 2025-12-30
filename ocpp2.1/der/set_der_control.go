package der

import (
	"github.com/xBlaz3kx/ocpp-go/ocpp2.1/types"
	"reflect"
)

// -------------------- SetDERControl (CSMS -> CS) --------------------

const SetDERControl = "SetDERControl"

// The field definition of the SetDERControlRequest request payload sent by the CSMS to the Charging Station.
type SetDERControlRequest struct {
	IsDefault         bool               `json:"isDefault" validate:"required"` // Indicates whether the DER control is set to default values.
	ControlID         string             `json:"controlId" validate:"required"` // The unique identifier of the DER control to be set.
	ControlType       DERControl         `json:"controlType" validate:"required,derControl"`
	Curve             *DERCurve          `json:"curve,omitempty" validate:"omitempty,dive"`
	Gradient          *Gradient          `json:"gradient,omitempty" validate:"omitempty,dive"`
	FreqDroop         *FreqDroop         `json:"freqDroop,omitempty" validate:"omitempty,dive"`
	FixedPFAbsorb     *FixedPF           `json:"fixedPFAbsorb,omitempty" validate:"omitempty,dive"`
	FixedPFInject     *FixedPF           `json:"fixedPFInject,omitempty" validate:"omitempty,dive"`
	LimitMaxDischarge *LimitMaxDischarge `json:"limitMaxDischarge,omitempty" validate:"omitempty,dive"`
	EnterService      *EnterService      `json:"enterService,omitempty" validate:"omitempty,dive"`
}

// This field definition of the SetDERControlResponse
type SetDERControlResponse struct {
	Status         DERControlStatus `json:"status" validate:"required,derControlStatus"`
	SuperseededIds []string         `json:"superseededIds,omitempty" validate:"omitempty,max=24"`
	StatusInfo     types.StatusInfo `json:"statusInfo,omitempty" validate:"omitempty"`
}

type SetDERControlFeature struct{}

func (f SetDERControlFeature) GetFeatureName() string {
	return SetDERControl
}

func (f SetDERControlFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(SetDERControlRequest{})
}

func (f SetDERControlFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(SetDERControlResponse{})
}

func (r SetDERControlRequest) GetFeatureName() string {
	return SetDERControl
}

func (c SetDERControlResponse) GetFeatureName() string {
	return SetDERControl
}

// Creates a new SetDERControlRequest, containing all required fields. Optional fields may be set afterwards.
func NewSetDERControlResponseRequest(isDefault bool, controlId string, controlType DERControl) *SetDERControlRequest {
	return &SetDERControlRequest{
		IsDefault:   isDefault,
		ControlID:   controlId,
		ControlType: controlType,
	}
}

// Creates a new SetDERControlResponse, containing all required fields. Optional fields may be set afterwards.
func NewSetDERControlResponseResponse(status DERControlStatus) *SetDERControlResponse {
	return &SetDERControlResponse{
		Status: status,
	}
}
