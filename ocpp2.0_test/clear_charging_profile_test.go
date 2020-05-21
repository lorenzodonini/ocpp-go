package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/smartcharging"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test
func (suite *OcppV2TestSuite) TestClearChargingProfileRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{smartcharging.ClearChargingProfileRequest{EvseID: newInt(1), ChargingProfile: &smartcharging.ClearChargingProfileType{ID: 1, ChargingProfilePurpose: types.ChargingProfilePurposeChargingStationMaxProfile, StackLevel: 1}}, true},
		{smartcharging.ClearChargingProfileRequest{EvseID: newInt(1), ChargingProfile: &smartcharging.ClearChargingProfileType{ID: 1, ChargingProfilePurpose: types.ChargingProfilePurposeChargingStationMaxProfile}}, true},
		{smartcharging.ClearChargingProfileRequest{EvseID: newInt(1), ChargingProfile: &smartcharging.ClearChargingProfileType{ID: 1}}, true},
		{smartcharging.ClearChargingProfileRequest{ChargingProfile: &smartcharging.ClearChargingProfileType{ID: 1}}, true},
		{smartcharging.ClearChargingProfileRequest{ChargingProfile: &smartcharging.ClearChargingProfileType{}}, true},
		{smartcharging.ClearChargingProfileRequest{}, true},
		{smartcharging.ClearChargingProfileRequest{EvseID: newInt(-1)}, false},
		{smartcharging.ClearChargingProfileRequest{ChargingProfile: &smartcharging.ClearChargingProfileType{ID: -1}}, false},
		{smartcharging.ClearChargingProfileRequest{ChargingProfile: &smartcharging.ClearChargingProfileType{ChargingProfilePurpose: "invalidChargingProfilePurposeType"}}, false},
		{smartcharging.ClearChargingProfileRequest{ChargingProfile: &smartcharging.ClearChargingProfileType{StackLevel: -1}}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestClearChargingProfileConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{smartcharging.ClearChargingProfileResponse{Status: smartcharging.ClearChargingProfileStatusAccepted}, true},
		{smartcharging.ClearChargingProfileResponse{Status: "invalidClearChargingProfileStatus"}, false},
		{smartcharging.ClearChargingProfileResponse{}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestClearChargingProfileE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	chargingProfileId := 1
	evseID := 1
	chargingProfilePurpose := types.ChargingProfilePurposeChargingStationMaxProfile
	stackLevel := 1
	status := smartcharging.ClearChargingProfileStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"evseId":%v,"chargingProfile":{"id":%v,"chargingProfilePurpose":"%v","stackLevel":%v}}]`,
		messageId, smartcharging.ClearChargingProfileFeatureName, evseID, chargingProfileId, chargingProfilePurpose, stackLevel)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	clearChargingProfileConfirmation := smartcharging.NewClearChargingProfileResponse(status)
	channel := NewMockWebSocket(wsId)

	handler := MockChargingStationSmartChargingHandler{}
	handler.On("OnClearChargingProfile", mock.Anything).Return(clearChargingProfileConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*smartcharging.ClearChargingProfileRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, evseID, *request.EvseID)
		assert.Equal(t, chargingProfileId, request.ChargingProfile.ID)
		assert.Equal(t, chargingProfilePurpose, request.ChargingProfile.ChargingProfilePurpose)
		assert.Equal(t, stackLevel, request.ChargingProfile.StackLevel)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.ClearChargingProfile(wsId, func(confirmation *smartcharging.ClearChargingProfileResponse, err error) {
		require.Nil(t, err)
		require.NotNil(t, confirmation)
		assert.Equal(t, status, confirmation.Status)
		resultChannel <- true
	}, func(request *smartcharging.ClearChargingProfileRequest) {
		request.EvseID = &evseID
		request.ChargingProfile = &smartcharging.ClearChargingProfileType{ID: chargingProfileId, ChargingProfilePurpose: chargingProfilePurpose, StackLevel: stackLevel}
	})
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV2TestSuite) TestClearChargingProfileInvalidEndpoint() {
	messageId := defaultMessageId
	evseID := 1
	chargingProfileId := 1
	stackLevel := 1
	chargingProfilePurpose := types.ChargingProfilePurposeChargingStationMaxProfile
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"evseId":%v,"chargingProfile":{"id":%v,"chargingProfilePurpose":"%v","stackLevel":%v}}]`,
		messageId, smartcharging.ClearChargingProfileFeatureName, evseID, chargingProfileId, chargingProfilePurpose, stackLevel)
	clearChargingProfileRequest := smartcharging.NewClearChargingProfileRequest()
	clearChargingProfileRequest.EvseID = &evseID
	clearChargingProfileRequest.ChargingProfile = &smartcharging.ClearChargingProfileType{ID: chargingProfileId, ChargingProfilePurpose: chargingProfilePurpose, StackLevel: stackLevel}
	testUnsupportedRequestFromChargingStation(suite, clearChargingProfileRequest, requestJson, messageId)
}
