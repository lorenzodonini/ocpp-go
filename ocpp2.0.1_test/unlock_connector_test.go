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
func (suite *OcppV2TestSuite) TestUnlockConnectorE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	evseID := 2
	connectorID := 1
	status := remotecontrol.UnlockStatusUnlocked
	statusInfo := types.NewStatusInfo("200", "")
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"evseId":%v,"connectorId":%v}]`,
		messageId, remotecontrol.UnlockConnectorFeatureName, evseID, connectorID)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v","statusInfo":{"reasonCode":"%v"}}]`,
		messageId, status, statusInfo.ReasonCode)
	triggerMessageResponse := remotecontrol.NewUnlockConnectorResponse(status)
	triggerMessageResponse.StatusInfo = statusInfo
	channel := NewMockWebSocket(wsId)

	handler := &MockChargingStationRemoteControlHandler{}
	handler.On("OnUnlockConnector", mock.Anything).Return(triggerMessageResponse, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*remotecontrol.UnlockConnectorRequest)
		require.True(t, ok)
		assert.Equal(t, evseID, request.EvseID)
		assert.Equal(t, connectorID, request.ConnectorID)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.UnlockConnector(wsId, func(response *remotecontrol.UnlockConnectorResponse, err error) {
		require.Nil(t, err)
		require.NotNil(t, response)
		assert.Equal(t, status, response.Status)
		assert.Equal(t, statusInfo.ReasonCode, response.StatusInfo.ReasonCode)
		resultChannel <- true
	}, evseID, connectorID)
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV2TestSuite) TestUnlockConnectorInvalidEndpoint() {
	messageId := defaultMessageId
	evseID := 2
	connectorID := 1
	request := remotecontrol.NewUnlockConnectorRequest(evseID, connectorID)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"evseId":%v,"connectorId":%v}]`,
		messageId, remotecontrol.UnlockConnectorFeatureName, evseID, connectorID)
	testUnsupportedRequestFromChargingStation(suite, request, requestJson, messageId)
}
