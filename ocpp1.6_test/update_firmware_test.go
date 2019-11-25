package ocpp16_test

import (
	"fmt"
	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"time"
)

// Test
func (suite *OcppV16TestSuite) TestUpdateFirmwareRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{ocpp16.UpdateFirmwareRequest{Location: "ftp:some/path", Retries: 10, RetryInterval: 10, RetrieveDate: ocpp16.NewDateTime(time.Now())}, true},
		{ocpp16.UpdateFirmwareRequest{Location: "ftp:some/path", Retries: 10, RetrieveDate: ocpp16.NewDateTime(time.Now())}, true},
		{ocpp16.UpdateFirmwareRequest{Location: "ftp:some/path", RetrieveDate: ocpp16.NewDateTime(time.Now())}, true},
		{ocpp16.UpdateFirmwareRequest{}, false},
		{ocpp16.UpdateFirmwareRequest{Location: "ftp:some/path"}, false},
		{ocpp16.UpdateFirmwareRequest{Location: "invalidUri", RetrieveDate: ocpp16.NewDateTime(time.Now())}, false},
		{ocpp16.UpdateFirmwareRequest{Location: "ftp:some/path", Retries: -1, RetrieveDate: ocpp16.NewDateTime(time.Now())}, false},
		{ocpp16.UpdateFirmwareRequest{Location: "ftp:some/path", RetryInterval: -1, RetrieveDate: ocpp16.NewDateTime(time.Now())}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV16TestSuite) TestUpdateFirmwareConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{ocpp16.UpdateFirmwareConfirmation{}, true},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV16TestSuite) TestUpdateFirmwareE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	location := "ftp:some/path"
	retries := 10
	retryInterval := 600
	retrieveDate := ocpp16.NewDateTime(time.Now())
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"location":"%v","retries":%v,"retrieveDate":"%v","retryInterval":%v}]`,
		messageId, ocpp16.UpdateFirmwareFeatureName, location, retries, retrieveDate.Format(ocpp16.ISO8601), retryInterval)
	responseJson := fmt.Sprintf(`[3,"%v",{}]`, messageId)
	updateFirmwareConfirmation := ocpp16.NewUpdateFirmwareConfirmation()
	channel := NewMockWebSocket(wsId)

	firmwareListener := MockChargePointFirmwareManagementListener{}
	firmwareListener.On("OnUpdateFirmware", mock.Anything).Return(updateFirmwareConfirmation, nil)
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	suite.chargePoint.SetFirmwareManagementListener(firmwareListener)
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	assert.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.centralSystem.UpdateFirmware(wsId, func(confirmation *ocpp16.UpdateFirmwareConfirmation, err error) {
		assert.Nil(t, err)
		assert.NotNil(t, confirmation)
		if confirmation == nil || err != nil {
			resultChannel <- false
		} else {
			resultChannel <- true
		}
	}, location, retrieveDate, func(request *ocpp16.UpdateFirmwareRequest) {
		request.RetryInterval = retryInterval
		request.Retries = retries
	})
	assert.Nil(t, err)
	if err == nil {
		result := <-resultChannel
		assert.True(t, result)
	}
}

func (suite *OcppV16TestSuite) TestUpdateFirmwareInvalidEndpoint() {
	messageId := defaultMessageId
	location := "ftp:some/path"
	retries := 10
	retryInterval := 600
	retrieveDate := ocpp16.NewDateTime(time.Now())
	localListVersionRequest := ocpp16.NewUpdateFirmwareRequest(location, retrieveDate)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"location":"%v","retries":%v,"retrieveDate":"%v","retryInterval":%v}]`,
		messageId, ocpp16.UpdateFirmwareFeatureName, location, retries, retrieveDate.Format(ocpp16.ISO8601), retryInterval)
	testUnsupportedRequestFromChargePoint(suite, localListVersionRequest, requestJson, messageId)
}
