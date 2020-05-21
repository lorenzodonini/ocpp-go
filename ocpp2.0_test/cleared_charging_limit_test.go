package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/smartcharging"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Tests
func (suite *OcppV2TestSuite) TestClearedChargingLimitRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{smartcharging.ClearedChargingLimitRequest{ChargingLimitSource: types.ChargingLimitSourceEMS, EvseID: newInt(0)}, true},
		{smartcharging.ClearedChargingLimitRequest{ChargingLimitSource: types.ChargingLimitSourceEMS}, true},
		{smartcharging.ClearedChargingLimitRequest{}, false},
		{smartcharging.ClearedChargingLimitRequest{ChargingLimitSource: types.ChargingLimitSourceEMS, EvseID: newInt(-1)}, false},
		{smartcharging.ClearedChargingLimitRequest{ChargingLimitSource: "invalidChargingLimitSource"}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestClearedChargingLimitConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{smartcharging.ClearedChargingLimitResponse{}, true},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestClearedChargingLimitE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := "1234"
	wsUrl := "someUrl"
	chargingLimitSource := types.ChargingLimitSourceEMS
	evseID := 42
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"chargingLimitSource":"%v","evseId":%v}]`, messageId, smartcharging.ClearedChargingLimitFeatureName, chargingLimitSource, evseID)
	responseJson := fmt.Sprintf(`[3,"%v",{}]`, messageId)
	clearedChargingLimitConfirmation := smartcharging.NewClearedChargingLimitResponse()
	channel := NewMockWebSocket(wsId)

	handler := MockCSMSSmartChargingHandler{}
	handler.On("OnClearedChargingLimit", mock.AnythingOfType("string"), mock.Anything).Return(clearedChargingLimitConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*smartcharging.ClearedChargingLimitRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, chargingLimitSource, request.ChargingLimitSource)
		assert.Equal(t, evseID, *request.EvseID)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	confirmation, err := suite.chargingStation.ClearedChargingLimit(chargingLimitSource, func(request *smartcharging.ClearedChargingLimitRequest) {
		request.EvseID = newInt(evseID)
	})
	require.Nil(t, err)
	require.NotNil(t, confirmation)
}

func (suite *OcppV2TestSuite) TestClearedChargingLimitInvalidEndpoint() {
	messageId := defaultMessageId
	chargingLimitSource := types.ChargingLimitSourceEMS
	evseID := 42
	clearedChargingLimitRequest := smartcharging.NewClearedChargingLimitRequest(chargingLimitSource)
	clearedChargingLimitRequest.EvseID = newInt(evseID)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"chargingLimitSource":"%v","evseId":%v}]`, messageId, smartcharging.ClearedChargingLimitFeatureName, chargingLimitSource, evseID)
	testUnsupportedRequestFromCentralSystem(suite, clearedChargingLimitRequest, requestJson, messageId)
}
