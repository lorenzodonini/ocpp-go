package smartcharging

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *smartChargingTestSuite) TestGetChargingProfilesRequestValidation() {
	t := suite.T()
	validChargingProfileCriterion := ChargingProfileCriterion{
		ChargingProfilePurpose: types.ChargingProfilePurposeTxDefaultProfile,
		StackLevel:             tests.NewInt(2),
		ChargingProfileID:      []int{1, 2},
		ChargingLimitSource:    []types.ChargingLimitSourceType{types.ChargingLimitSourceEMS},
	}
	var requestTable = []tests.GenericTestEntry{
		{GetChargingProfilesRequest{RequestID: 42, EvseID: tests.NewInt(1), ChargingProfile: validChargingProfileCriterion}, true},
		{GetChargingProfilesRequest{RequestID: 42, ChargingProfile: validChargingProfileCriterion}, true},
		{GetChargingProfilesRequest{EvseID: tests.NewInt(1), ChargingProfile: validChargingProfileCriterion}, true},
		{GetChargingProfilesRequest{ChargingProfile: validChargingProfileCriterion}, true},
		{GetChargingProfilesRequest{ChargingProfile: ChargingProfileCriterion{}}, true},
		{GetChargingProfilesRequest{}, true},
		{GetChargingProfilesRequest{RequestID: 42, EvseID: tests.NewInt(-1), ChargingProfile: validChargingProfileCriterion}, false},
		{GetChargingProfilesRequest{ChargingProfile: ChargingProfileCriterion{ChargingProfilePurpose: "invalidChargingProfilePurpose", StackLevel: tests.NewInt(2), ChargingProfileID: []int{1, 2}, ChargingLimitSource: []types.ChargingLimitSourceType{types.ChargingLimitSourceEMS}}}, false},
		{GetChargingProfilesRequest{ChargingProfile: ChargingProfileCriterion{ChargingProfilePurpose: types.ChargingProfilePurposeTxDefaultProfile, StackLevel: tests.NewInt(-1), ChargingProfileID: []int{1, 2}, ChargingLimitSource: []types.ChargingLimitSourceType{types.ChargingLimitSourceEMS}}}, false},
		{GetChargingProfilesRequest{ChargingProfile: ChargingProfileCriterion{ChargingProfilePurpose: types.ChargingProfilePurposeTxDefaultProfile, StackLevel: tests.NewInt(2), ChargingProfileID: []int{1, 2}, ChargingLimitSource: []types.ChargingLimitSourceType{types.ChargingLimitSourceEMS, types.ChargingLimitSourceCSO, types.ChargingLimitSourceSO, types.ChargingLimitSourceOther, types.ChargingLimitSourceEMS}}}, false},
		{GetChargingProfilesRequest{ChargingProfile: ChargingProfileCriterion{ChargingProfilePurpose: types.ChargingProfilePurposeTxDefaultProfile, StackLevel: tests.NewInt(2), ChargingProfileID: []int{1, 2}, ChargingLimitSource: []types.ChargingLimitSourceType{"invalidChargingLimitSource"}}}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *smartChargingTestSuite) TestGetChargingProfilesConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{GetChargingProfilesResponse{Status: GetChargingProfileStatusAccepted}, true},
		{GetChargingProfilesResponse{Status: GetChargingProfileStatusNoProfiles}, true},
		{GetChargingProfilesResponse{Status: "invalidGetChargingProfilesStatus"}, false},
		{GetChargingProfilesResponse{}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *smartChargingTestSuite) TestGetChargingProfilesFeature() {
	feature := GetChargingProfilesFeature{}
	suite.Equal(GetChargingProfilesFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(GetChargingProfilesRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(GetChargingProfilesResponse{}), feature.GetResponseType())
}

func (suite *smartChargingTestSuite) TestNewGetChargingProfilesRequest() {
	criterion := ChargingProfileCriterion{
		ChargingProfilePurpose: types.ChargingProfilePurposeTxDefaultProfile,
	}
	req := NewGetChargingProfilesRequest(criterion)
	suite.NotNil(req)
	suite.Equal(GetChargingProfilesFeatureName, req.GetFeatureName())
	suite.Equal(criterion, req.ChargingProfile)
}

func (suite *smartChargingTestSuite) TestNewGetChargingProfilesResponse() {
	resp := NewGetChargingProfilesResponse(GetChargingProfileStatusAccepted)
	suite.NotNil(resp)
	suite.Equal(GetChargingProfilesFeatureName, resp.GetFeatureName())
	suite.Equal(GetChargingProfileStatusAccepted, resp.Status)
}
