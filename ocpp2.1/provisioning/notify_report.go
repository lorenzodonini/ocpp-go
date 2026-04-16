package provisioning

import (
	"reflect"

	"gopkg.in/go-playground/validator.v9"

	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
)

// -------------------- Notify Report (CS -> CSMS) --------------------

const NotifyReportFeatureName = "NotifyReport"

// MutabilityType indicates the mutability of a variable.
type MutabilityType string

const (
	MutabilityReadOnly  MutabilityType = "ReadOnly"
	MutabilityWriteOnly MutabilityType = "WriteOnly"
	MutabilityReadWrite MutabilityType = "ReadWrite"
)

func isValidMutabilityType(fl validator.FieldLevel) bool {
	m := MutabilityType(fl.Field().String())
	switch m {
	case MutabilityReadOnly, MutabilityWriteOnly, MutabilityReadWrite:
		return true
	default:
		return false
	}
}

// DataType indicates the data type of a variable.
type DataType string

const (
	DataTypeString       DataType = "string"
	DataTypeDecimal      DataType = "decimal"
	DataTypeInteger      DataType = "integer"
	DataTypeDateTime     DataType = "dateTime"
	DataTypeBoolean      DataType = "boolean"
	DataTypeOptionList   DataType = "OptionList"
	DataTypeSequenceList DataType = "SequenceList"
	DataTypeMemberList   DataType = "MemberList"
)

func isValidDataType(fl validator.FieldLevel) bool {
	d := DataType(fl.Field().String())
	switch d {
	case DataTypeString, DataTypeDecimal, DataTypeInteger, DataTypeDateTime,
		DataTypeBoolean, DataTypeOptionList, DataTypeSequenceList, DataTypeMemberList:
		return true
	default:
		return false
	}
}

// VariableCharacteristics contains characteristics of a variable.
type VariableCharacteristics struct {
	DataType           DataType              `json:"dataType" validate:"required,dataType21"`
	Unit               string                `json:"unit,omitempty" validate:"omitempty,max=16"`
	MinLimit           *float64              `json:"minLimit,omitempty"`
	MaxLimit           *float64              `json:"maxLimit,omitempty"`
	MaxElements        *int                  `json:"maxElements,omitempty" validate:"omitempty,gte=1"`
	ValuesList         string                `json:"valuesList,omitempty" validate:"omitempty,max=1000"`
	SupportsMonitoring bool                  `json:"supportsMonitoring"`
	CustomData         *types.CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// VariableAttribute contains attributes of a variable.
type VariableAttribute struct {
	Type       types.Attribute       `json:"type,omitempty" validate:"omitempty,attribute21"`
	Value      string                `json:"value,omitempty" validate:"omitempty,max=2500"`
	Mutability MutabilityType        `json:"mutability,omitempty" validate:"omitempty,mutabilityType21"`
	Persistent bool                  `json:"persistent,omitempty"`
	Constant   bool                  `json:"constant,omitempty"`
	CustomData *types.CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// ReportData contains data for a single component-variable report.
type ReportData struct {
	Component               types.Component          `json:"component" validate:"required"`
	Variable                types.Variable           `json:"variable" validate:"required"`
	VariableAttribute       []VariableAttribute      `json:"variableAttribute" validate:"required,min=1,max=4,dive"`
	VariableCharacteristics *VariableCharacteristics `json:"variableCharacteristics,omitempty" validate:"omitempty"`
	CustomData              *types.CustomDataType    `json:"customData,omitempty" validate:"omitempty"`
}

// NotifyReportRequest is the payload for a NotifyReport request from CS to CSMS.
type NotifyReportRequest struct {
	RequestID   int                   `json:"requestId" validate:"gte=0"`
	GeneratedAt *types.DateTime       `json:"generatedAt" validate:"required"`
	SeqNo       int                   `json:"seqNo" validate:"gte=0"`
	Tbc         bool                  `json:"tbc,omitempty"`
	ReportData  []ReportData          `json:"reportData,omitempty" validate:"omitempty,dive"`
	CustomData  *types.CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// NotifyReportResponse is the payload for a NotifyReport response from CSMS to CS.
type NotifyReportResponse struct {
	CustomData *types.CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// NotifyReportFeature represents the NotifyReport feature.
type NotifyReportFeature struct{}

func (f NotifyReportFeature) GetFeatureName() string {
	return NotifyReportFeatureName
}

func (f NotifyReportFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(NotifyReportRequest{})
}

func (f NotifyReportFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(NotifyReportResponse{})
}

func (r NotifyReportRequest) GetFeatureName() string {
	return NotifyReportFeatureName
}

func (c NotifyReportResponse) GetFeatureName() string {
	return NotifyReportFeatureName
}

// NewNotifyReportRequest creates a new NotifyReportRequest.
func NewNotifyReportRequest(requestID int, generatedAt *types.DateTime, seqNo int) *NotifyReportRequest {
	return &NotifyReportRequest{RequestID: requestID, GeneratedAt: generatedAt, SeqNo: seqNo}
}

// NewNotifyReportResponse creates a new NotifyReportResponse.
func NewNotifyReportResponse() *NotifyReportResponse {
	return &NotifyReportResponse{}
}

func init() {
	_ = types.Validate.RegisterValidation("mutabilityType21", isValidMutabilityType)
	_ = types.Validate.RegisterValidation("dataType21", isValidDataType)
}
