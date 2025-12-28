package ocppj_test

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocppj"
	"github.com/lorenzodonini/ocpp-go/ws"
)

func (suite *OcppJTestSuite) TestNewServer() {
	s := ocppj.NewServer(nil, nil, nil)
	suite.Assert().NotNil(s)
}

func (suite *OcppJTestSuite) TestServerStart() {
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	suite.centralSystem.Start(8887, "/{ws}")
	suite.Assert().True(suite.serverDispatcher.IsRunning())
}

func (suite *OcppJTestSuite) TestServerNotStartedError() {
	mockChargePointId := "1234"
	// Start normally
	req := newMockRequest("somevalue")
	err := suite.centralSystem.SendRequest(mockChargePointId, req)
	suite.Require().Error(err, "ocppj server is not started, couldn't send request")
	suite.Assert().False(suite.serverDispatcher.IsRunning())
}

func (suite *OcppJTestSuite) TestServerStoppedError() {
	mockChargePointId := "1234"
	// Start server
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	suite.mockServer.On("Stop").Return(nil)
	suite.centralSystem.Start(8887, "/{ws}")
	// Stop server
	suite.centralSystem.Stop()
	// Send message. Expected error
	time.Sleep(20 * time.Millisecond)
	suite.Assert().False(suite.serverDispatcher.IsRunning())
	req := newMockRequest("somevalue")
	err := suite.centralSystem.SendRequest(mockChargePointId, req)
	suite.Assert().Error(err, "ocppj server is not started, couldn't send request")
}

// ----------------- SendRequest tests -----------------

func (suite *OcppJTestSuite) TestCentralSystemSendRequest() {
	mockChargePointId := "1234"
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	suite.mockServer.On("Write", mockChargePointId, mock.Anything).Return(nil)
	suite.centralSystem.Start(8887, "/{ws}")
	suite.serverDispatcher.CreateClient(mockChargePointId)
	mockRequest := newMockRequest("mockValue")
	err := suite.centralSystem.SendRequest(mockChargePointId, mockRequest)
	suite.Assert().Nil(err)
}

func (suite *OcppJTestSuite) TestCentralSystemSendInvalidRequest() {
	mockChargePointId := "1234"
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	suite.mockServer.On("Write", mockChargePointId, mock.Anything).Return(nil)
	suite.centralSystem.Start(8887, "/{ws}")
	suite.serverDispatcher.CreateClient(mockChargePointId)
	mockRequest := newMockRequest("")
	err := suite.centralSystem.SendRequest(mockChargePointId, mockRequest)
	suite.Assert().NotNil(err)
}

func (suite *OcppJTestSuite) TestCentralSystemSendRequestNoValidation() {
	mockChargePointId := "1234"
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	suite.mockServer.On("Write", mockChargePointId, mock.Anything).Return(nil)
	suite.centralSystem.Start(8887, "/{ws}")
	suite.serverDispatcher.CreateClient(mockChargePointId)
	mockRequest := newMockRequest("")
	// Temporarily disable message validation
	ocppj.SetMessageValidation(false)
	defer ocppj.SetMessageValidation(true)
	err := suite.centralSystem.SendRequest(mockChargePointId, mockRequest)
	suite.Assert().Nil(err)
}

func (suite *OcppJTestSuite) TestCentralSystemSendInvalidJsonRequest() {
	mockChargePointId := "1234"
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	suite.mockServer.On("Write", mockChargePointId, mock.Anything).Return(nil)
	suite.centralSystem.Start(8887, "/{ws}")
	suite.serverDispatcher.CreateClient(mockChargePointId)
	mockRequest := newMockRequest("somevalue")
	mockRequest.MockAny = make(chan int)
	err := suite.centralSystem.SendRequest(mockChargePointId, mockRequest)
	suite.Require().Error(err)
	suite.Assert().IsType(&json.UnsupportedTypeError{}, err)
}

