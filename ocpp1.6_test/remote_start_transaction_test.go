package ocpp16_test

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/stretchr/testify/mock"
)

// Test
func (suite *OcppV16TestSuite) TestRemoteStartTransactionRequestValidation() {
	chargingSchedule := types.NewChargingSchedule(types.ChargingRateUnitWatts, types.NewChargingSchedulePeriod(0, 10.0))
	chargingProfile := types.NewChargingProfile(1, 1, types.ChargingProfilePurposeChargePointMaxProfile, types.ChargingProfileKindAbsolute, chargingSchedule)
	var requestTable = []GenericTestEntry{
		{core.RemoteStartTransactionRequest{IdTag: "12345", ConnectorId: newInt(1), ChargingProfile: chargingProfile}, true},
		{core.RemoteStartTransactionRequest{IdTag: "12345", ConnectorId: newInt(1)}, true},
		{core.RemoteStartTransactionRequest{IdTag: "12345"}, true},
		{core.RemoteStartTransactionRequest{IdTag: "12345", ConnectorId: newInt(-1)}, false},
		{core.RemoteStartTransactionRequest{}, false},
		{core.RemoteStartTransactionRequest{IdTag: ">20..................", ConnectorId: newInt(1)}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV16TestSuite) TestRemoteStartTransactionConfirmationValidation() {
	var confirmationTable = []GenericTestEntry{
		{core.RemoteStartTransactionConfirmation{Status: types.RemoteStartStopStatusAccepted}, true},
		{core.RemoteStartTransactionConfirmation{Status: types.RemoteStartStopStatusRejected}, true},
		{core.RemoteStartTransactionConfirmation{Status: "invalidRemoteStartTransactionStatus"}, false},
		{core.RemoteStartTransactionConfirmation{}, false},
	}
	ExecuteGenericTestTable(suite, confirmationTable)
}

func (suite *OcppV16TestSuite) TestRemoteStartTransactionE2EMocked() {
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	idTag := "12345"
	connectorId := newInt(1)
	chargingProfileId := 1
	stackLevel := 1
	chargingProfilePurpose := types.ChargingProfilePurposeChargePointMaxProfile
	chargingProfileKind := types.ChargingProfileKindAbsolute
	chargingRateUnit := types.ChargingRateUnitWatts
	startPeriod := 0
	limit := 10.0
	status := types.RemoteStartStopStatusAccepted
	chargingSchedule := types.NewChargingSchedule(chargingRateUnit, types.NewChargingSchedulePeriod(startPeriod, limit))
	chargingProfile := types.NewChargingProfile(chargingProfileId, stackLevel, chargingProfilePurpose, chargingProfileKind, chargingSchedule)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"connectorId":%v,"idTag":"%v","chargingProfile":{"chargingProfileId":%v,"stackLevel":%v,"chargingProfilePurpose":"%v","chargingProfileKind":"%v","chargingSchedule":{"chargingRateUnit":"%v","chargingSchedulePeriod":[{"startPeriod":%v,"limit":%v}]}}}]`,
		messageId,
		core.RemoteStartTransactionFeatureName,
		*connectorId,
		idTag,
		chargingProfileId,
		stackLevel,
		chargingProfilePurpose,
		chargingProfileKind,
		chargingRateUnit,
		startPeriod,
		limit)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	RemoteStartTransactionConfirmation := core.NewRemoteStartTransactionConfirmation(status)
	channel := NewMockWebSocket(wsId)

	coreListener := &MockChargePointCoreListener{}
	coreListener.On("OnRemoteStartTransaction", mock.Anything).Return(RemoteStartTransactionConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*core.RemoteStartTransactionRequest)
		suite.Require().NotNil(request)
		suite.Require().True(ok)
		suite.Equal(*connectorId, *request.ConnectorId)
		suite.Equal(idTag, request.IdTag)
		suite.Require().NotNil(request.ChargingProfile)
		suite.Equal(chargingProfileId, request.ChargingProfile.ChargingProfileId)
		suite.Equal(stackLevel, request.ChargingProfile.StackLevel)
		suite.Equal(chargingProfilePurpose, request.ChargingProfile.ChargingProfilePurpose)
		suite.Equal(chargingProfileKind, request.ChargingProfile.ChargingProfileKind)
		suite.Equal(types.RecurrencyKindType(""), request.ChargingProfile.RecurrencyKind)
		suite.Nil(request.ChargingProfile.ValidFrom)
		suite.Nil(request.ChargingProfile.ValidTo)
		suite.Require().NotNil(request.ChargingProfile.ChargingSchedule)
		suite.Equal(chargingRateUnit, request.ChargingProfile.ChargingSchedule.ChargingRateUnit)
		suite.Require().Len(request.ChargingProfile.ChargingSchedule.ChargingSchedulePeriod, 1)
		suite.Equal(startPeriod, request.ChargingProfile.ChargingSchedule.ChargingSchedulePeriod[0].StartPeriod)
		suite.Equal(limit, request.ChargingProfile.ChargingSchedule.ChargingSchedulePeriod[0].Limit)
		suite.Nil(request.ChargingProfile.ChargingSchedule.ChargingSchedulePeriod[0].NumberPhases)
	})
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, coreListener, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	suite.Nil(err)
	resultChannel := make(chan bool, 1)
	err = suite.centralSystem.RemoteStartTransaction(wsId, func(confirmation *core.RemoteStartTransactionConfirmation, err error) {
		suite.Nil(err)
		suite.NotNil(confirmation)
		suite.Equal(status, confirmation.Status)
		resultChannel <- true
	}, idTag, func(request *core.RemoteStartTransactionRequest) {
		request.ConnectorId = connectorId
		request.ChargingProfile = chargingProfile
	})
	suite.Nil(err)
	result := <-resultChannel
	suite.True(result)
}

func (suite *OcppV16TestSuite) TestRemoteStartTransactionInvalidEndpoint() {
	messageId := defaultMessageId
	idTag := "12345"
	remoteStartTransactionRequest := core.NewRemoteStartTransactionRequest(idTag)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"idTag":"%v"}]`,
		messageId,
		core.RemoteStartTransactionFeatureName,
		idTag)
	testUnsupportedRequestFromChargePoint(suite, remoteStartTransactionRequest, requestJson, messageId)
}
