package ocpp2_test

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/transactions"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

// Test
func (suite *OcppV2TestSuite) TestTransactionInfoValidation() {
	var requestTable = []GenericTestEntry{
		{transactions.Transaction{TransactionID: "42", ChargingState: transactions.ChargingStateSuspendedEV, TimeSpentCharging: newInt(100), StoppedReason: transactions.ReasonLocal, RemoteStartID: newInt(7)}, true},
		{transactions.Transaction{TransactionID: "42", ChargingState: transactions.ChargingStateSuspendedEV, TimeSpentCharging: newInt(100), StoppedReason: transactions.ReasonLocal}, true},
		{transactions.Transaction{TransactionID: "42", ChargingState: transactions.ChargingStateSuspendedEV, TimeSpentCharging: newInt(100)}, true},
		{transactions.Transaction{TransactionID: "42", ChargingState: transactions.ChargingStateSuspendedEV}, true},
		{transactions.Transaction{TransactionID: "42"}, true},
		{transactions.Transaction{}, false},
		{transactions.Transaction{TransactionID: ">36..................................", ChargingState: transactions.ChargingStateSuspendedEV, TimeSpentCharging: newInt(100), StoppedReason: transactions.ReasonLocal, RemoteStartID: newInt(7)}, false},
		{transactions.Transaction{TransactionID: "42", ChargingState: "invalidChargingState", TimeSpentCharging: newInt(100), StoppedReason: transactions.ReasonLocal, RemoteStartID: newInt(7)}, false},
		{transactions.Transaction{TransactionID: "42", ChargingState: transactions.ChargingStateSuspendedEV, TimeSpentCharging: newInt(100), StoppedReason: "invalidReason", RemoteStartID: newInt(7)}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV2TestSuite) TestTransactionEventRequestValidation() {
	transactionInfo := transactions.Transaction{TransactionID: "42"}
	idToken := types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode}
	meterValue := types.MeterValue{Timestamp: *types.NewDateTime(time.Now()), SampledValue: []types.SampledValue{{Value: 64.0}}}
	var requestTable = []GenericTestEntry{
		{transactions.TransactionEventRequest{EventType: transactions.TransactionEventStarted, Timestamp: types.NewDateTime(time.Now()), TriggerReason: transactions.TriggerReasonAuthorized, SequenceNo: 1, Offline: true, NumberOfPhasesUsed: newInt(3), CableMaxCurrent: newInt(20), ReservationID: newInt(42), TransactionInfo: transactionInfo, IDToken: &idToken, Evse: &types.EVSE{ID: 1}, MeterValue: []types.MeterValue{meterValue}}, true},
		{transactions.TransactionEventRequest{EventType: transactions.TransactionEventStarted, Timestamp: types.NewDateTime(time.Now()), TriggerReason: transactions.TriggerReasonAuthorized, SequenceNo: 1, Offline: true, NumberOfPhasesUsed: newInt(3), CableMaxCurrent: newInt(20), ReservationID: newInt(42), TransactionInfo: transactionInfo, IDToken: &idToken, Evse: &types.EVSE{ID: 1}, MeterValue: []types.MeterValue{}}, true},
		{transactions.TransactionEventRequest{EventType: transactions.TransactionEventStarted, Timestamp: types.NewDateTime(time.Now()), TriggerReason: transactions.TriggerReasonAuthorized, SequenceNo: 1, Offline: true, NumberOfPhasesUsed: newInt(3), CableMaxCurrent: newInt(20), ReservationID: newInt(42), TransactionInfo: transactionInfo, IDToken: &idToken, Evse: &types.EVSE{ID: 1}}, true},
		{transactions.TransactionEventRequest{EventType: transactions.TransactionEventStarted, Timestamp: types.NewDateTime(time.Now()), TriggerReason: transactions.TriggerReasonAuthorized, SequenceNo: 1, Offline: true, NumberOfPhasesUsed: newInt(3), CableMaxCurrent: newInt(20), ReservationID: newInt(42), TransactionInfo: transactionInfo, IDToken: &idToken}, true},
		{transactions.TransactionEventRequest{EventType: transactions.TransactionEventStarted, Timestamp: types.NewDateTime(time.Now()), TriggerReason: transactions.TriggerReasonAuthorized, SequenceNo: 1, Offline: true, NumberOfPhasesUsed: newInt(3), CableMaxCurrent: newInt(20), ReservationID: newInt(42), TransactionInfo: transactionInfo}, true},
		{transactions.TransactionEventRequest{EventType: transactions.TransactionEventStarted, Timestamp: types.NewDateTime(time.Now()), TriggerReason: transactions.TriggerReasonAuthorized, SequenceNo: 1, Offline: true, NumberOfPhasesUsed: newInt(3), CableMaxCurrent: newInt(20), TransactionInfo: transactionInfo}, true},
		{transactions.TransactionEventRequest{EventType: transactions.TransactionEventStarted, Timestamp: types.NewDateTime(time.Now()), TriggerReason: transactions.TriggerReasonAuthorized, SequenceNo: 1, Offline: true, NumberOfPhasesUsed: newInt(3), TransactionInfo: transactionInfo}, true},
		{transactions.TransactionEventRequest{EventType: transactions.TransactionEventStarted, Timestamp: types.NewDateTime(time.Now()), TriggerReason: transactions.TriggerReasonAuthorized, SequenceNo: 1, Offline: true, TransactionInfo: transactionInfo}, true},
		{transactions.TransactionEventRequest{EventType: transactions.TransactionEventStarted, Timestamp: types.NewDateTime(time.Now()), TriggerReason: transactions.TriggerReasonAuthorized, SequenceNo: 1, TransactionInfo: transactionInfo}, true},
		{transactions.TransactionEventRequest{EventType: transactions.TransactionEventStarted, Timestamp: types.NewDateTime(time.Now()), TriggerReason: transactions.TriggerReasonAuthorized, TransactionInfo: transactionInfo}, true},
		{transactions.TransactionEventRequest{EventType: transactions.TransactionEventStarted, Timestamp: types.NewDateTime(time.Now()), TriggerReason: transactions.TriggerReasonAuthorized, TransactionInfo: transactionInfo, IDToken: &types.IdToken{Type: types.IdTokenTypeNoAuthorization}}, true},
		{transactions.TransactionEventRequest{EventType: transactions.TransactionEventStarted, Timestamp: types.NewDateTime(time.Now()), TriggerReason: transactions.TriggerReasonAuthorized, TransactionInfo: transactionInfo, IDToken: &types.IdToken{Type: types.IdTokenTypeKeyCode}}, false},
		{transactions.TransactionEventRequest{EventType: transactions.TransactionEventStarted, Timestamp: types.NewDateTime(time.Now()), TriggerReason: transactions.TriggerReasonAuthorized}, false},
		{transactions.TransactionEventRequest{EventType: transactions.TransactionEventStarted, Timestamp: types.NewDateTime(time.Now()), TransactionInfo: transactionInfo}, false},
		{transactions.TransactionEventRequest{EventType: transactions.TransactionEventStarted, TriggerReason: transactions.TriggerReasonAuthorized, TransactionInfo: transactionInfo}, false},
		{transactions.TransactionEventRequest{Timestamp: types.NewDateTime(time.Now()), TriggerReason: transactions.TriggerReasonAuthorized, TransactionInfo: transactionInfo}, false},
		{transactions.TransactionEventRequest{}, false},
		{transactions.TransactionEventRequest{EventType: "invalidEventType", Timestamp: types.NewDateTime(time.Now()), TriggerReason: transactions.TriggerReasonAuthorized, SequenceNo: 1, Offline: true, NumberOfPhasesUsed: newInt(3), CableMaxCurrent: newInt(20), ReservationID: newInt(42), TransactionInfo: transactionInfo, IDToken: &idToken, Evse: &types.EVSE{ID: 1}, MeterValue: []types.MeterValue{meterValue}}, false},
		{transactions.TransactionEventRequest{EventType: transactions.TransactionEventStarted, Timestamp: types.NewDateTime(time.Now()), TriggerReason: "invalidTriggerReason", SequenceNo: 1, Offline: true, NumberOfPhasesUsed: newInt(3), CableMaxCurrent: newInt(20), ReservationID: newInt(42), TransactionInfo: transactionInfo, IDToken: &idToken, Evse: &types.EVSE{ID: 1}, MeterValue: []types.MeterValue{meterValue}}, false},
		{transactions.TransactionEventRequest{EventType: transactions.TransactionEventStarted, Timestamp: types.NewDateTime(time.Now()), TriggerReason: transactions.TriggerReasonAuthorized, SequenceNo: -1, Offline: true, NumberOfPhasesUsed: newInt(3), CableMaxCurrent: newInt(20), ReservationID: newInt(42), TransactionInfo: transactionInfo, IDToken: &idToken, Evse: &types.EVSE{ID: 1}, MeterValue: []types.MeterValue{meterValue}}, false},
		{transactions.TransactionEventRequest{EventType: transactions.TransactionEventStarted, Timestamp: types.NewDateTime(time.Now()), TriggerReason: transactions.TriggerReasonAuthorized, SequenceNo: 1, Offline: true, NumberOfPhasesUsed: newInt(-1), CableMaxCurrent: newInt(20), ReservationID: newInt(42), TransactionInfo: transactionInfo, IDToken: &idToken, Evse: &types.EVSE{ID: 1}, MeterValue: []types.MeterValue{meterValue}}, false},
		{transactions.TransactionEventRequest{EventType: transactions.TransactionEventStarted, Timestamp: types.NewDateTime(time.Now()), TriggerReason: transactions.TriggerReasonAuthorized, SequenceNo: 1, Offline: true, NumberOfPhasesUsed: newInt(3), CableMaxCurrent: newInt(20), ReservationID: newInt(42), TransactionInfo: transactions.Transaction{}, IDToken: &idToken, Evse: &types.EVSE{ID: 1}, MeterValue: []types.MeterValue{meterValue}}, false},
		{transactions.TransactionEventRequest{EventType: transactions.TransactionEventStarted, Timestamp: types.NewDateTime(time.Now()), TriggerReason: transactions.TriggerReasonAuthorized, SequenceNo: 1, Offline: true, NumberOfPhasesUsed: newInt(3), CableMaxCurrent: newInt(20), ReservationID: newInt(42), TransactionInfo: transactionInfo, IDToken: &types.IdToken{}, Evse: &types.EVSE{ID: 1}, MeterValue: []types.MeterValue{meterValue}}, false},
		{transactions.TransactionEventRequest{EventType: transactions.TransactionEventStarted, Timestamp: types.NewDateTime(time.Now()), TriggerReason: transactions.TriggerReasonAuthorized, SequenceNo: 1, Offline: true, NumberOfPhasesUsed: newInt(3), CableMaxCurrent: newInt(20), ReservationID: newInt(42), TransactionInfo: transactionInfo, IDToken: &idToken, Evse: &types.EVSE{ID: -1}, MeterValue: []types.MeterValue{meterValue}}, false},
		{transactions.TransactionEventRequest{EventType: transactions.TransactionEventStarted, Timestamp: types.NewDateTime(time.Now()), TriggerReason: transactions.TriggerReasonAuthorized, SequenceNo: 1, Offline: true, NumberOfPhasesUsed: newInt(3), CableMaxCurrent: newInt(20), ReservationID: newInt(42), TransactionInfo: transactionInfo, IDToken: &idToken, Evse: &types.EVSE{ID: 1}, MeterValue: []types.MeterValue{{}}}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV2TestSuite) TestTransactionEventResponseValidation() {
	messageContent := types.MessageContent{Format: types.MessageFormatUTF8, Content: "dummyContent"}
	var responseTable = []GenericTestEntry{
		{transactions.TransactionEventResponse{TotalCost: newFloat(8.42), ChargingPriority: newInt(2), IDTokenInfo: types.NewIdTokenInfo(types.AuthorizationStatusAccepted), UpdatedPersonalMessage: &messageContent}, true},
		{transactions.TransactionEventResponse{TotalCost: newFloat(8.42), ChargingPriority: newInt(2), IDTokenInfo: types.NewIdTokenInfo(types.AuthorizationStatusAccepted)}, true},
		{transactions.TransactionEventResponse{TotalCost: newFloat(8.42), ChargingPriority: newInt(2)}, true},
		{transactions.TransactionEventResponse{TotalCost: newFloat(8.42)}, true},
		{transactions.TransactionEventResponse{}, true},
		{transactions.TransactionEventResponse{TotalCost: newFloat(-1.0), ChargingPriority: newInt(2), IDTokenInfo: types.NewIdTokenInfo(types.AuthorizationStatusAccepted), UpdatedPersonalMessage: &messageContent}, false},
		{transactions.TransactionEventResponse{TotalCost: newFloat(8.42), ChargingPriority: newInt(-10), IDTokenInfo: types.NewIdTokenInfo(types.AuthorizationStatusAccepted), UpdatedPersonalMessage: &messageContent}, false},
		{transactions.TransactionEventResponse{TotalCost: newFloat(8.42), ChargingPriority: newInt(10), IDTokenInfo: types.NewIdTokenInfo(types.AuthorizationStatusAccepted), UpdatedPersonalMessage: &messageContent}, false},
		{transactions.TransactionEventResponse{TotalCost: newFloat(8.42), ChargingPriority: newInt(2), IDTokenInfo: types.NewIdTokenInfo("invalidAuthorizationStatus"), UpdatedPersonalMessage: &messageContent}, false},
		{transactions.TransactionEventResponse{TotalCost: newFloat(8.42), ChargingPriority: newInt(2), IDTokenInfo: types.NewIdTokenInfo(types.AuthorizationStatusAccepted), UpdatedPersonalMessage: &types.MessageContent{}}, false},
	}
	ExecuteGenericTestTable(suite, responseTable)
}

func (suite *OcppV2TestSuite) TestTransactionEventE2EMocked() {
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	timestamp := types.NewDateTime(time.Now())
	eventType := transactions.TransactionEventEnded
	triggerReason := transactions.TriggerReasonEVDeparted
	seqNo := 10
	offline := false
	phases := newInt(3)
	cableMaxCurrent := newInt(20)
	reservationID := newInt(55)
	info := transactions.Transaction{TransactionID: "42", ChargingState: transactions.ChargingStateSuspendedEV, TimeSpentCharging: newInt(1000), StoppedReason: transactions.ReasonLocal, RemoteStartID: newInt(69)}
	idToken := types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode}
	evse := types.EVSE{ID: 1}
	meterValue := types.MeterValue{Timestamp: *types.NewDateTime(time.Now()), SampledValue: []types.SampledValue{{Value: 64.0}}}
	totalCost := newFloat(8.42)
	chargingPriority := newInt(2)
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
		suite.Require().True(ok)
		suite.Require().NotNil(request)
		suite.Equal(eventType, request.EventType)
		assertDateTimeEquality(suite, timestamp, request.Timestamp)
		suite.Equal(triggerReason, request.TriggerReason)
		suite.Equal(seqNo, request.SequenceNo)
		suite.Equal(offline, request.Offline)
		suite.Equal(*phases, *request.NumberOfPhasesUsed)
		suite.Equal(*cableMaxCurrent, *request.CableMaxCurrent)
		suite.Equal(*reservationID, *request.ReservationID)
		suite.Equal(eventType, request.EventType)
		suite.Equal(info.TransactionID, request.TransactionInfo.TransactionID)
		suite.Equal(info.StoppedReason, request.TransactionInfo.StoppedReason)
		suite.Equal(info.ChargingState, request.TransactionInfo.ChargingState)
		suite.Equal(*info.TimeSpentCharging, *request.TransactionInfo.TimeSpentCharging)
		suite.Equal(*info.RemoteStartID, *request.TransactionInfo.RemoteStartID)
		suite.Require().NotNil(request.IDToken)
		suite.Equal(idToken.IdToken, request.IDToken.IdToken)
		suite.Equal(idToken.Type, request.IDToken.Type)
		suite.Require().NotNil(request.Evse)
		suite.Equal(evse.ID, request.Evse.ID)
		suite.Require().Len(request.MeterValue, 1)
		assertDateTimeEquality(suite, &meterValue.Timestamp, &request.MeterValue[0].Timestamp)
		suite.Require().Len(request.MeterValue[0].SampledValue, 1)
		suite.Equal(meterValue.SampledValue[0].Value, request.MeterValue[0].SampledValue[0].Value)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	suite.Require().NoError(err)
	response, err := suite.chargingStation.TransactionEvent(eventType, timestamp, triggerReason, seqNo, info, func(request *transactions.TransactionEventRequest) {
		request.MeterValue = []types.MeterValue{meterValue}
		request.Evse = &evse
		request.IDToken = &idToken
		request.NumberOfPhasesUsed = phases
		request.CableMaxCurrent = cableMaxCurrent
		request.ReservationID = reservationID
		request.Offline = offline
	})
	suite.Require().NoError(err)
	suite.Require().NotNil(response)
	suite.Equal(*totalCost, *response.TotalCost)
	suite.Equal(*chargingPriority, *response.ChargingPriority)
	suite.Require().NotNil(response.IDTokenInfo)
	suite.Equal(idTokenInfo.Status, response.IDTokenInfo.Status)
	suite.Require().NotNil(response.UpdatedPersonalMessage)
	suite.Equal(messageContent.Format, response.UpdatedPersonalMessage.Format)
	suite.Equal(messageContent.Content, response.UpdatedPersonalMessage.Content)
}

func (suite *OcppV2TestSuite) TestTransactionEventInvalidEndpoint() {
	messageId := defaultMessageId
	timestamp := types.NewDateTime(time.Now())
	eventType := transactions.TransactionEventEnded
	triggerReason := transactions.TriggerReasonEVDeparted
	seqNo := 10
	phases := newInt(3)
	cableMaxCurrent := newInt(20)
	reservationID := newInt(55)
	info := transactions.Transaction{TransactionID: "42", ChargingState: transactions.ChargingStateSuspendedEV, TimeSpentCharging: newInt(1000), StoppedReason: transactions.ReasonLocal, RemoteStartID: newInt(69)}
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
