package types

import "gopkg.in/go-playground/validator.v9"

// Meter Value

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

	MeasurandCurrentExport        Measurand = "Current.Export"
	MeasurandCurrentImport        Measurand = "Current.Import"
	MeasurandCurrentOffered       Measurand = "Current.Offered"
	MeasurandCurrentExportOffered Measurand = "Current.Export.Offered"
	MeasurandCurrentExportMinimum Measurand = "Current.Export.Minimum"
	MeasurandCurrentImportOffered Measurand = "Current.Import.Offered"
	MeasurandCurrentImportMinimum Measurand = "Current.Import.Minimum"

	MeasurandEnergyActiveExportRegister                Measurand = "Energy.Active.Export.Register"
	MeasurandEnergyActiveImportRegister                Measurand = "Energy.Active.Import.Register"
	MeasurandEnergyReactiveExportRegister              Measurand = "Energy.Reactive.Export.Register"
	MeasurandEnergyReactiveImportRegister              Measurand = "Energy.Reactive.Import.Register"
	MeasurandEnergyActiveExportInterval                Measurand = "Energy.Active.Export.Interval"
	MeasurandEnergyActiveSetpointInterval              Measurand = "Energy.Active.Setpoint.Interval"
	MeasurandEnergyActiveImportInterval                Measurand = "Energy.Active.Import.Interval"
	MeasurandEnergyActiveImportCableLoss               Measurand = "Energy.Active.Import.CableLoss"
	MeasurandEnergyActiveImportLocalGenerationRegister Measurand = "Energy.Active.Import.LocalGeneration.Register"
	MeasurandEnergyActiveNet                           Measurand = "Energy.Active.Net"
	MeasurandEnergyReactiveExportInterval              Measurand = "Energy.Reactive.Export.Interval"
	MeasurandEnergyReactiveImportInterval              Measurand = "Energy.Reactive.Import.Interval"
	MeasurandEnergyReactiveNet                         Measurand = "Energy.Reactive.Net"
	MeasurandEnergyApparentNet                         Measurand = "Energy.Apparent.Net"
	MeasurandEnergyApparentImport                      Measurand = "Energy.Apparent.Import"
	MeasurandEnergyApparentExport                      Measurand = "Energy.Apparent.Export"
	MeasurandEnergyRequestMinimum                      Measurand = "EnergyRequest.Minimum"
	MeasurandEnergyRequestTarget                       Measurand = "EnergyRequest.Target"
	MeasurandEnergyRequestMaximum                      Measurand = "EnergyRequest.Maximum"
	MeasurandEnergyRequestMinimumV2X                   Measurand = "EnergyRequest.Minimum.V2X"
	MeasurandEnergyRequestMaximumV2X                   Measurand = "EnergyRequest.Maximum.V2X"
	MeasurandEnergyRequestBulk                         Measurand = "EnergyRequest.Bulk"

	MeasurandFrequency Measurand = "Frequency"

	MeasurandPowerActiveExport   Measurand = "Power.Active.Export"
	MeasurandPowerActiveImport   Measurand = "Power.Active.Import"
	MeasurandPowerFactor         Measurand = "Power.Factor"
	MeasurandPowerOffered        Measurand = "Power.Offered"
	MeasurandPowerReactiveExport Measurand = "Power.Reactive.Export"
	MeasurandPowerReactiveImport Measurand = "Power.Reactive.Import"
	MeasurandPowerActiveSetpoint Measurand = "Power.Active.Setpoint"
	MeasurandPowerActiveResidual Measurand = "Power.Active.Residual"
	MeasurandPowerExportMinimum  Measurand = "Power.Export.Minimum"
	MeasurandPowerExportOffered  Measurand = "Power.Export.Offered"
	MeasurandPowerImportOffered  Measurand = "Power.Import.Offered"
	MeasurandPowerImportMinimum  Measurand = "Power.Import.Minimum"

	MeasurandSoC                              Measurand = "SoC"
	MeasurandDisplayPresentSoC                          = Measurand("Display.PresentSOC")
	MeasurandDisplayMinimumSoC                          = Measurand("Display.MinimumSOC")
	MeasurandDisplayTargetSoC                           = Measurand("Display.TargetSOC")
	MeasurandDisplayMaximumSoC                          = Measurand("Display.MaximumSOC")
	MeasurandDisplayRemainingTimeToMinimumSoC           = Measurand("Display.RemainingTimeToMinimumSOC")
	MeasurandDisplayRemainingTimeToTargetSoC            = Measurand("Display.RemainingTimeToTargetSOC")
	MeasurandDisplayRemainingTimeToMaximumSoC           = Measurand("Display.RemainingTimeToMaximumSOC")
	MeasurandDisplayChargingComplete                    = Measurand("Display.ChargingComplete")
	MeasurandDisplayBatteryEnergyCapacity               = Measurand("Display.BatteryEnergyCapacity")
	MeasurandDisplayInletHot                            = Measurand("Display.InletHot")

	MeasurandTemperature Measurand = "Temperature"

	MeasurandVoltage        Measurand = "Voltage"
	MeasurandVoltageMinimum Measurand = "Voltage.Minimum"
	MeasurandVoltageMaximum Measurand = "Voltage.Maximum"

	PhaseL1          Phase    = "L1"
	PhaseL2          Phase    = "L2"
	PhaseL3          Phase    = "L3"
	PhaseN           Phase    = "N"
	PhaseL1N         Phase    = "L1-N"
	PhaseL2N         Phase    = "L2-N"
	PhaseL3N         Phase    = "L3-N"
	PhaseL1L2        Phase    = "L1-L2"
	PhaseL2L3        Phase    = "L2-L3"
	PhaseL3L1        Phase    = "L3-L1"
	LocationBody     Location = "Body"
	LocationCable    Location = "Cable"
	LocationEV       Location = "EV"
	LocationInlet    Location = "Inlet"
	LocationOutlet   Location = "Outlet"
	LocationUpstream Location = "Upstream"
)

