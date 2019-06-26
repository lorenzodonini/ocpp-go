package test

import (
	"errors"
	"github.com/lorenzodonini/go-ocpp/ocpp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (suite *OcppJTestSuite) TestCentralSystemSendRequest() {
	suite.mockServer.On("Write", mock.Anything).Return(nil)
	mockChargePointId := "1234"
	mockRequest := newMockRequest("mockValue")
	err := suite.centralSystem.SendRequest(mockChargePointId, mockRequest)
	assert.Nil(suite.T(), err)
}

func (suite *OcppJTestSuite) TestCentralSystemSendInvalidRequest() {
	suite.mockServer.On("Write", mock.Anything).Return(nil)
	mockChargePointId := "1234"
	mockRequest := newMockRequest("")
	err := suite.centralSystem.SendRequest(mockChargePointId, mockRequest)
	assert.NotNil(suite.T(), err)
}

func (suite *OcppJTestSuite) TestCentralSystemSendRequestPending() {
	suite.mockServer.On("Write", mock.Anything).Return(nil)
	mockChargePointId := "1234"
	mockRequest := newMockRequest("mockValue")
	err := suite.centralSystem.SendRequest(mockChargePointId, mockRequest)
	assert.Nil(suite.T(), err)
	suite.centralSystem.PendingRequests = map[string]ocpp.Request{} // Clearing map
	err = suite.centralSystem.SendRequest(mockChargePointId, mockRequest)
	assert.NotNil(suite.T(), err)
}

func (suite *OcppJTestSuite) TestCentralSystemSendRequestFailed() {
	suite.mockServer.On("Write", mock.Anything).Return(errors.New("mockError"))
	mockChargePointId := "1234"
	mockRequest := newMockRequest("mockValue")
	err := suite.centralSystem.SendRequest(mockChargePointId, mockRequest)
	assert.NotNil(suite.T(), err)
}
