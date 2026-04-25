package availability

import (
	"reflect"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *availabilityTestSuite) TestHeartbeatRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{HeartbeatRequest{}, true},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *availabilityTestSuite) TestHeartbeatResponseValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{HeartbeatResponse{CurrentTime: *types.NewDateTime(time.Now())}, true},
		{HeartbeatResponse{}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *availabilityTestSuite) TestHeartbeatFeature() {
	feature := HeartbeatFeature{}
	suite.Equal(HeartbeatFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(HeartbeatRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(HeartbeatResponse{}), feature.GetResponseType())
}

func (suite *availabilityTestSuite) TestNewHeartbeatRequest() {
	req := NewHeartbeatRequest()
	suite.NotNil(req)
	suite.Equal(HeartbeatFeatureName, req.GetFeatureName())
}

func (suite *availabilityTestSuite) TestNewHeartbeatResponse() {
	now := *types.NewDateTime(time.Now())
	resp := NewHeartbeatResponse(now)
	suite.NotNil(resp)
	suite.Equal(HeartbeatFeatureName, resp.GetFeatureName())
	suite.Equal(now, resp.CurrentTime)
}
