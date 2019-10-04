package ocpp16_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"time"
)

// Test
func (suite *OcppV16TestSuite) TestAuthorizeRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{ocpp16.AuthorizeRequest{IdTag: "12345"}, true},
		{ocpp16.AuthorizeRequest{}, false},
		{ocpp16.AuthorizeRequest{IdTag: ">20.................."}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV16TestSuite) TestAuthorizeConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{ocpp16.AuthorizeConfirmation{IdTagInfo: &ocpp16.IdTagInfo{ExpiryDate: ocpp16.NewDateTime(time.Now().Add(time.Hour * 8)), ParentIdTag: "00000", Status: ocpp16.AuthorizationStatusAccepted}}, true},
		{ocpp16.AuthorizeConfirmation{IdTagInfo: &ocpp16.IdTagInfo{Status: "invalidAuthorizationStatus"}}, false},
		{ocpp16.AuthorizeConfirmation{}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV16TestSuite) TestAuthorizeE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	idTag := "tag1"
	parentIdTag := "parentTag1"
	status := ocpp16.AuthorizationStatusAccepted
	expiryDate := ocpp16.NewDateTime(time.Now().Add(time.Hour * 8))
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"idTag":"%v"}]`, messageId, ocpp16.AuthorizeFeatureName, idTag)
	responseJson := fmt.Sprintf(`[3,"%v",{"idTagInfo":{"expiryDate":"%v","parentIdTag":"%v","status":"%v"}}]`, messageId, expiryDate.Time.Format(ocpp16.ISO8601), parentIdTag, status)
	authorizeConfirmation := ocpp16.NewAuthorizationConfirmation(&ocpp16.IdTagInfo{ExpiryDate: expiryDate, ParentIdTag: parentIdTag, Status: status})
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
	confirmation, err := suite.chargePoint.Authorize(idTag)
	assert.Nil(t, err)
	assert.NotNil(t, confirmation)
	assert.Equal(t, status, confirmation.IdTagInfo.Status)
	assert.Equal(t, parentIdTag, confirmation.IdTagInfo.ParentIdTag)
	assertDateTimeEquality(t, *expiryDate, *confirmation.IdTagInfo.ExpiryDate)
}

func (suite *OcppV16TestSuite) TestAuthorizeInvalidEndpoint() {
	messageId := defaultMessageId
	idTag := "tag1"
	authorizeRequest := ocpp16.NewAuthorizationRequest(idTag)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"idTag":"%v"}]`, messageId, ocpp16.AuthorizeFeatureName, idTag)
	testUnsupportedRequestFromCentralSystem(suite, authorizeRequest, requestJson, messageId)
}
