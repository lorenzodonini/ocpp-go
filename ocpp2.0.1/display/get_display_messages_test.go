package display

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *displayTestSuite) TestGetDisplayMessagesRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{GetDisplayMessagesRequest{RequestID: 1, Priority: MessagePriorityAlwaysFront, State: MessageStateCharging, ID: []int{2, 3}}, true},
		{GetDisplayMessagesRequest{RequestID: 1, Priority: MessagePriorityAlwaysFront, State: MessageStateCharging, ID: []int{}}, true},
		{GetDisplayMessagesRequest{RequestID: 1, Priority: MessagePriorityAlwaysFront, State: MessageStateCharging}, true},
		{GetDisplayMessagesRequest{RequestID: 1, Priority: MessagePriorityAlwaysFront}, true},
		{GetDisplayMessagesRequest{RequestID: 1, State: MessageStateCharging}, true},
		{GetDisplayMessagesRequest{RequestID: 1}, true},
		{GetDisplayMessagesRequest{}, true},
		{GetDisplayMessagesRequest{RequestID: -1}, false},
		{GetDisplayMessagesRequest{RequestID: 1, Priority: "invalidMessagePriority", State: MessageStateCharging, ID: []int{2, 3}}, false},
		{GetDisplayMessagesRequest{RequestID: 1, Priority: MessagePriorityAlwaysFront, State: "invalidMessageState", ID: []int{2, 3}}, false},
		{GetDisplayMessagesRequest{RequestID: 1, Priority: MessagePriorityAlwaysFront, State: MessageStateCharging, ID: []int{-2, 3}}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *displayTestSuite) TestGetDisplayMessagesConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{GetDisplayMessagesResponse{Status: MessageStatusAccepted}, true},
		{GetDisplayMessagesResponse{Status: MessageStatusUnknown}, true},
		{GetDisplayMessagesResponse{Status: "invalidMessageStatus"}, false},
		{GetDisplayMessagesResponse{}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *displayTestSuite) TestGetDisplayMessagesFeature() {
	feature := GetDisplayMessagesFeature{}
	suite.Equal(GetDisplayMessagesFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(GetDisplayMessagesRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(GetDisplayMessagesResponse{}), feature.GetResponseType())
}

func (suite *displayTestSuite) TestNewGetDisplayMessagesRequest() {
	req := NewGetDisplayMessagesRequest(7)
	suite.NotNil(req)
	suite.Equal(GetDisplayMessagesFeatureName, req.GetFeatureName())
	suite.Equal(7, req.RequestID)
}

func (suite *displayTestSuite) TestNewGetDisplayMessagesResponse() {
	resp := NewGetDisplayMessagesResponse(MessageStatusAccepted)
	suite.NotNil(resp)
	suite.Equal(GetDisplayMessagesFeatureName, resp.GetFeatureName())
	suite.Equal(MessageStatusAccepted, resp.Status)
}
