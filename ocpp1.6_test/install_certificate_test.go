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

func (suite *OcppV16TestSuite) TestInstallCertificateRequestValidation() {
	t := suite.T()
	var testTable = []GenericTestEntry{
		{certificates.InstallCertificateRequest{CertificateType: types.ManufacturerRootCertificate, Certificate: "0xdeadbeef"}, true},
		{certificates.InstallCertificateRequest{CertificateType: types.ManufacturerRootCertificate}, false},
		{certificates.InstallCertificateRequest{CertificateType: types.CentralSystemRootCertificate, Certificate: "0xdeadbeef"}, true},
		{certificates.InstallCertificateRequest{CertificateType: types.CentralSystemRootCertificate, Certificate: ""}, false},
		{certificates.InstallCertificateRequest{Certificate: "0xdeadbeef"}, false},
		{certificates.InstallCertificateRequest{}, false},
		{certificates.InstallCertificateRequest{CertificateType: "invalidCertificateUse", Certificate: "0xdeadbeef"}, false},
		{certificates.InstallCertificateRequest{CertificateType: types.ManufacturerRootCertificate, Certificate: newLongString(5501)}, false},
	}
	ExecuteGenericTestTable(t, testTable)
}

func (suite *OcppV16TestSuite) TestInstallCertificateConfirmationValidation() {
	t := suite.T()
	var testTable = []GenericTestEntry{
		{certificates.InstallCertificateResponse{Status: certificates.CertificateStatusAccepted}, true},
		{certificates.InstallCertificateResponse{Status: certificates.CertificateStatusRejected}, true},
		{certificates.InstallCertificateResponse{Status: certificates.CertificateStatusFailed}, true},
		{certificates.InstallCertificateResponse{}, false},
		{certificates.InstallCertificateResponse{Status: "invalidInstallCertificateStatus"}, false},
	}
	ExecuteGenericTestTable(t, testTable)
}

// Test
func (suite *OcppV16TestSuite) TestInstallCertificateE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	certificateType := types.CentralSystemRootCertificate
	status := certificates.CertificateStatusAccepted
	certificate := "0xdeadbeef"
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"certificateType":"%v","certificate":"%v"}]`, messageId, certificates.InstallCertificateFeatureName, certificateType, certificate)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	installCertificateResponse := certificates.NewInstallCertificateResponse(status)
	channel := NewMockWebSocket(wsId)

	// Setting handlers
	handler := mocks.NewMockCertificatesChargePointHandler(t)
	handler.EXPECT().OnInstallCertificate(mock.Anything).RunAndReturn(func(request *certificates.InstallCertificateRequest) (*certificates.InstallCertificateResponse, error) {
		assert.Equal(t, certificateType, request.CertificateType)
		assert.Equal(t, certificate, request.Certificate)
		return installCertificateResponse, nil
	})

	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	suite.chargePoint.SetCertificateHandler(handler)

	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	suite.chargePoint.SetCertificateHandler(handler)
	err := suite.chargePoint.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.centralSystem.InstallCertificate(wsId, func(response *certificates.InstallCertificateResponse, err error) {
		require.Nil(t, err)
		require.NotNil(t, response)
		assert.Equal(t, status, response.Status)
		resultChannel <- true
	}, certificateType, certificate)
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV16TestSuite) TestInstallCertificateInvalidEndpoint() {
	messageId := defaultMessageId
	certificateType := types.CentralSystemRootCertificate
	certificate := "0xdeadbeef"
	installCertificateRequest := certificates.NewInstallCertificateRequest(certificateType, certificate)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"certificateType":"%v","certificate":"%v"}]`, messageId, certificates.InstallCertificateFeatureName, certificateType, certificate)
	testUnsupportedRequestFromChargePoint(suite, installCertificateRequest, requestJson, messageId)
}
