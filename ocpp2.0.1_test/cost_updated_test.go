package ocpp2_test

import (
	"fmt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/tariffcost"
)

// Test
func (suite *OcppV2TestSuite) TestCostUpdatedRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{tariffcost.CostUpdatedRequest{TotalCost: 24.6, TransactionID: "1234"}, true},
		{tariffcost.CostUpdatedRequest{TotalCost: 24.6}, false},
		{tariffcost.CostUpdatedRequest{TransactionID: "1234"}, false},
		{tariffcost.CostUpdatedRequest{}, false},
		{tariffcost.CostUpdatedRequest{TotalCost: 24.6, TransactionID: ">36.................................."}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestCostUpdatedConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{tariffcost.CostUpdatedResponse{}, true},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestCostUpdatedE2EMocked() {
	t := suite.T()
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
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, totalCost, request.TotalCost)
		assert.Equal(t, transactionId, request.TransactionID)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.CostUpdated(wsId, func(confirmation *tariffcost.CostUpdatedResponse, err error) {
		require.Nil(t, err)
		require.NotNil(t, confirmation)
		resultChannel <- true
	}, totalCost, transactionId)
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV2TestSuite) TestCostUpdatedInvalidEndpoint() {
	messageId := defaultMessageId
	totalCost := 24.6
	transactionId := "1234"
	costUpdatedRequest := tariffcost.NewCostUpdatedRequest(totalCost, transactionId)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"totalCost":%v,"transactionId":"%v"}]`, messageId, tariffcost.CostUpdatedFeatureName, totalCost, transactionId)
	testUnsupportedRequestFromChargingStation(suite, costUpdatedRequest, requestJson, messageId)
}
