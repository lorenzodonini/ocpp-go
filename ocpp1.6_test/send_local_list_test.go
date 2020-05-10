package ocpp16_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"time"
)

// Test
func (suite *OcppV16TestSuite) TestSendLocalListRequestValidation() {
	t := suite.T()
	localAuthEntry := localauth.AuthorizationData{IdTag: "12345", IdTagInfo: &types.IdTagInfo{
		ExpiryDate:  types.NewDateTime(time.Now().Add(time.Hour * 8)),
		ParentIdTag: "000000",
		Status:      types.AuthorizationStatusAccepted,
	}}
	invalidAuthEntry := localauth.AuthorizationData{IdTag: "12345", IdTagInfo: &types.IdTagInfo{
		ExpiryDate:  types.NewDateTime(time.Now().Add(time.Hour * 8)),
		ParentIdTag: "000000",
		Status:      "invalidStatus",
	}}
	var requestTable = []GenericTestEntry{
		{localauth.SendLocalListRequest{UpdateType: localauth.UpdateTypeDifferential, ListVersion: 1, LocalAuthorizationList: []localauth.AuthorizationData{localAuthEntry}}, true},
		{localauth.SendLocalListRequest{UpdateType: localauth.UpdateTypeDifferential, ListVersion: 1, LocalAuthorizationList: []localauth.AuthorizationData{}}, true},
		{localauth.SendLocalListRequest{UpdateType: localauth.UpdateTypeDifferential, ListVersion: 1}, true},
		{localauth.SendLocalListRequest{UpdateType: localauth.UpdateTypeDifferential, ListVersion: 0}, true},
		{localauth.SendLocalListRequest{UpdateType: localauth.UpdateTypeDifferential}, true},
		{localauth.SendLocalListRequest{UpdateType: localauth.UpdateTypeDifferential, ListVersion: -1}, false},
		{localauth.SendLocalListRequest{UpdateType: localauth.UpdateTypeDifferential, ListVersion: 1, LocalAuthorizationList: []localauth.AuthorizationData{invalidAuthEntry}}, false},
		{localauth.SendLocalListRequest{UpdateType: "invalidUpdateType", ListVersion: 1}, false},
		{localauth.SendLocalListRequest{ListVersion: 1}, false},
		{localauth.SendLocalListRequest{}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV16TestSuite) TestSendLocalListConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{localauth.SendLocalListConfirmation{Status: localauth.UpdateStatusAccepted}, true},
		{localauth.SendLocalListConfirmation{Status: "invalidStatus"}, false},
		{localauth.SendLocalListConfirmation{}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV16TestSuite) TestSendLocalListE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	listVersion := 1
	updateType := localauth.UpdateTypeDifferential
	status := localauth.UpdateStatusAccepted
	localAuthEntry := localauth.AuthorizationData{IdTag: "12345", IdTagInfo: &types.IdTagInfo{
		ExpiryDate:  types.NewDateTime(time.Now().Add(time.Hour * 8)),
		ParentIdTag: "parentTag",
		Status:      types.AuthorizationStatusAccepted,
	}}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"listVersion":%v,"localAuthorizationList":[{"idTag":"%v","idTagInfo":{"expiryDate":"%v","parentIdTag":"%v","status":"%v"}}],"updateType":"%v"}]`,
		messageId, localauth.SendLocalListFeatureName, listVersion, localAuthEntry.IdTag, localAuthEntry.IdTagInfo.ExpiryDate.FormatTimestamp(), localAuthEntry.IdTagInfo.ParentIdTag, localAuthEntry.IdTagInfo.Status, updateType)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	sendLocalListConfirmation := localauth.NewSendLocalListConfirmation(status)
	channel := NewMockWebSocket(wsId)

	localAuthListListener := MockChargePointLocalAuthListListener{}
	localAuthListListener.On("OnSendLocalList", mock.Anything).Return(sendLocalListConfirmation, nil)
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	suite.chargePoint.SetLocalAuthListHandler(localAuthListListener)
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	assert.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.centralSystem.SendLocalList(wsId, func(confirmation *localauth.SendLocalListConfirmation, err error) {
		assert.Nil(t, err)
		assert.NotNil(t, confirmation)
		if confirmation != nil {
			assert.Equal(t, status, confirmation.Status)
			resultChannel <- true
		} else {
			resultChannel <- false
		}
	}, listVersion, updateType, func(request *localauth.SendLocalListRequest) {
		request.LocalAuthorizationList = []localauth.AuthorizationData{localAuthEntry}
	})
	assert.Nil(t, err)
	if err == nil {
		result := <-resultChannel
		assert.True(t, result)
	}
}

func (suite *OcppV16TestSuite) TestSendLocalListInvalidEndpoint() {
	messageId := defaultMessageId
	listVersion := 1
	updateType := localauth.UpdateTypeDifferential
	localAuthEntry := localauth.AuthorizationData{IdTag: "12345", IdTagInfo: &types.IdTagInfo{
		ExpiryDate:  types.NewDateTime(time.Now().Add(time.Hour * 8)),
		ParentIdTag: "parentTag",
		Status:      types.AuthorizationStatusAccepted,
	}}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"listVersion":%v,"localAuthorizationList":[{"idTag":"%v","idTagInfo":{"expiryDate":"%v","parentIdTag":"%v","status":"%v"}}],"updateType":"%v"}]`,
		messageId, localauth.SendLocalListFeatureName, listVersion, localAuthEntry.IdTag, localAuthEntry.IdTagInfo.ExpiryDate.FormatTimestamp(), localAuthEntry.IdTagInfo.ParentIdTag, localAuthEntry.IdTagInfo.Status, updateType)
	localListVersionRequest := localauth.NewSendLocalListRequest(listVersion, updateType)
	testUnsupportedRequestFromChargePoint(suite, localListVersionRequest, requestJson, messageId)
}
