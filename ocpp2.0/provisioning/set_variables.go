package provisioning

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
	"gopkg.in/go-playground/validator.v9"
)

// -------------------- Get Variable (CSMS -> CS) --------------------

const SetVariablesFeatureName = "SetVariables"

// SetVariableStatus indicates the result status of setting a variable in SetVariablesResponse.
type SetVariableStatus string

const (
	SetVariableStatusAccepted         SetVariableStatus = "Accepted"
	SetVariableStatusRejected         SetVariableStatus = "Rejected"
	SetVariableStatusUnknownComponent SetVariableStatus = "UnknownComponent"
	SetVariableStatusUnknownVariable  SetVariableStatus = "UnknownVariable"
	SetVariableStatusNotSupported     SetVariableStatus = "NotSupportedAttributeType"
	SetVariableStatusRebootRequired   SetVariableStatus = "RebootRequired"
)

func isValidSetVariableStatus(fl validator.FieldLevel) bool {
	status := SetVariableStatus(fl.Field().String())
	switch status {
	case SetVariableStatusAccepted, SetVariableStatusRejected, SetVariableStatusUnknownComponent, SetVariableStatusUnknownVariable, SetVariableStatusNotSupported, SetVariableStatusRebootRequired:
		return true
	default:
		return false
	}
}

type SetVariableData struct {
	AttributeType  types.Attribute `json:"attributeType,omitempty" validate:"omitempty,attribute"`
	AttributeValue string          `json:"attributeValue" validate:"required,max=1000"`
	Component      types.Component `json:"component" validate:"required"`
	Variable       types.Variable  `json:"variable" validate:"required"`
}

type SetVariableResult struct {
	AttributeType   types.Attribute   `json:"attributeType,omitempty" validate:"omitempty,attribute"`
	AttributeStatus SetVariableStatus `json:"attributeStatus" validate:"required,getVariableStatus"`
	Component       types.Component   `json:"component" validate:"required"`
	Variable        types.Variable    `json:"variable" validate:"required"`
	StatusInfo      *types.StatusInfo `json:"statusInfo,omitempty" validate:"omitempty"`
}

// The field definition of the SetVariables request payload sent by the CSMS to the Charging Station.
type SetVariablesRequest struct {
	SetVariableData []SetVariableData `json:"setVariableData" validate:"required,min=1,dive"` // List of Component-Variable pairs and attribute values to set.
}

// This field definition of the SetVariables response payload, sent by the Charging Station to the CSMS in response to a SetVariablesRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type SetVariablesResponse struct {
	SetVariableResult []SetVariableResult `json:"setVariableResult" validate:"required,min=1,dive"` //  List of result statuses per Component-Variable.
}

// A Charging Station can have a lot of variables that can be configured/changed by the CSMS.
//
// The CSO may trigger the CSMS to request setting one or more variables in a Charging Station.
// The CSMS sends a SetVariablesRequest to the Charging Station, to configured/change one or more variables.
// The Charging Station responds with a SetVariablesResponse indicating whether it was able to executed the change(s).
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

// Creates a new SetVariablesRequest, containing all required fields.  There are no optional fields for this message.
func NewSetVariablesRequest(variableData []SetVariableData) *SetVariablesRequest {
	return &SetVariablesRequest{variableData}
}

// Creates a new SetVariablesResponse, containing all required fields. There are no optional fields for this message.
func NewSetVariablesResponse(result []SetVariableResult) *SetVariablesResponse {
	return &SetVariablesResponse{result}
}

func init() {
	_ = types.Validate.RegisterValidation("setVariableStatus", isValidSetVariableStatus)
}
