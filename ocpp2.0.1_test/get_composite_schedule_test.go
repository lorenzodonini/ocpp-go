package ocpp2_test

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/smartcharging"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

// Test
func (suite *OcppV2TestSuite) TestGetCompositeScheduleRequestValidation() {
	var requestTable = []GenericTestEntry{
		{smartcharging.GetCompositeScheduleRequest{Duration: 600, EvseID: 1, ChargingRateUnit: types.ChargingRateUnitWatts}, true},
		{smartcharging.GetCompositeScheduleRequest{Duration: 600, EvseID: 1}, true},
		{smartcharging.GetCompositeScheduleRequest{EvseID: 1}, true},
		{smartcharging.GetCompositeScheduleRequest{}, true},
		{smartcharging.GetCompositeScheduleRequest{Duration: 600, EvseID: -1, ChargingRateUnit: types.ChargingRateUnitWatts}, false},
		{smartcharging.GetCompositeScheduleRequest{Duration: -1, EvseID: 1, ChargingRateUnit: types.ChargingRateUnitWatts}, false},
		{smartcharging.GetCompositeScheduleRequest{Duration: 600, EvseID: 1, ChargingRateUnit: "invalidChargingRateUnit"}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV2TestSuite) TestGetCompositeScheduleConfirmationValidation() {
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
	ExecuteGenericTestTable(suite, confirmationTable)
}

func (suite *OcppV2TestSuite) TestGetCompositeScheduleE2EMocked() {
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
	getCompositeScheduleConfirmation.Schedule = &compositeSchedule
	channel := NewMockWebSocket(wsId)

	handler := &MockChargingStationSmartChargingHandler{}
	handler.On("OnGetCompositeSchedule", mock.Anything).Return(getCompositeScheduleConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*smartcharging.GetCompositeScheduleRequest)
		suite.True(ok)
		suite.NotNil(request)
		suite.Equal(duration, request.Duration)
		suite.Equal(chargingRateUnit, request.ChargingRateUnit)
		suite.Equal(evseID, request.EvseID)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	suite.Nil(err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.GetCompositeSchedule(wsId, func(confirmation *smartcharging.GetCompositeScheduleResponse, err error) {
		suite.Require().Nil(err)
		suite.Require().NotNil(confirmation)
		suite.Equal(status, confirmation.Status)
		suite.Equal(statusInfo.ReasonCode, confirmation.StatusInfo.ReasonCode)
		suite.Require().NotNil(confirmation.Schedule)
		suite.Require().NotNil(confirmation.Schedule.StartDateTime)
		suite.Equal(compositeSchedule.StartDateTime.FormatTimestamp(), confirmation.Schedule.StartDateTime.FormatTimestamp())
		suite.Require().NotNil(confirmation.Schedule.ChargingSchedule)
		suite.Equal(chargingSchedule.ID, confirmation.Schedule.ChargingSchedule.ID)
		suite.Equal(chargingSchedule.ChargingRateUnit, confirmation.Schedule.ChargingSchedule.ChargingRateUnit)
		suite.Require().NotNil(confirmation.Schedule.ChargingSchedule.Duration)
		suite.Equal(*chargingSchedule.Duration, *confirmation.Schedule.ChargingSchedule.Duration)
		suite.Require().NotNil(confirmation.Schedule.ChargingSchedule.MinChargingRate)
		suite.Equal(*chargingSchedule.MinChargingRate, *confirmation.Schedule.ChargingSchedule.MinChargingRate)
		suite.Require().NotNil(confirmation.Schedule.ChargingSchedule.StartSchedule)
		suite.Equal(chargingSchedule.StartSchedule.FormatTimestamp(), confirmation.Schedule.ChargingSchedule.StartSchedule.FormatTimestamp())
		suite.Require().Len(confirmation.Schedule.ChargingSchedule.ChargingSchedulePeriod, len(chargingSchedule.ChargingSchedulePeriod))
		suite.Equal(chargingSchedule.ChargingSchedulePeriod[0].Limit, confirmation.Schedule.ChargingSchedule.ChargingSchedulePeriod[0].Limit)
		suite.Equal(chargingSchedule.ChargingSchedulePeriod[0].StartPeriod, confirmation.Schedule.ChargingSchedule.ChargingSchedulePeriod[0].StartPeriod)
		suite.Require().NotNil(confirmation.Schedule.ChargingSchedule.ChargingSchedulePeriod[0].NumberPhases)
		suite.Equal(*chargingSchedule.ChargingSchedulePeriod[0].NumberPhases, *confirmation.Schedule.ChargingSchedule.ChargingSchedulePeriod[0].NumberPhases)
		resultChannel <- true
	}, duration, evseID, func(request *smartcharging.GetCompositeScheduleRequest) {
		request.ChargingRateUnit = chargingRateUnit
	})
	suite.Nil(err)
	result := <-resultChannel
	suite.True(result)
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
