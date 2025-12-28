package ocpp16_test

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/smartcharging"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/stretchr/testify/mock"
)

// Test
func (suite *OcppV16TestSuite) TestSetChargingProfileRequestValidation() {
	chargingSchedule := types.NewChargingSchedule(types.ChargingRateUnitWatts, types.NewChargingSchedulePeriod(0, 10.0))
	chargingProfile := types.NewChargingProfile(1, 1, types.ChargingProfilePurposeChargePointMaxProfile, types.ChargingProfileKindAbsolute, chargingSchedule)
	requestTable := []GenericTestEntry{
		{smartcharging.SetChargingProfileRequest{ConnectorId: 1, ChargingProfile: chargingProfile}, true},
		{smartcharging.SetChargingProfileRequest{ChargingProfile: chargingProfile}, true},
		{smartcharging.SetChargingProfileRequest{}, false},
		{smartcharging.SetChargingProfileRequest{ConnectorId: 1}, false},
		{smartcharging.SetChargingProfileRequest{ConnectorId: -1, ChargingProfile: chargingProfile}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV16TestSuite) TestSetChargingProfileConfirmationValidation() {
	confirmationTable := []GenericTestEntry{
		{smartcharging.SetChargingProfileConfirmation{Status: smartcharging.ChargingProfileStatusAccepted}, true},
		{smartcharging.SetChargingProfileConfirmation{Status: "invalidChargingProfileStatus"}, false},
		{smartcharging.SetChargingProfileConfirmation{}, false},
	}
	ExecuteGenericTestTable(suite, confirmationTable)
}

func (suite *OcppV16TestSuite) TestSetChargingProfileE2EMocked() {
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	connectorId := 1
	chargingProfileId := 1
	stackLevel := 1
	chargingProfilePurpose := types.ChargingProfilePurposeChargePointMaxProfile
	chargingProfileKind := types.ChargingProfileKindAbsolute
	chargingRateUnit := types.ChargingRateUnitWatts
	startPeriod := 0
	limit := 10.0
	status := smartcharging.ChargingProfileStatusAccepted
	chargingSchedule := types.NewChargingSchedule(chargingRateUnit, types.NewChargingSchedulePeriod(startPeriod, limit))
	chargingProfile := types.NewChargingProfile(chargingProfileId, stackLevel, chargingProfilePurpose, chargingProfileKind, chargingSchedule)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"connectorId":%v,"csChargingProfiles":{"chargingProfileId":%v,"stackLevel":%v,"chargingProfilePurpose":"%v","chargingProfileKind":"%v","chargingSchedule":{"chargingRateUnit":"%v","chargingSchedulePeriod":[{"startPeriod":%v,"limit":%v}]}}}]`,
		messageId,
		smartcharging.SetChargingProfileFeatureName,
		connectorId,
		chargingProfileId,
		stackLevel,
		chargingProfilePurpose,
		chargingProfileKind,
		chargingRateUnit,
		startPeriod,
		limit)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	SetChargingProfileConfirmation := smartcharging.NewSetChargingProfileConfirmation(status)
	channel := NewMockWebSocket(wsId)

	smartChargingListener := &MockChargePointSmartChargingListener{}
	smartChargingListener.On("OnSetChargingProfile", mock.Anything).Return(SetChargingProfileConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*smartcharging.SetChargingProfileRequest)
		suite.Require().True(ok)
		suite.Require().NotNil(request)
		suite.Equal(connectorId, request.ConnectorId)
		suite.Equal(chargingProfileId, request.ChargingProfile.ChargingProfileId)
		suite.Equal(chargingProfileKind, request.ChargingProfile.ChargingProfileKind)
		suite.Equal(chargingProfilePurpose, request.ChargingProfile.ChargingProfilePurpose)
		suite.Equal(types.RecurrencyKindType(""), request.ChargingProfile.RecurrencyKind)
		suite.Equal(stackLevel, request.ChargingProfile.StackLevel)
		suite.Equal(0, request.ChargingProfile.TransactionId)
		suite.Nil(request.ChargingProfile.ValidFrom)
		suite.Nil(request.ChargingProfile.ValidTo)
		suite.Equal(chargingRateUnit, request.ChargingProfile.ChargingSchedule.ChargingRateUnit)
		suite.Nil(request.ChargingProfile.ChargingSchedule.MinChargingRate)
		suite.Nil(request.ChargingProfile.ChargingSchedule.Duration)
		suite.Nil(request.ChargingProfile.ChargingSchedule.StartSchedule)
		suite.Require().Len(request.ChargingProfile.ChargingSchedule.ChargingSchedulePeriod, 1)
		suite.Equal(limit, request.ChargingProfile.ChargingSchedule.ChargingSchedulePeriod[0].Limit)
		suite.Equal(startPeriod, request.ChargingProfile.ChargingSchedule.ChargingSchedulePeriod[0].StartPeriod)
		suite.Nil(request.ChargingProfile.ChargingSchedule.ChargingSchedulePeriod[0].NumberPhases)
	})
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	suite.chargePoint.SetSmartChargingHandler(smartChargingListener)
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	suite.Require().Nil(err)
	resultChannel := make(chan bool, 1)
	err = suite.centralSystem.SetChargingProfile(wsId, func(confirmation *smartcharging.SetChargingProfileConfirmation, err error) {
		suite.Require().Nil(err)
		suite.Require().NotNil(confirmation)
		suite.Equal(status, confirmation.Status)
		resultChannel <- true
	}, connectorId, chargingProfile)
	suite.Require().Nil(err)
	result := <-resultChannel
	suite.True(result)
}

func (suite *OcppV16TestSuite) TestSetChargingProfileInvalidEndpoint() {
	messageId := defaultMessageId
	connectorId := 1
	chargingProfileId := 1
	stackLevel := 1
	chargingProfilePurpose := types.ChargingProfilePurposeChargePointMaxProfile
	chargingProfileKind := types.ChargingProfileKindAbsolute
	chargingRateUnit := types.ChargingRateUnitWatts
	startPeriod := 0
	limit := 10.0
	chargingSchedule := types.NewChargingSchedule(chargingRateUnit, types.NewChargingSchedulePeriod(startPeriod, limit))
	chargingProfile := types.NewChargingProfile(chargingProfileId, stackLevel, chargingProfilePurpose, chargingProfileKind, chargingSchedule)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"connectorId":%v,"csChargingProfiles":{"chargingProfileId":%v,"stackLevel":%v,"chargingProfilePurpose":"%v","chargingProfileKind":"%v","chargingSchedule":{"chargingRateUnit":"%v","chargingSchedulePeriod":[{"startPeriod":%v,"limit":%v}]}}}]`,
		messageId,
		smartcharging.SetChargingProfileFeatureName,
		connectorId,
		chargingProfileId,
		stackLevel,
		chargingProfilePurpose,
		chargingProfileKind,
		chargingRateUnit,
		startPeriod,
		limit)
	SetChargingProfileRequest := smartcharging.NewSetChargingProfileRequest(connectorId, chargingProfile)
	testUnsupportedRequestFromChargePoint(suite, SetChargingProfileRequest, requestJson, messageId)
}
