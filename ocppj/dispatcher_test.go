package ocppj_test

import (
	"errors"
	"sync"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel/metric/noop"

	"github.com/xBlaz3kx/ocpp-go/ocpp"
	"github.com/xBlaz3kx/ocpp-go/ocppj"
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
	mockProfile := ocpp.NewProfile("mock", &MockFeature{})
	s.endpoint.AddProfile(mockProfile)
	s.queueMap = ocppj.NewFIFOQueueMap(10)
	s.dispatcher = ocppj.NewDefaultServerDispatcher(s.queueMap, noop.NewMeterProvider(), nil)
	s.state = ocppj.NewServerState(&s.mutex)
	s.dispatcher.SetPendingRequestState(s.state)
	s.websocketServer = MockWebsocketServer{}
	s.dispatcher.SetNetworkServer(&s.websocketServer)
}

func (s *ServerDispatcherTestSuite) TestServerSendRequest() {
	// Setup
	clientID := "client1"
	sent := make(chan bool, 1)
	s.websocketServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Run(func(args mock.Arguments) {
		id, _ := args.Get(0).(string)
		s.Assert().Equal(clientID, id)
		sent <- true
	}).Return(nil)
	timeout := time.Second * 1
	s.dispatcher.SetTimeout(timeout)
	s.dispatcher.SetOnRequestCanceled(func(cID string, rID string, request ocpp.Request, err *ocpp.Error) {
		s.Require().Fail("unexpected OnRequestCanceled")
	})
	s.dispatcher.Start()
	s.Require().True(s.dispatcher.IsRunning())
	// Simulate client connection
	s.dispatcher.CreateClient(clientID)
	// Create and send mock request
	req := newMockRequest("somevalue")
	call, err := s.endpoint.CreateCall(req)
	s.Require().NoError(err)
	requestID := call.UniqueId
	data, err := call.MarshalJSON()
	s.Require().NoError(err)
	bundle := ocppj.RequestBundle{Call: call, Data: data}
	err = s.dispatcher.SendRequest(clientID, bundle)
	s.Require().NoError(err)
	// Check underlying queue
	q, ok := s.queueMap.Get(clientID)
	s.Require().True(ok)
	s.Assert().False(q.IsEmpty())
	s.Assert().Equal(1, q.Size())
	// Wait for websocket to send message
	_, ok = <-sent
	s.Assert().True(ok)
	s.Assert().True(s.state.HasPendingRequest(clientID))
	// Complete request
	s.dispatcher.CompleteRequest(clientID, requestID)
	s.Assert().False(s.state.HasPendingRequest(clientID))
	s.Assert().True(q.IsEmpty())
	// Assert that no timeout is invoked
	time.Sleep(1300 * time.Millisecond)
}

func (s *ServerDispatcherTestSuite) TestServerRequestCanceled() {
	// Setup
	clientID := "client1"
	canceled := make(chan bool, 1)
	writeC := make(chan bool, 1)
	errMsg := "mockError"
	// Mock write error to trigger onRequestCanceled
	// This never starts a timeout
	s.websocketServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Run(func(args mock.Arguments) {
		id, _ := args.Get(0).(string)
		s.Assert().Equal(clientID, id)
		<-writeC
	}).Return(errors.New(errMsg))
	// Create mock request
	req := newMockRequest("somevalue")
	call, err := s.endpoint.CreateCall(req)
	s.Require().NoError(err)
	requestID := call.UniqueId
	data, err := call.MarshalJSON()
	s.Require().NoError(err)
	bundle := ocppj.RequestBundle{Call: call, Data: data}
	// Set canceled callback
	s.dispatcher.SetOnRequestCanceled(func(cID string, rID string, request ocpp.Request, err *ocpp.Error) {
		s.Assert().Equal(clientID, cID)
		s.Assert().Equal(requestID, rID)
		s.Assert().Equal(MockFeatureName, request.GetFeatureName())
		s.Assert().Equal(req, request)
		s.Assert().Equal(ocppj.InternalError, err.Code)
		s.Assert().Equal(errMsg, err.Description)
		canceled <- true
	})
	s.dispatcher.Start()
	s.Require().True(s.dispatcher.IsRunning())
	// Simulate client connection
	s.dispatcher.CreateClient(clientID)
	// Send mock request
	err = s.dispatcher.SendRequest(clientID, bundle)
	s.Require().NoError(err)
	// Check underlying queue
	time.Sleep(100 * time.Millisecond)
	q, ok := s.queueMap.Get(clientID)
	s.Require().True(ok)
	s.Assert().False(q.IsEmpty())
	s.Assert().Equal(1, q.Size())
	s.Assert().True(s.state.HasPendingRequest(clientID))
	// Signal that write can occur now, then check canceled request
	writeC <- true
	_, ok = <-canceled
	s.Require().True(ok)
	s.Assert().False(s.state.HasPendingRequest(clientID))
	s.Assert().True(q.IsEmpty())
}

