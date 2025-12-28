package ocpp16_test

import (
	"fmt"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/smartcharging"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/stretchr/testify/mock"
)

// Test
func (suite *OcppV16TestSuite) TestGetCompositeScheduleRequestValidation() {
	requestTable := []GenericTestEntry{
		{smartcharging.GetCompositeScheduleRequest{ConnectorId: 1, Duration: 600, ChargingRateUnit: types.ChargingRateUnitWatts}, true},
		{smartcharging.GetCompositeScheduleRequest{ConnectorId: 1, Duration: 600}, true},
		{smartcharging.GetCompositeScheduleRequest{ConnectorId: 1}, true},
		{smartcharging.GetCompositeScheduleRequest{ConnectorId: 0}, true},
		{smartcharging.GetCompositeScheduleRequest{}, true},
		{smartcharging.GetCompositeScheduleRequest{ConnectorId: -1, Duration: 600, ChargingRateUnit: types.ChargingRateUnitWatts}, false},
		{smartcharging.GetCompositeScheduleRequest{ConnectorId: 1, Duration: -1, ChargingRateUnit: types.ChargingRateUnitWatts}, false},
		{smartcharging.GetCompositeScheduleRequest{ConnectorId: 1, Duration: 600, ChargingRateUnit: "invalidChargingRateUnit"}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV16TestSuite) TestGetCompositeScheduleConfirmationValidation() {
	chargingSchedule := types.NewChargingSchedule(types.ChargingRateUnitWatts, types.NewChargingSchedulePeriod(0, 10.0))
	confirmationTable := []GenericTestEntry{
		{smartcharging.GetCompositeScheduleConfirmation{Status: smartcharging.GetCompositeScheduleStatusAccepted, ConnectorId: newInt(1), ScheduleStart: types.NewDateTime(time.Now()), ChargingSchedule: chargingSchedule}, true},
		{smartcharging.GetCompositeScheduleConfirmation{Status: smartcharging.GetCompositeScheduleStatusAccepted, ConnectorId: newInt(1), ScheduleStart: types.NewDateTime(time.Now())}, true},
		{smartcharging.GetCompositeScheduleConfirmation{Status: smartcharging.GetCompositeScheduleStatusAccepted, ConnectorId: newInt(1)}, true},
		{smartcharging.GetCompositeScheduleConfirmation{Status: smartcharging.GetCompositeScheduleStatusAccepted, ConnectorId: newInt(0)}, true},
		{smartcharging.GetCompositeScheduleConfirmation{Status: smartcharging.GetCompositeScheduleStatusAccepted}, true},
		{smartcharging.GetCompositeScheduleConfirmation{}, false},
		{smartcharging.GetCompositeScheduleConfirmation{Status: "invalidGetCompositeScheduleStatus"}, false},
		{smartcharging.GetCompositeScheduleConfirmation{Status: smartcharging.GetCompositeScheduleStatusAccepted, ConnectorId: newInt(-1)}, false},
		{smartcharging.GetCompositeScheduleConfirmation{Status: smartcharging.GetCompositeScheduleStatusAccepted, ConnectorId: newInt(1), ChargingSchedule: types.NewChargingSchedule(types.ChargingRateUnitWatts)}, false},
	}
	ExecuteGenericTestTable(suite, confirmationTable)
}

func (suite *OcppV16TestSuite) TestGetCompositeScheduleE2EMocked() {
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

	smartChargingListener := &MockChargePointSmartChargingListener{}
	smartChargingListener.On("OnGetCompositeSchedule", mock.Anything).Return(getCompositeScheduleConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*smartcharging.GetCompositeScheduleRequest)
		suite.Require().True(ok)
		suite.Require().NotNil(request)
		suite.Equal(connectorId, request.ConnectorId)
		suite.Equal(duration, request.Duration)
		suite.Equal(chargingRateUnit, request.ChargingRateUnit)
	})
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	suite.chargePoint.SetSmartChargingHandler(smartChargingListener)
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	suite.Require().Nil(err)
	resultChannel := make(chan bool, 1)
	err = suite.centralSystem.GetCompositeSchedule(wsId, func(confirmation *smartcharging.GetCompositeScheduleConfirmation, err error) {
		suite.Require().Nil(err)
		suite.Require().NotNil(confirmation)
		suite.Equal(status, confirmation.Status)
		suite.Require().NotNil(confirmation.ConnectorId)
		suite.Equal(connectorId, *confirmation.ConnectorId)
		suite.Equal(scheduleStart.FormatTimestamp(), confirmation.ScheduleStart.FormatTimestamp())
		suite.Equal(chargingSchedule.ChargingRateUnit, confirmation.ChargingSchedule.ChargingRateUnit)
		suite.Equal(chargingSchedule.Duration, confirmation.ChargingSchedule.Duration)
		suite.Equal(chargingSchedule.MinChargingRate, confirmation.ChargingSchedule.MinChargingRate)
		suite.Equal(chargingSchedule.StartSchedule, confirmation.ChargingSchedule.StartSchedule)
		suite.Equal(1, len(confirmation.ChargingSchedule.ChargingSchedulePeriod))
		suite.Equal(chargingSchedule.ChargingSchedulePeriod[0].StartPeriod, confirmation.ChargingSchedule.ChargingSchedulePeriod[0].StartPeriod)
		suite.Equal(chargingSchedule.ChargingSchedulePeriod[0].Limit, confirmation.ChargingSchedule.ChargingSchedulePeriod[0].Limit)
		suite.Equal(chargingSchedule.ChargingSchedulePeriod[0].NumberPhases, confirmation.ChargingSchedule.ChargingSchedulePeriod[0].NumberPhases)
		resultChannel <- true
	}, connectorId, duration, func(request *smartcharging.GetCompositeScheduleRequest) {
		request.ChargingRateUnit = chargingRateUnit
	})
	suite.Require().Nil(err)
	result := <-resultChannel
	suite.True(result)
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
