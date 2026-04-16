package provisioning

import (
	"reflect"

	"gopkg.in/go-playground/validator.v9"

	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
)

// -------------------- Set Variables (CSMS -> CS) --------------------

const SetVariablesFeatureName = "SetVariables"

// SetVariableStatus represents the status of setting a variable.
type SetVariableStatus string

const (
	SetVariableStatusAccepted                  SetVariableStatus = "Accepted"
	SetVariableStatusRejected                  SetVariableStatus = "Rejected"
	SetVariableStatusUnknownComponent          SetVariableStatus = "UnknownComponent"
	SetVariableStatusUnknownVariable           SetVariableStatus = "UnknownVariable"
	SetVariableStatusNotSupportedAttributeType SetVariableStatus = "NotSupportedAttributeType"
	SetVariableStatusRebootRequired            SetVariableStatus = "RebootRequired"
)

func isValidSetVariableStatus(fl validator.FieldLevel) bool {
	s := SetVariableStatus(fl.Field().String())
	switch s {
	case SetVariableStatusAccepted, SetVariableStatusRejected, SetVariableStatusUnknownComponent,
		SetVariableStatusUnknownVariable, SetVariableStatusNotSupportedAttributeType, SetVariableStatusRebootRequired:
		return true
	default:
		return false
	}
}

// SetVariableData contains data for setting a variable.
type SetVariableData struct {
	AttributeType  types.Attribute       `json:"attributeType,omitempty" validate:"omitempty,attribute21"`
	AttributeValue string                `json:"attributeValue" validate:"required,max=1000"`
	Component      types.Component       `json:"component" validate:"required"`
	Variable       types.Variable        `json:"variable" validate:"required"`
	CustomData     *types.CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// SetVariableResult contains the result of setting a variable.
type SetVariableResult struct {
	AttributeType       types.Attribute       `json:"attributeType,omitempty" validate:"omitempty,attribute21"`
	AttributeStatus     SetVariableStatus     `json:"attributeStatus" validate:"required,setVariableStatus21"`
	AttributeStatusInfo *types.StatusInfo     `json:"attributeStatusInfo,omitempty" validate:"omitempty"`
	Component           types.Component       `json:"component" validate:"required"`
	Variable            types.Variable        `json:"variable" validate:"required"`
	CustomData          *types.CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// SetVariablesRequest is the payload for a SetVariables request from CSMS to CS.
type SetVariablesRequest struct {
	SetVariableData []SetVariableData     `json:"setVariableData" validate:"required,min=1,dive"`
	CustomData      *types.CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// SetVariablesResponse is the payload for a SetVariables response from CS to CSMS.
type SetVariablesResponse struct {
	SetVariableResult []SetVariableResult   `json:"setVariableResult" validate:"required,min=1,dive"`
	CustomData        *types.CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// SetVariablesFeature represents the SetVariables feature.
type SetVariablesFeature struct{}

func (f SetVariablesFeature) GetFeatureName() string {
	return SetVariablesFeatureName
}

func (f SetVariablesFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(SetVariablesRequest{})
}

func (f SetVariablesFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(SetVariablesResponse{})
}

func (r SetVariablesRequest) GetFeatureName() string {
	return SetVariablesFeatureName
}

func (c SetVariablesResponse) GetFeatureName() string {
	return SetVariablesFeatureName
}

// NewSetVariablesRequest creates a new SetVariablesRequest.
func NewSetVariablesRequest(data []SetVariableData) *SetVariablesRequest {
	return &SetVariablesRequest{SetVariableData: data}
}

// NewSetVariablesResponse creates a new SetVariablesResponse.
func NewSetVariablesResponse(result []SetVariableResult) *SetVariablesResponse {
	return &SetVariablesResponse{SetVariableResult: result}
}

func init() {
	_ = types.Validate.RegisterValidation("setVariableStatus21", isValidSetVariableStatus)
}
