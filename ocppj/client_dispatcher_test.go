package ocppj_test

import (
	"context"
	"fmt"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocppj"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ClientDispatcherTestSuite struct {
	suite.Suite
	state           ocppj.ClientState
	queue           ocppj.RequestQueue
	dispatcher      ocppj.ClientDispatcher
	endpoint        ocppj.Client
	websocketClient MockWebsocketClient
}

func (c *ClientDispatcherTestSuite) SetupTest() {
	c.endpoint = ocppj.Client{Id: "client1"}
	mockProfile := ocpp.NewProfile("mock", MockFeature{})
	c.endpoint.AddProfile(mockProfile)
	c.queue = ocppj.NewFIFOClientQueue(10)
	c.dispatcher = ocppj.NewDefaultClientDispatcher(c.queue)
	c.state = ocppj.NewClientState()
	c.dispatcher.SetPendingRequestState(c.state)
	c.websocketClient = MockWebsocketClient{}
	c.dispatcher.SetNetworkSendHandler(c.websocketClient.Write)
	c.dispatcher.SetOnRequestCanceled(func(request ocppj.RequestBundle, err error) {
		if request.Callback != nil {
			go request.Callback(nil, err)
		}
	})
}

func (c *ClientDispatcherTestSuite) TestClientSendRequest() {
	t := c.T()
	// Setup
	sent := make(chan bool, 1)
	c.websocketClient.On("Write", mock.Anything).Run(func(args mock.Arguments) {
		sent <- true
	}).Return(nil)
	c.dispatcher.Start()
	require.True(t, c.dispatcher.IsRunning())
	// Create and send mock request
	req := newMockRequest("somevalue")
	call, err := c.endpoint.CreateCall(req)
	require.NoError(t, err)
	requestID := call.GetUniqueId()
	bundle := ocppj.RequestBundle{Call: call, Callback: nil, Context: context.TODO()}
	err = c.dispatcher.SendRequest(bundle)
	require.NoError(t, err)
	// Check underlying queue
	assert.False(t, c.queue.IsEmpty())
	assert.Equal(t, 1, c.queue.Size())
	// Wait for websocket to send message
	_, ok := <-sent
	assert.True(t, ok)
	assert.True(t, c.state.HasPendingRequest())
	// Complete request
	c.dispatcher.CompleteRequest(requestID)
	assert.False(t, c.state.HasPendingRequest())
	assert.True(t, c.queue.IsEmpty())
}

func (c *ClientDispatcherTestSuite) TestClientRequestFailed() {
	t := c.T()
	// Setup
	canceled := make(chan bool, 1)
	writeC := make(chan bool, 1)
	errMsg := "mockError"
	c.websocketClient.On("Write", mock.Anything).Run(func(args mock.Arguments) {
		<-writeC
	}).Return(fmt.Errorf(errMsg))
	// Create mock request
	req := newMockRequest("somevalue")
	call, err := c.endpoint.CreateCall(req)
	require.NoError(t, err)
	requestID := call.UniqueId
	bundle := ocppj.RequestBundle{Call: call, Callback: nil, Context: context.TODO()}
	// Start and set canceled callback
	c.dispatcher.SetOnRequestCanceled(func(reqBundle ocppj.RequestBundle, err error) {
		assert.Equal(t, requestID, reqBundle.Call.GetUniqueId())
		assert.Equal(t, MockFeatureName, reqBundle.Call.Payload.GetFeatureName())
		assert.Equal(t, req, reqBundle.Call.Payload)
		ocppErr, ok := err.(*ocpp.Error)
		require.True(t, ok)
		assert.Equal(t, ocppj.GenericError, ocppErr.Code)
		assert.Equal(t, errMsg, ocppErr.Description)
		canceled <- true
	})
	c.dispatcher.Start()
	require.True(t, c.dispatcher.IsRunning())
	// Send mock request
	err = c.dispatcher.SendRequest(bundle)
	require.NoError(t, err)
	// Check underlying queue
	time.Sleep(100 * time.Millisecond)
	assert.False(t, c.queue.IsEmpty())
	assert.Equal(t, 1, c.queue.Size())
	assert.True(t, c.state.HasPendingRequest())
	// Signal that write can occur now, then check canceled request
	writeC <- true
	_, ok := <-canceled
	require.True(t, ok)
	assert.False(t, c.state.HasPendingRequest())
	assert.True(t, c.queue.IsEmpty())
}

