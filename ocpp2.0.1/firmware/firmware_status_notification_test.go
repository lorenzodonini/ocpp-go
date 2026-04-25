package firmware

import (
	"reflect"
	"testing"

	"github.com/lorenzodonini/ocpp-go/tests"
	"github.com/stretchr/testify/suite"
)

type firmwareTestSuite struct {
	suite.Suite
}

func (suite *firmwareTestSuite) TestFirmwareStatusNotificationRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{FirmwareStatusNotificationRequest{Status: FirmwareStatusDownloaded, RequestID: tests.NewInt(42)}, true},
		{FirmwareStatusNotificationRequest{Status: FirmwareStatusDownloaded}, true},
		{FirmwareStatusNotificationRequest{RequestID: tests.NewInt(42)}, false},
		{FirmwareStatusNotificationRequest{}, false},
		{FirmwareStatusNotificationRequest{Status: FirmwareStatusDownloaded, RequestID: tests.NewInt(-1)}, false},
		{FirmwareStatusNotificationRequest{Status: "invalidFirmwareStatus"}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *firmwareTestSuite) TestFirmwareStatusNotificationConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{FirmwareStatusNotificationResponse{}, true},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *firmwareTestSuite) TestFirmwareStatusNotificationFeature() {
	feature := FirmwareStatusNotificationFeature{}
	suite.Equal(FirmwareStatusNotificationFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(FirmwareStatusNotificationRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(FirmwareStatusNotificationResponse{}), feature.GetResponseType())
}

func (suite *firmwareTestSuite) TestNewFirmwareStatusNotificationRequest() {
	req := NewFirmwareStatusNotificationRequest(FirmwareStatusDownloaded)
	suite.NotNil(req)
	suite.Equal(FirmwareStatusNotificationFeatureName, req.GetFeatureName())
	suite.Equal(FirmwareStatusDownloaded, req.Status)
}

func (suite *firmwareTestSuite) TestNewFirmwareStatusNotificationResponse() {
	resp := NewFirmwareStatusNotificationResponse()
	suite.NotNil(resp)
	suite.Equal(FirmwareStatusNotificationFeatureName, resp.GetFeatureName())
}

func TestFirmwareSuite(t *testing.T) {
	suite.Run(t, new(firmwareTestSuite))
}