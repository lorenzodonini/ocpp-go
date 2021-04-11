package ocppj_test

import (
	"sync"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/lorenzodonini/ocpp-go/ocppj"
)

type ClientStateTestSuite struct {
	suite.Suite
	state ocppj.ClientState
}

func (suite *ClientStateTestSuite) SetupTest() {
	suite.state = ocppj.NewSimpleClientState()
}

func (suite *ClientStateTestSuite) TestAddPendingRequest() {
	t := suite.T()
	requestID := "1234"
	req := newMockRequest("somevalue")
	require.False(t, suite.state.HasPendingRequest())
	suite.state.AddPendingRequest(requestID, req)
	require.True(t, suite.state.HasPendingRequest())
	r, exists := suite.state.GetPendingRequest(requestID)
	assert.True(t, exists)
	assert.Equal(t, req, r)
}

func (suite *ClientStateTestSuite) TestGetInvalidPendingRequest() {
	t := suite.T()
	requestID := "1234"
	suite.state.AddPendingRequest(requestID, newMockRequest("somevalue"))
	require.True(t, suite.state.HasPendingRequest())
	invalidRequestIDs := []string{"4321", "5678", "1230", "deadc0de"}
	// Nothing returned when querying for an unknown request ID
	for _, id := range invalidRequestIDs {
		r, exists := suite.state.GetPendingRequest(id)
		assert.False(t, exists)
		assert.Nil(t, r)
	}
}

func (suite *ClientStateTestSuite) TestAddMultiplePendingRequests() {
	t := suite.T()
	requestId1 := "1234"
	requestId2 := "5678"
	req1 := newMockRequest("somevalue1")
	req2 := newMockRequest("somevalue2")
	suite.state.AddPendingRequest(requestId1, req1)
	suite.state.AddPendingRequest(requestId2, req2)
	r, exists := suite.state.GetPendingRequest(requestId1)
	assert.True(t, exists)
	assert.NotNil(t, r)
	r, exists = suite.state.GetPendingRequest(requestId2)
	assert.False(t, exists)
	assert.Nil(t, r)
}

func (suite *ClientStateTestSuite) TestDeletePendingRequest() {
	t := suite.T()
	requestID := "1234"
	req := newMockRequest("somevalue")
	suite.state.AddPendingRequest(requestID, req)
	require.True(t, suite.state.HasPendingRequest())
	suite.state.DeletePendingRequest(requestID)
	// Previously added request is gone
	assert.False(t, suite.state.HasPendingRequest())
	r, exists := suite.state.GetPendingRequest(requestID)
	assert.False(t, exists)
	assert.Nil(t, r)
	// Deleting again has no effect
	suite.state.DeletePendingRequest(requestID)
	assert.False(t, suite.state.HasPendingRequest())
}

func (suite *ClientStateTestSuite) TestDeleteInvalidPendingRequest() {
	t := suite.T()
	requestID := "1234"
	req := newMockRequest("somevalue")
	suite.state.AddPendingRequest(requestID, req)
	require.True(t, suite.state.HasPendingRequest())
	suite.state.DeletePendingRequest("5678")
	// Previously added request is still there
	assert.True(t, suite.state.HasPendingRequest())
	r, exists := suite.state.GetPendingRequest(requestID)
	assert.True(t, exists)
	assert.NotNil(t, r)
}

func (suite *ClientStateTestSuite) TestClearPendingRequests() {
	t := suite.T()
	requestID := "1234"
	req := newMockRequest("somevalue")
	suite.state.AddPendingRequest(requestID, req)
	require.True(t, suite.state.HasPendingRequest())
	suite.state.ClearPendingRequests()
	// No more requests available in the struct
	assert.False(t, suite.state.HasPendingRequest())
}

type ServerStateTestSuite struct {
	suite.Suite
	mutex sync.RWMutex
	state ocppj.ServerState
}

func (suite *ServerStateTestSuite) SetupTest() {
	suite.state = ocppj.NewSimpleServerState(&suite.mutex)
}

