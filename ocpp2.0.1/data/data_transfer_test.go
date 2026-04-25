package data

import (
	"reflect"
	"testing"

	"github.com/lorenzodonini/ocpp-go/tests"
	"github.com/stretchr/testify/suite"
)

type dataTestSuite struct {
	suite.Suite
}

func (suite *dataTestSuite) TestDataTransferRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{DataTransferRequest{VendorID: "12345"}, true},
		{DataTransferRequest{VendorID: "12345", MessageID: "6789"}, true},
		{DataTransferRequest{VendorID: "12345", MessageID: "6789", Data: "mockData"}, true},
		{DataTransferRequest{}, false},
		{DataTransferRequest{VendorID: ">255............................................................................................................................................................................................................................................................"}, false},
		{DataTransferRequest{VendorID: "12345", MessageID: ">50................................................"}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *dataTestSuite) TestDataTransferConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{DataTransferResponse{Status: DataTransferStatusAccepted}, true},
		{DataTransferResponse{Status: DataTransferStatusRejected}, true},
		{DataTransferResponse{Status: DataTransferStatusUnknownMessageId}, true},
		{DataTransferResponse{Status: DataTransferStatusUnknownVendorId}, true},
		{DataTransferResponse{Status: "invalidDataTransferStatus"}, false},
		{DataTransferResponse{Status: DataTransferStatusAccepted, Data: "mockData"}, true},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *dataTestSuite) TestDataTransferFeature() {
	feature := DataTransferFeature{}
	suite.Equal(DataTransferFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(DataTransferRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(DataTransferResponse{}), feature.GetResponseType())
}

func (suite *dataTestSuite) TestNewDataTransferRequest() {
	req := NewDataTransferRequest("vendorX")
	suite.NotNil(req)
	suite.Equal(DataTransferFeatureName, req.GetFeatureName())
	suite.Equal("vendorX", req.VendorID)
}

func (suite *dataTestSuite) TestNewDataTransferResponse() {
	resp := NewDataTransferResponse(DataTransferStatusAccepted)
	suite.NotNil(resp)
	suite.Equal(DataTransferFeatureName, resp.GetFeatureName())
	suite.Equal(DataTransferStatusAccepted, resp.Status)
}

func TestDataSuite(t *testing.T) {
	suite.Run(t, new(dataTestSuite))
}
