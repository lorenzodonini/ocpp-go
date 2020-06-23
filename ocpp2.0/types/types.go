// Contains common and shared data types between OCPP 2.0 messages.
package types

import (
	"github.com/lorenzodonini/ocpp-go/ocppj"
	"gopkg.in/go-playground/validator.v9"
)

const (
	V2Subprotocol = "ocpp2.0"
)

type PropertyViolation struct {
	error
	Property string
}

func (e *PropertyViolation) Error() string {
	return ""
}

type AuthorizationStatus string

const (
	AuthorizationStatusAccepted           AuthorizationStatus = "Accepted"
	AuthorizationStatusBlocked            AuthorizationStatus = "Blocked"
	AuthorizationStatusExpired            AuthorizationStatus = "Expired"
	AuthorizationStatusInvalid            AuthorizationStatus = "Invalid"
	AuthorizationStatusConcurrentTx       AuthorizationStatus = "ConcurrentTx"
	AuthorizationStatusNoCredit           AuthorizationStatus = "NoCredit"
	AuthorizationStatusNotAllowedTypeEVSE AuthorizationStatus = "NotAllowedTypeEVS"
	AuthorizationStatusNotAtThisLocation  AuthorizationStatus = "NotAtThisLocation"
	AuthorizationStatusNotAtThisTime      AuthorizationStatus = "NotAtThisTime"
	AuthorizationStatusUnknown            AuthorizationStatus = "Unknown"
)

func isValidAuthorizationStatus(fl validator.FieldLevel) bool {
	status := AuthorizationStatus(fl.Field().String())
	switch status {
	case AuthorizationStatusAccepted, AuthorizationStatusBlocked, AuthorizationStatusExpired, AuthorizationStatusInvalid, AuthorizationStatusConcurrentTx, AuthorizationStatusNoCredit, AuthorizationStatusNotAllowedTypeEVSE, AuthorizationStatusNotAtThisLocation, AuthorizationStatusNotAtThisTime, AuthorizationStatusUnknown:
		return true
	default:
		return false
	}
}

// ID Token
type IdTokenType string

const (
	IdTokenTypeCentral         IdTokenType = "Central"
	IdTokenTypeEMAID           IdTokenType = "eMAID"
	IdTokenTypeISO14443        IdTokenType = "ISO14443"
	IdTokenTypeKeyCode         IdTokenType = "KeyCode"
	IdTokenTypeLocal           IdTokenType = "Local"
	IdTokenTypeNoAuthorization IdTokenType = "NoAuthorization"
	IdTokenTypeISO15693        IdTokenType = "ISO15693"
)

func isValidIdTokenType(fl validator.FieldLevel) bool {
	tokenType := IdTokenType(fl.Field().String())
	switch tokenType {
	case IdTokenTypeCentral, IdTokenTypeEMAID, IdTokenTypeISO14443, IdTokenTypeKeyCode, IdTokenTypeLocal, IdTokenTypeNoAuthorization, IdTokenTypeISO15693:
		return true
	default:
		return false
	}
}

type AdditionalInfo struct {
	AdditionalIdToken string `json:"additionalIdToken" validate:"required,max=36"`
	Type              string `json:"type" validate:"required,max=50"`
}

type IdToken struct {
	IdToken        string           `json:"idToken" validate:"required,max=36"`
	Type           IdTokenType      `json:"type" validate:"required,idTokenType"`
	AdditionalInfo []AdditionalInfo `json:"additionalInfo,omitempty" validate:"omitempty,dive"`
}

// Generic Device Model Status
type GenericDeviceModelStatus string

const (
	GenericDeviceModelStatusAccepted     GenericDeviceModelStatus = "Accepted"
	GenericDeviceModelStatusRejected     GenericDeviceModelStatus = "Rejected"
	GenericDeviceModelStatusNotSupported GenericDeviceModelStatus = "NotSupported"
)

