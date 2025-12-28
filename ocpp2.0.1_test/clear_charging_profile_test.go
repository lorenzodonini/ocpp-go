package ocpp2_test

import (
	"fmt"

	"github.com/stretchr/testify/mock"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/smartcharging"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

// Test
func (suite *OcppV2TestSuite) TestClearChargingProfileRequestValidation() {
	var requestTable = []GenericTestEntry{
		{smartcharging.ClearChargingProfileRequest{ChargingProfileID: newInt(1), ChargingProfileCriteria: &smartcharging.ClearChargingProfileType{EvseID: newInt(1), ChargingProfilePurpose: types.ChargingProfilePurposeChargingStationMaxProfile, StackLevel: newInt(1)}}, true},
		{smartcharging.ClearChargingProfileRequest{ChargingProfileID: newInt(1), ChargingProfileCriteria: &smartcharging.ClearChargingProfileType{EvseID: newInt(1), ChargingProfilePurpose: types.ChargingProfilePurposeChargingStationMaxProfile}}, true},
		{smartcharging.ClearChargingProfileRequest{ChargingProfileID: newInt(1), ChargingProfileCriteria: &smartcharging.ClearChargingProfileType{EvseID: newInt(1)}}, true},
		{smartcharging.ClearChargingProfileRequest{ChargingProfileCriteria: &smartcharging.ClearChargingProfileType{EvseID: newInt(1)}}, true},
		{smartcharging.ClearChargingProfileRequest{ChargingProfileCriteria: &smartcharging.ClearChargingProfileType{}}, true},
		{smartcharging.ClearChargingProfileRequest{}, true},
		{smartcharging.ClearChargingProfileRequest{ChargingProfileCriteria: &smartcharging.ClearChargingProfileType{EvseID: newInt(-1)}}, false},
		{smartcharging.ClearChargingProfileRequest{ChargingProfileCriteria: &smartcharging.ClearChargingProfileType{ChargingProfilePurpose: "invalidChargingProfilePurposeType"}}, false},
		{smartcharging.ClearChargingProfileRequest{ChargingProfileCriteria: &smartcharging.ClearChargingProfileType{StackLevel: newInt(-1)}}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV2TestSuite) TestClearChargingProfileConfirmationValidation() {
	var confirmationTable = []GenericTestEntry{
		{smartcharging.ClearChargingProfileResponse{Status: smartcharging.ClearChargingProfileStatusAccepted, StatusInfo: types.NewStatusInfo("200", "ok")}, true},
		{smartcharging.ClearChargingProfileResponse{Status: smartcharging.ClearChargingProfileStatusAccepted}, true},
		{smartcharging.ClearChargingProfileResponse{Status: "invalidClearChargingProfileStatus"}, false},
		{smartcharging.ClearChargingProfileResponse{}, false},
		{smartcharging.ClearChargingProfileResponse{StatusInfo: types.NewStatusInfo("", "")}, false},
	}
	ExecuteGenericTestTable(suite, confirmationTable)
}

func (suite *OcppV2TestSuite) TestClearChargingProfileE2EMocked() {
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	chargingProfileId := newInt(1)
	chargingProfileCriteria := smartcharging.ClearChargingProfileType{
		EvseID:                 newInt(1),
		ChargingProfilePurpose: types.ChargingProfilePurposeChargingStationMaxProfile,
		StackLevel:             newInt(1),
	}
	status := smartcharging.ClearChargingProfileStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"chargingProfileId":%v,"chargingProfileCriteria":{"evseId":%v,"chargingProfilePurpose":"%v","stackLevel":%v}}]`,
		messageId, smartcharging.ClearChargingProfileFeatureName, *chargingProfileId, *chargingProfileCriteria.EvseID, chargingProfileCriteria.ChargingProfilePurpose, *chargingProfileCriteria.StackLevel)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	clearChargingProfileConfirmation := smartcharging.NewClearChargingProfileResponse(status)
	channel := NewMockWebSocket(wsId)

	handler := &MockChargingStationSmartChargingHandler{}
	handler.On("OnClearChargingProfile", mock.Anything).Return(clearChargingProfileConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*smartcharging.ClearChargingProfileRequest)
		suite.Require().True(ok)
		suite.Require().NotNil(request)
		suite.Equal(*chargingProfileId, *request.ChargingProfileID)
		suite.Equal(*chargingProfileCriteria.EvseID, *request.ChargingProfileCriteria.EvseID)
		suite.Equal(chargingProfileCriteria.ChargingProfilePurpose, request.ChargingProfileCriteria.ChargingProfilePurpose)
		suite.Equal(*chargingProfileCriteria.StackLevel, *request.ChargingProfileCriteria.StackLevel)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	suite.Require().Nil(err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.ClearChargingProfile(wsId, func(confirmation *smartcharging.ClearChargingProfileResponse, err error) {
		suite.Require().Nil(err)
		suite.Require().NotNil(confirmation)
		suite.Equal(status, confirmation.Status)
		resultChannel <- true
	}, func(request *smartcharging.ClearChargingProfileRequest) {
		request.ChargingProfileID = chargingProfileId
		request.ChargingProfileCriteria = &chargingProfileCriteria
	})
	suite.Require().Nil(err)
	result := <-resultChannel
	suite.True(result)
}

func (suite *OcppV2TestSuite) TestClearChargingProfileInvalidEndpoint() {
	messageId := defaultMessageId
	chargingProfileId := 1
	chargingProfileCriteria := smartcharging.ClearChargingProfileType{
		EvseID:                 newInt(1),
		ChargingProfilePurpose: types.ChargingProfilePurposeChargingStationMaxProfile,
		StackLevel:             newInt(1),
	}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"chargingProfileId":%v,"chargingProfileCriteria":{"evseId":%v,"chargingProfilePurpose":"%v","stackLevel":%v}}]`,
		messageId, smartcharging.ClearChargingProfileFeatureName, chargingProfileId, *chargingProfileCriteria.EvseID, chargingProfileCriteria.ChargingProfilePurpose, *chargingProfileCriteria.StackLevel)
	clearChargingProfileRequest := smartcharging.NewClearChargingProfileRequest()
	clearChargingProfileRequest.ChargingProfileID = &chargingProfileId
	clearChargingProfileRequest.ChargingProfileCriteria = &chargingProfileCriteria
	testUnsupportedRequestFromChargingStation(suite, clearChargingProfileRequest, requestJson, messageId)
}
