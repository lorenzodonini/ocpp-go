package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test
func (suite *OcppV2TestSuite) TestCostUpdatedRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{ocpp2.CostUpdatedRequest{TotalCost: 24.6, TransactionID: "1234"}, true},
		{ocpp2.CostUpdatedRequest{TotalCost: 24.6}, false},
		{ocpp2.CostUpdatedRequest{TransactionID: "1234"}, false},
		{ocpp2.CostUpdatedRequest{}, false},
		{ocpp2.CostUpdatedRequest{TotalCost: 24.6, TransactionID: ">36.................................."}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestCostUpdatedConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{ocpp2.CostUpdatedConfirmation{}, true},
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
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"totalCost":%v,"transactionId":"%v"}]`, messageId, ocpp2.CostUpdatedFeatureName, totalCost, transactionId)
	responseJson := fmt.Sprintf(`[3,"%v",{}]`, messageId)
	costUpdatedConfirmation := ocpp2.NewCostUpdatedConfirmation()
	channel := NewMockWebSocket(wsId)

	coreListener := MockChargePointCoreListener{}
	coreListener.On("OnCostUpdated", mock.Anything).Return(costUpdatedConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*ocpp2.CostUpdatedRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, totalCost, request.TotalCost)
		assert.Equal(t, transactionId, request.TransactionID)
	})
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, coreListener, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.CostUpdated(wsId, func(confirmation *ocpp2.CostUpdatedConfirmation, err error) {
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
	costUpdatedRequest := ocpp2.NewCostUpdatedRequest(totalCost, transactionId)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"totalCost":%v,"transactionId":"%v"}]`, messageId, ocpp2.CostUpdatedFeatureName, totalCost, transactionId)
	testUnsupportedRequestFromChargePoint(suite, costUpdatedRequest, requestJson, messageId)
}
