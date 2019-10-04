package ocpp16_test

import (
	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	"time"
)

// Test
func (suite *OcppV16TestSuite) TestIdTagInfoValidation() {
	var testTable = []GenericTestEntry{
		{ocpp16.IdTagInfo{ExpiryDate: ocpp16.NewDateTime(time.Now()), ParentIdTag: "00000", Status: ocpp16.AuthorizationStatusAccepted}, true},
		{ocpp16.IdTagInfo{ExpiryDate: ocpp16.NewDateTime(time.Now()), Status: ocpp16.AuthorizationStatusAccepted}, true},
		{ocpp16.IdTagInfo{Status: ocpp16.AuthorizationStatusAccepted}, true},
		{ocpp16.IdTagInfo{Status: ocpp16.AuthorizationStatusBlocked}, true},
		{ocpp16.IdTagInfo{Status: ocpp16.AuthorizationStatusExpired}, true},
		{ocpp16.IdTagInfo{Status: ocpp16.AuthorizationStatusInvalid}, true},
		{ocpp16.IdTagInfo{Status: ocpp16.AuthorizationStatusConcurrentTx}, true},
		{ocpp16.IdTagInfo{Status: "invalidAuthorizationStatus"}, false},
		{ocpp16.IdTagInfo{}, false},
		{ocpp16.IdTagInfo{ExpiryDate: ocpp16.NewDateTime(time.Now()), ParentIdTag: ">20..................", Status: ocpp16.AuthorizationStatusAccepted}, false},
	}
	ExecuteGenericTestTable(suite.T(), testTable)
}

func (suite *OcppV16TestSuite) TestChargingSchedulePeriodValidation() {
	t := suite.T()
	var testTable = []GenericTestEntry{
		{ocpp16.ChargingSchedulePeriod{StartPeriod: 0, Limit: 10.0, NumberPhases: 3}, true},
		{ocpp16.ChargingSchedulePeriod{StartPeriod: 0, Limit: 10.0}, true},
		{ocpp16.ChargingSchedulePeriod{StartPeriod: 0}, true},
		{ocpp16.ChargingSchedulePeriod{}, true},
		{ocpp16.ChargingSchedulePeriod{StartPeriod: 0, Limit: -1.0}, false},
		{ocpp16.ChargingSchedulePeriod{StartPeriod: -1, Limit: 10.0}, false},
		{ocpp16.ChargingSchedulePeriod{StartPeriod: 0, Limit: 10.0, NumberPhases: -1}, false},
	}
	ExecuteGenericTestTable(t, testTable)
}

func (suite *OcppV16TestSuite) TestChargingScheduleValidation() {
	t := suite.T()
	chargingSchedulePeriods := make([]ocpp16.ChargingSchedulePeriod, 2)
	chargingSchedulePeriods[0] = ocpp16.NewChargingSchedulePeriod(0, 10.0)
	chargingSchedulePeriods[1] = ocpp16.NewChargingSchedulePeriod(100, 8.0)
	var testTable = []GenericTestEntry{
		{ocpp16.ChargingSchedule{Duration: 0, StartSchedule: ocpp16.NewDateTime(time.Now()), ChargingRateUnit: ocpp16.ChargingRateUnitWatts, ChargingSchedulePeriod: chargingSchedulePeriods, MinChargingRate: 1.0}, true},
		{ocpp16.ChargingSchedule{Duration: 0, ChargingRateUnit: ocpp16.ChargingRateUnitWatts, ChargingSchedulePeriod: chargingSchedulePeriods, MinChargingRate: 1.0}, true},
		{ocpp16.ChargingSchedule{Duration: 0, ChargingRateUnit: ocpp16.ChargingRateUnitWatts, ChargingSchedulePeriod: chargingSchedulePeriods}, true},
		{ocpp16.ChargingSchedule{Duration: 0, ChargingRateUnit: ocpp16.ChargingRateUnitWatts}, false},
		{ocpp16.ChargingSchedule{Duration: 0, ChargingSchedulePeriod: chargingSchedulePeriods}, false},
		{ocpp16.ChargingSchedule{Duration: -1, StartSchedule: ocpp16.NewDateTime(time.Now()), ChargingRateUnit: ocpp16.ChargingRateUnitWatts, ChargingSchedulePeriod: chargingSchedulePeriods, MinChargingRate: 1.0}, false},
		{ocpp16.ChargingSchedule{Duration: 0, StartSchedule: ocpp16.NewDateTime(time.Now()), ChargingRateUnit: ocpp16.ChargingRateUnitWatts, ChargingSchedulePeriod: chargingSchedulePeriods, MinChargingRate: -1.0}, false},
		{ocpp16.ChargingSchedule{Duration: 0, StartSchedule: ocpp16.NewDateTime(time.Now()), ChargingRateUnit: ocpp16.ChargingRateUnitWatts, ChargingSchedulePeriod: make([]ocpp16.ChargingSchedulePeriod, 0), MinChargingRate: 1.0}, false},
		{ocpp16.ChargingSchedule{Duration: -1, StartSchedule: ocpp16.NewDateTime(time.Now()), ChargingRateUnit: "invalidChargeRateUnit", ChargingSchedulePeriod: chargingSchedulePeriods, MinChargingRate: 1.0}, false},
	}
	ExecuteGenericTestTable(t, testTable)
}

