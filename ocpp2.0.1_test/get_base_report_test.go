package ocpp2_test

import (
	"fmt"

	"github.com/stretchr/testify/mock"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/provisioning"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

// Test
func (suite *OcppV2TestSuite) TestGetBaseReportRequestValidation() {
	var requestTable = []GenericTestEntry{
		{provisioning.GetBaseReportRequest{RequestID: 42, ReportBase: provisioning.ReportTypeConfigurationInventory}, true},
		{provisioning.GetBaseReportRequest{ReportBase: provisioning.ReportTypeConfigurationInventory}, true},
		{provisioning.GetBaseReportRequest{RequestID: 42}, false},
		{provisioning.GetBaseReportRequest{}, false},
		{provisioning.GetBaseReportRequest{RequestID: 42, ReportBase: "invalidReportType"}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV2TestSuite) TestGetBaseReportConfirmationValidation() {
	var confirmationTable = []GenericTestEntry{
		{provisioning.GetBaseReportResponse{Status: types.GenericDeviceModelStatusAccepted}, true},
		{provisioning.GetBaseReportResponse{Status: "invalidDeviceModelStatus"}, false},
		{provisioning.GetBaseReportResponse{}, false},
	}
	ExecuteGenericTestTable(suite, confirmationTable)
}

func (suite *OcppV2TestSuite) TestGetBaseReportE2EMocked() {
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

	handler := &MockChargingStationProvisioningHandler{}
	handler.On("OnGetBaseReport", mock.Anything).Return(getBaseReportConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*provisioning.GetBaseReportRequest)
		suite.Require().True(ok)
		suite.Require().NotNil(request)
		suite.Equal(requestID, request.RequestID)
		suite.Equal(reportBase, request.ReportBase)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	suite.Require().Nil(err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.GetBaseReport(wsId, func(confirmation *provisioning.GetBaseReportResponse, err error) {
		suite.Require().Nil(err)
		suite.Require().NotNil(confirmation)
		suite.Equal(status, confirmation.Status)
		resultChannel <- true
	}, requestID, reportBase)
	suite.Require().Nil(err)
	result := <-resultChannel
	suite.True(result)
}

func (suite *OcppV2TestSuite) TestGetBaseReportInvalidEndpoint() {
	messageId := defaultMessageId
	requestID := 42
	reportBase := provisioning.ReportTypeConfigurationInventory
	getBaseReportRequest := provisioning.NewGetBaseReportRequest(requestID, reportBase)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"requestId":%v,"reportBase":"%v"}]`, messageId, provisioning.GetBaseReportFeatureName, requestID, reportBase)
	testUnsupportedRequestFromChargingStation(suite, getBaseReportRequest, requestJson, messageId)
}
