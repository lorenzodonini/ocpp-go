package ocpp16_test

import (
	"fmt"
	"testing"

	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocppj"
	"github.com/stretchr/testify/mock"
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
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start("someUrl")
	suite.Require().Nil(err)
	resultChannel := make(chan error, 1)

	testCases := []struct {
		name        string
		confirmData interface{}
		expectedErr *ocpp.Error
	}{
		{
			name:        "ocurrence validation",
			confirmData: CustomData{Field1: "", Field2: 42},
			expectedErr: &ocpp.Error{Code: ocppj.OccurrenceConstraintViolationV16, Description: "Field CallResult.Payload.Data.Field1 required but not found for feature DataTransfer"},
		},
		{
			name:        "marshaling error",
			confirmData: make(chan struct{}),
			expectedErr: &ocpp.Error{Code: ocppj.GenericError, Description: "json: unsupported type: chan struct {}"},
		},
		{
			name:        "empty confirmation",
			confirmData: nil,
			expectedErr: &ocpp.Error{Code: ocppj.GenericError, Description: "empty confirmation to request 1234"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(*testing.T) {
			coreListener.ExpectedCalls = nil
			if tc.confirmData != nil {
				dataTransferConfirmation := core.NewDataTransferConfirmation(core.DataTransferStatusAccepted)
				dataTransferConfirmation.Data = tc.confirmData
				coreListener.On("OnDataTransfer", mock.Anything).Return(dataTransferConfirmation, nil)
			} else {
				coreListener.On("OnDataTransfer", mock.Anything).Return(nil, nil)
			}

			err = suite.centralSystem.DataTransfer(wsId, func(confirmation *core.DataTransferConfirmation, err error) {
				suite.Require().Nil(confirmation)
				suite.Require().Error(err)
				resultChannel <- err
			}, "vendor1")
			suite.Require().Nil(err)
			result := <-resultChannel
			suite.Require().IsType(&ocpp.Error{}, result)
			ocppErr = result.(*ocpp.Error)
			suite.Equal(tc.expectedErr.Code, ocppErr.Code)
			suite.Equal(tc.expectedErr.Description, ocppErr.Description)
		})
	}
}

func (suite *OcppV16TestSuite) TestCentralSystemSendResponseError() {
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
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start("someUrl")
	suite.Require().Nil(err)
	// Test 1: occurrence validation error
	dataTransferConfirmation := core.NewDataTransferConfirmation(core.DataTransferStatusAccepted)
	dataTransferConfirmation.Data = CustomData{Field1: "", Field2: 42}
	coreListener.On("OnDataTransfer", mock.AnythingOfType("string"), mock.Anything).Return(dataTransferConfirmation, nil)
	response, err = suite.chargePoint.DataTransfer("vendor1")
	suite.Require().Nil(response)
	suite.Require().Error(err)
	suite.Require().IsType(&ocpp.Error{}, err)
	ocppErr = err.(*ocpp.Error)
	suite.Equal(ocppj.OccurrenceConstraintViolationV16, ocppErr.Code)
	suite.Equal("Field CallResult.Payload.Data.Field1 required but not found for feature DataTransfer", ocppErr.Description)
	// Test 2: marshaling error
	dataTransferConfirmation = core.NewDataTransferConfirmation(core.DataTransferStatusAccepted)
	dataTransferConfirmation.Data = make(chan struct{})
	coreListener.ExpectedCalls = nil
	coreListener.On("OnDataTransfer", mock.AnythingOfType("string"), mock.Anything).Return(dataTransferConfirmation, nil)
	response, err = suite.chargePoint.DataTransfer("vendor1")
	suite.Require().Nil(response)
	suite.Require().Error(err)
	suite.Require().IsType(&ocpp.Error{}, err)
	ocppErr = err.(*ocpp.Error)
	suite.Equal(ocppj.GenericError, ocppErr.Code)
	suite.Equal("json: unsupported type: chan struct {}", ocppErr.Description)
	// Test 3: no results in callback
	coreListener.ExpectedCalls = nil
	coreListener.On("OnDataTransfer", mock.AnythingOfType("string"), mock.Anything).Return(nil, nil)
	response, err = suite.chargePoint.DataTransfer("vendor1")
	suite.Require().Nil(response)
	suite.Require().Error(err)
	suite.Require().IsType(&ocpp.Error{}, err)
	ocppErr = err.(*ocpp.Error)
	suite.Equal(ocppj.GenericError, ocppErr.Code)
	suite.Equal(fmt.Sprintf("empty confirmation to %s for request 1234", wsId), ocppErr.Description)
}

func (suite *OcppV16TestSuite) TestErrorCodes() {
	suite.Equal(ocppj.FormatViolationV16, ocppj.FormatErrorType(suite.ocppjCentralSystem))
	suite.Equal(ocppj.OccurrenceConstraintViolationV16, ocppj.OccurrenceConstraintErrorType(suite.ocppjCentralSystem))
}
