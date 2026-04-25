package smartcharging

import (
	"reflect"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *smartChargingTestSuite) TestNotifyEVChargingNeedsRequestValidation() {
	t := suite.T()
	chargingNeeds := ChargingNeeds{
		RequestedEnergyTransfer: EnergyTransferModeAC3Phase,
		DepartureTime:           types.NewDateTime(time.Now()),
		ACChargingParameters: &ACChargingParameters{
			EnergyAmount: 42,
			EVMinCurrent: 5,
			EVMaxCurrent: 10,
			EVMaxVoltage: 400,
		},
		DCChargingParameters: &DCChargingParameters{
			EVMaxCurrent:     0,
			EVMaxVoltage:     0,
			EnergyAmount:     tests.NewInt(42),
			EVMaxPower:       tests.NewInt(150),
			StateOfCharge:    tests.NewInt(50),
			EVEnergyCapacity: tests.NewInt(42),
			FullSoC:          tests.NewInt(100),
			BulkSoC:          tests.NewInt(80),
		},
	}
	var requestTable = []tests.GenericTestEntry{
		{NotifyEVChargingNeedsRequest{MaxScheduleTuples: tests.NewInt(5), EvseID: 1, ChargingNeeds: chargingNeeds}, true},
		{NotifyEVChargingNeedsRequest{EvseID: 1, ChargingNeeds: chargingNeeds}, true},
		{NotifyEVChargingNeedsRequest{EvseID: 1, ChargingNeeds: ChargingNeeds{RequestedEnergyTransfer: EnergyTransferModeAC3Phase}}, true},
		{NotifyEVChargingNeedsRequest{EvseID: 1, ChargingNeeds: ChargingNeeds{RequestedEnergyTransfer: EnergyTransferModeAC3Phase, ACChargingParameters: nil}}, true},
		{NotifyEVChargingNeedsRequest{EvseID: 1, ChargingNeeds: ChargingNeeds{RequestedEnergyTransfer: EnergyTransferModeDC, DCChargingParameters: nil}}, true},
		{NotifyEVChargingNeedsRequest{ChargingNeeds: chargingNeeds}, false},
		{NotifyEVChargingNeedsRequest{EvseID: 1}, false},
		{NotifyEVChargingNeedsRequest{}, false},
		{NotifyEVChargingNeedsRequest{MaxScheduleTuples: tests.NewInt(-1), EvseID: 1, ChargingNeeds: chargingNeeds}, false},
		{NotifyEVChargingNeedsRequest{EvseID: 1, ChargingNeeds: ChargingNeeds{RequestedEnergyTransfer: "invalidEnergyTransferMode"}}, false},
		{NotifyEVChargingNeedsRequest{EvseID: 1, ChargingNeeds: ChargingNeeds{RequestedEnergyTransfer: EnergyTransferModeAC3Phase, ACChargingParameters: &ACChargingParameters{EnergyAmount: -1}}}, false},
		{NotifyEVChargingNeedsRequest{EvseID: 1, ChargingNeeds: ChargingNeeds{RequestedEnergyTransfer: EnergyTransferModeDC, DCChargingParameters: &DCChargingParameters{EVMaxCurrent: -1}}}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *smartChargingTestSuite) TestDCChargingParametersValidation() {
	t := suite.T()
	var table = []tests.GenericTestEntry{
		{&DCChargingParameters{EVMaxCurrent: 0, EVMaxVoltage: 0, EnergyAmount: tests.NewInt(42), EVMaxPower: tests.NewInt(150), StateOfCharge: tests.NewInt(50), EVEnergyCapacity: tests.NewInt(42), FullSoC: tests.NewInt(100), BulkSoC: tests.NewInt(80)}, true},
		{&DCChargingParameters{EVMaxCurrent: 0, EVMaxVoltage: 0}, true},
		{&DCChargingParameters{}, true},
		{&DCChargingParameters{EVMaxCurrent: 0, EVMaxVoltage: -1, EnergyAmount: tests.NewInt(42), EVMaxPower: tests.NewInt(150), StateOfCharge: tests.NewInt(50), EVEnergyCapacity: tests.NewInt(42), FullSoC: tests.NewInt(100), BulkSoC: tests.NewInt(80)}, false},
		{&DCChargingParameters{EVMaxCurrent: -1, EVMaxVoltage: 0, EnergyAmount: tests.NewInt(42), EVMaxPower: tests.NewInt(150), StateOfCharge: tests.NewInt(50), EVEnergyCapacity: tests.NewInt(42), FullSoC: tests.NewInt(100), BulkSoC: tests.NewInt(80)}, false},
		{&DCChargingParameters{EVMaxCurrent: 0, EVMaxVoltage: 0, EnergyAmount: tests.NewInt(-1), EVMaxPower: tests.NewInt(150), StateOfCharge: tests.NewInt(50), EVEnergyCapacity: tests.NewInt(42), FullSoC: tests.NewInt(100), BulkSoC: tests.NewInt(80)}, false},
		{&DCChargingParameters{EVMaxCurrent: 0, EVMaxVoltage: 0, EnergyAmount: tests.NewInt(42), EVMaxPower: tests.NewInt(-1), StateOfCharge: tests.NewInt(50), EVEnergyCapacity: tests.NewInt(42), FullSoC: tests.NewInt(100), BulkSoC: tests.NewInt(80)}, false},
		{&DCChargingParameters{EVMaxCurrent: 0, EVMaxVoltage: 0, EnergyAmount: tests.NewInt(42), EVMaxPower: tests.NewInt(150), StateOfCharge: tests.NewInt(-1), EVEnergyCapacity: tests.NewInt(42), FullSoC: tests.NewInt(100), BulkSoC: tests.NewInt(80)}, false},
		{&DCChargingParameters{EVMaxCurrent: 0, EVMaxVoltage: 0, EnergyAmount: tests.NewInt(42), EVMaxPower: tests.NewInt(150), StateOfCharge: tests.NewInt(50), EVEnergyCapacity: tests.NewInt(-1), FullSoC: tests.NewInt(100), BulkSoC: tests.NewInt(80)}, false},
		{&DCChargingParameters{EVMaxCurrent: 0, EVMaxVoltage: 0, EnergyAmount: tests.NewInt(42), EVMaxPower: tests.NewInt(150), StateOfCharge: tests.NewInt(50), EVEnergyCapacity: tests.NewInt(42), FullSoC: tests.NewInt(-1), BulkSoC: tests.NewInt(80)}, false},
		{&DCChargingParameters{EVMaxCurrent: 0, EVMaxVoltage: 0, EnergyAmount: tests.NewInt(42), EVMaxPower: tests.NewInt(150), StateOfCharge: tests.NewInt(50), EVEnergyCapacity: tests.NewInt(42), FullSoC: tests.NewInt(100), BulkSoC: tests.NewInt(-1)}, false},
		{&DCChargingParameters{EVMaxCurrent: 0, EVMaxVoltage: 0, EnergyAmount: tests.NewInt(42), EVMaxPower: tests.NewInt(150), StateOfCharge: tests.NewInt(50), EVEnergyCapacity: tests.NewInt(42), FullSoC: tests.NewInt(101), BulkSoC: tests.NewInt(80)}, false},
		{&DCChargingParameters{EVMaxCurrent: 0, EVMaxVoltage: 0, EnergyAmount: tests.NewInt(42), EVMaxPower: tests.NewInt(150), StateOfCharge: tests.NewInt(101), EVEnergyCapacity: tests.NewInt(42), FullSoC: tests.NewInt(100), BulkSoC: tests.NewInt(80)}, false},
		{&DCChargingParameters{EVMaxCurrent: 0, EVMaxVoltage: 0, EnergyAmount: tests.NewInt(42), EVMaxPower: tests.NewInt(150), StateOfCharge: tests.NewInt(50), EVEnergyCapacity: tests.NewInt(42), FullSoC: tests.NewInt(100), BulkSoC: tests.NewInt(101)}, false},
	}
	tests.ExecuteGenericTestTable(t, table)
}

func (suite *smartChargingTestSuite) TestACChargingParametersValidation() {
	t := suite.T()
	var table = []tests.GenericTestEntry{
		{&ACChargingParameters{EnergyAmount: 42, EVMinCurrent: 6, EVMaxCurrent: 20, EVMaxVoltage: 400}, true},
		{&ACChargingParameters{}, true},
		{&ACChargingParameters{EnergyAmount: -1, EVMinCurrent: 0, EVMaxCurrent: 0, EVMaxVoltage: 0}, false},
		{&ACChargingParameters{EnergyAmount: 0, EVMinCurrent: -1, EVMaxCurrent: 0, EVMaxVoltage: 0}, false},
		{&ACChargingParameters{EnergyAmount: 0, EVMinCurrent: 0, EVMaxCurrent: -1, EVMaxVoltage: 0}, false},
		{&ACChargingParameters{EnergyAmount: 0, EVMinCurrent: 0, EVMaxCurrent: 0, EVMaxVoltage: -1}, false},
	}
	tests.ExecuteGenericTestTable(t, table)
}

func (suite *smartChargingTestSuite) TestNotifyEVChargingNeedsConfirmationValidation() {
	t := suite.T()
	var responseTable = []tests.GenericTestEntry{
		{NotifyEVChargingNeedsResponse{Status: EVChargingNeedsStatusAccepted, StatusInfo: types.NewStatusInfo("ok", "someInfo")}, true},
		{NotifyEVChargingNeedsResponse{Status: EVChargingNeedsStatusAccepted}, true},
		{NotifyEVChargingNeedsResponse{Status: EVChargingNeedsStatusRejected}, true},
		{NotifyEVChargingNeedsResponse{Status: EVChargingNeedsStatusProcessing}, true},
		{NotifyEVChargingNeedsResponse{}, false},
		{NotifyEVChargingNeedsResponse{Status: "invalidStatus"}, false},
		{NotifyEVChargingNeedsResponse{Status: EVChargingNeedsStatusAccepted, StatusInfo: types.NewStatusInfo("", "invalidStatusInfo")}, false},
	}
	tests.ExecuteGenericTestTable(t, responseTable)
}

func (suite *smartChargingTestSuite) TestNotifyEVChargingNeedsFeature() {
	feature := NotifyEVChargingNeedsFeature{}
	suite.Equal(NotifyEVChargingNeedsFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(NotifyEVChargingNeedsRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(NotifyEVChargingNeedsResponse{}), feature.GetResponseType())
}

func (suite *smartChargingTestSuite) TestNewNotifyEVChargingNeedsRequest() {
	needs := ChargingNeeds{RequestedEnergyTransfer: EnergyTransferModeAC3Phase}
	req := NewNotifyEVChargingNeedsRequest(1, needs)
	suite.NotNil(req)
	suite.Equal(NotifyEVChargingNeedsFeatureName, req.GetFeatureName())
	suite.Equal(1, req.EvseID)
	suite.Equal(needs, req.ChargingNeeds)
}

func (suite *smartChargingTestSuite) TestNewNotifyEVChargingNeedsResponse() {
	resp := NewNotifyEVChargingNeedsResponse(EVChargingNeedsStatusAccepted)
	suite.NotNil(resp)
	suite.Equal(NotifyEVChargingNeedsFeatureName, resp.GetFeatureName())
	suite.Equal(EVChargingNeedsStatusAccepted, resp.Status)
}
