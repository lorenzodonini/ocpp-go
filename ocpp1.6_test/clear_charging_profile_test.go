package ocpp16_test

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/smartcharging"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test
func (suite *OcppV16TestSuite) TestClearChargingProfileRequestValidation() {
	t := suite.T()
	requestTable := []GenericTestEntry{
		{smartcharging.ClearChargingProfileRequest{Id: newInt(1), ConnectorId: newInt(1), ChargingProfilePurpose: types.ChargingProfilePurposeChargePointMaxProfile, StackLevel: newInt(1)}, true},
		{smartcharging.ClearChargingProfileRequest{Id: newInt(1), ConnectorId: newInt(1), ChargingProfilePurpose: types.ChargingProfilePurposeChargePointMaxProfile}, true},
		{smartcharging.ClearChargingProfileRequest{Id: newInt(1), ConnectorId: newInt(1)}, true},
		{smartcharging.ClearChargingProfileRequest{Id: newInt(1)}, true},
		{smartcharging.ClearChargingProfileRequest{}, true},
		{smartcharging.ClearChargingProfileRequest{ConnectorId: newInt(-1)}, false},
		{smartcharging.ClearChargingProfileRequest{Id: newInt(-1)}, true},
		{smartcharging.ClearChargingProfileRequest{ChargingProfilePurpose: "invalidChargingProfilePurposeType"}, false},
		{smartcharging.ClearChargingProfileRequest{StackLevel: newInt(-1)}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV16TestSuite) TestClearChargingProfileConfirmationValidation() {
	t := suite.T()
	confirmationTable := []GenericTestEntry{
		{smartcharging.ClearChargingProfileConfirmation{Status: smartcharging.ClearChargingProfileStatusAccepted}, true},
		{smartcharging.ClearChargingProfileConfirmation{Status: "invalidClearChargingProfileStatus"}, false},
		{smartcharging.ClearChargingProfileConfirmation{}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV16TestSuite) TestClearChargingProfileE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	chargingProfileId := newInt(1)
	connectorId := newInt(1)
	chargingProfilePurpose := types.ChargingProfilePurposeChargePointMaxProfile
	stackLevel := newInt(1)
	status := smartcharging.ClearChargingProfileStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"id":%v,"connectorId":%v,"chargingProfilePurpose":"%v","stackLevel":%v}]`,
		messageId, smartcharging.ClearChargingProfileFeatureName, *chargingProfileId, *connectorId, chargingProfilePurpose, *stackLevel)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	ClearChargingProfileConfirmation := smartcharging.NewClearChargingProfileConfirmation(status)
	channel := NewMockWebSocket(wsId)

	smartChargingListener := &MockChargePointSmartChargingListener{}
	smartChargingListener.On("OnClearChargingProfile", mock.Anything).Return(ClearChargingProfileConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*smartcharging.ClearChargingProfileRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, *chargingProfileId, *request.Id)
		assert.Equal(t, *connectorId, *request.ConnectorId)
		assert.Equal(t, chargingProfilePurpose, request.ChargingProfilePurpose)
		assert.Equal(t, *stackLevel, *request.StackLevel)
	})
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	suite.chargePoint.SetSmartChargingHandler(smartChargingListener)
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	assert.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.centralSystem.ClearChargingProfile(wsId, func(confirmation *smartcharging.ClearChargingProfileConfirmation, err error) {
		require.Nil(t, err)
		require.NotNil(t, confirmation)
		assert.Equal(t, status, confirmation.Status)
		resultChannel <- true
	}, func(request *smartcharging.ClearChargingProfileRequest) {
		request.Id = chargingProfileId
		request.ConnectorId = connectorId
		request.ChargingProfilePurpose = chargingProfilePurpose
		request.StackLevel = stackLevel
	})
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV16TestSuite) TestClearChargingProfileInvalidEndpoint() {
	messageId := defaultMessageId
	chargingProfileId := newInt(1)
	connectorId := newInt(1)
	chargingProfilePurpose := types.ChargingProfilePurposeChargePointMaxProfile
	stackLevel := newInt(1)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"id":%v,"connectorId":%v,"chargingProfilePurpose":"%v","stackLevel":%v}]`,
		messageId, smartcharging.ClearChargingProfileFeatureName, *chargingProfileId, *connectorId, chargingProfilePurpose, *stackLevel)
	clearChargingProfileRequest := smartcharging.NewClearChargingProfileRequest()
	clearChargingProfileRequest.Id = chargingProfileId
	clearChargingProfileRequest.ConnectorId = connectorId
	clearChargingProfileRequest.ChargingProfilePurpose = chargingProfilePurpose
	clearChargingProfileRequest.StackLevel = stackLevel
	testUnsupportedRequestFromChargePoint(suite, clearChargingProfileRequest, requestJson, messageId)
}
