package diagnostics

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *diagnosticsTestSuite) TestGetMonitoringReportRequestValidation() {
	t := suite.T()
	componentVariables := []types.ComponentVariable{
		{
			Component: types.Component{Name: "component1", Instance: "instance1", EVSE: &types.EVSE{ID: 2, ConnectorID: tests.NewInt(2)}},
			Variable:  types.Variable{Name: "variable1", Instance: "instance1"},
		},
	}
	var requestTable = []tests.GenericTestEntry{
		{GetMonitoringReportRequest{RequestID: tests.NewInt(42), MonitoringCriteria: []MonitoringCriteriaType{MonitoringCriteriaThresholdMonitoring, MonitoringCriteriaDeltaMonitoring, MonitoringCriteriaPeriodicMonitoring}, ComponentVariable: componentVariables}, true},
		{GetMonitoringReportRequest{RequestID: tests.NewInt(42), MonitoringCriteria: []MonitoringCriteriaType{}, ComponentVariable: componentVariables}, true},
		{GetMonitoringReportRequest{RequestID: tests.NewInt(42), ComponentVariable: componentVariables}, true},
		{GetMonitoringReportRequest{RequestID: tests.NewInt(42), ComponentVariable: []types.ComponentVariable{}}, true},
		{GetMonitoringReportRequest{RequestID: tests.NewInt(42)}, true},
		{GetMonitoringReportRequest{}, true},
		{GetMonitoringReportRequest{RequestID: tests.NewInt(-1)}, false},
		{GetMonitoringReportRequest{MonitoringCriteria: []MonitoringCriteriaType{MonitoringCriteriaThresholdMonitoring, MonitoringCriteriaDeltaMonitoring, MonitoringCriteriaPeriodicMonitoring, MonitoringCriteriaThresholdMonitoring}}, false},
		{GetMonitoringReportRequest{MonitoringCriteria: []MonitoringCriteriaType{"invalidMonitoringCriteria"}}, false},
		{GetMonitoringReportRequest{ComponentVariable: []types.ComponentVariable{{Variable: types.Variable{Name: "variable1", Instance: "instance1"}}}}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *diagnosticsTestSuite) TestGetMonitoringReportConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{GetMonitoringReportResponse{Status: types.GenericDeviceModelStatusAccepted}, true},
		{GetMonitoringReportResponse{Status: "invalidDeviceModelStatus"}, false},
		{GetMonitoringReportResponse{}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *diagnosticsTestSuite) TestGetMonitoringReportFeature() {
	feature := GetMonitoringReportFeature{}
	suite.Equal(GetMonitoringReportFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(GetMonitoringReportRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(GetMonitoringReportResponse{}), feature.GetResponseType())
}

func (suite *diagnosticsTestSuite) TestNewGetMonitoringReportRequest() {
	req := NewGetMonitoringReportRequest()
	suite.NotNil(req)
	suite.Equal(GetMonitoringReportFeatureName, req.GetFeatureName())
}

func (suite *diagnosticsTestSuite) TestNewGetMonitoringReportResponse() {
	resp := NewGetMonitoringReportResponse(types.GenericDeviceModelStatusAccepted)
	suite.NotNil(resp)
	suite.Equal(GetMonitoringReportFeatureName, resp.GetFeatureName())
	suite.Equal(types.GenericDeviceModelStatusAccepted, resp.Status)
}
