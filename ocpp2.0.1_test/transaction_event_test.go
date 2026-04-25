package ocpp2_test

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/transactions"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

// Test
func (suite *OcppV2TestSuite) TestTransactionEventE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	timestamp := types.NewDateTime(time.Now())
	eventType := transactions.TransactionEventEnded
	triggerReason := transactions.TriggerReasonEVDeparted
	seqNo := 10
	offline := false
	phases := tests.NewInt(3)
	cableMaxCurrent := tests.NewInt(20)
	reservationID := tests.NewInt(55)
	info := transactions.Transaction{TransactionID: "42", ChargingState: transactions.ChargingStateSuspendedEV, TimeSpentCharging: tests.NewInt(1000), StoppedReason: transactions.ReasonLocal, RemoteStartID: tests.NewInt(69)}
	idToken := types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode}
	evse := types.EVSE{ID: 1}
	meterValue := types.MeterValue{Timestamp: *types.NewDateTime(time.Now()), SampledValue: []types.SampledValue{{Value: 64.0}}}
	totalCost := tests.NewFloat(8.42)
	chargingPriority := tests.NewInt(2)
	idTokenInfo := types.NewIdTokenInfo(types.AuthorizationStatusAccepted)
	messageContent := types.MessageContent{Format: types.MessageFormatUTF8, Content: "dummyContent"}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"eventType":"%v","timestamp":"%v","triggerReason":"%v","seqNo":%v,"numberOfPhasesUsed":%v,"cableMaxCurrent":%v,"reservationId":%v,"transactionInfo":{"transactionId":"%v","chargingState":"%v","timeSpentCharging":%v,"stoppedReason":"%v","remoteStartId":%v},"idToken":{"idToken":"%v","type":"%v"},"evse":{"id":%v},"meterValue":[{"timestamp":"%v","sampledValue":[{"value":%v}]}]}]`,
		messageId, transactions.TransactionEventFeatureName, eventType, timestamp.FormatTimestamp(), triggerReason, seqNo, *phases, *cableMaxCurrent, *reservationID, info.TransactionID, info.ChargingState, *info.TimeSpentCharging, info.StoppedReason, *info.RemoteStartID, idToken.IdToken, idToken.Type, evse.ID, meterValue.Timestamp.FormatTimestamp(), meterValue.SampledValue[0].Value)
	responseJson := fmt.Sprintf(`[3,"%v",{"totalCost":%v,"chargingPriority":%v,"idTokenInfo":{"status":"%v"},"updatedPersonalMessage":{"format":"%v","content":"%v"}}]`,
		messageId, *totalCost, *chargingPriority, idTokenInfo.Status, messageContent.Format, messageContent.Content)
	transactionResponse := transactions.NewTransactionEventResponse()
	transactionResponse.TotalCost = totalCost
	transactionResponse.ChargingPriority = chargingPriority
	transactionResponse.IDTokenInfo = idTokenInfo
	transactionResponse.UpdatedPersonalMessage = &messageContent
	channel := NewMockWebSocket(wsId)

	handler := &MockCSMSTransactionsHandler{}
	handler.On("OnTransactionEvent", mock.AnythingOfType("string"), mock.Anything).Return(transactionResponse, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*transactions.TransactionEventRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, eventType, request.EventType)
		assertDateTimeEquality(t, timestamp, request.Timestamp)
		assert.Equal(t, triggerReason, request.TriggerReason)
		assert.Equal(t, seqNo, request.SequenceNo)
		assert.Equal(t, offline, request.Offline)
		assert.Equal(t, *phases, *request.NumberOfPhasesUsed)
		assert.Equal(t, *cableMaxCurrent, *request.CableMaxCurrent)
		assert.Equal(t, *reservationID, *request.ReservationID)
		assert.Equal(t, eventType, request.EventType)
		assert.Equal(t, info.TransactionID, request.TransactionInfo.TransactionID)
		assert.Equal(t, info.StoppedReason, request.TransactionInfo.StoppedReason)
		assert.Equal(t, info.ChargingState, request.TransactionInfo.ChargingState)
		assert.Equal(t, *info.TimeSpentCharging, *request.TransactionInfo.TimeSpentCharging)
		assert.Equal(t, *info.RemoteStartID, *request.TransactionInfo.RemoteStartID)
		require.NotNil(t, request.IDToken)
		assert.Equal(t, idToken.IdToken, request.IDToken.IdToken)
		assert.Equal(t, idToken.Type, request.IDToken.Type)
		require.NotNil(t, request.Evse)
		assert.Equal(t, evse.ID, request.Evse.ID)
		require.Len(t, request.MeterValue, 1)
		assertDateTimeEquality(t, &meterValue.Timestamp, &request.MeterValue[0].Timestamp)
		require.Len(t, request.MeterValue[0].SampledValue, 1)
		assert.Equal(t, meterValue.SampledValue[0].Value, request.MeterValue[0].SampledValue[0].Value)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.NoError(t, err)
	response, err := suite.chargingStation.TransactionEvent(eventType, timestamp, triggerReason, seqNo, info, func(request *transactions.TransactionEventRequest) {
		request.MeterValue = []types.MeterValue{meterValue}
		request.Evse = &evse
		request.IDToken = &idToken
		request.NumberOfPhasesUsed = phases
		request.CableMaxCurrent = cableMaxCurrent
		request.ReservationID = reservationID
		request.Offline = offline
	})
	require.NoError(t, err)
	require.NotNil(t, response)
	assert.Equal(t, *totalCost, *response.TotalCost)
	assert.Equal(t, *chargingPriority, *response.ChargingPriority)
	require.NotNil(t, response.IDTokenInfo)
	assert.Equal(t, idTokenInfo.Status, response.IDTokenInfo.Status)
	require.NotNil(t, response.UpdatedPersonalMessage)
	assert.Equal(t, messageContent.Format, response.UpdatedPersonalMessage.Format)
	assert.Equal(t, messageContent.Content, response.UpdatedPersonalMessage.Content)
}

func (suite *OcppV2TestSuite) TestTransactionEventInvalidEndpoint() {
	messageId := defaultMessageId
	timestamp := types.NewDateTime(time.Now())
	eventType := transactions.TransactionEventEnded
	triggerReason := transactions.TriggerReasonEVDeparted
	seqNo := 10
	phases := tests.NewInt(3)
	cableMaxCurrent := tests.NewInt(20)
	reservationID := tests.NewInt(55)
	info := transactions.Transaction{TransactionID: "42", ChargingState: transactions.ChargingStateSuspendedEV, TimeSpentCharging: tests.NewInt(1000), StoppedReason: transactions.ReasonLocal, RemoteStartID: tests.NewInt(69)}
	idToken := types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode}
	evse := types.EVSE{ID: 1}
	meterValue := types.MeterValue{Timestamp: *types.NewDateTime(time.Now()), SampledValue: []types.SampledValue{{Value: 64.0}}}
	request := transactions.NewTransactionEventRequest(eventType, timestamp, triggerReason, seqNo, info)
	request.NumberOfPhasesUsed = phases
	request.CableMaxCurrent = cableMaxCurrent
	request.ReservationID = reservationID
	request.IDToken = &idToken
	request.Evse = &evse
	request.MeterValue = []types.MeterValue{meterValue}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"eventType":"%v","timestamp":"%v","triggerReason":"%v","seqNo":%v,"numberOfPhasesUsed":%v,"cableMaxCurrent":%v,"reservationId":%v,"transactionInfo":{"transactionId":"%v","chargingState":"%v","timeSpentCharging":%v,"stoppedReason":"%v","remoteStartId":%v},"idToken":{"idToken":"%v","type":"%v"},"evse":{"id":%v},"meterValue":[{"timestamp":"%v","sampledValue":[{"value":%v}]}]}]`,
		messageId, transactions.TransactionEventFeatureName, eventType, timestamp.FormatTimestamp(), triggerReason, seqNo, *phases, *cableMaxCurrent, *reservationID, info.TransactionID, info.ChargingState, *info.TimeSpentCharging, info.StoppedReason, *info.RemoteStartID, idToken.IdToken, idToken.Type, evse.ID, meterValue.Timestamp.FormatTimestamp(), meterValue.SampledValue[0].Value)
	testUnsupportedRequestFromCentralSystem(suite, request, requestJson, messageId)
}
