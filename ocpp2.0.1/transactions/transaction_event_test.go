package transactions

import (
	"reflect"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *transactionsTestSuite) TestTransactionInfoValidation() {
	var requestTable = []tests.GenericTestEntry{
		{Transaction{TransactionID: "42", ChargingState: ChargingStateSuspendedEV, TimeSpentCharging: tests.NewInt(100), StoppedReason: ReasonLocal, RemoteStartID: tests.NewInt(7)}, true},
		{Transaction{TransactionID: "42", ChargingState: ChargingStateSuspendedEV, TimeSpentCharging: tests.NewInt(100), StoppedReason: ReasonLocal}, true},
		{Transaction{TransactionID: "42", ChargingState: ChargingStateSuspendedEV, TimeSpentCharging: tests.NewInt(100)}, true},
		{Transaction{TransactionID: "42", ChargingState: ChargingStateSuspendedEV}, true},
		{Transaction{TransactionID: "42"}, true},
		{Transaction{}, false},
		{Transaction{TransactionID: ">36..................................", ChargingState: ChargingStateSuspendedEV, TimeSpentCharging: tests.NewInt(100), StoppedReason: ReasonLocal, RemoteStartID: tests.NewInt(7)}, false},
		{Transaction{TransactionID: "42", ChargingState: "invalidChargingState", TimeSpentCharging: tests.NewInt(100), StoppedReason: ReasonLocal, RemoteStartID: tests.NewInt(7)}, false},
		{Transaction{TransactionID: "42", ChargingState: ChargingStateSuspendedEV, TimeSpentCharging: tests.NewInt(100), StoppedReason: "invalidReason", RemoteStartID: tests.NewInt(7)}, false},
	}
	tests.ExecuteGenericTestTable(suite.T(), requestTable)
}

