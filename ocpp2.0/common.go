package ocpp2

import (
	"encoding/json"
	"github.com/lorenzodonini/ocpp-go/ocppj"
	"gopkg.in/go-playground/validator.v9"
	"strings"
	"time"
)

const (
	V2Subprotocol = "ocpp2.0"
)

type DateTime struct {
	time.Time
}

func NewDateTime(time time.Time) *DateTime {
	return &DateTime{Time: time}
}

var DateTimeFormat = time.RFC3339

func (dt *DateTime) UnmarshalJSON(input []byte) error {
	strInput := string(input)
	strInput = strings.Trim(strInput, `"`)
	if DateTimeFormat == "" {
		defaultTime := time.Time{}
		err := json.Unmarshal(input, defaultTime)
		if err != nil {
			return err
		}
		dt.Time = defaultTime.Local()
	} else {
		newTime, err := time.Parse(DateTimeFormat, strInput)
		if err != nil {
			return err
		}
		dt.Time = newTime
	}
	return nil
}

func (dt *DateTime) MarshalJSON() ([]byte, error) {
	if DateTimeFormat == "" {
		return json.Marshal(dt.Time)
	}
	timeStr := FormatTimestamp(dt.Time)
	return json.Marshal(timeStr)
}

func FormatTimestamp(t time.Time) string {
	return t.UTC().Format(DateTimeFormat)
}

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

