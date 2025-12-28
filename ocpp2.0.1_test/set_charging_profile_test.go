package ocpp2_test

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/smartcharging"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

// Test
func (suite *OcppV2TestSuite) TestSetChargingProfileRequestValidation() {
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
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV2TestSuite) TestSetChargingProfileResponseValidation() {
	var confirmationTable = []GenericTestEntry{
		{smartcharging.SetChargingProfileResponse{Status: smartcharging.ChargingProfileStatusAccepted, StatusInfo: types.NewStatusInfo("200", "")}, true},
		{smartcharging.SetChargingProfileResponse{Status: smartcharging.ChargingProfileStatusAccepted}, true},
		{smartcharging.SetChargingProfileResponse{}, false},
		{smartcharging.SetChargingProfileResponse{Status: "invalidChargingProfileStatus"}, false},
		{smartcharging.SetChargingProfileResponse{Status: smartcharging.ChargingProfileStatusAccepted, StatusInfo: types.NewStatusInfo("", "")}, false},
	}
	ExecuteGenericTestTable(suite, confirmationTable)
}

func (suite *OcppV2TestSuite) TestSetChargingProfileE2EMocked() {
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
		suite.Require().True(ok)
		suite.Require().NotNil(request)
		suite.Equal(evseID, request.EvseID)
		suite.Require().NotNil(request.ChargingProfile)
		suite.Equal(profile.ID, request.ChargingProfile.ID)
		suite.Equal(profile.StackLevel, request.ChargingProfile.StackLevel)
		suite.Equal(profile.ChargingProfilePurpose, request.ChargingProfile.ChargingProfilePurpose)
		suite.Equal(profile.ChargingProfileKind, request.ChargingProfile.ChargingProfileKind)
		suite.Equal(profile.ChargingProfileKind, request.ChargingProfile.ChargingProfileKind)
		suite.Equal(profile.ValidFrom.FormatTimestamp(), request.ChargingProfile.ValidFrom.FormatTimestamp())
		suite.Require().NotNil(request.ChargingProfile.ChargingSchedule)
		suite.Require().Len(request.ChargingProfile.ChargingSchedule, 1)
		suite.Equal(schedule.ID, request.ChargingProfile.ChargingSchedule[0].ID)
		suite.Equal(schedule.ChargingRateUnit, request.ChargingProfile.ChargingSchedule[0].ChargingRateUnit)
		suite.Require().NotNil(request.ChargingProfile.ChargingSchedule[0].ChargingSchedulePeriod)
		suite.Require().Len(request.ChargingProfile.ChargingSchedule[0].ChargingSchedulePeriod, 1)
		suite.Equal(period.StartPeriod, request.ChargingProfile.ChargingSchedule[0].ChargingSchedulePeriod[0].StartPeriod)
		suite.Equal(period.Limit, request.ChargingProfile.ChargingSchedule[0].ChargingSchedulePeriod[0].Limit)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	suite.Require().Nil(err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.SetChargingProfile(wsId, func(confirmation *smartcharging.SetChargingProfileResponse, err error) {
		suite.Require().Nil(err)
		suite.Require().NotNil(confirmation)
		suite.Equal(status, confirmation.Status)
		resultChannel <- true
	}, evseID, profile)
	suite.Require().Nil(err)
	result := <-resultChannel
	suite.True(result)
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
