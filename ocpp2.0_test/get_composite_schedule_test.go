package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"time"
)

// Test
func (suite *OcppV2TestSuite) TestGetCompositeScheduleRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{ocpp2.GetCompositeScheduleRequest{Duration: 600, EvseID: 1, ChargingRateUnit: ocpp2.ChargingRateUnitWatts}, true},
		{ocpp2.GetCompositeScheduleRequest{Duration: 600, EvseID: 1}, true},
		{ocpp2.GetCompositeScheduleRequest{EvseID: 1}, true},
		{ocpp2.GetCompositeScheduleRequest{}, true},
		{ocpp2.GetCompositeScheduleRequest{Duration: 600, EvseID: -1, ChargingRateUnit: ocpp2.ChargingRateUnitWatts}, false},
		{ocpp2.GetCompositeScheduleRequest{Duration: -1, EvseID: 1, ChargingRateUnit: ocpp2.ChargingRateUnitWatts}, false},
		{ocpp2.GetCompositeScheduleRequest{Duration: 600, EvseID: 1, ChargingRateUnit: "invalidChargingRateUnit"}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestGetCompositeScheduleConfirmationValidation() {
	t := suite.T()
	chargingSchedule := ocpp2.NewChargingSchedule(ocpp2.ChargingRateUnitWatts, ocpp2.NewChargingSchedulePeriod(0, 10.0))
	chargingSchedule.Duration = newInt(600)
	chargingSchedule.MinChargingRate = newFloat(6.0)
	chargingSchedule.StartSchedule = ocpp2.NewDateTime(time.Now())
	compositeSchedule := ocpp2.CompositeSchedule{
		StartDateTime:    ocpp2.NewDateTime(time.Now()),
		ChargingSchedule: chargingSchedule,
	}
	var confirmationTable = []GenericTestEntry{
		{ocpp2.GetCompositeScheduleConfirmation{Status: ocpp2.GetCompositeScheduleStatusAccepted, EvseID: 1, Schedule: &compositeSchedule}, true},
		{ocpp2.GetCompositeScheduleConfirmation{Status: ocpp2.GetCompositeScheduleStatusAccepted, EvseID: 1, Schedule: &ocpp2.CompositeSchedule{}}, true},
		{ocpp2.GetCompositeScheduleConfirmation{Status: ocpp2.GetCompositeScheduleStatusAccepted, EvseID: 1}, true},
		{ocpp2.GetCompositeScheduleConfirmation{Status: ocpp2.GetCompositeScheduleStatusAccepted}, true},
		{ocpp2.GetCompositeScheduleConfirmation{}, false},
		{ocpp2.GetCompositeScheduleConfirmation{Status: "invalidGetCompositeScheduleStatus"}, false},
		{ocpp2.GetCompositeScheduleConfirmation{Status: ocpp2.GetCompositeScheduleStatusAccepted, EvseID: -1}, false},
		{ocpp2.GetCompositeScheduleConfirmation{Status: ocpp2.GetCompositeScheduleStatusAccepted, EvseID: 1, Schedule: &ocpp2.CompositeSchedule{StartDateTime: ocpp2.NewDateTime(time.Now()), ChargingSchedule: ocpp2.NewChargingSchedule("invalidChargingRateUnit")}}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestGetCompositeScheduleE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	evseID := 1
	chargingRateUnit := ocpp2.ChargingRateUnitWatts
	duration := 600
	status := ocpp2.GetCompositeScheduleStatusAccepted
	scheduleStart := ocpp2.NewDateTime(time.Now())
	chargingSchedulePeriod := ocpp2.NewChargingSchedulePeriod(0, 10.0)
	chargingSchedulePeriod.NumberPhases = newInt(3)
	chargingSchedule := ocpp2.NewChargingSchedule(chargingRateUnit, chargingSchedulePeriod)
	chargingSchedule.Duration = newInt(600)
	chargingSchedule.StartSchedule = ocpp2.NewDateTime(time.Now())
	chargingSchedule.MinChargingRate = newFloat(6.0)
	compositeSchedule := ocpp2.CompositeSchedule{StartDateTime: scheduleStart, ChargingSchedule: chargingSchedule}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"duration":%v,"chargingRateUnit":"%v","evseId":%v}]`,
		messageId, ocpp2.GetCompositeScheduleFeatureName, duration, chargingRateUnit, evseID)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v","evseId":%v,"schedule":{"startDateTime":"%v","chargingSchedule":{"startSchedule":"%v","duration":%v,"chargingRateUnit":"%v","minChargingRate":%v,"chargingSchedulePeriod":[{"startPeriod":%v,"limit":%v,"numberPhases":%v}]}}}]`,
		messageId, status, evseID, ocpp2.FormatTimestamp(compositeSchedule.StartDateTime.Time), ocpp2.FormatTimestamp(chargingSchedule.StartSchedule.Time), *chargingSchedule.Duration, chargingSchedule.ChargingRateUnit, *chargingSchedule.MinChargingRate, chargingSchedulePeriod.StartPeriod, chargingSchedulePeriod.Limit, *chargingSchedulePeriod.NumberPhases)
	getCompositeScheduleConfirmation := ocpp2.NewGetCompositeScheduleConfirmation(status, evseID)
	getCompositeScheduleConfirmation.Schedule = &compositeSchedule
	channel := NewMockWebSocket(wsId)

	coreListener := MockChargePointCoreListener{}
	coreListener.On("OnGetCompositeSchedule", mock.Anything).Return(getCompositeScheduleConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*ocpp2.GetCompositeScheduleRequest)
		assert.True(t, ok)
		assert.NotNil(t, request)
		assert.Equal(t, duration, request.Duration)
		assert.Equal(t, chargingRateUnit, request.ChargingRateUnit)
		assert.Equal(t, evseID, request.EvseID)
	})
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, coreListener, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	assert.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.GetCompositeSchedule(wsId, func(confirmation *ocpp2.GetCompositeScheduleConfirmation, err error) {
		require.Nil(t, err)
		require.NotNil(t, confirmation)
		assert.Equal(t, status, confirmation.Status)
		assert.Equal(t, evseID, confirmation.EvseID)
		require.NotNil(t, confirmation.Schedule)
		require.NotNil(t, confirmation.Schedule.StartDateTime)
		assert.Equal(t, ocpp2.FormatTimestamp(compositeSchedule.StartDateTime.Time), ocpp2.FormatTimestamp(confirmation.Schedule.StartDateTime.Time))
		require.NotNil(t, confirmation.Schedule.ChargingSchedule)
		assert.Equal(t, chargingSchedule.ChargingRateUnit, confirmation.Schedule.ChargingSchedule.ChargingRateUnit)
		require.NotNil(t, confirmation.Schedule.ChargingSchedule.Duration)
		assert.Equal(t, *chargingSchedule.Duration, *confirmation.Schedule.ChargingSchedule.Duration)
		require.NotNil(t, confirmation.Schedule.ChargingSchedule.MinChargingRate)
		assert.Equal(t, *chargingSchedule.MinChargingRate, *confirmation.Schedule.ChargingSchedule.MinChargingRate)
		require.NotNil(t, confirmation.Schedule.ChargingSchedule.StartSchedule)
		assert.Equal(t, ocpp2.FormatTimestamp(chargingSchedule.StartSchedule.Time), ocpp2.FormatTimestamp(confirmation.Schedule.ChargingSchedule.StartSchedule.Time))
		require.Len(t, confirmation.Schedule.ChargingSchedule.ChargingSchedulePeriod, len(chargingSchedule.ChargingSchedulePeriod))
		assert.Equal(t, chargingSchedule.ChargingSchedulePeriod[0].Limit, confirmation.Schedule.ChargingSchedule.ChargingSchedulePeriod[0].Limit)
		assert.Equal(t, chargingSchedule.ChargingSchedulePeriod[0].StartPeriod, confirmation.Schedule.ChargingSchedule.ChargingSchedulePeriod[0].StartPeriod)
		require.NotNil(t, confirmation.Schedule.ChargingSchedule.ChargingSchedulePeriod[0].NumberPhases)
		assert.Equal(t, *chargingSchedule.ChargingSchedulePeriod[0].NumberPhases, *confirmation.Schedule.ChargingSchedule.ChargingSchedulePeriod[0].NumberPhases)
		resultChannel <- true
	}, duration, evseID, func(request *ocpp2.GetCompositeScheduleRequest) {
		request.ChargingRateUnit = chargingRateUnit
	})
	assert.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV2TestSuite) TestGetCompositeScheduleInvalidEndpoint() {
	messageId := defaultMessageId
	evseID := 1
	chargingRateUnit := ocpp2.ChargingRateUnitWatts
	duration := 600
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"duration":%v,"chargingRateUnit":"%v","evseId":%v}]`,
		messageId, ocpp2.GetCompositeScheduleFeatureName, duration, chargingRateUnit, evseID)
	getCompositeScheduleRequest := ocpp2.NewGetCompositeScheduleRequest(evseID, duration)
	getCompositeScheduleRequest.ChargingRateUnit = chargingRateUnit
	testUnsupportedRequestFromChargePoint(suite, getCompositeScheduleRequest, requestJson, messageId)
}
