package ocpp16_test

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocppj"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func (suite *OcppV16TestSuite) TestChargePointSendResponseError() {
	t := suite.T()
	wsId := "test_id"
	channel := NewMockWebSocket(wsId)
	var ocppErr *ocpp.Error
	// Setup internal communication and listeners
	coreListener := &MockChargePointCoreListener{}
	suite.chargePoint.SetCoreHandler(coreListener)
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
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start("someUrl")
	require.Nil(t, err)
	resultChannel := make(chan error, 1)
	// Test 1: occurrence validation error
	dataTransferConfirmation := core.NewDataTransferConfirmation(core.DataTransferStatusAccepted)
	dataTransferConfirmation.Data = CustomData{Field1: "", Field2: 42}
	coreListener.On("OnDataTransfer", mock.Anything).Return(dataTransferConfirmation, nil)
	err = suite.centralSystem.DataTransfer(wsId, func(confirmation *core.DataTransferConfirmation, err error) {
		require.Nil(t, confirmation)
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
	dataTransferConfirmation = core.NewDataTransferConfirmation(core.DataTransferStatusAccepted)
	dataTransferConfirmation.Data = make(chan struct{})
	coreListener.ExpectedCalls = nil
	coreListener.On("OnDataTransfer", mock.Anything).Return(dataTransferConfirmation, nil)
	err = suite.centralSystem.DataTransfer(wsId, func(confirmation *core.DataTransferConfirmation, err error) {
		require.Nil(t, confirmation)
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
	coreListener.ExpectedCalls = nil
	coreListener.On("OnDataTransfer", mock.Anything).Return(nil, nil)
	err = suite.centralSystem.DataTransfer(wsId, func(confirmation *core.DataTransferConfirmation, err error) {
		require.Nil(t, confirmation)
		require.Error(t, err)
		resultChannel <- err
	}, "vendor1")
	require.Nil(t, err)
	result = <-resultChannel
	require.IsType(t, &ocpp.Error{}, result)
	ocppErr = result.(*ocpp.Error)
	assert.Equal(t, ocppj.GenericError, ocppErr.Code)
	assert.Equal(t, "empty confirmation to request 1234", ocppErr.Description)
}

func (suite *OcppV16TestSuite) TestCentralSystemSendResponseError() {
	t := suite.T()
	wsId := "test_id"
	channel := NewMockWebSocket(wsId)
	var ocppErr *ocpp.Error
	var response *core.DataTransferConfirmation
	// Setup internal communication and listeners
	coreListener := &MockCentralSystemCoreListener{}
	suite.centralSystem.SetCoreHandler(coreListener)
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
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start("someUrl")
	require.Nil(t, err)
	// Test 1: occurrence validation error
	dataTransferConfirmation := core.NewDataTransferConfirmation(core.DataTransferStatusAccepted)
	dataTransferConfirmation.Data = CustomData{Field1: "", Field2: 42}
	coreListener.On("OnDataTransfer", mock.AnythingOfType("string"), mock.Anything).Return(dataTransferConfirmation, nil)
	response, err = suite.chargePoint.DataTransfer("vendor1")
	require.Nil(t, response)
	require.Error(t, err)
	require.IsType(t, &ocpp.Error{}, err)
	ocppErr = err.(*ocpp.Error)
	assert.Equal(t, ocppj.OccurrenceConstraintViolation, ocppErr.Code)
	assert.Equal(t, "Field CallResult.Payload.Data.Field1 required but not found for feature DataTransfer", ocppErr.Description)
	// Test 2: marshaling error
	dataTransferConfirmation = core.NewDataTransferConfirmation(core.DataTransferStatusAccepted)
	dataTransferConfirmation.Data = make(chan struct{})
	coreListener.ExpectedCalls = nil
	coreListener.On("OnDataTransfer", mock.AnythingOfType("string"), mock.Anything).Return(dataTransferConfirmation, nil)
	response, err = suite.chargePoint.DataTransfer("vendor1")
	require.Nil(t, response)
	require.Error(t, err)
	require.IsType(t, &ocpp.Error{}, err)
	ocppErr = err.(*ocpp.Error)
	assert.Equal(t, ocppj.GenericError, ocppErr.Code)
	assert.Equal(t, "json: unsupported type: chan struct {}", ocppErr.Description)
	// Test 3: no results in callback
	coreListener.ExpectedCalls = nil
	coreListener.On("OnDataTransfer", mock.AnythingOfType("string"), mock.Anything).Return(nil, nil)
	response, err = suite.chargePoint.DataTransfer("vendor1")
	require.Nil(t, response)
	require.Error(t, err)
	require.IsType(t, &ocpp.Error{}, err)
	ocppErr = err.(*ocpp.Error)
	assert.Equal(t, ocppj.GenericError, ocppErr.Code)
	assert.Equal(t, fmt.Sprintf("empty confirmation to %s for request 1234", wsId), ocppErr.Description)
}

func (suite *OcppV16TestSuite) TestErrorCodes() {
	t := suite.T()
	suite.mockWsServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	suite.centralSystem.Start(8887, "somePath")
	assert.Equal(t, ocppj.FormatViolationV16, ocppj.FormatErrorType(suite.centralSystem))
}
