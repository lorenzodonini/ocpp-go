package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/smartcharging"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"time"
)

func (suite *OcppV2TestSuite) TestNotifyEVChargingNeedsRequestValidation() {
	t := suite.T()
	chargingNeeds := smartcharging.ChargingNeeds{
		RequestedEnergyTransfer: smartcharging.EnergyTransferModeAC3Phase,
		DepartureTime:           types.NewDateTime(time.Now()),
		ACChargingParameters: &smartcharging.ACChargingParameters{
			EnergyAmount: 42,
			EVMinCurrent: 5,
			EVMaxCurrent: 10,
			EVMaxVoltage: 400,
		},
		DCChargingParameters: &smartcharging.DCChargingParameters{
			EVMaxCurrent:     0,
			EVMaxVoltage:     0,
			EnergyAmount:     newInt(42),
			EVMaxPower:       newInt(150),
			StateOfCharge:    newInt(50),
			EVEnergyCapacity: newInt(42),
			FullSoC:          newInt(100),
			BulkSoC:          newInt(80),
		},
	}
	var requestTable = []GenericTestEntry{
		{smartcharging.NotifyEVChargingNeedsRequest{MaxScheduleTuples: newInt(5), EvseID: 1, ChargingNeeds: chargingNeeds}, true},
		{smartcharging.NotifyEVChargingNeedsRequest{EvseID: 1, ChargingNeeds: chargingNeeds}, true},
		{smartcharging.NotifyEVChargingNeedsRequest{EvseID: 1, ChargingNeeds: smartcharging.ChargingNeeds{RequestedEnergyTransfer: smartcharging.EnergyTransferModeAC3Phase}}, true},
		{smartcharging.NotifyEVChargingNeedsRequest{EvseID: 1, ChargingNeeds: smartcharging.ChargingNeeds{RequestedEnergyTransfer: smartcharging.EnergyTransferModeAC3Phase, ACChargingParameters: nil}}, true},
		{smartcharging.NotifyEVChargingNeedsRequest{EvseID: 1, ChargingNeeds: smartcharging.ChargingNeeds{RequestedEnergyTransfer: smartcharging.EnergyTransferModeDC, DCChargingParameters: nil}}, true},
		{smartcharging.NotifyEVChargingNeedsRequest{ChargingNeeds: chargingNeeds}, false},
		{smartcharging.NotifyEVChargingNeedsRequest{EvseID: 1}, false},
		{smartcharging.NotifyEVChargingNeedsRequest{}, false},
		{smartcharging.NotifyEVChargingNeedsRequest{MaxScheduleTuples: newInt(-1), EvseID: 1, ChargingNeeds: chargingNeeds}, false},
		{smartcharging.NotifyEVChargingNeedsRequest{EvseID: 1, ChargingNeeds: smartcharging.ChargingNeeds{RequestedEnergyTransfer: "invalidEnergyTransferMode"}}, false},
		{smartcharging.NotifyEVChargingNeedsRequest{EvseID: 1, ChargingNeeds: smartcharging.ChargingNeeds{RequestedEnergyTransfer: smartcharging.EnergyTransferModeAC3Phase, ACChargingParameters: &smartcharging.ACChargingParameters{EnergyAmount: -1}}}, false},
		{smartcharging.NotifyEVChargingNeedsRequest{EvseID: 1, ChargingNeeds: smartcharging.ChargingNeeds{RequestedEnergyTransfer: smartcharging.EnergyTransferModeDC, DCChargingParameters: &smartcharging.DCChargingParameters{EVMaxCurrent: -1}}}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestDCChargingParametersValidation() {
	t := suite.T()
	var table = []GenericTestEntry{
		{&smartcharging.DCChargingParameters{EVMaxCurrent: 0, EVMaxVoltage: 0, EnergyAmount: newInt(42), EVMaxPower: newInt(150), StateOfCharge: newInt(50), EVEnergyCapacity: newInt(42), FullSoC: newInt(100), BulkSoC: newInt(80)}, true},
		{&smartcharging.DCChargingParameters{EVMaxCurrent: 0, EVMaxVoltage: 0}, true},
		{&smartcharging.DCChargingParameters{}, true},
		{&smartcharging.DCChargingParameters{EVMaxCurrent: 0, EVMaxVoltage: -1, EnergyAmount: newInt(42), EVMaxPower: newInt(150), StateOfCharge: newInt(50), EVEnergyCapacity: newInt(42), FullSoC: newInt(100), BulkSoC: newInt(80)}, false},
		{&smartcharging.DCChargingParameters{EVMaxCurrent: -1, EVMaxVoltage: 0, EnergyAmount: newInt(42), EVMaxPower: newInt(150), StateOfCharge: newInt(50), EVEnergyCapacity: newInt(42), FullSoC: newInt(100), BulkSoC: newInt(80)}, false},
		{&smartcharging.DCChargingParameters{EVMaxCurrent: 0, EVMaxVoltage: 0, EnergyAmount: newInt(-1), EVMaxPower: newInt(150), StateOfCharge: newInt(50), EVEnergyCapacity: newInt(42), FullSoC: newInt(100), BulkSoC: newInt(80)}, false},
		{&smartcharging.DCChargingParameters{EVMaxCurrent: 0, EVMaxVoltage: 0, EnergyAmount: newInt(42), EVMaxPower: newInt(-1), StateOfCharge: newInt(50), EVEnergyCapacity: newInt(42), FullSoC: newInt(100), BulkSoC: newInt(80)}, false},
		{&smartcharging.DCChargingParameters{EVMaxCurrent: 0, EVMaxVoltage: 0, EnergyAmount: newInt(42), EVMaxPower: newInt(150), StateOfCharge: newInt(-1), EVEnergyCapacity: newInt(42), FullSoC: newInt(100), BulkSoC: newInt(80)}, false},
		{&smartcharging.DCChargingParameters{EVMaxCurrent: 0, EVMaxVoltage: 0, EnergyAmount: newInt(42), EVMaxPower: newInt(150), StateOfCharge: newInt(50), EVEnergyCapacity: newInt(-1), FullSoC: newInt(100), BulkSoC: newInt(80)}, false},
		{&smartcharging.DCChargingParameters{EVMaxCurrent: 0, EVMaxVoltage: 0, EnergyAmount: newInt(42), EVMaxPower: newInt(150), StateOfCharge: newInt(50), EVEnergyCapacity: newInt(42), FullSoC: newInt(-1), BulkSoC: newInt(80)}, false},
		{&smartcharging.DCChargingParameters{EVMaxCurrent: 0, EVMaxVoltage: 0, EnergyAmount: newInt(42), EVMaxPower: newInt(150), StateOfCharge: newInt(50), EVEnergyCapacity: newInt(42), FullSoC: newInt(100), BulkSoC: newInt(-1)}, false},
		{&smartcharging.DCChargingParameters{EVMaxCurrent: 0, EVMaxVoltage: 0, EnergyAmount: newInt(42), EVMaxPower: newInt(150), StateOfCharge: newInt(50), EVEnergyCapacity: newInt(42), FullSoC: newInt(101), BulkSoC: newInt(80)}, false},
		{&smartcharging.DCChargingParameters{EVMaxCurrent: 0, EVMaxVoltage: 0, EnergyAmount: newInt(42), EVMaxPower: newInt(150), StateOfCharge: newInt(101), EVEnergyCapacity: newInt(42), FullSoC: newInt(100), BulkSoC: newInt(80)}, false},
		{&smartcharging.DCChargingParameters{EVMaxCurrent: 0, EVMaxVoltage: 0, EnergyAmount: newInt(42), EVMaxPower: newInt(150), StateOfCharge: newInt(50), EVEnergyCapacity: newInt(42), FullSoC: newInt(100), BulkSoC: newInt(101)}, false},
	}
	ExecuteGenericTestTable(t, table)
}

func (suite *OcppV2TestSuite) TestACChargingParametersValidation() {
	t := suite.T()
	var table = []GenericTestEntry{
		{&smartcharging.ACChargingParameters{EnergyAmount: 42, EVMinCurrent: 6, EVMaxCurrent: 20, EVMaxVoltage: 400}, true},
		{&smartcharging.ACChargingParameters{}, true},
		{&smartcharging.ACChargingParameters{EnergyAmount: -1, EVMinCurrent: 0, EVMaxCurrent: 0, EVMaxVoltage: 0}, false},
		{&smartcharging.ACChargingParameters{EnergyAmount: 0, EVMinCurrent: -1, EVMaxCurrent: 0, EVMaxVoltage: 0}, false},
		{&smartcharging.ACChargingParameters{EnergyAmount: 0, EVMinCurrent: 0, EVMaxCurrent: -1, EVMaxVoltage: 0}, false},
		{&smartcharging.ACChargingParameters{EnergyAmount: 0, EVMinCurrent: 0, EVMaxCurrent: 0, EVMaxVoltage: -1}, false},
	}
	ExecuteGenericTestTable(t, table)
}

func (suite *OcppV2TestSuite) TestNotifyEVChargingNeedsConfirmationValidation() {
	t := suite.T()
	var responseTable = []GenericTestEntry{
		{smartcharging.NotifyEVChargingNeedsResponse{Status: smartcharging.EVChargingNeedsStatusAccepted, StatusInfo: types.NewStatusInfo("ok", "someInfo")}, true},
		{smartcharging.NotifyEVChargingNeedsResponse{Status: smartcharging.EVChargingNeedsStatusAccepted}, true},
		{smartcharging.NotifyEVChargingNeedsResponse{Status: smartcharging.EVChargingNeedsStatusRejected}, true},
		{smartcharging.NotifyEVChargingNeedsResponse{Status: smartcharging.EVChargingNeedsStatusProcessing}, true},
		{smartcharging.NotifyEVChargingNeedsResponse{}, false},
		{smartcharging.NotifyEVChargingNeedsResponse{Status: "invalidStatus"}, false},
		{smartcharging.NotifyEVChargingNeedsResponse{Status: smartcharging.EVChargingNeedsStatusAccepted, StatusInfo: types.NewStatusInfo("", "invalidStatusInfo")}, false},
	}
	ExecuteGenericTestTable(t, responseTable)
}

func (suite *OcppV2TestSuite) TestNotifyEVChargingNeedsE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := "1234"
	wsUrl := "someUrl"
	maxScheduleTuples := newInt(5)
	evseID := 42
	acChargingParams := &smartcharging.ACChargingParameters{
		EnergyAmount: 42,
		EVMinCurrent: 5,
		EVMaxCurrent: 10,
		EVMaxVoltage: 400,
	}
	chargingNeeds := smartcharging.ChargingNeeds{
		RequestedEnergyTransfer: smartcharging.EnergyTransferModeAC3Phase,
		DepartureTime:           types.NewDateTime(time.Now()),
		ACChargingParameters:    acChargingParams,
		DCChargingParameters:    nil,
	}
	status := smartcharging.EVChargingNeedsStatusAccepted
	statusInfo := types.NewStatusInfo("ok", "someInfo")
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"maxScheduleTuples":%v,"evseId":%v,"chargingNeeds":{"requestedEnergyTransfer":"%v","departureTime":"%v","acChargingParameters":{"energyAmount":%v,"evMinCurrent":%v,"evMaxCurrent":%v,"evMaxVoltage":%v}}}]`,
		messageId, smartcharging.NotifyEVChargingNeedsFeatureName, *maxScheduleTuples, evseID, chargingNeeds.RequestedEnergyTransfer, chargingNeeds.DepartureTime.FormatTimestamp(), acChargingParams.EnergyAmount, acChargingParams.EVMinCurrent, acChargingParams.EVMaxCurrent, acChargingParams.EVMaxVoltage)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v","statusInfo":{"reasonCode":"%v","additionalInfo":"%v"}}]`,
		messageId, status, statusInfo.ReasonCode, statusInfo.AdditionalInfo)
	notifyEVChargingNeedsResponse := smartcharging.NewNotifyEVChargingNeedsResponse(status)
	notifyEVChargingNeedsResponse.StatusInfo = statusInfo
	channel := NewMockWebSocket(wsId)

	handler := MockCSMSSmartChargingHandler{}
	handler.On("OnNotifyEVChargingNeeds", mock.AnythingOfType("string"), mock.Anything).Return(notifyEVChargingNeedsResponse, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*smartcharging.NotifyEVChargingNeedsRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, *maxScheduleTuples, *request.MaxScheduleTuples)
		assert.Equal(t, evseID, request.EvseID)
		assert.Equal(t, chargingNeeds.RequestedEnergyTransfer, request.ChargingNeeds.RequestedEnergyTransfer)
		assert.Equal(t, chargingNeeds.DepartureTime.FormatTimestamp(), request.ChargingNeeds.DepartureTime.FormatTimestamp())
		assert.Equal(t, acChargingParams.EnergyAmount, request.ChargingNeeds.ACChargingParameters.EnergyAmount)
		assert.Equal(t, acChargingParams.EVMinCurrent, request.ChargingNeeds.ACChargingParameters.EVMinCurrent)
		assert.Equal(t, acChargingParams.EVMaxCurrent, request.ChargingNeeds.ACChargingParameters.EVMaxCurrent)
		assert.Equal(t, acChargingParams.EVMaxVoltage, request.ChargingNeeds.ACChargingParameters.EVMaxVoltage)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	response, err := suite.chargingStation.NotifyEVChargingNeeds(evseID, chargingNeeds, func(request *smartcharging.NotifyEVChargingNeedsRequest) {
		request.MaxScheduleTuples = maxScheduleTuples
	})
	require.Nil(t, err)
	require.NotNil(t, response)
	assert.Equal(t, status, response.Status)
	assert.Equal(t, statusInfo.ReasonCode, response.StatusInfo.ReasonCode)
	assert.Equal(t, statusInfo.AdditionalInfo, response.StatusInfo.AdditionalInfo)
}

