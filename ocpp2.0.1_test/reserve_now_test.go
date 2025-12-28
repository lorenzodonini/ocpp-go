package ocpp2_test

import (
	"fmt"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/reservation"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"

	"github.com/stretchr/testify/mock"
)

// Test
func (suite *OcppV2TestSuite) TestReserveNowRequestValidation() {
	var requestTable = []GenericTestEntry{
		{reservation.ReserveNowRequest{ID: 42, ExpiryDateTime: types.NewDateTime(time.Now()), ConnectorType: reservation.ConnectorTypeCCS1, EvseID: newInt(1), IdToken: types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode}, GroupIdToken: &types.IdToken{IdToken: "1234", Type: types.IdTokenTypeISO15693}}, true},
		{reservation.ReserveNowRequest{ID: 42, ExpiryDateTime: types.NewDateTime(time.Now()), ConnectorType: reservation.ConnectorTypeCCS1, EvseID: newInt(1), IdToken: types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode}}, true},
		{reservation.ReserveNowRequest{ID: 42, ExpiryDateTime: types.NewDateTime(time.Now()), ConnectorType: reservation.ConnectorTypeCCS1, IdToken: types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode}}, true},
		{reservation.ReserveNowRequest{ID: 42, ExpiryDateTime: types.NewDateTime(time.Now()), IdToken: types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode}}, true},
		{reservation.ReserveNowRequest{ExpiryDateTime: types.NewDateTime(time.Now()), IdToken: types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode}}, true},
		{reservation.ReserveNowRequest{ID: 42, ExpiryDateTime: types.NewDateTime(time.Now())}, false},
		{reservation.ReserveNowRequest{ID: 42, IdToken: types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode}}, false},
		{reservation.ReserveNowRequest{}, false},
		{reservation.ReserveNowRequest{ID: -1, ExpiryDateTime: types.NewDateTime(time.Now()), ConnectorType: reservation.ConnectorTypeCCS1, EvseID: newInt(1), IdToken: types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode}, GroupIdToken: &types.IdToken{IdToken: "1234", Type: types.IdTokenTypeISO15693}}, false},
		{reservation.ReserveNowRequest{ID: 42, ExpiryDateTime: types.NewDateTime(time.Now()), ConnectorType: "invalidConnectorType", EvseID: newInt(1), IdToken: types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode}, GroupIdToken: &types.IdToken{IdToken: "1234", Type: types.IdTokenTypeISO15693}}, false},
		{reservation.ReserveNowRequest{ID: 42, ExpiryDateTime: types.NewDateTime(time.Now()), ConnectorType: reservation.ConnectorTypeCCS1, EvseID: newInt(-1), IdToken: types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode}, GroupIdToken: &types.IdToken{IdToken: "1234", Type: types.IdTokenTypeISO15693}}, false},
		{reservation.ReserveNowRequest{ID: 42, ExpiryDateTime: types.NewDateTime(time.Now()), ConnectorType: reservation.ConnectorTypeCCS1, EvseID: newInt(1), IdToken: types.IdToken{IdToken: "1234", Type: "invalidIdToken"}, GroupIdToken: &types.IdToken{IdToken: "1234", Type: types.IdTokenTypeISO15693}}, false},
		{reservation.ReserveNowRequest{ID: 42, ExpiryDateTime: types.NewDateTime(time.Now()), ConnectorType: reservation.ConnectorTypeCCS1, EvseID: newInt(1), IdToken: types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode}, GroupIdToken: &types.IdToken{IdToken: "1234", Type: "invalidIdToken"}}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV2TestSuite) TestReserveNowConfirmationValidation() {
	var confirmationTable = []GenericTestEntry{
		{reservation.ReserveNowResponse{Status: reservation.ReserveNowStatusAccepted, StatusInfo: &types.StatusInfo{ReasonCode: "200"}}, true},
		{reservation.ReserveNowResponse{Status: reservation.ReserveNowStatusAccepted}, true},
		{reservation.ReserveNowResponse{}, false},
		{reservation.ReserveNowResponse{Status: "invalidReserveNowStatus"}, false},
		{reservation.ReserveNowResponse{Status: reservation.ReserveNowStatusAccepted, StatusInfo: &types.StatusInfo{}}, false},
	}
	ExecuteGenericTestTable(suite, confirmationTable)
}

