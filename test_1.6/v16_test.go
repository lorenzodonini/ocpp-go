package test_v16

import (
	"github.com/lorenzodonini/go-ocpp/ocpp"
	ocpp16 "github.com/lorenzodonini/go-ocpp/ocpp1.6"
	"github.com/lorenzodonini/go-ocpp/test"
	"github.com/stretchr/testify/suite"
	"testing"
)

// Tests
type OcppV16TestSuite struct {
	suite.Suite
	chargePoint   *ocpp.ChargePoint
	centralSystem *ocpp.CentralSystem
	mockServer    *test.MockWebsocketServer
	mockClient    *test.MockWebsocketClient
}

func (suite *OcppV16TestSuite) SetupTest() {
	coreProfile := ocpp.NewProfile("core", ocpp16.BootNotificationFeature{}, ocpp16.AuthorizeFeature{}, ocpp16.ChangeAvailabilityFeature{})
	mockClient := test.MockWebsocketClient{}
	mockServer := test.MockWebsocketServer{}
	suite.mockClient = &mockClient
	suite.mockServer = &mockServer
	suite.chargePoint = ocpp.NewChargePoint("test_id", suite.mockClient, coreProfile)
	suite.centralSystem = ocpp.NewCentralSystem(suite.mockServer, coreProfile)
}

//TODO: implement generic protocol tests

func TestOcpp16Protocol(t *testing.T) {
	suite.Run(t, new(OcppV16TestSuite))
}
