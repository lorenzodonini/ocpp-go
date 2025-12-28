package ocpp16_test

import (
	"encoding/json"
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/stretchr/testify/mock"
)

type CustomData struct {
	Field1 string `json:"field1" validate:"required"`
	Field2 int    `json:"field2" validate:"gt=0"`
}

func parseCustomData(req *core.DataTransferRequest) (CustomData, error) {
	jsonString, _ := json.Marshal(req.Data)
	var result CustomData
	err := json.Unmarshal(jsonString, &result)
	return result, err
}

// Test
func (suite *OcppV16TestSuite) TestDataTransferRequestValidation() {
	var requestTable = []GenericTestEntry{
		{core.DataTransferRequest{VendorId: "12345"}, true},
		{core.DataTransferRequest{VendorId: "12345", MessageId: "6789"}, true},
		{core.DataTransferRequest{VendorId: "12345", MessageId: "6789", Data: "mockData"}, true},
		{core.DataTransferRequest{}, false},
		{core.DataTransferRequest{VendorId: ">255............................................................................................................................................................................................................................................................"}, false},
		{core.DataTransferRequest{VendorId: "12345", MessageId: ">50................................................"}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV16TestSuite) TestDataTransferConfirmationValidation() {
	var confirmationTable = []GenericTestEntry{
		{core.DataTransferConfirmation{Status: core.DataTransferStatusAccepted}, true},
		{core.DataTransferConfirmation{Status: core.DataTransferStatusRejected}, true},
		{core.DataTransferConfirmation{Status: core.DataTransferStatusUnknownMessageId}, true},
		{core.DataTransferConfirmation{Status: core.DataTransferStatusUnknownVendorId}, true},
		{core.DataTransferConfirmation{Status: "invalidDataTransferStatus"}, false},
		{core.DataTransferConfirmation{Status: core.DataTransferStatusAccepted, Data: "mockData"}, true},
	}
	ExecuteGenericTestTable(suite, confirmationTable)
}

func (suite *OcppV16TestSuite) TestDataTransferFromChargePointE2EMocked() {
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	vendorId := "vendor1"
	data := CustomData{Field1: "dummyData", Field2: 42}
	status := core.DataTransferStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"vendorId":"%v","data":{"field1":"%v","field2":%v}}]`, messageId, core.DataTransferFeatureName, vendorId, data.Field1, data.Field2)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	dataTransferConfirmation := core.NewDataTransferConfirmation(status)
	channel := NewMockWebSocket(wsId)

	coreListener := &MockCentralSystemCoreListener{}
	coreListener.On("OnDataTransfer", mock.AnythingOfType("string"), mock.Anything).Return(dataTransferConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*core.DataTransferRequest)
		suite.Require().NotNil(request)
		suite.Require().True(ok)
		suite.Equal(vendorId, request.VendorId)
		suite.Require().NotNil(request.Data)
		customData, err := parseCustomData(request)
		suite.Require().Nil(err)
		suite.Require().NotNil(customData)
		suite.Equal(data.Field1, customData.Field1)
		suite.Equal(data.Field2, customData.Field2)
	})
	setupDefaultCentralSystemHandlers(suite, coreListener, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	suite.Require().Nil(err)
	confirmation, err := suite.chargePoint.DataTransfer(vendorId, func(request *core.DataTransferRequest) {
		request.Data = data
	})
	suite.Require().Nil(err)
	suite.Require().NotNil(confirmation)
	suite.Equal(status, confirmation.Status)
}

func (suite *OcppV16TestSuite) TestDataTransferFromCentralSystemE2EMocked() {
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	vendorId := "vendor1"
	data := CustomData{Field1: "dummyData", Field2: 42}
	status := core.DataTransferStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"vendorId":"%v","data":{"field1":"%v","field2":%v}}]`, messageId, core.DataTransferFeatureName, vendorId, data.Field1, data.Field2)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	dataTransferConfirmation := core.NewDataTransferConfirmation(status)
	channel := NewMockWebSocket(wsId)

	coreListener := &MockChargePointCoreListener{}
	coreListener.On("OnDataTransfer", mock.Anything).Return(dataTransferConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*core.DataTransferRequest)
		suite.Require().NotNil(request)
		suite.Require().True(ok)
		suite.Equal(vendorId, request.VendorId)
		suite.Require().NotNil(request.Data)
		customData, err := parseCustomData(request)
		suite.Require().Nil(err)
		suite.Require().NotNil(customData)
		suite.Equal(data.Field1, customData.Field1)
		suite.Equal(data.Field2, customData.Field2)
	})
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, coreListener, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	suite.Require().Nil(err)
	resultChannel := make(chan bool, 1)
	err = suite.centralSystem.DataTransfer(wsId, func(confirmation *core.DataTransferConfirmation, err error) {
		suite.Require().Nil(err)
		suite.Require().NotNil(confirmation)
		suite.Equal(status, confirmation.Status)
		resultChannel <- true
	}, vendorId, func(request *core.DataTransferRequest) {
		request.Data = data
	})
	suite.Require().Nil(err)
	result := <-resultChannel
	suite.True(result)
}
