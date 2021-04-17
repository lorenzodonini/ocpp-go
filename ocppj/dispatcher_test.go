package ocppj_test

import (
	"fmt"
	"sync"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocppj"
)

type ServerDispatcherTestSuite struct {
	suite.Suite
	mutex           sync.RWMutex
	state           ocppj.ServerState
	websocketServer MockWebsocketServer
	endpoint        ocppj.Server
	dispatcher      ocppj.ServerDispatcher
	queueMap        ocppj.ServerQueueMap
}

func (s *ServerDispatcherTestSuite) SetupTest() {
	s.endpoint = ocppj.Server{}
	mockProfile := ocpp.NewProfile("mock", MockFeature{})
	s.endpoint.AddProfile(mockProfile)
	s.queueMap = ocppj.NewFIFOQueueMap(10)
	s.dispatcher = ocppj.NewDefaultServerDispatcher(s.queueMap)
	s.state = ocppj.NewServerState(&s.mutex)
	s.dispatcher.SetPendingRequestState(s.state)
	s.websocketServer = MockWebsocketServer{}
	s.dispatcher.SetNetworkServer(&s.websocketServer)
}

func (s *ServerDispatcherTestSuite) TestSendRequest() {
	t := s.T()
	// Setup
	clientID := "client1"
	sent := make(chan bool, 1)
	s.websocketServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Run(func(args mock.Arguments) {
		id, _ := args.Get(0).(string)
		assert.Equal(t, clientID, id)
		sent <- true
	}).Return(nil)
	s.dispatcher.Start()
	require.True(t, s.dispatcher.IsRunning())
	// Create and send mock request
	req := newMockRequest("somevalue")
	call, err := s.endpoint.CreateCall(req)
	require.NoError(t, err)
	requestID := call.UniqueId
	data, err := call.MarshalJSON()
	require.NoError(t, err)
	bundle := ocppj.RequestBundle{Call: call, Data: data}
	err = s.dispatcher.SendRequest(clientID, bundle)
	require.NoError(t, err)
	// Check underlying queue
	q, ok := s.queueMap.Get(clientID)
	require.True(t, ok)
	assert.False(t, q.IsEmpty())
	assert.Equal(t, 1, q.Size())
	// Wait for websocket to send message
	_, ok = <-sent
	assert.True(t, ok)
	assert.True(t, s.state.HasPendingRequest(clientID))
	// Complete request
	s.dispatcher.CompleteRequest(clientID, requestID)
	assert.False(t, s.state.HasPendingRequest(clientID))
	assert.True(t, q.IsEmpty())
}

func (s *ServerDispatcherTestSuite) TestRequestCanceled() {
	t := s.T()
	// Setup
	clientID := "client1"
	canceled := make(chan bool, 1)
	writeC := make(chan bool, 1)
	s.websocketServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Run(func(args mock.Arguments) {
		id, _ := args.Get(0).(string)
		assert.Equal(t, clientID, id)
		_, _ = <-writeC
	}).Return(fmt.Errorf("mockError"))
	// Create mock request
	req := newMockRequest("somevalue")
	call, err := s.endpoint.CreateCall(req)
	require.NoError(t, err)
	requestID := call.UniqueId
	data, err := call.MarshalJSON()
	require.NoError(t, err)
	bundle := ocppj.RequestBundle{Call: call, Data: data}
	// Set canceled callback
	s.dispatcher.SetOnRequestCanceled(func(client string, request string) {
		assert.Equal(t, clientID, client)
		assert.Equal(t, requestID, request)
		canceled <- true
	})
	s.dispatcher.Start()
	require.True(t, s.dispatcher.IsRunning())
	// Send mock request
	err = s.dispatcher.SendRequest(clientID, bundle)
	require.NoError(t, err)
	// Check underlying queue
	time.Sleep(100 * time.Millisecond)
	q, ok := s.queueMap.Get(clientID)
	require.True(t, ok)
	assert.False(t, q.IsEmpty())
	assert.Equal(t, 1, q.Size())
	assert.True(t, s.state.HasPendingRequest(clientID))
	// Signal that write can occur now, then check canceled request
	writeC <- true
	_, ok = <-canceled
	require.True(t, ok)
	assert.False(t, s.state.HasPendingRequest(clientID))
	assert.True(t, q.IsEmpty())
}

func (s *ServerDispatcherTestSuite) TestDeleteClient() {
	t := s.T()
	// Setup
	clientID := "client1"
	sent := make(chan bool, 1)
	s.websocketServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Run(func(args mock.Arguments) {
		id, _ := args.Get(0).(string)
		assert.Equal(t, clientID, id)
		sent <- true
	}).Return(nil)
	s.dispatcher.Start()
	require.True(t, s.dispatcher.IsRunning())
	// Create and send mock request
	req := newMockRequest("somevalue")
	call, err := s.endpoint.CreateCall(req)
	require.NoError(t, err)
	data, err := call.MarshalJSON()
	require.NoError(t, err)
	bundle := ocppj.RequestBundle{Call: call, Data: data}
	err = s.dispatcher.SendRequest(clientID, bundle)
	require.NoError(t, err)
	// Wait for websocket to send message
	_, ok := <-sent
	assert.True(t, ok)
	// Delete client
	s.dispatcher.DeleteClient(clientID)
	// Pending request is still expected to be there
	assert.True(t, s.state.HasPendingRequest(clientID))
}

