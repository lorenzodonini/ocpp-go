package types

import (
	"testing"
	"time"

	"github.com/lorenzodonini/ocpp-go/tests"
	"github.com/stretchr/testify/suite"
)

type typesTestSuite struct {
	suite.Suite
}

func (suite *typesTestSuite) TestIdTokenInfoValidation() {
	var testTable = []tests.GenericTestEntry{
		{IdTokenInfo{Status: AuthorizationStatusAccepted, CacheExpiryDateTime: NewDateTime(time.Now()), ChargingPriority: 1, Language1: "l1", Language2: "l2", GroupIdToken: &GroupIdToken{IdToken: "1234", Type: IdTokenTypeCentral}, PersonalMessage: &MessageContent{Format: MessageFormatUTF8, Language: "en", Content: "random"}}, true},
		{IdTokenInfo{Status: AuthorizationStatusAccepted, CacheExpiryDateTime: NewDateTime(time.Now()), ChargingPriority: 1, Language1: "l1", Language2: "l2", GroupIdToken: &GroupIdToken{IdToken: "1234", Type: IdTokenTypeCentral}, PersonalMessage: &MessageContent{Format: MessageFormatUTF8, Content: "random"}}, true},
		{IdTokenInfo{Status: AuthorizationStatusAccepted, CacheExpiryDateTime: NewDateTime(time.Now()), ChargingPriority: 1, Language1: "l1", Language2: "l2", GroupIdToken: &GroupIdToken{IdToken: "1234", Type: IdTokenTypeCentral}}, true},
		{IdTokenInfo{Status: AuthorizationStatusAccepted, CacheExpiryDateTime: NewDateTime(time.Now()), ChargingPriority: 1, Language1: "l1", Language2: "l2"}, true},
		{IdTokenInfo{Status: AuthorizationStatusAccepted, CacheExpiryDateTime: NewDateTime(time.Now()), ChargingPriority: 1, Language1: "l1"}, true},
		{IdTokenInfo{Status: AuthorizationStatusAccepted, CacheExpiryDateTime: NewDateTime(time.Now()), ChargingPriority: 1}, true},
		{IdTokenInfo{Status: AuthorizationStatusAccepted, CacheExpiryDateTime: NewDateTime(time.Now())}, true},
		{IdTokenInfo{Status: AuthorizationStatusAccepted}, true},
		{IdTokenInfo{}, false},
		{IdTokenInfo{Status: AuthorizationStatusAccepted, CacheExpiryDateTime: NewDateTime(time.Now()), ChargingPriority: 1, Language1: "l1", Language2: "l2", GroupIdToken: &GroupIdToken{IdToken: "1234", Type: IdTokenTypeCentral}, PersonalMessage: &MessageContent{Format: "invalidFormat", Language: "en", Content: "random"}}, false},
		{IdTokenInfo{Status: AuthorizationStatusAccepted, CacheExpiryDateTime: NewDateTime(time.Now()), ChargingPriority: 1, Language1: "l1", Language2: "l2", GroupIdToken: &GroupIdToken{IdToken: "1234", Type: IdTokenTypeCentral}, PersonalMessage: &MessageContent{Format: MessageFormatUTF8, Language: "en", Content: ">512............................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................."}}, false},
		{IdTokenInfo{Status: AuthorizationStatusAccepted, CacheExpiryDateTime: NewDateTime(time.Now()), ChargingPriority: 1, Language1: "l1", Language2: "l2", GroupIdToken: &GroupIdToken{IdToken: "1234", Type: IdTokenTypeCentral}, PersonalMessage: &MessageContent{Format: MessageFormatUTF8, Language: "en"}}, false},
		{IdTokenInfo{Status: AuthorizationStatusAccepted, CacheExpiryDateTime: NewDateTime(time.Now()), ChargingPriority: 1, Language1: "l1", Language2: "l2", GroupIdToken: &GroupIdToken{IdToken: "1234", Type: IdTokenTypeCentral}, PersonalMessage: &MessageContent{}}, false},
		{IdTokenInfo{Status: AuthorizationStatusAccepted, CacheExpiryDateTime: NewDateTime(time.Now()), ChargingPriority: 1, Language1: "l1", Language2: "l2", GroupIdToken: &GroupIdToken{IdToken: "1234", Type: "invalidTokenType"}}, false},
		{IdTokenInfo{Status: AuthorizationStatusAccepted, CacheExpiryDateTime: NewDateTime(time.Now()), ChargingPriority: 1, Language1: "l1", Language2: "l2", GroupIdToken: &GroupIdToken{Type: IdTokenTypeCentral}}, false},
		{IdTokenInfo{Status: AuthorizationStatusAccepted, CacheExpiryDateTime: NewDateTime(time.Now()), ChargingPriority: 1, Language1: "l1", Language2: "l2", GroupIdToken: &GroupIdToken{IdToken: "1234"}}, false},
		{IdTokenInfo{Status: AuthorizationStatusAccepted, CacheExpiryDateTime: NewDateTime(time.Now()), ChargingPriority: 1, Language1: "l1", Language2: "l2", GroupIdToken: &GroupIdToken{}}, false},
		{IdTokenInfo{Status: AuthorizationStatusAccepted, CacheExpiryDateTime: NewDateTime(time.Now()), ChargingPriority: 1, Language1: "l1", Language2: ">8......."}, false},
		{IdTokenInfo{Status: AuthorizationStatusAccepted, CacheExpiryDateTime: NewDateTime(time.Now()), ChargingPriority: 1, Language1: ">8.......", Language2: "l2"}, false},
		{IdTokenInfo{Status: AuthorizationStatusAccepted, CacheExpiryDateTime: NewDateTime(time.Now()), ChargingPriority: -10}, false},
		{IdTokenInfo{Status: AuthorizationStatusAccepted, CacheExpiryDateTime: NewDateTime(time.Now()), ChargingPriority: 10}, false},
		{IdTokenInfo{Status: "invalidAuthStatus"}, false},
	}
	tests.ExecuteGenericTestTable(suite.T(), testTable)
}

