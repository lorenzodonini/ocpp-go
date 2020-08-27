package ocppj_test

import (
	"errors"
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocppj"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"time"
)

// SendRequest
func (suite *OcppJTestSuite) TestCentralSystemSendRequest() {
	mockChargePointId := "1234"
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	suite.mockServer.On("Write", mockChargePointId, mock.Anything).Return(nil)
	suite.centralSystem.Start(8887, "/{ws}")
	mockRequest := newMockRequest("mockValue")
	err := suite.centralSystem.SendRequest(mockChargePointId, mockRequest)
	assert.Nil(suite.T(), err)
}

func (suite *OcppJTestSuite) TestCentralSystemSendInvalidRequest() {
	mockChargePointId := "1234"
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	suite.mockServer.On("Write", mockChargePointId, mock.Anything).Return(nil)
	suite.centralSystem.Start(8887, "/{ws}")
	mockRequest := newMockRequest("")
	err := suite.centralSystem.SendRequest(mockChargePointId, mockRequest)
	assert.NotNil(suite.T(), err)
}

func (suite *OcppJTestSuite) TestCentralSystemSendRequestFailed() {
	t := suite.T()
	mockChargePointId := "1234"
	var callID string
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	suite.mockServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(errors.New("networkError")).Run(func(args mock.Arguments) {
		clientID := args.String(0)
		q, ok := suite.serverRequestMap.Get(clientID)
		require.True(t, ok)
		require.False(t, q.IsEmpty())
		req := q.Peek().(ocppj.RequestBundle)
		callID = req.Call.GetUniqueId()
		_, ok = suite.centralSystem.GetPendingRequest(callID)
		// Before anything is returned, the request must still be pending
		assert.True(t, ok)
	})
	suite.centralSystem.Start(8887, "/{ws}")
	mockRequest := newMockRequest("mockValue")
	err := suite.centralSystem.SendRequest(mockChargePointId, mockRequest)
	//TODO: currently the network error is not returned by SendRequest, but is only generated internally
	assert.Nil(t, err)
	// Assert that pending request was removed
	time.Sleep(500 * time.Millisecond)
	_, ok := suite.centralSystem.GetPendingRequest(callID)
	assert.False(t, ok)
}

// SendResponse
func (suite *OcppJTestSuite) TestCentralSystemSendConfirmation() {
	t := suite.T()
	mockChargePointId := "0101"
	mockUniqueId := "1234"
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	suite.mockServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(nil)
	suite.centralSystem.Start(8887, "/{ws}")
	mockConfirmation := newMockConfirmation("mockValue")
	err := suite.centralSystem.SendResponse(mockChargePointId, mockUniqueId, mockConfirmation)
	assert.Nil(t, err)
}

func (suite *OcppJTestSuite) TestCentralSystemSendInvalidConfirmation() {
	t := suite.T()
	mockChargePointId := "0101"
	mockUniqueId := "6789"
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	suite.mockServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(nil)
	suite.centralSystem.Start(8887, "/{ws}")
	mockConfirmation := newMockConfirmation("")
	// This is allowed. Endpoint doesn't keep track of incoming requests, but only outgoing ones
	err := suite.centralSystem.SendResponse(mockChargePointId, mockUniqueId, mockConfirmation)
	assert.NotNil(t, err)
}

func (suite *OcppJTestSuite) TestCentralSystemSendConfirmationFailed() {
	t := suite.T()
	mockChargePointId := "0101"
	mockUniqueId := "1234"
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	suite.mockServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(errors.New("networkError"))
	suite.centralSystem.Start(8887, "/{ws}")
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
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	suite.mockServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(nil)
	suite.centralSystem.Start(8887, "/{ws}")
	err := suite.centralSystem.SendError(mockChargePointId, mockUniqueId, ocppj.GenericError, mockDescription, nil)
	assert.Nil(t, err)
}

func (suite *OcppJTestSuite) TestCentralSystemSendInvalidError() {
	t := suite.T()
	mockChargePointId := "0101"
	mockUniqueId := "6789"
	mockDescription := "mockDescription"
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	suite.mockServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(nil)
	suite.centralSystem.Start(8887, "/{ws}")
	err := suite.centralSystem.SendError(mockChargePointId, mockUniqueId, "InvalidErrorCode", mockDescription, nil)
	assert.NotNil(t, err)
}

func (suite *OcppJTestSuite) TestCentralSystemSendErrorFailed() {
	t := suite.T()
	mockChargePointId := "0101"
	mockUniqueId := "1234"
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	suite.mockServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(errors.New("networkError"))
	suite.centralSystem.Start(8887, "/{ws}")
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
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	// Start central system
	suite.centralSystem.Start(8887, "somePath")
	// Set mocked request in queue and mark as pending
	addMockPendingRequest(suite, mockRequest, mockUniqueId, mockChargePointId)
	// Simulate charge point message
	channel := NewMockWebSocket(mockChargePointId)
	err := suite.mockServer.MessageHandler(channel, []byte(mockConfirmation))
	assert.Nil(t, err)
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
	suite.mockServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(nil)
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	// Start central system
	suite.centralSystem.Start(8887, "somePath")
	// Set mocked request in queue and mark as pending
	addMockPendingRequest(suite, mockRequest, mockUniqueId, mockChargePointId)
	// Simulate charge point message
	channel := NewMockWebSocket(mockChargePointId)
	err := suite.mockServer.MessageHandler(channel, []byte(mockError))
	assert.Nil(t, err)
}

func addMockPendingRequest(suite *OcppJTestSuite, mockRequest ocpp.Request, mockUniqueID string, mockChargePointID string) {
	mockCall, _ := suite.centralSystem.CreateCall(mockRequest)
	mockCall.UniqueId = mockUniqueID
	jsonMessage, _ := mockCall.MarshalJSON()
	requestBundle := ocppj.RequestBundle{
		Call: mockCall,
		Data: jsonMessage,
	}
	q := suite.serverRequestMap.GetOrCreate(mockChargePointID)
	_ = q.Push(requestBundle)
	suite.centralSystem.AddPendingRequest(mockUniqueID, mockRequest)
}
