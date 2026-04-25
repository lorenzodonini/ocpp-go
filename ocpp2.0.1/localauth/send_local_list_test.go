package localauth

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *localAuthTestSuite) TestSendLocalListRequestValidation() {
	t := suite.T()
	authData := AuthorizationData{
		IdToken: types.IdToken{
			IdToken:        "token1",
			Type:           types.IdTokenTypeKeyCode,
			AdditionalInfo: nil,
		},
		IdTokenInfo: types.NewIdTokenInfo(types.AuthorizationStatusAccepted),
	}
	var requestTable = []tests.GenericTestEntry{
		{SendLocalListRequest{VersionNumber: 42, UpdateType: UpdateTypeDifferential, LocalAuthorizationList: []AuthorizationData{authData}}, true},
		{SendLocalListRequest{VersionNumber: 42, UpdateType: UpdateTypeFull, LocalAuthorizationList: []AuthorizationData{authData}}, true},
		{SendLocalListRequest{VersionNumber: 42, UpdateType: UpdateTypeDifferential, LocalAuthorizationList: []AuthorizationData{}}, true},
		{SendLocalListRequest{VersionNumber: 42, UpdateType: UpdateTypeDifferential}, true},
		{SendLocalListRequest{UpdateType: UpdateTypeDifferential}, true},
		{SendLocalListRequest{}, false},
		{SendLocalListRequest{VersionNumber: -1, UpdateType: UpdateTypeDifferential, LocalAuthorizationList: []AuthorizationData{authData}}, false},
		{SendLocalListRequest{VersionNumber: 42, UpdateType: "invalidUpdateType", LocalAuthorizationList: []AuthorizationData{{IdToken: types.IdToken{IdToken: "tokenWithoutType"}}}}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *localAuthTestSuite) TestSendLocalListResponseValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{SendLocalListResponse{Status: SendLocalListStatusAccepted, StatusInfo: types.NewStatusInfo("200", "")}, true},
		{SendLocalListResponse{Status: SendLocalListStatusAccepted}, true},
		{SendLocalListResponse{}, false},
		{SendLocalListResponse{Status: "invalidStatus", StatusInfo: types.NewStatusInfo("200", "")}, false},
		{SendLocalListResponse{Status: SendLocalListStatusAccepted, StatusInfo: types.NewStatusInfo("", "")}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *localAuthTestSuite) TestSendLocalListFeature() {
	feature := SendLocalListFeature{}
	suite.Equal(SendLocalListFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(SendLocalListRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(SendLocalListResponse{}), feature.GetResponseType())
}

func (suite *localAuthTestSuite) TestNewSendLocalListRequest() {
	req := NewSendLocalListRequest(3, UpdateTypeDifferential)
	suite.NotNil(req)
	suite.Equal(SendLocalListFeatureName, req.GetFeatureName())
	suite.Equal(3, req.VersionNumber)
	suite.Equal(UpdateTypeDifferential, req.UpdateType)
}

func (suite *localAuthTestSuite) TestNewSendLocalListResponse() {
	resp := NewSendLocalListResponse(SendLocalListStatusAccepted)
	suite.NotNil(resp)
	suite.Equal(SendLocalListFeatureName, resp.GetFeatureName())
	suite.Equal(SendLocalListStatusAccepted, resp.Status)
}
