package ocpp16_test

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/remotetrigger"
	"github.com/stretchr/testify/mock"
)

// Test
func (suite *OcppV16TestSuite) TestTriggerMessageRequestValidation() {
	requestTable := []GenericTestEntry{
		{remotetrigger.TriggerMessageRequest{RequestedMessage: core.StatusNotificationFeatureName, ConnectorId: newInt(1)}, true},
		{remotetrigger.TriggerMessageRequest{RequestedMessage: core.StatusNotificationFeatureName}, true},
		{remotetrigger.TriggerMessageRequest{}, false},
		{remotetrigger.TriggerMessageRequest{RequestedMessage: core.StatusNotificationFeatureName, ConnectorId: newInt(0)}, true},
		{remotetrigger.TriggerMessageRequest{RequestedMessage: core.StatusNotificationFeatureName, ConnectorId: newInt(-1)}, false},
		{remotetrigger.TriggerMessageRequest{RequestedMessage: core.StartTransactionFeatureName}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV16TestSuite) TestTriggerMessageConfirmationValidation() {
	confirmationTable := []GenericTestEntry{
		{remotetrigger.TriggerMessageConfirmation{Status: remotetrigger.TriggerMessageStatusAccepted}, true},
		{remotetrigger.TriggerMessageConfirmation{Status: "invalidTriggerMessageStatus"}, false},
		{remotetrigger.TriggerMessageConfirmation{}, false},
	}
	ExecuteGenericTestTable(suite, confirmationTable)
}

func (suite *OcppV16TestSuite) TestTriggerMessageE2EMocked() {
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	connectorId := newInt(1)
	requestedMessage := remotetrigger.MessageTrigger(core.StatusNotificationFeatureName)
	status := remotetrigger.TriggerMessageStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"requestedMessage":"%v","connectorId":%v}]`, messageId, remotetrigger.TriggerMessageFeatureName, requestedMessage, *connectorId)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	TriggerMessageConfirmation := remotetrigger.NewTriggerMessageConfirmation(status)
	channel := NewMockWebSocket(wsId)

	remoteTriggerListener := &MockChargePointRemoteTriggerListener{}
	remoteTriggerListener.On("OnTriggerMessage", mock.Anything).Return(TriggerMessageConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*remotetrigger.TriggerMessageRequest)
		suite.Require().True(ok)
		suite.Require().NotNil(request)
		suite.Equal(requestedMessage, request.RequestedMessage)
		suite.Equal(*connectorId, *request.ConnectorId)
	})
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	suite.chargePoint.SetRemoteTriggerHandler(remoteTriggerListener)
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	suite.Require().Nil(err)
	resultChannel := make(chan bool, 1)
	err = suite.centralSystem.TriggerMessage(wsId, func(confirmation *remotetrigger.TriggerMessageConfirmation, err error) {
		suite.Require().Nil(err)
		suite.Require().NotNil(confirmation)
		suite.Equal(status, confirmation.Status)
		resultChannel <- true
	}, requestedMessage, func(request *remotetrigger.TriggerMessageRequest) {
		request.ConnectorId = connectorId
	})
	suite.Require().Nil(err)
	result := <-resultChannel
	suite.True(result)
}

func (suite *OcppV16TestSuite) TestTriggerMessageInvalidEndpoint() {
	messageId := defaultMessageId
	connectorId := 1
	requestedMessage := remotetrigger.MessageTrigger(core.StatusNotificationFeatureName)
	TriggerMessageRequest := remotetrigger.NewTriggerMessageRequest(requestedMessage)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"requestedMessage":"%v","connectorId":%v}]`, messageId, remotetrigger.TriggerMessageFeatureName, requestedMessage, connectorId)
	testUnsupportedRequestFromChargePoint(suite, TriggerMessageRequest, requestJson, messageId)
}
