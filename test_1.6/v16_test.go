package test_v16

import (
	"github.com/lorenzodonini/go-ocpp/ocppj"
	ocpp16 "github.com/lorenzodonini/go-ocpp/ocpp1.6"
	"github.com/lorenzodonini/go-ocpp/test"
	"github.com/stretchr/testify/suite"
	"testing"
)

// Tests
type OcppV16TestSuite struct {
	suite.Suite
	chargePoint   *ocppj.ChargePoint
	centralSystem *ocppj.CentralSystem
	mockServer    *test.MockWebsocketServer
	mockClient    *test.MockWebsocketClient
}

func (suite *OcppV16TestSuite) SetupTest() {
	coreProfile := ocppj.NewProfile("core", ocpp16.BootNotificationFeature{}, ocpp16.AuthorizeFeature{}, ocpp16.ChangeAvailabilityFeature{})
	mockClient := test.MockWebsocketClient{}
	mockServer := test.MockWebsocketServer{}
	suite.mockClient = &mockClient
	suite.mockServer = &mockServer
	suite.chargePoint = ocppj.NewChargePoint("test_id", suite.mockClient, coreProfile)
	suite.centralSystem = ocppj.NewCentralSystem(suite.mockServer, coreProfile)
}

//TODO: implement generic protocol tests

func TestOcpp16Protocol(t *testing.T) {
	suite.Run(t, new(OcppV16TestSuite))
}
