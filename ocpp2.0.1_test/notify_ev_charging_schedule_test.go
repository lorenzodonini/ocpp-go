package ocpp2_test

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/smartcharging"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *OcppV2TestSuite) TestNotifyEVChargingScheduleE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := "1234"
	wsUrl := "someUrl"
	timeBase := types.NewDateTime(time.Now())
	evseID := 42
	chargingSchedule := types.ChargingSchedule{
		ID:                     1,
		StartSchedule:          types.NewDateTime(time.Now()),
		Duration:               tests.NewInt(600),
		ChargingRateUnit:       types.ChargingRateUnitWatts,
		MinChargingRate:        tests.NewFloat(6.0),
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
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, timeBase.FormatTimestamp(), request.TimeBase.FormatTimestamp())
		assert.Equal(t, evseID, request.EvseID)
		assert.Equal(t, chargingSchedule.ID, request.ChargingSchedule.ID)
		assert.Equal(t, chargingSchedule.StartSchedule.FormatTimestamp(), request.ChargingSchedule.StartSchedule.FormatTimestamp())
		assert.Equal(t, *chargingSchedule.Duration, *request.ChargingSchedule.Duration)
		assert.Equal(t, *chargingSchedule.MinChargingRate, *request.ChargingSchedule.MinChargingRate)
		assert.Equal(t, *chargingSchedule.MinChargingRate, *request.ChargingSchedule.MinChargingRate)
		assert.Equal(t, chargingSchedule.ChargingRateUnit, request.ChargingSchedule.ChargingRateUnit)
		require.Len(t, request.ChargingSchedule.ChargingSchedulePeriod, len(request.ChargingSchedule.ChargingSchedulePeriod))
		assert.Equal(t, chargingSchedule.ChargingSchedulePeriod[0].StartPeriod, request.ChargingSchedule.ChargingSchedulePeriod[0].StartPeriod)
		assert.Equal(t, chargingSchedule.ChargingSchedulePeriod[0].Limit, request.ChargingSchedule.ChargingSchedulePeriod[0].Limit)
		assert.Nil(t, request.ChargingSchedule.ChargingSchedulePeriod[0].NumberPhases)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	response, err := suite.chargingStation.NotifyEVChargingSchedule(timeBase, evseID, chargingSchedule)
	require.Nil(t, err)
	require.NotNil(t, response)
	assert.Equal(t, status, response.Status)
	assert.Equal(t, statusInfo.ReasonCode, response.StatusInfo.ReasonCode)
	assert.Equal(t, statusInfo.AdditionalInfo, response.StatusInfo.AdditionalInfo)
}

func (suite *OcppV2TestSuite) TestNotifyEVChargingScheduleInvalidEndpoint() {
	messageId := defaultMessageId
	timeBase := types.NewDateTime(time.Now())
	evseID := 42
	chargingSchedule := types.ChargingSchedule{
		StartSchedule:          types.NewDateTime(time.Now()),
		Duration:               tests.NewInt(600),
		ChargingRateUnit:       types.ChargingRateUnitWatts,
		MinChargingRate:        tests.NewFloat(6.0),
		ChargingSchedulePeriod: []types.ChargingSchedulePeriod{types.NewChargingSchedulePeriod(0, 10.0)},
	}
	notifyEVChargingScheduleRequest := smartcharging.NewNotifyEVChargingScheduleRequest(timeBase, evseID, chargingSchedule)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"timeBase":"%v","evseId":%v,"chargingSchedule":{"startSchedule":"%v","duration":%v,"chargingRateUnit":"%v","minChargingRate":%v,"chargingSchedulePeriod":[{"startPeriod":%v,"limit":%v}]}}]`,
		messageId, smartcharging.NotifyEVChargingScheduleFeatureName, timeBase.FormatTimestamp(), evseID, chargingSchedule.StartSchedule.FormatTimestamp(), *chargingSchedule.Duration, chargingSchedule.ChargingRateUnit, *chargingSchedule.MinChargingRate, chargingSchedule.ChargingSchedulePeriod[0].StartPeriod, chargingSchedule.ChargingSchedulePeriod[0].Limit)
	testUnsupportedRequestFromCentralSystem(suite, notifyEVChargingScheduleRequest, requestJson, messageId)
}
