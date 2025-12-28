package ocpp2_test

import (
	"fmt"

	"github.com/stretchr/testify/mock"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/provisioning"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

// Test
func (suite *OcppV2TestSuite) TestSetVariablesRequestValidation() {
	component := types.Component{Name: "component1", Instance: "instance1", EVSE: &types.EVSE{ID: 2, ConnectorID: newInt(2)}}
	variable := types.Variable{Name: "variable1", Instance: "instance1"}

	var requestTable = []GenericTestEntry{
		{provisioning.SetVariablesRequest{SetVariableData: []provisioning.SetVariableData{{AttributeType: types.AttributeTarget, AttributeValue: "dummyValue", Component: component, Variable: variable}}}, true},
		{provisioning.SetVariablesRequest{SetVariableData: []provisioning.SetVariableData{{AttributeValue: "dummyValue", Component: component, Variable: variable}}}, true},
		{provisioning.SetVariablesRequest{SetVariableData: []provisioning.SetVariableData{{AttributeValue: "dummyValue", Component: types.Component{Name: "component1"}, Variable: variable}}}, true},
		{provisioning.SetVariablesRequest{SetVariableData: []provisioning.SetVariableData{{AttributeValue: "dummyValue", Component: component, Variable: types.Variable{Name: "variable1"}}}}, true},
		{provisioning.SetVariablesRequest{SetVariableData: []provisioning.SetVariableData{{Component: component, Variable: variable}}}, false},
		{provisioning.SetVariablesRequest{SetVariableData: []provisioning.SetVariableData{{AttributeValue: "dummyValue", Variable: variable}}}, false},
		{provisioning.SetVariablesRequest{SetVariableData: []provisioning.SetVariableData{{AttributeValue: "dummyValue", Component: component}}}, false},
		{provisioning.SetVariablesRequest{SetVariableData: []provisioning.SetVariableData{}}, false},
		{provisioning.SetVariablesRequest{}, false},
		{provisioning.SetVariablesRequest{SetVariableData: []provisioning.SetVariableData{{AttributeType: "invalidAttribute", AttributeValue: "dummyValue", Component: component, Variable: variable}}}, false},
		{provisioning.SetVariablesRequest{SetVariableData: []provisioning.SetVariableData{{AttributeType: types.AttributeTarget, AttributeValue: ">1000....................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................", Component: component, Variable: variable}}}, false},
		{provisioning.SetVariablesRequest{SetVariableData: []provisioning.SetVariableData{{AttributeType: types.AttributeTarget, AttributeValue: "dummyValue", Variable: variable}}}, false},
		{provisioning.SetVariablesRequest{SetVariableData: []provisioning.SetVariableData{{AttributeType: types.AttributeTarget, AttributeValue: "dummyValue", Component: component}}}, false},
		{provisioning.SetVariablesRequest{SetVariableData: []provisioning.SetVariableData{{AttributeType: types.AttributeTarget, AttributeValue: "dummyValue", Component: types.Component{Name: "component1", EVSE: &types.EVSE{ID: -1}}, Variable: variable}}}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV2TestSuite) TestSetVariablesResponseValidation() {
	component := types.Component{Name: "component1", Instance: "instance1", EVSE: &types.EVSE{ID: 2, ConnectorID: newInt(2)}}
	variable := types.Variable{Name: "variable1", Instance: "instance1"}
	var confirmationTable = []GenericTestEntry{
		{provisioning.SetVariablesResponse{SetVariableResult: []provisioning.SetVariableResult{{AttributeType: types.AttributeTarget, AttributeStatus: provisioning.SetVariableStatusAccepted, Component: component, Variable: variable, StatusInfo: types.NewStatusInfo("200", "")}}}, true},
		{provisioning.SetVariablesResponse{SetVariableResult: []provisioning.SetVariableResult{{AttributeType: types.AttributeTarget, AttributeStatus: provisioning.SetVariableStatusAccepted, Component: component, Variable: variable}}}, true},
		{provisioning.SetVariablesResponse{SetVariableResult: []provisioning.SetVariableResult{{AttributeStatus: provisioning.SetVariableStatusAccepted, Component: component, Variable: variable}}}, true},
		{provisioning.SetVariablesResponse{SetVariableResult: []provisioning.SetVariableResult{{Component: component, Variable: variable}}}, false},
		{provisioning.SetVariablesResponse{SetVariableResult: []provisioning.SetVariableResult{{AttributeStatus: provisioning.SetVariableStatusAccepted, Variable: variable}}}, false},
		{provisioning.SetVariablesResponse{SetVariableResult: []provisioning.SetVariableResult{{AttributeStatus: provisioning.SetVariableStatusAccepted, Component: component}}}, false},
		{provisioning.SetVariablesResponse{SetVariableResult: []provisioning.SetVariableResult{}}, false},
		{provisioning.SetVariablesResponse{}, false},
		{provisioning.SetVariablesResponse{SetVariableResult: []provisioning.SetVariableResult{{AttributeType: "invalidAttribute", AttributeStatus: provisioning.SetVariableStatusAccepted, Component: component, Variable: variable}}}, false},
		{provisioning.SetVariablesResponse{SetVariableResult: []provisioning.SetVariableResult{{AttributeType: types.AttributeTarget, AttributeStatus: "invalidStatus", Component: component, Variable: variable}}}, false},
		{provisioning.SetVariablesResponse{SetVariableResult: []provisioning.SetVariableResult{{AttributeType: types.AttributeTarget, AttributeStatus: provisioning.SetVariableStatusAccepted, Component: types.Component{}, Variable: variable}}}, false},
		{provisioning.SetVariablesResponse{SetVariableResult: []provisioning.SetVariableResult{{AttributeType: types.AttributeTarget, AttributeStatus: provisioning.SetVariableStatusAccepted, Component: component, Variable: types.Variable{}}}}, false},
		{provisioning.SetVariablesResponse{SetVariableResult: []provisioning.SetVariableResult{{AttributeType: types.AttributeTarget, AttributeStatus: provisioning.SetVariableStatusAccepted, Component: component, Variable: variable, StatusInfo: types.NewStatusInfo("", "")}}}, false},
	}
	ExecuteGenericTestTable(suite, confirmationTable)
}

func (suite *OcppV2TestSuite) TestSetVariablesE2EMocked() {
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	attributeType := types.AttributeTarget
	component := types.Component{Name: "component1", Instance: "instance1", EVSE: &types.EVSE{ID: 2, ConnectorID: newInt(2)}}
	variable := types.Variable{Name: "variable1", Instance: "instance1"}
	variableData := provisioning.SetVariableData{
		AttributeType:  attributeType,
		AttributeValue: "dummyValue",
		Component:      component,
		Variable:       variable,
	}
	statusInfo := types.NewStatusInfo("200", "")
	variableResult := provisioning.SetVariableResult{
		AttributeType:   attributeType,
		AttributeStatus: provisioning.SetVariableStatusAccepted,
		Component:       component,
		Variable:        variable,
		StatusInfo:      statusInfo,
	}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"setVariableData":[{"attributeType":"%v","attributeValue":"%v","component":{"name":"%v","instance":"%v","evse":{"id":%v,"connectorId":%v}},"variable":{"name":"%v","instance":"%v"}}]}]`,
		messageId, provisioning.SetVariablesFeatureName, variableData.AttributeType, variableData.AttributeValue, variableData.Component.Name, variableData.Component.Instance, variableData.Component.EVSE.ID, *variableData.Component.EVSE.ConnectorID, variableData.Variable.Name, variableData.Variable.Instance)
	responseJson := fmt.Sprintf(`[3,"%v",{"setVariableResult":[{"attributeType":"%v","attributeStatus":"%v","component":{"name":"%v","instance":"%v","evse":{"id":%v,"connectorId":%v}},"variable":{"name":"%v","instance":"%v"},"statusInfo":{"reasonCode":"%v"}}]}]`,
		messageId, variableResult.AttributeType, variableResult.AttributeStatus, variableResult.Component.Name, variableResult.Component.Instance, variableResult.Component.EVSE.ID, *variableResult.Component.EVSE.ConnectorID, variableResult.Variable.Name, variableResult.Variable.Instance, statusInfo.ReasonCode)
	getVariablesResponse := provisioning.NewSetVariablesResponse([]provisioning.SetVariableResult{variableResult})
	channel := NewMockWebSocket(wsId)

	handler := &MockChargingStationProvisioningHandler{}
	handler.On("OnSetVariables", mock.Anything).Return(getVariablesResponse, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*provisioning.SetVariablesRequest)
		suite.Require().True(ok)
		suite.Require().NotNil(request)
		suite.Require().Len(request.SetVariableData, 1)
		suite.Equal(variableData.AttributeType, request.SetVariableData[0].AttributeType)
		suite.Equal(variableData.AttributeValue, request.SetVariableData[0].AttributeValue)
		suite.Equal(variableData.Component.Name, request.SetVariableData[0].Component.Name)
		suite.Equal(variableData.Component.Instance, request.SetVariableData[0].Component.Instance)
		suite.Equal(variableData.Component.EVSE.ID, request.SetVariableData[0].Component.EVSE.ID)
		suite.Equal(*variableData.Component.EVSE.ConnectorID, *request.SetVariableData[0].Component.EVSE.ConnectorID)
		suite.Equal(variableData.Variable.Name, request.SetVariableData[0].Variable.Name)
		suite.Equal(variableData.Variable.Instance, request.SetVariableData[0].Variable.Instance)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	suite.Require().Nil(err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.SetVariables(wsId, func(response *provisioning.SetVariablesResponse, err error) {
		suite.Require().Nil(err)
		suite.Require().NotNil(response)
		suite.Require().Len(response.SetVariableResult, 1)
		suite.Equal(variableResult.AttributeStatus, response.SetVariableResult[0].AttributeStatus)
		suite.Equal(variableResult.AttributeType, response.SetVariableResult[0].AttributeType)
		suite.Equal(variableResult.Component.Name, response.SetVariableResult[0].Component.Name)
		suite.Equal(variableResult.Component.Instance, response.SetVariableResult[0].Component.Instance)
		suite.Equal(variableResult.Component.EVSE.ID, response.SetVariableResult[0].Component.EVSE.ID)
		suite.Equal(*variableResult.Component.EVSE.ConnectorID, *response.SetVariableResult[0].Component.EVSE.ConnectorID)
		suite.Equal(variableResult.Variable.Name, response.SetVariableResult[0].Variable.Name)
		suite.Equal(variableResult.Variable.Instance, response.SetVariableResult[0].Variable.Instance)
		suite.Require().NotNil(response.SetVariableResult[0].StatusInfo)
		suite.Equal(statusInfo.ReasonCode, response.SetVariableResult[0].StatusInfo.ReasonCode)
		resultChannel <- true
	}, []provisioning.SetVariableData{variableData})
	suite.Require().Nil(err)
	result := <-resultChannel
	suite.True(result)
}

func (suite *OcppV2TestSuite) TestSetVariablesInvalidEndpoint() {
	messageId := defaultMessageId
	attributeType := types.AttributeTarget
	component := types.Component{Name: "component1", Instance: "instance1", EVSE: &types.EVSE{ID: 2, ConnectorID: newInt(2)}}
	variable := types.Variable{Name: "variable1", Instance: "instance1"}
	variableData := provisioning.SetVariableData{
		AttributeType:  attributeType,
		AttributeValue: "dummyValue",
		Component:      component,
		Variable:       variable,
	}
	request := provisioning.NewSetVariablesRequest([]provisioning.SetVariableData{variableData})
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"setVariableData":[{"attributeType":"%v","attributeValue":"%v","component":{"name":"%v","instance":"%v","evse":{"id":%v,"connectorId":%v}},"variable":{"name":"%v","instance":"%v"}}]}]`,
		messageId, provisioning.SetVariablesFeatureName, variableData.AttributeType, variableData.AttributeValue, variableData.Component.Name, variableData.Component.Instance, variableData.Component.EVSE.ID, *variableData.Component.EVSE.ConnectorID, variableData.Variable.Name, variableData.Variable.Instance)
	testUnsupportedRequestFromChargingStation(suite, request, requestJson, messageId)
}
