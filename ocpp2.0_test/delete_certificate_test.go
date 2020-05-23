package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/iso15118"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test
func (suite *OcppV2TestSuite) TestDeleteCertificateRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{iso15118.DeleteCertificateRequest{CertificateHashData: types.CertificateHashData{HashAlgorithm: types.SHA256, IssuerNameHash: "hash00", IssuerKeyHash: "hash01", SerialNumber: "serial0"}}, true},
		{iso15118.DeleteCertificateRequest{}, false},
		{iso15118.DeleteCertificateRequest{CertificateHashData: types.CertificateHashData{HashAlgorithm: "invalidHashAlgorithm", IssuerNameHash: "hash00", IssuerKeyHash: "hash01", SerialNumber: "serial0"}}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestDeleteCertificateConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{iso15118.DeleteCertificateResponse{Status: iso15118.DeleteCertificateStatusAccepted}, true},
		{iso15118.DeleteCertificateResponse{Status: iso15118.DeleteCertificateStatusFailed}, true},
		{iso15118.DeleteCertificateResponse{Status: iso15118.DeleteCertificateStatusNotFound}, true},
		{iso15118.DeleteCertificateResponse{Status: "invalidDeleteCertificateStatus"}, false},
		{iso15118.DeleteCertificateResponse{}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestDeleteCertificateE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	certificateHashData := types.CertificateHashData{HashAlgorithm: types.SHA256, IssuerNameHash: "hash00", IssuerKeyHash: "hash01", SerialNumber: "serial0"}
	status := iso15118.DeleteCertificateStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"certificateHashData":{"hashAlgorithm":"%v","issuerNameHash":"%v","issuerKeyHash":"%v","serialNumber":"%v"}}]`,
		messageId, iso15118.DeleteCertificateFeatureName, certificateHashData.HashAlgorithm, certificateHashData.IssuerNameHash, certificateHashData.IssuerKeyHash, certificateHashData.SerialNumber)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	deleteCertificateConfirmation := iso15118.NewDeleteCertificateResponse(status)
	channel := NewMockWebSocket(wsId)

	handler := MockChargingStationIso15118Handler{}
	handler.On("OnDeleteCertificate", mock.Anything).Return(deleteCertificateConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*iso15118.DeleteCertificateRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, certificateHashData.HashAlgorithm, request.CertificateHashData.HashAlgorithm)
		assert.Equal(t, certificateHashData.IssuerNameHash, request.CertificateHashData.IssuerNameHash)
		assert.Equal(t, certificateHashData.IssuerKeyHash, request.CertificateHashData.IssuerKeyHash)
		assert.Equal(t, certificateHashData.SerialNumber, request.CertificateHashData.SerialNumber)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.DeleteCertificate(wsId, func(confirmation *iso15118.DeleteCertificateResponse, err error) {
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
	certificateHashData := types.CertificateHashData{HashAlgorithm: types.SHA256, IssuerNameHash: "hash00", IssuerKeyHash: "hash01", SerialNumber: "serial0"}
	deleteCertificateRequest := iso15118.NewDeleteCertificateRequest(certificateHashData)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"certificateHashData":{"hashAlgorithm":"%v","issuerNameHash":"%v","issuerKeyHash":"%v","serialNumber":"%v"}}]`,
		messageId, iso15118.DeleteCertificateFeatureName, certificateHashData.HashAlgorithm, certificateHashData.IssuerNameHash, certificateHashData.IssuerKeyHash, certificateHashData.SerialNumber)
	testUnsupportedRequestFromChargingStation(suite, deleteCertificateRequest, requestJson, messageId)
}