func isValidGenericDeviceModelStatus(fl validator.FieldLevel) bool {
	status := GenericDeviceModelStatus(fl.Field().String())
	switch status {
	case GenericDeviceModelStatusAccepted, GenericDeviceModelStatusRejected, GenericDeviceModelStatusNotSupported:
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

// Hash Algorithms
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

// OCSPRequestDataType
type OCSPRequestDataType struct {
	HashAlgorithm  HashAlgorithmType `json:"hashAlgorithm" validate:"required,hashAlgorithm"`
	IssuerNameHash string            `json:"issuerNameHash" validate:"required,max=128"`
	IssuerKeyHash  string            `json:"issuerKeyHash" validate:"required,max=128"`
	SerialNumber   string            `json:"serialNumber" validate:"required,max=20"`
	ResponderURL   string            `json:"responderURL,omitempty" validate:"max=512"`
}

// CertificateHashDataType
type CertificateHashData struct {
	HashAlgorithm  HashAlgorithmType `json:"hashAlgorithm" validate:"required,hashAlgorithm"`
	IssuerNameHash string            `json:"issuerNameHash" validate:"required,max=128"`
	IssuerKeyHash  string            `json:"issuerKeyHash" validate:"required,max=128"`
	SerialNumber   string            `json:"serialNumber" validate:"required,max=20"`
}

// CertificateStatus
type CertificateStatus string

const (
	CertificateStatusAccepted               CertificateStatus = "Accepted"
	CertificateStatusSignatureError         CertificateStatus = "SignatureError"
	CertificateStatusCertificateExpired     CertificateStatus = "CertificateExpired"
	CertificateStatusCertificateRevoked     CertificateStatus = "CertificateRevoked"
	CertificateStatusNoCertificateAvailable CertificateStatus = "NoCertificateAvailable"
	CertificateStatusCertChainError         CertificateStatus = "CertChainError"
	CertificateStatusContractCancelled      CertificateStatus = "ContractCancelled"
)

func isValidCertificateStatus(fl validator.FieldLevel) bool {
	status := CertificateStatus(fl.Field().String())
	switch status {
	case CertificateStatusAccepted, CertificateStatusCertChainError, CertificateStatusCertificateExpired, CertificateStatusSignatureError, CertificateStatusNoCertificateAvailable, CertificateStatusCertificateRevoked, CertificateStatusContractCancelled:
		return true
	default:
		return false
	}
}

// Certificate15118EVStatus
type Certificate15118EVStatus string

const (
	Certificate15188EVStatusAccepted Certificate15118EVStatus = "Accepted"
	Certificate15118EVStatusFailed   Certificate15118EVStatus = "Failed"
)

func isValidCertificate15118EVStatus(fl validator.FieldLevel) bool {
	status := Certificate15118EVStatus(fl.Field().String())
	switch status {
	case Certificate15188EVStatusAccepted, Certificate15118EVStatusFailed:
		return true
	default:
		return false
	}
}

// Indicates the type of the signed certificate that is returned.
// When omitted the certificate is used for both the 15118 connection (if implemented) and the Charging Station to CSMS connection.
// This field is required when a typeOfCertificate was included in the SignCertificateRequest that requested this certificate to be signed AND both the 15118 connection and the Charging Station connection are implemented.
type CertificateSigningUse string

const (
	ChargingStationCert CertificateSigningUse = "ChargingStationCertificate"
	V2GCertificate      CertificateSigningUse = "V2GCertificate"
)

func isValidCertificateSigningUse(fl validator.FieldLevel) bool {
	status := CertificateSigningUse(fl.Field().String())
	switch status {
	case ChargingStationCert, V2GCertificate:
		return true
	default:
		return false
	}
}

// Indicates the type of the requested certificate.
// It is used in GetInstalledCertificateIdsRequest and InstallCertificateRequest messages.
type CertificateUse string

const (
	V2GRootCertificate          CertificateUse = "V2GRootCertificate"
	MORootCertificate           CertificateUse = "MORootCertificate"
	CSOSubCA1                   CertificateUse = "CSOSubCA1"
	CSOSubCA2                   CertificateUse = "CSOSubCA2"
	CSMSRootCertificate         CertificateUse = "CSMSRootCertificate"
	ManufacturerRootCertificate CertificateUse = "ManufacturerRootCertificate"
)

func isValidCertificateUse(fl validator.FieldLevel) bool {
	use := CertificateUse(fl.Field().String())
	switch use {
	case V2GRootCertificate, MORootCertificate, CSOSubCA1, CSOSubCA2, CSMSRootCertificate, ManufacturerRootCertificate:
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
	Format   MessageFormatType `json:"format" validate:"required,messageFormat"`
	Language string            `json:"language,omitempty" validate:"max=8"`
	Content  string            `json:"content" validate:"required,max=512"`
}

type GroupIdToken struct {
	IdToken string      `json:"idToken" validate:"required,max=36"`
	Type    IdTokenType `json:"type" validate:"required,idTokenType"`
}

type IdTokenInfo struct {
	Status              AuthorizationStatus `json:"status" validate:"required,authorizationStatus"`
	CacheExpiryDateTime *DateTime           `json:"cacheExpiryDateTime,omitempty" validate:"omitempty"`
	ChargingPriority    int                 `json:"chargingPriority,omitempty" validate:"min=-9,max=9"`
	Language1           string              `json:"language1,omitempty" validate:"max=8"`
	Language2           string              `json:"language2,omitempty" validate:"max=8"`
	GroupIdToken        *GroupIdToken       `json:"groupIdToken,omitempty"`
	PersonalMessage     *MessageContent     `json:"personalMessage,omitempty"`
}

// NewIdTokenInfo creates an IdTokenInfo. Optional parameters may be set afterwards on the initialized struct.
func NewIdTokenInfo(status AuthorizationStatus) *IdTokenInfo {
	return &IdTokenInfo{Status: status}
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

// Charging Profiles

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
	ChargingProfileKindAbsolute                              ChargingProfileKindType    = "Absolute"
	ChargingProfileKindRecurring                             ChargingProfileKindType    = "Recurring"
	ChargingProfileKindRelative                              ChargingProfileKindType    = "Relative"
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
	purposeType := ChargingProfilePurposeType(fl.Field().String())
	switch purposeType {
	case ChargingProfilePurposeChargingStationExternalConstraints, ChargingProfilePurposeChargingStationMaxProfile, ChargingProfilePurposeTxDefaultProfile, ChargingProfilePurposeTxProfile:
		return true
	default:
		return false
	}
}

func isValidChargingProfileKind(fl validator.FieldLevel) bool {
	purposeType := ChargingProfileKindType(fl.Field().String())
	switch purposeType {
	case ChargingProfileKindAbsolute, ChargingProfileKindRecurring, ChargingProfileKindRelative:
		return true
	default:
		return false
	}
}

func isValidRecurrencyKind(fl validator.FieldLevel) bool {
	purposeType := RecurrencyKindType(fl.Field().String())
	switch purposeType {
	case RecurrencyKindDaily, RecurrencyKindWeekly:
		return true
	default:
		return false
	}
}

func isValidChargingRateUnit(fl validator.FieldLevel) bool {
	purposeType := ChargingRateUnitType(fl.Field().String())
	switch purposeType {
	case ChargingRateUnitWatts, ChargingRateUnitAmperes:
		return true
	default:
		return false
	}
}

func isValidChargingLimitSource(fl validator.FieldLevel) bool {
	chargingLimitSource := ChargingLimitSourceType(fl.Field().String())
	switch chargingLimitSource {
	case ChargingLimitSourceEMS, ChargingLimitSourceOther, ChargingLimitSourceSO, ChargingLimitSourceCSO:
		return true
	default:
		return false
	}
}

type ChargingSchedulePeriod struct {
	StartPeriod  int     `json:"startPeriod" validate:"gte=0"`
	Limit        float64 `json:"limit" validate:"gte=0"`
	NumberPhases *int    `json:"numberPhases,omitempty" validate:"omitempty,gte=0"`
}

func NewChargingSchedulePeriod(startPeriod int, limit float64) ChargingSchedulePeriod {
	return ChargingSchedulePeriod{StartPeriod: startPeriod, Limit: limit}
}

type ChargingSchedule struct {
	StartSchedule          *DateTime                `json:"startSchedule,omitempty" validate:"omitempty"`
	Duration               *int                     `json:"duration,omitempty" validate:"omitempty,gte=0"`
	ChargingRateUnit       ChargingRateUnitType     `json:"chargingRateUnit" validate:"required,chargingRateUnit"`
	MinChargingRate        *float64                 `json:"minChargingRate,omitempty" validate:"omitempty,gte=0"`
	ChargingSchedulePeriod []ChargingSchedulePeriod `json:"chargingSchedulePeriod" validate:"required,min=1"`
}

func NewChargingSchedule(chargingRateUnit ChargingRateUnitType, schedulePeriod ...ChargingSchedulePeriod) *ChargingSchedule {
	return &ChargingSchedule{ChargingRateUnit: chargingRateUnit, ChargingSchedulePeriod: schedulePeriod}
}

type ChargingProfile struct {
	ChargingProfileId      int                        `json:"chargingProfileId" validate:"gte=0"`
	TransactionId          int                        `json:"transactionId,omitempty"`
	StackLevel             int                        `json:"stackLevel" validate:"gt=0"`
	ChargingProfilePurpose ChargingProfilePurposeType `json:"chargingProfilePurpose" validate:"required,chargingProfilePurpose"`
	ChargingProfileKind    ChargingProfileKindType    `json:"chargingProfileKind" validate:"required,chargingProfileKind"`
	RecurrencyKind         RecurrencyKindType         `json:"recurrencyKind,omitempty" validate:"omitempty,recurrencyKind"`
	ValidFrom              *DateTime                  `json:"validFrom,omitempty"`
	ValidTo                *DateTime                  `json:"validTo,omitempty"`
	ChargingSchedule       *ChargingSchedule          `json:"chargingSchedule" validate:"required"`
}

func NewChargingProfile(chargingProfileId int, stackLevel int, chargingProfilePurpose ChargingProfilePurposeType, chargingProfileKind ChargingProfileKindType, schedule *ChargingSchedule) *ChargingProfile {
	return &ChargingProfile{ChargingProfileId: chargingProfileId, StackLevel: stackLevel, ChargingProfilePurpose: chargingProfilePurpose, ChargingProfileKind: chargingProfileKind, ChargingSchedule: schedule}
}

// Remote Start/Stop
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

// Meter Value

type ReadingContext string
type ValueFormat string
type Measurand string
type Phase string
type Location string

const (
	ReadingContextInterruptionBegin       ReadingContext = "Interruption.Begin"
	ReadingContextInterruptionEnd         ReadingContext = "Interruption.End"
	ReadingContextOther                   ReadingContext = "Other"
	ReadingContextSampleClock             ReadingContext = "Sample.Clock"
	ReadingContextSamplePeriodic          ReadingContext = "Sample.Periodic"
	ReadingContextTransactionBegin        ReadingContext = "Transaction.Begin"
	ReadingContextTransactionEnd          ReadingContext = "Transaction.End"
	ReadingContextTrigger                 ReadingContext = "Trigger"
	MeasurandCurrentExport                Measurand      = "Current.Export"
	MeasurandCurrentImport                Measurand      = "Current.Import"
	MeasurandCurrentOffered               Measurand      = "Current.Offered"
	MeasurandEnergyActiveExportRegister   Measurand      = "Energy.Active.Export.Register"
	MeasurandEnergyActiveImportRegister   Measurand      = "Energy.Active.Import.Register"
	MeasurandEnergyReactiveExportRegister Measurand      = "Energy.Reactive.Export.Register"
	MeasurandEnergyReactiveImportRegister Measurand      = "Energy.Reactive.Import.Register"
	MeasurandEnergyActiveExportInterval   Measurand      = "Energy.Active.Export.Interval"
	MeasurandEnergyActiveImportInterval   Measurand      = "Energy.Active.Import.Interval"
	MeasurandEnergyActiveNet              Measurand      = "Energy.Active.Net"
	MeasurandEnergyReactiveExportInterval Measurand      = "Energy.Reactive.Export.Interval"
	MeasurandEnergyReactiveImportInterval Measurand      = "Energy.Reactive.Import.Interval"
	MeasurandEnergyReactiveNet            Measurand      = "Energy.Reactive.Net"
	MeasurandEnergyApparentNet            Measurand      = "Energy.Apparent.Net"
	MeasurandEnergyApparentImport         Measurand      = "Energy.Apparent.Import"
	MeasurandEnergyApparentExport         Measurand      = "Energy.Apparent.Export"
	MeasurandFrequency                    Measurand      = "Frequency"
	MeasurandPowerActiveExport            Measurand      = "Power.Active.Export"
	MeasurandPowerActiveImport            Measurand      = "Power.Active.Import"
	MeasurandPowerFactor                  Measurand      = "Power.Factor"
	MeasurandPowerOffered                 Measurand      = "Power.Offered"
	MeasurandPowerReactiveExport          Measurand      = "Power.Reactive.Export"
	MeasurandPowerReactiveImport          Measurand      = "Power.Reactive.Import"
	MeasueandSoC                          Measurand      = "SoC"
	MeasurandTemperature                  Measurand      = "Temperature"
	MeasurandVoltage                      Measurand      = "Voltage"
	PhaseL1                               Phase          = "L1"
	PhaseL2                               Phase          = "L2"
	PhaseL3                               Phase          = "L3"
	PhaseN                                Phase          = "N"
	PhaseL1N                              Phase          = "L1-N"
	PhaseL2N                              Phase          = "L2-N"
	PhaseL3N                              Phase          = "L3-N"
	PhaseL1L2                             Phase          = "L1-L2"
	PhaseL2L3                             Phase          = "L2-L3"
	PhaseL3L1                             Phase          = "L3-L1"
	LocationBody                          Location       = "Body"
	LocationCable                         Location       = "Cable"
	LocationEV                            Location       = "EV"
	LocationInlet                         Location       = "Inlet"
	LocationOutlet                        Location       = "Outlet"
)

func isValidReadingContext(fl validator.FieldLevel) bool {
	readingContext := ReadingContext(fl.Field().String())
	switch readingContext {
	case ReadingContextInterruptionBegin, ReadingContextInterruptionEnd, ReadingContextOther, ReadingContextSampleClock, ReadingContextSamplePeriodic, ReadingContextTransactionBegin, ReadingContextTransactionEnd, ReadingContextTrigger:
		return true
	default:
		return false
	}
}

func isValidMeasurand(fl validator.FieldLevel) bool {
	measurand := Measurand(fl.Field().String())
	switch measurand {
	case MeasueandSoC, MeasurandCurrentExport, MeasurandCurrentImport, MeasurandCurrentOffered, MeasurandEnergyActiveExportInterval, MeasurandEnergyActiveExportRegister, MeasurandEnergyReactiveExportInterval, MeasurandEnergyReactiveExportRegister, MeasurandEnergyReactiveImportRegister, MeasurandEnergyReactiveImportInterval, MeasurandEnergyActiveImportInterval, MeasurandEnergyActiveImportRegister, MeasurandFrequency, MeasurandPowerActiveExport, MeasurandPowerActiveImport, MeasurandPowerReactiveImport, MeasurandPowerReactiveExport, MeasurandPowerOffered, MeasurandPowerFactor, MeasurandVoltage, MeasurandTemperature, MeasurandEnergyActiveNet, MeasurandEnergyApparentNet, MeasurandEnergyReactiveNet, MeasurandEnergyApparentImport, MeasurandEnergyApparentExport:
		return true
	default:
		return false
	}
}

func isValidPhase(fl validator.FieldLevel) bool {
	phase := Phase(fl.Field().String())
	switch phase {
	case PhaseL1, PhaseL2, PhaseL3, PhaseN, PhaseL1N, PhaseL2N, PhaseL3N, PhaseL1L2, PhaseL2L3, PhaseL3L1:
		return true
	default:
		return false
	}
}

func isValidLocation(fl validator.FieldLevel) bool {
	location := Location(fl.Field().String())
	switch location {
	case LocationBody, LocationCable, LocationEV, LocationInlet, LocationOutlet:
		return true
	default:
		return false
	}
}

type UnitOfMeasure struct {
	Unit       string `json:"unit,omitempty" validate:"omitempty,max=20"`
	Multiplier *int   `json:"multiplier,omitempty" validate:"omitempty,gte=0"`
}

//TODO: remove SignatureMethod (obsolete from 2.0.1 onwards)

// Enumeration of the cryptographic method used to create the digital signature.
// The list is expected to grow in future OCPP releases to allow other signature methods used by Smart Meters.
type SignatureMethod string

const (
	SignatureECDSAP256SHA256 SignatureMethod = "ECDSAP256SHA256" // The encoded data is hashed with the SHA-256 hash function, and the hash value is then signed with the ECDSA algorithm using the NIST P-256 elliptic curve.
	SignatureECDSAP384SHA384 SignatureMethod = "ECDSAP384SHA384" // The encoded data is hashed with the SHA-384 hash function, and the hash value is then signed with the ECDSA algorithm using the NIST P-384 elliptic curve.
	SignatureECDSA192SHA256  SignatureMethod = "ECDSA192SHA256"  // The encoded data is hashed with the SHA-256 hash function, and the hash value is then signed with the ECDSA algorithm using a 192-bit elliptic curve.
)

func isValidSignatureMethod(fl validator.FieldLevel) bool {
	signature := SignatureMethod(fl.Field().String())
	switch signature {
	case SignatureECDSA192SHA256, SignatureECDSAP256SHA256, SignatureECDSAP384SHA384:
		return true
	default:
		return false
	}
}

//TODO: remove EncodingMethod (obsolete from 2.0.1 onwards)

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

type SignedMeterValue struct {
	SignedMeterData string `json:"signedMeterData" validate:"required,max=2500"` // Base64 encoded, contains the signed data which might contain more then just the meter value. It can contain information like timestamps, reference to a customer etc.
	SigningMethod   string `json:"signingMethod" validate:"required,max=50"`     // Method used to create the digital signature.
	EncodingMethod  string `json:"encodingMethod" validate:"required,max=50"`    // Method used to encode the meter values before applying the digital signature algorithm.
	PublicKey       string `json:"publicKey" validate:"required,max=2500"`       // Base64 encoded, sending depends on configuration variable PublicKeyWithSignedMeterValue.
}

type SampledValue struct {
	Value            float64           `json:"value" validate:"required"`                             // Indicates the measured value.
	Context          ReadingContext    `json:"context,omitempty" validate:"omitempty,readingContext"` // Type of detail value: start, end or sample. Default = "Sample.Periodic"
	Measurand        Measurand         `json:"measurand,omitempty" validate:"omitempty,measurand"`    // Type of measurement. Default = "Energy.Active.Import.Register"
	Phase            Phase             `json:"phase,omitempty" validate:"omitempty,phase"`            // Indicates how the measured value is to be interpreted. For instance between L1 and neutral (L1-N) Please note that not all values of phase are applicable to all Measurands. When phase is absent, the measured value is interpreted as an overall value.
	Location         Location          `json:"location,omitempty" validate:"omitempty,location"`      // Indicates where the measured value has been sampled.
	SignedMeterValue *SignedMeterValue `json:"signedMeterValue,omitempty" validate:"omitempty"`       // Contains the MeterValueSignature with sign/encoding method information.
	UnitOfMeasure    *UnitOfMeasure    `json:"unitOfMeasure,omitempty" validate:"omitempty"`          // Represents a UnitOfMeasure including a multiplier.
}

type MeterValue struct {
	Timestamp    DateTime       `json:"timestamp" validate:"required"`
	SampledValue []SampledValue `json:"sampledValue" validate:"required,min=1,dive"`
}

// Validator used for validating all OCPP 2.0 messages.
// Any additional custom validations must be added to this object for automatic validation.
var Validate = ocppj.Validate

func init() {
	_ = Validate.RegisterValidation("idTokenType", isValidIdTokenType)
	_ = Validate.RegisterValidation("genericDeviceModelStatus", isValidGenericDeviceModelStatus)
	_ = Validate.RegisterValidation("genericStatus", isValidGenericStatus)
	_ = Validate.RegisterValidation("hashAlgorithm", isValidHashAlgorithmType)
	_ = Validate.RegisterValidation("certificateStatus", isValidCertificateStatus)
	_ = Validate.RegisterValidation("messageFormat", isValidMessageFormatType)
	_ = Validate.RegisterValidation("authorizationStatus", isValidAuthorizationStatus)
	_ = Validate.RegisterValidation("attribute", isValidAttribute)
	_ = Validate.RegisterValidation("chargingProfilePurpose", isValidChargingProfilePurpose)
	_ = Validate.RegisterValidation("chargingProfileKind", isValidChargingProfileKind)
	_ = Validate.RegisterValidation("recurrencyKind", isValidRecurrencyKind)
	_ = Validate.RegisterValidation("chargingRateUnit", isValidChargingRateUnit)
	_ = Validate.RegisterValidation("chargingLimitSource", isValidChargingLimitSource)
	_ = Validate.RegisterValidation("remoteStartStopStatus", isValidRemoteStartStopStatus)
	_ = Validate.RegisterValidation("readingContext", isValidReadingContext)
	_ = Validate.RegisterValidation("measurand", isValidMeasurand)
	_ = Validate.RegisterValidation("phase", isValidPhase)
	_ = Validate.RegisterValidation("location", isValidLocation)
	_ = Validate.RegisterValidation("signatureMethod", isValidSignatureMethod)
	_ = Validate.RegisterValidation("encodingMethod", isValidEncodingMethod)
	_ = Validate.RegisterValidation("certificateSigningUse", isValidCertificateSigningUse)
	_ = Validate.RegisterValidation("certificateUse", isValidCertificateUse)
	_ = Validate.RegisterValidation("15118EVCertificate", isValidCertificate15118EVStatus)
}
