package remotecontrol

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *remoteControlTestSuite) TestTriggerMessageRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{TriggerMessageRequest{RequestedMessage: MessageTriggerStatusNotification, Evse: &types.EVSE{ID: 1}}, true},
		{TriggerMessageRequest{RequestedMessage: MessageTriggerStatusNotification}, true},
		{TriggerMessageRequest{}, false},
		{TriggerMessageRequest{RequestedMessage: "invalidMessageTrigger", Evse: &types.EVSE{ID: 1}}, false},
		{TriggerMessageRequest{RequestedMessage: MessageTriggerStatusNotification, Evse: &types.EVSE{ID: -1}}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *remoteControlTestSuite) TestTriggerMessageResponseValidation() {
	t := suite.T()
	var responseTable = []tests.GenericTestEntry{
		{TriggerMessageResponse{Status: TriggerMessageStatusAccepted, StatusInfo: &types.StatusInfo{ReasonCode: "200"}}, true},
		{TriggerMessageResponse{Status: TriggerMessageStatusAccepted}, true},
		{TriggerMessageResponse{}, false},
		{TriggerMessageResponse{Status: "invalidTriggerMessageStatus", StatusInfo: &types.StatusInfo{ReasonCode: "200"}}, false},
		{TriggerMessageResponse{Status: TriggerMessageStatusAccepted, StatusInfo: &types.StatusInfo{}}, false},
	}
	tests.ExecuteGenericTestTable(t, responseTable)
}

func (suite *remoteControlTestSuite) TestTriggerMessageFeature() {
	feature := TriggerMessageFeature{}
	suite.Equal(TriggerMessageFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(TriggerMessageRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(TriggerMessageResponse{}), feature.GetResponseType())
}

func (suite *remoteControlTestSuite) TestNewTriggerMessageRequest() {
	req := NewTriggerMessageRequest(MessageTriggerStatusNotification)
	suite.NotNil(req)
	suite.Equal(TriggerMessageFeatureName, req.GetFeatureName())
	suite.Equal(MessageTriggerStatusNotification, req.RequestedMessage)
}

func (suite *remoteControlTestSuite) TestNewTriggerMessageResponse() {
	resp := NewTriggerMessageResponse(TriggerMessageStatusAccepted)
	suite.NotNil(resp)
	suite.Equal(TriggerMessageFeatureName, resp.GetFeatureName())
	suite.Equal(TriggerMessageStatusAccepted, resp.Status)
}
