package ocpp16_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test
func (suite *OcppV16TestSuite) TestRemoteStopTransactionRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{core.RemoteStopTransactionRequest{TransactionId: 1}, true},
		{core.RemoteStopTransactionRequest{}, true},
		{core.RemoteStopTransactionRequest{TransactionId: -1}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV16TestSuite) TestRemoteStopTransactionConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{core.RemoteStopTransactionConfirmation{Status: types.RemoteStartStopStatusAccepted}, true},
		{core.RemoteStopTransactionConfirmation{Status: types.RemoteStartStopStatusRejected}, true},
		{core.RemoteStopTransactionConfirmation{Status: "invalidRemoteStopTransactionStatus"}, false},
		{core.RemoteStopTransactionConfirmation{}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV16TestSuite) TestRemoteStopTransactionE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	transactionId := 1
	status := types.RemoteStartStopStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"transactionId":%v}]`, messageId, core.RemoteStopTransactionFeatureName, transactionId)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	RemoteStopTransactionConfirmation := core.NewRemoteStopTransactionConfirmation(status)
	channel := NewMockWebSocket(wsId)

	coreListener := MockChargePointCoreListener{}
	coreListener.On("OnRemoteStopTransaction", mock.Anything).Return(RemoteStopTransactionConfirmation, nil)
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, coreListener, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	assert.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.centralSystem.RemoteStopTransaction(wsId, func(confirmation *core.RemoteStopTransactionConfirmation, err error) {
		assert.Nil(t, err)
		assert.NotNil(t, confirmation)
		assert.Equal(t, status, confirmation.Status)
		resultChannel <- true
	}, transactionId)
	assert.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV16TestSuite) TestRemoteStopTransactionInvalidEndpoint() {
	messageId := defaultMessageId
	transactionId := 1
	RemoteStopTransactionRequest := core.NewRemoteStopTransactionRequest(transactionId)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"transactionId":%v}]`, messageId, core.RemoteStopTransactionFeatureName, transactionId)
	testUnsupportedRequestFromChargePoint(suite, RemoteStopTransactionRequest, requestJson, messageId)
}
