package ocpp2_test

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/data"

	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocppj"
	"github.com/stretchr/testify/mock"
)

func (suite *OcppV2TestSuite) TestChargePointSendResponseError() {
	wsId := "test_id"
	channel := NewMockWebSocket(wsId)
	var ocppErr *ocpp.Error
	// Setup internal communication and listeners
	dataListener := &MockChargingStationDataHandler{}
	suite.chargingStation.SetDataHandler(dataListener)
	suite.mockWsClient.On("Start", mock.AnythingOfType("string")).Return(nil).Run(func(args mock.Arguments) {
		// Notify server of incoming connection
		suite.mockWsServer.NewClientHandler(channel)
	})
	suite.mockWsClient.On("Write", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		rawMsg := args.Get(0)
		bytes := rawMsg.([]byte)
		err := suite.mockWsServer.MessageHandler(channel, bytes)
		suite.Nil(err)
	})
	suite.mockWsServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	suite.mockWsServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		rawMsg := args.Get(1)
		bytes := rawMsg.([]byte)
		err := suite.mockWsClient.MessageHandler(bytes)
		suite.NoError(err)
	})
	// Run Tests
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start("someUrl")
	suite.Require().Nil(err)
	resultChannel := make(chan error, 1)
	// Test 1: occurrence validation error
	dataTransferResponse := data.NewDataTransferResponse(data.DataTransferStatusAccepted)
	dataTransferResponse.Data = struct {
		Field1 string `validate:"required"`
	}{Field1: ""}
	dataListener.On("OnDataTransfer", mock.Anything).Return(dataTransferResponse, nil)
	err = suite.csms.DataTransfer(wsId, func(response *data.DataTransferResponse, err error) {
		suite.Require().Nil(response)
		suite.Require().Error(err)
		resultChannel <- err
	}, "vendor1")
	suite.Require().Nil(err)
	result := <-resultChannel
	suite.Require().IsType(&ocpp.Error{}, result)
	ocppErr = result.(*ocpp.Error)
	suite.Equal(ocppj.OccurrenceConstraintViolationV2, ocppErr.Code)
	suite.Equal("Field CallResult.Payload.Data.Field1 required but not found for feature DataTransfer", ocppErr.Description)
	// Test 2: marshaling error
	dataTransferResponse = data.NewDataTransferResponse(data.DataTransferStatusAccepted)
	dataTransferResponse.Data = make(chan struct{})
	dataListener.ExpectedCalls = nil
	dataListener.On("OnDataTransfer", mock.Anything).Return(dataTransferResponse, nil)
	err = suite.csms.DataTransfer(wsId, func(response *data.DataTransferResponse, err error) {
		suite.Require().Nil(response)
		suite.Require().Error(err)
		resultChannel <- err
	}, "vendor1")
	suite.Require().Nil(err)
	result = <-resultChannel
	suite.Require().IsType(&ocpp.Error{}, result)
	ocppErr = result.(*ocpp.Error)
	suite.Equal(ocppj.GenericError, ocppErr.Code)
	suite.Equal("json: unsupported type: chan struct {}", ocppErr.Description)
	// Test 3: no results in callback
	dataListener.ExpectedCalls = nil
	dataListener.On("OnDataTransfer", mock.Anything).Return(nil, nil)
	err = suite.csms.DataTransfer(wsId, func(response *data.DataTransferResponse, err error) {
		suite.Require().Nil(response)
		suite.Require().Error(err)
		resultChannel <- err
	}, "vendor1")
	suite.Require().Nil(err)
	result = <-resultChannel
	suite.Require().IsType(&ocpp.Error{}, result)
	ocppErr = result.(*ocpp.Error)
	suite.Equal(ocppj.GenericError, ocppErr.Code)
	suite.Equal("empty response to request 1234", ocppErr.Description)
}

func (suite *OcppV2TestSuite) TestCentralSystemSendResponseError() {
	wsId := "test_id"
	channel := NewMockWebSocket(wsId)
	var ocppErr *ocpp.Error
	var response *data.DataTransferResponse
	// Setup internal communication and listeners
	dataListener := &MockCSMSDataHandler{}
	suite.csms.SetDataHandler(dataListener)
	suite.mockWsClient.On("Start", mock.AnythingOfType("string")).Return(nil).Run(func(args mock.Arguments) {
		// Notify server of incoming connection
		suite.mockWsServer.NewClientHandler(channel)
	})
	suite.mockWsClient.On("Write", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		rawMsg := args.Get(0)
		bytes := rawMsg.([]byte)
		err := suite.mockWsServer.MessageHandler(channel, bytes)
		suite.Nil(err)
	})
	suite.mockWsServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	suite.mockWsServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		rawMsg := args.Get(1)
		bytes := rawMsg.([]byte)
		err := suite.mockWsClient.MessageHandler(bytes)
		suite.NoError(err)
	})
	// Run Tests
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start("someUrl")
	suite.Require().Nil(err)
	// Test 1: occurrence validation error
	dataTransferResponse := data.NewDataTransferResponse(data.DataTransferStatusAccepted)
	dataTransferResponse.Data = struct {
		Field1 string `validate:"required"`
	}{Field1: ""}
	dataListener.On("OnDataTransfer", mock.AnythingOfType("string"), mock.Anything).Return(dataTransferResponse, nil)
	response, err = suite.chargingStation.DataTransfer("vendor1")
	suite.Require().Nil(response)
	suite.Require().Error(err)
	suite.Require().IsType(&ocpp.Error{}, err)
	ocppErr = err.(*ocpp.Error)
	suite.Equal(ocppj.OccurrenceConstraintViolationV2, ocppErr.Code)
	suite.Equal("Field CallResult.Payload.Data.Field1 required but not found for feature DataTransfer", ocppErr.Description)
	// Test 2: marshaling error
	dataTransferResponse = data.NewDataTransferResponse(data.DataTransferStatusAccepted)
	dataTransferResponse.Data = make(chan struct{})
	dataListener.ExpectedCalls = nil
	dataListener.On("OnDataTransfer", mock.AnythingOfType("string"), mock.Anything).Return(dataTransferResponse, nil)
	response, err = suite.chargingStation.DataTransfer("vendor1")
	suite.Require().Nil(response)
	suite.Require().Error(err)
	suite.Require().IsType(&ocpp.Error{}, err)
	ocppErr = err.(*ocpp.Error)
	suite.Equal(ocppj.GenericError, ocppErr.Code)
	suite.Equal("json: unsupported type: chan struct {}", ocppErr.Description)
	// Test 3: no results in callback
	dataListener.ExpectedCalls = nil
	dataListener.On("OnDataTransfer", mock.AnythingOfType("string"), mock.Anything).Return(nil, nil)
	response, err = suite.chargingStation.DataTransfer("vendor1")
	suite.Require().Nil(response)
	suite.Require().Error(err)
	suite.Require().IsType(&ocpp.Error{}, err)
	ocppErr = err.(*ocpp.Error)
	suite.Equal(ocppj.GenericError, ocppErr.Code)
	suite.Equal(fmt.Sprintf("empty response to %s for request 1234", wsId), ocppErr.Description)
}

func (suite *OcppV2TestSuite) TestErrorCodes() {
	suite.Equal(ocppj.FormatViolationV2, ocppj.FormatErrorType(suite.ocppjServer))
	suite.Equal(ocppj.OccurrenceConstraintViolationV2, ocppj.OccurrenceConstraintErrorType(suite.ocppjServer))
}
