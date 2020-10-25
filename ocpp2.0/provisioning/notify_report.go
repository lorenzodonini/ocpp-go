package provisioning

import (
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Notify Report (CS -> CSMS) --------------------

const NotifyReportFeatureName = "NotifyReport"

// Mutability defines the mutability of an attribute.
type Mutability string

const (
	MutabilityReadOnly  Mutability = "ReadOnly"
	MutabilityWriteOnly Mutability = "WriteOnly"
	MutabilityReadWrite Mutability = "ReadWrite"
)

func isValidMutability(fl validator.FieldLevel) bool {
	status := Mutability(fl.Field().String())
	switch status {
	case MutabilityReadOnly, MutabilityWriteOnly, MutabilityReadWrite:
		return true
	default:
		return false
	}
}

// DataType defines the data type of a variable.
type DataType string

const (
	TypeString       DataType = "string"
	TypeDecimal      DataType = "decimal"
	TypeInteger      DataType = "integer"
	TypeDateTime     DataType = "dateTime"
	TypeBoolean      DataType = "boolean"
	TypeOptionList   DataType = "OptionList"
	TypeSequenceList DataType = "SequenceList"
	TypeMemberList   DataType = "MemberList"
)

func isValidDataType(fl validator.FieldLevel) bool {
	status := DataType(fl.Field().String())
	switch status {
	case TypeBoolean, TypeDateTime, TypeDecimal, TypeInteger, TypeString, TypeOptionList, TypeSequenceList, TypeMemberList:
		return true
	default:
		return false
	}
}

// VariableCharacteristics represents a fixed read-only parameters of a variable.
type VariableCharacteristics struct {
	Unit               string   `json:"unit,omitempty" validate:"max=16"`          // Unit of the variable. When the transmitted value has a unit, this field SHALL be included.
	DataType           DataType `json:"dataType" validate:"required,dataTypeEnum"` // Data type of this variable.
	MinLimit           *float64 `json:"minLimit,omitempty"`                        // Minimum possible value of this variable.
	MaxLimit           *float64 `json:"maxLimit,omitempty"`                        // Maximum possible value of this variable. When the datatype of this Variable is String, OptionList, SequenceList or MemberList, this field defines the maximum length of the (CSV) string.
	ValuesList         string   `json:"valuesList,omitempty" validate:"max=1000"`  // Allowed values when variable is Option/Member/SequenceList. This is a comma separated list.
	SupportsMonitoring bool     `json:"supportsMonitoring"`                        // Flag indicating if this variable supports monitoring.
}

// NewVariableCharacteristics returns a pointer to a new VariableCharacteristics struct.
func NewVariableCharacteristics(dataType DataType, supportsMonitoring bool) *VariableCharacteristics {
	return &VariableCharacteristics{DataType: dataType, SupportsMonitoring: supportsMonitoring}
}

// VariableAttribute describes the attribute data of a variable.
type VariableAttribute struct {
	Type       types.Attribute `json:"type,omitempty" validate:"omitempty,attribute"`        // Actual, MinSet, MaxSet, etc. Defaults to Actual if absent.
	Value      string          `json:"value,omitempty" validate:"max=2500"`                  // Value of the attribute. May only be omitted when mutability is set to 'WriteOnly'.
	Mutability Mutability      `json:"mutability,omitempty" validate:"omitempty,mutability"` // Defines the mutability of this attribute. Default is ReadWrite when omitted.
	Persistent bool            `json:"persistent,omitempty"`                                 // If true, value will be persistent across system reboots or power down. Default when omitted is false.
	Constant   bool            `json:"constant,omitempty"`                                   // If true, value that will never be changed by the Charging Station at runtime. Default when omitted is false.
}

// NewVariableAttribute creates a VariableAttribute struct, with all default values set.
func NewVariableAttribute() VariableAttribute {
	return VariableAttribute{
		Type:       types.AttributeActual,
		Mutability: MutabilityReadWrite,
	}
}

// ReportData is a struct to report components, variables and variable attributes and characteristics.
type ReportData struct {
	Component               types.Component          `json:"component" validate:"required"`
	Variable                types.Variable           `json:"variable" validate:"required"`
	VariableAttribute       []VariableAttribute      `json:"variableAttribute" validate:"required,min=1,max=4,dive"`
	VariableCharacteristics *VariableCharacteristics `json:"variableCharacteristics,omitempty" validate:"omitempty"`
}

// The field definition of the NotifyReport request payload sent by the Charging Station to the CSMS.
type NotifyReportRequest struct {
	RequestID   int             `json:"requestId" validate:"gte=0"`                     // The id of the GetMonitoringRequest that requested this report.
	GeneratedAt *types.DateTime `json:"generatedAt" validate:"required"`                // Timestamp of the moment this message was generated at the Charging Station.
	Tbc         bool            `json:"tbc,omitempty" validate:"omitempty"`             // “to be continued” indicator. Indicates whether another part of the monitoringData follows in an upcoming notifyMonitoringReportRequest message. Default value when omitted is false.
	SeqNo       int             `json:"seqNo" validate:"gte=0"`                         // Sequence number of this message. First message starts at 0.
	ReportData  []ReportData    `json:"reportData,omitempty" validate:"omitempty,dive"` // List of ReportData
}

// The field definition of the NotifyReport response payload, sent by the CSMS to the Charging Station in response to a NotifyReportRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type NotifyReportResponse struct {
}

// A Charging Station may send reports to the CSMS on demand, when requested to do so.
// After receiving a GetBaseReport from the CSMS, a Charging Station asynchronously sends the results
// in one or more NotifyReportRequest messages.
//
// The CSMS responds with NotifyReportResponse for each received request.
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

// Creates a new NotifyReportRequest, containing all required fields. Optional fields may be set afterwards.
func NewNotifyReportRequest(requestID int, generatedAt *types.DateTime, seqNo int) *NotifyReportRequest {
	return &NotifyReportRequest{RequestID: requestID, GeneratedAt: generatedAt, SeqNo: seqNo}
}

// Creates a new NotifyReportResponse. There are no optional fields for this message.
func NewNotifyReportResponse() *NotifyReportResponse {
	return &NotifyReportResponse{}
}

func init() {
	_ = types.Validate.RegisterValidation("mutability", isValidMutability)
	_ = types.Validate.RegisterValidation("dataTypeEnum", isValidDataType)
}
