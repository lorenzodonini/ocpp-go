package ocpp2_test

import (
	"fmt"

	"github.com/stretchr/testify/mock"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/display"
)

// Test
func (suite *OcppV2TestSuite) TestClearDisplayMessageRequestValidation() {
	var requestTable = []GenericTestEntry{
		{display.ClearDisplayRequest{ID: 42}, true},
		{display.ClearDisplayRequest{}, true},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV2TestSuite) TestClearDisplayMessageResponseValidation() {
	var confirmationTable = []GenericTestEntry{
		{display.ClearDisplayResponse{Status: display.ClearMessageStatusAccepted}, true},
		{display.ClearDisplayResponse{Status: display.ClearMessageStatusUnknown}, true},
		{display.ClearDisplayResponse{Status: "invalidClearMessageStatus"}, false},
		{display.ClearDisplayResponse{}, false},
	}
	ExecuteGenericTestTable(suite, confirmationTable)
}

func (suite *OcppV2TestSuite) TestClearDisplayMessageE2EMocked() {
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	displayMessageId := 42
	status := display.ClearMessageStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"id":%v}]`, messageId, display.ClearDisplayMessageFeatureName, displayMessageId)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	clearDisplayConfirmation := display.NewClearDisplayResponse(status)
	channel := NewMockWebSocket(wsId)

	handler := &MockChargingStationDisplayHandler{}
	handler.On("OnClearDisplay", mock.Anything).Return(clearDisplayConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*display.ClearDisplayRequest)
		suite.Require().True(ok)
		suite.Require().NotNil(request)
		suite.Equal(displayMessageId, request.ID)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	suite.Require().Nil(err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.ClearDisplay(wsId, func(confirmation *display.ClearDisplayResponse, err error) {
		suite.Require().Nil(err)
		suite.Require().NotNil(confirmation)
		suite.Equal(status, confirmation.Status)
		resultChannel <- true
	}, displayMessageId)
	suite.Require().Nil(err)
	result := <-resultChannel
	suite.True(result)
}

func (suite *OcppV2TestSuite) TestClearDisplayMessageInvalidEndpoint() {
	messageId := defaultMessageId
	displayMessageId := 42
	clearDisplayRequest := display.NewClearDisplayRequest(displayMessageId)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"id":%v}]`, messageId, display.ClearDisplayMessageFeatureName, displayMessageId)
	testUnsupportedRequestFromChargingStation(suite, clearDisplayRequest, requestJson, messageId)
}
