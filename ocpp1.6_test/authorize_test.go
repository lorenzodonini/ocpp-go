package ocpp16_test

import (
	"errors"
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

func (suite *OcppV16TestSuite) TestAuthorizeRequestFromJson() {
	t := suite.T()
	uniqueId := defaultMessageId
	idTag := "tag1"
	wsId := "test_id"
	channel := NewMockWebSocket(wsId)
	dataJson := fmt.Sprintf(`[2,"%v","Authorize",{"idTag":"%v"}]`, uniqueId, idTag)

	coreListener := MockCentralSystemCoreListener{}
	coreListener.On("OnAuthorize", mock.AnythingOfType("string"), mock.Anything).Return(nil, nil).Run(func(args mock.Arguments) {
		clientId := args.String(0)
		request := args.Get(1).(*ocpp16.AuthorizeRequest)
		assert.Equal(t, wsId, clientId)
		assert.Equal(t, idTag, request.IdTag)
	})
	setupDefaultCentralSystemHandlers(suite, coreListener, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: nil})
	suite.centralSystem.Start(8887, "somePath")
	suite.mockWsServer.NewClientHandler(channel)
	err := suite.mockWsServer.MessageHandler(channel, []byte(dataJson))
	assert.Nil(t, err)
}

func (suite *OcppV16TestSuite) TestAuthorizeRequestToJson() {
	t := suite.T()
	idTag := "tag2"
	wsUrl := "someUrl"
	wsId := "test_id"
	mockError := "mock"
	expectedJson := fmt.Sprintf(`[2,"%v","Authorize",{"idTag":"%v"}]`, defaultMessageId, idTag)

	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, rawWrittenMessage: []byte(expectedJson), writeReturnArgument: errors.New(mockError)})
	err := suite.chargePoint.Start(wsUrl)
	assert.Nil(t, err)
	confirmation, protoErr, err := suite.chargePoint.Authorize(idTag)
	assert.Nil(t, confirmation)
	assert.Nil(t, protoErr)
	assert.NotNil(t, err)
	assert.Equal(t, mockError, err.Error())
}

func (suite *OcppV16TestSuite) TestAuthorizeConfirmationFromJson() {
	t := suite.T()
	wsUrl := "someUrl"
	wsId := "test_id"
	uniqueId := defaultMessageId
	status := ocpp16.AuthorizationStatusAccepted
	parentIdTag := "parentTag1"
	expiryDate := ocpp16.DateTime{Time: time.Now().Add(time.Hour * 8)}
	expectedJson := fmt.Sprintf(`[3,"%v",{"idTagInfo":{"expiryDate":"%v","parentIdTag":"%v","status":"%v"}}]`, uniqueId, expiryDate.Format(ocpp16.ISO8601), parentIdTag, status)

	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: false})
	suite.ocppjChargePoint.SetConfirmationHandler(func(confirmation ocppj.Confirmation, requestId string) {
		authorizeConfirmation := getAuthorizeConfirmation(t, confirmation)
		assert.NotNil(t, authorizeConfirmation)
		assert.Equal(t, status, authorizeConfirmation.IdTagInfo.Status)
		assert.Equal(t, parentIdTag, authorizeConfirmation.IdTagInfo.ParentIdTag)
		assertDateTimeEquality(t, expiryDate, authorizeConfirmation.IdTagInfo.ExpiryDate)
	})
	// Mock pending request
	dummyRequest := ocpp16.AuthorizeRequest{}
	suite.ocppjChargePoint.Endpoint.AddPendingRequest(uniqueId, dummyRequest)
	// Run test
	err := suite.chargePoint.Start(wsUrl)
	assert.Nil(t, err)
	err = suite.mockWsClient.MessageHandler([]byte(expectedJson))
	assert.Nil(t, err)
}

func (suite *OcppV16TestSuite) TestAuthorizeConfirmationToJson() {
	t := suite.T()
	uniqueId := "1234"
	idTag := "tag1"
	parentIdTag := "parentTag1"
	wsId := "test_id"
	channel := NewMockWebSocket(wsId)
	expiryDate := ocpp16.DateTime{Time: time.Now().Add(time.Hour * 8)}
	status := ocpp16.AuthorizationStatusAccepted
	dummyRequest := fmt.Sprintf(`[2,"%v","Authorize",{"idTag":"%v"}]`, uniqueId, idTag)
	confirmation := ocpp16.AuthorizeConfirmation{IdTagInfo: ocpp16.IdTagInfo{Status: status, ParentIdTag: parentIdTag, ExpiryDate: expiryDate}}
	expectedJson := fmt.Sprintf(`[3,"%v",{"idTagInfo":{"expiryDate":"%v","parentIdTag":"%v","status":"%v"}}]`, uniqueId, expiryDate.Format(ocpp16.ISO8601), parentIdTag, status)

	coreListener := MockCentralSystemCoreListener{}
	coreListener.On("OnAuthorize", mock.AnythingOfType("string"), mock.Anything).Return(confirmation, nil)
	suite.centralSystem.SetCentralSystemCoreListener(coreListener)
	setupDefaultCentralSystemHandlers(suite, coreListener, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(expectedJson), forwardWrittenMessage: false})

	suite.centralSystem.Start(8887, "somePath")
	suite.mockWsServer.NewClientHandler(channel)
	err := suite.mockWsServer.MessageHandler(channel, []byte(dummyRequest))
	assert.Nil(t, err)
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
	// CS -> CP
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

// TODO: test invalid direction
