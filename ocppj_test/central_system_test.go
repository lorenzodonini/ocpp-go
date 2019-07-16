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
	suite.mockServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(errors.New("networkError"))
	mockRequest := newMockRequest("mockValue")
	err := suite.centralSystem.SendRequest(mockChargePointId, mockRequest)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "networkError", err.Error())
	//TODO: assert that pending request was removed
}

// SendConfirmation
func (suite *OcppJTestSuite) TestCentralSystemSendConfirmation() {
	t := suite.T()
	mockChargePointId := "0101"
	mockUniqueId := "1234"
	suite.mockServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(nil)
	mockConfirmation := newMockConfirmation("mockValue")
	err := suite.centralSystem.SendConfirmation(mockChargePointId, mockUniqueId, mockConfirmation)
	assert.Nil(t, err)
}

func (suite *OcppJTestSuite) TestCentralSystemSendInvalidConfirmation() {
	t := suite.T()
	mockChargePointId := "0101"
	mockUniqueId := "6789"
	suite.mockServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(nil)
	mockConfirmation := newMockConfirmation("")
	// This is allowed. Endpoint doesn't keep track of incoming requests, but only outgoing ones
	err := suite.centralSystem.SendConfirmation(mockChargePointId, mockUniqueId, mockConfirmation)
	assert.NotNil(t, err)
}

func (suite *OcppJTestSuite) TestCentralSystemSendConfirmationFailed() {
	t := suite.T()
	mockChargePointId := "0101"
	mockUniqueId := "1234"
	suite.mockServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(errors.New("networkError"))
	mockConfirmation := newMockConfirmation("mockValue")
	err := suite.centralSystem.SendConfirmation(mockChargePointId, mockUniqueId, mockConfirmation)
	assert.NotNil(t, err)
	assert.Equal(t, "networkError", err.Error())
}

// SendError
func (suite *OcppJTestSuite) TestCentralSystemSendError() {
	t := suite.T()
	mockChargePointId := "0101"
	mockUniqueId := "1234"
	mockDescription := "mockDescription"
	suite.mockServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(nil)
	err := suite.centralSystem.SendError(mockChargePointId, mockUniqueId, ocppj.GenericError, mockDescription, nil)
	assert.Nil(t, err)
}

func (suite *OcppJTestSuite) TestCentralSystemSendInvalidError() {
	t := suite.T()
	mockChargePointId := "0101"
	mockUniqueId := "6789"
	mockDescription := "mockDescription"
	suite.mockServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(nil)
	err := suite.centralSystem.SendError(mockChargePointId, mockUniqueId, "InvalidErrorCode", mockDescription, nil)
	assert.NotNil(t, err)
}

func (suite *OcppJTestSuite) TestCentralSystemSendErrorFailed() {
	t := suite.T()
	mockChargePointId := "0101"
	mockUniqueId := "1234"
	suite.mockServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(errors.New("networkError"))
	mockConfirmation := newMockConfirmation("mockValue")
	err := suite.centralSystem.SendConfirmation(mockChargePointId, mockUniqueId, mockConfirmation)
	assert.NotNil(t, err)
}

// Call Handlers
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
