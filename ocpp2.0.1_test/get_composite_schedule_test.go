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
	period := types.NewChargingSchedulePeriod(0, 10.0)
	compositeSchedule := smartcharging.CompositeSchedule{
		EvseID:                 1,
		Duration:               600,
		ScheduleStart:          types.NewDateTime(time.Now()),
		ChargingRateUnit:       types.ChargingRateUnitWatts,
		ChargingSchedulePeriod: []types.ChargingSchedulePeriod{period},
	}
	var confirmationTable = []GenericTestEntry{
		{smartcharging.GetCompositeScheduleResponse{Status: smartcharging.GetCompositeScheduleStatusAccepted, StatusInfo: types.NewStatusInfo("reasoncode", ""), Schedule: &compositeSchedule}, true},
		{smartcharging.GetCompositeScheduleResponse{Status: smartcharging.GetCompositeScheduleStatusAccepted, StatusInfo: types.NewStatusInfo("reasoncode", "")}, true},
		{smartcharging.GetCompositeScheduleResponse{Status: smartcharging.GetCompositeScheduleStatusAccepted}, true},
		{smartcharging.GetCompositeScheduleResponse{}, false},
		{smartcharging.GetCompositeScheduleResponse{Status: "invalidGetCompositeScheduleStatus"}, false},
		{smartcharging.GetCompositeScheduleResponse{Status: smartcharging.GetCompositeScheduleStatusAccepted, StatusInfo: types.NewStatusInfo("invalidreasoncodeasitslongerthan20", "")}, false},
		// Invalid: missing required chargingSchedulePeriod
		{smartcharging.GetCompositeScheduleResponse{Status: smartcharging.GetCompositeScheduleStatusAccepted, Schedule: &smartcharging.CompositeSchedule{EvseID: 1, Duration: 600, ScheduleStart: types.NewDateTime(time.Now()), ChargingRateUnit: types.ChargingRateUnitWatts}}, false},
		// Invalid: bad chargingRateUnit
		{smartcharging.GetCompositeScheduleResponse{Status: smartcharging.GetCompositeScheduleStatusAccepted, Schedule: &smartcharging.CompositeSchedule{EvseID: 1, Duration: 600, ScheduleStart: types.NewDateTime(time.Now()), ChargingRateUnit: "invalidChargingRateUnit", ChargingSchedulePeriod: []types.ChargingSchedulePeriod{period}}}, false},
		// Invalid: missing scheduleStart
		{smartcharging.GetCompositeScheduleResponse{Status: smartcharging.GetCompositeScheduleStatusAccepted, Schedule: &smartcharging.CompositeSchedule{EvseID: 1, Duration: 600, ChargingRateUnit: types.ChargingRateUnitWatts, ChargingSchedulePeriod: []types.ChargingSchedulePeriod{period}}}, false},
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
	statusInfo := types.NewStatusInfo("reasonCode", "")
	compositeSchedule := smartcharging.CompositeSchedule{
		EvseID:                 evseID,
		Duration:               duration,
		ScheduleStart:          scheduleStart,
		ChargingRateUnit:       chargingRateUnit,
		ChargingSchedulePeriod: []types.ChargingSchedulePeriod{chargingSchedulePeriod},
	}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"duration":%v,"chargingRateUnit":"%v","evseId":%v}]`,
		messageId, smartcharging.GetCompositeScheduleFeatureName, duration, chargingRateUnit, evseID)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v","statusInfo":{"reasonCode":"%v"},"schedule":{"evseId":%v,"duration":%v,"scheduleStart":"%v","chargingRateUnit":"%v","chargingSchedulePeriod":[{"startPeriod":%v,"limit":%v,"numberPhases":%v}]}}]`,
		messageId, status, statusInfo.ReasonCode, evseID, duration, scheduleStart.FormatTimestamp(), chargingRateUnit, chargingSchedulePeriod.StartPeriod, chargingSchedulePeriod.Limit, *chargingSchedulePeriod.NumberPhases)
	getCompositeScheduleConfirmation := smartcharging.NewGetCompositeScheduleResponse(status)
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
		assert.Equal(t, compositeSchedule.EvseID, confirmation.Schedule.EvseID)
		assert.Equal(t, compositeSchedule.Duration, confirmation.Schedule.Duration)
		require.NotNil(t, confirmation.Schedule.ScheduleStart)
		assert.Equal(t, compositeSchedule.ScheduleStart.FormatTimestamp(), confirmation.Schedule.ScheduleStart.FormatTimestamp())
		assert.Equal(t, compositeSchedule.ChargingRateUnit, confirmation.Schedule.ChargingRateUnit)
		require.Len(t, confirmation.Schedule.ChargingSchedulePeriod, len(compositeSchedule.ChargingSchedulePeriod))
		assert.Equal(t, chargingSchedulePeriod.StartPeriod, confirmation.Schedule.ChargingSchedulePeriod[0].StartPeriod)
		assert.Equal(t, chargingSchedulePeriod.Limit, confirmation.Schedule.ChargingSchedulePeriod[0].Limit)
		require.NotNil(t, confirmation.Schedule.ChargingSchedulePeriod[0].NumberPhases)
		assert.Equal(t, *chargingSchedulePeriod.NumberPhases, *confirmation.Schedule.ChargingSchedulePeriod[0].NumberPhases)
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
