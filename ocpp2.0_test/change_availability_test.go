package ocpp2_test

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0/availability"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func (suite *OcppV2TestSuite) TestChangeAvailabilityRequestValidation() {
	t := suite.T()
	var testTable = []GenericTestEntry{
		{availability.ChangeAvailabilityRequest{OperationalStatus: availability.OperationalStatusOperative, Evse: &types.EVSE{ID: 1, ConnectorID: newInt(1)}}, true},
		{availability.ChangeAvailabilityRequest{OperationalStatus: availability.OperationalStatusInoperative, Evse: &types.EVSE{ID: 1}}, true},
		{availability.ChangeAvailabilityRequest{OperationalStatus: availability.OperationalStatusInoperative}, true},
		{availability.ChangeAvailabilityRequest{OperationalStatus: availability.OperationalStatusOperative}, true},
		{availability.ChangeAvailabilityRequest{}, false},
		{availability.ChangeAvailabilityRequest{OperationalStatus: "invalidAvailabilityType"}, false},
		{availability.ChangeAvailabilityRequest{OperationalStatus: availability.OperationalStatusOperative, Evse: &types.EVSE{ID: -1}}, false},
	}
	ExecuteGenericTestTable(t, testTable)
}

func (suite *OcppV2TestSuite) TestChangeAvailabilityConfirmationValidation() {
	t := suite.T()
	var testTable = []GenericTestEntry{
		{availability.ChangeAvailabilityResponse{Status: availability.ChangeAvailabilityStatusAccepted}, true},
		{availability.ChangeAvailabilityResponse{Status: availability.ChangeAvailabilityStatusRejected}, true},
		{availability.ChangeAvailabilityResponse{Status: availability.ChangeAvailabilityStatusScheduled}, true},
		{availability.ChangeAvailabilityResponse{Status: "invalidAvailabilityStatus"}, false},
		{availability.ChangeAvailabilityResponse{}, false},
	}
	ExecuteGenericTestTable(t, testTable)
}

// Test
func (suite *OcppV2TestSuite) TestChangeAvailabilityE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	evse := types.EVSE{ID: 1, ConnectorID: newInt(1)}
	operationalStatus := availability.OperationalStatusOperative
	status := availability.ChangeAvailabilityStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"operationalStatus":"%v","evse":{"id":%v,"connectorId":%v}}]`,
		messageId, availability.ChangeAvailabilityFeatureName, operationalStatus, evse.ID, *evse.ConnectorID)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	changeAvailabilityConfirmation := availability.NewChangeAvailabilityResponse(status)
	channel := NewMockWebSocket(wsId)
	// Setting handlers
	handler := MockChargingStationAvailabilityHandler{}
	handler.On("OnChangeAvailability", mock.Anything).Return(changeAvailabilityConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*availability.ChangeAvailabilityRequest)
		require.True(t, ok)
		require.NotNil(t, request.Evse)
		assert.Equal(t, evse.ID, request.Evse.ID)
		assert.Equal(t, *evse.ConnectorID, *request.Evse.ConnectorID)
		assert.Equal(t, operationalStatus, request.OperationalStatus)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.ChangeAvailability(wsId, func(confirmation *availability.ChangeAvailabilityResponse, err error) {
		require.Nil(t, err)
		require.NotNil(t, confirmation)
		assert.Equal(t, status, confirmation.Status)
		resultChannel <- true
	}, operationalStatus, func(req *availability.ChangeAvailabilityRequest) {
		req.Evse = &evse
	})
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV2TestSuite) TestChangeAvailabilityInvalidEndpoint() {
	messageId := defaultMessageId
	evse := types.EVSE{ID: 1, ConnectorID: newInt(1)}
	operationalStatus := availability.OperationalStatusOperative
	changeAvailabilityRequest := availability.NewChangeAvailabilityRequest(operationalStatus)
	changeAvailabilityRequest.Evse = &evse
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"operationalStatus":"%v","evse":{"id":%v,"connectorId":%v}}]`,
		messageId, availability.ChangeAvailabilityFeatureName, operationalStatus, evse.ID, *evse.ConnectorID)
	testUnsupportedRequestFromChargingStation(suite, changeAvailabilityRequest, requestJson, messageId)
}
