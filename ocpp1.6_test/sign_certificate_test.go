package ocpp16_test

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/security"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6_test/mocks"
	"github.com/stretchr/testify/mock"
)

// Test
func (suite *OcppV16TestSuite) TestSignCertificateRequestValidation() {
	var requestTable = []GenericTestEntry{
		{security.SignCertificateRequest{CSR: "deadc0de", CertificateType: types.ChargingStationCert}, true},
		{security.SignCertificateRequest{CSR: "deadc0de"}, true},
		{security.SignCertificateRequest{}, false},
		{security.SignCertificateRequest{CSR: "deadc0de", CertificateType: "invalidType"}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV16TestSuite) TestSignCertificateConfirmationValidation() {
	var confirmationTable = []GenericTestEntry{
		{security.SignCertificateResponse{Status: types.GenericStatusAccepted}, true},
		{security.SignCertificateResponse{Status: types.GenericStatusAccepted}, true},
		{security.SignCertificateResponse{}, false},
		{security.SignCertificateResponse{Status: "invalidStatus"}, false},
	}
	ExecuteGenericTestTable(suite, confirmationTable)
}

func (suite *OcppV16TestSuite) TestSignCertificateE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	csr := "deadc0de"
	certificateType := types.ChargingStationCert
	status := types.GenericStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"csr":"%v","certificateType":"%v"}]`,
		messageId, security.SignCertificateFeatureName, csr, certificateType)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	signCertificateResponse := security.NewSignCertificateResponse(status)
	channel := NewMockWebSocket(wsId)

	handler := mocks.NewMockSecurityCentralSystemHandler(t)
	handler.EXPECT().OnSignCertificate(wsId, mock.Anything).RunAndReturn(func(s string, request *security.SignCertificateRequest) (*security.SignCertificateResponse, error) {
		suite.Equal(csr, request.CSR)
		suite.Equal(certificateType, request.CertificateType)
		return signCertificateResponse, nil
	})

	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	suite.centralSystem.SetSecurityHandler(handler)

	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	suite.Require().Nil(err)
	response, err := suite.chargePoint.SignCertificate(csr, func(request *security.SignCertificateRequest) {
		request.CertificateType = certificateType
	})
	suite.Require().Nil(err)
	suite.Require().NotNil(response)
	suite.Equal(status, response.Status)
}

func (suite *OcppV16TestSuite) TestSignCertificateInvalidEndpoint() {
	messageId := defaultMessageId
	csr := "deadc0de"
	certificateType := types.ChargingStationCert
	request := security.NewSignCertificateRequest(csr)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"csr":"%v","certificateType":"%v"}]`,
		messageId, security.SignCertificateFeatureName, csr, certificateType)
	testUnsupportedRequestFromCentralSystem(suite, request, requestJson, messageId)
}
