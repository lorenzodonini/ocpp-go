package diagnostics

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *diagnosticsTestSuite) TestCustomerInformationRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{CustomerInformationRequest{RequestID: 42, Report: true, Clear: true, CustomerIdentifier: "0001", IdToken: &types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode, AdditionalInfo: nil}, CustomerCertificate: &types.CertificateHashData{HashAlgorithm: types.SHA256, IssuerNameHash: "hash00", IssuerKeyHash: "hash01", SerialNumber: "serial0"}}, true},
		{CustomerInformationRequest{RequestID: 42, Report: true, Clear: true, CustomerIdentifier: "0001", IdToken: &types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode, AdditionalInfo: nil}}, true},
		{CustomerInformationRequest{RequestID: 42, Report: true, Clear: true, CustomerIdentifier: "0001"}, true},
		{CustomerInformationRequest{RequestID: 42, Report: true, Clear: true}, true},
		{CustomerInformationRequest{RequestID: 42, Report: true}, true},
		{CustomerInformationRequest{RequestID: 42, Clear: true}, true},
		{CustomerInformationRequest{Report: true, Clear: true}, true},
		{CustomerInformationRequest{}, true},
		{CustomerInformationRequest{RequestID: -1, Report: true, Clear: true}, false},
		{CustomerInformationRequest{RequestID: 42, Report: true, Clear: true, CustomerIdentifier: ">64.............................................................."}, false},
		{CustomerInformationRequest{RequestID: 42, Report: true, Clear: true, IdToken: &types.IdToken{IdToken: "1234", Type: "invalidTokenType", AdditionalInfo: nil}}, false},
		{CustomerInformationRequest{RequestID: 42, Report: true, Clear: true, CustomerCertificate: &types.CertificateHashData{HashAlgorithm: "invalidHasAlgorithm", IssuerNameHash: "hash00", IssuerKeyHash: "hash01", SerialNumber: "serial0"}}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *diagnosticsTestSuite) TestCustomerInformationConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{CustomerInformationResponse{Status: CustomerInformationStatusAccepted}, true},
		{CustomerInformationResponse{}, false},
		{CustomerInformationResponse{Status: "invalidCustomerInformationStatus"}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *diagnosticsTestSuite) TestCustomerInformationFeature() {
	feature := CustomerInformationFeature{}
	suite.Equal(CustomerInformationFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(CustomerInformationRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(CustomerInformationResponse{}), feature.GetResponseType())
}

func (suite *diagnosticsTestSuite) TestNewCustomerInformationRequest() {
	req := NewCustomerInformationRequest(42, true, false)
	suite.NotNil(req)
	suite.Equal(CustomerInformationFeatureName, req.GetFeatureName())
	suite.Equal(42, req.RequestID)
	suite.True(req.Report)
	suite.False(req.Clear)
}

func (suite *diagnosticsTestSuite) TestNewCustomerInformationResponse() {
	resp := NewCustomerInformationResponse(CustomerInformationStatusAccepted)
	suite.NotNil(resp)
	suite.Equal(CustomerInformationFeatureName, resp.GetFeatureName())
	suite.Equal(CustomerInformationStatusAccepted, resp.Status)
}
