package ocpp16_test

import (
	"fmt"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/stretchr/testify/mock"
)

// Test
func (suite *OcppV16TestSuite) TestStartTransactionRequestValidation() {
	var requestTable = []GenericTestEntry{
		{core.StartTransactionRequest{ConnectorId: 1, IdTag: "12345", MeterStart: 100, ReservationId: newInt(42), Timestamp: types.NewDateTime(time.Now())}, true},
		{core.StartTransactionRequest{ConnectorId: 1, IdTag: "12345", MeterStart: 100, Timestamp: types.NewDateTime(time.Now())}, true},
		{core.StartTransactionRequest{ConnectorId: 1, IdTag: "12345", Timestamp: types.NewDateTime(time.Now())}, true},
		{core.StartTransactionRequest{ConnectorId: 0, IdTag: "12345", MeterStart: 100, Timestamp: types.NewDateTime(time.Now())}, false},
		{core.StartTransactionRequest{ConnectorId: -1, IdTag: "12345", MeterStart: 100, Timestamp: types.NewDateTime(time.Now())}, false},
		{core.StartTransactionRequest{ConnectorId: 1, IdTag: ">20..................", MeterStart: 100, Timestamp: types.NewDateTime(time.Now())}, false},
		{core.StartTransactionRequest{IdTag: "12345", MeterStart: 100, Timestamp: types.NewDateTime(time.Now())}, false},
		{core.StartTransactionRequest{ConnectorId: 1, MeterStart: 100, Timestamp: types.NewDateTime(time.Now())}, false},
		{core.StartTransactionRequest{ConnectorId: 1, IdTag: "12345", MeterStart: 100}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV16TestSuite) TestStartTransactionConfirmationValidation() {
	var confirmationTable = []GenericTestEntry{
		{core.StartTransactionConfirmation{IdTagInfo: &types.IdTagInfo{ExpiryDate: types.NewDateTime(time.Now().Add(time.Hour * 8)), ParentIdTag: "00000", Status: types.AuthorizationStatusAccepted}, TransactionId: 10}, true},
		{core.StartTransactionConfirmation{IdTagInfo: &types.IdTagInfo{ExpiryDate: types.NewDateTime(time.Now().Add(time.Hour * 8)), ParentIdTag: "00000", Status: types.AuthorizationStatusAccepted}}, true},
		{core.StartTransactionConfirmation{IdTagInfo: &types.IdTagInfo{Status: "invalidAuthorizationStatus"}, TransactionId: 10}, false},
		{core.StartTransactionConfirmation{TransactionId: 10}, false},
	}
	ExecuteGenericTestTable(suite, confirmationTable)
}

func (suite *OcppV16TestSuite) TestStartTransactionE2EMocked() {
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	idTag := "tag1"
	meterStart := 100
	reservationId := newInt(42)
	connectorId := 1
	timestamp := types.NewDateTime(time.Now())
	parentIdTag := "parentTag1"
	status := types.AuthorizationStatusAccepted
	expiryDate := types.NewDateTime(time.Now().Add(time.Hour * 8))
	transactionId := 16
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"connectorId":%v,"idTag":"%v","meterStart":%v,"reservationId":%v,"timestamp":"%v"}]`,
		messageId, core.StartTransactionFeatureName, connectorId, idTag, meterStart, *reservationId, timestamp.FormatTimestamp())
	responseJson := fmt.Sprintf(`[3,"%v",{"idTagInfo":{"expiryDate":"%v","parentIdTag":"%v","status":"%v"},"transactionId":%v}]`, messageId, expiryDate.FormatTimestamp(), parentIdTag, status, transactionId)
	startTransactionConfirmation := core.NewStartTransactionConfirmation(&types.IdTagInfo{ExpiryDate: expiryDate, ParentIdTag: parentIdTag, Status: status}, transactionId)
	requestRaw := []byte(requestJson)
	responseRaw := []byte(responseJson)
	channel := NewMockWebSocket(wsId)

	coreListener := &MockCentralSystemCoreListener{}
	coreListener.On("OnStartTransaction", mock.AnythingOfType("string"), mock.Anything).Return(startTransactionConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*core.StartTransactionRequest)
		suite.Require().True(ok)
		suite.Require().NotNil(request)
		suite.Equal(connectorId, request.ConnectorId)
		suite.Equal(idTag, request.IdTag)
		suite.Equal(meterStart, request.MeterStart)
		suite.Equal(*reservationId, *request.ReservationId)
		assertDateTimeEquality(suite, *timestamp, *request.Timestamp)
	})
	setupDefaultCentralSystemHandlers(suite, coreListener, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: responseRaw, forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: requestRaw, forwardWrittenMessage: true})
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	suite.Require().Nil(err)
	confirmation, err := suite.chargePoint.StartTransaction(connectorId, idTag, meterStart, timestamp, func(request *core.StartTransactionRequest) {
		request.ReservationId = reservationId
	})
	suite.Require().Nil(err)
	suite.Require().NotNil(confirmation)
	suite.Equal(status, confirmation.IdTagInfo.Status)
	suite.Equal(parentIdTag, confirmation.IdTagInfo.ParentIdTag)
	assertDateTimeEquality(suite, *expiryDate, *confirmation.IdTagInfo.ExpiryDate)
}

func (suite *OcppV16TestSuite) TestStartTransactionInvalidEndpoint() {
	messageId := defaultMessageId
	idTag := "tag1"
	meterStart := 100
	reservationId := 42
	connectorId := 1
	timestamp := types.NewDateTime(time.Now())
	authorizeRequest := core.NewStartTransactionRequest(connectorId, idTag, meterStart, timestamp)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"connectorId":%v,"idTag":"%v","meterStart":%v,"reservationId":%v,"timestamp":"%v"}]`, messageId, core.StartTransactionFeatureName, connectorId, idTag, meterStart, reservationId, timestamp.FormatTimestamp())
	testUnsupportedRequestFromCentralSystem(suite, authorizeRequest, requestJson, messageId)
}
