package ocpp2_test

import (
	"fmt"

	"github.com/stretchr/testify/mock"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/authorization"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

// Test
func (suite *OcppV2TestSuite) TestAuthorizeRequestValidation() {
	var requestTable = []GenericTestEntry{
		{authorization.AuthorizeRequest{Certificate: "deadc0de", IdToken: types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode, AdditionalInfo: []types.AdditionalInfo{{AdditionalIdToken: "0000", Type: "someType"}}}, CertificateHashData: []types.OCSPRequestDataType{{SerialNumber: "serial0", HashAlgorithm: types.SHA256, IssuerNameHash: "hash0", IssuerKeyHash: "hash1", ResponderURL: "www.someurl.com"}}}, true},
		{authorization.AuthorizeRequest{Certificate: "deadc0de", IdToken: types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode, AdditionalInfo: []types.AdditionalInfo{{AdditionalIdToken: "0000", Type: "someType"}}}}, true},
		{authorization.AuthorizeRequest{IdToken: types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode, AdditionalInfo: []types.AdditionalInfo{{AdditionalIdToken: "0000", Type: "someType"}}}, CertificateHashData: []types.OCSPRequestDataType{{SerialNumber: "serial0", HashAlgorithm: types.SHA256, IssuerNameHash: "hash0", IssuerKeyHash: "hash1", ResponderURL: "www.someurl.com"}}}, true},
		{authorization.AuthorizeRequest{IdToken: types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode, AdditionalInfo: []types.AdditionalInfo{{AdditionalIdToken: "0000", Type: "someType"}}}, CertificateHashData: []types.OCSPRequestDataType{}}, true},
		{authorization.AuthorizeRequest{IdToken: types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode, AdditionalInfo: []types.AdditionalInfo{{AdditionalIdToken: "0000", Type: "someType"}}}}, true},
		{authorization.AuthorizeRequest{}, false},
		{authorization.AuthorizeRequest{Certificate: newLongString(5501), IdToken: types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode, AdditionalInfo: []types.AdditionalInfo{{AdditionalIdToken: "0000", Type: "someType"}}}}, false},
		{authorization.AuthorizeRequest{Certificate: "deadc0de", IdToken: types.IdToken{Type: types.IdTokenTypeKeyCode, AdditionalInfo: []types.AdditionalInfo{{AdditionalIdToken: "0000", Type: "someType"}}}}, false},
		{authorization.AuthorizeRequest{Certificate: "deadc0de", IdToken: types.IdToken{Type: types.IdTokenTypeKeyCode, AdditionalInfo: []types.AdditionalInfo{{AdditionalIdToken: "0000", Type: "someType"}}}, CertificateHashData: []types.OCSPRequestDataType{{HashAlgorithm: types.SHA256, IssuerNameHash: "hash0", IssuerKeyHash: "hash1"}}}, false},
		{authorization.AuthorizeRequest{Certificate: "deadc0de", IdToken: types.IdToken{Type: types.IdTokenTypeKeyCode, AdditionalInfo: []types.AdditionalInfo{{AdditionalIdToken: "0000", Type: "someType"}}}, CertificateHashData: []types.OCSPRequestDataType{{SerialNumber: "s0", HashAlgorithm: types.SHA256, IssuerNameHash: "h0", IssuerKeyHash: "h0.1"}, {SerialNumber: "s1", HashAlgorithm: types.SHA256, IssuerNameHash: "h1", IssuerKeyHash: "h1.1"}, {SerialNumber: "s2", HashAlgorithm: types.SHA256, IssuerNameHash: "h2", IssuerKeyHash: "h2.1"}, {SerialNumber: "s3", HashAlgorithm: types.SHA256, IssuerNameHash: "h3", IssuerKeyHash: "h3.1"}, {SerialNumber: "s4", HashAlgorithm: types.SHA256, IssuerNameHash: "h4", IssuerKeyHash: "h4.1"}}}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV2TestSuite) TestAuthorizeConfirmationValidation() {
	var confirmationTable = []GenericTestEntry{
		{authorization.AuthorizeResponse{CertificateStatus: authorization.CertificateStatusAccepted, IdTokenInfo: types.IdTokenInfo{Status: types.AuthorizationStatusAccepted}}, true},
		{authorization.AuthorizeResponse{CertificateStatus: authorization.CertificateStatusAccepted, IdTokenInfo: types.IdTokenInfo{Status: types.AuthorizationStatusAccepted}}, true},
		{authorization.AuthorizeResponse{IdTokenInfo: types.IdTokenInfo{Status: types.AuthorizationStatusAccepted}}, true},
		{authorization.AuthorizeResponse{}, false},
		{authorization.AuthorizeResponse{CertificateStatus: "invalidCertificateStatus", IdTokenInfo: types.IdTokenInfo{Status: types.AuthorizationStatusAccepted}}, false},
		{authorization.AuthorizeResponse{CertificateStatus: authorization.CertificateStatusAccepted, IdTokenInfo: types.IdTokenInfo{Status: "invalidTokenInfoStatus"}}, false},
	}
	ExecuteGenericTestTable(suite, confirmationTable)
}

