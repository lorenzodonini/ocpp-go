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
)

// ----------------- Start tests -----------------

func (suite *OcppJTestSuite) TestNewClient() {
	clientID := "mock_id"
	c := ocppj.NewClient(clientID, suite.mockClient, nil, nil)
	suite.Assert().NotNil(c)
	suite.Assert().Equal(clientID, c.Id)
}

func (suite *OcppJTestSuite) TestChargePointStart() {
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	err := suite.chargePoint.Start("someUrl")
	suite.Assert().Nil(err)
	suite.Assert().True(suite.clientDispatcher.IsRunning())
}

func (suite *OcppJTestSuite) TestChargePointStartFailed() {
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(fmt.Errorf("startError"))
	err := suite.chargePoint.Start("someUrl")
	suite.Assert().NotNil(err)
}

func (suite *OcppJTestSuite) TestClientNotStartedError() {
	// Start normally
	req := newMockRequest("somevalue")
	err := suite.chargePoint.SendRequest(req)
	suite.Require().NotNil(err)
	suite.Assert().Equal("ocppj client is not started, couldn't send request", err.Error())
	suite.Require().True(suite.clientRequestQueue.IsEmpty())
}

func (suite *OcppJTestSuite) TestClientStoppedError() {
	// Start client
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	suite.mockClient.On("Stop").Return(nil).Run(func(args mock.Arguments) {
		// Simulate websocket internal working
		suite.mockClient.DisconnectedHandler(nil)
	})
	call := suite.mockClient.On("IsConnected").Return(true)
	err := suite.chargePoint.Start("someUrl")
	suite.Require().NoError(err)
	// Stop client
	suite.chargePoint.Stop()
	// Send message. Expected error
	time.Sleep(20 * time.Millisecond)
	call.Return(false)
	suite.Assert().False(suite.clientDispatcher.IsRunning())
	req := newMockRequest("somevalue")
	err = suite.chargePoint.SendRequest(req)
	suite.Assert().Error(err, "ocppj client is not started, couldn't send request")
}

// ----------------- SendRequest tests -----------------

func (suite *OcppJTestSuite) TestChargePointSendRequest() {
	suite.mockClient.On("Write", mock.Anything).Return(nil)
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	_ = suite.chargePoint.Start("someUrl")
	mockRequest := newMockRequest("mockValue")
	err := suite.chargePoint.SendRequest(mockRequest)
	suite.Assert().Nil(err)
}

func (suite *OcppJTestSuite) TestChargePointSendInvalidRequest() {
	suite.mockClient.On("Write", mock.Anything).Return(nil)
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	_ = suite.chargePoint.Start("someUrl")
	mockRequest := newMockRequest("")
	err := suite.chargePoint.SendRequest(mockRequest)
	suite.Require().NotNil(err)
}

func (suite *OcppJTestSuite) TestChargePointSendRequestNoValidation() {
	suite.mockClient.On("Write", mock.Anything).Return(nil)
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	_ = suite.chargePoint.Start("someUrl")
	mockRequest := newMockRequest("")
	// Temporarily disable message validation
	ocppj.SetMessageValidation(false)
	defer ocppj.SetMessageValidation(true)
	err := suite.chargePoint.SendRequest(mockRequest)
	suite.Assert().Nil(err)
}

func (suite *OcppJTestSuite) TestChargePointSendInvalidJsonRequest() {
	suite.mockClient.On("Write", mock.Anything).Return(nil)
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	_ = suite.chargePoint.Start("someUrl")
	mockRequest := newMockRequest("somevalue")
	mockRequest.MockAny = make(chan int)
	err := suite.chargePoint.SendRequest(mockRequest)
	suite.Require().Error(err)
	suite.Assert().IsType(&json.UnsupportedTypeError{}, err)
}

