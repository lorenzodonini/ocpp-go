package ocpp2_test

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/availability"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

// Test
func (suite *OcppV2TestSuite) TestStatusNotificationRequestValidation() {
	var requestTable = []GenericTestEntry{
		{availability.StatusNotificationRequest{Timestamp: types.NewDateTime(time.Now()), ConnectorStatus: availability.ConnectorStatusAvailable, EvseID: 1, ConnectorID: 1}, true},
		{availability.StatusNotificationRequest{Timestamp: types.NewDateTime(time.Now()), ConnectorStatus: availability.ConnectorStatusAvailable, EvseID: 1}, true},
		{availability.StatusNotificationRequest{Timestamp: types.NewDateTime(time.Now()), ConnectorStatus: availability.ConnectorStatusAvailable}, true},
		{availability.StatusNotificationRequest{Timestamp: types.NewDateTime(time.Now())}, false},
		{availability.StatusNotificationRequest{ConnectorStatus: availability.ConnectorStatusAvailable}, false},
		{availability.StatusNotificationRequest{}, false},
		{availability.StatusNotificationRequest{Timestamp: types.NewDateTime(time.Now()), ConnectorStatus: "invalidConnectorStatus", EvseID: 1, ConnectorID: 1}, false},
		{availability.StatusNotificationRequest{Timestamp: types.NewDateTime(time.Now()), ConnectorStatus: availability.ConnectorStatusAvailable, EvseID: -1, ConnectorID: 1}, false},
		{availability.StatusNotificationRequest{Timestamp: types.NewDateTime(time.Now()), ConnectorStatus: availability.ConnectorStatusAvailable, EvseID: 1, ConnectorID: -1}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV2TestSuite) TestStatusNotificationResponseValidation() {
	var responseTable = []GenericTestEntry{
		{availability.StatusNotificationResponse{}, true},
	}
	ExecuteGenericTestTable(suite, responseTable)
}

func (suite *OcppV2TestSuite) TestStatusNotificationE2EMocked() {
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	timestamp := types.NewDateTime(time.Now())
	status := availability.ConnectorStatusAvailable
	evseID := 1
	connectorID := 1
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"timestamp":"%v","connectorStatus":"%v","evseId":%v,"connectorId":%v}]`,
		messageId, availability.StatusNotificationFeatureName, timestamp.FormatTimestamp(), status, evseID, connectorID)
	responseJson := fmt.Sprintf(`[3,"%v",{}]`,
		messageId)
	statusNotificationResponse := availability.NewStatusNotificationResponse()
	channel := NewMockWebSocket(wsId)

	handler := &MockCSMSAvailabilityHandler{}
	handler.On("OnStatusNotification", mock.AnythingOfType("string"), mock.Anything).Return(statusNotificationResponse, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*availability.StatusNotificationRequest)
		suite.Require().True(ok)
		suite.Require().NotNil(request)
		assertDateTimeEquality(suite, timestamp, request.Timestamp)
		suite.Equal(status, request.ConnectorStatus)
		suite.Equal(evseID, request.EvseID)
		suite.Equal(connectorID, request.ConnectorID)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	suite.Nil(err)
	response, err := suite.chargingStation.StatusNotification(timestamp, status, evseID, connectorID)
	suite.Nil(err)
	suite.NotNil(response)
}

func (suite *OcppV2TestSuite) TestStatusNotificationInvalidEndpoint() {
	messageId := defaultMessageId
	timestamp := types.NewDateTime(time.Now())
	status := availability.ConnectorStatusAvailable
	evseID := 1
	connectorID := 1
	request := availability.NewStatusNotificationRequest(timestamp, status, evseID, connectorID)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"timestamp":"%v","connectorStatus":"%v","evseId":%v,"connectorId":%v}]`,
		messageId, availability.StatusNotificationFeatureName, timestamp.FormatTimestamp(), status, evseID, connectorID)
	testUnsupportedRequestFromCentralSystem(suite, request, requestJson, messageId)
}
