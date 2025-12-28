package ocpp2_test

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/authorization"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"

	"github.com/stretchr/testify/mock"
)

// Test
func (suite *OcppV2TestSuite) TestClearCacheRequestValidation() {
	var requestTable = []GenericTestEntry{
		{authorization.ClearCacheRequest{}, true},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV2TestSuite) TestClearCacheConfirmationValidation() {
	var confirmationTable = []GenericTestEntry{
		{authorization.ClearCacheResponse{Status: authorization.ClearCacheStatusAccepted, StatusInfo: types.NewStatusInfo("200", "ok")}, true},
		{authorization.ClearCacheResponse{Status: authorization.ClearCacheStatusAccepted}, true},
		{authorization.ClearCacheResponse{Status: authorization.ClearCacheStatusRejected}, true},
		{authorization.ClearCacheResponse{Status: "invalidClearCacheStatus"}, false},
		{authorization.ClearCacheResponse{}, false},
	}
	ExecuteGenericTestTable(suite, confirmationTable)
}

func (suite *OcppV2TestSuite) TestClearCacheE2EMocked() {
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	status := authorization.ClearCacheStatusAccepted
	statusInfo := types.NewStatusInfo("200", "")
	requestJson := fmt.Sprintf(`[2,"%v","%v",{}]`, messageId, authorization.ClearCacheFeatureName)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v","statusInfo":{"reasonCode":"%v"}}]`, messageId, status, statusInfo.ReasonCode)
	clearCacheResponse := authorization.NewClearCacheResponse(status)
	clearCacheResponse.StatusInfo = statusInfo
	channel := NewMockWebSocket(wsId)

	handler := &MockChargingStationAuthorizationHandler{}
	handler.On("OnClearCache", mock.Anything).Return(clearCacheResponse, nil)
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	suite.Require().Nil(err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.ClearCache(wsId, func(confirmation *authorization.ClearCacheResponse, err error) {
		suite.Require().Nil(err)
		suite.Require().NotNil(confirmation)
		suite.Equal(status, confirmation.Status)
		resultChannel <- true
	})
	suite.Require().Nil(err)
	result := <-resultChannel
	suite.True(result)
}

func (suite *OcppV2TestSuite) TestClearCacheInvalidEndpoint() {
	messageId := defaultMessageId
	clearCacheRequest := authorization.NewClearCacheRequest()
	requestJson := fmt.Sprintf(`[2,"%v","%v",{}]`, messageId, authorization.ClearCacheFeatureName)
	testUnsupportedRequestFromChargingStation(suite, clearCacheRequest, requestJson, messageId)
}
