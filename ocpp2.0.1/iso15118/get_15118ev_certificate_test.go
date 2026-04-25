package iso15118

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *iso15118TestSuite) TestGet15118EVCertificateRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{Get15118EVCertificateRequest{SchemaVersion: "1.0", Action: CertificateActionInstall, ExiRequest: "deadbeef"}, true},
		{Get15118EVCertificateRequest{SchemaVersion: "1.0", Action: CertificateActionUpdate, ExiRequest: "deadbeef"}, true},
		{Get15118EVCertificateRequest{SchemaVersion: "1.0", Action: CertificateActionInstall}, false},
		{Get15118EVCertificateRequest{ExiRequest: "deadbeef"}, false},
		{Get15118EVCertificateRequest{}, false},
		{Get15118EVCertificateRequest{SchemaVersion: ">50................................................", Action: CertificateActionInstall, ExiRequest: "deadbeef"}, false},
		{Get15118EVCertificateRequest{SchemaVersion: "1.0", Action: "invalidCertificateAction", ExiRequest: "deadbeef"}, false},
		{Get15118EVCertificateRequest{SchemaVersion: "1.0", Action: CertificateActionInstall, ExiRequest: tests.NewLongString(5601)}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *iso15118TestSuite) TestGet15118EVCertificateConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{Get15118EVCertificateResponse{Status: types.Certificate15188EVStatusAccepted, ExiResponse: "deadbeef", StatusInfo: types.NewStatusInfo("200", "ok")}, true},
		{Get15118EVCertificateResponse{Status: types.Certificate15188EVStatusAccepted, ExiResponse: "deadbeef"}, true},
		{Get15118EVCertificateResponse{Status: types.Certificate15188EVStatusAccepted}, false},
		{Get15118EVCertificateResponse{ExiResponse: "deadbeef"}, false},
		{Get15118EVCertificateResponse{}, false},
		{Get15118EVCertificateResponse{Status: "invalidCertificateStatus", ExiResponse: "deadbeef", StatusInfo: types.NewStatusInfo("200", "ok")}, false},
		{Get15118EVCertificateResponse{Status: types.Certificate15188EVStatusAccepted, ExiResponse: tests.NewLongString(5601), StatusInfo: types.NewStatusInfo("200", "ok")}, false},
		{Get15118EVCertificateResponse{Status: types.Certificate15188EVStatusAccepted, ExiResponse: "deadbeef", StatusInfo: types.NewStatusInfo("", "")}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *iso15118TestSuite) TestGet15118EVCertificateFeature() {
	feature := Get15118EVCertificateFeature{}
	suite.Equal(Get15118EVCertificateFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(Get15118EVCertificateRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(Get15118EVCertificateResponse{}), feature.GetResponseType())
}

func (suite *iso15118TestSuite) TestNewGet15118EVCertificateRequest() {
	req := NewGet15118EVCertificateRequest("1.0", CertificateActionInstall, "exi-data")
	suite.NotNil(req)
	suite.Equal(Get15118EVCertificateFeatureName, req.GetFeatureName())
	suite.Equal("1.0", req.SchemaVersion)
	suite.Equal(CertificateActionInstall, req.Action)
	suite.Equal("exi-data", req.ExiRequest)
}

func (suite *iso15118TestSuite) TestNewGet15118EVCertificateResponse() {
	resp := NewGet15118EVCertificateResponse(types.Certificate15188EVStatusAccepted, "exi-response")
	suite.NotNil(resp)
	suite.Equal(Get15118EVCertificateFeatureName, resp.GetFeatureName())
	suite.Equal(types.Certificate15188EVStatusAccepted, resp.Status)
	suite.Equal("exi-response", resp.ExiResponse)
}
