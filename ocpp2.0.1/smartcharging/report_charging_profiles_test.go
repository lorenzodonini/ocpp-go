package smartcharging

import (
	"reflect"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *smartChargingTestSuite) TestReportChargingProfilesRequestValidation() {
	t := suite.T()
	chargingSchedule := types.ChargingSchedule{
		ID:                     1,
		StartSchedule:          types.NewDateTime(time.Now()),
		Duration:               tests.NewInt(600),
		ChargingRateUnit:       types.ChargingRateUnitWatts,
		MinChargingRate:        tests.NewFloat(6.0),
		ChargingSchedulePeriod: []types.ChargingSchedulePeriod{types.NewChargingSchedulePeriod(0, 10.0)},
	}
	chargingProfile := types.NewChargingProfile(1, 0, types.ChargingProfilePurposeTxDefaultProfile, types.ChargingProfileKindAbsolute, []types.ChargingSchedule{chargingSchedule})
	var requestTable = []tests.GenericTestEntry{
		{ReportChargingProfilesRequest{RequestID: 42, ChargingLimitSource: types.ChargingLimitSourceCSO, Tbc: true, EvseID: 1, ChargingProfile: []types.ChargingProfile{*chargingProfile}}, true},
		{ReportChargingProfilesRequest{RequestID: 42, ChargingLimitSource: types.ChargingLimitSourceCSO, EvseID: 1, ChargingProfile: []types.ChargingProfile{*chargingProfile}}, true},
		{ReportChargingProfilesRequest{RequestID: 42, ChargingLimitSource: types.ChargingLimitSourceCSO, ChargingProfile: []types.ChargingProfile{*chargingProfile}}, true},
		{ReportChargingProfilesRequest{ChargingLimitSource: types.ChargingLimitSourceCSO, ChargingProfile: []types.ChargingProfile{*chargingProfile}}, true},
		{ReportChargingProfilesRequest{ChargingLimitSource: types.ChargingLimitSourceCSO, ChargingProfile: []types.ChargingProfile{}}, false},
		{ReportChargingProfilesRequest{ChargingLimitSource: types.ChargingLimitSourceCSO}, false},
		{ReportChargingProfilesRequest{ChargingProfile: []types.ChargingProfile{*chargingProfile}}, false},
		{ReportChargingProfilesRequest{}, false},
		{ReportChargingProfilesRequest{RequestID: -1, ChargingLimitSource: types.ChargingLimitSourceCSO, Tbc: true, EvseID: 1, ChargingProfile: []types.ChargingProfile{*chargingProfile}}, false},
		{ReportChargingProfilesRequest{RequestID: 42, ChargingLimitSource: types.ChargingLimitSourceCSO, Tbc: true, EvseID: -1, ChargingProfile: []types.ChargingProfile{*chargingProfile}}, false},
		{ReportChargingProfilesRequest{RequestID: 42, ChargingLimitSource: "invalidChargingLimitSource", Tbc: true, EvseID: 1, ChargingProfile: []types.ChargingProfile{*chargingProfile}}, false},
		{ReportChargingProfilesRequest{RequestID: 42, ChargingLimitSource: types.ChargingLimitSourceCSO, Tbc: true, EvseID: 1, ChargingProfile: []types.ChargingProfile{
			*types.NewChargingProfile(1, -1, types.ChargingProfilePurposeTxDefaultProfile, types.ChargingProfileKindAbsolute, []types.ChargingSchedule{chargingSchedule})}}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *smartChargingTestSuite) TestReportChargingProfilesResponseValidation() {
	t := suite.T()
	var responseTable = []tests.GenericTestEntry{
		{ReportChargingProfilesResponse{}, true},
	}
	tests.ExecuteGenericTestTable(t, responseTable)
}

func (suite *smartChargingTestSuite) TestReportChargingProfilesFeature() {
	feature := ReportChargingProfilesFeature{}
	suite.Equal(ReportChargingProfilesFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(ReportChargingProfilesRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(ReportChargingProfilesResponse{}), feature.GetResponseType())
}

func (suite *smartChargingTestSuite) TestNewReportChargingProfilesRequest() {
	schedule := types.NewChargingSchedule(1, types.ChargingRateUnitWatts, types.NewChargingSchedulePeriod(0, 10.0))
	profile := types.NewChargingProfile(1, 0, types.ChargingProfilePurposeTxDefaultProfile, types.ChargingProfileKindAbsolute, []types.ChargingSchedule{*schedule})
	req := NewReportChargingProfilesRequest(42, types.ChargingLimitSourceCSO, 1, []types.ChargingProfile{*profile})
	suite.NotNil(req)
	suite.Equal(ReportChargingProfilesFeatureName, req.GetFeatureName())
	suite.Equal(42, req.RequestID)
	suite.Equal(types.ChargingLimitSourceCSO, req.ChargingLimitSource)
	suite.Equal(1, req.EvseID)
}

func (suite *smartChargingTestSuite) TestNewReportChargingProfilesResponse() {
	resp := NewReportChargingProfilesResponse()
	suite.NotNil(resp)
	suite.Equal(ReportChargingProfilesFeatureName, resp.GetFeatureName())
}
