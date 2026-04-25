package provisioning

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *provisioningTestSuite) TestResetRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{ResetRequest{Type: ResetTypeImmediate, EvseID: tests.NewInt(42)}, true},
		{ResetRequest{Type: ResetTypeOnIdle, EvseID: tests.NewInt(42)}, true},
		{ResetRequest{Type: ResetTypeImmediate}, true},
		{ResetRequest{}, false},
		{ResetRequest{Type: ResetTypeImmediate, EvseID: tests.NewInt(-1)}, false},
		{ResetRequest{Type: "invalidResetType", EvseID: tests.NewInt(42)}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *provisioningTestSuite) TestResetResponseValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{ResetResponse{Status: ResetStatusAccepted, StatusInfo: types.NewStatusInfo("200", "")}, true},
		{ResetResponse{Status: ResetStatusRejected, StatusInfo: types.NewStatusInfo("200", "")}, true},
		{ResetResponse{Status: ResetStatusScheduled, StatusInfo: types.NewStatusInfo("200", "")}, true},
		{ResetResponse{Status: ResetStatusAccepted}, true},
		{ResetResponse{}, false},
		{ResetResponse{Status: ResetStatusAccepted, StatusInfo: types.NewStatusInfo("", "")}, false},
		{ResetResponse{Status: "invalidResetStatus", StatusInfo: types.NewStatusInfo("200", "")}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *provisioningTestSuite) TestResetFeature() {
	feature := ResetFeature{}
	suite.Equal(ResetFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(ResetRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(ResetResponse{}), feature.GetResponseType())
}

func (suite *provisioningTestSuite) TestNewResetRequest() {
	req := NewResetRequest(ResetTypeImmediate)
	suite.NotNil(req)
	suite.Equal(ResetFeatureName, req.GetFeatureName())
	suite.Equal(ResetTypeImmediate, req.Type)
}

func (suite *provisioningTestSuite) TestNewResetResponse() {
	resp := NewResetResponse(ResetStatusAccepted)
	suite.NotNil(resp)
	suite.Equal(ResetFeatureName, resp.GetFeatureName())
	suite.Equal(ResetStatusAccepted, resp.Status)
}
