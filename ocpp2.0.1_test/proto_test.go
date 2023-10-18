package ocpp2_test

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/data"

	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocppj"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func (suite *OcppV2TestSuite) TestChargePointSendResponseError() {
	t := suite.T()
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
		assert.Nil(t, err)
	})
	suite.mockWsServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	suite.mockWsServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		rawMsg := args.Get(1)
		bytes := rawMsg.([]byte)
		err := suite.mockWsClient.MessageHandler(bytes)
		assert.NoError(t, err)
	})
	// Run Tests
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start("someUrl")
	require.Nil(t, err)
	resultChannel := make(chan error, 1)
	// Test 1: occurrence validation error
	dataTransferResponse := data.NewDataTransferResponse(data.DataTransferStatusAccepted)
	dataTransferResponse.Data = struct {
		Field1 string `validate:"required"`
	}{Field1: ""}
	dataListener.On("OnDataTransfer", mock.Anything).Return(dataTransferResponse, nil)
	err = suite.csms.DataTransfer(wsId, func(response *data.DataTransferResponse, err error) {
		require.Nil(t, response)
		require.Error(t, err)
		resultChannel <- err
	}, "vendor1")
	require.Nil(t, err)
	result := <-resultChannel
	require.IsType(t, &ocpp.Error{}, result)
	ocppErr = result.(*ocpp.Error)
	assert.Equal(t, ocppj.OccurrenceConstraintViolation, ocppErr.Code)
	assert.Equal(t, "Field CallResult.Payload.Data.Field1 required but not found for feature DataTransfer", ocppErr.Description)
	// Test 2: marshaling error
	dataTransferResponse = data.NewDataTransferResponse(data.DataTransferStatusAccepted)
	dataTransferResponse.Data = make(chan struct{})
	dataListener.ExpectedCalls = nil
	dataListener.On("OnDataTransfer", mock.Anything).Return(dataTransferResponse, nil)
	err = suite.csms.DataTransfer(wsId, func(response *data.DataTransferResponse, err error) {
		require.Nil(t, response)
		require.Error(t, err)
		resultChannel <- err
	}, "vendor1")
	require.Nil(t, err)
	result = <-resultChannel
	require.IsType(t, &ocpp.Error{}, result)
	ocppErr = result.(*ocpp.Error)
	assert.Equal(t, ocppj.GenericError, ocppErr.Code)
	assert.Equal(t, "json: unsupported type: chan struct {}", ocppErr.Description)
	// Test 3: no results in callback
	dataListener.ExpectedCalls = nil
	dataListener.On("OnDataTransfer", mock.Anything).Return(nil, nil)
	err = suite.csms.DataTransfer(wsId, func(response *data.DataTransferResponse, err error) {
		require.Nil(t, response)
		require.Error(t, err)
		resultChannel <- err
	}, "vendor1")
	require.Nil(t, err)
	result = <-resultChannel
	require.IsType(t, &ocpp.Error{}, result)
	ocppErr = result.(*ocpp.Error)
	assert.Equal(t, ocppj.GenericError, ocppErr.Code)
	assert.Equal(t, "empty response to request 1234", ocppErr.Description)
}

func (suite *OcppV2TestSuite) TestCentralSystemSendResponseError() {
	t := suite.T()
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
		assert.Nil(t, err)
	})
	suite.mockWsServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	suite.mockWsServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		rawMsg := args.Get(1)
		bytes := rawMsg.([]byte)
		err := suite.mockWsClient.MessageHandler(bytes)
		assert.NoError(t, err)
	})
	// Run Tests
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start("someUrl")
	require.Nil(t, err)
	// Test 1: occurrence validation error
	dataTransferResponse := data.NewDataTransferResponse(data.DataTransferStatusAccepted)
	dataTransferResponse.Data = struct {
		Field1 string `validate:"required"`
	}{Field1: ""}
	dataListener.On("OnDataTransfer", mock.AnythingOfType("string"), mock.Anything).Return(dataTransferResponse, nil)
	response, err = suite.chargingStation.DataTransfer("vendor1")
	require.Nil(t, response)
	require.Error(t, err)
	require.IsType(t, &ocpp.Error{}, err)
	ocppErr = err.(*ocpp.Error)
	assert.Equal(t, ocppj.OccurrenceConstraintViolation, ocppErr.Code)
	assert.Equal(t, "Field CallResult.Payload.Data.Field1 required but not found for feature DataTransfer", ocppErr.Description)
	// Test 2: marshaling error
	dataTransferResponse = data.NewDataTransferResponse(data.DataTransferStatusAccepted)
	dataTransferResponse.Data = make(chan struct{})
	dataListener.ExpectedCalls = nil
	dataListener.On("OnDataTransfer", mock.AnythingOfType("string"), mock.Anything).Return(dataTransferResponse, nil)
	response, err = suite.chargingStation.DataTransfer("vendor1")
	require.Nil(t, response)
	require.Error(t, err)
	require.IsType(t, &ocpp.Error{}, err)
	ocppErr = err.(*ocpp.Error)
	assert.Equal(t, ocppj.GenericError, ocppErr.Code)
	assert.Equal(t, "json: unsupported type: chan struct {}", ocppErr.Description)
	// Test 3: no results in callback
	dataListener.ExpectedCalls = nil
	dataListener.On("OnDataTransfer", mock.AnythingOfType("string"), mock.Anything).Return(nil, nil)
	response, err = suite.chargingStation.DataTransfer("vendor1")
	require.Nil(t, response)
	require.Error(t, err)
	require.IsType(t, &ocpp.Error{}, err)
	ocppErr = err.(*ocpp.Error)
	assert.Equal(t, ocppj.GenericError, ocppErr.Code)
	assert.Equal(t, fmt.Sprintf("empty response to %s for request 1234", wsId), ocppErr.Description)
}

func (suite *OcppV2TestSuite) TestErrorCodes() {
	t := suite.T()
	suite.mockWsServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	suite.csms.Start(8887, "somePath")
	assert.Equal(t, ocppj.FormatViolationV2, suite.csms.FormatError())
}
