package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test
func (suite *OcppV2TestSuite) TestCustomerInformationRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{ocpp2.CustomerInformationRequest{RequestID: 42, Report: true, Clear: true, CustomerIdentifier: "0001", IdToken: &ocpp2.IdToken{IdToken: "1234", Type: ocpp2.IdTokenTypeKeyCode, AdditionalInfo: nil}, CustomerCertificate: &ocpp2.CertificateHashData{HashAlgorithm: ocpp2.SHA256, IssuerNameHash: "hash00", IssuerKeyHash: "hash01", SerialNumber: "serial0"}}, true},
		{ocpp2.CustomerInformationRequest{RequestID: 42, Report: true, Clear: true, CustomerIdentifier: "0001", IdToken: &ocpp2.IdToken{IdToken: "1234", Type: ocpp2.IdTokenTypeKeyCode, AdditionalInfo: nil}}, true},
		{ocpp2.CustomerInformationRequest{RequestID: 42, Report: true, Clear: true, CustomerIdentifier: "0001"}, true},
		{ocpp2.CustomerInformationRequest{RequestID: 42, Report: true, Clear: true}, true},
		{ocpp2.CustomerInformationRequest{RequestID: 42, Report: true}, true},
		{ocpp2.CustomerInformationRequest{RequestID: 42, Clear: true}, true},
		{ocpp2.CustomerInformationRequest{Report: true, Clear: true}, true},
		{ocpp2.CustomerInformationRequest{}, true},
		{ocpp2.CustomerInformationRequest{RequestID: -1, Report: true, Clear: true}, false},
		{ocpp2.CustomerInformationRequest{RequestID: 42, Report: true, Clear: true, CustomerIdentifier: ">64.............................................................."}, false},
		{ocpp2.CustomerInformationRequest{RequestID: 42, Report: true, Clear: true, IdToken: &ocpp2.IdToken{IdToken: "1234", Type: "invalidTokenType", AdditionalInfo: nil}}, false},
		{ocpp2.CustomerInformationRequest{RequestID: 42, Report: true, Clear: true, CustomerCertificate: &ocpp2.CertificateHashData{HashAlgorithm: "invalidHasAlgorithm", IssuerNameHash: "hash00", IssuerKeyHash: "hash01", SerialNumber: "serial0"}}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestCustomerInformationConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{ocpp2.CustomerInformationConfirmation{Status: ocpp2.CustomerInformationStatusAccepted}, true},
		{ocpp2.CustomerInformationConfirmation{}, false},
		{ocpp2.CustomerInformationConfirmation{Status: "invalidCustomerInformationStatus"}, false},
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
	idToken := ocpp2.IdToken{IdToken: "1234", Type: ocpp2.IdTokenTypeKeyCode}
	customerCertificate := ocpp2.CertificateHashData{HashAlgorithm: ocpp2.SHA256, IssuerNameHash: "hash00", IssuerKeyHash: "hash01", SerialNumber: "serial0"}
	status := ocpp2.CustomerInformationStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"requestId":%v,"report":%v,"clear":%v,"customerIdentifier":"%v","idToken":{"idToken":"%v","type":"%v"},"customerCertificate":{"hashAlgorithm":"%v","issuerNameHash":"%v","issuerKeyHash":"%v","serialNumber":"%v"}}]`,
		messageId, ocpp2.CustomerInformationFeatureName, requestId, report, clear, customerId, idToken.IdToken, idToken.Type, customerCertificate.HashAlgorithm, customerCertificate.IssuerNameHash, customerCertificate.IssuerKeyHash, customerCertificate.SerialNumber)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	customerInformationConfirmation := ocpp2.NewCustomerInformationConfirmation(status)
	channel := NewMockWebSocket(wsId)

	coreListener := MockChargePointCoreListener{}
	coreListener.On("OnCustomerInformation", mock.Anything).Return(customerInformationConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*ocpp2.CustomerInformationRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, requestId, request.RequestID)
		assert.Equal(t, report, request.Report)
		assert.Equal(t, clear, request.Clear)
		assert.Equal(t, customerId, request.CustomerIdentifier)
		require.NotNil(t, request.IdToken)
		require.NotNil(t, request.CustomerCertificate)
	})
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, coreListener, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.CustomerInformation(wsId, func(confirmation *ocpp2.CustomerInformationConfirmation, err error) {
		require.Nil(t, err)
		require.NotNil(t, confirmation)
		require.Equal(t, status, confirmation.Status)
		resultChannel <- true
	}, requestId, report, clear, func(request *ocpp2.CustomerInformationRequest) {
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
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"requestId":%v,"report":%v,"clear":%v,"customerIdentifier":"%v"}]`, messageId, ocpp2.CustomerInformationFeatureName, requestId, report, clear, customerId)
	customerInformationRequest := ocpp2.NewCustomerInformationRequest(requestId, report, clear)
	testUnsupportedRequestFromChargePoint(suite, customerInformationRequest, requestJson, messageId)
}
