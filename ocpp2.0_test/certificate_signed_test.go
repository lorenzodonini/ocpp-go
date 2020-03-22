package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func (suite *OcppV2TestSuite) TestCertificateSignedRequestValidation() {
	t := suite.T()
	var testTable = []GenericTestEntry{
		{ocpp2.CertificateSignedRequest{Cert: []string{"sampleCert"}, TypeOfCertificate: ocpp2.ChargingStationCert}, true},
		{ocpp2.CertificateSignedRequest{Cert: []string{"sampleCert"}}, true},
		{ocpp2.CertificateSignedRequest{Cert: []string{}}, false},
		{ocpp2.CertificateSignedRequest{}, false},
		{ocpp2.CertificateSignedRequest{Cert: []string{">800............................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................."}}, false},
		{ocpp2.CertificateSignedRequest{Cert: []string{"sampleCert"}, TypeOfCertificate: "invalidCertificateType"}, false},
	}
	ExecuteGenericTestTable(t, testTable)
}

func (suite *OcppV2TestSuite) TestCertificateSignedConfirmationValidation() {
	t := suite.T()
	var testTable = []GenericTestEntry{
		{ocpp2.CertificateSignedConfirmation{Status: ocpp2.CertificateSignedStatusAccepted}, true},
		{ocpp2.CertificateSignedConfirmation{Status: ocpp2.CertificateSignedStatusRejected}, true},
		{ocpp2.CertificateSignedConfirmation{Status: "invalidCertificateSignedStatus"}, false},
		{ocpp2.CertificateSignedConfirmation{}, false},
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
	certificateType := ocpp2.ChargingStationCert
	status := ocpp2.CertificateSignedStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"cert":["%v"],"typeOfCertificate":"%v"}]`, messageId, ocpp2.CertificateSignedFeatureName, certificate, certificateType)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	certificateSignedConfirmation := ocpp2.NewCertificateSignedConfirmation(status)
	channel := NewMockWebSocket(wsId)
	// Setting handlers
	coreListener := MockChargePointCoreListener{}
	coreListener.On("OnCertificateSigned", mock.Anything).Return(certificateSignedConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*ocpp2.CertificateSignedRequest)
		require.True(t, ok)
		require.Len(t, request.Cert, 1)
		assert.Equal(t, certificate, request.Cert[0])
		assert.Equal(t, certificateType, request.TypeOfCertificate)
	})
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, coreListener, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.CertificateSigned(wsId, func(confirmation *ocpp2.CertificateSignedConfirmation, err error) {
		require.Nil(t, err)
		require.NotNil(t, confirmation)
		assert.Equal(t, status, confirmation.Status)
		resultChannel <- true
	}, []string{certificate}, func(request *ocpp2.CertificateSignedRequest) {
		request.TypeOfCertificate = certificateType
	})
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV2TestSuite) TestCertificateSignedInvalidEndpoint() {
	messageId := defaultMessageId
	certificate := "someX509Certificate"
	certificateType := ocpp2.ChargingStationCert
	certificateSignedRequest := ocpp2.NewCertificateSignedRequest([]string{certificate})
	certificateSignedRequest.TypeOfCertificate = certificateType
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"cert":["%v"],"typeOfCertificate":"%v"}]`, messageId, ocpp2.CertificateSignedFeatureName, certificate, certificateType)
	testUnsupportedRequestFromChargePoint(suite, certificateSignedRequest, requestJson, messageId)
}
