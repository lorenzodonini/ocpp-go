package der

import (
	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
	"github.com/lorenzodonini/ocpp-go/ocppj"
	"gopkg.in/go-playground/validator.v9"
)

type DERControlStatus string

const (
	DERControlStatusAccepted     = "Accepted"
	DERControlStatusRejected     = "Rejected"
	DERControlStatusNotSupported = "NotSupported"
	DERControlStatusNotFound     = "NotFound"
)

func isValidDERControlStatus(level validator.FieldLevel) bool {
	switch DERControlStatus(level.Field().String()) {
	case DERControlStatusAccepted, DERControlStatusRejected, DERControlStatusNotSupported, DERControlStatusNotFound:
		return true
	default:
		return false
	}
}

type DERControl string

const (
	DERControlEnterService            = "EnterService"
	DERControlFreqDroop               = "FreqDroop"
	DERControlFreqWatt                = "FreqWatt"
	DERControlFixedPFAbsorb           = "FixedPFAbsorb"
	DERControlFixedPFInject           = "FixedPFInject"
	DERControlFixedVar                = "FixedVar"
	DERControlGradients               = "Gradients"
	DERControlHFMustTrip              = "HFMustTrip"
	DERControlHFMayTrip               = "HFMayTrip"
	DERControlHVMustTrip              = "HVMustTrip"
	DERControlHVMomCess               = "HVMomCess"
	DERControlHVMayTrip               = "HVMayTrip"
	DERControlLimitMaxDischarge       = "LimitMaxDischarge"
	DERControlLFMustTrip              = "LFMustTrip"
	DERControlLVMustTrip              = "LVMustTrip"
	DERControlLVMomCess               = "LVMomCess"
	DERControlLVMayTrip               = "LVMayTrip"
	DERControlPowerMonitoringMustTrip = "PowerMonitoringMustTrip"
	DERControlVoltVar                 = "VoltVar"
	DERControlVoltWatt                = "VoltWatt"
	DERControlWattPF                  = "WattPF"
	DERControlWattVar                 = "WattVar"
)

func isValidDERControl(level validator.FieldLevel) bool {
	switch DERControl(level.Field().String()) {
	case DERControlEnterService, DERControlFreqDroop, DERControlFreqWatt, DERControlFixedPFAbsorb, DERControlFixedPFInject,
		DERControlFixedVar, DERControlGradients, DERControlHFMustTrip, DERControlHFMayTrip, DERControlHVMustTrip,
		DERControlHVMomCess, DERControlHVMayTrip, DERControlLimitMaxDischarge, DERControlLFMustTrip, DERControlLVMustTrip,
		DERControlLVMomCess, DERControlLVMayTrip, DERControlPowerMonitoringMustTrip, DERControlVoltVar,
		DERControlVoltWatt, DERControlWattPF, DERControlWattVar:
		return true
	default:
		return false
	}
}

type DERUnit string

const (
	DERUnitNotApplicable = "Not_Applicable"
	DERUnitPctMaxW       = "PctMaxW"
	DERUnitPctMaxVar     = "PctMaxVar"
	DERUnitPctWAvail     = "PctWAvail"
	DERUnitPctVarAvail   = "PctVarAvail"
	DERUnitPctEffectiveV = "PctEffectiveV"
)

func isValidDERUnit(level validator.FieldLevel) bool {
	switch DERUnit(level.Field().String()) {
	case DERUnitNotApplicable, DERUnitPctMaxW, DERUnitPctMaxVar, DERUnitPctWAvail, DERUnitPctVarAvail, DERUnitPctEffectiveV:
		return true
	default:
		return false
	}
}

func init() {
	_ = ocppj.Validate.RegisterValidation("derUnit", isValidDERUnit)
	_ = ocppj.Validate.RegisterValidation("derControlStatus", isValidDERControlStatus)
	_ = ocppj.Validate.RegisterValidation("derControl", isValidDERControl)
	_ = ocppj.Validate.RegisterValidation("powerDuringCessation", isValidPowerDuringCessation)
}

type DERCurve struct {
	Priority            int                  `json:"priority" validate:"required,gte=0"`
	YUnit               DERUnit              `json:"yUnit" validate:"required,derUnit"`
	ResponseTime        *float64             `json:"responseTime,omitempty" validate:"omitempty"`
	StartTime           *types.DateTime      `json:"startTime,omitempty" validate:"omitempty"`
	Duration            float64              `json:"duration,omitempty" validate:"omitempty"`
	Hysteresis          *Hysteresis          `json:"hysteresis,omitempty" validate:"omitempty,dive"`
	VoltageParams       *VoltageParams       `json:"voltageParams,omitempty" validate:"omitempty,dive"`
	ReactivePowerParams *ReactivePowerParams `json:"reactivePowerParams,omitempty" validate:"omitempty,dive"`
	CurveData           []DERCurvePoints     `json:"curveData" validate:"required,gte=1,max=10,dive"`
}

