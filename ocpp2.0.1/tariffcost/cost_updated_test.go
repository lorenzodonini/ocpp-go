package tariffcost

import (
	"reflect"
	"testing"

	"github.com/lorenzodonini/ocpp-go/tests"
	"github.com/stretchr/testify/suite"
)

type tariffCostTestSuite struct {
	suite.Suite
}

func (suite *tariffCostTestSuite) TestCostUpdatedRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{CostUpdatedRequest{TotalCost: 24.6, TransactionID: "1234"}, true},
		{CostUpdatedRequest{TotalCost: 24.6}, false},
		{CostUpdatedRequest{TransactionID: "1234"}, false},
		{CostUpdatedRequest{}, false},
		{CostUpdatedRequest{TotalCost: 24.6, TransactionID: ">36.................................."}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *tariffCostTestSuite) TestCostUpdatedConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{CostUpdatedResponse{}, true},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *tariffCostTestSuite) TestCostUpdatedFeature() {
	feature := CostUpdatedFeature{}
	suite.Equal(CostUpdatedFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(CostUpdatedRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(CostUpdatedResponse{}), feature.GetResponseType())
}

func (suite *tariffCostTestSuite) TestNewCostUpdatedRequest() {
	req := NewCostUpdatedRequest(9.99, "txn-42")
	suite.NotNil(req)
	suite.Equal(CostUpdatedFeatureName, req.GetFeatureName())
	suite.InDelta(9.99, req.TotalCost, 0.001)
	suite.Equal("txn-42", req.TransactionID)
}

func (suite *tariffCostTestSuite) TestNewCostUpdatedResponse() {
	resp := NewCostUpdatedResponse()
	suite.NotNil(resp)
	suite.Equal(CostUpdatedFeatureName, resp.GetFeatureName())
}

func TestTariffCostSuite(t *testing.T) {
	suite.Run(t, new(tariffCostTestSuite))
}