func (suite *OcppV2TestSuite) TestNotifyEVChargingNeedsInvalidEndpoint() {
	messageId := defaultMessageId
	maxScheduleTuples := newInt(5)
	evseID := 42
	acChargingParams := &smartcharging.ACChargingParameters{
		EnergyAmount: 42,
		EVMinCurrent: 5,
		EVMaxCurrent: 10,
		EVMaxVoltage: 400,
	}
	chargingNeeds := smartcharging.ChargingNeeds{
		RequestedEnergyTransfer: smartcharging.EnergyTransferModeAC3Phase,
		DepartureTime:           types.NewDateTime(time.Now()),
		ACChargingParameters:    acChargingParams,
		DCChargingParameters:    nil,
	}
	notifyEVChargingNeedsRequest := smartcharging.NewNotifyEVChargingNeedsRequest(evseID, chargingNeeds)
	notifyEVChargingNeedsRequest.MaxScheduleTuples = maxScheduleTuples
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"maxScheduleTuples":%v,"evseId":%v,"chargingNeeds":{"requestedEnergyTransfer":"%v","departureTime":"%v","acChargingParameters":{"energyAmount":%v,"evMinCurrent":%v,"evMaxCurrent":%v,"evMaxVoltage":%v}}}]`,
		messageId, smartcharging.NotifyEVChargingNeedsFeatureName, *maxScheduleTuples, evseID, chargingNeeds.RequestedEnergyTransfer, chargingNeeds.DepartureTime.FormatTimestamp(), acChargingParams.EnergyAmount, acChargingParams.EVMinCurrent, acChargingParams.EVMaxCurrent, acChargingParams.EVMaxVoltage)
	testUnsupportedRequestFromCentralSystem(suite, notifyEVChargingNeedsRequest, requestJson, messageId)
}
