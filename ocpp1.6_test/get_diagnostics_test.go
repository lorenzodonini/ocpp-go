package ocpp16_test

import (
	"fmt"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/firmware"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/stretchr/testify/mock"
)

// Test
func (suite *OcppV16TestSuite) TestGetDiagnosticsRequestValidation() {
	requestTable := []GenericTestEntry{
		{firmware.GetDiagnosticsRequest{Location: "ftp:some/path", Retries: newInt(10), RetryInterval: newInt(10), StartTime: types.NewDateTime(time.Now()), StopTime: types.NewDateTime(time.Now())}, true},
		{firmware.GetDiagnosticsRequest{Location: "ftp:some/path", Retries: newInt(10), RetryInterval: newInt(10), StartTime: types.NewDateTime(time.Now())}, true},
		{firmware.GetDiagnosticsRequest{Location: "ftp:some/path", Retries: newInt(10), RetryInterval: newInt(10)}, true},
		{firmware.GetDiagnosticsRequest{Location: "ftp:some/path", Retries: newInt(10)}, true},
		{firmware.GetDiagnosticsRequest{Location: "ftp:some/path"}, true},
		{firmware.GetDiagnosticsRequest{}, false},
		{firmware.GetDiagnosticsRequest{Location: "invalidUri"}, false},
		{firmware.GetDiagnosticsRequest{Location: "ftp:some/path", Retries: newInt(-1)}, false},
		{firmware.GetDiagnosticsRequest{Location: "ftp:some/path", RetryInterval: newInt(-1)}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV16TestSuite) TestGetDiagnosticsConfirmationValidation() {
	confirmationTable := []GenericTestEntry{
		{firmware.GetDiagnosticsConfirmation{FileName: "someFileName"}, true},
		{firmware.GetDiagnosticsConfirmation{FileName: ""}, true},
		{firmware.GetDiagnosticsConfirmation{}, true},
		{firmware.GetDiagnosticsConfirmation{FileName: ">255............................................................................................................................................................................................................................................................"}, false},
	}
	ExecuteGenericTestTable(suite, confirmationTable)
}

func (suite *OcppV16TestSuite) TestGetDiagnosticsE2EMocked() {
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	location := "ftp:some/path"
	fileName := "diagnostics.json"
	retries := newInt(10)
	retryInterval := newInt(600)
	startTime := types.NewDateTime(time.Now().Add(-10 * time.Hour * 24))
	stopTime := types.NewDateTime(time.Now())
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"location":"%v","retries":%v,"retryInterval":%v,"startTime":"%v","stopTime":"%v"}]`,
		messageId, firmware.GetDiagnosticsFeatureName, location, *retries, *retryInterval, startTime.FormatTimestamp(), stopTime.FormatTimestamp())
	responseJson := fmt.Sprintf(`[3,"%v",{"fileName":"%v"}]`, messageId, fileName)
	getDiagnosticsConfirmation := firmware.NewGetDiagnosticsConfirmation()
	getDiagnosticsConfirmation.FileName = fileName
	channel := NewMockWebSocket(wsId)

	firmwareListener := &MockChargePointFirmwareManagementListener{}
	firmwareListener.On("OnGetDiagnostics", mock.Anything).Return(getDiagnosticsConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*firmware.GetDiagnosticsRequest)
		suite.Require().NotNil(request)
		suite.Require().True(ok)
		suite.Equal(location, request.Location)
		suite.Require().NotNil(request.Retries)
		suite.Equal(*retries, *request.Retries)
		suite.Require().NotNil(request.RetryInterval)
		suite.Equal(*retryInterval, *request.RetryInterval)
		assertDateTimeEquality(suite, *startTime, *request.StartTime)
		assertDateTimeEquality(suite, *stopTime, *request.StopTime)
	})
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	suite.chargePoint.SetFirmwareManagementHandler(firmwareListener)
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	suite.Require().Nil(err)
	resultChannel := make(chan bool, 1)
	err = suite.centralSystem.GetDiagnostics(wsId, func(confirmation *firmware.GetDiagnosticsConfirmation, err error) {
		suite.Require().Nil(err)
		suite.Require().NotNil(confirmation)
		suite.Equal(fileName, confirmation.FileName)
		resultChannel <- true
	}, location, func(request *firmware.GetDiagnosticsRequest) {
		request.RetryInterval = retryInterval
		request.Retries = retries
		request.StartTime = startTime
		request.StopTime = stopTime
	})
	suite.Require().Nil(err)
	result := <-resultChannel
	suite.True(result)
}

func (suite *OcppV16TestSuite) TestGetDiagnosticsInvalidEndpoint() {
	messageId := defaultMessageId
	location := "ftp:some/path"
	retries := 10
	retryInterval := 600
	startTime := types.NewDateTime(time.Now().Add(-10 * time.Hour * 24))
	stopTime := types.NewDateTime(time.Now())
	localListVersionRequest := firmware.NewGetDiagnosticsRequest(location)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"location":"%v","retries":%v,"retryInterval":%v,"startTime":"%v","stopTime":"%v"}]`,
		messageId, firmware.GetDiagnosticsFeatureName, location, retries, retryInterval, startTime.FormatTimestamp(), stopTime.FormatTimestamp())
	testUnsupportedRequestFromChargePoint(suite, localListVersionRequest, requestJson, messageId)
}
