package firmware

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *firmwareTestSuite) TestPublishFirmwareRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{NewPublishFirmwareRequest("https://someurl", "deadbeef", 42), true},
		{PublishFirmwareRequest{Location: "http://someurl", Retries: tests.NewInt(5), Checksum: "deadbeef", RequestID: 42, RetryInterval: tests.NewInt(300)}, true},
		{PublishFirmwareRequest{Location: "http://someurl", Retries: tests.NewInt(5), Checksum: "deadbeef", RequestID: 42}, true},
		{PublishFirmwareRequest{Location: "http://someurl", Checksum: "deadbeef", RequestID: 42}, true},
		{PublishFirmwareRequest{Location: "http://someurl", Checksum: "deadbeef"}, true},
		{PublishFirmwareRequest{Location: "http://someurl"}, false},
		{PublishFirmwareRequest{Checksum: "deadbeef"}, false},
		{PublishFirmwareRequest{}, false},
		{PublishFirmwareRequest{Location: "http://someurl", Retries: tests.NewInt(5), Checksum: "deadbeef", RequestID: 42, RetryInterval: tests.NewInt(-1)}, false},
		{PublishFirmwareRequest{Location: "http://someurl", Retries: tests.NewInt(5), Checksum: "deadbeef", RequestID: -1, RetryInterval: tests.NewInt(300)}, false},
		{PublishFirmwareRequest{Location: "http://someurl", Retries: tests.NewInt(5), Checksum: ">32..............................", RequestID: 42, RetryInterval: tests.NewInt(300)}, false},
		{PublishFirmwareRequest{Location: "http://someurl", Retries: tests.NewInt(-1), Checksum: "deadbeef", RequestID: 42, RetryInterval: tests.NewInt(300)}, false},
		{PublishFirmwareRequest{Location: ">512.............................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................", Retries: tests.NewInt(5), Checksum: "deadbeef", RequestID: 42, RetryInterval: tests.NewInt(300)}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *firmwareTestSuite) TestPublishFirmwareResponseValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{PublishFirmwareResponse{Status: types.GenericStatusAccepted, StatusInfo: &types.StatusInfo{ReasonCode: "ok", AdditionalInfo: "someInfo"}}, true},
		{PublishFirmwareResponse{Status: types.GenericStatusAccepted}, true},
		{PublishFirmwareResponse{}, false},
		{PublishFirmwareResponse{Status: "invalidStatus"}, false},
		{PublishFirmwareResponse{Status: types.GenericStatusAccepted, StatusInfo: &types.StatusInfo{}}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *firmwareTestSuite) TestPublishFirmwareFeature() {
	feature := PublishFirmwareFeature{}
	suite.Equal(PublishFirmwareFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(PublishFirmwareRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(PublishFirmwareResponse{}), feature.GetResponseType())
}

func (suite *firmwareTestSuite) TestNewPublishFirmwareRequest() {
	req := NewPublishFirmwareRequest("https://example.com/fw", "abc123", 1)
	suite.NotNil(req)
	suite.Equal(PublishFirmwareFeatureName, req.GetFeatureName())
	suite.Equal("https://example.com/fw", req.Location)
	suite.Equal("abc123", req.Checksum)
	suite.Equal(1, req.RequestID)
}

func (suite *firmwareTestSuite) TestNewPublishFirmwareResponse() {
	resp := NewPublishFirmwareResponse(types.GenericStatusAccepted)
	suite.NotNil(resp)
	suite.Equal(PublishFirmwareFeatureName, resp.GetFeatureName())
	suite.Equal(types.GenericStatusAccepted, resp.Status)
}
