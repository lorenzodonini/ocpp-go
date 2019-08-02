package ocpp16_test

import (
	"fmt"
	"github.com/lorenzodonini/go-ocpp/ocpp1.6"
	"github.com/lorenzodonini/go-ocpp/ocppj"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

// Utility functions
func getAuthorizeRequest(t *testing.T, request ocppj.Request) *ocpp16.AuthorizeRequest {
	assert.NotNil(t, request)
	result, ok := request.(*ocpp16.AuthorizeRequest)
	assert.True(t, ok)
	assert.NotNil(t, result)
	assert.IsType(t, &ocpp16.AuthorizeRequest{}, result)
	return result
}

func getAuthorizeConfirmation(t *testing.T, confirmation ocppj.Confirmation) *ocpp16.AuthorizeConfirmation {
	assert.NotNil(t, confirmation)
	result, ok := confirmation.(*ocpp16.AuthorizeConfirmation)
	assert.True(t, ok)
	assert.NotNil(t, result)
	assert.IsType(t, &ocpp16.AuthorizeConfirmation{}, result)
	return result
}

// Test
func (suite *OcppV16TestSuite) TestAuthorizeRequestValidation() {
	t := suite.T()
	var requestTable = []RequestTestEntry{
		{ocpp16.AuthorizeRequest{IdTag: "12345"}, true},
		{ocpp16.AuthorizeRequest{}, false},
		{ocpp16.AuthorizeRequest{IdTag: ">20.................."}, false},
	}
	ExecuteRequestTestTable(t, requestTable)
}

func (suite *OcppV16TestSuite) TestAuthorizeConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []ConfirmationTestEntry{
		{ocpp16.AuthorizeConfirmation{IdTagInfo: ocpp16.IdTagInfo{ExpiryDate: ocpp16.DateTime{Time: time.Now().Add(time.Hour * 8)}, ParentIdTag: "00000", Status: ocpp16.AuthorizationStatusAccepted}}, true},
		{ocpp16.AuthorizeConfirmation{IdTagInfo: ocpp16.IdTagInfo{ParentIdTag: "00000", Status: ocpp16.AuthorizationStatusAccepted}}, true},
		{ocpp16.AuthorizeConfirmation{IdTagInfo: ocpp16.IdTagInfo{ExpiryDate: ocpp16.DateTime{Time: time.Now().Add(time.Hour * 8)}, Status: ocpp16.AuthorizationStatusAccepted}}, true},
		{ocpp16.AuthorizeConfirmation{IdTagInfo: ocpp16.IdTagInfo{Status: ocpp16.AuthorizationStatusAccepted}}, true},
		{ocpp16.AuthorizeConfirmation{IdTagInfo: ocpp16.IdTagInfo{Status: ocpp16.AuthorizationStatusBlocked}}, true},
		{ocpp16.AuthorizeConfirmation{IdTagInfo: ocpp16.IdTagInfo{Status: ocpp16.AuthorizationStatusExpired}}, true},
		{ocpp16.AuthorizeConfirmation{IdTagInfo: ocpp16.IdTagInfo{Status: ocpp16.AuthorizationStatusInvalid}}, true},
		{ocpp16.AuthorizeConfirmation{IdTagInfo: ocpp16.IdTagInfo{Status: ocpp16.AuthorizationStatusConcurrentTx}}, true},
		{ocpp16.AuthorizeConfirmation{IdTagInfo: ocpp16.IdTagInfo{ParentIdTag: ">20..................", Status: ocpp16.AuthorizationStatusAccepted}}, false},
		{ocpp16.AuthorizeConfirmation{IdTagInfo: ocpp16.IdTagInfo{ExpiryDate: ocpp16.DateTime{Time: time.Now().Add(time.Hour * -8)}, Status: ocpp16.AuthorizationStatusAccepted}}, false},
	}
	ExecuteConfirmationTestTable(t, confirmationTable)
}

