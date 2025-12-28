package ocpp2_test

import (
	"fmt"

	"github.com/stretchr/testify/mock"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/iso15118"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

func (suite *OcppV2TestSuite) TestInstallCertificateRequestValidation() {
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
		{iso15118.InstallCertificateRequest{CertificateType: types.V2GRootCertificate, Certificate: newLongString(5501)}, false},
	}
	ExecuteGenericTestTable(suite, testTable)
}

func (suite *OcppV2TestSuite) TestInstallCertificateConfirmationValidation() {
	var testTable = []GenericTestEntry{
		{iso15118.InstallCertificateResponse{Status: iso15118.CertificateStatusAccepted}, true},
		{iso15118.InstallCertificateResponse{Status: iso15118.CertificateStatusRejected}, true},
		{iso15118.InstallCertificateResponse{Status: iso15118.CertificateStatusFailed}, true},
		{iso15118.InstallCertificateResponse{}, false},
		{iso15118.InstallCertificateResponse{Status: "invalidInstallCertificateStatus"}, false},
	}
	ExecuteGenericTestTable(suite, testTable)
}

// Test
func (suite *OcppV2TestSuite) TestInstallCertificateE2EMocked() {
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	certificateType := types.CSMSRootCertificate
	status := iso15118.CertificateStatusAccepted
	certificate := "0xdeadbeef"
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"certificateType":"%v","certificate":"%v"}]`, messageId, iso15118.InstallCertificateFeatureName, certificateType, certificate)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	installCertificateResponse := iso15118.NewInstallCertificateResponse(status)
	channel := NewMockWebSocket(wsId)
	// Setting handlers
	handler := &MockChargingStationIso15118Handler{}
	handler.On("OnInstallCertificate", mock.Anything).Return(installCertificateResponse, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*iso15118.InstallCertificateRequest)
		suite.Require().True(ok)
		suite.Require().NotNil(request)
		suite.Equal(certificateType, request.CertificateType)
		suite.Equal(certificate, request.Certificate)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	suite.Require().Nil(err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.InstallCertificate(wsId, func(response *iso15118.InstallCertificateResponse, err error) {
		suite.Require().Nil(err)
		suite.Require().NotNil(response)
		suite.Equal(status, response.Status)
		resultChannel <- true
	}, certificateType, certificate)
	suite.Require().Nil(err)
	result := <-resultChannel
	suite.True(result)
}

func (suite *OcppV2TestSuite) TestInstallCertificateInvalidEndpoint() {
	messageId := defaultMessageId
	certificateType := types.CSMSRootCertificate
	certificate := "0xdeadbeef"
	installCertificateRequest := iso15118.NewInstallCertificateRequest(certificateType, certificate)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"certificateType":"%v","certificate":"%v"}]`, messageId, iso15118.InstallCertificateFeatureName, certificateType, certificate)
	testUnsupportedRequestFromChargingStation(suite, installCertificateRequest, requestJson, messageId)
}
