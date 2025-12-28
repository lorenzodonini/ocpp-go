package ocpp16_test

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/stretchr/testify/mock"
)

func (suite *OcppV16TestSuite) TestChangeAvailabilityRequestValidation() {
	var testTable = []GenericTestEntry{
		{core.ChangeAvailabilityRequest{ConnectorId: 0, Type: core.AvailabilityTypeOperative}, true},
		{core.ChangeAvailabilityRequest{ConnectorId: 0, Type: core.AvailabilityTypeInoperative}, true},
		{core.ChangeAvailabilityRequest{ConnectorId: 0}, false},
		{core.ChangeAvailabilityRequest{Type: core.AvailabilityTypeOperative}, true},
		{core.ChangeAvailabilityRequest{Type: "invalidAvailabilityType"}, false},
		{core.ChangeAvailabilityRequest{ConnectorId: -1, Type: core.AvailabilityTypeOperative}, false},
	}
	ExecuteGenericTestTable(suite, testTable)
}

func (suite *OcppV16TestSuite) TestChangeAvailabilityConfirmationValidation() {
	var testTable = []GenericTestEntry{
		{core.ChangeAvailabilityConfirmation{Status: core.AvailabilityStatusAccepted}, true},
		{core.ChangeAvailabilityConfirmation{Status: core.AvailabilityStatusRejected}, true},
		{core.ChangeAvailabilityConfirmation{Status: core.AvailabilityStatusScheduled}, true},
		{core.ChangeAvailabilityConfirmation{Status: "invalidAvailabilityStatus"}, false},
		{core.ChangeAvailabilityConfirmation{}, false},
	}
	ExecuteGenericTestTable(suite, testTable)
}

// Test
func (suite *OcppV16TestSuite) TestChangeAvailabilityE2EMocked() {
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	connectorId := 1
	availabilityType := core.AvailabilityTypeOperative
	status := core.AvailabilityStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"connectorId":%v,"type":"%v"}]`, messageId, core.ChangeAvailabilityFeatureName, connectorId, availabilityType)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	changeAvailabilityConfirmation := core.NewChangeAvailabilityConfirmation(status)
	channel := NewMockWebSocket(wsId)
	// Setting handlers
	coreListener := &MockChargePointCoreListener{}
	coreListener.On("OnChangeAvailability", mock.Anything).Return(changeAvailabilityConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*core.ChangeAvailabilityRequest)
		suite.Require().NotNil(request)
		suite.Require().True(ok)
		suite.Equal(connectorId, request.ConnectorId)
		suite.Equal(availabilityType, request.Type)
	})
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, coreListener, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	suite.Require().Nil(err)
	resultChannel := make(chan bool, 1)
	err = suite.centralSystem.ChangeAvailability(wsId, func(confirmation *core.ChangeAvailabilityConfirmation, err error) {
		suite.Require().NotNil(confirmation)
		suite.Require().Nil(err)
		suite.Equal(status, confirmation.Status)
		resultChannel <- true
	}, connectorId, availabilityType)
	suite.Require().Nil(err)
	result := <-resultChannel
	suite.True(result)
}

func (suite *OcppV16TestSuite) TestChangeAvailabilityInvalidEndpoint() {
	messageId := defaultMessageId
	connectorId := 1
	availabilityType := core.AvailabilityTypeOperative
	changeAvailabilityRequest := core.NewChangeAvailabilityRequest(connectorId, availabilityType)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"connectorId":%v,"type":"%v"}]`, messageId, core.ChangeAvailabilityFeatureName, connectorId, availabilityType)
	testUnsupportedRequestFromChargePoint(suite, changeAvailabilityRequest, requestJson, messageId)
}
