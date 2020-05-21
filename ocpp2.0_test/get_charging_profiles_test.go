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
func (suite *OcppV2TestSuite) TestGetChargingProfilesRequestValidation() {
	t := suite.T()
	validChargingProfileCriterion := smartcharging.ChargingProfileCriterion{
		ChargingProfilePurpose: types.ChargingProfilePurposeTxDefaultProfile,
		StackLevel:             newInt(2),
		ChargingProfileID:      []int{1, 2},
		ChargingLimitSource:    []types.ChargingLimitSourceType{types.ChargingLimitSourceEMS},
	}
	var requestTable = []GenericTestEntry{
		{smartcharging.GetChargingProfilesRequest{RequestID: newInt(42), EvseID: newInt(1), ChargingProfile: validChargingProfileCriterion}, true},
		{smartcharging.GetChargingProfilesRequest{RequestID: newInt(42), ChargingProfile: validChargingProfileCriterion}, true},
		{smartcharging.GetChargingProfilesRequest{EvseID: newInt(1), ChargingProfile: validChargingProfileCriterion}, true},
		{smartcharging.GetChargingProfilesRequest{ChargingProfile: validChargingProfileCriterion}, true},
		{smartcharging.GetChargingProfilesRequest{ChargingProfile: smartcharging.ChargingProfileCriterion{}}, true},
		{smartcharging.GetChargingProfilesRequest{}, true},
		{smartcharging.GetChargingProfilesRequest{RequestID: newInt(42), EvseID: newInt(-1), ChargingProfile: validChargingProfileCriterion}, false},
		{smartcharging.GetChargingProfilesRequest{RequestID: newInt(-1), EvseID: newInt(1), ChargingProfile: validChargingProfileCriterion}, false},
		{smartcharging.GetChargingProfilesRequest{ChargingProfile: smartcharging.ChargingProfileCriterion{ChargingProfilePurpose: "invalidChargingProfilePurpose", StackLevel: newInt(2), ChargingProfileID: []int{1, 2}, ChargingLimitSource: []types.ChargingLimitSourceType{types.ChargingLimitSourceEMS}}}, false},
		{smartcharging.GetChargingProfilesRequest{ChargingProfile: smartcharging.ChargingProfileCriterion{ChargingProfilePurpose: types.ChargingProfilePurposeTxDefaultProfile, StackLevel: newInt(-1), ChargingProfileID: []int{1, 2}, ChargingLimitSource: []types.ChargingLimitSourceType{types.ChargingLimitSourceEMS}}}, false},
		{smartcharging.GetChargingProfilesRequest{ChargingProfile: smartcharging.ChargingProfileCriterion{ChargingProfilePurpose: types.ChargingProfilePurposeTxDefaultProfile, StackLevel: newInt(2), ChargingProfileID: []int{1, 2}, ChargingLimitSource: []types.ChargingLimitSourceType{types.ChargingLimitSourceEMS, types.ChargingLimitSourceCSO, types.ChargingLimitSourceSO, types.ChargingLimitSourceOther, types.ChargingLimitSourceEMS}}}, false},
		{smartcharging.GetChargingProfilesRequest{ChargingProfile: smartcharging.ChargingProfileCriterion{ChargingProfilePurpose: types.ChargingProfilePurposeTxDefaultProfile, StackLevel: newInt(2), ChargingProfileID: []int{1, 2}, ChargingLimitSource: []types.ChargingLimitSourceType{"invalidChargingLimitSource"}}}, false},
		{smartcharging.GetChargingProfilesRequest{ChargingProfile: smartcharging.ChargingProfileCriterion{ChargingProfilePurpose: types.ChargingProfilePurposeTxDefaultProfile, StackLevel: newInt(2), ChargingProfileID: []int{-1}, ChargingLimitSource: []types.ChargingLimitSourceType{types.ChargingLimitSourceEMS}}}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestGetChargingProfilesConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{smartcharging.GetChargingProfilesResponse{Status: smartcharging.GetChargingProfileStatusAccepted}, true},
		{smartcharging.GetChargingProfilesResponse{Status: smartcharging.GetChargingProfileStatusNoProfiles}, true},
		{smartcharging.GetChargingProfilesResponse{Status: "invalidGetChargingProfilesStatus"}, false},
		{smartcharging.GetChargingProfilesResponse{}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestGetChargingProfilesE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	requestID := 42
	evseID := 1
	chargingProfileCriterion := smartcharging.ChargingProfileCriterion{
		ChargingProfilePurpose: types.ChargingProfilePurposeChargingStationMaxProfile,
		StackLevel:             newInt(1),
		ChargingProfileID:      []int{1, 2},
		ChargingLimitSource:    []types.ChargingLimitSourceType{types.ChargingLimitSourceEMS},
	}
	status := smartcharging.GetChargingProfileStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"requestId":%v,"evseId":%v,"chargingProfile":{"chargingProfilePurpose":"%v","stackLevel":%v,"chargingProfileId":[%v,%v],"chargingLimitSource":["%v"]}}]`,
		messageId, smartcharging.GetChargingProfilesFeatureName, requestID, evseID, chargingProfileCriterion.ChargingProfilePurpose, *chargingProfileCriterion.StackLevel, chargingProfileCriterion.ChargingProfileID[0], chargingProfileCriterion.ChargingProfileID[1], chargingProfileCriterion.ChargingLimitSource[0])
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	getChargingProfilesConfirmation := smartcharging.NewGetChargingProfilesResponse(status)
	channel := NewMockWebSocket(wsId)

	handler := MockChargingStationSmartChargingHandler{}
	handler.On("OnGetChargingProfiles", mock.Anything).Return(getChargingProfilesConfirmation, nil).Run(func(args mock.Arguments) {
		// Assert request message contents
		request, ok := args.Get(0).(*smartcharging.GetChargingProfilesRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, requestID, *request.RequestID)
		assert.Equal(t, evseID, *request.EvseID)
		assert.Equal(t, chargingProfileCriterion.ChargingProfilePurpose, request.ChargingProfile.ChargingProfilePurpose)
		assert.Equal(t, *chargingProfileCriterion.StackLevel, *request.ChargingProfile.StackLevel)
		require.Len(t, request.ChargingProfile.ChargingProfileID, len(chargingProfileCriterion.ChargingProfileID))
		assert.Equal(t, chargingProfileCriterion.ChargingProfileID[0], request.ChargingProfile.ChargingProfileID[0])
		assert.Equal(t, chargingProfileCriterion.ChargingProfileID[1], request.ChargingProfile.ChargingProfileID[1])
		require.Len(t, request.ChargingProfile.ChargingLimitSource, len(chargingProfileCriterion.ChargingLimitSource))
		assert.Equal(t, chargingProfileCriterion.ChargingLimitSource[0], request.ChargingProfile.ChargingLimitSource[0])
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.GetChargingProfiles(wsId, func(confirmation *smartcharging.GetChargingProfilesResponse, err error) {
		// Assert confirmation message contents
		require.Nil(t, err)
		require.NotNil(t, confirmation)
		assert.Equal(t, status, confirmation.Status)
		resultChannel <- true
	},
		chargingProfileCriterion,
		func(request *smartcharging.GetChargingProfilesRequest) {
			request.EvseID = &evseID
			request.RequestID = &requestID
		})
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV2TestSuite) TestGetChargingProfilesInvalidEndpoint() {
	messageId := defaultMessageId
	requestID := 42
	evseID := 1
	chargingProfileCriterion := smartcharging.ChargingProfileCriterion{
		ChargingProfilePurpose: types.ChargingProfilePurposeChargingStationMaxProfile,
		StackLevel:             newInt(1),
		ChargingProfileID:      []int{1, 2},
		ChargingLimitSource:    []types.ChargingLimitSourceType{types.ChargingLimitSourceEMS},
	}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"requestId":%v,"evseId":%v,"chargingProfile":{"chargingProfilePurpose":"%v","stackLevel":%v,"chargingProfileId":[%v,%v],"chargingLimitSource":["%v"]}}]`,
		messageId, smartcharging.GetChargingProfilesFeatureName, requestID, evseID, chargingProfileCriterion.ChargingProfilePurpose, *chargingProfileCriterion.StackLevel, chargingProfileCriterion.ChargingProfileID[0], chargingProfileCriterion.ChargingProfileID[1], chargingProfileCriterion.ChargingLimitSource[0])
	getChargingProfilesRequest := smartcharging.NewGetChargingProfilesRequest(chargingProfileCriterion)
	getChargingProfilesRequest.EvseID = &evseID
	getChargingProfilesRequest.RequestID = &requestID
	testUnsupportedRequestFromChargingStation(suite, getChargingProfilesRequest, requestJson, messageId)
}