func (suite *typesTestSuite) TestStatusInfo() {
	t := suite.T()
	var testTable = []tests.GenericTestEntry{
		{StatusInfo{ReasonCode: "okCode", AdditionalInfo: "someAdditionalInfo"}, true},
		{StatusInfo{ReasonCode: "okCode", AdditionalInfo: ""}, true},
		{StatusInfo{ReasonCode: "okCode"}, true},
		{StatusInfo{ReasonCode: ""}, false},
		{StatusInfo{}, false},
		{StatusInfo{ReasonCode: ">20.................."}, false},
		{StatusInfo{ReasonCode: "okCode", AdditionalInfo: ">512............................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................."}, false},
	}
	tests.ExecuteGenericTestTable(t, testTable)
}

func (suite *typesTestSuite) TestChargingSchedulePeriodValidation() {
	t := suite.T()
	var testTable = []tests.GenericTestEntry{
		{ChargingSchedulePeriod{StartPeriod: 0, Limit: 10.0, NumberPhases: tests.NewInt(3)}, true},
		{ChargingSchedulePeriod{StartPeriod: 0, Limit: 10.0}, true},
		{ChargingSchedulePeriod{StartPeriod: 0}, true},
		{ChargingSchedulePeriod{}, true},
		{ChargingSchedulePeriod{StartPeriod: 0, Limit: -1.0}, false},
		{ChargingSchedulePeriod{StartPeriod: -1, Limit: 10.0}, false},
		{ChargingSchedulePeriod{StartPeriod: 0, Limit: 10.0, NumberPhases: tests.NewInt(-1)}, false},
	}
	tests.ExecuteGenericTestTable(t, testTable)
}

