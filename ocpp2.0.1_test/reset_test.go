package ocpp2_test

import (
	"fmt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/provisioning"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

// Test
func (suite *OcppV2TestSuite) TestResetE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	resetType := provisioning.ResetTypeImmediate
	evseID := tests.NewInt(42)
	status := provisioning.ResetStatusAccepted
	statusInfo := types.NewStatusInfo("200", "")
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"type":"%v","evseId":%v}]`,
		messageId, provisioning.ResetFeatureName, resetType, *evseID)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v","statusInfo":{"reasonCode":"%v"}}]`, messageId, status, statusInfo.ReasonCode)
	resetResponse := provisioning.NewResetResponse(status)
	resetResponse.StatusInfo = statusInfo
	channel := NewMockWebSocket(wsId)

	handler := &MockChargingStationProvisioningHandler{}
	handler.On("OnReset", mock.Anything).Return(resetResponse, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*provisioning.ResetRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, resetType, request.Type)
		assert.Equal(t, *evseID, *request.EvseID)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.Reset(wsId, func(resp *provisioning.ResetResponse, err error) {
		require.Nil(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, status, resp.Status)
		assert.Equal(t, statusInfo.ReasonCode, resp.StatusInfo.ReasonCode)
		resultChannel <- true
	}, resetType, func(request *provisioning.ResetRequest) {
		request.EvseID = evseID
	})
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV2TestSuite) TestResetInvalidEndpoint() {
	messageId := defaultMessageId
	resetType := provisioning.ResetTypeImmediate
	evseID := tests.NewInt(42)
	resetRequest := provisioning.NewResetRequest(resetType)
	resetRequest.EvseID = evseID
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"type":"%v","evseId":%v}]`,
		messageId, provisioning.ResetFeatureName, resetType, *evseID)

	testUnsupportedRequestFromChargingStation(suite, resetRequest, requestJson, messageId)
}
