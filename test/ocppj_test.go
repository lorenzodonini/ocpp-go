package test

import (
	"fmt"
	"github.com/lorenzodonini/go-ocpp/ocpp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
	"testing"
)

// Mock Feature
const (
	MockFeatureName = "Mock"
)

type MockRequest struct {
	mock.Mock
	MockValue string `json:"mockValue" validate:"required,max=10"`
}

type MockConfirmation struct {
	mock.Mock
	MockValue string `json:"mockValue" validate:"required,min=5"`
}

type MockFeature struct {
	mock.Mock
}

func (f MockFeature) GetFeatureName() string {
	return MockFeatureName
}

func (f MockFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(MockRequest{})
}

func (f MockFeature) GetConfirmationType() reflect.Type {
	return reflect.TypeOf(MockConfirmation{})
}

func (r MockRequest) GetFeatureName() string {
	return MockFeatureName
}

func (c MockConfirmation) GetFeatureName() string {
	return MockFeatureName
}

func newMockRequest(value string) *MockRequest {
	return &MockRequest{MockValue: value}
}

func newMockConfirmation(value string) *MockConfirmation {
	return &MockConfirmation{MockValue: value}
}

// Tests
type OcppJTestSuite struct {
	suite.Suite
	chargePoint   *ocpp.ChargePoint
	centralSystem *ocpp.CentralSystem
	mockServer    *MockWebsocketServer
	mockClient    *MockWebsocketClient
}

func (suite *OcppJTestSuite) SetupTest() {
	mockProfile := ocpp.NewProfile("mock", MockFeature{})
	mockClient := MockWebsocketClient{}
	mockServer := MockWebsocketServer{}
	suite.mockClient = &mockClient
	suite.mockServer = &mockServer
	suite.chargePoint = ocpp.NewChargePoint("mock_id", suite.mockClient, mockProfile)
	suite.centralSystem = ocpp.NewCentralSystem(suite.mockServer, mockProfile)
}

// Protocol functions test
func (suite *OcppJTestSuite) TestGetProfile() {
	t := suite.T()
	profile, ok := suite.chargePoint.GetProfile("mock")
	assert.True(t, ok)
	assert.NotNil(t, profile)
	feature := profile.GetFeature(MockFeatureName)
	assert.NotNil(t, feature)
	assert.Equal(t, reflect.TypeOf(MockRequest{}), feature.GetRequestType())
	assert.Equal(t, reflect.TypeOf(MockConfirmation{}), feature.GetConfirmationType())
}

func (suite *OcppJTestSuite) TestGetProfileForFeature() {
	t := suite.T()
	profile, ok := suite.chargePoint.GetProfileForFeature(MockFeatureName)
	assert.True(t, ok)
	assert.NotNil(t, profile)
	assert.Equal(t, "mock", profile.Name)
}

//func (suite *OcppJTestSuite) TestAddFeature

func (suite *OcppJTestSuite) TestGetProfileForInvalidFeature() {
	t := suite.T()
	profile, ok := suite.chargePoint.GetProfileForFeature("test")
	assert.False(t, ok)
	assert.Nil(t, profile)
}

func (suite *OcppJTestSuite) TestCallValidation() {
	t := suite.T()
	mockLongValue := "somelongvalue"
	request := newMockRequest(mockLongValue)
	// Test invalid call
	call, err := suite.chargePoint.CreateCall(request)
	assert.Nil(t, call)
	assert.NotNil(t, err)
	assert.IsType(t, validator.ValidationErrors{}, err)
	errors := err.(validator.ValidationErrors)
	assert.Equal(t, 1, len(errors))
	validationError := errors[0]
	assert.Equal(t, "max", validationError.Tag())
}

func (suite *OcppJTestSuite) TestCallResultValidation() {
	t := suite.T()
	mockShortValue := "val"
	mockUniqueId := "123456"
	confirmation := newMockConfirmation(mockShortValue)
	// Test invalid call
	callResult, err := suite.chargePoint.CreateCallResult(confirmation, mockUniqueId)
	assert.Nil(t, callResult)
	assert.NotNil(t, err)
	assert.IsType(t, validator.ValidationErrors{}, err)
	errors := err.(validator.ValidationErrors)
	assert.Equal(t, 1, len(errors))
	validationError := errors[0]
	assert.Equal(t, "min", validationError.Tag())
}

func (suite *OcppJTestSuite) TestCreateCall() {
	t := suite.T()
	mockValue := "somevalue"
	request := newMockRequest(mockValue)
	call, err := suite.chargePoint.CreateCall(request)
	assert.Nil(t, err)
	CheckCall(call, t, MockFeatureName, call.UniqueId)
	message, ok := call.Payload.(*MockRequest)
	assert.True(t, ok)
	assert.NotNil(t, message)
	assert.Equal(t, mockValue, message.MockValue)
	// Check that request was stored as pending request
	pendingRequest, exists := suite.chargePoint.GetPendingRequest(call.UniqueId)
	assert.True(t, exists)
	assert.NotNil(t, pendingRequest)
	suite.chargePoint.DeletePendingRequest(call.UniqueId)
}

