package ocpp2_test

import (
	"github.com/lorenzodonini/ocpp-go/ocpp2.0"
	"time"
)

// Test
func (suite *OcppV16TestSuite) TestIdTagInfoValidation() {
	var testTable = []GenericTestEntry{
		{ocpp2.IdTagInfo{ExpiryDate: ocpp2.NewDateTime(time.Now()), ParentIdTag: "00000", Status: ocpp2.AuthorizationStatusAccepted}, true},
		{ocpp2.IdTagInfo{ExpiryDate: ocpp2.NewDateTime(time.Now()), Status: ocpp2.AuthorizationStatusAccepted}, true},
		{ocpp2.IdTagInfo{Status: ocpp2.AuthorizationStatusAccepted}, true},
		{ocpp2.IdTagInfo{Status: ocpp2.AuthorizationStatusBlocked}, true},
		{ocpp2.IdTagInfo{Status: ocpp2.AuthorizationStatusExpired}, true},
		{ocpp2.IdTagInfo{Status: ocpp2.AuthorizationStatusInvalid}, true},
		{ocpp2.IdTagInfo{Status: ocpp2.AuthorizationStatusConcurrentTx}, true},
		{ocpp2.IdTagInfo{Status: "invalidAuthorizationStatus"}, false},
		{ocpp2.IdTagInfo{}, false},
		{ocpp2.IdTagInfo{ExpiryDate: ocpp2.NewDateTime(time.Now()), ParentIdTag: ">20..................", Status: ocpp2.AuthorizationStatusAccepted}, false},
	}
	ExecuteGenericTestTable(suite.T(), testTable)
}

func (suite *OcppV16TestSuite) TestChargingSchedulePeriodValidation() {
	t := suite.T()
	var testTable = []GenericTestEntry{
		{ocpp2.ChargingSchedulePeriod{StartPeriod: 0, Limit: 10.0, NumberPhases: 3}, true},
		{ocpp2.ChargingSchedulePeriod{StartPeriod: 0, Limit: 10.0}, true},
		{ocpp2.ChargingSchedulePeriod{StartPeriod: 0}, true},
		{ocpp2.ChargingSchedulePeriod{}, true},
		{ocpp2.ChargingSchedulePeriod{StartPeriod: 0, Limit: -1.0}, false},
		{ocpp2.ChargingSchedulePeriod{StartPeriod: -1, Limit: 10.0}, false},
		{ocpp2.ChargingSchedulePeriod{StartPeriod: 0, Limit: 10.0, NumberPhases: -1}, false},
	}
	ExecuteGenericTestTable(t, testTable)
}

func (suite *OcppV16TestSuite) TestChargingScheduleValidation() {
	t := suite.T()
	chargingSchedulePeriods := make([]ocpp2.ChargingSchedulePeriod, 2)
	chargingSchedulePeriods[0] = ocpp2.NewChargingSchedulePeriod(0, 10.0)
	chargingSchedulePeriods[1] = ocpp2.NewChargingSchedulePeriod(100, 8.0)
	var testTable = []GenericTestEntry{
		{ocpp2.ChargingSchedule{Duration: 0, StartSchedule: ocpp2.NewDateTime(time.Now()), ChargingRateUnit: ocpp2.ChargingRateUnitWatts, ChargingSchedulePeriod: chargingSchedulePeriods, MinChargingRate: 1.0}, true},
		{ocpp2.ChargingSchedule{Duration: 0, ChargingRateUnit: ocpp2.ChargingRateUnitWatts, ChargingSchedulePeriod: chargingSchedulePeriods, MinChargingRate: 1.0}, true},
		{ocpp2.ChargingSchedule{Duration: 0, ChargingRateUnit: ocpp2.ChargingRateUnitWatts, ChargingSchedulePeriod: chargingSchedulePeriods}, true},
		{ocpp2.ChargingSchedule{Duration: 0, ChargingRateUnit: ocpp2.ChargingRateUnitWatts}, false},
		{ocpp2.ChargingSchedule{Duration: 0, ChargingSchedulePeriod: chargingSchedulePeriods}, false},
		{ocpp2.ChargingSchedule{Duration: -1, StartSchedule: ocpp2.NewDateTime(time.Now()), ChargingRateUnit: ocpp2.ChargingRateUnitWatts, ChargingSchedulePeriod: chargingSchedulePeriods, MinChargingRate: 1.0}, false},
		{ocpp2.ChargingSchedule{Duration: 0, StartSchedule: ocpp2.NewDateTime(time.Now()), ChargingRateUnit: ocpp2.ChargingRateUnitWatts, ChargingSchedulePeriod: chargingSchedulePeriods, MinChargingRate: -1.0}, false},
		{ocpp2.ChargingSchedule{Duration: 0, StartSchedule: ocpp2.NewDateTime(time.Now()), ChargingRateUnit: ocpp2.ChargingRateUnitWatts, ChargingSchedulePeriod: make([]ocpp2.ChargingSchedulePeriod, 0), MinChargingRate: 1.0}, false},
		{ocpp2.ChargingSchedule{Duration: -1, StartSchedule: ocpp2.NewDateTime(time.Now()), ChargingRateUnit: "invalidChargeRateUnit", ChargingSchedulePeriod: chargingSchedulePeriods, MinChargingRate: 1.0}, false},
	}
	ExecuteGenericTestTable(t, testTable)
}

