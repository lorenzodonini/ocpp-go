package ocpp2_test

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/smartcharging"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

// Tests
func (suite *OcppV2TestSuite) TestNotifyChargingLimitRequestValidation() {
	chargingSchedule := types.ChargingSchedule{
		ID:                     1,
		StartSchedule:          types.NewDateTime(time.Now()),
		Duration:               newInt(600),
		ChargingRateUnit:       types.ChargingRateUnitWatts,
		MinChargingRate:        newFloat(6.0),
		ChargingSchedulePeriod: []types.ChargingSchedulePeriod{types.NewChargingSchedulePeriod(0, 10.0)},
	}
	var requestTable = []GenericTestEntry{
		{smartcharging.NotifyChargingLimitRequest{EvseID: newInt(1), ChargingLimit: smartcharging.ChargingLimit{ChargingLimitSource: types.ChargingLimitSourceEMS, IsGridCritical: newBool(false)}, ChargingSchedule: []types.ChargingSchedule{chargingSchedule}}, true},
		{smartcharging.NotifyChargingLimitRequest{EvseID: newInt(1), ChargingLimit: smartcharging.ChargingLimit{ChargingLimitSource: types.ChargingLimitSourceEMS, IsGridCritical: newBool(false)}, ChargingSchedule: []types.ChargingSchedule{}}, true},
		{smartcharging.NotifyChargingLimitRequest{EvseID: newInt(1), ChargingLimit: smartcharging.ChargingLimit{ChargingLimitSource: types.ChargingLimitSourceEMS, IsGridCritical: newBool(false)}}, true},
		{smartcharging.NotifyChargingLimitRequest{ChargingLimit: smartcharging.ChargingLimit{ChargingLimitSource: types.ChargingLimitSourceEMS, IsGridCritical: newBool(false)}}, true},
		{smartcharging.NotifyChargingLimitRequest{ChargingLimit: smartcharging.ChargingLimit{ChargingLimitSource: types.ChargingLimitSourceEMS}}, true},
		{smartcharging.NotifyChargingLimitRequest{ChargingLimit: smartcharging.ChargingLimit{}}, false},
		{smartcharging.NotifyChargingLimitRequest{}, false},
		{smartcharging.NotifyChargingLimitRequest{ChargingLimit: smartcharging.ChargingLimit{ChargingLimitSource: "invalidChargingLimitSource", IsGridCritical: newBool(false)}}, false},
		{smartcharging.NotifyChargingLimitRequest{EvseID: newInt(-1), ChargingLimit: smartcharging.ChargingLimit{ChargingLimitSource: types.ChargingLimitSourceEMS, IsGridCritical: newBool(false)}}, false},
		{smartcharging.NotifyChargingLimitRequest{ChargingLimit: smartcharging.ChargingLimit{ChargingLimitSource: types.ChargingLimitSourceEMS, IsGridCritical: newBool(false)}, ChargingSchedule: []types.ChargingSchedule{{ChargingRateUnit: "invalidStruct"}}}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV2TestSuite) TestNotifyChargingLimitResponseValidation() {
	var responseTable = []GenericTestEntry{
		{smartcharging.NotifyChargingLimitResponse{}, true},
	}
	ExecuteGenericTestTable(suite, responseTable)
}

func (suite *OcppV2TestSuite) TestNotifyChargingLimitE2EMocked() {
	wsId := "test_id"
	messageId := "1234"
	wsUrl := "someUrl"
	evseID := newInt(42)
	chargingLimit := smartcharging.ChargingLimit{ChargingLimitSource: types.ChargingLimitSourceEMS, IsGridCritical: newBool(false)}
	chargingSchedule := types.ChargingSchedule{
		ID:                     1,
		StartSchedule:          types.NewDateTime(time.Now()),
		Duration:               newInt(600),
		ChargingRateUnit:       types.ChargingRateUnitWatts,
		MinChargingRate:        newFloat(6.0),
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
		suite.Require().True(ok)
		suite.Require().NotNil(request)
		suite.Equal(*evseID, *request.EvseID)
		suite.Equal(chargingLimit.ChargingLimitSource, request.ChargingLimit.ChargingLimitSource)
		suite.Require().NotNil(request.ChargingLimit.IsGridCritical)
		suite.Equal(chargingLimit.IsGridCritical, request.ChargingLimit.IsGridCritical)
		suite.Require().Len(request.ChargingSchedule, len(chargingSchedules))
		assertDateTimeEquality(suite, chargingSchedule.StartSchedule, request.ChargingSchedule[0].StartSchedule)
		suite.Equal(chargingSchedule.ID, request.ChargingSchedule[0].ID)
		suite.Equal(*chargingSchedule.Duration, *request.ChargingSchedule[0].Duration)
		suite.Equal(chargingSchedule.ChargingRateUnit, request.ChargingSchedule[0].ChargingRateUnit)
		suite.Equal(*chargingSchedule.MinChargingRate, *request.ChargingSchedule[0].MinChargingRate)
		suite.Require().Len(request.ChargingSchedule[0].ChargingSchedulePeriod, len(chargingSchedule.ChargingSchedulePeriod))
		suite.Equal(chargingSchedule.ChargingSchedulePeriod[0].StartPeriod, request.ChargingSchedule[0].ChargingSchedulePeriod[0].StartPeriod)
		suite.Equal(chargingSchedule.ChargingSchedulePeriod[0].Limit, request.ChargingSchedule[0].ChargingSchedulePeriod[0].Limit)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	suite.Require().Nil(err)
	r, err := suite.chargingStation.NotifyChargingLimit(chargingLimit, func(request *smartcharging.NotifyChargingLimitRequest) {
		request.EvseID = evseID
		request.ChargingSchedule = chargingSchedules
	})
	suite.Require().Nil(err)
	suite.Require().NotNil(r)
}

func (suite *OcppV2TestSuite) TestNotifyChargingLimitInvalidEndpoint() {
	messageId := defaultMessageId
	evseID := newInt(42)
	chargingLimit := smartcharging.ChargingLimit{ChargingLimitSource: types.ChargingLimitSourceEMS, IsGridCritical: newBool(false)}
	chargingSchedule := types.ChargingSchedule{
		StartSchedule:          types.NewDateTime(time.Now()),
		Duration:               newInt(600),
		ChargingRateUnit:       types.ChargingRateUnitWatts,
		MinChargingRate:        newFloat(6.0),
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
