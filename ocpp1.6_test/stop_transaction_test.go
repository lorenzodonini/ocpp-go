package ocpp16_test

import (
	"fmt"
	"github.com/lorenzodonini/go-ocpp/ocpp1.6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"time"
)

// Test
func (suite *OcppV16TestSuite) TestStopTransactionRequestValidation() {
	t := suite.T()
	transactionData := []ocpp16.MeterValue{{Timestamp: ocpp16.NewDateTime(time.Now()), SampledValue: []ocpp16.SampledValue{{Value: "value"}}}}
	var requestTable = []RequestTestEntry{
		{ocpp16.StopTransactionRequest{IdTag: "12345", MeterStop: 100, Timestamp: ocpp16.NewDateTime(time.Now()), TransactionId: 1, Reason: ocpp16.ReasonEVDisconnected, TransactionData: transactionData}, true},
		{ocpp16.StopTransactionRequest{IdTag: "12345", MeterStop: 100, Timestamp: ocpp16.NewDateTime(time.Now()), TransactionId: 1, Reason: ocpp16.ReasonEVDisconnected, TransactionData: []ocpp16.MeterValue{}}, true},
		{ocpp16.StopTransactionRequest{IdTag: "12345", MeterStop: 100, Timestamp: ocpp16.NewDateTime(time.Now()), TransactionId: 1, Reason: ocpp16.ReasonEVDisconnected}, true},
		{ocpp16.StopTransactionRequest{IdTag: "12345", MeterStop: 100, Timestamp: ocpp16.NewDateTime(time.Now()), TransactionId: 1}, true},
		{ocpp16.StopTransactionRequest{MeterStop: 100, Timestamp: ocpp16.NewDateTime(time.Now()), TransactionId: 1}, true},
		{ocpp16.StopTransactionRequest{MeterStop: 100, Timestamp: ocpp16.NewDateTime(time.Now())}, true},
		{ocpp16.StopTransactionRequest{Timestamp: ocpp16.NewDateTime(time.Now())}, true},
		{ocpp16.StopTransactionRequest{MeterStop: 100}, false},
		{ocpp16.StopTransactionRequest{IdTag: "12345", MeterStop: 100, Timestamp: ocpp16.NewDateTime(time.Now()), TransactionId: 1, Reason: "invalidReason"}, false},
		{ocpp16.StopTransactionRequest{IdTag: ">20..................", MeterStop: 100, Timestamp: ocpp16.NewDateTime(time.Now()), TransactionId: 1}, false},
		{ocpp16.StopTransactionRequest{MeterStop: -1, Timestamp: ocpp16.NewDateTime(time.Now()), TransactionId: 1}, false},
		{ocpp16.StopTransactionRequest{MeterStop: 100, Timestamp: ocpp16.NewDateTime(time.Now()), TransactionId: 1, TransactionData: []ocpp16.MeterValue{{Timestamp: ocpp16.NewDateTime(time.Now()), SampledValue: []ocpp16.SampledValue{}}}}, false},
	}
	ExecuteRequestTestTable(t, requestTable)
}

func (suite *OcppV16TestSuite) TestStopTransactionConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []ConfirmationTestEntry{
		{ocpp16.StopTransactionConfirmation{IdTagInfo: &ocpp16.IdTagInfo{ExpiryDate: ocpp16.DateTime{Time: time.Now().Add(time.Hour * 8)}, ParentIdTag: "00000", Status: ocpp16.AuthorizationStatusAccepted}}, true},
		{ocpp16.StopTransactionConfirmation{IdTagInfo: &ocpp16.IdTagInfo{ParentIdTag: "00000", Status: ocpp16.AuthorizationStatusAccepted}}, true},
		{ocpp16.StopTransactionConfirmation{IdTagInfo: &ocpp16.IdTagInfo{ExpiryDate: ocpp16.DateTime{Time: time.Now().Add(time.Hour * 8)}, Status: ocpp16.AuthorizationStatusAccepted}}, true},
		{ocpp16.StopTransactionConfirmation{IdTagInfo: &ocpp16.IdTagInfo{Status: ocpp16.AuthorizationStatusAccepted}}, true},
		{ocpp16.StopTransactionConfirmation{IdTagInfo: &ocpp16.IdTagInfo{Status: ocpp16.AuthorizationStatusBlocked}}, true},
		{ocpp16.StopTransactionConfirmation{IdTagInfo: &ocpp16.IdTagInfo{Status: ocpp16.AuthorizationStatusExpired}}, true},
		{ocpp16.StopTransactionConfirmation{IdTagInfo: &ocpp16.IdTagInfo{Status: ocpp16.AuthorizationStatusInvalid}}, true},
		{ocpp16.StopTransactionConfirmation{IdTagInfo: &ocpp16.IdTagInfo{Status: ocpp16.AuthorizationStatusConcurrentTx}}, true},
		{ocpp16.StopTransactionConfirmation{}, true},
		{ocpp16.StopTransactionConfirmation{IdTagInfo: &ocpp16.IdTagInfo{Status: "invalidAuthorizationStatus"}}, false},
		{ocpp16.StopTransactionConfirmation{IdTagInfo: &ocpp16.IdTagInfo{ParentIdTag: ">20..................", Status: ocpp16.AuthorizationStatusAccepted}}, false},
		{ocpp16.StopTransactionConfirmation{IdTagInfo: &ocpp16.IdTagInfo{ExpiryDate: ocpp16.DateTime{Time: time.Now().Add(time.Hour * -8)}, Status: ocpp16.AuthorizationStatusAccepted}}, false},
		{ocpp16.StopTransactionConfirmation{IdTagInfo: &ocpp16.IdTagInfo{}}, false},
	}
	ExecuteConfirmationTestTable(t, confirmationTable)
}