func (suite *OcppV16TestSuite) TestChargingProfileValidation() {
	t := suite.T()
	chargingSchedule := ocpp2.NewChargingSchedule(ocpp2.ChargingRateUnitWatts, ocpp2.NewChargingSchedulePeriod(0, 10.0), ocpp2.NewChargingSchedulePeriod(100, 8.0))
	var testTable = []GenericTestEntry{
		{ocpp2.ChargingProfile{ChargingProfileId: 1, TransactionId: 1, StackLevel: 1, ChargingProfilePurpose: ocpp2.ChargingProfilePurposeChargePointMaxProfile, ChargingProfileKind: ocpp2.ChargingProfileKindAbsolute, RecurrencyKind: ocpp2.RecurrencyKindDaily, ValidFrom: ocpp2.NewDateTime(time.Now()), ValidTo: ocpp2.NewDateTime(time.Now().Add(8 * time.Hour)), ChargingSchedule: chargingSchedule}, true},
		{ocpp2.ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfilePurpose: ocpp2.ChargingProfilePurposeChargePointMaxProfile, ChargingProfileKind: ocpp2.ChargingProfileKindAbsolute, ChargingSchedule: chargingSchedule}, true},
		{ocpp2.ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfilePurpose: ocpp2.ChargingProfilePurposeChargePointMaxProfile, ChargingProfileKind: ocpp2.ChargingProfileKindAbsolute}, false},
		{ocpp2.ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfilePurpose: ocpp2.ChargingProfilePurposeChargePointMaxProfile, ChargingSchedule: chargingSchedule}, false},
		{ocpp2.ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfileKind: ocpp2.ChargingProfileKindAbsolute, ChargingSchedule: chargingSchedule}, false},
		{ocpp2.ChargingProfile{ChargingProfileId: 1, ChargingProfilePurpose: ocpp2.ChargingProfilePurposeChargePointMaxProfile, ChargingProfileKind: ocpp2.ChargingProfileKindAbsolute, ChargingSchedule: chargingSchedule}, false},
		{ocpp2.ChargingProfile{StackLevel: 1, ChargingProfilePurpose: ocpp2.ChargingProfilePurposeChargePointMaxProfile, ChargingProfileKind: ocpp2.ChargingProfileKindAbsolute, ChargingSchedule: chargingSchedule}, true},
		{ocpp2.ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfilePurpose: ocpp2.ChargingProfilePurposeChargePointMaxProfile, ChargingProfileKind: "invalidChargingProfileKind", ChargingSchedule: chargingSchedule}, false},
		{ocpp2.ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfilePurpose: "invalidChargingProfilePurpose", ChargingProfileKind: ocpp2.ChargingProfileKindAbsolute, ChargingSchedule: chargingSchedule}, false},
		{ocpp2.ChargingProfile{ChargingProfileId: 1, StackLevel: 0, ChargingProfilePurpose: ocpp2.ChargingProfilePurposeChargePointMaxProfile, ChargingProfileKind: ocpp2.ChargingProfileKindAbsolute, ChargingSchedule: chargingSchedule}, false},
		{ocpp2.ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfilePurpose: ocpp2.ChargingProfilePurposeChargePointMaxProfile, ChargingProfileKind: ocpp2.ChargingProfileKindAbsolute, RecurrencyKind: "invalidRecurrencyKind", ChargingSchedule: chargingSchedule}, false},
		{ocpp2.ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfilePurpose: ocpp2.ChargingProfilePurposeChargePointMaxProfile, ChargingProfileKind: ocpp2.ChargingProfileKindAbsolute, ChargingSchedule: ocpp2.NewChargingSchedule(ocpp2.ChargingRateUnitWatts)}, false},
	}
	ExecuteGenericTestTable(t, testTable)
}

