package test

import (
	"errors"
	"fmt"
	"github.com/lorenzodonini/go-ocpp/ocpp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Start
func (suite *OcppJTestSuite) TestChargePointStart() {
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	err := suite.chargePoint.Start("someUrl")
	assert.Nil(suite.T(), err)
}

func (suite *OcppJTestSuite) TestChargePointStartFailed() {
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(errors.New("startError"))
	err := suite.chargePoint.Start("someUrl")
	assert.NotNil(suite.T(), err)
}

// SendRequest
func (suite *OcppJTestSuite) TestChargePointSendRequest() {
	suite.mockClient.On("Write", mock.Anything).Return(nil)
	mockRequest := newMockRequest("mockValue")
	err := suite.chargePoint.SendRequest(mockRequest)
	assert.Nil(suite.T(), err)
}

func (suite *OcppJTestSuite) TestChargePointSendInvalidRequest() {
	suite.mockClient.On("Write", mock.Anything).Return(nil)
	mockRequest := newMockRequest("")
	err := suite.chargePoint.SendRequest(mockRequest)
	assert.NotNil(suite.T(), err)
}

func (suite *OcppJTestSuite) TestChargePointSendRequestPending() {
	suite.mockClient.On("Write", mock.Anything).Return(nil)
	mockRequest := newMockRequest("mockValue")
	err := suite.chargePoint.SendRequest(mockRequest)
	assert.Nil(suite.T(), err)
	err = suite.chargePoint.SendRequest(mockRequest)
	assert.NotNil(suite.T(), err)
}

func (suite *OcppJTestSuite) TestChargePointSendRequestFailed() {
	suite.mockClient.On("Write", mock.Anything).Return(errors.New("mockError"))
	mockRequest := newMockRequest("mockValue")
	err := suite.chargePoint.SendRequest(mockRequest)
	assert.NotNil(suite.T(), err)
}

// SendMessage
func (suite *OcppJTestSuite) TestChargePointSendMessage() {
	t := suite.T()
	suite.mockClient.On("Write", mock.Anything).Return(nil)
	mockRequest := newMockRequest("mockValue")
	mockCall, err := suite.chargePoint.CreateCall(mockRequest)
	assert.Nil(t, err)
	assert.NotNil(t, mockCall)
	err = suite.chargePoint.SendMessage(mockCall)
	assert.Nil(t, err)
}

func (suite *OcppJTestSuite) TestChargePointSendInvalidMessage() {
	t := suite.T()
	mockCallId := "6789"
	mockInvalidId := "0000"
	suite.mockClient.On("Write", mock.Anything).Return(nil)
	mockRequest := newMockRequest("mockValue")
	suite.chargePoint.AddPendingRequest(mockCallId, mockRequest)
	mockConfirmation := newMockConfirmation("")
	mockCallResult := ocpp.CallResult{
		MessageTypeId: ocpp.CALL_RESULT,
		UniqueId:      mockInvalidId,
		Payload:       mockConfirmation,
	}
	err := suite.chargePoint.SendMessage(&mockCallResult)
	assert.NotNil(t, err)
}

func (suite *OcppJTestSuite) TestChargePointSendMessageFailed() {
	t := suite.T()
	mockCallId := "6789"
	suite.mockClient.On("Write", mock.Anything).Return(errors.New("mockError"))
	mockRequest := newMockRequest("mockValue")
	suite.chargePoint.AddPendingRequest(mockCallId, mockRequest)
	mockConfirmation := newMockConfirmation("mockValue")
	mockCallResult := ocpp.CallResult{
		MessageTypeId: ocpp.CALL_RESULT,
		UniqueId:      mockCallId,
		Payload:       mockConfirmation,
	}
	err := suite.chargePoint.SendMessage(&mockCallResult)
	assert.NotNil(t, err)
}

func (suite *OcppJTestSuite) TestChargePointCallHandler() {
	t := suite.T()
	mockUniqueId := "5678"
	mockValue := "someValue"
	mockRequest := fmt.Sprintf(`[2,"%v","%v",{"mockValue":"%v"}]`, mockUniqueId, MockFeatureName, mockValue)
	suite.chargePoint.SetCallHandler(func(call *ocpp.Call) {
		CheckCall(call, t, MockFeatureName, mockUniqueId)
	})
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil).Run(func(args mock.Arguments) {
		// Simulate central system message
		err := suite.mockClient.MessageHandler([]byte(mockRequest))
		assert.Nil(t, err)
	})
	err := suite.chargePoint.Start( "somePath")
	assert.Nil(t, err)
}

func (suite *OcppJTestSuite) TestChargePointCallResultHandler() {
	t := suite.T()
	mockUniqueId := "5678"
	mockValue := "someValue"
	mockRequest := newMockRequest("testValue")
	mockConfirmation := fmt.Sprintf(`[3,"%v",{"mockValue":"%v"}]`, mockUniqueId, mockValue)
	suite.chargePoint.SetCallResultHandler(func(callResult *ocpp.CallResult) {
		CheckCallResult(callResult, t, mockUniqueId)
	})
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil).Run(func(args mock.Arguments) {
		// Simulate central system message
		err := suite.mockClient.MessageHandler([]byte(mockConfirmation))
		assert.Nil(t, err)
	})
	suite.chargePoint.AddPendingRequest(mockUniqueId, mockRequest)
	err := suite.chargePoint.Start( "somePath")
	assert.Nil(t, err)
}

func (suite *OcppJTestSuite) TestChargePointCallErrorHandler() {
	t := suite.T()
	mockUniqueId := "5678"
	mockErrorCode := ocpp.GenericError
	mockErrorDescription := "Mock Description"
	mockValue := "someValue"
	mockErrorDetails := make(map[string]interface{})
	mockErrorDetails["details"] = "someValue"

	mockRequest := newMockRequest("testValue")
	mockError := fmt.Sprintf(`[4,"%v","%v","%v",{"details":"%v"}]`, mockUniqueId, mockErrorCode, mockErrorDescription, mockValue)
	suite.chargePoint.SetCallErrorHandler(func(callError *ocpp.CallError) {
		CheckCallError(t, callError, mockUniqueId, mockErrorCode, mockErrorDescription, mockErrorDetails)
	})
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil).Run(func(args mock.Arguments) {
		// Simulate central system message
		err := suite.mockClient.MessageHandler([]byte(mockError))
		assert.Nil(t, err)
	})
	suite.chargePoint.AddPendingRequest(mockUniqueId, mockRequest)
	err := suite.chargePoint.Start( "someUrl")
	assert.Nil(t, err)
}
