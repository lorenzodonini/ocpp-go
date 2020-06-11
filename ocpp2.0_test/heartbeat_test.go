package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/availability"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"time"
)

// Test
func (suite *OcppV2TestSuite) TestHeartbeatRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{availability.HeartbeatRequest{}, true},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestHeartbeatResponseValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{availability.HeartbeatResponse{CurrentTime: *types.NewDateTime(time.Now())}, true},
		{availability.HeartbeatResponse{}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestHeartbeatE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	currentTime := types.NewDateTime(time.Now())
	requestJson := fmt.Sprintf(`[2,"%v","%v",{}]`, messageId, availability.HeartbeatFeatureName)
	responseJson := fmt.Sprintf(`[3,"%v",{"currentTime":"%v"}]`, messageId, currentTime.FormatTimestamp())
	heartbeatResponse := availability.NewHeartbeatResponse(*currentTime)
	channel := NewMockWebSocket(wsId)

	handler := MockCSMSAvailabilityHandler{}
	handler.On("OnHeartbeat", mock.AnythingOfType("string"), mock.Anything).Return(heartbeatResponse, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*availability.HeartbeatRequest)
		require.True(t, ok)
		require.NotNil(t, request)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	assert.Nil(t, err)
	response, err := suite.chargingStation.Heartbeat()
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assertDateTimeEquality(t, currentTime, &response.CurrentTime)
}

func (suite *OcppV2TestSuite) TestHeartbeatInvalidEndpoint() {
	messageId := defaultMessageId
	heartbeatRequest := availability.NewHeartbeatRequest()
	requestJson := fmt.Sprintf(`[2,"%v","%v",{}]`, messageId, availability.HeartbeatFeatureName)
	testUnsupportedRequestFromCentralSystem(suite, heartbeatRequest, requestJson, messageId)
}
