package authorization

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

// Test
func (suite *authTestSuite) TestClearCacheRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{ClearCacheRequest{}, true},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *authTestSuite) TestClearCacheConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{ClearCacheResponse{Status: ClearCacheStatusAccepted, StatusInfo: types.NewStatusInfo("200", "ok")}, true},
		{ClearCacheResponse{Status: ClearCacheStatusAccepted}, true},
		{ClearCacheResponse{Status: ClearCacheStatusRejected}, true},
		{ClearCacheResponse{Status: "invalidClearCacheStatus"}, false},
		{ClearCacheResponse{}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *authTestSuite) TestClearCacheFeature() {
	feature := ClearCacheFeature{}
	suite.Equal(ClearCacheFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(ClearCacheRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(ClearCacheResponse{}), feature.GetResponseType())
}

func (suite *authTestSuite) TestNewClearCacheRequest() {
	req := NewClearCacheRequest()
	suite.NotNil(req)
	suite.Equal(ClearCacheFeatureName, req.GetFeatureName())
}

func (suite *authTestSuite) TestNewClearCacheResponse() {
	resp := NewClearCacheResponse(ClearCacheStatusAccepted)
	suite.NotNil(resp)
	suite.Equal(ClearCacheFeatureName, resp.GetFeatureName())
	suite.Equal(ClearCacheStatusAccepted, resp.Status)
}
