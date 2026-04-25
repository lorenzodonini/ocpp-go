package ocpp2_test

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/remotecontrol"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test
func (suite *OcppV2TestSuite) TestRequestStartTransactionE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	evseId := tests.NewInt(1)
	remoteStartID := 42
	idToken := types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode}
	schedule := []types.ChargingSchedule{
		{
			ID:               1,
			ChargingRateUnit: types.ChargingRateUnitAmperes,
			ChargingSchedulePeriod: []types.ChargingSchedulePeriod{
				{
					StartPeriod: 42,
					Limit:       16.0,
				},
			},
		},
	}
	chargingProfile := types.ChargingProfile{
		ID:                     1,
		StackLevel:             0,
		ChargingProfilePurpose: types.ChargingProfilePurposeTxProfile,
		ChargingProfileKind:    types.ChargingProfileKindAbsolute,
		ChargingSchedule:       schedule,
	}
	groupIdToken := types.IdToken{IdToken: "1234", Type: types.IdTokenTypeISO15693}
	status := remotecontrol.RequestStartStopStatusAccepted
	transactionId := "12345"
	statusInfo := types.StatusInfo{ReasonCode: "200"}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"evseId":%v,"remoteStartId":%v,"idToken":{"idToken":"%s","type":"%s"},"chargingProfile":{"id":%v,"stackLevel":%v,"chargingProfilePurpose":"%v","chargingProfileKind":"%v","chargingSchedule":[{"id":%v,"chargingRateUnit":"%v","chargingSchedulePeriod":[{"startPeriod":%v,"limit":%v}]}]},"groupIdToken":{"idToken":"%s","type":"%s"}}]`,
		messageId, remotecontrol.RequestStartTransactionFeatureName, *evseId, remoteStartID, idToken.IdToken, idToken.Type, chargingProfile.ID, chargingProfile.StackLevel, chargingProfile.ChargingProfilePurpose, chargingProfile.ChargingProfileKind, schedule[0].ID, schedule[0].ChargingRateUnit, schedule[0].ChargingSchedulePeriod[0].StartPeriod, schedule[0].ChargingSchedulePeriod[0].Limit, groupIdToken.IdToken, groupIdToken.Type)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v","transactionId":"%v","statusInfo":{"reasonCode":"%v"}}]`,
		messageId, status, transactionId, statusInfo.ReasonCode)
	requestStartTransactionResponse := remotecontrol.NewRequestStartTransactionResponse(status)
	requestStartTransactionResponse.TransactionID = transactionId
	requestStartTransactionResponse.StatusInfo = &statusInfo
	channel := NewMockWebSocket(wsId)

	handler := &MockChargingStationRemoteControlHandler{}
	handler.On("OnRequestStartTransaction", mock.Anything).Return(requestStartTransactionResponse, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*remotecontrol.RequestStartTransactionRequest)
		require.True(t, ok)
		assert.Equal(t, *evseId, *request.EvseID)
		assert.Equal(t, remoteStartID, request.RemoteStartID)
		assert.Equal(t, idToken.IdToken, request.IDToken.IdToken)
		assert.Equal(t, idToken.Type, request.IDToken.Type)
		assert.Equal(t, chargingProfile.ID, request.ChargingProfile.ID)
		assert.Equal(t, chargingProfile.ChargingProfilePurpose, request.ChargingProfile.ChargingProfilePurpose)
		assert.Equal(t, chargingProfile.ChargingProfileKind, request.ChargingProfile.ChargingProfileKind)
		require.Len(t, request.ChargingProfile.ChargingSchedule, len(chargingProfile.ChargingSchedule))
		s := request.ChargingProfile.ChargingSchedule[0]
		assert.Equal(t, chargingProfile.ChargingSchedule[0].ID, s.ID)
		assert.Equal(t, chargingProfile.ChargingSchedule[0].ChargingRateUnit, s.ChargingRateUnit)
		require.Len(t, s.ChargingSchedulePeriod, len(chargingProfile.ChargingSchedule[0].ChargingSchedulePeriod))
		assert.Equal(t, chargingProfile.ChargingSchedule[0].ChargingSchedulePeriod[0].Limit, s.ChargingSchedulePeriod[0].Limit)
		assert.Equal(t, chargingProfile.ChargingSchedule[0].ChargingSchedulePeriod[0].StartPeriod, s.ChargingSchedulePeriod[0].StartPeriod)
		require.NotNil(t, request.GroupIdToken)
		assert.Equal(t, groupIdToken.IdToken, request.GroupIdToken.IdToken)
		assert.Equal(t, groupIdToken.Type, request.GroupIdToken.Type)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.RequestStartTransaction(wsId, func(response *remotecontrol.RequestStartTransactionResponse, err error) {
		require.Nil(t, err)
		require.NotNil(t, response)
		assert.Equal(t, status, response.Status)
		assert.Equal(t, transactionId, response.TransactionID)
		assert.Equal(t, statusInfo.ReasonCode, response.StatusInfo.ReasonCode)
		resultChannel <- true
	}, remoteStartID, idToken, func(request *remotecontrol.RequestStartTransactionRequest) {
		request.EvseID = evseId
		request.ChargingProfile = &chargingProfile
		request.GroupIdToken = &groupIdToken
	})
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV2TestSuite) TestRequestStartTransactionInvalidEndpoint() {
	messageId := defaultMessageId
	evseId := tests.NewInt(1)
	remoteStartID := 42
	idToken := types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode}
	schedule := []types.ChargingSchedule{
		{
			ChargingRateUnit: types.ChargingRateUnitAmperes,
			ChargingSchedulePeriod: []types.ChargingSchedulePeriod{
				{
					StartPeriod: 42,
					Limit:       16.0,
				},
			},
		},
	}
	chargingProfile := types.ChargingProfile{
		ID:                     1,
		StackLevel:             0,
		ChargingProfilePurpose: types.ChargingProfilePurposeTxProfile,
		ChargingProfileKind:    types.ChargingProfileKindAbsolute,
		ChargingSchedule:       schedule,
	}
	groupIdToken := types.IdToken{IdToken: "1234", Type: types.IdTokenTypeISO15693}
	request := remotecontrol.RequestStartTransactionRequest{
		EvseID:          evseId,
		RemoteStartID:   remoteStartID,
		IDToken:         idToken,
		ChargingProfile: &chargingProfile,
		GroupIdToken:    &groupIdToken,
	}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"evseId":%v,"remoteStartId":%v,"idToken":{"idToken":"%s","type":"%s"},"chargingProfile":{"id":%v,"stackLevel":%v,"chargingProfilePurpose":"%v","chargingProfileKind":"%v","chargingSchedule":[{"id":%v,"chargingRateUnit":"%v","chargingSchedulePeriod":[{"startPeriod":%v,"limit":%v}]}]},"groupIdToken":{"idToken":"%s","type":"%s"}}]`,
		messageId, remotecontrol.RequestStartTransactionFeatureName, *evseId, remoteStartID, idToken.IdToken, idToken.Type, chargingProfile.ID, chargingProfile.StackLevel, chargingProfile.ChargingProfilePurpose, chargingProfile.ChargingProfileKind, schedule[0].ID, schedule[0].ChargingRateUnit, schedule[0].ChargingSchedulePeriod[0].StartPeriod, schedule[0].ChargingSchedulePeriod[0].Limit, groupIdToken.IdToken, groupIdToken.Type)
	testUnsupportedRequestFromChargingStation(suite, request, requestJson, messageId)
}
