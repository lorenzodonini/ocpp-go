package provisioning

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *provisioningTestSuite) TestGetVariablesRequestValidation() {
	t := suite.T()
	component := types.Component{Name: "component1", Instance: "instance1", EVSE: &types.EVSE{ID: 2, ConnectorID: tests.NewInt(2)}}
	variable := types.Variable{Name: "variable1", Instance: "instance1"}

	var requestTable = []tests.GenericTestEntry{
		{GetVariablesRequest{GetVariableData: []GetVariableData{{AttributeType: types.AttributeTarget, Component: component, Variable: variable}}}, true},
		{GetVariablesRequest{GetVariableData: []GetVariableData{{Component: component, Variable: variable}}}, true},
		{GetVariablesRequest{GetVariableData: []GetVariableData{{Component: types.Component{Name: "component1"}, Variable: variable}}}, true},
		{GetVariablesRequest{GetVariableData: []GetVariableData{{Component: component, Variable: types.Variable{Name: "variable1"}}}}, true},
		{GetVariablesRequest{GetVariableData: []GetVariableData{}}, false},
		{GetVariablesRequest{}, false},
		{GetVariablesRequest{GetVariableData: []GetVariableData{{AttributeType: "invalidAttribute", Component: component, Variable: variable}}}, false},
		{GetVariablesRequest{GetVariableData: []GetVariableData{{AttributeType: types.AttributeTarget, Variable: variable}}}, false},
		{GetVariablesRequest{GetVariableData: []GetVariableData{{AttributeType: types.AttributeTarget, Component: component}}}, false},
		{GetVariablesRequest{GetVariableData: []GetVariableData{{AttributeType: types.AttributeTarget, Component: types.Component{Name: "component1", EVSE: &types.EVSE{ID: -1}}, Variable: variable}}}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *provisioningTestSuite) TestGetVariablesConfirmationValidation() {
	t := suite.T()
	component := types.Component{Name: "component1", Instance: "instance1", EVSE: &types.EVSE{ID: 2, ConnectorID: tests.NewInt(2)}}
	variable := types.Variable{Name: "variable1", Instance: "instance1"}
	var confirmationTable = []tests.GenericTestEntry{
		{GetVariablesResponse{GetVariableResult: []GetVariableResult{{AttributeStatus: GetVariableStatusAccepted, AttributeType: types.AttributeTarget, AttributeValue: "dummyValue", Component: component, Variable: variable}}}, true},
		{GetVariablesResponse{GetVariableResult: []GetVariableResult{{AttributeStatus: GetVariableStatusAccepted, AttributeType: types.AttributeTarget, Component: component, Variable: variable}}}, true},
		{GetVariablesResponse{GetVariableResult: []GetVariableResult{{AttributeStatus: GetVariableStatusAccepted, Component: component, Variable: variable}}}, true},
		{GetVariablesResponse{GetVariableResult: []GetVariableResult{{Component: component, Variable: variable}}}, false},
		{GetVariablesResponse{GetVariableResult: []GetVariableResult{{AttributeStatus: GetVariableStatusAccepted, Variable: variable}}}, false},
		{GetVariablesResponse{GetVariableResult: []GetVariableResult{{AttributeStatus: GetVariableStatusAccepted, Component: component}}}, false},
		{GetVariablesResponse{GetVariableResult: []GetVariableResult{}}, false},
		{GetVariablesResponse{}, false},
		{GetVariablesResponse{GetVariableResult: []GetVariableResult{{AttributeStatus: GetVariableStatusAccepted, AttributeType: "invalidAttribute", AttributeValue: "dummyValue", Component: component, Variable: variable}}}, false},
		{GetVariablesResponse{GetVariableResult: []GetVariableResult{{AttributeStatus: "invalidStatus", AttributeType: types.AttributeTarget, AttributeValue: "dummyValue", Component: component, Variable: variable}}}, false},
		{GetVariablesResponse{GetVariableResult: []GetVariableResult{{AttributeStatus: GetVariableStatusAccepted, AttributeType: types.AttributeTarget, AttributeValue: ">1000....................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................", Component: component, Variable: variable}}}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *provisioningTestSuite) TestGetVariablesFeature() {
	feature := GetVariablesFeature{}
	suite.Equal(GetVariablesFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(GetVariablesRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(GetVariablesResponse{}), feature.GetResponseType())
}

func (suite *provisioningTestSuite) TestNewGetVariablesRequest() {
	data := []GetVariableData{
		{Component: types.Component{Name: "comp1"}, Variable: types.Variable{Name: "var1"}},
	}
	req := NewGetVariablesRequest(data)
	suite.NotNil(req)
	suite.Equal(GetVariablesFeatureName, req.GetFeatureName())
	suite.Equal(data, req.GetVariableData)
}

func (suite *provisioningTestSuite) TestNewGetVariablesResponse() {
	result := []GetVariableResult{
		{AttributeStatus: GetVariableStatusAccepted, Component: types.Component{Name: "comp1"}, Variable: types.Variable{Name: "var1"}},
	}
	resp := NewGetVariablesResponse(result)
	suite.NotNil(resp)
	suite.Equal(GetVariablesFeatureName, resp.GetFeatureName())
	suite.Equal(result, resp.GetVariableResult)
}
