package ocpp16_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (suite *OcppV16TestSuite) TestResetRequestValidation() {
	t := suite.T()
	var testTable = []GenericTestEntry{
		{core.ResetRequest{Type: core.ResetTypeHard}, true},
		{core.ResetRequest{Type: core.ResetTypeSoft}, true},
		{core.ResetRequest{Type: "invalidResetType"}, false},
		{core.ResetRequest{}, false},
	}
	ExecuteGenericTestTable(t, testTable)
}

func (suite *OcppV16TestSuite) TestResetConfirmationValidation() {
	t := suite.T()
	var testTable = []GenericTestEntry{
		{core.ResetConfirmation{Status: core.ResetStatusAccepted}, true},
		{core.ResetConfirmation{Status: core.ResetStatusRejected}, true},
		{core.ResetConfirmation{Status: "invalidResetStatus"}, false},
		{core.ResetConfirmation{}, false},
	}
	ExecuteGenericTestTable(t, testTable)
}

// Test
func (suite *OcppV16TestSuite) TestResetE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	resetType := core.ResetTypeSoft
	status := core.ResetStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"type":"%v"}]`, messageId, core.ResetFeatureName, resetType)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	resetConfirmation := core.NewResetConfirmation(status)
	channel := NewMockWebSocket(wsId)
	// Setting handlers
	coreListener := MockChargePointCoreListener{}
	coreListener.On("OnReset", mock.Anything).Return(resetConfirmation, nil)
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, coreListener, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	assert.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.centralSystem.Reset(wsId, func(confirmation *core.ResetConfirmation, err error) {
		assert.NotNil(t, confirmation)
		assert.Nil(t, err)
		assert.Equal(t, status, confirmation.Status)
		resultChannel <- true
	}, resetType)
	assert.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV16TestSuite) TestResetInvalidEndpoint() {
	messageId := defaultMessageId
	resetType := core.ResetTypeSoft
	resetRequest := core.NewResetRequest(resetType)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"type":"%v"}]`, messageId, core.ResetFeatureName, resetType)
	testUnsupportedRequestFromChargePoint(suite, resetRequest, requestJson, messageId)
}
