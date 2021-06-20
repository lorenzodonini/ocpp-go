package ocpp2_test

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/diagnostics"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

// Test
func (suite *OcppV2TestSuite) TestGetLogRequestValidation() {
	t := suite.T()
	logParameters := diagnostics.LogParameters{
		RemoteLocation:  "ftp://someurl/diagnostics/1",
		OldestTimestamp: types.NewDateTime(time.Now().Add(-2 * time.Hour)),
		LatestTimestamp: types.NewDateTime(time.Now()),
	}
	var requestTable = []GenericTestEntry{
		{diagnostics.GetLogRequest{LogType: diagnostics.LogTypeDiagnostics, RequestID: 1, Retries: newInt(5), RetryInterval: newInt(120), Log: logParameters}, true},
		{diagnostics.GetLogRequest{LogType: diagnostics.LogTypeDiagnostics, RequestID: 1, Retries: newInt(5), Log: logParameters}, true},
		{diagnostics.GetLogRequest{LogType: diagnostics.LogTypeDiagnostics, RequestID: 1, Log: logParameters}, true},
		{diagnostics.GetLogRequest{LogType: diagnostics.LogTypeDiagnostics, Log: logParameters}, true},
		{diagnostics.GetLogRequest{LogType: diagnostics.LogTypeDiagnostics}, false},
		{diagnostics.GetLogRequest{Log: logParameters}, false},
		{diagnostics.GetLogRequest{}, false},
		{diagnostics.GetLogRequest{LogType: "invalidLogType", RequestID: 1, Retries: newInt(5), RetryInterval: newInt(120), Log: logParameters}, false},
		{diagnostics.GetLogRequest{LogType: diagnostics.LogTypeDiagnostics, RequestID: -1, Retries: newInt(5), RetryInterval: newInt(120), Log: logParameters}, false},
		{diagnostics.GetLogRequest{LogType: diagnostics.LogTypeDiagnostics, RequestID: 1, Retries: newInt(-1), RetryInterval: newInt(120), Log: logParameters}, false},
		{diagnostics.GetLogRequest{LogType: diagnostics.LogTypeDiagnostics, RequestID: 1, Retries: newInt(5), RetryInterval: newInt(-1), Log: logParameters}, false},
		{diagnostics.GetLogRequest{LogType: diagnostics.LogTypeDiagnostics, RequestID: 1, Retries: newInt(5), RetryInterval: newInt(120), Log: diagnostics.LogParameters{RemoteLocation: ".invalidUrl.", OldestTimestamp: nil, LatestTimestamp: nil}}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestGetLogConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{diagnostics.GetLogResponse{Status: diagnostics.LogStatusAccepted, Filename: "testFileName.log"}, true},
		{diagnostics.GetLogResponse{Status: diagnostics.LogStatusAccepted}, true},
		{diagnostics.GetLogResponse{Status: diagnostics.LogStatusRejected}, true},
		{diagnostics.GetLogResponse{Status: diagnostics.LogStatusAcceptedCanceled}, true},
		{diagnostics.GetLogResponse{}, false},
		{diagnostics.GetLogResponse{Status: "invalidLogStatus"}, false},
		{diagnostics.GetLogResponse{Status: diagnostics.LogStatusAccepted, Filename: ">256............................................................................................................................................................................................................................................................."}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestGetLogE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	logParameters := diagnostics.LogParameters{
		RemoteLocation:  "ftp://someurl/diagnostics/1",
		OldestTimestamp: types.NewDateTime(time.Now().Add(-2 * time.Hour)),
		LatestTimestamp: types.NewDateTime(time.Now()),
	}
	logType := diagnostics.LogTypeDiagnostics
	requestID := 42
	retries := newInt(5)
	retryInterval := newInt(120)
	status := diagnostics.LogStatusAccepted
	filename := "someFileName.log"
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"logType":"%v","requestId":%v,"retries":%v,"retryInterval":%v,"log":{"remoteLocation":"%v","oldestTimestamp":"%v","latestTimestamp":"%v"}}]`,
		messageId, diagnostics.GetLogFeatureName, logType, requestID, *retries, *retryInterval, logParameters.RemoteLocation, logParameters.OldestTimestamp.FormatTimestamp(), logParameters.LatestTimestamp.FormatTimestamp())
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v","filename":"%v"}]`, messageId, status, filename)
	getLogConfirmation := diagnostics.NewGetLogResponse(status)
	getLogConfirmation.Filename = filename
	channel := NewMockWebSocket(wsId)

	handler := &MockChargingStationDiagnosticsHandler{}
	handler.On("OnGetLog", mock.Anything).Return(getLogConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*diagnostics.GetLogRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, logType, request.LogType)
		assert.Equal(t, requestID, request.RequestID)
		assert.Equal(t, *retries, *request.Retries)
		assert.Equal(t, *retryInterval, *request.RetryInterval)
		assert.Equal(t, logParameters.RemoteLocation, request.Log.RemoteLocation)
		assert.Equal(t, logParameters.LatestTimestamp.FormatTimestamp(), request.Log.LatestTimestamp.FormatTimestamp())
		assert.Equal(t, logParameters.OldestTimestamp.FormatTimestamp(), request.Log.OldestTimestamp.FormatTimestamp())
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.GetLog(wsId, func(confirmation *diagnostics.GetLogResponse, err error) {
		require.Nil(t, err)
		require.NotNil(t, confirmation)
		assert.Equal(t, status, confirmation.Status)
		assert.Equal(t, filename, confirmation.Filename)
		resultChannel <- true
	}, logType, requestID, logParameters, func(request *diagnostics.GetLogRequest) {
		request.Retries = retries
		request.RetryInterval = retryInterval
	})
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV2TestSuite) TestGetLogInvalidEndpoint() {
	messageId := defaultMessageId
	logParameters := diagnostics.LogParameters{
		RemoteLocation:  "ftp://someurl/diagnostics/1",
		OldestTimestamp: types.NewDateTime(time.Now().Add(-2 * time.Hour)),
		LatestTimestamp: types.NewDateTime(time.Now()),
	}
	logType := diagnostics.LogTypeDiagnostics
	requestID := 42
	retries := newInt(5)
	retryInterval := newInt(120)
	getLogRequest := diagnostics.NewGetLogRequest(logType, requestID, logParameters)
	getLogRequest.Retries = retries
	getLogRequest.RetryInterval = retryInterval
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"logType":"%v","requestId":%v,"retries":%v,"retryInterval":%v,"log":{"remoteLocation":"%v","oldestTimestamp":"%v","latestTimestamp":"%v"}}]`,
		messageId, diagnostics.GetLogFeatureName, logType, requestID, *retries, *retryInterval, logParameters.RemoteLocation, logParameters.OldestTimestamp.FormatTimestamp(), logParameters.LatestTimestamp.FormatTimestamp())
	testUnsupportedRequestFromChargingStation(suite, getLogRequest, requestJson, messageId)
}