func (suite *OcppJTestSuite) TestChargePointInvalidMessageHook() {
	// Prepare invalid payload
	mockID := "1234"
	mockPayload := map[string]interface{}{
		"mockValue": float64(1234),
	}
	serializedPayload, err := json.Marshal(mockPayload)
	suite.Require().NoError(err)
	invalidMessage := fmt.Sprintf("[2,\"%v\",\"%s\",%v]", mockID, MockFeatureName, string(serializedPayload))
	expectedError := fmt.Sprintf("[4,\"%v\",\"%v\",\"%v\",{}]", mockID, ocppj.FormatErrorType(suite.chargePoint), "json: cannot unmarshal number into Go struct field MockRequest.mockValue of type string")
	writeHook := suite.mockClient.On("Write", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		data := args.Get(0).([]byte)
		suite.Assert().Equal(expectedError, string(data))
	})
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	// Setup hook 1
	suite.chargePoint.SetInvalidMessageHook(func(err *ocpp.Error, rawMessage string, parsedFields []interface{}) *ocpp.Error {
		// Verify the correct fields are passed to the hook. Content is very low-level, since parsing failed
		suite.Assert().Equal(float64(ocppj.CALL), parsedFields[0])
		suite.Assert().Equal(mockID, parsedFields[1])
		suite.Assert().Equal(MockFeatureName, parsedFields[2])
		suite.Assert().Equal(mockPayload, parsedFields[3])
		return nil
	})
	_ = suite.chargePoint.Start("someUrl")
	// Trigger incoming invalid CALL
	err = suite.mockClient.MessageHandler([]byte(invalidMessage))
	ocppErr, ok := err.(*ocpp.Error)
	suite.Require().True(ok)
	suite.Assert().Equal(ocppj.FormatErrorType(suite.chargePoint), ocppErr.Code)
	// Setup hook 2
	mockError := ocpp.NewError(ocppj.InternalError, "custom error", mockID)
	expectedError = fmt.Sprintf("[4,\"%v\",\"%v\",\"%v\",{}]", mockError.MessageId, mockError.Code, mockError.Description)
	writeHook.Run(func(args mock.Arguments) {
		data := args.Get(0).([]byte)
		suite.Assert().Equal(expectedError, string(data))
	})
	suite.chargePoint.SetInvalidMessageHook(func(err *ocpp.Error, rawMessage string, parsedFields []interface{}) *ocpp.Error {
		// Verify the correct fields are passed to the hook. Content is very low-level, since parsing failed
		suite.Assert().Equal(float64(ocppj.CALL), parsedFields[0])
		suite.Assert().Equal(mockID, parsedFields[1])
		suite.Assert().Equal(MockFeatureName, parsedFields[2])
		suite.Assert().Equal(mockPayload, parsedFields[3])
		return mockError
	})
	// Trigger incoming invalid CALL that returns custom error
	err = suite.mockClient.MessageHandler([]byte(invalidMessage))
	ocppErr, ok = err.(*ocpp.Error)
	suite.Require().True(ok)
	suite.Assert().Equal(mockError.Code, ocppErr.Code)
	suite.Assert().Equal(mockError.Description, ocppErr.Description)
	suite.Assert().Equal(mockError.MessageId, ocppErr.MessageId)
}

func (suite *OcppJTestSuite) TestChargePointSendInvalidCall() {
	suite.mockClient.On("Write", mock.Anything).Return(nil)
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	_ = suite.chargePoint.Start("someUrl")
	mockRequest := newMockRequest("somevalue")
	// Delete existing profiles and test error
	suite.chargePoint.Profiles = []*ocpp.Profile{}
	err := suite.chargePoint.SendRequest(mockRequest)
	suite.Assert().Error(err, fmt.Sprintf("Couldn't create Call for unsupported action %v", mockRequest.GetFeatureName()))
}

func (suite *OcppJTestSuite) TestChargePointSendRequestFailed() {
	var callID string
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	suite.mockClient.On("Write", mock.Anything).Return(fmt.Errorf("networkError")).Run(func(args mock.Arguments) {
		suite.Require().False(suite.clientRequestQueue.IsEmpty())
		req := suite.clientRequestQueue.Peek().(ocppj.RequestBundle)
		callID = req.Call.GetUniqueId()
		_, ok := suite.chargePoint.RequestState.GetPendingRequest(callID)
		// Before anything is returned, the request must still be pending
		suite.Assert().True(ok)
	})
	_ = suite.chargePoint.Start("someUrl")
	mockRequest := newMockRequest("mockValue")
	err := suite.chargePoint.SendRequest(mockRequest)
	// TODO: currently the network error is not returned by SendRequest, but is only generated internally
	suite.Assert().Nil(err)
	// Assert that pending request was removed
	time.Sleep(500 * time.Millisecond)
	_, ok := suite.chargePoint.RequestState.GetPendingRequest(callID)
	suite.Assert().False(ok)
}

// ----------------- SendResponse tests -----------------

func (suite *OcppJTestSuite) TestChargePointSendConfirmation() {
	mockUniqueId := "1234"
	suite.mockClient.On("Write", mock.Anything).Return(nil)
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	_ = suite.chargePoint.Start("someUrl")
	mockConfirmation := newMockConfirmation("mockValue")
	// This is allowed. Endpoint doesn't keep track of incoming requests, but only outgoing ones
	err := suite.chargePoint.SendResponse(mockUniqueId, mockConfirmation)
	suite.Assert().Nil(err)
}

