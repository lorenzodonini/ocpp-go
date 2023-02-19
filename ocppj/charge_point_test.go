package ocppj_test

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocppj"
)

// ----------------- Start tests -----------------

func (suite *OcppJTestSuite) TestNewClient() {
	clientID := "mock_id"
	c := ocppj.NewClient(clientID, suite.mockClient, nil, nil)
	assert.NotNil(suite.T(), c)
	assert.Equal(suite.T(), clientID, c.Id)
}

func (suite *OcppJTestSuite) TestChargePointStart() {
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	err := suite.chargePoint.Start("someUrl")
	assert.Nil(suite.T(), err)
	assert.True(suite.T(), suite.clientDispatcher.IsRunning())
}

func (suite *OcppJTestSuite) TestChargePointStartFailed() {
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(fmt.Errorf("startError"))
	err := suite.chargePoint.Start("someUrl")
	assert.NotNil(suite.T(), err)
}

func (suite *OcppJTestSuite) TestClientNotStartedError() {
	t := suite.T()
	// Start normally
	req := newMockRequest("somevalue")
	_, err := suite.chargePoint.SendRequest(req, nil, nil)
	require.NotNil(t, err)
	assert.Equal(t, "ocppj client is not started, couldn't send request", err.Error())
	require.True(t, suite.clientRequestQueue.IsEmpty())
}

func (suite *OcppJTestSuite) TestClientStoppedError() {
	t := suite.T()
	// Start client
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	suite.mockClient.On("Stop").Return(nil).Run(func(args mock.Arguments) {
		// Simulate websocket internal working
		suite.mockClient.DisconnectedHandler(nil)
	})
	err := suite.chargePoint.Start("someUrl")
	require.NoError(t, err)
	// Stop client
	suite.chargePoint.Stop()
	// Send message. Expected error
	time.Sleep(20 * time.Millisecond)
	assert.False(t, suite.clientDispatcher.IsRunning())
	req := newMockRequest("somevalue")
	_, err = suite.chargePoint.SendRequest(req, nil, nil)
	assert.Error(t, err, "ocppj client is not started, couldn't send request")
}

// ----------------- SendRequest tests -----------------

func (suite *OcppJTestSuite) TestChargePointSendRequest() {
	suite.mockClient.On("Write", mock.Anything).Return(nil)
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	_ = suite.chargePoint.Start("someUrl")
	mockRequest := newMockRequest("mockValue")
	_, err := suite.chargePoint.SendRequest(mockRequest, nil, nil)
	assert.NoError(suite.T(), err)
}

func (suite *OcppJTestSuite) TestChargePointSendInvalidRequest() {
	suite.mockClient.On("Write", mock.Anything).Return(nil)
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	_ = suite.chargePoint.Start("someUrl")
	mockRequest := newMockRequest("")
	// Send request, then wait for error
	_, err := suite.chargePoint.SendRequest(mockRequest, nil, nil)
	assert.Error(suite.T(), err)
}

func (suite *OcppJTestSuite) TestChargePointSendRequestNoValidation() {
	suite.mockClient.On("Write", mock.Anything).Return(nil)
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	_ = suite.chargePoint.Start("someUrl")
	mockRequest := newMockRequest("")
	// Temporarily disable message validation
	ocppj.SetMessageValidation(false)
	defer ocppj.SetMessageValidation(true)
	_, err := suite.chargePoint.SendRequest(mockRequest, nil, nil)
	assert.NoError(suite.T(), err)
}

func (suite *OcppJTestSuite) TestChargePointSendInvalidJsonRequest() {
	t := suite.T()
	suite.mockClient.On("Write", mock.Anything).Return(nil)
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	_ = suite.chargePoint.Start("someUrl")
	mockRequest := newMockRequest("somevalue")
	invalidType := make(chan int)
	mockRequest.MockAny = invalidType // Will fail when attempting to serialize to json
	// Expected callback with error
	errC := make(chan error, 1)
	cb := func(response ocpp.Response, err error) {
		assert.Nil(t, response)
		errC <- err
	}
	// Send request, then wait for error
	_, err := suite.chargePoint.SendRequest(mockRequest, cb, nil)
	require.NoError(suite.T(), err)
	err = <-errC
	require.Error(t, err)
	// Make sure the native error was a type marshaling error
	nativeErr := json.UnsupportedTypeError{Type: reflect.TypeOf(invalidType)}
	assert.ErrorContains(suite.T(), err, nativeErr.Error())
}

