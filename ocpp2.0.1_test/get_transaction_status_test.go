package ocpp2_test

import (
	"fmt"

	"github.com/stretchr/testify/mock"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/transactions"
)

// Test
func (suite *OcppV2TestSuite) TestGetTransactionStatusRequestValidation() {
	var requestTable = []GenericTestEntry{
		{transactions.GetTransactionStatusRequest{}, true},
		{transactions.GetTransactionStatusRequest{TransactionID: "12345"}, true},
		{transactions.GetTransactionStatusRequest{TransactionID: ">36.................................."}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV2TestSuite) TestGetTransactionStatusResponseValidation() {
	var confirmationTable = []GenericTestEntry{
		{transactions.GetTransactionStatusResponse{OngoingIndicator: newBool(true), MessagesInQueue: true}, true},
		{transactions.GetTransactionStatusResponse{MessagesInQueue: true}, true},
		{transactions.GetTransactionStatusResponse{}, true},
	}
	ExecuteGenericTestTable(suite, confirmationTable)
}

func (suite *OcppV2TestSuite) TestGetTransactionStatusE2EMocked() {
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	transactionID := "12345"
	messagesInQueue := false
	ongoingIndicator := newBool(true)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"transactionId":"%v"}]`, messageId, transactions.GetTransactionStatusFeatureName, transactionID)
	responseJson := fmt.Sprintf(`[3,"%v",{"ongoingIndicator":%v,"messagesInQueue":%v}]`, messageId, *ongoingIndicator, messagesInQueue)
	getTransactionStatusResponse := transactions.NewGetTransactionStatusResponse(messagesInQueue)
	getTransactionStatusResponse.OngoingIndicator = ongoingIndicator
	channel := NewMockWebSocket(wsId)

	handler := &MockChargingStationTransactionHandler{}
	handler.On("OnGetTransactionStatus", mock.Anything).Return(getTransactionStatusResponse, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*transactions.GetTransactionStatusRequest)
		suite.Require().True(ok)
		suite.Require().NotNil(request)
		suite.Equal(transactionID, request.TransactionID)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	suite.Require().Nil(err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.GetTransactionStatus(wsId, func(response *transactions.GetTransactionStatusResponse, err error) {
		suite.Require().Nil(err)
		suite.Require().NotNil(response)
		suite.Equal(messagesInQueue, response.MessagesInQueue)
		suite.Require().NotNil(response.OngoingIndicator)
		suite.Require().Equal(*ongoingIndicator, *response.OngoingIndicator)
		resultChannel <- true
	}, func(request *transactions.GetTransactionStatusRequest) {
		request.TransactionID = transactionID
	})
	suite.Require().Nil(err)
	result := <-resultChannel
	suite.True(result)
}

func (suite *OcppV2TestSuite) TestGetTransactionStatusInvalidEndpoint() {
	messageId := defaultMessageId
	transactionID := "12345"
	getTransactionStatusRequest := transactions.NewGetTransactionStatusRequest()
	getTransactionStatusRequest.TransactionID = transactionID
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"transactionId":"%v"}]`, messageId, transactions.GetTransactionStatusFeatureName, transactionID)
	testUnsupportedRequestFromChargingStation(suite, getTransactionStatusRequest, requestJson, messageId)
}
