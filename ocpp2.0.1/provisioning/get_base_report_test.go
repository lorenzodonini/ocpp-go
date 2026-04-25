package provisioning

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *provisioningTestSuite) TestGetBaseReportRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{GetBaseReportRequest{RequestID: 42, ReportBase: ReportTypeConfigurationInventory}, true},
		{GetBaseReportRequest{ReportBase: ReportTypeConfigurationInventory}, true},
		{GetBaseReportRequest{RequestID: 42}, false},
		{GetBaseReportRequest{}, false},
		{GetBaseReportRequest{RequestID: 42, ReportBase: "invalidReportType"}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *provisioningTestSuite) TestGetBaseReportConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{GetBaseReportResponse{Status: types.GenericDeviceModelStatusAccepted}, true},
		{GetBaseReportResponse{Status: "invalidDeviceModelStatus"}, false},
		{GetBaseReportResponse{}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *provisioningTestSuite) TestGetBaseReportFeature() {
	feature := GetBaseReportFeature{}
	suite.Equal(GetBaseReportFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(GetBaseReportRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(GetBaseReportResponse{}), feature.GetResponseType())
}

func (suite *provisioningTestSuite) TestNewGetBaseReportRequest() {
	req := NewGetBaseReportRequest(1, ReportTypeConfigurationInventory)
	suite.NotNil(req)
	suite.Equal(GetBaseReportFeatureName, req.GetFeatureName())
	suite.Equal(1, req.RequestID)
	suite.Equal(ReportTypeConfigurationInventory, req.ReportBase)
}

func (suite *provisioningTestSuite) TestNewGetBaseReportResponse() {
	resp := NewGetBaseReportResponse(types.GenericDeviceModelStatusAccepted)
	suite.NotNil(resp)
	suite.Equal(GetBaseReportFeatureName, resp.GetFeatureName())
	suite.Equal(types.GenericDeviceModelStatusAccepted, resp.Status)
}
