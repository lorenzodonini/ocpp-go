package test_v16

import (
	"fmt"
	"github.com/lorenzodonini/go-ocpp/ocpp"
	"github.com/lorenzodonini/go-ocpp/ocpp/1.6"
	"github.com/lorenzodonini/go-ocpp/test"
	"github.com/lorenzodonini/go-ocpp/ws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

// Utility functions
func GetAuthorizeRequest(t *testing.T, request ocpp.Request) *v16.AuthorizeRequest {
	assert.NotNil(t, request)
	result := request.(*v16.AuthorizeRequest)
	assert.NotNil(t, result)
	assert.IsType(t, &v16.AuthorizeRequest{}, result)
	return result
}

func GetAuthorizeConfirmation(t *testing.T, confirmation ocpp.Confirmation) *v16.AuthorizeConfirmation {
	assert.NotNil(t, confirmation)
	result := confirmation.(*v16.AuthorizeConfirmation)
	assert.NotNil(t, result)
	assert.IsType(t, &v16.AuthorizeConfirmation{}, result)
	return result
}

// Test
func (suite *OcppV16TestSuite) TestAuthorizeRequestValidation() {
	t := suite.T()
	var requestTable = []test.RequestTestEntry{
		{v16.AuthorizeRequest{IdTag: "12345"}, true},
		{v16.AuthorizeRequest{}, false},
		{v16.AuthorizeRequest{IdTag: ">20.................."}, false},
	}
	test.ExecuteRequestTestTable(t, requestTable)
}

func (suite *OcppV16TestSuite) TestAuthorizeConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []test.ConfirmationTestEntry{
		{v16.AuthorizeConfirmation{IdTagInfo: v16.IdTagInfo{ExpiryDate: time.Now().Add(time.Hour * 8), ParentIdTag: "00000", Status: v16.AuthorizationStatusAccepted}}, true},
		{v16.AuthorizeConfirmation{IdTagInfo: v16.IdTagInfo{ParentIdTag: "00000", Status: v16.AuthorizationStatusAccepted}}, true},
		{v16.AuthorizeConfirmation{IdTagInfo: v16.IdTagInfo{ExpiryDate: time.Now().Add(time.Hour * 8), Status: v16.AuthorizationStatusAccepted}}, true},
		{v16.AuthorizeConfirmation{IdTagInfo: v16.IdTagInfo{Status: v16.AuthorizationStatusAccepted}}, true},
		{v16.AuthorizeConfirmation{IdTagInfo: v16.IdTagInfo{Status: v16.AuthorizationStatusBlocked}}, true},
		{v16.AuthorizeConfirmation{IdTagInfo: v16.IdTagInfo{Status: v16.AuthorizationStatusExpired}}, true},
		{v16.AuthorizeConfirmation{IdTagInfo: v16.IdTagInfo{Status: v16.AuthorizationStatusInvalid}}, true},
		{v16.AuthorizeConfirmation{IdTagInfo: v16.IdTagInfo{Status: v16.AuthorizationStatusConcurrentTx}}, true},
		{v16.AuthorizeConfirmation{IdTagInfo: v16.IdTagInfo{ParentIdTag: ">20..................", Status: v16.AuthorizationStatusAccepted}}, false},
		{v16.AuthorizeConfirmation{IdTagInfo: v16.IdTagInfo{ExpiryDate: time.Now().Add(time.Hour * -8), Status: v16.AuthorizationStatusAccepted}}, false},
	}
	test.ExecuteConfirmationTestTable(t, confirmationTable)
}

func (suite *OcppV16TestSuite) TestAuthorizeRequestFromJson() {
	t := suite.T()
	uniqueId := "1234"
	idTag := "tag1"
	dataJson := fmt.Sprintf(`[2,"%v","Authorize",{"idTag":"%v"}]`, uniqueId, idTag)
	call := test.ParseCall(&suite.centralSystem.Endpoint, dataJson, t)
	test.CheckCall(call, t, v16.AuthorizeFeatureName, uniqueId)
	request := GetAuthorizeRequest(t, call.Payload)
	assert.Equal(t, idTag, request.IdTag)
}