type DERCurvePoints struct {
	X float64 `json:"x" validate:"required"` // X value of the curve point, e.g., frequency or voltage.
	Y float64 `json:"y" validate:"required"` // Y value of the curve point, e.g., active or reactive power.
}

type DERCurveGet struct {
	Id           string   `json:"id" validate:"required,max=36"`
	IsDefault    bool     `json:"isDefault" validate:"required"`
	IsSuperseded bool     `json:"isSuperseded" validate:"required"`
	DERCurve     DERCurve `json:"derCurve" validate:"required,dive"`
}

type Hysteresis struct {
	HysteresisHigh     *float64 `json:"hysteresisHigh,omitempty" validate:"omitempty"`
	HysteresisLow      *float64 `json:"hysteresisLow,omitempty" validate:"omitempty"`
	HysteresisDelay    *float64 `json:"hysteresisDelay,omitempty" validate:"omitempty"`
	HysteresisGradient *float64 `json:"hysteresisGradient,omitempty" validate:"omitempty"`
}

type PowerDuringCessation string

const (
	PowerDuringCessationActive   PowerDuringCessation = "Active"
	PowerDuringCessationReactive PowerDuringCessation = "Reactive"
)

func isValidPowerDuringCessation(level validator.FieldLevel) bool {
	switch PowerDuringCessation(level.Field().String()) {
	case PowerDuringCessationActive, PowerDuringCessationReactive:
		return true
	default:
		return false
	}
}

type VoltageParams struct {
	Hv10MinMeanValue     *float64              `json:"hv10MinMeanValue,omitempty" validate:"omitempty"`
	Hv10MinMeanTripDelay *float64              `json:"hv10MinMeanTripDelay,omitempty" validate:"omitempty"`
	PowerDuringCessation *PowerDuringCessation `json:"powerDuring,omitempty" validate:"omitempty,powerDuringCessation"`
}

type ReactivePowerParams struct {
	VRef                       *float64 `json:"vRef,omitempty" validate:"omitempty"`
	AutonomousVRefEnable       *bool    `json:"autonomousVRefEnable,omitempty" validate:"omitempty"`
	AutonomousVRefTimeConstant *float64 `json:"autonomousVRefTimeConstant,omitempty" validate:"omitempty"`
}

type Gradient struct {
	Priority     int     `json:"priority" validate:"required,gte=0"`
	Gradient     float64 `json:"gradient" validate:"required"`
	SoftGradient float64 `json:"softGradient" validate:"required"`
}

type GradientGet struct {
	Id           string   `json:"id" validate:"required,max=36"`
	IsDefault    bool     `json:"isDefault" validate:"required"`
	IsSuperseded bool     `json:"isSuperseded" validate:"required"`
	Gradient     Gradient `json:"gradient" validate:"required,dive"`
}

type FreqDroop struct {
	Priority       int      `json:"priority" validate:"required,gte=0"`
	OverFrequency  float64  `json:"overFreq" validate:"required"`
	UnderFrequency float64  `json:"underFreq" validate:"required"`
	OverDroop      float64  `json:"overDroop" validate:"required"`
	UnderDroop     float64  `json:"underDroop" validate:"required"`
	ResponseTime   *float64 `json:"responseTime,omitempty" validate:"omitempty"`
	Duration       *float64 `json:"duration,omitempty" validate:"omitempty"`
}

type FreqDroopGet struct {
	Id           string    `json:"id" validate:"required,max=36"`
	IsDefault    bool      `json:"isDefault" validate:"required"`
	IsSuperseded bool      `json:"isSuperseded" validate:"required"`
	FreqDroop    FreqDroop `json:"freqDroop" validate:"required,dive"`
}

