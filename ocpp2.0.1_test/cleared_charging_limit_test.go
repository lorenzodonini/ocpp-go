package ocpp2_test

import (
	"fmt"

	"github.com/stretchr/testify/mock"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/smartcharging"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

// Tests
func (suite *OcppV2TestSuite) TestClearedChargingLimitRequestValidation() {
	var requestTable = []GenericTestEntry{
		{smartcharging.ClearedChargingLimitRequest{ChargingLimitSource: types.ChargingLimitSourceEMS, EvseID: newInt(0)}, true},
		{smartcharging.ClearedChargingLimitRequest{ChargingLimitSource: types.ChargingLimitSourceEMS}, true},
		{smartcharging.ClearedChargingLimitRequest{}, false},
		{smartcharging.ClearedChargingLimitRequest{ChargingLimitSource: types.ChargingLimitSourceEMS, EvseID: newInt(-1)}, false},
		{smartcharging.ClearedChargingLimitRequest{ChargingLimitSource: "invalidChargingLimitSource"}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV2TestSuite) TestClearedChargingLimitConfirmationValidation() {
	var confirmationTable = []GenericTestEntry{
		{smartcharging.ClearedChargingLimitResponse{}, true},
	}
	ExecuteGenericTestTable(suite, confirmationTable)
}

func (suite *OcppV2TestSuite) TestClearedChargingLimitE2EMocked() {
	wsId := "test_id"
	messageId := "1234"
	wsUrl := "someUrl"
	chargingLimitSource := types.ChargingLimitSourceEMS
	evseID := 42
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"chargingLimitSource":"%v","evseId":%v}]`, messageId, smartcharging.ClearedChargingLimitFeatureName, chargingLimitSource, evseID)
	responseJson := fmt.Sprintf(`[3,"%v",{}]`, messageId)
	clearedChargingLimitConfirmation := smartcharging.NewClearedChargingLimitResponse()
	channel := NewMockWebSocket(wsId)

	handler := &MockCSMSSmartChargingHandler{}
	handler.On("OnClearedChargingLimit", mock.AnythingOfType("string"), mock.Anything).Return(clearedChargingLimitConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*smartcharging.ClearedChargingLimitRequest)
		suite.Require().True(ok)
		suite.Require().NotNil(request)
		suite.Equal(chargingLimitSource, request.ChargingLimitSource)
		suite.Equal(evseID, *request.EvseID)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	suite.Require().Nil(err)
	confirmation, err := suite.chargingStation.ClearedChargingLimit(chargingLimitSource, func(request *smartcharging.ClearedChargingLimitRequest) {
		request.EvseID = newInt(evseID)
	})
	suite.Require().Nil(err)
	suite.Require().NotNil(confirmation)
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
