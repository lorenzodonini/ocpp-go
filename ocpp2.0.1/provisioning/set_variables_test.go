package provisioning

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *provisioningTestSuite) TestSetVariablesRequestValidation() {
	t := suite.T()
	component := types.Component{Name: "component1", Instance: "instance1", EVSE: &types.EVSE{ID: 2, ConnectorID: tests.NewInt(2)}}
	variable := types.Variable{Name: "variable1", Instance: "instance1"}

	var requestTable = []tests.GenericTestEntry{
		{SetVariablesRequest{SetVariableData: []SetVariableData{{AttributeType: types.AttributeTarget, AttributeValue: "dummyValue", Component: component, Variable: variable}}}, true},
		{SetVariablesRequest{SetVariableData: []SetVariableData{{AttributeValue: "dummyValue", Component: component, Variable: variable}}}, true},
		{SetVariablesRequest{SetVariableData: []SetVariableData{{AttributeValue: "dummyValue", Component: types.Component{Name: "component1"}, Variable: variable}}}, true},
		{SetVariablesRequest{SetVariableData: []SetVariableData{{AttributeValue: "dummyValue", Component: component, Variable: types.Variable{Name: "variable1"}}}}, true},
		{SetVariablesRequest{SetVariableData: []SetVariableData{{Component: component, Variable: variable}}}, false},
		{SetVariablesRequest{SetVariableData: []SetVariableData{{AttributeValue: "dummyValue", Variable: variable}}}, false},
		{SetVariablesRequest{SetVariableData: []SetVariableData{{AttributeValue: "dummyValue", Component: component}}}, false},
		{SetVariablesRequest{SetVariableData: []SetVariableData{}}, false},
		{SetVariablesRequest{}, false},
		{SetVariablesRequest{SetVariableData: []SetVariableData{{AttributeType: "invalidAttribute", AttributeValue: "dummyValue", Component: component, Variable: variable}}}, false},
		{SetVariablesRequest{SetVariableData: []SetVariableData{{AttributeType: types.AttributeTarget, AttributeValue: ">1000....................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................", Component: component, Variable: variable}}}, false},
		{SetVariablesRequest{SetVariableData: []SetVariableData{{AttributeType: types.AttributeTarget, AttributeValue: "dummyValue", Variable: variable}}}, false},
		{SetVariablesRequest{SetVariableData: []SetVariableData{{AttributeType: types.AttributeTarget, AttributeValue: "dummyValue", Component: component}}}, false},
		{SetVariablesRequest{SetVariableData: []SetVariableData{{AttributeType: types.AttributeTarget, AttributeValue: "dummyValue", Component: types.Component{Name: "component1", EVSE: &types.EVSE{ID: -1}}, Variable: variable}}}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *provisioningTestSuite) TestSetVariablesResponseValidation() {
	t := suite.T()
	component := types.Component{Name: "component1", Instance: "instance1", EVSE: &types.EVSE{ID: 2, ConnectorID: tests.NewInt(2)}}
	variable := types.Variable{Name: "variable1", Instance: "instance1"}
	var confirmationTable = []tests.GenericTestEntry{
		{SetVariablesResponse{SetVariableResult: []SetVariableResult{{AttributeType: types.AttributeTarget, AttributeStatus: SetVariableStatusAccepted, Component: component, Variable: variable, StatusInfo: types.NewStatusInfo("200", "")}}}, true},
		{SetVariablesResponse{SetVariableResult: []SetVariableResult{{AttributeType: types.AttributeTarget, AttributeStatus: SetVariableStatusAccepted, Component: component, Variable: variable}}}, true},
		{SetVariablesResponse{SetVariableResult: []SetVariableResult{{AttributeStatus: SetVariableStatusAccepted, Component: component, Variable: variable}}}, true},
		{SetVariablesResponse{SetVariableResult: []SetVariableResult{{Component: component, Variable: variable}}}, false},
		{SetVariablesResponse{SetVariableResult: []SetVariableResult{{AttributeStatus: SetVariableStatusAccepted, Variable: variable}}}, false},
		{SetVariablesResponse{SetVariableResult: []SetVariableResult{{AttributeStatus: SetVariableStatusAccepted, Component: component}}}, false},
		{SetVariablesResponse{SetVariableResult: []SetVariableResult{}}, false},
		{SetVariablesResponse{}, false},
		{SetVariablesResponse{SetVariableResult: []SetVariableResult{{AttributeType: "invalidAttribute", AttributeStatus: SetVariableStatusAccepted, Component: component, Variable: variable}}}, false},
		{SetVariablesResponse{SetVariableResult: []SetVariableResult{{AttributeType: types.AttributeTarget, AttributeStatus: "invalidStatus", Component: component, Variable: variable}}}, false},
		{SetVariablesResponse{SetVariableResult: []SetVariableResult{{AttributeType: types.AttributeTarget, AttributeStatus: SetVariableStatusAccepted, Component: types.Component{}, Variable: variable}}}, false},
		{SetVariablesResponse{SetVariableResult: []SetVariableResult{{AttributeType: types.AttributeTarget, AttributeStatus: SetVariableStatusAccepted, Component: component, Variable: types.Variable{}}}}, false},
		{SetVariablesResponse{SetVariableResult: []SetVariableResult{{AttributeType: types.AttributeTarget, AttributeStatus: SetVariableStatusAccepted, Component: component, Variable: variable, StatusInfo: types.NewStatusInfo("", "")}}}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *provisioningTestSuite) TestSetVariablesFeature() {
	feature := SetVariablesFeature{}
	suite.Equal(SetVariablesFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(SetVariablesRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(SetVariablesResponse{}), feature.GetResponseType())
}

func (suite *provisioningTestSuite) TestNewSetVariablesRequest() {
	component := types.Component{Name: "comp1"}
	variable := types.Variable{Name: "var1"}
	data := []SetVariableData{
		{AttributeValue: "val1", Component: component, Variable: variable},
	}
	req := NewSetVariablesRequest(data)
	suite.NotNil(req)
	suite.Equal(SetVariablesFeatureName, req.GetFeatureName())
	suite.Equal(data, req.SetVariableData)
}

func (suite *provisioningTestSuite) TestNewSetVariablesResponse() {
	component := types.Component{Name: "comp1"}
	variable := types.Variable{Name: "var1"}
	result := []SetVariableResult{
		{AttributeStatus: SetVariableStatusAccepted, Component: component, Variable: variable},
	}
	resp := NewSetVariablesResponse(result)
	suite.NotNil(resp)
	suite.Equal(SetVariablesFeatureName, resp.GetFeatureName())
	suite.Equal(result, resp.SetVariableResult)
}
