// Contains common and shared data types between OCPP 2.1 messages.
package types

import (
	"gopkg.in/go-playground/validator.v9"

	"github.com/lorenzodonini/ocpp-go/ocppj"
)

const (
	V2Subprotocol = "ocpp2.1"
)

type PropertyViolation struct {
	error
	Property string
}

func (e *PropertyViolation) Error() string {
	return ""
}

// Generic Device Model Status
type GenericDeviceModelStatus string

const (
	GenericDeviceModelStatusAccepted       GenericDeviceModelStatus = "Accepted"
	GenericDeviceModelStatusRejected       GenericDeviceModelStatus = "Rejected"
	GenericDeviceModelStatusNotSupported   GenericDeviceModelStatus = "NotSupported"
	GenericDeviceModelStatusEmptyResultSet GenericDeviceModelStatus = "EmptyResultSet" // If the combination of received criteria result in an empty result set.
)

func isValidGenericDeviceModelStatus(fl validator.FieldLevel) bool {
	status := GenericDeviceModelStatus(fl.Field().String())
	switch status {
	case GenericDeviceModelStatusAccepted, GenericDeviceModelStatusRejected, GenericDeviceModelStatusNotSupported, GenericDeviceModelStatusEmptyResultSet:
		return true
	default:
		return false
	}
}

// Generic Status
type GenericStatus string

const (
	GenericStatusAccepted GenericStatus = "Accepted"
	GenericStatusRejected GenericStatus = "Rejected"
)

func isValidGenericStatus(fl validator.FieldLevel) bool {
	status := GenericStatus(fl.Field().String())
	switch status {
	case GenericStatusAccepted, GenericStatusRejected:
		return true
	default:
		return false
	}
}

// ID Token Info
type MessageFormatType string

const (
	MessageFormatASCII MessageFormatType = "ASCII"
	MessageFormatHTML  MessageFormatType = "HTML"
	MessageFormatURI   MessageFormatType = "URI"
	MessageFormatUTF8  MessageFormatType = "UTF8"
)

func isValidMessageFormatType(fl validator.FieldLevel) bool {
	algorithm := MessageFormatType(fl.Field().String())
	switch algorithm {
	case MessageFormatASCII, MessageFormatHTML, MessageFormatURI, MessageFormatUTF8:
		return true
	default:
		return false
	}
}

type MessageContent struct {
	Format   MessageFormatType `json:"format" validate:"required,messageFormat21"`
	Language string            `json:"language,omitempty" validate:"max=8"`
	Content  string            `json:"content" validate:"required,max=1024"`
}

// StatusInfo is an element providing more information about the message status.
type StatusInfo struct {
	ReasonCode     string `json:"reasonCode" validate:"required,max=20"`                  // A predefined code for the reason why the status is returned in this response. The string is case- insensitive.
	AdditionalInfo string `json:"additionalInfo,omitempty" validate:"omitempty,max=1024"` // Additional text to provide detailed information.
}

// NewStatusInfo creates a StatusInfo struct.
// If no additional info need to be set, an empty string may be passed.
func NewStatusInfo(reasonCode string, additionalInfo string) *StatusInfo {
	return &StatusInfo{ReasonCode: reasonCode, AdditionalInfo: additionalInfo}
}

// EVSE represents the Electric Vehicle Supply Equipment, formerly referred to as connector(s).
type EVSE struct {
	ID          int  `json:"id" validate:"gte=0"`                              // The EVSE Identifier. When 0, the ID references the Charging Station as a whole.
	ConnectorID *int `json:"connectorId,omitempty" validate:"omitempty,gte=0"` // An id to designate a specific connector (on an EVSE) by connector index number.
}

// Component represents a physical or logical component.
type Component struct {
	Name     string `json:"name" validate:"required,max=50"`                // Name of the component. Name should be taken from the list of standardized component names whenever possible. Case Insensitive. strongly advised to use Camel Case.
	Instance string `json:"instance,omitempty" validate:"omitempty,max=50"` // Name of instance in case the component exists as multiple instances. Case Insensitive. strongly advised to use Camel Case.
	EVSE     *EVSE  `json:"evse,omitempty" validate:"omitempty"`            // Specifies the EVSE when component is located at EVSE level, also specifies the connector when component is located at Connector level.
}

// Variable is a reference key to a component-variable.
type Variable struct {
	Name     string `json:"name" validate:"required,max=50"`                // Name of the variable. Name should be taken from the list of standardized variable names whenever possible. Case Insensitive. strongly advised to use Camel Case.
	Instance string `json:"instance,omitempty" validate:"omitempty,max=50"` // Name of instance in case the variable exists as multiple instances. Case Insensitive. strongly advised to use Camel Case.
}

