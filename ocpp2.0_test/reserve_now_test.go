package ocpp2_test

import (
	"fmt"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0/reservation"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test
func (suite *OcppV2TestSuite) TestReserveNowRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{reservation.ReserveNowRequest{ID: 42, ExpiryDateTime: types.NewDateTime(time.Now()), ConnectorType: reservation.ConnectorTypeCCS1, EvseID: newInt(1), IdToken: types.IdTokenTypeKeyCode, GroupIdToken: types.IdTokenTypeISO15693}, true},
		{reservation.ReserveNowRequest{ID: 42, ExpiryDateTime: types.NewDateTime(time.Now()), ConnectorType: reservation.ConnectorTypeCCS1, EvseID: newInt(1), IdToken: types.IdTokenTypeKeyCode}, true},
		{reservation.ReserveNowRequest{ID: 42, ExpiryDateTime: types.NewDateTime(time.Now()), ConnectorType: reservation.ConnectorTypeCCS1, IdToken: types.IdTokenTypeKeyCode}, true},
		{reservation.ReserveNowRequest{ID: 42, ExpiryDateTime: types.NewDateTime(time.Now()), IdToken: types.IdTokenTypeKeyCode}, true},
		{reservation.ReserveNowRequest{ExpiryDateTime: types.NewDateTime(time.Now()), IdToken: types.IdTokenTypeKeyCode}, true},
		{reservation.ReserveNowRequest{ID: 42, ExpiryDateTime: types.NewDateTime(time.Now())}, false},
		{reservation.ReserveNowRequest{ID: 42, IdToken: types.IdTokenTypeKeyCode}, false},
		{reservation.ReserveNowRequest{}, false},
		{reservation.ReserveNowRequest{ID: -1, ExpiryDateTime: types.NewDateTime(time.Now()), ConnectorType: reservation.ConnectorTypeCCS1, EvseID: newInt(1), IdToken: types.IdTokenTypeKeyCode, GroupIdToken: types.IdTokenTypeISO15693}, false},
		{reservation.ReserveNowRequest{ID: 42, ExpiryDateTime: types.NewDateTime(time.Now()), ConnectorType: "invalidConnectorType", EvseID: newInt(1), IdToken: types.IdTokenTypeKeyCode, GroupIdToken: types.IdTokenTypeISO15693}, false},
		{reservation.ReserveNowRequest{ID: 42, ExpiryDateTime: types.NewDateTime(time.Now()), ConnectorType: reservation.ConnectorTypeCCS1, EvseID: newInt(-1), IdToken: types.IdTokenTypeKeyCode, GroupIdToken: types.IdTokenTypeISO15693}, false},
		{reservation.ReserveNowRequest{ID: 42, ExpiryDateTime: types.NewDateTime(time.Now()), ConnectorType: reservation.ConnectorTypeCCS1, EvseID: newInt(1), IdToken: "invalidIdToken", GroupIdToken: types.IdTokenTypeISO15693}, false},
		{reservation.ReserveNowRequest{ID: 42, ExpiryDateTime: types.NewDateTime(time.Now()), ConnectorType: reservation.ConnectorTypeCCS1, EvseID: newInt(1), IdToken: types.IdTokenTypeKeyCode, GroupIdToken: "invalidIdToken"}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestReserveNowConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{reservation.ReserveNowResponse{Status: reservation.ReserveNowStatusAccepted, StatusInfo: &types.StatusInfo{ReasonCode: "200"}}, true},
		{reservation.ReserveNowResponse{Status: reservation.ReserveNowStatusAccepted}, true},
		{reservation.ReserveNowResponse{}, false},
		{reservation.ReserveNowResponse{Status: "invalidReserveNowStatus"}, false},
		{reservation.ReserveNowResponse{Status: reservation.ReserveNowStatusAccepted, StatusInfo: &types.StatusInfo{}}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestReserveNowE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	id := 42
	expiryDateTime := types.NewDateTime(time.Now())
	connectorType := reservation.ConnectorTypeCCS1
	evseID := newInt(1)
	idToken := types.IdTokenTypeKeyCode
	groupIdToken := types.IdTokenTypeISO15693
	status := reservation.ReserveNowStatusAccepted
	statusInfo := types.StatusInfo{ReasonCode: "200"}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"id":%v,"expiryDateTime":"%v","connectorType":"%v","evseId":%v,"idToken":"%v","groupIdToken":"%v"}]`,
		messageId, reservation.ReserveNowFeatureName, id, expiryDateTime.FormatTimestamp(), connectorType, *evseID, idToken, groupIdToken)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v","statusInfo":{"reasonCode":"%v"}}]`,
		messageId, status, statusInfo.ReasonCode)
	reserveNowResponse := reservation.NewReserveNowResponse(status)
	reserveNowResponse.StatusInfo = &statusInfo
	channel := NewMockWebSocket(wsId)

	handler := MockChargingStationReservationHandler{}
	handler.On("OnReserveNow", mock.Anything).Return(reserveNowResponse, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*reservation.ReserveNowRequest)
		require.True(t, ok)
		assert.Equal(t, id, request.ID)
		assert.Equal(t, expiryDateTime.FormatTimestamp(), request.ExpiryDateTime.FormatTimestamp())
		assert.Equal(t, connectorType, request.ConnectorType)
		assert.Equal(t, *evseID, *request.EvseID)
		assert.Equal(t, idToken, request.IdToken)
		assert.Equal(t, groupIdToken, request.GroupIdToken)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.ReserveNow(wsId, func(resp *reservation.ReserveNowResponse, err error) {
		require.Nil(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, status, resp.Status)
		assert.Equal(t, statusInfo.ReasonCode, resp.StatusInfo.ReasonCode)
		resultChannel <- true
	}, id, expiryDateTime, idToken, func(request *reservation.ReserveNowRequest) {
		request.ConnectorType = connectorType
		request.EvseID = evseID
		request.GroupIdToken = groupIdToken
	})
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV2TestSuite) TestReserveNowInvalidEndpoint() {
	messageId := defaultMessageId
	id := 42
	expiryDateTime := types.NewDateTime(time.Now())
	connectorType := reservation.ConnectorTypeCCS1
	evseID := newInt(1)
	idToken := types.IdTokenTypeKeyCode
	groupIdToken := types.IdTokenTypeISO15693
	reserveNowRequest := reservation.ReserveNowRequest{
		ID:             id,
		ExpiryDateTime: expiryDateTime,
		ConnectorType:  connectorType,
		EvseID:         evseID,
		IdToken:        idToken,
		GroupIdToken:   groupIdToken,
	}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"id":%v,"expiryDateTime":"%v","connectorType":"%v","evseId":%v,"idToken":"%v","groupIdToken":"%v"}]`,
		messageId, reservation.ReserveNowFeatureName, id, expiryDateTime.FormatTimestamp(), connectorType, *evseID, idToken, groupIdToken)
	testUnsupportedRequestFromChargingStation(suite, reserveNowRequest, requestJson, messageId)
}
