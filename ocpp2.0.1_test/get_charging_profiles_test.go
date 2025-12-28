package ocpp2_test

import (
	"fmt"

	"github.com/stretchr/testify/mock"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/smartcharging"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

// Test
func (suite *OcppV2TestSuite) TestGetChargingProfilesRequestValidation() {
	validChargingProfileCriterion := smartcharging.ChargingProfileCriterion{
		ChargingProfilePurpose: types.ChargingProfilePurposeTxDefaultProfile,
		StackLevel:             newInt(2),
		ChargingProfileID:      []int{1, 2},
		ChargingLimitSource:    []types.ChargingLimitSourceType{types.ChargingLimitSourceEMS},
	}
	var requestTable = []GenericTestEntry{
		{smartcharging.GetChargingProfilesRequest{RequestID: 42, EvseID: newInt(1), ChargingProfile: validChargingProfileCriterion}, true},
		{smartcharging.GetChargingProfilesRequest{RequestID: 42, ChargingProfile: validChargingProfileCriterion}, true},
		{smartcharging.GetChargingProfilesRequest{EvseID: newInt(1), ChargingProfile: validChargingProfileCriterion}, true},
		{smartcharging.GetChargingProfilesRequest{ChargingProfile: validChargingProfileCriterion}, true},
		{smartcharging.GetChargingProfilesRequest{ChargingProfile: smartcharging.ChargingProfileCriterion{}}, true},
		{smartcharging.GetChargingProfilesRequest{}, true},
		{smartcharging.GetChargingProfilesRequest{RequestID: 42, EvseID: newInt(-1), ChargingProfile: validChargingProfileCriterion}, false},
		{smartcharging.GetChargingProfilesRequest{ChargingProfile: smartcharging.ChargingProfileCriterion{ChargingProfilePurpose: "invalidChargingProfilePurpose", StackLevel: newInt(2), ChargingProfileID: []int{1, 2}, ChargingLimitSource: []types.ChargingLimitSourceType{types.ChargingLimitSourceEMS}}}, false},
		{smartcharging.GetChargingProfilesRequest{ChargingProfile: smartcharging.ChargingProfileCriterion{ChargingProfilePurpose: types.ChargingProfilePurposeTxDefaultProfile, StackLevel: newInt(-1), ChargingProfileID: []int{1, 2}, ChargingLimitSource: []types.ChargingLimitSourceType{types.ChargingLimitSourceEMS}}}, false},
		{smartcharging.GetChargingProfilesRequest{ChargingProfile: smartcharging.ChargingProfileCriterion{ChargingProfilePurpose: types.ChargingProfilePurposeTxDefaultProfile, StackLevel: newInt(2), ChargingProfileID: []int{1, 2}, ChargingLimitSource: []types.ChargingLimitSourceType{types.ChargingLimitSourceEMS, types.ChargingLimitSourceCSO, types.ChargingLimitSourceSO, types.ChargingLimitSourceOther, types.ChargingLimitSourceEMS}}}, false},
		{smartcharging.GetChargingProfilesRequest{ChargingProfile: smartcharging.ChargingProfileCriterion{ChargingProfilePurpose: types.ChargingProfilePurposeTxDefaultProfile, StackLevel: newInt(2), ChargingProfileID: []int{1, 2}, ChargingLimitSource: []types.ChargingLimitSourceType{"invalidChargingLimitSource"}}}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV2TestSuite) TestGetChargingProfilesConfirmationValidation() {
	var confirmationTable = []GenericTestEntry{
		{smartcharging.GetChargingProfilesResponse{Status: smartcharging.GetChargingProfileStatusAccepted}, true},
		{smartcharging.GetChargingProfilesResponse{Status: smartcharging.GetChargingProfileStatusNoProfiles}, true},
		{smartcharging.GetChargingProfilesResponse{Status: "invalidGetChargingProfilesStatus"}, false},
		{smartcharging.GetChargingProfilesResponse{}, false},
	}
	ExecuteGenericTestTable(suite, confirmationTable)
}

func (suite *OcppV2TestSuite) TestGetChargingProfilesE2EMocked() {
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

	handler := &MockChargingStationSmartChargingHandler{}
	handler.On("OnGetChargingProfiles", mock.Anything).Return(getChargingProfilesConfirmation, nil).Run(func(args mock.Arguments) {
		// Assert request message contents
		request, ok := args.Get(0).(*smartcharging.GetChargingProfilesRequest)
		suite.Require().True(ok)
		suite.Require().NotNil(request)
		suite.Equal(requestID, request.RequestID)
		suite.Equal(evseID, *request.EvseID)
		suite.Equal(chargingProfileCriterion.ChargingProfilePurpose, request.ChargingProfile.ChargingProfilePurpose)
		suite.Equal(*chargingProfileCriterion.StackLevel, *request.ChargingProfile.StackLevel)
		suite.Require().Len(request.ChargingProfile.ChargingProfileID, len(chargingProfileCriterion.ChargingProfileID))
		suite.Equal(chargingProfileCriterion.ChargingProfileID[0], request.ChargingProfile.ChargingProfileID[0])
		suite.Equal(chargingProfileCriterion.ChargingProfileID[1], request.ChargingProfile.ChargingProfileID[1])
		suite.Require().Len(request.ChargingProfile.ChargingLimitSource, len(chargingProfileCriterion.ChargingLimitSource))
		suite.Equal(chargingProfileCriterion.ChargingLimitSource[0], request.ChargingProfile.ChargingLimitSource[0])
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	suite.Require().Nil(err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.GetChargingProfiles(wsId, func(confirmation *smartcharging.GetChargingProfilesResponse, err error) {
		// Assert confirmation message contents
		suite.Require().Nil(err)
		suite.Require().NotNil(confirmation)
		suite.Equal(status, confirmation.Status)
		resultChannel <- true
	},
		chargingProfileCriterion,
		func(request *smartcharging.GetChargingProfilesRequest) {
			request.EvseID = &evseID
			request.RequestID = requestID
		})
	suite.Require().Nil(err)
	result := <-resultChannel
	suite.True(result)
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
	getChargingProfilesRequest.RequestID = requestID
	testUnsupportedRequestFromChargingStation(suite, getChargingProfilesRequest, requestJson, messageId)
}
