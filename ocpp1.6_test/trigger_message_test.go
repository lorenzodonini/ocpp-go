package ocpp16_test

import (
	"fmt"
	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test
func (suite *OcppV16TestSuite) TestTriggerMessageRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{ocpp16.TriggerMessageRequest{RequestedMessage: ocpp16.StatusNotificationFeatureName, ConnectorId: 1}, true},
		{ocpp16.TriggerMessageRequest{RequestedMessage: ocpp16.StatusNotificationFeatureName}, true},
		{ocpp16.TriggerMessageRequest{}, false},
		{ocpp16.TriggerMessageRequest{RequestedMessage: ocpp16.StatusNotificationFeatureName, ConnectorId: -1}, false},
		{ocpp16.TriggerMessageRequest{RequestedMessage: ocpp16.StartTransactionFeatureName}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV16TestSuite) TestTriggerMessageConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{ocpp16.TriggerMessageConfirmation{Status: ocpp16.TriggerMessageStatusAccepted}, true},
		{ocpp16.TriggerMessageConfirmation{Status: "invalidTriggerMessageStatus"}, false},
		{ocpp16.TriggerMessageConfirmation{}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV16TestSuite) TestTriggerMessageE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	connectorId := 1
	requestedMessage := ocpp16.MessageTrigger(ocpp16.StatusNotificationFeatureName)
	status := ocpp16.TriggerMessageStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"requestedMessage":"%v","connectorId":%v}]`, messageId, ocpp16.TriggerMessageFeatureName, requestedMessage, connectorId)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	TriggerMessageConfirmation := ocpp16.NewTriggerMessageConfirmation(status)
	channel := NewMockWebSocket(wsId)

	remoteTriggerListener := MockChargePointRemoteTriggerListener{}
	remoteTriggerListener.On("OnTriggerMessage", mock.Anything).Return(TriggerMessageConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*ocpp16.TriggerMessageRequest)
		assert.True(t, ok)
		assert.Equal(t, requestedMessage, request.RequestedMessage)
		assert.Equal(t, connectorId, request.ConnectorId)
	})
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	suite.chargePoint.SetRemoteTriggerListener(remoteTriggerListener)
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	assert.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.centralSystem.TriggerMessage(wsId, func(confirmation *ocpp16.TriggerMessageConfirmation, err error) {
		assert.Nil(t, err)
		assert.NotNil(t, confirmation)
		assert.Equal(t, status, confirmation.Status)
		resultChannel <- true
	}, requestedMessage, func(request *ocpp16.TriggerMessageRequest) {
		request.ConnectorId = connectorId
	})
	assert.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV16TestSuite) TestTriggerMessageInvalidEndpoint() {
	messageId := defaultMessageId
	connectorId := 1
	requestedMessage := ocpp16.MessageTrigger(ocpp16.StatusNotificationFeatureName)
	TriggerMessageRequest := ocpp16.NewTriggerMessageRequest(requestedMessage)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"requestedMessage":"%v","connectorId":%v}]`, messageId, ocpp16.TriggerMessageFeatureName, requestedMessage, connectorId)
	testUnsupportedRequestFromChargePoint(suite, TriggerMessageRequest, requestJson, messageId)
}
