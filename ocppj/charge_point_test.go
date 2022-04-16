package ocppj_test

import (
	"encoding/json"
	"fmt"
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
	c := ocppj.NewClient(clientID, nil, nil, nil)
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
	err := suite.chargePoint.SendRequest(req)
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
	err = suite.chargePoint.SendRequest(req)
	assert.Error(t, err, "ocppj client is not started, couldn't send request")
}

// ----------------- SendRequest tests -----------------

func (suite *OcppJTestSuite) TestChargePointSendRequest() {
	suite.mockClient.On("Write", mock.Anything).Return(nil)
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	_ = suite.chargePoint.Start("someUrl")
	mockRequest := newMockRequest("mockValue")
	err := suite.chargePoint.SendRequest(mockRequest)
	assert.Nil(suite.T(), err)
}

func (suite *OcppJTestSuite) TestChargePointSendInvalidRequest() {
	suite.mockClient.On("Write", mock.Anything).Return(nil)
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	_ = suite.chargePoint.Start("someUrl")
	mockRequest := newMockRequest("")
	err := suite.chargePoint.SendRequest(mockRequest)
	assert.NotNil(suite.T(), err)
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
	assert.Nil(suite.T(), err)
}

func (suite *OcppJTestSuite) TestChargePointSendInvalidJsonRequest() {
	suite.mockClient.On("Write", mock.Anything).Return(nil)
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	_ = suite.chargePoint.Start("someUrl")
	mockRequest := newMockRequest("somevalue")
	mockRequest.MockAny = make(chan int)
	err := suite.chargePoint.SendRequest(mockRequest)
	require.Error(suite.T(), err)
	assert.IsType(suite.T(), &json.UnsupportedTypeError{}, err)
}

func (suite *OcppJTestSuite) TestChargePointSendInvalidCall() {
	suite.mockClient.On("Write", mock.Anything).Return(nil)
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	_ = suite.chargePoint.Start("someUrl")
	mockRequest := newMockRequest("somevalue")
	// Delete existing profiles and test error
	suite.chargePoint.Profiles = []*ocpp.Profile{}
	err := suite.chargePoint.SendRequest(mockRequest)
	assert.Error(suite.T(), err, fmt.Sprintf("Couldn't create Call for unsupported action %v", mockRequest.GetFeatureName()))
}

func (suite *OcppJTestSuite) TestChargePointSendRequestFailed() {
	t := suite.T()
	var callID string
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	suite.mockClient.On("Write", mock.Anything).Return(fmt.Errorf("networkError")).Run(func(args mock.Arguments) {
		require.False(t, suite.clientRequestQueue.IsEmpty())
		req := suite.clientRequestQueue.Peek().(ocppj.RequestBundle)
		callID = req.Call.GetUniqueId()
		_, ok := suite.chargePoint.RequestState.GetPendingRequest(callID)
		// Before anything is returned, the request must still be pending
		assert.True(t, ok)
	})
	_ = suite.chargePoint.Start("someUrl")
	mockRequest := newMockRequest("mockValue")
	err := suite.chargePoint.SendRequest(mockRequest)
	//TODO: currently the network error is not returned by SendRequest, but is only generated internally
	assert.Nil(t, err)
	// Assert that pending request was removed
	time.Sleep(500 * time.Millisecond)
	_, ok := suite.chargePoint.RequestState.GetPendingRequest(callID)
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
	mockUniqueId := "5678"
	mockValue := "someValue"
	mockRequest := newMockRequest("testValue")
	mockConfirmation := fmt.Sprintf(`[3,"%v",{"mockValue":"%v"}]`, mockUniqueId, mockValue)
	suite.chargePoint.SetResponseHandler(func(confirmation ocpp.Response, requestId string) {
		assert.Equal(t, mockUniqueId, requestId)
		assert.NotNil(t, confirmation)
	})
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	suite.chargePoint.RequestState.AddPendingRequest(mockUniqueId, mockRequest) // Manually add a pending request, so that response is not rejected
	err := suite.chargePoint.Start("somePath")
	assert.Nil(t, err)
	// Simulate central system message
	err = suite.mockClient.MessageHandler([]byte(mockConfirmation))
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
	suite.chargePoint.SetErrorHandler(func(err *ocpp.Error, details interface{}) {
		assert.Equal(t, mockUniqueId, err.MessageId)
		assert.Equal(t, mockErrorCode, err.Code)
		assert.Equal(t, mockErrorDescription, err.Description)
		assert.Equal(t, mockErrorDetails, details)
	})
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	suite.chargePoint.RequestState.AddPendingRequest(mockUniqueId, mockRequest) // Manually add a pending request, so that response is not rejected
	err := suite.chargePoint.Start("someUrl")
	assert.Nil(t, err)
	// Simulate central system message
	err = suite.mockClient.MessageHandler([]byte(mockError))
	assert.Nil(t, err)
}

// ----------------- Queue processing tests -----------------

