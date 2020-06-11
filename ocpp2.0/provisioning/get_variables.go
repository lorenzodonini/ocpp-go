package provisioning

import (
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Get Variable (CSMS -> CS) --------------------

const GetVariablesFeatureName = "GetVariables"

// GetVariableStatus indicates the result status of getting a variable in GetVariablesResponse.
type GetVariableStatus string

const (
	GetVariableStatusAccepted         GetVariableStatus = "Accepted"
	GetVariableStatusRejected         GetVariableStatus = "Rejected"
	GetVariableStatusUnknownComponent GetVariableStatus = "UnknownComponent"
	GetVariableStatusUnknownVariable  GetVariableStatus = "UnknownVariable"
	GetVariableStatusNotSupported     GetVariableStatus = "NotSupportedAttributeType"
)

func isValidGetVariableStatus(fl validator.FieldLevel) bool {
	status := GetVariableStatus(fl.Field().String())
	switch status {
	case GetVariableStatusAccepted, GetVariableStatusRejected, GetVariableStatusUnknownComponent, GetVariableStatusUnknownVariable, GetVariableStatusNotSupported:
		return true
	default:
		return false
	}
}

type VariableData struct {
	AttributeType types.Attribute `json:"attributeType,omitempty" validate:"omitempty,attribute"`
	Component     types.Component `json:"component" validate:"required"`
	Variable      types.Variable  `json:"variable" validate:"required"`
}

type VariableResult struct {
	AttributeStatus GetVariableStatus `json:"attributeStatus" validate:"required,getVariableStatus"`
	AttributeType   types.Attribute   `json:"attributeType,omitempty" validate:"omitempty,attribute"`
	AttributeValue  string            `json:"attributeValue,omitempty" validate:"omitempty,max=1000"`
	Component       types.Component   `json:"component" validate:"required"`
	Variable        types.Variable    `json:"variable" validate:"required"`
}

// The field definition of the GetVariables request payload sent by the CSMS to the Charging Station.
type GetVariablesRequest struct {
	GetVariableData   []VariableData            `json:"getVariableData" validate:"required,min=1,dive"`
}

// This field definition of the GetVariables response payload, sent by the Charging Station to the CSMS in response to a GetVariablesRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type GetVariablesResponse struct {
	GetVariableResult []VariableResult               `json:"getVariableResult" validate:"required,min=1,dive"`
}

// The CSO may trigger the CSMS to request to request for a number of variables in a Charging Station.
// The CSMS request the Charging Station for a number of variables (of one or more components) with GetVariablesRequest with a list of requested variables.
// The Charging Station responds with a GetVariablesResponse with the requested variables.
//
// It is not possible to get all attributes of all variables in one call.
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

// Creates a new GetVariablesRequest, containing all required fields.  There are no optional fields for this message.
func NewGetVariablesRequest(variableData []VariableData) *GetVariablesRequest {
	return &GetVariablesRequest{variableData}
}

// Creates a new GetVariablesResponse, containing all required fields. There are no optional fields for this message.
func NewGetVariablesResponse(result []VariableResult) *GetVariablesResponse {
	return &GetVariablesResponse{result}
}

func init() {
	_ = types.Validate.RegisterValidation("getVariableStatus", isValidGetVariableStatus)
}
