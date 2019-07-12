package ocppj_test

import (
	"errors"
	"fmt"
	"github.com/lorenzodonini/go-ocpp/ocppj"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// SendRequest
func (suite *OcppJTestSuite) TestCentralSystemSendRequest() {
	mockChargePointId := "1234"
	suite.mockServer.On("Write", mockChargePointId, mock.Anything).Return(nil)
	mockRequest := newMockRequest("mockValue")
	err := suite.centralSystem.SendRequest(mockChargePointId, mockRequest)
	assert.Nil(suite.T(), err)
}

func (suite *OcppJTestSuite) TestCentralSystemSendInvalidRequest() {
	mockChargePointId := "1234"
	suite.mockServer.On("Write", mockChargePointId, mock.Anything).Return(nil)
	mockRequest := newMockRequest("")
	err := suite.centralSystem.SendRequest(mockChargePointId, mockRequest)
	assert.NotNil(suite.T(), err)
}

func (suite *OcppJTestSuite) TestCentralSystemSendRequestPending() {
	mockChargePointId := "1234"
	suite.mockServer.On("Write", mockChargePointId, mock.Anything).Return(nil)
	mockRequest := newMockRequest("mockValue")
	err := suite.centralSystem.SendRequest(mockChargePointId, mockRequest)
	assert.Nil(suite.T(), err)
	err = suite.centralSystem.SendRequest(mockChargePointId, mockRequest)
	assert.NotNil(suite.T(), err)
}

func (suite *OcppJTestSuite) TestCentralSystemSendRequestFailed() {
	mockChargePointId := "1234"
	suite.mockServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(errors.New("mockError"))
	mockRequest := newMockRequest("mockValue")
	err := suite.centralSystem.SendRequest(mockChargePointId, mockRequest)
	assert.NotNil(suite.T(), err)
}

// SendMessage
func (suite *OcppJTestSuite) TestCentralSystemSendMessage() {
	t := suite.T()
	mockChargePointId := "1234"
	suite.mockServer.On("Write", mockChargePointId, mock.Anything).Return(nil)
	mockRequest := newMockRequest("mockValue")
	mockCall, err := suite.centralSystem.CreateCall(mockRequest)
	assert.Nil(t, err)
	assert.NotNil(t, mockCall)
	err = suite.centralSystem.SendMessage(mockChargePointId, mockCall)
	assert.Nil(t, err)
}

func (suite *OcppJTestSuite) TestCentralSystemSendInvalidMessage() {
	t := suite.T()
	mockChargePointId := "1234"
	mockCallId := "6789"
	suite.mockServer.On("Write", mockChargePointId, mock.Anything).Return(nil)
	mockRequest := newMockRequest("mockValue")
	suite.centralSystem.AddPendingRequest(mockCallId, mockRequest)
	mockConfirmation := newMockConfirmation("")
	mockCallResult := ocppj.CallResult{
		MessageTypeId: ocppj.CALL_RESULT,
		UniqueId:      mockChargePointId,
		Payload:       mockConfirmation,
	}
	err := suite.centralSystem.SendMessage(mockChargePointId, &mockCallResult)
	assert.NotNil(t, err)
}

func (suite *OcppJTestSuite) TestCentralSystemSenddessageFailed() {
	t := suite.T()
	mockChargePointId := "1234"
	mockCallId := "6789"
	suite.mockServer.On("Write", mockChargePointId, mock.Anything).Return(errors.New("mockError"))
	mockRequest := newMockRequest("mockValue")
	suite.centralSystem.AddPendingRequest(mockCallId, mockRequest)
	mockConfirmation := newMockConfirmation("mockValue")
	mockCallResult := ocppj.CallResult{
		MessageTypeId: ocppj.CALL_RESULT,
		UniqueId:      mockChargePointId,
		Payload:       mockConfirmation,
	}
	err := suite.centralSystem.SendMessage(mockChargePointId, &mockCallResult)
	assert.NotNil(t, err)
}

func (suite *OcppJTestSuite) TestCentralSystemCallHandler() {
	t := suite.T()
	mockChargePointId := "1234"
	mockUniqueId := "5678"
	mockValue := "someValue"
	mockRequest := fmt.Sprintf(`[2,"%v","%v",{"mockValue":"%v"}]`, mockUniqueId, MockFeatureName, mockValue)
	suite.centralSystem.SetCallHandler(func(chargePointId string, call *ocppj.Call) {
		assert.Equal(t, mockChargePointId, chargePointId)
		CheckCall(call, t, MockFeatureName, mockUniqueId)
	})
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return().Run(func(args mock.Arguments) {
		// Simulate charge point message
		channel := NewMockWebSocket(mockChargePointId)
		err := suite.mockServer.MessageHandler(channel, []byte(mockRequest))
		assert.Nil(t, err)
	})
	suite.centralSystem.Start(8887, "somePath")
}

func (suite *OcppJTestSuite) TestCentralSystemCallResultHandler() {
	t := suite.T()
	mockChargePointId := "1234"
	mockUniqueId := "5678"
	mockValue := "someValue"
	mockRequest := newMockRequest("testValue")
	mockConfirmation := fmt.Sprintf(`[3,"%v",{"mockValue":"%v"}]`, mockUniqueId, mockValue)
	suite.centralSystem.SetCallResultHandler(func(chargePointId string, callResult *ocppj.CallResult) {
		assert.Equal(t, mockChargePointId, chargePointId)
		CheckCallResult(callResult, t, mockUniqueId)
	})
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return().Run(func(args mock.Arguments) {
		// Simulate charge point message
		channel := NewMockWebSocket(mockChargePointId)
		err := suite.mockServer.MessageHandler(channel, []byte(mockConfirmation))
		assert.Nil(t, err)
	})
	suite.centralSystem.AddPendingRequest(mockUniqueId, mockRequest)
	suite.centralSystem.Start(8887, "somePath")
}

func (suite *OcppJTestSuite) TestCentralSystemCallErrorHandler() {
	t := suite.T()
	mockChargePointId := "1234"
	mockUniqueId := "5678"
	mockErrorCode := ocppj.GenericError
	mockErrorDescription := "Mock Description"
	mockValue := "someValue"
	mockErrorDetails := make(map[string]interface{})
	mockErrorDetails["details"] = "someValue"

	mockRequest := newMockRequest("testValue")
	mockError := fmt.Sprintf(`[4,"%v","%v","%v",{"details":"%v"}]`, mockUniqueId, mockErrorCode, mockErrorDescription, mockValue)
	suite.centralSystem.SetCallErrorHandler(func(chargePointId string, callError *ocppj.CallError) {
		assert.Equal(t, mockChargePointId, chargePointId)
		CheckCallError(t, callError, mockUniqueId, mockErrorCode, mockErrorDescription, mockErrorDetails)
	})
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return().Run(func(args mock.Arguments) {
		// Simulate charge point message
		channel := NewMockWebSocket(mockChargePointId)
		err := suite.mockServer.MessageHandler(channel, []byte(mockError))
		assert.Nil(t, err)
	})
	suite.centralSystem.AddPendingRequest(mockUniqueId, mockRequest)
	suite.centralSystem.Start(8887, "somePath")
}