func (suite *OcppJTestSuite) TestChargePointSendConfirmationNoValidation() {
	mockUniqueId := "6789"
	suite.mockClient.On("Write", mock.Anything).Return(nil)
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	_ = suite.chargePoint.Start("someUrl")
	mockConfirmation := newMockConfirmation("")
	// Temporarily disable message validation
	ocppj.SetMessageValidation(false)
	defer ocppj.SetMessageValidation(true)
	// This is allowed. Endpoint doesn't keep track of incoming requests, but only outgoing ones
	err := suite.chargePoint.SendResponse(mockUniqueId, mockConfirmation)
	suite.Assert().Nil(err)
}

func (suite *OcppJTestSuite) TestChargePointSendInvalidConfirmation() {
	mockUniqueId := "6789"
	suite.mockClient.On("Write", mock.Anything).Return(nil)
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	_ = suite.chargePoint.Start("someUrl")
	mockConfirmation := newMockConfirmation("")
	// This is allowed. Endpoint doesn't keep track of incoming requests, but only outgoing ones
	err := suite.chargePoint.SendResponse(mockUniqueId, mockConfirmation)
	suite.Assert().NotNil(err)
}

func (suite *OcppJTestSuite) TestChargePointSendConfirmationFailed() {
	mockUniqueId := "1234"
	suite.mockClient.On("Write", mock.Anything).Return(fmt.Errorf("networkError"))
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	_ = suite.chargePoint.Start("someUrl")
	mockConfirmation := newMockConfirmation("mockValue")
	err := suite.chargePoint.SendResponse(mockUniqueId, mockConfirmation)
	suite.Assert().NotNil(err)
	expectedErr := fmt.Sprintf("ocpp message (%v): GenericError - networkError", mockUniqueId)
	suite.Assert().ErrorContains(err, expectedErr)
}

// ----------------- SendError tests -----------------

func (suite *OcppJTestSuite) TestChargePointSendError() {
	mockUniqueId := "1234"
	mockDescription := "mockDescription"
	suite.mockClient.On("Write", mock.Anything).Return(nil)
	err := suite.chargePoint.SendError(mockUniqueId, ocppj.GenericError, mockDescription, nil)
	suite.Assert().Nil(err)
}

func (suite *OcppJTestSuite) TestChargePointSendInvalidError() {
	mockUniqueId := "6789"
	mockDescription := "mockDescription"
	suite.mockClient.On("Write", mock.Anything).Return(nil)
	err := suite.chargePoint.SendError(mockUniqueId, "InvalidErrorCode", mockDescription, nil)
	suite.Assert().NotNil(err)
}

func (suite *OcppJTestSuite) TestChargePointSendErrorFailed() {
	mockUniqueId := "1234"
	suite.mockClient.On("Write", mock.Anything).Return(fmt.Errorf("networkError"))
	mockConfirmation := newMockConfirmation("mockValue")
	err := suite.chargePoint.SendResponse(mockUniqueId, mockConfirmation)
	suite.Assert().NotNil(err)
	expectedErr := fmt.Sprintf("ocpp message (%v): GenericError - networkError", mockUniqueId)
	suite.Assert().ErrorContains(err, expectedErr)
}

