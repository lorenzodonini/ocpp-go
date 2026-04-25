package provisioning

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *provisioningTestSuite) TestGetReportRequestValidation() {
	t := suite.T()
	componentVariables := []types.ComponentVariable{
		{
			Component: types.Component{Name: "component1", Instance: "instance1", EVSE: &types.EVSE{ID: 2, ConnectorID: tests.NewInt(2)}},
			Variable:  types.Variable{Name: "variable1", Instance: "instance1"},
		},
	}
	var requestTable = []tests.GenericTestEntry{
		{GetReportRequest{RequestID: tests.NewInt(42), ComponentCriteria: []ComponentCriterion{ComponentCriterionActive, ComponentCriterionEnabled, ComponentCriterionAvailable, ComponentCriterionProblem}, ComponentVariable: componentVariables}, true},
		{GetReportRequest{RequestID: tests.NewInt(42), ComponentCriteria: []ComponentCriterion{ComponentCriterionActive, ComponentCriterionEnabled, ComponentCriterionAvailable, ComponentCriterionProblem}, ComponentVariable: []types.ComponentVariable{}}, true},
		{GetReportRequest{RequestID: tests.NewInt(42), ComponentCriteria: []ComponentCriterion{ComponentCriterionActive, ComponentCriterionEnabled, ComponentCriterionAvailable, ComponentCriterionProblem}}, true},
		{GetReportRequest{RequestID: tests.NewInt(42), ComponentCriteria: []ComponentCriterion{}}, true},
		{GetReportRequest{RequestID: tests.NewInt(42)}, true},
		{GetReportRequest{}, true},
		{GetReportRequest{RequestID: tests.NewInt(-1)}, false},
		{GetReportRequest{ComponentCriteria: []ComponentCriterion{"invalidComponentCriterion"}}, false},
		{GetReportRequest{ComponentCriteria: []ComponentCriterion{ComponentCriterionActive, ComponentCriterionEnabled, ComponentCriterionAvailable, ComponentCriterionProblem, ComponentCriterionActive}}, false},
		{GetReportRequest{ComponentVariable: []types.ComponentVariable{{Variable: types.Variable{Name: "variable1", Instance: "instance1"}}}}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *provisioningTestSuite) TestGetReportConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{GetReportResponse{Status: types.GenericDeviceModelStatusAccepted}, true},
		{GetReportResponse{Status: "invalidDeviceModelStatus"}, false},
		{GetReportResponse{}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *provisioningTestSuite) TestGetReportFeature() {
	feature := GetReportFeature{}
	suite.Equal(GetReportFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(GetReportRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(GetReportResponse{}), feature.GetResponseType())
}

func (suite *provisioningTestSuite) TestNewGetReportRequest() {
	req := NewGetReportRequest()
	suite.NotNil(req)
	suite.Equal(GetReportFeatureName, req.GetFeatureName())
}

func (suite *provisioningTestSuite) TestNewGetReportResponse() {
	resp := NewGetReportResponse(types.GenericDeviceModelStatusAccepted)
	suite.NotNil(resp)
	suite.Equal(GetReportFeatureName, resp.GetFeatureName())
	suite.Equal(types.GenericDeviceModelStatusAccepted, resp.Status)
}
