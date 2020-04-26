package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"time"
)

// Test
func (suite *OcppV2TestSuite) TestGetLogRequestValidation() {
	t := suite.T()
	logParameters := ocpp2.LogParameters{
		RemoteLocation:  "ftp://someurl/diagnostics/1",
		OldestTimestamp: ocpp2.NewDateTime(time.Now().Add(-2 * time.Hour)),
		LatestTimestamp: ocpp2.NewDateTime(time.Now()),
	}
	var requestTable = []GenericTestEntry{
		{ocpp2.GetLogRequest{LogType: ocpp2.LogTypeDiagnostics, RequestID: 1, Retries: newInt(5), RetryInterval: newInt(120), Log: logParameters}, true},
		{ocpp2.GetLogRequest{LogType: ocpp2.LogTypeDiagnostics, RequestID: 1, Retries: newInt(5), Log: logParameters}, true},
		{ocpp2.GetLogRequest{LogType: ocpp2.LogTypeDiagnostics, RequestID: 1, Log: logParameters}, true},
		{ocpp2.GetLogRequest{LogType: ocpp2.LogTypeDiagnostics, Log: logParameters}, true},
		{ocpp2.GetLogRequest{LogType: ocpp2.LogTypeDiagnostics}, false},
		{ocpp2.GetLogRequest{Log: logParameters}, false},
		{ocpp2.GetLogRequest{}, false},
		{ocpp2.GetLogRequest{LogType: "invalidLogType", RequestID: 1, Retries: newInt(5), RetryInterval: newInt(120), Log: logParameters}, false},
		{ocpp2.GetLogRequest{LogType: ocpp2.LogTypeDiagnostics, RequestID: -1, Retries: newInt(5), RetryInterval: newInt(120), Log: logParameters}, false},
		{ocpp2.GetLogRequest{LogType: ocpp2.LogTypeDiagnostics, RequestID: 1, Retries: newInt(-1), RetryInterval: newInt(120), Log: logParameters}, false},
		{ocpp2.GetLogRequest{LogType: ocpp2.LogTypeDiagnostics, RequestID: 1, Retries: newInt(5), RetryInterval: newInt(-1), Log: logParameters}, false},
		{ocpp2.GetLogRequest{LogType: ocpp2.LogTypeDiagnostics, RequestID: 1, Retries: newInt(5), RetryInterval: newInt(120), Log: ocpp2.LogParameters{RemoteLocation:  ".invalidUrl.", OldestTimestamp: nil, LatestTimestamp: nil}}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestGetLogConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{ocpp2.GetLogConfirmation{Status: ocpp2.LogStatusAccepted, Filename: "testFileName.log"}, true},
		{ocpp2.GetLogConfirmation{Status: ocpp2.LogStatusAccepted}, true},
		{ocpp2.GetLogConfirmation{Status: ocpp2.LogStatusRejected}, true},
		{ocpp2.GetLogConfirmation{Status: ocpp2.LogStatusAcceptedCanceled}, true},
		{ocpp2.GetLogConfirmation{}, false},
		{ocpp2.GetLogConfirmation{Status: "invalidLogStatus"}, false},
		{ocpp2.GetLogConfirmation{Status: ocpp2.LogStatusAccepted, Filename: ">256............................................................................................................................................................................................................................................................."}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestGetLogE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	logParameters := ocpp2.LogParameters{
		RemoteLocation:  "ftp://someurl/diagnostics/1",
		OldestTimestamp: ocpp2.NewDateTime(time.Now().Add(-2 * time.Hour)),
		LatestTimestamp: ocpp2.NewDateTime(time.Now()),
	}
	logType := ocpp2.LogTypeDiagnostics
	requestID := 42
	retries := newInt(5)
	retryInterval := newInt(120)
	status := ocpp2.LogStatusAccepted
	filename := "someFileName.log"
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"logType":"%v","requestId":%v,"retries":%v,"retryInterval":%v,"log":{"remoteLocation":"%v","oldestTimestamp":"%v","latestTimestamp":"%v"}}]`,
		messageId, ocpp2.GetLogFeatureName, logType, requestID, *retries, *retryInterval, logParameters.RemoteLocation, ocpp2.FormatTimestamp(logParameters.OldestTimestamp.Time), ocpp2.FormatTimestamp(logParameters.LatestTimestamp.Time))
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v","filename":"%v"}]`, messageId, status, filename)
	getLogConfirmation := ocpp2.NewGetLogConfirmation(status)
	getLogConfirmation.Filename = filename
	channel := NewMockWebSocket(wsId)

	coreListener := MockChargePointCoreListener{}
	coreListener.On("OnGetLog", mock.Anything).Return(getLogConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*ocpp2.GetLogRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, logType, request.LogType)
		assert.Equal(t, requestID, request.RequestID)
		assert.Equal(t, *retries, *request.Retries)
		assert.Equal(t, *retryInterval, *request.RetryInterval)
		assert.Equal(t, logParameters.RemoteLocation, request.Log.RemoteLocation)
		assert.Equal(t, ocpp2.FormatTimestamp(logParameters.LatestTimestamp.Time), ocpp2.FormatTimestamp(request.Log.LatestTimestamp.Time))
		assert.Equal(t, ocpp2.FormatTimestamp(logParameters.OldestTimestamp.Time), ocpp2.FormatTimestamp(request.Log.OldestTimestamp.Time))
	})
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, coreListener, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.GetLog(wsId, func(confirmation *ocpp2.GetLogConfirmation, err error) {
		require.Nil(t, err)
		require.NotNil(t, confirmation)
		assert.Equal(t, status, confirmation.Status)
		assert.Equal(t, filename, confirmation.Filename)
		resultChannel <- true
	}, logType, requestID, logParameters, func(request *ocpp2.GetLogRequest) {
		request.Retries = retries
		request.RetryInterval = retryInterval
	})
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV2TestSuite) TestGetLogInvalidEndpoint() {
	messageId := defaultMessageId
	logParameters := ocpp2.LogParameters{
		RemoteLocation:  "ftp://someurl/diagnostics/1",
		OldestTimestamp: ocpp2.NewDateTime(time.Now().Add(-2 * time.Hour)),
		LatestTimestamp: ocpp2.NewDateTime(time.Now()),
	}
	logType := ocpp2.LogTypeDiagnostics
	requestID := 42
	retries := newInt(5)
	retryInterval := newInt(120)
	getLogRequest := ocpp2.NewGetLogRequest(logType, requestID, logParameters)
	getLogRequest.Retries = retries
	getLogRequest.RetryInterval = retryInterval
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"logType":"%v","requestId":%v,"retries":%v,"retryInterval":%v,"log":{"remoteLocation":"%v","oldestTimestamp":"%v","latestTimestamp":"%v"}}]`,
		messageId, ocpp2.GetLogFeatureName, logType, requestID, *retries, *retryInterval, logParameters.RemoteLocation, ocpp2.FormatTimestamp(logParameters.OldestTimestamp.Time), ocpp2.FormatTimestamp(logParameters.LatestTimestamp.Time))
	testUnsupportedRequestFromChargePoint(suite, getLogRequest, requestJson, messageId)
}
