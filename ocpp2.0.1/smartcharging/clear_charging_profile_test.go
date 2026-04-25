package smartcharging

import (
	"reflect"
	"testing"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
	"github.com/stretchr/testify/suite"
)

type smartChargingTestSuite struct {
	suite.Suite
}

func (suite *smartChargingTestSuite) TestClearChargingProfileRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{ClearChargingProfileRequest{ChargingProfileID: tests.NewInt(1), ChargingProfileCriteria: &ClearChargingProfileType{EvseID: tests.NewInt(1), ChargingProfilePurpose: types.ChargingProfilePurposeChargingStationMaxProfile, StackLevel: tests.NewInt(1)}}, true},
		{ClearChargingProfileRequest{ChargingProfileID: tests.NewInt(1), ChargingProfileCriteria: &ClearChargingProfileType{EvseID: tests.NewInt(1), ChargingProfilePurpose: types.ChargingProfilePurposeChargingStationMaxProfile}}, true},
		{ClearChargingProfileRequest{ChargingProfileID: tests.NewInt(1), ChargingProfileCriteria: &ClearChargingProfileType{EvseID: tests.NewInt(1)}}, true},
		{ClearChargingProfileRequest{ChargingProfileCriteria: &ClearChargingProfileType{EvseID: tests.NewInt(1)}}, true},
		{ClearChargingProfileRequest{ChargingProfileCriteria: &ClearChargingProfileType{}}, true},
		{ClearChargingProfileRequest{}, true},
		{ClearChargingProfileRequest{ChargingProfileCriteria: &ClearChargingProfileType{EvseID: tests.NewInt(-1)}}, false},
		{ClearChargingProfileRequest{ChargingProfileCriteria: &ClearChargingProfileType{ChargingProfilePurpose: "invalidChargingProfilePurposeType"}}, false},
		{ClearChargingProfileRequest{ChargingProfileCriteria: &ClearChargingProfileType{StackLevel: tests.NewInt(-1)}}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *smartChargingTestSuite) TestClearChargingProfileConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{ClearChargingProfileResponse{Status: ClearChargingProfileStatusAccepted, StatusInfo: types.NewStatusInfo("200", "ok")}, true},
		{ClearChargingProfileResponse{Status: ClearChargingProfileStatusAccepted}, true},
		{ClearChargingProfileResponse{Status: "invalidClearChargingProfileStatus"}, false},
		{ClearChargingProfileResponse{}, false},
		{ClearChargingProfileResponse{StatusInfo: types.NewStatusInfo("", "")}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *smartChargingTestSuite) TestClearChargingProfileFeature() {
	feature := ClearChargingProfileFeature{}
	suite.Equal(ClearChargingProfileFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(ClearChargingProfileRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(ClearChargingProfileResponse{}), feature.GetResponseType())
}

func (suite *smartChargingTestSuite) TestNewClearChargingProfileRequest() {
	req := NewClearChargingProfileRequest()
	suite.NotNil(req)
	suite.Equal(ClearChargingProfileFeatureName, req.GetFeatureName())
}

func (suite *smartChargingTestSuite) TestNewClearChargingProfileResponse() {
	resp := NewClearChargingProfileResponse(ClearChargingProfileStatusAccepted)
	suite.NotNil(resp)
	suite.Equal(ClearChargingProfileFeatureName, resp.GetFeatureName())
	suite.Equal(ClearChargingProfileStatusAccepted, resp.Status)
}

func TestSmartChargingSuite(t *testing.T) {
	suite.Run(t, new(smartChargingTestSuite))
}