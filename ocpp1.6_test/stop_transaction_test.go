package ocpp16_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"time"
)

// Test
func (suite *OcppV16TestSuite) TestStopTransactionRequestValidation() {
	t := suite.T()
	transactionData := []types.MeterValue{{Timestamp: types.NewDateTime(time.Now()), SampledValue: []types.SampledValue{{Value: "value"}}}}
	var requestTable = []GenericTestEntry{
		{core.StopTransactionRequest{IdTag: "12345", MeterStop: 100, Timestamp: types.NewDateTime(time.Now()), TransactionId: 1, Reason: core.ReasonEVDisconnected, TransactionData: transactionData}, true},
		{core.StopTransactionRequest{IdTag: "12345", MeterStop: 100, Timestamp: types.NewDateTime(time.Now()), TransactionId: 1, Reason: core.ReasonEVDisconnected, TransactionData: []types.MeterValue{}}, true},
		{core.StopTransactionRequest{IdTag: "12345", MeterStop: 100, Timestamp: types.NewDateTime(time.Now()), TransactionId: 1, Reason: core.ReasonEVDisconnected}, true},
		{core.StopTransactionRequest{IdTag: "12345", MeterStop: 100, Timestamp: types.NewDateTime(time.Now()), TransactionId: 1}, true},
		{core.StopTransactionRequest{MeterStop: 100, Timestamp: types.NewDateTime(time.Now()), TransactionId: 1}, true},
		{core.StopTransactionRequest{MeterStop: 100, Timestamp: types.NewDateTime(time.Now())}, true},
		{core.StopTransactionRequest{Timestamp: types.NewDateTime(time.Now())}, true},
		{core.StopTransactionRequest{MeterStop: 100}, false},
		{core.StopTransactionRequest{IdTag: "12345", MeterStop: 100, Timestamp: types.NewDateTime(time.Now()), TransactionId: 1, Reason: "invalidReason"}, false},
		{core.StopTransactionRequest{IdTag: ">20..................", MeterStop: 100, Timestamp: types.NewDateTime(time.Now()), TransactionId: 1}, false},
		{core.StopTransactionRequest{MeterStop: -1, Timestamp: types.NewDateTime(time.Now()), TransactionId: 1}, false},
		{core.StopTransactionRequest{MeterStop: 100, Timestamp: types.NewDateTime(time.Now()), TransactionId: 1, TransactionData: []types.MeterValue{{Timestamp: types.NewDateTime(time.Now()), SampledValue: []types.SampledValue{}}}}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV16TestSuite) TestStopTransactionConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{core.StopTransactionConfirmation{IdTagInfo: &types.IdTagInfo{ExpiryDate: types.NewDateTime(time.Now().Add(time.Hour * 8)), ParentIdTag: "00000", Status: types.AuthorizationStatusAccepted}}, true},
		{core.StopTransactionConfirmation{}, true},
		{core.StopTransactionConfirmation{IdTagInfo: &types.IdTagInfo{Status: "invalidAuthorizationStatus"}}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV16TestSuite) TestStopTransactionE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	idTag := "tag1"
	mockValue := "value"
	mockUnit := types.UnitOfMeasureKW
	meterStop := 100
	transactionId := 42
	timestamp := types.NewDateTime(time.Now())
	meterValues := []types.MeterValue{{Timestamp: types.NewDateTime(time.Now()), SampledValue: []types.SampledValue{{Value: mockValue, Unit: mockUnit}}}}
	parentIdTag := "parentTag1"
	status := types.AuthorizationStatusAccepted
	expiryDate := types.NewDateTime(time.Now().Add(time.Hour * 8))
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"idTag":"%v","meterStop":%v,"timestamp":"%v","transactionId":%v,"transactionData":[{"timestamp":"%v","sampledValue":[{"value":"%v","unit":"%v"}]}]}]`, messageId, core.StopTransactionFeatureName, idTag, meterStop, timestamp.FormatTimestamp(), transactionId, timestamp.FormatTimestamp(), mockValue, mockUnit)
	responseJson := fmt.Sprintf(`[3,"%v",{"idTagInfo":{"expiryDate":"%v","parentIdTag":"%v","status":"%v"}}]`, messageId, expiryDate.FormatTimestamp(), parentIdTag, status)
	stopTransactionConfirmation := core.NewStopTransactionConfirmation()
	stopTransactionConfirmation.IdTagInfo = &types.IdTagInfo{ExpiryDate: expiryDate, ParentIdTag: parentIdTag, Status: status}
	requestRaw := []byte(requestJson)
	responseRaw := []byte(responseJson)
	channel := NewMockWebSocket(wsId)

	coreListener := MockCentralSystemCoreListener{}
	coreListener.On("OnStopTransaction", mock.AnythingOfType("string"), mock.Anything).Return(stopTransactionConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*core.StopTransactionRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, meterStop, request.MeterStop)
		assert.Equal(t, transactionId, request.TransactionId)
		assert.Equal(t, idTag, request.IdTag)
		assertDateTimeEquality(t, *timestamp, *request.Timestamp)
		require.Len(t, request.TransactionData, 1)
		assertDateTimeEquality(t, *timestamp, *request.TransactionData[0].Timestamp)
		require.Len(t, request.TransactionData[0].SampledValue, 1)
		sv := request.TransactionData[0].SampledValue[0]
		assert.Equal(t, mockValue, sv.Value)
		assert.Equal(t, mockUnit, sv.Unit)
	})
	setupDefaultCentralSystemHandlers(suite, coreListener, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: responseRaw, forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: requestRaw, forwardWrittenMessage: true})
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	require.Nil(t, err)
	confirmation, err := suite.chargePoint.StopTransaction(meterStop, timestamp, transactionId, func(request *core.StopTransactionRequest) {
		request.IdTag = idTag
		request.TransactionData = meterValues
	})
	require.Nil(t, err)
	require.NotNil(t, confirmation)
	assert.Equal(t, status, confirmation.IdTagInfo.Status)
	assert.Equal(t, parentIdTag, confirmation.IdTagInfo.ParentIdTag)
	assertDateTimeEquality(t, *expiryDate, *confirmation.IdTagInfo.ExpiryDate)
}

func (suite *OcppV16TestSuite) TestStopTransactionInvalidEndpoint() {
	messageId := defaultMessageId
	idTag := "tag1"
	mockValue := "value"
	mockUnit := types.UnitOfMeasureKW
	meterStop := 100
	transactionId := 42
	timestamp := types.NewDateTime(time.Now())
	stopTransactionRequest := core.NewStopTransactionRequest(meterStop, timestamp, transactionId)
	stopTransactionRequest.IdTag = idTag
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"idTag":"%v","meterStop":%v,"timestamp":"%v","transactionId":%v,"transactionData":[{"timestamp":"%v","sampledValue":[{"value":"%v","unit":"%v"}]}]}]`, messageId, core.StopTransactionFeatureName, idTag, meterStop, timestamp.FormatTimestamp(), transactionId, timestamp.FormatTimestamp(), mockValue, mockUnit)
	testUnsupportedRequestFromCentralSystem(suite, stopTransactionRequest, requestJson, messageId)
}
