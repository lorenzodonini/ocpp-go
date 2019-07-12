package ocpp16_test

import (
	"fmt"
	"github.com/lorenzodonini/go-ocpp/ocpp1.6"
	"github.com/lorenzodonini/go-ocpp/ocppj"
	"github.com/lorenzodonini/go-ocpp/ws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

// Utility functions
func GetAuthorizeRequest(t *testing.T, request ocppj.Request) *ocpp16.AuthorizeRequest {
	assert.NotNil(t, request)
	result := request.(*ocpp16.AuthorizeRequest)
	assert.NotNil(t, result)
	assert.IsType(t, &ocpp16.AuthorizeRequest{}, result)
	return result
}

func GetAuthorizeConfirmation(t *testing.T, confirmation ocppj.Confirmation) *ocpp16.AuthorizeConfirmation {
	assert.NotNil(t, confirmation)
	result := confirmation.(*ocpp16.AuthorizeConfirmation)
	assert.NotNil(t, result)
	assert.IsType(t, &ocpp16.AuthorizeConfirmation{}, result)
	return result
}

// Test
func (suite *OcppV16TestSuite) TestAuthorizeRequestValidation() {
	t := suite.T()
	var requestTable = []RequestTestEntry{
		{ocpp16.AuthorizeRequest{IdTag: "12345"}, true},
		{ocpp16.AuthorizeRequest{}, false},
		{ocpp16.AuthorizeRequest{IdTag: ">20.................."}, false},
	}
	ExecuteRequestTestTable(t, requestTable)
}

func (suite *OcppV16TestSuite) TestAuthorizeConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []ConfirmationTestEntry{
		{ocpp16.AuthorizeConfirmation{IdTagInfo: ocpp16.IdTagInfo{ExpiryDate: time.Now().Add(time.Hour * 8), ParentIdTag: "00000", Status: ocpp16.AuthorizationStatusAccepted}}, true},
		{ocpp16.AuthorizeConfirmation{IdTagInfo: ocpp16.IdTagInfo{ParentIdTag: "00000", Status: ocpp16.AuthorizationStatusAccepted}}, true},
		{ocpp16.AuthorizeConfirmation{IdTagInfo: ocpp16.IdTagInfo{ExpiryDate: time.Now().Add(time.Hour * 8), Status: ocpp16.AuthorizationStatusAccepted}}, true},
		{ocpp16.AuthorizeConfirmation{IdTagInfo: ocpp16.IdTagInfo{Status: ocpp16.AuthorizationStatusAccepted}}, true},
		{ocpp16.AuthorizeConfirmation{IdTagInfo: ocpp16.IdTagInfo{Status: ocpp16.AuthorizationStatusBlocked}}, true},
		{ocpp16.AuthorizeConfirmation{IdTagInfo: ocpp16.IdTagInfo{Status: ocpp16.AuthorizationStatusExpired}}, true},
		{ocpp16.AuthorizeConfirmation{IdTagInfo: ocpp16.IdTagInfo{Status: ocpp16.AuthorizationStatusInvalid}}, true},
		{ocpp16.AuthorizeConfirmation{IdTagInfo: ocpp16.IdTagInfo{Status: ocpp16.AuthorizationStatusConcurrentTx}}, true},
		{ocpp16.AuthorizeConfirmation{IdTagInfo: ocpp16.IdTagInfo{ParentIdTag: ">20..................", Status: ocpp16.AuthorizationStatusAccepted}}, false},
		{ocpp16.AuthorizeConfirmation{IdTagInfo: ocpp16.IdTagInfo{ExpiryDate: time.Now().Add(time.Hour * -8), Status: ocpp16.AuthorizationStatusAccepted}}, false},
	}
	ExecuteConfirmationTestTable(t, confirmationTable)
}

func (suite *OcppV16TestSuite) TestAuthorizeRequestFromJson() {
	t := suite.T()
	uniqueId := "1234"
	idTag := "tag1"
	dataJson := fmt.Sprintf(`[2,"%v","Authorize",{"idTag":"%v"}]`, uniqueId, idTag)
	call := ParseCall(&suite.centralSystem.Endpoint, dataJson, t)
	CheckCall(call, t, ocpp16.AuthorizeFeatureName, uniqueId)
	request := GetAuthorizeRequest(t, call.Payload)
	assert.Equal(t, idTag, request.IdTag)
}

func (suite *OcppV16TestSuite) TestAuthorizeRequestToJson() {
	t := suite.T()
	idTag := "tag2"
	request := ocpp16.AuthorizeRequest{IdTag: idTag}
	call, err := suite.chargePoint.CreateCall(request)
	assert.Nil(t, err)
	uniqueId := call.GetUniqueId()
	assert.NotNil(t, call)
	err = Validate.Struct(call)
	assert.Nil(t, err)
	jsonData, err := call.MarshalJSON()
	assert.Nil(t, err)
	assert.NotNil(t, jsonData)
	expectedJson := fmt.Sprintf(`[2,"%v","Authorize",{"idTag":"%v"}]`, uniqueId, idTag)
	assert.Equal(t, []byte(expectedJson), jsonData)
}

