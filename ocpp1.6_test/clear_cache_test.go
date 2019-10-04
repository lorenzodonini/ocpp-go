package ocpp16_test

import (
	"fmt"
	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test
func (suite *OcppV16TestSuite) TestClearCacheRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{ocpp16.ClearCacheRequest{}, true},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV16TestSuite) TestClearCacheConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{ocpp16.ClearCacheConfirmation{Status: ocpp16.ClearCacheStatusAccepted}, true},
		{ocpp16.ClearCacheConfirmation{Status: ocpp16.ClearCacheStatusRejected}, true},
		{ocpp16.ClearCacheConfirmation{Status: "invalidClearCacheStatus"}, false},
		{ocpp16.ClearCacheConfirmation{}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV16TestSuite) TestClearCacheE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	status := ocpp16.ClearCacheStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{}]`, messageId, ocpp16.ClearCacheFeatureName)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	clearCacheConfirmation := ocpp16.NewClearCacheConfirmation(status)
	channel := NewMockWebSocket(wsId)

	coreListener := MockChargePointCoreListener{}
	coreListener.On("OnClearCache", mock.Anything).Return(clearCacheConfirmation, nil)
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, coreListener, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	assert.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.centralSystem.ClearCache(wsId, func(confirmation *ocpp16.ClearCacheConfirmation, err error) {
		assert.Nil(t, err)
		assert.NotNil(t, confirmation)
		assert.Equal(t, status, confirmation.Status)
		resultChannel <- true
	})
	assert.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV16TestSuite) TestClearCacheInvalidEndpoint() {
	messageId := defaultMessageId
	clearCacheRequest := ocpp16.NewClearCacheRequest()
	requestJson := fmt.Sprintf(`[2,"%v","%v",{}]`, messageId, ocpp16.ClearCacheFeatureName)
	testUnsupportedRequestFromChargePoint(suite, clearCacheRequest, requestJson, messageId)
}
