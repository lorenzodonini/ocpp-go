package ocpp16_test

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/stretchr/testify/mock"
)

// Test
func (suite *OcppV16TestSuite) TestClearCacheRequestValidation() {
	var requestTable = []GenericTestEntry{
		{core.ClearCacheRequest{}, true},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV16TestSuite) TestClearCacheConfirmationValidation() {
	var confirmationTable = []GenericTestEntry{
		{core.ClearCacheConfirmation{Status: core.ClearCacheStatusAccepted}, true},
		{core.ClearCacheConfirmation{Status: core.ClearCacheStatusRejected}, true},
		{core.ClearCacheConfirmation{Status: "invalidClearCacheStatus"}, false},
		{core.ClearCacheConfirmation{}, false},
	}
	ExecuteGenericTestTable(suite, confirmationTable)
}

func (suite *OcppV16TestSuite) TestClearCacheE2EMocked() {
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	status := core.ClearCacheStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{}]`, messageId, core.ClearCacheFeatureName)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	clearCacheConfirmation := core.NewClearCacheConfirmation(status)
	channel := NewMockWebSocket(wsId)

	coreListener := &MockChargePointCoreListener{}
	coreListener.On("OnClearCache", mock.Anything).Return(clearCacheConfirmation, nil)
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, coreListener, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	suite.Require().Nil(err)
	resultChannel := make(chan bool, 1)
	err = suite.centralSystem.ClearCache(wsId, func(confirmation *core.ClearCacheConfirmation, err error) {
		suite.Require().Nil(err)
		suite.Require().NotNil(confirmation)
		suite.Equal(status, confirmation.Status)
		resultChannel <- true
	})
	suite.Require().Nil(err)
	result := <-resultChannel
	suite.True(result)
}

func (suite *OcppV16TestSuite) TestClearCacheInvalidEndpoint() {
	messageId := defaultMessageId
	clearCacheRequest := core.NewClearCacheRequest()
	requestJson := fmt.Sprintf(`[2,"%v","%v",{}]`, messageId, core.ClearCacheFeatureName)
	testUnsupportedRequestFromChargePoint(suite, clearCacheRequest, requestJson, messageId)
}
