package test

import (
	"fmt"
	"github.com/lorenzodonini/go-ocpp/ocpp"
	"github.com/lorenzodonini/go-ocpp/ocpp/1.6"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// Utility functions
func GetAuthorizeRequest(t* testing.T, request ocpp.Request) *v16.AuthorizeRequest {
	assert.NotNil(t, request)
	result := request.(*v16.AuthorizeRequest)
	assert.NotNil(t, result)
	assert.IsType(t, &v16.AuthorizeRequest{}, result)
	return result
}

func GetAuthorizeConfirmation(t* testing.T, confirmation ocpp.Confirmation) *v16.AuthorizeConfirmation {
	assert.NotNil(t, confirmation)
	result := confirmation.(*v16.AuthorizeConfirmation)
	assert.NotNil(t, result)
	assert.IsType(t, &v16.AuthorizeConfirmation{}, result)
	return result
}

// Test
func (suite *OcppV16TestSuite) TestAuthorizeRequestValidation() {
	t := suite.T()
	var requestTable = []RequestTestEntry{
		{v16.AuthorizeRequest{IdTag: "12345"}, true},
		{v16.AuthorizeRequest{}, false},
		{v16.AuthorizeRequest{IdTag: ">20.................."}, false},
	}
	executeRequestTestTable(t, requestTable)
}

func (suite *OcppV16TestSuite) TestAuthorizeConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []ConfirmationTestEntry {
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
	executeConfirmationTestTable(t, confirmationTable)
}

func (suite *OcppV16TestSuite) TestAuthorizeRequestFromJson() {
	t := suite.T()
	uniqueId := "1234"
	idTag := "tag1"
	dataJson := fmt.Sprintf(`[2,"%v","Authorize",{"idTag":"%v"}]`, uniqueId, idTag)
	call := ParseCall(&suite.centralSystem.Endpoint, dataJson, t)
	CheckCall(call, t, v16.AuthorizeFeatureName, uniqueId)
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
	err = validate.Struct(call)
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
	status := v16.AuthorizationStatusAccepted
	confirmation := v16.AuthorizeConfirmation{IdTagInfo: v16.IdTagInfo{Status: status, ParentIdTag: parentIdTag, ExpiryDate: expiryDate}}
	callResult, err := suite.centralSystem.CreateCallResult(confirmation, uniqueId)
	assert.Nil(t, err)
	assert.NotNil(t, callResult)
	err = validate.Struct(callResult)
	assert.Nil(t, err)
	jsonData, err := callResult.MarshalJSON()
	assert.Nil(t, err)
	assert.NotNil(t, jsonData)
	expectedJson := fmt.Sprintf(`[3,"%v",{"idTagInfo":{"expiryDate":"%v","parentIdTag":"%v","status":"%v"}}]`, uniqueId, expiryDate.Format(time.RFC3339Nano), parentIdTag, status)
	assert.Equal(t, []byte(expectedJson), jsonData)
}
