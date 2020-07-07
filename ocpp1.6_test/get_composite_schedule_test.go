package ocpp16_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/smartcharging"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"time"
)

// Test
func (suite *OcppV16TestSuite) TestGetCompositeScheduleRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{smartcharging.GetCompositeScheduleRequest{ConnectorId: 1, Duration: 600, ChargingRateUnit: types.ChargingRateUnitWatts}, true},
		{smartcharging.GetCompositeScheduleRequest{ConnectorId: 1, Duration: 600}, true},
		{smartcharging.GetCompositeScheduleRequest{ConnectorId: 1}, true},
		{smartcharging.GetCompositeScheduleRequest{}, true},
		{smartcharging.GetCompositeScheduleRequest{ConnectorId: -1, Duration: 600, ChargingRateUnit: types.ChargingRateUnitWatts}, false},
		{smartcharging.GetCompositeScheduleRequest{ConnectorId: 1, Duration: -1, ChargingRateUnit: types.ChargingRateUnitWatts}, false},
		{smartcharging.GetCompositeScheduleRequest{ConnectorId: 1, Duration: 600, ChargingRateUnit: "invalidChargingRateUnit"}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV16TestSuite) TestGetCompositeScheduleConfirmationValidation() {
	t := suite.T()
	chargingSchedule := types.NewChargingSchedule(types.ChargingRateUnitWatts, types.NewChargingSchedulePeriod(0, 10.0))
	var confirmationTable = []GenericTestEntry{
		{smartcharging.GetCompositeScheduleConfirmation{Status: smartcharging.GetCompositeScheduleStatusAccepted, ConnectorId: newInt(1), ScheduleStart: types.NewDateTime(time.Now()), ChargingSchedule: chargingSchedule}, true},
		{smartcharging.GetCompositeScheduleConfirmation{Status: smartcharging.GetCompositeScheduleStatusAccepted, ConnectorId: newInt(1), ScheduleStart: types.NewDateTime(time.Now())}, true},
		{smartcharging.GetCompositeScheduleConfirmation{Status: smartcharging.GetCompositeScheduleStatusAccepted, ConnectorId: newInt(1)}, true},
		{smartcharging.GetCompositeScheduleConfirmation{Status: smartcharging.GetCompositeScheduleStatusAccepted}, true},
		{smartcharging.GetCompositeScheduleConfirmation{}, false},
		{smartcharging.GetCompositeScheduleConfirmation{Status: "invalidGetCompositeScheduleStatus"}, false},
		{smartcharging.GetCompositeScheduleConfirmation{Status: smartcharging.GetCompositeScheduleStatusAccepted, ConnectorId: newInt(-1)}, false},
		{smartcharging.GetCompositeScheduleConfirmation{Status: smartcharging.GetCompositeScheduleStatusAccepted, ConnectorId: newInt(1), ChargingSchedule: types.NewChargingSchedule(types.ChargingRateUnitWatts)}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV16TestSuite) TestGetCompositeScheduleE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	connectorId := 1
	chargingRateUnit := types.ChargingRateUnitWatts
	startPeriod := 0
	limit := 10.0
	duration := 600
	status := smartcharging.GetCompositeScheduleStatusAccepted
	scheduleStart := types.NewDateTime(time.Now())
	chargingSchedule := types.NewChargingSchedule(chargingRateUnit, types.NewChargingSchedulePeriod(startPeriod, limit))
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"connectorId":%v,"duration":%v,"chargingRateUnit":"%v"}]`,
		messageId, smartcharging.GetCompositeScheduleFeatureName, connectorId, duration, chargingRateUnit)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v","connectorId":%v,"scheduleStart":"%v","chargingSchedule":{"chargingRateUnit":"%v","chargingSchedulePeriod":[{"startPeriod":%v,"limit":%v}]}}]`,
		messageId, status, connectorId, scheduleStart.FormatTimestamp(), chargingRateUnit, startPeriod, limit)
	getCompositeScheduleConfirmation := smartcharging.NewGetCompositeScheduleConfirmation(status)
	getCompositeScheduleConfirmation.ChargingSchedule = chargingSchedule
	getCompositeScheduleConfirmation.ScheduleStart = scheduleStart
	getCompositeScheduleConfirmation.ConnectorId = &connectorId
	channel := NewMockWebSocket(wsId)

	smartChargingListener := MockChargePointSmartChargingListener{}
	smartChargingListener.On("OnGetCompositeSchedule", mock.Anything).Return(getCompositeScheduleConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*smartcharging.GetCompositeScheduleRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, connectorId, request.ConnectorId)
		assert.Equal(t, duration, request.Duration)
		assert.Equal(t, chargingRateUnit, request.ChargingRateUnit)
	})
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	suite.chargePoint.SetSmartChargingHandler(smartChargingListener)
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.centralSystem.GetCompositeSchedule(wsId, func(confirmation *smartcharging.GetCompositeScheduleConfirmation, err error) {
		require.Nil(t, err)
		require.NotNil(t, confirmation)
		assert.Equal(t, status, confirmation.Status)
		require.NotNil(t, confirmation.ConnectorId)
		assert.Equal(t, connectorId, *confirmation.ConnectorId)
		assert.Equal(t, scheduleStart.FormatTimestamp(), confirmation.ScheduleStart.FormatTimestamp())
		assert.Equal(t, chargingSchedule.ChargingRateUnit, confirmation.ChargingSchedule.ChargingRateUnit)
		assert.Equal(t, chargingSchedule.Duration, confirmation.ChargingSchedule.Duration)
		assert.Equal(t, chargingSchedule.MinChargingRate, confirmation.ChargingSchedule.MinChargingRate)
		assert.Equal(t, chargingSchedule.StartSchedule, confirmation.ChargingSchedule.StartSchedule)
		assert.Equal(t, 1, len(confirmation.ChargingSchedule.ChargingSchedulePeriod))
		assert.Equal(t, chargingSchedule.ChargingSchedulePeriod[0].StartPeriod, confirmation.ChargingSchedule.ChargingSchedulePeriod[0].StartPeriod)
		assert.Equal(t, chargingSchedule.ChargingSchedulePeriod[0].Limit, confirmation.ChargingSchedule.ChargingSchedulePeriod[0].Limit)
		assert.Equal(t, chargingSchedule.ChargingSchedulePeriod[0].NumberPhases, confirmation.ChargingSchedule.ChargingSchedulePeriod[0].NumberPhases)
		resultChannel <- true
	}, connectorId, duration, func(request *smartcharging.GetCompositeScheduleRequest) {
		request.ChargingRateUnit = chargingRateUnit
	})
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV16TestSuite) TestGetCompositeScheduleInvalidEndpoint() {
	messageId := defaultMessageId
	connectorId := 1
	chargingRateUnit := types.ChargingRateUnitWatts
	duration := 600
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"connectorId":%v,"duration":%v,"chargingRateUnit":"%v"}]`,
		messageId, smartcharging.GetCompositeScheduleFeatureName, connectorId, duration, chargingRateUnit)
	GetCompositeScheduleRequest := smartcharging.NewGetCompositeScheduleRequest(connectorId, duration)
	testUnsupportedRequestFromChargePoint(suite, GetCompositeScheduleRequest, requestJson, messageId)
}
