package ocpp2_test

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/firmware"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

// Test
func (suite *OcppV2TestSuite) TestUpdateFirmwareRequestValidation() {
	fw := firmware.Firmware{
		Location:           "https://someurl",
		RetrieveDateTime:   types.NewDateTime(time.Now()),
		InstallDateTime:    types.NewDateTime(time.Now()),
		SigningCertificate: "1337c0de",
		Signature:          "deadc0de",
	}
	var requestTable = []GenericTestEntry{
		{firmware.UpdateFirmwareRequest{Retries: newInt(5), RetryInterval: newInt(300), RequestID: 42, Firmware: fw}, true},
		{firmware.UpdateFirmwareRequest{Retries: newInt(5), RequestID: 42, Firmware: fw}, true},
		{firmware.UpdateFirmwareRequest{RequestID: 42, Firmware: fw}, true},
		{firmware.UpdateFirmwareRequest{Firmware: fw}, true},
		{firmware.UpdateFirmwareRequest{Retries: newInt(5), RetryInterval: newInt(300), RequestID: 42, Firmware: firmware.Firmware{Location: "https://someurl", RetrieveDateTime: types.NewDateTime(time.Now())}}, true},
		{firmware.UpdateFirmwareRequest{}, false},
		{firmware.UpdateFirmwareRequest{Retries: newInt(-1), RetryInterval: newInt(300), RequestID: 42, Firmware: fw}, false},
		{firmware.UpdateFirmwareRequest{Retries: newInt(5), RetryInterval: newInt(-1), RequestID: 42, Firmware: fw}, false},
		{firmware.UpdateFirmwareRequest{Retries: newInt(5), RetryInterval: newInt(300), RequestID: -1, Firmware: fw}, false},
		{firmware.UpdateFirmwareRequest{Retries: newInt(5), RetryInterval: newInt(300), RequestID: 42, Firmware: firmware.Firmware{RetrieveDateTime: types.NewDateTime(time.Now()), InstallDateTime: types.NewDateTime(time.Now()), SigningCertificate: "1337c0de", Signature: "deadc0de"}}, false},
		{firmware.UpdateFirmwareRequest{Retries: newInt(5), RetryInterval: newInt(300), RequestID: 42, Firmware: firmware.Firmware{Location: "https://someurl", InstallDateTime: types.NewDateTime(time.Now()), SigningCertificate: "1337c0de", Signature: "deadc0de"}}, false},
		{firmware.UpdateFirmwareRequest{Retries: newInt(5), RetryInterval: newInt(300), RequestID: 42, Firmware: firmware.Firmware{Location: ">512.............................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................", RetrieveDateTime: types.NewDateTime(time.Now()), InstallDateTime: types.NewDateTime(time.Now()), SigningCertificate: "1337c0de", Signature: "deadc0de"}}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV2TestSuite) TestUpdateFirmwareResponseValidation() {
	var responseTable = []GenericTestEntry{
		{firmware.UpdateFirmwareResponse{Status: firmware.UpdateFirmwareStatusAccepted, StatusInfo: &types.StatusInfo{ReasonCode: "ok", AdditionalInfo: "someInfo"}}, true},
		{firmware.UpdateFirmwareResponse{Status: firmware.UpdateFirmwareStatusAccepted, StatusInfo: &types.StatusInfo{ReasonCode: "ok"}}, true},
		{firmware.UpdateFirmwareResponse{Status: firmware.UpdateFirmwareStatusAccepted}, true},
		{firmware.UpdateFirmwareResponse{}, false},
		{firmware.UpdateFirmwareResponse{Status: "invalidFirmwareUpdateStatus"}, false},
		{firmware.UpdateFirmwareResponse{Status: firmware.UpdateFirmwareStatusAccepted, StatusInfo: &types.StatusInfo{}}, false},
	}
	ExecuteGenericTestTable(suite, responseTable)
}

func (suite *OcppV2TestSuite) TestUpdateFirmwareE2EMocked() {
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	fw := firmware.Firmware{
		Location:           "https://someurl",
		RetrieveDateTime:   types.NewDateTime(time.Now()),
		InstallDateTime:    types.NewDateTime(time.Now()),
		SigningCertificate: "1337c0de",
		Signature:          "deadc0de",
	}
	retries := newInt(5)
	requestID := 42
	retryInterval := newInt(300)
	status := firmware.UpdateFirmwareStatusAccepted
	statusInfo := types.StatusInfo{ReasonCode: "ok", AdditionalInfo: "someInfo"}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"retries":%v,"retryInterval":%v,"requestId":%v,"firmware":{"location":"%v","retrieveDateTime":"%v","installDateTime":"%v","signingCertificate":"%v","signature":"%v"}}]`,
		messageId, firmware.UpdateFirmwareFeatureName, *retries, *retryInterval, requestID, fw.Location, fw.RetrieveDateTime.FormatTimestamp(), fw.InstallDateTime.FormatTimestamp(), fw.SigningCertificate, fw.Signature)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v","statusInfo":{"reasonCode":"%v","additionalInfo":"%v"}}]`,
		messageId, status, statusInfo.ReasonCode, statusInfo.AdditionalInfo)
	updateFirmwareResponse := firmware.NewUpdateFirmwareResponse(status)
	updateFirmwareResponse.StatusInfo = &statusInfo
	channel := NewMockWebSocket(wsId)

	handler := &MockChargingStationFirmwareHandler{}
	handler.On("OnUpdateFirmware", mock.Anything).Return(updateFirmwareResponse, nil).Run(func(args mock.Arguments) {
		request := args.Get(0).(*firmware.UpdateFirmwareRequest)
		suite.Require().NotNil(request)
		suite.Require().NotNil(request.Retries)
		suite.Equal(*retries, *request.Retries)
		suite.Require().NotNil(request.RetryInterval)
		suite.Equal(*retryInterval, *request.RetryInterval)
		suite.Equal(requestID, request.RequestID)
		suite.Equal(fw.Location, request.Firmware.Location)
		assertDateTimeEquality(suite, fw.RetrieveDateTime, request.Firmware.RetrieveDateTime)
		assertDateTimeEquality(suite, fw.InstallDateTime, request.Firmware.InstallDateTime)
		suite.Equal(fw.SigningCertificate, request.Firmware.SigningCertificate)
		suite.Equal(fw.Signature, request.Firmware.Signature)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	suite.Nil(err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.UpdateFirmware(wsId, func(resp *firmware.UpdateFirmwareResponse, err error) {
		suite.Nil(err)
		suite.Require().NotNil(resp)
		suite.Equal(status, resp.Status)
		suite.Require().NotNil(resp.StatusInfo)
		suite.Equal(statusInfo.ReasonCode, resp.StatusInfo.ReasonCode)
		suite.Equal(statusInfo.AdditionalInfo, resp.StatusInfo.AdditionalInfo)
		resultChannel <- true
	}, requestID, fw, func(request *firmware.UpdateFirmwareRequest) {
		request.Retries = retries
		request.RetryInterval = retryInterval
	})
	suite.Require().Nil(err)
	result := <-resultChannel
	suite.True(result)
}

func (suite *OcppV2TestSuite) TestUpdateFirmwareInvalidEndpoint() {
	messageId := defaultMessageId
	fw := firmware.Firmware{
		Location:           "https://someurl",
		RetrieveDateTime:   types.NewDateTime(time.Now()),
		InstallDateTime:    types.NewDateTime(time.Now()),
		SigningCertificate: "1337c0de",
		Signature:          "deadc0de",
	}
	retries := newInt(5)
	requestID := 42
	retryInterval := newInt(300)
	request := firmware.NewUpdateFirmwareRequest(requestID, fw)
	request.Retries = retries
	request.RetryInterval = retryInterval
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"retries":%v,"retryInterval":%v,"requestId":%v,"firmware":{"location":"%v","retrieveDateTime":"%v","installDateTime":"%v","signingCertificate":"%v","signature":"%v"}}]`,
		messageId, firmware.UpdateFirmwareFeatureName, *retries, *retryInterval, requestID, fw.Location, fw.RetrieveDateTime.FormatTimestamp(), fw.InstallDateTime.FormatTimestamp(), fw.SigningCertificate, fw.Signature)
	testUnsupportedRequestFromChargingStation(suite, request, requestJson, messageId)
}