type LimitMaxDischarge struct {
	Priority                int             `json:"priority" validate:"required,gte=0"`
	PctMaxDischargePower    float64         `json:"pctMaxDischargePower,omitempty" validate:"omitempty"`         // Percentage of maximum discharge power.
	StartTime               *types.DateTime `json:"startTime,omitempty" validate:"omitempty"`                    // The time at which the limit starts.
	Duration                *float64        `json:"duration,omitempty" validate:"omitempty"`                     // Duration of the limit in seconds.
	PowerMonitoringMustTrip *DERCurve       `json:"powerMonitoringMustTrip,omitempty" validate:"omitempty,dive"` // Optional DER curve for power monitoring must trip.
}
type LimitMaxDischargeGet struct {
	Id                string            `json:"id" validate:"required,max=36"`
	IsDefault         bool              `json:"isDefault" validate:"required"`
	IsSuperseded      bool              `json:"isSuperseded" validate:"required"`
	LimitMaxDischarge LimitMaxDischarge `json:"limitMaxDischarge" validate:"required,dive"`
}

type EnterService struct {
	Priority      int      `json:"priority" validate:"required,gte=0"`
	HighVoltage   float64  `json:"highVoltage" validate:"required"`
	LowVoltage    float64  `json:"lowVoltage" validate:"required"`
	HighFrequency float64  `json:"highFreq" validate:"required"`
	LowFrequency  float64  `json:"lowFreq" validate:"required"`
	Delay         *float64 `json:"delay,omitempty" validate:"omitempty"`
	RandomDelay   *float64 `json:"randomDelay,omitempty" validate:"omitempty"`
	RampRate      *float64 `json:"rampRate,omitempty" validate:"omitempty"`
}
type EnterServiceGet struct {
	Id           string       `json:"id" validate:"required,max=36"`
	IsDefault    bool         `json:"isDefault" validate:"required"`
	IsSuperseded bool         `json:"isSuperseded" validate:"required"`
	EnterService EnterService `json:"enterService" validate:"required,dive"`
}

type FixedPF struct {
	Priority     int             `json:"priority" validate:"required,gte=0"`
	Displacement float64         `json:"displacement" validate:"required"`
	Excitation   bool            `json:"excitation" validate:"required"`
	StartTime    *types.DateTime `json:"startTime,omitempty" validate:"omitempty"`
	Duration     *float64        `json:"duration,omitempty" validate:"omitempty"`
}

type FixedPFGet struct {
	Id           string  `json:"id" validate:"required,max=36"`
	IsDefault    bool    `json:"isDefault" validate:"required"`
	IsSuperseded bool    `json:"isSuperseded" validate:"required"`
	FixedPF      FixedPF `json:"fixedPF" validate:"required,dive"`
}

type FixedVar struct {
	Priority  int             `json:"priority" validate:"required,gte=0"`
	Setpoint  float64         `json:"setpoint" validate:"required"`
	Unit      DERUnit         `json:"unit" validate:"required,derUnit"`
	StartTime *types.DateTime `json:"startTime,omitempty" validate:"omitempty"`
	Duration  *float64        `json:"duration,omitempty" validate:"omitempty"`
}

type FixedVarGet struct {
	Id           string   `json:"id" validate:"required,max=36"`
	IsDefault    bool     `json:"isDefault" validate:"required"`
	IsSuperseded bool     `json:"isSuperseded" validate:"required"`
	FixedVar     FixedVar `json:"fixedVar" validate:"required,dive"`
}

type GridEventFault string

const (
	GridEventFaultOverVoltage      GridEventFault = "OverVoltage"
	GridEventFaultUnderVoltage     GridEventFault = "UnderVoltage"
	GridEventFaultOverFrequency    GridEventFault = "OverFrequency"
	GridEventFaultUnderFrequency   GridEventFault = "UnderFrequency"
	GridEventFaultVoltageImbalance GridEventFault = "VoltageImbalance"
	GridEventFaultLowInputPower    GridEventFault = "LowInputPower"
	GridEventFaultOverCurrent      GridEventFault = "OverCurrent"
	GridEventFaultPhaseRotation    GridEventFault = "PhaseRotation"
	GridEventFaultRemoteEmergency  GridEventFault = "RemoteEmergency"
	GridEventFaultCurrentImbalance GridEventFault = "CurrentImbalance"
)

func isValidGridEventFault(level validator.FieldLevel) bool {
	switch GridEventFault(level.Field().String()) {
	case GridEventFaultOverVoltage,
		GridEventFaultUnderVoltage,
		GridEventFaultOverFrequency,
		GridEventFaultUnderFrequency,
		GridEventFaultVoltageImbalance,
		GridEventFaultLowInputPower,
		GridEventFaultOverCurrent,
		GridEventFaultPhaseRotation,
		GridEventFaultRemoteEmergency,
		GridEventFaultCurrentImbalance:
		return true
	default:
		return false
	}
}

func init() {
	// Register the custom validation function for GridEventFault
	_ = ocppj.Validate.RegisterValidation("gridEventFault", isValidGridEventFault)
}