func (suite *OcppJTestSuite) TestChargePointHandleFailedResponse() {
	msgC := make(chan []byte, 1)
	mockUniqueID := "1234"
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	suite.mockClient.On("Write", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		data, ok := args.Get(0).([]byte)
		suite.Require().True(ok)
		msgC <- data
	})
	var callResult *ocppj.CallResult
	var callError *ocppj.CallError
	var err error
	// 1. occurrence validation error
	mockField := "CallResult.Payload.MockValue"
	mockResponse := newMockConfirmation("")
	callResult, err = suite.chargePoint.CreateCallResult(mockResponse, mockUniqueID)
	suite.Require().Error(err)
	suite.Require().Nil(callResult)
	suite.chargePoint.HandleFailedResponseError(mockUniqueID, err, mockResponse.GetFeatureName())
	rawResponse := <-msgC
	expectedErr := fmt.Sprintf(`[4,"%v","%v","Field %s required but not found for feature %s",{}]`, mockUniqueID, ocppj.OccurrenceConstraintErrorType(suite.chargePoint), mockField, mockResponse.GetFeatureName())
	suite.Assert().Equal(expectedErr, string(rawResponse))
	// 2. property constraint validation error
	val := "len4"
	minParamLength := "5"
	mockResponse = newMockConfirmation(val)
	callResult, err = suite.chargePoint.CreateCallResult(mockResponse, mockUniqueID)
	suite.Require().Error(err)
	suite.Require().Nil(callResult)
	suite.chargePoint.HandleFailedResponseError(mockUniqueID, err, mockResponse.GetFeatureName())
	rawResponse = <-msgC
	expectedErr = fmt.Sprintf(`[4,"%v","%v","Field %s must be minimum %s, but was %d for feature %s",{}]`,
		mockUniqueID, ocppj.PropertyConstraintViolation, mockField, minParamLength, len(val), mockResponse.GetFeatureName())
	suite.Assert().Equal(expectedErr, string(rawResponse))
	// 3. profile not supported
	mockUnsupportedResponse := &MockUnsupportedResponse{MockValue: "someValue"}
	callResult, err = suite.chargePoint.CreateCallResult(mockUnsupportedResponse, mockUniqueID)
	suite.Require().Error(err)
	suite.Require().Nil(callResult)
	suite.chargePoint.HandleFailedResponseError(mockUniqueID, err, mockUnsupportedResponse.GetFeatureName())
	rawResponse = <-msgC
	expectedErr = fmt.Sprintf(`[4,"%v","%v","couldn't create Call Result for unsupported action %s",{}]`,
		mockUniqueID, ocppj.NotSupported, mockUnsupportedResponse.GetFeatureName())
	suite.Assert().Equal(expectedErr, string(rawResponse))
	// 4. ocpp error validation failed
	invalidErrorCode := "InvalidErrorCode"
	callError, err = suite.chargePoint.CreateCallError(mockUniqueID, ocpp.ErrorCode(invalidErrorCode), "", nil)
	suite.Require().Error(err)
	suite.Require().Nil(callError)
	suite.chargePoint.HandleFailedResponseError(mockUniqueID, err, "")
	rawResponse = <-msgC
	expectedErr = fmt.Sprintf(`[4,"%v","%v","Key: 'CallError.ErrorCode' Error:Field validation for 'ErrorCode' failed on the 'errorCode' tag",{}]`,
		mockUniqueID, ocppj.GenericError)
	suite.Assert().Equal(expectedErr, string(rawResponse))
	// 5. marshaling err
	err = suite.chargePoint.SendError(mockUniqueID, ocppj.SecurityError, "", make(chan struct{}))
	suite.Require().Error(err)
	suite.chargePoint.HandleFailedResponseError(mockUniqueID, err, "")
	rawResponse = <-msgC
	expectedErr = fmt.Sprintf(`[4,"%v","%v","json: unsupported type: chan struct {}",{}]`, mockUniqueID, ocppj.GenericError)
	suite.Assert().Equal(expectedErr, string(rawResponse))
	// 6. network error
	rawErr := "client is currently not connected, cannot send data"
	err = ocpp.NewError(ocppj.GenericError, rawErr, mockUniqueID)
	suite.chargePoint.HandleFailedResponseError(mockUniqueID, err, "")
	rawResponse = <-msgC
	expectedErr = fmt.Sprintf(`[4,"%v","%v","%s",{}]`, mockUniqueID, ocppj.GenericError, rawErr)
	suite.Assert().Equal(expectedErr, string(rawResponse))
}

// ----------------- Call Handlers tests -----------------

func (suite *OcppJTestSuite) TestChargePointCallHandler() {
	mockUniqueId := "5678"
	mockValue := "someValue"
	mockRequest := fmt.Sprintf(`[2,"%v","%v",{"mockValue":"%v"}]`, mockUniqueId, MockFeatureName, mockValue)
	suite.chargePoint.SetRequestHandler(func(request ocpp.Request, requestId string, action string) {
		suite.Assert().Equal(mockUniqueId, requestId)
		suite.Assert().Equal(MockFeatureName, action)
		suite.Assert().NotNil(request)
	})
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil).Run(func(args mock.Arguments) {
		// Simulate central system message
		err := suite.mockClient.MessageHandler([]byte(mockRequest))
		suite.Assert().Nil(err)
	})
	err := suite.chargePoint.Start("somePath")
	suite.Assert().Nil(err)
}

func (suite *OcppJTestSuite) TestChargePointCallResultHandler() {
	mockUniqueId := "5678"
	mockValue := "someValue"
	mockRequest := newMockRequest("testValue")
	mockConfirmation := fmt.Sprintf(`[3,"%v",{"mockValue":"%v"}]`, mockUniqueId, mockValue)
	suite.chargePoint.SetResponseHandler(func(confirmation ocpp.Response, requestId string) {
		suite.Assert().Equal(mockUniqueId, requestId)
		suite.Assert().NotNil(confirmation)
	})
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	suite.chargePoint.RequestState.AddPendingRequest(mockUniqueId, mockRequest) // Manually add a pending request, so that response is not rejected
	err := suite.chargePoint.Start("somePath")
	suite.Assert().Nil(err)
	// Simulate central system message
	err = suite.mockClient.MessageHandler([]byte(mockConfirmation))
	suite.Assert().Nil(err)
}