// ComponentVariable is used to report components, variables and variable attributes and characteristics.
type ComponentVariable struct {
	Component Component `json:"component" validate:"required"` // Component for which a report of Variable is requested.
	Variable  Variable  `json:"variable" validate:"required"`  // Variable for which report is requested.
}

// Attribute is an enumeration type used when requesting a variable value.
type Attribute string

const (
	AttributeActual Attribute = "Actual" // The actual value of the variable.
	AttributeTarget Attribute = "Target" // The target value for this variable.
	AttributeMinSet Attribute = "MinSet" // The minimal allowed value for this variable.
	AttributeMaxSet Attribute = "MaxSet" // The maximum allowed value for this variable
)

func isValidAttribute(fl validator.FieldLevel) bool {
	purposeType := Attribute(fl.Field().String())
	switch purposeType {
	case AttributeActual, AttributeTarget, AttributeMinSet, AttributeMaxSet:
		return true
	default:
		return false
	}
}

//TODO: remove SignatureMethod (obsolete from 2.0.1 onwards)

// Enumeration of the method used to encode the meter value into binary data before applying the digital signature algorithm.
// If the EncodingMethod is set to Other, the CSMS MAY try to determine the encoding method from the encodedMeterValue field.
type EncodingMethod string

const (
	EncodingOther              EncodingMethod = "Other"                // Encoding method is not included in the enumeration.
	EncodingDLMSMessage        EncodingMethod = "DLMS Message"         // The data is encoded in a digitally signed DLMS message, as described in the DLMS Green Book 8.
	EncodingCOSEMProtectedData EncodingMethod = "COSEM Protected Data" // The data is encoded according to the COSEM data protection methods, as described in the DLMS Blue Book 12.
	EncodingEDL                EncodingMethod = "EDL"                  // The data is encoded in the format used by EDL meters.
)

func isValidEncodingMethod(fl validator.FieldLevel) bool {
	encoding := EncodingMethod(fl.Field().String())
	switch encoding {
	case EncodingCOSEMProtectedData, EncodingEDL, EncodingDLMSMessage, EncodingOther:
		return true
	default:
		return false
	}
}

// Validator used for validating all OCPP 2.1 messages.
// Any additional custom validations must be added to this object for automatic validation.
var Validate = ocppj.Validate

func init() {
	_ = Validate.RegisterValidation("idTokenType21", isValidIdTokenType)
	_ = Validate.RegisterValidation("genericDeviceModelStatus21", isValidGenericDeviceModelStatus)
	_ = Validate.RegisterValidation("genericStatus21", isValidGenericStatus)
	_ = Validate.RegisterValidation("hashAlgorithm21", isValidHashAlgorithmType)
	_ = Validate.RegisterValidation("messageFormat21", isValidMessageFormatType)
	_ = Validate.RegisterValidation("authorizationStatus21", isValidAuthorizationStatus)
	_ = Validate.RegisterValidation("attribute21", isValidAttribute)
	_ = Validate.RegisterValidation("chargingProfilePurpose21", isValidChargingProfilePurpose)
	_ = Validate.RegisterValidation("chargingProfileKind21", isValidChargingProfileKind)
	_ = Validate.RegisterValidation("recurrencyKind21", isValidRecurrencyKind)
	_ = Validate.RegisterValidation("chargingRateUnit21", isValidChargingRateUnit)
	_ = Validate.RegisterValidation("chargingLimitSource", isValidChargingLimitSource)
	_ = Validate.RegisterValidation("remoteStartStopStatus21", isValidRemoteStartStopStatus)
	_ = Validate.RegisterValidation("readingContext21", isValidReadingContext)
	_ = Validate.RegisterValidation("measurand21", isValidMeasurand)
	_ = Validate.RegisterValidation("phase21", isValidPhase)
	_ = Validate.RegisterValidation("location21", isValidLocation)
	_ = Validate.RegisterValidation("signatureMethod21", isValidSignatureMethod)
	_ = Validate.RegisterValidation("encodingMethod21", isValidEncodingMethod)
	_ = Validate.RegisterValidation("certificateSigningUse21", isValidCertificateSigningUse)
	_ = Validate.RegisterValidation("certificateUse21", isValidCertificateUse)
	_ = Validate.RegisterValidation("15118EVCertificate21", isValidCertificate15118EVStatus)
	_ = Validate.RegisterValidation("costKind21", isValidCostKind)
	_ = Validate.RegisterValidation("operationMode21", isValidOperationMode)

	Validate.RegisterStructValidation(isValidIdToken, IdToken{})
	Validate.RegisterStructValidation(isValidGroupIdToken, GroupIdToken{})
}
