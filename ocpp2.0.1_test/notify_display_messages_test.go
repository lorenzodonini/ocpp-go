package ocpp2_test

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/display"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

// Test
func (suite *OcppV2TestSuite) TestNotifyDisplayMessagesRequestValidation() {
	messageInfo := display.MessageInfo{ID: 42, Priority: display.MessagePriorityAlwaysFront, State: display.MessageStateIdle, Message: types.MessageContent{Format: types.MessageFormatUTF8, Content: "hello world"}}
	var requestTable = []GenericTestEntry{
		{display.NotifyDisplayMessagesRequest{RequestID: 42, Tbc: false, MessageInfo: []display.MessageInfo{messageInfo}}, true},
		{display.NotifyDisplayMessagesRequest{RequestID: 42, Tbc: false, MessageInfo: []display.MessageInfo{}}, true},
		{display.NotifyDisplayMessagesRequest{RequestID: 42, Tbc: false}, true},
		{display.NotifyDisplayMessagesRequest{RequestID: 42}, true},
		{display.NotifyDisplayMessagesRequest{}, true},
		{display.NotifyDisplayMessagesRequest{RequestID: -1}, false},
		{display.NotifyDisplayMessagesRequest{RequestID: 42, MessageInfo: []display.MessageInfo{{ID: 42, Priority: "invalidPriority", State: display.MessageStateIdle, Message: types.MessageContent{Format: types.MessageFormatUTF8, Content: "hello world"}}}}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV2TestSuite) TestNotifyDisplayMessagesConfirmationValidation() {
	var confirmationTable = []GenericTestEntry{
		{display.NotifyDisplayMessagesResponse{}, true},
	}
	ExecuteGenericTestTable(suite, confirmationTable)
}

func (suite *OcppV2TestSuite) TestNotifyDisplayMessagesE2EMocked() {
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	requestID := 42
	tbc := false
	messageInfo := display.MessageInfo{ID: 42, Priority: display.MessagePriorityAlwaysFront, State: display.MessageStateIdle, StartDateTime: types.NewDateTime(time.Now()), Message: types.MessageContent{Format: types.MessageFormatUTF8, Content: "hello world"}}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"requestId":%v,"messageInfo":[{"id":%v,"priority":"%v","state":"%v","startDateTime":"%v","message":{"format":"%v","content":"%v"}}]}]`,
		messageId, display.NotifyDisplayMessagesFeatureName, requestID, messageInfo.ID, messageInfo.Priority, messageInfo.State, messageInfo.StartDateTime.FormatTimestamp(), messageInfo.Message.Format, messageInfo.Message.Content)
	responseJson := fmt.Sprintf(`[3,"%v",{}]`, messageId)
	response := display.NewNotifyDisplayMessagesResponse()
	channel := NewMockWebSocket(wsId)

	handler := &MockCSMSDisplayHandler{}
	handler.On("OnNotifyDisplayMessages", mock.AnythingOfType("string"), mock.Anything).Return(response, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*display.NotifyDisplayMessagesRequest)
		suite.Require().True(ok)
		suite.Equal(requestID, request.RequestID)
		suite.Equal(tbc, request.Tbc)
		suite.Require().Len(request.MessageInfo, 1)
		suite.Equal(messageInfo.ID, request.MessageInfo[0].ID)
		suite.Equal(messageInfo.Priority, request.MessageInfo[0].Priority)
		suite.Equal(messageInfo.State, request.MessageInfo[0].State)
		assertDateTimeEquality(suite, messageInfo.StartDateTime, request.MessageInfo[0].StartDateTime)
		suite.Equal(messageInfo.Message.Format, request.MessageInfo[0].Message.Format)
		suite.Equal(messageInfo.Message.Content, request.MessageInfo[0].Message.Content)
		suite.Equal(messageInfo.Message.Language, request.MessageInfo[0].Message.Language)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	suite.Require().Nil(err)
	r, err := suite.chargingStation.NotifyDisplayMessages(requestID, func(request *display.NotifyDisplayMessagesRequest) {
		request.MessageInfo = []display.MessageInfo{messageInfo}
	})
	suite.Nil(err)
	suite.NotNil(r)
}

func (suite *OcppV2TestSuite) TestNotifyDisplayMessagesInvalidEndpoint() {
	messageId := defaultMessageId
	requestID := 42
	messageInfo := display.MessageInfo{ID: 42, Priority: display.MessagePriorityAlwaysFront, State: display.MessageStateIdle, StartDateTime: types.NewDateTime(time.Now()), Message: types.MessageContent{Format: types.MessageFormatUTF8, Content: "hello world"}}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"requestId":%v,"messageInfo":[{"id":%v,"priority":"%v","state":"%v","startDateTime":"%v","message":{"format":"%v","content":"%v"}}]}]`,
		messageId, display.NotifyDisplayMessagesFeatureName, requestID, messageInfo.ID, messageInfo.Priority, messageInfo.State, messageInfo.StartDateTime.FormatTimestamp(), messageInfo.Message.Format, messageInfo.Message.Content)
	req := display.NewNotifyDisplayMessagesRequest(requestID)
	testUnsupportedRequestFromCentralSystem(suite, req, requestJson, messageId)
}
