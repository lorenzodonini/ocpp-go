package localauth

import (
	"reflect"
	"testing"

	"github.com/lorenzodonini/ocpp-go/tests"
	"github.com/stretchr/testify/suite"
)

type localAuthTestSuite struct {
	suite.Suite
}

func (suite *localAuthTestSuite) TestGetLocalListVersionRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{GetLocalListVersionRequest{}, true},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *localAuthTestSuite) TestGetLocalListVersionConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{GetLocalListVersionResponse{VersionNumber: 1}, true},
		{GetLocalListVersionResponse{VersionNumber: 0}, true},
		{GetLocalListVersionResponse{}, true},
		{GetLocalListVersionResponse{VersionNumber: -1}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *localAuthTestSuite) TestGetLocalListVersionFeature() {
	feature := GetLocalListVersionFeature{}
	suite.Equal(GetLocalListVersionFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(GetLocalListVersionRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(GetLocalListVersionResponse{}), feature.GetResponseType())
}

func (suite *localAuthTestSuite) TestNewGetLocalListVersionRequest() {
	req := NewGetLocalListVersionRequest()
	suite.NotNil(req)
	suite.Equal(GetLocalListVersionFeatureName, req.GetFeatureName())
}

func (suite *localAuthTestSuite) TestNewGetLocalListVersionResponse() {
	resp := NewGetLocalListVersionResponse(5)
	suite.NotNil(resp)
	suite.Equal(GetLocalListVersionFeatureName, resp.GetFeatureName())
	suite.Equal(5, resp.VersionNumber)
}

func TestLocalAuthSuite(t *testing.T) {
	suite.Run(t, new(localAuthTestSuite))
}