package ocpp16_test

import (
	"fmt"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/firmware"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/stretchr/testify/mock"
)

// Test
func (suite *OcppV16TestSuite) TestUpdateFirmwareRequestValidation() {
	requestTable := []GenericTestEntry{
		{firmware.UpdateFirmwareRequest{Location: "ftp:some/path", Retries: newInt(10), RetryInterval: newInt(10), RetrieveDate: types.NewDateTime(time.Now())}, true},
		{firmware.UpdateFirmwareRequest{Location: "ftp:some/path", Retries: newInt(10), RetrieveDate: types.NewDateTime(time.Now())}, true},
		{firmware.UpdateFirmwareRequest{Location: "ftp:some/path", RetrieveDate: types.NewDateTime(time.Now())}, true},
		{firmware.UpdateFirmwareRequest{}, false},
		{firmware.UpdateFirmwareRequest{Location: "ftp:some/path"}, false},
		{firmware.UpdateFirmwareRequest{Location: "invalidUri", RetrieveDate: types.NewDateTime(time.Now())}, false},
		{firmware.UpdateFirmwareRequest{Location: "ftp:some/path", Retries: newInt(-1), RetrieveDate: types.NewDateTime(time.Now())}, false},
		{firmware.UpdateFirmwareRequest{Location: "ftp:some/path", RetryInterval: newInt(-1), RetrieveDate: types.NewDateTime(time.Now())}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV16TestSuite) TestUpdateFirmwareConfirmationValidation() {
	confirmationTable := []GenericTestEntry{
		{firmware.UpdateFirmwareConfirmation{}, true},
	}
	ExecuteGenericTestTable(suite, confirmationTable)
}

func (suite *OcppV16TestSuite) TestUpdateFirmwareE2EMocked() {
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	location := "ftp:some/path"
	retries := newInt(10)
	retryInterval := newInt(600)
	retrieveDate := types.NewDateTime(time.Now())
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"location":"%v","retries":%v,"retrieveDate":"%v","retryInterval":%v}]`,
		messageId, firmware.UpdateFirmwareFeatureName, location, *retries, retrieveDate.FormatTimestamp(), *retryInterval)
	responseJson := fmt.Sprintf(`[3,"%v",{}]`, messageId)
	updateFirmwareConfirmation := firmware.NewUpdateFirmwareConfirmation()
	channel := NewMockWebSocket(wsId)

	firmwareListener := &MockChargePointFirmwareManagementListener{}
	firmwareListener.On("OnUpdateFirmware", mock.Anything).Return(updateFirmwareConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*firmware.UpdateFirmwareRequest)
		suite.Require().NotNil(request)
		suite.Require().True(ok)
		suite.Equal(location, request.Location)
		suite.NotNil(request.Retries)
		suite.Equal(*retries, *request.Retries)
		suite.NotNil(request.RetryInterval)
		suite.Equal(*retryInterval, *request.RetryInterval)
		assertDateTimeEquality(suite, *retrieveDate, *request.RetrieveDate)
	})
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	suite.chargePoint.SetFirmwareManagementHandler(firmwareListener)
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	suite.Nil(err)
	resultChannel := make(chan bool, 1)
	err = suite.centralSystem.UpdateFirmware(wsId, func(confirmation *firmware.UpdateFirmwareConfirmation, err error) {
		suite.Require().Nil(err)
		suite.Require().NotNil(confirmation)
		resultChannel <- true
	}, location, retrieveDate, func(request *firmware.UpdateFirmwareRequest) {
		request.RetryInterval = retryInterval
		request.Retries = retries
	})
	suite.Require().Nil(err)
	result := <-resultChannel
	suite.True(result)
}

func (suite *OcppV16TestSuite) TestUpdateFirmwareInvalidEndpoint() {
	messageId := defaultMessageId
	location := "ftp:some/path"
	retries := 10
	retryInterval := 600
	retrieveDate := types.NewDateTime(time.Now())
	localListVersionRequest := firmware.NewUpdateFirmwareRequest(location, retrieveDate)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"location":"%v","retries":%v,"retrieveDate":"%v","retryInterval":%v}]`,
		messageId, firmware.UpdateFirmwareFeatureName, location, retries, retrieveDate.FormatTimestamp(), retryInterval)
	testUnsupportedRequestFromChargePoint(suite, localListVersionRequest, requestJson, messageId)
}
