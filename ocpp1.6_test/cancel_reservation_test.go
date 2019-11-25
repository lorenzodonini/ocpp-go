package ocpp16_test

import (
	"fmt"
	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test
func (suite *OcppV16TestSuite) TestCancelReservationRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{ocpp16.CancelReservationRequest{ReservationId: 42}, true},
		{ocpp16.CancelReservationRequest{}, true},
		{ocpp16.CancelReservationRequest{ReservationId: -1}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV16TestSuite) TestCancelReservationConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{ocpp16.CancelReservationConfirmation{Status: ocpp16.CancelReservationStatusAccepted}, true},
		{ocpp16.CancelReservationConfirmation{Status: "invalidCancelReservationStatus"}, false},
		{ocpp16.CancelReservationConfirmation{}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV16TestSuite) TestCancelReservationE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	reservationId := 42
	status := ocpp16.CancelReservationStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"reservationId":%v}]`, messageId, ocpp16.CancelReservationFeatureName, reservationId)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	cancelReservationConfirmation := ocpp16.NewCancelReservationConfirmation(status)
	channel := NewMockWebSocket(wsId)

	reservationListener := MockChargePointReservationListener{}
	reservationListener.On("OnCancelReservation", mock.Anything).Return(cancelReservationConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*ocpp16.CancelReservationRequest)
		assert.True(t, ok)
		assert.Equal(t, reservationId, request.ReservationId)
	})
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	suite.chargePoint.SetReservationListener(reservationListener)
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	assert.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.centralSystem.CancelReservation(wsId, func(confirmation *ocpp16.CancelReservationConfirmation, err error) {
		assert.Nil(t, err)
		assert.NotNil(t, confirmation)
		assert.Equal(t, status, confirmation.Status)
		resultChannel <- true
	}, reservationId)
	assert.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV16TestSuite) TestCancelReservationInvalidEndpoint() {
	messageId := defaultMessageId
	reservationId := 42
	cancelReservationRequest := ocpp16.NewCancelReservationRequest(reservationId)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"reservationId":%v}]`, messageId, ocpp16.CancelReservationFeatureName, reservationId)
	testUnsupportedRequestFromChargePoint(suite, cancelReservationRequest, requestJson, messageId)
}
