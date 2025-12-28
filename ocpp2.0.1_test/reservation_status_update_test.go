package ocpp2_test

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/reservation"

	"github.com/stretchr/testify/mock"
)

// Test
func (suite *OcppV2TestSuite) TestReservationStatusUpdateRequestValidation() {
	var requestTable = []GenericTestEntry{
		{reservation.ReservationStatusUpdateRequest{ReservationID: 42, Status: reservation.ReservationUpdateStatusExpired}, true},
		{reservation.ReservationStatusUpdateRequest{ReservationID: 42, Status: reservation.ReservationUpdateStatusRemoved}, true},
		{reservation.ReservationStatusUpdateRequest{Status: reservation.ReservationUpdateStatusExpired}, true},
		{reservation.ReservationStatusUpdateRequest{}, false},
		{reservation.ReservationStatusUpdateRequest{ReservationID: 42}, false},
		{reservation.ReservationStatusUpdateRequest{ReservationID: -1, Status: reservation.ReservationUpdateStatusExpired}, false},
		{reservation.ReservationStatusUpdateRequest{ReservationID: 42, Status: "invalidReservationStatus"}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV2TestSuite) TestReservationStatusUpdateConfirmationValidation() {
	var confirmationTable = []GenericTestEntry{
		{reservation.ReservationStatusUpdateResponse{}, true},
	}
	ExecuteGenericTestTable(suite, confirmationTable)
}

func (suite *OcppV2TestSuite) TestReservationStatusUpdateE2EMocked() {
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	status := reservation.ReservationUpdateStatusExpired
	reservationID := 42
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"reservationId":%v,"reservationUpdateStatus":"%v"}]`,
		messageId, reservation.ReservationStatusUpdateFeatureName, reservationID, status)
	responseJson := fmt.Sprintf(`[3,"%v",{}]`, messageId)
	dummyResponse := reservation.NewReservationStatusUpdateResponse()
	channel := NewMockWebSocket(wsId)

	handler := &MockCSMSReservationHandler{}
	handler.On("OnReservationStatusUpdate", mock.AnythingOfType("string"), mock.Anything).Return(dummyResponse, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*reservation.ReservationStatusUpdateRequest)
		suite.Require().True(ok)
		suite.Equal(reservationID, request.ReservationID)
		suite.Equal(status, request.Status)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	suite.Require().Nil(err)
	confirmation, err := suite.chargingStation.ReservationStatusUpdate(reservationID, status)
	suite.Nil(err)
	suite.NotNil(confirmation)
}

func (suite *OcppV2TestSuite) TestReservationStatusUpdateInvalidEndpoint() {
	messageId := defaultMessageId
	status := reservation.ReservationUpdateStatusExpired
	reservationID := 42
	request := reservation.NewReservationStatusUpdateRequest(reservationID, status)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"reservationId":%v,"reservationUpdateStatus":"%v"}]`,
		messageId, reservation.ReservationStatusUpdateFeatureName, reservationID, status)
	testUnsupportedRequestFromCentralSystem(suite, request, requestJson, messageId)
}
