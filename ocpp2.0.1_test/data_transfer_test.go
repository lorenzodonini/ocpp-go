package ocpp2_test

import (
	"fmt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/data"
)

// Test
func (suite *OcppV2TestSuite) TestDataTransferRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{data.DataTransferRequest{VendorID: "12345"}, true},
		{data.DataTransferRequest{VendorID: "12345", MessageID: "6789"}, true},
		{data.DataTransferRequest{VendorID: "12345", MessageID: "6789", Data: "mockData"}, true},
		{data.DataTransferRequest{}, false},
		{data.DataTransferRequest{VendorID: ">255............................................................................................................................................................................................................................................................"}, false},
		{data.DataTransferRequest{VendorID: "12345", MessageID: ">50................................................"}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestDataTransferConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{data.DataTransferResponse{Status: data.DataTransferStatusAccepted}, true},
		{data.DataTransferResponse{Status: data.DataTransferStatusRejected}, true},
		{data.DataTransferResponse{Status: data.DataTransferStatusUnknownMessageId}, true},
		{data.DataTransferResponse{Status: data.DataTransferStatusUnknownVendorId}, true},
		{data.DataTransferResponse{Status: "invalidDataTransferStatus"}, false},
		{data.DataTransferResponse{Status: data.DataTransferStatusAccepted, Data: "mockData"}, true},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestDataTransferFromChargePointE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	vendorId := "vendor1"
	status := data.DataTransferStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"vendorId":"%v"}]`, messageId, data.DataTransferFeatureName, vendorId)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	dataTransferConfirmation := data.NewDataTransferResponse(status)
	channel := NewMockWebSocket(wsId)

	handler := &MockCSMSDataHandler{}
	handler.On("OnDataTransfer", mock.AnythingOfType("string"), mock.Anything).Return(dataTransferConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*data.DataTransferRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, vendorId, request.VendorID)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	assert.Nil(t, err)
	confirmation, err := suite.chargingStation.DataTransfer(vendorId)
	assert.Nil(t, err)
	assert.NotNil(t, confirmation)
	assert.Equal(t, status, confirmation.Status)
}

func (suite *OcppV2TestSuite) TestDataTransferFromCentralSystemE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	vendorId := "vendor1"
	status := data.DataTransferStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"vendorId":"%v"}]`, messageId, data.DataTransferFeatureName, vendorId)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	dataTransferConfirmation := data.NewDataTransferResponse(status)
	channel := NewMockWebSocket(wsId)

	handler := &MockChargingStationDataHandler{}
	handler.On("OnDataTransfer", mock.Anything).Return(dataTransferConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*data.DataTransferRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, vendorId, request.VendorID)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	assert.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.DataTransfer(wsId, func(confirmation *data.DataTransferResponse, err error) {
		assert.Nil(t, err)
		assert.NotNil(t, confirmation)
		assert.Equal(t, status, confirmation.Status)
		resultChannel <- true
	}, vendorId)
	assert.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}
