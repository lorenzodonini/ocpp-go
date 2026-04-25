package remotecontrol

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *remoteControlTestSuite) TestRequestStopTransactionRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{RequestStopTransactionRequest{TransactionID: "12345"}, true},
		{RequestStopTransactionRequest{}, false},
		{RequestStopTransactionRequest{TransactionID: ">36.................................."}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *remoteControlTestSuite) TestRequestStopTransactionConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{RequestStopTransactionResponse{Status: RequestStartStopStatusAccepted, StatusInfo: &types.StatusInfo{ReasonCode: "200"}}, true},
		{RequestStopTransactionResponse{Status: RequestStartStopStatusAccepted}, true},
		{RequestStopTransactionResponse{Status: RequestStartStopStatusRejected}, true},
		{RequestStopTransactionResponse{}, false},
		{RequestStopTransactionResponse{Status: "invalidRequestStartStopStatus", StatusInfo: &types.StatusInfo{ReasonCode: "200"}}, false},
		{RequestStopTransactionResponse{Status: RequestStartStopStatusAccepted, StatusInfo: &types.StatusInfo{}}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *remoteControlTestSuite) TestRequestStopTransactionFeature() {
	feature := RequestStopTransactionFeature{}
	suite.Equal(RequestStopTransactionFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(RequestStopTransactionRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(RequestStopTransactionResponse{}), feature.GetResponseType())
}

func (suite *remoteControlTestSuite) TestNewRequestStopTransactionRequest() {
	req := NewRequestStopTransactionRequest("txn-1")
	suite.NotNil(req)
	suite.Equal(RequestStopTransactionFeatureName, req.GetFeatureName())
	suite.Equal("txn-1", req.TransactionID)
}

func (suite *remoteControlTestSuite) TestNewRequestStopTransactionResponse() {
	resp := NewRequestStopTransactionResponse(RequestStartStopStatusAccepted)
	suite.NotNil(resp)
	suite.Equal(RequestStopTransactionFeatureName, resp.GetFeatureName())
	suite.Equal(RequestStartStopStatusAccepted, resp.Status)
}