func (suite *OcppJTestSuite) TestChargePointSendInvalidCall() {
	suite.mockClient.On("Write", mock.Anything).Return(nil)
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	_ = suite.chargePoint.Start("someUrl")
	mockRequest := newMockRequest("somevalue")
	// Delete existing profiles and test error
	suite.chargePoint.Profiles = []*ocpp.Profile{}
	_, err := suite.chargePoint.SendRequest(mockRequest, nil, nil)
	assert.Error(suite.T(), err, fmt.Sprintf("Couldn't create Call for unsupported action %v", mockRequest.GetFeatureName()))
}

func (suite *OcppJTestSuite) TestChargePointSendRequestFailed() {
	t := suite.T()
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	suite.mockClient.On("Write", mock.Anything).Return(fmt.Errorf("networkError")).Run(func(args mock.Arguments) {
		require.False(t, suite.clientRequestQueue.IsEmpty())
	})
	_ = suite.chargePoint.Start("someUrl")
	mockRequest := newMockRequest("mockValue")
	// Expected callback with error
	errC := make(chan error, 1)
	cb := func(response ocpp.Response, err error) {
		assert.Nil(t, response)
		errC <- err
	}
	// Send request, then wait for error
	requestID, err := suite.chargePoint.SendRequest(mockRequest, cb, nil)
	//TODO: currently the network error is not returned by SendRequest, but is only generated internally
	require.NoError(t, err)
	err = <-errC
	require.Error(t, err)
	assert.ErrorContains(t, err, "networkError")
	// Assert that pending request was removed
	_, ok := suite.chargePoint.RequestState.GetPendingRequest(requestID)
	assert.False(t, ok)
}

// ----------------- SendResponse tests -----------------

func (suite *OcppJTestSuite) TestChargePointSendConfirmation() {
	t := suite.T()
	mockUniqueId := "1234"
	suite.mockClient.On("Write", mock.Anything).Return(nil)
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	_ = suite.chargePoint.Start("someUrl")
	mockConfirmation := newMockConfirmation("mockValue")
	// This is allowed. Endpoint doesn't keep track of incoming requests, but only outgoing ones
	err := suite.chargePoint.SendResponse(mockUniqueId, mockConfirmation)
	assert.Nil(t, err)
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
	assert.Nil(suite.T(), err)
}

func (suite *OcppJTestSuite) TestChargePointSendInvalidConfirmation() {
	t := suite.T()
	mockUniqueId := "6789"
	suite.mockClient.On("Write", mock.Anything).Return(nil)
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	_ = suite.chargePoint.Start("someUrl")
	mockConfirmation := newMockConfirmation("")
	// This is allowed. Endpoint doesn't keep track of incoming requests, but only outgoing ones
	err := suite.chargePoint.SendResponse(mockUniqueId, mockConfirmation)
	assert.NotNil(t, err)
}

func (suite *OcppJTestSuite) TestChargePointSendConfirmationFailed() {
	t := suite.T()
	mockUniqueId := "1234"
	suite.mockClient.On("Write", mock.Anything).Return(fmt.Errorf("networkError"))
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	_ = suite.chargePoint.Start("someUrl")
	mockConfirmation := newMockConfirmation("mockValue")
	err := suite.chargePoint.SendResponse(mockUniqueId, mockConfirmation)
	assert.NotNil(t, err)
	assert.Equal(t, "networkError", err.Error())
}

// ----------------- SendError tests -----------------

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
	suite.mockClient.On("Write", mock.Anything).Return(fmt.Errorf("networkError"))
	mockConfirmation := newMockConfirmation("mockValue")
	err := suite.chargePoint.SendResponse(mockUniqueId, mockConfirmation)
	assert.NotNil(t, err)
	assert.Equal(t, "networkError", err.Error())
}

// ----------------- Call Handlers tests -----------------

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
	mockRequest := newMockRequest("testValue")
	sentC := make(chan struct{}, 1)
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	suite.mockClient.On("Write", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		// Trigger event as soon as the request was sent out
		sentC <- struct{}{}
	})
	err := suite.chargePoint.Start("somePath")
	assert.Nil(t, err)
	// Create callback with delivery confirmation channel
	responseC := make(chan ocpp.Response, 1)
	cb := func(response ocpp.Response, err error) {
		responseC <- response
	}
	// Send request, then wait for network to send it out
	requestID, err := suite.chargePoint.SendRequest(mockRequest, cb, nil)
	<-sentC
	// Simulate central system message
	mockValue := "someValue"
	mockConfirmation := fmt.Sprintf(`[3,"%v",{"mockValue":"%v"}]`, requestID, mockValue)
	err = suite.mockClient.MessageHandler([]byte(mockConfirmation))
	require.Nil(t, err)
	// Retrieve async response and check equality
	response, ok := <-responseC
	require.True(t, ok)
	require.NotNil(t, response)
	mockResponse, ok := response.(*MockConfirmation)
	require.True(t, ok)
	assert.Equal(t, mockValue, mockResponse.MockValue)
}

