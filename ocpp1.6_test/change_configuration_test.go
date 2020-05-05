package ocpp16_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test
func (suite *OcppV16TestSuite) TestChangeConfigurationRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{core.ChangeConfigurationRequest{Key: "someKey", Value: "someValue"}, true},
		{core.ChangeConfigurationRequest{Key: "someKey"}, false},
		{core.ChangeConfigurationRequest{Value: "someValue"}, false},
		{core.ChangeConfigurationRequest{}, false},
		{core.ChangeConfigurationRequest{Key: ">50................................................", Value: "someValue"}, false},
		{core.ChangeConfigurationRequest{Key: "someKey", Value: ">500................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................."}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV16TestSuite) TestChangeConfigurationConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{core.ChangeConfigurationConfirmation{Status: core.ConfigurationStatusAccepted}, true},
		{core.ChangeConfigurationConfirmation{Status: core.ConfigurationStatusRejected}, true},
		{core.ChangeConfigurationConfirmation{Status: core.ConfigurationStatusRebootRequired}, true},
		{core.ChangeConfigurationConfirmation{Status: core.ConfigurationStatusNotSupported}, true},
		{core.ChangeConfigurationConfirmation{Status: "invalidConfigurationStatus"}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV16TestSuite) TestChangeConfigurationE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	key := "someKey"
	value := "someValue"
	status := core.ConfigurationStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"key":"%v","value":"%v"}]`, messageId, core.ChangeConfigurationFeatureName, key, value)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	changeConfigurationConfirmation := core.NewChangeConfigurationConfirmation(status)
	channel := NewMockWebSocket(wsId)

	coreListener := MockChargePointCoreListener{}
	coreListener.On("OnChangeConfiguration", mock.Anything).Return(changeConfigurationConfirmation, nil)
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, coreListener, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	assert.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.centralSystem.ChangeConfiguration(wsId, func(confirmation *core.ChangeConfigurationConfirmation, err error) {
		assert.Nil(t, err)
		assert.NotNil(t, confirmation)
		assert.Equal(t, status, confirmation.Status)
		resultChannel <- true
	}, key, value)
	assert.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV16TestSuite) TestChangeConfigurationInvalidEndpoint() {
	messageId := defaultMessageId
	key := "someKey"
	value := "someValue"
	changeConfigurationRequest := core.NewChangeConfigurationRequest(key, value)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"key":"%v","value":"%v"}]`, messageId, core.ChangeConfigurationFeatureName, key, value)
	testUnsupportedRequestFromChargePoint(suite, changeConfigurationRequest, requestJson, messageId)
}