func (suite *OcppV16TestSuite) TestSampledValueValidation() {
	t := suite.T()
	var testTable = []GenericTestEntry{
		{ocpp2.SampledValue{Value: "value", Context: ocpp2.ReadingContextTransactionEnd, Format: ocpp2.ValueFormatRaw, Measurand: ocpp2.MeasurandPowerActiveExport, Phase: ocpp2.PhaseL2, Location: ocpp2.LocationBody, Unit: ocpp2.UnitOfMeasureKW}, true},
		{ocpp2.SampledValue{Value: "value", Context: ocpp2.ReadingContextTransactionEnd, Format: ocpp2.ValueFormatRaw, Measurand: ocpp2.MeasurandPowerActiveExport, Phase: ocpp2.PhaseL2, Location: ocpp2.LocationBody}, true},
		{ocpp2.SampledValue{Value: "value", Context: ocpp2.ReadingContextTransactionEnd, Format: ocpp2.ValueFormatRaw, Measurand: ocpp2.MeasurandPowerActiveExport, Phase: ocpp2.PhaseL2}, true},
		{ocpp2.SampledValue{Value: "value", Context: ocpp2.ReadingContextTransactionEnd, Format: ocpp2.ValueFormatRaw, Measurand: ocpp2.MeasurandPowerActiveExport}, true},
		{ocpp2.SampledValue{Value: "value", Context: ocpp2.ReadingContextTransactionEnd, Format: ocpp2.ValueFormatRaw}, true},
		{ocpp2.SampledValue{Value: "value", Context: ocpp2.ReadingContextTransactionEnd}, true},
		{ocpp2.SampledValue{Value: "value"}, true},
		{ocpp2.SampledValue{Value: "value", Context: "invalidContext"}, false},
		{ocpp2.SampledValue{Value: "value", Format: "invalidFormat"}, false},
		{ocpp2.SampledValue{Value: "value", Measurand: "invalidMeasurand"}, false},
		{ocpp2.SampledValue{Value: "value", Phase: "invalidPhase"}, false},
		{ocpp2.SampledValue{Value: "value", Location: "invalidLocation"}, false},
		{ocpp2.SampledValue{Value: "value", Unit: "invalidUnit"}, false},
	}
	ExecuteGenericTestTable(t, testTable)
}

func (suite *OcppV16TestSuite) TestMeterValueValidation() {
	var testTable = []GenericTestEntry{
		{ocpp2.MeterValue{Timestamp: ocpp2.NewDateTime(time.Now()), SampledValue: []ocpp2.SampledValue{{Value: "value"}, {Value: "value2", Unit: ocpp2.UnitOfMeasureKW}}}, true},
		{ocpp2.MeterValue{Timestamp: ocpp2.NewDateTime(time.Now()), SampledValue: []ocpp2.SampledValue{{Value: "value"}}}, true},
		{ocpp2.MeterValue{Timestamp: ocpp2.NewDateTime(time.Now()), SampledValue: []ocpp2.SampledValue{}}, false},
		{ocpp2.MeterValue{Timestamp: ocpp2.NewDateTime(time.Now())}, false},
		{ocpp2.MeterValue{SampledValue: []ocpp2.SampledValue{{Value: "value"}}}, false},
	}
	ExecuteGenericTestTable(suite.T(), testTable)
}