func (suite *OcppJTestSuite) TestChargePointCallErrorHandler() {
	mockUniqueId := "5678"
	mockErrorCode := ocppj.GenericError
	mockErrorDescription := "Mock Description"
	mockValue := "someValue"
	mockErrorDetails := make(map[string]interface{})
	mockErrorDetails["details"] = "someValue"

	mockRequest := newMockRequest("testValue")
	mockError := fmt.Sprintf(`[4,"%v","%v","%v",{"details":"%v"}]`, mockUniqueId, mockErrorCode, mockErrorDescription, mockValue)
	suite.chargePoint.SetErrorHandler(func(err *ocpp.Error, details interface{}) {
		suite.Assert().Equal(mockUniqueId, err.MessageId)
		suite.Assert().Equal(mockErrorCode, err.Code)
		suite.Assert().Equal(mockErrorDescription, err.Description)
		suite.Assert().Equal(mockErrorDetails, details)
	})
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	suite.chargePoint.RequestState.AddPendingRequest(mockUniqueId, mockRequest) // Manually add a pending request, so that response is not rejected
	err := suite.chargePoint.Start("someUrl")
	suite.Assert().Nil(err)
	// Simulate central system message
	err = suite.mockClient.MessageHandler([]byte(mockError))
	suite.Assert().Nil(err)
}

// ----------------- Queue processing tests -----------------

func (suite *OcppJTestSuite) TestClientEnqueueRequest() {
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	suite.mockClient.On("Write", mock.Anything).Return(nil)
	// Start normally
	err := suite.chargePoint.Start("someUrl")
	suite.Require().Nil(err)
	req := newMockRequest("somevalue")
	err = suite.chargePoint.SendRequest(req)
	suite.Require().Nil(err)
	time.Sleep(500 * time.Millisecond)
	// Message was sent, but element should still be in queue
	suite.Require().False(suite.clientRequestQueue.IsEmpty())
	suite.Assert().Equal(1, suite.clientRequestQueue.Size())
	// Analyze enqueued bundle
	peeked := suite.clientRequestQueue.Peek()
	suite.Require().NotNil(peeked)
	bundle, ok := peeked.(ocppj.RequestBundle)
	suite.Require().True(ok)
	suite.Require().NotNil(bundle)
	suite.Assert().Equal(req.GetFeatureName(), bundle.Call.Action)
	marshaled, err := bundle.Call.MarshalJSON()
	suite.Require().Nil(err)
	suite.Assert().Equal(marshaled, bundle.Data)
}

func (suite *OcppJTestSuite) TestClientEnqueueMultipleRequests() {
	messagesToQueue := 5
	sentMessages := 0
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	suite.mockClient.On("Write", mock.Anything).Run(func(args mock.Arguments) {
		sentMessages += 1
	}).Return(nil)
	// Start normally
	err := suite.chargePoint.Start("someUrl")
	suite.Require().Nil(err)
	for i := 0; i < messagesToQueue; i++ {
		req := newMockRequest(fmt.Sprintf("request-%v", i))
		err = suite.chargePoint.SendRequest(req)
		suite.Require().Nil(err)
	}
	time.Sleep(500 * time.Millisecond)
	// Only one message was sent, but all elements should still be in queue
	suite.Assert().Equal(1, sentMessages)
	suite.Require().False(suite.clientRequestQueue.IsEmpty())
	suite.Assert().Equal(messagesToQueue, suite.clientRequestQueue.Size())
	// Analyze enqueued bundle
	var i int
	for !suite.clientRequestQueue.IsEmpty() {
		popped := suite.clientRequestQueue.Pop()
		suite.Require().NotNil(popped)
		bundle, ok := popped.(ocppj.RequestBundle)
		suite.Require().True(ok)
		suite.Require().NotNil(bundle)
		suite.Assert().Equal(MockFeatureName, bundle.Call.Action)
		i++
	}
	suite.Assert().Equal(messagesToQueue, i)
}

