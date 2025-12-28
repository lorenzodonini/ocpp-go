package ocpp2_test

import (
	"fmt"

	"github.com/stretchr/testify/mock"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/firmware"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

// Test
func (suite *OcppV2TestSuite) TestPublishFirmwareRequestValidation() {
	var requestTable = []GenericTestEntry{
		{firmware.NewPublishFirmwareRequest("https://someurl", "deadbeef", 42), true},
		{firmware.PublishFirmwareRequest{Location: "http://someurl", Retries: newInt(5), Checksum: "deadbeef", RequestID: 42, RetryInterval: newInt(300)}, true},
		{firmware.PublishFirmwareRequest{Location: "http://someurl", Retries: newInt(5), Checksum: "deadbeef", RequestID: 42}, true},
		{firmware.PublishFirmwareRequest{Location: "http://someurl", Checksum: "deadbeef", RequestID: 42}, true},
		{firmware.PublishFirmwareRequest{Location: "http://someurl", Checksum: "deadbeef"}, true},
		{firmware.PublishFirmwareRequest{Location: "http://someurl"}, false},
		{firmware.PublishFirmwareRequest{Checksum: "deadbeef"}, false},
		{firmware.PublishFirmwareRequest{}, false},
		{firmware.PublishFirmwareRequest{Location: "http://someurl", Retries: newInt(5), Checksum: "deadbeef", RequestID: 42, RetryInterval: newInt(-1)}, false},
		{firmware.PublishFirmwareRequest{Location: "http://someurl", Retries: newInt(5), Checksum: "deadbeef", RequestID: -1, RetryInterval: newInt(300)}, false},
		{firmware.PublishFirmwareRequest{Location: "http://someurl", Retries: newInt(5), Checksum: ">32..............................", RequestID: 42, RetryInterval: newInt(300)}, false},
		{firmware.PublishFirmwareRequest{Location: "http://someurl", Retries: newInt(-1), Checksum: "deadbeef", RequestID: 42, RetryInterval: newInt(300)}, false},
		{firmware.PublishFirmwareRequest{Location: ">512.............................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................", Retries: newInt(5), Checksum: "deadbeef", RequestID: 42, RetryInterval: newInt(300)}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV2TestSuite) TestPublishFirmwareResponseValidation() {
	var confirmationTable = []GenericTestEntry{
		{firmware.PublishFirmwareResponse{Status: types.GenericStatusAccepted, StatusInfo: &types.StatusInfo{ReasonCode: "ok", AdditionalInfo: "someInfo"}}, true},
		{firmware.PublishFirmwareResponse{Status: types.GenericStatusAccepted}, true},
		{firmware.PublishFirmwareResponse{}, false},
		{firmware.PublishFirmwareResponse{Status: "invalidStatus"}, false},
		{firmware.PublishFirmwareResponse{Status: types.GenericStatusAccepted, StatusInfo: &types.StatusInfo{}}, false},
	}
	ExecuteGenericTestTable(suite, confirmationTable)
}

func (suite *OcppV2TestSuite) TestPublishFirmwareE2EMocked() {
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	location := "http://someUrl"
	retries := newInt(5)
	checksum := "deadc0d3"
	requestID := 42
	retryInterval := newInt(300)
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
		suite.Equal(location, request.Location)
		suite.Equal(*retries, *request.Retries)
		suite.Equal(checksum, request.Checksum)
		suite.Equal(requestID, request.RequestID)
		suite.Equal(*retryInterval, *request.RetryInterval)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	suite.Nil(err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.PublishFirmware(wsId, func(resp *firmware.PublishFirmwareResponse, err error) {
		suite.Nil(err)
		suite.Require().NotNil(resp)
		suite.Equal(publishFirmwareResponse.Status, resp.Status)
		suite.Require().NotNil(resp.StatusInfo)
		suite.Equal(publishFirmwareResponse.StatusInfo.ReasonCode, resp.StatusInfo.ReasonCode)
		suite.Equal(publishFirmwareResponse.StatusInfo.AdditionalInfo, resp.StatusInfo.AdditionalInfo)
		resultChannel <- true
	}, location, checksum, requestID, func(request *firmware.PublishFirmwareRequest) {
		request.Retries = retries
		request.RetryInterval = retryInterval
	})
	suite.Require().Nil(err)
	if err == nil {
		result := <-resultChannel
		suite.True(result)
	}
}

func (suite *OcppV2TestSuite) TestPublishFirmwareInvalidEndpoint() {
	messageId := defaultMessageId
	location := "http://someUrl"
	retries := newInt(5)
	checksum := "deadc0d3"
	requestID := 42
	retryInterval := newInt(300)
	request := firmware.NewPublishFirmwareRequest(location, checksum, requestID)
	request.Retries = retries
	request.RetryInterval = retryInterval
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"location":"%v","retries":%v,"checksum":"%v","requestId":%v,"retryInterval":%v}]`,
		messageId, firmware.PublishFirmwareFeatureName, location, *retries, checksum, requestID, *retryInterval)
	testUnsupportedRequestFromChargingStation(suite, request, requestJson, messageId)
}
