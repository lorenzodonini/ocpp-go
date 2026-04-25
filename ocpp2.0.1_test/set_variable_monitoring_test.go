package ocpp2_test

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/diagnostics"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test
func (suite *OcppV2TestSuite) TestSetVariableMonitoringE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	monitoringData := diagnostics.SetMonitoringData{ID: tests.NewInt(2), Transaction: false, Value: 42.0, Type: diagnostics.MonitorUpperThreshold, Severity: 5, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}
	monitoringResult := diagnostics.SetMonitoringResult{ID: tests.NewInt(2), Status: diagnostics.SetMonitoringStatusAccepted, Type: diagnostics.MonitorUpperThreshold, Severity: 5, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}, StatusInfo: types.NewStatusInfo("200", "")}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"setMonitoringData":[{"id":%v,"value":%v,"type":"%v","severity":%v,"component":{"name":"%v"},"variable":{"name":"%v"}}]}]`,
		messageId, diagnostics.SetVariableMonitoringFeatureName, *monitoringData.ID, monitoringData.Value, monitoringData.Type, monitoringData.Severity, monitoringData.Component.Name, monitoringData.Variable.Name)
	responseJson := fmt.Sprintf(`[3,"%v",{"setMonitoringResult":[{"id":%v,"status":"%v","type":"%v","severity":%v,"component":{"name":"%v"},"variable":{"name":"%v"},"statusInfo":{"reasonCode":"%v"}}]}]`,
		messageId, *monitoringResult.ID, monitoringResult.Status, monitoringResult.Type, monitoringResult.Severity, monitoringResult.Component.Name, monitoringResult.Variable.Name, monitoringResult.StatusInfo.ReasonCode)
	setMonitoringVariableResponse := diagnostics.NewSetVariableMonitoringResponse([]diagnostics.SetMonitoringResult{monitoringResult})
	channel := NewMockWebSocket(wsId)

	handler := &MockChargingStationDiagnosticsHandler{}
	handler.On("OnSetVariableMonitoring", mock.Anything).Return(setMonitoringVariableResponse, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*diagnostics.SetVariableMonitoringRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		require.NotNil(t, request.MonitoringData)
		require.Len(t, request.MonitoringData, 1)
		assert.Equal(t, *monitoringData.ID, *request.MonitoringData[0].ID)
		assert.Equal(t, monitoringData.Transaction, request.MonitoringData[0].Transaction)
		assert.Equal(t, monitoringData.Value, request.MonitoringData[0].Value)
		assert.Equal(t, monitoringData.Type, request.MonitoringData[0].Type)
		assert.Equal(t, monitoringData.Severity, request.MonitoringData[0].Severity)
		assert.Equal(t, monitoringData.Component.Name, request.MonitoringData[0].Component.Name)
		assert.Equal(t, monitoringData.Variable.Name, request.MonitoringData[0].Variable.Name)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.SetVariableMonitoring(wsId, func(response *diagnostics.SetVariableMonitoringResponse, err error) {
		require.Nil(t, err)
		require.NotNil(t, response)
		require.NotNil(t, response.MonitoringResult)
		require.Len(t, response.MonitoringResult, 1)
		assert.Equal(t, *monitoringResult.ID, *response.MonitoringResult[0].ID)
		assert.Equal(t, monitoringResult.Status, response.MonitoringResult[0].Status)
		assert.Equal(t, monitoringResult.Type, response.MonitoringResult[0].Type)
		assert.Equal(t, monitoringResult.Severity, response.MonitoringResult[0].Severity)
		assert.Equal(t, monitoringResult.Component.Name, response.MonitoringResult[0].Component.Name)
		assert.Equal(t, monitoringResult.Variable, response.MonitoringResult[0].Variable)
		require.NotNil(t, response.MonitoringResult[0].StatusInfo)
		assert.Equal(t, monitoringResult.StatusInfo.ReasonCode, response.MonitoringResult[0].StatusInfo.ReasonCode)
		resultChannel <- true
	}, []diagnostics.SetMonitoringData{monitoringData})
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV2TestSuite) TestSetVariableMonitoringInvalidEndpoint() {
	messageId := defaultMessageId
	monitoringData := diagnostics.SetMonitoringData{ID: tests.NewInt(2), Transaction: false, Value: 42.0, Type: diagnostics.MonitorUpperThreshold, Severity: 5, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}
	request := diagnostics.NewSetVariableMonitoringRequest([]diagnostics.SetMonitoringData{monitoringData})
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"setMonitoringData":[{"id":%v,"value":%v,"type":"%v","severity":%v,"component":{"name":"%v"},"variable":{"name":"%v"}}]}]`,
		messageId, diagnostics.SetVariableMonitoringFeatureName, *monitoringData.ID, monitoringData.Value, monitoringData.Type, monitoringData.Severity, monitoringData.Component.Name, monitoringData.Variable.Name)
	testUnsupportedRequestFromChargingStation(suite, request, requestJson, messageId)
}
