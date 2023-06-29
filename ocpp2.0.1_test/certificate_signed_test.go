package ocpp2_test

import (
	"fmt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/security"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

func (suite *OcppV2TestSuite) TestCertificateSignedRequestValidation() {
	t := suite.T()
	var testTable = []GenericTestEntry{
		{security.CertificateSignedRequest{CertificateChain: "sampleCert", TypeOfCertificate: types.ChargingStationCert}, true},
		{security.CertificateSignedRequest{CertificateChain: "sampleCert"}, true},
		{security.CertificateSignedRequest{CertificateChain: ""}, false},
		{security.CertificateSignedRequest{}, false},
		{security.CertificateSignedRequest{CertificateChain: newLongString(100001)}, false},
		{security.CertificateSignedRequest{CertificateChain: "sampleCert", TypeOfCertificate: "invalidCertificateType"}, false},
	}
	ExecuteGenericTestTable(t, testTable)
}

func (suite *OcppV2TestSuite) TestCertificateSignedConfirmationValidation() {
	t := suite.T()
	var testTable = []GenericTestEntry{
		{security.CertificateSignedResponse{Status: security.CertificateSignedStatusAccepted, StatusInfo: types.NewStatusInfo("200", "ok")}, true},
		{security.CertificateSignedResponse{Status: security.CertificateSignedStatusAccepted}, true},
		{security.CertificateSignedResponse{Status: security.CertificateSignedStatusRejected}, true},
		{security.CertificateSignedResponse{Status: "invalidCertificateSignedStatus"}, false},
		{security.CertificateSignedResponse{Status: security.CertificateSignedStatusAccepted, StatusInfo: types.NewStatusInfo("", "")}, false},
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
	certificateChain := "someX509CertificateChain"
	certificateType := types.ChargingStationCert
	status := security.CertificateSignedStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"certificateChain":"%v","certificateType":"%v"}]`,
		messageId, security.CertificateSignedFeatureName, certificateChain, certificateType)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	certificateSignedConfirmation := security.NewCertificateSignedResponse(status)
	channel := NewMockWebSocket(wsId)
	// Setting handlers
	handler := &MockChargingStationSecurityHandler{}
	handler.On("OnCertificateSigned", mock.Anything).Return(certificateSignedConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*security.CertificateSignedRequest)
		require.True(t, ok)
		assert.Equal(t, certificateChain, request.CertificateChain)
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
	}, certificateChain, func(request *security.CertificateSignedRequest) {
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
	certificateSignedRequest := security.NewCertificateSignedRequest(certificate)
	certificateSignedRequest.TypeOfCertificate = certificateType
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"certificateChain":"%v","certificateType":"%v"}]`, messageId, security.CertificateSignedFeatureName, certificate, certificateType)
	testUnsupportedRequestFromChargingStation(suite, certificateSignedRequest, requestJson, messageId)
}
