package ocpp2_test

import (
	"fmt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/firmware"
)

// Test
func (suite *OcppV2TestSuite) TestUnpublishFirmwareRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{firmware.UnpublishFirmwareRequest{Checksum: "deadc0de"}, true},
		{firmware.UnpublishFirmwareRequest{}, false},
		{firmware.UnpublishFirmwareRequest{Checksum: ">32.............................."}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestUnpublishFirmwareResponseValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{firmware.UnpublishFirmwareResponse{Status: firmware.UnpublishFirmwareStatusUnpublished}, true},
		{firmware.UnpublishFirmwareResponse{Status: firmware.UnpublishFirmwareStatusNoFirmware}, true},
		{firmware.UnpublishFirmwareResponse{Status: firmware.UnpublishFirmwareStatusDownloadOngoing}, true},
		{firmware.UnpublishFirmwareResponse{}, false},
		{firmware.UnpublishFirmwareResponse{Status: "invalidUnpublishFirmwareStatus"}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestUnpublishFirmwareE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	checksum := "deadc0de"
	status := firmware.UnpublishFirmwareStatusUnpublished
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"checksum":"%v"}]`,
		messageId, firmware.UnpublishFirmwareFeatureName, checksum)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`,
		messageId, status)
	unpublishFirmwareResponse := firmware.NewUnpublishFirmwareResponse(status)
	channel := NewMockWebSocket(wsId)

	handler := &MockChargingStationFirmwareHandler{}
	handler.On("OnUnpublishFirmware", mock.Anything).Return(unpublishFirmwareResponse, nil).Run(func(args mock.Arguments) {
		request := args.Get(0).(*firmware.UnpublishFirmwareRequest)
		assert.Equal(t, checksum, request.Checksum)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	assert.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.UnpublishFirmware(wsId, func(resp *firmware.UnpublishFirmwareResponse, err error) {
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, unpublishFirmwareResponse.Status, resp.Status)
		resultChannel <- true
	}, checksum)
	require.Nil(t, err)
	if err == nil {
		result := <-resultChannel
		assert.True(t, result)
	}
}

func (suite *OcppV2TestSuite) TestUnpublishFirmwareInvalidEndpoint() {
	messageId := defaultMessageId
	checksum := "deadc0de"
	request := firmware.NewUnpublishFirmwareRequest(checksum)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"checksum":"%v"}]`,
		messageId, firmware.UnpublishFirmwareFeatureName, checksum)
	testUnsupportedRequestFromChargingStation(suite, request, requestJson, messageId)
}