func (suite *OcppV2TestSuite) TestReserveNowE2EMocked() {
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	id := 42
	expiryDateTime := types.NewDateTime(time.Now())
	connectorType := reservation.ConnectorTypeCCS1
	evseID := newInt(1)
	idToken := types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode}
	groupIdToken := types.IdToken{IdToken: "1234", Type: types.IdTokenTypeISO15693}
	status := reservation.ReserveNowStatusAccepted
	statusInfo := types.StatusInfo{ReasonCode: "200"}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"id":%v,"expiryDateTime":"%v","connectorType":"%v","evseId":%v,"idToken":{"idToken":"%s","type":"%s"},"groupIdToken":{"idToken":"%s","type":"%s"}}]`,
		messageId, reservation.ReserveNowFeatureName, id, expiryDateTime.FormatTimestamp(), connectorType, *evseID, idToken.IdToken, idToken.Type, groupIdToken.IdToken, groupIdToken.Type)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v","statusInfo":{"reasonCode":"%v"}}]`,
		messageId, status, statusInfo.ReasonCode)
	reserveNowResponse := reservation.NewReserveNowResponse(status)
	reserveNowResponse.StatusInfo = &statusInfo
	channel := NewMockWebSocket(wsId)

	handler := &MockChargingStationReservationHandler{}
	handler.On("OnReserveNow", mock.Anything).Return(reserveNowResponse, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*reservation.ReserveNowRequest)
		suite.Require().True(ok)
		suite.Equal(id, request.ID)
		suite.Equal(expiryDateTime.FormatTimestamp(), request.ExpiryDateTime.FormatTimestamp())
		suite.Equal(connectorType, request.ConnectorType)
		suite.Equal(*evseID, *request.EvseID)
		suite.Equal(idToken.IdToken, request.IdToken.IdToken)
		suite.Equal(idToken.Type, request.IdToken.Type)
		suite.Require().NotNil(request.GroupIdToken)
		suite.Equal(groupIdToken.IdToken, request.GroupIdToken.IdToken)
		suite.Equal(groupIdToken.Type, request.GroupIdToken.Type)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	suite.Require().Nil(err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.ReserveNow(wsId, func(resp *reservation.ReserveNowResponse, err error) {
		suite.Require().Nil(err)
		suite.Require().NotNil(resp)
		suite.Equal(status, resp.Status)
		suite.Equal(statusInfo.ReasonCode, resp.StatusInfo.ReasonCode)
		resultChannel <- true
	}, id, expiryDateTime, idToken, func(request *reservation.ReserveNowRequest) {
		request.ConnectorType = connectorType
		request.EvseID = evseID
		request.GroupIdToken = &groupIdToken
	})
	suite.Require().Nil(err)
	result := <-resultChannel
	suite.True(result)
}

func (suite *OcppV2TestSuite) TestReserveNowInvalidEndpoint() {
	messageId := defaultMessageId
	id := 42
	expiryDateTime := types.NewDateTime(time.Now())
	connectorType := reservation.ConnectorTypeCCS1
	evseID := newInt(1)
	idToken := types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode}
	groupIdToken := types.IdToken{IdToken: "1234", Type: types.IdTokenTypeISO15693}
	reserveNowRequest := reservation.ReserveNowRequest{
		ID:             id,
		ExpiryDateTime: expiryDateTime,
		ConnectorType:  connectorType,
		EvseID:         evseID,
		IdToken:        idToken,
		GroupIdToken:   &groupIdToken,
	}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"id":%v,"expiryDateTime":"%v","connectorType":"%v","evseId":%v,"idToken":{"idToken":"%s","type":"%s"},"groupIdToken":{"idToken":"%s","type":"%s"}}]`,
		messageId, reservation.ReserveNowFeatureName, id, expiryDateTime.FormatTimestamp(), connectorType, *evseID, idToken.IdToken, idToken.Type, groupIdToken.IdToken, groupIdToken.Type)
	testUnsupportedRequestFromChargingStation(suite, reserveNowRequest, requestJson, messageId)
}