func (suite *typesTestSuite) TestChargingScheduleValidation() {
	t := suite.T()
	chargingSchedulePeriods := make([]ChargingSchedulePeriod, 2)
	chargingSchedulePeriods[0] = NewChargingSchedulePeriod(0, 10.0)
	chargingSchedulePeriods[1] = NewChargingSchedulePeriod(100, 8.0)
	var testTable = []tests.GenericTestEntry{
		{ChargingSchedule{Duration: tests.NewInt(0), StartSchedule: NewDateTime(time.Now()), ChargingRateUnit: ChargingRateUnitWatts, ChargingSchedulePeriod: chargingSchedulePeriods, MinChargingRate: tests.NewFloat(1.0)}, true},
		{ChargingSchedule{Duration: tests.NewInt(0), ChargingRateUnit: ChargingRateUnitWatts, ChargingSchedulePeriod: chargingSchedulePeriods, MinChargingRate: tests.NewFloat(1.0)}, true},
		{ChargingSchedule{Duration: tests.NewInt(0), ChargingRateUnit: ChargingRateUnitWatts, ChargingSchedulePeriod: chargingSchedulePeriods}, true},
		{ChargingSchedule{Duration: tests.NewInt(0), ChargingRateUnit: ChargingRateUnitWatts}, false},
		{ChargingSchedule{Duration: tests.NewInt(0), ChargingSchedulePeriod: chargingSchedulePeriods}, false},
		{ChargingSchedule{Duration: tests.NewInt(-1), StartSchedule: NewDateTime(time.Now()), ChargingRateUnit: ChargingRateUnitWatts, ChargingSchedulePeriod: chargingSchedulePeriods, MinChargingRate: tests.NewFloat(1.0)}, false},
		{ChargingSchedule{Duration: tests.NewInt(0), StartSchedule: NewDateTime(time.Now()), ChargingRateUnit: ChargingRateUnitWatts, ChargingSchedulePeriod: chargingSchedulePeriods, MinChargingRate: tests.NewFloat(-1.0)}, false},
		{ChargingSchedule{Duration: tests.NewInt(0), StartSchedule: NewDateTime(time.Now()), ChargingRateUnit: ChargingRateUnitWatts, ChargingSchedulePeriod: make([]ChargingSchedulePeriod, 0), MinChargingRate: tests.NewFloat(1.0)}, false},
		{ChargingSchedule{Duration: tests.NewInt(-1), StartSchedule: NewDateTime(time.Now()), ChargingRateUnit: "invalidChargeRateUnit", ChargingSchedulePeriod: chargingSchedulePeriods, MinChargingRate: tests.NewFloat(1.0)}, false},
	}
	tests.ExecuteGenericTestTable(t, testTable)
}

func (suite *typesTestSuite) TestComponentVariableValidation() {
	t := suite.T()
	var testTable = []tests.GenericTestEntry{
		{ComponentVariable{Component: Component{Name: "component1", Instance: "instance1", EVSE: &EVSE{ID: 2, ConnectorID: tests.NewInt(2)}}, Variable: Variable{Name: "variable1", Instance: "instance1"}}, true},
		{ComponentVariable{Component: Component{Name: "component1", Instance: "instance1", EVSE: &EVSE{ID: 2}}, Variable: Variable{Name: "variable1", Instance: "instance1"}}, true},
		{ComponentVariable{Component: Component{Name: "component1", EVSE: &EVSE{ID: 2}}, Variable: Variable{Name: "variable1", Instance: "instance1"}}, true},
		{ComponentVariable{Component: Component{Name: "component1", EVSE: &EVSE{ID: 2}}, Variable: Variable{Name: "variable1"}}, true},
		{ComponentVariable{Component: Component{Name: "component1", EVSE: &EVSE{}}, Variable: Variable{Name: "variable1"}}, true},
		{ComponentVariable{Component: Component{Name: "component1"}, Variable: Variable{Name: "variable1"}}, true},
		{ComponentVariable{Component: Component{Name: "component1"}, Variable: Variable{}}, false},
		{ComponentVariable{Component: Component{}, Variable: Variable{Name: "variable1"}}, false},
		{ComponentVariable{Variable: Variable{Name: "variable1"}}, false},
		{ComponentVariable{Component: Component{Name: "component1"}}, false},
		{ComponentVariable{Component: Component{Name: ">50................................................", Instance: "instance1", EVSE: &EVSE{ID: 2, ConnectorID: tests.NewInt(2)}}, Variable: Variable{Name: "variable1", Instance: "instance1"}}, false},
		{ComponentVariable{Component: Component{Name: "component1", Instance: ">50................................................", EVSE: &EVSE{ID: 2, ConnectorID: tests.NewInt(2)}}, Variable: Variable{Name: "variable1", Instance: "instance1"}}, false},
		{ComponentVariable{Component: Component{Name: "component1", Instance: "instance1", EVSE: &EVSE{ID: 2, ConnectorID: tests.NewInt(2)}}, Variable: Variable{Name: ">50................................................", Instance: "instance1"}}, false},
		{ComponentVariable{Component: Component{Name: "component1", Instance: "instance1", EVSE: &EVSE{ID: 2, ConnectorID: tests.NewInt(2)}}, Variable: Variable{Name: "variable1", Instance: ">50................................................"}}, false},
		{ComponentVariable{Component: Component{Name: "component1", Instance: "instance1", EVSE: &EVSE{ID: 2, ConnectorID: tests.NewInt(-2)}}, Variable: Variable{Name: "variable1", Instance: "instance1"}}, false},
		{ComponentVariable{Component: Component{Name: "component1", Instance: "instance1", EVSE: &EVSE{ID: -2, ConnectorID: tests.NewInt(2)}}, Variable: Variable{Name: "variable1", Instance: "instance1"}}, false},
	}
	tests.ExecuteGenericTestTable(t, testTable)
}

