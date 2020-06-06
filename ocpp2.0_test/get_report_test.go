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
func (suite *OcppV2TestSuite) TestGetReportRequestValidation() {
	t := suite.T()
	componentVariables := []types.ComponentVariable{
		{
			Component: types.Component{Name: "component1", Instance: "instance1", EVSE: &types.EVSE{ID: 2, ConnectorID: newInt(2)}},
			Variable:  types.Variable{Name: "variable1", Instance: "instance1"},
		},
	}
	var requestTable = []GenericTestEntry{
		{provisioning.GetReportRequest{RequestID: newInt(42), ComponentCriteria: []provisioning.ComponentCriterion{provisioning.ComponentCriterionActive, provisioning.ComponentCriterionEnabled, provisioning.ComponentCriterionAvailable, provisioning.ComponentCriterionProblem}, ComponentVariable: componentVariables}, true},
		{provisioning.GetReportRequest{RequestID: newInt(42), ComponentCriteria: []provisioning.ComponentCriterion{provisioning.ComponentCriterionActive, provisioning.ComponentCriterionEnabled, provisioning.ComponentCriterionAvailable, provisioning.ComponentCriterionProblem}, ComponentVariable: []types.ComponentVariable{}}, true},
		{provisioning.GetReportRequest{RequestID: newInt(42), ComponentCriteria: []provisioning.ComponentCriterion{provisioning.ComponentCriterionActive, provisioning.ComponentCriterionEnabled, provisioning.ComponentCriterionAvailable, provisioning.ComponentCriterionProblem}}, true},
		{provisioning.GetReportRequest{RequestID: newInt(42), ComponentCriteria: []provisioning.ComponentCriterion{}}, true},
		{provisioning.GetReportRequest{RequestID: newInt(42)}, true},
		{provisioning.GetReportRequest{}, true},
		{provisioning.GetReportRequest{RequestID: newInt(-1)}, false},
		{provisioning.GetReportRequest{ComponentCriteria: []provisioning.ComponentCriterion{"invalidComponentCriterion"}}, false},
		{provisioning.GetReportRequest{ComponentCriteria: []provisioning.ComponentCriterion{provisioning.ComponentCriterionActive, provisioning.ComponentCriterionEnabled, provisioning.ComponentCriterionAvailable, provisioning.ComponentCriterionProblem, provisioning.ComponentCriterionActive}}, false},
		{provisioning.GetReportRequest{ComponentVariable: []types.ComponentVariable{{Variable: types.Variable{Name: "variable1", Instance: "instance1"}}}}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestGetReportConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{provisioning.GetReportResponse{Status: types.GenericDeviceModelStatusAccepted}, true},
		{provisioning.GetReportResponse{Status: "invalidDeviceModelStatus"}, false},
		{provisioning.GetReportResponse{}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestGetReportE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	requestID := newInt(42)
	componentCriteria := []provisioning.ComponentCriterion{provisioning.ComponentCriterionActive, provisioning.ComponentCriterionAvailable}
	componentVariable := types.ComponentVariable{
		Component: types.Component{Name: "component1", Instance: "instance1", EVSE: &types.EVSE{ID: 2, ConnectorID: newInt(2)}},
		Variable:  types.Variable{Name: "variable1", Instance: "instance1"},
	}
	componentVariables := []types.ComponentVariable{componentVariable}
	status := types.GenericDeviceModelStatusAccepted

	requestJson := fmt.Sprintf(`[2,"%v","%v",{"requestId":%v,"componentCriteria":["%v","%v"],"componentVariable":[{"component":{"name":"%v","instance":"%v","evse":{"id":%v,"connectorId":%v}},"variable":{"name":"%v","instance":"%v"}}]}]`,
		messageId, provisioning.GetReportFeatureName, *requestID, componentCriteria[0], componentCriteria[1], componentVariable.Component.Name, componentVariable.Component.Instance, componentVariable.Component.EVSE.ID, *componentVariable.Component.EVSE.ConnectorID, componentVariable.Variable.Name, componentVariable.Variable.Instance)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	getReportConfirmation := provisioning.NewGetReportResponse(status)
	channel := NewMockWebSocket(wsId)

	handler := MockChargingStationProvisioningHandler{}
	handler.On("OnGetReport", mock.Anything).Return(getReportConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*provisioning.GetReportRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, requestID, request.RequestID)
		require.Len(t, request.ComponentCriteria, len(componentCriteria))
		assert.Equal(t, componentCriteria[0], request.ComponentCriteria[0])
		assert.Equal(t, componentCriteria[1], request.ComponentCriteria[1])
		require.Len(t, request.ComponentVariable, len(componentVariables))
		assert.Equal(t, componentVariable.Component.Name, request.ComponentVariable[0].Component.Name)
		assert.Equal(t, componentVariable.Component.Instance, request.ComponentVariable[0].Component.Instance)
		require.NotNil(t, request.ComponentVariable[0].Component.EVSE)
		assert.Equal(t, componentVariable.Component.EVSE.ID, request.ComponentVariable[0].Component.EVSE.ID)
		assert.Equal(t, *componentVariable.Component.EVSE.ConnectorID, *request.ComponentVariable[0].Component.EVSE.ConnectorID)
		assert.Equal(t, componentVariable.Variable.Name, request.ComponentVariable[0].Variable.Name)
		assert.Equal(t, componentVariable.Variable.Instance, request.ComponentVariable[0].Variable.Instance)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.GetReport(wsId, func(confirmation *provisioning.GetReportResponse, err error) {
		require.Nil(t, err)
		require.NotNil(t, confirmation)
		assert.Equal(t, status, confirmation.Status)
		resultChannel <- true
	}, func(request *provisioning.GetReportRequest) {
		request.RequestID = requestID
		request.ComponentCriteria = componentCriteria
		request.ComponentVariable = componentVariables
	})
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV2TestSuite) TestGetReportInvalidEndpoint() {
	messageId := defaultMessageId
	requestID := newInt(42)
	componentCriteria := []provisioning.ComponentCriterion{provisioning.ComponentCriterionActive, provisioning.ComponentCriterionAvailable}
	componentVariable := types.ComponentVariable{
		Component: types.Component{Name: "component1", Instance: "instance1", EVSE: &types.EVSE{ID: 2, ConnectorID: newInt(2)}},
		Variable:  types.Variable{Name: "variable1", Instance: "instance1"},
	}
	componentVariables := []types.ComponentVariable{componentVariable}
	getReportRequest := provisioning.NewGetReportRequest()
	getReportRequest.RequestID = requestID
	getReportRequest.ComponentCriteria = componentCriteria
	getReportRequest.ComponentVariable = componentVariables
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"requestId":%v,"componentCriteria":["%v","%v"],"componentVariable":[{"component":{"name":"%v","instance":"%v","evse":{"id":%v,"connectorId":%v}},"variable":{"name":"%v","instance":"%v"}}]}]`,
		messageId, provisioning.GetReportFeatureName, *requestID, componentCriteria[0], componentCriteria[1], componentVariable.Component.Name, componentVariable.Component.Instance, componentVariable.Component.EVSE.ID, *componentVariable.Component.EVSE.ConnectorID, componentVariable.Variable.Name, componentVariable.Variable.Instance)

	testUnsupportedRequestFromChargingStation(suite, getReportRequest, requestJson, messageId)
}
