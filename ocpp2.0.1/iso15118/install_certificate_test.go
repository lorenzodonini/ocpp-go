package iso15118

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *iso15118TestSuite) TestInstallCertificateRequestValidation() {
	t := suite.T()
	var testTable = []tests.GenericTestEntry{
		{InstallCertificateRequest{CertificateType: types.V2GRootCertificate, Certificate: "0xdeadbeef"}, true},
		{InstallCertificateRequest{CertificateType: types.MORootCertificate, Certificate: "0xdeadbeef"}, true},
		{InstallCertificateRequest{CertificateType: types.CSOSubCA1, Certificate: "0xdeadbeef"}, true},
		{InstallCertificateRequest{CertificateType: types.CSOSubCA2, Certificate: "0xdeadbeef"}, true},
		{InstallCertificateRequest{CertificateType: types.CSMSRootCertificate, Certificate: "0xdeadbeef"}, true},
		{InstallCertificateRequest{CertificateType: types.ManufacturerRootCertificate, Certificate: "0xdeadbeef"}, true},
		{InstallCertificateRequest{CertificateType: types.ManufacturerRootCertificate}, false},
		{InstallCertificateRequest{Certificate: "0xdeadbeef"}, false},
		{InstallCertificateRequest{}, false},
		{InstallCertificateRequest{CertificateType: "invalidCertificateUse", Certificate: "0xdeadbeef"}, false},
		{InstallCertificateRequest{CertificateType: types.V2GRootCertificate, Certificate: tests.NewLongString(5501)}, false},
	}
	tests.ExecuteGenericTestTable(t, testTable)
}

func (suite *iso15118TestSuite) TestInstallCertificateConfirmationValidation() {
	t := suite.T()
	var testTable = []tests.GenericTestEntry{
		{InstallCertificateResponse{Status: CertificateStatusAccepted}, true},
		{InstallCertificateResponse{Status: CertificateStatusRejected}, true},
		{InstallCertificateResponse{Status: CertificateStatusFailed}, true},
		{InstallCertificateResponse{}, false},
		{InstallCertificateResponse{Status: "invalidInstallCertificateStatus"}, false},
	}
	tests.ExecuteGenericTestTable(t, testTable)
}

func (suite *iso15118TestSuite) TestInstallCertificateFeature() {
	feature := InstallCertificateFeature{}
	suite.Equal(InstallCertificateFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(InstallCertificateRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(InstallCertificateResponse{}), feature.GetResponseType())
}

func (suite *iso15118TestSuite) TestNewInstallCertificateRequest() {
	req := NewInstallCertificateRequest(types.V2GRootCertificate, "cert-data")
	suite.NotNil(req)
	suite.Equal(InstallCertificateFeatureName, req.GetFeatureName())
	suite.Equal(types.V2GRootCertificate, req.CertificateType)
	suite.Equal("cert-data", req.Certificate)
}

func (suite *iso15118TestSuite) TestNewInstallCertificateResponse() {
	resp := NewInstallCertificateResponse(CertificateStatusAccepted)
	suite.NotNil(resp)
	suite.Equal(InstallCertificateFeatureName, resp.GetFeatureName())
	suite.Equal(CertificateStatusAccepted, resp.Status)
}