func (suite *typesTestSuite) TestConsumptionCostValidation() {
	var testTable = []tests.GenericTestEntry{
		{NewConsumptionCost(1.0, []CostType{{CostKind: CostKindRelativePricePercentage, Amount: 7, AmountMultiplier: tests.NewInt(3)}}), true},
		{NewConsumptionCost(1.0, []CostType{{CostKind: CostKindRelativePricePercentage, Amount: 7, AmountMultiplier: tests.NewInt(-3)}}), true},
		{NewConsumptionCost(1.0, []CostType{{CostKind: CostKindRelativePricePercentage, Amount: 7}}), true},
		{NewConsumptionCost(1.0, []CostType{{CostKind: CostKindRelativePricePercentage}}), true},
		{ConsumptionCost{Cost: []CostType{{CostKind: CostKindRelativePricePercentage}}}, true},
		{NewConsumptionCost(1.0, []CostType{{}}), false},
		{NewConsumptionCost(1.0, []CostType{{CostKind: CostKindRelativePricePercentage, Amount: 7, AmountMultiplier: tests.NewInt(4)}}), false},
		{NewConsumptionCost(1.0, []CostType{{CostKind: CostKindRelativePricePercentage, Amount: 7, AmountMultiplier: tests.NewInt(-4)}}), false},
		{NewConsumptionCost(1.0, []CostType{{CostKind: CostKindRelativePricePercentage, Amount: -1, AmountMultiplier: tests.NewInt(3)}}), false},
		{NewConsumptionCost(1.0, []CostType{{CostKind: "invalidCostKind", Amount: 7, AmountMultiplier: tests.NewInt(3)}}), false},
		{NewConsumptionCost(1.0, []CostType{{CostKind: CostKindRelativePricePercentage, Amount: 7}, {CostKind: CostKindRelativePricePercentage, Amount: 7}, {CostKind: CostKindRelativePricePercentage, Amount: 7}, {CostKind: CostKindRelativePricePercentage, Amount: 7}}), false},
	}
	tests.ExecuteGenericTestTable(suite.T(), testTable)
}

func (suite *typesTestSuite) TestSalesTariffEntryValidation() {
	dummyCostType := NewConsumptionCost(1.0, []CostType{{CostKind: CostKindRelativePricePercentage, Amount: 7}})
	var testTable = []tests.GenericTestEntry{
		{SalesTariffEntry{EPriceLevel: tests.NewInt(8), RelativeTimeInterval: RelativeTimeInterval{Start: 500, Duration: tests.NewInt(1200)}, ConsumptionCost: []ConsumptionCost{dummyCostType}}, true},
		{SalesTariffEntry{EPriceLevel: tests.NewInt(8), RelativeTimeInterval: RelativeTimeInterval{Start: 500}}, true},
		{SalesTariffEntry{EPriceLevel: tests.NewInt(8), RelativeTimeInterval: RelativeTimeInterval{}}, true},
		{SalesTariffEntry{RelativeTimeInterval: RelativeTimeInterval{}}, true},
		{SalesTariffEntry{}, true},
		{SalesTariffEntry{EPriceLevel: tests.NewInt(-1), RelativeTimeInterval: RelativeTimeInterval{Start: 500, Duration: tests.NewInt(1200)}, ConsumptionCost: []ConsumptionCost{dummyCostType}}, false},
		{SalesTariffEntry{EPriceLevel: tests.NewInt(8), RelativeTimeInterval: RelativeTimeInterval{Start: 500, Duration: tests.NewInt(-1)}, ConsumptionCost: []ConsumptionCost{dummyCostType}}, false},
		{SalesTariffEntry{EPriceLevel: tests.NewInt(8), RelativeTimeInterval: RelativeTimeInterval{Start: 500, Duration: tests.NewInt(1200)}, ConsumptionCost: []ConsumptionCost{dummyCostType, dummyCostType, dummyCostType, dummyCostType}}, false},
		{SalesTariffEntry{EPriceLevel: tests.NewInt(8), RelativeTimeInterval: RelativeTimeInterval{Start: 500, Duration: tests.NewInt(1200)}, ConsumptionCost: []ConsumptionCost{NewConsumptionCost(1.0, []CostType{{}})}}, false},
	}
	tests.ExecuteGenericTestTable(suite.T(), testTable)
}

