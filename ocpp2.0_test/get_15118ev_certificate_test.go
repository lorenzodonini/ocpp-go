package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/iso15118"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test
func (suite *OcppV2TestSuite) TestGet15118EVCertificateRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{iso15118.Get15118EVCertificateRequest{SchemaVersion: "1.0", ExiRequest: "deadbeef"}, true},
		{iso15118.Get15118EVCertificateRequest{SchemaVersion: "1.0"}, false},
		{iso15118.Get15118EVCertificateRequest{ExiRequest: "deadbeef"}, false},
		{iso15118.Get15118EVCertificateRequest{}, false},
		{iso15118.Get15118EVCertificateRequest{SchemaVersion: ">50................................................", ExiRequest: "deadbeef"}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestGet15118EVCertificateConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{iso15118.Get15118EVCertificateResponse{Status: types.Certificate15188EVStatusAccepted, ExiResponse: "deadbeef", ContractSignatureCertificateChain: iso15118.CertificateChain{Certificate: "deadcode"}, SaProvisioningCertificateChain: iso15118.CertificateChain{Certificate: "deadcode2"}}, true},
		{iso15118.Get15118EVCertificateResponse{Status: types.Certificate15188EVStatusAccepted, ExiResponse: "deadbeef", ContractSignatureCertificateChain: iso15118.CertificateChain{Certificate: "deadcode", ChildCertificate: []string{"c1", "c2", "c3", "c4"}}, SaProvisioningCertificateChain: iso15118.CertificateChain{Certificate: "deadcode2"}}, true},
		{iso15118.Get15118EVCertificateResponse{Status: types.Certificate15188EVStatusAccepted, ExiResponse: "deadbeef", ContractSignatureCertificateChain: iso15118.CertificateChain{Certificate: "deadcode"}}, false},
		{iso15118.Get15118EVCertificateResponse{Status: types.Certificate15188EVStatusAccepted, ExiResponse: "deadbeef", SaProvisioningCertificateChain: iso15118.CertificateChain{Certificate: "deadcode2"}}, false},
		{iso15118.Get15118EVCertificateResponse{Status: types.Certificate15188EVStatusAccepted, ContractSignatureCertificateChain: iso15118.CertificateChain{Certificate: "deadcode"}, SaProvisioningCertificateChain: iso15118.CertificateChain{Certificate: "deadcode2"}}, false},
		{iso15118.Get15118EVCertificateResponse{ExiResponse: "deadbeef", ContractSignatureCertificateChain: iso15118.CertificateChain{Certificate: "deadcode"}, SaProvisioningCertificateChain: iso15118.CertificateChain{Certificate: "deadcode2"}}, false},
		{iso15118.Get15118EVCertificateResponse{}, false},
		{iso15118.Get15118EVCertificateResponse{Status: "invalidCertificateStatus", ExiResponse: "deadbeef", ContractSignatureCertificateChain: iso15118.CertificateChain{Certificate: "deadcode"}, SaProvisioningCertificateChain: iso15118.CertificateChain{Certificate: "deadcode2"}}, false},
		{iso15118.Get15118EVCertificateResponse{Status: types.Certificate15188EVStatusAccepted, ExiResponse: "deadbeef", ContractSignatureCertificateChain: iso15118.CertificateChain{Certificate: "deadcode"}, SaProvisioningCertificateChain: iso15118.CertificateChain{}}, false},
		{iso15118.Get15118EVCertificateResponse{Status: types.Certificate15188EVStatusAccepted, ExiResponse: "deadbeef", ContractSignatureCertificateChain: iso15118.CertificateChain{Certificate: "deadcode", ChildCertificate: []string{"c1", "c2", "c3", "c4", "c5"}}, SaProvisioningCertificateChain: iso15118.CertificateChain{Certificate: "deadcode2"}}, false},
		{iso15118.Get15118EVCertificateResponse{Status: types.Certificate15188EVStatusAccepted, ExiResponse: "deadbeef", ContractSignatureCertificateChain: iso15118.CertificateChain{Certificate: "deadcode", ChildCertificate: []string{"c1", "c2", "c3", ""}}, SaProvisioningCertificateChain: iso15118.CertificateChain{Certificate: "deadcode2"}}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestGet15118EVCertificateE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	status := types.Certificate15188EVStatusAccepted
	schemaVersion := "1.0"
	exiRequest := "deadbeef"
	exiResponse := "deadbeef2"
	contractSignatureCertificateChain := iso15118.CertificateChain{Certificate: "deadcode"}
	saProvisioningCertificateChain := iso15118.CertificateChain{Certificate: "deadcode2"}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"15118SchemaVersion":"%v","exiRequest":"%v"}]`, messageId, iso15118.Get15118EVCertificateFeatureName, schemaVersion, exiRequest)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v","exiResponse":"%v","contractSignatureCertificateChain":{"certificate":"%v"},"saProvisioningCertificateChain":{"certificate":"%v"}}]`,
		messageId, status, exiResponse, contractSignatureCertificateChain.Certificate, saProvisioningCertificateChain.Certificate)
	get15118EVCertificateConfirmation := iso15118.NewGet15118EVCertificateResponse(status, exiResponse, contractSignatureCertificateChain, saProvisioningCertificateChain)
	channel := NewMockWebSocket(wsId)

	handler := MockCSMSIso15118Handler{}
	handler.On("OnGet15118EVCertificate", mock.AnythingOfType("string"), mock.Anything).Return(get15118EVCertificateConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*iso15118.Get15118EVCertificateRequest)
		require.True(t, ok)
		assert.Equal(t, schemaVersion, request.SchemaVersion)
		assert.Equal(t, exiRequest, request.ExiRequest)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	confirmation, err := suite.chargingStation.Get15118EVCertificate(schemaVersion, exiRequest)
	require.Nil(t, err)
	require.NotNil(t, confirmation)
	assert.Equal(t, status, confirmation.Status)
	assert.Equal(t, exiResponse, confirmation.ExiResponse)
	assert.Equal(t, contractSignatureCertificateChain.Certificate, confirmation.ContractSignatureCertificateChain.Certificate)
	assert.Equal(t, saProvisioningCertificateChain.Certificate, confirmation.SaProvisioningCertificateChain.Certificate)
}

func (suite *OcppV2TestSuite) TestGet15118EVCertificateInvalidEndpoint() {
	messageId := defaultMessageId
	schemaVersion := "1.0"
	exiRequest := "deadbeef"
	get15118EVCertificateRequest := iso15118.NewGet15118EVCertificateRequest(schemaVersion, exiRequest)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"15118SchemaVersion":"%v","exiRequest":"%v"}]`, messageId, iso15118.Get15118EVCertificateFeatureName, schemaVersion, exiRequest)
	testUnsupportedRequestFromCentralSystem(suite, get15118EVCertificateRequest, requestJson, messageId)
}
