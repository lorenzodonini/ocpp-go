package ocpp16_test

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test
func (suite *OcppV16TestSuite) TestGetLocalListVersionRequestValidation() {
	t := suite.T()
	requestTable := []GenericTestEntry{
		{localauth.GetLocalListVersionRequest{}, true},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV16TestSuite) TestGetLocalListVersionConfirmationValidation() {
	t := suite.T()
	confirmationTable := []GenericTestEntry{
		{localauth.GetLocalListVersionConfirmation{ListVersion: 1}, true},
		{localauth.GetLocalListVersionConfirmation{ListVersion: 0}, true},
		{localauth.GetLocalListVersionConfirmation{}, true},
		{localauth.GetLocalListVersionConfirmation{ListVersion: -1}, true},
		{localauth.GetLocalListVersionConfirmation{ListVersion: -2}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV16TestSuite) TestGetLocalListVersionE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	listVersion := 1
	requestJson := fmt.Sprintf(`[2,"%v","%v",{}]`, messageId, localauth.GetLocalListVersionFeatureName)
	responseJson := fmt.Sprintf(`[3,"%v",{"listVersion":%v}]`, messageId, listVersion)
	localListVersionConfirmation := localauth.NewGetLocalListVersionConfirmation(listVersion)
	channel := NewMockWebSocket(wsId)

	localAuthListListener := MockChargePointLocalAuthListListener{}
	localAuthListListener.On("OnGetLocalListVersion", mock.Anything).Return(localListVersionConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*localauth.GetLocalListVersionRequest)
		require.NotNil(t, request)
		require.True(t, ok)
	})
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	suite.chargePoint.SetLocalAuthListHandler(&localAuthListListener)
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.centralSystem.GetLocalListVersion(wsId, func(confirmation *localauth.GetLocalListVersionConfirmation, err error) {
		require.Nil(t, err)
		require.NotNil(t, confirmation)
		assert.Equal(t, listVersion, confirmation.ListVersion)
		resultChannel <- true
	})
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV16TestSuite) TestGetLocalListVersionInvalidEndpoint() {
	messageId := defaultMessageId
	localListVersionRequest := localauth.NewGetLocalListVersionRequest()
	requestJson := fmt.Sprintf(`[2,"%v","%v",{}]`, messageId, localauth.GetLocalListVersionFeatureName)
	testUnsupportedRequestFromChargePoint(suite, localListVersionRequest, requestJson, messageId)
}