func (suite *OcppJTestSuite) TestClientEnqueueRequest() {
	t := suite.T()
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	suite.mockClient.On("Write", mock.Anything).Return(nil)
	// Start normally
	err := suite.chargePoint.Start("someUrl")
	require.Nil(t, err)
	req := newMockRequest("somevalue")
	err = suite.chargePoint.SendRequest(req)
	require.Nil(t, err)
	time.Sleep(500 * time.Millisecond)
	// Message was sent, but element should still be in queue
	require.False(t, suite.clientRequestQueue.IsEmpty())
	assert.Equal(t, 1, suite.clientRequestQueue.Size())
	// Analyze enqueued bundle
	peeked := suite.clientRequestQueue.Peek()
	require.NotNil(t, peeked)
	bundle, ok := peeked.(ocppj.RequestBundle)
	require.True(t, ok)
	require.NotNil(t, bundle)
	assert.Equal(t, req.GetFeatureName(), bundle.Call.Action)
	marshaled, err := bundle.Call.MarshalJSON()
	require.Nil(t, err)
	assert.Equal(t, marshaled, bundle.Data)
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
		err = suite.chargePoint.SendRequest(req)
		require.Nil(t, err)
	}
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
		err = suite.chargePoint.SendRequest(req)
		require.Nil(t, err)
	}
	// Queue is now full. Trying to send an additional message should throw an error
	req := newMockRequest("full")
	err = suite.chargePoint.SendRequest(req)
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
			req := newMockRequest(fmt.Sprintf("someReq"))
			err = suite.chargePoint.SendRequest(req)
			require.Nil(t, err)
		}()
	}
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
			peeked := suite.clientRequestQueue.Peek()
			bundle, _ := peeked.(ocppj.RequestBundle)
			require.NotNil(t, bundle)
			assert.Equal(t, call.UniqueId, bundle.Call.UniqueId)
			req, _ := call.Payload.(*MockRequest)
			// Send response back to client
			var data []byte
			var err error
			v, _ := strconv.Atoi(req.MockValue)
			if v%2 == 0 {
				// Send CallResult
				resp := newMockConfirmation("someResp")
				res, err := suite.chargePoint.CreateCallResult(resp, call.GetUniqueId())
				require.Nil(t, err)
				data, err = res.MarshalJSON()
				require.Nil(t, err)
			} else {
				// Send CallError
				res, err := suite.chargePoint.CreateCallError(call.GetUniqueId(), ocppj.GenericError, fmt.Sprintf("error-%v", req.MockValue), nil)
				require.Nil(t, err)
				data, err = res.MarshalJSON()
				require.Nil(t, err)
			}
			fmt.Printf("sending mocked response to message %v\n", call.GetUniqueId())
			err = suite.mockClient.MessageHandler(data) // Triggers ocppMessageHandler
			require.Nil(t, err)
			// Make sure the top queue element was popped
			mutex.Lock()
			processedMessages += 1
			peeked = suite.clientRequestQueue.Peek()
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
			err = suite.chargePoint.SendRequest(req)
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
		err = suite.chargePoint.SendRequest(req)
		require.NoError(t, err)
	}
	// Wait for trigger disconnect after a few responses were returned
	_ = <-triggerC
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
		err = suite.chargePoint.SendRequest(req)
		require.NoError(t, err)
	}
	// Wait for trigger disconnect after a few responses were returned
	_ = <-triggerC
	suite.mockClient.DisconnectedHandler(disconnectError)
	// One message was sent, but all others are still in queue
	time.Sleep(200 * time.Millisecond)
	assert.True(t, suite.clientDispatcher.IsPaused())
	// Wait for some more time and then reconnect
	time.Sleep(500 * time.Millisecond)
	suite.mockClient.ReconnectedHandler()
	assert.False(t, suite.clientDispatcher.IsPaused())
	assert.True(t, suite.clientDispatcher.IsRunning())
	assert.False(t, suite.clientRequestQueue.IsEmpty())
	// Wait until remaining messages are sent
	_ = <-triggerC
	assert.False(t, suite.clientDispatcher.IsPaused())
	assert.True(t, suite.clientDispatcher.IsRunning())
	assert.Equal(t, messagesToQueue, sentMessages)
	assert.True(t, suite.clientRequestQueue.IsEmpty())
	assert.False(t, state.HasPendingRequest())
}

// TestClientResponseTimeout ensures that upon a response timeout, the client dispatcher:
//
//  - cancels the current pending request
//	- fires an error, which is returned to the caller
func (suite *OcppJTestSuite) TestClientResponseTimeout() {
	t := suite.T()
	requestID := ""
	req := newMockRequest("test")
	timeoutC := make(chan bool, 1)
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	suite.mockClient.On("Write", mock.Anything).Run(func(args mock.Arguments) {
		data := args.Get(0).([]byte)
		call := ParseCall(&suite.chargePoint.Endpoint, suite.chargePoint.RequestState, string(data), t)
		require.NotNil(t, call)
		requestID = call.UniqueId
	}).Return(nil)
	suite.clientDispatcher.SetOnRequestCanceled(func(rID string, request ocpp.Request, err *ocpp.Error) {
		assert.Equal(t, requestID, rID)
		assert.Equal(t, MockFeatureName, request.GetFeatureName())
		assert.Equal(t, req, request)
		assert.Error(t, err)
		timeoutC <- true
	})
	// Sets a low response timeout for testing purposes
	suite.clientDispatcher.SetTimeout(500 * time.Millisecond)
	// Start normally and send a message
	err := suite.chargePoint.Start("someUrl")
	require.NoError(t, err)
	err = suite.chargePoint.SendRequest(req)
	require.NoError(t, err)
	// Wait for request to be enqueued, then check state
	time.Sleep(50 * time.Millisecond)
	state := suite.chargePoint.RequestState
	assert.False(t, suite.clientRequestQueue.IsEmpty())
	assert.True(t, suite.clientDispatcher.IsRunning())
	assert.Equal(t, 1, suite.clientRequestQueue.Size())
	assert.True(t, state.HasPendingRequest())
	// Wait for timeout error to be thrown
	_ = <-timeoutC
	assert.True(t, suite.clientRequestQueue.IsEmpty())
	assert.True(t, suite.clientDispatcher.IsRunning())
	assert.False(t, state.HasPendingRequest())
}
