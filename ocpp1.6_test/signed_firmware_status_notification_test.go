package ocpp16_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/firmware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test
func (suite *OcppV16TestSuite) TestSignedFirmwareStatusNotificationRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{firmware.SignedFirmwareStatusNotificationRequest{Status: firmware.SignedFirmwareStatusDownloaded, RequestId: 10}, true},
		{firmware.SignedFirmwareStatusNotificationRequest{}, false},
		{firmware.SignedFirmwareStatusNotificationRequest{Status: "invalidFirmwareStatus"}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV16TestSuite) TestSignedFirmwareStatusNotificationConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{firmware.SignedFirmwareStatusNotificationConfirmation{}, true},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV16TestSuite) TestSignedFirmwareStatusNotificationE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	requestId := 10
	status := firmware.SignedFirmwareStatusDownloaded
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"status":"%v","requestId":%d}]`, messageId, firmware.SignedFirmwareStatusNotificationFeatureName, status, requestId)
	responseJson := fmt.Sprintf(`[3,"%v",{}]`, messageId)
	signedFirmwareStatusNotificationConfirmation := firmware.NewSignedFirmwareStatusNotificationConfirmation()
	channel := NewMockWebSocket(wsId)

	firmwareListener := MockCentralSystemFirmwareManagementListener{}
	firmwareListener.On("OnSignedFirmwareStatusNotification", mock.AnythingOfType("string"), mock.Anything).Return(signedFirmwareStatusNotificationConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*firmware.SignedFirmwareStatusNotificationRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, status, request.Status)
	})
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	suite.centralSystem.SetFirmwareManagementHandler(firmwareListener)
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	require.Nil(t, err)
	confirmation, err := suite.chargePoint.SignedFirmwareStatusNotification(status, requestId)
	require.Nil(t, err)
	require.NotNil(t, confirmation)
}

func (suite *OcppV16TestSuite) TestSignedFirmwareStatusNotificationInvalidEndpoint() {
	messageId := defaultMessageId
	status := firmware.SignedFirmwareStatusDownloaded
	requestId := 10
	signedFirmwareStatusRequest := firmware.NewSignedFirmwareStatusNotificationRequest(status, requestId)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"status":"%v","requestId":%d}]`, messageId, firmware.SignedFirmwareStatusNotificationFeatureName, status, requestId)
	testUnsupportedRequestFromCentralSystem(suite, signedFirmwareStatusRequest, requestJson, messageId)
}
