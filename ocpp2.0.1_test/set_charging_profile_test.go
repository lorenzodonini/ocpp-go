package ocpp2_test

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/smartcharging"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

// Test
func (suite *OcppV2TestSuite) TestSetChargingProfileRequestValidation() {
	t := suite.T()
	schedule := types.NewChargingSchedule(1, types.ChargingRateUnitWatts, types.NewChargingSchedulePeriod(0, 200.0))
	chargingProfile := types.NewChargingProfile(
		1,
		0,
		types.ChargingProfilePurposeChargingStationMaxProfile,
		types.ChargingProfileKindAbsolute,
		[]types.ChargingSchedule{*schedule})
	var requestTable = []GenericTestEntry{
		{smartcharging.SetChargingProfileRequest{EvseID: 1, ChargingProfile: chargingProfile}, true},
		{smartcharging.SetChargingProfileRequest{ChargingProfile: chargingProfile}, true},
		{smartcharging.SetChargingProfileRequest{}, false},
		{smartcharging.SetChargingProfileRequest{EvseID: 1, ChargingProfile: types.NewChargingProfile(1, -1, types.ChargingProfilePurposeChargingStationMaxProfile, types.ChargingProfileKindAbsolute, []types.ChargingSchedule{*schedule})}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestSetChargingProfileResponseValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{smartcharging.SetChargingProfileResponse{Status: smartcharging.ChargingProfileStatusAccepted, StatusInfo: types.NewStatusInfo("200", "")}, true},
		{smartcharging.SetChargingProfileResponse{Status: smartcharging.ChargingProfileStatusAccepted}, true},
		{smartcharging.SetChargingProfileResponse{}, false},
		{smartcharging.SetChargingProfileResponse{Status: "invalidChargingProfileStatus"}, false},
		{smartcharging.SetChargingProfileResponse{Status: smartcharging.ChargingProfileStatusAccepted, StatusInfo: types.NewStatusInfo("", "")}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestSetChargingProfileE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	evseID := 1
	period := types.NewChargingSchedulePeriod(0, 200.0)
	schedule := types.NewChargingSchedule(
		1,
		types.ChargingRateUnitWatts,
		period)
	profile := types.NewChargingProfile(
		1,
		7,
		types.ChargingProfilePurposeChargingStationMaxProfile,
		types.ChargingProfileKindAbsolute,
		[]types.ChargingSchedule{*schedule})
	profile.ValidFrom = types.NewDateTime(time.Now())
	status := smartcharging.ChargingProfileStatusAccepted
	statusInfo := types.NewStatusInfo("200", "")
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"evseId":%v,"chargingProfile":{"id":%v,"stackLevel":%v,"chargingProfilePurpose":"%v","chargingProfileKind":"%v","validFrom":"%v","chargingSchedule":[{"id":%v,"chargingRateUnit":"%v","chargingSchedulePeriod":[{"startPeriod":%v,"limit":%v}]}]}}]`,
		messageId, smartcharging.SetChargingProfileFeatureName, evseID, profile.ID, profile.StackLevel, profile.ChargingProfilePurpose, profile.ChargingProfileKind, profile.ValidFrom.FormatTimestamp(), schedule.ID, schedule.ChargingRateUnit, period.StartPeriod, period.Limit)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v","statusInfo":{"reasonCode":"%v"}}]`, messageId, status, statusInfo.ReasonCode)
	setChargingProfileResponse := smartcharging.NewSetChargingProfileResponse(status)
	setChargingProfileResponse.StatusInfo = statusInfo
	channel := NewMockWebSocket(wsId)

	handler := &MockChargingStationSmartChargingHandler{}
	handler.On("OnSetChargingProfile", mock.Anything).Return(setChargingProfileResponse, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*smartcharging.SetChargingProfileRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, evseID, request.EvseID)
		require.NotNil(t, request.ChargingProfile)
		assert.Equal(t, profile.ID, request.ChargingProfile.ID)
		assert.Equal(t, profile.StackLevel, request.ChargingProfile.StackLevel)
		assert.Equal(t, profile.ChargingProfilePurpose, request.ChargingProfile.ChargingProfilePurpose)
		assert.Equal(t, profile.ChargingProfileKind, request.ChargingProfile.ChargingProfileKind)
		assert.Equal(t, profile.ChargingProfileKind, request.ChargingProfile.ChargingProfileKind)
		assert.Equal(t, profile.ValidFrom.FormatTimestamp(), request.ChargingProfile.ValidFrom.FormatTimestamp())
		require.NotNil(t, request.ChargingProfile.ChargingSchedule)
		assert.Len(t, request.ChargingProfile.ChargingSchedule, 1)
		assert.Equal(t, schedule.ID, request.ChargingProfile.ChargingSchedule[0].ID)
		assert.Equal(t, schedule.ChargingRateUnit, request.ChargingProfile.ChargingSchedule[0].ChargingRateUnit)
		require.NotNil(t, request.ChargingProfile.ChargingSchedule[0].ChargingSchedulePeriod)
		assert.Len(t, request.ChargingProfile.ChargingSchedule[0].ChargingSchedulePeriod, 1)
		assert.Equal(t, period.StartPeriod, request.ChargingProfile.ChargingSchedule[0].ChargingSchedulePeriod[0].StartPeriod)
		assert.Equal(t, period.Limit, request.ChargingProfile.ChargingSchedule[0].ChargingSchedulePeriod[0].Limit)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.SetChargingProfile(wsId, func(confirmation *smartcharging.SetChargingProfileResponse, err error) {
		require.Nil(t, err)
		require.NotNil(t, confirmation)
		assert.Equal(t, status, confirmation.Status)
		resultChannel <- true
	}, evseID, profile)
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV2TestSuite) TestSetChargingProfileInvalidEndpoint() {
	messageId := defaultMessageId
	evseID := 1
	period := types.NewChargingSchedulePeriod(0, 200.0)
	schedule := types.NewChargingSchedule(
		1,
		types.ChargingRateUnitWatts,
		period)
	profile := types.NewChargingProfile(
		1,
		7,
		types.ChargingProfilePurposeChargingStationMaxProfile,
		types.ChargingProfileKindAbsolute,
		[]types.ChargingSchedule{*schedule})
	profile.ValidFrom = types.NewDateTime(time.Now())
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"evseId":%v,"chargingProfile":{"id":%v,"stackLevel":%v,"chargingProfilePurpose":"%v","chargingProfileKind":"%v","validFrom":"%v","chargingSchedule":[{"id":%v,"chargingRateUnit":"%v","chargingSchedulePeriod":[{"startPeriod":%v,"limit":%v}]}]}}]`,
		messageId, smartcharging.SetChargingProfileFeatureName, evseID, profile.ID, profile.StackLevel, profile.ChargingProfilePurpose, profile.ChargingProfileKind, profile.ValidFrom.FormatTimestamp(), schedule.ID, schedule.ChargingRateUnit, period.StartPeriod, period.Limit)
	request := smartcharging.NewSetChargingProfileRequest(evseID, profile)
	testUnsupportedRequestFromChargingStation(suite, request, requestJson, messageId)
}
