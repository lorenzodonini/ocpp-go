package ocpp2_test

import (
	"fmt"

	"github.com/stretchr/testify/mock"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/diagnostics"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

// Test
func (suite *OcppV2TestSuite) TestSetMonitoringLevelRequestValidation() {
	var requestTable = []GenericTestEntry{
		{diagnostics.SetMonitoringLevelRequest{Severity: 0}, true},
		{diagnostics.SetMonitoringLevelRequest{Severity: 1}, true},
		{diagnostics.SetMonitoringLevelRequest{Severity: 2}, true},
		{diagnostics.SetMonitoringLevelRequest{Severity: 3}, true},
		{diagnostics.SetMonitoringLevelRequest{Severity: 4}, true},
		{diagnostics.SetMonitoringLevelRequest{Severity: 5}, true},
		{diagnostics.SetMonitoringLevelRequest{Severity: 6}, true},
		{diagnostics.SetMonitoringLevelRequest{Severity: 7}, true},
		{diagnostics.SetMonitoringLevelRequest{Severity: 8}, true},
		{diagnostics.SetMonitoringLevelRequest{Severity: 9}, true},
		{diagnostics.SetMonitoringLevelRequest{}, true},
		{diagnostics.SetMonitoringLevelRequest{Severity: -1}, false},
		{diagnostics.SetMonitoringLevelRequest{Severity: 10}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV2TestSuite) TestSetMonitoringLevelConfirmationValidation() {
	var confirmationTable = []GenericTestEntry{
		{diagnostics.SetMonitoringLevelResponse{Status: types.GenericDeviceModelStatusAccepted, StatusInfo: types.NewStatusInfo("200", "")}, true},
		{diagnostics.SetMonitoringLevelResponse{Status: types.GenericDeviceModelStatusAccepted}, true},
		{diagnostics.SetMonitoringLevelResponse{Status: "invalidDeviceModelStatus"}, false},
		{diagnostics.SetMonitoringLevelResponse{}, false},
	}
	ExecuteGenericTestTable(suite, confirmationTable)
}

func (suite *OcppV2TestSuite) TestSetMonitoringLevelE2EMocked() {
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	severity := 3
	status := types.GenericDeviceModelStatusAccepted
	statusInfo := types.NewStatusInfo("200", "")
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"severity":%v}]`,
		messageId, diagnostics.SetMonitoringLevelFeatureName, severity)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v","statusInfo":{"reasonCode":"%v"}}]`,
		messageId, status, statusInfo.ReasonCode)
	setMonitoringLevelResponse := diagnostics.NewSetMonitoringLevelResponse(status)
	setMonitoringLevelResponse.StatusInfo = statusInfo
	channel := NewMockWebSocket(wsId)

	handler := &MockChargingStationDiagnosticsHandler{}
	handler.On("OnSetMonitoringLevel", mock.Anything).Return(setMonitoringLevelResponse, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*diagnostics.SetMonitoringLevelRequest)
		suite.Require().True(ok)
		suite.Require().NotNil(request)
		suite.Equal(severity, request.Severity)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	suite.Require().Nil(err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.SetMonitoringLevel(wsId, func(response *diagnostics.SetMonitoringLevelResponse, err error) {
		suite.Require().Nil(err)
		suite.Require().NotNil(response)
		suite.Equal(status, response.Status)
		suite.Equal(statusInfo.ReasonCode, response.StatusInfo.ReasonCode)
		suite.Equal(statusInfo.AdditionalInfo, response.StatusInfo.AdditionalInfo)
		resultChannel <- true
	}, severity)
	suite.Require().Nil(err)
	result := <-resultChannel
	suite.True(result)
}

func (suite *OcppV2TestSuite) TestSetMonitoringLevelInvalidEndpoint() {
	messageId := defaultMessageId
	severity := 3
	request := diagnostics.NewSetMonitoringLevelRequest(severity)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"severity":%v}]`,
		messageId, diagnostics.SetMonitoringLevelFeatureName, severity)
	testUnsupportedRequestFromChargingStation(suite, request, requestJson, messageId)
}
