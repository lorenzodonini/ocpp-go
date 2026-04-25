package diagnostics

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *diagnosticsTestSuite) TestSetVariableMonitoringRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{SetVariableMonitoringRequest{MonitoringData: []SetMonitoringData{{ID: tests.NewInt(2), Transaction: true, Value: 42.0, Type: MonitorUpperThreshold, Severity: 5, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}}}, true},
		{SetVariableMonitoringRequest{MonitoringData: []SetMonitoringData{{Transaction: true, Value: 42.0, Type: MonitorUpperThreshold, Severity: 5, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}}}, true},
		{SetVariableMonitoringRequest{MonitoringData: []SetMonitoringData{{Value: 42.0, Type: MonitorUpperThreshold, Severity: 5, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}}}, true},
		{SetVariableMonitoringRequest{MonitoringData: []SetMonitoringData{{Type: MonitorUpperThreshold, Severity: 5, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}}}, true},
		{SetVariableMonitoringRequest{MonitoringData: []SetMonitoringData{{Type: MonitorUpperThreshold, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}}}, true},
		{SetVariableMonitoringRequest{MonitoringData: []SetMonitoringData{{Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}}}, false},
		{SetVariableMonitoringRequest{MonitoringData: []SetMonitoringData{{Type: MonitorUpperThreshold, Variable: types.Variable{Name: "variable1"}}}}, false},
		{SetVariableMonitoringRequest{MonitoringData: []SetMonitoringData{{Type: MonitorUpperThreshold, Component: types.Component{Name: "component1"}}}}, false},
		{SetVariableMonitoringRequest{MonitoringData: []SetMonitoringData{{}}}, false},
		{SetVariableMonitoringRequest{MonitoringData: []SetMonitoringData{}}, false},
		{SetVariableMonitoringRequest{}, false},
		{SetVariableMonitoringRequest{MonitoringData: []SetMonitoringData{{ID: tests.NewInt(2), Transaction: true, Value: 42.0, Type: "invalidType", Severity: 5, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}}}, false},
		{SetVariableMonitoringRequest{MonitoringData: []SetMonitoringData{{ID: tests.NewInt(2), Transaction: true, Value: 42.0, Type: MonitorUpperThreshold, Severity: -1, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}}}, false},
		{SetVariableMonitoringRequest{MonitoringData: []SetMonitoringData{{ID: tests.NewInt(2), Transaction: true, Value: 42.0, Type: MonitorUpperThreshold, Severity: 10, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}}}, false},
		{SetVariableMonitoringRequest{MonitoringData: []SetMonitoringData{{ID: tests.NewInt(2), Transaction: true, Value: 42.0, Type: MonitorUpperThreshold, Severity: 5, Component: types.Component{}, Variable: types.Variable{}}}}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *diagnosticsTestSuite) TestSetVariableMonitoringResponseValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{SetVariableMonitoringResponse{MonitoringResult: []SetMonitoringResult{{ID: tests.NewInt(2), Status: SetMonitoringStatusAccepted, Type: MonitorUpperThreshold, Severity: 5, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}, StatusInfo: types.NewStatusInfo("200", "")}}}, true},
		{SetVariableMonitoringResponse{MonitoringResult: []SetMonitoringResult{{ID: tests.NewInt(2), Status: SetMonitoringStatusAccepted, Type: MonitorUpperThreshold, Severity: 5, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}}}, true},
		{SetVariableMonitoringResponse{MonitoringResult: []SetMonitoringResult{{Status: SetMonitoringStatusAccepted, Type: MonitorUpperThreshold, Severity: 5, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}}}, true},
		{SetVariableMonitoringResponse{MonitoringResult: []SetMonitoringResult{{Status: SetMonitoringStatusAccepted, Type: MonitorUpperThreshold, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}}}, true},
		{SetVariableMonitoringResponse{MonitoringResult: []SetMonitoringResult{{Status: SetMonitoringStatusAccepted, Type: MonitorUpperThreshold, Component: types.Component{Name: "component1"}}}}, false},
		{SetVariableMonitoringResponse{MonitoringResult: []SetMonitoringResult{{Status: SetMonitoringStatusAccepted, Type: MonitorUpperThreshold, Variable: types.Variable{Name: "variable1"}}}}, false},
		{SetVariableMonitoringResponse{MonitoringResult: []SetMonitoringResult{{Status: SetMonitoringStatusAccepted, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}}}, false},
		{SetVariableMonitoringResponse{MonitoringResult: []SetMonitoringResult{{Type: MonitorUpperThreshold, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}}}, false},
		{SetVariableMonitoringResponse{}, false},
		{SetVariableMonitoringResponse{MonitoringResult: []SetMonitoringResult{{ID: tests.NewInt(2), Status: "invalidStatus", Type: MonitorUpperThreshold, Severity: 5, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}, StatusInfo: types.NewStatusInfo("200", "")}}}, false},
		{SetVariableMonitoringResponse{MonitoringResult: []SetMonitoringResult{{ID: tests.NewInt(2), Status: SetMonitoringStatusAccepted, Type: "invalidType", Severity: 5, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}, StatusInfo: types.NewStatusInfo("200", "")}}}, false},
		{SetVariableMonitoringResponse{MonitoringResult: []SetMonitoringResult{{ID: tests.NewInt(2), Status: SetMonitoringStatusAccepted, Type: MonitorUpperThreshold, Severity: -1, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}, StatusInfo: types.NewStatusInfo("200", "")}}}, false},
		{SetVariableMonitoringResponse{MonitoringResult: []SetMonitoringResult{{ID: tests.NewInt(2), Status: SetMonitoringStatusAccepted, Type: MonitorUpperThreshold, Severity: 10, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}, StatusInfo: types.NewStatusInfo("200", "")}}}, false},
		{SetVariableMonitoringResponse{MonitoringResult: []SetMonitoringResult{{ID: tests.NewInt(2), Status: SetMonitoringStatusAccepted, Type: MonitorUpperThreshold, Severity: 5, Component: types.Component{Name: ""}, Variable: types.Variable{Name: "variable1"}, StatusInfo: types.NewStatusInfo("200", "")}}}, false},
		{SetVariableMonitoringResponse{MonitoringResult: []SetMonitoringResult{{ID: tests.NewInt(2), Status: SetMonitoringStatusAccepted, Type: MonitorUpperThreshold, Severity: 5, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: ""}, StatusInfo: types.NewStatusInfo("200", "")}}}, false},
		{SetVariableMonitoringResponse{MonitoringResult: []SetMonitoringResult{{ID: tests.NewInt(2), Status: SetMonitoringStatusAccepted, Type: MonitorUpperThreshold, Severity: 5, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}, StatusInfo: types.NewStatusInfo("", "")}}}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *diagnosticsTestSuite) TestSetVariableMonitoringFeature() {
	feature := SetVariableMonitoringFeature{}
	suite.Equal(SetVariableMonitoringFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(SetVariableMonitoringRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(SetVariableMonitoringResponse{}), feature.GetResponseType())
}

func (suite *diagnosticsTestSuite) TestNewSetVariableMonitoringRequest() {
	data := []SetMonitoringData{
		{Type: MonitorUpperThreshold, Component: types.Component{Name: "comp1"}, Variable: types.Variable{Name: "var1"}},
	}
	req := NewSetVariableMonitoringRequest(data)
	suite.NotNil(req)
	suite.Equal(SetVariableMonitoringFeatureName, req.GetFeatureName())
	suite.Equal(data, req.MonitoringData)
}

func (suite *diagnosticsTestSuite) TestNewSetVariableMonitoringResponse() {
	result := []SetMonitoringResult{
		{Status: SetMonitoringStatusAccepted, Type: MonitorUpperThreshold, Component: types.Component{Name: "comp1"}, Variable: types.Variable{Name: "var1"}},
	}
	resp := NewSetVariableMonitoringResponse(result)
	suite.NotNil(resp)
	suite.Equal(SetVariableMonitoringFeatureName, resp.GetFeatureName())
	suite.Equal(result, resp.MonitoringResult)
}