func (suite *OcppJTestSuite) TestCentralSystemInvalidMessageHook() {
	mockChargePointId := "1234"
	mockChargePoint := NewMockWebSocket(mockChargePointId)
	// Prepare invalid payload
	mockID := "1234"
	mockPayload := map[string]interface{}{
		"mockValue": float64(1234),
	}
	serializedPayload, err := json.Marshal(mockPayload)
	suite.Require().NoError(err)
	invalidMessage := fmt.Sprintf("[2,\"%v\",\"%s\",%v]", mockID, MockFeatureName, string(serializedPayload))
	expectedError := fmt.Sprintf("[4,\"%v\",\"%v\",\"%v\",{}]", mockID, ocppj.FormatErrorType(suite.centralSystem), "json: cannot unmarshal number into Go struct field MockRequest.mockValue of type string")
	writeHook := suite.mockServer.On("Write", mockChargePointId, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		data := args.Get(1).([]byte)
		suite.Assert().Equal(expectedError, string(data))
	})
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	// Setup hook 1
	suite.centralSystem.SetInvalidMessageHook(func(client ws.Channel, err *ocpp.Error, rawMessage string, parsedFields []interface{}) *ocpp.Error {
		suite.Assert().Equal(mockChargePoint.ID(), client.ID())
		// Verify the correct fields are passed to the hook. Content is very low-level, since parsing failed
		suite.Assert().Equal(float64(ocppj.CALL), parsedFields[0])
		suite.Assert().Equal(mockID, parsedFields[1])
		suite.Assert().Equal(MockFeatureName, parsedFields[2])
		suite.Assert().Equal(mockPayload, parsedFields[3])
		return nil
	})
	suite.centralSystem.Start(8887, "/{ws}")
	// Trigger incoming invalid CALL
	err = suite.mockServer.MessageHandler(mockChargePoint, []byte(invalidMessage))
	ocppErr, ok := err.(*ocpp.Error)
	suite.Require().True(ok)
	suite.Assert().Equal(ocppj.FormatErrorType(suite.centralSystem), ocppErr.Code)
	// Setup hook 2
	mockError := ocpp.NewError(ocppj.InternalError, "custom error", mockID)
	expectedError = fmt.Sprintf("[4,\"%v\",\"%v\",\"%v\",{}]", mockError.MessageId, mockError.Code, mockError.Description)
	writeHook.Run(func(args mock.Arguments) {
		data := args.Get(1).([]byte)
		suite.Assert().Equal(expectedError, string(data))
	})
	suite.centralSystem.SetInvalidMessageHook(func(client ws.Channel, err *ocpp.Error, rawMessage string, parsedFields []interface{}) *ocpp.Error {
		suite.Assert().Equal(mockChargePoint.ID(), client.ID())
		// Verify the correct fields are passed to the hook. Content is very low-level, since parsing failed
		suite.Assert().Equal(float64(ocppj.CALL), parsedFields[0])
		suite.Assert().Equal(mockID, parsedFields[1])
		suite.Assert().Equal(MockFeatureName, parsedFields[2])
		suite.Assert().Equal(mockPayload, parsedFields[3])
		return mockError
	})
	// Trigger incoming invalid CALL that returns custom error
	err = suite.mockServer.MessageHandler(mockChargePoint, []byte(invalidMessage))
	ocppErr, ok = err.(*ocpp.Error)
	suite.Require().True(ok)
	suite.Assert().Equal(mockError.Code, ocppErr.Code)
	suite.Assert().Equal(mockError.Description, ocppErr.Description)
	suite.Assert().Equal(mockError.MessageId, ocppErr.MessageId)
}

func (suite *OcppJTestSuite) TestServerSendInvalidCall() {
	mockChargePointId := "1234"
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	suite.mockServer.On("Write", mockChargePointId, mock.Anything).Return(nil)
	suite.centralSystem.Start(8887, "/{ws}")
	suite.serverDispatcher.CreateClient(mockChargePointId)
	mockRequest := newMockRequest("somevalue")
	// Delete existing profiles and test error
	suite.centralSystem.Profiles = []*ocpp.Profile{}
	err := suite.centralSystem.SendRequest(mockChargePointId, mockRequest)
	suite.Assert().Error(err, fmt.Sprintf("Couldn't create Call for unsupported action %v", mockRequest.GetFeatureName()))
}

