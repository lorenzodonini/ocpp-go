package ocpp16_test

import (
	"fmt"
	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test
func (suite *OcppV16TestSuite) TestSetChargingProfileRequestValidation() {
	t := suite.T()
	chargingSchedule := ocpp16.NewChargingSchedule(ocpp16.ChargingRateUnitWatts, ocpp16.NewChargingSchedulePeriod(0, 10.0))
	chargingProfile := ocpp16.NewChargingProfile(1, 1, ocpp16.ChargingProfilePurposeChargePointMaxProfile, ocpp16.ChargingProfileKindAbsolute, chargingSchedule)
	var requestTable = []GenericTestEntry{
		{ocpp16.SetChargingProfileRequest{ConnectorId: 1, ChargingProfile: chargingProfile}, true},
		{ocpp16.SetChargingProfileRequest{ChargingProfile: chargingProfile}, true},
		{ocpp16.SetChargingProfileRequest{}, false},
		{ocpp16.SetChargingProfileRequest{ConnectorId: 1}, false},
		{ocpp16.SetChargingProfileRequest{ConnectorId: -1, ChargingProfile: chargingProfile}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV16TestSuite) TestSetChargingProfileConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{ocpp16.SetChargingProfileConfirmation{Status: ocpp16.ChargingProfileStatusAccepted}, true},
		{ocpp16.SetChargingProfileConfirmation{Status: "invalidChargingProfileStatus"}, false},
		{ocpp16.SetChargingProfileConfirmation{}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV16TestSuite) TestSetChargingProfileE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	connectorId := 1
	chargingProfileId := 1
	stackLevel := 1
	chargingProfilePurpose := ocpp16.ChargingProfilePurposeChargePointMaxProfile
	chargingProfileKind := ocpp16.ChargingProfileKindAbsolute
	chargingRateUnit := ocpp16.ChargingRateUnitWatts
	startPeriod := 0
	limit := 10.0
	status := ocpp16.ChargingProfileStatusAccepted
	chargingSchedule := ocpp16.NewChargingSchedule(chargingRateUnit, ocpp16.NewChargingSchedulePeriod(startPeriod, limit))
	chargingProfile := ocpp16.NewChargingProfile(chargingProfileId, stackLevel, chargingProfilePurpose, chargingProfileKind, chargingSchedule)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"connectorId":%v,"chargingProfile":{"chargingProfileId":%v,"stackLevel":%v,"chargingProfilePurpose":"%v","chargingProfileKind":"%v","chargingSchedule":{"chargingRateUnit":"%v","chargingSchedulePeriod":[{"startPeriod":%v,"limit":%v}]}}}]`,
		messageId,
		ocpp16.SetChargingProfileFeatureName,
		connectorId,
		chargingProfileId,
		stackLevel,
		chargingProfilePurpose,
		chargingProfileKind,
		chargingRateUnit,
		startPeriod,
		limit)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	SetChargingProfileConfirmation := ocpp16.NewSetChargingProfileConfirmation(status)
	channel := NewMockWebSocket(wsId)

	smartChargingListener := MockChargePointSmartChargingListener{}
	smartChargingListener.On("OnSetChargingProfile", mock.Anything).Return(SetChargingProfileConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*ocpp16.SetChargingProfileRequest)
		assert.True(t, ok)
		assert.NotNil(t, request)
		assert.Equal(t, connectorId, request.ConnectorId)
		assert.Equal(t, chargingProfileId, request.ChargingProfile.ChargingProfileId)
		assert.Equal(t, chargingProfileKind, request.ChargingProfile.ChargingProfileKind)
		assert.Equal(t, chargingProfilePurpose, request.ChargingProfile.ChargingProfilePurpose)
		assert.Equal(t, ocpp16.RecurrencyKindType(""), request.ChargingProfile.RecurrencyKind)
		assert.Equal(t, stackLevel, request.ChargingProfile.StackLevel)
		assert.Equal(t, 0, request.ChargingProfile.TransactionId)
		assert.Nil(t, request.ChargingProfile.ValidFrom)
		assert.Nil(t, request.ChargingProfile.ValidTo)
		assert.Equal(t, chargingRateUnit, request.ChargingProfile.ChargingSchedule.ChargingRateUnit)
		assert.Equal(t, 0.0, request.ChargingProfile.ChargingSchedule.MinChargingRate)
		assert.Equal(t, 0, request.ChargingProfile.ChargingSchedule.Duration)
		assert.Nil(t, request.ChargingProfile.ChargingSchedule.StartSchedule)
		assert.Equal(t, 1, len(request.ChargingProfile.ChargingSchedule.ChargingSchedulePeriod))
		assert.Equal(t, limit, request.ChargingProfile.ChargingSchedule.ChargingSchedulePeriod[0].Limit)
		assert.Equal(t, startPeriod, request.ChargingProfile.ChargingSchedule.ChargingSchedulePeriod[0].StartPeriod)
		assert.Equal(t, 0, request.ChargingProfile.ChargingSchedule.ChargingSchedulePeriod[0].NumberPhases)
	})
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	suite.chargePoint.SetSmartChargingListener(smartChargingListener)
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	assert.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.centralSystem.SetChargingProfile(wsId, func(confirmation *ocpp16.SetChargingProfileConfirmation, err error) {
		assert.Nil(t, err)
		assert.NotNil(t, confirmation)
		assert.Equal(t, status, confirmation.Status)
		resultChannel <- true
	}, connectorId, chargingProfile)
	assert.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV16TestSuite) TestSetChargingProfileInvalidEndpoint() {
	messageId := defaultMessageId
	connectorId := 1
	chargingProfileId := 1
	stackLevel := 1
	chargingProfilePurpose := ocpp16.ChargingProfilePurposeChargePointMaxProfile
	chargingProfileKind := ocpp16.ChargingProfileKindAbsolute
	chargingRateUnit := ocpp16.ChargingRateUnitWatts
	startPeriod := 0
	limit := 10.0
	chargingSchedule := ocpp16.NewChargingSchedule(chargingRateUnit, ocpp16.NewChargingSchedulePeriod(startPeriod, limit))
	chargingProfile := ocpp16.NewChargingProfile(chargingProfileId, stackLevel, chargingProfilePurpose, chargingProfileKind, chargingSchedule)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"connectorId":%v,"chargingProfile":{"chargingProfileId":%v,"stackLevel":%v,"chargingProfilePurpose":"%v","chargingProfileKind":"%v","chargingSchedule":{"chargingRateUnit":"%v","chargingSchedulePeriod":[{"startPeriod":%v,"limit":%v}]}}}]`,
		messageId,
		ocpp16.SetChargingProfileFeatureName,
		connectorId,
		chargingProfileId,
		stackLevel,
		chargingProfilePurpose,
		chargingProfileKind,
		chargingRateUnit,
		startPeriod,
		limit)
	SetChargingProfileRequest := ocpp16.NewSetChargingProfileRequest(connectorId, chargingProfile)
	testUnsupportedRequestFromChargePoint(suite, SetChargingProfileRequest, requestJson, messageId)
}