func (suite *OcppJTestSuite) TestClientRequestQueueFull() {
	messagesToQueue := queueCapacity
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	suite.mockClient.On("Write", mock.Anything).Return(nil)
	// Start normally
	err := suite.chargePoint.Start("someUrl")
	suite.Require().Nil(err)
	for i := 0; i < messagesToQueue; i++ {
		req := newMockRequest(fmt.Sprintf("request-%v", i))
		err = suite.chargePoint.SendRequest(req)
		suite.Require().Nil(err)
	}
	// Queue is now full. Trying to send an additional message should throw an error
	req := newMockRequest("full")
	err = suite.chargePoint.SendRequest(req)
	suite.Require().NotNil(err)
	suite.Assert().Equal("request queue is full, cannot push new element", err.Error())
}

func (suite *OcppJTestSuite) TestClientParallelRequests() {
	messagesToQueue := 10
	sentMessages := 0
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	suite.mockClient.On("Write", mock.Anything).Run(func(args mock.Arguments) {
		sentMessages += 1
	}).Return(nil)
	// Start normally
	err := suite.chargePoint.Start("someUrl")
	suite.Require().Nil(err)
	for i := 0; i < messagesToQueue; i++ {
		go func() {
			req := newMockRequest("someReq")
			err = suite.chargePoint.SendRequest(req)
			suite.Require().Nil(err)
		}()
	}
	time.Sleep(1000 * time.Millisecond)
	// Only one message was sent, but all element should still be in queue
	suite.Require().False(suite.clientRequestQueue.IsEmpty())
	suite.Assert().Equal(messagesToQueue, suite.clientRequestQueue.Size())
}

// TestClientRequestFlow tests a typical flow with multiple request-responses.
//
// Requests are sent concurrently and a response to each message is sent from the mocked server endpoint.
// Both CallResult and CallError messages are returned to test all message types.
func (suite *OcppJTestSuite) TestClientRequestFlow() {
	t := suite.T()
	var mutex sync.Mutex
	messagesToQueue := 10
	processedMessages := 0
	sendResponseTrigger := make(chan *ocppj.Call, 1)
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	suite.mockClient.On("Write", mock.Anything).Run(func(args mock.Arguments) {
		data := args.Get(0).([]byte)
		call := ParseCall(&suite.chargePoint.Endpoint, suite.chargePoint.RequestState, string(data), t)
		suite.Require().NotNil(call)
		sendResponseTrigger <- call
	}).Return(nil)
	// Mocked response generator
	var wg sync.WaitGroup
	wg.Add(messagesToQueue)
	go func() {
		for {
			call, ok := <-sendResponseTrigger
			if !ok {
				// Test completed, quitting
				return
			}
			// Get original request to generate meaningful response
			peeked := suite.clientRequestQueue.Peek()
			bundle, _ := peeked.(ocppj.RequestBundle)
			suite.Require().NotNil(bundle)
			suite.Assert().Equal(call.UniqueId, bundle.Call.UniqueId)
			req, _ := call.Payload.(*MockRequest)
			// Send response back to client
			var data []byte
			var err error
			v, _ := strconv.Atoi(req.MockValue)
			if v%2 == 0 {
				// Send CallResult
				resp := newMockConfirmation("someResp")
				res, err := suite.chargePoint.CreateCallResult(resp, call.GetUniqueId())
				suite.Require().Nil(err)
				data, err = res.MarshalJSON()
				suite.Require().Nil(err)
			} else {
				// Send CallError
				res, err := suite.chargePoint.CreateCallError(call.GetUniqueId(), ocppj.GenericError, fmt.Sprintf("error-%v", req.MockValue), nil)
				suite.Require().Nil(err)
				data, err = res.MarshalJSON()
				suite.Require().Nil(err)
			}
			fmt.Printf("sending mocked response to message %v\n", call.GetUniqueId())
			err = suite.mockClient.MessageHandler(data) // Triggers ocppMessageHandler
			suite.Require().Nil(err)
			// Make sure the top queue element was popped
			mutex.Lock()
			processedMessages += 1
			peeked = suite.clientRequestQueue.Peek()
			if peeked != nil {
				bundle, _ := peeked.(ocppj.RequestBundle)
				suite.Require().NotNil(bundle)
				suite.Assert().NotEqual(call.UniqueId, bundle.Call.UniqueId)
			}
			mutex.Unlock()
			wg.Done()
		}
	}()
	// Start client normally
	err := suite.chargePoint.Start("someUrl")
	suite.Require().Nil(err)
	for i := 0; i < messagesToQueue; i++ {
		go func(j int) {
			req := newMockRequest(fmt.Sprintf("%v", j))
			err = suite.chargePoint.SendRequest(req)
			suite.Require().Nil(err)
		}(i)
	}
	// Wait for processing to complete
	wg.Wait()
	close(sendResponseTrigger)
	suite.Assert().True(suite.clientRequestQueue.IsEmpty())
}

