package ocpp16_test

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/firmware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test
func (suite *OcppV16TestSuite) TestDiagnosticsStatusNotificationRequestValidation() {
	t := suite.T()
	requestTable := []GenericTestEntry{
		{firmware.DiagnosticsStatusNotificationRequest{Status: firmware.DiagnosticsStatusUploaded}, true},
		{firmware.DiagnosticsStatusNotificationRequest{}, false},
		{firmware.DiagnosticsStatusNotificationRequest{Status: "invalidDiagnosticsStatus"}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV16TestSuite) TestDiagnosticsStatusNotificationConfirmationValidation() {
	t := suite.T()
	confirmationTable := []GenericTestEntry{
		{firmware.DiagnosticsStatusNotificationConfirmation{}, true},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV16TestSuite) TestDiagnosticsStatusNotificationE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	status := firmware.DiagnosticsStatusUploaded
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"status":"%v"}]`, messageId, firmware.DiagnosticsStatusNotificationFeatureName, status)
	responseJson := fmt.Sprintf(`[3,"%v",{}]`, messageId)
	diagnosticsStatusNotificationConfirmation := firmware.NewDiagnosticsStatusNotificationConfirmation()
	channel := NewMockWebSocket(wsId)

	firmwareListener := MockCentralSystemFirmwareManagementListener{}
	firmwareListener.On("OnDiagnosticsStatusNotification", mock.AnythingOfType("string"), mock.Anything).Return(diagnosticsStatusNotificationConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*firmware.DiagnosticsStatusNotificationRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, status, request.Status)
	})
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	suite.centralSystem.SetFirmwareManagementHandler(&firmwareListener)
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	require.Nil(t, err)
	confirmation, err := suite.chargePoint.DiagnosticsStatusNotification(status)
	require.Nil(t, err)
	require.NotNil(t, confirmation)
}

func (suite *OcppV16TestSuite) TestDiagnosticsStatusNotificationInvalidEndpoint() {
	messageId := defaultMessageId
	status := firmware.DiagnosticsStatusUploaded
	diagnosticsStatusRequest := firmware.NewDiagnosticsStatusNotificationRequest(status)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"status":"%v"}]`, messageId, firmware.DiagnosticsStatusNotificationFeatureName, status)
	testUnsupportedRequestFromCentralSystem(suite, diagnosticsStatusRequest, requestJson, messageId)
}
