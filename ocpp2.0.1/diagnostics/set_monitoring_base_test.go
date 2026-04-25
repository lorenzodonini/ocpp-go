package diagnostics

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *diagnosticsTestSuite) TestSetMonitoringBaseRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{SetMonitoringBaseRequest{MonitoringBase: MonitoringBaseAll}, true},
		{SetMonitoringBaseRequest{MonitoringBase: MonitoringBaseFactoryDefault}, true},
		{SetMonitoringBaseRequest{MonitoringBase: MonitoringBaseHardWiredOnly}, true},
		{SetMonitoringBaseRequest{MonitoringBase: "invalidMonitoringBase"}, false},
		{SetMonitoringBaseRequest{}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *diagnosticsTestSuite) TestSetMonitoringBaseConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{SetMonitoringBaseResponse{Status: types.GenericDeviceModelStatusAccepted, StatusInfo: types.NewStatusInfo("200", "")}, true},
		{SetMonitoringBaseResponse{Status: types.GenericDeviceModelStatusAccepted}, true},
		{SetMonitoringBaseResponse{Status: "invalidDeviceModelStatus"}, false},
		{SetMonitoringBaseResponse{}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *diagnosticsTestSuite) TestSetMonitoringBaseFeature() {
	feature := SetMonitoringBaseFeature{}
	suite.Equal(SetMonitoringBaseFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(SetMonitoringBaseRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(SetMonitoringBaseResponse{}), feature.GetResponseType())
}

func (suite *diagnosticsTestSuite) TestNewSetMonitoringBaseRequest() {
	req := NewSetMonitoringBaseRequest(MonitoringBaseAll)
	suite.NotNil(req)
	suite.Equal(SetMonitoringBaseFeatureName, req.GetFeatureName())
	suite.Equal(MonitoringBaseAll, req.MonitoringBase)
}

func (suite *diagnosticsTestSuite) TestNewSetMonitoringBaseResponse() {
	resp := NewSetMonitoringBaseResponse(types.GenericDeviceModelStatusAccepted)
	suite.NotNil(resp)
	suite.Equal(SetMonitoringBaseFeatureName, resp.GetFeatureName())
	suite.Equal(types.GenericDeviceModelStatusAccepted, resp.Status)
}
