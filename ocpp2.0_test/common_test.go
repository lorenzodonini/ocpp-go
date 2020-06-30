package ocpp2_test

import (
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/display"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
	"time"
)

// Utility functions
func newInt(i int) *int {
	return &i
}

func newFloat(f float64) *float64 {
	return &f
}

func newBool(b bool) *bool {
	return &b
}

// Test
func (suite *OcppV2TestSuite) TestIdTokenInfoValidation() {
	var testTable = []GenericTestEntry{
		{types.IdTokenInfo{Status: types.AuthorizationStatusAccepted, CacheExpiryDateTime: types.NewDateTime(time.Now()), ChargingPriority: 1, Language1: "l1", Language2: "l2", GroupIdToken: &types.GroupIdToken{IdToken: "1234", Type: types.IdTokenTypeCentral}, PersonalMessage: &types.MessageContent{Format: types.MessageFormatUTF8, Language: "en", Content: "random"}}, true},
		{types.IdTokenInfo{Status: types.AuthorizationStatusAccepted, CacheExpiryDateTime: types.NewDateTime(time.Now()), ChargingPriority: 1, Language1: "l1", Language2: "l2", GroupIdToken: &types.GroupIdToken{IdToken: "1234", Type: types.IdTokenTypeCentral}, PersonalMessage: &types.MessageContent{Format: types.MessageFormatUTF8, Content: "random"}}, true},
		{types.IdTokenInfo{Status: types.AuthorizationStatusAccepted, CacheExpiryDateTime: types.NewDateTime(time.Now()), ChargingPriority: 1, Language1: "l1", Language2: "l2", GroupIdToken: &types.GroupIdToken{IdToken: "1234", Type: types.IdTokenTypeCentral}}, true},
		{types.IdTokenInfo{Status: types.AuthorizationStatusAccepted, CacheExpiryDateTime: types.NewDateTime(time.Now()), ChargingPriority: 1, Language1: "l1", Language2: "l2"}, true},
		{types.IdTokenInfo{Status: types.AuthorizationStatusAccepted, CacheExpiryDateTime: types.NewDateTime(time.Now()), ChargingPriority: 1, Language1: "l1"}, true},
		{types.IdTokenInfo{Status: types.AuthorizationStatusAccepted, CacheExpiryDateTime: types.NewDateTime(time.Now()), ChargingPriority: 1}, true},
		{types.IdTokenInfo{Status: types.AuthorizationStatusAccepted, CacheExpiryDateTime: types.NewDateTime(time.Now())}, true},
		{types.IdTokenInfo{Status: types.AuthorizationStatusAccepted}, true},
		{types.IdTokenInfo{}, false},
		{types.IdTokenInfo{Status: types.AuthorizationStatusAccepted, CacheExpiryDateTime: types.NewDateTime(time.Now()), ChargingPriority: 1, Language1: "l1", Language2: "l2", GroupIdToken: &types.GroupIdToken{IdToken: "1234", Type: types.IdTokenTypeCentral}, PersonalMessage: &types.MessageContent{Format: "invalidFormat", Language: "en", Content: "random"}}, false},
		{types.IdTokenInfo{Status: types.AuthorizationStatusAccepted, CacheExpiryDateTime: types.NewDateTime(time.Now()), ChargingPriority: 1, Language1: "l1", Language2: "l2", GroupIdToken: &types.GroupIdToken{IdToken: "1234", Type: types.IdTokenTypeCentral}, PersonalMessage: &types.MessageContent{Format: types.MessageFormatUTF8, Language: "en", Content: ">512............................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................."}}, false},
		{types.IdTokenInfo{Status: types.AuthorizationStatusAccepted, CacheExpiryDateTime: types.NewDateTime(time.Now()), ChargingPriority: 1, Language1: "l1", Language2: "l2", GroupIdToken: &types.GroupIdToken{IdToken: "1234", Type: types.IdTokenTypeCentral}, PersonalMessage: &types.MessageContent{Format: types.MessageFormatUTF8, Language: "en"}}, false},
		{types.IdTokenInfo{Status: types.AuthorizationStatusAccepted, CacheExpiryDateTime: types.NewDateTime(time.Now()), ChargingPriority: 1, Language1: "l1", Language2: "l2", GroupIdToken: &types.GroupIdToken{IdToken: "1234", Type: types.IdTokenTypeCentral}, PersonalMessage: &types.MessageContent{}}, false},
		{types.IdTokenInfo{Status: types.AuthorizationStatusAccepted, CacheExpiryDateTime: types.NewDateTime(time.Now()), ChargingPriority: 1, Language1: "l1", Language2: "l2", GroupIdToken: &types.GroupIdToken{IdToken: "1234", Type: "invalidTokenType"}}, false},
		{types.IdTokenInfo{Status: types.AuthorizationStatusAccepted, CacheExpiryDateTime: types.NewDateTime(time.Now()), ChargingPriority: 1, Language1: "l1", Language2: "l2", GroupIdToken: &types.GroupIdToken{Type: types.IdTokenTypeCentral}}, false},
		{types.IdTokenInfo{Status: types.AuthorizationStatusAccepted, CacheExpiryDateTime: types.NewDateTime(time.Now()), ChargingPriority: 1, Language1: "l1", Language2: "l2", GroupIdToken: &types.GroupIdToken{IdToken: "1234"}}, false},
		{types.IdTokenInfo{Status: types.AuthorizationStatusAccepted, CacheExpiryDateTime: types.NewDateTime(time.Now()), ChargingPriority: 1, Language1: "l1", Language2: "l2", GroupIdToken: &types.GroupIdToken{}}, false},
		{types.IdTokenInfo{Status: types.AuthorizationStatusAccepted, CacheExpiryDateTime: types.NewDateTime(time.Now()), ChargingPriority: 1, Language1: "l1", Language2: ">8......."}, false},
		{types.IdTokenInfo{Status: types.AuthorizationStatusAccepted, CacheExpiryDateTime: types.NewDateTime(time.Now()), ChargingPriority: 1, Language1: ">8.......", Language2: "l2"}, false},
		{types.IdTokenInfo{Status: types.AuthorizationStatusAccepted, CacheExpiryDateTime: types.NewDateTime(time.Now()), ChargingPriority: -10}, false},
		{types.IdTokenInfo{Status: types.AuthorizationStatusAccepted, CacheExpiryDateTime: types.NewDateTime(time.Now()), ChargingPriority: 10}, false},
		{types.IdTokenInfo{Status: "invalidAuthStatus"}, false},
	}
	ExecuteGenericTestTable(suite.T(), testTable)
}