func (suite *OcppV2TestSuite) TestAuthorizeE2EMocked() {
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	certificate := "deadc0de"
	additionalInfo := types.AdditionalInfo{AdditionalIdToken: "at1", Type: "some"}
	idToken := types.IdToken{IdToken: "tok1", Type: types.IdTokenTypeKeyCode, AdditionalInfo: []types.AdditionalInfo{additionalInfo}}
	certHashData := types.OCSPRequestDataType{HashAlgorithm: types.SHA256, IssuerNameHash: "h0", IssuerKeyHash: "h0.1", SerialNumber: "s0", ResponderURL: "http://www.test.org"}
	status := types.AuthorizationStatusAccepted
	certificateStatus := authorization.CertificateStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"certificate":"%v","idToken":{"idToken":"%v","type":"%v","additionalInfo":[{"additionalIdToken":"%v","type":"%v"}]},"iso15118CertificateHashData":[{"hashAlgorithm":"%v","issuerNameHash":"%v","issuerKeyHash":"%v","serialNumber":"%v","responderURL":"%v"}]}]`,
		messageId, authorization.AuthorizeFeatureName, certificate, idToken.IdToken, idToken.Type, additionalInfo.AdditionalIdToken, additionalInfo.Type, certHashData.HashAlgorithm, certHashData.IssuerNameHash, certHashData.IssuerKeyHash, certHashData.SerialNumber, certHashData.ResponderURL)
	responseJson := fmt.Sprintf(`[3,"%v",{"certificateStatus":"%v","idTokenInfo":{"status":"%v"}}]`,
		messageId, certificateStatus, status)
	authorizeConfirmation := authorization.NewAuthorizationResponse(types.IdTokenInfo{Status: status})
	authorizeConfirmation.CertificateStatus = certificateStatus
	requestRaw := []byte(requestJson)
	responseRaw := []byte(responseJson)
	channel := NewMockWebSocket(wsId)

	handler := &MockCSMSAuthorizationHandler{}
	handler.On("OnAuthorize", mock.AnythingOfType("string"), mock.Anything).Return(authorizeConfirmation, nil).Run(func(args mock.Arguments) {
		request := args.Get(1).(*authorization.AuthorizeRequest)
		suite.Equal(certificate, request.Certificate)
		suite.Equal(idToken.IdToken, request.IdToken.IdToken)
		suite.Equal(idToken.Type, request.IdToken.Type)
		suite.Require().Len(request.IdToken.AdditionalInfo, 1)
		suite.Equal(idToken.AdditionalInfo[0].AdditionalIdToken, request.IdToken.AdditionalInfo[0].AdditionalIdToken)
		suite.Equal(idToken.AdditionalInfo[0].Type, request.IdToken.AdditionalInfo[0].Type)
		suite.Require().Len(request.CertificateHashData, 1)
		suite.Equal(certHashData.HashAlgorithm, request.CertificateHashData[0].HashAlgorithm)
		suite.Equal(certHashData.IssuerNameHash, request.CertificateHashData[0].IssuerNameHash)
		suite.Equal(certHashData.IssuerKeyHash, request.CertificateHashData[0].IssuerKeyHash)
		suite.Equal(certHashData.SerialNumber, request.CertificateHashData[0].SerialNumber)
		suite.Equal(certHashData.ResponderURL, request.CertificateHashData[0].ResponderURL)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: responseRaw, forwardWrittenMessage: true}, handler)
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: requestRaw, forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	suite.Require().Nil(err)
	response, err := suite.chargingStation.Authorize(idToken.IdToken, idToken.Type, func(request *authorization.AuthorizeRequest) {
		request.IdToken.AdditionalInfo = []types.AdditionalInfo{additionalInfo}
		request.Certificate = certificate
		request.CertificateHashData = []types.OCSPRequestDataType{certHashData}
	})
	suite.Require().Nil(err)
	suite.Require().NotNil(response)
	suite.Equal(certificateStatus, response.CertificateStatus)
	suite.Equal(status, response.IdTokenInfo.Status)
}

func (suite *OcppV2TestSuite) TestAuthorizeInvalidEndpoint() {
	messageId := defaultMessageId
	certificate := "deadc0de"
	additionalInfo := types.AdditionalInfo{AdditionalIdToken: "at1", Type: "some"}
	idToken := types.IdToken{IdToken: "tok1", Type: types.IdTokenTypeKeyCode, AdditionalInfo: []types.AdditionalInfo{additionalInfo}}
	certHashData := types.OCSPRequestDataType{HashAlgorithm: types.SHA256, IssuerNameHash: "h0", IssuerKeyHash: "h0.1", SerialNumber: "s0", ResponderURL: "http://www.test.org"}
	authorizeRequest := authorization.NewAuthorizationRequest(idToken.IdToken, idToken.Type)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"certificate":"%v","idToken":{"idToken":"%v","type":"%v","additionalInfo":[{"additionalIdToken":"%v","type":"%v"}]},"iso15118CertificateHashData":[{"hashAlgorithm":"%v","issuerNameHash":"%v","issuerKeyHash":"%v","serialNumber":"%v","responderURL":"%v"}]}]`,
		messageId, authorization.AuthorizeFeatureName, certificate, idToken.IdToken, idToken.Type, additionalInfo.AdditionalIdToken, additionalInfo.Type, certHashData.HashAlgorithm, certHashData.IssuerNameHash, certHashData.IssuerKeyHash, certHashData.SerialNumber, certHashData.ResponderURL)
	testUnsupportedRequestFromCentralSystem(suite, authorizeRequest, requestJson, messageId)
}
