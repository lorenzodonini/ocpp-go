package ocpp16_test

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"time"
)

// Utility functions
func newInt(i int) *int {
	return &i
}

func newFloat(f float64) *float64 {
	return &f
}

// Test
func (suite *OcppV16TestSuite) TestIdTagInfoValidation() {
	var testTable = []GenericTestEntry{
		{types.IdTagInfo{ExpiryDate: types.NewDateTime(time.Now()), ParentIdTag: "00000", Status: types.AuthorizationStatusAccepted}, true},
		{types.IdTagInfo{ExpiryDate: types.NewDateTime(time.Now()), Status: types.AuthorizationStatusAccepted}, true},
		{types.IdTagInfo{Status: types.AuthorizationStatusAccepted}, true},
		{types.IdTagInfo{Status: types.AuthorizationStatusBlocked}, true},
		{types.IdTagInfo{Status: types.AuthorizationStatusExpired}, true},
		{types.IdTagInfo{Status: types.AuthorizationStatusInvalid}, true},
		{types.IdTagInfo{Status: types.AuthorizationStatusConcurrentTx}, true},
		{types.IdTagInfo{Status: "invalidAuthorizationStatus"}, false},
		{types.IdTagInfo{}, false},
		{types.IdTagInfo{ExpiryDate: types.NewDateTime(time.Now()), ParentIdTag: ">20..................", Status: types.AuthorizationStatusAccepted}, false},
	}
	ExecuteGenericTestTable(suite.T(), testTable)
}

func (suite *OcppV16TestSuite) TestChargingSchedulePeriodValidation() {
	t := suite.T()
	var testTable = []GenericTestEntry{
		{types.ChargingSchedulePeriod{StartPeriod: 0, Limit: 10.0, NumberPhases: newInt(3)}, true},
		{types.ChargingSchedulePeriod{StartPeriod: 0, Limit: 10.0}, true},
		{types.ChargingSchedulePeriod{StartPeriod: 0}, true},
		{types.ChargingSchedulePeriod{}, true},
		{types.ChargingSchedulePeriod{StartPeriod: 0, Limit: -1.0}, false},
		{types.ChargingSchedulePeriod{StartPeriod: -1, Limit: 10.0}, false},
		{types.ChargingSchedulePeriod{StartPeriod: 0, Limit: 10.0, NumberPhases: newInt(-1)}, false},
	}
	ExecuteGenericTestTable(t, testTable)
}

func (suite *OcppV16TestSuite) TestChargingScheduleValidation() {
	t := suite.T()
	chargingSchedulePeriods := make([]types.ChargingSchedulePeriod, 2)
	chargingSchedulePeriods[0] = types.NewChargingSchedulePeriod(0, 10.0)
	chargingSchedulePeriods[1] = types.NewChargingSchedulePeriod(100, 8.0)
	var testTable = []GenericTestEntry{
		{types.ChargingSchedule{Duration: newInt(0), StartSchedule: types.NewDateTime(time.Now()), ChargingRateUnit: types.ChargingRateUnitWatts, ChargingSchedulePeriod: chargingSchedulePeriods, MinChargingRate: newFloat(1.0)}, true},
		{types.ChargingSchedule{Duration: newInt(0), ChargingRateUnit: types.ChargingRateUnitWatts, ChargingSchedulePeriod: chargingSchedulePeriods, MinChargingRate: newFloat(1.0)}, true},
		{types.ChargingSchedule{Duration: newInt(0), ChargingRateUnit: types.ChargingRateUnitWatts, ChargingSchedulePeriod: chargingSchedulePeriods}, true},
		{types.ChargingSchedule{Duration: newInt(0), ChargingRateUnit: types.ChargingRateUnitWatts}, false},
		{types.ChargingSchedule{Duration: newInt(0), ChargingSchedulePeriod: chargingSchedulePeriods}, false},
		{types.ChargingSchedule{Duration: newInt(-1), StartSchedule: types.NewDateTime(time.Now()), ChargingRateUnit: types.ChargingRateUnitWatts, ChargingSchedulePeriod: chargingSchedulePeriods, MinChargingRate: newFloat(1.0)}, false},
		{types.ChargingSchedule{Duration: newInt(0), StartSchedule: types.NewDateTime(time.Now()), ChargingRateUnit: types.ChargingRateUnitWatts, ChargingSchedulePeriod: chargingSchedulePeriods, MinChargingRate: newFloat(-1.0)}, false},
		{types.ChargingSchedule{Duration: newInt(0), StartSchedule: types.NewDateTime(time.Now()), ChargingRateUnit: types.ChargingRateUnitWatts, ChargingSchedulePeriod: make([]types.ChargingSchedulePeriod, 0), MinChargingRate: newFloat(1.0)}, false},
		{types.ChargingSchedule{Duration: newInt(0), StartSchedule: types.NewDateTime(time.Now()), ChargingRateUnit: "invalidChargeRateUnit", ChargingSchedulePeriod: chargingSchedulePeriods, MinChargingRate: newFloat(1.0)}, false},
	}
	ExecuteGenericTestTable(t, testTable)
}

