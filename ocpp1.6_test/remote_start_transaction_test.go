package ocpp16_test

import (
	"fmt"
	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test
func (suite *OcppV16TestSuite) TestRemoteStartTransactionRequestValidation() {
	t := suite.T()
	chargingSchedule := ocpp16.NewChargingSchedule(ocpp16.ChargingRateUnitWatts, ocpp16.NewChargingSchedulePeriod(0, 10.0))
	chargingProfile := ocpp16.NewChargingProfile(1, 1, ocpp16.ChargingProfilePurposeChargePointMaxProfile, ocpp16.ChargingProfileKindAbsolute, chargingSchedule)
	var requestTable = []GenericTestEntry{
		{ocpp16.RemoteStartTransactionRequest{IdTag: "12345", ConnectorId: 1, ChargingProfile: chargingProfile}, true},
		{ocpp16.RemoteStartTransactionRequest{IdTag: "12345", ConnectorId: 1}, true},
		{ocpp16.RemoteStartTransactionRequest{IdTag: "12345"}, true},
		{ocpp16.RemoteStartTransactionRequest{IdTag: "12345", ConnectorId: -1}, false},
		{ocpp16.RemoteStartTransactionRequest{}, false},
		{ocpp16.RemoteStartTransactionRequest{IdTag: ">20..................", ConnectorId: 1}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV16TestSuite) TestRemoteStartTransactionConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{ocpp16.RemoteStartTransactionConfirmation{Status: ocpp16.RemoteStartStopStatusAccepted}, true},
		{ocpp16.RemoteStartTransactionConfirmation{Status: ocpp16.RemoteStartStopStatusRejected}, true},
		{ocpp16.RemoteStartTransactionConfirmation{Status: "invalidRemoteStartTransactionStatus"}, false},
		{ocpp16.RemoteStartTransactionConfirmation{}, false},
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
	chargingProfilePurpose := ocpp16.ChargingProfilePurposeChargePointMaxProfile
	chargingProfileKind := ocpp16.ChargingProfileKindAbsolute
	chargingRateUnit := ocpp16.ChargingRateUnitWatts
	startPeriod := 0
	limit := 10.0
	status := ocpp16.RemoteStartStopStatusAccepted
	chargingSchedule := ocpp16.NewChargingSchedule(chargingRateUnit, ocpp16.NewChargingSchedulePeriod(startPeriod, limit))
	chargingProfile := ocpp16.NewChargingProfile(chargingProfileId, stackLevel, chargingProfilePurpose, chargingProfileKind, chargingSchedule)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"connectorId":%v,"idTag":"%v","chargingProfile":{"chargingProfileId":%v,"stackLevel":%v,"chargingProfilePurpose":"%v","chargingProfileKind":"%v","chargingSchedule":{"chargingRateUnit":"%v","chargingSchedulePeriod":[{"startPeriod":%v,"limit":%v}]}}}]`,
		messageId,
		ocpp16.RemoteStartTransactionFeatureName,
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
	RemoteStartTransactionConfirmation := ocpp16.NewRemoteStartTransactionConfirmation(status)
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
	err = suite.centralSystem.RemoteStartTransaction(wsId, func(confirmation *ocpp16.RemoteStartTransactionConfirmation, err error) {
		assert.Nil(t, err)
		assert.NotNil(t, confirmation)
		assert.Equal(t, status, confirmation.Status)
		resultChannel <- true
	}, idTag, func(request *ocpp16.RemoteStartTransactionRequest) {
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
	remoteStartTransactionRequest := ocpp16.NewRemoteStartTransactionRequest(idTag)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"idTag":"%v"}]`,
		messageId,
		ocpp16.RemoteStartTransactionFeatureName,
		idTag)
	testUnsupportedRequestFromChargePoint(suite, remoteStartTransactionRequest, requestJson, messageId)
}
