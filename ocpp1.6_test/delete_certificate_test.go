package ocpp16_test

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/certificates"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6_test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test
func (suite *OcppV16TestSuite) TestDeleteCertificateRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{certificates.DeleteCertificateRequest{CertificateHashData: types.CertificateHashData{HashAlgorithm: types.SHA256, IssuerNameHash: "hash00", IssuerKeyHash: "hash01", SerialNumber: "serial0"}}, true},
		{certificates.DeleteCertificateRequest{}, false},
		{certificates.DeleteCertificateRequest{CertificateHashData: types.CertificateHashData{HashAlgorithm: "invalidHashAlgorithm", IssuerNameHash: "hash00", IssuerKeyHash: "hash01", SerialNumber: "serial0"}}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV16TestSuite) TestDeleteCertificateConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{certificates.DeleteCertificateResponse{Status: certificates.DeleteCertificateStatusAccepted}, true},
		{certificates.DeleteCertificateResponse{Status: certificates.DeleteCertificateStatusFailed}, true},
		{certificates.DeleteCertificateResponse{Status: certificates.DeleteCertificateStatusNotFound}, true},
		{certificates.DeleteCertificateResponse{Status: "invalidDeleteCertificateStatus"}, false},
		{certificates.DeleteCertificateResponse{}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV16TestSuite) TestDeleteCertificateE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	certificateHashData := types.CertificateHashData{HashAlgorithm: types.SHA256, IssuerNameHash: "hash00", IssuerKeyHash: "hash01", SerialNumber: "serial0"}
	status := certificates.DeleteCertificateStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"certificateHashData":{"hashAlgorithm":"%v","issuerNameHash":"%v","issuerKeyHash":"%v","serialNumber":"%v"}}]`,
		messageId, certificates.DeleteCertificateFeatureName, certificateHashData.HashAlgorithm, certificateHashData.IssuerNameHash, certificateHashData.IssuerKeyHash, certificateHashData.SerialNumber)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	deleteCertificateConfirmation := certificates.NewDeleteCertificateResponse(status)
	channel := NewMockWebSocket(wsId)

	handler := mocks.NewMockCertificatesChargePointHandler(t)
	handler.EXPECT().OnDeleteCertificate(mock.Anything).RunAndReturn(func(request *certificates.DeleteCertificateRequest) (*certificates.DeleteCertificateResponse, error) {
		assert.Equal(t, certificateHashData.HashAlgorithm, request.CertificateHashData.HashAlgorithm)
		assert.Equal(t, certificateHashData.IssuerNameHash, request.CertificateHashData.IssuerNameHash)
		assert.Equal(t, certificateHashData.IssuerKeyHash, request.CertificateHashData.IssuerKeyHash)
		assert.Equal(t, certificateHashData.SerialNumber, request.CertificateHashData.SerialNumber)
		return deleteCertificateConfirmation, nil
	})

	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	suite.chargePoint.SetCertificateHandler(handler)

	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	require.Nil(t, err)

	resultChannel := make(chan bool, 1)
	err = suite.centralSystem.DeleteCertificate(wsId, func(confirmation *certificates.DeleteCertificateResponse, err error) {
		require.Nil(t, err)
		require.NotNil(t, confirmation)
		assert.Equal(t, status, confirmation.Status)
		resultChannel <- true
	}, certificateHashData)
	require.Nil(t, err)

	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV16TestSuite) TestDeleteCertificateInvalidEndpoint() {
	messageId := defaultMessageId
	certificateHashData := types.CertificateHashData{HashAlgorithm: types.SHA256, IssuerNameHash: "hash00", IssuerKeyHash: "hash01", SerialNumber: "serial0"}
	deleteCertificateRequest := certificates.NewDeleteCertificateRequest(certificateHashData)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"certificateHashData":{"hashAlgorithm":"%v","issuerNameHash":"%v","issuerKeyHash":"%v","serialNumber":"%v"}}]`,
		messageId, certificates.DeleteCertificateFeatureName, certificateHashData.HashAlgorithm, certificateHashData.IssuerNameHash, certificateHashData.IssuerKeyHash, certificateHashData.SerialNumber)
	testUnsupportedRequestFromChargePoint(suite, deleteCertificateRequest, requestJson, messageId)
}
