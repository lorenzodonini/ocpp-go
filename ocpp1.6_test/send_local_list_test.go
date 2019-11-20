package ocpp16_test

import (
	"fmt"
	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"time"
)

// Test
func (suite *OcppV16TestSuite) TestSendLocalListRequestValidation() {
	t := suite.T()
	localAuthEntry := ocpp16.AuthorizationData{IdTag: "12345", IdTagInfo: &ocpp16.IdTagInfo{
		ExpiryDate:  ocpp16.NewDateTime(time.Now().Add(time.Hour * 8)),
		ParentIdTag: "000000",
		Status:      ocpp16.AuthorizationStatusAccepted,
	}}
	invalidAuthEntry := ocpp16.AuthorizationData{IdTag: "12345", IdTagInfo: &ocpp16.IdTagInfo{
		ExpiryDate:  ocpp16.NewDateTime(time.Now().Add(time.Hour * 8)),
		ParentIdTag: "000000",
		Status:      "invalidStatus",
	}}
	var requestTable = []GenericTestEntry{
		{ocpp16.SendLocalListRequest{UpdateType: ocpp16.UpdateTypeDifferential, ListVersion: 1, LocalAuthorizationList: []ocpp16.AuthorizationData{localAuthEntry}}, true},
		{ocpp16.SendLocalListRequest{UpdateType: ocpp16.UpdateTypeDifferential, ListVersion: 1, LocalAuthorizationList: []ocpp16.AuthorizationData{}}, true},
		{ocpp16.SendLocalListRequest{UpdateType: ocpp16.UpdateTypeDifferential, ListVersion: 1}, true},
		{ocpp16.SendLocalListRequest{UpdateType: ocpp16.UpdateTypeDifferential, ListVersion: 0}, true},
		{ocpp16.SendLocalListRequest{UpdateType: ocpp16.UpdateTypeDifferential}, true},
		{ocpp16.SendLocalListRequest{UpdateType: ocpp16.UpdateTypeDifferential, ListVersion: -1}, false},
		{ocpp16.SendLocalListRequest{UpdateType: ocpp16.UpdateTypeDifferential, ListVersion: 1, LocalAuthorizationList: []ocpp16.AuthorizationData{invalidAuthEntry}}, false},
		{ocpp16.SendLocalListRequest{UpdateType: "invalidUpdateType", ListVersion: 1}, false},
		{ocpp16.SendLocalListRequest{ListVersion: 1}, false},
		{ocpp16.SendLocalListRequest{}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV16TestSuite) TestSendLocalListConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{ocpp16.SendLocalListConfirmation{Status: ocpp16.UpdateStatusAccepted}, true},
		{ocpp16.SendLocalListConfirmation{Status: "invalidStatus"}, false},
		{ocpp16.SendLocalListConfirmation{}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV16TestSuite) TestSendLocalListE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	listVersion := 1
	updateType := ocpp16.UpdateTypeDifferential
	status := ocpp16.UpdateStatusAccepted
	localAuthEntry := ocpp16.AuthorizationData{IdTag: "12345", IdTagInfo: &ocpp16.IdTagInfo{
		ExpiryDate:  ocpp16.NewDateTime(time.Now().Add(time.Hour * 8)),
		ParentIdTag: "parentTag",
		Status:      ocpp16.AuthorizationStatusAccepted,
	}}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"listVersion":%v,"localAuthorizationList":[{"idTag":"%v","idTagInfo":{"expiryDate":"%v","parentIdTag":"%v","status":"%v"}}],"updateType":"%v"}]`,
		messageId, ocpp16.SendLocalListFeatureName, listVersion, localAuthEntry.IdTag, localAuthEntry.IdTagInfo.ExpiryDate.Format(ocpp16.ISO8601), localAuthEntry.IdTagInfo.ParentIdTag, localAuthEntry.IdTagInfo.Status, updateType)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	sendLocalListConfirmation := ocpp16.NewSendLocalListConfirmation(status)
	channel := NewMockWebSocket(wsId)

	localAuthListListener := MockChargePointLocalAuthListListener{}
	localAuthListListener.On("OnSendLocalList", mock.Anything).Return(sendLocalListConfirmation, nil)
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	suite.chargePoint.SetLocalAuthListListener(localAuthListListener)
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	assert.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.centralSystem.SendLocalList(wsId, func(confirmation *ocpp16.SendLocalListConfirmation, err error) {
		assert.Nil(t, err)
		assert.NotNil(t, confirmation)
		if confirmation != nil {
			assert.Equal(t, status, confirmation.Status)
			resultChannel <- true
		} else {
			resultChannel <- false
		}
	}, listVersion, updateType, func(request *ocpp16.SendLocalListRequest) {
		request.LocalAuthorizationList = []ocpp16.AuthorizationData{localAuthEntry}
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
	updateType := ocpp16.UpdateTypeDifferential
	localAuthEntry := ocpp16.AuthorizationData{IdTag: "12345", IdTagInfo: &ocpp16.IdTagInfo{
		ExpiryDate:  ocpp16.NewDateTime(time.Now().Add(time.Hour * 8)),
		ParentIdTag: "parentTag",
		Status:      ocpp16.AuthorizationStatusAccepted,
	}}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"listVersion":%v,"localAuthorizationList":[{"idTag":"%v","idTagInfo":{"expiryDate":"%v","parentIdTag":"%v","status":"%v"}}],"updateType":"%v"}]`,
		messageId, ocpp16.SendLocalListFeatureName, listVersion, localAuthEntry.IdTag, localAuthEntry.IdTagInfo.ExpiryDate.Format(ocpp16.ISO8601), localAuthEntry.IdTagInfo.ParentIdTag, localAuthEntry.IdTagInfo.Status, updateType)
	localListVersionRequest := ocpp16.NewSendLocalListRequest(listVersion, updateType)
	testUnsupportedRequestFromChargePoint(suite, localListVersionRequest, requestJson, messageId)
}
