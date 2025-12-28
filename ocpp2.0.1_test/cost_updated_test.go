package ocpp2_test

import (
	"fmt"

	"github.com/stretchr/testify/mock"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/tariffcost"
)

// Test
func (suite *OcppV2TestSuite) TestCostUpdatedRequestValidation() {
	var requestTable = []GenericTestEntry{
		{tariffcost.CostUpdatedRequest{TotalCost: 24.6, TransactionID: "1234"}, true},
		{tariffcost.CostUpdatedRequest{TotalCost: 24.6}, false},
		{tariffcost.CostUpdatedRequest{TransactionID: "1234"}, false},
		{tariffcost.CostUpdatedRequest{}, false},
		{tariffcost.CostUpdatedRequest{TotalCost: 24.6, TransactionID: ">36.................................."}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV2TestSuite) TestCostUpdatedConfirmationValidation() {
	var confirmationTable = []GenericTestEntry{
		{tariffcost.CostUpdatedResponse{}, true},
	}
	ExecuteGenericTestTable(suite, confirmationTable)
}

func (suite *OcppV2TestSuite) TestCostUpdatedE2EMocked() {
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	totalCost := 24.6
	transactionId := "1234"
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"totalCost":%v,"transactionId":"%v"}]`, messageId, tariffcost.CostUpdatedFeatureName, totalCost, transactionId)
	responseJson := fmt.Sprintf(`[3,"%v",{}]`, messageId)
	costUpdatedConfirmation := tariffcost.NewCostUpdatedResponse()
	channel := NewMockWebSocket(wsId)

	handler := &MockChargingStationTariffCostHandler{}
	handler.On("OnCostUpdated", mock.Anything).Return(costUpdatedConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*tariffcost.CostUpdatedRequest)
		suite.Require().True(ok)
		suite.Require().NotNil(request)
		suite.Equal(totalCost, request.TotalCost)
		suite.Equal(transactionId, request.TransactionID)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	suite.Require().Nil(err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.CostUpdated(wsId, func(confirmation *tariffcost.CostUpdatedResponse, err error) {
		suite.Require().Nil(err)
		suite.Require().NotNil(confirmation)
		resultChannel <- true
	}, totalCost, transactionId)
	suite.Require().Nil(err)
	result := <-resultChannel
	suite.True(result)
}

func (suite *OcppV2TestSuite) TestCostUpdatedInvalidEndpoint() {
	messageId := defaultMessageId
	totalCost := 24.6
	transactionId := "1234"
	costUpdatedRequest := tariffcost.NewCostUpdatedRequest(totalCost, transactionId)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"totalCost":%v,"transactionId":"%v"}]`, messageId, tariffcost.CostUpdatedFeatureName, totalCost, transactionId)
	testUnsupportedRequestFromChargingStation(suite, costUpdatedRequest, requestJson, messageId)
}
