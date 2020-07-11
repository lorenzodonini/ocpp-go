package ocpp16_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/firmware"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"time"
)

// Test
func (suite *OcppV16TestSuite) TestUpdateFirmwareRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{firmware.UpdateFirmwareRequest{Location: "ftp:some/path", Retries: newInt(10), RetryInterval: newInt(10), RetrieveDate: types.NewDateTime(time.Now())}, true},
		{firmware.UpdateFirmwareRequest{Location: "ftp:some/path", Retries: newInt(10), RetrieveDate: types.NewDateTime(time.Now())}, true},
		{firmware.UpdateFirmwareRequest{Location: "ftp:some/path", RetrieveDate: types.NewDateTime(time.Now())}, true},
		{firmware.UpdateFirmwareRequest{}, false},
		{firmware.UpdateFirmwareRequest{Location: "ftp:some/path"}, false},
		{firmware.UpdateFirmwareRequest{Location: "invalidUri", RetrieveDate: types.NewDateTime(time.Now())}, false},
		{firmware.UpdateFirmwareRequest{Location: "ftp:some/path", Retries: newInt(-1), RetrieveDate: types.NewDateTime(time.Now())}, false},
		{firmware.UpdateFirmwareRequest{Location: "ftp:some/path", RetryInterval: newInt(-1), RetrieveDate: types.NewDateTime(time.Now())}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV16TestSuite) TestUpdateFirmwareConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{firmware.UpdateFirmwareConfirmation{}, true},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV16TestSuite) TestUpdateFirmwareE2EMocked() {
	t := suite.T()
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

	firmwareListener := MockChargePointFirmwareManagementListener{}
	firmwareListener.On("OnUpdateFirmware", mock.Anything).Return(updateFirmwareConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*firmware.UpdateFirmwareRequest)
		require.NotNil(t, request)
		require.True(t, ok)
		assert.Equal(t, location, request.Location)
		assert.NotNil(t, request.Retries)
		assert.Equal(t, *retries, *request.Retries)
		assert.NotNil(t, request.RetryInterval)
		assert.Equal(t, *retryInterval, *request.RetryInterval)
		assertDateTimeEquality(t, *retrieveDate, *request.RetrieveDate)
	})
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	suite.chargePoint.SetFirmwareManagementHandler(firmwareListener)
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	assert.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.centralSystem.UpdateFirmware(wsId, func(confirmation *firmware.UpdateFirmwareConfirmation, err error) {
		require.Nil(t, err)
		require.NotNil(t, confirmation)
		resultChannel <- true
	}, location, retrieveDate, func(request *firmware.UpdateFirmwareRequest) {
		request.RetryInterval = retryInterval
		request.Retries = retries
	})
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
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
