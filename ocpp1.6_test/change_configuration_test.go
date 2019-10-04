package ocpp16_test

import (
	"fmt"
	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test
func (suite *OcppV16TestSuite) TestChangeConfigurationRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{ocpp16.ChangeConfigurationRequest{Key: "someKey", Value: "someValue"}, true},
		{ocpp16.ChangeConfigurationRequest{Key: "someKey"}, false},
		{ocpp16.ChangeConfigurationRequest{Value: "someValue"}, false},
		{ocpp16.ChangeConfigurationRequest{}, false},
		{ocpp16.ChangeConfigurationRequest{Key: ">50................................................", Value: "someValue"}, false},
		{ocpp16.ChangeConfigurationRequest{Key: "someKey", Value: ">500................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................."}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV16TestSuite) TestChangeConfigurationConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{ocpp16.ChangeConfigurationConfirmation{Status: ocpp16.ConfigurationStatusAccepted}, true},
		{ocpp16.ChangeConfigurationConfirmation{Status: ocpp16.ConfigurationStatusRejected}, true},
		{ocpp16.ChangeConfigurationConfirmation{Status: ocpp16.ConfigurationStatusRebootRequired}, true},
		{ocpp16.ChangeConfigurationConfirmation{Status: ocpp16.ConfigurationStatusNotSupported}, true},
		{ocpp16.ChangeConfigurationConfirmation{Status: "invalidConfigurationStatus"}, false},
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
	status := ocpp16.ConfigurationStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"key":"%v","value":"%v"}]`, messageId, ocpp16.ChangeConfigurationFeatureName, key, value)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	changeConfigurationConfirmation := ocpp16.NewChangeConfigurationConfirmation(status)
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
	err = suite.centralSystem.ChangeConfiguration(wsId, func(confirmation *ocpp16.ChangeConfigurationConfirmation, err error) {
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
	changeConfigurationRequest := ocpp16.NewChangeConfigurationRequest(key, value)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"key":"%v","value":"%v"}]`, messageId, ocpp16.ChangeConfigurationFeatureName, key, value)
	testUnsupportedRequestFromChargePoint(suite, changeConfigurationRequest, requestJson, messageId)
}
