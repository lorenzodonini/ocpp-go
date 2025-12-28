package ocpp2_test

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/security"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

// Test
func (suite *OcppV2TestSuite) TestSecurityEventNotificationRequestValidation() {
	var requestTable = []GenericTestEntry{
		{security.SecurityEventNotificationRequest{Type: "type1", Timestamp: types.NewDateTime(time.Now()), TechInfo: "someTechInfo"}, true},
		{security.SecurityEventNotificationRequest{Type: "type1", Timestamp: types.NewDateTime(time.Now())}, true},
		{security.SecurityEventNotificationRequest{Type: "type1"}, false},
		{security.SecurityEventNotificationRequest{}, false},
		{security.SecurityEventNotificationRequest{Type: "", Timestamp: types.NewDateTime(time.Now()), TechInfo: "someTechInfo"}, false},
		{security.SecurityEventNotificationRequest{Type: ">50................................................", Timestamp: types.NewDateTime(time.Now()), TechInfo: "someTechInfo"}, false},
		{security.SecurityEventNotificationRequest{Type: "type1", Timestamp: types.NewDateTime(time.Now()), TechInfo: ">255............................................................................................................................................................................................................................................................"}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV2TestSuite) TestSecurityEventNotificationConfirmationValidation() {
	var confirmationTable = []GenericTestEntry{
		{security.SecurityEventNotificationResponse{}, true},
	}
	ExecuteGenericTestTable(suite, confirmationTable)
}

func (suite *OcppV2TestSuite) TestSecurityEventNotificationE2EMocked() {
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	typ := "type1"
	timestamp := types.NewDateTime(time.Now())
	techInfo := "someTechInfo"
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"type":"%v","timestamp":"%v","techInfo":"%v"}]`,
		messageId, security.SecurityEventNotificationFeatureName, typ, timestamp.FormatTimestamp(), techInfo)
	responseJson := fmt.Sprintf(`[3,"%v",{}]`, messageId)
	securityEventNotificationResponse := security.NewSecurityEventNotificationResponse()
	channel := NewMockWebSocket(wsId)

	handler := &MockCSMSSecurityHandler{}
	handler.On("OnSecurityEventNotification", mock.AnythingOfType("string"), mock.Anything).Return(securityEventNotificationResponse, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*security.SecurityEventNotificationRequest)
		suite.Require().True(ok)
		suite.Require().NotNil(request)
		suite.Equal(typ, request.Type)
		suite.Equal(timestamp.FormatTimestamp(), request.Timestamp.FormatTimestamp())
		suite.Equal(techInfo, request.TechInfo)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	suite.Require().Nil(err)
	response, err := suite.chargingStation.SecurityEventNotification(typ, timestamp, func(request *security.SecurityEventNotificationRequest) {
		request.TechInfo = techInfo
	})
	suite.Require().Nil(err)
	suite.Require().NotNil(response)
}

func (suite *OcppV2TestSuite) TestSecurityEventNotificationInvalidEndpoint() {
	messageId := defaultMessageId
	typ := "type1"
	timestamp := types.NewDateTime(time.Now())
	techInfo := "someTechInfo"
	request := security.NewSecurityEventNotificationRequest(typ, timestamp)
	request.TechInfo = techInfo
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"type":"%v","timestamp":"%v","techInfo":"%v"}]`,
		messageId, security.SecurityEventNotificationFeatureName, typ, timestamp.FormatTimestamp(), techInfo)
	testUnsupportedRequestFromCentralSystem(suite, request, requestJson, messageId)
}
