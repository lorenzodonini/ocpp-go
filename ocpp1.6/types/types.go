// Contains common and shared data types between OCPP 1.6 messages.
package types

import (
	"github.com/lorenzodonini/ocpp-go/ocppj"
	"gopkg.in/go-playground/validator.v9"
)

const (
	V16Subprotocol = "ocpp1.6"
)

type PropertyViolation struct {
	Property string
}

func (e *PropertyViolation) Error() string {
	return ""
}

type AuthorizationStatus string

const (
	AuthorizationStatusAccepted     AuthorizationStatus = "Accepted"
	AuthorizationStatusBlocked      AuthorizationStatus = "Blocked"
	AuthorizationStatusExpired      AuthorizationStatus = "Expired"
	AuthorizationStatusInvalid      AuthorizationStatus = "Invalid"
	AuthorizationStatusConcurrentTx AuthorizationStatus = "ConcurrentTx"
)

func isValidAuthorizationStatus(fl validator.FieldLevel) bool {
	status := AuthorizationStatus(fl.Field().String())
	switch status {
	case AuthorizationStatusAccepted, AuthorizationStatusBlocked, AuthorizationStatusExpired, AuthorizationStatusInvalid, AuthorizationStatusConcurrentTx:
		return true
	default:
		return false
	}
}

type IdTagInfo struct {
	ExpiryDate  *DateTime           `json:"expiryDate,omitempty" validate:"omitempty"`
	ParentIdTag string              `json:"parentIdTag,omitempty" validate:"omitempty,max=20"`
	Status      AuthorizationStatus `json:"status" validate:"required,authorizationStatus16"`
}

func NewIdTagInfo(status AuthorizationStatus) *IdTagInfo {
	return &IdTagInfo{Status: status}
}

// Charging Profiles
type ChargingProfilePurposeType string
type ChargingProfileKindType string
type RecurrencyKindType string
type ChargingRateUnitType string

const (
	ChargingProfilePurposeChargePointMaxProfile ChargingProfilePurposeType = "ChargePointMaxProfile"
	ChargingProfilePurposeTxDefaultProfile      ChargingProfilePurposeType = "TxDefaultProfile"
	ChargingProfilePurposeTxProfile             ChargingProfilePurposeType = "TxProfile"
	ChargingProfileKindAbsolute                 ChargingProfileKindType    = "Absolute"
	ChargingProfileKindRecurring                ChargingProfileKindType    = "Recurring"
	ChargingProfileKindRelative                 ChargingProfileKindType    = "Relative"
	RecurrencyKindDaily                         RecurrencyKindType         = "Daily"
	RecurrencyKindWeekly                        RecurrencyKindType         = "Weekly"
	ChargingRateUnitWatts                       ChargingRateUnitType       = "W"
	ChargingRateUnitAmperes                     ChargingRateUnitType       = "A"
)

