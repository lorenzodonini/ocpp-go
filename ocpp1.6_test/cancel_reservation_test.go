package ocpp16_test

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/reservation"
	"github.com/stretchr/testify/mock"
)

// Test
func (suite *OcppV16TestSuite) TestCancelReservationRequestValidation() {
	requestTable := []GenericTestEntry{
		{reservation.CancelReservationRequest{ReservationId: 42}, true},
		{reservation.CancelReservationRequest{}, true},
		{reservation.CancelReservationRequest{ReservationId: -1}, true},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV16TestSuite) TestCancelReservationConfirmationValidation() {
	confirmationTable := []GenericTestEntry{
		{reservation.CancelReservationConfirmation{Status: reservation.CancelReservationStatusAccepted}, true},
		{reservation.CancelReservationConfirmation{Status: "invalidCancelReservationStatus"}, false},
		{reservation.CancelReservationConfirmation{}, false},
	}
	ExecuteGenericTestTable(suite, confirmationTable)
}

func (suite *OcppV16TestSuite) TestCancelReservationE2EMocked() {
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	reservationId := 42
	status := reservation.CancelReservationStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"reservationId":%v}]`, messageId, reservation.CancelReservationFeatureName, reservationId)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	cancelReservationConfirmation := reservation.NewCancelReservationConfirmation(status)
	channel := NewMockWebSocket(wsId)

	reservationListener := &MockChargePointReservationListener{}
	reservationListener.On("OnCancelReservation", mock.Anything).Return(cancelReservationConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*reservation.CancelReservationRequest)
		suite.Require().NotNil(request)
		suite.Require().True(ok)
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
	err = suite.centralSystem.CancelReservation(wsId, func(confirmation *reservation.CancelReservationConfirmation, err error) {
		suite.Require().Nil(err)
		suite.Require().NotNil(confirmation)
		suite.Equal(status, confirmation.Status)
		resultChannel <- true
	}, reservationId)
	suite.Require().Nil(err)
	result := <-resultChannel
	suite.True(result)
}

func (suite *OcppV16TestSuite) TestCancelReservationInvalidEndpoint() {
	messageId := defaultMessageId
	reservationId := 42
	cancelReservationRequest := reservation.NewCancelReservationRequest(reservationId)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"reservationId":%v}]`, messageId, reservation.CancelReservationFeatureName, reservationId)
	testUnsupportedRequestFromChargePoint(suite, cancelReservationRequest, requestJson, messageId)
}
