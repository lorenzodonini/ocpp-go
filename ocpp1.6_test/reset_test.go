package ocpp16_test

import (
	"fmt"
	"github.com/lorenzodonini/go-ocpp/ocpp1.6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (suite *OcppV16TestSuite) TestResetRequestValidation() {
	t := suite.T()
	var testTable = []GenericTestEntry{
		{ocpp16.ResetRequest{Type: ocpp16.ResetTypeHard}, true},
		{ocpp16.ResetRequest{Type: ocpp16.ResetTypeSoft}, true},
		{ocpp16.ResetRequest{Type: "invalidResetType"}, false},
		{ocpp16.ResetRequest{}, false},
	}
	ExecuteGenericTestTable(t, testTable)
}

func (suite *OcppV16TestSuite) TestResetConfirmationValidation() {
	t := suite.T()
	var testTable = []GenericTestEntry{
		{ocpp16.ResetConfirmation{Status: ocpp16.ResetStatusAccepted}, true},
		{ocpp16.ResetConfirmation{Status: ocpp16.ResetStatusRejected}, true},
		{ocpp16.ResetConfirmation{Status: "invalidResetStatus"}, false},
		{ocpp16.ResetConfirmation{}, false},
	}
	ExecuteGenericTestTable(t, testTable)
}

// Test
func (suite *OcppV16TestSuite) TestResetE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	resetType := ocpp16.ResetTypeSoft
	status := ocpp16.ResetStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"type":"%v"}]`, messageId, ocpp16.ResetFeatureName, resetType)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	resetConfirmation := ocpp16.NewResetConfirmation(status)
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
	err = suite.centralSystem.Reset(wsId, func(confirmation *ocpp16.ResetConfirmation, err error) {
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
	resetType := ocpp16.ResetTypeSoft
	resetRequest := ocpp16.NewResetRequest(resetType)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"type":"%v"}]`, messageId, ocpp16.ResetFeatureName, resetType)
	testUnsupportedRequestFromChargePoint(suite, resetRequest, requestJson, messageId)
}
