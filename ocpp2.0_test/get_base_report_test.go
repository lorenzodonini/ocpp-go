package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test
func (suite *OcppV2TestSuite) TestGetBaseReportRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{ocpp2.GetBaseReportRequest{RequestID: 42, ReportBase: ocpp2.ReportTypeConfigurationInventory}, true},
		{ocpp2.GetBaseReportRequest{ReportBase: ocpp2.ReportTypeConfigurationInventory}, true},
		{ocpp2.GetBaseReportRequest{RequestID: 42}, false},
		{ocpp2.GetBaseReportRequest{}, false},
		{ocpp2.GetBaseReportRequest{RequestID: -1, ReportBase: ocpp2.ReportTypeConfigurationInventory}, false},
		{ocpp2.GetBaseReportRequest{RequestID: 42, ReportBase: "invalidReportType"}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestGetBaseReportConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{ocpp2.GetBaseReportConfirmation{Status: ocpp2.GenericDeviceModelStatusAccepted}, true},
		{ocpp2.GetBaseReportConfirmation{Status: "invalidDeviceModelStatus"}, false},
		{ocpp2.GetBaseReportConfirmation{}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestGetBaseReportE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	requestID := 42
	reportBase := ocpp2.ReportTypeConfigurationInventory
	status := ocpp2.GenericDeviceModelStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"requestId":%v,"reportBase":"%v"}]`, messageId, ocpp2.GetBaseReportFeatureName, requestID, reportBase)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	getBaseReportConfirmation := ocpp2.NewGetBaseReportConfirmation(status)
	channel := NewMockWebSocket(wsId)

	coreListener := MockChargePointCoreListener{}
	coreListener.On("OnGetBaseReport", mock.Anything).Return(getBaseReportConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*ocpp2.GetBaseReportRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, requestID, request.RequestID)
		assert.Equal(t, reportBase, request.ReportBase)
	})
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, coreListener, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.GetBaseReport(wsId, func(confirmation *ocpp2.GetBaseReportConfirmation, err error) {
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
	reportBase := ocpp2.ReportTypeConfigurationInventory
	getBaseReportRequest := ocpp2.NewGetBaseReportRequest(requestID, reportBase)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"requestId":%v,"reportBase":"%v"}]`, messageId, ocpp2.GetBaseReportFeatureName, requestID, reportBase)
	testUnsupportedRequestFromChargePoint(suite, getBaseReportRequest, requestJson, messageId)
}
