package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test
func (suite *OcppV2TestSuite) TestDeleteCertificateRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{ocpp2.DeleteCertificateRequest{CertificateHashData: ocpp2.CertificateHashData{HashAlgorithm: ocpp2.SHA256, IssuerNameHash: "hash00", IssuerKeyHash: "hash01", SerialNumber: "serial0"}}, true},
		{ocpp2.DeleteCertificateRequest{}, false},
		{ocpp2.DeleteCertificateRequest{CertificateHashData: ocpp2.CertificateHashData{HashAlgorithm: "invalidHashAlgorithm", IssuerNameHash: "hash00", IssuerKeyHash: "hash01", SerialNumber: "serial0"}}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestDeleteCertificateConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{ocpp2.DeleteCertificateConfirmation{Status: ocpp2.DeleteCertificateStatusAccepted}, true},
		{ocpp2.DeleteCertificateConfirmation{Status: ocpp2.DeleteCertificateStatusFailed}, true},
		{ocpp2.DeleteCertificateConfirmation{Status: ocpp2.DeleteCertificateStatusNotFound}, true},
		{ocpp2.DeleteCertificateConfirmation{Status: "invalidDeleteCertificateStatus"}, false},
		{ocpp2.DeleteCertificateConfirmation{}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestDeleteCertificateE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	certificateHashData := ocpp2.CertificateHashData{HashAlgorithm: ocpp2.SHA256, IssuerNameHash: "hash00", IssuerKeyHash: "hash01", SerialNumber: "serial0"}
	status := ocpp2.DeleteCertificateStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"certificateHashData":{"hashAlgorithm":"%v","issuerNameHash":"%v","issuerKeyHash":"%v","serialNumber":"%v"}}]`,
		messageId, ocpp2.DeleteCertificateFeatureName, certificateHashData.HashAlgorithm, certificateHashData.IssuerNameHash, certificateHashData.IssuerKeyHash, certificateHashData.SerialNumber)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	deleteCertificateConfirmation := ocpp2.NewDeleteCertificateConfirmation(status)
	channel := NewMockWebSocket(wsId)

	coreListener := MockChargePointCoreListener{}
	coreListener.On("OnDeleteCertificate", mock.Anything).Return(deleteCertificateConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*ocpp2.DeleteCertificateRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, certificateHashData.HashAlgorithm, request.CertificateHashData.HashAlgorithm)
		assert.Equal(t, certificateHashData.IssuerNameHash, request.CertificateHashData.IssuerNameHash)
		assert.Equal(t, certificateHashData.IssuerKeyHash, request.CertificateHashData.IssuerKeyHash)
		assert.Equal(t, certificateHashData.SerialNumber, request.CertificateHashData.SerialNumber)
	})
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, coreListener, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.DeleteCertificate(wsId, func(confirmation *ocpp2.DeleteCertificateConfirmation, err error) {
		require.Nil(t, err)
		require.NotNil(t, confirmation)
		assert.Equal(t, status, confirmation.Status)
		resultChannel <- true
	}, certificateHashData)
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV2TestSuite) TestDeleteCertificateInvalidEndpoint() {
	messageId := defaultMessageId
	certificateHashData := ocpp2.CertificateHashData{HashAlgorithm: ocpp2.SHA256, IssuerNameHash: "hash00", IssuerKeyHash: "hash01", SerialNumber: "serial0"}
	deleteCertificateRequest := ocpp2.NewDeleteCertificateRequest(certificateHashData)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"certificateHashData":{"hashAlgorithm":"%v","issuerNameHash":"%v","issuerKeyHash":"%v","serialNumber":"%v"}}]`,
		messageId, ocpp2.DeleteCertificateFeatureName, certificateHashData.HashAlgorithm, certificateHashData.IssuerNameHash, certificateHashData.IssuerKeyHash, certificateHashData.SerialNumber)
	testUnsupportedRequestFromChargePoint(suite, deleteCertificateRequest, requestJson, messageId)
}