func (suite *transactionsTestSuite) TestTransactionEventRequestValidation() {
	t := suite.T()
	transactionInfo := Transaction{TransactionID: "42"}
	idToken := types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode}
	meterValue := types.MeterValue{Timestamp: *types.NewDateTime(time.Now()), SampledValue: []types.SampledValue{{Value: 64.0}}}
	var requestTable = []tests.GenericTestEntry{
		{TransactionEventRequest{EventType: TransactionEventStarted, Timestamp: types.NewDateTime(time.Now()), TriggerReason: TriggerReasonAuthorized, SequenceNo: 1, Offline: true, NumberOfPhasesUsed: tests.NewInt(3), CableMaxCurrent: tests.NewInt(20), ReservationID: tests.NewInt(42), TransactionInfo: transactionInfo, IDToken: &idToken, Evse: &types.EVSE{ID: 1}, MeterValue: []types.MeterValue{meterValue}}, true},
		{TransactionEventRequest{EventType: TransactionEventStarted, Timestamp: types.NewDateTime(time.Now()), TriggerReason: TriggerReasonAuthorized, SequenceNo: 1, Offline: true, NumberOfPhasesUsed: tests.NewInt(3), CableMaxCurrent: tests.NewInt(20), ReservationID: tests.NewInt(42), TransactionInfo: transactionInfo, IDToken: &idToken, Evse: &types.EVSE{ID: 1}, MeterValue: []types.MeterValue{}}, true},
		{TransactionEventRequest{EventType: TransactionEventStarted, Timestamp: types.NewDateTime(time.Now()), TriggerReason: TriggerReasonAuthorized, SequenceNo: 1, Offline: true, NumberOfPhasesUsed: tests.NewInt(3), CableMaxCurrent: tests.NewInt(20), ReservationID: tests.NewInt(42), TransactionInfo: transactionInfo, IDToken: &idToken, Evse: &types.EVSE{ID: 1}}, true},
		{TransactionEventRequest{EventType: TransactionEventStarted, Timestamp: types.NewDateTime(time.Now()), TriggerReason: TriggerReasonAuthorized, SequenceNo: 1, Offline: true, NumberOfPhasesUsed: tests.NewInt(3), CableMaxCurrent: tests.NewInt(20), ReservationID: tests.NewInt(42), TransactionInfo: transactionInfo, IDToken: &idToken}, true},
		{TransactionEventRequest{EventType: TransactionEventStarted, Timestamp: types.NewDateTime(time.Now()), TriggerReason: TriggerReasonAuthorized, SequenceNo: 1, Offline: true, NumberOfPhasesUsed: tests.NewInt(3), CableMaxCurrent: tests.NewInt(20), ReservationID: tests.NewInt(42), TransactionInfo: transactionInfo}, true},
		{TransactionEventRequest{EventType: TransactionEventStarted, Timestamp: types.NewDateTime(time.Now()), TriggerReason: TriggerReasonAuthorized, SequenceNo: 1, Offline: true, NumberOfPhasesUsed: tests.NewInt(3), CableMaxCurrent: tests.NewInt(20), TransactionInfo: transactionInfo}, true},
		{TransactionEventRequest{EventType: TransactionEventStarted, Timestamp: types.NewDateTime(time.Now()), TriggerReason: TriggerReasonAuthorized, SequenceNo: 1, Offline: true, NumberOfPhasesUsed: tests.NewInt(3), TransactionInfo: transactionInfo}, true},
		{TransactionEventRequest{EventType: TransactionEventStarted, Timestamp: types.NewDateTime(time.Now()), TriggerReason: TriggerReasonAuthorized, SequenceNo: 1, Offline: true, TransactionInfo: transactionInfo}, true},
		{TransactionEventRequest{EventType: TransactionEventStarted, Timestamp: types.NewDateTime(time.Now()), TriggerReason: TriggerReasonAuthorized, SequenceNo: 1, TransactionInfo: transactionInfo}, true},
		{TransactionEventRequest{EventType: TransactionEventStarted, Timestamp: types.NewDateTime(time.Now()), TriggerReason: TriggerReasonAuthorized, TransactionInfo: transactionInfo}, true},
		{TransactionEventRequest{EventType: TransactionEventStarted, Timestamp: types.NewDateTime(time.Now()), TriggerReason: TriggerReasonAuthorized, TransactionInfo: transactionInfo, IDToken: &types.IdToken{Type: types.IdTokenTypeNoAuthorization}}, true},
		{TransactionEventRequest{EventType: TransactionEventStarted, Timestamp: types.NewDateTime(time.Now()), TriggerReason: TriggerReasonAuthorized, TransactionInfo: transactionInfo, IDToken: &types.IdToken{Type: types.IdTokenTypeKeyCode}}, false},
		{TransactionEventRequest{EventType: TransactionEventStarted, Timestamp: types.NewDateTime(time.Now()), TriggerReason: TriggerReasonAuthorized}, false},
		{TransactionEventRequest{EventType: TransactionEventStarted, Timestamp: types.NewDateTime(time.Now()), TransactionInfo: transactionInfo}, false},
		{TransactionEventRequest{EventType: TransactionEventStarted, TriggerReason: TriggerReasonAuthorized, TransactionInfo: transactionInfo}, false},
		{TransactionEventRequest{Timestamp: types.NewDateTime(time.Now()), TriggerReason: TriggerReasonAuthorized, TransactionInfo: transactionInfo}, false},
		{TransactionEventRequest{}, false},
		{TransactionEventRequest{EventType: "invalidEventType", Timestamp: types.NewDateTime(time.Now()), TriggerReason: TriggerReasonAuthorized, SequenceNo: 1, Offline: true, NumberOfPhasesUsed: tests.NewInt(3), CableMaxCurrent: tests.NewInt(20), ReservationID: tests.NewInt(42), TransactionInfo: transactionInfo, IDToken: &idToken, Evse: &types.EVSE{ID: 1}, MeterValue: []types.MeterValue{meterValue}}, false},
		{TransactionEventRequest{EventType: TransactionEventStarted, Timestamp: types.NewDateTime(time.Now()), TriggerReason: "invalidTriggerReason", SequenceNo: 1, Offline: true, NumberOfPhasesUsed: tests.NewInt(3), CableMaxCurrent: tests.NewInt(20), ReservationID: tests.NewInt(42), TransactionInfo: transactionInfo, IDToken: &idToken, Evse: &types.EVSE{ID: 1}, MeterValue: []types.MeterValue{meterValue}}, false},
		{TransactionEventRequest{EventType: TransactionEventStarted, Timestamp: types.NewDateTime(time.Now()), TriggerReason: TriggerReasonAuthorized, SequenceNo: -1, Offline: true, NumberOfPhasesUsed: tests.NewInt(3), CableMaxCurrent: tests.NewInt(20), ReservationID: tests.NewInt(42), TransactionInfo: transactionInfo, IDToken: &idToken, Evse: &types.EVSE{ID: 1}, MeterValue: []types.MeterValue{meterValue}}, false},
		{TransactionEventRequest{EventType: TransactionEventStarted, Timestamp: types.NewDateTime(time.Now()), TriggerReason: TriggerReasonAuthorized, SequenceNo: 1, Offline: true, NumberOfPhasesUsed: tests.NewInt(-1), CableMaxCurrent: tests.NewInt(20), ReservationID: tests.NewInt(42), TransactionInfo: transactionInfo, IDToken: &idToken, Evse: &types.EVSE{ID: 1}, MeterValue: []types.MeterValue{meterValue}}, false},
		{TransactionEventRequest{EventType: TransactionEventStarted, Timestamp: types.NewDateTime(time.Now()), TriggerReason: TriggerReasonAuthorized, SequenceNo: 1, Offline: true, NumberOfPhasesUsed: tests.NewInt(3), CableMaxCurrent: tests.NewInt(20), ReservationID: tests.NewInt(42), TransactionInfo: Transaction{}, IDToken: &idToken, Evse: &types.EVSE{ID: 1}, MeterValue: []types.MeterValue{meterValue}}, false},
		{TransactionEventRequest{EventType: TransactionEventStarted, Timestamp: types.NewDateTime(time.Now()), TriggerReason: TriggerReasonAuthorized, SequenceNo: 1, Offline: true, NumberOfPhasesUsed: tests.NewInt(3), CableMaxCurrent: tests.NewInt(20), ReservationID: tests.NewInt(42), TransactionInfo: transactionInfo, IDToken: &types.IdToken{}, Evse: &types.EVSE{ID: 1}, MeterValue: []types.MeterValue{meterValue}}, false},
		{TransactionEventRequest{EventType: TransactionEventStarted, Timestamp: types.NewDateTime(time.Now()), TriggerReason: TriggerReasonAuthorized, SequenceNo: 1, Offline: true, NumberOfPhasesUsed: tests.NewInt(3), CableMaxCurrent: tests.NewInt(20), ReservationID: tests.NewInt(42), TransactionInfo: transactionInfo, IDToken: &idToken, Evse: &types.EVSE{ID: -1}, MeterValue: []types.MeterValue{meterValue}}, false},
		{TransactionEventRequest{EventType: TransactionEventStarted, Timestamp: types.NewDateTime(time.Now()), TriggerReason: TriggerReasonAuthorized, SequenceNo: 1, Offline: true, NumberOfPhasesUsed: tests.NewInt(3), CableMaxCurrent: tests.NewInt(20), ReservationID: tests.NewInt(42), TransactionInfo: transactionInfo, IDToken: &idToken, Evse: &types.EVSE{ID: 1}, MeterValue: []types.MeterValue{{}}}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *transactionsTestSuite) TestTransactionEventFeature() {
	feature := TransactionEventFeature{}
	suite.Equal(TransactionEventFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(TransactionEventRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(TransactionEventResponse{}), feature.GetResponseType())
}

func (suite *transactionsTestSuite) TestNewTransactionEventRequest() {
	ts := types.NewDateTime(time.Now())
	txInfo := Transaction{TransactionID: "tx-42"}
	req := NewTransactionEventRequest(TransactionEventStarted, ts, TriggerReasonAuthorized, 1, txInfo)
	suite.NotNil(req)
	suite.Equal(TransactionEventFeatureName, req.GetFeatureName())
	suite.Equal(TransactionEventStarted, req.EventType)
	suite.Equal(ts, req.Timestamp)
	suite.Equal(TriggerReasonAuthorized, req.TriggerReason)
	suite.Equal(1, req.SequenceNo)
	suite.Equal(txInfo, req.TransactionInfo)
}

func (suite *transactionsTestSuite) TestNewTransactionEventResponse() {
	resp := NewTransactionEventResponse()
	suite.NotNil(resp)
	suite.Equal(TransactionEventFeatureName, resp.GetFeatureName())
}

func (suite *transactionsTestSuite) TestTransactionEventResponseValidation() {
	t := suite.T()
	messageContent := types.MessageContent{Format: types.MessageFormatUTF8, Content: "dummyContent"}
	var responseTable = []tests.GenericTestEntry{
		{TransactionEventResponse{TotalCost: tests.NewFloat(8.42), ChargingPriority: tests.NewInt(2), IDTokenInfo: types.NewIdTokenInfo(types.AuthorizationStatusAccepted), UpdatedPersonalMessage: &messageContent}, true},
		{TransactionEventResponse{TotalCost: tests.NewFloat(8.42), ChargingPriority: tests.NewInt(2), IDTokenInfo: types.NewIdTokenInfo(types.AuthorizationStatusAccepted)}, true},
		{TransactionEventResponse{TotalCost: tests.NewFloat(8.42), ChargingPriority: tests.NewInt(2)}, true},
		{TransactionEventResponse{TotalCost: tests.NewFloat(8.42)}, true},
		{TransactionEventResponse{}, true},
		{TransactionEventResponse{TotalCost: tests.NewFloat(-1.0), ChargingPriority: tests.NewInt(2), IDTokenInfo: types.NewIdTokenInfo(types.AuthorizationStatusAccepted), UpdatedPersonalMessage: &messageContent}, false},
		{TransactionEventResponse{TotalCost: tests.NewFloat(8.42), ChargingPriority: tests.NewInt(-10), IDTokenInfo: types.NewIdTokenInfo(types.AuthorizationStatusAccepted), UpdatedPersonalMessage: &messageContent}, false},
		{TransactionEventResponse{TotalCost: tests.NewFloat(8.42), ChargingPriority: tests.NewInt(10), IDTokenInfo: types.NewIdTokenInfo(types.AuthorizationStatusAccepted), UpdatedPersonalMessage: &messageContent}, false},
		{TransactionEventResponse{TotalCost: tests.NewFloat(8.42), ChargingPriority: tests.NewInt(2), IDTokenInfo: types.NewIdTokenInfo("invalidAuthorizationStatus"), UpdatedPersonalMessage: &messageContent}, false},
		{TransactionEventResponse{TotalCost: tests.NewFloat(8.42), ChargingPriority: tests.NewInt(2), IDTokenInfo: types.NewIdTokenInfo(types.AuthorizationStatusAccepted), UpdatedPersonalMessage: &types.MessageContent{}}, false},
	}
	tests.ExecuteGenericTestTable(t, responseTable)
}
