package ocppj_test

import (
	"sync"

	"github.com/stretchr/testify/suite"

	"github.com/lorenzodonini/ocpp-go/ocppj"
)

type ClientStateTestSuite struct {
	suite.Suite
	state ocppj.ClientState
}

func (suite *ClientStateTestSuite) SetupTest() {
	suite.state = ocppj.NewClientState()
}

func (suite *ClientStateTestSuite) TestAddPendingRequest() {
	requestID := "1234"
	req := newMockRequest("somevalue")
	suite.Require().False(suite.state.HasPendingRequest())
	suite.state.AddPendingRequest(requestID, req)
	suite.Require().True(suite.state.HasPendingRequest())
	r, exists := suite.state.GetPendingRequest(requestID)
	suite.Assert().True(exists)
	suite.Assert().Equal(req, r)
}

func (suite *ClientStateTestSuite) TestGetInvalidPendingRequest() {
	requestID := "1234"
	suite.state.AddPendingRequest(requestID, newMockRequest("somevalue"))
	suite.Require().True(suite.state.HasPendingRequest())
	invalidRequestIDs := []string{"4321", "5678", "1230", "deadc0de"}
	// Nothing returned when querying for an unknown request ID
	for _, id := range invalidRequestIDs {
		r, exists := suite.state.GetPendingRequest(id)
		suite.Assert().False(exists)
		suite.Assert().Nil(r)
	}
}

func (suite *ClientStateTestSuite) TestAddMultiplePendingRequests() {
	requestId1 := "1234"
	requestId2 := "5678"
	req1 := newMockRequest("somevalue1")
	req2 := newMockRequest("somevalue2")
	suite.state.AddPendingRequest(requestId1, req1)
	suite.state.AddPendingRequest(requestId2, req2)
	r, exists := suite.state.GetPendingRequest(requestId1)
	suite.Assert().True(exists)
	suite.Assert().NotNil(r)
	r, exists = suite.state.GetPendingRequest(requestId2)
	suite.Assert().False(exists)
	suite.Assert().Nil(r)
}

func (suite *ClientStateTestSuite) TestDeletePendingRequest() {
	requestID := "1234"
	req := newMockRequest("somevalue")
	suite.state.AddPendingRequest(requestID, req)
	suite.Require().True(suite.state.HasPendingRequest())
	suite.state.DeletePendingRequest(requestID)
	// Previously added request is gone
	suite.Assert().False(suite.state.HasPendingRequest())
	r, exists := suite.state.GetPendingRequest(requestID)
	suite.Assert().False(exists)
	suite.Assert().Nil(r)
	// Deleting again has no effect
	suite.state.DeletePendingRequest(requestID)
	suite.Assert().False(suite.state.HasPendingRequest())
}

func (suite *ClientStateTestSuite) TestDeleteInvalidPendingRequest() {
	requestID := "1234"
	req := newMockRequest("somevalue")
	suite.state.AddPendingRequest(requestID, req)
	suite.Require().True(suite.state.HasPendingRequest())
	suite.state.DeletePendingRequest("5678")
	// Previously added request is still there
	suite.Assert().True(suite.state.HasPendingRequest())
	r, exists := suite.state.GetPendingRequest(requestID)
	suite.Assert().True(exists)
	suite.Assert().NotNil(r)
}

func (suite *ClientStateTestSuite) TestClearPendingRequests() {
	requestID := "1234"
	req := newMockRequest("somevalue")
	suite.state.AddPendingRequest(requestID, req)
	suite.Require().True(suite.state.HasPendingRequest())
	suite.state.ClearPendingRequests()
	// No more requests available in the struct
	suite.Assert().False(suite.state.HasPendingRequest())
}

type ServerStateTestSuite struct {
	suite.Suite
	mutex sync.RWMutex
	state ocppj.ServerState
}

func (suite *ServerStateTestSuite) SetupTest() {
	suite.state = ocppj.NewServerState(&suite.mutex)
}

