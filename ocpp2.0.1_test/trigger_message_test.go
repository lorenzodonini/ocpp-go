package ocpp2_test

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/remotecontrol"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test
func (suite *OcppV2TestSuite) TestTriggerMessageRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{remotecontrol.TriggerMessageRequest{RequestedMessage: remotecontrol.MessageTriggerStatusNotification, Evse: &types.EVSE{ID: 1}}, true},
		{remotecontrol.TriggerMessageRequest{RequestedMessage: remotecontrol.MessageTriggerStatusNotification}, true},
		{remotecontrol.TriggerMessageRequest{}, false},
		{remotecontrol.TriggerMessageRequest{RequestedMessage: "invalidMessageTrigger", Evse: &types.EVSE{ID: 1}}, false},
		{remotecontrol.TriggerMessageRequest{RequestedMessage: remotecontrol.MessageTriggerStatusNotification, Evse: &types.EVSE{ID: -1}}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestTriggerMessageResponseValidation() {
	t := suite.T()
	var responseTable = []GenericTestEntry{
		{remotecontrol.TriggerMessageResponse{Status: remotecontrol.TriggerMessageStatusAccepted, StatusInfo: &types.StatusInfo{ReasonCode: "200"}}, true},
		{remotecontrol.TriggerMessageResponse{Status: remotecontrol.TriggerMessageStatusAccepted}, true},
		{remotecontrol.TriggerMessageResponse{}, false},
		{remotecontrol.TriggerMessageResponse{Status: "invalidTriggerMessageStatus", StatusInfo: &types.StatusInfo{ReasonCode: "200"}}, false},
		{remotecontrol.TriggerMessageResponse{Status: remotecontrol.TriggerMessageStatusAccepted, StatusInfo: &types.StatusInfo{}}, false},
	}
	ExecuteGenericTestTable(t, responseTable)
}

func (suite *OcppV2TestSuite) TestTriggerMessageE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	requestedMessage := remotecontrol.MessageTriggerStatusNotification
	evse := types.EVSE{ID: 1}
	status := remotecontrol.TriggerMessageStatusAccepted
	statusInfo := types.NewStatusInfo("200", "")
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"requestedMessage":"%v","evse":{"id":%v}}]`,
		messageId, remotecontrol.TriggerMessageFeatureName, requestedMessage, evse.ID)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v","statusInfo":{"reasonCode":"%v"}}]`,
		messageId, status, statusInfo.ReasonCode)
	triggerMessageResponse := remotecontrol.NewTriggerMessageResponse(status)
	triggerMessageResponse.StatusInfo = statusInfo
	channel := NewMockWebSocket(wsId)

	handler := &MockChargingStationRemoteControlHandler{}
	handler.On("OnTriggerMessage", mock.Anything).Return(triggerMessageResponse, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*remotecontrol.TriggerMessageRequest)
		require.True(t, ok)
		assert.Equal(t, requestedMessage, request.RequestedMessage)
		require.NotNil(t, request.Evse)
		assert.Equal(t, evse.ID, request.Evse.ID)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.TriggerMessage(wsId, func(response *remotecontrol.TriggerMessageResponse, err error) {
		require.Nil(t, err)
		require.NotNil(t, response)
		assert.Equal(t, status, response.Status)
		assert.Equal(t, statusInfo.ReasonCode, response.StatusInfo.ReasonCode)
		resultChannel <- true
	}, requestedMessage, func(request *remotecontrol.TriggerMessageRequest) {
		request.Evse = &evse
	})
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV2TestSuite) TestTriggerMessageInvalidEndpoint() {
	messageId := defaultMessageId
	requestedMessage := remotecontrol.MessageTriggerStatusNotification
	evse := types.EVSE{ID: 1}
	request := remotecontrol.NewTriggerMessageRequest(requestedMessage)
	request.Evse = &evse
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"requestedMessage":"%v","evse":{"id":%v}}]`,
		messageId, remotecontrol.TriggerMessageFeatureName, requestedMessage, evse.ID)
	testUnsupportedRequestFromChargingStation(suite, request, requestJson, messageId)
}
