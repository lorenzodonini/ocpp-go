package ocpp16_test

import (
	"fmt"
	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"time"
)

// Test
func (suite *OcppV16TestSuite) TestGetCompositeScheduleRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{ocpp16.GetCompositeScheduleRequest{ConnectorId: 1, Duration: 600, ChargingRateUnit: ocpp16.ChargingRateUnitWatts}, true},
		{ocpp16.GetCompositeScheduleRequest{ConnectorId: 1, Duration: 600}, true},
		{ocpp16.GetCompositeScheduleRequest{ConnectorId: 1}, true},
		{ocpp16.GetCompositeScheduleRequest{}, true},
		{ocpp16.GetCompositeScheduleRequest{ConnectorId: -1, Duration: 600, ChargingRateUnit: ocpp16.ChargingRateUnitWatts}, false},
		{ocpp16.GetCompositeScheduleRequest{ConnectorId: 1, Duration: -1, ChargingRateUnit: ocpp16.ChargingRateUnitWatts}, false},
		{ocpp16.GetCompositeScheduleRequest{ConnectorId: 1, Duration: 600, ChargingRateUnit: "invalidChargingRateUnit"}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV16TestSuite) TestGetCompositeScheduleConfirmationValidation() {
	t := suite.T()
	chargingSchedule := ocpp16.NewChargingSchedule(ocpp16.ChargingRateUnitWatts, ocpp16.NewChargingSchedulePeriod(0, 10.0))
	var confirmationTable = []GenericTestEntry{
		{ocpp16.GetCompositeScheduleConfirmation{Status: ocpp16.GetCompositeScheduleStatusAccepted, ConnectorId: 1, ScheduleStart: ocpp16.NewDateTime(time.Now()), ChargingSchedule: chargingSchedule}, true},
		{ocpp16.GetCompositeScheduleConfirmation{Status: ocpp16.GetCompositeScheduleStatusAccepted, ConnectorId: 1, ScheduleStart: ocpp16.NewDateTime(time.Now())}, true},
		{ocpp16.GetCompositeScheduleConfirmation{Status: ocpp16.GetCompositeScheduleStatusAccepted, ConnectorId: 1}, true},
		{ocpp16.GetCompositeScheduleConfirmation{Status: ocpp16.GetCompositeScheduleStatusAccepted}, true},
		{ocpp16.GetCompositeScheduleConfirmation{}, false},
		{ocpp16.GetCompositeScheduleConfirmation{Status: "invalidGetCompositeScheduleStatus"}, false},
		{ocpp16.GetCompositeScheduleConfirmation{Status: ocpp16.GetCompositeScheduleStatusAccepted, ConnectorId: -1}, false},
		{ocpp16.GetCompositeScheduleConfirmation{Status: ocpp16.GetCompositeScheduleStatusAccepted, ConnectorId: 1, ChargingSchedule: ocpp16.NewChargingSchedule(ocpp16.ChargingRateUnitWatts)}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV16TestSuite) TestGetCompositeScheduleE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	connectorId := 1
	chargingRateUnit := ocpp16.ChargingRateUnitWatts
	startPeriod := 0
	limit := 10.0
	duration := 600
	status := ocpp16.GetCompositeScheduleStatusAccepted
	scheduleStart := ocpp16.NewDateTime(time.Now())
	chargingSchedule := ocpp16.NewChargingSchedule(chargingRateUnit, ocpp16.NewChargingSchedulePeriod(startPeriod, limit))
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"connectorId":%v,"duration":%v,"chargingRateUnit":"%v"}]`,
		messageId, ocpp16.GetCompositeScheduleFeatureName, connectorId, duration, chargingRateUnit)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v","connectorId":%v,"scheduleStart":"%v","chargingSchedule":{"chargingRateUnit":"%v","chargingSchedulePeriod":[{"startPeriod":%v,"limit":%v}]}}]`,
		messageId, status, connectorId, scheduleStart.Format(ocpp16.ISO8601), chargingRateUnit, startPeriod, limit)
	getCompositeScheduleConfirmation := ocpp16.NewGetCompositeScheduleConfirmation(status)
	getCompositeScheduleConfirmation.ChargingSchedule = chargingSchedule
	getCompositeScheduleConfirmation.ScheduleStart = scheduleStart
	getCompositeScheduleConfirmation.ConnectorId = connectorId
	channel := NewMockWebSocket(wsId)

	smartChargingListener := MockChargePointSmartChargingListener{}
	smartChargingListener.On("OnGetCompositeSchedule", mock.Anything).Return(getCompositeScheduleConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*ocpp16.GetCompositeScheduleRequest)
		assert.True(t, ok)
		assert.NotNil(t, request)
		assert.Equal(t, connectorId, request.ConnectorId)
		assert.Equal(t, duration, request.Duration)
		assert.Equal(t, chargingRateUnit, request.ChargingRateUnit)
	})
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	suite.chargePoint.SetSmartChargingListener(smartChargingListener)
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	assert.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.centralSystem.GetCompositeSchedule(wsId, func(confirmation *ocpp16.GetCompositeScheduleConfirmation, err error) {
		if !assert.Nil(t, err) || !assert.NotNil(t, confirmation) {
			resultChannel <- false
		} else {
			assert.Equal(t, status, confirmation.Status)
			assert.Equal(t, connectorId, confirmation.ConnectorId)
			assert.Equal(t, scheduleStart.Format(ocpp16.ISO8601), confirmation.ScheduleStart.Format(ocpp16.ISO8601))
			assert.Equal(t, chargingSchedule.ChargingRateUnit, confirmation.ChargingSchedule.ChargingRateUnit)
			assert.Equal(t, chargingSchedule.Duration, confirmation.ChargingSchedule.Duration)
			assert.Equal(t, chargingSchedule.MinChargingRate, confirmation.ChargingSchedule.MinChargingRate)
			assert.Equal(t, chargingSchedule.StartSchedule, confirmation.ChargingSchedule.StartSchedule)
			assert.Equal(t, 1, len(confirmation.ChargingSchedule.ChargingSchedulePeriod))
			assert.Equal(t, chargingSchedule.ChargingSchedulePeriod[0].StartPeriod, confirmation.ChargingSchedule.ChargingSchedulePeriod[0].StartPeriod)
			assert.Equal(t, chargingSchedule.ChargingSchedulePeriod[0].Limit, confirmation.ChargingSchedule.ChargingSchedulePeriod[0].Limit)
			assert.Equal(t, chargingSchedule.ChargingSchedulePeriod[0].NumberPhases, confirmation.ChargingSchedule.ChargingSchedulePeriod[0].NumberPhases)
			resultChannel <- true
		}
	}, connectorId, duration, func(request *ocpp16.GetCompositeScheduleRequest) {
		request.ChargingRateUnit = chargingRateUnit
	})
	assert.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV16TestSuite) TestGetCompositeScheduleInvalidEndpoint() {
	messageId := defaultMessageId
	connectorId := 1
	chargingRateUnit := ocpp16.ChargingRateUnitWatts
	duration := 600
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"connectorId":%v,"duration":%v,"chargingRateUnit":"%v"}]`,
		messageId, ocpp16.GetCompositeScheduleFeatureName, connectorId, duration, chargingRateUnit)
	GetCompositeScheduleRequest := ocpp16.NewGetCompositeScheduleRequest(connectorId, duration)
	testUnsupportedRequestFromChargePoint(suite, GetCompositeScheduleRequest, requestJson, messageId)
}