// TestClientDisconnected ensures that upon disconnection, the client keeps its internal state
// and the internal queue does not change.
func (suite *OcppJTestSuite) TestClientDisconnected() {
	t := suite.T()
	messagesToQueue := 8
	sentMessages := 0
	writeC := make(chan *ocppj.Call, 1)
	triggerC := make(chan bool, 1)
	disconnectError := fmt.Errorf("some error")
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	suite.mockClient.On("Write", mock.Anything).Run(func(args mock.Arguments) {
		sentMessages += 1
		data := args.Get(0).([]byte)
		call := ParseCall(&suite.chargePoint.Endpoint, suite.chargePoint.RequestState, string(data), t)
		suite.Require().NotNil(call)
		writeC <- call
	}).Return(nil)
	// Start normally
	err := suite.chargePoint.Start("someUrl")
	suite.Require().Nil(err)
	// Start mocked response routine
	go func() {
		counter := 0
		for {
			call, ok := <-writeC
			if !ok {
				return
			}
			// Trigger request completion after some artificial delay
			time.Sleep(50 * time.Millisecond)
			suite.clientDispatcher.CompleteRequest(call.UniqueId)
			counter++
			if counter == (messagesToQueue / 2) {
				triggerC <- true
			}
		}
	}()
	// Send some messages
	for i := 0; i < messagesToQueue; i++ {
		req := newMockRequest(fmt.Sprintf("%v", i))
		err = suite.chargePoint.SendRequest(req)
		suite.Require().NoError(err)
	}
	// Wait for trigger disconnect after a few responses were returned
	<-triggerC
	suite.Assert().False(suite.clientDispatcher.IsPaused())
	suite.mockClient.DisconnectedHandler(disconnectError)
	time.Sleep(200 * time.Millisecond)
	// Not all messages were sent, some are still in queue
	suite.Assert().True(suite.clientDispatcher.IsPaused())
	suite.Assert().True(suite.clientDispatcher.IsRunning())
	currentSize := suite.clientRequestQueue.Size()
	currentSent := sentMessages
	// Wait for some more time and double-check
	time.Sleep(500 * time.Millisecond)
	suite.Assert().True(suite.clientDispatcher.IsPaused())
	suite.Assert().True(suite.clientDispatcher.IsRunning())
	suite.Assert().Equal(currentSize, suite.clientRequestQueue.Size())
	suite.Assert().Equal(currentSent, sentMessages)
	suite.Assert().Less(currentSize, messagesToQueue)
	suite.Assert().Less(sentMessages, messagesToQueue)
}

// TestClientReconnected ensures that upon reconnection, the client retains its internal state
// and resumes sending requests.
func (suite *OcppJTestSuite) TestClientReconnected() {
	t := suite.T()
	messagesToQueue := 8
	sentMessages := 0
	writeC := make(chan *ocppj.Call, 1)
	triggerC := make(chan bool, 1)
	disconnectError := fmt.Errorf("some error")
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	suite.mockClient.On("Write", mock.Anything).Run(func(args mock.Arguments) {
		sentMessages += 1
		data := args.Get(0).([]byte)
		call := ParseCall(&suite.chargePoint.Endpoint, suite.chargePoint.RequestState, string(data), t)
		suite.Require().NotNil(call)
		writeC <- call
	}).Return(nil)
	isConnectedCall := suite.mockClient.On("IsConnected").Return(true)
	// Start normally
	err := suite.chargePoint.Start("someUrl")
	suite.Require().Nil(err)
	suite.Assert().True(suite.chargePoint.IsConnected())
	// Start mocked response routine
	go func() {
		counter := 0
		for {
			call, ok := <-writeC
			if !ok {
				return
			}
			// Trigger request completion after some artificial delay
			time.Sleep(50 * time.Millisecond)
			suite.clientDispatcher.CompleteRequest(call.UniqueId)
			counter++
			if counter == (messagesToQueue/2) || counter == messagesToQueue {
				triggerC <- true
			}
		}
	}()
	// Get the pending request state struct
	state := suite.chargePoint.RequestState
	suite.Assert().False(state.HasPendingRequest())
	// Send some messages
	for i := 0; i < messagesToQueue; i++ {
		req := newMockRequest(fmt.Sprintf("%v", i))
		err = suite.chargePoint.SendRequest(req)
		suite.Require().NoError(err)
	}
	// Wait for trigger disconnect after a few responses were returned
	<-triggerC
	isConnectedCall.Return(false)
	suite.mockClient.DisconnectedHandler(disconnectError)
	// One message was sent, but all others are still in queue
	time.Sleep(200 * time.Millisecond)
	suite.Assert().True(suite.clientDispatcher.IsPaused())
	suite.Assert().False(suite.chargePoint.IsConnected())
	// Wait for some more time and then reconnect
	time.Sleep(500 * time.Millisecond)
	isConnectedCall.Return(true)
	suite.mockClient.ReconnectedHandler()
	suite.Assert().False(suite.clientDispatcher.IsPaused())
	suite.Assert().True(suite.clientDispatcher.IsRunning())
	suite.Assert().False(suite.clientRequestQueue.IsEmpty())
	suite.Assert().True(suite.chargePoint.IsConnected())
	// Wait until remaining messages are sent
	<-triggerC
	suite.Assert().False(suite.clientDispatcher.IsPaused())
	suite.Assert().True(suite.clientDispatcher.IsRunning())
	suite.Assert().Equal(messagesToQueue, sentMessages)
	suite.Assert().True(suite.clientRequestQueue.IsEmpty())
	suite.Assert().False(state.HasPendingRequest())
	suite.Assert().True(suite.chargePoint.IsConnected())
}

