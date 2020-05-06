package ocppj_test

import (
	"errors"
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocppj"
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

// SendResponse
func (suite *OcppJTestSuite) TestCentralSystemSendConfirmation() {
	t := suite.T()
	mockChargePointId := "0101"
	mockUniqueId := "1234"
	suite.mockServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(nil)
	mockConfirmation := newMockConfirmation("mockValue")
	err := suite.centralSystem.SendResponse(mockChargePointId, mockUniqueId, mockConfirmation)
	assert.Nil(t, err)
}

func (suite *OcppJTestSuite) TestCentralSystemSendInvalidConfirmation() {
	t := suite.T()
	mockChargePointId := "0101"
	mockUniqueId := "6789"
	suite.mockServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(nil)
	mockConfirmation := newMockConfirmation("")
	// This is allowed. Endpoint doesn't keep track of incoming requests, but only outgoing ones
	err := suite.centralSystem.SendResponse(mockChargePointId, mockUniqueId, mockConfirmation)
	assert.NotNil(t, err)
}

func (suite *OcppJTestSuite) TestCentralSystemSendConfirmationFailed() {
	t := suite.T()
	mockChargePointId := "0101"
	mockUniqueId := "1234"
	suite.mockServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(errors.New("networkError"))
	mockConfirmation := newMockConfirmation("mockValue")
	err := suite.centralSystem.SendResponse(mockChargePointId, mockUniqueId, mockConfirmation)
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
	err := suite.centralSystem.SendResponse(mockChargePointId, mockUniqueId, mockConfirmation)
	assert.NotNil(t, err)
}

// Call Handlers
func (suite *OcppJTestSuite) TestCentralSystemRequestHandler() {
	t := suite.T()
	mockChargePointId := "1234"
	mockUniqueId := "5678"
	mockValue := "someValue"
	mockRequest := fmt.Sprintf(`[2,"%v","%v",{"mockValue":"%v"}]`, mockUniqueId, MockFeatureName, mockValue)
	suite.centralSystem.SetRequestHandler(func(chargePointId string, request ocpp.Request, requestId string, action string) {
		assert.Equal(t, mockChargePointId, chargePointId)
		assert.Equal(t, mockUniqueId, requestId)
		assert.Equal(t, MockFeatureName, action)
		assert.NotNil(t, request)
	})
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return().Run(func(args mock.Arguments) {
		// Simulate charge point message
		channel := NewMockWebSocket(mockChargePointId)
		err := suite.mockServer.MessageHandler(channel, []byte(mockRequest))
		assert.Nil(t, err)
	})
	suite.centralSystem.Start(8887, "somePath")
}

func (suite *OcppJTestSuite) TestCentralSystemConfirmationHandler() {
	t := suite.T()
	mockChargePointId := "1234"
	mockUniqueId := "5678"
	mockValue := "someValue"
	mockRequest := newMockRequest("testValue")
	mockConfirmation := fmt.Sprintf(`[3,"%v",{"mockValue":"%v"}]`, mockUniqueId, mockValue)
	suite.centralSystem.SetResponseHandler(func(chargePointId string, confirmation ocpp.Response, requestId string) {
		assert.Equal(t, mockChargePointId, chargePointId)
		assert.Equal(t, mockUniqueId, requestId)
		assert.NotNil(t, confirmation)
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

func (suite *OcppJTestSuite) TestCentralSystemErrorHandler() {
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
	suite.centralSystem.SetErrorHandler(func(chargePointId string, err *ocpp.Error, details interface{}) {
		assert.Equal(t, mockChargePointId, chargePointId)
		assert.Equal(t, mockUniqueId, err.MessageId)
		assert.Equal(t, mockErrorCode, err.Code)
		assert.Equal(t, mockErrorDescription, err.Description)
		assert.Equal(t, mockErrorDetails, details)
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
