package test

import (
	"github.com/lorenzodonini/go-ocpp/ocpp"
	"github.com/lorenzodonini/go-ocpp/ocpp/1.6"
	"github.com/stretchr/testify/suite"
	"gopkg.in/go-playground/validator.v9"
	"testing"
)

// Tests
type OcppV16TestSuite struct {
	suite.Suite
	chargePoint   *ocpp.ChargePoint
	centralSystem *ocpp.CentralSystem
	mockServer    *MockWebsocketServer
	mockClient    *MockWebsocketClient
}

func (suite *OcppV16TestSuite) SetupTest() {
	coreProfile := ocpp.NewProfile("core", v16.BootNotificationFeature{}, v16.AuthorizeFeature{}, v16.ChangeAvailabilityFeature{})
	mockClient := MockWebsocketClient{}
	mockServer := MockWebsocketServer{}
	suite.mockClient = &mockClient
	suite.mockServer = &mockServer
	suite.chargePoint = ocpp.NewChargePoint("test_id", suite.mockClient, coreProfile)
	suite.centralSystem = ocpp.NewCentralSystem(suite.mockServer, coreProfile)
}

var validate = validator.New()

//TODO: implement generic protocol tests

func TestOcpp16Protocol(t *testing.T) {
	suite.Run(t, new(OcppV16TestSuite))
}
