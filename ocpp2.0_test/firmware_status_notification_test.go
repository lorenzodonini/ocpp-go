package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/firmware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test
func (suite *OcppV2TestSuite) TestFirmwareStatusNotificationRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{firmware.FirmwareStatusNotificationRequest{Status: firmware.FirmwareStatusDownloaded, RequestID: 42}, true},
		{firmware.FirmwareStatusNotificationRequest{Status: firmware.FirmwareStatusDownloaded}, true},
		{firmware.FirmwareStatusNotificationRequest{RequestID: 42}, false},
		{firmware.FirmwareStatusNotificationRequest{}, false},
		{firmware.FirmwareStatusNotificationRequest{Status: firmware.FirmwareStatusDownloaded, RequestID: -1}, false},
		{firmware.FirmwareStatusNotificationRequest{Status: "invalidFirmwareStatus"}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestFirmwareStatusNotificationConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{firmware.FirmwareStatusNotificationResponse{}, true},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestFirmwareStatusNotificationE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	status := firmware.FirmwareStatusDownloaded
	requestID := 42
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"status":"%v","requestId":%v}]`, messageId, firmware.FirmwareStatusNotificationFeatureName, status, requestID)
	responseJson := fmt.Sprintf(`[3,"%v",{}]`, messageId)
	firmwareStatusNotificationConfirmation := firmware.NewFirmwareStatusNotificationResponse()
	channel := NewMockWebSocket(wsId)

	handler := MockCSMSFirmwareHandler{}
	handler.On("OnFirmwareStatusNotification", mock.AnythingOfType("string"), mock.Anything).Return(firmwareStatusNotificationConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*firmware.FirmwareStatusNotificationRequest)
		require.True(t, ok)
		assert.Equal(t, status, request.Status)
		assert.Equal(t, requestID, request.RequestID)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	confirmation, err := suite.chargingStation.FirmwareStatusNotification(status, requestID)
	assert.Nil(t, err)
	assert.NotNil(t, confirmation)
}

func (suite *OcppV2TestSuite) TestFirmwareStatusNotificationInvalidEndpoint() {
	messageId := defaultMessageId
	status := firmware.FirmwareStatusDownloaded
	requestID := 42
	firmwareStatusRequest := firmware.NewFirmwareStatusNotificationRequest(status, requestID)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"status":"%v","requestId":%v}]`, messageId, firmware.FirmwareStatusNotificationFeatureName, status, requestID)
	testUnsupportedRequestFromCentralSystem(suite, firmwareStatusRequest, requestJson, messageId)
}