func (s *ServerDispatcherTestSuite) TestCreateClient() {
	// Setup
	clientID := "client1"
	s.dispatcher.Start()
	s.Require().True(s.dispatcher.IsRunning())
	// No client state created yet
	_, ok := s.queueMap.Get(clientID)
	s.Assert().False(ok)
	// Create client state
	s.dispatcher.CreateClient(clientID)
	_, ok = s.queueMap.Get(clientID)
	s.Assert().True(ok)
	s.Assert().False(s.state.HasPendingRequest(clientID))
}

func (s *ServerDispatcherTestSuite) TestDeleteClient() {
	// Setup
	clientID := "client1"
	sent := make(chan bool, 1)
	s.websocketServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Run(func(args mock.Arguments) {
		id, _ := args.Get(0).(string)
		s.Assert().Equal(clientID, id)
		sent <- true
	}).Return(nil)
	s.dispatcher.Start()
	s.Require().True(s.dispatcher.IsRunning())
	// Simulate client connection
	s.dispatcher.CreateClient(clientID)
	// Create and send mock request
	req := newMockRequest("somevalue")
	call, err := s.endpoint.CreateCall(req)
	s.Require().NoError(err)
	data, err := call.MarshalJSON()
	s.Require().NoError(err)
	bundle := ocppj.RequestBundle{Call: call, Data: data}
	err = s.dispatcher.SendRequest(clientID, bundle)
	s.Require().NoError(err)
	// Wait for websocket to send message
	_, ok := <-sent
	s.Assert().True(ok)
	// Delete client
	s.dispatcher.DeleteClient(clientID)
	// Pending request is still expected to be there
	s.Assert().True(s.state.HasPendingRequest(clientID))
}

func (s *ServerDispatcherTestSuite) TestServerDispatcherTimeout() {
	// Setup
	clientID := "client1"
	canceled := make(chan bool, 1)
	s.websocketServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Run(func(args mock.Arguments) {
		id, _ := args.Get(0).(string)
		s.Assert().Equal(clientID, id)
	}).Return(nil)
	// Create mock request
	req := newMockRequest("somevalue")
	call, err := s.endpoint.CreateCall(req)
	s.Require().NoError(err)
	requestID := call.UniqueId
	data, err := call.MarshalJSON()
	s.Require().NoError(err)
	bundle := ocppj.RequestBundle{Call: call, Data: data}
	// Set canceled callback
	s.dispatcher.SetOnRequestCanceled(func(cID string, rID string, request ocpp.Request, err *ocpp.Error) {
		s.Assert().Equal(clientID, cID)
		s.Assert().Equal(requestID, rID)
		s.Assert().Equal(MockFeatureName, request.GetFeatureName())
		s.Assert().Equal(req, request)
		s.Assert().Equal(ocppj.GenericError, err.Code)
		s.Assert().Equal("Request timed out", err.Description)
		canceled <- true
	})
	// Set timeout and start
	timeout := time.Second * 1
	s.dispatcher.SetTimeout(timeout)
	s.dispatcher.Start()
	s.Require().True(s.dispatcher.IsRunning())
	// Simulate client connection
	s.dispatcher.CreateClient(clientID)
	// Send mock request
	startTime := time.Now()
	err = s.dispatcher.SendRequest(clientID, bundle)
	s.Require().NoError(err)
	// Wait for timeout, canceled callback will be invoked
	_, ok := <-canceled
	s.Assert().True(ok)
	elapsed := time.Since(startTime)
	s.Assert().GreaterOrEqual(elapsed.Seconds(), timeout.Seconds())
	clientQ, _ := s.queueMap.Get(clientID)
	s.Assert().True(clientQ.IsEmpty())
}

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
	mockProfile := ocpp.NewProfile("mock", &MockFeature{})
	c.endpoint.AddProfile(mockProfile)
	c.queue = ocppj.NewFIFOClientQueue(10)
	c.dispatcher = ocppj.NewDefaultClientDispatcher(c.queue, nil)
	c.state = ocppj.NewClientState()
	c.dispatcher.SetPendingRequestState(c.state)
	c.websocketClient = MockWebsocketClient{}
	c.dispatcher.SetNetworkClient(&c.websocketClient)
}

