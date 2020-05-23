package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/availability"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func (suite *OcppV2TestSuite) TestChangeAvailabilityRequestValidation() {
	t := suite.T()
	var testTable = []GenericTestEntry{
		{availability.ChangeAvailabilityRequest{EvseID: 0, OperationalStatus: availability.OperationalStatusOperative}, true},
		{availability.ChangeAvailabilityRequest{EvseID: 0, OperationalStatus: availability.OperationalStatusInoperative}, true},
		{availability.ChangeAvailabilityRequest{EvseID: 0}, false},
		{availability.ChangeAvailabilityRequest{OperationalStatus: availability.OperationalStatusOperative}, true},
		{availability.ChangeAvailabilityRequest{OperationalStatus: "invalidAvailabilityType"}, false},
		{availability.ChangeAvailabilityRequest{EvseID: -1, OperationalStatus: availability.OperationalStatusOperative}, false},
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
	evseID := 1
	operationalStatus := availability.OperationalStatusOperative
	status := availability.ChangeAvailabilityStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"evseId":%v,"operationalStatus":"%v"}]`, messageId, availability.ChangeAvailabilityFeatureName, evseID, operationalStatus)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	changeAvailabilityConfirmation := availability.NewChangeAvailabilityResponse(status)
	channel := NewMockWebSocket(wsId)
	// Setting handlers
	handler := MockChargingStationAvailabilityHandler{}
	handler.On("OnChangeAvailability", mock.Anything).Return(changeAvailabilityConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*availability.ChangeAvailabilityRequest)
		require.True(t, ok)
		assert.Equal(t, evseID, request.EvseID)
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
	}, evseID, operationalStatus)
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV2TestSuite) TestChangeAvailabilityInvalidEndpoint() {
	messageId := defaultMessageId
	evseID := 1
	operationalStatus := availability.OperationalStatusOperative
	changeAvailabilityRequest := availability.NewChangeAvailabilityRequest(evseID, operationalStatus)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"evseId":%v,"operationalStatus":"%v"}]`, messageId, availability.ChangeAvailabilityFeatureName, evseID, operationalStatus)
	testUnsupportedRequestFromChargingStation(suite, changeAvailabilityRequest, requestJson, messageId)
}
