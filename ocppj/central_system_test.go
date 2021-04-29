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
	"github.com/lorenzodonini/ocpp-go/ws"
)

func (suite *OcppJTestSuite) TestNewServer() {
	s := ocppj.NewServer(nil, nil, nil)
	assert.NotNil(suite.T(), s)
}

func (suite *OcppJTestSuite) TestServerStart() {
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	suite.centralSystem.Start(8887, "/{ws}")
	assert.True(suite.T(), suite.serverDispatcher.IsRunning())
}

func (suite *OcppJTestSuite) TestServerNotStartedError() {
	t := suite.T()
	mockChargePointId := "1234"
	// Start normally
	req := newMockRequest("somevalue")
	err := suite.centralSystem.SendRequest(mockChargePointId, req)
	require.Error(t, err, "ocppj server is not started, couldn't send request")
	assert.False(t, suite.serverDispatcher.IsRunning())
}

func (suite *OcppJTestSuite) TestServerStoppedError() {
	t := suite.T()
	mockChargePointId := "1234"
	// Start server
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	suite.mockServer.On("Stop").Return(nil)
	suite.centralSystem.Start(8887, "/{ws}")
	// Stop server
	suite.centralSystem.Stop()
	// Send message. Expected error
	assert.False(t, suite.serverDispatcher.IsRunning())
	req := newMockRequest("somevalue")
	err := suite.centralSystem.SendRequest(mockChargePointId, req)
	assert.Error(t, err, "ocppj server is not started, couldn't send request")
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
	assert.Nil(suite.T(), err)
}

func (suite *OcppJTestSuite) TestCentralSystemSendInvalidRequest() {
	mockChargePointId := "1234"
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	suite.mockServer.On("Write", mockChargePointId, mock.Anything).Return(nil)
	suite.centralSystem.Start(8887, "/{ws}")
	suite.serverDispatcher.CreateClient(mockChargePointId)
	mockRequest := newMockRequest("")
	err := suite.centralSystem.SendRequest(mockChargePointId, mockRequest)
	assert.NotNil(suite.T(), err)
}

func (suite *OcppJTestSuite) TestServerSendInvalidJsonRequest() {
	mockChargePointId := "1234"
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	suite.mockServer.On("Write", mockChargePointId, mock.Anything).Return(nil)
	suite.centralSystem.Start(8887, "/{ws}")
	suite.serverDispatcher.CreateClient(mockChargePointId)
	mockRequest := newMockRequest("somevalue")
	mockRequest.MockAny = make(chan int)
	err := suite.centralSystem.SendRequest(mockChargePointId, mockRequest)
	require.Error(suite.T(), err)
	assert.IsType(suite.T(), &json.UnsupportedTypeError{}, err)
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
	assert.Error(suite.T(), err, fmt.Sprintf("Couldn't create Call for unsupported action %v", mockRequest.GetFeatureName()))
}

func (suite *OcppJTestSuite) TestCentralSystemSendRequestFailed() {
	t := suite.T()
	mockChargePointId := "1234"
	var callID string
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	suite.mockServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(fmt.Errorf("networkError")).Run(func(args mock.Arguments) {
		clientID := args.String(0)
		q, ok := suite.serverRequestMap.Get(clientID)
		require.True(t, ok)
		require.False(t, q.IsEmpty())
		req := q.Peek().(ocppj.RequestBundle)
		callID = req.Call.GetUniqueId()
		// Before error is returned, the request must still be pending
		_, ok = suite.centralSystem.RequestState.GetClientState(mockChargePointId).GetPendingRequest(callID)
		assert.True(t, ok)
	})
	suite.centralSystem.Start(8887, "/{ws}")
	suite.serverDispatcher.CreateClient(mockChargePointId)
	mockRequest := newMockRequest("mockValue")
	err := suite.centralSystem.SendRequest(mockChargePointId, mockRequest)
	//TODO: currently the network error is not returned by SendRequest, but is only generated internally
	assert.Nil(t, err)
	// Assert that pending request was removed
	time.Sleep(500 * time.Millisecond)
	_, ok := suite.centralSystem.RequestState.GetClientState(mockChargePointId).GetPendingRequest(callID)
	assert.False(t, ok)
}

// ----------------- SendResponse tests -----------------

func (suite *OcppJTestSuite) TestCentralSystemSendConfirmation() {
	t := suite.T()
	mockChargePointId := "0101"
	mockUniqueId := "1234"
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	suite.mockServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(nil)
	suite.centralSystem.Start(8887, "/{ws}")
	suite.serverDispatcher.CreateClient(mockChargePointId)
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
	suite.serverDispatcher.CreateClient(mockChargePointId)
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
	suite.mockServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(fmt.Errorf("networkError"))
	suite.centralSystem.Start(8887, "/{ws}")
	suite.serverDispatcher.CreateClient(mockChargePointId)
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
	suite.serverDispatcher.CreateClient(mockChargePointId)
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
	suite.serverDispatcher.CreateClient(mockChargePointId)
	err := suite.centralSystem.SendError(mockChargePointId, mockUniqueId, "InvalidErrorCode", mockDescription, nil)
	assert.NotNil(t, err)
}

