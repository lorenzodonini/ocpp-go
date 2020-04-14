package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test
func (suite *OcppV2TestSuite) TestGetLocalListVersionRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{ocpp2.GetLocalListVersionRequest{}, true},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestGetLocalListVersionConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{ocpp2.GetLocalListVersionConfirmation{VersionNumber: 1}, true},
		{ocpp2.GetLocalListVersionConfirmation{VersionNumber: 0}, true},
		{ocpp2.GetLocalListVersionConfirmation{}, true},
		{ocpp2.GetLocalListVersionConfirmation{VersionNumber: -1}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestGetLocalListVersionE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	listVersion := 1
	requestJson := fmt.Sprintf(`[2,"%v","%v",{}]`, messageId, ocpp2.GetLocalListVersionFeatureName)
	responseJson := fmt.Sprintf(`[3,"%v",{"versionNumber":%v}]`, messageId, listVersion)
	localListVersionConfirmation := ocpp2.NewGetLocalListVersionConfirmation(listVersion)
	channel := NewMockWebSocket(wsId)

	coreListener := MockChargePointCoreListener{}
	coreListener.On("OnGetLocalListVersion", mock.Anything).Return(localListVersionConfirmation, nil)
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, coreListener, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	assert.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.GetLocalListVersion(wsId, func(confirmation *ocpp2.GetLocalListVersionConfirmation, err error) {
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
	localListVersionRequest := ocpp2.NewGetLocalListVersionRequest()
	requestJson := fmt.Sprintf(`[2,"%v","%v",{}]`, messageId, ocpp2.GetLocalListVersionFeatureName)
	testUnsupportedRequestFromChargePoint(suite, localListVersionRequest, requestJson, messageId)
}
