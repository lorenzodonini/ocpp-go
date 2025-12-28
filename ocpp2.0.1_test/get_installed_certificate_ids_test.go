package ocpp2_test

import (
	"fmt"

	"github.com/stretchr/testify/mock"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/iso15118"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

func (suite *OcppV2TestSuite) TestGetInstalledCertificateIdsRequestValidation() {
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
	ExecuteGenericTestTable(suite, testTable)
}

func (suite *OcppV2TestSuite) TestGetInstalledCertificateIdsConfirmationValidation() {
	var testTable = []GenericTestEntry{
		{iso15118.GetInstalledCertificateIdsResponse{Status: iso15118.GetInstalledCertificateStatusAccepted, CertificateHashDataChain: []types.CertificateHashDataChain{{CertificateType: types.CSMSRootCertificate, CertificateHashData: types.CertificateHashData{HashAlgorithm: types.SHA256, IssuerNameHash: "name0", IssuerKeyHash: "key0", SerialNumber: "serial0"}}}}, true},
		{iso15118.GetInstalledCertificateIdsResponse{Status: iso15118.GetInstalledCertificateStatusNotFound, CertificateHashDataChain: []types.CertificateHashDataChain{{CertificateType: types.CSMSRootCertificate, CertificateHashData: types.CertificateHashData{HashAlgorithm: types.SHA256, IssuerNameHash: "name0", IssuerKeyHash: "key0", SerialNumber: "serial0"}}}}, true},
		{iso15118.GetInstalledCertificateIdsResponse{Status: iso15118.GetInstalledCertificateStatusAccepted, CertificateHashDataChain: []types.CertificateHashDataChain{}}, true},
		{iso15118.GetInstalledCertificateIdsResponse{Status: iso15118.GetInstalledCertificateStatusAccepted}, true},
		{iso15118.GetInstalledCertificateIdsResponse{}, false},
		{iso15118.GetInstalledCertificateIdsResponse{Status: "invalidGetInstalledCertificateStatus"}, false},
		{iso15118.GetInstalledCertificateIdsResponse{Status: iso15118.GetInstalledCertificateStatusAccepted, CertificateHashDataChain: []types.CertificateHashDataChain{{CertificateType: types.CSMSRootCertificate, CertificateHashData: types.CertificateHashData{HashAlgorithm: "invalidHashAlgorithm", IssuerNameHash: "name0", IssuerKeyHash: "key0", SerialNumber: "serial0"}}}}, false},
	}
	ExecuteGenericTestTable(suite, testTable)
}

// Test
func (suite *OcppV2TestSuite) TestGetInstalledCertificateIdsE2EMocked() {
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
		suite.Require().True(ok)
		suite.Require().NotNil(request)
		suite.Equal(certificateTypes, request.CertificateTypes)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	suite.Require().Nil(err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.GetInstalledCertificateIds(wsId, func(confirmation *iso15118.GetInstalledCertificateIdsResponse, err error) {
		suite.Require().Nil(err)
		suite.Require().NotNil(confirmation)
		suite.Equal(status, confirmation.Status)
		suite.Require().Len(confirmation.CertificateHashDataChain, len(certificateHashDataChain))
		suite.Equal(certificateHashDataChain[0].CertificateHashData.HashAlgorithm, confirmation.CertificateHashDataChain[0].CertificateHashData.HashAlgorithm)
		suite.Equal(certificateHashDataChain[0].CertificateHashData.IssuerNameHash, confirmation.CertificateHashDataChain[0].CertificateHashData.IssuerNameHash)
		suite.Equal(certificateHashDataChain[0].CertificateHashData.IssuerKeyHash, confirmation.CertificateHashDataChain[0].CertificateHashData.IssuerKeyHash)
		suite.Equal(certificateHashDataChain[0].CertificateHashData.SerialNumber, confirmation.CertificateHashDataChain[0].CertificateHashData.SerialNumber)
		resultChannel <- true
	}, func(request *iso15118.GetInstalledCertificateIdsRequest) {
		request.CertificateTypes = certificateTypes
	})
	suite.Require().Nil(err)
	result := <-resultChannel
	suite.True(result)
}

func (suite *OcppV2TestSuite) TestGetInstalledCertificateIdsInvalidEndpoint() {
	messageId := defaultMessageId
	certificateTypes := []types.CertificateUse{types.CSMSRootCertificate}
	GetInstalledCertificateIdsRequest := iso15118.NewGetInstalledCertificateIdsRequest()
	GetInstalledCertificateIdsRequest.CertificateTypes = certificateTypes
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"certificateType":["%v"]}]`, messageId, iso15118.GetInstalledCertificateIdsFeatureName, certificateTypes[0])
	testUnsupportedRequestFromChargingStation(suite, GetInstalledCertificateIdsRequest, requestJson, messageId)
}
