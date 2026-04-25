package ocpp2_test

import (
	"fmt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/firmware"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

// Test
func (suite *OcppV2TestSuite) TestPublishFirmwareE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	location := "http://someUrl"
	retries := tests.NewInt(5)
	checksum := "deadc0d3"
	requestID := 42
	retryInterval := tests.NewInt(300)
	status := types.GenericStatusAccepted
	statusInfo := types.StatusInfo{ReasonCode: "ok", AdditionalInfo: "someInfo"}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"location":"%v","retries":%v,"checksum":"%v","requestId":%v,"retryInterval":%v}]`,
		messageId, firmware.PublishFirmwareFeatureName, location, *retries, checksum, requestID, *retryInterval)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v","statusInfo":{"reasonCode":"%v","additionalInfo":"%v"}}]`,
		messageId, status, statusInfo.ReasonCode, statusInfo.AdditionalInfo)
	publishFirmwareResponse := firmware.NewPublishFirmwareResponse(status)
	publishFirmwareResponse.StatusInfo = &statusInfo
	channel := NewMockWebSocket(wsId)

	handler := &MockChargingStationFirmwareHandler{}
	handler.On("OnPublishFirmware", mock.Anything).Return(publishFirmwareResponse, nil).Run(func(args mock.Arguments) {
		request := args.Get(0).(*firmware.PublishFirmwareRequest)
		assert.Equal(t, location, request.Location)
		assert.Equal(t, *retries, *request.Retries)
		assert.Equal(t, checksum, request.Checksum)
		assert.Equal(t, requestID, request.RequestID)
		assert.Equal(t, *retryInterval, *request.RetryInterval)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	assert.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.PublishFirmware(wsId, func(resp *firmware.PublishFirmwareResponse, err error) {
		assert.Nil(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, publishFirmwareResponse.Status, resp.Status)
		require.NotNil(t, resp.StatusInfo)
		assert.Equal(t, publishFirmwareResponse.StatusInfo.ReasonCode, resp.StatusInfo.ReasonCode)
		assert.Equal(t, publishFirmwareResponse.StatusInfo.AdditionalInfo, resp.StatusInfo.AdditionalInfo)
		resultChannel <- true
	}, location, checksum, requestID, func(request *firmware.PublishFirmwareRequest) {
		request.Retries = retries
		request.RetryInterval = retryInterval
	})
	require.Nil(t, err)
	if err == nil {
		result := <-resultChannel
		assert.True(t, result)
	}
}

func (suite *OcppV2TestSuite) TestPublishFirmwareInvalidEndpoint() {
	messageId := defaultMessageId
	location := "http://someUrl"
	retries := tests.NewInt(5)
	checksum := "deadc0d3"
	requestID := 42
	retryInterval := tests.NewInt(300)
	request := firmware.NewPublishFirmwareRequest(location, checksum, requestID)
	request.Retries = retries
	request.RetryInterval = retryInterval
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"location":"%v","retries":%v,"checksum":"%v","requestId":%v,"retryInterval":%v}]`,
		messageId, firmware.PublishFirmwareFeatureName, location, *retries, checksum, requestID, *retryInterval)
	testUnsupportedRequestFromChargingStation(suite, request, requestJson, messageId)
}
