package ocpp16_test

import (
	"fmt"
	"github.com/lorenzodonini/go-ocpp/ocpp1.6"
	"github.com/lorenzodonini/go-ocpp/ocppj"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

// Utility functions
func getChangeAvailabilityRequest(t *testing.T, request ocppj.Request) *ocpp16.ChangeAvailabilityRequest {
	assert.NotNil(t, request)
	result := request.(*ocpp16.ChangeAvailabilityRequest)
	assert.NotNil(t, result)
	assert.IsType(t, &ocpp16.ChangeAvailabilityRequest{}, result)
	return result
}

func getChangeAvailabilityConfirmation(t *testing.T, confirmation ocppj.Confirmation) *ocpp16.ChangeAvailabilityConfirmation {
	assert.NotNil(t, confirmation)
	result := confirmation.(*ocpp16.ChangeAvailabilityConfirmation)
	assert.NotNil(t, result)
	assert.IsType(t, &ocpp16.ChangeAvailabilityConfirmation{}, result)
	return result
}

func (suite *OcppV16TestSuite) TestChangeAvailabilityRequestValidation() {
	t := suite.T()
	var testTable = []RequestTestEntry{
		{ocpp16.ChangeAvailabilityRequest{ConnectorId: 0, Type: ocpp16.AvailabilityTypeOperative}, true},
		{ocpp16.ChangeAvailabilityRequest{ConnectorId: 0, Type: ocpp16.AvailabilityTypeInoperative}, true},
		{ocpp16.ChangeAvailabilityRequest{ConnectorId: 0}, false},
		{ocpp16.ChangeAvailabilityRequest{Type: ocpp16.AvailabilityTypeOperative}, true},
		{ocpp16.ChangeAvailabilityRequest{ConnectorId: -1, Type: ocpp16.AvailabilityTypeOperative}, false},
	}
	ExecuteRequestTestTable(t, testTable)
}

func (suite *OcppV16TestSuite) TestChangeAvailabilityConfirmationValidation() {
	t := suite.T()
	var testTable = []ConfirmationTestEntry{
		{ocpp16.ChangeAvailabilityConfirmation{Status: ocpp16.AvailabilityStatusAccepted}, true},
		{ocpp16.ChangeAvailabilityConfirmation{Status: ocpp16.AvailabilityStatusRejected}, true},
		{ocpp16.ChangeAvailabilityConfirmation{Status: ocpp16.AvailabilityStatusScheduled}, true},
		{ocpp16.ChangeAvailabilityConfirmation{}, false},
	}
	ExecuteConfirmationTestTable(t, testTable)
}

// Test
func (suite *OcppV16TestSuite) TestChangeAvailabilityE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	connectorId := 1
	availabilityType := ocpp16.AvailabilityTypeOperative
	status := ocpp16.AvailabilityStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"connectorId":%v,"type":"%v"}]`, messageId, ocpp16.ChangeAvailabilityFeatureName, connectorId, availabilityType)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	changeAvailabilityConfirmation := ocpp16.NewChangeAvailabilityConfirmation(status)
	channel := NewMockWebSocket(wsId)
	// Setting handlers
	coreListener := MockChargePointCoreListener{}
	coreListener.On("OnChangeAvailability", mock.Anything).Return(changeAvailabilityConfirmation, nil)
	setupDefaultChargePointHandlers(suite, coreListener, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	assert.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.centralSystem.ChangeAvailability(wsId, func(confirmation *ocpp16.ChangeAvailabilityConfirmation, callError *ocppj.ProtoError) {
		assert.NotNil(t, confirmation)
		assert.Nil(t, callError)
		assert.Equal(t, status, confirmation.Status)
		resultChannel <- true
	}, connectorId, availabilityType)
	assert.Nil(t, err)
	result := <- resultChannel
	assert.True(t, result)
}

func (suite *OcppV16TestSuite) TestChangeAvailabilityInvalidEndpoint() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	connectorId := 1
	availabilityType := ocpp16.AvailabilityTypeOperative
	expectedError := fmt.Sprintf("unsupported action %v on charge point, cannot send request", ocpp16.ChangeAvailabilityFeatureName)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"connectorId":%v,"type":"%v"}]`, messageId, ocpp16.ChangeAvailabilityFeatureName, connectorId, availabilityType)

	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: false, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: false})
	// Run test
	err := suite.chargePoint.Start(wsUrl)
	assert.Nil(t, err)
	changeAvailabilityRequest := ocpp16.NewChangeAvailabilityRequest(connectorId, availabilityType)
	err = suite.chargePoint.SendRequestAsync(changeAvailabilityRequest, func(confirmation ocppj.Confirmation, callError *ocppj.ProtoError) {
		t.Fail()
	})
	assert.Error(t, err)
	assert.Equal(t, expectedError, err.Error())
}

func (suite *OcppV16TestSuite) TestChangeAvailabilityInvalidEndpointResponse() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	connectorId := 1
	availabilityType := ocpp16.AvailabilityTypeOperative
	errorDescription := fmt.Sprintf("unsupported action %v on central system", ocpp16.ChangeAvailabilityFeatureName)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"connectorId":%v,"type":"%v"}]`, messageId, ocpp16.ChangeAvailabilityFeatureName, connectorId, availabilityType)
	errorJson := fmt.Sprintf(`[4,"%v","%v","%v",null]`, messageId, ocppj.NotSupported, errorDescription)
	channel := NewMockWebSocket(wsId)

	coreListener := MockCentralSystemCoreListener{}
	setupDefaultCentralSystemHandlers(suite, coreListener, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(errorJson), forwardWrittenMessage: true})
	suite.ocppjCentralSystem.SetErrorHandler(func(chargePointId string, errorCode ocppj.ErrorCode, description string, details interface{}, requestId string) {
		assert.Equal(t, messageId, requestId)
		assert.Equal(t, wsId, chargePointId)
		assert.Equal(t, ocppj.NotSupported, errorCode)
		assert.Equal(t, errorDescription, description)
		assert.Nil(t, details)
	})
	// Mock pending request
	pendingRequest := ocpp16.NewChangeAvailabilityRequest(connectorId, availabilityType)
	suite.ocppjChargePoint.AddPendingRequest(messageId, pendingRequest)
	// Run test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	assert.Nil(t, err)
	err = suite.mockWsServer.MessageHandler(channel, []byte(requestJson))
	assert.Nil(t, err)
}
