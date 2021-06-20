package ocpp2_test

import (
	"fmt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/provisioning"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

// Test
func (suite *OcppV2TestSuite) TestGetVariablesRequestValidation() {
	t := suite.T()
	component := types.Component{Name: "component1", Instance: "instance1", EVSE: &types.EVSE{ID: 2, ConnectorID: newInt(2)}}
	variable := types.Variable{Name: "variable1", Instance: "instance1"}

	var requestTable = []GenericTestEntry{
		{provisioning.GetVariablesRequest{GetVariableData: []provisioning.GetVariableData{{AttributeType: types.AttributeTarget, Component: component, Variable: variable}}}, true},
		{provisioning.GetVariablesRequest{GetVariableData: []provisioning.GetVariableData{{Component: component, Variable: variable}}}, true},
		{provisioning.GetVariablesRequest{GetVariableData: []provisioning.GetVariableData{{Component: types.Component{Name: "component1"}, Variable: variable}}}, true},
		{provisioning.GetVariablesRequest{GetVariableData: []provisioning.GetVariableData{{Component: component, Variable: types.Variable{Name: "variable1"}}}}, true},
		{provisioning.GetVariablesRequest{GetVariableData: []provisioning.GetVariableData{}}, false},
		{provisioning.GetVariablesRequest{}, false},
		{provisioning.GetVariablesRequest{GetVariableData: []provisioning.GetVariableData{{AttributeType: "invalidAttribute", Component: component, Variable: variable}}}, false},
		{provisioning.GetVariablesRequest{GetVariableData: []provisioning.GetVariableData{{AttributeType: types.AttributeTarget, Variable: variable}}}, false},
		{provisioning.GetVariablesRequest{GetVariableData: []provisioning.GetVariableData{{AttributeType: types.AttributeTarget, Component: component}}}, false},
		{provisioning.GetVariablesRequest{GetVariableData: []provisioning.GetVariableData{{AttributeType: types.AttributeTarget, Component: types.Component{Name: "component1", EVSE: &types.EVSE{ID: -1}}, Variable: variable}}}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestGetVariablesConfirmationValidation() {
	t := suite.T()
	component := types.Component{Name: "component1", Instance: "instance1", EVSE: &types.EVSE{ID: 2, ConnectorID: newInt(2)}}
	variable := types.Variable{Name: "variable1", Instance: "instance1"}
	var confirmationTable = []GenericTestEntry{
		{provisioning.GetVariablesResponse{GetVariableResult: []provisioning.GetVariableResult{{AttributeStatus: provisioning.GetVariableStatusAccepted, AttributeType: types.AttributeTarget, AttributeValue: "dummyValue", Component: component, Variable: variable}}}, true},
		{provisioning.GetVariablesResponse{GetVariableResult: []provisioning.GetVariableResult{{AttributeStatus: provisioning.GetVariableStatusAccepted, AttributeType: types.AttributeTarget, Component: component, Variable: variable}}}, true},
		{provisioning.GetVariablesResponse{GetVariableResult: []provisioning.GetVariableResult{{AttributeStatus: provisioning.GetVariableStatusAccepted, Component: component, Variable: variable}}}, true},
		{provisioning.GetVariablesResponse{GetVariableResult: []provisioning.GetVariableResult{{Component: component, Variable: variable}}}, false},
		{provisioning.GetVariablesResponse{GetVariableResult: []provisioning.GetVariableResult{{AttributeStatus: provisioning.GetVariableStatusAccepted, Variable: variable}}}, false},
		{provisioning.GetVariablesResponse{GetVariableResult: []provisioning.GetVariableResult{{AttributeStatus: provisioning.GetVariableStatusAccepted, Component: component}}}, false},
		{provisioning.GetVariablesResponse{GetVariableResult: []provisioning.GetVariableResult{}}, false},
		{provisioning.GetVariablesResponse{}, false},
		{provisioning.GetVariablesResponse{GetVariableResult: []provisioning.GetVariableResult{{AttributeStatus: provisioning.GetVariableStatusAccepted, AttributeType: "invalidAttribute", AttributeValue: "dummyValue", Component: component, Variable: variable}}}, false},
		{provisioning.GetVariablesResponse{GetVariableResult: []provisioning.GetVariableResult{{AttributeStatus: "invalidStatus", AttributeType: types.AttributeTarget, AttributeValue: "dummyValue", Component: component, Variable: variable}}}, false},
		{provisioning.GetVariablesResponse{GetVariableResult: []provisioning.GetVariableResult{{AttributeStatus: provisioning.GetVariableStatusAccepted, AttributeType: types.AttributeTarget, AttributeValue: ">1000....................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................", Component: component, Variable: variable}}}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestGetVariablesE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	attributeType := types.AttributeTarget
	component := types.Component{Name: "component1", Instance: "instance1", EVSE: &types.EVSE{ID: 2, ConnectorID: newInt(2)}}
	variable := types.Variable{Name: "variable1", Instance: "instance1"}
	variableData := provisioning.GetVariableData{
		AttributeType: attributeType,
		Component:     component,
		Variable:      variable,
	}
	variableResult := provisioning.GetVariableResult{
		AttributeStatus: provisioning.GetVariableStatusAccepted,
		AttributeType:   attributeType,
		AttributeValue:  "dummyValue",
		Component:       component,
		Variable:        variable,
	}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"getVariableData":[{"attributeType":"%v","component":{"name":"%v","instance":"%v","evse":{"id":%v,"connectorId":%v}},"variable":{"name":"%v","instance":"%v"}}]}]`,
		messageId, provisioning.GetVariablesFeatureName, variableData.AttributeType, variableData.Component.Name, variableData.Component.Instance, variableData.Component.EVSE.ID, *variableData.Component.EVSE.ConnectorID, variableData.Variable.Name, variableData.Variable.Instance)
	responseJson := fmt.Sprintf(`[3,"%v",{"getVariableResult":[{"attributeStatus":"%v","attributeType":"%v","attributeValue":"%v","component":{"name":"%v","instance":"%v","evse":{"id":%v,"connectorId":%v}},"variable":{"name":"%v","instance":"%v"}}]}]`,
		messageId, variableResult.AttributeStatus, variableResult.AttributeType, variableResult.AttributeValue, variableResult.Component.Name, variableResult.Component.Instance, variableResult.Component.EVSE.ID, *variableResult.Component.EVSE.ConnectorID, variableResult.Variable.Name, variableResult.Variable.Instance)
	getVariablesResponse := provisioning.NewGetVariablesResponse([]provisioning.GetVariableResult{variableResult})
	channel := NewMockWebSocket(wsId)

	handler := &MockChargingStationProvisioningHandler{}
	handler.On("OnGetVariables", mock.Anything).Return(getVariablesResponse, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*provisioning.GetVariablesRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		require.Len(t, request.GetVariableData, 1)
		assert.Equal(t, variableData.AttributeType, request.GetVariableData[0].AttributeType)
		assert.Equal(t, variableData.Component.Name, request.GetVariableData[0].Component.Name)
		assert.Equal(t, variableData.Component.Instance, request.GetVariableData[0].Component.Instance)
		assert.Equal(t, variableData.Component.EVSE.ID, request.GetVariableData[0].Component.EVSE.ID)
		assert.Equal(t, *variableData.Component.EVSE.ConnectorID, *request.GetVariableData[0].Component.EVSE.ConnectorID)
		assert.Equal(t, variableData.Variable.Name, request.GetVariableData[0].Variable.Name)
		assert.Equal(t, variableData.Variable.Instance, request.GetVariableData[0].Variable.Instance)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.GetVariables(wsId, func(response *provisioning.GetVariablesResponse, err error) {
		require.Nil(t, err)
		require.NotNil(t, response)
		require.Len(t, response.GetVariableResult, 1)
		assert.Equal(t, variableResult.AttributeStatus, response.GetVariableResult[0].AttributeStatus)
		assert.Equal(t, variableResult.AttributeType, response.GetVariableResult[0].AttributeType)
		assert.Equal(t, variableResult.AttributeValue, response.GetVariableResult[0].AttributeValue)
		assert.Equal(t, variableResult.Component.Name, response.GetVariableResult[0].Component.Name)
		assert.Equal(t, variableResult.Component.Instance, response.GetVariableResult[0].Component.Instance)
		assert.Equal(t, variableResult.Component.EVSE.ID, response.GetVariableResult[0].Component.EVSE.ID)
		assert.Equal(t, *variableResult.Component.EVSE.ConnectorID, *response.GetVariableResult[0].Component.EVSE.ConnectorID)
		assert.Equal(t, variableResult.Variable.Name, response.GetVariableResult[0].Variable.Name)
		assert.Equal(t, variableResult.Variable.Instance, response.GetVariableResult[0].Variable.Instance)
		resultChannel <- true
	}, []provisioning.GetVariableData{variableData})
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV2TestSuite) TestGetVariablesInvalidEndpoint() {
	messageId := defaultMessageId
	attributeType := types.AttributeTarget
	component := types.Component{Name: "component1", Instance: "instance1", EVSE: &types.EVSE{ID: 2, ConnectorID: newInt(2)}}
	variable := types.Variable{Name: "variable1", Instance: "instance1"}
	variableData := provisioning.GetVariableData{
		AttributeType: attributeType,
		Component:     component,
		Variable:      variable,
	}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"getVariableData":[{"attributeType":"%v","component":{"name":"%v","instance":"%v","evse":{"id":%v,"connectorId":%v}},"variable":{"name":"%v","instance":"%v"}}]}]`,
		messageId, provisioning.GetVariablesFeatureName, variableData.AttributeType, variableData.Component.Name, variableData.Component.Instance, variableData.Component.EVSE.ID, *variableData.Component.EVSE.ConnectorID, variableData.Variable.Name, variableData.Variable.Instance)
	getVariablesRequest := provisioning.NewGetVariablesRequest([]provisioning.GetVariableData{variableData})

	testUnsupportedRequestFromChargingStation(suite, getVariablesRequest, requestJson, messageId)
}
