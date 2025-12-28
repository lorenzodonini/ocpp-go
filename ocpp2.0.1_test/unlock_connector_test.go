package ocpp2_test

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/remotecontrol"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"

	"github.com/stretchr/testify/mock"
)

// Test
func (suite *OcppV2TestSuite) TestUnlockConnectorRequestValidation() {
	var requestTable = []GenericTestEntry{
		{remotecontrol.UnlockConnectorRequest{EvseID: 2, ConnectorID: 1}, true},
		{remotecontrol.UnlockConnectorRequest{EvseID: 2}, true},
		{remotecontrol.UnlockConnectorRequest{}, true},
		{remotecontrol.UnlockConnectorRequest{EvseID: -1, ConnectorID: 1}, false},
		{remotecontrol.UnlockConnectorRequest{EvseID: 2, ConnectorID: -1}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV2TestSuite) TestUnlockConnectorResponseValidation() {
	var responseTable = []GenericTestEntry{
		{remotecontrol.UnlockConnectorResponse{Status: remotecontrol.UnlockStatusUnlocked, StatusInfo: &types.StatusInfo{ReasonCode: "200"}}, true},
		{remotecontrol.UnlockConnectorResponse{Status: remotecontrol.UnlockStatusUnlocked}, true},
		{remotecontrol.UnlockConnectorResponse{}, false},
		{remotecontrol.UnlockConnectorResponse{Status: "invalidUnlockStatus", StatusInfo: &types.StatusInfo{ReasonCode: "200"}}, false},
		{remotecontrol.UnlockConnectorResponse{Status: remotecontrol.UnlockStatusUnlocked, StatusInfo: &types.StatusInfo{}}, false},
	}
	ExecuteGenericTestTable(suite, responseTable)
}

func (suite *OcppV2TestSuite) TestUnlockConnectorE2EMocked() {
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
		suite.Require().True(ok)
		suite.Equal(evseID, request.EvseID)
		suite.Equal(connectorID, request.ConnectorID)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	suite.Require().Nil(err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.UnlockConnector(wsId, func(response *remotecontrol.UnlockConnectorResponse, err error) {
		suite.Require().Nil(err)
		suite.Require().NotNil(response)
		suite.Equal(status, response.Status)
		suite.Equal(statusInfo.ReasonCode, response.StatusInfo.ReasonCode)
		resultChannel <- true
	}, evseID, connectorID)
	suite.Require().Nil(err)
	result := <-resultChannel
	suite.True(result)
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