func (suite *OcppV16TestSuite) TestStopTransactionE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	idTag := "tag1"
	mockValue := "value"
	mockUnit := ocpp16.UnitOfMeasureKW
	meterStop := 100
	transactionId := 42
	timestamp := ocpp16.NewDateTime(time.Now())
	meterValues := []ocpp16.MeterValue{{Timestamp: ocpp16.NewDateTime(time.Now()), SampledValue: []ocpp16.SampledValue{{Value: mockValue, Unit: mockUnit}}}}
	parentIdTag := "parentTag1"
	status := ocpp16.AuthorizationStatusAccepted
	expiryDate := ocpp16.DateTime{Time: time.Now().Add(time.Hour * 8)}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"idTag":"%v","meterStop":%v,"timestamp":"%v","transactionId":%v,"transactionData":[{"timestamp":"%v","sampledValue":[{"value":"%v","unit":"%v"}]}]}]`, messageId, ocpp16.StopTransactionFeatureName, idTag, meterStop, timestamp.Format(ocpp16.ISO8601), transactionId, timestamp.Format(ocpp16.ISO8601), mockValue, mockUnit)
	responseJson := fmt.Sprintf(`[3,"%v",{"idTagInfo":{"expiryDate":"%v","parentIdTag":"%v","status":"%v"}}]`, messageId, expiryDate.Format(ocpp16.ISO8601), parentIdTag, status)
	stopTransactionConfirmation := ocpp16.NewStopTransactionConfirmation()
	stopTransactionConfirmation.IdTagInfo = &ocpp16.IdTagInfo{ExpiryDate: expiryDate, ParentIdTag: parentIdTag, Status: status}
	requestRaw := []byte(requestJson)
	responseRaw := []byte(responseJson)
	channel := NewMockWebSocket(wsId)

	coreListener := MockCentralSystemCoreListener{}
	coreListener.On("OnStopTransaction", mock.AnythingOfType("string"), mock.Anything).Return(stopTransactionConfirmation, nil).Run(func(args mock.Arguments) {
		request := args.Get(1).(*ocpp16.StopTransactionRequest)
		assert.Equal(t, meterStop, request.MeterStop)
		assert.Equal(t, transactionId, request.TransactionId)
		assert.Equal(t, idTag, request.IdTag)
		assertDateTimeEquality(t, *timestamp, *request.Timestamp)
		assert.Equal(t, 1, len(request.TransactionData))
		mv := request.TransactionData[0]
		assertDateTimeEquality(t, *timestamp, *mv.Timestamp)
		assert.Equal(t, 1, len(mv.SampledValue))
		sv := mv.SampledValue[0]
		assert.Equal(t, mockValue, sv.Value)
		assert.Equal(t, mockUnit, sv.Unit)
	})
	setupDefaultCentralSystemHandlers(suite, coreListener, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: responseRaw, forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: requestRaw, forwardWrittenMessage: true})
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	assert.Nil(t, err)
	confirmation, err := suite.chargePoint.StopTransaction(meterStop, timestamp, transactionId, func(request *ocpp16.StopTransactionRequest) {
		request.IdTag = idTag
		request.TransactionData = meterValues
	})
	assert.Nil(t, err)
	assert.NotNil(t, confirmation)
	assert.Equal(t, status, confirmation.IdTagInfo.Status)
	assert.Equal(t, parentIdTag, confirmation.IdTagInfo.ParentIdTag)
	assertDateTimeEquality(t, expiryDate, confirmation.IdTagInfo.ExpiryDate)
}

func (suite *OcppV16TestSuite) TestStopTransactionInvalidEndpoint() {
	messageId := defaultMessageId
	idTag := "tag1"
	mockValue := "value"
	mockUnit := ocpp16.UnitOfMeasureKW
	meterStop := 100
	transactionId := 42
	timestamp := ocpp16.NewDateTime(time.Now())
	stopTransactionRequest := ocpp16.NewStopTransactionRequest(meterStop, timestamp, transactionId)
	stopTransactionRequest.IdTag = idTag
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"idTag":"%v","meterStop":%v,"timestamp":"%v","transactionId":%v,"transactionData":[{"timestamp":"%v","sampledValue":[{"value":"%v","unit":"%v"}]}]}]`, messageId, ocpp16.StopTransactionFeatureName, idTag, meterStop, timestamp.Format(ocpp16.ISO8601), transactionId, timestamp.Format(ocpp16.ISO8601), mockValue, mockUnit)
	testUnsupportedRequestFromCentralSystem(suite, stopTransactionRequest, requestJson, messageId)
}
