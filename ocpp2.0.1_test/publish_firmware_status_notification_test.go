package ocpp2_test

import (
	"fmt"

	"github.com/stretchr/testify/mock"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/firmware"
)

// Test
func (suite *OcppV2TestSuite) TestPublishFirmwareStatusNotificationRequestValidation() {
	var requestTable = []GenericTestEntry{
		{firmware.PublishFirmwareStatusNotificationRequest{Status: firmware.PublishFirmwareStatusPublished, Location: []string{"http://someUri"}, RequestID: newInt(42)}, true},
		{firmware.PublishFirmwareStatusNotificationRequest{Status: firmware.PublishFirmwareStatusPublished, Location: []string{"http://someUri"}}, true},
		{firmware.PublishFirmwareStatusNotificationRequest{Status: firmware.PublishFirmwareStatusChecksumVerified, Location: []string{}}, true},
		{firmware.PublishFirmwareStatusNotificationRequest{Status: firmware.PublishFirmwareStatusChecksumVerified}, true},
		{firmware.PublishFirmwareStatusNotificationRequest{Status: firmware.PublishFirmwareStatusDownloaded}, true},
		{firmware.PublishFirmwareStatusNotificationRequest{Status: firmware.PublishFirmwareStatusDownloadFailed}, true},
		{firmware.PublishFirmwareStatusNotificationRequest{Status: firmware.PublishFirmwareStatusDownloading}, true},
		{firmware.PublishFirmwareStatusNotificationRequest{Status: firmware.PublishFirmwareStatusDownloadScheduled}, true},
		{firmware.PublishFirmwareStatusNotificationRequest{Status: firmware.PublishFirmwareStatusDownloadPaused}, true},
		{firmware.PublishFirmwareStatusNotificationRequest{Status: firmware.PublishFirmwareStatusInvalidChecksum}, true},
		{firmware.PublishFirmwareStatusNotificationRequest{Status: firmware.PublishFirmwareStatusIdle}, true},
		{firmware.PublishFirmwareStatusNotificationRequest{Status: firmware.PublishFirmwareStatusPublishFailed}, true},
		{firmware.PublishFirmwareStatusNotificationRequest{}, false},
		{firmware.PublishFirmwareStatusNotificationRequest{Status: "invalidStatus"}, false},
		{firmware.PublishFirmwareStatusNotificationRequest{Status: firmware.PublishFirmwareStatusPublished, Location: []string{"http://someUri"}, RequestID: newInt(-1)}, false},
		{firmware.PublishFirmwareStatusNotificationRequest{Status: firmware.PublishFirmwareStatusPublished, Location: []string{"http://someUri>512..............................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................."}, RequestID: newInt(42)}, false},
		//TODO: add test for empty location field with published status
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV2TestSuite) TestPublishFirmwareStatusNotificationResponseValidation() {
	var confirmationTable = []GenericTestEntry{
		{firmware.PublishFirmwareStatusNotificationResponse{}, true},
	}
	ExecuteGenericTestTable(suite, confirmationTable)
}

func (suite *OcppV2TestSuite) TestPublishFirmwareStatusNotificationE2EMocked() {
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	status := firmware.PublishFirmwareStatusPublished
	requestID := newInt(42)
	location := []string{"https://someUri"}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"status":"%v","location":["%v"],"requestId":%v}]`,
		messageId, firmware.PublishFirmwareStatusNotificationFeatureName, status, location[0], *requestID)
	responseJson := fmt.Sprintf(`[3,"%v",{}]`, messageId)
	publishFirmwareStatusNotificationResponse := firmware.NewPublishFirmwareStatusNotificationResponse()
	channel := NewMockWebSocket(wsId)

	handler := &MockCSMSFirmwareHandler{}
	handler.On("OnPublishFirmwareStatusNotification", mock.AnythingOfType("string"), mock.Anything).Return(publishFirmwareStatusNotificationResponse, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*firmware.PublishFirmwareStatusNotificationRequest)
		suite.Require().True(ok)
		suite.Equal(status, request.Status)
		suite.Require().Len(request.Location, len(location))
		suite.Equal(location[0], request.Location[0])
		suite.Require().NotNil(request.RequestID)
		suite.Equal(*requestID, *request.RequestID)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	suite.Require().Nil(err)
	response, err := suite.chargingStation.PublishFirmwareStatusNotification(status, func(request *firmware.PublishFirmwareStatusNotificationRequest) {
		request.Location = location
		request.RequestID = requestID
	})
	suite.Nil(err)
	suite.NotNil(response)
}

func (suite *OcppV2TestSuite) TestPublishFirmwareStatusNotificationInvalidEndpoint() {
	messageId := defaultMessageId
	status := firmware.PublishFirmwareStatusPublished
	requestID := newInt(42)
	location := []string{"https://someUri"}
	request := firmware.NewPublishFirmwareStatusNotificationRequest(status)
	request.Location = location
	request.RequestID = requestID
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"status":"%v","location":["%v"],"requestId":%v}]`,
		messageId, firmware.PublishFirmwareStatusNotificationFeatureName, status, location[0], *requestID)
	testUnsupportedRequestFromCentralSystem(suite, request, requestJson, messageId)
}
