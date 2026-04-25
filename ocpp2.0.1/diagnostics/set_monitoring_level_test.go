package diagnostics

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *diagnosticsTestSuite) TestSetMonitoringLevelRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{SetMonitoringLevelRequest{Severity: 0}, true},
		{SetMonitoringLevelRequest{Severity: 1}, true},
		{SetMonitoringLevelRequest{Severity: 2}, true},
		{SetMonitoringLevelRequest{Severity: 3}, true},
		{SetMonitoringLevelRequest{Severity: 4}, true},
		{SetMonitoringLevelRequest{Severity: 5}, true},
		{SetMonitoringLevelRequest{Severity: 6}, true},
		{SetMonitoringLevelRequest{Severity: 7}, true},
		{SetMonitoringLevelRequest{Severity: 8}, true},
		{SetMonitoringLevelRequest{Severity: 9}, true},
		{SetMonitoringLevelRequest{}, true},
		{SetMonitoringLevelRequest{Severity: -1}, false},
		{SetMonitoringLevelRequest{Severity: 10}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *diagnosticsTestSuite) TestSetMonitoringLevelConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{SetMonitoringLevelResponse{Status: types.GenericDeviceModelStatusAccepted, StatusInfo: types.NewStatusInfo("200", "")}, true},
		{SetMonitoringLevelResponse{Status: types.GenericDeviceModelStatusAccepted}, true},
		{SetMonitoringLevelResponse{Status: "invalidDeviceModelStatus"}, false},
		{SetMonitoringLevelResponse{}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *diagnosticsTestSuite) TestSetMonitoringLevelFeature() {
	feature := SetMonitoringLevelFeature{}
	suite.Equal(SetMonitoringLevelFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(SetMonitoringLevelRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(SetMonitoringLevelResponse{}), feature.GetResponseType())
}

func (suite *diagnosticsTestSuite) TestNewSetMonitoringLevelRequest() {
	req := NewSetMonitoringLevelRequest(5)
	suite.NotNil(req)
	suite.Equal(SetMonitoringLevelFeatureName, req.GetFeatureName())
	suite.Equal(5, req.Severity)
}

func (suite *diagnosticsTestSuite) TestNewSetMonitoringLevelResponse() {
	resp := NewSetMonitoringLevelResponse(types.GenericDeviceModelStatusAccepted)
	suite.NotNil(resp)
	suite.Equal(SetMonitoringLevelFeatureName, resp.GetFeatureName())
	suite.Equal(types.GenericDeviceModelStatusAccepted, resp.Status)
}
