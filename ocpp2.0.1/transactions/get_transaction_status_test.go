package transactions

import (
	"reflect"
	"testing"

	"github.com/lorenzodonini/ocpp-go/tests"
	"github.com/stretchr/testify/suite"
)

type transactionsTestSuite struct {
	suite.Suite
}

func (suite *transactionsTestSuite) TestGetTransactionStatusRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{GetTransactionStatusRequest{}, true},
		{GetTransactionStatusRequest{TransactionID: "12345"}, true},
		{GetTransactionStatusRequest{TransactionID: ">36.................................."}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *transactionsTestSuite) TestGetTransactionStatusResponseValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{GetTransactionStatusResponse{OngoingIndicator: tests.NewBool(true), MessagesInQueue: true}, true},
		{GetTransactionStatusResponse{MessagesInQueue: true}, true},
		{GetTransactionStatusResponse{}, true},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *transactionsTestSuite) TestGetTransactionStatusFeature() {
	feature := GetTransactionStatusFeature{}
	suite.Equal(GetTransactionStatusFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(GetTransactionStatusRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(GetTransactionStatusResponse{}), feature.GetResponseType())
}

func (suite *transactionsTestSuite) TestNewGetTransactionStatusRequest() {
	req := NewGetTransactionStatusRequest()
	suite.NotNil(req)
	suite.Equal(GetTransactionStatusFeatureName, req.GetFeatureName())
}

func (suite *transactionsTestSuite) TestNewGetTransactionStatusResponse() {
	resp := NewGetTransactionStatusResponse(true)
	suite.NotNil(resp)
	suite.Equal(GetTransactionStatusFeatureName, resp.GetFeatureName())
	suite.True(resp.MessagesInQueue)
}

func TestTransactionsSuite(t *testing.T) {
	suite.Run(t, new(transactionsTestSuite))
}