package diagnostics

import (
	"reflect"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *diagnosticsTestSuite) TestGetLogRequestValidation() {
	t := suite.T()
	logParameters := LogParameters{
		RemoteLocation:  "ftp://someurl/diagnostics/1",
		OldestTimestamp: types.NewDateTime(time.Now().Add(-2 * time.Hour)),
		LatestTimestamp: types.NewDateTime(time.Now()),
	}
	var requestTable = []tests.GenericTestEntry{
		{GetLogRequest{LogType: LogTypeDiagnostics, RequestID: 1, Retries: tests.NewInt(5), RetryInterval: tests.NewInt(120), Log: logParameters}, true},
		{GetLogRequest{LogType: LogTypeDiagnostics, RequestID: 1, Retries: tests.NewInt(5), Log: logParameters}, true},
		{GetLogRequest{LogType: LogTypeDiagnostics, RequestID: 1, Log: logParameters}, true},
		{GetLogRequest{LogType: LogTypeDiagnostics, Log: logParameters}, true},
		{GetLogRequest{LogType: LogTypeDiagnostics}, false},
		{GetLogRequest{Log: logParameters}, false},
		{GetLogRequest{}, false},
		{GetLogRequest{LogType: "invalidLogType", RequestID: 1, Retries: tests.NewInt(5), RetryInterval: tests.NewInt(120), Log: logParameters}, false},
		{GetLogRequest{LogType: LogTypeDiagnostics, RequestID: -1, Retries: tests.NewInt(5), RetryInterval: tests.NewInt(120), Log: logParameters}, false},
		{GetLogRequest{LogType: LogTypeDiagnostics, RequestID: 1, Retries: tests.NewInt(-1), RetryInterval: tests.NewInt(120), Log: logParameters}, false},
		{GetLogRequest{LogType: LogTypeDiagnostics, RequestID: 1, Retries: tests.NewInt(5), RetryInterval: tests.NewInt(-1), Log: logParameters}, false},
		{GetLogRequest{LogType: LogTypeDiagnostics, RequestID: 1, Retries: tests.NewInt(5), RetryInterval: tests.NewInt(120), Log: LogParameters{RemoteLocation: ".invalidUrl.", OldestTimestamp: nil, LatestTimestamp: nil}}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *diagnosticsTestSuite) TestGetLogConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{GetLogResponse{Status: LogStatusAccepted, Filename: "testFileName.log"}, true},
		{GetLogResponse{Status: LogStatusAccepted}, true},
		{GetLogResponse{Status: LogStatusRejected}, true},
		{GetLogResponse{Status: LogStatusAcceptedCanceled}, true},
		{GetLogResponse{}, false},
		{GetLogResponse{Status: "invalidLogStatus"}, false},
		{GetLogResponse{Status: LogStatusAccepted, Filename: ">256............................................................................................................................................................................................................................................................."}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *diagnosticsTestSuite) TestGetLogFeature() {
	feature := GetLogFeature{}
	suite.Equal(GetLogFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(GetLogRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(GetLogResponse{}), feature.GetResponseType())
}

func (suite *diagnosticsTestSuite) TestNewGetLogRequest() {
	logParams := LogParameters{RemoteLocation: "ftp://server/logs"}
	req := NewGetLogRequest(LogTypeDiagnostics, 1, logParams)
	suite.NotNil(req)
	suite.Equal(GetLogFeatureName, req.GetFeatureName())
	suite.Equal(LogTypeDiagnostics, req.LogType)
	suite.Equal(1, req.RequestID)
	suite.Equal(logParams, req.Log)
}

func (suite *diagnosticsTestSuite) TestNewGetLogResponse() {
	resp := NewGetLogResponse(LogStatusAccepted)
	suite.NotNil(resp)
	suite.Equal(GetLogFeatureName, resp.GetFeatureName())
	suite.Equal(LogStatusAccepted, resp.Status)
}
