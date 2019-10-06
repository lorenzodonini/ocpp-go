package ocpp16_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"time"
)

// Test
func (suite *OcppV16TestSuite) TestHeartbeatRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{ocpp16.HeartbeatRequest{}, true},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV16TestSuite) TestHeartbeatConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{ocpp16.HeartbeatConfirmation{CurrentTime: ocpp16.NewDateTime(time.Now())}, true},
		{ocpp16.HeartbeatConfirmation{}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV16TestSuite) TestHeartbeatE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	currentTime := ocpp16.NewDateTime(time.Now())
	requestJson := fmt.Sprintf(`[2,"%v","%v",{}]`, messageId, ocpp16.HeartbeatFeatureName)
	responseJson := fmt.Sprintf(`[3,"%v",{"currentTime":"%v"}]`, messageId, currentTime.Time.Format(ocpp16.ISO8601))
	heartbeatConfirmation := ocpp16.NewHeartbeatConfirmation(currentTime)
	channel := NewMockWebSocket(wsId)

	coreListener := MockCentralSystemCoreListener{}
	coreListener.On("OnHeartbeat", mock.AnythingOfType("string"), mock.Anything).Return(heartbeatConfirmation, nil)
	setupDefaultCentralSystemHandlers(suite, coreListener, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	assert.Nil(t, err)
	confirmation, err := suite.chargePoint.Heartbeat()
	assert.Nil(t, err)
	assert.NotNil(t, confirmation)
	assertDateTimeEquality(t, *currentTime, *confirmation.CurrentTime)
}

func (suite *OcppV16TestSuite) TestHeartbeatInvalidEndpoint() {
	messageId := defaultMessageId
	heartbeatRequest := ocpp16.NewHeartbeatRequest()
	requestJson := fmt.Sprintf(`[2,"%v","%v",{}]`, messageId, ocpp16.HeartbeatFeatureName)
	testUnsupportedRequestFromCentralSystem(suite, heartbeatRequest, requestJson, messageId)
}
