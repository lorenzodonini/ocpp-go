package ocpp2_test

import (
	"fmt"

	"github.com/stretchr/testify/mock"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/localauth"
)

// Test
func (suite *OcppV2TestSuite) TestGetLocalListVersionRequestValidation() {
	var requestTable = []GenericTestEntry{
		{localauth.GetLocalListVersionRequest{}, true},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV2TestSuite) TestGetLocalListVersionConfirmationValidation() {
	var confirmationTable = []GenericTestEntry{
		{localauth.GetLocalListVersionResponse{VersionNumber: 1}, true},
		{localauth.GetLocalListVersionResponse{VersionNumber: 0}, true},
		{localauth.GetLocalListVersionResponse{}, true},
		{localauth.GetLocalListVersionResponse{VersionNumber: -1}, false},
	}
	ExecuteGenericTestTable(suite, confirmationTable)
}

func (suite *OcppV2TestSuite) TestGetLocalListVersionE2EMocked() {
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	listVersion := 1
	requestJson := fmt.Sprintf(`[2,"%v","%v",{}]`, messageId, localauth.GetLocalListVersionFeatureName)
	responseJson := fmt.Sprintf(`[3,"%v",{"versionNumber":%v}]`, messageId, listVersion)
	localListVersionConfirmation := localauth.NewGetLocalListVersionResponse(listVersion)
	channel := NewMockWebSocket(wsId)

	handler := &MockChargingStationLocalAuthHandler{}
	handler.On("OnGetLocalListVersion", mock.Anything).Return(localListVersionConfirmation, nil)
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	suite.Nil(err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.GetLocalListVersion(wsId, func(confirmation *localauth.GetLocalListVersionResponse, err error) {
		suite.Nil(err)
		suite.Require().NotNil(confirmation)
		suite.Equal(listVersion, confirmation.VersionNumber)
		resultChannel <- true
	})
	suite.Nil(err)
	if err == nil {
		result := <-resultChannel
		suite.True(result)
	}
}

func (suite *OcppV2TestSuite) TestGetLocalListVersionInvalidEndpoint() {
	messageId := defaultMessageId
	localListVersionRequest := localauth.NewGetLocalListVersionRequest()
	requestJson := fmt.Sprintf(`[2,"%v","%v",{}]`, messageId, localauth.GetLocalListVersionFeatureName)
	testUnsupportedRequestFromChargingStation(suite, localListVersionRequest, requestJson, messageId)
}