func isValidChargingProfilePurpose(fl validator.FieldLevel) bool {
	purposeType := ChargingProfilePurposeType(fl.Field().String())
	switch purposeType {
	case ChargingProfilePurposeChargePointMaxProfile, ChargingProfilePurposeTxDefaultProfile, ChargingProfilePurposeTxProfile:
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

type ChargingSchedulePeriod struct {
	StartPeriod  int     `json:"startPeriod" validate:"gte=0"`
	Limit        float64 `json:"limit" validate:"gte=0"`
	NumberPhases *int    `json:"numberPhases,omitempty" validate:"omitempty,gte=0"`
}

func NewChargingSchedulePeriod(startPeriod int, limit float64) ChargingSchedulePeriod {
	return ChargingSchedulePeriod{StartPeriod: startPeriod, Limit: limit}
}

type ChargingSchedule struct {
	Duration               *int                     `json:"duration,omitempty" validate:"omitempty,gte=0"`
	StartSchedule          *DateTime                `json:"startSchedule,omitempty"`
	ChargingRateUnit       ChargingRateUnitType     `json:"chargingRateUnit" validate:"required,chargingRateUnit16"`
	ChargingSchedulePeriod []ChargingSchedulePeriod `json:"chargingSchedulePeriod" validate:"required,min=1"`
	MinChargingRate        *float64                 `json:"minChargingRate,omitempty" validate:"omitempty,gte=0"`
}

func NewChargingSchedule(chargingRateUnit ChargingRateUnitType, schedulePeriod ...ChargingSchedulePeriod) *ChargingSchedule {
	return &ChargingSchedule{ChargingRateUnit: chargingRateUnit, ChargingSchedulePeriod: schedulePeriod}
}

type ChargingProfile struct {
	ChargingProfileId      int                        `json:"chargingProfileId"`
	TransactionId          int                        `json:"transactionId,omitempty"`
	StackLevel             int                        `json:"stackLevel" validate:"gte=0"`
	ChargingProfilePurpose ChargingProfilePurposeType `json:"chargingProfilePurpose" validate:"required,chargingProfilePurpose16"`
	ChargingProfileKind    ChargingProfileKindType    `json:"chargingProfileKind" validate:"required,chargingProfileKind16"`
	RecurrencyKind         RecurrencyKindType         `json:"recurrencyKind,omitempty" validate:"omitempty,recurrencyKind16"`
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
	MeasurandSoC                          Measurand      = "SoC"
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
	UnitOfMeasureCelcius                  UnitOfMeasure  = "Celcius"
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
	case MeasurandSoC, MeasurandCurrentExport, MeasurandCurrentImport, MeasurandCurrentOffered, MeasurandEnergyActiveExportInterval, MeasurandEnergyActiveExportRegister, MeasurandEnergyReactiveExportInterval, MeasurandEnergyReactiveExportRegister, MeasurandEnergyReactiveImportRegister, MeasurandEnergyReactiveImportInterval, MeasurandEnergyActiveImportInterval, MeasurandEnergyActiveImportRegister, MeasurandFrequency, MeasurandPowerActiveExport, MeasurandPowerActiveImport, MeasurandPowerReactiveImport, MeasurandPowerReactiveExport, MeasurandPowerOffered, MeasurandPowerFactor, MeasurandVoltage, MeasurandTemperature, MeasurandRPM:
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
	case UnitOfMeasureA, UnitOfMeasureWh, UnitOfMeasureKWh, UnitOfMeasureVarh, UnitOfMeasureKvarh, UnitOfMeasureW, UnitOfMeasureKW, UnitOfMeasureVA, UnitOfMeasureKVA, UnitOfMeasureVar, UnitOfMeasureKvar, UnitOfMeasureV, UnitOfMeasureCelsius, UnitOfMeasureCelcius, UnitOfMeasureFahrenheit, UnitOfMeasureK, UnitOfMeasurePercent:
		return true
	default:
		return false
	}
}

type SampledValue struct {
	Value     string         `json:"value" validate:"required"`
	Context   ReadingContext `json:"context,omitempty" validate:"omitempty,readingContext16"`
	Format    ValueFormat    `json:"format,omitempty" validate:"omitempty,valueFormat"`
	Measurand Measurand      `json:"measurand,omitempty" validate:"omitempty,measurand16"`
	Phase     Phase          `json:"phase,omitempty" validate:"omitempty,phase16"`
	Location  Location       `json:"location,omitempty" validate:"omitempty,location16"`
	Unit      UnitOfMeasure  `json:"unit,omitempty" validate:"omitempty,unitOfMeasure"`
}

type MeterValue struct {
	Timestamp    *DateTime      `json:"timestamp" validate:"required"`
	SampledValue []SampledValue `json:"sampledValue" validate:"required,min=1,dive"`
}

// Initialize validator
var Validate = ocppj.Validate

func init() {
	_ = Validate.RegisterValidation("authorizationStatus16", isValidAuthorizationStatus)
	_ = Validate.RegisterValidation("chargingProfilePurpose16", isValidChargingProfilePurpose)
	_ = Validate.RegisterValidation("chargingProfileKind16", isValidChargingProfileKind)
	_ = Validate.RegisterValidation("recurrencyKind16", isValidRecurrencyKind)
	_ = Validate.RegisterValidation("chargingRateUnit16", isValidChargingRateUnit)
	_ = Validate.RegisterValidation("remoteStartStopStatus16", isValidRemoteStartStopStatus)
	_ = Validate.RegisterValidation("readingContext16", isValidReadingContext)
	_ = Validate.RegisterValidation("valueFormat", isValidValueFormat)
	_ = Validate.RegisterValidation("measurand16", isValidMeasurand)
	_ = Validate.RegisterValidation("phase16", isValidPhase)
	_ = Validate.RegisterValidation("location16", isValidLocation)
	_ = Validate.RegisterValidation("unitOfMeasure", isValidUnitOfMeasure)
	_ = Validate.RegisterValidation("certificateSigningUse16", isValidCertificateSigningUse)
	_ = Validate.RegisterValidation("certificateUse16", isValidCertificateUse)
	_ = Validate.RegisterValidation("genericStatus16", isValidGenericStatus)
}