func (suite *OcppV2TestSuite) TestChargingSchedulePeriodValidation() {
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

func (suite *OcppV2TestSuite) TestChargingScheduleValidation() {
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
		{types.ChargingSchedule{Duration: newInt(-1), StartSchedule: types.NewDateTime(time.Now()), ChargingRateUnit: "invalidChargeRateUnit", ChargingSchedulePeriod: chargingSchedulePeriods, MinChargingRate: newFloat(1.0)}, false},
	}
	ExecuteGenericTestTable(t, testTable)
}

func (suite *OcppV2TestSuite) TestComponentVariableValidation() {
	t := suite.T()
	var testTable = []GenericTestEntry{
		{types.ComponentVariable{Component: types.Component{Name: "component1", Instance: "instance1", EVSE: &types.EVSE{ID: 2, ConnectorID: newInt(2)}}, Variable: types.Variable{Name: "variable1", Instance: "instance1"}}, true},
		{types.ComponentVariable{Component: types.Component{Name: "component1", Instance: "instance1", EVSE: &types.EVSE{ID: 2}}, Variable: types.Variable{Name: "variable1", Instance: "instance1"}}, true},
		{types.ComponentVariable{Component: types.Component{Name: "component1", EVSE: &types.EVSE{ID: 2}}, Variable: types.Variable{Name: "variable1", Instance: "instance1"}}, true},
		{types.ComponentVariable{Component: types.Component{Name: "component1", EVSE: &types.EVSE{ID: 2}}, Variable: types.Variable{Name: "variable1"}}, true},
		{types.ComponentVariable{Component: types.Component{Name: "component1", EVSE: &types.EVSE{}}, Variable: types.Variable{Name: "variable1"}}, true},
		{types.ComponentVariable{Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}, true},
		{types.ComponentVariable{Component: types.Component{Name: "component1"}, Variable: types.Variable{}}, false},
		{types.ComponentVariable{Component: types.Component{}, Variable: types.Variable{Name: "variable1"}}, false},
		{types.ComponentVariable{Variable: types.Variable{Name: "variable1"}}, false},
		{types.ComponentVariable{Component: types.Component{Name: "component1"}}, false},
		{types.ComponentVariable{Component: types.Component{Name: ">50................................................", Instance: "instance1", EVSE: &types.EVSE{ID: 2, ConnectorID: newInt(2)}}, Variable: types.Variable{Name: "variable1", Instance: "instance1"}}, false},
		{types.ComponentVariable{Component: types.Component{Name: "component1", Instance: ">50................................................", EVSE: &types.EVSE{ID: 2, ConnectorID: newInt(2)}}, Variable: types.Variable{Name: "variable1", Instance: "instance1"}}, false},
		{types.ComponentVariable{Component: types.Component{Name: "component1", Instance: "instance1", EVSE: &types.EVSE{ID: 2, ConnectorID: newInt(2)}}, Variable: types.Variable{Name: ">50................................................", Instance: "instance1"}}, false},
		{types.ComponentVariable{Component: types.Component{Name: "component1", Instance: "instance1", EVSE: &types.EVSE{ID: 2, ConnectorID: newInt(2)}}, Variable: types.Variable{Name: "variable1", Instance: ">50................................................"}}, false},
		{types.ComponentVariable{Component: types.Component{Name: "component1", Instance: "instance1", EVSE: &types.EVSE{ID: 2, ConnectorID: newInt(-2)}}, Variable: types.Variable{Name: "variable1", Instance: "instance1"}}, false},
		{types.ComponentVariable{Component: types.Component{Name: "component1", Instance: "instance1", EVSE: &types.EVSE{ID: -2, ConnectorID: newInt(2)}}, Variable: types.Variable{Name: "variable1", Instance: "instance1"}}, false},
	}
	ExecuteGenericTestTable(t, testTable)
}