func isValidReadingContext(fl validator.FieldLevel) bool {
	readingContext := ReadingContext(fl.Field().String())
	switch readingContext {
	case ReadingContextInterruptionBegin, ReadingContextInterruptionEnd,
		ReadingContextOther, ReadingContextSampleClock, ReadingContextSamplePeriodic,
		ReadingContextTransactionBegin, ReadingContextTransactionEnd, ReadingContextTrigger:
		return true
	default:
		return false
	}
}

func isValidMeasurand(fl validator.FieldLevel) bool {
	measurand := Measurand(fl.Field().String())
	switch measurand {
	case MeasurandSoC, MeasurandCurrentExport, MeasurandCurrentImport, MeasurandCurrentOffered, MeasurandEnergyActiveExportInterval,
		MeasurandEnergyActiveExportRegister, MeasurandEnergyReactiveExportInterval, MeasurandEnergyReactiveExportRegister, MeasurandEnergyReactiveImportRegister,
		MeasurandEnergyReactiveImportInterval, MeasurandEnergyActiveImportInterval, MeasurandEnergyActiveImportRegister, MeasurandFrequency, MeasurandPowerActiveExport,
		MeasurandPowerActiveImport, MeasurandPowerReactiveImport, MeasurandPowerReactiveExport, MeasurandPowerOffered, MeasurandPowerFactor, MeasurandVoltage,
		MeasurandTemperature, MeasurandEnergyActiveNet, MeasurandEnergyApparentNet, MeasurandEnergyReactiveNet, MeasurandEnergyApparentImport,
		MeasurandEnergyApparentExport, MeasurandEnergyActiveSetpointInterval, MeasurandEnergyActiveImportCableLoss, MeasurandEnergyActiveImportLocalGenerationRegister,
		MeasurandEnergyRequestMinimum, MeasurandEnergyRequestTarget, MeasurandEnergyRequestMaximum, MeasurandEnergyRequestMinimumV2X, MeasurandPowerActiveSetpoint,
		MeasurandPowerActiveResidual, MeasurandEnergyRequestBulk, MeasurandDisplayPresentSoC, MeasurandDisplayMinimumSoC, MeasurandDisplayTargetSoC,
		MeasurandPowerExportMinimum, MeasurandDisplayMaximumSoC, MeasurandEnergyRequestMaximumV2X, MeasurandVoltageMinimum, MeasurandVoltageMaximum,
		MeasurandCurrentExportOffered, MeasurandPowerExportOffered, MeasurandDisplayRemainingTimeToMinimumSoC, MeasurandDisplayChargingComplete, MeasurandCurrentExportMinimum,
		MeasurandPowerImportOffered, MeasurandDisplayRemainingTimeToTargetSoC, MeasurandDisplayBatteryEnergyCapacity, MeasurandCurrentImportOffered,
		MeasurandPowerImportMinimum, MeasurandDisplayRemainingTimeToMaximumSoC, MeasurandDisplayInletHot, MeasurandCurrentImportMinimum:
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
	case LocationBody, LocationCable, LocationEV, LocationInlet, LocationOutlet, LocationUpstream:
		return true
	default:
		return false
	}
}

type UnitOfMeasure struct {
	Unit       string `json:"unit,omitempty" validate:"omitempty,max=20"`
	Multiplier *int   `json:"multiplier,omitempty" validate:"omitempty,gte=0"`
}

type SignedMeterValue struct {
	SignedMeterData string `json:"signedMeterData" validate:"required,max=32768"`       // Base64 encoded, contains the signed data which might contain more then just the meter value. It can contain information like timestamps, reference to a customer etc.
	SigningMethod   string `json:"signingMethod,omitempty" validate:"omitempty,max=50"` // Method used to create the digital signature.
	EncodingMethod  string `json:"encodingMethod" validate:"required,max=50"`           // Method used to encode the meter values before applying the digital signature algorithm.
	PublicKey       string `json:"publicKey,omitempty" validate:"omitempty,max=2500"`   // Base64 encoded, sending depends on configuration variable PublicKeyWithSignedMeterValue.
}

type SampledValue struct {
	Value            float64           `json:"value"`                                                   // Indicates the measured value. This value is required.
	Context          ReadingContext    `json:"context,omitempty" validate:"omitempty,readingContext21"` // Type of detail value: start, end or sample. Default = "Sample.Periodic"
	Measurand        Measurand         `json:"measurand,omitempty" validate:"omitempty,measurand21"`    // Type of measurement. Default = "Energy.Active.Import.Register"
	Phase            Phase             `json:"phase,omitempty" validate:"omitempty,phase21"`            // Indicates how the measured value is to be interpreted. For instance between L1 and neutral (L1-N) Please note that not all values of phase are applicable to all Measurands. When phase is absent, the measured value is interpreted as an overall value.
	Location         Location          `json:"location,omitempty" validate:"omitempty,location21"`      // Indicates where the measured value has been sampled.
	SignedMeterValue *SignedMeterValue `json:"signedMeterValue,omitempty" validate:"omitempty"`         // Contains the MeterValueSignature with sign/encoding method information.
	UnitOfMeasure    *UnitOfMeasure    `json:"unitOfMeasure,omitempty" validate:"omitempty"`            // Represents a UnitOfMeasure including a multiplier.
}

type MeterValue struct {
	Timestamp    DateTime       `json:"timestamp" validate:"required"`
	SampledValue []SampledValue `json:"sampledValue" validate:"required,min=1,dive"`
}