func (suite *OcppV16TestSuite) TestChargingProfileValidation() {
	t := suite.T()
	chargingSchedule := types.NewChargingSchedule(types.ChargingRateUnitWatts, types.NewChargingSchedulePeriod(0, 10.0), types.NewChargingSchedulePeriod(100, 8.0))
	var testTable = []GenericTestEntry{
		{types.ChargingProfile{ChargingProfileId: 1, TransactionId: 1, StackLevel: 1, ChargingProfilePurpose: types.ChargingProfilePurposeChargePointMaxProfile, ChargingProfileKind: types.ChargingProfileKindAbsolute, RecurrencyKind: types.RecurrencyKindDaily, ValidFrom: types.NewDateTime(time.Now()), ValidTo: types.NewDateTime(time.Now().Add(8 * time.Hour)), ChargingSchedule: chargingSchedule}, true},
		{types.ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfilePurpose: types.ChargingProfilePurposeChargePointMaxProfile, ChargingProfileKind: types.ChargingProfileKindAbsolute, ChargingSchedule: chargingSchedule}, true},
		{types.ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfilePurpose: types.ChargingProfilePurposeChargePointMaxProfile, ChargingProfileKind: types.ChargingProfileKindAbsolute}, false},
		{types.ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfilePurpose: types.ChargingProfilePurposeChargePointMaxProfile, ChargingSchedule: chargingSchedule}, false},
		{types.ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfileKind: types.ChargingProfileKindAbsolute, ChargingSchedule: chargingSchedule}, false},
		{types.ChargingProfile{ChargingProfileId: 1, ChargingProfilePurpose: types.ChargingProfilePurposeChargePointMaxProfile, ChargingProfileKind: types.ChargingProfileKindAbsolute, ChargingSchedule: chargingSchedule}, true},
		{types.ChargingProfile{StackLevel: 1, ChargingProfilePurpose: types.ChargingProfilePurposeChargePointMaxProfile, ChargingProfileKind: types.ChargingProfileKindAbsolute, ChargingSchedule: chargingSchedule}, true},
		{types.ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfilePurpose: types.ChargingProfilePurposeChargePointMaxProfile, ChargingProfileKind: "invalidChargingProfileKind", ChargingSchedule: chargingSchedule}, false},
		{types.ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfilePurpose: "invalidChargingProfilePurpose", ChargingProfileKind: types.ChargingProfileKindAbsolute, ChargingSchedule: chargingSchedule}, false},
		{types.ChargingProfile{ChargingProfileId: 1, StackLevel: 0, ChargingProfilePurpose: types.ChargingProfilePurposeChargePointMaxProfile, ChargingProfileKind: types.ChargingProfileKindAbsolute, ChargingSchedule: chargingSchedule}, true},
		{types.ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfilePurpose: types.ChargingProfilePurposeChargePointMaxProfile, ChargingProfileKind: types.ChargingProfileKindAbsolute, RecurrencyKind: "invalidRecurrencyKind", ChargingSchedule: chargingSchedule}, false},
		{types.ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfilePurpose: types.ChargingProfilePurposeChargePointMaxProfile, ChargingProfileKind: types.ChargingProfileKindAbsolute, ChargingSchedule: types.NewChargingSchedule(types.ChargingRateUnitWatts)}, false},
	}
	ExecuteGenericTestTable(t, testTable)
}

func (suite *OcppV16TestSuite) TestSampledValueValidation() {
	t := suite.T()
	var testTable = []GenericTestEntry{
		{types.SampledValue{Value: "value", Context: types.ReadingContextTransactionEnd, Format: types.ValueFormatRaw, Measurand: types.MeasurandPowerActiveExport, Phase: types.PhaseL2, Location: types.LocationBody, Unit: types.UnitOfMeasureKW}, true},
		{types.SampledValue{Value: "value", Context: types.ReadingContextTransactionEnd, Format: types.ValueFormatRaw, Measurand: types.MeasurandPowerActiveExport, Phase: types.PhaseL2, Location: types.LocationBody}, true},
		{types.SampledValue{Value: "value", Context: types.ReadingContextTransactionEnd, Format: types.ValueFormatRaw, Measurand: types.MeasurandPowerActiveExport, Phase: types.PhaseL2}, true},
		{types.SampledValue{Value: "value", Context: types.ReadingContextTransactionEnd, Format: types.ValueFormatRaw, Measurand: types.MeasurandPowerActiveExport}, true},
		{types.SampledValue{Value: "value", Context: types.ReadingContextTransactionEnd, Format: types.ValueFormatRaw}, true},
		{types.SampledValue{Value: "value", Context: types.ReadingContextTransactionEnd}, true},
		{types.SampledValue{Value: "value"}, true},
		{types.SampledValue{Value: "value", Context: "invalidContext"}, false},
		{types.SampledValue{Value: "value", Format: "invalidFormat"}, false},
		{types.SampledValue{Value: "value", Measurand: "invalidMeasurand"}, false},
		{types.SampledValue{Value: "value", Phase: "invalidPhase"}, false},
		{types.SampledValue{Value: "value", Location: "invalidLocation"}, false},
		{types.SampledValue{Value: "value", Unit: "invalidUnit"}, false},
	}
	ExecuteGenericTestTable(t, testTable)
}

func (suite *OcppV16TestSuite) TestMeterValueValidation() {
	var testTable = []GenericTestEntry{
		{types.MeterValue{Timestamp: types.NewDateTime(time.Now()), SampledValue: []types.SampledValue{{Value: "value"}, {Value: "value2", Unit: types.UnitOfMeasureKW}}}, true},
		{types.MeterValue{Timestamp: types.NewDateTime(time.Now()), SampledValue: []types.SampledValue{{Value: "value"}}}, true},
		{types.MeterValue{Timestamp: types.NewDateTime(time.Now()), SampledValue: []types.SampledValue{}}, false},
		{types.MeterValue{Timestamp: types.NewDateTime(time.Now())}, false},
		{types.MeterValue{SampledValue: []types.SampledValue{{Value: "value"}}}, false},
	}
	ExecuteGenericTestTable(suite.T(), testTable)
}
