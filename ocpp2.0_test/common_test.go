package ocpp2_test

import (
	"github.com/lorenzodonini/ocpp-go/ocpp2.0"
	"time"
)

// Utility functions
func newInt(i int) *int{
	return &i
}

// Test
func (suite *OcppV2TestSuite) TestIdTokenInfoValidation() {
	var testTable = []GenericTestEntry{
		{ocpp2.IdTokenInfo{Status: ocpp2.AuthorizationStatusAccepted, CacheExpiryDateTime: ocpp2.NewDateTime(time.Now()), ChargingPriority: 1, Language1: "l1", Language2: "l2", GroupIdToken: &ocpp2.GroupIdToken{IdToken: "1234", Type: ocpp2.IdTokenTypeCentral}, PersonalMessage: &ocpp2.MessageContent{Format: ocpp2.MessageFormatUTF8, Language: "en", Content: "random"}}, true},
		{ocpp2.IdTokenInfo{Status: ocpp2.AuthorizationStatusAccepted, CacheExpiryDateTime: ocpp2.NewDateTime(time.Now()), ChargingPriority: 1, Language1: "l1", Language2: "l2", GroupIdToken: &ocpp2.GroupIdToken{IdToken: "1234", Type: ocpp2.IdTokenTypeCentral}, PersonalMessage: &ocpp2.MessageContent{Format: ocpp2.MessageFormatUTF8, Content: "random"}}, true},
		{ocpp2.IdTokenInfo{Status: ocpp2.AuthorizationStatusAccepted, CacheExpiryDateTime: ocpp2.NewDateTime(time.Now()), ChargingPriority: 1, Language1: "l1", Language2: "l2", GroupIdToken: &ocpp2.GroupIdToken{IdToken: "1234", Type: ocpp2.IdTokenTypeCentral}}, true},
		{ocpp2.IdTokenInfo{Status: ocpp2.AuthorizationStatusAccepted, CacheExpiryDateTime: ocpp2.NewDateTime(time.Now()), ChargingPriority: 1, Language1: "l1", Language2: "l2"}, true},
		{ocpp2.IdTokenInfo{Status: ocpp2.AuthorizationStatusAccepted, CacheExpiryDateTime: ocpp2.NewDateTime(time.Now()), ChargingPriority: 1, Language1: "l1"}, true},
		{ocpp2.IdTokenInfo{Status: ocpp2.AuthorizationStatusAccepted, CacheExpiryDateTime: ocpp2.NewDateTime(time.Now()), ChargingPriority: 1}, true},
		{ocpp2.IdTokenInfo{Status: ocpp2.AuthorizationStatusAccepted, CacheExpiryDateTime: ocpp2.NewDateTime(time.Now())}, true},
		{ocpp2.IdTokenInfo{Status: ocpp2.AuthorizationStatusAccepted}, true},
		{ocpp2.IdTokenInfo{}, false},
		{ocpp2.IdTokenInfo{Status: ocpp2.AuthorizationStatusAccepted, CacheExpiryDateTime: ocpp2.NewDateTime(time.Now()), ChargingPriority: 1, Language1: "l1", Language2: "l2", GroupIdToken: &ocpp2.GroupIdToken{IdToken: "1234", Type: ocpp2.IdTokenTypeCentral}, PersonalMessage: &ocpp2.MessageContent{Format: "invalidFormat", Language: "en", Content: "random"}}, false},
		{ocpp2.IdTokenInfo{Status: ocpp2.AuthorizationStatusAccepted, CacheExpiryDateTime: ocpp2.NewDateTime(time.Now()), ChargingPriority: 1, Language1: "l1", Language2: "l2", GroupIdToken: &ocpp2.GroupIdToken{IdToken: "1234", Type: ocpp2.IdTokenTypeCentral}, PersonalMessage: &ocpp2.MessageContent{Format: ocpp2.MessageFormatUTF8, Language: "en", Content: ">512............................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................."}}, false},
		{ocpp2.IdTokenInfo{Status: ocpp2.AuthorizationStatusAccepted, CacheExpiryDateTime: ocpp2.NewDateTime(time.Now()), ChargingPriority: 1, Language1: "l1", Language2: "l2", GroupIdToken: &ocpp2.GroupIdToken{IdToken: "1234", Type: ocpp2.IdTokenTypeCentral}, PersonalMessage: &ocpp2.MessageContent{Format: ocpp2.MessageFormatUTF8, Language: "en"}}, false},
		{ocpp2.IdTokenInfo{Status: ocpp2.AuthorizationStatusAccepted, CacheExpiryDateTime: ocpp2.NewDateTime(time.Now()), ChargingPriority: 1, Language1: "l1", Language2: "l2", GroupIdToken: &ocpp2.GroupIdToken{IdToken: "1234", Type: ocpp2.IdTokenTypeCentral}, PersonalMessage: &ocpp2.MessageContent{}}, false},
		{ocpp2.IdTokenInfo{Status: ocpp2.AuthorizationStatusAccepted, CacheExpiryDateTime: ocpp2.NewDateTime(time.Now()), ChargingPriority: 1, Language1: "l1", Language2: "l2", GroupIdToken: &ocpp2.GroupIdToken{IdToken: "1234", Type: "invalidTokenType"}}, false},
		{ocpp2.IdTokenInfo{Status: ocpp2.AuthorizationStatusAccepted, CacheExpiryDateTime: ocpp2.NewDateTime(time.Now()), ChargingPriority: 1, Language1: "l1", Language2: "l2", GroupIdToken: &ocpp2.GroupIdToken{Type: ocpp2.IdTokenTypeCentral}}, false},
		{ocpp2.IdTokenInfo{Status: ocpp2.AuthorizationStatusAccepted, CacheExpiryDateTime: ocpp2.NewDateTime(time.Now()), ChargingPriority: 1, Language1: "l1", Language2: "l2", GroupIdToken: &ocpp2.GroupIdToken{IdToken: "1234"}}, false},
		{ocpp2.IdTokenInfo{Status: ocpp2.AuthorizationStatusAccepted, CacheExpiryDateTime: ocpp2.NewDateTime(time.Now()), ChargingPriority: 1, Language1: "l1", Language2: "l2", GroupIdToken: &ocpp2.GroupIdToken{}}, false},
		{ocpp2.IdTokenInfo{Status: ocpp2.AuthorizationStatusAccepted, CacheExpiryDateTime: ocpp2.NewDateTime(time.Now()), ChargingPriority: 1, Language1: "l1", Language2: ">8......."}, false},
		{ocpp2.IdTokenInfo{Status: ocpp2.AuthorizationStatusAccepted, CacheExpiryDateTime: ocpp2.NewDateTime(time.Now()), ChargingPriority: 1, Language1: ">8.......", Language2: "l2"}, false},
		{ocpp2.IdTokenInfo{Status: ocpp2.AuthorizationStatusAccepted, CacheExpiryDateTime: ocpp2.NewDateTime(time.Now()), ChargingPriority: -10}, false},
		{ocpp2.IdTokenInfo{Status: ocpp2.AuthorizationStatusAccepted, CacheExpiryDateTime: ocpp2.NewDateTime(time.Now()), ChargingPriority: 10}, false},
		{ocpp2.IdTokenInfo{Status: "invalidAuthStatus"}, false},
	}
	ExecuteGenericTestTable(suite.T(), testTable)
}

