package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/localauth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test
func (suite *OcppV2TestSuite) TestGetLocalListVersionRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{localauth.GetLocalListVersionRequest{}, true},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestGetLocalListVersionConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{localauth.GetLocalListVersionResponse{VersionNumber: 1}, true},
		{localauth.GetLocalListVersionResponse{VersionNumber: 0}, true},
		{localauth.GetLocalListVersionResponse{}, true},
		{localauth.GetLocalListVersionResponse{VersionNumber: -1}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestGetLocalListVersionE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	listVersion := 1
	requestJson := fmt.Sprintf(`[2,"%v","%v",{}]`, messageId, localauth.GetLocalListVersionFeatureName)
	responseJson := fmt.Sprintf(`[3,"%v",{"versionNumber":%v}]`, messageId, listVersion)
	localListVersionConfirmation := localauth.NewGetLocalListVersionResponse(listVersion)
	channel := NewMockWebSocket(wsId)

	handler := MockChargingStationLocalAuthHandler{}
	handler.On("OnGetLocalListVersion", mock.Anything).Return(localListVersionConfirmation, nil)
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	assert.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.GetLocalListVersion(wsId, func(confirmation *localauth.GetLocalListVersionResponse, err error) {
		assert.Nil(t, err)
		require.NotNil(t, confirmation)
		assert.Equal(t, listVersion, confirmation.VersionNumber)
		resultChannel <- true
	})
	assert.Nil(t, err)
	if err == nil {
		result := <-resultChannel
		assert.True(t, result)
	}
}

func (suite *OcppV2TestSuite) TestGetLocalListVersionInvalidEndpoint() {
	messageId := defaultMessageId
	localListVersionRequest := localauth.NewGetLocalListVersionRequest()
	requestJson := fmt.Sprintf(`[2,"%v","%v",{}]`, messageId, localauth.GetLocalListVersionFeatureName)
	testUnsupportedRequestFromChargingStation(suite, localListVersionRequest, requestJson, messageId)
}
