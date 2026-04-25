package diagnostics

import (
	"reflect"
	"testing"

	"github.com/lorenzodonini/ocpp-go/tests"
	"github.com/stretchr/testify/suite"
)

type diagnosticsTestSuite struct {
	suite.Suite
}

func (suite *diagnosticsTestSuite) TestClearVariableMonitoringRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{ClearVariableMonitoringRequest{ID: []int{0, 2, 15}}, true},
		{ClearVariableMonitoringRequest{ID: []int{0}}, true},
		{ClearVariableMonitoringRequest{ID: []int{}}, false},
		{ClearVariableMonitoringRequest{}, false},
		{ClearVariableMonitoringRequest{ID: []int{-1}}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *diagnosticsTestSuite) TestClearVariableMonitoringConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{ClearVariableMonitoringResponse{ClearMonitoringResult: []ClearMonitoringResult{{ID: 2, Status: ClearMonitoringStatusAccepted}}}, true},
		{ClearVariableMonitoringResponse{ClearMonitoringResult: []ClearMonitoringResult{{ID: 2}}}, false},
		{ClearVariableMonitoringResponse{ClearMonitoringResult: []ClearMonitoringResult{}}, false},
		{ClearVariableMonitoringResponse{}, false},
		{ClearVariableMonitoringResponse{ClearMonitoringResult: []ClearMonitoringResult{{ID: -1, Status: ClearMonitoringStatusAccepted}}}, false},
		{ClearVariableMonitoringResponse{ClearMonitoringResult: []ClearMonitoringResult{{ID: 2, Status: "invalidClearMonitoringStatus"}}}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *diagnosticsTestSuite) TestClearVariableMonitoringFeature() {
	feature := ClearVariableMonitoringFeature{}
	suite.Equal(ClearVariableMonitoringFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(ClearVariableMonitoringRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(ClearVariableMonitoringResponse{}), feature.GetResponseType())
}

func (suite *diagnosticsTestSuite) TestNewClearVariableMonitoringRequest() {
	ids := []int{1, 2, 3}
	req := NewClearVariableMonitoringRequest(ids)
	suite.NotNil(req)
	suite.Equal(ClearVariableMonitoringFeatureName, req.GetFeatureName())
	suite.Equal(ids, req.ID)
}

func (suite *diagnosticsTestSuite) TestNewClearVariableMonitoringResponse() {
	result := []ClearMonitoringResult{{ID: 1, Status: ClearMonitoringStatusAccepted}}
	resp := NewClearVariableMonitoringResponse(result)
	suite.NotNil(resp)
	suite.Equal(ClearVariableMonitoringFeatureName, resp.GetFeatureName())
	suite.Equal(result, resp.ClearMonitoringResult)
}

func TestDiagnosticsSuite(t *testing.T) {
	suite.Run(t, new(diagnosticsTestSuite))
}