func (suite *ServerStateTestSuite) TestAddPendingRequests() {
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
	suite.Require().True(suite.state.HasPendingRequests())
	for _, r := range requests {
		suite.Assert().True(suite.state.HasPendingRequest(r.clientID))
		req, exists := suite.state.GetClientState(r.clientID).GetPendingRequest(r.requestID)
		suite.Assert().True(exists)
		suite.Assert().Equal(r.request, req)
	}
}

func (suite *ServerStateTestSuite) TestGetInvalidPendingRequest() {
	requestID := "1234"
	clientID := "client1"
	suite.state.AddPendingRequest(clientID, requestID, newMockRequest("somevalue"))
	suite.Require().True(suite.state.HasPendingRequest(clientID))
	invalidRequestIDs := []string{"4321", "5678", "1230", "deadc0de"}
	// Nothing returned when querying for an unknown request ID
	for _, id := range invalidRequestIDs {
		r, exists := suite.state.GetClientState(clientID).GetPendingRequest(id)
		suite.Assert().False(exists)
		suite.Assert().Nil(r)
	}
}

func (suite *ServerStateTestSuite) TestClearClientPendingRequests() {
	client1 := "client1"
	client2 := "client2"
	suite.state.AddPendingRequest(client1, "1234", newMockRequest("somevalue1"))
	suite.state.AddPendingRequest(client2, "5678", newMockRequest("somevalue2"))
	suite.Require().True(suite.state.HasPendingRequest(client1))
	suite.state.ClearClientPendingRequest(client1)
	// Request for client1 is deleted
	suite.Assert().False(suite.state.HasPendingRequest(client1))
	r, exists := suite.state.GetClientState(client1).GetPendingRequest("1234")
	suite.Assert().False(exists)
	suite.Assert().Nil(r)
	// Request for client2 is safe and sound
	suite.Assert().True(suite.state.HasPendingRequest(client2))
}

func (suite *ServerStateTestSuite) TestClearAllPendingRequests() {
	client1 := "client1"
	client2 := "client2"
	suite.state.AddPendingRequest(client1, "1234", newMockRequest("somevalue1"))
	suite.state.AddPendingRequest(client2, "5678", newMockRequest("somevalue2"))
	suite.Require().True(suite.state.HasPendingRequests())
	suite.state.ClearAllPendingRequests()
	suite.Assert().False(suite.state.HasPendingRequests())
	// No more requests available in the struct
	suite.Assert().False(suite.state.HasPendingRequest(client1))
	suite.Assert().False(suite.state.HasPendingRequest(client2))
}

func (suite *ServerStateTestSuite) TestDeletePendingRequest() {
	client1 := "client1"
	client2 := "client2"
	suite.state.AddPendingRequest(client1, "1234", newMockRequest("somevalue1"))
	suite.state.AddPendingRequest(client2, "5678", newMockRequest("somevalue2"))
	suite.Require().True(suite.state.HasPendingRequest(client1))
	suite.Require().True(suite.state.HasPendingRequest(client2))
	suite.state.DeletePendingRequest(client1, "1234")
	// Previously added request for client1 is gone
	suite.Assert().False(suite.state.HasPendingRequest(client1))
	r, exists := suite.state.GetClientState(client1).GetPendingRequest("1234")
	suite.Assert().False(exists)
	suite.Assert().Nil(r)
	// Deleting again has no effect
	suite.state.DeletePendingRequest(client1, "1234")
	suite.Assert().False(suite.state.HasPendingRequest(client1))
	// Previously added request for client2 is unaffected
	suite.Assert().True(suite.state.HasPendingRequest(client2))
}

func (suite *ServerStateTestSuite) TestDeleteInvalidPendingRequest() {
	client1 := "client1"
	suite.state.AddPendingRequest(client1, "1234", newMockRequest("somevalue1"))
	suite.Require().True(suite.state.HasPendingRequest(client1))
	suite.state.DeletePendingRequest(client1, "5678")
	// Previously added request is still there
	suite.Assert().True(suite.state.HasPendingRequest(client1))
	r, exists := suite.state.GetClientState(client1).GetPendingRequest("1234")
	suite.Assert().True(exists)
	suite.Assert().NotNil(r)
}
