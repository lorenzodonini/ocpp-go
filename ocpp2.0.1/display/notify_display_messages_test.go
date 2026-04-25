package display

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *displayTestSuite) TestNotifyDisplayMessagesRequestValidation() {
	t := suite.T()
	messageInfo := MessageInfo{ID: 42, Priority: MessagePriorityAlwaysFront, State: MessageStateIdle, Message: types.MessageContent{Format: types.MessageFormatUTF8, Content: "hello world"}}
	var requestTable = []tests.GenericTestEntry{
		{NotifyDisplayMessagesRequest{RequestID: 42, Tbc: false, MessageInfo: []MessageInfo{messageInfo}}, true},
		{NotifyDisplayMessagesRequest{RequestID: 42, Tbc: false, MessageInfo: []MessageInfo{}}, true},
		{NotifyDisplayMessagesRequest{RequestID: 42, Tbc: false}, true},
		{NotifyDisplayMessagesRequest{RequestID: 42}, true},
		{NotifyDisplayMessagesRequest{}, true},
		{NotifyDisplayMessagesRequest{RequestID: -1}, false},
		{NotifyDisplayMessagesRequest{RequestID: 42, MessageInfo: []MessageInfo{{ID: 42, Priority: "invalidPriority", State: MessageStateIdle, Message: types.MessageContent{Format: types.MessageFormatUTF8, Content: "hello world"}}}}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *displayTestSuite) TestNotifyDisplayMessagesConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{NotifyDisplayMessagesResponse{}, true},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *displayTestSuite) TestNotifyDisplayMessagesFeature() {
	feature := NotifyDisplayMessagesFeature{}
	suite.Equal(NotifyDisplayMessagesFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(NotifyDisplayMessagesRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(NotifyDisplayMessagesResponse{}), feature.GetResponseType())
}

func (suite *displayTestSuite) TestNewNotifyDisplayMessagesRequest() {
	req := NewNotifyDisplayMessagesRequest(5)
	suite.NotNil(req)
	suite.Equal(NotifyDisplayMessagesFeatureName, req.GetFeatureName())
	suite.Equal(5, req.RequestID)
}

func (suite *displayTestSuite) TestNewNotifyDisplayMessagesResponse() {
	resp := NewNotifyDisplayMessagesResponse()
	suite.NotNil(resp)
	suite.Equal(NotifyDisplayMessagesFeatureName, resp.GetFeatureName())
}
