package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/diagnostics"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test
func (suite *OcppV2TestSuite) TestCustomerInformationRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{diagnostics.CustomerInformationRequest{RequestID: 42, Report: true, Clear: true, CustomerIdentifier: "0001", IdToken: &types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode, AdditionalInfo: nil}, CustomerCertificate: &types.CertificateHashData{HashAlgorithm: types.SHA256, IssuerNameHash: "hash00", IssuerKeyHash: "hash01", SerialNumber: "serial0"}}, true},
		{diagnostics.CustomerInformationRequest{RequestID: 42, Report: true, Clear: true, CustomerIdentifier: "0001", IdToken: &types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode, AdditionalInfo: nil}}, true},
		{diagnostics.CustomerInformationRequest{RequestID: 42, Report: true, Clear: true, CustomerIdentifier: "0001"}, true},
		{diagnostics.CustomerInformationRequest{RequestID: 42, Report: true, Clear: true}, true},
		{diagnostics.CustomerInformationRequest{RequestID: 42, Report: true}, true},
		{diagnostics.CustomerInformationRequest{RequestID: 42, Clear: true}, true},
		{diagnostics.CustomerInformationRequest{Report: true, Clear: true}, true},
		{diagnostics.CustomerInformationRequest{}, true},
		{diagnostics.CustomerInformationRequest{RequestID: -1, Report: true, Clear: true}, false},
		{diagnostics.CustomerInformationRequest{RequestID: 42, Report: true, Clear: true, CustomerIdentifier: ">64.............................................................."}, false},
		{diagnostics.CustomerInformationRequest{RequestID: 42, Report: true, Clear: true, IdToken: &types.IdToken{IdToken: "1234", Type: "invalidTokenType", AdditionalInfo: nil}}, false},
		{diagnostics.CustomerInformationRequest{RequestID: 42, Report: true, Clear: true, CustomerCertificate: &types.CertificateHashData{HashAlgorithm: "invalidHasAlgorithm", IssuerNameHash: "hash00", IssuerKeyHash: "hash01", SerialNumber: "serial0"}}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestCustomerInformationConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{diagnostics.CustomerInformationResponse{Status: diagnostics.CustomerInformationStatusAccepted}, true},
		{diagnostics.CustomerInformationResponse{}, false},
		{diagnostics.CustomerInformationResponse{Status: "invalidCustomerInformationStatus"}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestCustomerInformationE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	requestId := 42
	report := true
	clear := true
	customerId := "0001"
	idToken := types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode}
	customerCertificate := types.CertificateHashData{HashAlgorithm: types.SHA256, IssuerNameHash: "hash00", IssuerKeyHash: "hash01", SerialNumber: "serial0"}
	status := diagnostics.CustomerInformationStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"requestId":%v,"report":%v,"clear":%v,"customerIdentifier":"%v","idToken":{"idToken":"%v","type":"%v"},"customerCertificate":{"hashAlgorithm":"%v","issuerNameHash":"%v","issuerKeyHash":"%v","serialNumber":"%v"}}]`,
		messageId, diagnostics.CustomerInformationFeatureName, requestId, report, clear, customerId, idToken.IdToken, idToken.Type, customerCertificate.HashAlgorithm, customerCertificate.IssuerNameHash, customerCertificate.IssuerKeyHash, customerCertificate.SerialNumber)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	customerInformationConfirmation := diagnostics.NewCustomerInformationResponse(status)
	channel := NewMockWebSocket(wsId)

	handler := MockChargingStationDiagnosticsHandler{}
	handler.On("OnCustomerInformation", mock.Anything).Return(customerInformationConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*diagnostics.CustomerInformationRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, requestId, request.RequestID)
		assert.Equal(t, report, request.Report)
		assert.Equal(t, clear, request.Clear)
		assert.Equal(t, customerId, request.CustomerIdentifier)
		require.NotNil(t, request.IdToken)
		require.NotNil(t, request.CustomerCertificate)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.CustomerInformation(wsId, func(confirmation *diagnostics.CustomerInformationResponse, err error) {
		require.Nil(t, err)
		require.NotNil(t, confirmation)
		require.Equal(t, status, confirmation.Status)
		resultChannel <- true
	}, requestId, report, clear, func(request *diagnostics.CustomerInformationRequest) {
		request.CustomerIdentifier = customerId
		request.IdToken = &idToken
		request.CustomerCertificate = &customerCertificate
	})
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV2TestSuite) TestCustomerInformationInvalidEndpoint() {
	messageId := defaultMessageId
	requestId := 42
	report := true
	clear := true
	customerId := "0001"
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"requestId":%v,"report":%v,"clear":%v,"customerIdentifier":"%v"}]`, messageId, diagnostics.CustomerInformationFeatureName, requestId, report, clear, customerId)
	customerInformationRequest := diagnostics.NewCustomerInformationRequest(requestId, report, clear)
	testUnsupportedRequestFromChargingStation(suite, customerInformationRequest, requestJson, messageId)
}