func NewIdTokenInfo(status AuthorizationStatus) *IdTokenInfo {
	return &IdTokenInfo{Status: status}
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
type UnitOfMeasure string

const (
	ReadingContextInterruptionBegin       ReadingContext = "Interruption.Begin"
	ReadingContextInterruptionEnd         ReadingContext = "Interruption.End"
	ReadingContextOther                   ReadingContext = "Other"
	ReadingContextSampleClock             ReadingContext = "Sample.Clock"
	ReadingContextSamplePeriodic          ReadingContext = "Sample.Periodic"
	ReadingContextTransactionBegin        ReadingContext = "Transaction.Begin"
	ReadingContextTransactionEnd          ReadingContext = "Transaction.End"
	ReadingContextTrigger                 ReadingContext = "Trigger"
	ValueFormatRaw                        ValueFormat    = "Raw"
	ValueFormatSignedData                 ValueFormat    = "SignedData"
	MeasurandCurrentExport                Measurand      = "Current.Export"
	MeasurandCurrentImport                Measurand      = "Current.Import"
	MeasurandCurrentOffered               Measurand      = "Current.Offered"
	MeasurandEnergyActiveExportRegister   Measurand      = "Energy.Active.Export.Register"
	MeasurandEnergyActiveImportRegister   Measurand      = "Energy.Active.Import.Register"
	MeasurandEnergyReactiveExportRegister Measurand      = "Energy.Reactive.Export.Register"
	MeasurandEnergyReactiveImportRegister Measurand      = "Energy.Reactive.Import.Register"
	MeasurandEnergyActiveExportInterval   Measurand      = "Energy.Active.Export.Interval"
	MeasurandEnergyActiveImportInterval   Measurand      = "Energy.Active.Import.Interval"
	MeasurandEnergyReactiveExportInterval Measurand      = "Energy.Reactive.Export.Interval"
	MeasurandEnergyReactiveImportInterval Measurand      = "Energy.Reactive.Import.Interval"
	MeasurandFrequency                    Measurand      = "Frequency"
	MeasurandPowerActiveExport            Measurand      = "Power.Active.Export"
	MeasurandPowerActiveImport            Measurand      = "Power.Active.Import"
	MeasurandPowerFactor                  Measurand      = "Power.Factor"
	MeasurandPowerOffered                 Measurand      = "Power.Offered"
	MeasurandPowerReactiveExport          Measurand      = "Power.Reactive.Export"
	MeasurandPowerReactiveImport          Measurand      = "Power.Reactive.Import"
	MeasurandRPM                          Measurand      = "RPM"
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
	UnitOfMeasureWh                       UnitOfMeasure  = "Wh"
	UnitOfMeasureKWh                      UnitOfMeasure  = "kWh"
	UnitOfMeasureVarh                     UnitOfMeasure  = "varh"
	UnitOfMeasureKvarh                    UnitOfMeasure  = "kvarh"
	UnitOfMeasureW                        UnitOfMeasure  = "W"
	UnitOfMeasureKW                       UnitOfMeasure  = "kW"
	UnitOfMeasureVA                       UnitOfMeasure  = "VA"
	UnitOfMeasureKVA                      UnitOfMeasure  = "kVA"
	UnitOfMeasureVar                      UnitOfMeasure  = "var"
	UnitOfMeasureKvar                     UnitOfMeasure  = "kvar"
	UnitOfMeasureA                        UnitOfMeasure  = "A"
	UnitOfMeasureV                        UnitOfMeasure  = "V"
	UnitOfMeasureCelsius                  UnitOfMeasure  = "Celsius"
	UnitOfMeasureFahrenheit               UnitOfMeasure  = "Fahrenheit"
	UnitOfMeasureK                        UnitOfMeasure  = "K"
	UnitOfMeasurePercent                  UnitOfMeasure  = "Percent"
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

func isValidValueFormat(fl validator.FieldLevel) bool {
	valueFormat := ValueFormat(fl.Field().String())
	switch valueFormat {
	case ValueFormatRaw, ValueFormatSignedData:
		return true
	default:
		return false
	}
}

func isValidMeasurand(fl validator.FieldLevel) bool {
	measurand := Measurand(fl.Field().String())
	switch measurand {
	case MeasueandSoC, MeasurandCurrentExport, MeasurandCurrentImport, MeasurandCurrentOffered, MeasurandEnergyActiveExportInterval, MeasurandEnergyActiveExportRegister, MeasurandEnergyReactiveExportInterval, MeasurandEnergyReactiveExportRegister, MeasurandEnergyReactiveImportRegister, MeasurandEnergyReactiveImportInterval, MeasurandEnergyActiveImportInterval, MeasurandEnergyActiveImportRegister, MeasurandFrequency, MeasurandPowerActiveExport, MeasurandPowerActiveImport, MeasurandPowerReactiveImport, MeasurandPowerReactiveExport, MeasurandPowerOffered, MeasurandPowerFactor, MeasurandVoltage, MeasurandTemperature, MeasurandRPM:
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

func isValidUnitOfMeasure(fl validator.FieldLevel) bool {
	unitOfMeasure := UnitOfMeasure(fl.Field().String())
	switch unitOfMeasure {
	case UnitOfMeasureA, UnitOfMeasureWh, UnitOfMeasureKWh, UnitOfMeasureVarh, UnitOfMeasureKvarh, UnitOfMeasureW, UnitOfMeasureKW, UnitOfMeasureVA, UnitOfMeasureKVA, UnitOfMeasureVar, UnitOfMeasureKvar, UnitOfMeasureV, UnitOfMeasureCelsius, UnitOfMeasureFahrenheit, UnitOfMeasureK, UnitOfMeasurePercent:
		return true
	default:
		return false
	}
}

type SampledValue struct {
	Value     string         `json:"value" validate:"required"`
	Context   ReadingContext `json:"context,omitempty" validate:"omitempty,readingContext"`
	Format    ValueFormat    `json:"format,omitempty" validate:"omitempty,valueFormat"`
	Measurand Measurand      `json:"measurand,omitempty" validate:"omitempty,measurand"`
	Phase     Phase          `json:"phase,omitempty" validate:"omitempty,phase"`
	Location  Location       `json:"location,omitempty" validate:"omitempty,location"`
	Unit      UnitOfMeasure  `json:"unit,omitempty" validate:"omitempty,unitOfMeasure"`
}

type MeterValue struct {
	Timestamp    *DateTime      `json:"timestamp" validate:"required"`
	SampledValue []SampledValue `json:"sampledValue" validate:"required,min=1,dive"`
}

// DateTime Validation
func dateTimeIsNull(dateTime *DateTime) bool {
	return dateTime != nil && dateTime.IsZero()
}

func validateDateTimeGt(dateTime *DateTime, than time.Time) bool {
	return dateTime != nil && dateTime.After(than)
}

func validateDateTimeNow(dateTime DateTime) bool {
	dur := time.Now().Sub(dateTime.Time).Minutes()
	return dur < 1
}

func validateDateTimeLt(dateTime DateTime, than time.Time) bool {
	return dateTime.Before(than)
}

// Initialize validator
var Validate = ocppj.Validate

func init() {
	_ = Validate.RegisterValidation("idTokenType", isValidIdTokenType)
	_ = Validate.RegisterValidation("genericDeviceModelStatus", isValidGenericDeviceModelStatus)
	_ = Validate.RegisterValidation("genericStatus", isValidGenericStatus)
	_ = Validate.RegisterValidation("hashAlgorithm", isValidHashAlgorithmType)
	_ = Validate.RegisterValidation("certificateStatus", isValidCertificateStatus)
	_ = Validate.RegisterValidation("messageFormat", isValidMessageFormatType)
	_ = Validate.RegisterValidation("authorizationStatus", isValidAuthorizationStatus)
	_ = Validate.RegisterValidation("chargingProfilePurpose", isValidChargingProfilePurpose)
	_ = Validate.RegisterValidation("chargingProfileKind", isValidChargingProfileKind)
	_ = Validate.RegisterValidation("recurrencyKind", isValidRecurrencyKind)
	_ = Validate.RegisterValidation("chargingRateUnit", isValidChargingRateUnit)
	_ = Validate.RegisterValidation("chargingLimitSource", isValidChargingLimitSource)
	_ = Validate.RegisterValidation("remoteStartStopStatus", isValidRemoteStartStopStatus)
	_ = Validate.RegisterValidation("readingContext", isValidReadingContext)
	_ = Validate.RegisterValidation("valueFormat", isValidValueFormat)
	_ = Validate.RegisterValidation("measurand", isValidMeasurand)
	_ = Validate.RegisterValidation("phase", isValidPhase)
	_ = Validate.RegisterValidation("location", isValidLocation)
	_ = Validate.RegisterValidation("unitOfMeasure", isValidUnitOfMeasure)
	_ = Validate.RegisterValidation("certificateSigningUse", isValidCertificateSigningUse)
	_ = Validate.RegisterValidation("15118EVCertificate", isValidCertificate15118EVStatus)
}
