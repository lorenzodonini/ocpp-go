package ocpp2_test

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0/diagnostics"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test
func (suite *OcppV2TestSuite) TestSetMonitoringBaseRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{diagnostics.SetMonitoringBaseRequest{MonitoringBase: diagnostics.MonitoringBaseAll}, true},
		{diagnostics.SetMonitoringBaseRequest{MonitoringBase: diagnostics.MonitoringBaseFactoryDefault}, true},
		{diagnostics.SetMonitoringBaseRequest{MonitoringBase: diagnostics.MonitoringBaseHardWiredOnly}, true},
		{diagnostics.SetMonitoringBaseRequest{MonitoringBase: "invalidMonitoringBase"}, false},
		{diagnostics.SetMonitoringBaseRequest{}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestSetMonitoringBaseConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{diagnostics.SetMonitoringBaseResponse{Status: types.GenericDeviceModelStatusAccepted, StatusInfo: types.NewStatusInfo("200", "")}, true},
		{diagnostics.SetMonitoringBaseResponse{Status: types.GenericDeviceModelStatusAccepted}, true},
		{diagnostics.SetMonitoringBaseResponse{Status: "invalidDeviceModelStatus"}, false},
		{diagnostics.SetMonitoringBaseResponse{}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestSetMonitoringBaseE2EMocked() {
	t := suite.T()
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

	handler := MockChargingStationDiagnosticsHandler{}
	handler.On("OnSetMonitoringBase", mock.Anything).Return(setMonitoringBaseResponse, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*diagnostics.SetMonitoringBaseRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, monitoringBase, request.MonitoringBase)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.SetMonitoringBase(wsId, func(response *diagnostics.SetMonitoringBaseResponse, err error) {
		require.Nil(t, err)
		require.NotNil(t, response)
		assert.Equal(t, status, response.Status)
		assert.Equal(t, statusInfo.ReasonCode, response.StatusInfo.ReasonCode)
		assert.Equal(t, statusInfo.AdditionalInfo, response.StatusInfo.AdditionalInfo)
		resultChannel <- true
	}, monitoringBase)
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV2TestSuite) TestSetMonitoringBaseInvalidEndpoint() {
	messageId := defaultMessageId
	monitoringBase := diagnostics.MonitoringBaseAll
	request := diagnostics.NewSetMonitoringBaseRequest(monitoringBase)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"monitoringBase":"%v"}]`,
		messageId, diagnostics.SetMonitoringBaseFeatureName, monitoringBase)
	testUnsupportedRequestFromChargingStation(suite, request, requestJson, messageId)
}