func (suite *OcppV2TestSuite) TestChargingProfileValidation() {
	t := suite.T()
	chargingSchedule := types.NewChargingSchedule(types.ChargingRateUnitWatts, types.NewChargingSchedulePeriod(0, 10.0), types.NewChargingSchedulePeriod(100, 8.0))
	var testTable = []GenericTestEntry{
		{types.ChargingProfile{ChargingProfileId: 1, TransactionId: 1, StackLevel: 1, ChargingProfilePurpose: types.ChargingProfilePurposeChargingStationMaxProfile, ChargingProfileKind: types.ChargingProfileKindAbsolute, RecurrencyKind: types.RecurrencyKindDaily, ValidFrom: types.NewDateTime(time.Now()), ValidTo: types.NewDateTime(time.Now().Add(8 * time.Hour)), ChargingSchedule: chargingSchedule}, true},
		{types.ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfilePurpose: types.ChargingProfilePurposeChargingStationMaxProfile, ChargingProfileKind: types.ChargingProfileKindAbsolute, ChargingSchedule: chargingSchedule}, true},
		{types.ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfilePurpose: types.ChargingProfilePurposeChargingStationMaxProfile, ChargingProfileKind: types.ChargingProfileKindAbsolute}, false},
		{types.ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfilePurpose: types.ChargingProfilePurposeChargingStationMaxProfile, ChargingSchedule: chargingSchedule}, false},
		{types.ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfileKind: types.ChargingProfileKindAbsolute, ChargingSchedule: chargingSchedule}, false},
		{types.ChargingProfile{ChargingProfileId: 1, ChargingProfilePurpose: types.ChargingProfilePurposeChargingStationMaxProfile, ChargingProfileKind: types.ChargingProfileKindAbsolute, ChargingSchedule: chargingSchedule}, false},
		{types.ChargingProfile{StackLevel: 1, ChargingProfilePurpose: types.ChargingProfilePurposeChargingStationMaxProfile, ChargingProfileKind: types.ChargingProfileKindAbsolute, ChargingSchedule: chargingSchedule}, true},
		{types.ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfilePurpose: types.ChargingProfilePurposeChargingStationMaxProfile, ChargingProfileKind: "invalidChargingProfileKind", ChargingSchedule: chargingSchedule}, false},
		{types.ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfilePurpose: "invalidChargingProfilePurpose", ChargingProfileKind: types.ChargingProfileKindAbsolute, ChargingSchedule: chargingSchedule}, false},
		{types.ChargingProfile{ChargingProfileId: 1, StackLevel: 0, ChargingProfilePurpose: types.ChargingProfilePurposeChargingStationMaxProfile, ChargingProfileKind: types.ChargingProfileKindAbsolute, ChargingSchedule: chargingSchedule}, false},
		{types.ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfilePurpose: types.ChargingProfilePurposeChargingStationMaxProfile, ChargingProfileKind: types.ChargingProfileKindAbsolute, RecurrencyKind: "invalidRecurrencyKind", ChargingSchedule: chargingSchedule}, false},
		{types.ChargingProfile{ChargingProfileId: 1, StackLevel: 1, ChargingProfilePurpose: types.ChargingProfilePurposeChargingStationMaxProfile, ChargingProfileKind: types.ChargingProfileKindAbsolute, ChargingSchedule: types.NewChargingSchedule(types.ChargingRateUnitWatts)}, false},
	}
	ExecuteGenericTestTable(t, testTable)
}

