package ocpp2_test

import (
	"fmt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0/provisioning"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
)

// Test
func (suite *OcppV2TestSuite) TestGetBaseReportRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{provisioning.GetBaseReportRequest{RequestID: 42, ReportBase: provisioning.ReportTypeConfigurationInventory}, true},
		{provisioning.GetBaseReportRequest{ReportBase: provisioning.ReportTypeConfigurationInventory}, true},
		{provisioning.GetBaseReportRequest{RequestID: 42}, false},
		{provisioning.GetBaseReportRequest{}, false},
		{provisioning.GetBaseReportRequest{RequestID: 42, ReportBase: "invalidReportType"}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestGetBaseReportConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{provisioning.GetBaseReportResponse{Status: types.GenericDeviceModelStatusAccepted}, true},
		{provisioning.GetBaseReportResponse{Status: "invalidDeviceModelStatus"}, false},
		{provisioning.GetBaseReportResponse{}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestGetBaseReportE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	requestID := 42
	reportBase := provisioning.ReportTypeConfigurationInventory
	status := types.GenericDeviceModelStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"requestId":%v,"reportBase":"%v"}]`, messageId, provisioning.GetBaseReportFeatureName, requestID, reportBase)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	getBaseReportConfirmation := provisioning.NewGetBaseReportResponse(status)
	channel := NewMockWebSocket(wsId)

	handler := MockChargingStationProvisioningHandler{}
	handler.On("OnGetBaseReport", mock.Anything).Return(getBaseReportConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*provisioning.GetBaseReportRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, requestID, request.RequestID)
		assert.Equal(t, reportBase, request.ReportBase)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.GetBaseReport(wsId, func(confirmation *provisioning.GetBaseReportResponse, err error) {
		require.Nil(t, err)
		require.NotNil(t, confirmation)
		assert.Equal(t, status, confirmation.Status)
		resultChannel <- true
	}, requestID, reportBase)
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV2TestSuite) TestGetBaseReportInvalidEndpoint() {
	messageId := defaultMessageId
	requestID := 42
	reportBase := provisioning.ReportTypeConfigurationInventory
	getBaseReportRequest := provisioning.NewGetBaseReportRequest(requestID, reportBase)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"requestId":%v,"reportBase":"%v"}]`, messageId, provisioning.GetBaseReportFeatureName, requestID, reportBase)
	testUnsupportedRequestFromChargingStation(suite, getBaseReportRequest, requestJson, messageId)
}