func (c *ClientDispatcherTestSuite) TestClientSendRequest() {
	// Setup
	sent := make(chan bool, 1)
	c.websocketClient.On("Write", mock.Anything).Run(func(args mock.Arguments) {
		sent <- true
	}).Return(nil)
	c.dispatcher.Start()
	c.Require().True(c.dispatcher.IsRunning())
	// Create and send mock request
	req := newMockRequest("somevalue")
	call, err := c.endpoint.CreateCall(req)
	c.Require().NoError(err)
	requestID := call.UniqueId
	data, err := call.MarshalJSON()
	c.Require().NoError(err)
	bundle := ocppj.RequestBundle{Call: call, Data: data}
	err = c.dispatcher.SendRequest(bundle)
	c.Require().NoError(err)
	// Check underlying queue
	c.Assert().False(c.queue.IsEmpty())
	c.Assert().Equal(1, c.queue.Size())
	// Wait for websocket to send message
	_, ok := <-sent
	c.Assert().True(ok)
	c.Assert().True(c.state.HasPendingRequest())
	// Complete request
	c.dispatcher.CompleteRequest(requestID)
	c.Assert().False(c.state.HasPendingRequest())
	c.Assert().True(c.queue.IsEmpty())

}

func (c *ClientDispatcherTestSuite) TestClientSendEvent() {
	// Setup
	sent := make(chan bool, 1)
	c.websocketClient.On("Write", mock.Anything).Run(func(args mock.Arguments) {
		sent <- true
	}).Return(nil)
	c.dispatcher.Start()
	c.Require().True(c.dispatcher.IsRunning())

	// Create and send mock request
	req := newMockRequest("somevalue")

	call, err := c.endpoint.CreateSend(req)
	c.Require().NoError(err)

	data, err := call.MarshalJSON()
	c.Require().NoError(err)

	bundle := ocppj.RequestBundle{Call: (*ocppj.Call)(call), Data: data}
	err = c.dispatcher.SendRequest(bundle)
	c.Require().NoError(err)

	// Check underlying queue
	c.False(c.queue.IsEmpty())
	c.Equal(1, c.queue.Size())

	// Wait for websocket to send message
	_, ok := <-sent
	c.True(ok)
	c.True(c.state.HasPendingRequest())
}

func (c *ClientDispatcherTestSuite) TestClientRequestCanceled() {
	// Setup
	canceled := make(chan bool, 1)
	writeC := make(chan bool, 1)
	errMsg := "mockError"
	c.websocketClient.On("Write", mock.Anything).Run(func(args mock.Arguments) {
		<-writeC
	}).Return(errors.New(errMsg))
	// Create mock request
	req := newMockRequest("somevalue")
	call, err := c.endpoint.CreateCall(req)
	c.Require().NoError(err)
	requestID := call.UniqueId
	data, err := call.MarshalJSON()
	c.Require().NoError(err)
	bundle := ocppj.RequestBundle{Call: call, Data: data}
	// Set canceled callback
	c.dispatcher.SetOnRequestCanceled(func(rID string, request ocpp.Request, err *ocpp.Error) {
		c.Assert().Equal(requestID, rID)
		c.Assert().Equal(MockFeatureName, request.GetFeatureName())
		c.Assert().Equal(req, request)
		c.Assert().Equal(ocppj.InternalError, err.Code)
		c.Assert().Equal(errMsg, err.Description)
		canceled <- true
	})
	c.dispatcher.Start()
	c.Require().True(c.dispatcher.IsRunning())
	// Send mock request
	err = c.dispatcher.SendRequest(bundle)
	c.Require().NoError(err)
	// Check underlying queue
	time.Sleep(100 * time.Millisecond)
	c.Assert().False(c.queue.IsEmpty())
	c.Assert().Equal(1, c.queue.Size())
	c.Assert().True(c.state.HasPendingRequest())
	// Signal that write can occur now, then check canceled request
	writeC <- true
	_, ok := <-canceled
	c.Require().True(ok)
	c.Assert().False(c.state.HasPendingRequest())
	c.Assert().True(c.queue.IsEmpty())
}

func (c *ClientDispatcherTestSuite) TestClientDispatcherTimeout() {
	// Setup
	writeC := make(chan bool, 1)
	timeout := make(chan bool, 1)
	c.websocketClient.On("Write", mock.Anything).Run(func(args mock.Arguments) {
		writeC <- true
	}).Return(nil)
	// Create mock request
	req := newMockRequest("somevalue")
	call, err := c.endpoint.CreateCall(req)
	c.Require().NoError(err)
	requestID := call.UniqueId
	data, err := call.MarshalJSON()
	c.Require().NoError(err)
	bundle := ocppj.RequestBundle{Call: call, Data: data}
	// Set low timeout to trigger OnRequestCanceled callback
	c.dispatcher.SetTimeout(1 * time.Second)
	c.dispatcher.SetOnRequestCanceled(func(rID string, request ocpp.Request, err *ocpp.Error) {
		c.Assert().Equal(requestID, rID)
		c.Assert().Equal(MockFeatureName, request.GetFeatureName())
		c.Assert().Equal(req, request)
		c.Assert().Equal(ocppj.GenericError, err.Code)
		c.Assert().Equal("Request timed out", err.Description)
		timeout <- true
	})
	c.dispatcher.Start()
	c.Require().True(c.dispatcher.IsRunning())
	// Send mocked request
	err = c.dispatcher.SendRequest(bundle)
	c.Require().NoError(err)
	// Check status after sending request
	<-writeC
	c.Assert().True(c.state.HasPendingRequest())
	// Wait for timeout
	_, ok := <-timeout
	c.Assert().True(ok)
	c.Assert().False(c.state.HasPendingRequest())
	c.Assert().True(c.queue.IsEmpty())
}