func (suite *OcppV2TestSuite) TestSignedMeterValue() {
	t := suite.T()
	var testTable = []GenericTestEntry{
		{types.SignedMeterValue{SignedMeterData: "0xdeadbeef", SigningMethod: "ECDSAP256SHA256", EncodingMethod: "DLMS Message", PublicKey: "0xd34dc0de"}, true},
		{types.SignedMeterValue{SignedMeterData: "0xdeadbeef", SigningMethod: "ECDSAP256SHA256", EncodingMethod: "DLMS Message"}, false},
		{types.SignedMeterValue{SignedMeterData: "0xdeadbeef", SigningMethod: "ECDSAP256SHA256", PublicKey: "0xd34dc0de"}, false},
		{types.SignedMeterValue{SignedMeterData: "0xdeadbeef", EncodingMethod: "DLMS Message", PublicKey: "0xd34dc0de"}, false},
		{types.SignedMeterValue{SigningMethod: "ECDSAP256SHA256", EncodingMethod: "DLMS Message", PublicKey: "0xd34dc0de"}, false},
		{types.SignedMeterValue{SignedMeterData: ">2500................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................", SigningMethod: "ECDSAP256SHA256", EncodingMethod: "DLMS Message", PublicKey: "0xd34dc0de"}, false},
		{types.SignedMeterValue{SignedMeterData: "0xdeadbeef", SigningMethod: ">50................................................", EncodingMethod: "DLMS Message", PublicKey: "0xd34dc0de"}, false},
		{types.SignedMeterValue{SignedMeterData: "0xdeadbeef", SigningMethod: "ECDSAP256SHA256", EncodingMethod: ">50................................................", PublicKey: "0xd34dc0de"}, false},
		{types.SignedMeterValue{SignedMeterData: "0xdeadbeef", SigningMethod: "ECDSAP256SHA256", EncodingMethod: "DLMS Message", PublicKey: ">2500................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................"}, false},
	}
	ExecuteGenericTestTable(t, testTable)
}

