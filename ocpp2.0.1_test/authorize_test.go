package ocpp2_test

import (
	"fmt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/authorization"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

func (suite *OcppV2TestSuite) TestAuthorizeE2EMocked() {
	t := suite.T()
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
		assert.Equal(t, certificate, request.Certificate)
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
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: responseRaw, forwardWrittenMessage: true}, handler)
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: requestRaw, forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	response, err := suite.chargingStation.Authorize(idToken.IdToken, idToken.Type, func(request *authorization.AuthorizeRequest) {
		request.IdToken.AdditionalInfo = []types.AdditionalInfo{additionalInfo}
		request.Certificate = certificate
		request.CertificateHashData = []types.OCSPRequestDataType{certHashData}
	})
	require.Nil(t, err)
	require.NotNil(t, response)
	assert.Equal(t, certificateStatus, response.CertificateStatus)
	assert.Equal(t, status, response.IdTokenInfo.Status)
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
