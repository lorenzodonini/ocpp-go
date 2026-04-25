package diagnostics

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *diagnosticsTestSuite) TestLogStatusNotificationRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{LogStatusNotificationRequest{Status: UploadLogStatusUploading, RequestID: 42}, true},
		{LogStatusNotificationRequest{Status: UploadLogStatusUploadFailure, RequestID: 42}, true},
		{LogStatusNotificationRequest{Status: UploadLogStatusUploaded, RequestID: 42}, true},
		{LogStatusNotificationRequest{Status: UploadLogStatusPermissionDenied, RequestID: 42}, true},
		{LogStatusNotificationRequest{Status: UploadLogStatusNotSupportedOp, RequestID: 42}, true},
		{LogStatusNotificationRequest{Status: UploadLogStatusIdle, RequestID: 42}, true},
		{LogStatusNotificationRequest{Status: UploadLogStatusBadMessage, RequestID: 42}, true},
		{LogStatusNotificationRequest{Status: UploadLogStatusIdle}, true},
		{LogStatusNotificationRequest{RequestID: 42}, false},
		{LogStatusNotificationRequest{}, false},
		{LogStatusNotificationRequest{Status: UploadLogStatusIdle, RequestID: -1}, false},
		{LogStatusNotificationRequest{Status: "invalidUploadLogStatus", RequestID: 42}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *diagnosticsTestSuite) TestLogStatusNotificationConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{LogStatusNotificationResponse{}, true},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *diagnosticsTestSuite) TestLogStatusNotificationFeature() {
	feature := LogStatusNotificationFeature{}
	suite.Equal(LogStatusNotificationFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(LogStatusNotificationRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(LogStatusNotificationResponse{}), feature.GetResponseType())
}

func (suite *diagnosticsTestSuite) TestNewLogStatusNotificationRequest() {
	req := NewLogStatusNotificationRequest(UploadLogStatusUploading, 42)
	suite.NotNil(req)
	suite.Equal(LogStatusNotificationFeatureName, req.GetFeatureName())
	suite.Equal(UploadLogStatusUploading, req.Status)
	suite.Equal(42, req.RequestID)
}

func (suite *diagnosticsTestSuite) TestNewLogStatusNotificationResponse() {
	resp := NewLogStatusNotificationResponse()
	suite.NotNil(resp)
	suite.Equal(LogStatusNotificationFeatureName, resp.GetFeatureName())
}
