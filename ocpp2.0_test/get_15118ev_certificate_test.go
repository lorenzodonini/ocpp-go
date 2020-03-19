package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test
func (suite *OcppV2TestSuite) TestGet15118EVCertificateRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{ocpp2.Get15118EVCertificateRequest{SchemaVersion: "1.0", ExiRequest: "deadbeef"}, true},
		{ocpp2.Get15118EVCertificateRequest{SchemaVersion: "1.0"}, false},
		{ocpp2.Get15118EVCertificateRequest{ExiRequest: "deadbeef"}, false},
		{ocpp2.Get15118EVCertificateRequest{}, false},
		{ocpp2.Get15118EVCertificateRequest{SchemaVersion: ">50................................................", ExiRequest: "deadbeef"}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestGet15118EVCertificateConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{ocpp2.Get15118EVCertificateConfirmation{Status: ocpp2.Certificate15188EVStatusAccepted, ExiResponse: "deadbeef", ContractSignatureCertificateChain: ocpp2.CertificateChain{Certificate: "deadcode"}, SaProvisioningCertificateChain: ocpp2.CertificateChain{Certificate: "deadcode2"}}, true},
		{ocpp2.Get15118EVCertificateConfirmation{Status: ocpp2.Certificate15188EVStatusAccepted, ExiResponse: "deadbeef", ContractSignatureCertificateChain: ocpp2.CertificateChain{Certificate: "deadcode", ChildCertificate: []string{"c1", "c2", "c3", "c4"}}, SaProvisioningCertificateChain: ocpp2.CertificateChain{Certificate: "deadcode2"}}, true},
		{ocpp2.Get15118EVCertificateConfirmation{Status: ocpp2.Certificate15188EVStatusAccepted, ExiResponse: "deadbeef", ContractSignatureCertificateChain: ocpp2.CertificateChain{Certificate: "deadcode"}}, false},
		{ocpp2.Get15118EVCertificateConfirmation{Status: ocpp2.Certificate15188EVStatusAccepted, ExiResponse: "deadbeef", SaProvisioningCertificateChain: ocpp2.CertificateChain{Certificate: "deadcode2"}}, false},
		{ocpp2.Get15118EVCertificateConfirmation{Status: ocpp2.Certificate15188EVStatusAccepted, ContractSignatureCertificateChain: ocpp2.CertificateChain{Certificate: "deadcode"}, SaProvisioningCertificateChain: ocpp2.CertificateChain{Certificate: "deadcode2"}}, false},
		{ocpp2.Get15118EVCertificateConfirmation{ExiResponse: "deadbeef", ContractSignatureCertificateChain: ocpp2.CertificateChain{Certificate: "deadcode"}, SaProvisioningCertificateChain: ocpp2.CertificateChain{Certificate: "deadcode2"}}, false},
		{ocpp2.Get15118EVCertificateConfirmation{}, false},
		{ocpp2.Get15118EVCertificateConfirmation{Status: "invalidCertificateStatus", ExiResponse: "deadbeef", ContractSignatureCertificateChain: ocpp2.CertificateChain{Certificate: "deadcode"}, SaProvisioningCertificateChain: ocpp2.CertificateChain{Certificate: "deadcode2"}}, false},
		{ocpp2.Get15118EVCertificateConfirmation{Status: ocpp2.Certificate15188EVStatusAccepted, ExiResponse: "deadbeef", ContractSignatureCertificateChain: ocpp2.CertificateChain{Certificate: "deadcode"}, SaProvisioningCertificateChain: ocpp2.CertificateChain{}}, false},
		{ocpp2.Get15118EVCertificateConfirmation{Status: ocpp2.Certificate15188EVStatusAccepted, ExiResponse: "deadbeef", ContractSignatureCertificateChain: ocpp2.CertificateChain{Certificate: "deadcode", ChildCertificate: []string{"c1", "c2", "c3", "c4", "c5"}}, SaProvisioningCertificateChain: ocpp2.CertificateChain{Certificate: "deadcode2"}}, false},
		{ocpp2.Get15118EVCertificateConfirmation{Status: ocpp2.Certificate15188EVStatusAccepted, ExiResponse: "deadbeef", ContractSignatureCertificateChain: ocpp2.CertificateChain{Certificate: "deadcode", ChildCertificate: []string{"c1", "c2", "c3", ""}}, SaProvisioningCertificateChain: ocpp2.CertificateChain{Certificate: "deadcode2"}}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestGet15118EVCertificateE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	status := ocpp2.Certificate15188EVStatusAccepted
	schemaVersion := "1.0"
	exiRequest := "deadbeef"
	exiResponse := "deadbeef2"
	contractSignatureCertificateChain := ocpp2.CertificateChain{Certificate: "deadcode"}
	saProvisioningCertificateChain := ocpp2.CertificateChain{Certificate: "deadcode2"}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"15118SchemaVersion":"%v","exiRequest":"%v"}]`, messageId, ocpp2.Get15118EVCertificateFeatureName, schemaVersion, exiRequest)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v","exiResponse":"%v","contractSignatureCertificateChain":{"certificate":"%v"},"saProvisioningCertificateChain":{"certificate":"%v"}}]`,
		messageId, status, exiResponse, contractSignatureCertificateChain.Certificate, saProvisioningCertificateChain.Certificate)
	Get15118EVCertificateConfirmation := ocpp2.NewGet15118EVCertificateConfirmation(status, exiResponse, contractSignatureCertificateChain, saProvisioningCertificateChain)
	channel := NewMockWebSocket(wsId)

	coreListener := MockCentralSystemCoreListener{}
	coreListener.On("OnGet15118EVCertificate", mock.AnythingOfType("string"), mock.Anything).Return(Get15118EVCertificateConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*ocpp2.Get15118EVCertificateRequest)
		require.True(t, ok)
		assert.Equal(t, schemaVersion, request.SchemaVersion)
		assert.Equal(t, exiRequest, request.ExiRequest)
	})
	setupDefaultCentralSystemHandlers(suite, coreListener, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	require.Nil(t, err)
	confirmation, err := suite.chargePoint.Get15118EVCertificate(schemaVersion, exiRequest)
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
	firmwareStatusRequest := ocpp2.NewGet15118EVCertificateRequest(schemaVersion, exiRequest)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"15118SchemaVersion":"%v","exiRequest":"%v"}]`, messageId, ocpp2.Get15118EVCertificateFeatureName, schemaVersion, exiRequest)
	testUnsupportedRequestFromCentralSystem(suite, firmwareStatusRequest, requestJson, messageId)
}
