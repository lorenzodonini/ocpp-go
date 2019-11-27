package ocpp16_test

import (
	"fmt"
	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test
func (suite *OcppV16TestSuite) TestClearChargingProfileRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{ocpp16.ClearChargingProfileRequest{Id: 1, ConnectorId: 1, ChargingProfilePurpose: ocpp16.ChargingProfilePurposeChargePointMaxProfile, StackLevel: 1}, true},
		{ocpp16.ClearChargingProfileRequest{Id: 1, ConnectorId: 1, ChargingProfilePurpose: ocpp16.ChargingProfilePurposeChargePointMaxProfile}, true},
		{ocpp16.ClearChargingProfileRequest{Id: 1, ConnectorId: 1}, true},
		{ocpp16.ClearChargingProfileRequest{Id: 1}, true},
		{ocpp16.ClearChargingProfileRequest{}, true},
		{ocpp16.ClearChargingProfileRequest{ConnectorId: -1}, false},
		{ocpp16.ClearChargingProfileRequest{Id: -1}, false},
		{ocpp16.ClearChargingProfileRequest{ChargingProfilePurpose: "invalidChargingProfilePurposeType"}, false},
		{ocpp16.ClearChargingProfileRequest{StackLevel: -1}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV16TestSuite) TestClearChargingProfileConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{ocpp16.ClearChargingProfileConfirmation{Status: ocpp16.ClearChargingProfileStatusAccepted}, true},
		{ocpp16.ClearChargingProfileConfirmation{Status: "invalidClearChargingProfileStatus"}, false},
		{ocpp16.ClearChargingProfileConfirmation{}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV16TestSuite) TestClearChargingProfileE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	chargingProfileId := 1
	connectorId := 1
	chargingProfilePurpose := ocpp16.ChargingProfilePurposeChargePointMaxProfile
	stackLevel := 1
	status := ocpp16.ClearChargingProfileStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"id":%v,"connectorId":%v,"chargingProfilePurpose":"%v","stackLevel":%v}]`,
		messageId, ocpp16.ClearChargingProfileFeatureName, chargingProfileId, connectorId, chargingProfilePurpose, stackLevel)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	ClearChargingProfileConfirmation := ocpp16.NewClearChargingProfileConfirmation(status)
	channel := NewMockWebSocket(wsId)

	smartChargingListener := MockChargePointSmartChargingListener{}
	smartChargingListener.On("OnClearChargingProfile", mock.Anything).Return(ClearChargingProfileConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*ocpp16.ClearChargingProfileRequest)
		assert.True(t, ok)
		assert.NotNil(t, request)
		assert.Equal(t, chargingProfileId, request.Id)
		assert.Equal(t, connectorId, request.ConnectorId)
		assert.Equal(t, chargingProfilePurpose, request.ChargingProfilePurpose)
		assert.Equal(t, stackLevel, request.StackLevel)
	})
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	suite.chargePoint.SetSmartChargingListener(smartChargingListener)
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	assert.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.centralSystem.ClearChargingProfile(wsId, func(confirmation *ocpp16.ClearChargingProfileConfirmation, err error) {
		if !assert.Nil(t, err) || !assert.NotNil(t, confirmation) {
			resultChannel <- false
			return
		}
		assert.Equal(t, status, confirmation.Status)
		resultChannel <- true
	}, func(request *ocpp16.ClearChargingProfileRequest) {
		request.Id = chargingProfileId
		request.ConnectorId = connectorId
		request.ChargingProfilePurpose = chargingProfilePurpose
		request.StackLevel = stackLevel
	})
	assert.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV16TestSuite) TestClearChargingProfileInvalidEndpoint() {
	messageId := defaultMessageId
	connectorId := 1
	chargingProfileId := 1
	stackLevel := 1
	chargingProfilePurpose := ocpp16.ChargingProfilePurposeChargePointMaxProfile
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"id":%v,"connectorId":%v,"chargingProfilePurpose":"%v","stackLevel":%v}]`,
		messageId, ocpp16.ClearChargingProfileFeatureName, chargingProfileId, connectorId, chargingProfilePurpose, stackLevel)
	clearChargingProfileRequest := ocpp16.NewClearChargingProfileRequest()
	clearChargingProfileRequest.Id = chargingProfileId
	clearChargingProfileRequest.ConnectorId = connectorId
	clearChargingProfileRequest.ChargingProfilePurpose = chargingProfilePurpose
	clearChargingProfileRequest.StackLevel = stackLevel
	testUnsupportedRequestFromChargePoint(suite, clearChargingProfileRequest, requestJson, messageId)
}
