package ocpp2_test

import (
	"fmt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/iso15118"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

func (suite *OcppV2TestSuite) TestGetInstalledCertificateIdsRequestValidation() {
	t := suite.T()
	var testTable = []GenericTestEntry{
		{iso15118.GetInstalledCertificateIdsRequest{CertificateTypes: []types.CertificateUse{types.V2GRootCertificate}}, true},
		{iso15118.GetInstalledCertificateIdsRequest{CertificateTypes: []types.CertificateUse{types.MORootCertificate}}, true},
		{iso15118.GetInstalledCertificateIdsRequest{CertificateTypes: []types.CertificateUse{types.CSOSubCA1}}, true},
		{iso15118.GetInstalledCertificateIdsRequest{CertificateTypes: []types.CertificateUse{types.CSOSubCA2}}, true},
		{iso15118.GetInstalledCertificateIdsRequest{CertificateTypes: []types.CertificateUse{types.CSMSRootCertificate}}, true},
		{iso15118.GetInstalledCertificateIdsRequest{CertificateTypes: []types.CertificateUse{types.ManufacturerRootCertificate}}, true},
		{iso15118.GetInstalledCertificateIdsRequest{}, true},
		{iso15118.GetInstalledCertificateIdsRequest{CertificateTypes: []types.CertificateUse{"invalidCertificateUse"}}, false},
	}
	ExecuteGenericTestTable(t, testTable)
}

func (suite *OcppV2TestSuite) TestGetInstalledCertificateIdsConfirmationValidation() {
	t := suite.T()
	var testTable = []GenericTestEntry{
		{iso15118.GetInstalledCertificateIdsResponse{Status: iso15118.GetInstalledCertificateStatusAccepted, CertificateHashDataChain: []types.CertificateHashDataChain{{CertificateType: types.CSMSRootCertificate, CertificateHashData: types.CertificateHashData{HashAlgorithm: types.SHA256, IssuerNameHash: "name0", IssuerKeyHash: "key0", SerialNumber: "serial0"}}}}, true},
		{iso15118.GetInstalledCertificateIdsResponse{Status: iso15118.GetInstalledCertificateStatusNotFound, CertificateHashDataChain: []types.CertificateHashDataChain{{CertificateType: types.CSMSRootCertificate, CertificateHashData: types.CertificateHashData{HashAlgorithm: types.SHA256, IssuerNameHash: "name0", IssuerKeyHash: "key0", SerialNumber: "serial0"}}}}, true},
		{iso15118.GetInstalledCertificateIdsResponse{Status: iso15118.GetInstalledCertificateStatusAccepted, CertificateHashDataChain: []types.CertificateHashDataChain{}}, true},
		{iso15118.GetInstalledCertificateIdsResponse{Status: iso15118.GetInstalledCertificateStatusAccepted}, true},
		{iso15118.GetInstalledCertificateIdsResponse{}, false},
		{iso15118.GetInstalledCertificateIdsResponse{Status: "invalidGetInstalledCertificateStatus"}, false},
		{iso15118.GetInstalledCertificateIdsResponse{Status: iso15118.GetInstalledCertificateStatusAccepted, CertificateHashDataChain: []types.CertificateHashDataChain{{CertificateType: types.CSMSRootCertificate, CertificateHashData: types.CertificateHashData{HashAlgorithm: "invalidHashAlgorithm", IssuerNameHash: "name0", IssuerKeyHash: "key0", SerialNumber: "serial0"}}}}, false},
	}
	ExecuteGenericTestTable(t, testTable)
}

// Test
func (suite *OcppV2TestSuite) TestGetInstalledCertificateIdsE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	certificateTypes := []types.CertificateUse{types.CSMSRootCertificate}
	status := iso15118.GetInstalledCertificateStatusAccepted
	certificateHashDataChain := []types.CertificateHashDataChain{
		{CertificateType: types.CSMSRootCertificate, CertificateHashData: types.CertificateHashData{HashAlgorithm: types.SHA256, IssuerNameHash: "name0", IssuerKeyHash: "key0", SerialNumber: "serial0"}},
	}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"certificateType":["%v"]}]`, messageId, iso15118.GetInstalledCertificateIdsFeatureName, certificateTypes[0])
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v","certificateHashDataChain":[{"certificateType":"%v","certificateHashData":{"hashAlgorithm":"%v","issuerNameHash":"%v","issuerKeyHash":"%v","serialNumber":"%v"}}]}]`,
		messageId, status, certificateHashDataChain[0].CertificateType, certificateHashDataChain[0].CertificateHashData.HashAlgorithm, certificateHashDataChain[0].CertificateHashData.IssuerNameHash, certificateHashDataChain[0].CertificateHashData.IssuerKeyHash, certificateHashDataChain[0].CertificateHashData.SerialNumber)
	getInstalledCertificateIdsConfirmation := iso15118.NewGetInstalledCertificateIdsResponse(status)
	getInstalledCertificateIdsConfirmation.CertificateHashDataChain = certificateHashDataChain
	channel := NewMockWebSocket(wsId)
	// Setting handlers
	handler := &MockChargingStationIso15118Handler{}
	handler.On("OnGetInstalledCertificateIds", mock.Anything).Return(getInstalledCertificateIdsConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*iso15118.GetInstalledCertificateIdsRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, certificateTypes, request.CertificateTypes)
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
		require.Len(t, confirmation.CertificateHashDataChain, len(certificateHashDataChain))
		assert.Equal(t, certificateHashDataChain[0].CertificateHashData.HashAlgorithm, confirmation.CertificateHashDataChain[0].CertificateHashData.HashAlgorithm)
		assert.Equal(t, certificateHashDataChain[0].CertificateHashData.IssuerNameHash, confirmation.CertificateHashDataChain[0].CertificateHashData.IssuerNameHash)
		assert.Equal(t, certificateHashDataChain[0].CertificateHashData.IssuerKeyHash, confirmation.CertificateHashDataChain[0].CertificateHashData.IssuerKeyHash)
		assert.Equal(t, certificateHashDataChain[0].CertificateHashData.SerialNumber, confirmation.CertificateHashDataChain[0].CertificateHashData.SerialNumber)
		resultChannel <- true
	}, func(request *iso15118.GetInstalledCertificateIdsRequest) {
		request.CertificateTypes = certificateTypes
	})
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV2TestSuite) TestGetInstalledCertificateIdsInvalidEndpoint() {
	messageId := defaultMessageId
	certificateTypes := []types.CertificateUse{types.CSMSRootCertificate}
	GetInstalledCertificateIdsRequest := iso15118.NewGetInstalledCertificateIdsRequest()
	GetInstalledCertificateIdsRequest.CertificateTypes = certificateTypes
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"certificateType":["%v"]}]`, messageId, iso15118.GetInstalledCertificateIdsFeatureName, certificateTypes[0])
	testUnsupportedRequestFromChargingStation(suite, GetInstalledCertificateIdsRequest, requestJson, messageId)
}
