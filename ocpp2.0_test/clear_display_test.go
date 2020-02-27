package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test
func (suite *OcppV2TestSuite) TestClearDisplayRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{ocpp2.ClearDisplayRequest{ID: 42}, true},
		{ocpp2.ClearDisplayRequest{}, false},
		{ocpp2.ClearDisplayRequest{ID: -1}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestClearDisplayConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{ocpp2.ClearDisplayConfirmation{Status: ocpp2.ClearMessageStatusAccepted}, true},
		{ocpp2.ClearDisplayConfirmation{Status: ocpp2.ClearMessageStatusUnknown}, true},
		{ocpp2.ClearDisplayConfirmation{Status: "invalidClearMessageStatus"}, false},
		{ocpp2.ClearDisplayConfirmation{}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestClearDisplayE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	displayMessageId := 42
	status := ocpp2.ClearMessageStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"id":%v}]`, messageId, ocpp2.ClearDisplayFeatureName, displayMessageId)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	ClearDisplayConfirmation := ocpp2.NewClearDisplayConfirmation(status)
	channel := NewMockWebSocket(wsId)

	coreListener := MockChargePointCoreListener{}
	coreListener.On("OnClearDisplay", mock.Anything).Return(ClearDisplayConfirmation, nil)
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, coreListener, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.ClearDisplay(wsId, func(confirmation *ocpp2.ClearDisplayConfirmation, err error) {
		require.Nil(t, err)
		require.NotNil(t, confirmation)
		assert.Equal(t, status, confirmation.Status)
		resultChannel <- true
	}, displayMessageId)
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV2TestSuite) TestClearDisplayInvalidEndpoint() {
	messageId := defaultMessageId
	displayMessageId := 42
	ClearDisplayRequest := ocpp2.NewClearDisplayRequest(displayMessageId)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"id":%v}]`, messageId, ocpp2.ClearDisplayFeatureName, displayMessageId)
	testUnsupportedRequestFromChargePoint(suite, ClearDisplayRequest, requestJson, messageId)
}
