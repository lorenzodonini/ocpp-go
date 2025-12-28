package ocpp16_test

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/firmware"
	"github.com/stretchr/testify/mock"
)

// Test
func (suite *OcppV16TestSuite) TestDiagnosticsStatusNotificationRequestValidation() {
	requestTable := []GenericTestEntry{
		{firmware.DiagnosticsStatusNotificationRequest{Status: firmware.DiagnosticsStatusUploaded}, true},
		{firmware.DiagnosticsStatusNotificationRequest{}, false},
		{firmware.DiagnosticsStatusNotificationRequest{Status: "invalidDiagnosticsStatus"}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV16TestSuite) TestDiagnosticsStatusNotificationConfirmationValidation() {
	confirmationTable := []GenericTestEntry{
		{firmware.DiagnosticsStatusNotificationConfirmation{}, true},
	}
	ExecuteGenericTestTable(suite, confirmationTable)
}

func (suite *OcppV16TestSuite) TestDiagnosticsStatusNotificationE2EMocked() {
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	status := firmware.DiagnosticsStatusUploaded
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"status":"%v"}]`, messageId, firmware.DiagnosticsStatusNotificationFeatureName, status)
	responseJson := fmt.Sprintf(`[3,"%v",{}]`, messageId)
	diagnosticsStatusNotificationConfirmation := firmware.NewDiagnosticsStatusNotificationConfirmation()
	channel := NewMockWebSocket(wsId)

	firmwareListener := &MockCentralSystemFirmwareManagementListener{}
	firmwareListener.On("OnDiagnosticsStatusNotification", mock.AnythingOfType("string"), mock.Anything).Return(diagnosticsStatusNotificationConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*firmware.DiagnosticsStatusNotificationRequest)
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
	confirmation, err := suite.chargePoint.DiagnosticsStatusNotification(status)
	suite.Require().Nil(err)
	suite.Require().NotNil(confirmation)
}

func (suite *OcppV16TestSuite) TestDiagnosticsStatusNotificationInvalidEndpoint() {
	messageId := defaultMessageId
	status := firmware.DiagnosticsStatusUploaded
	diagnosticsStatusRequest := firmware.NewDiagnosticsStatusNotificationRequest(status)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"status":"%v"}]`, messageId, firmware.DiagnosticsStatusNotificationFeatureName, status)
	testUnsupportedRequestFromCentralSystem(suite, diagnosticsStatusRequest, requestJson, messageId)
}
