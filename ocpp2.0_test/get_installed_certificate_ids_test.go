package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/iso15118"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func (suite *OcppV2TestSuite) TestGetInstalledCertificateIdsRequestValidation() {
	t := suite.T()
	var testTable = []GenericTestEntry{
		{iso15118.GetInstalledCertificateIdsRequest{TypeOfCertificate: types.V2GRootCertificate}, true},
		{iso15118.GetInstalledCertificateIdsRequest{TypeOfCertificate: types.MORootCertificate}, true},
		{iso15118.GetInstalledCertificateIdsRequest{TypeOfCertificate: types.CSOSubCA1}, true},
		{iso15118.GetInstalledCertificateIdsRequest{TypeOfCertificate: types.CSOSubCA2}, true},
		{iso15118.GetInstalledCertificateIdsRequest{TypeOfCertificate: types.CSMSRootCertificate}, true},
		{iso15118.GetInstalledCertificateIdsRequest{TypeOfCertificate: types.ManufacturerRootCertificate}, true},
		{iso15118.GetInstalledCertificateIdsRequest{}, false},
		{iso15118.GetInstalledCertificateIdsRequest{TypeOfCertificate: "invalidCertificateUse"}, false},
	}
	ExecuteGenericTestTable(t, testTable)
}

func (suite *OcppV2TestSuite) TestGetInstalledCertificateIdsConfirmationValidation() {
	t := suite.T()
	var testTable = []GenericTestEntry{
		{iso15118.GetInstalledCertificateIdsResponse{Status: iso15118.GetInstalledCertificateStatusAccepted, CertificateHashData: []types.CertificateHashData{{HashAlgorithm: types.SHA256, IssuerNameHash: "name0", IssuerKeyHash: "key0", SerialNumber: "serial0"}}}, true},
		{iso15118.GetInstalledCertificateIdsResponse{Status: iso15118.GetInstalledCertificateStatusNotFound, CertificateHashData: []types.CertificateHashData{{HashAlgorithm: types.SHA256, IssuerNameHash: "name0", IssuerKeyHash: "key0", SerialNumber: "serial0"}}}, true},
		{iso15118.GetInstalledCertificateIdsResponse{Status: iso15118.GetInstalledCertificateStatusAccepted, CertificateHashData: []types.CertificateHashData{}}, true},
		{iso15118.GetInstalledCertificateIdsResponse{Status: iso15118.GetInstalledCertificateStatusAccepted}, true},
		{iso15118.GetInstalledCertificateIdsResponse{}, false},
		{iso15118.GetInstalledCertificateIdsResponse{Status: "invalidGetInstalledCertificateStatus"}, false},
		{iso15118.GetInstalledCertificateIdsResponse{Status: iso15118.GetInstalledCertificateStatusAccepted, CertificateHashData: []types.CertificateHashData{{HashAlgorithm: "invalidHashAlgorithm", IssuerNameHash: "name0", IssuerKeyHash: "key0", SerialNumber: "serial0"}}}, false},
	}
	ExecuteGenericTestTable(t, testTable)
}

// Test
func (suite *OcppV2TestSuite) TestGetInstalledCertificateIdsE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	certificateType := types.CSMSRootCertificate
	status := iso15118.GetInstalledCertificateStatusAccepted
	certificateHashData := []types.CertificateHashData{
		{HashAlgorithm: types.SHA256, IssuerNameHash: "name0", IssuerKeyHash: "key0", SerialNumber: "serial0"},
	}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"typeOfCertificate":"%v"}]`, messageId, iso15118.GetInstalledCertificateIdsFeatureName, certificateType)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v","certificateHashData":[{"hashAlgorithm":"%v","issuerNameHash":"%v","issuerKeyHash":"%v","serialNumber":"%v"}]}]`,
		messageId, status, certificateHashData[0].HashAlgorithm, certificateHashData[0].IssuerNameHash, certificateHashData[0].IssuerKeyHash, certificateHashData[0].SerialNumber)
	getInstalledCertificateIdsConfirmation := iso15118.NewGetInstalledCertificateIdsResponse(status)
	getInstalledCertificateIdsConfirmation.CertificateHashData = certificateHashData
	channel := NewMockWebSocket(wsId)
	// Setting handlers
	handler := MockChargingStationIso15118Handler{}
	handler.On("OnGetInstalledCertificateIds", mock.Anything).Return(getInstalledCertificateIdsConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*iso15118.GetInstalledCertificateIdsRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, certificateType, request.TypeOfCertificate)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.GetInstalledCertificateIds(wsId, func(confirmation *iso15118.GetInstalledCertificateIdsResponse, err error) {
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
	certificateType := types.CSMSRootCertificate
	GetInstalledCertificateIdsRequest := iso15118.NewGetInstalledCertificateIdsRequest(certificateType)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"typeOfCertificate":"%v"}]`, messageId, iso15118.GetInstalledCertificateIdsFeatureName, certificateType)
	testUnsupportedRequestFromChargingStation(suite, GetInstalledCertificateIdsRequest, requestJson, messageId)
}
