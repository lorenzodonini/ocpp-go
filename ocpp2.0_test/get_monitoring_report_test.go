package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test
func (suite *OcppV2TestSuite) TestGetMonitoringReportRequestValidation() {
	t := suite.T()
	componentVariables := []ocpp2.ComponentVariable{
		{
			Component: ocpp2.Component{ Name: "component1", Instance: "instance1", EVSE: &ocpp2.EVSE{ID: 2, ConnectorID: newInt(2)}},
			Variable:  ocpp2.Variable{ Name: "variable1", Instance: "instance1"},
		},
	}
	var requestTable = []GenericTestEntry{
		{ocpp2.GetMonitoringReportRequest{RequestID: newInt(42), MonitoringCriteria: []ocpp2.MonitoringCriteriaType{ocpp2.MonitoringCriteriaThresholdMonitoring, ocpp2.MonitoringCriteriaDeltaMonitoring, ocpp2.MonitoringCriteriaPeriodicMonitoring}, ComponentVariable: componentVariables}, true},
		{ocpp2.GetMonitoringReportRequest{RequestID: newInt(42), MonitoringCriteria: []ocpp2.MonitoringCriteriaType{}, ComponentVariable: componentVariables}, true},
		{ocpp2.GetMonitoringReportRequest{RequestID: newInt(42), ComponentVariable: componentVariables}, true},
		{ocpp2.GetMonitoringReportRequest{RequestID: newInt(42), ComponentVariable: []ocpp2.ComponentVariable{}}, true},
		{ocpp2.GetMonitoringReportRequest{RequestID: newInt(42)}, true},
		{ocpp2.GetMonitoringReportRequest{}, true},
		{ocpp2.GetMonitoringReportRequest{RequestID: newInt(-1)}, false},
		{ocpp2.GetMonitoringReportRequest{MonitoringCriteria: []ocpp2.MonitoringCriteriaType{ocpp2.MonitoringCriteriaThresholdMonitoring, ocpp2.MonitoringCriteriaDeltaMonitoring, ocpp2.MonitoringCriteriaPeriodicMonitoring, ocpp2.MonitoringCriteriaThresholdMonitoring}}, false},
		{ocpp2.GetMonitoringReportRequest{MonitoringCriteria: []ocpp2.MonitoringCriteriaType{"invalidMonitoringCriteria"}}, false},
		{ocpp2.GetMonitoringReportRequest{ComponentVariable: []ocpp2.ComponentVariable{ { Variable: ocpp2.Variable{ Name: "variable1", Instance: "instance1"}}}}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestGetMonitoringReportConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{ocpp2.GetMonitoringReportConfirmation{Status: ocpp2.GenericDeviceModelStatusAccepted}, true},
		{ocpp2.GetMonitoringReportConfirmation{Status: "invalidDeviceModelStatus"}, false},
		{ocpp2.GetMonitoringReportConfirmation{}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestGetMonitoringReportE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	requestID := newInt(42)
	monitoringCriteria := []ocpp2.MonitoringCriteriaType{ocpp2.MonitoringCriteriaThresholdMonitoring, ocpp2.MonitoringCriteriaPeriodicMonitoring}
	componentVariable := ocpp2.ComponentVariable{
		Component: ocpp2.Component{ Name: "component1", Instance: "instance1", EVSE: &ocpp2.EVSE{ID: 2, ConnectorID: newInt(2)}},
		Variable:  ocpp2.Variable{ Name: "variable1", Instance: "instance1"},
	}
	componentVariables := []ocpp2.ComponentVariable{componentVariable}
	status := ocpp2.GenericDeviceModelStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"requestId":%v,"monitoringCriteria":["%v","%v"],"componentVariable":[{"component":{"name":"%v","instance":"%v","evse":{"id":%v,"connectorId":%v}},"variable":{"name":"%v","instance":"%v"}}]}]`,
		messageId, ocpp2.GetMonitoringReportFeatureName, *requestID, monitoringCriteria[0], monitoringCriteria[1], componentVariable.Component.Name, componentVariable.Component.Instance, componentVariable.Component.EVSE.ID, *componentVariable.Component.EVSE.ConnectorID, componentVariable.Variable.Name, componentVariable.Variable.Instance)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	getMonitoringReportConfirmation := ocpp2.NewGetMonitoringReportConfirmation(status)
	channel := NewMockWebSocket(wsId)

	coreListener := MockChargePointCoreListener{}
	coreListener.On("OnGetMonitoringReport", mock.Anything).Return(getMonitoringReportConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*ocpp2.GetMonitoringReportRequest)
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
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, coreListener, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.GetMonitoringReport(wsId, func(confirmation *ocpp2.GetMonitoringReportConfirmation, err error) {
		require.Nil(t, err)
		require.NotNil(t, confirmation)
		assert.Equal(t, status, confirmation.Status)
		resultChannel <- true
	}, func(request *ocpp2.GetMonitoringReportRequest) {
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
	monitoringCriteria := []ocpp2.MonitoringCriteriaType{ocpp2.MonitoringCriteriaThresholdMonitoring, ocpp2.MonitoringCriteriaPeriodicMonitoring}
	componentVariable := ocpp2.ComponentVariable{
		Component: ocpp2.Component{ Name: "component1", Instance: "instance1", EVSE: &ocpp2.EVSE{ID: 2, ConnectorID: newInt(2)}},
		Variable:  ocpp2.Variable{ Name: "variable1", Instance: "instance1"},
	}
	GetMonitoringReportRequest := ocpp2.NewGetMonitoringReportRequest()
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"requestId":%v,"monitoringCriteria":["%v","%v"],"componentVariable":[{"component":{"name":"%v","instance":"%v","evse":{"id":%v,"connectorId":%v}},"variable":{"name":"%v","instance":"%v"}}]}]`,
		messageId, ocpp2.GetMonitoringReportFeatureName, *requestID, monitoringCriteria[0], monitoringCriteria[1], componentVariable.Component.Name, componentVariable.Component.Instance, componentVariable.Component.EVSE.ID, *componentVariable.Component.EVSE.ConnectorID, componentVariable.Variable.Name, componentVariable.Variable.Instance)
	testUnsupportedRequestFromChargePoint(suite, GetMonitoringReportRequest, requestJson, messageId)
}