func (suite *OcppV16TestSuite) TestAuthorizeRequestToJson() {
	t := suite.T()
	idTag := "tag2"
	request := v16.AuthorizeRequest{IdTag: idTag}
	call, err := suite.chargePoint.CreateCall(request)
	assert.Nil(t, err)
	uniqueId := call.GetUniqueId()
	assert.NotNil(t, call)
	err = test.Validate.Struct(call)
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
	rawTime := time.Now().Add(time.Hour * 8).Format(v16.ISO8601)
	expiryDate, err := time.Parse(v16.ISO8601, rawTime)
	assert.Nil(t, err)
	parentIdTag := "parentTag1"
	status := v16.AuthorizationStatusAccepted
	dummyRequest := v16.AuthorizeRequest{}
	dataJson := fmt.Sprintf(`[3,"%v",{"idTagInfo":{"expiryDate":"%v","parentIdTag":"%v","status":"%v"}}]`, uniqueId, expiryDate.Format(v16.ISO8601), parentIdTag, status)
	suite.chargePoint.Endpoint.AddPendingRequest(uniqueId, dummyRequest)
	callResult := test.ParseCallResult(&suite.chargePoint.Endpoint, dataJson, t)
	test.CheckCallResult(callResult, t, uniqueId)
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
	status := v16.AuthorizationStatusAccepted
	confirmation := v16.AuthorizeConfirmation{IdTagInfo: v16.IdTagInfo{Status: status, ParentIdTag: parentIdTag, ExpiryDate: expiryDate}}
	callResult, err := suite.centralSystem.CreateCallResult(confirmation, uniqueId)
	assert.Nil(t, err)
	assert.NotNil(t, callResult)
	err = test.Validate.Struct(callResult)
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
	status := v16.AuthorizationStatusAccepted
	expiryDate := time.Now().Add(time.Hour * 8)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"idTag":"%v"}]`, messageId, v16.AuthorizeFeatureName, idTag)
	responseJson := fmt.Sprintf(`[3,"%v",{"idTagInfo":{"expiryDate":"%v","parentIdTag":"%v","status":"%v"}}]`, messageId, expiryDate.Format(time.RFC3339Nano), parentIdTag, status)
	requestRaw := []byte(requestJson)
	responseRaw := []byte(responseJson)
	channel := test.NewMockWebSocket(wsId)
	// Setting server handlers
	suite.mockServer.SetNewClientHandler(func(ws ws.Channel) {
		assert.NotNil(t, ws)
		assert.Equal(t, wsId, ws.GetId())
	})
	suite.mockServer.SetMessageHandler(func(ws ws.Channel, data []byte) error {
		assert.Equal(t, requestRaw, data)
		jsonData := string(data)
		assert.Equal(t, requestJson, jsonData)
		call := test.ParseCall(&suite.chargePoint.Endpoint, jsonData, t)
		test.CheckCall(call, t, v16.AuthorizeFeatureName, messageId)
		suite.chargePoint.AddPendingRequest(messageId, call.Payload)
		// TODO: generate the response dynamically
		err := suite.mockClient.MessageHandler(responseRaw)
		assert.Nil(t, err)
		return nil
	})
	// Setting client handlers
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return().Run(func(args mock.Arguments) {
		u := args.String(0)
		assert.Equal(t, wsUrl, u)
		suite.mockServer.NewClientHandler(channel)
	})
	suite.mockClient.SetMessageHandler(func(data []byte) error {
		assert.Equal(t, responseRaw, data)
		jsonData := string(data)
		assert.Equal(t, responseJson, jsonData)
		callResult := test.ParseCallResult(&suite.chargePoint.Endpoint, jsonData, t)
		test.CheckCallResult(callResult, t, messageId)
		return nil
	})
	suite.mockClient.On("Write", mock.Anything).Return().Run(func(args mock.Arguments) {
		data := args.Get(0)
		bytes := data.([]byte)
		assert.NotNil(t, bytes)
		err := suite.mockServer.MessageHandler(channel, bytes)
		assert.Nil(t, err)
	})
	// Test Run
	err := suite.mockClient.Start(wsUrl)
	assert.Nil(t, err)
	suite.mockClient.Write(requestRaw)
}
