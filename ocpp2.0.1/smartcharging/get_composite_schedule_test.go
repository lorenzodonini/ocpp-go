package smartcharging

import (
	"reflect"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *smartChargingTestSuite) TestGetCompositeScheduleRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{GetCompositeScheduleRequest{Duration: 600, EvseID: 1, ChargingRateUnit: types.ChargingRateUnitWatts}, true},
		{GetCompositeScheduleRequest{Duration: 600, EvseID: 1}, true},
		{GetCompositeScheduleRequest{EvseID: 1}, true},
		{GetCompositeScheduleRequest{}, true},
		{GetCompositeScheduleRequest{Duration: 600, EvseID: -1, ChargingRateUnit: types.ChargingRateUnitWatts}, false},
		{GetCompositeScheduleRequest{Duration: -1, EvseID: 1, ChargingRateUnit: types.ChargingRateUnitWatts}, false},
		{GetCompositeScheduleRequest{Duration: 600, EvseID: 1, ChargingRateUnit: "invalidChargingRateUnit"}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *smartChargingTestSuite) TestGetCompositeScheduleConfirmationValidation() {
	t := suite.T()
	chargingSchedule := types.NewChargingSchedule(1, types.ChargingRateUnitWatts, types.NewChargingSchedulePeriod(0, 10.0))
	chargingSchedule.Duration = tests.NewInt(600)
	chargingSchedule.MinChargingRate = tests.NewFloat(6.0)
	chargingSchedule.StartSchedule = types.NewDateTime(time.Now())
	compositeSchedule := CompositeSchedule{
		StartDateTime:    types.NewDateTime(time.Now()),
		ChargingSchedule: chargingSchedule,
	}
	var confirmationTable = []tests.GenericTestEntry{
		{GetCompositeScheduleResponse{Status: GetCompositeScheduleStatusAccepted, StatusInfo: types.NewStatusInfo("reasoncode", ""), Schedule: &compositeSchedule}, true},
		{GetCompositeScheduleResponse{Status: GetCompositeScheduleStatusAccepted, StatusInfo: types.NewStatusInfo("reasoncode", ""), Schedule: &CompositeSchedule{}}, true},
		{GetCompositeScheduleResponse{Status: GetCompositeScheduleStatusAccepted, StatusInfo: types.NewStatusInfo("reasoncode", "")}, true},
		{GetCompositeScheduleResponse{Status: GetCompositeScheduleStatusAccepted}, true},
		{GetCompositeScheduleResponse{}, false},
		{GetCompositeScheduleResponse{Status: "invalidGetCompositeScheduleStatus"}, false},
		{GetCompositeScheduleResponse{Status: GetCompositeScheduleStatusAccepted, StatusInfo: types.NewStatusInfo("invalidreasoncodeasitslongerthan20", "")}, false},
		{GetCompositeScheduleResponse{Status: GetCompositeScheduleStatusAccepted, StatusInfo: types.NewStatusInfo("", ""), Schedule: &CompositeSchedule{StartDateTime: types.NewDateTime(time.Now()), ChargingSchedule: types.NewChargingSchedule(1, "invalidChargingRateUnit")}}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *smartChargingTestSuite) TestGetCompositeScheduleFeature() {
	feature := GetCompositeScheduleFeature{}
	suite.Equal(GetCompositeScheduleFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(GetCompositeScheduleRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(GetCompositeScheduleResponse{}), feature.GetResponseType())
}

func (suite *smartChargingTestSuite) TestNewGetCompositeScheduleRequest() {
	req := NewGetCompositeScheduleRequest(600, 1)
	suite.NotNil(req)
	suite.Equal(GetCompositeScheduleFeatureName, req.GetFeatureName())
	suite.Equal(600, req.Duration)
	suite.Equal(1, req.EvseID)
}

func (suite *smartChargingTestSuite) TestNewGetCompositeScheduleResponse() {
	resp := NewGetCompositeScheduleResponse(GetCompositeScheduleStatusAccepted)
	suite.NotNil(resp)
	suite.Equal(GetCompositeScheduleFeatureName, resp.GetFeatureName())
	suite.Equal(GetCompositeScheduleStatusAccepted, resp.Status)
}
