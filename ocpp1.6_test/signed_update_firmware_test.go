package ocpp16_test

import (
	"fmt"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/securefirmware"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6_test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test
func (suite *OcppV16TestSuite) TestSignedUpdateFirmwareRequestValidation() {
	t := suite.T()
	fw := securefirmware.Firmware{
		Location:           "https://someurl",
		RetrieveDateTime:   types.NewDateTime(time.Now()),
		InstallDateTime:    types.NewDateTime(time.Now()),
		SigningCertificate: "1337c0de",
		Signature:          "deadc0de",
	}
	var requestTable = []GenericTestEntry{
		{securefirmware.SignedUpdateFirmwareRequest{Retries: newInt(5), RetryInterval: newInt(300), RequestID: 42, Firmware: fw}, true},
		{securefirmware.SignedUpdateFirmwareRequest{Retries: newInt(5), RequestID: 42, Firmware: fw}, true},
		{securefirmware.SignedUpdateFirmwareRequest{RequestID: 42, Firmware: fw}, true},
		{securefirmware.SignedUpdateFirmwareRequest{Firmware: fw}, true},
		{securefirmware.SignedUpdateFirmwareRequest{Retries: newInt(5), RetryInterval: newInt(300), RequestID: 42, Firmware: securefirmware.Firmware{Location: "https://someurl", RetrieveDateTime: types.NewDateTime(time.Now())}}, true},
		{securefirmware.SignedUpdateFirmwareRequest{}, false},
		{securefirmware.SignedUpdateFirmwareRequest{Retries: newInt(-1), RetryInterval: newInt(300), RequestID: 42, Firmware: fw}, false},
		{securefirmware.SignedUpdateFirmwareRequest{Retries: newInt(5), RetryInterval: newInt(-1), RequestID: 42, Firmware: fw}, false},
		{securefirmware.SignedUpdateFirmwareRequest{Retries: newInt(5), RetryInterval: newInt(300), RequestID: -1, Firmware: fw}, false},
		{securefirmware.SignedUpdateFirmwareRequest{Retries: newInt(5), RetryInterval: newInt(300), RequestID: 42, Firmware: securefirmware.Firmware{RetrieveDateTime: types.NewDateTime(time.Now()), InstallDateTime: types.NewDateTime(time.Now()), SigningCertificate: "1337c0de", Signature: "deadc0de"}}, false},
		{securefirmware.SignedUpdateFirmwareRequest{Retries: newInt(5), RetryInterval: newInt(300), RequestID: 42, Firmware: securefirmware.Firmware{Location: "https://someurl", InstallDateTime: types.NewDateTime(time.Now()), SigningCertificate: "1337c0de", Signature: "deadc0de"}}, false},
		{securefirmware.SignedUpdateFirmwareRequest{Retries: newInt(5), RetryInterval: newInt(300), RequestID: 42, Firmware: securefirmware.Firmware{Location: ">512.............................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................", RetrieveDateTime: types.NewDateTime(time.Now()), InstallDateTime: types.NewDateTime(time.Now()), SigningCertificate: "1337c0de", Signature: "deadc0de"}}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV16TestSuite) TestSignedUpdateFirmwareResponseValidation() {
	t := suite.T()
	var responseTable = []GenericTestEntry{
		{securefirmware.SignedUpdateFirmwareResponse{Status: securefirmware.UpdateFirmwareStatusAccepted}, true},
		{securefirmware.SignedUpdateFirmwareResponse{Status: securefirmware.UpdateFirmwareStatusAccepted}, true},
		{securefirmware.SignedUpdateFirmwareResponse{Status: securefirmware.UpdateFirmwareStatusAccepted}, true},
		{securefirmware.SignedUpdateFirmwareResponse{}, false},
		{securefirmware.SignedUpdateFirmwareResponse{Status: "invalidFirmwareUpdateStatus"}, false},
	}
	ExecuteGenericTestTable(t, responseTable)
}

func (suite *OcppV16TestSuite) TestSignedUpdateFirmwareE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	fw := securefirmware.Firmware{
		Location:           "https://someurl",
		RetrieveDateTime:   types.NewDateTime(time.Now()),
		InstallDateTime:    types.NewDateTime(time.Now()),
		SigningCertificate: "1337c0de",
		Signature:          "deadc0de",
	}
	retries := newInt(5)
	requestID := 42
	retryInterval := newInt(300)
	status := securefirmware.UpdateFirmwareStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"retries":%v,"retryInterval":%v,"requestId":%v,"firmware":{"location":"%v","retrieveDateTime":"%v","installDateTime":"%v","signingCertificate":"%v","signature":"%v"}}]`,
		messageId, securefirmware.SignedUpdateFirmwareFeatureName, *retries, *retryInterval, requestID, fw.Location, fw.RetrieveDateTime.FormatTimestamp(), fw.InstallDateTime.FormatTimestamp(), fw.SigningCertificate, fw.Signature)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	updateFirmwareResponse := securefirmware.NewSignedUpdateFirmwareResponse(status)
	channel := NewMockWebSocket(wsId)

	handler := mocks.NewMockSecureFirmwareChargePointHandler(t)
	handler.EXPECT().OnSignedUpdateFirmware(mock.Anything).RunAndReturn(func(request *securefirmware.SignedUpdateFirmwareRequest) (*securefirmware.SignedUpdateFirmwareResponse, error) {
		assert.Equal(t, *retries, *request.Retries)
		assert.Equal(t, *retryInterval, *request.RetryInterval)
		assert.Equal(t, requestID, request.RequestID)
		assert.Equal(t, fw.Location, request.Firmware.Location)
		assertDateTimeEquality(t, *fw.RetrieveDateTime, *request.Firmware.RetrieveDateTime)
		assertDateTimeEquality(t, *fw.InstallDateTime, *request.Firmware.InstallDateTime)
		assert.Equal(t, fw.SigningCertificate, request.Firmware.SigningCertificate)
		assert.Equal(t, fw.Signature, request.Firmware.Signature)

		return updateFirmwareResponse, nil
	})

	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	suite.chargePoint.SetSecureFirmwareHandler(handler)

	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	assert.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.centralSystem.SignedUpdateFirmware(wsId, func(resp *securefirmware.SignedUpdateFirmwareResponse, err error) {
		assert.Nil(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, status, resp.Status)
		resultChannel <- true
	}, requestID, fw, func(request *securefirmware.SignedUpdateFirmwareRequest) {
		request.Retries = retries
		request.RetryInterval = retryInterval
	})
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV16TestSuite) TestSignedUpdateFirmwareInvalidEndpoint() {
	messageId := defaultMessageId
	fw := securefirmware.Firmware{
		Location:           "https://someurl",
		RetrieveDateTime:   types.NewDateTime(time.Now()),
		InstallDateTime:    types.NewDateTime(time.Now()),
		SigningCertificate: "1337c0de",
		Signature:          "deadc0de",
	}
	retries := newInt(5)
	requestID := 42
	retryInterval := newInt(300)
	request := securefirmware.NewSignedUpdateFirmwareRequest(requestID, fw)
	request.Retries = retries
	request.RetryInterval = retryInterval
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"retries":%v,"retryInterval":%v,"requestId":%v,"firmware":{"location":"%v","retrieveDateTime":"%v","installDateTime":"%v","signingCertificate":"%v","signature":"%v"}}]`,
		messageId, securefirmware.SignedUpdateFirmwareFeatureName, *retries, *retryInterval, requestID, fw.Location, fw.RetrieveDateTime.FormatTimestamp(), fw.InstallDateTime.FormatTimestamp(), fw.SigningCertificate, fw.Signature)
	testUnsupportedRequestFromChargePoint(suite, request, requestJson, messageId)
}
