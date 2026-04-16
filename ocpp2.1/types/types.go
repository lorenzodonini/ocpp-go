// Contains common and shared data types between OCPP 2.1 messages.
package types

import (
	"gopkg.in/go-playground/validator.v9"

	"github.com/lorenzodonini/ocpp-go/ocppj"
)

const (
	V21Subprotocol = "ocpp2.1"
)

// PropertyViolation is returned in case of a validation error.
type PropertyViolation struct {
	error
	Property string
}

func (e *PropertyViolation) Error() string {
	return ""
}

// CustomDataType allows for vendor-specific extensions to OCPP messages.
// This class does not get 'AdditionalProperties = false' in the schema generation,
// so it can be extended with arbitrary JSON properties to allow adding custom data.
type CustomDataType struct {
	VendorId string `json:"vendorId" validate:"required,max=255"`
}

// NewCustomData creates a new CustomDataType with the given vendor ID.
func NewCustomData(vendorId string) *CustomDataType {
	return &CustomDataType{VendorId: vendorId}
}

// AuthorizationStatus represents the result of an authorization request.
type AuthorizationStatus string

const (
	AuthorizationStatusAccepted           AuthorizationStatus = "Accepted"
	AuthorizationStatusBlocked            AuthorizationStatus = "Blocked"
	AuthorizationStatusExpired            AuthorizationStatus = "Expired"
	AuthorizationStatusInvalid            AuthorizationStatus = "Invalid"
	AuthorizationStatusConcurrentTx       AuthorizationStatus = "ConcurrentTx"
	AuthorizationStatusNoCredit           AuthorizationStatus = "NoCredit"
	AuthorizationStatusNotAllowedTypeEVSE AuthorizationStatus = "NotAllowedTypeEVSE"
	AuthorizationStatusNotAtThisLocation  AuthorizationStatus = "NotAtThisLocation"
	AuthorizationStatusNotAtThisTime      AuthorizationStatus = "NotAtThisTime"
	AuthorizationStatusUnknown            AuthorizationStatus = "Unknown"
)

func isValidAuthorizationStatus(fl validator.FieldLevel) bool {
	status := AuthorizationStatus(fl.Field().String())
	switch status {
	case AuthorizationStatusAccepted, AuthorizationStatusBlocked, AuthorizationStatusExpired,
		AuthorizationStatusInvalid, AuthorizationStatusConcurrentTx, AuthorizationStatusNoCredit,
		AuthorizationStatusNotAllowedTypeEVSE, AuthorizationStatusNotAtThisLocation,
		AuthorizationStatusNotAtThisTime, AuthorizationStatusUnknown:
		return true
	default:
		return false
	}
}

// IdTokenType represents the type of an identification token.
// In OCPP 2.1, this is a string field with max length 20, allowing for extensibility.
type IdTokenType string

const (
	IdTokenTypeCentral         IdTokenType = "Central"
	IdTokenTypeEMAID           IdTokenType = "eMAID"
	IdTokenTypeISO14443        IdTokenType = "ISO14443"
	IdTokenTypeISO15693        IdTokenType = "ISO15693"
	IdTokenTypeKeyCode         IdTokenType = "KeyCode"
	IdTokenTypeLocal           IdTokenType = "Local"
	IdTokenTypeMacAddress      IdTokenType = "MacAddress"
	IdTokenTypeNoAuthorization IdTokenType = "NoAuthorization"
)

func isValidIdTokenType(fl validator.FieldLevel) bool {
	tokenType := IdTokenType(fl.Field().String())
	switch tokenType {
	case IdTokenTypeCentral, IdTokenTypeEMAID, IdTokenTypeISO14443, IdTokenTypeISO15693,
		IdTokenTypeKeyCode, IdTokenTypeLocal, IdTokenTypeMacAddress, IdTokenTypeNoAuthorization:
		return true
	default:
		// In 2.1, IdToken type is extensible (max=20), so we allow other values
		return len(fl.Field().String()) <= 20
	}
}

func isValidIdToken(sl validator.StructLevel) {
	idToken := sl.Current().Interface().(IdToken)
	// validate required idToken value except `NoAuthorization` type
	switch idToken.Type {
	case IdTokenTypeCentral, IdTokenTypeEMAID, IdTokenTypeISO14443, IdTokenTypeISO15693,
		IdTokenTypeKeyCode, IdTokenTypeLocal, IdTokenTypeMacAddress:
		if idToken.IdToken == "" {
			sl.ReportError(idToken.IdToken, "IdToken", "IdToken", "required", "")
		}
	}
}