func (suite *ServerStateTestSuite) TestAddPendingRequests() {
	t := suite.T()
	type testClientRequest struct {
		clientID  string
		requestID string
		request   *MockRequest
	}
	requests := []testClientRequest{
		{"client1", "0001", newMockRequest("somevalue1")},
		{"client2", "0002", newMockRequest("somevalue2")},
		{"client3", "0003", newMockRequest("somevalue3")},
	}
	for _, r := range requests {
		suite.state.AddPendingRequest(r.clientID, r.requestID, r.request)
	}
	require.True(t, suite.state.HasPendingRequests())
	for _, r := range requests {
		assert.True(t, suite.state.HasPendingRequest(r.clientID))
		req, exists := suite.state.GetClientState(r.clientID).GetPendingRequest(r.requestID)
		assert.True(t, exists)
		assert.Equal(t, r.request, req)
	}
}

func (suite *ServerStateTestSuite) TestGetInvalidPendingRequest() {
	t := suite.T()
	requestID := "1234"
	clientID := "client1"
	suite.state.AddPendingRequest(clientID, requestID, newMockRequest("somevalue"))
	require.True(t, suite.state.HasPendingRequest(clientID))
	invalidRequestIDs := []string{"4321", "5678", "1230", "deadc0de"}
	// Nothing returned when querying for an unknown request ID
	for _, id := range invalidRequestIDs {
		r, exists := suite.state.GetClientState(clientID).GetPendingRequest(id)
		assert.False(t, exists)
		assert.Nil(t, r)
	}
}

func (suite *ServerStateTestSuite) TestClearClientPendingRequests() {
	t := suite.T()
	client1 := "client1"
	client2 := "client2"
	suite.state.AddPendingRequest(client1, "1234", newMockRequest("somevalue1"))
	suite.state.AddPendingRequest(client2, "5678", newMockRequest("somevalue2"))
	require.True(t, suite.state.HasPendingRequest(client1))
	suite.state.ClearClientPendingRequest(client1)
	// Request for client1 is deleted
	assert.False(t, suite.state.HasPendingRequest(client1))
	r, exists := suite.state.GetClientState(client1).GetPendingRequest("1234")
	assert.False(t, exists)
	assert.Nil(t, r)
	// Request for client2 is safe and sound
	assert.True(t, suite.state.HasPendingRequest(client2))
}

func (suite *ServerStateTestSuite) TestClearAllPendingRequests() {
	t := suite.T()
	client1 := "client1"
	client2 := "client2"
	suite.state.AddPendingRequest(client1, "1234", newMockRequest("somevalue1"))
	suite.state.AddPendingRequest(client2, "5678", newMockRequest("somevalue2"))
	require.True(t, suite.state.HasPendingRequests())
	suite.state.ClearAllPendingRequests()
	assert.False(t, suite.state.HasPendingRequests())
	// No more requests available in the struct
	assert.False(t, suite.state.HasPendingRequest(client1))
	assert.False(t, suite.state.HasPendingRequest(client2))
}

func (suite *ServerStateTestSuite) TestDeletePendingRequest() {
	t := suite.T()
	client1 := "client1"
	client2 := "client2"
	suite.state.AddPendingRequest(client1, "1234", newMockRequest("somevalue1"))
	suite.state.AddPendingRequest(client2, "5678", newMockRequest("somevalue2"))
	require.True(t, suite.state.HasPendingRequest(client1))
	require.True(t, suite.state.HasPendingRequest(client2))
	suite.state.DeletePendingRequest(client1, "1234")
	// Previously added request for client1 is gone
	assert.False(t, suite.state.HasPendingRequest(client1))
	r, exists := suite.state.GetClientState(client1).GetPendingRequest("1234")
	assert.False(t, exists)
	assert.Nil(t, r)
	// Deleting again has no effect
	suite.state.DeletePendingRequest(client1, "1234")
	assert.False(t, suite.state.HasPendingRequest(client1))
	// Previously added request for client2 is unaffected
	assert.True(t, suite.state.HasPendingRequest(client2))
}

func (suite *ServerStateTestSuite) TestDeleteInvalidPendingRequest() {
	t := suite.T()
	client1 := "client1"
	suite.state.AddPendingRequest(client1, "1234", newMockRequest("somevalue1"))
	require.True(t, suite.state.HasPendingRequest(client1))
	suite.state.DeletePendingRequest(client1, "5678")
	// Previously added request is still there
	assert.True(t, suite.state.HasPendingRequest(client1))
	r, exists := suite.state.GetClientState(client1).GetPendingRequest("1234")
	assert.True(t, exists)
	assert.NotNil(t, r)
}
