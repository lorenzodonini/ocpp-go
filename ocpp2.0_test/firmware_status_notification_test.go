package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test
func (suite *OcppV2TestSuite) TestFirmwareStatusNotificationRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{ocpp2.FirmwareStatusNotificationRequest{Status: ocpp2.FirmwareStatusDownloaded, RequestID: 42}, true},
		{ocpp2.FirmwareStatusNotificationRequest{Status: ocpp2.FirmwareStatusDownloaded}, true},
		{ocpp2.FirmwareStatusNotificationRequest{RequestID: 42}, false},
		{ocpp2.FirmwareStatusNotificationRequest{}, false},
		{ocpp2.FirmwareStatusNotificationRequest{Status: ocpp2.FirmwareStatusDownloaded, RequestID: -1},  false},
		{ocpp2.FirmwareStatusNotificationRequest{Status: "invalidFirmwareStatus"}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestFirmwareStatusNotificationConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{ocpp2.FirmwareStatusNotificationConfirmation{}, true},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestFirmwareStatusNotificationE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	status := ocpp2.FirmwareStatusDownloaded
	requestID := 42
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"status":"%v","requestId":%v}]`, messageId, ocpp2.FirmwareStatusNotificationFeatureName, status, requestID)
	responseJson := fmt.Sprintf(`[3,"%v",{}]`, messageId)
	firmwareStatusNotificationConfirmation := ocpp2.NewFirmwareStatusNotificationConfirmation()
	channel := NewMockWebSocket(wsId)

	coreListener := MockCentralSystemCoreListener{}
	coreListener.On("OnFirmwareStatusNotification", mock.AnythingOfType("string"), mock.Anything).Return(firmwareStatusNotificationConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*ocpp2.FirmwareStatusNotificationRequest)
		assert.True(t, ok)
		assert.Equal(t, status, request.Status)
	})
	setupDefaultCentralSystemHandlers(suite, coreListener, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	assert.Nil(t, err)
	confirmation, err := suite.chargePoint.FirmwareStatusNotification(status, requestID)
	assert.Nil(t, err)
	assert.NotNil(t, confirmation)
}

func (suite *OcppV2TestSuite) TestFirmwareStatusNotificationInvalidEndpoint() {
	messageId := defaultMessageId
	status := ocpp2.FirmwareStatusDownloaded
	requestID := 42
	firmwareStatusRequest := ocpp2.NewFirmwareStatusNotificationRequest(status, requestID)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"status":"%v","requestId":%v}]`, messageId, ocpp2.FirmwareStatusNotificationFeatureName, status, requestID)
	testUnsupportedRequestFromCentralSystem(suite, firmwareStatusRequest, requestJson, messageId)
}
