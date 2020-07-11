package ocpp16_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/reservation"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"time"
)

// Test
func (suite *OcppV16TestSuite) TestReserveNowRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{reservation.ReserveNowRequest{ConnectorId: 1, ExpiryDate: types.NewDateTime(time.Now()), IdTag: "12345", ReservationId: 42, ParentIdTag: "9999"}, true},
		{reservation.ReserveNowRequest{ConnectorId: 1, ExpiryDate: types.NewDateTime(time.Now()), IdTag: "12345", ReservationId: 42}, true},
		{reservation.ReserveNowRequest{ExpiryDate: types.NewDateTime(time.Now()), IdTag: "12345", ReservationId: 42}, true},
		{reservation.ReserveNowRequest{ConnectorId: 1, ExpiryDate: types.NewDateTime(time.Now()), IdTag: "12345"}, true},
		{reservation.ReserveNowRequest{ConnectorId: 1, ExpiryDate: types.NewDateTime(time.Now())}, false},
		{reservation.ReserveNowRequest{ConnectorId: 1, IdTag: "12345"}, false},
		{reservation.ReserveNowRequest{ConnectorId: -1, ExpiryDate: types.NewDateTime(time.Now()), IdTag: "12345", ReservationId: 42}, false},
		{reservation.ReserveNowRequest{ConnectorId: 1, ExpiryDate: types.NewDateTime(time.Now()), IdTag: "12345", ReservationId: -1}, false},
		{reservation.ReserveNowRequest{ConnectorId: 1, ExpiryDate: types.NewDateTime(time.Now()), IdTag: ">20.................."}, false},
		{reservation.ReserveNowRequest{ConnectorId: 1, ExpiryDate: types.NewDateTime(time.Now()), IdTag: "12345", ReservationId: 42, ParentIdTag: ">20.................."}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV16TestSuite) TestReserveNowConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{reservation.ReserveNowConfirmation{Status: reservation.ReservationStatusAccepted}, true},
		{reservation.ReserveNowConfirmation{Status: "invalidReserveNowStatus"}, false},
		{reservation.ReserveNowConfirmation{}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV16TestSuite) TestReserveNowE2EMocked() {
	t := suite.T()
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

	reservationListener := MockChargePointReservationListener{}
	reservationListener.On("OnReserveNow", mock.Anything).Return(ReserveNowConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*reservation.ReserveNowRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, connectorId, request.ConnectorId)
		require.NotNil(t, request.ExpiryDate)
		assertDateTimeEquality(t, *expiryDate, *request.ExpiryDate)
		assert.Equal(t, idTag, request.IdTag)
		assert.Equal(t, parentIdTag, request.ParentIdTag)
		assert.Equal(t, reservationId, request.ReservationId)
	})
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	suite.chargePoint.SetReservationHandler(reservationListener)
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.centralSystem.ReserveNow(wsId, func(confirmation *reservation.ReserveNowConfirmation, err error) {
		require.Nil(t, err)
		require.NotNil(t, confirmation)
		assert.Equal(t, status, confirmation.Status)
		resultChannel <- true
	}, connectorId, expiryDate, idTag, reservationId, func(request *reservation.ReserveNowRequest) {
		request.ParentIdTag = parentIdTag
	})
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
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