func (suite *OcppJTestSuite) TestCentralSystemSendRequestFailed() {
	mockChargePointId := "1234"
	var callID string
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	suite.mockServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(fmt.Errorf("networkError")).Run(func(args mock.Arguments) {
		clientID := args.String(0)
		q, ok := suite.serverRequestMap.Get(clientID)
		suite.Require().True(ok)
		suite.Require().False(q.IsEmpty())
		req := q.Peek().(ocppj.RequestBundle)
		callID = req.Call.GetUniqueId()
		// Before error is returned, the request must still be pending
		_, ok = suite.centralSystem.RequestState.GetClientState(mockChargePointId).GetPendingRequest(callID)
		suite.Assert().True(ok)
	})
	suite.centralSystem.Start(8887, "/{ws}")
	suite.serverDispatcher.CreateClient(mockChargePointId)
	mockRequest := newMockRequest("mockValue")
	err := suite.centralSystem.SendRequest(mockChargePointId, mockRequest)
	// TODO: currently the network error is not returned by SendRequest, but is only generated internally
	suite.Assert().Nil(err)
	// Assert that pending request was removed
	time.Sleep(500 * time.Millisecond)
	_, ok := suite.centralSystem.RequestState.GetClientState(mockChargePointId).GetPendingRequest(callID)
	suite.Assert().False(ok)
}

// ----------------- SendResponse tests -----------------

func (suite *OcppJTestSuite) TestCentralSystemSendConfirmation() {
	mockChargePointId := "0101"
	mockUniqueId := "1234"
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	suite.mockServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(nil)
	suite.centralSystem.Start(8887, "/{ws}")
	suite.serverDispatcher.CreateClient(mockChargePointId)
	mockConfirmation := newMockConfirmation("mockValue")
	err := suite.centralSystem.SendResponse(mockChargePointId, mockUniqueId, mockConfirmation)
	suite.Assert().Nil(err)
}

func (suite *OcppJTestSuite) TestCentralSystemSendInvalidConfirmation() {
	mockChargePointId := "0101"
	mockUniqueId := "6789"
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	suite.mockServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(nil)
	suite.centralSystem.Start(8887, "/{ws}")
	suite.serverDispatcher.CreateClient(mockChargePointId)
	mockConfirmation := newMockConfirmation("")
	// This is allowed. Endpoint doesn't keep track of incoming requests, but only outgoing ones
	err := suite.centralSystem.SendResponse(mockChargePointId, mockUniqueId, mockConfirmation)
	suite.Assert().NotNil(err)
}

func (suite *OcppJTestSuite) TestCentralSystemSendConfirmationNoValidation() {
	mockChargePointId := "0101"
	mockUniqueId := "6789"
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	suite.mockServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(nil)
	suite.centralSystem.Start(8887, "/{ws}")
	suite.serverDispatcher.CreateClient(mockChargePointId)
	mockConfirmation := newMockConfirmation("")
	// Temporarily disable message validation
	ocppj.SetMessageValidation(false)
	defer ocppj.SetMessageValidation(true)
	// This is allowed. Endpoint doesn't keep track of incoming requests, but only outgoing ones
	err := suite.centralSystem.SendResponse(mockChargePointId, mockUniqueId, mockConfirmation)
	suite.Assert().Nil(err)
}

func (suite *OcppJTestSuite) TestCentralSystemSendConfirmationFailed() {
	mockChargePointId := "0101"
	mockUniqueId := "1234"
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	suite.mockServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(fmt.Errorf("networkError"))
	suite.centralSystem.Start(8887, "/{ws}")
	suite.serverDispatcher.CreateClient(mockChargePointId)
	mockConfirmation := newMockConfirmation("mockValue")
	err := suite.centralSystem.SendResponse(mockChargePointId, mockUniqueId, mockConfirmation)
	suite.Assert().NotNil(err)
	expectedErr := fmt.Sprintf("ocpp message (%v): GenericError - networkError", mockUniqueId)
	suite.Assert().ErrorContains(err, expectedErr)
}

// SendError
func (suite *OcppJTestSuite) TestCentralSystemSendError() {
	mockChargePointId := "0101"
	mockUniqueId := "1234"
	mockDescription := "mockDescription"
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	suite.mockServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(nil)
	suite.centralSystem.Start(8887, "/{ws}")
	suite.serverDispatcher.CreateClient(mockChargePointId)
	err := suite.centralSystem.SendError(mockChargePointId, mockUniqueId, ocppj.GenericError, mockDescription, nil)
	suite.Assert().Nil(err)
}

