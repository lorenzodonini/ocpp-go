package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test
func (suite *OcppV2TestSuite) TestDataTransferRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{ocpp2.DataTransferRequest{VendorId: "12345"}, true},
		{ocpp2.DataTransferRequest{VendorId: "12345", MessageId: "6789"}, true},
		{ocpp2.DataTransferRequest{VendorId: "12345", MessageId: "6789", Data: "mockData"}, true},
		{ocpp2.DataTransferRequest{}, false},
		{ocpp2.DataTransferRequest{VendorId: ">255............................................................................................................................................................................................................................................................"}, false},
		{ocpp2.DataTransferRequest{VendorId: "12345", MessageId: ">50................................................"}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestDataTransferConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{ocpp2.DataTransferConfirmation{Status: ocpp2.DataTransferStatusAccepted}, true},
		{ocpp2.DataTransferConfirmation{Status: ocpp2.DataTransferStatusRejected}, true},
		{ocpp2.DataTransferConfirmation{Status: ocpp2.DataTransferStatusUnknownMessageId}, true},
		{ocpp2.DataTransferConfirmation{Status: ocpp2.DataTransferStatusUnknownVendorId}, true},
		{ocpp2.DataTransferConfirmation{Status: "invalidDataTransferStatus"}, false},
		{ocpp2.DataTransferConfirmation{Status: ocpp2.DataTransferStatusAccepted, Data: "mockData"}, true},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestDataTransferFromChargePointE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	vendorId := "vendor1"
	status := ocpp2.DataTransferStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"vendorId":"%v"}]`, messageId, ocpp2.DataTransferFeatureName, vendorId)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	dataTransferConfirmation := ocpp2.NewDataTransferConfirmation(status)
	channel := NewMockWebSocket(wsId)

	coreListener := MockCentralSystemCoreListener{}
	coreListener.On("OnDataTransfer", mock.AnythingOfType("string"), mock.Anything).Return(dataTransferConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*ocpp2.DataTransferRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, vendorId, request.VendorId)
	})
	setupDefaultCentralSystemHandlers(suite, coreListener, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	assert.Nil(t, err)
	confirmation, err := suite.chargePoint.DataTransfer(vendorId)
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
	status := ocpp2.DataTransferStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"vendorId":"%v"}]`, messageId, ocpp2.DataTransferFeatureName, vendorId)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	dataTransferConfirmation := ocpp2.NewDataTransferConfirmation(status)
	channel := NewMockWebSocket(wsId)

	coreListener := MockChargePointCoreListener{}
	coreListener.On("OnDataTransfer", mock.Anything).Return(dataTransferConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*ocpp2.DataTransferRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, vendorId, request.VendorId)
	})
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, coreListener, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	assert.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.DataTransfer(wsId, func(confirmation *ocpp2.DataTransferConfirmation, err error) {
		assert.Nil(t, err)
		assert.NotNil(t, confirmation)
		assert.Equal(t, status, confirmation.Status)
		resultChannel <- true
	}, vendorId)
	assert.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}
