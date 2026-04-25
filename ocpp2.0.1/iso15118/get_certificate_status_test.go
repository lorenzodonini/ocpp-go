package iso15118

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *iso15118TestSuite) TestGetCertificateStatusRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{GetCertificateStatusRequest{OcspRequestData: types.OCSPRequestDataType{HashAlgorithm: types.SHA256, IssuerNameHash: "hash00", IssuerKeyHash: "hash01", SerialNumber: "serial0", ResponderURL: "http://someUrl"}}, true},
		{GetCertificateStatusRequest{}, false},
		{GetCertificateStatusRequest{OcspRequestData: types.OCSPRequestDataType{HashAlgorithm: "invalidHashAlgorithm", IssuerNameHash: "hash00", IssuerKeyHash: "hash01", SerialNumber: "serial0", ResponderURL: "http://someUrl"}}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *iso15118TestSuite) TestGetCertificateStatusConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{GetCertificateStatusResponse{Status: types.GenericStatusAccepted, OcspResult: "deadbeef"}, true},
		{GetCertificateStatusResponse{Status: types.GenericStatusAccepted}, true},
		{GetCertificateStatusResponse{Status: types.GenericStatusRejected}, true},
		{GetCertificateStatusResponse{Status: "invalidGenericStatus"}, false},
		{GetCertificateStatusResponse{Status: types.GenericStatusAccepted, OcspResult: tests.NewLongString(5501)}, false},
		{GetCertificateStatusResponse{}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *iso15118TestSuite) TestGetCertificateStatusFeature() {
	feature := GetCertificateStatusFeature{}
	suite.Equal(GetCertificateStatusFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(GetCertificateStatusRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(GetCertificateStatusResponse{}), feature.GetResponseType())
}

func (suite *iso15118TestSuite) TestNewGetCertificateStatusRequest() {
	ocspData := types.OCSPRequestDataType{
		HashAlgorithm:  types.SHA256,
		IssuerNameHash: "hash0",
		IssuerKeyHash:  "hash1",
		SerialNumber:   "serial0",
		ResponderURL:   "http://example.com",
	}
	req := NewGetCertificateStatusRequest(ocspData)
	suite.NotNil(req)
	suite.Equal(GetCertificateStatusFeatureName, req.GetFeatureName())
	suite.Equal(ocspData, req.OcspRequestData)
}

func (suite *iso15118TestSuite) TestNewGetCertificateStatusResponse() {
	resp := NewGetCertificateStatusResponse(types.GenericStatusAccepted)
	suite.NotNil(resp)
	suite.Equal(GetCertificateStatusFeatureName, resp.GetFeatureName())
	suite.Equal(types.GenericStatusAccepted, resp.Status)
}
