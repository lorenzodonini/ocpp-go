package ocpp16_test

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/stretchr/testify/mock"
)

func (suite *OcppV16TestSuite) TestUnlockConnectorRequestValidation() {
	var testTable = []GenericTestEntry{
		{core.UnlockConnectorRequest{ConnectorId: 1}, true},
		{core.UnlockConnectorRequest{ConnectorId: -1}, false},
		{core.UnlockConnectorRequest{}, false},
	}
	ExecuteGenericTestTable(suite, testTable)
}

func (suite *OcppV16TestSuite) TestUnlockConnectorConfirmationValidation() {
	var testTable = []GenericTestEntry{
		{core.UnlockConnectorConfirmation{Status: core.UnlockStatusUnlocked}, true},
		{core.UnlockConnectorConfirmation{Status: "invalidUnlockStatus"}, false},
		{core.UnlockConnectorConfirmation{}, false},
	}
	ExecuteGenericTestTable(suite, testTable)
}

// Test
func (suite *OcppV16TestSuite) TestUnlockConnectorE2EMocked() {
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	connectorId := 1
	status := core.UnlockStatusUnlocked
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"connectorId":%v}]`, messageId, core.UnlockConnectorFeatureName, connectorId)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	unlockConnectorConfirmation := core.NewUnlockConnectorConfirmation(status)
	channel := NewMockWebSocket(wsId)
	// Setting handlers
	coreListener := &MockChargePointCoreListener{}
	coreListener.On("OnUnlockConnector", mock.Anything).Return(unlockConnectorConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*core.UnlockConnectorRequest)
		suite.Require().NotNil(request)
		suite.Require().True(ok)
		suite.Equal(connectorId, request.ConnectorId)
	})
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, coreListener, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	suite.Require().Nil(err)
	resultChannel := make(chan bool, 1)
	err = suite.centralSystem.UnlockConnector(wsId, func(confirmation *core.UnlockConnectorConfirmation, err error) {
		suite.Require().NotNil(confirmation)
		suite.Require().Nil(err)
		suite.Equal(status, confirmation.Status)
		resultChannel <- true
	}, connectorId)
	suite.Require().Nil(err)
	result := <-resultChannel
	suite.True(result)
}

func (suite *OcppV16TestSuite) TestUnlockConnectorInvalidEndpoint() {
	messageId := defaultMessageId
	connectorId := 1
	changeAvailabilityRequest := core.NewUnlockConnectorRequest(connectorId)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"connectorId":%v}]`, messageId, core.UnlockConnectorFeatureName, connectorId)
	testUnsupportedRequestFromChargePoint(suite, changeAvailabilityRequest, requestJson, messageId)
}