// AdditionalInfo contains additional information for authorization.
type AdditionalInfo struct {
	AdditionalIdToken string          `json:"additionalIdToken" validate:"required,max=255"`
	Type              string          `json:"type" validate:"required,max=50"`
	CustomData        *CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// IdToken contains a case insensitive identifier to use for authorization.
type IdToken struct {
	IdToken        string           `json:"idToken" validate:"max=255"`
	Type           IdTokenType      `json:"type" validate:"required,max=20"`
	AdditionalInfo []AdditionalInfo `json:"additionalInfo,omitempty" validate:"omitempty,dive"`
	CustomData     *CustomDataType  `json:"customData,omitempty" validate:"omitempty"`
}

// GenericDeviceModelStatus represents a generic status for device model operations.
type GenericDeviceModelStatus string

const (
	GenericDeviceModelStatusAccepted       GenericDeviceModelStatus = "Accepted"
	GenericDeviceModelStatusRejected       GenericDeviceModelStatus = "Rejected"
	GenericDeviceModelStatusNotSupported   GenericDeviceModelStatus = "NotSupported"
	GenericDeviceModelStatusEmptyResultSet GenericDeviceModelStatus = "EmptyResultSet"
)

func isValidGenericDeviceModelStatus(fl validator.FieldLevel) bool {
	status := GenericDeviceModelStatus(fl.Field().String())
	switch status {
	case GenericDeviceModelStatusAccepted, GenericDeviceModelStatusRejected,
		GenericDeviceModelStatusNotSupported, GenericDeviceModelStatusEmptyResultSet:
		return true
	default:
		return false
	}
}

// GenericStatus represents a generic status response.
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

// HashAlgorithmType represents supported hash algorithms.
type HashAlgorithmType string

const (
	SHA256 HashAlgorithmType = "SHA256"
	SHA384 HashAlgorithmType = "SHA384"
	SHA512 HashAlgorithmType = "SHA512"
)

func isValidHashAlgorithmType(fl validator.FieldLevel) bool {
	algorithm := HashAlgorithmType(fl.Field().String())
	switch algorithm {
	case SHA256, SHA384, SHA512:
		return true
	default:
		return false
	}
}

// OCSPRequestDataType contains data needed to request an OCSP certificate status.
type OCSPRequestDataType struct {
	HashAlgorithm  HashAlgorithmType `json:"hashAlgorithm" validate:"required,hashAlgorithm21"`
	IssuerNameHash string            `json:"issuerNameHash" validate:"required,max=128"`
	IssuerKeyHash  string            `json:"issuerKeyHash" validate:"required,max=128"`
	SerialNumber   string            `json:"serialNumber" validate:"required,max=40"`
	ResponderURL   string            `json:"responderURL" validate:"required,max=2000"`
	CustomData     *CustomDataType   `json:"customData,omitempty" validate:"omitempty"`
}

// CertificateHashData contains hash data of a certificate.
type CertificateHashData struct {
	HashAlgorithm  HashAlgorithmType `json:"hashAlgorithm" validate:"required,hashAlgorithm21"`
	IssuerNameHash string            `json:"issuerNameHash" validate:"required,max=128"`
	IssuerKeyHash  string            `json:"issuerKeyHash" validate:"required,max=128"`
	SerialNumber   string            `json:"serialNumber" validate:"required,max=40"`
	CustomData     *CustomDataType   `json:"customData,omitempty" validate:"omitempty"`
}

// CertificateHashDataChain contains certificate hash data with child certificates.
type CertificateHashDataChain struct {
	CertificateType          CertificateUse        `json:"certificateType" validate:"required,certificateUse21"`
	CertificateHashData      CertificateHashData   `json:"certificateHashData" validate:"required"`
	ChildCertificateHashData []CertificateHashData `json:"childCertificateHashData,omitempty" validate:"omitempty,max=4,dive"`
	CustomData               *CustomDataType       `json:"customData,omitempty" validate:"omitempty"`
}

// Certificate15118EVStatus represents the status of a 15118 certificate operation.
type Certificate15118EVStatus string

const (
	Certificate15118EVStatusAccepted Certificate15118EVStatus = "Accepted"
	Certificate15118EVStatusFailed   Certificate15118EVStatus = "Failed"
)

func isValidCertificate15118EVStatus(fl validator.FieldLevel) bool {
	status := Certificate15118EVStatus(fl.Field().String())
	switch status {
	case Certificate15118EVStatusAccepted, Certificate15118EVStatusFailed:
		return true
	default:
		return false
	}
}

// AuthorizeCertificateStatus represents the certificate status in an AuthorizeResponse.
type AuthorizeCertificateStatus string

const (
	AuthorizeCertificateStatusAccepted               AuthorizeCertificateStatus = "Accepted"
	AuthorizeCertificateStatusSignatureError         AuthorizeCertificateStatus = "SignatureError"
	AuthorizeCertificateStatusCertificateExpired     AuthorizeCertificateStatus = "CertificateExpired"
	AuthorizeCertificateStatusCertificateRevoked     AuthorizeCertificateStatus = "CertificateRevoked"
	AuthorizeCertificateStatusNoCertificateAvailable AuthorizeCertificateStatus = "NoCertificateAvailable"
	AuthorizeCertificateStatusCertChainError         AuthorizeCertificateStatus = "CertChainError"
	AuthorizeCertificateStatusContractCancelled      AuthorizeCertificateStatus = "ContractCancelled"
)

func isValidAuthorizeCertificateStatus(fl validator.FieldLevel) bool {
	status := AuthorizeCertificateStatus(fl.Field().String())
	switch status {
	case AuthorizeCertificateStatusAccepted, AuthorizeCertificateStatusSignatureError,
		AuthorizeCertificateStatusCertificateExpired, AuthorizeCertificateStatusCertificateRevoked,
		AuthorizeCertificateStatusNoCertificateAvailable, AuthorizeCertificateStatusCertChainError,
		AuthorizeCertificateStatusContractCancelled:
		return true
	default:
		return false
	}
}

// CertificateSigningUse indicates the type of signed certificate.
type CertificateSigningUse string

const (
	ChargingStationCert CertificateSigningUse = "ChargingStationCertificate"
	V2GCertificate      CertificateSigningUse = "V2GCertificate"
	V2G20Certificate    CertificateSigningUse = "V2G20Certificate"
)

func isValidCertificateSigningUse(fl validator.FieldLevel) bool {
	status := CertificateSigningUse(fl.Field().String())
	switch status {
	case ChargingStationCert, V2GCertificate, V2G20Certificate:
		return true
	default:
		return false
	}
}

// CertificateUse indicates the type of requested certificate.
type CertificateUse string

const (
	V2GRootCertificate          CertificateUse = "V2GRootCertificate"
	MORootCertificate           CertificateUse = "MORootCertificate"
	CSMSRootCertificate         CertificateUse = "CSMSRootCertificate"
	V2GCertificateChain         CertificateUse = "V2GCertificateChain"
	ManufacturerRootCertificate CertificateUse = "ManufacturerRootCertificate"
	OEMRootCertificate          CertificateUse = "OEMRootCertificate" // New in 2.1
)

func isValidCertificateUse(fl validator.FieldLevel) bool {
	use := CertificateUse(fl.Field().String())
	switch use {
	case V2GRootCertificate, MORootCertificate, CSMSRootCertificate,
		V2GCertificateChain, ManufacturerRootCertificate, OEMRootCertificate:
		return true
	default:
		return false
	}
}

// MessageFormatType represents the format of a message.
type MessageFormatType string

const (
	MessageFormatASCII  MessageFormatType = "ASCII"
	MessageFormatHTML   MessageFormatType = "HTML"
	MessageFormatURI    MessageFormatType = "URI"
	MessageFormatUTF8   MessageFormatType = "UTF8"
	MessageFormatQRCode MessageFormatType = "QRCODE"
)

func isValidMessageFormatType(fl validator.FieldLevel) bool {
	format := MessageFormatType(fl.Field().String())
	switch format {
	case MessageFormatASCII, MessageFormatHTML, MessageFormatURI, MessageFormatUTF8, MessageFormatQRCode:
		return true
	default:
		return false
	}
}

// MessageContent contains a message with format information.
type MessageContent struct {
	Format     MessageFormatType `json:"format" validate:"required,messageFormat21"`
	Language   string            `json:"language,omitempty" validate:"omitempty,max=8"`
	Content    string            `json:"content" validate:"required,max=1024"`
	CustomData *CustomDataType   `json:"customData,omitempty" validate:"omitempty"`
}

// GroupIdToken is used to group multiple IdTokens together.
type GroupIdToken struct {
	IdToken    string          `json:"idToken" validate:"max=255"`
	Type       IdTokenType     `json:"type" validate:"required,max=20"`
	CustomData *CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

func isValidGroupIdToken(sl validator.StructLevel) {
	groupIdToken := sl.Current().Interface().(GroupIdToken)
	switch groupIdToken.Type {
	case IdTokenTypeCentral, IdTokenTypeEMAID, IdTokenTypeISO14443, IdTokenTypeISO15693,
		IdTokenTypeKeyCode, IdTokenTypeLocal, IdTokenTypeMacAddress:
		if groupIdToken.IdToken == "" {
			sl.ReportError(groupIdToken.IdToken, "IdToken", "IdToken", "required", "")
		}
	}
}

// IdTokenInfo contains information about an IdToken.
type IdTokenInfo struct {
	Status              AuthorizationStatus `json:"status" validate:"required,authorizationStatus21"`
	CacheExpiryDateTime *DateTime           `json:"cacheExpiryDateTime,omitempty" validate:"omitempty"`
	ChargingPriority    int                 `json:"chargingPriority,omitempty" validate:"min=-9,max=9"`
	Language1           string              `json:"language1,omitempty" validate:"omitempty,max=8"`
	Language2           string              `json:"language2,omitempty" validate:"omitempty,max=8"`
	EvseId              []int               `json:"evseId,omitempty" validate:"omitempty,min=1,dive,gte=0"`
	GroupIdToken        *IdToken            `json:"groupIdToken,omitempty" validate:"omitempty"`
	PersonalMessage     *MessageContent     `json:"personalMessage,omitempty" validate:"omitempty"`
	CustomData          *CustomDataType     `json:"customData,omitempty" validate:"omitempty"`
}

// NewIdTokenInfo creates an IdTokenInfo with the given status.
func NewIdTokenInfo(status AuthorizationStatus) *IdTokenInfo {
	return &IdTokenInfo{Status: status}
}

// StatusInfo provides additional information about a status.
type StatusInfo struct {
	ReasonCode     string          `json:"reasonCode" validate:"required,max=20"`
	AdditionalInfo string          `json:"additionalInfo,omitempty" validate:"omitempty,max=1024"`
	CustomData     *CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// NewStatusInfo creates a StatusInfo with the given parameters.
func NewStatusInfo(reasonCode string, additionalInfo string) *StatusInfo {
	return &StatusInfo{ReasonCode: reasonCode, AdditionalInfo: additionalInfo}
}

// EVSE represents an Electric Vehicle Supply Equipment.
type EVSE struct {
	ID          int             `json:"id" validate:"gte=0"`
	ConnectorID *int            `json:"connectorId,omitempty" validate:"omitempty,gte=0"`
	CustomData  *CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// Component represents a physical or logical component.
type Component struct {
	Name       string          `json:"name" validate:"required,max=50"`
	Instance   string          `json:"instance,omitempty" validate:"omitempty,max=50"`
	EVSE       *EVSE           `json:"evse,omitempty" validate:"omitempty"`
	CustomData *CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// Variable is a reference to a component variable.
type Variable struct {
	Name       string          `json:"name" validate:"required,max=50"`
	Instance   string          `json:"instance,omitempty" validate:"omitempty,max=50"`
	CustomData *CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// ComponentVariable is used to report components, variables and their attributes.
type ComponentVariable struct {
	Component  Component       `json:"component" validate:"required"`
	Variable   *Variable       `json:"variable,omitempty" validate:"omitempty"`
	CustomData *CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// Attribute is an enumeration type used when requesting a variable value.
type Attribute string

const (
	AttributeActual Attribute = "Actual"
	AttributeTarget Attribute = "Target"
	AttributeMinSet Attribute = "MinSet"
	AttributeMaxSet Attribute = "MaxSet"
)

func isValidAttribute(fl validator.FieldLevel) bool {
	attr := Attribute(fl.Field().String())
	switch attr {
	case AttributeActual, AttributeTarget, AttributeMinSet, AttributeMaxSet:
		return true
	default:
		return false
	}
}

// CostKind represents the type of cost.
type CostKind string

const (
	CostKindCarbonDioxideEmission         CostKind = "CarbonDioxideEmission"
	CostKindRelativePricePercentage       CostKind = "RelativePricePercentage"
	CostKindRenewableGenerationPercentage CostKind = "RenewableGenerationPercentage"
)

func isValidCostKind(fl validator.FieldLevel) bool {
	kind := CostKind(fl.Field().String())
	switch kind {
	case CostKindCarbonDioxideEmission, CostKindRelativePricePercentage, CostKindRenewableGenerationPercentage:
		return true
	default:
		return false
	}
}

// RelativeTimeInterval defines a time interval based on relative times.
type RelativeTimeInterval struct {
	Start      int             `json:"start"`
	Duration   *int            `json:"duration,omitempty" validate:"omitempty,gte=0"`
	CustomData *CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// CostType contains cost details.
type CostType struct {
	CostKind         CostKind        `json:"costKind" validate:"required,costKind21"`
	Amount           int             `json:"amount" validate:"gte=0"`
	AmountMultiplier *int            `json:"amountMultiplier,omitempty" validate:"omitempty,min=-3,max=3"`
	CustomData       *CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// ConsumptionCost contains price information and/or alternative costs.
type ConsumptionCost struct {
	StartValue float64         `json:"startValue"`
	Cost       []CostType      `json:"cost" validate:"required,max=3,dive"`
	CustomData *CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// NewConsumptionCost creates a new ConsumptionCost.
func NewConsumptionCost(startValue float64, cost []CostType) ConsumptionCost {
	return ConsumptionCost{StartValue: startValue, Cost: cost}
}

// SalesTariffEntry describes details for one time interval of a SalesTariff.
type SalesTariffEntry struct {
	EPriceLevel          *int                 `json:"ePriceLevel,omitempty" validate:"omitempty,gte=0"`
	RelativeTimeInterval RelativeTimeInterval `json:"relativeTimeInterval" validate:"required"`
	ConsumptionCost      []ConsumptionCost    `json:"consumptionCost,omitempty" validate:"omitempty,max=3,dive"`
	CustomData           *CustomDataType      `json:"customData,omitempty" validate:"omitempty"`
}

// SalesTariff represents a sales tariff associated with a charging schedule.
type SalesTariff struct {
	ID                     int                `json:"id"`
	SalesTariffDescription string             `json:"salesTariffDescription,omitempty" validate:"omitempty,max=32"`
	NumEPriceLevels        *int               `json:"numEPriceLevels,omitempty"`
	SalesTariffEntry       []SalesTariffEntry `json:"salesTariffEntry" validate:"required,min=1,max=1024,dive"`
	CustomData             *CustomDataType    `json:"customData,omitempty" validate:"omitempty"`
}

// NewSalesTariff creates a new SalesTariff.
func NewSalesTariff(id int, salesTariffEntries []SalesTariffEntry) *SalesTariff {
	return &SalesTariff{ID: id, SalesTariffEntry: salesTariffEntries}
}

// Charging Profile types

type ChargingProfilePurposeType string
type ChargingProfileKindType string
type RecurrencyKindType string
type ChargingRateUnitType string
type ChargingLimitSourceType string

const (
	ChargingProfilePurposeChargingStationExternalConstraints ChargingProfilePurposeType = "ChargingStationExternalConstraints"
	ChargingProfilePurposeChargingStationMaxProfile          ChargingProfilePurposeType = "ChargingStationMaxProfile"
	ChargingProfilePurposeTxDefaultProfile                   ChargingProfilePurposeType = "TxDefaultProfile"
	ChargingProfilePurposeTxProfile                          ChargingProfilePurposeType = "TxProfile"
	ChargingProfilePurposePriorityCharging                   ChargingProfilePurposeType = "PriorityCharging"
	ChargingProfilePurposeLocalGeneration                    ChargingProfilePurposeType = "LocalGeneration"
	ChargingProfileKindAbsolute                              ChargingProfileKindType    = "Absolute"
	ChargingProfileKindRecurring                             ChargingProfileKindType    = "Recurring"
	ChargingProfileKindRelative                              ChargingProfileKindType    = "Relative"
	ChargingProfileKindDynamic                               ChargingProfileKindType    = "Dynamic"
	RecurrencyKindDaily                                      RecurrencyKindType         = "Daily"
	RecurrencyKindWeekly                                     RecurrencyKindType         = "Weekly"
	ChargingRateUnitWatts                                    ChargingRateUnitType       = "W"
	ChargingRateUnitAmperes                                  ChargingRateUnitType       = "A"
	ChargingLimitSourceEMS                                   ChargingLimitSourceType    = "EMS"
	ChargingLimitSourceOther                                 ChargingLimitSourceType    = "Other"
	ChargingLimitSourceSO                                    ChargingLimitSourceType    = "SO"
	ChargingLimitSourceCSO                                   ChargingLimitSourceType    = "CSO"
)

func isValidChargingProfilePurpose(fl validator.FieldLevel) bool {
	purpose := ChargingProfilePurposeType(fl.Field().String())
	switch purpose {
	case ChargingProfilePurposeChargingStationExternalConstraints, ChargingProfilePurposeChargingStationMaxProfile,
		ChargingProfilePurposeTxDefaultProfile, ChargingProfilePurposeTxProfile,
		ChargingProfilePurposePriorityCharging, ChargingProfilePurposeLocalGeneration:
		return true
	default:
		return false
	}
}

func isValidChargingProfileKind(fl validator.FieldLevel) bool {
	kind := ChargingProfileKindType(fl.Field().String())
	switch kind {
	case ChargingProfileKindAbsolute, ChargingProfileKindRecurring, ChargingProfileKindRelative, ChargingProfileKindDynamic:
		return true
	default:
		return false
	}
}

func isValidRecurrencyKind(fl validator.FieldLevel) bool {
	kind := RecurrencyKindType(fl.Field().String())
	switch kind {
	case RecurrencyKindDaily, RecurrencyKindWeekly:
		return true
	default:
		return false
	}
}

func isValidChargingRateUnit(fl validator.FieldLevel) bool {
	unit := ChargingRateUnitType(fl.Field().String())
	switch unit {
	case ChargingRateUnitWatts, ChargingRateUnitAmperes:
		return true
	default:
		return false
	}
}

func isValidChargingLimitSource(fl validator.FieldLevel) bool {
	source := ChargingLimitSourceType(fl.Field().String())
	switch source {
	case ChargingLimitSourceEMS, ChargingLimitSourceOther, ChargingLimitSourceSO, ChargingLimitSourceCSO:
		return true
	default:
		return false
	}
}

// ChargingSchedulePeriod defines a period in a charging schedule.
type ChargingSchedulePeriod struct {
	StartPeriod            int             `json:"startPeriod" validate:"gte=0"`
	Limit                  float64         `json:"limit"`
	NumberPhases           *int            `json:"numberPhases,omitempty" validate:"omitempty,gte=0,lte=3"`
	PhaseToUse             *int            `json:"phaseToUse,omitempty" validate:"omitempty,gte=1,lte=3"`
	DischargeLimit         *float64        `json:"dischargeLimit,omitempty"`
	Setpoint               *float64        `json:"setpoint,omitempty"`
	SetpointReactive       *float64        `json:"setpointReactive,omitempty"`
	SetpointReactiveL2     *float64        `json:"setpointReactive_L2,omitempty"`
	SetpointReactiveL3     *float64        `json:"setpointReactive_L3,omitempty"`
	PreconditioningRequest *bool           `json:"preconditioningRequest,omitempty"`
	OperationMode          *OperationMode  `json:"operationMode,omitempty"`
	CustomData             *CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// NewChargingSchedulePeriod creates a new ChargingSchedulePeriod.
func NewChargingSchedulePeriod(startPeriod int, limit float64) ChargingSchedulePeriod {
	return ChargingSchedulePeriod{StartPeriod: startPeriod, Limit: limit}
}

// OperationMode represents the operation mode for a charging schedule period.
type OperationMode string

const (
	OperationModeIdle               OperationMode = "Idle"
	OperationModeChargingOnly       OperationMode = "ChargingOnly"
	OperationModeCentralSetpoint    OperationMode = "CentralSetpoint"
	OperationModeExternalSetpoint   OperationMode = "ExternalSetpoint"
	OperationModeExternalLimits     OperationMode = "ExternalLimits"
	OperationModeCentralFrequency   OperationMode = "CentralFrequency"
	OperationModeLocalFrequency     OperationMode = "LocalFrequency"
	OperationModeLocalLoadBalancing OperationMode = "LocalLoadBalancing"
)

func isValidOperationMode(fl validator.FieldLevel) bool {
	mode := OperationMode(fl.Field().String())
	switch mode {
	case OperationModeIdle, OperationModeChargingOnly, OperationModeCentralSetpoint,
		OperationModeExternalSetpoint, OperationModeExternalLimits, OperationModeCentralFrequency,
		OperationModeLocalFrequency, OperationModeLocalLoadBalancing:
		return true
	default:
		return false
	}
}

// ChargingSchedule represents a charging schedule.
type ChargingSchedule struct {
	ID                     int                      `json:"id" validate:"gte=0"`
	StartSchedule          *DateTime                `json:"startSchedule,omitempty"`
	Duration               *int                     `json:"duration,omitempty" validate:"omitempty,gte=0"`
	ChargingRateUnit       ChargingRateUnitType     `json:"chargingRateUnit" validate:"required,chargingRateUnit21"`
	MinChargingRate        *float64                 `json:"minChargingRate,omitempty" validate:"omitempty,gte=0"`
	ChargingSchedulePeriod []ChargingSchedulePeriod `json:"chargingSchedulePeriod" validate:"required,min=1,max=1024,dive"`
	SalesTariff            *SalesTariff             `json:"salesTariff,omitempty"`
	PowerTolerance         *float64                 `json:"powerTolerance,omitempty"`
	SignatureId            *int                     `json:"signatureId,omitempty"`
	DigestValue            string                   `json:"digestValue,omitempty" validate:"omitempty,max=88"`
	CustomData             *CustomDataType          `json:"customData,omitempty" validate:"omitempty"`
}

// NewChargingSchedule creates a new ChargingSchedule.
func NewChargingSchedule(id int, chargingRateUnit ChargingRateUnitType, periods ...ChargingSchedulePeriod) *ChargingSchedule {
	return &ChargingSchedule{ID: id, ChargingRateUnit: chargingRateUnit, ChargingSchedulePeriod: periods}
}

// ChargingProfile represents a charging profile.
type ChargingProfile struct {
	ID                     int                        `json:"id" validate:"gte=0"`
	StackLevel             int                        `json:"stackLevel" validate:"gte=0"`
	ChargingProfilePurpose ChargingProfilePurposeType `json:"chargingProfilePurpose" validate:"required,chargingProfilePurpose21"`
	ChargingProfileKind    ChargingProfileKindType    `json:"chargingProfileKind" validate:"required,chargingProfileKind21"`
	RecurrencyKind         RecurrencyKindType         `json:"recurrencyKind,omitempty" validate:"omitempty,recurrencyKind21"`
	ValidFrom              *DateTime                  `json:"validFrom,omitempty"`
	ValidTo                *DateTime                  `json:"validTo,omitempty"`
	TransactionID          string                     `json:"transactionId,omitempty" validate:"omitempty,max=36"`
	ChargingSchedule       []ChargingSchedule         `json:"chargingSchedule" validate:"required,min=1,max=3,dive"`
	CustomData             *CustomDataType            `json:"customData,omitempty" validate:"omitempty"`
}

// NewChargingProfile creates a new ChargingProfile.
func NewChargingProfile(id int, stackLevel int, purpose ChargingProfilePurposeType, kind ChargingProfileKindType, schedules []ChargingSchedule) *ChargingProfile {
	return &ChargingProfile{ID: id, StackLevel: stackLevel, ChargingProfilePurpose: purpose, ChargingProfileKind: kind, ChargingSchedule: schedules}
}

// RemoteStartStopStatus represents the status of a remote start/stop request.
type RemoteStartStopStatus string

const (
	RemoteStartStopStatusAccepted RemoteStartStopStatus = "Accepted"
	RemoteStartStopStatusRejected RemoteStartStopStatus = "Rejected"
)

func isValidRemoteStartStopStatus(fl validator.FieldLevel) bool {
	status := RemoteStartStopStatus(fl.Field().String())
	switch status {
	case RemoteStartStopStatusAccepted, RemoteStartStopStatusRejected:
		return true
	default:
		return false
	}
}

// EnergyTransferMode represents the energy transfer mode for charging.
type EnergyTransferMode string

const (
	EnergyTransferModeACSinglePhase EnergyTransferMode = "AC_single_phase"
	EnergyTransferModeACTwoPhase    EnergyTransferMode = "AC_two_phase"
	EnergyTransferModeACThreePhase  EnergyTransferMode = "AC_three_phase"
	EnergyTransferModeDC            EnergyTransferMode = "DC"
	EnergyTransferModeACBPT         EnergyTransferMode = "AC_BPT"
	EnergyTransferModeACBPTDER      EnergyTransferMode = "AC_BPT_DER"
	EnergyTransferModeACDER         EnergyTransferMode = "AC_DER"
	EnergyTransferModeDCBPT         EnergyTransferMode = "DC_BPT"
	EnergyTransferModeDCACDP        EnergyTransferMode = "DC_ACDP"
	EnergyTransferModeDCACDPBPT     EnergyTransferMode = "DC_ACDP_BPT"
	EnergyTransferModeWPT           EnergyTransferMode = "WPT"
)

func isValidEnergyTransferMode(fl validator.FieldLevel) bool {
	mode := EnergyTransferMode(fl.Field().String())
	switch mode {
	case EnergyTransferModeACSinglePhase, EnergyTransferModeACTwoPhase, EnergyTransferModeACThreePhase,
		EnergyTransferModeDC, EnergyTransferModeACBPT, EnergyTransferModeACBPTDER, EnergyTransferModeACDER,
		EnergyTransferModeDCBPT, EnergyTransferModeDCACDP, EnergyTransferModeDCACDPBPT, EnergyTransferModeWPT:
		return true
	default:
		return false
	}
}

// Meter Value types

type ReadingContext string
type Measurand string
type Phase string
type Location string

const (
	ReadingContextInterruptionBegin ReadingContext = "Interruption.Begin"
	ReadingContextInterruptionEnd   ReadingContext = "Interruption.End"
	ReadingContextOther             ReadingContext = "Other"
	ReadingContextSampleClock       ReadingContext = "Sample.Clock"
	ReadingContextSamplePeriodic    ReadingContext = "Sample.Periodic"
	ReadingContextTransactionBegin  ReadingContext = "Transaction.Begin"
	ReadingContextTransactionEnd    ReadingContext = "Transaction.End"
	ReadingContextTrigger           ReadingContext = "Trigger"

	// Original 2.0.1 measurands
	MeasurandCurrentExport                Measurand = "Current.Export"
	MeasurandCurrentImport                Measurand = "Current.Import"
	MeasurandCurrentOffered               Measurand = "Current.Offered"
	MeasurandEnergyActiveExportRegister   Measurand = "Energy.Active.Export.Register"
	MeasurandEnergyActiveImportRegister   Measurand = "Energy.Active.Import.Register"
	MeasurandEnergyReactiveExportRegister Measurand = "Energy.Reactive.Export.Register"
	MeasurandEnergyReactiveImportRegister Measurand = "Energy.Reactive.Import.Register"
	MeasurandEnergyActiveExportInterval   Measurand = "Energy.Active.Export.Interval"
	MeasurandEnergyActiveImportInterval   Measurand = "Energy.Active.Import.Interval"
	MeasurandEnergyActiveNet              Measurand = "Energy.Active.Net"
	MeasurandEnergyReactiveExportInterval Measurand = "Energy.Reactive.Export.Interval"
	MeasurandEnergyReactiveImportInterval Measurand = "Energy.Reactive.Import.Interval"
	MeasurandEnergyReactiveNet            Measurand = "Energy.Reactive.Net"
	MeasurandEnergyApparentNet            Measurand = "Energy.Apparent.Net"
	MeasurandEnergyApparentImport         Measurand = "Energy.Apparent.Import"
	MeasurandEnergyApparentExport         Measurand = "Energy.Apparent.Export"
	MeasurandFrequency                    Measurand = "Frequency"
	MeasurandPowerActiveExport            Measurand = "Power.Active.Export"
	MeasurandPowerActiveImport            Measurand = "Power.Active.Import"
	MeasurandPowerFactor                  Measurand = "Power.Factor"
	MeasurandPowerOffered                 Measurand = "Power.Offered"
	MeasurandPowerReactiveExport          Measurand = "Power.Reactive.Export"
	MeasurandPowerReactiveImport          Measurand = "Power.Reactive.Import"
	MeasurandSoC                          Measurand = "SoC"
	MeasurandVoltage                      Measurand = "Voltage"

	// New 2.1 measurands
	MeasurandCurrentExportOffered               Measurand = "Current.Export.Offered"
	MeasurandCurrentExportMinimum               Measurand = "Current.Export.Minimum"
	MeasurandCurrentImportOffered               Measurand = "Current.Import.Offered"
	MeasurandCurrentImportMinimum               Measurand = "Current.Import.Minimum"
	MeasurandDisplayPresentSOC                  Measurand = "Display.PresentSOC"
	MeasurandDisplayMinimumSOC                  Measurand = "Display.MinimumSOC"
	MeasurandDisplayTargetSOC                   Measurand = "Display.TargetSOC"
	MeasurandDisplayMaximumSOC                  Measurand = "Display.MaximumSOC"
	MeasurandDisplayRemainingTimeToMinimumSOC   Measurand = "Display.RemainingTimeToMinimumSOC"
	MeasurandDisplayRemainingTimeToTargetSOC    Measurand = "Display.RemainingTimeToTargetSOC"
	MeasurandDisplayRemainingTimeToMaximumSOC   Measurand = "Display.RemainingTimeToMaximumSOC"
	MeasurandDisplayChargingComplete            Measurand = "Display.ChargingComplete"
	MeasurandDisplayBatteryEnergyCapacity       Measurand = "Display.BatteryEnergyCapacity"
	MeasurandDisplayInletHot                    Measurand = "Display.InletHot"
	MeasurandEnergyActiveImportCableLoss        Measurand = "Energy.Active.Import.CableLoss"
	MeasurandEnergyActiveImportLocalGenRegister Measurand = "Energy.Active.Import.LocalGeneration.Register"
	MeasurandEnergyActiveSetpointInterval       Measurand = "Energy.Active.Setpoint.Interval"
	MeasurandEnergyRequestTarget                Measurand = "EnergyRequest.Target"
	MeasurandEnergyRequestMinimum               Measurand = "EnergyRequest.Minimum"
	MeasurandEnergyRequestMaximum               Measurand = "EnergyRequest.Maximum"
	MeasurandEnergyRequestMinimumV2X            Measurand = "EnergyRequest.Minimum.V2X"
	MeasurandEnergyRequestMaximumV2X            Measurand = "EnergyRequest.Maximum.V2X"
	MeasurandEnergyRequestBulk                  Measurand = "EnergyRequest.Bulk"
	MeasurandPowerActiveSetpoint                Measurand = "Power.Active.Setpoint"
	MeasurandPowerActiveResidual                Measurand = "Power.Active.Residual"
	MeasurandPowerExportMinimum                 Measurand = "Power.Export.Minimum"
	MeasurandPowerExportOffered                 Measurand = "Power.Export.Offered"
	MeasurandPowerImportOffered                 Measurand = "Power.Import.Offered"
	MeasurandPowerImportMinimum                 Measurand = "Power.Import.Minimum"
	MeasurandVoltageMinimum                     Measurand = "Voltage.Minimum"
	MeasurandVoltageMaximum                     Measurand = "Voltage.Maximum"

	PhaseL1   Phase = "L1"
	PhaseL2   Phase = "L2"
	PhaseL3   Phase = "L3"
	PhaseN    Phase = "N"
	PhaseL1N  Phase = "L1-N"
	PhaseL2N  Phase = "L2-N"
	PhaseL3N  Phase = "L3-N"
	PhaseL1L2 Phase = "L1-L2"
	PhaseL2L3 Phase = "L2-L3"
	PhaseL3L1 Phase = "L3-L1"

	LocationBody     Location = "Body"
	LocationCable    Location = "Cable"
	LocationEV       Location = "EV"
	LocationInlet    Location = "Inlet"
	LocationOutlet   Location = "Outlet"
	LocationUpstream Location = "Upstream" // New in 2.1
)

func isValidReadingContext(fl validator.FieldLevel) bool {
	ctx := ReadingContext(fl.Field().String())
	switch ctx {
	case ReadingContextInterruptionBegin, ReadingContextInterruptionEnd, ReadingContextOther,
		ReadingContextSampleClock, ReadingContextSamplePeriodic, ReadingContextTransactionBegin,
		ReadingContextTransactionEnd, ReadingContextTrigger:
		return true
	default:
		return false
	}
}

func isValidMeasurand(fl validator.FieldLevel) bool {
	m := Measurand(fl.Field().String())
	switch m {
	case MeasurandCurrentExport, MeasurandCurrentImport, MeasurandCurrentOffered,
		MeasurandEnergyActiveExportRegister, MeasurandEnergyActiveImportRegister,
		MeasurandEnergyReactiveExportRegister, MeasurandEnergyReactiveImportRegister,
		MeasurandEnergyActiveExportInterval, MeasurandEnergyActiveImportInterval,
		MeasurandEnergyActiveNet, MeasurandEnergyReactiveExportInterval,
		MeasurandEnergyReactiveImportInterval, MeasurandEnergyReactiveNet,
		MeasurandEnergyApparentNet, MeasurandEnergyApparentImport, MeasurandEnergyApparentExport,
		MeasurandFrequency, MeasurandPowerActiveExport, MeasurandPowerActiveImport,
		MeasurandPowerFactor, MeasurandPowerOffered, MeasurandPowerReactiveExport,
		MeasurandPowerReactiveImport, MeasurandSoC, MeasurandVoltage,
		// New 2.1 measurands
		MeasurandCurrentExportOffered, MeasurandCurrentExportMinimum,
		MeasurandCurrentImportOffered, MeasurandCurrentImportMinimum,
		MeasurandDisplayPresentSOC, MeasurandDisplayMinimumSOC, MeasurandDisplayTargetSOC,
		MeasurandDisplayMaximumSOC, MeasurandDisplayRemainingTimeToMinimumSOC,
		MeasurandDisplayRemainingTimeToTargetSOC, MeasurandDisplayRemainingTimeToMaximumSOC,
		MeasurandDisplayChargingComplete, MeasurandDisplayBatteryEnergyCapacity,
		MeasurandDisplayInletHot, MeasurandEnergyActiveImportCableLoss,
		MeasurandEnergyActiveImportLocalGenRegister, MeasurandEnergyActiveSetpointInterval,
		MeasurandEnergyRequestTarget, MeasurandEnergyRequestMinimum, MeasurandEnergyRequestMaximum,
		MeasurandEnergyRequestMinimumV2X, MeasurandEnergyRequestMaximumV2X, MeasurandEnergyRequestBulk,
		MeasurandPowerActiveSetpoint, MeasurandPowerActiveResidual, MeasurandPowerExportMinimum,
		MeasurandPowerExportOffered, MeasurandPowerImportOffered, MeasurandPowerImportMinimum,
		MeasurandVoltageMinimum, MeasurandVoltageMaximum:
		return true
	default:
		return false
	}
}

func isValidPhase(fl validator.FieldLevel) bool {
	p := Phase(fl.Field().String())
	switch p {
	case PhaseL1, PhaseL2, PhaseL3, PhaseN, PhaseL1N, PhaseL2N, PhaseL3N, PhaseL1L2, PhaseL2L3, PhaseL3L1:
		return true
	default:
		return false
	}
}

func isValidLocation(fl validator.FieldLevel) bool {
	loc := Location(fl.Field().String())
	switch loc {
	case LocationBody, LocationCable, LocationEV, LocationInlet, LocationOutlet, LocationUpstream:
		return true
	default:
		return false
	}
}

// UnitOfMeasure represents a unit of measurement with multiplier.
type UnitOfMeasure struct {
	Unit       string          `json:"unit,omitempty" validate:"omitempty,max=20"`
	Multiplier *int            `json:"multiplier,omitempty"`
	CustomData *CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// SignedMeterValue represents a signed version of a meter value.
type SignedMeterValue struct {
	SignedMeterData string          `json:"signedMeterData" validate:"required,max=32768"`
	SigningMethod   string          `json:"signingMethod,omitempty" validate:"omitempty,max=50"`
	EncodingMethod  string          `json:"encodingMethod" validate:"required,max=50"`
	PublicKey       string          `json:"publicKey,omitempty" validate:"omitempty,max=2500"`
	CustomData      *CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// SampledValue represents a single sampled value in meter values.
type SampledValue struct {
	Value            float64           `json:"value"`
	Context          ReadingContext    `json:"context,omitempty" validate:"omitempty,readingContext21"`
	Measurand        Measurand         `json:"measurand,omitempty" validate:"omitempty,measurand21"`
	Phase            Phase             `json:"phase,omitempty" validate:"omitempty,phase21"`
	Location         Location          `json:"location,omitempty" validate:"omitempty,location21"`
	SignedMeterValue *SignedMeterValue `json:"signedMeterValue,omitempty"`
	UnitOfMeasure    *UnitOfMeasure    `json:"unitOfMeasure,omitempty"`
	CustomData       *CustomDataType   `json:"customData,omitempty" validate:"omitempty"`
}

// MeterValue represents a collection of sampled values.
type MeterValue struct {
	Timestamp    DateTime        `json:"timestamp" validate:"required"`
	SampledValue []SampledValue  `json:"sampledValue" validate:"required,min=1,dive"`
	CustomData   *CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// ChargingState represents the current charging state of a transaction.
type ChargingState string

const (
	ChargingStateEVConnected   ChargingState = "EVConnected"
	ChargingStateCharging      ChargingState = "Charging"
	ChargingStateSuspendedEV   ChargingState = "SuspendedEV"
	ChargingStateSuspendedEVSE ChargingState = "SuspendedEVSE"
	ChargingStateIdle          ChargingState = "Idle"
)

func isValidChargingState(fl validator.FieldLevel) bool {
	state := ChargingState(fl.Field().String())
	switch state {
	case ChargingStateEVConnected, ChargingStateCharging, ChargingStateSuspendedEV,
		ChargingStateSuspendedEVSE, ChargingStateIdle:
		return true
	default:
		return false
	}
}

// TransactionEvent represents the type of a transaction event.
type TransactionEvent string

const (
	TransactionEventEnded   TransactionEvent = "Ended"
	TransactionEventStarted TransactionEvent = "Started"
	TransactionEventUpdated TransactionEvent = "Updated"
)

func isValidTransactionEvent(fl validator.FieldLevel) bool {
	event := TransactionEvent(fl.Field().String())
	switch event {
	case TransactionEventEnded, TransactionEventStarted, TransactionEventUpdated:
		return true
	default:
		return false
	}
}

// TriggerReason represents the reason for sending a TransactionEvent.
type TriggerReason string

const (
	TriggerReasonAbnormalCondition    TriggerReason = "AbnormalCondition"
	TriggerReasonAuthorized           TriggerReason = "Authorized"
	TriggerReasonCablePluggedIn       TriggerReason = "CablePluggedIn"
	TriggerReasonChargingRateChanged  TriggerReason = "ChargingRateChanged"
	TriggerReasonChargingStateChanged TriggerReason = "ChargingStateChanged"
	TriggerReasonCostLimitReached     TriggerReason = "CostLimitReached"
	TriggerReasonDeauthorized         TriggerReason = "Deauthorized"
	TriggerReasonEnergyLimitReached   TriggerReason = "EnergyLimitReached"
	TriggerReasonEVCommunicationLost  TriggerReason = "EVCommunicationLost"
	TriggerReasonEVConnectTimeout     TriggerReason = "EVConnectTimeout"
	TriggerReasonEVDeparted           TriggerReason = "EVDeparted"
	TriggerReasonEVDetected           TriggerReason = "EVDetected"
	TriggerReasonLimitSet             TriggerReason = "LimitSet"
	TriggerReasonMeterValueClock      TriggerReason = "MeterValueClock"
	TriggerReasonMeterValuePeriodic   TriggerReason = "MeterValuePeriodic"
	TriggerReasonOperationModeChanged TriggerReason = "OperationModeChanged"
	TriggerReasonRemoteStart          TriggerReason = "RemoteStart"
	TriggerReasonRemoteStop           TriggerReason = "RemoteStop"
	TriggerReasonResetCommand         TriggerReason = "ResetCommand"
	TriggerReasonRunningCost          TriggerReason = "RunningCost"
	TriggerReasonSignedDataReceived   TriggerReason = "SignedDataReceived"
	TriggerReasonSoCLimitReached      TriggerReason = "SoCLimitReached"
	TriggerReasonStopAuthorized       TriggerReason = "StopAuthorized"
	TriggerReasonTariffChanged        TriggerReason = "TariffChanged"
	TriggerReasonTariffNotAccepted    TriggerReason = "TariffNotAccepted"
	TriggerReasonTimeLimitReached     TriggerReason = "TimeLimitReached"
	TriggerReasonTrigger              TriggerReason = "Trigger"
	TriggerReasonTxResumed            TriggerReason = "TxResumed"
	TriggerReasonUnlockCommand        TriggerReason = "UnlockCommand"
)

func isValidTriggerReason(fl validator.FieldLevel) bool {
	reason := TriggerReason(fl.Field().String())
	switch reason {
	case TriggerReasonAbnormalCondition, TriggerReasonAuthorized, TriggerReasonCablePluggedIn,
		TriggerReasonChargingRateChanged, TriggerReasonChargingStateChanged, TriggerReasonCostLimitReached,
		TriggerReasonDeauthorized, TriggerReasonEnergyLimitReached, TriggerReasonEVCommunicationLost,
		TriggerReasonEVConnectTimeout, TriggerReasonEVDeparted, TriggerReasonEVDetected,
		TriggerReasonLimitSet, TriggerReasonMeterValueClock, TriggerReasonMeterValuePeriodic,
		TriggerReasonOperationModeChanged, TriggerReasonRemoteStart, TriggerReasonRemoteStop,
		TriggerReasonResetCommand, TriggerReasonRunningCost, TriggerReasonSignedDataReceived,
		TriggerReasonSoCLimitReached, TriggerReasonStopAuthorized, TriggerReasonTariffChanged,
		TriggerReasonTariffNotAccepted, TriggerReasonTimeLimitReached, TriggerReasonTrigger,
		TriggerReasonTxResumed, TriggerReasonUnlockCommand:
		return true
	default:
		return false
	}
}

// StoppedReason represents the reason a transaction was stopped.
type StoppedReason string

const (
	StoppedReasonDeAuthorized              StoppedReason = "DeAuthorized"
	StoppedReasonEmergencyStop             StoppedReason = "EmergencyStop"
	StoppedReasonEnergyLimitReached        StoppedReason = "EnergyLimitReached"
	StoppedReasonEVDisconnected            StoppedReason = "EVDisconnected"
	StoppedReasonGroundFault               StoppedReason = "GroundFault"
	StoppedReasonImmediateReset            StoppedReason = "ImmediateReset"
	StoppedReasonMasterPass                StoppedReason = "MasterPass"
	StoppedReasonLocal                     StoppedReason = "Local"
	StoppedReasonLocalOutOfCredit          StoppedReason = "LocalOutOfCredit"
	StoppedReasonOther                     StoppedReason = "Other"
	StoppedReasonOvercurrentFault          StoppedReason = "OvercurrentFault"
	StoppedReasonPowerLoss                 StoppedReason = "PowerLoss"
	StoppedReasonPowerQuality              StoppedReason = "PowerQuality"
	StoppedReasonReboot                    StoppedReason = "Reboot"
	StoppedReasonRemote                    StoppedReason = "Remote"
	StoppedReasonSOCLimitReached           StoppedReason = "SOCLimitReached"
	StoppedReasonStoppedByEV               StoppedReason = "StoppedByEV"
	StoppedReasonTimeLimitReached          StoppedReason = "TimeLimitReached"
	StoppedReasonTimeout                   StoppedReason = "Timeout"
	StoppedReasonReqEnergyTransferRejected StoppedReason = "ReqEnergyTransferRejected"
)

func isValidStoppedReason(fl validator.FieldLevel) bool {
	reason := StoppedReason(fl.Field().String())
	switch reason {
	case StoppedReasonDeAuthorized, StoppedReasonEmergencyStop, StoppedReasonEnergyLimitReached,
		StoppedReasonEVDisconnected, StoppedReasonGroundFault, StoppedReasonImmediateReset,
		StoppedReasonMasterPass, StoppedReasonLocal, StoppedReasonLocalOutOfCredit,
		StoppedReasonOther, StoppedReasonOvercurrentFault, StoppedReasonPowerLoss,
		StoppedReasonPowerQuality, StoppedReasonReboot, StoppedReasonRemote,
		StoppedReasonSOCLimitReached, StoppedReasonStoppedByEV, StoppedReasonTimeLimitReached,
		StoppedReasonTimeout, StoppedReasonReqEnergyTransferRejected:
		return true
	default:
		return false
	}
}

// PreconditioningStatus represents the preconditioning status of the EV battery.
type PreconditioningStatus string

const (
	PreconditioningStatusUnknown         PreconditioningStatus = "Unknown"
	PreconditioningStatusReady           PreconditioningStatus = "Ready"
	PreconditioningStatusNotReady        PreconditioningStatus = "NotReady"
	PreconditioningStatusPreconditioning PreconditioningStatus = "Preconditioning"
)

func isValidPreconditioningStatus(fl validator.FieldLevel) bool {
	status := PreconditioningStatus(fl.Field().String())
	switch status {
	case PreconditioningStatusUnknown, PreconditioningStatusReady,
		PreconditioningStatusNotReady, PreconditioningStatusPreconditioning:
		return true
	default:
		return false
	}
}

// TransactionLimit represents cost, energy, time or SoC limits for a transaction.
type TransactionLimit struct {
	MaxCost    *float64        `json:"maxCost,omitempty"`
	MaxEnergy  *float64        `json:"maxEnergy,omitempty"`
	MaxTime    *int            `json:"maxTime,omitempty"`
	MaxSoC     *int            `json:"maxSoC,omitempty" validate:"omitempty,gte=0,lte=100"`
	CustomData *CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// Transaction contains information about a transaction.
type Transaction struct {
	TransactionID     string            `json:"transactionId" validate:"required,max=36"`
	ChargingState     ChargingState     `json:"chargingState,omitempty" validate:"omitempty,chargingState21"`
	TimeSpentCharging *int              `json:"timeSpentCharging,omitempty"`
	StoppedReason     StoppedReason     `json:"stoppedReason,omitempty" validate:"omitempty,stoppedReason21"`
	RemoteStartID     *int              `json:"remoteStartId,omitempty"`
	OperationMode     OperationMode     `json:"operationMode,omitempty" validate:"omitempty,operationMode21"`
	TariffID          string            `json:"tariffId,omitempty" validate:"omitempty,max=60"`
	TransactionLimit  *TransactionLimit `json:"transactionLimit,omitempty" validate:"omitempty"`
	CustomData        *CustomDataType   `json:"customData,omitempty" validate:"omitempty"`
}

// MessageTrigger represents the type of message to be triggered.
type MessageTrigger string

const (
	MessageTriggerBootNotification                  MessageTrigger = "BootNotification"
	MessageTriggerLogStatusNotification             MessageTrigger = "LogStatusNotification"
	MessageTriggerFirmwareStatusNotification        MessageTrigger = "FirmwareStatusNotification"
	MessageTriggerHeartbeat                         MessageTrigger = "Heartbeat"
	MessageTriggerMeterValues                       MessageTrigger = "MeterValues"
	MessageTriggerSignChargingStationCertificate    MessageTrigger = "SignChargingStationCertificate"
	MessageTriggerSignV2GCertificate                MessageTrigger = "SignV2GCertificate"
	MessageTriggerSignV2G20Certificate              MessageTrigger = "SignV2G20Certificate"
	MessageTriggerStatusNotification                MessageTrigger = "StatusNotification"
	MessageTriggerTransactionEvent                  MessageTrigger = "TransactionEvent"
	MessageTriggerSignCombinedCertificate           MessageTrigger = "SignCombinedCertificate"
	MessageTriggerPublishFirmwareStatusNotification MessageTrigger = "PublishFirmwareStatusNotification"
	MessageTriggerCustomTrigger                     MessageTrigger = "CustomTrigger"
)

func isValidMessageTrigger(fl validator.FieldLevel) bool {
	trigger := MessageTrigger(fl.Field().String())
	switch trigger {
	case MessageTriggerBootNotification, MessageTriggerLogStatusNotification,
		MessageTriggerFirmwareStatusNotification, MessageTriggerHeartbeat,
		MessageTriggerMeterValues, MessageTriggerSignChargingStationCertificate,
		MessageTriggerSignV2GCertificate, MessageTriggerSignV2G20Certificate,
		MessageTriggerStatusNotification, MessageTriggerTransactionEvent,
		MessageTriggerSignCombinedCertificate, MessageTriggerPublishFirmwareStatusNotification,
		MessageTriggerCustomTrigger:
		return true
	default:
		return false
	}
}

// TriggerMessageStatus indicates whether the Charging Station will send the requested notification.
type TriggerMessageStatus string

const (
	TriggerMessageStatusAccepted       TriggerMessageStatus = "Accepted"
	TriggerMessageStatusRejected       TriggerMessageStatus = "Rejected"
	TriggerMessageStatusNotImplemented TriggerMessageStatus = "NotImplemented"
)

func isValidTriggerMessageStatus(fl validator.FieldLevel) bool {
	status := TriggerMessageStatus(fl.Field().String())
	switch status {
	case TriggerMessageStatusAccepted, TriggerMessageStatusRejected, TriggerMessageStatusNotImplemented:
		return true
	default:
		return false
	}
}

// UnlockStatus indicates whether the Charging Station has unlocked the connector.
type UnlockStatus string

const (
	UnlockStatusUnlocked                     UnlockStatus = "Unlocked"
	UnlockStatusUnlockFailed                 UnlockStatus = "UnlockFailed"
	UnlockStatusOngoingAuthorizedTransaction UnlockStatus = "OngoingAuthorizedTransaction"
	UnlockStatusUnknownConnector             UnlockStatus = "UnknownConnector"
)

func isValidUnlockStatus(fl validator.FieldLevel) bool {
	status := UnlockStatus(fl.Field().String())
	switch status {
	case UnlockStatusUnlocked, UnlockStatusUnlockFailed,
		UnlockStatusOngoingAuthorizedTransaction, UnlockStatusUnknownConnector:
		return true
	default:
		return false
	}
}

// CostDimension represents the type of cost dimension.
type CostDimension string

const (
	CostDimensionEnergy       CostDimension = "Energy"
	CostDimensionMaxCurrent   CostDimension = "MaxCurrent"
	CostDimensionMinCurrent   CostDimension = "MinCurrent"
	CostDimensionMaxPower     CostDimension = "MaxPower"
	CostDimensionMinPower     CostDimension = "MinPower"
	CostDimensionIdleTime     CostDimension = "IdleTIme" // Note: schema typo preserved
	CostDimensionChargingTime CostDimension = "ChargingTime"
)

func isValidCostDimension(fl validator.FieldLevel) bool {
	dim := CostDimension(fl.Field().String())
	switch dim {
	case CostDimensionEnergy, CostDimensionMaxCurrent, CostDimensionMinCurrent,
		CostDimensionMaxPower, CostDimensionMinPower, CostDimensionIdleTime,
		CostDimensionChargingTime:
		return true
	default:
		return false
	}
}

// TariffCost represents the type of cost: normal, minimum or maximum.
type TariffCost string

const (
	TariffCostNormal TariffCost = "NormalCost"
	TariffCostMin    TariffCost = "MinCost"
	TariffCostMax    TariffCost = "MaxCost"
)

func isValidTariffCost(fl validator.FieldLevel) bool {
	cost := TariffCost(fl.Field().String())
	switch cost {
	case TariffCostNormal, TariffCostMin, TariffCostMax:
		return true
	default:
		return false
	}
}

// TaxRate represents a tax percentage.
type TaxRate struct {
	Type       string          `json:"type" validate:"required,max=20"`
	Tax        float64         `json:"tax" validate:"required"`
	Stack      *int            `json:"stack,omitempty" validate:"omitempty,gte=0"`
	CustomData *CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// Price represents a price with and without tax.
type Price struct {
	ExclTax    *float64        `json:"exclTax,omitempty"`
	InclTax    *float64        `json:"inclTax,omitempty"`
	TaxRates   []TaxRate       `json:"taxRates,omitempty" validate:"omitempty,min=1,max=5,dive"`
	CustomData *CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// TotalPrice represents the total cost with and without tax.
type TotalPrice struct {
	ExclTax    *float64        `json:"exclTax,omitempty"`
	InclTax    *float64        `json:"inclTax,omitempty"`
	CustomData *CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// TotalUsage represents the calculated usage during a transaction.
type TotalUsage struct {
	Energy          float64         `json:"energy" validate:"required"`
	ChargingTime    int             `json:"chargingTime" validate:"required"`
	IdleTime        int             `json:"idleTime" validate:"required"`
	ReservationTime *int            `json:"reservationTime,omitempty"`
	CustomData      *CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// TotalCost represents the cost calculated during a transaction.
type TotalCost struct {
	Currency         string          `json:"currency" validate:"required,max=3"`
	TypeOfCost       TariffCost      `json:"typeOfCost,omitempty" validate:"omitempty,tariffCost21"`
	Fixed            *Price          `json:"fixed,omitempty" validate:"omitempty"`
	Energy           *Price          `json:"energy,omitempty" validate:"omitempty"`
	ChargingTime     *Price          `json:"chargingTime,omitempty" validate:"omitempty"`
	IdleTime         *Price          `json:"idleTime,omitempty" validate:"omitempty"`
	ReservationTime  *Price          `json:"reservationTime,omitempty" validate:"omitempty"`
	ReservationFixed *Price          `json:"reservationFixed,omitempty" validate:"omitempty"`
	Total            *TotalPrice     `json:"total,omitempty" validate:"omitempty"`
	CustomData       *CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// CostDimensionValue represents the volume consumed of a cost dimension.
type CostDimensionValue struct {
	Type       CostDimension   `json:"type" validate:"required,costDimension21"`
	Volume     float64         `json:"volume" validate:"required"`
	CustomData *CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// ChargingPeriod represents a period during charging with associated costs.
type ChargingPeriod struct {
	StartPeriod *DateTime            `json:"startPeriod" validate:"required"`
	Dimensions  []CostDimensionValue `json:"dimensions,omitempty" validate:"omitempty,min=1,dive"`
	TariffID    string               `json:"tariffId,omitempty" validate:"omitempty,max=60"`
	CustomData  *CustomDataType      `json:"customData,omitempty" validate:"omitempty"`
}

// CostDetails represents the cost calculated by Charging Station based on a tariff.
type CostDetails struct {
	TotalCost          TotalCost        `json:"totalCost" validate:"required"`
	TotalUsage         TotalUsage       `json:"totalUsage" validate:"required"`
	ChargingPeriods    []ChargingPeriod `json:"chargingPeriods,omitempty" validate:"omitempty,min=1,dive"`
	FailureToCalculate bool             `json:"failureToCalculate,omitempty"`
	FailureReason      string           `json:"failureReason,omitempty" validate:"omitempty,max=500"`
	CustomData         *CustomDataType  `json:"customData,omitempty" validate:"omitempty"`
}

// Validate is the validator used for all OCPP 2.1 messages.
var Validate = ocppj.Validate

func init() {
	// Register all validators with "21" suffix to avoid conflicts with 2.0.1
	_ = Validate.RegisterValidation("authorizationStatus21", isValidAuthorizationStatus)
	_ = Validate.RegisterValidation("idTokenType21", isValidIdTokenType)
	_ = Validate.RegisterValidation("genericDeviceModelStatus21", isValidGenericDeviceModelStatus)
	_ = Validate.RegisterValidation("genericStatus21", isValidGenericStatus)
	_ = Validate.RegisterValidation("hashAlgorithm21", isValidHashAlgorithmType)
	_ = Validate.RegisterValidation("messageFormat21", isValidMessageFormatType)
	_ = Validate.RegisterValidation("certificateSigningUse21", isValidCertificateSigningUse)
	_ = Validate.RegisterValidation("certificateUse21", isValidCertificateUse)
	_ = Validate.RegisterValidation("15118EVCertificate21", isValidCertificate15118EVStatus)
	_ = Validate.RegisterValidation("attribute21", isValidAttribute)
	_ = Validate.RegisterValidation("costKind21", isValidCostKind)
	_ = Validate.RegisterValidation("chargingProfilePurpose21", isValidChargingProfilePurpose)
	_ = Validate.RegisterValidation("chargingProfileKind21", isValidChargingProfileKind)
	_ = Validate.RegisterValidation("recurrencyKind21", isValidRecurrencyKind)
	_ = Validate.RegisterValidation("chargingRateUnit21", isValidChargingRateUnit)
	_ = Validate.RegisterValidation("chargingLimitSource21", isValidChargingLimitSource)
	_ = Validate.RegisterValidation("operationMode21", isValidOperationMode)
	_ = Validate.RegisterValidation("remoteStartStopStatus21", isValidRemoteStartStopStatus)
	_ = Validate.RegisterValidation("authorizeCertificateStatus21", isValidAuthorizeCertificateStatus)
	_ = Validate.RegisterValidation("energyTransferMode21", isValidEnergyTransferMode)
	_ = Validate.RegisterValidation("chargingState21", isValidChargingState)
	_ = Validate.RegisterValidation("transactionEvent21", isValidTransactionEvent)
	_ = Validate.RegisterValidation("triggerReason21", isValidTriggerReason)
	_ = Validate.RegisterValidation("stoppedReason21", isValidStoppedReason)
	_ = Validate.RegisterValidation("preconditioningStatus21", isValidPreconditioningStatus)
	_ = Validate.RegisterValidation("readingContext21", isValidReadingContext)
	_ = Validate.RegisterValidation("measurand21", isValidMeasurand)
	_ = Validate.RegisterValidation("phase21", isValidPhase)
	_ = Validate.RegisterValidation("location21", isValidLocation)
	_ = Validate.RegisterValidation("messageTrigger21", isValidMessageTrigger)
	_ = Validate.RegisterValidation("triggerMessageStatus21", isValidTriggerMessageStatus)
	_ = Validate.RegisterValidation("unlockStatus21", isValidUnlockStatus)
	_ = Validate.RegisterValidation("costDimension21", isValidCostDimension)
	_ = Validate.RegisterValidation("tariffCost21", isValidTariffCost)

	// Register struct validators
	Validate.RegisterStructValidation(isValidIdToken, IdToken{})
	Validate.RegisterStructValidation(isValidGroupIdToken, GroupIdToken{})
}