func (suite *OcppV2TestSuite) TestChargingSchedulePeriodValidation() {
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

func (suite *OcppV2TestSuite) TestChargingScheduleValidation() {
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

func (suite *OcppV2TestSuite) TestChargingProfileValidation() {
	t := suite.T()
	chargingSchedule := ocpp2.NewChargingSchedule(ocpp2.ChargingRateUnitWatts, ocpp2.NewChargingSchedulePeriod(0, 10.0), ocpp2.NewChargingSchedulePeriod(100, 8.0))
	var testTable = []GenericTestEntry{
		{ocpp2.ChargingProfile{ChargingProfileId: 1, TransactionId: 1, StackLevel: 1, ChargingProfilePurpose: ocpp2.ChargingProfilePurposeChargingStationMaxProfile, ChargingProfileKind: ocpp2.ChargingProfileKindAbsolute, RecurrencyKind: ocpp2.RecurrencyKindDaily, ValidFrom: ocpp2.NewDateTime(time.Now()), ValidTo: ocpp2.NewDateTime(time.Now().Add(8 * time.Hour)), ChargingSchedule: chargingSchedule}, true},
		{ocpp2.ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfilePurpose: ocpp2.ChargingProfilePurposeChargingStationMaxProfile, ChargingProfileKind: ocpp2.ChargingProfileKindAbsolute, ChargingSchedule: chargingSchedule}, true},
		{ocpp2.ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfilePurpose: ocpp2.ChargingProfilePurposeChargingStationMaxProfile, ChargingProfileKind: ocpp2.ChargingProfileKindAbsolute}, false},
		{ocpp2.ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfilePurpose: ocpp2.ChargingProfilePurposeChargingStationMaxProfile, ChargingSchedule: chargingSchedule}, false},
		{ocpp2.ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfileKind: ocpp2.ChargingProfileKindAbsolute, ChargingSchedule: chargingSchedule}, false},
		{ocpp2.ChargingProfile{ChargingProfileId: 1, ChargingProfilePurpose: ocpp2.ChargingProfilePurposeChargingStationMaxProfile, ChargingProfileKind: ocpp2.ChargingProfileKindAbsolute, ChargingSchedule: chargingSchedule}, false},
		{ocpp2.ChargingProfile{StackLevel: 1, ChargingProfilePurpose: ocpp2.ChargingProfilePurposeChargingStationMaxProfile, ChargingProfileKind: ocpp2.ChargingProfileKindAbsolute, ChargingSchedule: chargingSchedule}, true},
		{ocpp2.ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfilePurpose: ocpp2.ChargingProfilePurposeChargingStationMaxProfile, ChargingProfileKind: "invalidChargingProfileKind", ChargingSchedule: chargingSchedule}, false},
		{ocpp2.ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfilePurpose: "invalidChargingProfilePurpose", ChargingProfileKind: ocpp2.ChargingProfileKindAbsolute, ChargingSchedule: chargingSchedule}, false},
		{ocpp2.ChargingProfile{ChargingProfileId: 1, StackLevel: 0, ChargingProfilePurpose: ocpp2.ChargingProfilePurposeChargingStationMaxProfile, ChargingProfileKind: ocpp2.ChargingProfileKindAbsolute, ChargingSchedule: chargingSchedule}, false},
		{ocpp2.ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfilePurpose: ocpp2.ChargingProfilePurposeChargingStationMaxProfile, ChargingProfileKind: ocpp2.ChargingProfileKindAbsolute, RecurrencyKind: "invalidRecurrencyKind", ChargingSchedule: chargingSchedule}, false},
		{ocpp2.ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfilePurpose: ocpp2.ChargingProfilePurposeChargingStationMaxProfile, ChargingProfileKind: ocpp2.ChargingProfileKindAbsolute, ChargingSchedule: ocpp2.NewChargingSchedule(ocpp2.ChargingRateUnitWatts)}, false},
	}
	ExecuteGenericTestTable(t, testTable)
}

func (suite *OcppV2TestSuite) TestSampledValueValidation() {
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

func (suite *OcppV2TestSuite) TestMeterValueValidation() {
	var testTable = []GenericTestEntry{
		{ocpp2.MeterValue{Timestamp: ocpp2.NewDateTime(time.Now()), SampledValue: []ocpp2.SampledValue{{Value: "value"}, {Value: "value2", Unit: ocpp2.UnitOfMeasureKW}}}, true},
		{ocpp2.MeterValue{Timestamp: ocpp2.NewDateTime(time.Now()), SampledValue: []ocpp2.SampledValue{{Value: "value"}}}, true},
		{ocpp2.MeterValue{Timestamp: ocpp2.NewDateTime(time.Now()), SampledValue: []ocpp2.SampledValue{}}, false},
		{ocpp2.MeterValue{Timestamp: ocpp2.NewDateTime(time.Now())}, false},
		{ocpp2.MeterValue{SampledValue: []ocpp2.SampledValue{{Value: "value"}}}, false},
	}
	ExecuteGenericTestTable(suite.T(), testTable)
}