func (suite *typesTestSuite) TestSalesTariffValidation() {
	dummySalesTariffEntry := SalesTariffEntry{}
	var testTable = []tests.GenericTestEntry{
		{SalesTariff{ID: 1, SalesTariffDescription: "someDesc", NumEPriceLevels: tests.NewInt(1), SalesTariffEntry: []SalesTariffEntry{dummySalesTariffEntry}}, true},
		{SalesTariff{ID: 1, NumEPriceLevels: tests.NewInt(1), SalesTariffEntry: []SalesTariffEntry{dummySalesTariffEntry}}, true},
		{SalesTariff{ID: 1, SalesTariffEntry: []SalesTariffEntry{dummySalesTariffEntry}}, true},
		{SalesTariff{SalesTariffEntry: []SalesTariffEntry{dummySalesTariffEntry}}, true},
		{SalesTariff{SalesTariffEntry: []SalesTariffEntry{}}, false},
		{SalesTariff{}, false},
		{SalesTariff{ID: 1, SalesTariffDescription: ">32..............................", NumEPriceLevels: tests.NewInt(1), SalesTariffEntry: []SalesTariffEntry{dummySalesTariffEntry}}, false},
		{SalesTariff{ID: 1, SalesTariffDescription: "someDesc", NumEPriceLevels: tests.NewInt(1), SalesTariffEntry: []SalesTariffEntry{{EPriceLevel: tests.NewInt(-1)}}}, false},
	}
	tests.ExecuteGenericTestTable(suite.T(), testTable)
}