func (suite *OcppJTestSuite) TestCreateCallResult() {
	t := suite.T()
	mockValue := "someothervalue"
	mockUniqueId := "123456"
	confirmation := newMockConfirmation(mockValue)
	callResult, err := suite.chargePoint.CreateCallResult(confirmation, mockUniqueId)
	assert.Nil(t, err)
	CheckCallResult(callResult, t, mockUniqueId)
	message, ok := callResult.Payload.(*MockConfirmation)
	assert.True(t, ok)
	assert.NotNil(t, message)
	assert.Equal(t, mockValue, message.MockValue)

}

func (suite *OcppJTestSuite) TestCreateCallError() {
	t := suite.T()
	mockUniqueId := "123456"
	mockDescription := "somedescription"
	mockDetailString := "somedetailstring"
	type MockDetails struct {
		DetailString string
	}
	mockDetails := MockDetails{DetailString: mockDetailString}
	callError := suite.chargePoint.CreateCallError(mockUniqueId, ocpp.GenericError, mockDescription, mockDetails)
	assert.NotNil(t, callError)
	CheckCallError(t, callError, mockUniqueId, ocpp.GenericError, mockDescription, mockDetails)
}

func (suite *OcppJTestSuite) TestParseMessageInvalidLength() {
	t := suite.T()
	mockMessage := make([]interface{}, 2)
	messageId := "12345"
	// Test invalid message length
	mockMessage[0] = ocpp.CALL 	// Message Type ID
	mockMessage[1] = messageId 	// Unique ID
	message, protoErr := suite.chargePoint.ParseMessage(mockMessage)
	assert.Nil(t, message)
	assert.NotNil(t, protoErr)
	assert.Equal(t, "", protoErr.MessageId)
	assert.Equal(t, ocpp.FormationViolation, protoErr.ErrorCode)
	assert.Equal(t, "Invalid message. Expected array length >= 3", protoErr.Error.Error())
}

func (suite *OcppJTestSuite) TestParseMessageInvalidTypeId() {
	t := suite.T()
	mockMessage := make([]interface{}, 3)
	invalidTypeId := "2"
	messageId := "12345"
	// Test invalid message length
	mockMessage[0] = invalidTypeId 	// Message Type ID
	mockMessage[1] = messageId 		// Unique ID
	message, protoErr := suite.chargePoint.ParseMessage(mockMessage)
	assert.Nil(t, message)
	assert.NotNil(t, protoErr)
	assert.Equal(t, "", protoErr.MessageId)
	assert.Equal(t, ocpp.FormationViolation, protoErr.ErrorCode)
	assert.Equal(t, fmt.Sprintf("Invalid element %v at 0, expected message type (int)", invalidTypeId), protoErr.Error.Error())
}

func (suite *OcppJTestSuite) TestParseMessageInvalidMessageId() {
	t := suite.T()
	mockMessage := make([]interface{}, 3)
	invalidMessageId := 12345
	// Test invalid message length
	mockMessage[0] = float64(ocpp.CALL)		// Message Type ID
	mockMessage[1] = invalidMessageId 		// Unique ID
	message, protoErr := suite.chargePoint.ParseMessage(mockMessage)
	assert.Nil(t, message)
	assert.NotNil(t, protoErr)
	assert.Equal(t, "", protoErr.MessageId)
	assert.Equal(t, ocpp.FormationViolation, protoErr.ErrorCode)
	assert.Equal(t, fmt.Sprintf("Invalid element %v at 1, expected unique ID (string)", invalidMessageId), protoErr.Error.Error())
}

func (suite *OcppJTestSuite) TestParseMessageUnknownTypeId() {
	t := suite.T()
	mockMessage := make([]interface{}, 3)
	messageId := "12345"
	invalidTypeId := 1
	// Test invalid message length
	mockMessage[0] = float64(invalidTypeId)		// Message Type ID
	mockMessage[1] = messageId 					// Unique ID
	message, protoErr := suite.chargePoint.ParseMessage(mockMessage)
	assert.Nil(t, message)
	assert.NotNil(t, protoErr)
	assert.Equal(t, messageId, protoErr.MessageId)
	assert.Equal(t, ocpp.FormationViolation, protoErr.ErrorCode)
	assert.Equal(t, fmt.Sprintf("Invalid message type ID %v", invalidTypeId), protoErr.Error.Error())
}

//TODO: implement generic protocol tests

func TestOcppJ(t *testing.T) {
	suite.Run(t, new(OcppJTestSuite))
}