func (suite *OcppJTestSuite) TestCentralSystemSendErrorFailed() {
	t := suite.T()
	mockChargePointId := "0101"
	mockUniqueId := "1234"
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	suite.mockServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(fmt.Errorf("networkError"))
	suite.centralSystem.Start(8887, "/{ws}")
	suite.serverDispatcher.CreateClient(mockChargePointId)
	mockConfirmation := newMockConfirmation("mockValue")
	err := suite.centralSystem.SendResponse(mockChargePointId, mockUniqueId, mockConfirmation)
	assert.NotNil(t, err)
}

// Handlers

func (suite *OcppJTestSuite) TestCentralSystemNewClientHandler() {
	t := suite.T()
	mockClientID := "1234"
	connectedC := make(chan bool, 1)
	suite.centralSystem.SetNewClientHandler(func(client ws.Channel) {
		assert.Equal(t, mockClientID, client.ID())
		connectedC <- true
	})
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return()
	// Internal ocppj <-> websocket handlers are registered on start
	suite.centralSystem.Start(8887, "somePath")
	// Simulate client connection
	channel := NewMockWebSocket(mockClientID)
	suite.mockServer.NewClientHandler(channel)
	ok, _ := <-connectedC
	assert.True(t, ok)
	// Client state was created
	_, ok = suite.serverRequestMap.Get(mockClientID)
	assert.True(t, ok)
}

func (suite *OcppJTestSuite) TestCentralSystemDisconnectedHandler() {
	t := suite.T()
	mockClientID := "1234"
	connectedC := make(chan bool, 1)
	disconnectedC := make(chan bool, 1)
	suite.centralSystem.SetNewClientHandler(func(client ws.Channel) {
		assert.Equal(t, mockClientID, client.ID())
		connectedC <- true
	})
	suite.centralSystem.SetDisconnectedClientHandler(func(client ws.Channel) {
		assert.Equal(t, mockClientID, client.ID())
		disconnectedC <- true
	})
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return()
	// Internal ocppj <-> websocket handlers are registered on start
	suite.centralSystem.Start(8887, "somePath")
	// Simulate client connection
	channel := NewMockWebSocket(mockClientID)
	suite.mockServer.NewClientHandler(channel)
	ok, _ := <-connectedC
	assert.True(t, ok)
	// Simulate client disconnection
	suite.mockServer.DisconnectedClientHandler(channel)
	ok, _ = <-disconnectedC
	assert.True(t, ok)
}

func (suite *OcppJTestSuite) TestCentralSystemRequestHandler() {
	t := suite.T()
	mockChargePointId := "1234"
	mockUniqueId := "5678"
	mockValue := "someValue"
	mockRequest := fmt.Sprintf(`[2,"%v","%v",{"mockValue":"%v"}]`, mockUniqueId, MockFeatureName, mockValue)
	suite.centralSystem.SetRequestHandler(func(chargePoint ws.Channel, request ocpp.Request, requestId string, action string) {
		assert.Equal(t, mockChargePointId, chargePoint.ID())
		assert.Equal(t, mockUniqueId, requestId)
		assert.Equal(t, MockFeatureName, action)
		assert.NotNil(t, request)
	})
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return()
	suite.centralSystem.Start(8887, "somePath")
	suite.serverDispatcher.CreateClient(mockChargePointId)
	// Simulate charge point message
	channel := NewMockWebSocket(mockChargePointId)
	err := suite.mockServer.MessageHandler(channel, []byte(mockRequest))
	assert.Nil(t, err)
}

