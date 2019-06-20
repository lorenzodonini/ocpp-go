package test

import (
	"fmt"
	"github.com/lorenzodonini/go-ocpp/ocpp"
	"github.com/lorenzodonini/go-ocpp/ocpp/1.6"
	"github.com/lorenzodonini/go-ocpp/ws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

// Utility functions
func GetChangeAvailabilityRequest(t *testing.T, request ocpp.Request) *v16.ChangeAvailabilityRequest {
	assert.NotNil(t, request)
	result := request.(*v16.ChangeAvailabilityRequest)
	assert.NotNil(t, result)
	assert.IsType(t, &v16.ChangeAvailabilityRequest{}, result)
	return result
}

func GetChangeAvailabilityConfirmation(t *testing.T, confirmation ocpp.Confirmation) *v16.ChangeAvailabilityConfirmation {
	assert.NotNil(t, confirmation)
	result := confirmation.(*v16.ChangeAvailabilityConfirmation)
	assert.NotNil(t, result)
	assert.IsType(t, &v16.ChangeAvailabilityConfirmation{}, result)
	return result
}

func (suite *OcppV16TestSuite) TestChangeAvailabilityRequestValidation() {
	t := suite.T()
	var testTable = []RequestTestEntry{
		{v16.ChangeAvailabilityRequest{ConnectorId: 0, Type: v16.AvailabilityTypeOperative}, true},
		{v16.ChangeAvailabilityRequest{ConnectorId: 0, Type: v16.AvailabilityTypeInoperative}, true},
		{v16.ChangeAvailabilityRequest{ConnectorId: 0}, false},
		{v16.ChangeAvailabilityRequest{Type: v16.AvailabilityTypeOperative}, true},
		{v16.ChangeAvailabilityRequest{ConnectorId: -1, Type: v16.AvailabilityTypeOperative}, false},
	}
	executeRequestTestTable(t, testTable)
}

func (suite *OcppV16TestSuite) TestChangeAvailabilityConfirmationValidation() {
	t := suite.T()
	var testTable = []ConfirmationTestEntry{
		{v16.ChangeAvailabilityConfirmation{Status: v16.AvailabilityStatusAccepted}, true},
		{v16.ChangeAvailabilityConfirmation{Status: v16.AvailabilityStatusRejected}, true},
		{v16.ChangeAvailabilityConfirmation{Status: v16.AvailabilityStatusScheduled}, true},
		{v16.ChangeAvailabilityConfirmation{}, false},
	}
	executeConfirmationTestTable(t, testTable)
}

// Test
func (suite *OcppV16TestSuite) TestChangeAvailabilityRequestFromJson() {
	t := suite.T()
	uniqueId := "1234"
	connectorId := 1
	availabilityType := v16.AvailabilityTypeOperative
	dataJson := fmt.Sprintf(`[2,"%v","ChangeAvailability",{"connectorId":%v,"type":"%v"}]`, uniqueId, connectorId, availabilityType)
	call := ParseCall(&suite.centralSystem.Endpoint, dataJson, t)
	CheckCall(call, t, v16.ChangeAvailabilityFeatureName, uniqueId)
	request := GetChangeAvailabilityRequest(t, call.Payload)
	assert.Equal(t, connectorId, request.ConnectorId)
	assert.Equal(t, availabilityType, request.Type)
}

func (suite *OcppV16TestSuite) TestChangeAvailabilityRequestToJson() {
	t := suite.T()
	connectorId := 1
	availabilityType := v16.AvailabilityTypeOperative
	request := v16.ChangeAvailabilityRequest{ConnectorId: connectorId, Type: availabilityType}
	call, err := suite.chargePoint.CreateCall(request)
	assert.Nil(t, err)
	uniqueId := call.GetUniqueId()
	assert.NotNil(t, call)
	err = validate.Struct(call)
	assert.Nil(t, err)
	jsonData, err := call.MarshalJSON()
	assert.Nil(t, err)
	assert.NotNil(t, jsonData)
	expectedJson := fmt.Sprintf(`[2,"%v","ChangeAvailability",{"connectorId":%v,"type":"%v"}]`, uniqueId, connectorId, availabilityType)
	assert.Equal(t, []byte(expectedJson), jsonData)
}

func (suite *OcppV16TestSuite) TestChangeAvailabilityConfirmationFromJson() {
	t := suite.T()
	uniqueId := "5678"
	status := v16.AvailabilityStatusAccepted
	dummyRequest := v16.ChangeAvailabilityRequest{}
	dataJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, uniqueId, status)
	suite.chargePoint.Endpoint.AddPendingRequest(uniqueId, dummyRequest)
	callResult := ParseCallResult(&suite.chargePoint.Endpoint, dataJson, t)
	CheckCallResult(callResult, t, uniqueId)
	confirmation := GetChangeAvailabilityConfirmation(t, callResult.Payload)
	assert.Equal(t, status, confirmation.Status)
}

