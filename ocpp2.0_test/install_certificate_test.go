package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/iso15118"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func (suite *OcppV2TestSuite) TestInstallCertificateRequestValidation() {
	t := suite.T()
	var testTable = []GenericTestEntry{
		{iso15118.InstallCertificateRequest{CertificateType: types.V2GRootCertificate, Certificate: "0xdeadbeef"}, true},
		{iso15118.InstallCertificateRequest{CertificateType: types.MORootCertificate, Certificate: "0xdeadbeef"}, true},
		{iso15118.InstallCertificateRequest{CertificateType: types.CSOSubCA1, Certificate: "0xdeadbeef"}, true},
		{iso15118.InstallCertificateRequest{CertificateType: types.CSOSubCA2, Certificate: "0xdeadbeef"}, true},
		{iso15118.InstallCertificateRequest{CertificateType: types.CSMSRootCertificate, Certificate: "0xdeadbeef"}, true},
		{iso15118.InstallCertificateRequest{CertificateType: types.ManufacturerRootCertificate, Certificate: "0xdeadbeef"}, true},
		{iso15118.InstallCertificateRequest{CertificateType: types.ManufacturerRootCertificate}, false},
		{iso15118.InstallCertificateRequest{Certificate: "0xdeadbeef"}, false},
		{iso15118.InstallCertificateRequest{}, false},
		{iso15118.InstallCertificateRequest{CertificateType: "invalidCertificateUse", Certificate: "0xdeadbeef"}, false},
		{iso15118.InstallCertificateRequest{CertificateType: types.V2GRootCertificate, Certificate: ">800............................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................."}, false},
	}
	ExecuteGenericTestTable(t, testTable)
}

func (suite *OcppV2TestSuite) TestInstallCertificateConfirmationValidation() {
	t := suite.T()
	var testTable = []GenericTestEntry{
		{iso15118.InstallCertificateResponse{Status: types.CertificateStatusAccepted}, true},
		{iso15118.InstallCertificateResponse{}, false},
		{iso15118.InstallCertificateResponse{Status: "invalidCertificateStatus"}, false},
	}
	ExecuteGenericTestTable(t, testTable)
}

// Test
func (suite *OcppV2TestSuite) TestInstallCertificateE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	certificateType := types.CSMSRootCertificate
	status := types.CertificateStatusAccepted
	certificate := "0xdeadbeef"
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"certificateType":"%v","certificate":"%v"}]`, messageId, iso15118.InstallCertificateFeatureName, certificateType, certificate)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	installCertificateResponse := iso15118.NewInstallCertificateResponse(status)
	channel := NewMockWebSocket(wsId)
	// Setting handlers
	handler := MockChargingStationIso15118Handler{}
	handler.On("OnInstallCertificate", mock.Anything).Return(installCertificateResponse, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*iso15118.InstallCertificateRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, certificateType, request.CertificateType)
		assert.Equal(t, certificate, request.Certificate)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.InstallCertificate(wsId, func(response *iso15118.InstallCertificateResponse, err error) {
		require.Nil(t, err)
		require.NotNil(t, response)
		assert.Equal(t, status, response.Status)
		resultChannel <- true
	}, certificateType, certificate)
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV2TestSuite) TestInstallCertificateInvalidEndpoint() {
	messageId := defaultMessageId
	certificateType := types.CSMSRootCertificate
	certificate := "0xdeadbeef"
	installCertificateRequest := iso15118.NewInstallCertificateRequest(certificateType, certificate)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"certificateType":"%v","certificate":"%v"}]`, messageId, iso15118.InstallCertificateFeatureName, certificateType, certificate)
	testUnsupportedRequestFromChargingStation(suite, installCertificateRequest, requestJson, messageId)
}
