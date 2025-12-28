package ocpp2_test

import (
	"fmt"

	"github.com/stretchr/testify/mock"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/firmware"
)

// Test
func (suite *OcppV2TestSuite) TestFirmwareStatusNotificationRequestValidation() {
	var requestTable = []GenericTestEntry{
		{firmware.FirmwareStatusNotificationRequest{Status: firmware.FirmwareStatusDownloaded, RequestID: newInt(42)}, true},
		{firmware.FirmwareStatusNotificationRequest{Status: firmware.FirmwareStatusDownloaded}, true},
		{firmware.FirmwareStatusNotificationRequest{RequestID: newInt(42)}, false},
		{firmware.FirmwareStatusNotificationRequest{}, false},
		{firmware.FirmwareStatusNotificationRequest{Status: firmware.FirmwareStatusDownloaded, RequestID: newInt(-1)}, false},
		{firmware.FirmwareStatusNotificationRequest{Status: "invalidFirmwareStatus"}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV2TestSuite) TestFirmwareStatusNotificationConfirmationValidation() {
	var confirmationTable = []GenericTestEntry{
		{firmware.FirmwareStatusNotificationResponse{}, true},
	}
	ExecuteGenericTestTable(suite, confirmationTable)
}

func (suite *OcppV2TestSuite) TestFirmwareStatusNotificationE2EMocked() {
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	status := firmware.FirmwareStatusDownloaded
	requestID := 42
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"status":"%v","requestId":%v}]`, messageId, firmware.FirmwareStatusNotificationFeatureName, status, requestID)
	responseJson := fmt.Sprintf(`[3,"%v",{}]`, messageId)
	firmwareStatusNotificationConfirmation := firmware.NewFirmwareStatusNotificationResponse()
	channel := NewMockWebSocket(wsId)

	handler := &MockCSMSFirmwareHandler{}
	handler.On("OnFirmwareStatusNotification", mock.AnythingOfType("string"), mock.Anything).Return(firmwareStatusNotificationConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*firmware.FirmwareStatusNotificationRequest)
		suite.Require().True(ok)
		suite.Equal(status, request.Status)
		suite.Equal(requestID, *request.RequestID)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	suite.Require().Nil(err)
	response, err := suite.chargingStation.FirmwareStatusNotification(status, func(request *firmware.FirmwareStatusNotificationRequest) {
		request.RequestID = &requestID
	})
	suite.Nil(err)
	suite.NotNil(response)
}

func (suite *OcppV2TestSuite) TestFirmwareStatusNotificationInvalidEndpoint() {
	messageId := defaultMessageId
	status := firmware.FirmwareStatusDownloaded
	requestID := 42
	firmwareStatusRequest := firmware.NewFirmwareStatusNotificationRequest(status)
	firmwareStatusRequest.RequestID = &requestID
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"status":"%v","requestId":%v}]`, messageId, firmware.FirmwareStatusNotificationFeatureName, status, requestID)
	testUnsupportedRequestFromCentralSystem(suite, firmwareStatusRequest, requestJson, messageId)
}
