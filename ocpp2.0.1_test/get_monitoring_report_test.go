package ocpp2_test

import (
	"fmt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/diagnostics"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

// Test
func (suite *OcppV2TestSuite) TestGetMonitoringReportRequestValidation() {
	t := suite.T()
	componentVariables := []types.ComponentVariable{
		{
			Component: types.Component{Name: "component1", Instance: "instance1", EVSE: &types.EVSE{ID: 2, ConnectorID: newInt(2)}},
			Variable:  types.Variable{Name: "variable1", Instance: "instance1"},
		},
	}
	var requestTable = []GenericTestEntry{
		{diagnostics.GetMonitoringReportRequest{RequestID: newInt(42), MonitoringCriteria: []diagnostics.MonitoringCriteriaType{diagnostics.MonitoringCriteriaThresholdMonitoring, diagnostics.MonitoringCriteriaDeltaMonitoring, diagnostics.MonitoringCriteriaPeriodicMonitoring}, ComponentVariable: componentVariables}, true},
		{diagnostics.GetMonitoringReportRequest{RequestID: newInt(42), MonitoringCriteria: []diagnostics.MonitoringCriteriaType{}, ComponentVariable: componentVariables}, true},
		{diagnostics.GetMonitoringReportRequest{RequestID: newInt(42), ComponentVariable: componentVariables}, true},
		{diagnostics.GetMonitoringReportRequest{RequestID: newInt(42), ComponentVariable: []types.ComponentVariable{}}, true},
		{diagnostics.GetMonitoringReportRequest{RequestID: newInt(42)}, true},
		{diagnostics.GetMonitoringReportRequest{}, true},
		{diagnostics.GetMonitoringReportRequest{RequestID: newInt(-1)}, false},
		{diagnostics.GetMonitoringReportRequest{MonitoringCriteria: []diagnostics.MonitoringCriteriaType{diagnostics.MonitoringCriteriaThresholdMonitoring, diagnostics.MonitoringCriteriaDeltaMonitoring, diagnostics.MonitoringCriteriaPeriodicMonitoring, diagnostics.MonitoringCriteriaThresholdMonitoring}}, false},
		{diagnostics.GetMonitoringReportRequest{MonitoringCriteria: []diagnostics.MonitoringCriteriaType{"invalidMonitoringCriteria"}}, false},
		{diagnostics.GetMonitoringReportRequest{ComponentVariable: []types.ComponentVariable{{Variable: types.Variable{Name: "variable1", Instance: "instance1"}}}}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestGetMonitoringReportConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{diagnostics.GetMonitoringReportResponse{Status: types.GenericDeviceModelStatusAccepted}, true},
		{diagnostics.GetMonitoringReportResponse{Status: "invalidDeviceModelStatus"}, false},
		{diagnostics.GetMonitoringReportResponse{}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestGetMonitoringReportE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	requestID := newInt(42)
	monitoringCriteria := []diagnostics.MonitoringCriteriaType{diagnostics.MonitoringCriteriaThresholdMonitoring, diagnostics.MonitoringCriteriaPeriodicMonitoring}
	componentVariable := types.ComponentVariable{
		Component: types.Component{Name: "component1", Instance: "instance1", EVSE: &types.EVSE{ID: 2, ConnectorID: newInt(2)}},
		Variable:  types.Variable{Name: "variable1", Instance: "instance1"},
	}
	componentVariables := []types.ComponentVariable{componentVariable}
	status := types.GenericDeviceModelStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"requestId":%v,"monitoringCriteria":["%v","%v"],"componentVariable":[{"component":{"name":"%v","instance":"%v","evse":{"id":%v,"connectorId":%v}},"variable":{"name":"%v","instance":"%v"}}]}]`,
		messageId, diagnostics.GetMonitoringReportFeatureName, *requestID, monitoringCriteria[0], monitoringCriteria[1], componentVariable.Component.Name, componentVariable.Component.Instance, componentVariable.Component.EVSE.ID, *componentVariable.Component.EVSE.ConnectorID, componentVariable.Variable.Name, componentVariable.Variable.Instance)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	getMonitoringReportConfirmation := diagnostics.NewGetMonitoringReportResponse(status)
	channel := NewMockWebSocket(wsId)

	handler := &MockChargingStationDiagnosticsHandler{}
	handler.On("OnGetMonitoringReport", mock.Anything).Return(getMonitoringReportConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*diagnostics.GetMonitoringReportRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, *requestID, *request.RequestID)
		require.Len(t, request.MonitoringCriteria, len(monitoringCriteria))
		assert.Equal(t, monitoringCriteria[0], request.MonitoringCriteria[0])
		assert.Equal(t, monitoringCriteria[1], request.MonitoringCriteria[1])
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
	err = suite.csms.GetMonitoringReport(wsId, func(confirmation *diagnostics.GetMonitoringReportResponse, err error) {
		require.Nil(t, err)
		require.NotNil(t, confirmation)
		assert.Equal(t, status, confirmation.Status)
		resultChannel <- true
	}, func(request *diagnostics.GetMonitoringReportRequest) {
		request.RequestID = requestID
		request.MonitoringCriteria = monitoringCriteria
		request.ComponentVariable = componentVariables
	})
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV2TestSuite) TestGetMonitoringReportInvalidEndpoint() {
	messageId := defaultMessageId
	requestID := newInt(42)
	monitoringCriteria := []diagnostics.MonitoringCriteriaType{diagnostics.MonitoringCriteriaThresholdMonitoring, diagnostics.MonitoringCriteriaPeriodicMonitoring}
	componentVariable := types.ComponentVariable{
		Component: types.Component{Name: "component1", Instance: "instance1", EVSE: &types.EVSE{ID: 2, ConnectorID: newInt(2)}},
		Variable:  types.Variable{Name: "variable1", Instance: "instance1"},
	}
	GetMonitoringReportRequest := diagnostics.NewGetMonitoringReportRequest()
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"requestId":%v,"monitoringCriteria":["%v","%v"],"componentVariable":[{"component":{"name":"%v","instance":"%v","evse":{"id":%v,"connectorId":%v}},"variable":{"name":"%v","instance":"%v"}}]}]`,
		messageId, diagnostics.GetMonitoringReportFeatureName, *requestID, monitoringCriteria[0], monitoringCriteria[1], componentVariable.Component.Name, componentVariable.Component.Instance, componentVariable.Component.EVSE.ID, *componentVariable.Component.EVSE.ConnectorID, componentVariable.Variable.Name, componentVariable.Variable.Instance)
	testUnsupportedRequestFromChargingStation(suite, GetMonitoringReportRequest, requestJson, messageId)
}
