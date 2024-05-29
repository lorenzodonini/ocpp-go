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
func (suite *OcppV16TestSuite) TestSetChargingProfileRequestValidation() {
	t := suite.T()
	chargingSchedule := types.NewChargingSchedule(types.ChargingRateUnitWatts, types.NewChargingSchedulePeriod(0, 10.0))
	chargingProfile := types.NewChargingProfile(1, 1, types.ChargingProfilePurposeChargePointMaxProfile, types.ChargingProfileKindAbsolute, chargingSchedule)
	requestTable := []GenericTestEntry{
		{smartcharging.SetChargingProfileRequest{ConnectorId: 1, ChargingProfile: chargingProfile}, true},
		{smartcharging.SetChargingProfileRequest{ChargingProfile: chargingProfile}, true},
		{smartcharging.SetChargingProfileRequest{}, false},
		{smartcharging.SetChargingProfileRequest{ConnectorId: 1}, false},
		{smartcharging.SetChargingProfileRequest{ConnectorId: -1, ChargingProfile: chargingProfile}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV16TestSuite) TestSetChargingProfileConfirmationValidation() {
	t := suite.T()
	confirmationTable := []GenericTestEntry{
		{smartcharging.SetChargingProfileConfirmation{Status: smartcharging.ChargingProfileStatusAccepted}, true},
		{smartcharging.SetChargingProfileConfirmation{Status: "invalidChargingProfileStatus"}, false},
		{smartcharging.SetChargingProfileConfirmation{}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV16TestSuite) TestSetChargingProfileE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	connectorId := 1
	chargingProfileId := 1
	stackLevel := 1
	chargingProfilePurpose := types.ChargingProfilePurposeChargePointMaxProfile
	chargingProfileKind := types.ChargingProfileKindAbsolute
	chargingRateUnit := types.ChargingRateUnitWatts
	startPeriod := 0
	limit := 10.0
	status := smartcharging.ChargingProfileStatusAccepted
	chargingSchedule := types.NewChargingSchedule(chargingRateUnit, types.NewChargingSchedulePeriod(startPeriod, limit))
	chargingProfile := types.NewChargingProfile(chargingProfileId, stackLevel, chargingProfilePurpose, chargingProfileKind, chargingSchedule)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"connectorId":%v,"csChargingProfiles":{"chargingProfileId":%v,"stackLevel":%v,"chargingProfilePurpose":"%v","chargingProfileKind":"%v","chargingSchedule":{"chargingRateUnit":"%v","chargingSchedulePeriod":[{"startPeriod":%v,"limit":%v}]}}}]`,
		messageId,
		smartcharging.SetChargingProfileFeatureName,
		connectorId,
		chargingProfileId,
		stackLevel,
		chargingProfilePurpose,
		chargingProfileKind,
		chargingRateUnit,
		startPeriod,
		limit)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	SetChargingProfileConfirmation := smartcharging.NewSetChargingProfileConfirmation(status)
	channel := NewMockWebSocket(wsId)

	smartChargingListener := &MockChargePointSmartChargingListener{}
	smartChargingListener.On("OnSetChargingProfile", mock.Anything).Return(SetChargingProfileConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*smartcharging.SetChargingProfileRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, connectorId, request.ConnectorId)
		assert.Equal(t, chargingProfileId, request.ChargingProfile.ChargingProfileId)
		assert.Equal(t, chargingProfileKind, request.ChargingProfile.ChargingProfileKind)
		assert.Equal(t, chargingProfilePurpose, request.ChargingProfile.ChargingProfilePurpose)
		assert.Equal(t, types.RecurrencyKindType(""), request.ChargingProfile.RecurrencyKind)
		assert.Equal(t, stackLevel, request.ChargingProfile.StackLevel)
		assert.Equal(t, 0, request.ChargingProfile.TransactionId)
		assert.Nil(t, request.ChargingProfile.ValidFrom)
		assert.Nil(t, request.ChargingProfile.ValidTo)
		assert.Equal(t, chargingRateUnit, request.ChargingProfile.ChargingSchedule.ChargingRateUnit)
		assert.Nil(t, request.ChargingProfile.ChargingSchedule.MinChargingRate)
		assert.Nil(t, request.ChargingProfile.ChargingSchedule.Duration)
		assert.Nil(t, request.ChargingProfile.ChargingSchedule.StartSchedule)
		require.Len(t, request.ChargingProfile.ChargingSchedule.ChargingSchedulePeriod, 1)
		assert.Equal(t, limit, request.ChargingProfile.ChargingSchedule.ChargingSchedulePeriod[0].Limit)
		assert.Equal(t, startPeriod, request.ChargingProfile.ChargingSchedule.ChargingSchedulePeriod[0].StartPeriod)
		assert.Nil(t, request.ChargingProfile.ChargingSchedule.ChargingSchedulePeriod[0].NumberPhases)
	})
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	suite.chargePoint.SetSmartChargingHandler(smartChargingListener)
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.centralSystem.SetChargingProfile(wsId, func(confirmation *smartcharging.SetChargingProfileConfirmation, err error) {
		require.Nil(t, err)
		require.NotNil(t, confirmation)
		assert.Equal(t, status, confirmation.Status)
		resultChannel <- true
	}, connectorId, chargingProfile)
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV16TestSuite) TestSetChargingProfileInvalidEndpoint() {
	messageId := defaultMessageId
	connectorId := 1
	chargingProfileId := 1
	stackLevel := 1
	chargingProfilePurpose := types.ChargingProfilePurposeChargePointMaxProfile
	chargingProfileKind := types.ChargingProfileKindAbsolute
	chargingRateUnit := types.ChargingRateUnitWatts
	startPeriod := 0
	limit := 10.0
	chargingSchedule := types.NewChargingSchedule(chargingRateUnit, types.NewChargingSchedulePeriod(startPeriod, limit))
	chargingProfile := types.NewChargingProfile(chargingProfileId, stackLevel, chargingProfilePurpose, chargingProfileKind, chargingSchedule)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"connectorId":%v,"csChargingProfiles":{"chargingProfileId":%v,"stackLevel":%v,"chargingProfilePurpose":"%v","chargingProfileKind":"%v","chargingSchedule":{"chargingRateUnit":"%v","chargingSchedulePeriod":[{"startPeriod":%v,"limit":%v}]}}}]`,
		messageId,
		smartcharging.SetChargingProfileFeatureName,
		connectorId,
		chargingProfileId,
		stackLevel,
		chargingProfilePurpose,
		chargingProfileKind,
		chargingRateUnit,
		startPeriod,
		limit)
	SetChargingProfileRequest := smartcharging.NewSetChargingProfileRequest(connectorId, chargingProfile)
	testUnsupportedRequestFromChargePoint(suite, SetChargingProfileRequest, requestJson, messageId)
}
