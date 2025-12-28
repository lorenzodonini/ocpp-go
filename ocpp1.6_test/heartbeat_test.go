package ocpp16_test

import (
	"fmt"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/stretchr/testify/mock"
)

// Test
func (suite *OcppV16TestSuite) TestHeartbeatRequestValidation() {
	var requestTable = []GenericTestEntry{
		{core.HeartbeatRequest{}, true},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV16TestSuite) TestHeartbeatConfirmationValidation() {
	var confirmationTable = []GenericTestEntry{
		{core.HeartbeatConfirmation{CurrentTime: types.NewDateTime(time.Now())}, true},
		{core.HeartbeatConfirmation{}, false},
	}
	ExecuteGenericTestTable(suite, confirmationTable)
}

func (suite *OcppV16TestSuite) TestHeartbeatE2EMocked() {
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	currentTime := types.NewDateTime(time.Now())
	requestJson := fmt.Sprintf(`[2,"%v","%v",{}]`, messageId, core.HeartbeatFeatureName)
	responseJson := fmt.Sprintf(`[3,"%v",{"currentTime":"%v"}]`, messageId, currentTime.FormatTimestamp())
	heartbeatConfirmation := core.NewHeartbeatConfirmation(currentTime)
	channel := NewMockWebSocket(wsId)

	coreListener := &MockCentralSystemCoreListener{}
	coreListener.On("OnHeartbeat", mock.AnythingOfType("string"), mock.Anything).Return(heartbeatConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*core.HeartbeatRequest)
		suite.Require().NotNil(request)
		suite.Require().True(ok)
	})
	setupDefaultCentralSystemHandlers(suite, coreListener, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	suite.Require().Nil(err)
	confirmation, err := suite.chargePoint.Heartbeat()
	suite.Require().Nil(err)
	suite.Require().NotNil(confirmation)
	assertDateTimeEquality(suite, *currentTime, *confirmation.CurrentTime)
}

func (suite *OcppV16TestSuite) TestHeartbeatInvalidEndpoint() {
	messageId := defaultMessageId
	heartbeatRequest := core.NewHeartbeatRequest()
	requestJson := fmt.Sprintf(`[2,"%v","%v",{}]`, messageId, core.HeartbeatFeatureName)
	testUnsupportedRequestFromCentralSystem(suite, heartbeatRequest, requestJson, messageId)
}
