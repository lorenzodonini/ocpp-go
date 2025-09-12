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
func (suite *OcppV2TestSuite) TestGetCompositeScheduleRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{smartcharging.GetCompositeScheduleRequest{Duration: 600, EvseID: 1, ChargingRateUnit: types.ChargingRateUnitWatts}, true},
		{smartcharging.GetCompositeScheduleRequest{Duration: 600, EvseID: 1}, true},
		{smartcharging.GetCompositeScheduleRequest{EvseID: 1}, true},
		{smartcharging.GetCompositeScheduleRequest{}, true},
		{smartcharging.GetCompositeScheduleRequest{Duration: 600, EvseID: -1, ChargingRateUnit: types.ChargingRateUnitWatts}, false},
		{smartcharging.GetCompositeScheduleRequest{Duration: -1, EvseID: 1, ChargingRateUnit: types.ChargingRateUnitWatts}, false},
		{smartcharging.GetCompositeScheduleRequest{Duration: 600, EvseID: 1, ChargingRateUnit: "invalidChargingRateUnit"}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestGetCompositeScheduleConfirmationValidation() {
	t := suite.T()
	chargingSchedule := types.NewChargingSchedule(1, types.ChargingRateUnitWatts, types.NewChargingSchedulePeriod(0, 10.0))
	chargingSchedule.Duration = newInt(600)
	chargingSchedule.MinChargingRate = newFloat(6.0)
	chargingSchedule.StartSchedule = types.NewDateTime(time.Now())
	compositeSchedule := smartcharging.CompositeSchedule{
		StartDateTime:    types.NewDateTime(time.Now()),
		ChargingSchedule: chargingSchedule,
	}
	var confirmationTable = []GenericTestEntry{
		{smartcharging.GetCompositeScheduleResponse{Status: smartcharging.GetCompositeScheduleStatusAccepted, StatusInfo: types.NewStatusInfo("reasoncode", ""), Schedule: &compositeSchedule}, true},
		{smartcharging.GetCompositeScheduleResponse{Status: smartcharging.GetCompositeScheduleStatusAccepted, StatusInfo: types.NewStatusInfo("reasoncode", ""), Schedule: &smartcharging.CompositeSchedule{}}, true},
		{smartcharging.GetCompositeScheduleResponse{Status: smartcharging.GetCompositeScheduleStatusAccepted, StatusInfo: types.NewStatusInfo("reasoncode", "")}, true},
		{smartcharging.GetCompositeScheduleResponse{Status: smartcharging.GetCompositeScheduleStatusAccepted}, true},
		{smartcharging.GetCompositeScheduleResponse{}, false},
		{smartcharging.GetCompositeScheduleResponse{Status: "invalidGetCompositeScheduleStatus"}, false},
		{smartcharging.GetCompositeScheduleResponse{Status: smartcharging.GetCompositeScheduleStatusAccepted, StatusInfo: types.NewStatusInfo("invalidreasoncodeasitslongerthan20", "")}, false},
		{smartcharging.GetCompositeScheduleResponse{Status: smartcharging.GetCompositeScheduleStatusAccepted, StatusInfo: types.NewStatusInfo("", ""), Schedule: &smartcharging.CompositeSchedule{StartDateTime: types.NewDateTime(time.Now()), ChargingSchedule: types.NewChargingSchedule(1, "invalidChargingRateUnit")}}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestGetCompositeScheduleE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	evseID := 1
	chargingRateUnit := types.ChargingRateUnitWatts
	duration := 600
	status := smartcharging.GetCompositeScheduleStatusAccepted
	scheduleStart := types.NewDateTime(time.Now())
	chargingSchedulePeriod := types.NewChargingSchedulePeriod(0, 10.0)
	chargingSchedulePeriod.NumberPhases = newInt(3)
	chargingSchedule := types.NewChargingSchedule(1, chargingRateUnit, chargingSchedulePeriod)
	chargingSchedule.Duration = newInt(600)
	chargingSchedule.StartSchedule = types.NewDateTime(time.Now())
	chargingSchedule.MinChargingRate = newFloat(6.0)
	statusInfo := types.NewStatusInfo("reasonCode", "")
	compositeSchedule := smartcharging.CompositeSchedule{StartDateTime: scheduleStart, ChargingSchedule: chargingSchedule}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"duration":%v,"chargingRateUnit":"%v","evseId":%v}]`,
		messageId, smartcharging.GetCompositeScheduleFeatureName, duration, chargingRateUnit, evseID)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v","statusInfo":{"reasonCode":"%v"},"schedule":{"startDateTime":"%v","chargingSchedule":{"id":%v,"startSchedule":"%v","duration":%v,"chargingRateUnit":"%v","minChargingRate":%v,"chargingSchedulePeriod":[{"startPeriod":%v,"limit":%v,"numberPhases":%v}]}}}]`,
		messageId, status, statusInfo.ReasonCode, compositeSchedule.StartDateTime.FormatTimestamp(), chargingSchedule.ID, chargingSchedule.StartSchedule.FormatTimestamp(), *chargingSchedule.Duration, chargingSchedule.ChargingRateUnit, *chargingSchedule.MinChargingRate, chargingSchedulePeriod.StartPeriod, chargingSchedulePeriod.Limit, *chargingSchedulePeriod.NumberPhases)
	getCompositeScheduleConfirmation := smartcharging.NewGetCompositeScheduleResponse(status)
	getCompositeScheduleConfirmation.StatusInfo = statusInfo
	getCompositeScheduleConfirmation.Schedule = &compositeSchedule
	channel := NewMockWebSocket(wsId)

	handler := &MockChargingStationSmartChargingHandler{}
	handler.On("OnGetCompositeSchedule", mock.Anything).Return(getCompositeScheduleConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*smartcharging.GetCompositeScheduleRequest)
		assert.True(t, ok)
		assert.NotNil(t, request)
		assert.Equal(t, duration, request.Duration)
		assert.Equal(t, chargingRateUnit, request.ChargingRateUnit)
		assert.Equal(t, evseID, request.EvseID)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	assert.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.GetCompositeSchedule(wsId, func(confirmation *smartcharging.GetCompositeScheduleResponse, err error) {
		require.Nil(t, err)
		require.NotNil(t, confirmation)
		assert.Equal(t, status, confirmation.Status)
		assert.Equal(t, statusInfo.ReasonCode, confirmation.StatusInfo.ReasonCode)
		require.NotNil(t, confirmation.Schedule)
		require.NotNil(t, confirmation.Schedule.StartDateTime)
		assert.Equal(t, compositeSchedule.StartDateTime.FormatTimestamp(), confirmation.Schedule.StartDateTime.FormatTimestamp())
		require.NotNil(t, confirmation.Schedule.ChargingSchedule)
		assert.Equal(t, chargingSchedule.ID, confirmation.Schedule.ChargingSchedule.ID)
		assert.Equal(t, chargingSchedule.ChargingRateUnit, confirmation.Schedule.ChargingSchedule.ChargingRateUnit)
		require.NotNil(t, confirmation.Schedule.ChargingSchedule.Duration)
		assert.Equal(t, *chargingSchedule.Duration, *confirmation.Schedule.ChargingSchedule.Duration)
		require.NotNil(t, confirmation.Schedule.ChargingSchedule.MinChargingRate)
		assert.Equal(t, *chargingSchedule.MinChargingRate, *confirmation.Schedule.ChargingSchedule.MinChargingRate)
		require.NotNil(t, confirmation.Schedule.ChargingSchedule.StartSchedule)
		assert.Equal(t, chargingSchedule.StartSchedule.FormatTimestamp(), confirmation.Schedule.ChargingSchedule.StartSchedule.FormatTimestamp())
		require.Len(t, confirmation.Schedule.ChargingSchedule.ChargingSchedulePeriod, len(chargingSchedule.ChargingSchedulePeriod))
		assert.Equal(t, chargingSchedule.ChargingSchedulePeriod[0].Limit, confirmation.Schedule.ChargingSchedule.ChargingSchedulePeriod[0].Limit)
		assert.Equal(t, chargingSchedule.ChargingSchedulePeriod[0].StartPeriod, confirmation.Schedule.ChargingSchedule.ChargingSchedulePeriod[0].StartPeriod)
		require.NotNil(t, confirmation.Schedule.ChargingSchedule.ChargingSchedulePeriod[0].NumberPhases)
		assert.Equal(t, *chargingSchedule.ChargingSchedulePeriod[0].NumberPhases, *confirmation.Schedule.ChargingSchedule.ChargingSchedulePeriod[0].NumberPhases)
		resultChannel <- true
	}, duration, evseID, func(request *smartcharging.GetCompositeScheduleRequest) {
		request.ChargingRateUnit = chargingRateUnit
	})
	assert.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV2TestSuite) TestGetCompositeScheduleInvalidEndpoint() {
	messageId := defaultMessageId
	evseID := 1
	chargingRateUnit := types.ChargingRateUnitWatts
	duration := 600
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"duration":%v,"chargingRateUnit":"%v","evseId":%v}]`,
		messageId, smartcharging.GetCompositeScheduleFeatureName, duration, chargingRateUnit, evseID)
	getCompositeScheduleRequest := smartcharging.NewGetCompositeScheduleRequest(evseID, duration)
	getCompositeScheduleRequest.ChargingRateUnit = chargingRateUnit
	testUnsupportedRequestFromChargingStation(suite, getCompositeScheduleRequest, requestJson, messageId)
}