func (suite *OcppV2TestSuite) TestSampledValueValidation() {
	t := suite.T()
	signedMeterValue := types.SignedMeterValue{
		SignedMeterData: "0xdeadbeef",
		SigningMethod:   "ECDSAP256SHA256",
		EncodingMethod:  "DLMS Message",
		PublicKey:       "0xd34dc0de",
	}
	var testTable = []GenericTestEntry{
		{types.SampledValue{Value: 3.14, Context: types.ReadingContextTransactionEnd, Measurand: types.MeasurandPowerActiveExport, Phase: types.PhaseL2, Location: types.LocationBody, SignedMeterValue: &signedMeterValue, UnitOfMeasure: &types.UnitOfMeasure{Unit: "kW", Multiplier: newInt(0)}}, true},
		{types.SampledValue{Value: 3.14, Context: types.ReadingContextTransactionEnd, Measurand: types.MeasurandPowerActiveExport, Phase: types.PhaseL2, Location: types.LocationBody, SignedMeterValue: &signedMeterValue}, true},
		{types.SampledValue{Value: 3.14, Context: types.ReadingContextTransactionEnd, Measurand: types.MeasurandPowerActiveExport, Phase: types.PhaseL2, Location: types.LocationBody}, true},
		{types.SampledValue{Value: 3.14, Context: types.ReadingContextTransactionEnd, Measurand: types.MeasurandPowerActiveExport, Phase: types.PhaseL2}, true},
		{types.SampledValue{Value: 3.14, Context: types.ReadingContextTransactionEnd, Measurand: types.MeasurandPowerActiveExport}, true},
		{types.SampledValue{Value: 3.14, Context: types.ReadingContextTransactionEnd}, true},
		{types.SampledValue{Value: 3.14, Context: types.ReadingContextTransactionEnd}, true},
		{types.SampledValue{Value: 3.14}, true},
		{types.SampledValue{Value: -3.14}, true},
		{types.SampledValue{Value: 3.14, Context: "invalidContext"}, false},
		{types.SampledValue{Value: 3.14, Measurand: "invalidMeasurand"}, false},
		{types.SampledValue{Value: 3.14, Phase: "invalidPhase"}, false},
		{types.SampledValue{Value: 3.14, Location: "invalidLocation"}, false},
		{types.SampledValue{Value: 3.14, SignedMeterValue: &types.SignedMeterValue{}}, false},
		{types.SampledValue{Value: 3.14, UnitOfMeasure: &types.UnitOfMeasure{Unit: "invalidUnit>20......."}}, false},
	}
	ExecuteGenericTestTable(t, testTable)
}

func (suite *OcppV2TestSuite) TestMeterValueValidation() {
	var testTable = []GenericTestEntry{
		{types.MeterValue{Timestamp: types.DateTime{Time: time.Now()}, SampledValue: []types.SampledValue{ {Value: 3.14, Context: types.ReadingContextTransactionEnd, Measurand: types.MeasurandPowerActiveExport, Phase: types.PhaseL2, Location: types.LocationBody}}}, true},
		{types.MeterValue{SampledValue: []types.SampledValue{ {Value: 3.14, Context: types.ReadingContextTransactionEnd, Measurand: types.MeasurandPowerActiveExport, Phase: types.PhaseL2, Location: types.LocationBody}}}, true},
		{types.MeterValue{SampledValue: []types.SampledValue{}}, false},
		{types.MeterValue{}, false},
		{types.MeterValue{Timestamp: types.DateTime{Time: time.Now()}, SampledValue: []types.SampledValue{ {Value: 3.14, Context: "invalidContext", Measurand: types.MeasurandPowerActiveExport, Phase: types.PhaseL2, Location: types.LocationBody}}}, false},
	}
	ExecuteGenericTestTable(suite.T(), testTable)
}