func (suite *OcppJTestSuite) TestChargePointCallErrorHandler() {
	t := suite.T()
	mockValue := "someValue"
	mockRequest := newMockRequest("testValue")
	sentC := make(chan struct{}, 1)
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	suite.mockClient.On("Write", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		// Trigger event as soon as the request was sent out
		sentC <- struct{}{}
	})
	err := suite.chargePoint.Start("someUrl")
	assert.Nil(t, err)
	// Create callback with error channel
	errC := make(chan error, 1)
	cb := func(response ocpp.Response, err error) {
		errC <- err
	}
	// Send request, then wait for network to send it out
	requestID, err := suite.chargePoint.SendRequest(mockRequest, cb, nil)
	<-sentC
	// Simulate central system error message
	mockErrorCode := ocppj.GenericError
	mockErrorDescription := "Mock Description"
	mockErrorDetails := make(map[string]interface{})
	mockErrorDetails["details"] = "someValue"
	mockError := fmt.Sprintf(`[4,"%v","%v","%v",{"details":"%v"}]`, requestID, mockErrorCode, mockErrorDescription, mockValue)
	err = suite.mockClient.MessageHandler([]byte(mockError))
	assert.Nil(t, err)
	// Retrieve async error and check equality
	err, ok := <-errC
	require.True(t, ok)
	require.NotNil(t, err)
	ocppErr, ok := err.(*ocpp.Error)
	require.True(t, ok)
	assert.Equal(t, requestID, ocppErr.MessageId)
	assert.Equal(t, mockErrorCode, ocppErr.Code)
	assert.Equal(t, mockErrorDescription, ocppErr.Description)
}

// ----------------- Queue processing tests -----------------

func (suite *OcppJTestSuite) TestClientEnqueueRequest() {
	t := suite.T()
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	sentDataC := make(chan []byte, 1)
	suite.mockClient.On("Write", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		rawData := args.Get(0)
		rawBytes, ok := rawData.([]byte)
		require.True(t, ok)
		sentDataC <- rawBytes
	})
	// Start normally
	err := suite.chargePoint.Start("someUrl")
	require.Nil(t, err)
	req := newMockRequest("somevalue")
	requestID, err := suite.chargePoint.SendRequest(req, nil, nil)
	require.Nil(t, err)
	// Message was sent, but element should still be in queue
	sentData, ok := <-sentDataC
	require.True(t, ok)
	require.False(t, suite.clientRequestQueue.IsEmpty())
	assert.Equal(t, 1, suite.clientRequestQueue.Size())
	// Analyze enqueued bundle
	peeked := suite.clientRequestQueue.Peek()
	require.NotNil(t, peeked)
	bundle, ok := peeked.(ocppj.RequestBundle)
	require.True(t, ok)
	require.NotNil(t, bundle)
	assert.Equal(t, req.GetFeatureName(), bundle.Call.Action)
	assert.Equal(t, requestID, bundle.Call.GetUniqueId())
	marshaled, err := bundle.Call.MarshalJSON()
	require.Nil(t, err)
	assert.Equal(t, marshaled, sentData)
}

func (suite *OcppJTestSuite) TestClientEnqueueMultipleRequests() {
	t := suite.T()
	messagesToQueue := 5
	sentMessages := 0
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	suite.mockClient.On("Write", mock.Anything).Run(func(args mock.Arguments) {
		sentMessages += 1
	}).Return(nil)
	// Start normally
	err := suite.chargePoint.Start("someUrl")
	require.Nil(t, err)
	for i := 0; i < messagesToQueue; i++ {
		req := newMockRequest(fmt.Sprintf("request-%v", i))
		_, err = suite.chargePoint.SendRequest(req, nil, nil)
		require.Nil(t, err)
	}
	// Wait for queue to fill up and first message to be sent out
	time.Sleep(500 * time.Millisecond)
	// Only one message was sent, but all elements should still be in queue
	assert.Equal(t, 1, sentMessages)
	require.False(t, suite.clientRequestQueue.IsEmpty())
	assert.Equal(t, messagesToQueue, suite.clientRequestQueue.Size())
	// Analyze enqueued bundle
	var i = 0
	for !suite.clientRequestQueue.IsEmpty() {
		popped := suite.clientRequestQueue.Pop()
		require.NotNil(t, popped)
		bundle, ok := popped.(ocppj.RequestBundle)
		require.True(t, ok)
		require.NotNil(t, bundle)
		assert.Equal(t, MockFeatureName, bundle.Call.Action)
		payload, ok := bundle.Call.Payload.(*MockRequest)
		require.True(t, ok)
		assert.Equal(t, fmt.Sprintf("request-%v", i), payload.MockValue)
		i++
	}
	assert.Equal(t, messagesToQueue, i)
}

