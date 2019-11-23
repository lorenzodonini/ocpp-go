package ocpp16_test

import (
	"fmt"
	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"time"
)

// Test
func (suite *OcppV16TestSuite) TestGetDiagnosticsRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{ocpp16.GetDiagnosticsRequest{Location: "ftp:some/path", Retries: 10, RetryInterval: 10, StartTime: ocpp16.NewDateTime(time.Now()), EndTime: ocpp16.NewDateTime(time.Now())}, true},
		{ocpp16.GetDiagnosticsRequest{Location: "ftp:some/path", Retries: 10, RetryInterval: 10, StartTime: ocpp16.NewDateTime(time.Now())}, true},
		{ocpp16.GetDiagnosticsRequest{Location: "ftp:some/path", Retries: 10, RetryInterval: 10}, true},
		{ocpp16.GetDiagnosticsRequest{Location: "ftp:some/path", Retries: 10}, true},
		{ocpp16.GetDiagnosticsRequest{Location: "ftp:some/path"}, true},
		{ocpp16.GetDiagnosticsRequest{}, false},
		{ocpp16.GetDiagnosticsRequest{Location: "invalidUri"}, false},
		{ocpp16.GetDiagnosticsRequest{Location: "ftp:some/path", Retries: -1}, false},
		{ocpp16.GetDiagnosticsRequest{Location: "ftp:some/path", RetryInterval: -1}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV16TestSuite) TestGetDiagnosticsConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{ocpp16.GetDiagnosticsConfirmation{FileName: "someFileName"}, true},
		{ocpp16.GetDiagnosticsConfirmation{FileName: ""}, true},
		{ocpp16.GetDiagnosticsConfirmation{}, true},
		{ocpp16.GetDiagnosticsConfirmation{FileName: ">255............................................................................................................................................................................................................................................................"}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV16TestSuite) TestGetDiagnosticsE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	location := "ftp:some/path"
	fileName := "diagnostics.json"
	retries := 10
	retryInterval := 600
	startTime := ocpp16.NewDateTime(time.Now().Add(-10 * time.Hour * 24))
	endTime := ocpp16.NewDateTime(time.Now())
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"location":"%v","retries":%v,"retryInterval":%v,"startTime":"%v","endTime":"%v"}]`,
		messageId, ocpp16.GetDiagnosticsFeatureName, location, retries, retryInterval, startTime.Format(ocpp16.ISO8601), endTime.Format(ocpp16.ISO8601))
	responseJson := fmt.Sprintf(`[3,"%v",{"fileName":"%v"}]`, messageId, fileName)
	getDiagnosticsConfirmation := ocpp16.NewGetDiagnosticsConfirmation()
	getDiagnosticsConfirmation.FileName = fileName
	channel := NewMockWebSocket(wsId)

	firmwareListener := MockChargePointFirmwareManagementListener{}
	firmwareListener.On("OnGetDiagnostics", mock.Anything).Return(getDiagnosticsConfirmation, nil)
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	suite.chargePoint.SetFirmwareManagementListener(firmwareListener)
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	assert.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.centralSystem.GetDiagnostics(wsId, func(confirmation *ocpp16.GetDiagnosticsConfirmation, err error) {
		assert.Nil(t, err)
		assert.NotNil(t, confirmation)
		if confirmation != nil {
			assert.Equal(t, fileName, confirmation.FileName)
			resultChannel <- true
		} else {
			resultChannel <- false
		}
	}, location, func(request *ocpp16.GetDiagnosticsRequest) {
		request.RetryInterval = retryInterval
		request.Retries = retries
		request.StartTime = startTime
		request.EndTime = endTime
	})
	assert.Nil(t, err)
	if err == nil {
		result := <-resultChannel
		assert.True(t, result)
	}
}

func (suite *OcppV16TestSuite) TestGetDiagnosticsInvalidEndpoint() {
	messageId := defaultMessageId
	location := "ftp:some/path"
	retries := 10
	retryInterval := 600
	startTime := ocpp16.NewDateTime(time.Now().Add(-10 * time.Hour * 24))
	endTime := ocpp16.NewDateTime(time.Now())
	localListVersionRequest := ocpp16.NewGetDiagnosticsRequest(location)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"location":"%v","retries":%v,"retryInterval":%v,"startTime":"%v","endTime":"%v"}]`,
		messageId, ocpp16.GetDiagnosticsFeatureName, location, retries, retryInterval, startTime.Format(ocpp16.ISO8601), endTime.Format(ocpp16.ISO8601))
	testUnsupportedRequestFromChargePoint(suite, localListVersionRequest, requestJson, messageId)
}