func (suite *OcppV16TestSuite) TestChargingProfileValidation() {
	t := suite.T()
	chargingSchedule := ocpp16.NewChargingSchedule(ocpp16.ChargingRateUnitWatts, ocpp16.NewChargingSchedulePeriod(0, 10.0), ocpp16.NewChargingSchedulePeriod(100, 8.0))
	var testTable = []GenericTestEntry{
		{ocpp16.ChargingProfile{ChargingProfileId: 1, TransactionId: 1, StackLevel: 1, ChargingProfilePurpose: ocpp16.ChargingProfilePurposeChargePointMaxProfile, ChargingProfileKind: ocpp16.ChargingProfileKindAbsolute, RecurrencyKind: ocpp16.RecurrencyKindDaily, ValidFrom: ocpp16.NewDateTime(time.Now()), ValidTo: ocpp16.NewDateTime(time.Now().Add(8 * time.Hour)), ChargingSchedule: chargingSchedule}, true},
		{ocpp16.ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfilePurpose: ocpp16.ChargingProfilePurposeChargePointMaxProfile, ChargingProfileKind: ocpp16.ChargingProfileKindAbsolute, ChargingSchedule: chargingSchedule}, true},
		{ocpp16.ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfilePurpose: ocpp16.ChargingProfilePurposeChargePointMaxProfile, ChargingProfileKind: ocpp16.ChargingProfileKindAbsolute}, false},
		{ocpp16.ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfilePurpose: ocpp16.ChargingProfilePurposeChargePointMaxProfile, ChargingSchedule: chargingSchedule}, false},
		{ocpp16.ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfileKind: ocpp16.ChargingProfileKindAbsolute, ChargingSchedule: chargingSchedule}, false},
		{ocpp16.ChargingProfile{ChargingProfileId: 1, ChargingProfilePurpose: ocpp16.ChargingProfilePurposeChargePointMaxProfile, ChargingProfileKind: ocpp16.ChargingProfileKindAbsolute, ChargingSchedule: chargingSchedule}, false},
		{ocpp16.ChargingProfile{StackLevel: 1, ChargingProfilePurpose: ocpp16.ChargingProfilePurposeChargePointMaxProfile, ChargingProfileKind: ocpp16.ChargingProfileKindAbsolute, ChargingSchedule: chargingSchedule}, true},
		{ocpp16.ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfilePurpose: ocpp16.ChargingProfilePurposeChargePointMaxProfile, ChargingProfileKind: "invalidChargingProfileKind", ChargingSchedule: chargingSchedule}, false},
		{ocpp16.ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfilePurpose: "invalidChargingProfilePurpose", ChargingProfileKind: ocpp16.ChargingProfileKindAbsolute, ChargingSchedule: chargingSchedule}, false},
		{ocpp16.ChargingProfile{ChargingProfileId: 1, StackLevel: 0, ChargingProfilePurpose: ocpp16.ChargingProfilePurposeChargePointMaxProfile, ChargingProfileKind: ocpp16.ChargingProfileKindAbsolute, ChargingSchedule: chargingSchedule}, false},
		{ocpp16.ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfilePurpose: ocpp16.ChargingProfilePurposeChargePointMaxProfile, ChargingProfileKind: ocpp16.ChargingProfileKindAbsolute, RecurrencyKind: "invalidRecurrencyKind", ChargingSchedule: chargingSchedule}, false},
		{ocpp16.ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfilePurpose: ocpp16.ChargingProfilePurposeChargePointMaxProfile, ChargingProfileKind: ocpp16.ChargingProfileKindAbsolute, ChargingSchedule: ocpp16.NewChargingSchedule(ocpp16.ChargingRateUnitWatts)}, false},
	}
	ExecuteGenericTestTable(t, testTable)
}

