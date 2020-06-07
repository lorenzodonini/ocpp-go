package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/transactions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test
func (suite *OcppV2TestSuite) TestGetTransactionStatusRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{transactions.GetTransactionStatusRequest{}, true},
		{transactions.GetTransactionStatusRequest{TransactionID: "12345"}, true},
		{transactions.GetTransactionStatusRequest{TransactionID: ">36.................................."}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestGetTransactionStatusResponseValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{transactions.GetTransactionStatusResponse{OngoingIndicator: newBool(true), MessageInQueue: true}, true},
		{transactions.GetTransactionStatusResponse{MessageInQueue: true}, true},
		{transactions.GetTransactionStatusResponse{}, true},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestGetTransactionStatusE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	transactionID := "12345"
	messageInQueue := false
	ongoingIndicator := newBool(true)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"transactionId":"%v"}]`, messageId, transactions.GetTransactionStatusFeatureName, transactionID)
	responseJson := fmt.Sprintf(`[3,"%v",{"ongoingIndicator":%v,"messageInQueue":%v}]`, messageId, *ongoingIndicator, messageInQueue)
	getTransactionStatusResponse := transactions.NewGetTransactionStatusResponse(messageInQueue)
	getTransactionStatusResponse.OngoingIndicator = ongoingIndicator
	channel := NewMockWebSocket(wsId)

	handler := MockChargingStationTransactionHandler{}
	handler.On("OnGetTransactionStatus", mock.Anything).Return(getTransactionStatusResponse, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*transactions.GetTransactionStatusRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, transactionID, request.TransactionID)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.GetTransactionStatus(wsId, func(response *transactions.GetTransactionStatusResponse, err error) {
		require.Nil(t, err)
		require.NotNil(t, response)
		assert.Equal(t, messageInQueue, response.MessageInQueue)
		require.NotNil(t, response.OngoingIndicator)
		require.Equal(t, *ongoingIndicator, *response.OngoingIndicator)
		resultChannel <- true
	}, func(request *transactions.GetTransactionStatusRequest) {
		request.TransactionID = transactionID
	})
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV2TestSuite) TestGetTransactionStatusInvalidEndpoint() {
	messageId := defaultMessageId
	transactionID := "12345"
	getTransactionStatusRequest := transactions.NewGetTransactionStatusRequest()
	getTransactionStatusRequest.TransactionID = transactionID
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"transactionId":"%v"}]`, messageId, transactions.GetTransactionStatusFeatureName, transactionID)
	testUnsupportedRequestFromChargingStation(suite, getTransactionStatusRequest, requestJson, messageId)
}
