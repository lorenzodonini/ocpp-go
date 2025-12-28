package ocpp2_test

import (
	"fmt"

	"github.com/stretchr/testify/mock"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/security"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

// Test
func (suite *OcppV2TestSuite) TestSignCertificateRequestValidation() {
	var requestTable = []GenericTestEntry{
		{security.SignCertificateRequest{CSR: "deadc0de", CertificateType: types.ChargingStationCert}, true},
		{security.SignCertificateRequest{CSR: "deadc0de", CertificateType: types.V2GCertificate}, true},
		{security.SignCertificateRequest{CSR: "deadc0de"}, true},
		{security.SignCertificateRequest{}, false},
		{security.SignCertificateRequest{CSR: "deadc0de", CertificateType: "invalidType"}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV2TestSuite) TestSignCertificateConfirmationValidation() {
	var confirmationTable = []GenericTestEntry{
		{security.SignCertificateResponse{Status: types.GenericStatusAccepted, StatusInfo: types.NewStatusInfo("200", "")}, true},
		{security.SignCertificateResponse{Status: types.GenericStatusAccepted}, true},
		{security.SignCertificateResponse{}, false},
		{security.SignCertificateResponse{Status: types.GenericStatusAccepted, StatusInfo: types.NewStatusInfo("", "")}, false},
		{security.SignCertificateResponse{Status: "invalidStatus", StatusInfo: types.NewStatusInfo("200", "")}, false},
	}
	ExecuteGenericTestTable(suite, confirmationTable)
}

func (suite *OcppV2TestSuite) TestSignCertificateE2EMocked() {
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	csr := "deadc0de"
	certificateType := types.ChargingStationCert
	status := types.GenericStatusAccepted
	statusInfo := types.NewStatusInfo("200", "")
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"csr":"%v","certificateType":"%v"}]`,
		messageId, security.SignCertificateFeatureName, csr, certificateType)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v","statusInfo":{"reasonCode":"%v"}}]`,
		messageId, status, statusInfo.ReasonCode)
	signCertificateResponse := security.NewSignCertificateResponse(status)
	signCertificateResponse.StatusInfo = statusInfo
	channel := NewMockWebSocket(wsId)

	handler := &MockCSMSSecurityHandler{}
	handler.On("OnSignCertificate", mock.AnythingOfType("string"), mock.Anything).Return(signCertificateResponse, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*security.SignCertificateRequest)
		suite.Require().True(ok)
		suite.Require().NotNil(request)
		suite.Equal(csr, request.CSR)
		suite.Equal(certificateType, request.CertificateType)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	suite.Require().Nil(err)
	response, err := suite.chargingStation.SignCertificate(csr, func(request *security.SignCertificateRequest) {
		request.CertificateType = certificateType
	})
	suite.Require().Nil(err)
	suite.Require().NotNil(response)
	suite.Equal(status, response.Status)
	suite.Equal(statusInfo.ReasonCode, response.StatusInfo.ReasonCode)
	suite.Equal(statusInfo.AdditionalInfo, response.StatusInfo.AdditionalInfo)
}

func (suite *OcppV2TestSuite) TestSignCertificateInvalidEndpoint() {
	messageId := defaultMessageId
	csr := "deadc0de"
	certificateType := types.ChargingStationCert
	request := security.NewSignCertificateRequest(csr)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"csr":"%v","certificateType":"%v"}]`,
		messageId, security.SignCertificateFeatureName, csr, certificateType)
	testUnsupportedRequestFromCentralSystem(suite, request, requestJson, messageId)
}