// TestClientResponseTimeout ensures that upon a response timeout, the client dispatcher:
//
//   - cancels the current pending request
//   - fires an error, which is returned to the caller
func (suite *OcppJTestSuite) TestClientResponseTimeout() {
	t := suite.T()
	requestID := ""
	req := newMockRequest("test")
	timeoutC := make(chan bool, 1)
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	suite.mockClient.On("Write", mock.Anything).Run(func(args mock.Arguments) {
		data := args.Get(0).([]byte)
		call := ParseCall(&suite.chargePoint.Endpoint, suite.chargePoint.RequestState, string(data), t)
		suite.Require().NotNil(call)
		requestID = call.UniqueId
	}).Return(nil)
	suite.clientDispatcher.SetOnRequestCanceled(func(rID string, request ocpp.Request, err *ocpp.Error) {
		suite.Assert().Equal(requestID, rID)
		suite.Assert().Equal(MockFeatureName, request.GetFeatureName())
		suite.Assert().Equal(req, request)
		suite.Assert().Error(err)
		timeoutC <- true
	})
	// Sets a low response timeout for testing purposes
	suite.clientDispatcher.SetTimeout(500 * time.Millisecond)
	// Start normally and send a message
	err := suite.chargePoint.Start("someUrl")
	suite.Require().NoError(err)
	err = suite.chargePoint.SendRequest(req)
	suite.Require().NoError(err)
	// Wait for request to be enqueued, then check state
	time.Sleep(50 * time.Millisecond)
	state := suite.chargePoint.RequestState
	suite.Assert().False(suite.clientRequestQueue.IsEmpty())
	suite.Assert().True(suite.clientDispatcher.IsRunning())
	suite.Assert().Equal(1, suite.clientRequestQueue.Size())
	suite.Assert().True(state.HasPendingRequest())
	// Wait for timeout error to be thrown
	<-timeoutC
	suite.Assert().True(suite.clientRequestQueue.IsEmpty())
	suite.Assert().True(suite.clientDispatcher.IsRunning())
	suite.Assert().False(state.HasPendingRequest())
}

func (suite *OcppJTestSuite) TestStopDisconnectedClient() {
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	suite.mockClient.On("Write", mock.Anything).Return(nil)
	suite.mockClient.On("Stop").Return(nil)
	call := suite.mockClient.On("IsConnected").Return(true)
	// Start normally
	err := suite.chargePoint.Start("someUrl")
	suite.Require().NoError(err)
	// Trigger network disconnect
	disconnectError := fmt.Errorf("some error")
	suite.chargePoint.SetOnDisconnectedHandler(func(err error) {
		suite.Require().Errorf(err, disconnectError.Error())
	})
	call.Return(false)
	suite.mockClient.DisconnectedHandler(disconnectError)
	time.Sleep(100 * time.Millisecond)
	// Dispatcher should be paused
	suite.Assert().True(suite.clientDispatcher.IsPaused())
	suite.Assert().False(suite.chargePoint.IsConnected())
	// Stop client while reconnecting
	suite.chargePoint.Stop()
	time.Sleep(50 * time.Millisecond)
	suite.Assert().True(suite.clientDispatcher.IsPaused())
	suite.Assert().False(suite.chargePoint.IsConnected())
	// Attempt stopping client again
	suite.chargePoint.Stop()
	time.Sleep(50 * time.Millisecond)
	suite.Assert().True(suite.clientDispatcher.IsPaused())
	suite.Assert().False(suite.chargePoint.IsConnected())
}
