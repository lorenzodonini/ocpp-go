package ocpp2_test

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/diagnostics"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"

	"github.com/stretchr/testify/mock"
)

// Test
func (suite *OcppV2TestSuite) TestSetVariableMonitoringRequestValidation() {
	var requestTable = []GenericTestEntry{
		{diagnostics.SetVariableMonitoringRequest{MonitoringData: []diagnostics.SetMonitoringData{{ID: newInt(2), Transaction: true, Value: 42.0, Type: diagnostics.MonitorUpperThreshold, Severity: 5, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}}}, true},
		{diagnostics.SetVariableMonitoringRequest{MonitoringData: []diagnostics.SetMonitoringData{{Transaction: true, Value: 42.0, Type: diagnostics.MonitorUpperThreshold, Severity: 5, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}}}, true},
		{diagnostics.SetVariableMonitoringRequest{MonitoringData: []diagnostics.SetMonitoringData{{Value: 42.0, Type: diagnostics.MonitorUpperThreshold, Severity: 5, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}}}, true},
		{diagnostics.SetVariableMonitoringRequest{MonitoringData: []diagnostics.SetMonitoringData{{Type: diagnostics.MonitorUpperThreshold, Severity: 5, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}}}, true},
		{diagnostics.SetVariableMonitoringRequest{MonitoringData: []diagnostics.SetMonitoringData{{Type: diagnostics.MonitorUpperThreshold, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}}}, true},
		{diagnostics.SetVariableMonitoringRequest{MonitoringData: []diagnostics.SetMonitoringData{{Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}}}, false},
		{diagnostics.SetVariableMonitoringRequest{MonitoringData: []diagnostics.SetMonitoringData{{Type: diagnostics.MonitorUpperThreshold, Variable: types.Variable{Name: "variable1"}}}}, false},
		{diagnostics.SetVariableMonitoringRequest{MonitoringData: []diagnostics.SetMonitoringData{{Type: diagnostics.MonitorUpperThreshold, Component: types.Component{Name: "component1"}}}}, false},
		{diagnostics.SetVariableMonitoringRequest{MonitoringData: []diagnostics.SetMonitoringData{{}}}, false},
		{diagnostics.SetVariableMonitoringRequest{MonitoringData: []diagnostics.SetMonitoringData{}}, false},
		{diagnostics.SetVariableMonitoringRequest{}, false},
		{diagnostics.SetVariableMonitoringRequest{MonitoringData: []diagnostics.SetMonitoringData{{ID: newInt(2), Transaction: true, Value: 42.0, Type: "invalidType", Severity: 5, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}}}, false},
		{diagnostics.SetVariableMonitoringRequest{MonitoringData: []diagnostics.SetMonitoringData{{ID: newInt(2), Transaction: true, Value: 42.0, Type: diagnostics.MonitorUpperThreshold, Severity: -1, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}}}, false},
		{diagnostics.SetVariableMonitoringRequest{MonitoringData: []diagnostics.SetMonitoringData{{ID: newInt(2), Transaction: true, Value: 42.0, Type: diagnostics.MonitorUpperThreshold, Severity: 10, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}}}, false},
		{diagnostics.SetVariableMonitoringRequest{MonitoringData: []diagnostics.SetMonitoringData{{ID: newInt(2), Transaction: true, Value: 42.0, Type: diagnostics.MonitorUpperThreshold, Severity: 5, Component: types.Component{}, Variable: types.Variable{}}}}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV2TestSuite) TestSetVariableMonitoringResponseValidation() {
	var confirmationTable = []GenericTestEntry{
		{diagnostics.SetVariableMonitoringResponse{MonitoringResult: []diagnostics.SetMonitoringResult{{ID: newInt(2), Status: diagnostics.SetMonitoringStatusAccepted, Type: diagnostics.MonitorUpperThreshold, Severity: 5, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}, StatusInfo: types.NewStatusInfo("200", "")}}}, true},
		{diagnostics.SetVariableMonitoringResponse{MonitoringResult: []diagnostics.SetMonitoringResult{{ID: newInt(2), Status: diagnostics.SetMonitoringStatusAccepted, Type: diagnostics.MonitorUpperThreshold, Severity: 5, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}}}, true},
		{diagnostics.SetVariableMonitoringResponse{MonitoringResult: []diagnostics.SetMonitoringResult{{Status: diagnostics.SetMonitoringStatusAccepted, Type: diagnostics.MonitorUpperThreshold, Severity: 5, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}}}, true},
		{diagnostics.SetVariableMonitoringResponse{MonitoringResult: []diagnostics.SetMonitoringResult{{Status: diagnostics.SetMonitoringStatusAccepted, Type: diagnostics.MonitorUpperThreshold, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}}}, true},
		{diagnostics.SetVariableMonitoringResponse{MonitoringResult: []diagnostics.SetMonitoringResult{{Status: diagnostics.SetMonitoringStatusAccepted, Type: diagnostics.MonitorUpperThreshold, Component: types.Component{Name: "component1"}}}}, false},
		{diagnostics.SetVariableMonitoringResponse{MonitoringResult: []diagnostics.SetMonitoringResult{{Status: diagnostics.SetMonitoringStatusAccepted, Type: diagnostics.MonitorUpperThreshold, Variable: types.Variable{Name: "variable1"}}}}, false},
		{diagnostics.SetVariableMonitoringResponse{MonitoringResult: []diagnostics.SetMonitoringResult{{Status: diagnostics.SetMonitoringStatusAccepted, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}}}, false},
		{diagnostics.SetVariableMonitoringResponse{MonitoringResult: []diagnostics.SetMonitoringResult{{Type: diagnostics.MonitorUpperThreshold, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}}}, false},
		{diagnostics.SetVariableMonitoringResponse{}, false},
		{diagnostics.SetVariableMonitoringResponse{MonitoringResult: []diagnostics.SetMonitoringResult{{ID: newInt(2), Status: "invalidStatus", Type: diagnostics.MonitorUpperThreshold, Severity: 5, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}, StatusInfo: types.NewStatusInfo("200", "")}}}, false},
		{diagnostics.SetVariableMonitoringResponse{MonitoringResult: []diagnostics.SetMonitoringResult{{ID: newInt(2), Status: diagnostics.SetMonitoringStatusAccepted, Type: "invalidType", Severity: 5, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}, StatusInfo: types.NewStatusInfo("200", "")}}}, false},
		{diagnostics.SetVariableMonitoringResponse{MonitoringResult: []diagnostics.SetMonitoringResult{{ID: newInt(2), Status: diagnostics.SetMonitoringStatusAccepted, Type: diagnostics.MonitorUpperThreshold, Severity: -1, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}, StatusInfo: types.NewStatusInfo("200", "")}}}, false},
		{diagnostics.SetVariableMonitoringResponse{MonitoringResult: []diagnostics.SetMonitoringResult{{ID: newInt(2), Status: diagnostics.SetMonitoringStatusAccepted, Type: diagnostics.MonitorUpperThreshold, Severity: 10, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}, StatusInfo: types.NewStatusInfo("200", "")}}}, false},
		{diagnostics.SetVariableMonitoringResponse{MonitoringResult: []diagnostics.SetMonitoringResult{{ID: newInt(2), Status: diagnostics.SetMonitoringStatusAccepted, Type: diagnostics.MonitorUpperThreshold, Severity: 5, Component: types.Component{Name: ""}, Variable: types.Variable{Name: "variable1"}, StatusInfo: types.NewStatusInfo("200", "")}}}, false},
		{diagnostics.SetVariableMonitoringResponse{MonitoringResult: []diagnostics.SetMonitoringResult{{ID: newInt(2), Status: diagnostics.SetMonitoringStatusAccepted, Type: diagnostics.MonitorUpperThreshold, Severity: 5, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: ""}, StatusInfo: types.NewStatusInfo("200", "")}}}, false},
		{diagnostics.SetVariableMonitoringResponse{MonitoringResult: []diagnostics.SetMonitoringResult{{ID: newInt(2), Status: diagnostics.SetMonitoringStatusAccepted, Type: diagnostics.MonitorUpperThreshold, Severity: 5, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}, StatusInfo: types.NewStatusInfo("", "")}}}, false},
	}
	ExecuteGenericTestTable(suite, confirmationTable)
}