func (suite *OcppV16TestSuite) TestAuthorizeConfirmationFromJson() {
	t := suite.T()
	uniqueId := "5678"
	rawTime := time.Now().Add(time.Hour * 8).Format(ocpp16.ISO8601)
	expiryDate, err := time.Parse(ocpp16.ISO8601, rawTime)
	assert.Nil(t, err)
	parentIdTag := "parentTag1"
	status := ocpp16.AuthorizationStatusAccepted
	dummyRequest := ocpp16.AuthorizeRequest{}
	dataJson := fmt.Sprintf(`[3,"%v",{"idTagInfo":{"expiryDate":"%v","parentIdTag":"%v","status":"%v"}}]`, uniqueId, expiryDate.Format(ocpp16.ISO8601), parentIdTag, status)
	suite.chargePoint.Endpoint.AddPendingRequest(uniqueId, dummyRequest)
	callResult := ParseCallResult(&suite.chargePoint.Endpoint, dataJson, t)
	CheckCallResult(callResult, t, uniqueId)
	confirmation := GetAuthorizeConfirmation(t, callResult.Payload)
	assert.Equal(t, status, confirmation.IdTagInfo.Status)
	assert.Equal(t, parentIdTag, confirmation.IdTagInfo.ParentIdTag)
	assert.Equal(t, expiryDate, confirmation.IdTagInfo.ExpiryDate)
}

func (suite *OcppV16TestSuite) TestAuthorizeConfirmationToJson() {
	t := suite.T()
	uniqueId := "1234"
	parentIdTag := "parentTag1"
	expiryDate := time.Now().Add(time.Hour * 8)
	status := ocpp16.AuthorizationStatusAccepted
	confirmation := ocpp16.AuthorizeConfirmation{IdTagInfo: ocpp16.IdTagInfo{Status: status, ParentIdTag: parentIdTag, ExpiryDate: expiryDate}}
	callResult, err := suite.centralSystem.CreateCallResult(confirmation, uniqueId)
	assert.Nil(t, err)
	assert.NotNil(t, callResult)
	err = Validate.Struct(callResult)
	assert.Nil(t, err)
	jsonData, err := callResult.MarshalJSON()
	assert.Nil(t, err)
	assert.NotNil(t, jsonData)
	expectedJson := fmt.Sprintf(`[3,"%v",{"idTagInfo":{"expiryDate":"%v","parentIdTag":"%v","status":"%v"}}]`, uniqueId, expiryDate.Format(time.RFC3339Nano), parentIdTag, status)
	assert.Equal(t, []byte(expectedJson), jsonData)
}

func (suite *OcppV16TestSuite) TestAuthorizeE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := "1234"
	wsUrl := "someUrl"
	idTag := "tag1"
	parentIdTag := "parentTag1"
	status := ocpp16.AuthorizationStatusAccepted
	expiryDate := time.Now().Add(time.Hour * 8)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"idTag":"%v"}]`, messageId, ocpp16.AuthorizeFeatureName, idTag)
	responseJson := fmt.Sprintf(`[3,"%v",{"idTagInfo":{"expiryDate":"%v","parentIdTag":"%v","status":"%v"}}]`, messageId, expiryDate.Format(time.RFC3339Nano), parentIdTag, status)
	requestRaw := []byte(requestJson)
	responseRaw := []byte(responseJson)
	channel := NewMockWebSocket(wsId)
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
		CheckCall(call, t, ocpp16.AuthorizeFeatureName, messageId)
		suite.chargePoint.AddPendingRequest(messageId, call.Payload)
		// TODO: generate the response dynamically
		err := suite.mockClient.MessageHandler(responseRaw)
		assert.Nil(t, err)
		return nil
	})
	// Setting client handlers
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil).Run(func(args mock.Arguments) {
		u := args.String(0)
		assert.Equal(t, wsUrl, u)
		suite.mockServer.NewClientHandler(channel)
	})
	suite.mockClient.SetMessageHandler(func(data []byte) error {
		assert.Equal(t, responseRaw, data)
		jsonData := string(data)
		assert.Equal(t, responseJson, jsonData)
		callResult := ParseCallResult(&suite.chargePoint.Endpoint, jsonData, t)
		CheckCallResult(callResult, t, messageId)
		return nil
	})
	suite.mockClient.On("Write", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		data := args.Get(0)
		bytes := data.([]byte)
		assert.NotNil(t, bytes)
		err := suite.mockServer.MessageHandler(channel, bytes)
		assert.Nil(t, err)
	})
	// Test Run
	err := suite.mockClient.Start(wsUrl)
	assert.Nil(t, err)
	err = suite.mockClient.Write(requestRaw)
	assert.Nil(t, err)
}