func (suite *OcppJTestSuite) TestCentralSystemSendInvalidError() {
	mockChargePointId := "0101"
	mockUniqueId := "6789"
	mockDescription := "mockDescription"
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	suite.mockServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(nil)
	suite.centralSystem.Start(8887, "/{ws}")
	suite.serverDispatcher.CreateClient(mockChargePointId)
	err := suite.centralSystem.SendError(mockChargePointId, mockUniqueId, "InvalidErrorCode", mockDescription, nil)
	suite.Assert().NotNil(err)
}

func (suite *OcppJTestSuite) TestCentralSystemSendErrorFailed() {
	mockChargePointId := "0101"
	mockUniqueId := "1234"
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	suite.mockServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(fmt.Errorf("networkError"))
	suite.centralSystem.Start(8887, "/{ws}")
	suite.serverDispatcher.CreateClient(mockChargePointId)
	mockConfirmation := newMockConfirmation("mockValue")
	err := suite.centralSystem.SendResponse(mockChargePointId, mockUniqueId, mockConfirmation)
	suite.Assert().NotNil(err)
	expectedErr := fmt.Sprintf("ocpp message (%v): GenericError - networkError", mockUniqueId)
	suite.Assert().ErrorContains(err, expectedErr)
}

func (suite *OcppJTestSuite) TestCentralSystemHandleFailedResponse() {
	msgC := make(chan []byte, 1)
	mockChargePointID := "0101"
	mockUniqueID := "1234"
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	suite.mockServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		data, ok := args.Get(1).([]byte)
		suite.Require().True(ok)
		msgC <- data
	})
	suite.centralSystem.Start(8887, "/{ws}")
	suite.serverDispatcher.CreateClient(mockChargePointID)
	var callResult *ocppj.CallResult
	var callError *ocppj.CallError
	var err error
	// 1. occurrence validation error
	mockField := "CallResult.Payload.MockValue"
	mockResponse := newMockConfirmation("")
	callResult, err = suite.centralSystem.CreateCallResult(mockResponse, mockUniqueID)
	suite.Require().Error(err)
	suite.Require().Nil(callResult)
	suite.centralSystem.HandleFailedResponseError(mockChargePointID, mockUniqueID, err, mockResponse.GetFeatureName())
	rawResponse := <-msgC
	expectedErr := fmt.Sprintf(`[4,"%v","%v","Field %s required but not found for feature %s",{}]`, mockUniqueID, ocppj.OccurrenceConstraintErrorType(suite.centralSystem), mockField, mockResponse.GetFeatureName())
	suite.Assert().Equal(expectedErr, string(rawResponse))
	// 2. property constraint validation error
	val := "len4"
	minParamLength := "5"
	mockResponse = newMockConfirmation(val)
	callResult, err = suite.centralSystem.CreateCallResult(mockResponse, mockUniqueID)
	suite.Require().Error(err)
	suite.Require().Nil(callResult)
	suite.centralSystem.HandleFailedResponseError(mockChargePointID, mockUniqueID, err, mockResponse.GetFeatureName())
	rawResponse = <-msgC
	expectedErr = fmt.Sprintf(`[4,"%v","%v","Field %s must be minimum %s, but was %d for feature %s",{}]`,
		mockUniqueID, ocppj.PropertyConstraintViolation, mockField, minParamLength, len(val), mockResponse.GetFeatureName())
	suite.Assert().Equal(expectedErr, string(rawResponse))
	// 3. profile not supported
	mockUnsupportedResponse := &MockUnsupportedResponse{MockValue: "someValue"}
	callResult, err = suite.centralSystem.CreateCallResult(mockUnsupportedResponse, mockUniqueID)
	suite.Require().Error(err)
	suite.Require().Nil(callResult)
	suite.centralSystem.HandleFailedResponseError(mockChargePointID, mockUniqueID, err, mockUnsupportedResponse.GetFeatureName())
	rawResponse = <-msgC
	expectedErr = fmt.Sprintf(`[4,"%v","%v","couldn't create Call Result for unsupported action %s",{}]`,
		mockUniqueID, ocppj.NotSupported, mockUnsupportedResponse.GetFeatureName())
	suite.Assert().Equal(expectedErr, string(rawResponse))
	// 4. ocpp error validation failed
	invalidErrorCode := "InvalidErrorCode"
	callError, err = suite.centralSystem.CreateCallError(mockUniqueID, ocpp.ErrorCode(invalidErrorCode), "", nil)
	suite.Require().Error(err)
	suite.Require().Nil(callError)
	suite.centralSystem.HandleFailedResponseError(mockChargePointID, mockUniqueID, err, "")
	rawResponse = <-msgC
	expectedErr = fmt.Sprintf(`[4,"%v","%v","Key: 'CallError.ErrorCode' Error:Field validation for 'ErrorCode' failed on the 'errorCode' tag",{}]`,
		mockUniqueID, ocppj.GenericError)
	suite.Assert().Equal(expectedErr, string(rawResponse))
	// 5. marshaling err
	err = suite.centralSystem.SendError(mockChargePointID, mockUniqueID, ocppj.SecurityError, "", make(chan struct{}))
	suite.Require().Error(err)
	suite.centralSystem.HandleFailedResponseError(mockChargePointID, mockUniqueID, err, "")
	rawResponse = <-msgC
	expectedErr = fmt.Sprintf(`[4,"%v","%v","json: unsupported type: chan struct {}",{}]`, mockUniqueID, ocppj.GenericError)
	suite.Assert().Equal(expectedErr, string(rawResponse))
	// 6. network error
	rawErr := fmt.Sprintf("couldn't write to websocket. No socket with id %s is open", mockChargePointID)
	err = ocpp.NewError(ocppj.GenericError, rawErr, mockUniqueID)
	suite.centralSystem.HandleFailedResponseError(mockChargePointID, mockUniqueID, err, "")
	rawResponse = <-msgC
	expectedErr = fmt.Sprintf(`[4,"%v","%v","%s",{}]`, mockUniqueID, ocppj.GenericError, rawErr)
	suite.Assert().Equal(expectedErr, string(rawResponse))
}

