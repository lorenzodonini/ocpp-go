package firmware

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *firmwareTestSuite) TestPublishFirmwareStatusNotificationRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{PublishFirmwareStatusNotificationRequest{Status: PublishFirmwareStatusPublished, Location: []string{"http://someUri"}, RequestID: tests.NewInt(42)}, true},
		{PublishFirmwareStatusNotificationRequest{Status: PublishFirmwareStatusPublished, Location: []string{"http://someUri"}}, true},
		{PublishFirmwareStatusNotificationRequest{Status: PublishFirmwareStatusChecksumVerified, Location: []string{}}, true},
		{PublishFirmwareStatusNotificationRequest{Status: PublishFirmwareStatusChecksumVerified}, true},
		{PublishFirmwareStatusNotificationRequest{Status: PublishFirmwareStatusDownloaded}, true},
		{PublishFirmwareStatusNotificationRequest{Status: PublishFirmwareStatusDownloadFailed}, true},
		{PublishFirmwareStatusNotificationRequest{Status: PublishFirmwareStatusDownloading}, true},
		{PublishFirmwareStatusNotificationRequest{Status: PublishFirmwareStatusDownloadScheduled}, true},
		{PublishFirmwareStatusNotificationRequest{Status: PublishFirmwareStatusDownloadPaused}, true},
		{PublishFirmwareStatusNotificationRequest{Status: PublishFirmwareStatusInvalidChecksum}, true},
		{PublishFirmwareStatusNotificationRequest{Status: PublishFirmwareStatusIdle}, true},
		{PublishFirmwareStatusNotificationRequest{Status: PublishFirmwareStatusPublishFailed}, true},
		{PublishFirmwareStatusNotificationRequest{}, false},
		{PublishFirmwareStatusNotificationRequest{Status: "invalidStatus"}, false},
		{PublishFirmwareStatusNotificationRequest{Status: PublishFirmwareStatusPublished, Location: []string{"http://someUri"}, RequestID: tests.NewInt(-1)}, false},
		{PublishFirmwareStatusNotificationRequest{Status: PublishFirmwareStatusPublished, Location: []string{"http://someUri>512..............................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................."}, RequestID: tests.NewInt(42)}, false},
		//TODO: add test for empty location field with published status
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *firmwareTestSuite) TestPublishFirmwareStatusNotificationResponseValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{PublishFirmwareStatusNotificationResponse{}, true},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *firmwareTestSuite) TestPublishFirmwareStatusNotificationFeature() {
	feature := PublishFirmwareStatusNotificationFeature{}
	suite.Equal(PublishFirmwareStatusNotificationFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(PublishFirmwareStatusNotificationRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(PublishFirmwareStatusNotificationResponse{}), feature.GetResponseType())
}

func (suite *firmwareTestSuite) TestNewPublishFirmwareStatusNotificationRequest() {
	req := NewPublishFirmwareStatusNotificationRequest(PublishFirmwareStatusPublished)
	suite.NotNil(req)
	suite.Equal(PublishFirmwareStatusNotificationFeatureName, req.GetFeatureName())
	suite.Equal(PublishFirmwareStatusPublished, req.Status)
}

func (suite *firmwareTestSuite) TestNewPublishFirmwareStatusNotificationResponse() {
	resp := NewPublishFirmwareStatusNotificationResponse()
	suite.NotNil(resp)
	suite.Equal(PublishFirmwareStatusNotificationFeatureName, resp.GetFeatureName())
}
