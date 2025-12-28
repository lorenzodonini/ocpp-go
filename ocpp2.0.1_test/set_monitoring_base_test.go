package ocpp2_test

import (
	"fmt"

	"github.com/stretchr/testify/mock"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/diagnostics"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

// Test
func (suite *OcppV2TestSuite) TestSetMonitoringBaseRequestValidation() {
	var requestTable = []GenericTestEntry{
		{diagnostics.SetMonitoringBaseRequest{MonitoringBase: diagnostics.MonitoringBaseAll}, true},
		{diagnostics.SetMonitoringBaseRequest{MonitoringBase: diagnostics.MonitoringBaseFactoryDefault}, true},
		{diagnostics.SetMonitoringBaseRequest{MonitoringBase: diagnostics.MonitoringBaseHardWiredOnly}, true},
		{diagnostics.SetMonitoringBaseRequest{MonitoringBase: "invalidMonitoringBase"}, false},
		{diagnostics.SetMonitoringBaseRequest{}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV2TestSuite) TestSetMonitoringBaseConfirmationValidation() {
	var confirmationTable = []GenericTestEntry{
		{diagnostics.SetMonitoringBaseResponse{Status: types.GenericDeviceModelStatusAccepted, StatusInfo: types.NewStatusInfo("200", "")}, true},
		{diagnostics.SetMonitoringBaseResponse{Status: types.GenericDeviceModelStatusAccepted}, true},
		{diagnostics.SetMonitoringBaseResponse{Status: "invalidDeviceModelStatus"}, false},
		{diagnostics.SetMonitoringBaseResponse{}, false},
	}
	ExecuteGenericTestTable(suite, confirmationTable)
}

func (suite *OcppV2TestSuite) TestSetMonitoringBaseE2EMocked() {
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	monitoringBase := diagnostics.MonitoringBaseAll
	status := types.GenericDeviceModelStatusAccepted
	statusInfo := types.NewStatusInfo("200", "")
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"monitoringBase":"%v"}]`,
		messageId, diagnostics.SetMonitoringBaseFeatureName, monitoringBase)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v","statusInfo":{"reasonCode":"%v"}}]`,
		messageId, status, statusInfo.ReasonCode)
	setMonitoringBaseResponse := diagnostics.NewSetMonitoringBaseResponse(status)
	setMonitoringBaseResponse.StatusInfo = statusInfo
	channel := NewMockWebSocket(wsId)

	handler := &MockChargingStationDiagnosticsHandler{}
	handler.On("OnSetMonitoringBase", mock.Anything).Return(setMonitoringBaseResponse, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*diagnostics.SetMonitoringBaseRequest)
		suite.Require().True(ok)
		suite.Require().NotNil(request)
		suite.Equal(monitoringBase, request.MonitoringBase)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	suite.Require().Nil(err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.SetMonitoringBase(wsId, func(response *diagnostics.SetMonitoringBaseResponse, err error) {
		suite.Require().Nil(err)
		suite.Require().NotNil(response)
		suite.Equal(status, response.Status)
		suite.Equal(statusInfo.ReasonCode, response.StatusInfo.ReasonCode)
		suite.Equal(statusInfo.AdditionalInfo, response.StatusInfo.AdditionalInfo)
		resultChannel <- true
	}, monitoringBase)
	suite.Require().Nil(err)
	result := <-resultChannel
	suite.True(result)
}

func (suite *OcppV2TestSuite) TestSetMonitoringBaseInvalidEndpoint() {
	messageId := defaultMessageId
	monitoringBase := diagnostics.MonitoringBaseAll
	request := diagnostics.NewSetMonitoringBaseRequest(monitoringBase)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"monitoringBase":"%v"}]`,
		messageId, diagnostics.SetMonitoringBaseFeatureName, monitoringBase)
	testUnsupportedRequestFromChargingStation(suite, request, requestJson, messageId)
}
