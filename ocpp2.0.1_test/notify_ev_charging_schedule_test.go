package ocpp2_test

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/smartcharging"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

func (suite *OcppV2TestSuite) TestNotifyEVChargingScheduleRequestValidation() {
	chargingSchedule := types.ChargingSchedule{
		StartSchedule:          types.NewDateTime(time.Now()),
		Duration:               newInt(600),
		ChargingRateUnit:       types.ChargingRateUnitWatts,
		MinChargingRate:        newFloat(6.0),
		ChargingSchedulePeriod: []types.ChargingSchedulePeriod{types.NewChargingSchedulePeriod(0, 10.0)},
	}
	var requestTable = []GenericTestEntry{
		// {ChargingRateUnit: "invalidStruct"}
		{smartcharging.NotifyEVChargingScheduleRequest{TimeBase: types.NewDateTime(time.Now()), EvseID: 1, ChargingSchedule: chargingSchedule}, true},
		{smartcharging.NotifyEVChargingScheduleRequest{TimeBase: types.NewDateTime(time.Now()), EvseID: 1}, false},
		{smartcharging.NotifyEVChargingScheduleRequest{TimeBase: types.NewDateTime(time.Now()), ChargingSchedule: chargingSchedule}, false},
		{smartcharging.NotifyEVChargingScheduleRequest{EvseID: 1}, false},
		{smartcharging.NotifyEVChargingScheduleRequest{}, false},
		{smartcharging.NotifyEVChargingScheduleRequest{TimeBase: types.NewDateTime(time.Now()), EvseID: -1, ChargingSchedule: chargingSchedule}, false},
		{smartcharging.NotifyEVChargingScheduleRequest{TimeBase: types.NewDateTime(time.Now()), EvseID: -1, ChargingSchedule: types.ChargingSchedule{ChargingRateUnit: "invalidStruct"}}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV2TestSuite) TestNotifyEVChargingScheduleResponseValidation() {
	var responseTable = []GenericTestEntry{
		{smartcharging.NotifyEVChargingScheduleResponse{Status: types.GenericStatusAccepted, StatusInfo: types.NewStatusInfo("ok", "someInfo")}, true},
		{smartcharging.NotifyEVChargingScheduleResponse{Status: types.GenericStatusRejected, StatusInfo: types.NewStatusInfo("ok", "someInfo")}, true},
		{smartcharging.NotifyEVChargingScheduleResponse{Status: types.GenericStatusAccepted}, true},
		{smartcharging.NotifyEVChargingScheduleResponse{}, false},
		{smartcharging.NotifyEVChargingScheduleResponse{Status: "invalidStatus"}, false},
		{smartcharging.NotifyEVChargingScheduleResponse{Status: types.GenericStatusAccepted, StatusInfo: types.NewStatusInfo("", "invalidStatusInfo")}, false},
	}
	ExecuteGenericTestTable(suite, responseTable)
}

func (suite *OcppV2TestSuite) TestNotifyEVChargingScheduleE2EMocked() {
	wsId := "test_id"
	messageId := "1234"
	wsUrl := "someUrl"
	timeBase := types.NewDateTime(time.Now())
	evseID := 42
	chargingSchedule := types.ChargingSchedule{
		ID:                     1,
		StartSchedule:          types.NewDateTime(time.Now()),
		Duration:               newInt(600),
		ChargingRateUnit:       types.ChargingRateUnitWatts,
		MinChargingRate:        newFloat(6.0),
		ChargingSchedulePeriod: []types.ChargingSchedulePeriod{types.NewChargingSchedulePeriod(0, 10.0)},
	}
	status := types.GenericStatusAccepted
	statusInfo := types.NewStatusInfo("ok", "someInfo")
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"timeBase":"%v","evseId":%v,"chargingSchedule":{"id":%v,"startSchedule":"%v","duration":%v,"chargingRateUnit":"%v","minChargingRate":%v,"chargingSchedulePeriod":[{"startPeriod":%v,"limit":%v}]}}]`,
		messageId, smartcharging.NotifyEVChargingScheduleFeatureName, timeBase.FormatTimestamp(), evseID, chargingSchedule.ID, chargingSchedule.StartSchedule.FormatTimestamp(), *chargingSchedule.Duration, chargingSchedule.ChargingRateUnit, *chargingSchedule.MinChargingRate, chargingSchedule.ChargingSchedulePeriod[0].StartPeriod, chargingSchedule.ChargingSchedulePeriod[0].Limit)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v","statusInfo":{"reasonCode":"%v","additionalInfo":"%v"}}]`,
		messageId, status, statusInfo.ReasonCode, statusInfo.AdditionalInfo)
	notifyEVChargingScheduleResponse := smartcharging.NewNotifyEVChargingScheduleResponse(status)
	notifyEVChargingScheduleResponse.StatusInfo = statusInfo
	channel := NewMockWebSocket(wsId)

	handler := &MockCSMSSmartChargingHandler{}
	handler.On("OnNotifyEVChargingSchedule", mock.AnythingOfType("string"), mock.Anything).Return(notifyEVChargingScheduleResponse, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*smartcharging.NotifyEVChargingScheduleRequest)
		suite.Require().True(ok)
		suite.Require().NotNil(request)
		suite.Equal(timeBase.FormatTimestamp(), request.TimeBase.FormatTimestamp())
		suite.Equal(evseID, request.EvseID)
		suite.Equal(chargingSchedule.ID, request.ChargingSchedule.ID)
		suite.Equal(chargingSchedule.StartSchedule.FormatTimestamp(), request.ChargingSchedule.StartSchedule.FormatTimestamp())
		suite.Equal(*chargingSchedule.Duration, *request.ChargingSchedule.Duration)
		suite.Equal(*chargingSchedule.MinChargingRate, *request.ChargingSchedule.MinChargingRate)
		suite.Equal(*chargingSchedule.MinChargingRate, *request.ChargingSchedule.MinChargingRate)
		suite.Equal(chargingSchedule.ChargingRateUnit, request.ChargingSchedule.ChargingRateUnit)
		suite.Require().Len(request.ChargingSchedule.ChargingSchedulePeriod, len(request.ChargingSchedule.ChargingSchedulePeriod))
		suite.Equal(chargingSchedule.ChargingSchedulePeriod[0].StartPeriod, request.ChargingSchedule.ChargingSchedulePeriod[0].StartPeriod)
		suite.Equal(chargingSchedule.ChargingSchedulePeriod[0].Limit, request.ChargingSchedule.ChargingSchedulePeriod[0].Limit)
		suite.Nil(request.ChargingSchedule.ChargingSchedulePeriod[0].NumberPhases)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	suite.Require().Nil(err)
	response, err := suite.chargingStation.NotifyEVChargingSchedule(timeBase, evseID, chargingSchedule)
	suite.Require().Nil(err)
	suite.Require().NotNil(response)
	suite.Equal(status, response.Status)
	suite.Equal(statusInfo.ReasonCode, response.StatusInfo.ReasonCode)
	suite.Equal(statusInfo.AdditionalInfo, response.StatusInfo.AdditionalInfo)
}