func (c *ClientDispatcherTestSuite) TestClientPauseDispatcher() {
	// Create mock request
	timeout := make(chan bool, 1)
	c.websocketClient.On("Write", mock.Anything).Return(nil)
	req := newMockRequest("somevalue")
	call, err := c.endpoint.CreateCall(req)
	c.Require().NoError(err)
	requestID := call.UniqueId
	data, err := call.MarshalJSON()
	c.Require().NoError(err)
	bundle := ocppj.RequestBundle{Call: call, Data: data}
	// Set timeout to test pause functionality
	c.dispatcher.SetTimeout(500 * time.Millisecond)
	// The callback will only be triggered at the end of the test case
	c.dispatcher.SetOnRequestCanceled(func(rID string, request ocpp.Request, err *ocpp.Error) {
		c.Assert().Equal(requestID, rID)
		c.Assert().Equal(MockFeatureName, request.GetFeatureName())
		c.Assert().Equal(req, request)
		timeout <- true
	})
	c.dispatcher.Start()
	c.Require().True(c.dispatcher.IsRunning())
	err = c.dispatcher.SendRequest(bundle)
	c.Require().NoError(err)
	// Pause and attempt retransmission 2 times
	for i := 0; i < 2; i++ {
		time.Sleep(200 * time.Millisecond)
		// Pause dispatcher
		c.dispatcher.Pause()
		c.Assert().True(c.dispatcher.IsPaused())
		// Elapsed time since start ~ 1 second, no timeout should be triggered (set to 0.5 seconds)
		time.Sleep(800 * time.Millisecond)
		c.Assert().True(c.state.HasPendingRequest())
		c.Assert().False(c.queue.IsEmpty())
		// Resume and restart transmission timer
		c.dispatcher.Resume()
		c.Assert().False(c.dispatcher.IsPaused())
	}
	// Wait for timeout
	_, ok := <-timeout
	c.Assert().True(ok)
	c.Assert().False(c.state.HasPendingRequest())
	c.Assert().True(c.queue.IsEmpty())
}

func (c *ClientDispatcherTestSuite) TestClientSendPausedDispatcher() {
	// Create mock request
	c.websocketClient.On("Write", mock.Anything).Run(func(args mock.Arguments) {
		c.Require().Fail("write should never be called")
	}).Return(nil)
	// Set timeout (unused for this test)
	c.dispatcher.SetTimeout(1 * time.Second)
	// The callback will only be triggered at the end of the test case
	c.dispatcher.SetOnRequestCanceled(func(rID string, request ocpp.Request, err *ocpp.Error) {
		c.Require().Fail("unexpected OnRequestCanceled")
	})
	c.dispatcher.Start()
	c.Require().True(c.dispatcher.IsRunning())
	// Pause, then send request
	c.dispatcher.Pause()
	c.Assert().False(c.state.HasPendingRequest())
	c.Assert().True(c.queue.IsEmpty())
	requestIDs := []string{}
	requestNumber := 2
	for i := 0; i < requestNumber; i++ {
		req := newMockRequest("somevalue")
		call, err := c.endpoint.CreateCall(req)
		c.Require().NoError(err)
		requestID := call.UniqueId
		data, err := call.MarshalJSON()
		c.Require().NoError(err)
		bundle := ocppj.RequestBundle{Call: call, Data: data}
		err = c.dispatcher.SendRequest(bundle)
		c.Require().NoError(err)
		requestIDs = append(requestIDs, requestID)
	}
	time.Sleep(500 * time.Millisecond)
	// Request is queued
	c.Assert().Equal(requestNumber, c.queue.Size())
	c.Assert().False(c.state.HasPendingRequest())
	// After waiting for some time, no timeout was triggered and no pending requests
	time.Sleep(1 * time.Second)
	c.Assert().Equal(requestNumber, c.queue.Size())
	c.Assert().False(c.state.HasPendingRequest())
}
