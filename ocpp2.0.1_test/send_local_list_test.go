package ocpp2_test

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"

	"github.com/stretchr/testify/mock"
)

// Test
func (suite *OcppV2TestSuite) TestSendLocalListRequestValidation() {
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
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV2TestSuite) TestSendLocalListResponseValidation() {
	var confirmationTable = []GenericTestEntry{
		{localauth.SendLocalListResponse{Status: localauth.SendLocalListStatusAccepted, StatusInfo: types.NewStatusInfo("200", "")}, true},
		{localauth.SendLocalListResponse{Status: localauth.SendLocalListStatusAccepted}, true},
		{localauth.SendLocalListResponse{}, false},
		{localauth.SendLocalListResponse{Status: "invalidStatus", StatusInfo: types.NewStatusInfo("200", "")}, false},
		{localauth.SendLocalListResponse{Status: localauth.SendLocalListStatusAccepted, StatusInfo: types.NewStatusInfo("", "")}, false},
	}
	ExecuteGenericTestTable(suite, confirmationTable)
}

func (suite *OcppV2TestSuite) TestSendLocalListE2EMocked() {
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
		suite.Require().True(ok)
		suite.Require().NotNil(request)
		suite.Equal(versionNumber, request.VersionNumber)
		suite.Equal(updateType, request.UpdateType)
		suite.Require().NotNil(request.LocalAuthorizationList)
		suite.Require().Len(request.LocalAuthorizationList, 1)
		suite.Equal(authData.IdToken.IdToken, request.LocalAuthorizationList[0].IdToken.IdToken)
		suite.Equal(authData.IdToken.Type, request.LocalAuthorizationList[0].IdToken.Type)
		suite.Require().NotNil(request.LocalAuthorizationList[0].IdTokenInfo)
		suite.Equal(authData.IdTokenInfo.Status, request.LocalAuthorizationList[0].IdTokenInfo.Status)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	suite.Nil(err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.SendLocalList(wsId, func(response *localauth.SendLocalListResponse, err error) {
		suite.Nil(err)
		suite.Require().NotNil(response)
		suite.Equal(status, response.Status)
		suite.Require().NotNil(response.StatusInfo)
		suite.Equal(statusInfo.ReasonCode, response.StatusInfo.ReasonCode)
		suite.Equal(statusInfo.AdditionalInfo, response.StatusInfo.AdditionalInfo)
		resultChannel <- true
	}, versionNumber, updateType, func(request *localauth.SendLocalListRequest) {
		request.LocalAuthorizationList = []localauth.AuthorizationData{authData}
	})
	suite.Nil(err)
	if err == nil {
		result := <-resultChannel
		suite.True(result)
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
