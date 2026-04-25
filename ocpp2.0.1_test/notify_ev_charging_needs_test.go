package ocpp2_test

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/smartcharging"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *OcppV2TestSuite) TestNotifyEVChargingNeedsE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := "1234"
	wsUrl := "someUrl"
	maxScheduleTuples := tests.NewInt(5)
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

	handler := &MockCSMSSmartChargingHandler{}
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
	maxScheduleTuples := tests.NewInt(5)
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
