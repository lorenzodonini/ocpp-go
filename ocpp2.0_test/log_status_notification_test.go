package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/diagnostics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test
func (suite *OcppV2TestSuite) TestLogStatusNotificationRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{diagnostics.LogStatusNotificationRequest{Status: diagnostics.UploadLogStatusUploading, RequestID: 42}, true},
		{diagnostics.LogStatusNotificationRequest{Status: diagnostics.UploadLogStatusUploadFailure, RequestID: 42}, true},
		{diagnostics.LogStatusNotificationRequest{Status: diagnostics.UploadLogStatusUploaded, RequestID: 42}, true},
		{diagnostics.LogStatusNotificationRequest{Status: diagnostics.UploadLogStatusPermissionDenied, RequestID: 42}, true},
		{diagnostics.LogStatusNotificationRequest{Status: diagnostics.UploadLogStatusNotSupportedOp, RequestID: 42}, true},
		{diagnostics.LogStatusNotificationRequest{Status: diagnostics.UploadLogStatusIdle, RequestID: 42}, true},
		{diagnostics.LogStatusNotificationRequest{Status: diagnostics.UploadLogStatusBadMessage, RequestID: 42}, true},
		{diagnostics.LogStatusNotificationRequest{Status: diagnostics.UploadLogStatusIdle}, true},
		{diagnostics.LogStatusNotificationRequest{RequestID: 42}, false},
		{diagnostics.LogStatusNotificationRequest{}, false},
		{diagnostics.LogStatusNotificationRequest{Status: diagnostics.UploadLogStatusIdle, RequestID: -1}, false},
		{diagnostics.LogStatusNotificationRequest{Status: "invalidUploadLogStatus", RequestID: 42}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestLogStatusNotificationConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{diagnostics.LogStatusNotificationResponse{}, true},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestLogStatusNotificationE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	status := diagnostics.UploadLogStatusIdle
	requestID := 42
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"status":"%v","requestId":%v}]`, messageId, diagnostics.LogStatusNotificationFeatureName, status, requestID)
	responseJson := fmt.Sprintf(`[3,"%v",{}]`, messageId)
	logStatusNotificationResponse := diagnostics.NewLogStatusNotificationResponse()
	channel := NewMockWebSocket(wsId)

	handler := MockCSMSDiagnosticsHandler{}
	handler.On("OnLogStatusNotification", mock.AnythingOfType("string"), mock.Anything).Return(logStatusNotificationResponse, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*diagnostics.LogStatusNotificationRequest)
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
	confirmation, err := suite.chargingStation.LogStatusNotification(status, requestID)
	assert.Nil(t, err)
	assert.NotNil(t, confirmation)
}

func (suite *OcppV2TestSuite) TestLogStatusNotificationInvalidEndpoint() {
	messageId := defaultMessageId
	status := diagnostics.UploadLogStatusIdle
	requestID := 42
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"status":"%v","requestId":%v}]`, messageId, diagnostics.LogStatusNotificationFeatureName, status, requestID)
	req := diagnostics.NewLogStatusNotificationRequest(status, requestID)
	testUnsupportedRequestFromCentralSystem(suite, req, requestJson, messageId)
}
