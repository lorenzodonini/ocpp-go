package ocpp16_test

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/firmware"
	"github.com/stretchr/testify/mock"
)

// Test
func (suite *OcppV16TestSuite) TestFirmwareStatusNotificationRequestValidation() {
	requestTable := []GenericTestEntry{
		{firmware.FirmwareStatusNotificationRequest{Status: firmware.FirmwareStatusDownloaded}, true},
		{firmware.FirmwareStatusNotificationRequest{}, false},
		{firmware.FirmwareStatusNotificationRequest{Status: "invalidFirmwareStatus"}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV16TestSuite) TestFirmwareStatusNotificationConfirmationValidation() {
	confirmationTable := []GenericTestEntry{
		{firmware.FirmwareStatusNotificationConfirmation{}, true},
	}
	ExecuteGenericTestTable(suite, confirmationTable)
}

func (suite *OcppV16TestSuite) TestFirmwareStatusNotificationE2EMocked() {
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	status := firmware.FirmwareStatusDownloaded
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"status":"%v"}]`, messageId, firmware.FirmwareStatusNotificationFeatureName, status)
	responseJson := fmt.Sprintf(`[3,"%v",{}]`, messageId)
	firmwareStatusNotificationConfirmation := firmware.NewFirmwareStatusNotificationConfirmation()
	channel := NewMockWebSocket(wsId)

	firmwareListener := &MockCentralSystemFirmwareManagementListener{}
	firmwareListener.On("OnFirmwareStatusNotification", mock.AnythingOfType("string"), mock.Anything).Return(firmwareStatusNotificationConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*firmware.FirmwareStatusNotificationRequest)
		suite.Require().True(ok)
		suite.Require().NotNil(request)
		suite.Equal(status, request.Status)
	})
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	suite.centralSystem.SetFirmwareManagementHandler(firmwareListener)
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	suite.Require().Nil(err)
	confirmation, err := suite.chargePoint.FirmwareStatusNotification(status)
	suite.Require().Nil(err)
	suite.Require().NotNil(confirmation)
}

func (suite *OcppV16TestSuite) TestFirmwareStatusNotificationInvalidEndpoint() {
	messageId := defaultMessageId
	status := firmware.FirmwareStatusDownloaded
	firmwareStatusRequest := firmware.NewFirmwareStatusNotificationRequest(status)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"status":"%v"}]`, messageId, firmware.FirmwareStatusNotificationFeatureName, status)
	testUnsupportedRequestFromCentralSystem(suite, firmwareStatusRequest, requestJson, messageId)
}
