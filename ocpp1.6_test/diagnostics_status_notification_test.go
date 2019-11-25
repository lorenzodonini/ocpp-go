package ocpp16_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test
func (suite *OcppV16TestSuite) TestDiagnosticsStatusNotificationRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{ocpp16.DiagnosticsStatusNotificationRequest{Status: ocpp16.DiagnosticsStatusUploaded}, true},
		{ocpp16.DiagnosticsStatusNotificationRequest{}, false},
		{ocpp16.DiagnosticsStatusNotificationRequest{Status: "invalidDiagnosticsStatus"}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV16TestSuite) TestDiagnosticsStatusNotificationConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{ocpp16.DiagnosticsStatusNotificationConfirmation{}, true},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV16TestSuite) TestDiagnosticsStatusNotificationE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	status := ocpp16.DiagnosticsStatusUploaded
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"status":"%v"}]`, messageId, ocpp16.DiagnosticsStatusNotificationFeatureName, status)
	responseJson := fmt.Sprintf(`[3,"%v",{}]`, messageId)
	diagnosticsStatusNotificationConfirmation := ocpp16.NewDiagnosticsStatusNotificationConfirmation()
	channel := NewMockWebSocket(wsId)

	firmwareListener := MockCentralSystemFirmwareManagementListener{}
	firmwareListener.On("OnDiagnosticsStatusNotification", mock.AnythingOfType("string"), mock.Anything).Return(diagnosticsStatusNotificationConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*ocpp16.DiagnosticsStatusNotificationRequest)
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
	confirmation, err := suite.chargePoint.DiagnosticsStatusNotification(status)
	assert.Nil(t, err)
	assert.NotNil(t, confirmation)
}

func (suite *OcppV16TestSuite) TestDiagnosticsStatusNotificationInvalidEndpoint() {
	messageId := defaultMessageId
	status := ocpp16.DiagnosticsStatusUploaded
	diagnosticsStatusRequest := ocpp16.NewDiagnosticsStatusNotificationRequest(status)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"status":"%v"}]`, messageId, ocpp16.DiagnosticsStatusNotificationFeatureName, status)
	testUnsupportedRequestFromCentralSystem(suite, diagnosticsStatusRequest, requestJson, messageId)
}