func (suite *OcppJTestSuite) TestClientRequestQueueFull() {
	t := suite.T()
	messagesToQueue := queueCapacity
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	suite.mockClient.On("Write", mock.Anything).Return(nil)
	// Start normally
	err := suite.chargePoint.Start("someUrl")
	require.Nil(t, err)
	for i := 0; i < messagesToQueue; i++ {
		req := newMockRequest(fmt.Sprintf("request-%v", i))
		_, err = suite.chargePoint.SendRequest(req, nil, nil)
		require.Nil(t, err)
	}
	// Queue is now full. Trying to send an additional message should throw an error
	req := newMockRequest("full")
	_, err = suite.chargePoint.SendRequest(req, nil, nil)
	require.NotNil(t, err)
	assert.Equal(t, "request queue is full, cannot push new element", err.Error())
}

func (suite *OcppJTestSuite) TestClientParallelRequests() {
	t := suite.T()
	messagesToQueue := 10
	sentMessages := 0
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	suite.mockClient.On("Write", mock.Anything).Run(func(args mock.Arguments) {
		sentMessages += 1
	}).Return(nil)
	// Start normally
	err := suite.chargePoint.Start("someUrl")
	require.Nil(t, err)
	for i := 0; i < messagesToQueue; i++ {
		go func() {
			req := newMockRequest("someReq")
			_, err = suite.chargePoint.SendRequest(req, nil, nil)
			require.Nil(t, err)
		}()
	}
	// Wait for queue to fill up and first message to be sent out
	time.Sleep(1000 * time.Millisecond)
	// Only one message was sent, but all element should still be in queue
	require.False(t, suite.clientRequestQueue.IsEmpty())
	assert.Equal(t, messagesToQueue, suite.clientRequestQueue.Size())
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
		require.NotNil(t, call)
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
			req, _ := call.Payload.(*MockRequest)
			// Send response back to client
			var data []byte
			var err error
			var callResult *ocppj.CallResult
			var callError *ocppj.CallError
			v, _ := strconv.Atoi(req.MockValue)
			if v%2 == 0 {
				// Send CallResult
				resp := newMockConfirmation("someResp")
				callResult, err = suite.chargePoint.CreateCallResult(resp, call.GetUniqueId())
				require.Nil(t, err)
				data, err = callResult.MarshalJSON()
				require.Nil(t, err)
			} else {
				// Send CallError
				callError, err = suite.chargePoint.CreateCallError(call.GetUniqueId(), ocppj.GenericError, fmt.Sprintf("error-%v", req.MockValue), nil)
				require.Nil(t, err)
				data, err = callError.MarshalJSON()
				require.Nil(t, err)
			}
			fmt.Printf("sending mocked response to message %v\n", call.GetUniqueId())
			err = suite.mockClient.MessageHandler(data) // Triggers ocppMessageHandler
			require.Nil(t, err)
			// Make sure the top queue element was popped
			mutex.Lock()
			processedMessages += 1
			peeked := suite.clientRequestQueue.Peek()
			if peeked != nil {
				bundle, _ := peeked.(ocppj.RequestBundle)
				require.NotNil(t, bundle)
				assert.NotEqual(t, call.UniqueId, bundle.Call.UniqueId)
			}
			mutex.Unlock()
			wg.Done()
		}
	}()
	// Start client normally
	err := suite.chargePoint.Start("someUrl")
	require.Nil(t, err)
	for i := 0; i < messagesToQueue; i++ {
		go func(j int) {
			req := newMockRequest(fmt.Sprintf("%v", j))
			_, err = suite.chargePoint.SendRequest(req, nil, nil)
			require.Nil(t, err)
		}(i)
	}
	// Wait for processing to complete
	wg.Wait()
	close(sendResponseTrigger)
	assert.True(t, suite.clientRequestQueue.IsEmpty())
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
		require.NotNil(t, call)
		writeC <- call
	}).Return(nil)
	// Start normally
	err := suite.chargePoint.Start("someUrl")
	require.Nil(t, err)
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
		_, err = suite.chargePoint.SendRequest(req, nil, nil)
		require.NoError(t, err)
	}
	// Wait for trigger disconnect after a few responses were returned
	<-triggerC
	assert.False(t, suite.clientDispatcher.IsPaused())
	suite.mockClient.DisconnectedHandler(disconnectError)
	time.Sleep(200 * time.Millisecond)
	// Not all messages were sent, some are still in queue
	assert.True(t, suite.clientDispatcher.IsPaused())
	assert.True(t, suite.clientDispatcher.IsRunning())
	currentSize := suite.clientRequestQueue.Size()
	currentSent := sentMessages
	// Wait for some more time and double-check
	time.Sleep(500 * time.Millisecond)
	assert.True(t, suite.clientDispatcher.IsPaused())
	assert.True(t, suite.clientDispatcher.IsRunning())
	assert.Equal(t, currentSize, suite.clientRequestQueue.Size())
	assert.Equal(t, currentSent, sentMessages)
	assert.Less(t, currentSize, messagesToQueue)
	assert.Less(t, sentMessages, messagesToQueue)
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
		require.NotNil(t, call)
		writeC <- call
	}).Return(nil)
	isConnectedCall := suite.mockClient.On("IsConnected").Return(true)
	// Start normally
	err := suite.chargePoint.Start("someUrl")
	require.Nil(t, err)
	assert.True(t, suite.chargePoint.IsConnected())
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
	assert.False(t, state.HasPendingRequest())
	// Send some messages
	for i := 0; i < messagesToQueue; i++ {
		req := newMockRequest(fmt.Sprintf("%v", i))
		_, err = suite.chargePoint.SendRequest(req, nil, nil)
		require.NoError(t, err)
	}
	// Wait for trigger disconnect after a few responses were returned
	<-triggerC
	isConnectedCall.Return(false)
	suite.mockClient.DisconnectedHandler(disconnectError)
	// One message was sent, but all others are still in queue
	time.Sleep(200 * time.Millisecond)
	assert.True(t, suite.clientDispatcher.IsPaused())
	assert.False(t, suite.chargePoint.IsConnected())
	// Wait for some more time and then reconnect
	time.Sleep(500 * time.Millisecond)
	isConnectedCall.Return(true)
	suite.mockClient.ReconnectedHandler()
	time.Sleep(50 * time.Millisecond)
	assert.False(t, suite.clientDispatcher.IsPaused())
	assert.True(t, suite.clientDispatcher.IsRunning())
	assert.False(t, suite.clientRequestQueue.IsEmpty())
	assert.True(t, suite.chargePoint.IsConnected())
	// Wait until remaining messages are sent
	<-triggerC
	assert.False(t, suite.clientDispatcher.IsPaused())
	assert.True(t, suite.clientDispatcher.IsRunning())
	assert.Equal(t, messagesToQueue, sentMessages)
	assert.True(t, suite.clientRequestQueue.IsEmpty())
	assert.False(t, state.HasPendingRequest())
	assert.True(t, suite.chargePoint.IsConnected())
}

