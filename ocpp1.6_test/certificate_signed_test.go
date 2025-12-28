package ocpp16_test

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/security"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6_test/mocks"
	"github.com/stretchr/testify/mock"
)

func (suite *OcppV16TestSuite) TestCertificateSignedRequestValidation() {
	var testTable = []GenericTestEntry{
		{security.CertificateSignedRequest{CertificateChain: "sampleCert"}, true},
		{security.CertificateSignedRequest{CertificateChain: ""}, false},
		{security.CertificateSignedRequest{}, false},
		{security.CertificateSignedRequest{CertificateChain: newLongString(100001)}, false},
	}
	ExecuteGenericTestTable(suite, testTable)
}

func (suite *OcppV16TestSuite) TestCertificateSignedConfirmationValidation() {
	var testTable = []GenericTestEntry{
		{security.CertificateSignedResponse{Status: security.CertificateSignedStatusAccepted}, true},
		{security.CertificateSignedResponse{Status: security.CertificateSignedStatusAccepted}, true},
		{security.CertificateSignedResponse{Status: security.CertificateSignedStatusRejected}, true},
		{security.CertificateSignedResponse{Status: "invalidCertificateSignedStatus"}, false},
		{security.CertificateSignedResponse{}, false},
	}
	ExecuteGenericTestTable(suite, testTable)
}

// Test
func (suite *OcppV16TestSuite) TestCertificateSignedE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	certificateChain := "someX509CertificateChain"
	status := security.CertificateSignedStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"certificateChain":"%v"}]`, messageId, security.CertificateSignedFeatureName, certificateChain)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	certificateSignedConfirmation := security.NewCertificateSignedResponse(status)
	channel := NewMockWebSocket(wsId)

	// Setting handlers
	handler := mocks.NewMockSecurityChargePointHandler(t)
	handler.EXPECT().OnCertificateSigned(mock.Anything).RunAndReturn(func(request *security.CertificateSignedRequest) (*security.CertificateSignedResponse, error) {
		suite.Equal(certificateChain, request.CertificateChain)
		return certificateSignedConfirmation, nil
	})

	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})

	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	suite.chargePoint.SetSecurityHandler(handler)
	err := suite.chargePoint.Start(wsUrl)
	suite.Require().Nil(err)
	resultChannel := make(chan bool, 1)
	err = suite.centralSystem.CertificateSigned(wsId, func(confirmation *security.CertificateSignedResponse, err error) {
		suite.Require().Nil(err)
		suite.Require().NotNil(confirmation)
		suite.Equal(status, confirmation.Status)
		resultChannel <- true
	}, certificateChain, func(request *security.CertificateSignedRequest) {
		request.CertificateChain = certificateChain
	})
	suite.Require().Nil(err)
	result := <-resultChannel
	suite.True(result)
}

func (suite *OcppV16TestSuite) TestCertificateSignedInvalidEndpoint() {
	messageId := defaultMessageId
	certificateChain := "someX509CertificateChain"
	certificateSignedRequest := security.NewCertificateSignedRequest(certificateChain)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"certificateChain":"%v"}]`, messageId, security.CertificateSignedFeatureName, certificateChain)
	testUnsupportedRequestFromChargePoint(suite, certificateSignedRequest, requestJson, messageId)
}
