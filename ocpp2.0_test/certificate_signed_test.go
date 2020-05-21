package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/security"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func (suite *OcppV2TestSuite) TestCertificateSignedRequestValidation() {
	t := suite.T()
	var testTable = []GenericTestEntry{
		{security.CertificateSignedRequest{Cert: []string{"sampleCert"}, TypeOfCertificate: types.ChargingStationCert}, true},
		{security.CertificateSignedRequest{Cert: []string{"sampleCert"}}, true},
		{security.CertificateSignedRequest{Cert: []string{}}, false},
		{security.CertificateSignedRequest{}, false},
		{security.CertificateSignedRequest{Cert: []string{">800............................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................."}}, false},
		{security.CertificateSignedRequest{Cert: []string{"sampleCert"}, TypeOfCertificate: "invalidCertificateType"}, false},
	}
	ExecuteGenericTestTable(t, testTable)
}

func (suite *OcppV2TestSuite) TestCertificateSignedConfirmationValidation() {
	t := suite.T()
	var testTable = []GenericTestEntry{
		{security.CertificateSignedResponse{Status: security.CertificateSignedStatusAccepted}, true},
		{security.CertificateSignedResponse{Status: security.CertificateSignedStatusRejected}, true},
		{security.CertificateSignedResponse{Status: "invalidCertificateSignedStatus"}, false},
		{security.CertificateSignedResponse{}, false},
	}
	ExecuteGenericTestTable(t, testTable)
}

// Test
func (suite *OcppV2TestSuite) TestCertificateSignedE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	certificate := "someX509Certificate"
	certificateType := types.ChargingStationCert
	status := security.CertificateSignedStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"cert":["%v"],"typeOfCertificate":"%v"}]`, messageId, security.CertificateSignedFeatureName, certificate, certificateType)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	certificateSignedConfirmation := security.NewCertificateSignedResponse(status)
	channel := NewMockWebSocket(wsId)
	// Setting handlers
	handler := MockChargingStationSecurityHandler{}
	handler.On("OnCertificateSigned", mock.Anything).Return(certificateSignedConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*security.CertificateSignedRequest)
		require.True(t, ok)
		require.Len(t, request.Cert, 1)
		assert.Equal(t, certificate, request.Cert[0])
		assert.Equal(t, certificateType, request.TypeOfCertificate)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.CertificateSigned(wsId, func(confirmation *security.CertificateSignedResponse, err error) {
		require.Nil(t, err)
		require.NotNil(t, confirmation)
		assert.Equal(t, status, confirmation.Status)
		resultChannel <- true
	}, []string{certificate}, func(request *security.CertificateSignedRequest) {
		request.TypeOfCertificate = certificateType
	})
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV2TestSuite) TestCertificateSignedInvalidEndpoint() {
	messageId := defaultMessageId
	certificate := "someX509Certificate"
	certificateType := types.ChargingStationCert
	certificateSignedRequest := security.NewCertificateSignedRequest([]string{certificate})
	certificateSignedRequest.TypeOfCertificate = certificateType
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"cert":["%v"],"typeOfCertificate":"%v"}]`, messageId, security.CertificateSignedFeatureName, certificate, certificateType)
	testUnsupportedRequestFromChargingStation(suite, certificateSignedRequest, requestJson, messageId)
}
