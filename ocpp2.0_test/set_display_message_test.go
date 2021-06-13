package ocpp2_test

import (
	"fmt"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0/display"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test
func (suite *OcppV2TestSuite) TestSetDisplayMessageRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{display.SetDisplayMessageRequest{Message: display.MessageInfo{ID: 42, Priority: display.MessagePriorityAlwaysFront, State: display.MessageStateIdle, StartDateTime: types.NewDateTime(time.Now()), Message: types.MessageContent{Format: types.MessageFormatUTF8, Content: "hello world"}}}, true},
		{display.SetDisplayMessageRequest{}, false},
		{display.SetDisplayMessageRequest{Message: display.MessageInfo{ID: 42, Priority: "invalidPriority", State: display.MessageStateIdle, StartDateTime: types.NewDateTime(time.Now()), Message: types.MessageContent{Format: types.MessageFormatUTF8, Content: "hello world"}}}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestSetDisplayMessageConfirmationValidation() {
	t := suite.T()
	var responseTable = []GenericTestEntry{
		{display.SetDisplayMessageResponse{Status: display.DisplayMessageStatusAccepted, StatusInfo: types.NewStatusInfo("200", "")}, true},
		{display.SetDisplayMessageResponse{Status: display.DisplayMessageStatusAccepted}, true},
		{display.SetDisplayMessageResponse{Status: display.DisplayMessageStatusNotSupportedMessageFormat}, true},
		{display.SetDisplayMessageResponse{Status: display.DisplayMessageStatusNotSupportedState}, true},
		{display.SetDisplayMessageResponse{Status: display.DisplayMessageStatusNotSupportedPriority}, true},
		{display.SetDisplayMessageResponse{Status: display.DisplayMessageStatusRejected}, true},
		{display.SetDisplayMessageResponse{Status: display.DisplayMessageStatusUnknownTransaction}, true},
		{display.SetDisplayMessageResponse{Status: "invalidDisplayMessageStatus"}, false},
		{display.SetDisplayMessageResponse{}, false},
	}
	ExecuteGenericTestTable(t, responseTable)
}

func (suite *OcppV2TestSuite) TestSetDisplayMessageE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	message := display.MessageInfo{
		ID:            42,
		Priority:      display.MessagePriorityAlwaysFront,
		State:         display.MessageStateIdle,
		StartDateTime: types.NewDateTime(time.Now()),
		Message: types.MessageContent{
			Format:  types.MessageFormatUTF8,
			Content: "hello world",
		},
	}
	status := display.DisplayMessageStatusAccepted
	statusInfo := types.NewStatusInfo("200", "")
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"message":{"id":%v,"priority":"%v","state":"%v","startDateTime":"%v","message":{"format":"%v","content":"%v"}}}]`,
		messageId, display.SetDisplayMessageFeatureName, message.ID, message.Priority, message.State, message.StartDateTime.FormatTimestamp(), message.Message.Format, message.Message.Content)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v","statusInfo":{"reasonCode":"%v"}}]`,
		messageId, status, statusInfo.ReasonCode)
	setDisplayResponse := display.NewSetDisplayMessageResponse(status)
	setDisplayResponse.StatusInfo = statusInfo
	channel := NewMockWebSocket(wsId)

	handler := MockChargingStationDisplayHandler{}
	handler.On("OnSetDisplayMessage", mock.Anything).Return(setDisplayResponse, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*display.SetDisplayMessageRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, message.ID, request.Message.ID)
		assert.Equal(t, message.Priority, request.Message.Priority)
		assert.Equal(t, message.State, request.Message.State)
		assertDateTimeEquality(t, message.StartDateTime, request.Message.StartDateTime)
		assert.Equal(t, message.Message.Format, request.Message.Message.Format)
		assert.Equal(t, message.Message.Content, request.Message.Message.Content)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.SetDisplayMessage(wsId, func(response *display.SetDisplayMessageResponse, err error) {
		require.Nil(t, err)
		require.NotNil(t, response)
		assert.Equal(t, status, response.Status)
		assert.Equal(t, statusInfo.ReasonCode, response.StatusInfo.ReasonCode)
		assert.Equal(t, statusInfo.AdditionalInfo, response.StatusInfo.AdditionalInfo)
		resultChannel <- true
	}, message)
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV2TestSuite) TestSetDisplayMessageInvalidEndpoint() {
	messageId := defaultMessageId
	message := display.MessageInfo{
		ID:            42,
		Priority:      display.MessagePriorityAlwaysFront,
		State:         display.MessageStateIdle,
		StartDateTime: types.NewDateTime(time.Now()),
		Message: types.MessageContent{
			Format:  types.MessageFormatUTF8,
			Content: "hello world",
		},
	}
	request := display.NewSetDisplayMessageRequest(message)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"message":{"id":%v,"priority":"%v","state":"%v","startDateTime":"%v","message":{"format":"%v","content":"%v"}}}]`,
		messageId, display.SetDisplayMessageFeatureName, message.ID, message.Priority, message.State, message.StartDateTime.FormatTimestamp(), message.Message.Format, message.Message.Content)
	testUnsupportedRequestFromChargingStation(suite, request, requestJson, messageId)
}
