package ocpp16_test

import (
	"fmt"
	"github.com/lorenzodonini/go-ocpp/ocpp1.6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"time"
)

// Test
func (suite *OcppV16TestSuite) TestStartTransactionRequestValidation() {
	t := suite.T()
	var requestTable = []RequestTestEntry{
		{ocpp16.StartTransactionRequest{ConnectorId: 1, IdTag: "12345", MeterStart: 100, ReservationId: 42, Timestamp: ocpp16.DateTime{Time: time.Now()}}, true},
		{ocpp16.StartTransactionRequest{ConnectorId: 1, IdTag: "12345", MeterStart: 100, Timestamp: ocpp16.DateTime{Time: time.Now()}}, true},
		{ocpp16.StartTransactionRequest{ConnectorId: 0, IdTag: "12345", MeterStart: 100, Timestamp: ocpp16.DateTime{Time: time.Now()}}, false},
		{ocpp16.StartTransactionRequest{ConnectorId: -1, IdTag: "12345", MeterStart: 100, Timestamp: ocpp16.DateTime{Time: time.Now()}}, false},
		{ocpp16.StartTransactionRequest{ConnectorId: 1, IdTag: ">20..................", MeterStart: 100, Timestamp: ocpp16.DateTime{Time: time.Now()}}, false},
		{ocpp16.StartTransactionRequest{IdTag: "12345", MeterStart: 100, Timestamp: ocpp16.DateTime{Time: time.Now()}}, false},
		{ocpp16.StartTransactionRequest{ConnectorId: 1, MeterStart: 100, Timestamp: ocpp16.DateTime{Time: time.Now()}}, false},
		{ocpp16.StartTransactionRequest{ConnectorId: 1, IdTag: "12345", Timestamp: ocpp16.DateTime{Time: time.Now()}}, false},
		{ocpp16.StartTransactionRequest{ConnectorId: 1, IdTag: "12345", MeterStart: 100}, false},
	}
	ExecuteRequestTestTable(t, requestTable)
}

func (suite *OcppV16TestSuite) TestStartTransactionConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []ConfirmationTestEntry{
		{ocpp16.StartTransactionConfirmation{IdTagInfo: ocpp16.IdTagInfo{ExpiryDate: ocpp16.DateTime{Time: time.Now().Add(time.Hour * 8)}, ParentIdTag: "00000", Status: ocpp16.AuthorizationStatusAccepted}, TransactionId: 10}, true},
		{ocpp16.StartTransactionConfirmation{IdTagInfo: ocpp16.IdTagInfo{ParentIdTag: "00000", Status: ocpp16.AuthorizationStatusAccepted}, TransactionId: 10}, true},
		{ocpp16.StartTransactionConfirmation{IdTagInfo: ocpp16.IdTagInfo{ExpiryDate: ocpp16.DateTime{Time: time.Now().Add(time.Hour * 8)}, Status: ocpp16.AuthorizationStatusAccepted}, TransactionId: 10}, true},
		{ocpp16.StartTransactionConfirmation{IdTagInfo: ocpp16.IdTagInfo{Status: ocpp16.AuthorizationStatusAccepted}, TransactionId: 10}, true},
		{ocpp16.StartTransactionConfirmation{IdTagInfo: ocpp16.IdTagInfo{Status: ocpp16.AuthorizationStatusBlocked}, TransactionId: 10}, true},
		{ocpp16.StartTransactionConfirmation{IdTagInfo: ocpp16.IdTagInfo{Status: ocpp16.AuthorizationStatusExpired}, TransactionId: 10}, true},
		{ocpp16.StartTransactionConfirmation{IdTagInfo: ocpp16.IdTagInfo{Status: ocpp16.AuthorizationStatusInvalid}, TransactionId: 10}, true},
		{ocpp16.StartTransactionConfirmation{IdTagInfo: ocpp16.IdTagInfo{Status: ocpp16.AuthorizationStatusConcurrentTx}, TransactionId: 10}, true},
		{ocpp16.StartTransactionConfirmation{IdTagInfo: ocpp16.IdTagInfo{Status: "invalidAuthorizationStatus"}, TransactionId: 10}, false},
		{ocpp16.StartTransactionConfirmation{IdTagInfo: ocpp16.IdTagInfo{ParentIdTag: ">20..................", Status: ocpp16.AuthorizationStatusAccepted}, TransactionId: 10}, false},
		{ocpp16.StartTransactionConfirmation{IdTagInfo: ocpp16.IdTagInfo{ExpiryDate: ocpp16.DateTime{Time: time.Now().Add(time.Hour * -8)}, Status: ocpp16.AuthorizationStatusAccepted}, TransactionId: 10}, false},
		{ocpp16.StartTransactionConfirmation{IdTagInfo: ocpp16.IdTagInfo{}, TransactionId: 10}, false},
		{ocpp16.StartTransactionConfirmation{TransactionId: 10}, false},
		{ocpp16.StartTransactionConfirmation{IdTagInfo: ocpp16.IdTagInfo{ExpiryDate: ocpp16.DateTime{Time: time.Now().Add(time.Hour * 8)}, ParentIdTag: "00000", Status: ocpp16.AuthorizationStatusAccepted}}, false},
	}
	ExecuteConfirmationTestTable(t, confirmationTable)
}

