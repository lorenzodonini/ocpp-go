package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func (suite *OcppV2TestSuite) TestChangeAvailabilityRequestValidation() {
	t := suite.T()
	var testTable = []GenericTestEntry{
		{ocpp2.ChangeAvailabilityRequest{EvseID: 0, OperationalStatus: ocpp2.OperationalStatusOperative}, true},
		{ocpp2.ChangeAvailabilityRequest{EvseID: 0, OperationalStatus: ocpp2.OperationalStatusInoperative}, true},
		{ocpp2.ChangeAvailabilityRequest{EvseID: 0}, false},
		{ocpp2.ChangeAvailabilityRequest{OperationalStatus: ocpp2.OperationalStatusOperative}, true},
		{ocpp2.ChangeAvailabilityRequest{OperationalStatus: "invalidAvailabilityType"}, false},
		{ocpp2.ChangeAvailabilityRequest{EvseID: -1, OperationalStatus: ocpp2.OperationalStatusOperative}, false},
	}
	ExecuteGenericTestTable(t, testTable)
}

func (suite *OcppV2TestSuite) TestChangeAvailabilityConfirmationValidation() {
	t := suite.T()
	var testTable = []GenericTestEntry{
		{ocpp2.ChangeAvailabilityConfirmation{Status: ocpp2.ChangeAvailabilityStatusAccepted}, true},
		{ocpp2.ChangeAvailabilityConfirmation{Status: ocpp2.ChangeAvailabilityStatusRejected}, true},
		{ocpp2.ChangeAvailabilityConfirmation{Status: ocpp2.ChangeAvailabilityStatusScheduled}, true},
		{ocpp2.ChangeAvailabilityConfirmation{Status: "invalidAvailabilityStatus"}, false},
		{ocpp2.ChangeAvailabilityConfirmation{}, false},
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
	operationalStatus := ocpp2.OperationalStatusOperative
	status := ocpp2.ChangeAvailabilityStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"evseId":%v,"operationalStatus":"%v"}]`, messageId, ocpp2.ChangeAvailabilityFeatureName, evseID, operationalStatus)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	changeAvailabilityConfirmation := ocpp2.NewChangeAvailabilityConfirmation(status)
	channel := NewMockWebSocket(wsId)
	// Setting handlers
	coreListener := MockChargePointCoreListener{}
	coreListener.On("OnChangeAvailability", mock.Anything).Return(changeAvailabilityConfirmation, nil)
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, coreListener, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	assert.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.ChangeAvailability(wsId, func(confirmation *ocpp2.ChangeAvailabilityConfirmation, err error) {
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
	operationalStatus := ocpp2.OperationalStatusOperative
	changeAvailabilityRequest := ocpp2.NewChangeAvailabilityRequest(evseID, operationalStatus)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"evseId":%v,"operationalStatus":"%v"}]`, messageId, ocpp2.ChangeAvailabilityFeatureName, evseID, operationalStatus)
	testUnsupportedRequestFromChargePoint(suite, changeAvailabilityRequest, requestJson, messageId)
}
