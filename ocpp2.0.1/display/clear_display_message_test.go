package display

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *displayTestSuite) TestClearDisplayMessageRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{ClearDisplayRequest{ID: 42}, true},
		{ClearDisplayRequest{}, true},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *displayTestSuite) TestClearDisplayMessageResponseValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{ClearDisplayResponse{Status: ClearMessageStatusAccepted}, true},
		{ClearDisplayResponse{Status: ClearMessageStatusUnknown}, true},
		{ClearDisplayResponse{Status: "invalidClearMessageStatus"}, false},
		{ClearDisplayResponse{}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *displayTestSuite) TestClearDisplayFeature() {
	feature := ClearDisplayFeature{}
	suite.Equal(ClearDisplayMessageFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(ClearDisplayRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(ClearDisplayResponse{}), feature.GetResponseType())
}

func (suite *displayTestSuite) TestNewClearDisplayRequest() {
	req := NewClearDisplayRequest(42)
	suite.NotNil(req)
	suite.Equal(ClearDisplayMessageFeatureName, req.GetFeatureName())
	suite.Equal(42, req.ID)
}

func (suite *displayTestSuite) TestNewClearDisplayResponse() {
	resp := NewClearDisplayResponse(ClearMessageStatusAccepted)
	suite.NotNil(resp)
	suite.Equal(ClearDisplayMessageFeatureName, resp.GetFeatureName())
	suite.Equal(ClearMessageStatusAccepted, resp.Status)
}
