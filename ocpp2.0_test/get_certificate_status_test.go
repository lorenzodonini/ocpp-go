package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test
func (suite *OcppV2TestSuite) TestGetCertificateStatusRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{ocpp2.GetCertificateStatusRequest{OcspRequestData: ocpp2.OCSPRequestDataType{HashAlgorithm: ocpp2.SHA256, IssuerNameHash: "hash00", IssuerKeyHash: "hash01", SerialNumber: "serial0", ResponderURL: "http://someUrl"}}, true},
		{ocpp2.GetCertificateStatusRequest{}, false},
		{ocpp2.GetCertificateStatusRequest{OcspRequestData: ocpp2.OCSPRequestDataType{HashAlgorithm: "invalidHashAlgorithm", IssuerNameHash: "hash00", IssuerKeyHash: "hash01", SerialNumber: "serial0", ResponderURL: "http://someUrl"}}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestGetCertificateStatusConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{ocpp2.GetCertificateStatusConfirmation{Status: ocpp2.GenericStatusAccepted, OcspResult: "deadbeef"}, true},
		{ocpp2.GetCertificateStatusConfirmation{Status: ocpp2.GenericStatusAccepted}, true},
		{ocpp2.GetCertificateStatusConfirmation{Status: ocpp2.GenericStatusRejected}, true},
		{ocpp2.GetCertificateStatusConfirmation{Status: "invalidGenericStatus"}, false},
		{ocpp2.GetCertificateStatusConfirmation{}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestGetCertificateStatusE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	ocspData := ocpp2.OCSPRequestDataType{HashAlgorithm: ocpp2.SHA256, IssuerNameHash: "hash00", IssuerKeyHash: "hash01", SerialNumber: "serial0", ResponderURL: "http://someUrl"}
	ocspResult := "deadbeef"
	status := ocpp2.GenericStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"ocspRequestData":{"hashAlgorithm":"%v","issuerNameHash":"%v","issuerKeyHash":"%v","serialNumber":"%v","responderURL":"%v"}}]`,
		messageId, ocpp2.GetCertificateStatusFeatureName, ocspData.HashAlgorithm, ocspData.IssuerNameHash, ocspData.IssuerKeyHash, ocspData.SerialNumber, ocspData.ResponderURL)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v","ocspResult":"%v"}]`, messageId, status, ocspResult)
	getCertificateStatusConfirmation := ocpp2.NewGetCertificateStatusConfirmation(status)
	getCertificateStatusConfirmation.OcspResult = ocspResult
	channel := NewMockWebSocket(wsId)

	coreListener := MockCentralSystemCoreListener{}
	coreListener.On("OnGetCertificateStatus", mock.AnythingOfType("string"), mock.Anything).Return(getCertificateStatusConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*ocpp2.GetCertificateStatusRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, ocspData.HashAlgorithm, request.OcspRequestData.HashAlgorithm)
		assert.Equal(t, ocspData.IssuerNameHash, request.OcspRequestData.IssuerNameHash)
		assert.Equal(t, ocspData.IssuerKeyHash, request.OcspRequestData.IssuerKeyHash)
		assert.Equal(t, ocspData.SerialNumber, request.OcspRequestData.SerialNumber)
		assert.Equal(t, ocspData.ResponderURL, request.OcspRequestData.ResponderURL)
	})
	setupDefaultCentralSystemHandlers(suite, coreListener, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	require.Nil(t, err)
	confirmation, err := suite.chargePoint.GetCertificateStatus(ocspData)
	require.Nil(t, err)
	require.NotNil(t, confirmation)
	assert.Equal(t, status, confirmation.Status)
	assert.Equal(t, ocspResult, confirmation.OcspResult)
}

func (suite *OcppV2TestSuite) TestGetCertificateStatusInvalidEndpoint() {
	messageId := defaultMessageId
	ocspData := ocpp2.OCSPRequestDataType{HashAlgorithm: ocpp2.SHA256, IssuerNameHash: "hash00", IssuerKeyHash: "hash01", SerialNumber: "serial0", ResponderURL: "http://someUrl"}
	getCertificateStatusRequest := ocpp2.NewGetCertificateStatusRequest(ocspData)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"ocspRequestData":{"hashAlgorithm":"%v","issuerNameHash":"%v","issuerKeyHash":"%v","serialNumber":"%v","responderURL":"%v"}}]`,
		messageId, ocpp2.GetCertificateStatusFeatureName, ocspData.HashAlgorithm, ocspData.IssuerNameHash, ocspData.IssuerKeyHash, ocspData.SerialNumber, ocspData.ResponderURL)
	testUnsupportedRequestFromCentralSystem(suite, getCertificateStatusRequest, requestJson, messageId)
}