func (c *ClientDispatcherTestSuite) TestClientRequestContextCanceled() {
	t := c.T()
	// Setup
	callbackC := make(chan struct{}, 1)
	errMsg := "Request canceled by user"
	c.websocketClient.On("Write", mock.Anything).Return(nil)
	// Create mock request
	req := newMockRequest("somevalue")
	call, err := c.endpoint.CreateCall(req)
	require.NoError(t, err)
	requestID := call.UniqueId
	ctx, cancelFn := context.WithCancel(context.TODO())
	customCallback := func(response ocpp.Response, err error) {
		require.Nil(t, response)
		require.NotNil(t, err)
		ocppErr, ok := err.(*ocpp.Error)
		require.True(t, ok)
		assert.Equal(t, requestID, ocppErr.MessageId)
		assert.Equal(t, ocppj.GenericError, ocppErr.Code)
		assert.Equal(t, errMsg, ocppErr.Description)
		// Assert context
		assert.Equal(t, context.Canceled, ctx.Err())
		callbackC <- struct{}{}
	}
	bundle := ocppj.RequestBundle{Call: call, Callback: customCallback, Context: ctx}
	// Start and send mock request
	c.dispatcher.Start()
	require.True(t, c.dispatcher.IsRunning())
	err = c.dispatcher.SendRequest(bundle)
	require.NoError(t, err)
	// Check underlying queue
	time.Sleep(100 * time.Millisecond)
	assert.False(t, c.queue.IsEmpty())
	assert.Equal(t, 1, c.queue.Size())
	assert.True(t, c.state.HasPendingRequest())
	// Cancel request and handle callback
	cancelFn()
	_, ok := <-callbackC
	require.True(t, ok)
	assert.False(t, c.state.HasPendingRequest())
	assert.True(t, c.queue.IsEmpty())
}

func (c *ClientDispatcherTestSuite) TestClientDispatcherTimeout() {
	t := c.T()
	// Setup
	writeC := make(chan bool, 1)
	timeout := make(chan bool, 1)
	c.websocketClient.On("Write", mock.Anything).Run(func(args mock.Arguments) {
		writeC <- true
	}).Return(nil)
	// Create mock request
	req := newMockRequest("somevalue")
	call, err := c.endpoint.CreateCall(req)
	require.NoError(t, err)
	requestID := call.UniqueId
	bundle := ocppj.RequestBundle{Call: call, Callback: nil, Context: context.TODO()}
	// Set low timeout to trigger OnRequestCanceled callback
	c.dispatcher.SetOnRequestCanceled(func(reqBundle ocppj.RequestBundle, err error) {
		assert.Equal(t, requestID, reqBundle.Call.GetUniqueId())
		assert.Equal(t, MockFeatureName, reqBundle.Call.Payload.GetFeatureName())
		assert.Equal(t, req, reqBundle.Call.Payload)
		ocppErr, ok := err.(*ocpp.Error)
		require.True(t, ok)
		assert.Equal(t, ocppj.GenericError, ocppErr.Code)
		assert.Equal(t, "Request timed out", ocppErr.Description)
		timeout <- true
	})
	c.dispatcher.SetTimeout(1 * time.Second)
	c.dispatcher.Start()
	require.True(t, c.dispatcher.IsRunning())
	// Send mocked request
	err = c.dispatcher.SendRequest(bundle)
	require.NoError(t, err)
	// Check status after sending request
	<-writeC
	assert.True(t, c.state.HasPendingRequest())
	// Wait for timeout
	_, ok := <-timeout
	assert.True(t, ok)
	assert.False(t, c.state.HasPendingRequest())
	assert.True(t, c.queue.IsEmpty())
}

func (c *ClientDispatcherTestSuite) TestClientDispatcherContextTimeout() {
	t := c.T()
	// Setup
	callbackC := make(chan struct{}, 1)
	errMsg := "Request timed out"
	c.websocketClient.On("Write", mock.Anything).Return(nil)
	// Create mock request
	req := newMockRequest("somevalue")
	call, err := c.endpoint.CreateCall(req)
	require.NoError(t, err)
	requestID := call.UniqueId
	customCallback := func(response ocpp.Response, err error) {
		require.Nil(t, response)
		require.NotNil(t, err)
		ocppErr, ok := err.(*ocpp.Error)
		require.True(t, ok)
		assert.Equal(t, requestID, ocppErr.MessageId)
		assert.Equal(t, ocppj.GenericError, ocppErr.Code)
		assert.Equal(t, errMsg, ocppErr.Description)
		callbackC <- struct{}{}
	}
	bundle := ocppj.RequestBundle{Call: call, Callback: customCallback, Context: context.TODO()}
	// Set low timeout, start and send mock request
	c.dispatcher.SetTimeout(500 * time.Millisecond)
	c.dispatcher.Start()
	require.True(t, c.dispatcher.IsRunning())
	err = c.dispatcher.SendRequest(bundle)
	require.NoError(t, err)
	// Check underlying queue
	time.Sleep(100 * time.Millisecond)
	assert.False(t, c.queue.IsEmpty())
	assert.Equal(t, 1, c.queue.Size())
	assert.True(t, c.state.HasPendingRequest())
	// Cancel request and handle callback
	_, ok := <-callbackC
	require.True(t, ok)
	assert.False(t, c.state.HasPendingRequest())
	assert.True(t, c.queue.IsEmpty())
}