// ----------------- Handlers tests -----------------

func (suite *OcppJTestSuite) TestCentralSystemNewClientHandler() {
	mockClientID := "1234"
	connectedC := make(chan bool, 1)
	suite.centralSystem.SetNewClientHandler(func(client ws.Channel) {
		suite.Assert().Equal(mockClientID, client.ID())
		connectedC <- true
	})
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return()
	// Internal ocppj <-> websocket handlers are registered on start
	suite.centralSystem.Start(8887, "somePath")
	// Simulate client connection
	channel := NewMockWebSocket(mockClientID)
	suite.mockServer.NewClientHandler(channel)
	ok := <-connectedC
	suite.Assert().True(ok)
	// client state was created
	_, ok = suite.serverRequestMap.Get(mockClientID)
	suite.Assert().True(ok)
}

func (suite *OcppJTestSuite) TestCentralSystemDisconnectedHandler() {
	mockClientID := "1234"
	connectedC := make(chan bool, 1)
	disconnectedC := make(chan bool, 1)
	suite.centralSystem.SetNewClientHandler(func(client ws.Channel) {
		suite.Assert().Equal(mockClientID, client.ID())
		connectedC <- true
	})
	suite.centralSystem.SetDisconnectedClientHandler(func(client ws.Channel) {
		suite.Assert().Equal(mockClientID, client.ID())
		disconnectedC <- true
	})
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return()
	// Internal ocppj <-> websocket handlers are registered on start
	suite.centralSystem.Start(8887, "somePath")
	// Simulate client connection
	channel := NewMockWebSocket(mockClientID)
	suite.mockServer.NewClientHandler(channel)
	ok := <-connectedC
	suite.Assert().True(ok)
	// Simulate client disconnection
	suite.mockServer.DisconnectedClientHandler(channel)
	ok = <-disconnectedC
	suite.Assert().True(ok)
}

