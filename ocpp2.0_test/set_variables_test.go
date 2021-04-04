package ocpp2_test

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0/provisioning"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test
func (suite *OcppV2TestSuite) TestSetVariablesRequestValidation() {
	t := suite.T()
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
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestSetVariablesResponseValidation() {
	t := suite.T()
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
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestSetVariablesE2EMocked() {
	t := suite.T()
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

	handler := MockChargingStationProvisioningHandler{}
	handler.On("OnSetVariables", mock.Anything).Return(getVariablesResponse, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*provisioning.SetVariablesRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		require.Len(t, request.SetVariableData, 1)
		assert.Equal(t, variableData.AttributeType, request.SetVariableData[0].AttributeType)
		assert.Equal(t, variableData.AttributeValue, request.SetVariableData[0].AttributeValue)
		assert.Equal(t, variableData.Component.Name, request.SetVariableData[0].Component.Name)
		assert.Equal(t, variableData.Component.Instance, request.SetVariableData[0].Component.Instance)
		assert.Equal(t, variableData.Component.EVSE.ID, request.SetVariableData[0].Component.EVSE.ID)
		assert.Equal(t, *variableData.Component.EVSE.ConnectorID, *request.SetVariableData[0].Component.EVSE.ConnectorID)
		assert.Equal(t, variableData.Variable.Name, request.SetVariableData[0].Variable.Name)
		assert.Equal(t, variableData.Variable.Instance, request.SetVariableData[0].Variable.Instance)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.SetVariables(wsId, func(response *provisioning.SetVariablesResponse, err error) {
		require.Nil(t, err)
		require.NotNil(t, response)
		require.Len(t, response.SetVariableResult, 1)
		assert.Equal(t, variableResult.AttributeStatus, response.SetVariableResult[0].AttributeStatus)
		assert.Equal(t, variableResult.AttributeType, response.SetVariableResult[0].AttributeType)
		assert.Equal(t, variableResult.Component.Name, response.SetVariableResult[0].Component.Name)
		assert.Equal(t, variableResult.Component.Instance, response.SetVariableResult[0].Component.Instance)
		assert.Equal(t, variableResult.Component.EVSE.ID, response.SetVariableResult[0].Component.EVSE.ID)
		assert.Equal(t, *variableResult.Component.EVSE.ConnectorID, *response.SetVariableResult[0].Component.EVSE.ConnectorID)
		assert.Equal(t, variableResult.Variable.Name, response.SetVariableResult[0].Variable.Name)
		assert.Equal(t, variableResult.Variable.Instance, response.SetVariableResult[0].Variable.Instance)
		require.NotNil(t, response.SetVariableResult[0].StatusInfo)
		assert.Equal(t, statusInfo.ReasonCode, response.SetVariableResult[0].StatusInfo.ReasonCode)
		resultChannel <- true
	}, []provisioning.SetVariableData{variableData})
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
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
