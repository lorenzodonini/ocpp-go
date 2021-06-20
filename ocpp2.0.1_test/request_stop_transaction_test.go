package ocpp2_test

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/remotecontrol"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test
func (suite *OcppV2TestSuite) TestRequestStopTransactionRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{remotecontrol.RequestStopTransactionRequest{TransactionID: "12345"}, true},
		{remotecontrol.RequestStopTransactionRequest{}, false},
		{remotecontrol.RequestStopTransactionRequest{TransactionID: ">36.................................."}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestRequestStopTransactionConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{remotecontrol.RequestStopTransactionResponse{Status: remotecontrol.RequestStartStopStatusAccepted, StatusInfo: &types.StatusInfo{ReasonCode: "200"}}, true},
		{remotecontrol.RequestStopTransactionResponse{Status: remotecontrol.RequestStartStopStatusAccepted}, true},
		{remotecontrol.RequestStopTransactionResponse{Status: remotecontrol.RequestStartStopStatusRejected}, true},
		{remotecontrol.RequestStopTransactionResponse{}, false},
		{remotecontrol.RequestStopTransactionResponse{Status: "invalidRequestStartStopStatus", StatusInfo: &types.StatusInfo{ReasonCode: "200"}}, false},
		{remotecontrol.RequestStopTransactionResponse{Status: remotecontrol.RequestStartStopStatusAccepted, StatusInfo: &types.StatusInfo{}}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestRequestStopTransactionE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	transactionId := "12345"
	status := remotecontrol.RequestStartStopStatusAccepted
	statusInfo := types.StatusInfo{ReasonCode: "200"}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"transactionId":"%v"}]`,
		messageId, remotecontrol.RequestStopTransactionFeatureName, transactionId)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v","statusInfo":{"reasonCode":"%v"}}]`,
		messageId, status, statusInfo.ReasonCode)
	RequestStopTransactionResponse := remotecontrol.NewRequestStopTransactionResponse(status)
	RequestStopTransactionResponse.StatusInfo = &statusInfo
	channel := NewMockWebSocket(wsId)

	handler := &MockChargingStationRemoteControlHandler{}
	handler.On("OnRequestStopTransaction", mock.Anything).Return(RequestStopTransactionResponse, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*remotecontrol.RequestStopTransactionRequest)
		require.True(t, ok)
		assert.Equal(t, transactionId, request.TransactionID)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.RequestStopTransaction(wsId, func(response *remotecontrol.RequestStopTransactionResponse, err error) {
		require.Nil(t, err)
		require.NotNil(t, response)
		assert.Equal(t, status, response.Status)
		assert.Equal(t, statusInfo.ReasonCode, response.StatusInfo.ReasonCode)
		resultChannel <- true
	}, transactionId)
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV2TestSuite) TestRequestStopTransactionInvalidEndpoint() {
	messageId := defaultMessageId
	transactionId := "12345"
	request := remotecontrol.RequestStopTransactionRequest{
		TransactionID: transactionId,
	}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"transactionId":"%v"}]`,
		messageId, remotecontrol.RequestStopTransactionFeatureName, transactionId)
	testUnsupportedRequestFromChargingStation(suite, request, requestJson, messageId)
}
