package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test
func (suite *OcppV2TestSuite) TestClearCacheRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{ocpp2.ClearCacheRequest{}, true},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestClearCacheConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{ocpp2.ClearCacheConfirmation{Status: ocpp2.ClearCacheStatusAccepted}, true},
		{ocpp2.ClearCacheConfirmation{Status: ocpp2.ClearCacheStatusRejected}, true},
		{ocpp2.ClearCacheConfirmation{Status: "invalidClearCacheStatus"}, false},
		{ocpp2.ClearCacheConfirmation{}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestClearCacheE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	status := ocpp2.ClearCacheStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{}]`, messageId, ocpp2.ClearCacheFeatureName)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	clearCacheConfirmation := ocpp2.NewClearCacheConfirmation(status)
	channel := NewMockWebSocket(wsId)

	coreListener := MockChargePointCoreListener{}
	coreListener.On("OnClearCache", mock.Anything).Return(clearCacheConfirmation, nil)
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, coreListener, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.ClearCache(wsId, func(confirmation *ocpp2.ClearCacheConfirmation, err error) {
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
	clearCacheRequest := ocpp2.NewClearCacheRequest()
	requestJson := fmt.Sprintf(`[2,"%v","%v",{}]`, messageId, ocpp2.ClearCacheFeatureName)
	testUnsupportedRequestFromChargePoint(suite, clearCacheRequest, requestJson, messageId)
}
