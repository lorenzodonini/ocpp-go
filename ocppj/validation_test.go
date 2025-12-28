package ocppj_test

import (
	"testing"

	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocppj"
	"github.com/stretchr/testify/suite"
)

type ValidationTestSuite struct {
	suite.Suite
	endpoint ocppj.Endpoint
}

func (suite *ValidationTestSuite) SetupTest() {
	mockProfile := ocpp.NewProfile("mock", &MockFeature{})
	suite.endpoint = ocppj.Endpoint{}
	suite.endpoint.AddProfile(mockProfile)
	// Ensure validation is enabled by default for each test
	ocppj.SetMessageValidation(true)
}

func (suite *ValidationTestSuite) TearDownTest() {
	// Reset validation to enabled state after each test
	ocppj.SetMessageValidation(true)
}

func (suite *ValidationTestSuite) TestSetMessageValidation_DefaultEnabled() {
	// Test with invalid request (empty MockValue violates "required" constraint)
	invalidRequest := &MockRequest{MockValue: ""}
	call, err := suite.endpoint.CreateCall(invalidRequest)
	suite.Require().Error(err, "validation should fail for invalid request when enabled")
	suite.Assert().Nil(call, "call should be nil when validation fails")

	// Test with valid request
	validRequest := &MockRequest{MockValue: "valid"}
	call, err = suite.endpoint.CreateCall(validRequest)
	suite.Require().NoError(err, "validation should pass for valid request")
	suite.Assert().NotNil(call, "call should not be nil when validation passes")
}

func (suite *ValidationTestSuite) TestSetMessageValidation_DisableValidation() {
	// Disable validation
	ocppj.SetMessageValidation(false)

	// Test with invalid request (empty MockValue violates "required" constraint)
	invalidRequest := &MockRequest{MockValue: ""}
	call, err := suite.endpoint.CreateCall(invalidRequest)
	suite.Require().NoError(err, "validation should be skipped when disabled")
	suite.Assert().NotNil(call, "call should be created even with invalid request when validation is disabled")

	// Test with invalid request (too long MockValue violates "max=10" constraint)
	invalidRequestLong := &MockRequest{MockValue: "this is way too long"}
	call, err = suite.endpoint.CreateCall(invalidRequestLong)
	suite.Require().NoError(err, "validation should be skipped when disabled")
	suite.Assert().NotNil(call, "call should be created even with invalid request when validation is disabled")

	// Test with invalid CallResult (empty MockValue violates "required" and "min=5" constraints)
	invalidConfirmation := &MockConfirmation{MockValue: ""}
	callResult, err := suite.endpoint.CreateCallResult(invalidConfirmation, "test-id")
	suite.Require().NoError(err, "validation should be skipped when disabled")
	suite.Assert().NotNil(callResult, "callResult should be created even with invalid confirmation when validation is disabled")

	// Test with invalid CallResult (too short MockValue violates "min=5" constraint)
	invalidConfirmationShort := &MockConfirmation{MockValue: "min"}
	callResult, err = suite.endpoint.CreateCallResult(invalidConfirmationShort, "test-id")
	suite.Require().NoError(err, "validation should be skipped when disabled")
	suite.Assert().NotNil(callResult, "callResult should be created even with invalid confirmation when validation is disabled")
}

func (suite *ValidationTestSuite) TestSetMessageValidation_EnableValidation() {
	// First disable validation
	ocppj.SetMessageValidation(false)

	// Verify invalid request passes when disabled
	invalidRequest := &MockRequest{MockValue: ""}
	call, err := suite.endpoint.CreateCall(invalidRequest)
	suite.Require().NoError(err, "validation should be skipped when disabled")
	suite.Assert().NotNil(call)

	// Now enable validation
	ocppj.SetMessageValidation(true)

	// Same invalid request should now fail
	call, err = suite.endpoint.CreateCall(invalidRequest)
	suite.Require().Error(err, "validation should fail for invalid request when enabled")
	suite.Assert().Nil(call, "call should be nil when validation fails")

	// Valid request should still pass
	validRequest := &MockRequest{MockValue: "valid"}
	call, err = suite.endpoint.CreateCall(validRequest)
	suite.Require().NoError(err, "validation should pass for valid request")
	suite.Assert().NotNil(call, "call should not be nil when validation passes")
}

func (suite *ValidationTestSuite) TestSetMessageValidation_ToggleValidation() {
	invalidRequest := &MockRequest{MockValue: ""}
	validRequest := &MockRequest{MockValue: "valid"}

	// Test 1: Enable -> Disable -> Enable
	ocppj.SetMessageValidation(true)
	call, err := suite.endpoint.CreateCall(invalidRequest)
	suite.Require().Error(err, "should fail when enabled")
	suite.Assert().Nil(call)

	ocppj.SetMessageValidation(false)
	call, err = suite.endpoint.CreateCall(invalidRequest)
	suite.Require().NoError(err, "should pass when disabled")
	suite.Assert().NotNil(call)

	ocppj.SetMessageValidation(true)
	call, err = suite.endpoint.CreateCall(invalidRequest)
	suite.Require().Error(err, "should fail when enabled again")
	suite.Assert().Nil(call)

	// Test 2: Valid request should always pass regardless of validation setting
	ocppj.SetMessageValidation(true)
	call, err = suite.endpoint.CreateCall(validRequest)
	suite.Require().NoError(err, "valid request should pass when enabled")
	suite.Assert().NotNil(call)

	ocppj.SetMessageValidation(false)
	call, err = suite.endpoint.CreateCall(validRequest)
	suite.Require().NoError(err, "valid request should pass when disabled")
	suite.Assert().NotNil(call)
}

func TestValidationSuite(t *testing.T) {
	if !testing.Short() {
		t.Skip("")
	}
	suite.Run(t, new(ValidationTestSuite))
}