func (suite *OcppV16TestSuite) TestAuthorizeE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	idTag := "tag1"
	parentIdTag := "parentTag1"
	status := ocpp16.AuthorizationStatusAccepted
	expiryDate := ocpp16.DateTime{Time: time.Now().Add(time.Hour * 8)}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"idTag":"%v"}]`, messageId, ocpp16.AuthorizeFeatureName, idTag)
	responseJson := fmt.Sprintf(`[3,"%v",{"idTagInfo":{"expiryDate":"%v","parentIdTag":"%v","status":"%v"}}]`, messageId, expiryDate.Time.Format(ocpp16.ISO8601), parentIdTag, status)
	authorizeConfirmation := ocpp16.NewAuthorizationConfirmation(ocpp16.IdTagInfo{ExpiryDate: expiryDate, ParentIdTag: parentIdTag, Status: status})
	requestRaw := []byte(requestJson)
	responseRaw := []byte(responseJson)
	channel := NewMockWebSocket(wsId)

	coreListener := MockCentralSystemCoreListener{}
	coreListener.On("OnAuthorize", mock.AnythingOfType("string"), mock.Anything).Return(authorizeConfirmation, nil)
	setupDefaultCentralSystemHandlers(suite, coreListener, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: responseRaw, forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: requestRaw, forwardWrittenMessage: true})
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	assert.Nil(t, err)
	confirmation, protoErr, err := suite.chargePoint.Authorize(idTag)
	assert.Nil(t, err)
	assert.Nil(t, protoErr)
	assert.NotNil(t, confirmation)
	assert.Equal(t, status, confirmation.IdTagInfo.Status)
	assert.Equal(t, parentIdTag, confirmation.IdTagInfo.ParentIdTag)
	assertDateTimeEquality(t, expiryDate, confirmation.IdTagInfo.ExpiryDate)
}

func (suite *OcppV16TestSuite) TestAuthorizeInvalidEndpoint() {
	t := suite.T()
	wsId := "test_id"
	messageId := "1234"
	idTag := "tag1"
	expectedError := fmt.Sprintf("unsupported action %v on central system, cannot send request", ocpp16.AuthorizeFeatureName)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"idTag":"%v"}]`, messageId, ocpp16.AuthorizeFeatureName, idTag)

	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: false})
	// Run test
	authorizeRequest := ocpp16.NewAuthorizationRequest(idTag)
	suite.centralSystem.Start(8887, "somePath")
	err := suite.centralSystem.SendRequestAsync(wsId, authorizeRequest, func(confirmation ocppj.Confirmation, callError *ocppj.ProtoError) {
		t.Fail()
	})
	assert.Error(t, err)
	assert.Equal(t, expectedError, err.Error())
}

func (suite *OcppV16TestSuite) TestAuthorizeInvalidEndpointResponse() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	idTag := "tag1"
	errorDescription := fmt.Sprintf("unsupported action %v on charge point", ocpp16.AuthorizeFeatureName)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"idTag":"%v"}]`, messageId, ocpp16.AuthorizeFeatureName, idTag)
	errorJson := fmt.Sprintf(`[4,"%v","%v","%v",null]`, messageId, ocppj.NotSupported, errorDescription)
	channel := NewMockWebSocket(wsId)

	coreListener := MockChargePointCoreListener{}
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: false})
	setupDefaultChargePointHandlers(suite, coreListener, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(errorJson), forwardWrittenMessage: true})
	suite.ocppjCentralSystem.SetErrorHandler(func(chargePointId string, errorCode ocppj.ErrorCode, description string, details interface{}, requestId string) {
		assert.Equal(t, messageId, requestId)
		assert.Equal(t, wsId, chargePointId)
		assert.Equal(t, ocppj.NotSupported, errorCode)
		assert.Equal(t, errorDescription, description)
		assert.Nil(t, details)
	})
	// Mock pending request
	pendingRequest := ocpp16.NewAuthorizationRequest(idTag)
	suite.ocppjCentralSystem.AddPendingRequest(messageId, pendingRequest)
	// Run test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	assert.Nil(t, err)
	err = suite.mockWsClient.MessageHandler([]byte(requestJson))
	assert.Nil(t, err)
}