func (suite *OcppV16TestSuite) TestStartTransactionE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	idTag := "tag1"
	meterStart := 100
	reservationId := 42
	connectorId := 1
	timestamp := ocpp16.DateTime{Time: time.Now()}
	parentIdTag := "parentTag1"
	status := ocpp16.AuthorizationStatusAccepted
	expiryDate := ocpp16.DateTime{Time: time.Now().Add(time.Hour * 8)}
	transactionId := 16
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"connectorId":%v,"idTag":"%v","meterStart":%v,"reservationId":%v,"timestamp":"%v"}]`, messageId, ocpp16.StartTransactionFeatureName, connectorId, idTag, meterStart, reservationId, timestamp.Format(ocpp16.ISO8601))
	responseJson := fmt.Sprintf(`[3,"%v",{"idTagInfo":{"expiryDate":"%v","parentIdTag":"%v","status":"%v"},"transactionId":%v}]`, messageId, expiryDate.Time.Format(ocpp16.ISO8601), parentIdTag, status, transactionId)
	startTransactionConfirmation := ocpp16.NewStartTransactionConfirmation(ocpp16.IdTagInfo{ExpiryDate: expiryDate, ParentIdTag: parentIdTag, Status: status}, transactionId)
	requestRaw := []byte(requestJson)
	responseRaw := []byte(responseJson)
	channel := NewMockWebSocket(wsId)

	coreListener := MockCentralSystemCoreListener{}
	coreListener.On("OnStartTransaction", mock.AnythingOfType("string"), mock.Anything).Return(startTransactionConfirmation, nil).Run(func(args mock.Arguments) {
		request := args.Get(1).(*ocpp16.StartTransactionRequest)
		assert.Equal(t, connectorId, request.ConnectorId)
		assert.Equal(t, idTag, request.IdTag)
		assert.Equal(t, meterStart, request.MeterStart)
		assert.Equal(t, reservationId, request.ReservationId)
		assertDateTimeEquality(t, timestamp, request.Timestamp)
	})
	setupDefaultCentralSystemHandlers(suite, coreListener, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: responseRaw, forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: requestRaw, forwardWrittenMessage: true})
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	assert.Nil(t, err)
	confirmation, err := suite.chargePoint.StartTransaction(connectorId, idTag, meterStart, timestamp, func(request *ocpp16.StartTransactionRequest) {
		request.ReservationId = reservationId
	})
	assert.Nil(t, err)
	assert.NotNil(t, confirmation)
	assert.Equal(t, status, confirmation.IdTagInfo.Status)
	assert.Equal(t, parentIdTag, confirmation.IdTagInfo.ParentIdTag)
	assertDateTimeEquality(t, expiryDate, confirmation.IdTagInfo.ExpiryDate)
}

func (suite *OcppV16TestSuite) TestStartTransactionInvalidEndpoint() {
	messageId := defaultMessageId
	idTag := "tag1"
	meterStart := 100
	reservationId := 42
	connectorId := 1
	timestamp := ocpp16.DateTime{Time: time.Now()}
	authorizeRequest := ocpp16.NewStartTransactionRequest(connectorId, idTag, meterStart, timestamp)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"connectorId":%v,"idTag":"%v","meterStart":%v,"reservationId":%v,"timestamp":"%v"}]`, messageId, ocpp16.StartTransactionFeatureName, connectorId, idTag, meterStart, reservationId, timestamp.Format(ocpp16.ISO8601))
	testUnsupportedRequestFromCentralSystem(suite, authorizeRequest, requestJson, messageId)
}
