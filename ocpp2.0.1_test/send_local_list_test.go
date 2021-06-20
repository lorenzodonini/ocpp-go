package ocpp2_test

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test
func (suite *OcppV2TestSuite) TestSendLocalListRequestValidation() {
	t := suite.T()
	authData := localauth.AuthorizationData{
		IdToken: types.IdToken{
			IdToken:        "token1",
			Type:           types.IdTokenTypeKeyCode,
			AdditionalInfo: nil,
		},
		IdTokenInfo: types.NewIdTokenInfo(types.AuthorizationStatusAccepted),
	}
	var requestTable = []GenericTestEntry{
		{localauth.SendLocalListRequest{VersionNumber: 42, UpdateType: localauth.UpdateTypeDifferential, LocalAuthorizationList: []localauth.AuthorizationData{authData}}, true},
		{localauth.SendLocalListRequest{VersionNumber: 42, UpdateType: localauth.UpdateTypeFull, LocalAuthorizationList: []localauth.AuthorizationData{authData}}, true},
		{localauth.SendLocalListRequest{VersionNumber: 42, UpdateType: localauth.UpdateTypeDifferential, LocalAuthorizationList: []localauth.AuthorizationData{}}, true},
		{localauth.SendLocalListRequest{VersionNumber: 42, UpdateType: localauth.UpdateTypeDifferential}, true},
		{localauth.SendLocalListRequest{UpdateType: localauth.UpdateTypeDifferential}, true},
		{localauth.SendLocalListRequest{}, false},
		{localauth.SendLocalListRequest{VersionNumber: -1, UpdateType: localauth.UpdateTypeDifferential, LocalAuthorizationList: []localauth.AuthorizationData{authData}}, false},
		{localauth.SendLocalListRequest{VersionNumber: 42, UpdateType: "invalidUpdateType", LocalAuthorizationList: []localauth.AuthorizationData{{IdToken: types.IdToken{IdToken: "tokenWithoutType"}}}}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestSendLocalListResponseValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{localauth.SendLocalListResponse{Status: localauth.SendLocalListStatusAccepted, StatusInfo: types.NewStatusInfo("200", "")}, true},
		{localauth.SendLocalListResponse{Status: localauth.SendLocalListStatusAccepted}, true},
		{localauth.SendLocalListResponse{}, false},
		{localauth.SendLocalListResponse{Status: "invalidStatus", StatusInfo: types.NewStatusInfo("200", "")}, false},
		{localauth.SendLocalListResponse{Status: localauth.SendLocalListStatusAccepted, StatusInfo: types.NewStatusInfo("", "")}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestSendLocalListE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	versionNumber := 1
	updateType := localauth.UpdateTypeDifferential
	authData := localauth.AuthorizationData{
		IdToken: types.IdToken{
			IdToken:        "token1",
			Type:           types.IdTokenTypeKeyCode,
			AdditionalInfo: nil,
		},
		IdTokenInfo: types.NewIdTokenInfo(types.AuthorizationStatusAccepted),
	}
	status := localauth.SendLocalListStatusAccepted
	statusInfo := types.NewStatusInfo("200", "")
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"versionNumber":%v,"updateType":"%v","localAuthorizationList":[{"idTokenInfo":{"status":"%v"},"idToken":{"idToken":"%v","type":"%v"}}]}]`,
		messageId, localauth.SendLocalListFeatureName, versionNumber, updateType, authData.IdTokenInfo.Status, authData.IdToken.IdToken, authData.IdToken.Type)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v","statusInfo":{"reasonCode":"%v"}}]`,
		messageId, status, statusInfo.ReasonCode)
	sendLocalListResponse := localauth.NewSendLocalListResponse(status)
	sendLocalListResponse.StatusInfo = statusInfo
	channel := NewMockWebSocket(wsId)

	handler := &MockChargingStationLocalAuthHandler{}
	handler.On("OnSendLocalList", mock.Anything).Return(sendLocalListResponse, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*localauth.SendLocalListRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, versionNumber, request.VersionNumber)
		assert.Equal(t, updateType, request.UpdateType)
		require.NotNil(t, request.LocalAuthorizationList)
		require.Len(t, request.LocalAuthorizationList, 1)
		assert.Equal(t, authData.IdToken.IdToken, request.LocalAuthorizationList[0].IdToken.IdToken)
		assert.Equal(t, authData.IdToken.Type, request.LocalAuthorizationList[0].IdToken.Type)
		require.NotNil(t, request.LocalAuthorizationList[0].IdTokenInfo)
		assert.Equal(t, authData.IdTokenInfo.Status, request.LocalAuthorizationList[0].IdTokenInfo.Status)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	assert.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.SendLocalList(wsId, func(response *localauth.SendLocalListResponse, err error) {
		assert.Nil(t, err)
		require.NotNil(t, response)
		assert.Equal(t, status, response.Status)
		require.NotNil(t, response.StatusInfo)
		assert.Equal(t, statusInfo.ReasonCode, response.StatusInfo.ReasonCode)
		assert.Equal(t, statusInfo.AdditionalInfo, response.StatusInfo.AdditionalInfo)
		resultChannel <- true
	}, versionNumber, updateType, func(request *localauth.SendLocalListRequest) {
		request.LocalAuthorizationList = []localauth.AuthorizationData{authData}
	})
	assert.Nil(t, err)
	if err == nil {
		result := <-resultChannel
		assert.True(t, result)
	}
}

func (suite *OcppV2TestSuite) TestSendLocalListInvalidEndpoint() {
	messageId := defaultMessageId
	versionNumber := 1
	updateType := localauth.UpdateTypeDifferential
	authData := localauth.AuthorizationData{
		IdToken: types.IdToken{
			IdToken:        "token1",
			Type:           types.IdTokenTypeKeyCode,
			AdditionalInfo: nil,
		},
		IdTokenInfo: types.NewIdTokenInfo(types.AuthorizationStatusAccepted),
	}
	localListVersionRequest := localauth.NewSendLocalListRequest(versionNumber, updateType)
	localListVersionRequest.LocalAuthorizationList = []localauth.AuthorizationData{authData}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"versionNumber":%v,"updateType":"%v","localAuthorizationList":[{"idTokenInfo":{"status":"%v"},"idToken":{"idToken":"%v","type":"%v"}}]}]`,
		messageId, localauth.SendLocalListFeatureName, versionNumber, updateType, authData.IdTokenInfo.Status, authData.IdToken.IdToken, authData.IdToken.Type)
	testUnsupportedRequestFromChargingStation(suite, localListVersionRequest, requestJson, messageId)
}