type ClientDispatcherTestSuite struct {
	suite.Suite
	mutex           sync.Mutex
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
	c.dispatcher.SetNetworkClient(&c.websocketClient)
}

func (c *ClientDispatcherTestSuite) TestSendRequest() {
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
	requestID := call.UniqueId
	data, err := call.MarshalJSON()
	require.NoError(t, err)
	bundle := ocppj.RequestBundle{Call: call, Data: data}
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

func (c *ClientDispatcherTestSuite) TestRequestCanceled() {
	t := c.T()
	// Setup
	canceled := make(chan bool, 1)
	writeC := make(chan bool, 1)
	c.websocketClient.On("Write", mock.Anything).Run(func(args mock.Arguments) {
		_, _ = <-writeC
	}).Return(fmt.Errorf("mockError"))
	// Create mock request
	req := newMockRequest("somevalue")
	call, err := c.endpoint.CreateCall(req)
	require.NoError(t, err)
	requestID := call.UniqueId
	data, err := call.MarshalJSON()
	require.NoError(t, err)
	bundle := ocppj.RequestBundle{Call: call, Data: data}
	// Set canceled callback
	c.dispatcher.SetOnRequestCanceled(func(request string) {
		assert.Equal(t, requestID, request)
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

func (c *ClientDispatcherTestSuite) TestDispatcherTimeout() {
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
	data, err := call.MarshalJSON()
	require.NoError(t, err)
	bundle := ocppj.RequestBundle{Call: call, Data: data}
	// Set low timeout to trigger OnRequestCanceled callback
	c.dispatcher.SetTimeout(1 * time.Second)
	c.dispatcher.SetOnRequestCanceled(func(reqID string) {
		assert.Equal(t, requestID, reqID)
		timeout <- true
	})
	c.dispatcher.Start()
	require.True(t, c.dispatcher.IsRunning())
	// Send mocked request
	err = c.dispatcher.SendRequest(bundle)
	require.NoError(t, err)
	// Check status after sending request
	_, _ = <-writeC
	assert.True(t, c.state.HasPendingRequest())
	// Wait for timeout
	_, ok := <-timeout
	assert.True(t, ok)
	assert.False(t, c.state.HasPendingRequest())
	assert.True(t, c.queue.IsEmpty())
}

func (c *ClientDispatcherTestSuite) TestPauseDispatcher() {
	t := c.T()
	// Create mock request
	timeout := make(chan bool, 1)
	c.websocketClient.On("Write", mock.Anything).Return(nil)
	req := newMockRequest("somevalue")
	call, err := c.endpoint.CreateCall(req)
	require.NoError(t, err)
	requestID := call.UniqueId
	data, err := call.MarshalJSON()
	require.NoError(t, err)
	bundle := ocppj.RequestBundle{Call: call, Data: data}
	// Set timeout to test pause functionality
	c.dispatcher.SetTimeout(500 * time.Millisecond)
	// The callback will only be triggered at the end of the test case
	c.dispatcher.SetOnRequestCanceled(func(reqID string) {
		assert.Equal(t, requestID, reqID)
		timeout <- true
	})
	c.dispatcher.Start()
	require.True(t, c.dispatcher.IsRunning())
	err = c.dispatcher.SendRequest(bundle)
	require.NoError(t, err)
	// Pause and attempt retransmission 2 times
	for i := 0; i < 2; i++ {
		time.Sleep(200 * time.Millisecond)
		// Pause dispatcher
		c.dispatcher.Pause()
		assert.True(t, c.dispatcher.IsPaused())
		// Elapsed time since start ~ 1 second, no timeout should be triggered (set to 0.5 seconds)
		time.Sleep(800 * time.Millisecond)
		assert.True(t, c.state.HasPendingRequest())
		assert.False(t, c.queue.IsEmpty())
		// Resume and restart transmission timer
		c.dispatcher.Resume()
		assert.False(t, c.dispatcher.IsPaused())
	}
	// Wait for timeout
	_, ok := <-timeout
	assert.True(t, ok)
	assert.False(t, c.state.HasPendingRequest())
	assert.True(t, c.queue.IsEmpty())
}

func (c *ClientDispatcherTestSuite) TestSendWhilePausedDispatcher() {
	t := c.T()
	// Create mock request
	c.websocketClient.On("Write", mock.Anything).Run(func(args mock.Arguments) {
		require.Fail(t, "write should never be called")
	}).Return(nil)
	// Set timeout (unused for this test)
	c.dispatcher.SetTimeout(1 * time.Second)
	// The callback will only be triggered at the end of the test case
	c.dispatcher.SetOnRequestCanceled(func(reqID string) {
		require.Fail(t, "unexpected OnRequestCanceled")
	})
	c.dispatcher.Start()
	require.True(t, c.dispatcher.IsRunning())
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
		data, err := call.MarshalJSON()
		require.NoError(t, err)
		bundle := ocppj.RequestBundle{Call: call, Data: data}
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