func (suite *OcppJTestSuite) TestCentralSystemRequestHandler() {
	mockChargePointId := "1234"
	mockUniqueId := "5678"
	mockValue := "someValue"
	mockRequest := fmt.Sprintf(`[2,"%v","%v",{"mockValue":"%v"}]`, mockUniqueId, MockFeatureName, mockValue)
	suite.centralSystem.SetRequestHandler(func(chargePoint ws.Channel, request ocpp.Request, requestId string, action string) {
		suite.Assert().Equal(mockChargePointId, chargePoint.ID())
		suite.Assert().Equal(mockUniqueId, requestId)
		suite.Assert().Equal(MockFeatureName, action)
		suite.Assert().NotNil(request)
	})
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return()
	suite.centralSystem.Start(8887, "somePath")
	suite.serverDispatcher.CreateClient(mockChargePointId)
	// Simulate charge point message
	channel := NewMockWebSocket(mockChargePointId)
	err := suite.mockServer.MessageHandler(channel, []byte(mockRequest))
	suite.Assert().Nil(err)
}

func (suite *OcppJTestSuite) TestCentralSystemConfirmationHandler() {
	mockChargePointId := "1234"
	mockUniqueId := "5678"
	mockValue := "someValue"
	mockRequest := newMockRequest("testValue")
	mockConfirmation := fmt.Sprintf(`[3,"%v",{"mockValue":"%v"}]`, mockUniqueId, mockValue)
	suite.centralSystem.SetResponseHandler(func(chargePoint ws.Channel, confirmation ocpp.Response, requestId string) {
		suite.Assert().Equal(mockChargePointId, chargePoint.ID())
		suite.Assert().Equal(mockUniqueId, requestId)
		suite.Assert().NotNil(confirmation)
	})
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	// Start central system
	suite.centralSystem.Start(8887, "somePath")
	suite.serverDispatcher.CreateClient(mockChargePointId)
	// Set mocked request in queue and mark as pending
	addMockPendingRequest(suite, mockRequest, mockUniqueId, mockChargePointId)
	// Simulate charge point message
	channel := NewMockWebSocket(mockChargePointId)
	err := suite.mockServer.MessageHandler(channel, []byte(mockConfirmation))
	suite.Assert().Nil(err)
}

func (suite *OcppJTestSuite) TestCentralSystemErrorHandler() {
	mockChargePointId := "1234"
	mockUniqueId := "5678"
	mockErrorCode := ocppj.GenericError
	mockErrorDescription := "Mock Description"
	mockValue := "someValue"
	mockErrorDetails := make(map[string]interface{})
	mockErrorDetails["details"] = "someValue"
	mockRequest := newMockRequest("testValue")
	mockError := fmt.Sprintf(`[4,"%v","%v","%v",{"details":"%v"}]`, mockUniqueId, mockErrorCode, mockErrorDescription, mockValue)
	suite.centralSystem.SetErrorHandler(func(chargePoint ws.Channel, err *ocpp.Error, details interface{}) {
		suite.Assert().Equal(mockChargePointId, chargePoint.ID())
		suite.Assert().Equal(mockUniqueId, err.MessageId)
		suite.Assert().Equal(mockErrorCode, err.Code)
		suite.Assert().Equal(mockErrorDescription, err.Description)
		suite.Assert().Equal(mockErrorDetails, details)
	})
	suite.mockServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(nil)
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	// Start central system
	suite.centralSystem.Start(8887, "somePath")
	suite.serverDispatcher.CreateClient(mockChargePointId)
	// Set mocked request in queue and mark as pending
	addMockPendingRequest(suite, mockRequest, mockUniqueId, mockChargePointId)
	// Simulate charge point message
	channel := NewMockWebSocket(mockChargePointId)
	err := suite.mockServer.MessageHandler(channel, []byte(mockError))
	suite.Assert().Nil(err)
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
	suite.centralSystem.RequestState.AddPendingRequest(mockChargePointID, mockUniqueID, mockRequest)
}

// ----------------- Queue processing tests -----------------

