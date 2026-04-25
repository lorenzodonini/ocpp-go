package firmware

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *firmwareTestSuite) TestUnpublishFirmwareRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{UnpublishFirmwareRequest{Checksum: "deadc0de"}, true},
		{UnpublishFirmwareRequest{}, false},
		{UnpublishFirmwareRequest{Checksum: ">32.............................."}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *firmwareTestSuite) TestUnpublishFirmwareResponseValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{UnpublishFirmwareResponse{Status: UnpublishFirmwareStatusUnpublished}, true},
		{UnpublishFirmwareResponse{Status: UnpublishFirmwareStatusNoFirmware}, true},
		{UnpublishFirmwareResponse{Status: UnpublishFirmwareStatusDownloadOngoing}, true},
		{UnpublishFirmwareResponse{}, false},
		{UnpublishFirmwareResponse{Status: "invalidUnpublishFirmwareStatus"}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *firmwareTestSuite) TestUnpublishFirmwareFeature() {
	feature := UnpublishFirmwareFeature{}
	suite.Equal(UnpublishFirmwareFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(UnpublishFirmwareRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(UnpublishFirmwareResponse{}), feature.GetResponseType())
}

func (suite *firmwareTestSuite) TestNewUnpublishFirmwareRequest() {
	req := NewUnpublishFirmwareRequest("deadc0de")
	suite.NotNil(req)
	suite.Equal(UnpublishFirmwareFeatureName, req.GetFeatureName())
	suite.Equal("deadc0de", req.Checksum)
}

func (suite *firmwareTestSuite) TestNewUnpublishFirmwareResponse() {
	resp := NewUnpublishFirmwareResponse(UnpublishFirmwareStatusUnpublished)
	suite.NotNil(resp)
	suite.Equal(UnpublishFirmwareFeatureName, resp.GetFeatureName())
	suite.Equal(UnpublishFirmwareStatusUnpublished, resp.Status)
}