func (suite *OcppV2TestSuite) TestNotifyEVChargingScheduleInvalidEndpoint() {
	messageId := defaultMessageId
	timeBase := types.NewDateTime(time.Now())
	evseID := 42
	chargingSchedule := types.ChargingSchedule{
		StartSchedule:          types.NewDateTime(time.Now()),
		Duration:               newInt(600),
		ChargingRateUnit:       types.ChargingRateUnitWatts,
		MinChargingRate:        newFloat(6.0),
		ChargingSchedulePeriod: []types.ChargingSchedulePeriod{types.NewChargingSchedulePeriod(0, 10.0)},
	}
	notifyEVChargingScheduleRequest := smartcharging.NewNotifyEVChargingScheduleRequest(timeBase, evseID, chargingSchedule)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"timeBase":"%v","evseId":%v,"chargingSchedule":{"startSchedule":"%v","duration":%v,"chargingRateUnit":"%v","minChargingRate":%v,"chargingSchedulePeriod":[{"startPeriod":%v,"limit":%v}]}}]`,
		messageId, smartcharging.NotifyEVChargingScheduleFeatureName, timeBase.FormatTimestamp(), evseID, chargingSchedule.StartSchedule.FormatTimestamp(), *chargingSchedule.Duration, chargingSchedule.ChargingRateUnit, *chargingSchedule.MinChargingRate, chargingSchedule.ChargingSchedulePeriod[0].StartPeriod, chargingSchedule.ChargingSchedulePeriod[0].Limit)
	testUnsupportedRequestFromCentralSystem(suite, notifyEVChargingScheduleRequest, requestJson, messageId)
}
