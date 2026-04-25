package smartcharging

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *smartChargingTestSuite) TestSetChargingProfileRequestValidation() {
	t := suite.T()
	schedule := types.NewChargingSchedule(1, types.ChargingRateUnitWatts, types.NewChargingSchedulePeriod(0, 200.0))
	chargingProfile := types.NewChargingProfile(
		1,
		0,
		types.ChargingProfilePurposeChargingStationMaxProfile,
		types.ChargingProfileKindAbsolute,
		[]types.ChargingSchedule{*schedule})
	var requestTable = []tests.GenericTestEntry{
		{SetChargingProfileRequest{EvseID: 1, ChargingProfile: chargingProfile}, true},
		{SetChargingProfileRequest{ChargingProfile: chargingProfile}, true},
		{SetChargingProfileRequest{}, false},
		{SetChargingProfileRequest{EvseID: 1, ChargingProfile: types.NewChargingProfile(1, -1, types.ChargingProfilePurposeChargingStationMaxProfile, types.ChargingProfileKindAbsolute, []types.ChargingSchedule{*schedule})}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *smartChargingTestSuite) TestSetChargingProfileResponseValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{SetChargingProfileResponse{Status: ChargingProfileStatusAccepted, StatusInfo: types.NewStatusInfo("200", "")}, true},
		{SetChargingProfileResponse{Status: ChargingProfileStatusAccepted}, true},
		{SetChargingProfileResponse{}, false},
		{SetChargingProfileResponse{Status: "invalidChargingProfileStatus"}, false},
		{SetChargingProfileResponse{Status: ChargingProfileStatusAccepted, StatusInfo: types.NewStatusInfo("", "")}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *smartChargingTestSuite) TestSetChargingProfileFeature() {
	feature := SetChargingProfileFeature{}
	suite.Equal(SetChargingProfileFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(SetChargingProfileRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(SetChargingProfileResponse{}), feature.GetResponseType())
}

func (suite *smartChargingTestSuite) TestNewSetChargingProfileRequest() {
	schedule := types.NewChargingSchedule(1, types.ChargingRateUnitWatts, types.NewChargingSchedulePeriod(0, 10.0))
	profile := types.NewChargingProfile(1, 0, types.ChargingProfilePurposeTxDefaultProfile, types.ChargingProfileKindAbsolute, []types.ChargingSchedule{*schedule})
	req := NewSetChargingProfileRequest(1, profile)
	suite.NotNil(req)
	suite.Equal(SetChargingProfileFeatureName, req.GetFeatureName())
	suite.Equal(1, req.EvseID)
	suite.Equal(profile, req.ChargingProfile)
}

func (suite *smartChargingTestSuite) TestNewSetChargingProfileResponse() {
	resp := NewSetChargingProfileResponse(ChargingProfileStatusAccepted)
	suite.NotNil(resp)
	suite.Equal(SetChargingProfileFeatureName, resp.GetFeatureName())
	suite.Equal(ChargingProfileStatusAccepted, resp.Status)
}