func (c *ClientDispatcherTestSuite) TestClientPauseDispatcher() {
	t := c.T()
	// Create mock request
	timeout := make(chan bool, 1)
	c.websocketClient.On("Write", mock.Anything).Return(nil)
	req := newMockRequest("somevalue")
	call, err := c.endpoint.CreateCall(req)
	require.NoError(t, err)
	// The callback will only be triggered at the end of the test case
	customCallback := func(response ocpp.Response, err error) {
		require.Nil(t, response)
		require.NotNil(t, err)
		timeout <- true
	}
	bundle := ocppj.RequestBundle{Call: call, Callback: customCallback, Context: context.TODO()}
	// Set timeout to test pause functionality, then start
	c.dispatcher.SetTimeout(500 * time.Millisecond)
	c.dispatcher.Start()
	require.True(t, c.dispatcher.IsRunning())
	// The callback will only be triggered at the end of the test case
	err = c.dispatcher.SendRequest(bundle)
	require.NoError(t, err)
	// Pause and attempt retransmission 2 times
	for i := 0; i < 2; i++ {
		time.Sleep(200 * time.Millisecond)
		// Pause dispatcher
		c.dispatcher.Pause()
		time.Sleep(100 * time.Millisecond)
		assert.True(t, c.dispatcher.IsPaused())
		// Elapsed time since start ~ 1 second, no timeout should be triggered (set to 0.5 seconds)
		time.Sleep(700 * time.Millisecond)
		assert.True(t, c.state.HasPendingRequest())
		assert.False(t, c.queue.IsEmpty())
		// Resume and restart transmission timer
		c.dispatcher.Resume()
		time.Sleep(50 * time.Millisecond)
		assert.False(t, c.dispatcher.IsPaused())
	}
	// Wait for timeout
	_, ok := <-timeout
	assert.True(t, ok)
	assert.False(t, c.state.HasPendingRequest())
	assert.True(t, c.queue.IsEmpty())
}

func (c *ClientDispatcherTestSuite) TestClientSendPausedDispatcher() {
	t := c.T()
	// Create mock request
	c.websocketClient.On("Write", mock.Anything).Run(func(args mock.Arguments) {
		require.Fail(t, "write should never be called")
	}).Return(nil)
	// Set timeout (unused for this test)
	c.dispatcher.SetTimeout(1 * time.Second)
	// Start
	c.dispatcher.Start()
	require.True(t, c.dispatcher.IsRunning())
	// The callback will only be triggered at the end of the test case
	c.dispatcher.SetOnRequestCanceled(func(reqBundle ocppj.RequestBundle, err error) {
		require.Fail(t, "unexpected onRequestCanceled")
	})
	// Pause, then send request
	c.dispatcher.Pause()
	assert.False(t, c.state.HasPendingRequest())
	assert.True(t, c.queue.IsEmpty())
	requestIDs := []string{}
	requestNumber := 2
	for i := 0; i < requestNumber; i++ {
		req := newMockRequest("somevalue")
		call, err := c.endpoint.CreateCall(req)
		require.NoError(t, err)
		requestID := call.UniqueId
		bundle := ocppj.RequestBundle{Call: call, Callback: nil, Context: context.TODO()}
		err = c.dispatcher.SendRequest(bundle)
		require.NoError(t, err)
		requestIDs = append(requestIDs, requestID)
	}
	time.Sleep(500 * time.Millisecond)
	// Request is queued
	assert.Equal(t, requestNumber, c.queue.Size())
	assert.False(t, c.state.HasPendingRequest())
	// After waiting for some time, no timeout was triggered and no pending requests
	time.Sleep(1 * time.Second)
	assert.Equal(t, requestNumber, c.queue.Size())
	assert.False(t, c.state.HasPendingRequest())
}