func (suite *OcppJTestSuite) TestServerEnqueueRequest() {
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	suite.mockServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(nil)
	// Start normally
	suite.centralSystem.Start(8887, "/{ws}")
	mockChargePointId := "1234"
	suite.serverDispatcher.CreateClient(mockChargePointId)
	// Simulate request
	req := newMockRequest("somevalue")
	err := suite.centralSystem.SendRequest(mockChargePointId, req)
	suite.Require().Nil(err)
	time.Sleep(500 * time.Millisecond)
	// Message was sent, but element should still be in queue
	q, ok := suite.serverRequestMap.Get(mockChargePointId)
	suite.Require().True(ok)
	suite.Assert().False(q.IsEmpty())
	suite.Assert().Equal(1, q.Size())
	// Analyze enqueued bundle
	peeked := q.Peek()
	suite.Require().NotNil(peeked)
	bundle, ok := peeked.(ocppj.RequestBundle)
	suite.Require().True(ok)
	suite.Require().NotNil(bundle)
	suite.Assert().Equal(req.GetFeatureName(), bundle.Call.Action)
	marshaled, err := bundle.Call.MarshalJSON()
	suite.Require().Nil(err)
	suite.Assert().Equal(marshaled, bundle.Data)
}

func (suite *OcppJTestSuite) TestEnqueueMultipleRequests() {
	messagesToQueue := 5
	sentMessages := 0
	mockChargePointId := "1234"
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	suite.mockServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Run(func(args mock.Arguments) {
		sentMessages += 1
	}).Return(nil)
	// Start normally
	suite.centralSystem.Start(8887, "/{ws}")
	suite.serverDispatcher.CreateClient(mockChargePointId)
	for i := 0; i < messagesToQueue; i++ {
		req := newMockRequest(fmt.Sprintf("request-%v", i))
		err := suite.centralSystem.SendRequest(mockChargePointId, req)
		suite.Require().Nil(err)
	}
	time.Sleep(500 * time.Millisecond)
	// Only one message was sent, but all elements should still be in queue
	suite.Assert().Equal(1, sentMessages)
	q, ok := suite.serverRequestMap.Get(mockChargePointId)
	suite.Require().True(ok)
	suite.Assert().False(q.IsEmpty())
	suite.Assert().Equal(messagesToQueue, q.Size())
	// Analyze enqueued bundle
	var i int
	for !q.IsEmpty() {
		popped := q.Pop()
		suite.Require().NotNil(popped)
		bundle, ok := popped.(ocppj.RequestBundle)
		suite.Require().True(ok)
		suite.Require().NotNil(bundle)
		suite.Assert().Equal(MockFeatureName, bundle.Call.Action)
		i++
	}
	suite.Assert().Equal(messagesToQueue, i)
}

func (suite *OcppJTestSuite) TestRequestQueueFull() {
	messagesToQueue := queueCapacity
	mockChargePointId := "1234"
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	suite.mockServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(nil)
	// Start normally
	suite.centralSystem.Start(8887, "/{ws}")
	suite.serverDispatcher.CreateClient(mockChargePointId)
	for i := 0; i < messagesToQueue; i++ {
		req := newMockRequest(fmt.Sprintf("request-%v", i))
		err := suite.centralSystem.SendRequest(mockChargePointId, req)
		suite.Require().Nil(err)
	}
	// Queue is now full. Trying to send an additional message should throw an error
	req := newMockRequest("full")
	err := suite.centralSystem.SendRequest(mockChargePointId, req)
	suite.Require().NotNil(err)
	suite.Assert().Equal("request queue is full, cannot push new element", err.Error())
}

func (suite *OcppJTestSuite) TestParallelRequests() {
	messagesToQueue := 10
	sentMessages := 0
	mockChargePointId := "1234"
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	suite.mockServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Run(func(args mock.Arguments) {
		sentMessages += 1
	}).Return(nil)
	// Start normally
	suite.centralSystem.Start(8887, "/{ws}")
	suite.serverDispatcher.CreateClient(mockChargePointId)
	for i := 0; i < messagesToQueue; i++ {
		go func() {
			req := newMockRequest("someReq")
			err := suite.centralSystem.SendRequest(mockChargePointId, req)
			suite.Require().Nil(err)
		}()
	}
	time.Sleep(1000 * time.Millisecond)
	// Only one message was sent, but all elements should still be in queue
	q, ok := suite.serverRequestMap.Get(mockChargePointId)
	suite.Require().True(ok)
	suite.Assert().False(q.IsEmpty())
	suite.Assert().Equal(messagesToQueue, q.Size())
	suite.Assert().Equal(1, sentMessages)
}