func (suite *OcppJTestSuite) TestCentralSystemConfirmationHandler() {
	t := suite.T()
	mockChargePointId := "1234"
	mockUniqueId := "5678"
	mockValue := "someValue"
	mockRequest := newMockRequest("testValue")
	mockConfirmation := fmt.Sprintf(`[3,"%v",{"mockValue":"%v"}]`, mockUniqueId, mockValue)
	suite.centralSystem.SetResponseHandler(func(chargePoint ws.Channel, confirmation ocpp.Response, requestId string) {
		assert.Equal(t, mockChargePointId, chargePoint.ID())
		assert.Equal(t, mockUniqueId, requestId)
		assert.NotNil(t, confirmation)
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
	suite.centralSystem.SetErrorHandler(func(chargePoint ws.Channel, err *ocpp.Error, details interface{}) {
		assert.Equal(t, mockChargePointId, chargePoint.ID())
		assert.Equal(t, mockUniqueId, err.MessageId)
		assert.Equal(t, mockErrorCode, err.Code)
		assert.Equal(t, mockErrorDescription, err.Description)
		assert.Equal(t, mockErrorDetails, details)
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
	suite.centralSystem.RequestState.AddPendingRequest(mockChargePointID, mockUniqueID, mockRequest)
}

// ----------------- Queue processing tests -----------------

func (suite *OcppJTestSuite) TestServerEnqueueRequest() {
	t := suite.T()
	suite.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	suite.mockServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(nil)
	// Start normally
	suite.centralSystem.Start(8887, "/{ws}")
	mockChargePointId := "1234"
	suite.serverDispatcher.CreateClient(mockChargePointId)
	// Simulate request
	req := newMockRequest("somevalue")
	err := suite.centralSystem.SendRequest(mockChargePointId, req)
	require.Nil(t, err)
	time.Sleep(500 * time.Millisecond)
	// Message was sent, but element should still be in queue
	q, ok := suite.serverRequestMap.Get(mockChargePointId)
	require.True(t, ok)
	assert.False(t, q.IsEmpty())
	assert.Equal(t, 1, q.Size())
	// Analyze enqueued bundle
	peeked := q.Peek()
	require.NotNil(t, peeked)
	bundle, ok := peeked.(ocppj.RequestBundle)
	require.True(t, ok)
	require.NotNil(t, bundle)
	assert.Equal(t, req.GetFeatureName(), bundle.Call.Action)
	marshaled, err := bundle.Call.MarshalJSON()
	require.Nil(t, err)
	assert.Equal(t, marshaled, bundle.Data)
}

func (suite *OcppJTestSuite) TestEnqueueMultipleRequests() {
	t := suite.T()
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
		require.Nil(t, err)
	}
	time.Sleep(500 * time.Millisecond)
	// Only one message was sent, but all elements should still be in queue
	assert.Equal(t, 1, sentMessages)
	q, ok := suite.serverRequestMap.Get(mockChargePointId)
	require.True(t, ok)
	assert.False(t, q.IsEmpty())
	assert.Equal(t, messagesToQueue, q.Size())
	// Analyze enqueued bundle
	var i = 0
	for !q.IsEmpty() {
		popped := q.Pop()
		require.NotNil(t, popped)
		bundle, ok := popped.(ocppj.RequestBundle)
		require.True(t, ok)
		require.NotNil(t, bundle)
		assert.Equal(t, MockFeatureName, bundle.Call.Action)
		i++
	}
	assert.Equal(t, messagesToQueue, i)
}

func (suite *OcppJTestSuite) TestRequestQueueFull() {
	t := suite.T()
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
		require.Nil(t, err)
	}
	// Queue is now full. Trying to send an additional message should throw an error
	req := newMockRequest("full")
	err := suite.centralSystem.SendRequest(mockChargePointId, req)
	require.NotNil(t, err)
	assert.Equal(t, "request queue is full, cannot push new element", err.Error())
}

func (suite *OcppJTestSuite) TestParallelRequests() {
	t := suite.T()
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
			req := newMockRequest(fmt.Sprintf("someReq"))
			err := suite.centralSystem.SendRequest(mockChargePointId, req)
			require.Nil(t, err)
		}()
	}
	time.Sleep(1000 * time.Millisecond)
	// Only one message was sent, but all elements should still be in queue
	q, ok := suite.serverRequestMap.Get(mockChargePointId)
	require.True(t, ok)
	assert.False(t, q.IsEmpty())
	assert.Equal(t, messagesToQueue, q.Size())
	assert.Equal(t, 1, sentMessages)
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
		require.NotNil(t, call)
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
			require.True(t, ok)
			assert.False(t, q.IsEmpty())
			peeked := q.Peek()
			bundle, _ := peeked.(ocppj.RequestBundle)
			require.NotNil(t, bundle)
			assert.Equal(t, call.UniqueId, bundle.Call.UniqueId)
			req, _ := call.Payload.(*MockRequest)
			// Send response back to server
			var data []byte
			var err error
			v, _ := strconv.Atoi(req.MockValue)
			if v%2 == 0 {
				// Send CallResult
				resp := newMockConfirmation("someResp")
				res, err := suite.centralSystem.CreateCallResult(resp, call.GetUniqueId())
				require.Nil(t, err)
				data, err = res.MarshalJSON()
				require.Nil(t, err)
			} else {
				// Send CallError
				res := suite.centralSystem.CreateCallError(call.GetUniqueId(), ocppj.GenericError, fmt.Sprintf("error-%v", req.MockValue), nil)
				data, err = res.MarshalJSON()
				require.Nil(t, err)
			}
			wsChannel := mockChargePoints[d.clientID]
			err = suite.mockServer.MessageHandler(wsChannel, data) // Triggers ocppMessageHandler
			require.Nil(t, err)
			// Make sure the top queue element was popped
			mutex.Lock()
			processedMessages += 1
			peeked = q.Peek()
			if peeked != nil {
				bundle, _ := peeked.(ocppj.RequestBundle)
				require.NotNil(t, bundle)
				assert.NotEqual(t, call.UniqueId, bundle.Call.UniqueId)
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
			require.Nil(t, err)
		}(i, chargePointTarget)
	}
	// Wait for processing to complete
	wg.Wait()
	close(sendResponseTrigger)
	q, _ := suite.serverRequestMap.Get(mockChargePoint1)
	assert.True(t, q.IsEmpty())
	q, _ = suite.serverRequestMap.Get(mockChargePoint2)
	assert.True(t, q.IsEmpty())
}
