package iso15118

import (
	"reflect"
	"testing"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
	"github.com/stretchr/testify/suite"
)

type iso15118TestSuite struct {
	suite.Suite
}

func (suite *iso15118TestSuite) TestDeleteCertificateRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{DeleteCertificateRequest{CertificateHashData: types.CertificateHashData{HashAlgorithm: types.SHA256, IssuerNameHash: "hash00", IssuerKeyHash: "hash01", SerialNumber: "serial0"}}, true},
		{DeleteCertificateRequest{}, false},
		{DeleteCertificateRequest{CertificateHashData: types.CertificateHashData{HashAlgorithm: "invalidHashAlgorithm", IssuerNameHash: "hash00", IssuerKeyHash: "hash01", SerialNumber: "serial0"}}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *iso15118TestSuite) TestDeleteCertificateConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{DeleteCertificateResponse{Status: DeleteCertificateStatusAccepted}, true},
		{DeleteCertificateResponse{Status: DeleteCertificateStatusFailed}, true},
		{DeleteCertificateResponse{Status: DeleteCertificateStatusNotFound}, true},
		{DeleteCertificateResponse{Status: "invalidDeleteCertificateStatus"}, false},
		{DeleteCertificateResponse{}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *iso15118TestSuite) TestDeleteCertificateFeature() {
	feature := DeleteCertificateFeature{}
	suite.Equal(DeleteCertificateFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(DeleteCertificateRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(DeleteCertificateResponse{}), feature.GetResponseType())
}

func (suite *iso15118TestSuite) TestNewDeleteCertificateRequest() {
	hashData := types.CertificateHashData{
		HashAlgorithm:  types.SHA256,
		IssuerNameHash: "hash0",
		IssuerKeyHash:  "hash1",
		SerialNumber:   "serial0",
	}
	req := NewDeleteCertificateRequest(hashData)
	suite.NotNil(req)
	suite.Equal(DeleteCertificateFeatureName, req.GetFeatureName())
	suite.Equal(hashData, req.CertificateHashData)
}

func (suite *iso15118TestSuite) TestNewDeleteCertificateResponse() {
	resp := NewDeleteCertificateResponse(DeleteCertificateStatusAccepted)
	suite.NotNil(resp)
	suite.Equal(DeleteCertificateFeatureName, resp.GetFeatureName())
	suite.Equal(DeleteCertificateStatusAccepted, resp.Status)
}

func TestIso15118Suite(t *testing.T) {
	suite.Run(t, new(iso15118TestSuite))
}