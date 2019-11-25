package ocpp16_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test
func (suite *OcppV16TestSuite) TestFirmwareStatusNotificationRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{ocpp16.FirmwareStatusNotificationRequest{Status: ocpp16.FirmwareStatusDownloaded}, true},
		{ocpp16.FirmwareStatusNotificationRequest{}, false},
		{ocpp16.FirmwareStatusNotificationRequest{Status: "invalidFirmwareStatus"}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV16TestSuite) TestFirmwareStatusNotificationConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{ocpp16.FirmwareStatusNotificationConfirmation{}, true},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV16TestSuite) TestFirmwareStatusNotificationE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	status := ocpp16.FirmwareStatusDownloaded
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"status":"%v"}]`, messageId, ocpp16.FirmwareStatusNotificationFeatureName, status)
	responseJson := fmt.Sprintf(`[3,"%v",{}]`, messageId)
	firmwareStatusNotificationConfirmation := ocpp16.NewFirmwareStatusNotificationConfirmation()
	channel := NewMockWebSocket(wsId)

	firmwareListener := MockCentralSystemFirmwareManagementListener{}
	firmwareListener.On("OnFirmwareStatusNotification", mock.AnythingOfType("string"), mock.Anything).Return(firmwareStatusNotificationConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*ocpp16.FirmwareStatusNotificationRequest)
		assert.True(t, ok)
		assert.Equal(t, status, request.Status)
	})
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	suite.centralSystem.SetFirmwareManagementListener(firmwareListener)
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	assert.Nil(t, err)
	confirmation, err := suite.chargePoint.FirmwareStatusNotification(status)
	assert.Nil(t, err)
	assert.NotNil(t, confirmation)
}

func (suite *OcppV16TestSuite) TestFirmwareStatusNotificationInvalidEndpoint() {
	messageId := defaultMessageId
	status := ocpp16.FirmwareStatusDownloaded
	firmwareStatusRequest := ocpp16.NewFirmwareStatusNotificationRequest(status)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"status":"%v"}]`, messageId, ocpp16.FirmwareStatusNotificationFeatureName, status)
	testUnsupportedRequestFromCentralSystem(suite, firmwareStatusRequest, requestJson, messageId)
}
