package ocppj_test

import (
	"errors"
	"fmt"
	"github.com/lorenzodonini/go-ocpp/ocpp"
	"github.com/lorenzodonini/go-ocpp/ocppj"
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
	suite.mockClient.On("Write", mock.Anything).Return(errors.New("networkError"))
	mockRequest := newMockRequest("mockValue")
	err := suite.chargePoint.SendRequest(mockRequest)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "networkError", err.Error())
	//TODO: assert that pending request was removed
}

// SendConfirmation
func (suite *OcppJTestSuite) TestChargePointSendConfirmation() {
	t := suite.T()
	mockUniqueId := "1234"
	suite.mockClient.On("Write", mock.Anything).Return(nil)
	mockConfirmation := newMockConfirmation("mockValue")
	err := suite.chargePoint.SendConfirmation(mockUniqueId, mockConfirmation)
	assert.Nil(t, err)
}

func (suite *OcppJTestSuite) TestChargePointSendInvalidConfirmation() {
	t := suite.T()
	mockUniqueId := "6789"
	suite.mockClient.On("Write", mock.Anything).Return(nil)
	mockConfirmation := newMockConfirmation("")
	// This is allowed. Endpoint doesn't keep track of incoming requests, but only outgoing ones
	err := suite.chargePoint.SendConfirmation(mockUniqueId, mockConfirmation)
	assert.NotNil(t, err)
}

func (suite *OcppJTestSuite) TestChargePointSendConfirmationFailed() {
	t := suite.T()
	mockUniqueId := "1234"
	suite.mockClient.On("Write", mock.Anything).Return(errors.New("networkError"))
	mockConfirmation := newMockConfirmation("mockValue")
	err := suite.chargePoint.SendConfirmation(mockUniqueId, mockConfirmation)
	assert.NotNil(t, err)
	assert.Equal(t, "networkError", err.Error())
}

// SendError
func (suite *OcppJTestSuite) TestChargePointSendError() {
	t := suite.T()
	mockUniqueId := "1234"
	mockDescription := "mockDescription"
	suite.mockClient.On("Write", mock.Anything).Return(nil)
	err := suite.chargePoint.SendError(mockUniqueId, ocppj.GenericError, mockDescription, nil)
	assert.Nil(t, err)
}

func (suite *OcppJTestSuite) TestChargePointSendInvalidError() {
	t := suite.T()
	mockUniqueId := "6789"
	mockDescription := "mockDescription"
	suite.mockClient.On("Write", mock.Anything).Return(nil)
	err := suite.chargePoint.SendError(mockUniqueId, "InvalidErrorCode", mockDescription, nil)
	assert.NotNil(t, err)
}

func (suite *OcppJTestSuite) TestChargePointSendErrorFailed() {
	t := suite.T()
	mockUniqueId := "1234"
	suite.mockClient.On("Write", mock.Anything).Return(errors.New("networkError"))
	mockConfirmation := newMockConfirmation("mockValue")
	err := suite.chargePoint.SendConfirmation(mockUniqueId, mockConfirmation)
	assert.NotNil(t, err)
	assert.Equal(t, "networkError", err.Error())
}

// Call Handlers
func (suite *OcppJTestSuite) TestChargePointCallHandler() {
	t := suite.T()
	mockUniqueId := "5678"
	mockValue := "someValue"
	mockRequest := fmt.Sprintf(`[2,"%v","%v",{"mockValue":"%v"}]`, mockUniqueId, MockFeatureName, mockValue)
	suite.chargePoint.SetRequestHandler(func(request ocpp.Request, requestId string, action string) {
		assert.Equal(t, mockUniqueId, requestId)
		assert.Equal(t, MockFeatureName, action)
		assert.NotNil(t, request)
	})
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil).Run(func(args mock.Arguments) {
		// Simulate central system message
		err := suite.mockClient.MessageHandler([]byte(mockRequest))
		assert.Nil(t, err)
	})
	err := suite.chargePoint.Start("somePath")
	assert.Nil(t, err)
}

func (suite *OcppJTestSuite) TestChargePointCallResultHandler() {
	t := suite.T()
	mockUniqueId := "5678"
	mockValue := "someValue"
	mockRequest := newMockRequest("testValue")
	mockConfirmation := fmt.Sprintf(`[3,"%v",{"mockValue":"%v"}]`, mockUniqueId, mockValue)
	suite.chargePoint.SetConfirmationHandler(func(confirmation ocpp.Confirmation, requestId string) {
		assert.Equal(t, mockUniqueId, requestId)
		assert.NotNil(t, confirmation)
	})
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil).Run(func(args mock.Arguments) {
		// Simulate central system message
		err := suite.mockClient.MessageHandler([]byte(mockConfirmation))
		assert.Nil(t, err)
	})
	suite.chargePoint.AddPendingRequest(mockUniqueId, mockRequest)
	err := suite.chargePoint.Start("somePath")
	assert.Nil(t, err)
}

func (suite *OcppJTestSuite) TestChargePointCallErrorHandler() {
	t := suite.T()
	mockUniqueId := "5678"
	mockErrorCode := ocppj.GenericError
	mockErrorDescription := "Mock Description"
	mockValue := "someValue"
	mockErrorDetails := make(map[string]interface{})
	mockErrorDetails["details"] = "someValue"

	mockRequest := newMockRequest("testValue")
	mockError := fmt.Sprintf(`[4,"%v","%v","%v",{"details":"%v"}]`, mockUniqueId, mockErrorCode, mockErrorDescription, mockValue)
	suite.chargePoint.SetErrorHandler(func(errorCode ocppj.ErrorCode, description string, details interface{}, requestId string) {
		assert.Equal(t, mockUniqueId, requestId)
		assert.Equal(t, mockErrorCode, errorCode)
		assert.Equal(t, mockErrorDescription, description)
		assert.Equal(t, mockErrorDetails, details)
	})
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil).Run(func(args mock.Arguments) {
		// Simulate central system message
		err := suite.mockClient.MessageHandler([]byte(mockError))
		assert.Nil(t, err)
	})
	suite.chargePoint.AddPendingRequest(mockUniqueId, mockRequest)
	err := suite.chargePoint.Start("someUrl")
	assert.Nil(t, err)
}
