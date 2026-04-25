package smartcharging

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *smartChargingTestSuite) TestClearedChargingLimitRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{ClearedChargingLimitRequest{ChargingLimitSource: types.ChargingLimitSourceEMS, EvseID: tests.NewInt(0)}, true},
		{ClearedChargingLimitRequest{ChargingLimitSource: types.ChargingLimitSourceEMS}, true},
		{ClearedChargingLimitRequest{}, false},
		{ClearedChargingLimitRequest{ChargingLimitSource: types.ChargingLimitSourceEMS, EvseID: tests.NewInt(-1)}, false},
		{ClearedChargingLimitRequest{ChargingLimitSource: "invalidChargingLimitSource"}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *smartChargingTestSuite) TestClearedChargingLimitConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{ClearedChargingLimitResponse{}, true},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *smartChargingTestSuite) TestClearedChargingLimitFeature() {
	feature := ClearedChargingLimitFeature{}
	suite.Equal(ClearedChargingLimitFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(ClearedChargingLimitRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(ClearedChargingLimitResponse{}), feature.GetResponseType())
}

func (suite *smartChargingTestSuite) TestNewClearedChargingLimitRequest() {
	req := NewClearedChargingLimitRequest(types.ChargingLimitSourceEMS)
	suite.NotNil(req)
	suite.Equal(ClearedChargingLimitFeatureName, req.GetFeatureName())
	suite.Equal(types.ChargingLimitSourceEMS, req.ChargingLimitSource)
}

func (suite *smartChargingTestSuite) TestNewClearedChargingLimitResponse() {
	resp := NewClearedChargingLimitResponse()
	suite.NotNil(resp)
	suite.Equal(ClearedChargingLimitFeatureName, resp.GetFeatureName())
}
