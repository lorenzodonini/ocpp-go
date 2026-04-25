package security

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *securityTestSuite) TestSignCertificateRequestValidation() {
	var requestTable = []tests.GenericTestEntry{
		{SignCertificateRequest{CSR: "deadc0de", CertificateType: types.ChargingStationCert}, true},
		{SignCertificateRequest{CSR: "deadc0de", CertificateType: types.V2GCertificate}, true},
		{SignCertificateRequest{CSR: "deadc0de"}, true},
		{SignCertificateRequest{}, false},
		{SignCertificateRequest{CSR: "deadc0de", CertificateType: "invalidType"}, false},
	}
	tests.ExecuteGenericTestTable(suite.T(), requestTable)
}

func (suite *securityTestSuite) TestSignCertificateConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{SignCertificateResponse{Status: types.GenericStatusAccepted, StatusInfo: types.NewStatusInfo("200", "")}, true},
		{SignCertificateResponse{Status: types.GenericStatusAccepted}, true},
		{SignCertificateResponse{}, false},
		{SignCertificateResponse{Status: types.GenericStatusAccepted, StatusInfo: types.NewStatusInfo("", "")}, false},
		{SignCertificateResponse{Status: "invalidStatus", StatusInfo: types.NewStatusInfo("200", "")}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *securityTestSuite) TestSignCertificateFeature() {
	feature := SignCertificateFeature{}
	suite.Equal(SignCertificateFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(SignCertificateRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(SignCertificateResponse{}), feature.GetResponseType())
}

func (suite *securityTestSuite) TestNewSignCertificateRequest() {
	req := NewSignCertificateRequest("csr-data")
	suite.NotNil(req)
	suite.Equal(SignCertificateFeatureName, req.GetFeatureName())
	suite.Equal("csr-data", req.CSR)
}

func (suite *securityTestSuite) TestNewSignCertificateResponse() {
	resp := NewSignCertificateResponse(types.GenericStatusAccepted)
	suite.NotNil(resp)
	suite.Equal(SignCertificateFeatureName, resp.GetFeatureName())
	suite.Equal(types.GenericStatusAccepted, resp.Status)
}
