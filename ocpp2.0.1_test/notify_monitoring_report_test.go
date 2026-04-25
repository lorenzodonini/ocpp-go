package ocpp2_test

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/diagnostics"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

// Test
func (suite *OcppV2TestSuite) TestNotifyMonitoringReportE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	requestID := 42
	tbc := true
	seqNo := 0
	generatedAt := types.NewDateTime(time.Now())
	varMonitoring := diagnostics.NewVariableMonitoring(1, false, 42.42, diagnostics.MonitorPeriodic, 0)
	monitoringData := diagnostics.MonitoringData{
		Component:          types.Component{Name: "component1"},
		Variable:           types.Variable{Name: "variable1"},
		VariableMonitoring: []diagnostics.VariableMonitoring{varMonitoring},
	}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"requestId":%v,"tbc":%v,"seqNo":%v,"generatedAt":"%v","monitor":[{"component":{"name":"%v"},"variable":{"name":"%v"},"variableMonitoring":[{"id":%v,"transaction":%v,"value":%v,"type":"%v","severity":%v}]}]}]`,
		messageId, diagnostics.NotifyMonitoringReportFeatureName, requestID, tbc, seqNo, generatedAt.FormatTimestamp(), monitoringData.Component.Name, monitoringData.Variable.Name, varMonitoring.ID, varMonitoring.Transaction, varMonitoring.Value, varMonitoring.Type, varMonitoring.Severity)
	responseJson := fmt.Sprintf(`[3,"%v",{}]`, messageId)
	response := diagnostics.NewNotifyMonitoringReportResponse()
	channel := NewMockWebSocket(wsId)

	handler := &MockCSMSDiagnosticsHandler{}
	handler.On("OnNotifyMonitoringReport", mock.AnythingOfType("string"), mock.Anything).Return(response, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*diagnostics.NotifyMonitoringReportRequest)
		require.True(t, ok)
		assert.Equal(t, requestID, request.RequestID)
		assert.Equal(t, tbc, request.Tbc)
		assert.Equal(t, seqNo, request.SeqNo)
		assertDateTimeEquality(t, generatedAt, request.GeneratedAt)
		require.Len(t, request.Monitor, 1)
		assert.Equal(t, monitoringData.Component.Name, request.Monitor[0].Component.Name)
		assert.Equal(t, monitoringData.Variable.Name, request.Monitor[0].Variable.Name)
		require.Len(t, request.Monitor[0].VariableMonitoring, len(monitoringData.VariableMonitoring))
		assert.Equal(t, monitoringData.VariableMonitoring[0].ID, request.Monitor[0].VariableMonitoring[0].ID)
		assert.Equal(t, monitoringData.VariableMonitoring[0].Transaction, request.Monitor[0].VariableMonitoring[0].Transaction)
		assert.Equal(t, monitoringData.VariableMonitoring[0].Type, request.Monitor[0].VariableMonitoring[0].Type)
		assert.Equal(t, monitoringData.VariableMonitoring[0].Value, request.Monitor[0].VariableMonitoring[0].Value)
		assert.Equal(t, monitoringData.VariableMonitoring[0].Severity, request.Monitor[0].VariableMonitoring[0].Severity)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	r, err := suite.chargingStation.NotifyMonitoringReport(requestID, seqNo, generatedAt, []diagnostics.MonitoringData{monitoringData}, func(request *diagnostics.NotifyMonitoringReportRequest) {
		request.Tbc = tbc
	})
	assert.Nil(t, err)
	assert.NotNil(t, r)
}

func (suite *OcppV2TestSuite) TestNotifyMonitoringReportInvalidEndpoint() {
	messageId := defaultMessageId
	requestID := 42
	tbc := true
	seqNo := 0
	generatedAt := types.NewDateTime(time.Now())
	varMonitoring := diagnostics.NewVariableMonitoring(1, false, 42.42, diagnostics.MonitorPeriodic, 0)
	monitoringData := diagnostics.MonitoringData{
		Component:          types.Component{Name: "component1"},
		Variable:           types.Variable{Name: "variable1"},
		VariableMonitoring: []diagnostics.VariableMonitoring{varMonitoring},
	}
	req := diagnostics.NewNotifyMonitoringReportRequest(requestID, seqNo, generatedAt, []diagnostics.MonitoringData{monitoringData})
	req.Tbc = tbc
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"requestId":%v,"tbc":%v,"seqNo":%v,"generatedAt":"%v","monitor":[{"component":{"name":"%v"},"variable":{"name":"%v"},"variableMonitoring":[{"id":%v,"transaction":%v,"value":%v,"type":"%v","severity":%v}]}]}]`,
		messageId, diagnostics.NotifyMonitoringReportFeatureName, requestID, tbc, seqNo, generatedAt.FormatTimestamp(), monitoringData.Component.Name, monitoringData.Variable.Name, varMonitoring.ID, varMonitoring.Transaction, varMonitoring.Value, varMonitoring.Type, varMonitoring.Severity)
	testUnsupportedRequestFromCentralSystem(suite, req, requestJson, messageId)
}