func (suite *OcppV16TestSuite) TestSampledValueValidation() {
	t := suite.T()
	var testTable = []GenericTestEntry{
		{ocpp16.SampledValue{Value: "value", Context: ocpp16.ReadingContextTransactionEnd, Format: ocpp16.ValueFormatRaw, Measurand: ocpp16.MeasurandPowerActiveExport, Phase: ocpp16.PhaseL2, Location: ocpp16.LocationBody, Unit: ocpp16.UnitOfMeasureKW}, true},
		{ocpp16.SampledValue{Value: "value", Context: ocpp16.ReadingContextTransactionEnd, Format: ocpp16.ValueFormatRaw, Measurand: ocpp16.MeasurandPowerActiveExport, Phase: ocpp16.PhaseL2, Location: ocpp16.LocationBody}, true},
		{ocpp16.SampledValue{Value: "value", Context: ocpp16.ReadingContextTransactionEnd, Format: ocpp16.ValueFormatRaw, Measurand: ocpp16.MeasurandPowerActiveExport, Phase: ocpp16.PhaseL2}, true},
		{ocpp16.SampledValue{Value: "value", Context: ocpp16.ReadingContextTransactionEnd, Format: ocpp16.ValueFormatRaw, Measurand: ocpp16.MeasurandPowerActiveExport}, true},
		{ocpp16.SampledValue{Value: "value", Context: ocpp16.ReadingContextTransactionEnd, Format: ocpp16.ValueFormatRaw}, true},
		{ocpp16.SampledValue{Value: "value", Context: ocpp16.ReadingContextTransactionEnd}, true},
		{ocpp16.SampledValue{Value: "value"}, true},
		{ocpp16.SampledValue{Value: "value", Context: "invalidContext"}, false},
		{ocpp16.SampledValue{Value: "value", Format: "invalidFormat"}, false},
		{ocpp16.SampledValue{Value: "value", Measurand: "invalidMeasurand"}, false},
		{ocpp16.SampledValue{Value: "value", Phase: "invalidPhase"}, false},
		{ocpp16.SampledValue{Value: "value", Location: "invalidLocation"}, false},
		{ocpp16.SampledValue{Value: "value", Unit: "invalidUnit"}, false},
	}
	ExecuteGenericTestTable(t, testTable)
}

func (suite *OcppV16TestSuite) TestMeterValueValidation() {
	var testTable = []GenericTestEntry{
		{ocpp16.MeterValue{Timestamp: ocpp16.NewDateTime(time.Now()), SampledValue: []ocpp16.SampledValue{{Value: "value"}, {Value: "value2", Unit: ocpp16.UnitOfMeasureKW}}}, true},
		{ocpp16.MeterValue{Timestamp: ocpp16.NewDateTime(time.Now()), SampledValue: []ocpp16.SampledValue{{Value: "value"}}}, true},
		{ocpp16.MeterValue{Timestamp: ocpp16.NewDateTime(time.Now()), SampledValue: []ocpp16.SampledValue{}}, false},
		{ocpp16.MeterValue{Timestamp: ocpp16.NewDateTime(time.Now())}, false},
		{ocpp16.MeterValue{SampledValue: []ocpp16.SampledValue{{Value: "value"}}}, false},
	}
	ExecuteGenericTestTable(suite.T(), testTable)
}
