package iso15118

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *iso15118TestSuite) TestGetInstalledCertificateIdsRequestValidation() {
	t := suite.T()
	var testTable = []tests.GenericTestEntry{
		{GetInstalledCertificateIdsRequest{CertificateTypes: []types.CertificateUse{types.V2GRootCertificate}}, true},
		{GetInstalledCertificateIdsRequest{CertificateTypes: []types.CertificateUse{types.MORootCertificate}}, true},
		{GetInstalledCertificateIdsRequest{CertificateTypes: []types.CertificateUse{types.CSOSubCA1}}, true},
		{GetInstalledCertificateIdsRequest{CertificateTypes: []types.CertificateUse{types.CSOSubCA2}}, true},
		{GetInstalledCertificateIdsRequest{CertificateTypes: []types.CertificateUse{types.CSMSRootCertificate}}, true},
		{GetInstalledCertificateIdsRequest{CertificateTypes: []types.CertificateUse{types.ManufacturerRootCertificate}}, true},
		{GetInstalledCertificateIdsRequest{}, true},
		{GetInstalledCertificateIdsRequest{CertificateTypes: []types.CertificateUse{"invalidCertificateUse"}}, false},
	}
	tests.ExecuteGenericTestTable(t, testTable)
}

func (suite *iso15118TestSuite) TestGetInstalledCertificateIdsConfirmationValidation() {
	t := suite.T()
	var testTable = []tests.GenericTestEntry{
		{GetInstalledCertificateIdsResponse{Status: GetInstalledCertificateStatusAccepted, CertificateHashDataChain: []types.CertificateHashDataChain{{CertificateType: types.CSMSRootCertificate, CertificateHashData: types.CertificateHashData{HashAlgorithm: types.SHA256, IssuerNameHash: "name0", IssuerKeyHash: "key0", SerialNumber: "serial0"}}}}, true},
		{GetInstalledCertificateIdsResponse{Status: GetInstalledCertificateStatusNotFound, CertificateHashDataChain: []types.CertificateHashDataChain{{CertificateType: types.CSMSRootCertificate, CertificateHashData: types.CertificateHashData{HashAlgorithm: types.SHA256, IssuerNameHash: "name0", IssuerKeyHash: "key0", SerialNumber: "serial0"}}}}, true},
		{GetInstalledCertificateIdsResponse{Status: GetInstalledCertificateStatusAccepted, CertificateHashDataChain: []types.CertificateHashDataChain{}}, true},
		{GetInstalledCertificateIdsResponse{Status: GetInstalledCertificateStatusAccepted}, true},
		{GetInstalledCertificateIdsResponse{}, false},
		{GetInstalledCertificateIdsResponse{Status: "invalidGetInstalledCertificateStatus"}, false},
		{GetInstalledCertificateIdsResponse{Status: GetInstalledCertificateStatusAccepted, CertificateHashDataChain: []types.CertificateHashDataChain{{CertificateType: types.CSMSRootCertificate, CertificateHashData: types.CertificateHashData{HashAlgorithm: "invalidHashAlgorithm", IssuerNameHash: "name0", IssuerKeyHash: "key0", SerialNumber: "serial0"}}}}, false},
	}
	tests.ExecuteGenericTestTable(t, testTable)
}

func (suite *iso15118TestSuite) TestGetInstalledCertificateIdsFeature() {
	feature := GetInstalledCertificateIdsFeature{}
	suite.Equal(GetInstalledCertificateIdsFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(GetInstalledCertificateIdsRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(GetInstalledCertificateIdsResponse{}), feature.GetResponseType())
}

func (suite *iso15118TestSuite) TestNewGetInstalledCertificateIdsRequest() {
	req := NewGetInstalledCertificateIdsRequest()
	suite.NotNil(req)
	suite.Equal(GetInstalledCertificateIdsFeatureName, req.GetFeatureName())
}

func (suite *iso15118TestSuite) TestNewGetInstalledCertificateIdsResponse() {
	resp := NewGetInstalledCertificateIdsResponse(GetInstalledCertificateStatusAccepted)
	suite.NotNil(resp)
	suite.Equal(GetInstalledCertificateIdsFeatureName, resp.GetFeatureName())
	suite.Equal(GetInstalledCertificateStatusAccepted, resp.Status)
}
