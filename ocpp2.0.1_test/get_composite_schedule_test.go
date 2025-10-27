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
	compositeSchedule := smartcharging.CompositeSchedule{
		ChargingSchedulePeriod: []types.ChargingSchedulePeriod{
			{
				StartPeriod:  0,
				Limit:        32.0,
				NumberPhases: newInt(1),
				PhaseToUse:   newInt(3),
			},
			{
				StartPeriod:  3600,
				Limit:        16.0,
				NumberPhases: newInt(1),
				PhaseToUse:   newInt(3),
			},
		},
		ScheduleStart:    types.DateTime{Time: time.Now()},
		EvseID:           1,
		Duration:         18000,
		ChargingRateUnit: types.ChargingRateUnitAmperes,
	}
	var confirmationTable = []GenericTestEntry{
		{smartcharging.GetCompositeScheduleResponse{Status: smartcharging.GetCompositeScheduleStatusAccepted, StatusInfo: types.NewStatusInfo("reasoncode", ""), Schedule: &compositeSchedule}, true},
		{smartcharging.GetCompositeScheduleResponse{Status: smartcharging.GetCompositeScheduleStatusAccepted, StatusInfo: types.NewStatusInfo("reasoncode", "")}, true},
		{smartcharging.GetCompositeScheduleResponse{Status: smartcharging.GetCompositeScheduleStatusAccepted}, true},
		{smartcharging.GetCompositeScheduleResponse{}, false},
		{smartcharging.GetCompositeScheduleResponse{Status: "invalidGetCompositeScheduleStatus"}, false},
		{smartcharging.GetCompositeScheduleResponse{Status: smartcharging.GetCompositeScheduleStatusAccepted, StatusInfo: types.NewStatusInfo("invalidreasoncodeasitslongerthan20", "")}, false},
		{smartcharging.GetCompositeScheduleResponse{Status: smartcharging.GetCompositeScheduleStatusAccepted, StatusInfo: types.NewStatusInfo("", ""), Schedule: &smartcharging.CompositeSchedule{ChargingSchedulePeriod: []types.ChargingSchedulePeriod{{StartPeriod: 0, Limit: 32}}, ScheduleStart: types.DateTime{Time: time.Now()}, EvseID: 1, Duration: 86400, ChargingRateUnit: "invalidChargingRateUnit"}}, false},
		{smartcharging.GetCompositeScheduleResponse{Status: smartcharging.GetCompositeScheduleStatusAccepted, StatusInfo: types.NewStatusInfo("", ""), Schedule: &smartcharging.CompositeSchedule{ChargingSchedulePeriod: []types.ChargingSchedulePeriod{}, ScheduleStart: types.DateTime{Time: time.Now()}, EvseID: 1, Duration: 86400, ChargingRateUnit: types.ChargingRateUnitAmperes}}, false},
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
	duration := 18000
	status := smartcharging.GetCompositeScheduleStatusAccepted
	scheduleStart := types.NewDateTime(time.Now())
	numberPhases := 1
	phasesToUse := 3
	period1 := types.ChargingSchedulePeriod{
		StartPeriod:  0,
		Limit:        10.0,
		NumberPhases: &numberPhases,
		PhaseToUse:   &phasesToUse,
	}
	period2 := types.ChargingSchedulePeriod{
		StartPeriod:  300,
		Limit:        8.0,
		NumberPhases: &numberPhases,
		PhaseToUse:   &phasesToUse,
	}
	compositeSchedule := smartcharging.CompositeSchedule{
		ChargingSchedulePeriod: []types.ChargingSchedulePeriod{period1, period2},
		ScheduleStart:          *scheduleStart,
		EvseID:                 evseID,
		Duration:               duration,
		ChargingRateUnit:       chargingRateUnit,
	}

	statusInfo := types.NewStatusInfo("reasonCode", "")
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"duration":%v,"chargingRateUnit":"%v","evseId":%v}]`,
		messageId, smartcharging.GetCompositeScheduleFeatureName, duration, chargingRateUnit, evseID)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v","statusInfo":{"reasonCode":"%v"},"schedule":{"chargingSchedulePeriod":[{"startPeriod":%v,"limit":%v,"numberPhases":%v,"phaseToUse":%v},{"startPeriod":%v,"limit":%v,"numberPhases":%v,"phaseToUse":%v}],"evseId":%v,"duration":%v,"scheduleStart":"%v","chargingRateUnit":"%v"}}]`,
		messageId, status, statusInfo.ReasonCode, period1.StartPeriod, period1.Limit, numberPhases, phasesToUse, period2.StartPeriod, period2.Limit, numberPhases, phasesToUse, evseID, duration, compositeSchedule.ScheduleStart.FormatTimestamp(), chargingRateUnit)
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
		require.NotNil(t, confirmation.Schedule.ScheduleStart)
		assert.Equal(t, compositeSchedule.ScheduleStart.FormatTimestamp(), confirmation.Schedule.ScheduleStart.FormatTimestamp())
		assert.Equal(t, evseID, confirmation.Schedule.EvseID)
		assert.Equal(t, duration, confirmation.Schedule.Duration)
		require.Len(t, confirmation.Schedule.ChargingSchedulePeriod, len(compositeSchedule.ChargingSchedulePeriod))
		for i := 0; i < len(compositeSchedule.ChargingSchedulePeriod); i++ {
			refPeriod := compositeSchedule.ChargingSchedulePeriod[i]
			confPeriod := confirmation.Schedule.ChargingSchedulePeriod[i]

			assert.Equal(t, refPeriod.StartPeriod, confPeriod.StartPeriod)
			assert.Equal(t, refPeriod.Limit, confPeriod.Limit)
			require.NotNil(t, confPeriod.NumberPhases)
			assert.Equal(t, *refPeriod.NumberPhases, *confPeriod.NumberPhases)
			require.NotNil(t, confPeriod.PhaseToUse)
			assert.Equal(t, *refPeriod.PhaseToUse, *confPeriod.PhaseToUse)
		}
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