func (suite *OcppV2TestSuite) TestSetVariableMonitoringE2EMocked() {
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	monitoringData := diagnostics.SetMonitoringData{ID: newInt(2), Transaction: false, Value: 42.0, Type: diagnostics.MonitorUpperThreshold, Severity: 5, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}
	monitoringResult := diagnostics.SetMonitoringResult{ID: newInt(2), Status: diagnostics.SetMonitoringStatusAccepted, Type: diagnostics.MonitorUpperThreshold, Severity: 5, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}, StatusInfo: types.NewStatusInfo("200", "")}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"setMonitoringData":[{"id":%v,"value":%v,"type":"%v","severity":%v,"component":{"name":"%v"},"variable":{"name":"%v"}}]}]`,
		messageId, diagnostics.SetVariableMonitoringFeatureName, *monitoringData.ID, monitoringData.Value, monitoringData.Type, monitoringData.Severity, monitoringData.Component.Name, monitoringData.Variable.Name)
	responseJson := fmt.Sprintf(`[3,"%v",{"setMonitoringResult":[{"id":%v,"status":"%v","type":"%v","severity":%v,"component":{"name":"%v"},"variable":{"name":"%v"},"statusInfo":{"reasonCode":"%v"}}]}]`,
		messageId, *monitoringResult.ID, monitoringResult.Status, monitoringResult.Type, monitoringResult.Severity, monitoringResult.Component.Name, monitoringResult.Variable.Name, monitoringResult.StatusInfo.ReasonCode)
	setMonitoringVariableResponse := diagnostics.NewSetVariableMonitoringResponse([]diagnostics.SetMonitoringResult{monitoringResult})
	channel := NewMockWebSocket(wsId)

	handler := &MockChargingStationDiagnosticsHandler{}
	handler.On("OnSetVariableMonitoring", mock.Anything).Return(setMonitoringVariableResponse, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*diagnostics.SetVariableMonitoringRequest)
		suite.Require().True(ok)
		suite.Require().NotNil(request)
		suite.Require().NotNil(request.MonitoringData)
		suite.Require().Len(request.MonitoringData, 1)
		suite.Equal(*monitoringData.ID, *request.MonitoringData[0].ID)
		suite.Equal(monitoringData.Transaction, request.MonitoringData[0].Transaction)
		suite.Equal(monitoringData.Value, request.MonitoringData[0].Value)
		suite.Equal(monitoringData.Type, request.MonitoringData[0].Type)
		suite.Equal(monitoringData.Severity, request.MonitoringData[0].Severity)
		suite.Equal(monitoringData.Component.Name, request.MonitoringData[0].Component.Name)
		suite.Equal(monitoringData.Variable.Name, request.MonitoringData[0].Variable.Name)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	suite.Require().Nil(err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.SetVariableMonitoring(wsId, func(response *diagnostics.SetVariableMonitoringResponse, err error) {
		suite.Require().Nil(err)
		suite.Require().NotNil(response)
		suite.Require().NotNil(response.MonitoringResult)
		suite.Require().Len(response.MonitoringResult, 1)
		suite.Equal(*monitoringResult.ID, *response.MonitoringResult[0].ID)
		suite.Equal(monitoringResult.Status, response.MonitoringResult[0].Status)
		suite.Equal(monitoringResult.Type, response.MonitoringResult[0].Type)
		suite.Equal(monitoringResult.Severity, response.MonitoringResult[0].Severity)
		suite.Equal(monitoringResult.Component.Name, response.MonitoringResult[0].Component.Name)
		suite.Equal(monitoringResult.Variable, response.MonitoringResult[0].Variable)
		suite.Require().NotNil(response.MonitoringResult[0].StatusInfo)
		suite.Equal(monitoringResult.StatusInfo.ReasonCode, response.MonitoringResult[0].StatusInfo.ReasonCode)
		resultChannel <- true
	}, []diagnostics.SetMonitoringData{monitoringData})
	suite.Require().Nil(err)
	result := <-resultChannel
	suite.True(result)
}

func (suite *OcppV2TestSuite) TestSetVariableMonitoringInvalidEndpoint() {
	messageId := defaultMessageId
	monitoringData := diagnostics.SetMonitoringData{ID: newInt(2), Transaction: false, Value: 42.0, Type: diagnostics.MonitorUpperThreshold, Severity: 5, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}
	request := diagnostics.NewSetVariableMonitoringRequest([]diagnostics.SetMonitoringData{monitoringData})
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"setMonitoringData":[{"id":%v,"value":%v,"type":"%v","severity":%v,"component":{"name":"%v"},"variable":{"name":"%v"}}]}]`,
		messageId, diagnostics.SetVariableMonitoringFeatureName, *monitoringData.ID, monitoringData.Value, monitoringData.Type, monitoringData.Severity, monitoringData.Component.Name, monitoringData.Variable.Name)
	testUnsupportedRequestFromChargingStation(suite, request, requestJson, messageId)
}
