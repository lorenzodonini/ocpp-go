package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test
func (suite *OcppV2TestSuite) TestAuthorizeRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{ocpp2.AuthorizeRequest{EvseID: []int{4,2}, IdToken: ocpp2.IdToken{IdToken: "1234", Type: ocpp2.IdTokenTypeKeyCode, AdditionalInfo: []ocpp2.AdditionalInfo{{AdditionalIdToken: "0000", Type: "someType"}}}, CertificateHashData: []ocpp2.OCSPRequestDataType{{SerialNumber: "serial0", HashAlgorithm: ocpp2.SHA256, IssuerNameHash: "hash0", IssuerKeyHash: "hash1", ResponderURL: "www.someurl.com"}}}, true},
		{ocpp2.AuthorizeRequest{EvseID: []int{4,2}, IdToken: ocpp2.IdToken{IdToken: "1234", Type: ocpp2.IdTokenTypeKeyCode, AdditionalInfo: []ocpp2.AdditionalInfo{{AdditionalIdToken: "0000", Type: "someType"}}}}, true},
		{ocpp2.AuthorizeRequest{IdToken: ocpp2.IdToken{IdToken: "1234", Type: ocpp2.IdTokenTypeKeyCode, AdditionalInfo: []ocpp2.AdditionalInfo{{AdditionalIdToken: "0000", Type: "someType"}}}, CertificateHashData: []ocpp2.OCSPRequestDataType{{SerialNumber: "serial0", HashAlgorithm: ocpp2.SHA256, IssuerNameHash: "hash0", IssuerKeyHash: "hash1", ResponderURL: "www.someurl.com"}}}, true},
		{ocpp2.AuthorizeRequest{IdToken: ocpp2.IdToken{IdToken: "1234", Type: ocpp2.IdTokenTypeKeyCode, AdditionalInfo: []ocpp2.AdditionalInfo{{AdditionalIdToken: "0000", Type: "someType"}}}, CertificateHashData: []ocpp2.OCSPRequestDataType{}}, true},
		{ocpp2.AuthorizeRequest{EvseID: []int{}, IdToken: ocpp2.IdToken{IdToken: "1234", Type: ocpp2.IdTokenTypeKeyCode, AdditionalInfo: []ocpp2.AdditionalInfo{{AdditionalIdToken: "0000", Type: "someType"}}}}, true},
		{ocpp2.AuthorizeRequest{}, false},
		{ocpp2.AuthorizeRequest{IdToken: ocpp2.IdToken{Type: ocpp2.IdTokenTypeKeyCode, AdditionalInfo: []ocpp2.AdditionalInfo{{AdditionalIdToken: "0000", Type: "someType"}}}}, false},
		{ocpp2.AuthorizeRequest{IdToken: ocpp2.IdToken{Type: ocpp2.IdTokenTypeKeyCode, AdditionalInfo: []ocpp2.AdditionalInfo{{AdditionalIdToken: "0000", Type: "someType"}}}, CertificateHashData: []ocpp2.OCSPRequestDataType{{HashAlgorithm: ocpp2.SHA256, IssuerNameHash: "hash0", IssuerKeyHash: "hash1"}}}, false},
		{ocpp2.AuthorizeRequest{IdToken: ocpp2.IdToken{Type: ocpp2.IdTokenTypeKeyCode, AdditionalInfo: []ocpp2.AdditionalInfo{{AdditionalIdToken: "0000", Type: "someType"}}}, CertificateHashData: []ocpp2.OCSPRequestDataType{{SerialNumber: "s0", HashAlgorithm: ocpp2.SHA256, IssuerNameHash: "h0", IssuerKeyHash: "h0.1"},{SerialNumber: "s1", HashAlgorithm: ocpp2.SHA256, IssuerNameHash: "h1", IssuerKeyHash: "h1.1"},{SerialNumber: "s2", HashAlgorithm: ocpp2.SHA256, IssuerNameHash: "h2", IssuerKeyHash: "h2.1"},{SerialNumber: "s3", HashAlgorithm: ocpp2.SHA256, IssuerNameHash: "h3", IssuerKeyHash: "h3.1"},{SerialNumber: "s4", HashAlgorithm: ocpp2.SHA256, IssuerNameHash: "h4", IssuerKeyHash: "h4.1"}}}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestAuthorizeConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{ocpp2.AuthorizeConfirmation{CertificateStatus: ocpp2.CertificateStatusAccepted, EvseID: []int{4,2}, IdTokenInfo: ocpp2.IdTokenInfo{Status: ocpp2.AuthorizationStatusAccepted}}, true},
		{ocpp2.AuthorizeConfirmation{CertificateStatus: ocpp2.CertificateStatusAccepted, IdTokenInfo: ocpp2.IdTokenInfo{Status: ocpp2.AuthorizationStatusAccepted}}, true},
		{ocpp2.AuthorizeConfirmation{IdTokenInfo: ocpp2.IdTokenInfo{Status: ocpp2.AuthorizationStatusAccepted}}, true},
		{ocpp2.AuthorizeConfirmation{}, false},
		{ocpp2.AuthorizeConfirmation{CertificateStatus:"invalidCertificateStatus", EvseID: []int{4,2}, IdTokenInfo: ocpp2.IdTokenInfo{Status: ocpp2.AuthorizationStatusAccepted}}, false},
		{ocpp2.AuthorizeConfirmation{CertificateStatus:"invalidCertificateStatus", EvseID: []int{4,2}, IdTokenInfo: ocpp2.IdTokenInfo{Status: "invalidTokenInfoStatus"}}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestAuthorizeE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	evseIds := []int{4,2}
	additionalInfo := ocpp2.AdditionalInfo{AdditionalIdToken: "at1", Type: "some"}
	idToken := ocpp2.IdToken{IdToken: "tok1", Type: ocpp2.IdTokenTypeKeyCode, AdditionalInfo: []ocpp2.AdditionalInfo{additionalInfo}}
	certHashData := ocpp2.OCSPRequestDataType{HashAlgorithm: ocpp2.SHA256, IssuerNameHash: "h0", IssuerKeyHash: "h0.1", SerialNumber: "s0", ResponderURL: "http://www.test.org"}
	status := ocpp2.AuthorizationStatusAccepted
	certificateStatus := ocpp2.CertificateStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"evseId":[%v,%v],"idToken":{"idToken":"%v","type":"%v","additionalInfo":[{"additionalIdToken":"%v","type":"%v"}]},"15118CertificateHashData":[{"hashAlgorithm":"%v","issuerNameHash":"%v","issuerKeyHash":"%v","serialNumber":"%v","responderURL":"%v"}]}]`,
		messageId, ocpp2.AuthorizeFeatureName, evseIds[0], evseIds[1], idToken.IdToken, idToken.Type, additionalInfo.AdditionalIdToken, additionalInfo.Type, certHashData.HashAlgorithm, certHashData.IssuerNameHash, certHashData.IssuerKeyHash, certHashData.SerialNumber, certHashData.ResponderURL)
	responseJson := fmt.Sprintf(`[3,"%v",{"certificateStatus":"%v","evseId":[%v,%v],"idTokenInfo":{"status":"%v"}}]`,
		messageId, certificateStatus, evseIds[0], evseIds[1], status)
	authorizeConfirmation := ocpp2.NewAuthorizationConfirmation(ocpp2.IdTokenInfo{Status: status})
	authorizeConfirmation.EvseID = evseIds
	authorizeConfirmation.CertificateStatus = certificateStatus
	requestRaw := []byte(requestJson)
	responseRaw := []byte(responseJson)
	channel := NewMockWebSocket(wsId)

	coreListener := MockCentralSystemCoreListener{}
	coreListener.On("OnAuthorize", mock.AnythingOfType("string"), mock.Anything).Return(authorizeConfirmation, nil).Run(func(args mock.Arguments) {
		request := args.Get(1).(*ocpp2.AuthorizeRequest)
		require.Len(t, request.EvseID, 2)
		assert.Equal(t, evseIds[0], request.EvseID[0])
		assert.Equal(t, evseIds[1], request.EvseID[1])
		assert.Equal(t, idToken.IdToken, request.IdToken.IdToken)
		assert.Equal(t, idToken.Type, request.IdToken.Type)
		require.Len(t, request.IdToken.AdditionalInfo, 1)
		assert.Equal(t, idToken.AdditionalInfo[0].AdditionalIdToken, request.IdToken.AdditionalInfo[0].AdditionalIdToken)
		assert.Equal(t, idToken.AdditionalInfo[0].Type, request.IdToken.AdditionalInfo[0].Type)
		require.Len(t, request.CertificateHashData, 1)
		assert.Equal(t, certHashData.HashAlgorithm, request.CertificateHashData[0].HashAlgorithm)
		assert.Equal(t, certHashData.IssuerNameHash, request.CertificateHashData[0].IssuerNameHash)
		assert.Equal(t, certHashData.IssuerKeyHash, request.CertificateHashData[0].IssuerKeyHash)
		assert.Equal(t, certHashData.SerialNumber, request.CertificateHashData[0].SerialNumber)
		assert.Equal(t, certHashData.ResponderURL, request.CertificateHashData[0].ResponderURL)
	})
	setupDefaultCentralSystemHandlers(suite, coreListener, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: responseRaw, forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: requestRaw, forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	require.Nil(t, err)
	confirmation, err := suite.chargePoint.Authorize(idToken.IdToken, idToken.Type, func(request *ocpp2.AuthorizeRequest) {
		request.IdToken.AdditionalInfo = []ocpp2.AdditionalInfo{additionalInfo}
		request.EvseID = evseIds
		request.CertificateHashData = []ocpp2.OCSPRequestDataType{certHashData}
	})
	require.Nil(t, err)
	require.NotNil(t, confirmation)
	require.Len(t, confirmation.EvseID, 2)
	assert.Equal(t, evseIds[0], confirmation.EvseID[0])
	assert.Equal(t, evseIds[1], confirmation.EvseID[1])
	assert.Equal(t, certificateStatus, confirmation.CertificateStatus)
	assert.Equal(t, status, confirmation.IdTokenInfo.Status)
}

func (suite *OcppV2TestSuite) TestAuthorizeInvalidEndpoint() {
	messageId := defaultMessageId
	evseIds := []int{4,2}
	additionalInfo := ocpp2.AdditionalInfo{AdditionalIdToken: "at1", Type: "some"}
	idToken := ocpp2.IdToken{IdToken: "tok1", Type: ocpp2.IdTokenTypeKeyCode, AdditionalInfo: []ocpp2.AdditionalInfo{additionalInfo}}
	certHashData := ocpp2.OCSPRequestDataType{HashAlgorithm: ocpp2.SHA256, IssuerNameHash: "h0", IssuerKeyHash: "h0.1", SerialNumber: "s0", ResponderURL: "http://www.test.org"}
	authorizeRequest := ocpp2.NewAuthorizationRequest(idToken.IdToken, idToken.Type)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"evseId":[%v,%v],"idToken":{"idToken":"%v","type":"%v","additionalInfo":[{"additionalIdToken":"%v","type":"%v"}]},"15118CertificateHashData":[{"hashAlgorithm":"%v","issuerNameHash":"%v","issuerKeyHash":"%v","serialNumber":"%v","responderURL":"%v"}]}]`,
		messageId, ocpp2.AuthorizeFeatureName, evseIds[0], evseIds[1], idToken.IdToken, idToken.Type, additionalInfo.AdditionalIdToken, additionalInfo.Type, certHashData.HashAlgorithm, certHashData.IssuerNameHash, certHashData.IssuerKeyHash, certHashData.SerialNumber, certHashData.ResponderURL)
	testUnsupportedRequestFromCentralSystem(suite, authorizeRequest, requestJson, messageId)
}