func (suite *typesTestSuite) TestChargingProfileValidation() {
	t := suite.T()
	chargingSchedule := NewChargingSchedule(1, ChargingRateUnitWatts, NewChargingSchedulePeriod(0, 10.0), NewChargingSchedulePeriod(100, 8.0))
	var testTable = []tests.GenericTestEntry{
		{ChargingProfile{ID: 1, StackLevel: 1, ChargingProfilePurpose: ChargingProfilePurposeChargingStationMaxProfile, ChargingProfileKind: ChargingProfileKindAbsolute, RecurrencyKind: RecurrencyKindDaily, ValidFrom: NewDateTime(time.Now()), ValidTo: NewDateTime(time.Now().Add(8 * time.Hour)), TransactionID: "d34d", ChargingSchedule: []ChargingSchedule{*chargingSchedule}}, true},
		{ChargingProfile{ID: 1, StackLevel: 1, ChargingProfilePurpose: ChargingProfilePurposeChargingStationMaxProfile, ChargingProfileKind: ChargingProfileKindAbsolute, ChargingSchedule: []ChargingSchedule{*chargingSchedule}}, true},
		{ChargingProfile{StackLevel: 1, ChargingProfilePurpose: ChargingProfilePurposeChargingStationMaxProfile, ChargingProfileKind: ChargingProfileKindAbsolute, ChargingSchedule: []ChargingSchedule{*chargingSchedule}}, true},
		{ChargingProfile{ID: 1, ChargingProfilePurpose: ChargingProfilePurposeChargingStationMaxProfile, ChargingProfileKind: ChargingProfileKindAbsolute, ChargingSchedule: []ChargingSchedule{*chargingSchedule}}, true},
		{ChargingProfile{ChargingProfilePurpose: ChargingProfilePurposeChargingStationMaxProfile, ChargingProfileKind: ChargingProfileKindAbsolute, ChargingSchedule: []ChargingSchedule{*chargingSchedule}}, true},
		{ChargingProfile{ID: 1, StackLevel: 1, ChargingProfilePurpose: ChargingProfilePurposeChargingStationMaxProfile, ChargingProfileKind: ChargingProfileKindAbsolute, ChargingSchedule: []ChargingSchedule{}}, false},
		{ChargingProfile{ID: 1, StackLevel: 1, ChargingProfilePurpose: ChargingProfilePurposeChargingStationMaxProfile, ChargingProfileKind: ChargingProfileKindAbsolute}, false},
		{ChargingProfile{ID: 1, StackLevel: 1, ChargingProfilePurpose: ChargingProfilePurposeChargingStationMaxProfile, ChargingSchedule: []ChargingSchedule{*chargingSchedule}}, false},
		{ChargingProfile{ID: 1, StackLevel: 1, ChargingProfileKind: ChargingProfileKindAbsolute, ChargingSchedule: []ChargingSchedule{*chargingSchedule}}, false},
		{ChargingProfile{ID: 1, StackLevel: 1, ChargingProfilePurpose: ChargingProfilePurposeChargingStationMaxProfile, ChargingProfileKind: "invalidChargingProfileKind", ChargingSchedule: []ChargingSchedule{*chargingSchedule}}, false},
		{ChargingProfile{ID: 1, StackLevel: 1, ChargingProfilePurpose: "invalidChargingProfilePurpose", ChargingProfileKind: ChargingProfileKindAbsolute, ChargingSchedule: []ChargingSchedule{*chargingSchedule}}, false},
		{ChargingProfile{ID: 1, StackLevel: -1, ChargingProfilePurpose: ChargingProfilePurposeChargingStationMaxProfile, ChargingProfileKind: ChargingProfileKindAbsolute, ChargingSchedule: []ChargingSchedule{*chargingSchedule}}, false},
		{ChargingProfile{ID: 1, StackLevel: 1, ChargingProfilePurpose: ChargingProfilePurposeChargingStationMaxProfile, ChargingProfileKind: ChargingProfileKindAbsolute, RecurrencyKind: "invalidRecurrencyKind", ChargingSchedule: []ChargingSchedule{*chargingSchedule}}, false},
		{ChargingProfile{ID: 1, StackLevel: 1, ChargingProfilePurpose: ChargingProfilePurposeChargingStationMaxProfile, ChargingProfileKind: ChargingProfileKindAbsolute, ChargingSchedule: []ChargingSchedule{*NewChargingSchedule(1, ChargingRateUnitWatts)}}, false},
		{ChargingProfile{ID: 1, StackLevel: 1, ChargingProfilePurpose: ChargingProfilePurposeChargingStationMaxProfile, ChargingProfileKind: ChargingProfileKindAbsolute, ChargingSchedule: []ChargingSchedule{*chargingSchedule, *chargingSchedule, *chargingSchedule, *chargingSchedule}}, false},
	}
	tests.ExecuteGenericTestTable(t, testTable)
}

func (suite *typesTestSuite) TestSignedMeterValue() {
	t := suite.T()
	var testTable = []tests.GenericTestEntry{
		{SignedMeterValue{SignedMeterData: "0xdeadbeef", SigningMethod: "ECDSAP256SHA256", EncodingMethod: "DLMS Message", PublicKey: "0xd34dc0de"}, true},
		{SignedMeterValue{SignedMeterData: "0xdeadbeef", SigningMethod: "ECDSAP256SHA256", EncodingMethod: "DLMS Message"}, false},
		{SignedMeterValue{SignedMeterData: "0xdeadbeef", SigningMethod: "ECDSAP256SHA256", PublicKey: "0xd34dc0de"}, false},
		{SignedMeterValue{SignedMeterData: "0xdeadbeef", EncodingMethod: "DLMS Message", PublicKey: "0xd34dc0de"}, false},
		{SignedMeterValue{SigningMethod: "ECDSAP256SHA256", EncodingMethod: "DLMS Message", PublicKey: "0xd34dc0de"}, false},
		{SignedMeterValue{SignedMeterData: ">2500................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................", SigningMethod: "ECDSAP256SHA256", EncodingMethod: "DLMS Message", PublicKey: "0xd34dc0de"}, false},
		{SignedMeterValue{SignedMeterData: "0xdeadbeef", SigningMethod: ">50................................................", EncodingMethod: "DLMS Message", PublicKey: "0xd34dc0de"}, false},
		{SignedMeterValue{SignedMeterData: "0xdeadbeef", SigningMethod: "ECDSAP256SHA256", EncodingMethod: ">50................................................", PublicKey: "0xd34dc0de"}, false},
		{SignedMeterValue{SignedMeterData: "0xdeadbeef", SigningMethod: "ECDSAP256SHA256", EncodingMethod: "DLMS Message", PublicKey: ">2500................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................"}, false},
	}
	tests.ExecuteGenericTestTable(t, testTable)
}

