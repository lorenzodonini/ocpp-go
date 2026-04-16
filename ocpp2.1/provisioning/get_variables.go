package provisioning

import (
	"reflect"

	"gopkg.in/go-playground/validator.v9"

	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
)

// -------------------- Get Variables (CSMS -> CS) --------------------

const GetVariablesFeatureName = "GetVariables"

// GetVariableStatus represents the status of getting a variable.
type GetVariableStatus string

const (
	GetVariableStatusAccepted                  GetVariableStatus = "Accepted"
	GetVariableStatusRejected                  GetVariableStatus = "Rejected"
	GetVariableStatusUnknownComponent          GetVariableStatus = "UnknownComponent"
	GetVariableStatusUnknownVariable           GetVariableStatus = "UnknownVariable"
	GetVariableStatusNotSupportedAttributeType GetVariableStatus = "NotSupportedAttributeType"
)

func isValidGetVariableStatus(fl validator.FieldLevel) bool {
	s := GetVariableStatus(fl.Field().String())
	switch s {
	case GetVariableStatusAccepted, GetVariableStatusRejected, GetVariableStatusUnknownComponent,
		GetVariableStatusUnknownVariable, GetVariableStatusNotSupportedAttributeType:
		return true
	default:
		return false
	}
}

// GetVariableData contains data for getting a variable.
type GetVariableData struct {
	Component     types.Component       `json:"component" validate:"required"`
	Variable      types.Variable        `json:"variable" validate:"required"`
	AttributeType types.Attribute       `json:"attributeType,omitempty" validate:"omitempty,attribute21"`
	CustomData    *types.CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// GetVariableResult contains the result of getting a variable.
type GetVariableResult struct {
	AttributeStatus     GetVariableStatus     `json:"attributeStatus" validate:"required,getVariableStatus21"`
	AttributeType       types.Attribute       `json:"attributeType,omitempty" validate:"omitempty,attribute21"`
	AttributeValue      string                `json:"attributeValue,omitempty" validate:"omitempty,max=2500"`
	AttributeStatusInfo *types.StatusInfo     `json:"attributeStatusInfo,omitempty" validate:"omitempty"`
	Component           types.Component       `json:"component" validate:"required"`
	Variable            types.Variable        `json:"variable" validate:"required"`
	CustomData          *types.CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// GetVariablesRequest is the payload for a GetVariables request from CSMS to CS.
type GetVariablesRequest struct {
	GetVariableData []GetVariableData     `json:"getVariableData" validate:"required,min=1,dive"`
	CustomData      *types.CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// GetVariablesResponse is the payload for a GetVariables response from CS to CSMS.
type GetVariablesResponse struct {
	GetVariableResult []GetVariableResult   `json:"getVariableResult" validate:"required,min=1,dive"`
	CustomData        *types.CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// GetVariablesFeature represents the GetVariables feature.
type GetVariablesFeature struct{}

func (f GetVariablesFeature) GetFeatureName() string {
	return GetVariablesFeatureName
}

func (f GetVariablesFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(GetVariablesRequest{})
}

func (f GetVariablesFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(GetVariablesResponse{})
}

func (r GetVariablesRequest) GetFeatureName() string {
	return GetVariablesFeatureName
}

func (c GetVariablesResponse) GetFeatureName() string {
	return GetVariablesFeatureName
}

// NewGetVariablesRequest creates a new GetVariablesRequest.
func NewGetVariablesRequest(data []GetVariableData) *GetVariablesRequest {
	return &GetVariablesRequest{GetVariableData: data}
}

// NewGetVariablesResponse creates a new GetVariablesResponse.
func NewGetVariablesResponse(result []GetVariableResult) *GetVariablesResponse {
	return &GetVariablesResponse{GetVariableResult: result}
}

func init() {
	_ = types.Validate.RegisterValidation("getVariableStatus21", isValidGetVariableStatus)
}
