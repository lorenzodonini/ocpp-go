package ocpp16_test

import (
	"fmt"
	ocpp16 "github.com/lorenzodonini/go-ocpp/ocpp1.6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test
func (suite *OcppV16TestSuite) TestDataTransferRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{ocpp16.DataTransferRequest{VendorId: "12345"}, true},
		{ocpp16.DataTransferRequest{VendorId: "12345", MessageId: "6789"}, true},
		{ocpp16.DataTransferRequest{VendorId: "12345", MessageId: "6789", Data: "mockData"}, true},
		{ocpp16.DataTransferRequest{}, false},
		{ocpp16.DataTransferRequest{VendorId: ">255............................................................................................................................................................................................................................................................"}, false},
		{ocpp16.DataTransferRequest{VendorId: "12345", MessageId: ">50................................................"}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV16TestSuite) TestDataTransferConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{ocpp16.DataTransferConfirmation{Status: ocpp16.DataTransferStatusAccepted}, true},
		{ocpp16.DataTransferConfirmation{Status: ocpp16.DataTransferStatusRejected}, true},
		{ocpp16.DataTransferConfirmation{Status: ocpp16.DataTransferStatusUnknownMessageId}, true},
		{ocpp16.DataTransferConfirmation{Status: ocpp16.DataTransferStatusUnknownVendorId}, true},
		{ocpp16.DataTransferConfirmation{Status: "invalidDataTransferStatus"}, false},
		{ocpp16.DataTransferConfirmation{Status: ocpp16.DataTransferStatusAccepted, Data: "mockData"}, true},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV16TestSuite) TestDataTransferFromChargePointE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	vendorId := "vendor1"
	status := ocpp16.DataTransferStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"vendorId":"%v"}]`, messageId, ocpp16.DataTransferFeatureName, vendorId)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	dataTransferConfirmation := ocpp16.NewDataTransferConfirmation(status)
	channel := NewMockWebSocket(wsId)

	coreListener := MockCentralSystemCoreListener{}
	coreListener.On("OnDataTransfer", mock.AnythingOfType("string"), mock.Anything).Return(dataTransferConfirmation, nil)
	setupDefaultCentralSystemHandlers(suite, coreListener, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	assert.Nil(t, err)
	confirmation, err := suite.chargePoint.DataTransfer(vendorId)
	assert.Nil(t, err)
	assert.NotNil(t, confirmation)
	assert.Equal(t, status, confirmation.Status)
}

func (suite *OcppV16TestSuite) TestDataTransferFromCentralSystemE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	vendorId := "vendor1"
	status := ocpp16.DataTransferStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"vendorId":"%v"}]`, messageId, ocpp16.DataTransferFeatureName, vendorId)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	dataTransferConfirmation := ocpp16.NewDataTransferConfirmation(status)
	channel := NewMockWebSocket(wsId)

	coreListener := MockChargePointCoreListener{}
	coreListener.On("OnDataTransfer", mock.Anything).Return(dataTransferConfirmation, nil)
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, coreListener, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	assert.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.centralSystem.DataTransfer(wsId, func(confirmation *ocpp16.DataTransferConfirmation, err error) {
		assert.Nil(t, err)
		assert.NotNil(t, confirmation)
		assert.Equal(t, status, confirmation.Status)
		resultChannel <- true
	}, vendorId)
	assert.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}