// TestClientResponseTimeout ensures that upon a response timeout, the client dispatcher:
//
//   - cancels the current pending request
//   - fires an error, which is returned to the caller
func (suite *OcppJTestSuite) TestClientResponseTimeout() {
	t := suite.T()
	req := newMockRequest("test")
	timeoutC := make(chan bool, 1)
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	suite.mockClient.On("Write", mock.Anything).Return(nil)
	// Sets a low response timeout for testing purposes
	suite.clientDispatcher.SetTimeout(500 * time.Millisecond)
	// Set timeout callback
	timeoutCb := func(response ocpp.Response, err error) {
		assert.Nil(t, response)
		assert.Error(t, err)
		timeoutC <- true
	}
	// Start normally and send a message
	err := suite.chargePoint.Start("someUrl")
	require.NoError(t, err)
	_, err = suite.chargePoint.SendRequest(req, timeoutCb, context.TODO())
	require.NoError(t, err)
	// Wait for request to be enqueued, then check state
	time.Sleep(50 * time.Millisecond)
	state := suite.chargePoint.RequestState
	assert.False(t, suite.clientRequestQueue.IsEmpty())
	assert.True(t, suite.clientDispatcher.IsRunning())
	assert.Equal(t, 1, suite.clientRequestQueue.Size())
	assert.True(t, state.HasPendingRequest())
	// Wait for timeout error to be thrown
	<-timeoutC
	assert.True(t, suite.clientRequestQueue.IsEmpty())
	assert.True(t, suite.clientDispatcher.IsRunning())
	assert.False(t, state.HasPendingRequest())
}
