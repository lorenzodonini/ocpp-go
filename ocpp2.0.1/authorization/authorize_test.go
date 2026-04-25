package authorization

import (
	"reflect"
	"testing"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
	"github.com/stretchr/testify/suite"
)

type authTestSuite struct {
	suite.Suite
}

// Test
func (suite *authTestSuite) TestAuthorizeRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{AuthorizeRequest{Certificate: "deadc0de", IdToken: types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode, AdditionalInfo: []types.AdditionalInfo{{AdditionalIdToken: "0000", Type: "someType"}}}, CertificateHashData: []types.OCSPRequestDataType{{SerialNumber: "serial0", HashAlgorithm: types.SHA256, IssuerNameHash: "hash0", IssuerKeyHash: "hash1", ResponderURL: "www.someurl.com"}}}, true},
		{AuthorizeRequest{Certificate: "deadc0de", IdToken: types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode, AdditionalInfo: []types.AdditionalInfo{{AdditionalIdToken: "0000", Type: "someType"}}}}, true},
		{AuthorizeRequest{IdToken: types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode, AdditionalInfo: []types.AdditionalInfo{{AdditionalIdToken: "0000", Type: "someType"}}}, CertificateHashData: []types.OCSPRequestDataType{{SerialNumber: "serial0", HashAlgorithm: types.SHA256, IssuerNameHash: "hash0", IssuerKeyHash: "hash1", ResponderURL: "www.someurl.com"}}}, true},
		{AuthorizeRequest{IdToken: types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode, AdditionalInfo: []types.AdditionalInfo{{AdditionalIdToken: "0000", Type: "someType"}}}, CertificateHashData: []types.OCSPRequestDataType{}}, true},
		{AuthorizeRequest{IdToken: types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode, AdditionalInfo: []types.AdditionalInfo{{AdditionalIdToken: "0000", Type: "someType"}}}}, true},
		{AuthorizeRequest{}, false},
		{AuthorizeRequest{Certificate: tests.NewLongString(5501), IdToken: types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode, AdditionalInfo: []types.AdditionalInfo{{AdditionalIdToken: "0000", Type: "someType"}}}}, false},
		{AuthorizeRequest{Certificate: "deadc0de", IdToken: types.IdToken{Type: types.IdTokenTypeKeyCode, AdditionalInfo: []types.AdditionalInfo{{AdditionalIdToken: "0000", Type: "someType"}}}}, false},
		{AuthorizeRequest{Certificate: "deadc0de", IdToken: types.IdToken{Type: types.IdTokenTypeKeyCode, AdditionalInfo: []types.AdditionalInfo{{AdditionalIdToken: "0000", Type: "someType"}}}, CertificateHashData: []types.OCSPRequestDataType{{HashAlgorithm: types.SHA256, IssuerNameHash: "hash0", IssuerKeyHash: "hash1"}}}, false},
		{AuthorizeRequest{Certificate: "deadc0de", IdToken: types.IdToken{Type: types.IdTokenTypeKeyCode, AdditionalInfo: []types.AdditionalInfo{{AdditionalIdToken: "0000", Type: "someType"}}}, CertificateHashData: []types.OCSPRequestDataType{{SerialNumber: "s0", HashAlgorithm: types.SHA256, IssuerNameHash: "h0", IssuerKeyHash: "h0.1"}, {SerialNumber: "s1", HashAlgorithm: types.SHA256, IssuerNameHash: "h1", IssuerKeyHash: "h1.1"}, {SerialNumber: "s2", HashAlgorithm: types.SHA256, IssuerNameHash: "h2", IssuerKeyHash: "h2.1"}, {SerialNumber: "s3", HashAlgorithm: types.SHA256, IssuerNameHash: "h3", IssuerKeyHash: "h3.1"}, {SerialNumber: "s4", HashAlgorithm: types.SHA256, IssuerNameHash: "h4", IssuerKeyHash: "h4.1"}}}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *authTestSuite) TestAuthorizeConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{AuthorizeResponse{CertificateStatus: CertificateStatusAccepted, IdTokenInfo: types.IdTokenInfo{Status: types.AuthorizationStatusAccepted}}, true},
		{AuthorizeResponse{CertificateStatus: CertificateStatusAccepted, IdTokenInfo: types.IdTokenInfo{Status: types.AuthorizationStatusAccepted}}, true},
		{AuthorizeResponse{IdTokenInfo: types.IdTokenInfo{Status: types.AuthorizationStatusAccepted}}, true},
		{AuthorizeResponse{}, false},
		{AuthorizeResponse{CertificateStatus: "invalidCertificateStatus", IdTokenInfo: types.IdTokenInfo{Status: types.AuthorizationStatusAccepted}}, false},
		{AuthorizeResponse{CertificateStatus: CertificateStatusAccepted, IdTokenInfo: types.IdTokenInfo{Status: "invalidTokenInfoStatus"}}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *authTestSuite) TestAuthorizeFeature() {
	feature := AuthorizeFeature{}
	suite.Equal(AuthorizeFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(AuthorizeRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(AuthorizeResponse{}), feature.GetResponseType())
}

func (suite *authTestSuite) TestNewAuthorizationRequest() {
	req := NewAuthorizationRequest("1234", types.IdTokenTypeKeyCode)
	suite.NotNil(req)
	suite.Equal(AuthorizeFeatureName, req.GetFeatureName())
	suite.Equal("1234", req.IdToken.IdToken)
	suite.Equal(types.IdTokenTypeKeyCode, req.IdToken.Type)
}

func (suite *authTestSuite) TestNewAuthorizationResponse() {
	idTokenInfo := types.IdTokenInfo{Status: types.AuthorizationStatusAccepted}
	resp := NewAuthorizationResponse(idTokenInfo)
	suite.NotNil(resp)
	suite.Equal(AuthorizeFeatureName, resp.GetFeatureName())
	suite.Equal(idTokenInfo, resp.IdTokenInfo)
}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(authTestSuite))
}
