package ocpp16_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (suite *OcppV16TestSuite) TestUnlockConnectorRequestValidation() {
	t := suite.T()
	var testTable = []GenericTestEntry{
		{core.UnlockConnectorRequest{ConnectorId: 1}, true},
		{core.UnlockConnectorRequest{ConnectorId: -1}, false},
		{core.UnlockConnectorRequest{}, false},
	}
	ExecuteGenericTestTable(t, testTable)
}

func (suite *OcppV16TestSuite) TestUnlockConnectorConfirmationValidation() {
	t := suite.T()
	var testTable = []GenericTestEntry{
		{core.UnlockConnectorConfirmation{Status: core.UnlockStatusUnlocked}, true},
		{core.UnlockConnectorConfirmation{Status: "invalidUnlockStatus"}, false},
		{core.UnlockConnectorConfirmation{}, false},
	}
	ExecuteGenericTestTable(t, testTable)
}

// Test
func (suite *OcppV16TestSuite) TestUnlockConnectorE2EMocked() {
	t := suite.T()
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
	coreListener := MockChargePointCoreListener{}
	coreListener.On("OnUnlockConnector", mock.Anything).Return(unlockConnectorConfirmation, nil)
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, coreListener, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	assert.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.centralSystem.UnlockConnector(wsId, func(confirmation *core.UnlockConnectorConfirmation, err error) {
		assert.NotNil(t, confirmation)
		assert.Nil(t, err)
		assert.Equal(t, status, confirmation.Status)
		resultChannel <- true
	}, connectorId)
	assert.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV16TestSuite) TestUnlockConnectorInvalidEndpoint() {
	messageId := defaultMessageId
	connectorId := 1
	changeAvailabilityRequest := core.NewUnlockConnectorRequest(connectorId)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"connectorId":%v}]`, messageId, core.UnlockConnectorFeatureName, connectorId)
	testUnsupportedRequestFromChargePoint(suite, changeAvailabilityRequest, requestJson, messageId)
}
