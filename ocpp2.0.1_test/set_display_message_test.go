package ocpp2_test

import (
	"fmt"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/display"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"

	"github.com/stretchr/testify/mock"
)

// Test
func (suite *OcppV2TestSuite) TestSetDisplayMessageRequestValidation() {
	var requestTable = []GenericTestEntry{
		{display.SetDisplayMessageRequest{Message: display.MessageInfo{ID: 42, Priority: display.MessagePriorityAlwaysFront, State: display.MessageStateIdle, StartDateTime: types.NewDateTime(time.Now()), Message: types.MessageContent{Format: types.MessageFormatUTF8, Content: "hello world"}}}, true},
		{display.SetDisplayMessageRequest{}, false},
		{display.SetDisplayMessageRequest{Message: display.MessageInfo{ID: 42, Priority: "invalidPriority", State: display.MessageStateIdle, StartDateTime: types.NewDateTime(time.Now()), Message: types.MessageContent{Format: types.MessageFormatUTF8, Content: "hello world"}}}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV2TestSuite) TestSetDisplayMessageConfirmationValidation() {
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
	ExecuteGenericTestTable(suite, responseTable)
}

func (suite *OcppV2TestSuite) TestSetDisplayMessageE2EMocked() {
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

	handler := &MockChargingStationDisplayHandler{}
	handler.On("OnSetDisplayMessage", mock.Anything).Return(setDisplayResponse, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*display.SetDisplayMessageRequest)
		suite.Require().True(ok)
		suite.Require().NotNil(request)
		suite.Equal(message.ID, request.Message.ID)
		suite.Equal(message.Priority, request.Message.Priority)
		suite.Equal(message.State, request.Message.State)
		assertDateTimeEquality(suite, message.StartDateTime, request.Message.StartDateTime)
		suite.Equal(message.Message.Format, request.Message.Message.Format)
		suite.Equal(message.Message.Content, request.Message.Message.Content)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	suite.Require().Nil(err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.SetDisplayMessage(wsId, func(response *display.SetDisplayMessageResponse, err error) {
		suite.Require().Nil(err)
		suite.Require().NotNil(response)
		suite.Equal(status, response.Status)
		suite.Equal(statusInfo.ReasonCode, response.StatusInfo.ReasonCode)
		suite.Equal(statusInfo.AdditionalInfo, response.StatusInfo.AdditionalInfo)
		resultChannel <- true
	}, message)
	suite.Require().Nil(err)
	result := <-resultChannel
	suite.True(result)
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
