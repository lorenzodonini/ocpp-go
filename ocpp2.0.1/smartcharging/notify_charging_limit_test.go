package smartcharging

import (
	"reflect"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *smartChargingTestSuite) TestNotifyChargingLimitRequestValidation() {
	t := suite.T()
	chargingSchedule := types.ChargingSchedule{
		ID:                     1,
		StartSchedule:          types.NewDateTime(time.Now()),
		Duration:               tests.NewInt(600),
		ChargingRateUnit:       types.ChargingRateUnitWatts,
		MinChargingRate:        tests.NewFloat(6.0),
		ChargingSchedulePeriod: []types.ChargingSchedulePeriod{types.NewChargingSchedulePeriod(0, 10.0)},
	}
	var requestTable = []tests.GenericTestEntry{
		{NotifyChargingLimitRequest{EvseID: tests.NewInt(1), ChargingLimit: ChargingLimit{ChargingLimitSource: types.ChargingLimitSourceEMS, IsGridCritical: tests.NewBool(false)}, ChargingSchedule: []types.ChargingSchedule{chargingSchedule}}, true},
		{NotifyChargingLimitRequest{EvseID: tests.NewInt(1), ChargingLimit: ChargingLimit{ChargingLimitSource: types.ChargingLimitSourceEMS, IsGridCritical: tests.NewBool(false)}, ChargingSchedule: []types.ChargingSchedule{}}, true},
		{NotifyChargingLimitRequest{EvseID: tests.NewInt(1), ChargingLimit: ChargingLimit{ChargingLimitSource: types.ChargingLimitSourceEMS, IsGridCritical: tests.NewBool(false)}}, true},
		{NotifyChargingLimitRequest{ChargingLimit: ChargingLimit{ChargingLimitSource: types.ChargingLimitSourceEMS, IsGridCritical: tests.NewBool(false)}}, true},
		{NotifyChargingLimitRequest{ChargingLimit: ChargingLimit{ChargingLimitSource: types.ChargingLimitSourceEMS}}, true},
		{NotifyChargingLimitRequest{ChargingLimit: ChargingLimit{}}, false},
		{NotifyChargingLimitRequest{}, false},
		{NotifyChargingLimitRequest{ChargingLimit: ChargingLimit{ChargingLimitSource: "invalidChargingLimitSource", IsGridCritical: tests.NewBool(false)}}, false},
		{NotifyChargingLimitRequest{EvseID: tests.NewInt(-1), ChargingLimit: ChargingLimit{ChargingLimitSource: types.ChargingLimitSourceEMS, IsGridCritical: tests.NewBool(false)}}, false},
		{NotifyChargingLimitRequest{ChargingLimit: ChargingLimit{ChargingLimitSource: types.ChargingLimitSourceEMS, IsGridCritical: tests.NewBool(false)}, ChargingSchedule: []types.ChargingSchedule{{ChargingRateUnit: "invalidStruct"}}}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *smartChargingTestSuite) TestNotifyChargingLimitResponseValidation() {
	t := suite.T()
	var responseTable = []tests.GenericTestEntry{
		{NotifyChargingLimitResponse{}, true},
	}
	tests.ExecuteGenericTestTable(t, responseTable)
}

func (suite *smartChargingTestSuite) TestNotifyChargingLimitFeature() {
	feature := NotifyChargingLimitFeature{}
	suite.Equal(NotifyChargingLimitFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(NotifyChargingLimitRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(NotifyChargingLimitResponse{}), feature.GetResponseType())
}

func (suite *smartChargingTestSuite) TestNewNotifyChargingLimitRequest() {
	limit := ChargingLimit{ChargingLimitSource: types.ChargingLimitSourceEMS}
	req := NewNotifyChargingLimitRequest(limit)
	suite.NotNil(req)
	suite.Equal(NotifyChargingLimitFeatureName, req.GetFeatureName())
	suite.Equal(limit, req.ChargingLimit)
}

func (suite *smartChargingTestSuite) TestNewNotifyChargingLimitResponse() {
	resp := NewNotifyChargingLimitResponse()
	suite.NotNil(resp)
	suite.Equal(NotifyChargingLimitFeatureName, resp.GetFeatureName())
}
