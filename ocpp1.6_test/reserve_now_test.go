package ocpp16_test

import (
	"fmt"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/reservation"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/stretchr/testify/mock"
)

// Test
func (suite *OcppV16TestSuite) TestReserveNowRequestValidation() {
	requestTable := []GenericTestEntry{
		{reservation.ReserveNowRequest{ConnectorId: 1, ExpiryDate: types.NewDateTime(time.Now()), IdTag: "12345", ReservationId: 42, ParentIdTag: "9999"}, true},
		{reservation.ReserveNowRequest{ConnectorId: 1, ExpiryDate: types.NewDateTime(time.Now()), IdTag: "12345", ReservationId: 42}, true},
		{reservation.ReserveNowRequest{ExpiryDate: types.NewDateTime(time.Now()), IdTag: "12345", ReservationId: 42}, true},
		{reservation.ReserveNowRequest{ConnectorId: 1, ExpiryDate: types.NewDateTime(time.Now()), IdTag: "12345"}, true},
		{reservation.ReserveNowRequest{ConnectorId: 1, ExpiryDate: types.NewDateTime(time.Now())}, false},
		{reservation.ReserveNowRequest{ConnectorId: 1, IdTag: "12345"}, false},
		{reservation.ReserveNowRequest{ConnectorId: -1, ExpiryDate: types.NewDateTime(time.Now()), IdTag: "12345", ReservationId: 42}, false},
		{reservation.ReserveNowRequest{ConnectorId: 1, ExpiryDate: types.NewDateTime(time.Now()), IdTag: "12345", ReservationId: -1}, true},
		{reservation.ReserveNowRequest{ConnectorId: 1, ExpiryDate: types.NewDateTime(time.Now()), IdTag: ">20.................."}, false},
		{reservation.ReserveNowRequest{ConnectorId: 1, ExpiryDate: types.NewDateTime(time.Now()), IdTag: "12345", ReservationId: 42, ParentIdTag: ">20.................."}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV16TestSuite) TestReserveNowConfirmationValidation() {
	confirmationTable := []GenericTestEntry{
		{reservation.ReserveNowConfirmation{Status: reservation.ReservationStatusAccepted}, true},
		{reservation.ReserveNowConfirmation{Status: "invalidReserveNowStatus"}, false},
		{reservation.ReserveNowConfirmation{}, false},
	}
	ExecuteGenericTestTable(suite, confirmationTable)
}

func (suite *OcppV16TestSuite) TestReserveNowE2EMocked() {
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	idTag := "12345"
	parentIdTag := "00000"
	connectorId := 1
	reservationId := 42
	expiryDate := types.NewDateTime(time.Now())
	status := reservation.ReservationStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"connectorId":%v,"expiryDate":"%v","idTag":"%v","parentIdTag":"%v","reservationId":%v}]`,
		messageId, reservation.ReserveNowFeatureName, connectorId, expiryDate.FormatTimestamp(), idTag, parentIdTag, reservationId)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	ReserveNowConfirmation := reservation.NewReserveNowConfirmation(status)
	channel := NewMockWebSocket(wsId)

	reservationListener := &MockChargePointReservationListener{}
	reservationListener.On("OnReserveNow", mock.Anything).Return(ReserveNowConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*reservation.ReserveNowRequest)
		suite.Require().True(ok)
		suite.Require().NotNil(request)
		suite.Equal(connectorId, request.ConnectorId)
		suite.Require().NotNil(request.ExpiryDate)
		assertDateTimeEquality(suite, *expiryDate, *request.ExpiryDate)
		suite.Equal(idTag, request.IdTag)
		suite.Equal(parentIdTag, request.ParentIdTag)
		suite.Equal(reservationId, request.ReservationId)
	})
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	suite.chargePoint.SetReservationHandler(reservationListener)
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	suite.Require().Nil(err)
	resultChannel := make(chan bool, 1)
	err = suite.centralSystem.ReserveNow(wsId, func(confirmation *reservation.ReserveNowConfirmation, err error) {
		suite.Require().Nil(err)
		suite.Require().NotNil(confirmation)
		suite.Equal(status, confirmation.Status)
		resultChannel <- true
	}, connectorId, expiryDate, idTag, reservationId, func(request *reservation.ReserveNowRequest) {
		request.ParentIdTag = parentIdTag
	})
	suite.Require().Nil(err)
	result := <-resultChannel
	suite.True(result)
}

func (suite *OcppV16TestSuite) TestReserveNowInvalidEndpoint() {
	messageId := defaultMessageId
	idTag := "12345"
	parentIdTag := "00000"
	connectorId := 1
	reservationId := 42
	expiryDate := types.NewDateTime(time.Now())
	reserveNowRequest := reservation.NewReserveNowRequest(connectorId, expiryDate, idTag, reservationId)
	reserveNowRequest.ParentIdTag = parentIdTag
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"connectorId":%v,"expiryDate":"%v","idTag":"%v","parentIdTag":"%v","reservationId":%v}]`,
		messageId, reservation.ReserveNowFeatureName, connectorId, expiryDate.FormatTimestamp(), idTag, parentIdTag, reservationId)
	testUnsupportedRequestFromChargePoint(suite, reserveNowRequest, requestJson, messageId)
}
