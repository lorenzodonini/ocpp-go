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

func (s *ServerDispatcherTestSuite) TestServerSendRequest() {
	t := s.T()
	// Setup
	clientID := "client1"
	sent := make(chan bool, 1)
	s.websocketServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Run(func(args mock.Arguments) {
		id, _ := args.Get(0).(string)
		assert.Equal(t, clientID, id)
		sent <- true
	}).Return(nil)
	timeout := time.Second * 1
	s.dispatcher.SetTimeout(timeout)
	s.dispatcher.SetOnRequestCanceled(func(cID string, rID string, request ocpp.Request, err *ocpp.Error) {
		require.Fail(t, "unexpected OnRequestCanceled")
	})
	s.dispatcher.Start()
	require.True(t, s.dispatcher.IsRunning())
	// Simulate client connection
	s.dispatcher.CreateClient(clientID)
	// Create and send mock request
	req := newMockRequest("somevalue")
	call, err := s.endpoint.CreateCall(req)
	require.NoError(t, err)
	requestID := call.UniqueId
	bundle := ocppj.RequestBundle{Call: call}
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
	// Assert that no timeout is invoked
	time.Sleep(1300 * time.Millisecond)
}

func (s *ServerDispatcherTestSuite) TestServerRequestCanceled() {
	t := s.T()
	// Setup
	clientID := "client1"
	canceled := make(chan bool, 1)
	writeC := make(chan bool, 1)
	errMsg := "mockError"
	// Mock write error to trigger onRequestCanceled
	// This never starts a timeout
	s.websocketServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Run(func(args mock.Arguments) {
		id, _ := args.Get(0).(string)
		assert.Equal(t, clientID, id)
		<-writeC
	}).Return(fmt.Errorf(errMsg))
	// Create mock request
	req := newMockRequest("somevalue")
	call, err := s.endpoint.CreateCall(req)
	require.NoError(t, err)
	requestID := call.UniqueId
	bundle := ocppj.RequestBundle{Call: call}
	// Set canceled callback
	s.dispatcher.SetOnRequestCanceled(func(cID string, rID string, request ocpp.Request, err *ocpp.Error) {
		assert.Equal(t, clientID, cID)
		assert.Equal(t, requestID, rID)
		assert.Equal(t, MockFeatureName, request.GetFeatureName())
		assert.Equal(t, req, request)
		assert.Equal(t, ocppj.InternalError, err.Code)
		assert.Equal(t, errMsg, err.Description)
		canceled <- true
	})
	s.dispatcher.Start()
	require.True(t, s.dispatcher.IsRunning())
	// Simulate client connection
	s.dispatcher.CreateClient(clientID)
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

func (s *ServerDispatcherTestSuite) TestCreateClient() {
	t := s.T()
	// Setup
	clientID := "client1"
	s.dispatcher.Start()
	require.True(t, s.dispatcher.IsRunning())
	// No client state created yet
	_, ok := s.queueMap.Get(clientID)
	assert.False(t, ok)
	// Create client state
	s.dispatcher.CreateClient(clientID)
	_, ok = s.queueMap.Get(clientID)
	assert.True(t, ok)
	assert.False(t, s.state.HasPendingRequest(clientID))
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
	// Simulate client connection
	s.dispatcher.CreateClient(clientID)
	// Create and send mock request
	req := newMockRequest("somevalue")
	call, err := s.endpoint.CreateCall(req)
	require.NoError(t, err)
	bundle := ocppj.RequestBundle{Call: call}
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

func (s *ServerDispatcherTestSuite) TestServerDispatcherTimeout() {
	t := s.T()
	// Setup
	clientID := "client1"
	canceled := make(chan bool, 1)
	s.websocketServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Run(func(args mock.Arguments) {
		id, _ := args.Get(0).(string)
		assert.Equal(t, clientID, id)
	}).Return(nil)
	// Create mock request
	req := newMockRequest("somevalue")
	call, err := s.endpoint.CreateCall(req)
	require.NoError(t, err)
	requestID := call.UniqueId
	bundle := ocppj.RequestBundle{Call: call}
	// Set canceled callback
	s.dispatcher.SetOnRequestCanceled(func(cID string, rID string, request ocpp.Request, err *ocpp.Error) {
		assert.Equal(t, clientID, cID)
		assert.Equal(t, requestID, rID)
		assert.Equal(t, MockFeatureName, request.GetFeatureName())
		assert.Equal(t, req, request)
		assert.Equal(t, ocppj.GenericError, err.Code)
		assert.Equal(t, "Request timed out", err.Description)
		canceled <- true
	})
	// Set timeout and start
	timeout := time.Second * 1
	s.dispatcher.SetTimeout(timeout)
	s.dispatcher.Start()
	require.True(t, s.dispatcher.IsRunning())
	// Simulate client connection
	s.dispatcher.CreateClient(clientID)
	// Send mock request
	startTime := time.Now()
	err = s.dispatcher.SendRequest(clientID, bundle)
	require.NoError(t, err)
	// Wait for timeout, canceled callback will be invoked
	_, ok := <-canceled
	assert.True(t, ok)
	elapsed := time.Since(startTime)
	assert.GreaterOrEqual(t, elapsed.Seconds(), timeout.Seconds())
	clientQ, _ := s.queueMap.Get(clientID)
	assert.True(t, clientQ.IsEmpty())
}
