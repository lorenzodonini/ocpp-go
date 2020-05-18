package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/authorization"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test
func (suite *OcppV2TestSuite) TestClearCacheRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{authorization.ClearCacheRequest{}, true},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestClearCacheConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{authorization.ClearCacheConfirmation{Status: authorization.ClearCacheStatusAccepted}, true},
		{authorization.ClearCacheConfirmation{Status: authorization.ClearCacheStatusRejected}, true},
		{authorization.ClearCacheConfirmation{Status: "invalidClearCacheStatus"}, false},
		{authorization.ClearCacheConfirmation{}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestClearCacheE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	status := authorization.ClearCacheStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{}]`, messageId, authorization.ClearCacheFeatureName)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	clearCacheConfirmation := authorization.NewClearCacheConfirmation(status)
	channel := NewMockWebSocket(wsId)

	handler := MockChargingStationAuthorizationHandler{}
	handler.On("OnClearCache", mock.Anything).Return(clearCacheConfirmation, nil)
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	suite.chargingStation.SetAuthorizationHandler(handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.ClearCache(wsId, func(confirmation *authorization.ClearCacheConfirmation, err error) {
		require.Nil(t, err)
		require.NotNil(t, confirmation)
		assert.Equal(t, status, confirmation.Status)
		resultChannel <- true
	})
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV2TestSuite) TestClearCacheInvalidEndpoint() {
	messageId := defaultMessageId
	clearCacheRequest := authorization.NewClearCacheRequest()
	requestJson := fmt.Sprintf(`[2,"%v","%v",{}]`, messageId, authorization.ClearCacheFeatureName)
	testUnsupportedRequestFromChargePoint(suite, clearCacheRequest, requestJson, messageId)
}
