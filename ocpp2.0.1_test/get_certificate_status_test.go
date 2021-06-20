package ocpp2_test

import (
	"fmt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/iso15118"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

// Test
func (suite *OcppV2TestSuite) TestGetCertificateStatusRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{iso15118.GetCertificateStatusRequest{OcspRequestData: types.OCSPRequestDataType{HashAlgorithm: types.SHA256, IssuerNameHash: "hash00", IssuerKeyHash: "hash01", SerialNumber: "serial0", ResponderURL: "http://someUrl"}}, true},
		{iso15118.GetCertificateStatusRequest{}, false},
		{iso15118.GetCertificateStatusRequest{OcspRequestData: types.OCSPRequestDataType{HashAlgorithm: "invalidHashAlgorithm", IssuerNameHash: "hash00", IssuerKeyHash: "hash01", SerialNumber: "serial0", ResponderURL: "http://someUrl"}}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestGetCertificateStatusConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{iso15118.GetCertificateStatusResponse{Status: types.GenericStatusAccepted, OcspResult: "deadbeef"}, true},
		{iso15118.GetCertificateStatusResponse{Status: types.GenericStatusAccepted}, true},
		{iso15118.GetCertificateStatusResponse{Status: types.GenericStatusRejected}, true},
		{iso15118.GetCertificateStatusResponse{Status: "invalidGenericStatus"}, false},
		{iso15118.GetCertificateStatusResponse{Status: types.GenericStatusAccepted, OcspResult: newLongString(5501)}, false},
		{iso15118.GetCertificateStatusResponse{}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestGetCertificateStatusE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	ocspData := types.OCSPRequestDataType{HashAlgorithm: types.SHA256, IssuerNameHash: "hash00", IssuerKeyHash: "hash01", SerialNumber: "serial0", ResponderURL: "http://someUrl"}
	ocspResult := "deadbeef"
	status := types.GenericStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"ocspRequestData":{"hashAlgorithm":"%v","issuerNameHash":"%v","issuerKeyHash":"%v","serialNumber":"%v","responderURL":"%v"}}]`,
		messageId, iso15118.GetCertificateStatusFeatureName, ocspData.HashAlgorithm, ocspData.IssuerNameHash, ocspData.IssuerKeyHash, ocspData.SerialNumber, ocspData.ResponderURL)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v","ocspResult":"%v"}]`, messageId, status, ocspResult)
	getCertificateStatusConfirmation := iso15118.NewGetCertificateStatusResponse(status)
	getCertificateStatusConfirmation.OcspResult = ocspResult
	channel := NewMockWebSocket(wsId)

	handler := &MockCSMSIso15118Handler{}
	handler.On("OnGetCertificateStatus", mock.AnythingOfType("string"), mock.Anything).Return(getCertificateStatusConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*iso15118.GetCertificateStatusRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, ocspData.HashAlgorithm, request.OcspRequestData.HashAlgorithm)
		assert.Equal(t, ocspData.IssuerNameHash, request.OcspRequestData.IssuerNameHash)
		assert.Equal(t, ocspData.IssuerKeyHash, request.OcspRequestData.IssuerKeyHash)
		assert.Equal(t, ocspData.SerialNumber, request.OcspRequestData.SerialNumber)
		assert.Equal(t, ocspData.ResponderURL, request.OcspRequestData.ResponderURL)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	confirmation, err := suite.chargingStation.GetCertificateStatus(ocspData)
	require.Nil(t, err)
	require.NotNil(t, confirmation)
	assert.Equal(t, status, confirmation.Status)
	assert.Equal(t, ocspResult, confirmation.OcspResult)
}

func (suite *OcppV2TestSuite) TestGetCertificateStatusInvalidEndpoint() {
	messageId := defaultMessageId
	ocspData := types.OCSPRequestDataType{HashAlgorithm: types.SHA256, IssuerNameHash: "hash00", IssuerKeyHash: "hash01", SerialNumber: "serial0", ResponderURL: "http://someUrl"}
	getCertificateStatusRequest := iso15118.NewGetCertificateStatusRequest(ocspData)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"ocspRequestData":{"hashAlgorithm":"%v","issuerNameHash":"%v","issuerKeyHash":"%v","serialNumber":"%v","responderURL":"%v"}}]`,
		messageId, iso15118.GetCertificateStatusFeatureName, ocspData.HashAlgorithm, ocspData.IssuerNameHash, ocspData.IssuerKeyHash, ocspData.SerialNumber, ocspData.ResponderURL)
	testUnsupportedRequestFromCentralSystem(suite, getCertificateStatusRequest, requestJson, messageId)
}
