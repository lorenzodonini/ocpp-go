package display

import (
	"reflect"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *displayTestSuite) TestSetDisplayMessageRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{SetDisplayMessageRequest{Message: MessageInfo{ID: 42, Priority: MessagePriorityAlwaysFront, State: MessageStateIdle, StartDateTime: types.NewDateTime(time.Now()), Message: types.MessageContent{Format: types.MessageFormatUTF8, Content: "hello world"}}}, true},
		{SetDisplayMessageRequest{}, false},
		{SetDisplayMessageRequest{Message: MessageInfo{ID: 42, Priority: "invalidPriority", State: MessageStateIdle, StartDateTime: types.NewDateTime(time.Now()), Message: types.MessageContent{Format: types.MessageFormatUTF8, Content: "hello world"}}}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *displayTestSuite) TestSetDisplayMessageConfirmationValidation() {
	t := suite.T()
	var responseTable = []tests.GenericTestEntry{
		{SetDisplayMessageResponse{Status: DisplayMessageStatusAccepted, StatusInfo: types.NewStatusInfo("200", "")}, true},
		{SetDisplayMessageResponse{Status: DisplayMessageStatusAccepted}, true},
		{SetDisplayMessageResponse{Status: DisplayMessageStatusNotSupportedMessageFormat}, true},
		{SetDisplayMessageResponse{Status: DisplayMessageStatusNotSupportedState}, true},
		{SetDisplayMessageResponse{Status: DisplayMessageStatusNotSupportedPriority}, true},
		{SetDisplayMessageResponse{Status: DisplayMessageStatusRejected}, true},
		{SetDisplayMessageResponse{Status: DisplayMessageStatusUnknownTransaction}, true},
		{SetDisplayMessageResponse{Status: "invalidDisplayMessageStatus"}, false},
		{SetDisplayMessageResponse{}, false},
	}
	tests.ExecuteGenericTestTable(t, responseTable)
}

func (suite *displayTestSuite) TestSetDisplayMessageFeature() {
	feature := SetDisplayMessageFeature{}
	suite.Equal(SetDisplayMessageFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(SetDisplayMessageRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(SetDisplayMessageResponse{}), feature.GetResponseType())
}

func (suite *displayTestSuite) TestNewSetDisplayMessageRequest() {
	msg := MessageInfo{
		ID:       1,
		Priority: MessagePriorityAlwaysFront,
		Message:  types.MessageContent{Format: types.MessageFormatUTF8, Content: "hello"},
	}
	req := NewSetDisplayMessageRequest(msg)
	suite.NotNil(req)
	suite.Equal(SetDisplayMessageFeatureName, req.GetFeatureName())
	suite.Equal(msg, req.Message)
}

func (suite *displayTestSuite) TestNewSetDisplayMessageResponse() {
	resp := NewSetDisplayMessageResponse(DisplayMessageStatusAccepted)
	suite.NotNil(resp)
	suite.Equal(SetDisplayMessageFeatureName, resp.GetFeatureName())
	suite.Equal(DisplayMessageStatusAccepted, resp.Status)
}

