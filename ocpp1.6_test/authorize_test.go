package ocpp16_test

import (
	"fmt"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test
func (suite *OcppV16TestSuite) TestAuthorizeRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{core.AuthorizeRequest{IdTag: "12345"}, true},
		{core.AuthorizeRequest{}, false},
		{core.AuthorizeRequest{IdTag: ">20.................."}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV16TestSuite) TestAuthorizeConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{core.AuthorizeConfirmation{IdTagInfo: &types.IdTagInfo{ExpiryDate: types.NewDateTime(time.Now().Add(time.Hour * 8)), ParentIdTag: "00000", Status: types.AuthorizationStatusAccepted}}, true},
		{core.AuthorizeConfirmation{IdTagInfo: &types.IdTagInfo{Status: "invalidAuthorizationStatus"}}, false},
		{core.AuthorizeConfirmation{}, false},
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
	status := types.AuthorizationStatusAccepted
	expiryDate := types.NewDateTime(time.Now().Add(time.Hour * 8))
	idTagInfo := types.IdTagInfo{ExpiryDate: expiryDate, ParentIdTag: parentIdTag, Status: status}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"idTag":"%v"}]`, messageId, core.AuthorizeFeatureName, idTag)
	responseJson := fmt.Sprintf(`[3,"%v",{"idTagInfo":{"expiryDate":"%v","parentIdTag":"%v","status":"%v"}}]`, messageId, expiryDate.FormatTimestamp(), parentIdTag, status)
	authorizeConfirmation := core.NewAuthorizationConfirmation(&idTagInfo)
	requestRaw := []byte(requestJson)
	responseRaw := []byte(responseJson)
	channel := NewMockWebSocket(wsId)

	coreListener := &MockCentralSystemCoreListener{}
	coreListener.On("OnAuthorize", mock.AnythingOfType("string"), mock.Anything).Return(authorizeConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*core.AuthorizeRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, idTag, request.IdTag)
	})
	setupDefaultCentralSystemHandlers(suite, coreListener, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: responseRaw, forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: requestRaw, forwardWrittenMessage: true})
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	require.Nil(t, err)
	confirmation, err := suite.chargePoint.Authorize(idTag)
	require.Nil(t, err)
	require.NotNil(t, confirmation)
	assert.Equal(t, status, confirmation.IdTagInfo.Status)
	assert.Equal(t, parentIdTag, confirmation.IdTagInfo.ParentIdTag)
	assertDateTimeEquality(t, *expiryDate, *confirmation.IdTagInfo.ExpiryDate)
}

func (suite *OcppV16TestSuite) TestAuthorizeInvalidEndpoint() {
	messageId := defaultMessageId
	idTag := "tag1"
	authorizeRequest := core.NewAuthorizationRequest(idTag)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"idTag":"%v"}]`, messageId, core.AuthorizeFeatureName, idTag)
	testUnsupportedRequestFromCentralSystem(suite, authorizeRequest, requestJson, messageId)
}
