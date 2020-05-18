package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/display"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test
func (suite *OcppV2TestSuite) TestClearDisplayRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{display.ClearDisplayRequest{ID: 42}, true},
		{display.ClearDisplayRequest{}, false},
		{display.ClearDisplayRequest{ID: -1}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestClearDisplayConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{display.ClearDisplayConfirmation{Status: display.ClearMessageStatusAccepted}, true},
		{display.ClearDisplayConfirmation{Status: display.ClearMessageStatusUnknown}, true},
		{display.ClearDisplayConfirmation{Status: "invalidClearMessageStatus"}, false},
		{display.ClearDisplayConfirmation{}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestClearDisplayE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	displayMessageId := 42
	status := display.ClearMessageStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"id":%v}]`, messageId, display.ClearDisplayFeatureName, displayMessageId)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	clearDisplayConfirmation := display.NewClearDisplayConfirmation(status)
	channel := NewMockWebSocket(wsId)

	handler := MockChargingStationDisplayHandler{}
	handler.On("OnClearDisplay", mock.Anything).Return(clearDisplayConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*display.ClearDisplayRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, displayMessageId, request.ID)
	})
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	suite.chargingStation.SetDisplayHandler(handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.ClearDisplay(wsId, func(confirmation *display.ClearDisplayConfirmation, err error) {
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
	clearDisplayRequest := display.NewClearDisplayRequest(displayMessageId)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"id":%v}]`, messageId, display.ClearDisplayFeatureName, displayMessageId)
	testUnsupportedRequestFromChargePoint(suite, clearDisplayRequest, requestJson, messageId)
}
