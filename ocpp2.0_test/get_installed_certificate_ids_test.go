package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func (suite *OcppV2TestSuite) TestGetInstalledCertificateIdsRequestValidation() {
	t := suite.T()
	var testTable = []GenericTestEntry{
		{ocpp2.GetInstalledCertificateIdsRequest{TypeOfCertificate: ocpp2.V2GRootCertificate}, true},
		{ocpp2.GetInstalledCertificateIdsRequest{TypeOfCertificate: ocpp2.MORootCertificate}, true},
		{ocpp2.GetInstalledCertificateIdsRequest{TypeOfCertificate: ocpp2.CSOSubCA1}, true},
		{ocpp2.GetInstalledCertificateIdsRequest{TypeOfCertificate: ocpp2.CSOSubCA2}, true},
		{ocpp2.GetInstalledCertificateIdsRequest{TypeOfCertificate: ocpp2.CSMSRootCertificate}, true},
		{ocpp2.GetInstalledCertificateIdsRequest{TypeOfCertificate: ocpp2.ManufacturerRootCertificate}, true},
		{ocpp2.GetInstalledCertificateIdsRequest{}, false},
		{ocpp2.GetInstalledCertificateIdsRequest{TypeOfCertificate: "invalidCertificateUse"}, false},
	}
	ExecuteGenericTestTable(t, testTable)
}

func (suite *OcppV2TestSuite) TestGetInstalledCertificateIdsConfirmationValidation() {
	t := suite.T()
	var testTable = []GenericTestEntry{
		{ocpp2.GetInstalledCertificateIdsConfirmation{Status: ocpp2.GetInstalledCertificateStatusAccepted, CertificateHashData: []ocpp2.CertificateHashData{{HashAlgorithm: ocpp2.SHA256, IssuerNameHash: "name0", IssuerKeyHash: "key0", SerialNumber: "serial0"}}}, true},
		{ocpp2.GetInstalledCertificateIdsConfirmation{Status: ocpp2.GetInstalledCertificateStatusNotFound, CertificateHashData: []ocpp2.CertificateHashData{{HashAlgorithm: ocpp2.SHA256, IssuerNameHash: "name0", IssuerKeyHash: "key0", SerialNumber: "serial0"}}}, true},
		{ocpp2.GetInstalledCertificateIdsConfirmation{Status: ocpp2.GetInstalledCertificateStatusAccepted, CertificateHashData: []ocpp2.CertificateHashData{}}, true},
		{ocpp2.GetInstalledCertificateIdsConfirmation{Status: ocpp2.GetInstalledCertificateStatusAccepted}, true},
		{ocpp2.GetInstalledCertificateIdsConfirmation{}, false},
		{ocpp2.GetInstalledCertificateIdsConfirmation{Status: "invalidGetInstalledCertificateStatus"}, false},
		{ocpp2.GetInstalledCertificateIdsConfirmation{Status: ocpp2.GetInstalledCertificateStatusAccepted, CertificateHashData: []ocpp2.CertificateHashData{{HashAlgorithm: "invalidHashAlgorithm", IssuerNameHash: "name0", IssuerKeyHash: "key0", SerialNumber: "serial0"}}}, false},
	}
	ExecuteGenericTestTable(t, testTable)
}

// Test
func (suite *OcppV2TestSuite) TestGetInstalledCertificateIdsE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	certificateType := ocpp2.CSMSRootCertificate
	status := ocpp2.GetInstalledCertificateStatusAccepted
	certificateHashData := []ocpp2.CertificateHashData{
		{ HashAlgorithm: ocpp2.SHA256, IssuerNameHash: "name0", IssuerKeyHash: "key0", SerialNumber: "serial0" },
	}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"typeOfCertificate":"%v"}]`, messageId, ocpp2.GetInstalledCertificateIdsFeatureName, certificateType)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v","certificateHashData":[{"hashAlgorithm":"%v","issuerNameHash":"%v","issuerKeyHash":"%v","serialNumber":"%v"}]}]`,
		messageId, status, certificateHashData[0].HashAlgorithm, certificateHashData[0].IssuerNameHash, certificateHashData[0].IssuerKeyHash, certificateHashData[0].SerialNumber)
	getInstalledCertificateIdsConfirmation := ocpp2.NewGetInstalledCertificateIdsConfirmation(status)
	getInstalledCertificateIdsConfirmation.CertificateHashData = certificateHashData
	channel := NewMockWebSocket(wsId)
	// Setting handlers
	coreListener := MockChargePointCoreListener{}
	coreListener.On("OnGetInstalledCertificateIds", mock.Anything).Return(getInstalledCertificateIdsConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*ocpp2.GetInstalledCertificateIdsRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, certificateType, request.TypeOfCertificate)
	})
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, coreListener, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.GetInstalledCertificateIds(wsId, func(confirmation *ocpp2.GetInstalledCertificateIdsConfirmation, err error) {
		require.Nil(t, err)
		require.NotNil(t, confirmation)
		assert.Equal(t, status, confirmation.Status)
		require.Len(t, confirmation.CertificateHashData, len(certificateHashData))
		assert.Equal(t, certificateHashData[0].HashAlgorithm, confirmation.CertificateHashData[0].HashAlgorithm)
		assert.Equal(t, certificateHashData[0].IssuerNameHash, confirmation.CertificateHashData[0].IssuerNameHash)
		assert.Equal(t, certificateHashData[0].IssuerKeyHash, confirmation.CertificateHashData[0].IssuerKeyHash)
		assert.Equal(t, certificateHashData[0].SerialNumber, confirmation.CertificateHashData[0].SerialNumber)
		resultChannel <- true
	}, certificateType)
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV2TestSuite) TestGetInstalledCertificateIdsInvalidEndpoint() {
	messageId := defaultMessageId
	certificateType := ocpp2.CSMSRootCertificate
	GetInstalledCertificateIdsRequest := ocpp2.NewGetInstalledCertificateIdsRequest(certificateType)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"typeOfCertificate":"%v"}]`, messageId, ocpp2.GetInstalledCertificateIdsFeatureName, certificateType)
	testUnsupportedRequestFromChargePoint(suite, GetInstalledCertificateIdsRequest, requestJson, messageId)
}