func (suite *OcppV16TestSuite) TestChangeAvailabilityConfirmationToJson() {
	t := suite.T()
	uniqueId := "1234"
	status := v16.AvailabilityStatusAccepted
	confirmation := v16.ChangeAvailabilityConfirmation{Status: status}
	callResult, err := suite.centralSystem.CreateCallResult(confirmation, uniqueId)
	assert.Nil(t, err)
	assert.NotNil(t, callResult)
	err = validate.Struct(callResult)
	assert.Nil(t, err)
	jsonData, err := callResult.MarshalJSON()
	assert.Nil(t, err)
	assert.NotNil(t, jsonData)
	expectedJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, uniqueId, status)
	assert.Equal(t, []byte(expectedJson), jsonData)
}

func (suite *OcppV16TestSuite) TestChangeAvailabilityE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := "1234"
	wsUrl := "someUrl"
	connectorId := 1
	availabilityType := v16.AvailabilityTypeOperative
	status := v16.AvailabilityStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"connectorId":%v,"type":"%v"}]`, messageId, v16.ChangeAvailabilityFeatureName, connectorId, availabilityType)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	requestRaw := []byte(requestJson)
	responseRaw := []byte(responseJson)
	channel := MockWebSocket{id: wsId}
	// Setting server handlers
	suite.mockServer.SetNewClientHandler(func(ws ws.Channel) {
		assert.NotNil(t, ws)
		assert.Equal(t, wsId, ws.GetId())
	})
	suite.mockServer.SetMessageHandler(func(ws ws.Channel, data []byte) error {
		assert.Equal(t, requestRaw, data)
		jsonData := string(data)
		assert.Equal(t, requestJson, jsonData)
		call := ParseCall(&suite.chargePoint.Endpoint, jsonData, t)
		CheckCall(call, t, v16.ChangeAvailabilityFeatureName, messageId)
		suite.chargePoint.AddPendingRequest(messageId, call.Payload)
		// TODO: generate the response dynamically
		err := suite.mockClient.messageHandler(responseRaw)
		assert.Nil(t, err)
		return nil
	})
	// Setting client handlers
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return().Run(func(args mock.Arguments) {
		u := args.String(0)
		assert.Equal(t, wsUrl, u)
		suite.mockServer.newClientHandler(channel)
	})
	suite.mockClient.SetMessageHandler(func(data []byte) error {
		assert.Equal(t, responseRaw, data)
		jsonData := string(data)
		assert.Equal(t, responseJson, jsonData)
		callResult := ParseCallResult(&suite.chargePoint.Endpoint, jsonData, t)
		CheckCallResult(callResult, t, messageId)
		return nil
	})
	suite.mockClient.On("Write", mock.Anything).Return().Run(func(args mock.Arguments) {
		data := args.Get(0)
		bytes := data.([]byte)
		assert.NotNil(t, bytes)
		err := suite.mockServer.messageHandler(channel, bytes)
		assert.Nil(t, err)
	})
	// Test Run
	err := suite.mockClient.Start(wsUrl)
	assert.Nil(t, err)
	suite.mockClient.Write(requestRaw)
}