func (suite *OcppV2TestSuite) TestMessageInfoValidation() {
	var testTable = []GenericTestEntry{
		{display.MessageInfo{ID: 42, Priority: display.MessagePriorityAlwaysFront, State: display.MessageStateIdle, StartDateTime: types.NewDateTime(time.Now()), EndDateTime: types.NewDateTime(time.Now().Add(1*time.Hour)), TransactionID: "123456", Message: types.MessageContent{Format: types.MessageFormatUTF8, Content: "hello world"}, Display: &types.Component{Name: "name1"}}, true},
		{display.MessageInfo{ID: 42, Priority: display.MessagePriorityAlwaysFront, State: display.MessageStateIdle, StartDateTime: types.NewDateTime(time.Now()), EndDateTime: types.NewDateTime(time.Now().Add(1*time.Hour)), TransactionID: "123456", Message: types.MessageContent{Format: types.MessageFormatUTF8, Content: "hello world"}}, true},
		{display.MessageInfo{ID: 42, Priority: display.MessagePriorityAlwaysFront, State: display.MessageStateIdle, StartDateTime: types.NewDateTime(time.Now()), TransactionID: "123456", Message: types.MessageContent{Format: types.MessageFormatUTF8, Content: "hello world"}}, true},
		{display.MessageInfo{ID: 42, Priority: display.MessagePriorityAlwaysFront, State: display.MessageStateIdle, TransactionID: "123456", Message: types.MessageContent{Format: types.MessageFormatUTF8, Content: "hello world"}}, true},
		{display.MessageInfo{ID: 42, Priority: display.MessagePriorityAlwaysFront, State: display.MessageStateIdle, Message: types.MessageContent{Format: types.MessageFormatUTF8, Content: "hello world"}}, true},
		{display.MessageInfo{ID: 42, Priority: display.MessagePriorityAlwaysFront, Message: types.MessageContent{Format: types.MessageFormatUTF8, Content: "hello world"}}, true},
		{display.MessageInfo{ID: 42, Priority: display.MessagePriorityAlwaysFront, State: display.MessageStateIdle}, false},
		{display.MessageInfo{ID: 42, Priority: display.MessagePriorityAlwaysFront, State: display.MessageStateIdle, Message: types.MessageContent{Format: types.MessageFormatUTF8}}, false},
		{display.MessageInfo{ID: 42, Priority: display.MessagePriorityAlwaysFront, State: "invalidState", Message: types.MessageContent{Format: types.MessageFormatUTF8, Content: "hello world"}}, false},
		{display.MessageInfo{ID: 42, State: display.MessageStateIdle, Message: types.MessageContent{Format: types.MessageFormatUTF8, Content: "hello world"}}, false},
		{display.MessageInfo{ID: 42, Priority: "invalidPriority", State: display.MessageStateIdle, Message: types.MessageContent{Format: types.MessageFormatUTF8, Content: "hello world"}}, false},
		{display.MessageInfo{ID: -1, Priority: display.MessagePriorityAlwaysFront, State: display.MessageStateIdle, Message: types.MessageContent{Format: types.MessageFormatUTF8, Content: "hello world"}}, false},
		{display.MessageInfo{ID: 42, Priority: display.MessagePriorityAlwaysFront, State: display.MessageStateIdle, TransactionID: ">36..................................", Message: types.MessageContent{Format: types.MessageFormatUTF8, Content: "hello world"}}, false},
		{display.MessageInfo{ID: 42, Priority: display.MessagePriorityAlwaysFront, State: display.MessageStateIdle, StartDateTime: types.NewDateTime(time.Now()), EndDateTime: types.NewDateTime(time.Now().Add(1*time.Hour)), TransactionID: "123456", Message: types.MessageContent{Format: types.MessageFormatUTF8, Content: "hello world"}, Display: &types.Component{}}, false},
	}
	ExecuteGenericTestTable(suite.T(), testTable)
}
