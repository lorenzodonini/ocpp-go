package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test
func (suite *OcppV2TestSuite) TestGetDisplayMessagesRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{ocpp2.GetDisplayMessagesRequest{RequestID: 1, Priority: ocpp2.MessagePriorityAlwaysFront, State: ocpp2.MessageStateCharging, ID: []int{2, 3}}, true},
		{ocpp2.GetDisplayMessagesRequest{RequestID: 1, Priority: ocpp2.MessagePriorityAlwaysFront, State: ocpp2.MessageStateCharging, ID: []int{}}, true},
		{ocpp2.GetDisplayMessagesRequest{RequestID: 1, Priority: ocpp2.MessagePriorityAlwaysFront, State: ocpp2.MessageStateCharging}, true},
		{ocpp2.GetDisplayMessagesRequest{RequestID: 1, Priority: ocpp2.MessagePriorityAlwaysFront}, true},
		{ocpp2.GetDisplayMessagesRequest{RequestID: 1, State: ocpp2.MessageStateCharging}, true},
		{ocpp2.GetDisplayMessagesRequest{RequestID: 1}, true},
		{ocpp2.GetDisplayMessagesRequest{}, true},
		{ocpp2.GetDisplayMessagesRequest{RequestID: -1}, false},
		{ocpp2.GetDisplayMessagesRequest{RequestID: 1, Priority: "invalidMessagePriority", State: ocpp2.MessageStateCharging, ID: []int{2, 3}}, false},
		{ocpp2.GetDisplayMessagesRequest{RequestID: 1, Priority: ocpp2.MessagePriorityAlwaysFront, State: "invalidMessageState", ID: []int{2, 3}}, false},
		{ocpp2.GetDisplayMessagesRequest{RequestID: 1, Priority: ocpp2.MessagePriorityAlwaysFront, State: ocpp2.MessageStateCharging, ID: []int{-2, 3}}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestGetDisplayMessagesConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{ocpp2.GetDisplayMessagesConfirmation{Status: ocpp2.MessageStatusAccepted}, true},
		{ocpp2.GetDisplayMessagesConfirmation{Status: ocpp2.MessageStatusUnknown}, true},
		{ocpp2.GetDisplayMessagesConfirmation{Status: "invalidMessageStatus"}, false},
		{ocpp2.GetDisplayMessagesConfirmation{}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestGetDisplayMessagesE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	requestId := 42
	messageIds := []int{2,3}
	priority := ocpp2.MessagePriorityInFront
	state := ocpp2.MessageStateCharging
	status := ocpp2.MessageStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"requestId":%v,"priority":"%v","state":"%v","id":[%v,%v]}]`,
		messageId, ocpp2.GetDisplayMessagesFeatureName, requestId, priority, state, messageIds[0], messageIds[1])
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	getDisplayMessagesConfirmation := ocpp2.NewGetDisplayMessagesConfirmation(status)
	channel := NewMockWebSocket(wsId)

	coreListener := MockChargePointCoreListener{}
	coreListener.On("OnGetDisplayMessages", mock.Anything).Return(getDisplayMessagesConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*ocpp2.GetDisplayMessagesRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, requestId, request.RequestID)
		assert.Equal(t, priority, request.Priority)
		assert.Equal(t, state, request.State)
		require.Len(t, request.ID, len(messageIds))
		assert.Equal(t, messageIds[0], request.ID[0])
		assert.Equal(t, messageIds[1], request.ID[1])
	})
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, coreListener, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.GetDisplayMessages(wsId, func(confirmation *ocpp2.GetDisplayMessagesConfirmation, err error) {
		require.Nil(t, err)
		require.NotNil(t, confirmation)
		assert.Equal(t, status, confirmation.Status)
		resultChannel <- true
	}, requestId, func(request *ocpp2.GetDisplayMessagesRequest) {
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
	messageIds := []int{2,3}
	priority := ocpp2.MessagePriorityInFront
	state := ocpp2.MessageStateCharging
	getDisplayMessagesRequest := ocpp2.NewGetDisplayMessagesRequest(requestId)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"requestId":%v,"priority":"%v","state":"%v","id":[%v,%v]}]`,
		messageId, ocpp2.GetDisplayMessagesFeatureName, requestId, priority, state, messageIds[0], messageIds[1])
	testUnsupportedRequestFromChargePoint(suite, getDisplayMessagesRequest, requestJson, messageId)
}
