package ocpp2_test

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/reservation"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test
func (suite *OcppV2TestSuite) TestCancelReservationRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{reservation.CancelReservationRequest{ReservationID: 42}, true},
		{reservation.CancelReservationRequest{}, true},
		{reservation.CancelReservationRequest{ReservationID: -1}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestCancelReservationConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{reservation.CancelReservationResponse{Status: reservation.CancelReservationStatusAccepted, StatusInfo: types.NewStatusInfo("200", "ok")}, true},
		{reservation.CancelReservationResponse{Status: reservation.CancelReservationStatusAccepted}, true},
		{reservation.CancelReservationResponse{Status: "invalidCancelReservationStatus"}, false},
		{reservation.CancelReservationResponse{}, false},
		{reservation.CancelReservationResponse{Status: reservation.CancelReservationStatusAccepted, StatusInfo: types.NewStatusInfo("", "")}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestCancelReservationE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	reservationId := 42
	status := reservation.CancelReservationStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"reservationId":%v}]`, messageId, reservation.CancelReservationFeatureName, reservationId)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	cancelReservationConfirmation := reservation.NewCancelReservationResponse(status)
	channel := NewMockWebSocket(wsId)

	handler := &MockChargingStationReservationHandler{}
	handler.On("OnCancelReservation", mock.Anything).Return(cancelReservationConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*reservation.CancelReservationRequest)
		require.True(t, ok)
		assert.Equal(t, reservationId, request.ReservationID)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.CancelReservation(wsId, func(confirmation *reservation.CancelReservationResponse, err error) {
		require.Nil(t, err)
		assert.NotNil(t, confirmation)
		require.Equal(t, status, confirmation.Status)
		resultChannel <- true
	}, reservationId)
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV2TestSuite) TestCancelReservationInvalidEndpoint() {
	messageId := defaultMessageId
	reservationId := 42
	cancelReservationRequest := reservation.NewCancelReservationRequest(reservationId)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"reservationId":%v}]`, messageId, reservation.CancelReservationFeatureName, reservationId)
	testUnsupportedRequestFromChargingStation(suite, cancelReservationRequest, requestJson, messageId)
}
