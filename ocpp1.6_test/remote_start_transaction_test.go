package ocpp16_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test
func (suite *OcppV16TestSuite) TestRemoteStartTransactionRequestValidation() {
	t := suite.T()
	chargingSchedule := types.NewChargingSchedule(types.ChargingRateUnitWatts, types.NewChargingSchedulePeriod(0, 10.0))
	chargingProfile := types.NewChargingProfile(1, 1, types.ChargingProfilePurposeChargePointMaxProfile, types.ChargingProfileKindAbsolute, chargingSchedule)
	var requestTable = []GenericTestEntry{
		{core.RemoteStartTransactionRequest{IdTag: "12345", ConnectorId: 1, ChargingProfile: chargingProfile}, true},
		{core.RemoteStartTransactionRequest{IdTag: "12345", ConnectorId: 1}, true},
		{core.RemoteStartTransactionRequest{IdTag: "12345"}, true},
		{core.RemoteStartTransactionRequest{IdTag: "12345", ConnectorId: -1}, false},
		{core.RemoteStartTransactionRequest{}, false},
		{core.RemoteStartTransactionRequest{IdTag: ">20..................", ConnectorId: 1}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV16TestSuite) TestRemoteStartTransactionConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{core.RemoteStartTransactionConfirmation{Status: types.RemoteStartStopStatusAccepted}, true},
		{core.RemoteStartTransactionConfirmation{Status: types.RemoteStartStopStatusRejected}, true},
		{core.RemoteStartTransactionConfirmation{Status: "invalidRemoteStartTransactionStatus"}, false},
		{core.RemoteStartTransactionConfirmation{}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV16TestSuite) TestRemoteStartTransactionE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	idTag := "12345"
	connectorId := 1
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
		connectorId,
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

	coreListener := MockChargePointCoreListener{}
	coreListener.On("OnRemoteStartTransaction", mock.Anything).Return(RemoteStartTransactionConfirmation, nil)
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, coreListener, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	assert.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.centralSystem.RemoteStartTransaction(wsId, func(confirmation *core.RemoteStartTransactionConfirmation, err error) {
		assert.Nil(t, err)
		assert.NotNil(t, confirmation)
		assert.Equal(t, status, confirmation.Status)
		resultChannel <- true
	}, idTag, func(request *core.RemoteStartTransactionRequest) {
		request.ConnectorId = connectorId
		request.ChargingProfile = chargingProfile
	})
	assert.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
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
