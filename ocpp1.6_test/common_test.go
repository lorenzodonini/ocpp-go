package ocpp16_test

import (
	ocpp16 "github.com/lorenzodonini/go-ocpp/ocpp1.6"
	"time"
)

// Test
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
		{ocpp16.ChargingSchedule{Duration: 0, StartSchedule: ocpp16.DateTime{Time: time.Now()}, ChargingRateUnit: ocpp16.ChargingRateUnitWatts, ChargingSchedulePeriod: chargingSchedulePeriods, MinChargingRate: 1.0}, true},
		{ocpp16.ChargingSchedule{Duration: 0, ChargingRateUnit: ocpp16.ChargingRateUnitWatts, ChargingSchedulePeriod: chargingSchedulePeriods, MinChargingRate: 1.0}, true},
		{ocpp16.ChargingSchedule{Duration: 0, ChargingRateUnit: ocpp16.ChargingRateUnitWatts, ChargingSchedulePeriod: chargingSchedulePeriods}, true},
		{ocpp16.ChargingSchedule{Duration: 0, ChargingRateUnit: ocpp16.ChargingRateUnitWatts}, false},
		{ocpp16.ChargingSchedule{Duration: 0, ChargingSchedulePeriod: chargingSchedulePeriods}, false},
		{ocpp16.ChargingSchedule{Duration: -1, StartSchedule: ocpp16.DateTime{Time: time.Now()}, ChargingRateUnit: ocpp16.ChargingRateUnitWatts, ChargingSchedulePeriod: chargingSchedulePeriods, MinChargingRate: 1.0}, false},
		{ocpp16.ChargingSchedule{Duration: 0, StartSchedule: ocpp16.DateTime{Time: time.Now()}, ChargingRateUnit: ocpp16.ChargingRateUnitWatts, ChargingSchedulePeriod: chargingSchedulePeriods, MinChargingRate: -1.0}, false},
		{ocpp16.ChargingSchedule{Duration: 0, StartSchedule: ocpp16.DateTime{Time: time.Now()}, ChargingRateUnit: ocpp16.ChargingRateUnitWatts, ChargingSchedulePeriod: make([]ocpp16.ChargingSchedulePeriod, 0), MinChargingRate: 1.0}, false},
		{ocpp16.ChargingSchedule{Duration: -1, StartSchedule: ocpp16.DateTime{Time: time.Now()}, ChargingRateUnit: "invalidChargeRateUnit", ChargingSchedulePeriod: chargingSchedulePeriods, MinChargingRate: 1.0}, false},
	}
	ExecuteGenericTestTable(t, testTable)
}

func (suite *OcppV16TestSuite) TestChargingProfileValidation() {
	t := suite.T()
	chargingSchedule := ocpp16.NewChargingSchedule(ocpp16.ChargingRateUnitWatts, ocpp16.NewChargingSchedulePeriod(0, 10.0), ocpp16.NewChargingSchedulePeriod(100, 8.0))
	var testTable = []GenericTestEntry{
		{ocpp16.ChargingProfile{ChargingProfileId: 1, TransactionId: 1, StackLevel: 1, ChargingProfilePurpose: ocpp16.ChargingProfilePurposeChargePointMaxProfile, ChargingProfileKind: ocpp16.ChargingProfileKindAbsolute, RecurrencyKind: ocpp16.RecurrencyKindDaily, ValidFrom: ocpp16.DateTime{Time: time.Now()}, ValidTo: ocpp16.DateTime{Time: time.Now().Add(8 * time.Hour)}, ChargingSchedule: chargingSchedule}, true},
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
