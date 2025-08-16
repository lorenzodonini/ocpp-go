package types

import (
	"gopkg.in/go-playground/validator.v9"
)

// Charging Profiles
type ChargingProfilePurposeType string
type ChargingProfileKindType string
type RecurrencyKindType string
type ChargingRateUnitType string
type ChargingLimitSourceType string

type OperationMode string

const (
	ChargingProfilePurposeChargingStationExternalConstraints ChargingProfilePurposeType = "ChargingStationExternalConstraints"
	ChargingProfilePurposeChargingStationMaxProfile          ChargingProfilePurposeType = "ChargingStationMaxProfile"
	ChargingProfilePurposeTxDefaultProfile                   ChargingProfilePurposeType = "TxDefaultProfile"
	ChargingProfilePurposeTxProfile                          ChargingProfilePurposeType = "TxProfile"

	ChargingProfileKindAbsolute  ChargingProfileKindType = "Absolute"
	ChargingProfileKindRecurring ChargingProfileKindType = "Recurring"
	ChargingProfileKindRelative  ChargingProfileKindType = "Relative"

	RecurrencyKindDaily  RecurrencyKindType = "Daily"
	RecurrencyKindWeekly RecurrencyKindType = "Weekly"

	ChargingRateUnitWatts   ChargingRateUnitType = "W"
	ChargingRateUnitAmperes ChargingRateUnitType = "A"

	ChargingLimitSourceEMS   ChargingLimitSourceType = "EMS"
	ChargingLimitSourceOther ChargingLimitSourceType = "Other"
	ChargingLimitSourceSO    ChargingLimitSourceType = "SO"
	ChargingLimitSourceCSO   ChargingLimitSourceType = "CSO"

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
	operationMode := OperationMode(fl.Field().String())
	switch operationMode {
	case OperationModeIdle, OperationModeChargingOnly, OperationModeCentralSetpoint, OperationModeExternalSetpoint,
		OperationModeExternalLimits, OperationModeCentralFrequency, OperationModeLocalFrequency,
		OperationModeLocalLoadBalancing:
		return true
	default:
		return false
	}
}

