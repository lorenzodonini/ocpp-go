package ocpp16_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test
func (suite *OcppV16TestSuite) TestGetLocalListVersionRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{auth.GetLocalListVersionRequest{}, true},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV16TestSuite) TestGetLocalListVersionConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{auth.GetLocalListVersionConfirmation{ListVersion: 1}, true},
		{auth.GetLocalListVersionConfirmation{ListVersion: 0}, true},
		{auth.GetLocalListVersionConfirmation{}, true},
		{auth.GetLocalListVersionConfirmation{ListVersion: -1}, true},
		{auth.GetLocalListVersionConfirmation{ListVersion: -2}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV16TestSuite) TestGetLocalListVersionE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	listVersion := 1
	requestJson := fmt.Sprintf(`[2,"%v","%v",{}]`, messageId, auth.GetLocalListVersionFeatureName)
	responseJson := fmt.Sprintf(`[3,"%v",{"listVersion":%v}]`, messageId, listVersion)
	localListVersionConfirmation := auth.NewGetLocalListVersionConfirmation(listVersion)
	channel := NewMockWebSocket(wsId)

	localAuthListListener := MockChargePointLocalAuthListListener{}
	localAuthListListener.On("OnGetLocalListVersion", mock.Anything).Return(localListVersionConfirmation, nil)
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	suite.chargePoint.SetLocalAuthListHandler(localAuthListListener)
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	assert.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.centralSystem.GetLocalListVersion(wsId, func(confirmation *auth.GetLocalListVersionConfirmation, err error) {
		assert.Nil(t, err)
		assert.NotNil(t, confirmation)
		if confirmation != nil {
			assert.Equal(t, listVersion, confirmation.ListVersion)
			resultChannel <- true
		} else {
			resultChannel <- false
		}
	})
	assert.Nil(t, err)
	if err == nil {
		result := <-resultChannel
		assert.True(t, result)
	}
}

func (suite *OcppV16TestSuite) TestGetLocalListVersionInvalidEndpoint() {
	messageId := defaultMessageId
	localListVersionRequest := auth.NewGetLocalListVersionRequest()
	requestJson := fmt.Sprintf(`[2,"%v","%v",{}]`, messageId, auth.GetLocalListVersionFeatureName)
	testUnsupportedRequestFromChargePoint(suite, localListVersionRequest, requestJson, messageId)
}
