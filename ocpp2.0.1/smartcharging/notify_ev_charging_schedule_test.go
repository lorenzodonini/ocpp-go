package smartcharging

import (
	"reflect"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *smartChargingTestSuite) TestNotifyEVChargingScheduleRequestValidation() {
	t := suite.T()
	chargingSchedule := types.ChargingSchedule{
		StartSchedule:          types.NewDateTime(time.Now()),
		Duration:               tests.NewInt(600),
		ChargingRateUnit:       types.ChargingRateUnitWatts,
		MinChargingRate:        tests.NewFloat(6.0),
		ChargingSchedulePeriod: []types.ChargingSchedulePeriod{types.NewChargingSchedulePeriod(0, 10.0)},
	}
	var requestTable = []tests.GenericTestEntry{
		// {ChargingRateUnit: "invalidStruct"}
		{NotifyEVChargingScheduleRequest{TimeBase: types.NewDateTime(time.Now()), EvseID: 1, ChargingSchedule: chargingSchedule}, true},
		{NotifyEVChargingScheduleRequest{TimeBase: types.NewDateTime(time.Now()), EvseID: 1}, false},
		{NotifyEVChargingScheduleRequest{TimeBase: types.NewDateTime(time.Now()), ChargingSchedule: chargingSchedule}, false},
		{NotifyEVChargingScheduleRequest{EvseID: 1}, false},
		{NotifyEVChargingScheduleRequest{}, false},
		{NotifyEVChargingScheduleRequest{TimeBase: types.NewDateTime(time.Now()), EvseID: -1, ChargingSchedule: chargingSchedule}, false},
		{NotifyEVChargingScheduleRequest{TimeBase: types.NewDateTime(time.Now()), EvseID: -1, ChargingSchedule: types.ChargingSchedule{ChargingRateUnit: "invalidStruct"}}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *smartChargingTestSuite) TestNotifyEVChargingScheduleResponseValidation() {
	t := suite.T()
	var responseTable = []tests.GenericTestEntry{
		{NotifyEVChargingScheduleResponse{Status: types.GenericStatusAccepted, StatusInfo: types.NewStatusInfo("ok", "someInfo")}, true},
		{NotifyEVChargingScheduleResponse{Status: types.GenericStatusRejected, StatusInfo: types.NewStatusInfo("ok", "someInfo")}, true},
		{NotifyEVChargingScheduleResponse{Status: types.GenericStatusAccepted}, true},
		{NotifyEVChargingScheduleResponse{}, false},
		{NotifyEVChargingScheduleResponse{Status: "invalidStatus"}, false},
		{NotifyEVChargingScheduleResponse{Status: types.GenericStatusAccepted, StatusInfo: types.NewStatusInfo("", "invalidStatusInfo")}, false},
	}
	tests.ExecuteGenericTestTable(t, responseTable)
}

func (suite *smartChargingTestSuite) TestNotifyEVChargingScheduleFeature() {
	feature := NotifyEVChargingScheduleFeature{}
	suite.Equal(NotifyEVChargingScheduleFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(NotifyEVChargingScheduleRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(NotifyEVChargingScheduleResponse{}), feature.GetResponseType())
}

func (suite *smartChargingTestSuite) TestNewNotifyEVChargingScheduleRequest() {
	ts := types.NewDateTime(time.Now())
	schedule := types.NewChargingSchedule(1, types.ChargingRateUnitWatts, types.NewChargingSchedulePeriod(0, 10.0))
	req := NewNotifyEVChargingScheduleRequest(ts, 1, *schedule)
	suite.NotNil(req)
	suite.Equal(NotifyEVChargingScheduleFeatureName, req.GetFeatureName())
	suite.Equal(ts, req.TimeBase)
	suite.Equal(1, req.EvseID)
}

func (suite *smartChargingTestSuite) TestNewNotifyEVChargingScheduleResponse() {
	resp := NewNotifyEVChargingScheduleResponse(types.GenericStatusAccepted)
	suite.NotNil(resp)
	suite.Equal(NotifyEVChargingScheduleFeatureName, resp.GetFeatureName())
	suite.Equal(types.GenericStatusAccepted, resp.Status)
}
