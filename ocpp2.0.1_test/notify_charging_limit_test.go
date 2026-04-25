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

// Tests
func (suite *OcppV2TestSuite) TestNotifyChargingLimitE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := "1234"
	wsUrl := "someUrl"
	evseID := tests.NewInt(42)
	chargingLimit := smartcharging.ChargingLimit{ChargingLimitSource: types.ChargingLimitSourceEMS, IsGridCritical: tests.NewBool(false)}
	chargingSchedule := types.ChargingSchedule{
		ID:                     1,
		StartSchedule:          types.NewDateTime(time.Now()),
		Duration:               tests.NewInt(600),
		ChargingRateUnit:       types.ChargingRateUnitWatts,
		MinChargingRate:        tests.NewFloat(6.0),
		ChargingSchedulePeriod: []types.ChargingSchedulePeriod{types.NewChargingSchedulePeriod(0, 10.0)},
	}
	chargingSchedules := []types.ChargingSchedule{chargingSchedule}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"evseId":%v,"chargingLimit":{"chargingLimitSource":"%v","isGridCritical":%v},"chargingSchedule":[{"id":%v,"startSchedule":"%v","duration":%v,"chargingRateUnit":"%v","minChargingRate":%v,"chargingSchedulePeriod":[{"startPeriod":%v,"limit":%v}]}]}]`,
		messageId, smartcharging.NotifyChargingLimitFeatureName, *evseID, chargingLimit.ChargingLimitSource, *chargingLimit.IsGridCritical, chargingSchedule.ID, chargingSchedule.StartSchedule.FormatTimestamp(), *chargingSchedule.Duration, chargingSchedule.ChargingRateUnit, *chargingSchedule.MinChargingRate, chargingSchedule.ChargingSchedulePeriod[0].StartPeriod, chargingSchedule.ChargingSchedulePeriod[0].Limit)
	responseJson := fmt.Sprintf(`[3,"%v",{}]`, messageId)
	response := smartcharging.NewNotifyChargingLimitResponse()
	channel := NewMockWebSocket(wsId)

	handler := &MockCSMSSmartChargingHandler{}
	handler.On("OnNotifyChargingLimit", mock.AnythingOfType("string"), mock.Anything).Return(response, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*smartcharging.NotifyChargingLimitRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, *evseID, *request.EvseID)
		assert.Equal(t, chargingLimit.ChargingLimitSource, request.ChargingLimit.ChargingLimitSource)
		require.NotNil(t, request.ChargingLimit.IsGridCritical)
		assert.Equal(t, chargingLimit.IsGridCritical, request.ChargingLimit.IsGridCritical)
		require.Len(t, request.ChargingSchedule, len(chargingSchedules))
		assertDateTimeEquality(t, chargingSchedule.StartSchedule, request.ChargingSchedule[0].StartSchedule)
		assert.Equal(t, chargingSchedule.ID, request.ChargingSchedule[0].ID)
		assert.Equal(t, *chargingSchedule.Duration, *request.ChargingSchedule[0].Duration)
		assert.Equal(t, chargingSchedule.ChargingRateUnit, request.ChargingSchedule[0].ChargingRateUnit)
		assert.Equal(t, *chargingSchedule.MinChargingRate, *request.ChargingSchedule[0].MinChargingRate)
		assert.Len(t, request.ChargingSchedule[0].ChargingSchedulePeriod, len(chargingSchedule.ChargingSchedulePeriod))
		assert.Equal(t, chargingSchedule.ChargingSchedulePeriod[0].StartPeriod, request.ChargingSchedule[0].ChargingSchedulePeriod[0].StartPeriod)
		assert.Equal(t, chargingSchedule.ChargingSchedulePeriod[0].Limit, request.ChargingSchedule[0].ChargingSchedulePeriod[0].Limit)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	r, err := suite.chargingStation.NotifyChargingLimit(chargingLimit, func(request *smartcharging.NotifyChargingLimitRequest) {
		request.EvseID = evseID
		request.ChargingSchedule = chargingSchedules
	})
	require.Nil(t, err)
	require.NotNil(t, r)
}

func (suite *OcppV2TestSuite) TestNotifyChargingLimitInvalidEndpoint() {
	messageId := defaultMessageId
	evseID := tests.NewInt(42)
	chargingLimit := smartcharging.ChargingLimit{ChargingLimitSource: types.ChargingLimitSourceEMS, IsGridCritical: tests.NewBool(false)}
	chargingSchedule := types.ChargingSchedule{
		StartSchedule:          types.NewDateTime(time.Now()),
		Duration:               tests.NewInt(600),
		ChargingRateUnit:       types.ChargingRateUnitWatts,
		MinChargingRate:        tests.NewFloat(6.0),
		ChargingSchedulePeriod: []types.ChargingSchedulePeriod{types.NewChargingSchedulePeriod(0, 10.0)},
	}
	chargingSchedules := []types.ChargingSchedule{chargingSchedule}
	request := smartcharging.NewNotifyChargingLimitRequest(chargingLimit)
	request.EvseID = evseID
	request.ChargingSchedule = chargingSchedules
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"evseId":%v,"chargingLimit":{"chargingLimitSource":"%v","isGridCritical":%v},"chargingSchedule":[{"startSchedule":"%v","duration":%v,"chargingRateUnit":"%v","minChargingRate":%v,"chargingSchedulePeriod":[{"startPeriod":%v,"limit":%v}]}]}]`,
		messageId, smartcharging.NotifyChargingLimitFeatureName, *evseID, chargingLimit.ChargingLimitSource, *chargingLimit.IsGridCritical, chargingSchedule.StartSchedule.FormatTimestamp(), *chargingSchedule.Duration, chargingSchedule.ChargingRateUnit, *chargingSchedule.MinChargingRate, chargingSchedule.ChargingSchedulePeriod[0].StartPeriod, chargingSchedule.ChargingSchedulePeriod[0].Limit)
	testUnsupportedRequestFromCentralSystem(suite, request, requestJson, messageId)
}
