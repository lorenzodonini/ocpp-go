package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/display"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test
func (suite *OcppV2TestSuite) TestGetDisplayMessagesRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{display.GetDisplayMessagesRequest{RequestID: 1, Priority: display.MessagePriorityAlwaysFront, State: display.MessageStateCharging, ID: []int{2, 3}}, true},
		{display.GetDisplayMessagesRequest{RequestID: 1, Priority: display.MessagePriorityAlwaysFront, State: display.MessageStateCharging, ID: []int{}}, true},
		{display.GetDisplayMessagesRequest{RequestID: 1, Priority: display.MessagePriorityAlwaysFront, State: display.MessageStateCharging}, true},
		{display.GetDisplayMessagesRequest{RequestID: 1, Priority: display.MessagePriorityAlwaysFront}, true},
		{display.GetDisplayMessagesRequest{RequestID: 1, State: display.MessageStateCharging}, true},
		{display.GetDisplayMessagesRequest{RequestID: 1}, true},
		{display.GetDisplayMessagesRequest{}, true},
		{display.GetDisplayMessagesRequest{RequestID: -1}, false},
		{display.GetDisplayMessagesRequest{RequestID: 1, Priority: "invalidMessagePriority", State: display.MessageStateCharging, ID: []int{2, 3}}, false},
		{display.GetDisplayMessagesRequest{RequestID: 1, Priority: display.MessagePriorityAlwaysFront, State: "invalidMessageState", ID: []int{2, 3}}, false},
		{display.GetDisplayMessagesRequest{RequestID: 1, Priority: display.MessagePriorityAlwaysFront, State: display.MessageStateCharging, ID: []int{-2, 3}}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestGetDisplayMessagesConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{display.GetDisplayMessagesConfirmation{Status: display.MessageStatusAccepted}, true},
		{display.GetDisplayMessagesConfirmation{Status: display.MessageStatusUnknown}, true},
		{display.GetDisplayMessagesConfirmation{Status: "invalidMessageStatus"}, false},
		{display.GetDisplayMessagesConfirmation{}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestGetDisplayMessagesE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	requestId := 42
	messageIds := []int{2, 3}
	priority := display.MessagePriorityInFront
	state := display.MessageStateCharging
	status := display.MessageStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"requestId":%v,"priority":"%v","state":"%v","id":[%v,%v]}]`,
		messageId, display.GetDisplayMessagesFeatureName, requestId, priority, state, messageIds[0], messageIds[1])
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	getDisplayMessagesConfirmation := display.NewGetDisplayMessagesConfirmation(status)
	channel := NewMockWebSocket(wsId)

	handler := MockChargingStationDisplayHandler{}
	handler.On("OnGetDisplayMessages", mock.Anything).Return(getDisplayMessagesConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*display.GetDisplayMessagesRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, requestId, request.RequestID)
		assert.Equal(t, priority, request.Priority)
		assert.Equal(t, state, request.State)
		require.Len(t, request.ID, len(messageIds))
		assert.Equal(t, messageIds[0], request.ID[0])
		assert.Equal(t, messageIds[1], request.ID[1])
	})
	setupDefaultCSMSHandlers(suite, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.GetDisplayMessages(wsId, func(confirmation *display.GetDisplayMessagesConfirmation, err error) {
		require.Nil(t, err)
		require.NotNil(t, confirmation)
		assert.Equal(t, status, confirmation.Status)
		resultChannel <- true
	}, requestId, func(request *display.GetDisplayMessagesRequest) {
		request.Priority = priority
		request.State = state
		request.ID = messageIds
	})
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV2TestSuite) TestGetDisplayMessagesInvalidEndpoint() {
	messageId := defaultMessageId
	requestId := 42
	messageIds := []int{2, 3}
	priority := display.MessagePriorityInFront
	state := display.MessageStateCharging
	getDisplayMessagesRequest := display.NewGetDisplayMessagesRequest(requestId)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"requestId":%v,"priority":"%v","state":"%v","id":[%v,%v]}]`,
		messageId, display.GetDisplayMessagesFeatureName, requestId, priority, state, messageIds[0], messageIds[1])
	testUnsupportedRequestFromChargingStation(suite, getDisplayMessagesRequest, requestJson, messageId)
}
