package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test
func (suite *OcppV2TestSuite) TestCancelReservationRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{ocpp2.CancelReservationRequest{ReservationId: 42}, true},
		{ocpp2.CancelReservationRequest{}, true},
		{ocpp2.CancelReservationRequest{ReservationId: -1}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestCancelReservationConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{ocpp2.CancelReservationConfirmation{Status: ocpp2.CancelReservationStatusAccepted}, true},
		{ocpp2.CancelReservationConfirmation{Status: "invalidCancelReservationStatus"}, false},
		{ocpp2.CancelReservationConfirmation{}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestCancelReservationE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	reservationId := 42
	status := ocpp2.CancelReservationStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"reservationId":%v}]`, messageId, ocpp2.CancelReservationFeatureName, reservationId)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	cancelReservationConfirmation := ocpp2.NewCancelReservationConfirmation(status)
	channel := NewMockWebSocket(wsId)

	coreListener := MockChargePointCoreListener{}
	coreListener.On("OnCancelReservation", mock.Anything).Return(cancelReservationConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*ocpp2.CancelReservationRequest)
		require.True(t, ok)
		assert.Equal(t, reservationId, request.ReservationId)
	})
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, coreListener, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.CancelReservation(wsId, func(confirmation *ocpp2.CancelReservationConfirmation, err error) {
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
	cancelReservationRequest := ocpp2.NewCancelReservationRequest(reservationId)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"reservationId":%v}]`, messageId, ocpp2.CancelReservationFeatureName, reservationId)
	testUnsupportedRequestFromChargePoint(suite, cancelReservationRequest, requestJson, messageId)
}