// TestRequestFlow tests a typical flow with multiple request-responses, sent to different clients.
//
// Requests are sent concurrently and a response to each message is sent from the mocked client endpoint.
// Both CallResult and CallError messages are returned to test all message types.
func (suite *OcppJTestSuite) TestServerRequestFlow() {
	t := suite.T()
	var mutex sync.Mutex
	messagesToQueue := 10
	processedMessages := 0
	mockChargePoint1 := "cp1"
	mockChargePoint2 := "cp2"
	mockChargePoints := map[string]ws.Channel{
		mockChargePoint1: NewMockWebSocket(mockChargePoint1),
		mockChargePoint2: NewMockWebSocket(mockChargePoint2),
	}
	type triggerData struct {
		clientID string
		call     *ocppj.Call
	}
	sendResponseTrigger := make(chan triggerData, 1)
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	suite.mockServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Run(func(args mock.Arguments) {
		wsID := args.String(0)
		data := args.Get(1).([]byte)
		state := suite.centralSystem.RequestState.GetClientState(wsID)
		call := ParseCall(&suite.centralSystem.Endpoint, state, string(data), t)
		suite.Require().NotNil(call)
		sendResponseTrigger <- triggerData{clientID: wsID, call: call}
	}).Return(nil)
	// Mocked response generator
	var wg sync.WaitGroup
	wg.Add(messagesToQueue * 2)
	go func() {
		for {
			d, ok := <-sendResponseTrigger
			if !ok {
				// Test completed, quitting
				return
			}
			// Get original request to generate meaningful response
			call := d.call
			q, ok := suite.serverRequestMap.Get(d.clientID)
			suite.Require().True(ok)
			suite.Assert().False(q.IsEmpty())
			peeked := q.Peek()
			bundle, _ := peeked.(ocppj.RequestBundle)
			suite.Require().NotNil(bundle)
			suite.Assert().Equal(call.UniqueId, bundle.Call.UniqueId)
			req, _ := call.Payload.(*MockRequest)
			// Send response back to server
			var data []byte
			var err error
			v, _ := strconv.Atoi(req.MockValue)
			if v%2 == 0 {
				// Send CallResult
				resp := newMockConfirmation("someResp")
				res, err := suite.centralSystem.CreateCallResult(resp, call.GetUniqueId())
				suite.Require().Nil(err)
				data, err = res.MarshalJSON()
				suite.Require().Nil(err)
			} else {
				// Send CallError
				res, err := suite.centralSystem.CreateCallError(call.GetUniqueId(), ocppj.GenericError, fmt.Sprintf("error-%v", req.MockValue), nil)
				suite.Require().Nil(err)
				data, err = res.MarshalJSON()
				suite.Require().Nil(err)
			}
			wsChannel := mockChargePoints[d.clientID]
			err = suite.mockServer.MessageHandler(wsChannel, data) // Triggers ocppMessageHandler
			suite.Require().Nil(err)
			// Make sure the top queue element was popped
			mutex.Lock()
			processedMessages += 1
			peeked = q.Peek()
			if peeked != nil {
				bundle, _ := peeked.(ocppj.RequestBundle)
				suite.Require().NotNil(bundle)
				suite.Assert().NotEqual(call.UniqueId, bundle.Call.UniqueId)
			}
			mutex.Unlock()
			wg.Done()
		}
	}()
	// Start server normally
	suite.centralSystem.Start(8887, "/{ws}")
	suite.serverDispatcher.CreateClient(mockChargePoint1)
	suite.serverDispatcher.CreateClient(mockChargePoint2)
	for i := 0; i < messagesToQueue*2; i++ {
		// Select a source client
		var chargePointTarget string
		if i%2 == 0 {
			chargePointTarget = mockChargePoint1
		} else {
			chargePointTarget = mockChargePoint2
		}
		go func(j int, clientID string) {
			req := newMockRequest(fmt.Sprintf("%v", j))
			err := suite.centralSystem.SendRequest(clientID, req)
			suite.Require().Nil(err)
		}(i, chargePointTarget)
	}
	// Wait for processing to complete
	wg.Wait()
	close(sendResponseTrigger)
	q, _ := suite.serverRequestMap.Get(mockChargePoint1)
	suite.Assert().True(q.IsEmpty())
	q, _ = suite.serverRequestMap.Get(mockChargePoint2)
	suite.Assert().True(q.IsEmpty())
}
