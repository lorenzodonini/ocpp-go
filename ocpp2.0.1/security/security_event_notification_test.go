package security

import (
	"reflect"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *securityTestSuite) TestSecurityEventNotificationRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{SecurityEventNotificationRequest{Type: "type1", Timestamp: types.NewDateTime(time.Now()), TechInfo: "someTechInfo"}, true},
		{SecurityEventNotificationRequest{Type: "type1", Timestamp: types.NewDateTime(time.Now())}, true},
		{SecurityEventNotificationRequest{Type: "type1"}, false},
		{SecurityEventNotificationRequest{}, false},
		{SecurityEventNotificationRequest{Type: "", Timestamp: types.NewDateTime(time.Now()), TechInfo: "someTechInfo"}, false},
		{SecurityEventNotificationRequest{Type: ">50................................................", Timestamp: types.NewDateTime(time.Now()), TechInfo: "someTechInfo"}, false},
		{SecurityEventNotificationRequest{Type: "type1", Timestamp: types.NewDateTime(time.Now()), TechInfo: ">255............................................................................................................................................................................................................................................................"}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *securityTestSuite) TestSecurityEventNotificationConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{SecurityEventNotificationResponse{}, true},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *securityTestSuite) TestSecurityEventNotificationFeature() {
	feature := SecurityEventNotificationFeature{}
	suite.Equal(SecurityEventNotificationFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(SecurityEventNotificationRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(SecurityEventNotificationResponse{}), feature.GetResponseType())
}

func (suite *securityTestSuite) TestNewSecurityEventNotificationRequest() {
	ts := types.NewDateTime(time.Now())
	req := NewSecurityEventNotificationRequest("FirmwareUpdated", ts)
	suite.NotNil(req)
	suite.Equal(SecurityEventNotificationFeatureName, req.GetFeatureName())
	suite.Equal("FirmwareUpdated", req.Type)
	suite.Equal(ts, req.Timestamp)
}

func (suite *securityTestSuite) TestNewSecurityEventNotificationResponse() {
	resp := NewSecurityEventNotificationResponse()
	suite.NotNil(resp)
	suite.Equal(SecurityEventNotificationFeatureName, resp.GetFeatureName())
}