func isValidChargingProfilePurpose(fl validator.FieldLevel) bool {
	purposeType := ChargingProfilePurposeType(fl.Field().String())
	switch purposeType {
	case ChargingProfilePurposeChargingStationExternalConstraints, ChargingProfilePurposeChargingStationMaxProfile,
		ChargingProfilePurposeTxDefaultProfile, ChargingProfilePurposeTxProfile:
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
	StartPeriod        int                  `json:"startPeriod" validate:"gte=0"`
	NumberPhases       *int                 `json:"numberPhases,omitempty" validate:"omitempty,gte=0,lte=3"`
	Limit              float64              `json:"limit" validate:"gte=0"`
	LimitL2            *float64             `json:"limit_L2,omitempty" validate:"omitempty,gte=0"`
	LimitL3            *float64             `json:"limit_L3,omitempty" validate:"omitempty,gte=0"`
	PhaseToUse         *int                 `json:"phaseToUse,omitempty" validate:"omitempty,gte=0,lte=3"`
	DischargeLimit     *float64             `json:"dischargeLimit,omitempty" validate:"omitempty,lte=0"`
	DischargeLimitL2   *float64             `json:"dischargeLimit_L2,omitempty" validate:"omitempty,lte=0"`
	DischargeLimitL3   *float64             `json:"dischargeLimit_L3,omitempty" validate:"omitempty,lte=0"`
	SetPoint           *float64             `json:"setpoint,omitempty" validate:"omitempty"`
	SetPointL2         *float64             `json:"setpoint_L2,omitempty" validate:"omitempty"`
	SetPointL3         *float64             `json:"setpoint_L3,omitempty" validate:"omitempty"`
	SetpointReactive   *float64             `json:"setpointReactive,omitempty" validate:"omitempty"`
	SetpointReactiveL2 *float64             `json:"setpointReactive_L2,omitempty" validate:"omitempty"`
	SetpointReactiveL3 *float64             `json:"setpointReactive_L3,omitempty" validate:"omitempty"`
	OperationMode      OperationMode        `json:"operationMode,omitempty" validate:"omitempty,operationMode21"`
	V2xFreqWattCurve   []V2xFreqWattCurve   `json:"v2xFreqWattCurve,omitempty" validate:"omitempty,max=20,dive"`
	V2xSignalWattPoint []V2XSignalWattPoint `json:"v2xSignalWattPoint,omitempty" validate:"omitempty,max=20,dive"`
}

type V2xFreqWattCurve struct {
	Frequency float64 `json:"frequency" validate:"gte=0"`
	Power     float64 `json:"power" validate:"gte=0"` // Power in Watts, positive for export, negative for import.
}

type V2XSignalWattPoint struct {
	Signal int     `json:"signal" validate:"required"`
	Power  float64 `json:"power" validate:"required"`
}

func NewChargingSchedulePeriod(startPeriod int, limit float64) ChargingSchedulePeriod {
	return ChargingSchedulePeriod{StartPeriod: startPeriod, Limit: limit}
}

type ChargingSchedule struct {
	ID                     int                      `json:"id" validate:"gte=0"` // Identifies the ChargingSchedule.
	StartSchedule          *DateTime                `json:"startSchedule,omitempty" validate:"omitempty"`
	Duration               *int                     `json:"duration,omitempty" validate:"omitempty,gte=0"`
	ChargingRateUnit       ChargingRateUnitType     `json:"chargingRateUnit" validate:"required,chargingRateUnit21"`
	MinChargingRate        *float64                 `json:"minChargingRate,omitempty" validate:"omitempty,gte=0"`
	ChargingSchedulePeriod []ChargingSchedulePeriod `json:"chargingSchedulePeriod" validate:"required,min=1,max=1024"`
	SalesTariff            *SalesTariff             `json:"salesTariff,omitempty" validate:"omitempty"` // Sales tariff associated with this charging schedule.
	PowerTolerance         *float64                 `json:"powerTolerance,omitempty" validate:"omitempty"`
	SignatureId            *int                     `json:"signatureId,omitempty" validate:"omitempty,gte=0"`
	DigestValue            *string                  `json:"digestValue,omitempty" validate:"omitempty,max=88"`
	UseLocalTime           bool                     `json:"useLocalTime,omitempty" validate:"omitempty"`
	RandomizedDelay        *int                     `json:"randomizedDelay,omitempty" validate:"omitempty,gte=0"`
}

func NewChargingSchedule(id int, chargingRateUnit ChargingRateUnitType, schedulePeriod ...ChargingSchedulePeriod) *ChargingSchedule {
	return &ChargingSchedule{ID: id, ChargingRateUnit: chargingRateUnit, ChargingSchedulePeriod: schedulePeriod}
}

type ChargingProfile struct {
	ID                          int                        `json:"id" validate:"gte=0"`
	StackLevel                  int                        `json:"stackLevel" validate:"gte=0"`
	ChargingProfilePurpose      ChargingProfilePurposeType `json:"chargingProfilePurpose" validate:"required,chargingProfilePurpose21"`
	ChargingProfileKind         ChargingProfileKindType    `json:"chargingProfileKind" validate:"required,chargingProfileKind21"`
	RecurrencyKind              RecurrencyKindType         `json:"recurrencyKind,omitempty" validate:"omitempty,recurrencyKind21"`
	ValidFrom                   *DateTime                  `json:"validFrom,omitempty"`
	ValidTo                     *DateTime                  `json:"validTo,omitempty"`
	TransactionID               string                     `json:"transactionId,omitempty" validate:"omitempty,max=36"`
	MaxOfflineDuration          *int                       `json:"maxOfflineDuration,omitempty" validate:"omitempty"`
	InvalidAfterOfflineDuration bool                       `json:"invalidAfterOfflineDuration,omitempty" validate:"omitempty"`
	DynUpdateInterval           *int                       `json:"dynUpdateInterval,omitempty" validate:"omitempty"`
	DynUpdateTime               *DateTime                  `json:"dynUpdateTime,omitempty" validate:"omitempty"`
	PriceScheduleSignature      *string                    `json:"priceScheduleSignature,omitempty" validate:"omitempty,max=256"`
	ChargingSchedule            []ChargingSchedule         `json:"chargingSchedule" validate:"required,min=1,max=3,dive"`
}

func NewChargingProfile(id int, stackLevel int, chargingProfilePurpose ChargingProfilePurposeType, chargingProfileKind ChargingProfileKindType, schedule []ChargingSchedule) *ChargingProfile {
	return &ChargingProfile{ID: id, StackLevel: stackLevel, ChargingProfilePurpose: chargingProfilePurpose, ChargingProfileKind: chargingProfileKind, ChargingSchedule: schedule}
}

type EnergyTransferMode string

const (
	EnergyTransferModeDC       EnergyTransferMode = "DC"              // DC charging.
	EnergyTransferModeAC1Phase EnergyTransferMode = "AC_single_phase" // AC single phase charging according to IEC 62196.
	EnergyTransferModeAC2Phase EnergyTransferMode = "AC_two_phase"    // AC two phase charging according to IEC 62196.
	EnergyTransferModeAC3Phase EnergyTransferMode = "AC_three_phase"  // AC three phase charging according to IEC 62196.
)
