package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test
func (suite *OcppV2TestSuite) TestGetChargingProfilesRequestValidation() {
	t := suite.T()
	validChargingProfileCriterion := ocpp2.ChargingProfileCriterion{
		ChargingProfilePurpose: ocpp2.ChargingProfilePurposeTxDefaultProfile,
		StackLevel:             newInt(2),
		ChargingProfileID:      []int{1,2},
		ChargingLimitSource:    []ocpp2.ChargingLimitSourceType{ocpp2.ChargingLimitSourceEMS},
	}
	var requestTable = []GenericTestEntry{
		{ocpp2.GetChargingProfilesRequest{RequestID: newInt(42), EvseID: newInt(1), ChargingProfile: validChargingProfileCriterion}, true},
		{ocpp2.GetChargingProfilesRequest{RequestID: newInt(42), ChargingProfile: validChargingProfileCriterion}, true},
		{ocpp2.GetChargingProfilesRequest{EvseID: newInt(1), ChargingProfile: validChargingProfileCriterion}, true},
		{ocpp2.GetChargingProfilesRequest{ChargingProfile: validChargingProfileCriterion}, true},
		{ocpp2.GetChargingProfilesRequest{ChargingProfile: ocpp2.ChargingProfileCriterion{}}, true},
		{ocpp2.GetChargingProfilesRequest{}, true},
		{ocpp2.GetChargingProfilesRequest{RequestID: newInt(42), EvseID: newInt(-1), ChargingProfile: validChargingProfileCriterion}, false},
		{ocpp2.GetChargingProfilesRequest{RequestID: newInt(-1), EvseID: newInt(1), ChargingProfile: validChargingProfileCriterion}, false},
		{ocpp2.GetChargingProfilesRequest{ChargingProfile: ocpp2.ChargingProfileCriterion{ChargingProfilePurpose: "invalidChargingProfilePurpose", StackLevel: newInt(2), ChargingProfileID: []int{1,2}, ChargingLimitSource:[]ocpp2.ChargingLimitSourceType{ocpp2.ChargingLimitSourceEMS}}}, false},
		{ocpp2.GetChargingProfilesRequest{ChargingProfile: ocpp2.ChargingProfileCriterion{ChargingProfilePurpose: ocpp2.ChargingProfilePurposeTxDefaultProfile, StackLevel: newInt(-1), ChargingProfileID: []int{1,2}, ChargingLimitSource:[]ocpp2.ChargingLimitSourceType{ocpp2.ChargingLimitSourceEMS}}}, false},
		{ocpp2.GetChargingProfilesRequest{ChargingProfile: ocpp2.ChargingProfileCriterion{ChargingProfilePurpose: ocpp2.ChargingProfilePurposeTxDefaultProfile, StackLevel: newInt(2), ChargingProfileID: []int{1,2}, ChargingLimitSource:[]ocpp2.ChargingLimitSourceType{ocpp2.ChargingLimitSourceEMS, ocpp2.ChargingLimitSourceCSO, ocpp2.ChargingLimitSourceSO, ocpp2.ChargingLimitSourceOther, ocpp2.ChargingLimitSourceEMS}}}, false},
		{ocpp2.GetChargingProfilesRequest{ChargingProfile: ocpp2.ChargingProfileCriterion{ChargingProfilePurpose: ocpp2.ChargingProfilePurposeTxDefaultProfile, StackLevel: newInt(2), ChargingProfileID: []int{1,2}, ChargingLimitSource:[]ocpp2.ChargingLimitSourceType{"invalidChargingLimitSource"}}}, false},
		{ocpp2.GetChargingProfilesRequest{ChargingProfile: ocpp2.ChargingProfileCriterion{ChargingProfilePurpose: ocpp2.ChargingProfilePurposeTxDefaultProfile, StackLevel: newInt(2), ChargingProfileID: []int{-1}, ChargingLimitSource:[]ocpp2.ChargingLimitSourceType{ocpp2.ChargingLimitSourceEMS}}}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestGetChargingProfilesConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{ocpp2.GetChargingProfilesConfirmation{Status: ocpp2.GetChargingProfileStatusAccepted}, true},
		{ocpp2.GetChargingProfilesConfirmation{Status: ocpp2.GetChargingProfileStatusNoProfiles}, true},
		{ocpp2.GetChargingProfilesConfirmation{Status: "invalidGetChargingProfilesStatus"}, false},
		{ocpp2.GetChargingProfilesConfirmation{}, false},
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
	chargingProfileCriterion := ocpp2.ChargingProfileCriterion{
		ChargingProfilePurpose: ocpp2.ChargingProfilePurposeChargingStationMaxProfile,
		StackLevel:             newInt(1),
		ChargingProfileID:      []int{1,2},
		ChargingLimitSource:    []ocpp2.ChargingLimitSourceType{ocpp2.ChargingLimitSourceEMS},
	}
	status := ocpp2.GetChargingProfileStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"requestId":%v,"evseId":%v,"chargingProfile":{"chargingProfilePurpose":"%v","stackLevel":%v,"chargingProfileId":[%v,%v],"chargingLimitSource":["%v"]}}]`,
		messageId, ocpp2.GetChargingProfilesFeatureName, requestID, evseID, chargingProfileCriterion.ChargingProfilePurpose, *chargingProfileCriterion.StackLevel, chargingProfileCriterion.ChargingProfileID[0], chargingProfileCriterion.ChargingProfileID[1], chargingProfileCriterion.ChargingLimitSource[0])
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	getChargingProfilesConfirmation := ocpp2.NewGetChargingProfilesConfirmation(status)
	channel := NewMockWebSocket(wsId)

	coreListener := MockChargePointCoreListener{}
	coreListener.On("OnGetChargingProfiles", mock.Anything).Return(getChargingProfilesConfirmation, nil).Run(func(args mock.Arguments) {
		// Assert request message contents
		request, ok := args.Get(0).(*ocpp2.GetChargingProfilesRequest)
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
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, coreListener, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.GetChargingProfiles(wsId, func(confirmation *ocpp2.GetChargingProfilesConfirmation, err error) {
		// Assert confirmation message contents
		require.Nil(t, err)
		require.NotNil(t, confirmation)
		assert.Equal(t, status, confirmation.Status)
		resultChannel <- true
	},
	chargingProfileCriterion,
	func(request *ocpp2.GetChargingProfilesRequest) {
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
	chargingProfileCriterion := ocpp2.ChargingProfileCriterion{
		ChargingProfilePurpose: ocpp2.ChargingProfilePurposeChargingStationMaxProfile,
		StackLevel:             newInt(1),
		ChargingProfileID:      []int{1,2},
		ChargingLimitSource:    []ocpp2.ChargingLimitSourceType{ocpp2.ChargingLimitSourceEMS},
	}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"requestId":%v,"evseId":%v,"chargingProfile":{"chargingProfilePurpose":"%v","stackLevel":%v,"chargingProfileId":[%v,%v],"chargingLimitSource":["%v"]}}]`,
		messageId, ocpp2.GetChargingProfilesFeatureName, requestID, evseID, chargingProfileCriterion.ChargingProfilePurpose, *chargingProfileCriterion.StackLevel, chargingProfileCriterion.ChargingProfileID[0], chargingProfileCriterion.ChargingProfileID[1], chargingProfileCriterion.ChargingLimitSource[0])
	getChargingProfilesRequest := ocpp2.NewGetChargingProfilesRequest(chargingProfileCriterion)
	getChargingProfilesRequest.EvseID = &evseID
	getChargingProfilesRequest.RequestID = &requestID
	testUnsupportedRequestFromChargePoint(suite, getChargingProfilesRequest, requestJson, messageId)
}