func (suite *typesTestSuite) TestSampledValueValidation() {
	t := suite.T()
	signedMeterValue := SignedMeterValue{
		SignedMeterData: "0xdeadbeef",
		SigningMethod:   "ECDSAP256SHA256",
		EncodingMethod:  "DLMS Message",
		PublicKey:       "0xd34dc0de",
	}
	var testTable = []tests.GenericTestEntry{
		{SampledValue{Value: 3.14, Context: ReadingContextTransactionEnd, Measurand: MeasurandPowerActiveExport, Phase: PhaseL2, Location: LocationBody, SignedMeterValue: &signedMeterValue, UnitOfMeasure: &UnitOfMeasure{Unit: "kW", Multiplier: tests.NewInt(0)}}, true},
		{SampledValue{Value: 3.14, Context: ReadingContextTransactionEnd, Measurand: MeasurandPowerActiveExport, Phase: PhaseL2, Location: LocationBody, SignedMeterValue: &signedMeterValue}, true},
		{SampledValue{Value: 3.14, Context: ReadingContextTransactionEnd, Measurand: MeasurandPowerActiveExport, Phase: PhaseL2, Location: LocationBody}, true},
		{SampledValue{Value: 3.14, Context: ReadingContextTransactionEnd, Measurand: MeasurandPowerActiveExport, Phase: PhaseL2}, true},
		{SampledValue{Value: 3.14, Context: ReadingContextTransactionEnd, Measurand: MeasurandPowerActiveExport}, true},
		{SampledValue{Value: 3.14, Context: ReadingContextTransactionEnd}, true},
		{SampledValue{Value: 3.14, Context: ReadingContextTransactionEnd}, true},
		{SampledValue{Value: 3.14}, true},
		{SampledValue{Value: -3.14}, true},
		{SampledValue{}, true},
		{SampledValue{Value: 3.14, Context: "invalidContext"}, false},
		{SampledValue{Value: 3.14, Measurand: "invalidMeasurand"}, false},
		{SampledValue{Value: 3.14, Phase: "invalidPhase"}, false},
		{SampledValue{Value: 3.14, Location: "invalidLocation"}, false},
		{SampledValue{Value: 3.14, SignedMeterValue: &SignedMeterValue{}}, false},
		{SampledValue{Value: 3.14, UnitOfMeasure: &UnitOfMeasure{Unit: "invalidUnit>20......."}}, false},
	}
	tests.ExecuteGenericTestTable(t, testTable)
}

func (suite *typesTestSuite) TestMeterValueValidation() {
	var testTable = []tests.GenericTestEntry{
		{MeterValue{Timestamp: DateTime{Time: time.Now()}, SampledValue: []SampledValue{{Value: 3.14, Context: ReadingContextTransactionEnd, Measurand: MeasurandPowerActiveExport, Phase: PhaseL2, Location: LocationBody}}}, true},
		{MeterValue{SampledValue: []SampledValue{{Value: 3.14, Context: ReadingContextTransactionEnd, Measurand: MeasurandPowerActiveExport, Phase: PhaseL2, Location: LocationBody}}}, true},
		{MeterValue{SampledValue: []SampledValue{}}, false},
		{MeterValue{}, false},
		{MeterValue{Timestamp: DateTime{Time: time.Now()}, SampledValue: []SampledValue{{Value: 3.14, Context: "invalidContext", Measurand: MeasurandPowerActiveExport, Phase: PhaseL2, Location: LocationBody}}}, false},
	}
	tests.ExecuteGenericTestTable(suite.T(), testTable)
}

func TestTypes(t *testing.T) {
	suite.Run(t, new(typesTestSuite))
}
