package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newInt(i int) *int{
	return &i
}

// Test
func (suite *OcppV2TestSuite) TestClearChargingProfileRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{ocpp2.ClearChargingProfileRequest{EvseID: newInt(1), ChargingProfile: &ocpp2.ClearChargingProfileType{ID: 1, ChargingProfilePurpose: ocpp2.ChargingProfilePurposeChargingStationMaxProfile, StackLevel: 1}}, true},
		{ocpp2.ClearChargingProfileRequest{EvseID: newInt(1), ChargingProfile: &ocpp2.ClearChargingProfileType{ID: 1, ChargingProfilePurpose: ocpp2.ChargingProfilePurposeChargingStationMaxProfile}}, true},
		{ocpp2.ClearChargingProfileRequest{EvseID: newInt(1), ChargingProfile: &ocpp2.ClearChargingProfileType{ID: 1}}, true},
		{ocpp2.ClearChargingProfileRequest{ChargingProfile: &ocpp2.ClearChargingProfileType{ID: 1}}, true},
		{ocpp2.ClearChargingProfileRequest{ChargingProfile: &ocpp2.ClearChargingProfileType{}}, true},
		{ocpp2.ClearChargingProfileRequest{}, true},
		{ocpp2.ClearChargingProfileRequest{EvseID: newInt(-1)}, false},
		{ocpp2.ClearChargingProfileRequest{ChargingProfile: &ocpp2.ClearChargingProfileType{ID: -1}}, false},
		{ocpp2.ClearChargingProfileRequest{ChargingProfile: &ocpp2.ClearChargingProfileType{ChargingProfilePurpose: "invalidChargingProfilePurposeType"}}, false},
		{ocpp2.ClearChargingProfileRequest{ChargingProfile: &ocpp2.ClearChargingProfileType{StackLevel: -1}}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestClearChargingProfileConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{ocpp2.ClearChargingProfileConfirmation{Status: ocpp2.ClearChargingProfileStatusAccepted}, true},
		{ocpp2.ClearChargingProfileConfirmation{Status: "invalidClearChargingProfileStatus"}, false},
		{ocpp2.ClearChargingProfileConfirmation{}, false},
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
	chargingProfilePurpose := ocpp2.ChargingProfilePurposeChargingStationMaxProfile
	stackLevel := 1
	status := ocpp2.ClearChargingProfileStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"evseId":%v,"chargingProfile":{"id":%v,"chargingProfilePurpose":"%v","stackLevel":%v}}]`,
		messageId, ocpp2.ClearChargingProfileFeatureName, evseID, chargingProfileId, chargingProfilePurpose, stackLevel)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	ClearChargingProfileConfirmation := ocpp2.NewClearChargingProfileConfirmation(status)
	channel := NewMockWebSocket(wsId)

	coreListener := MockChargePointCoreListener{}
	coreListener.On("OnClearChargingProfile", mock.Anything).Return(ClearChargingProfileConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*ocpp2.ClearChargingProfileRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, evseID, *request.EvseID)
		assert.Equal(t, chargingProfileId, request.ChargingProfile.ID)
		assert.Equal(t, chargingProfilePurpose, request.ChargingProfile.ChargingProfilePurpose)
		assert.Equal(t, stackLevel, request.ChargingProfile.StackLevel)
	})
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, coreListener, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.ClearChargingProfile(wsId, func(confirmation *ocpp2.ClearChargingProfileConfirmation, err error) {
		require.Nil(t, err)
		require.NotNil(t, confirmation)
		assert.Equal(t, status, confirmation.Status)
		resultChannel <- true
	}, func(request *ocpp2.ClearChargingProfileRequest) {
		request.EvseID = &evseID
		request.ChargingProfile = &ocpp2.ClearChargingProfileType{ID: chargingProfileId, ChargingProfilePurpose: chargingProfilePurpose, StackLevel: stackLevel}
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
	chargingProfilePurpose := ocpp2.ChargingProfilePurposeChargingStationMaxProfile
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"evseId":%v,"chargingProfile":{"id":%v,"chargingProfilePurpose":"%v","stackLevel":%v}}]`,
		messageId, ocpp2.ClearChargingProfileFeatureName, evseID, chargingProfileId, chargingProfilePurpose, stackLevel)
	clearChargingProfileRequest := ocpp2.NewClearChargingProfileRequest()
	clearChargingProfileRequest.EvseID = &evseID
	clearChargingProfileRequest.ChargingProfile = &ocpp2.ClearChargingProfileType{ID: chargingProfileId, ChargingProfilePurpose: chargingProfilePurpose, StackLevel: stackLevel}
	testUnsupportedRequestFromChargePoint(suite, clearChargingProfileRequest, requestJson, messageId)
}
