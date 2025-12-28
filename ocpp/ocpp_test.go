package ocpp

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
)

// Mock feature for testing
type TestFeature struct {
	name         string
	requestType  reflect.Type
	responseType reflect.Type
}

func (f *TestFeature) GetFeatureName() string {
	return f.name
}

func (f *TestFeature) GetRequestType() reflect.Type {
	return f.requestType
}

func (f *TestFeature) GetResponseType() reflect.Type {
	return f.responseType
}

// Mock request and response for testing ParseRequest and ParseResponse
type TestRequest struct {
	Value string
}

func (r *TestRequest) GetFeatureName() string {
	return "TestFeature"
}

type TestResponse struct {
	Result int
}

func (r *TestResponse) GetFeatureName() string {
	return "TestFeature"
}

type ProfileTestSuite struct {
	suite.Suite
}

func TestProfileSuite(t *testing.T) {
	suite.Run(t, new(ProfileTestSuite))
}

func (suite *ProfileTestSuite) TestNewProfile_BasicCreation() {
	tests := []struct {
		name           string
		profileName    string
		features       []Feature
		expectedLength int
	}{
		{
			name:           "empty name",
			profileName:    "",
			features:       nil,
			expectedLength: 0,
		},
		{
			name:           "with name, no features",
			profileName:    "test-profile",
			features:       nil,
			expectedLength: 0,
		},
		{
			name:           "empty profile name",
			profileName:    "empty-profile",
			features:       nil,
			expectedLength: 0,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			profile := NewProfile(tt.profileName, tt.features...)
			suite.Assert().NotNil(profile)
			suite.Assert().Equal(tt.profileName, profile.Name)
			suite.Assert().NotNil(profile.features)
			suite.Assert().Equal(tt.expectedLength, len(profile.features))
			suite.Assert().False(profile.SupportsFeature("nonexistent"))
		})
	}
}

func (suite *ProfileTestSuite) TestNewProfile_WithFeatures() {
	tests := []struct {
		name             string
		profileName      string
		features         []*TestFeature
		expectedLength   int
		expectedFeatures []struct {
			name         string
			requestType  reflect.Type
			responseType reflect.Type
		}
		shouldNotExist []string
	}{
		{
			name:        "single feature",
			profileName: "single-feature-profile",
			features: []*TestFeature{
				{
					name:         "TestFeature",
					requestType:  reflect.TypeOf(""),
					responseType: reflect.TypeOf(0),
				},
			},
			expectedLength: 1,
			expectedFeatures: []struct {
				name         string
				requestType  reflect.Type
				responseType reflect.Type
			}{
				{
					name:         "TestFeature",
					requestType:  reflect.TypeOf(""),
					responseType: reflect.TypeOf(0),
				},
			},
			shouldNotExist: []string{"NonexistentFeature"},
		},
		{
			name:        "multiple features",
			profileName: "multi-feature-profile",
			features: []*TestFeature{
				{
					name:         "Feature1",
					requestType:  reflect.TypeOf(""),
					responseType: reflect.TypeOf(0),
				},
				{
					name:         "Feature2",
					requestType:  reflect.TypeOf(0),
					responseType: reflect.TypeOf(""),
				},
				{
					name:         "Feature3",
					requestType:  reflect.TypeOf(true),
					responseType: reflect.TypeOf(1.0),
				},
			},
			expectedLength: 3,
			expectedFeatures: []struct {
				name         string
				requestType  reflect.Type
				responseType reflect.Type
			}{
				{
					name:         "Feature1",
					requestType:  reflect.TypeOf(""),
					responseType: reflect.TypeOf(0),
				},
				{
					name:         "Feature2",
					requestType:  reflect.TypeOf(0),
					responseType: reflect.TypeOf(""),
				},
				{
					name:         "Feature3",
					requestType:  reflect.TypeOf(true),
					responseType: reflect.TypeOf(1.0),
				},
			},
			shouldNotExist: []string{"NonexistentFeature"},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Convert []*TestFeature to []Feature
			features := make([]Feature, len(tt.features))
			for i, f := range tt.features {
				features[i] = f
			}

			profile := NewProfile(tt.profileName, features...)
			suite.Assert().NotNil(profile)
			suite.Assert().Equal(tt.profileName, profile.Name)
			suite.Assert().NotNil(profile.features)
			suite.Assert().Equal(tt.expectedLength, len(profile.features))

			// Verify expected features exist and are correct
			for _, expected := range tt.expectedFeatures {
				suite.Assert().True(profile.SupportsFeature(expected.name))
				retrievedFeature := profile.GetFeature(expected.name)
				suite.Require().NotNil(retrievedFeature)
				suite.Assert().Equal(expected.name, retrievedFeature.GetFeatureName())
				suite.Assert().Equal(expected.requestType, retrievedFeature.GetRequestType())
				suite.Assert().Equal(expected.responseType, retrievedFeature.GetResponseType())
			}

			// Verify features that should not exist
			for _, name := range tt.shouldNotExist {
				suite.Assert().False(profile.SupportsFeature(name))
				suite.Assert().Nil(profile.GetFeature(name))
			}
		})
	}
}

func (suite *ProfileTestSuite) TestNewProfile_DuplicateFeatures() {
	tests := []struct {
		name            string
		profileName     string
		features        []*TestFeature
		expectedLength  int
		expectedFeature struct {
			name         string
			requestType  reflect.Type
			responseType reflect.Type
		}
	}{
		{
			name:        "duplicate features overwrite",
			profileName: "duplicate-feature-profile",
			features: []*TestFeature{
				{
					name:         "DuplicateFeature",
					requestType:  reflect.TypeOf(""),
					responseType: reflect.TypeOf(0),
				},
				{
					name:         "DuplicateFeature",
					requestType:  reflect.TypeOf(0),
					responseType: reflect.TypeOf(""),
				},
			},
			expectedLength: 1,
			expectedFeature: struct {
				name         string
				requestType  reflect.Type
				responseType reflect.Type
			}{
				name:         "DuplicateFeature",
				requestType:  reflect.TypeOf(0),
				responseType: reflect.TypeOf(""),
			},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Convert []*TestFeature to []Feature
			features := make([]Feature, len(tt.features))
			for i, f := range tt.features {
				features[i] = f
			}

			// When duplicate features are added, the last one should overwrite the previous one
			profile := NewProfile(tt.profileName, features...)
			suite.Assert().NotNil(profile)
			suite.Assert().Equal(tt.profileName, profile.Name)
			suite.Assert().NotNil(profile.features)
			suite.Assert().Equal(tt.expectedLength, len(profile.features))

			retrievedFeature := profile.GetFeature(tt.expectedFeature.name)
			suite.Require().NotNil(retrievedFeature)
			// Should be the last feature added (overwrites previous)
			suite.Assert().Equal(tt.expectedFeature.requestType, retrievedFeature.GetRequestType())
			suite.Assert().Equal(tt.expectedFeature.responseType, retrievedFeature.GetResponseType())
		})
	}
}

func (suite *ProfileTestSuite) TestProfile_AddFeature() {
	tests := []struct {
		name             string
		profileName      string
		initialFeatures  []*TestFeature
		featuresToAdd    []*TestFeature
		expectedLength   int
		expectedFeatures []string
	}{
		{
			name:            "add to empty profile",
			profileName:     "test-profile",
			initialFeatures: nil,
			featuresToAdd: []*TestFeature{
				{
					name:         "Feature1",
					requestType:  reflect.TypeOf(""),
					responseType: reflect.TypeOf(0),
				},
			},
			expectedLength:   1,
			expectedFeatures: []string{"Feature1"},
		},
		{
			name:        "add multiple features",
			profileName: "test-profile",
			initialFeatures: []*TestFeature{
				{
					name:         "ExistingFeature",
					requestType:  reflect.TypeOf(""),
					responseType: reflect.TypeOf(0),
				},
			},
			featuresToAdd: []*TestFeature{
				{
					name:         "NewFeature1",
					requestType:  reflect.TypeOf(0),
					responseType: reflect.TypeOf(""),
				},
				{
					name:         "NewFeature2",
					requestType:  reflect.TypeOf(true),
					responseType: reflect.TypeOf(1.0),
				},
			},
			expectedLength:   3,
			expectedFeatures: []string{"ExistingFeature", "NewFeature1", "NewFeature2"},
		},
		{
			name:        "overwrite existing feature",
			profileName: "test-profile",
			initialFeatures: []*TestFeature{
				{
					name:         "Feature1",
					requestType:  reflect.TypeOf(""),
					responseType: reflect.TypeOf(0),
				},
			},
			featuresToAdd: []*TestFeature{
				{
					name:         "Feature1",
					requestType:  reflect.TypeOf(0),
					responseType: reflect.TypeOf(""),
				},
			},
			expectedLength:   1,
			expectedFeatures: []string{"Feature1"},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Create profile with initial features
			initialFeatures := make([]Feature, len(tt.initialFeatures))
			for i, f := range tt.initialFeatures {
				initialFeatures[i] = f
			}
			profile := NewProfile(tt.profileName, initialFeatures...)

			// Add new features
			for _, feature := range tt.featuresToAdd {
				profile.AddFeature(feature)
			}

			suite.Assert().Equal(tt.expectedLength, len(profile.features))
			for _, expectedName := range tt.expectedFeatures {
				suite.Assert().True(profile.SupportsFeature(expectedName))
			}

			// Verify last added feature overwrites if duplicate
			if len(tt.featuresToAdd) > 0 {
				lastFeature := tt.featuresToAdd[len(tt.featuresToAdd)-1]
				retrieved := profile.GetFeature(lastFeature.name)
				suite.Require().NotNil(retrieved)
				suite.Assert().Equal(lastFeature.requestType, retrieved.GetRequestType())
			}
		})
	}
}

func (suite *ProfileTestSuite) TestProfile_SupportsFeature() {
	tests := []struct {
		name            string
		profileName     string
		features        []*TestFeature
		featureToCheck  string
		expectedSupport bool
	}{
		{
			name:        "feature exists",
			profileName: "test-profile",
			features: []*TestFeature{
				{
					name:         "Feature1",
					requestType:  reflect.TypeOf(""),
					responseType: reflect.TypeOf(0),
				},
			},
			featureToCheck:  "Feature1",
			expectedSupport: true,
		},
		{
			name:        "feature does not exist",
			profileName: "test-profile",
			features: []*TestFeature{
				{
					name:         "Feature1",
					requestType:  reflect.TypeOf(""),
					responseType: reflect.TypeOf(0),
				},
			},
			featureToCheck:  "NonexistentFeature",
			expectedSupport: false,
		},
		{
			name:            "empty profile",
			profileName:     "test-profile",
			features:        nil,
			featureToCheck:  "AnyFeature",
			expectedSupport: false,
		},
		{
			name:        "case sensitive check",
			profileName: "test-profile",
			features: []*TestFeature{
				{
					name:         "Feature1",
					requestType:  reflect.TypeOf(""),
					responseType: reflect.TypeOf(0),
				},
			},
			featureToCheck:  "feature1",
			expectedSupport: false,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			features := make([]Feature, len(tt.features))
			for i, f := range tt.features {
				features[i] = f
			}
			profile := NewProfile(tt.profileName, features...)
			suite.Assert().Equal(tt.expectedSupport, profile.SupportsFeature(tt.featureToCheck))
		})
	}
}

func (suite *ProfileTestSuite) TestProfile_GetFeature() {
	tests := []struct {
		name            string
		profileName     string
		features        []*TestFeature
		featureToGet    string
		expectedFeature *TestFeature
		shouldBeNil     bool
	}{
		{
			name:        "get existing feature",
			profileName: "test-profile",
			features: []*TestFeature{
				{
					name:         "Feature1",
					requestType:  reflect.TypeOf(""),
					responseType: reflect.TypeOf(0),
				},
			},
			featureToGet: "Feature1",
			expectedFeature: &TestFeature{
				name:         "Feature1",
				requestType:  reflect.TypeOf(""),
				responseType: reflect.TypeOf(0),
			},
			shouldBeNil: false,
		},
		{
			name:        "get nonexistent feature",
			profileName: "test-profile",
			features: []*TestFeature{
				{
					name:         "Feature1",
					requestType:  reflect.TypeOf(""),
					responseType: reflect.TypeOf(0),
				},
			},
			featureToGet:    "NonexistentFeature",
			expectedFeature: nil,
			shouldBeNil:     true,
		},
		{
			name:            "get from empty profile",
			profileName:     "test-profile",
			features:        nil,
			featureToGet:    "AnyFeature",
			expectedFeature: nil,
			shouldBeNil:     true,
		},
		{
			name:        "get from multiple features",
			profileName: "test-profile",
			features: []*TestFeature{
				{
					name:         "Feature1",
					requestType:  reflect.TypeOf(""),
					responseType: reflect.TypeOf(0),
				},
				{
					name:         "Feature2",
					requestType:  reflect.TypeOf(0),
					responseType: reflect.TypeOf(""),
				},
			},
			featureToGet: "Feature2",
			expectedFeature: &TestFeature{
				name:         "Feature2",
				requestType:  reflect.TypeOf(0),
				responseType: reflect.TypeOf(""),
			},
			shouldBeNil: false,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			features := make([]Feature, len(tt.features))
			for i, f := range tt.features {
				features[i] = f
			}
			profile := NewProfile(tt.profileName, features...)
			retrieved := profile.GetFeature(tt.featureToGet)

			if tt.shouldBeNil {
				suite.Assert().Nil(retrieved)
			} else {
				suite.Require().NotNil(retrieved)
				suite.Assert().Equal(tt.expectedFeature.name, retrieved.GetFeatureName())
				suite.Assert().Equal(tt.expectedFeature.requestType, retrieved.GetRequestType())
				suite.Assert().Equal(tt.expectedFeature.responseType, retrieved.GetResponseType())
			}
		})
	}
}

func (suite *ProfileTestSuite) TestProfile_ParseRequest() {
	tests := []struct {
		name            string
		profileName     string
		features        []*TestFeature
		featureName     string
		rawRequest      interface{}
		requestParser   func(raw interface{}, requestType reflect.Type) (Request, error)
		expectedRequest Request
		expectedError   string
	}{
		{
			name:        "successful parse",
			profileName: "test-profile",
			features: []*TestFeature{
				{
					name:         "TestFeature",
					requestType:  reflect.TypeOf(&TestRequest{}),
					responseType: reflect.TypeOf(&TestResponse{}),
				},
			},
			featureName: "TestFeature",
			rawRequest:  map[string]interface{}{"Value": "test"},
			requestParser: func(raw interface{}, requestType reflect.Type) (Request, error) {
				suite.Assert().Equal(reflect.TypeOf(&TestRequest{}), requestType)
				return &TestRequest{Value: "test"}, nil
			},
			expectedRequest: &TestRequest{Value: "test"},
			expectedError:   "",
		},
		{
			name:        "feature not found",
			profileName: "test-profile",
			features: []*TestFeature{
				{
					name:         "OtherFeature",
					requestType:  reflect.TypeOf(&TestRequest{}),
					responseType: reflect.TypeOf(&TestResponse{}),
				},
			},
			featureName: "NonexistentFeature",
			rawRequest:  map[string]interface{}{"Value": "test"},
			requestParser: func(raw interface{}, requestType reflect.Type) (Request, error) {
				return nil, nil
			},
			expectedRequest: nil,
			expectedError:   "Feature NonexistentFeature not found",
		},
		{
			name:        "empty profile",
			profileName: "test-profile",
			features:    nil,
			featureName: "AnyFeature",
			rawRequest:  map[string]interface{}{"Value": "test"},
			requestParser: func(raw interface{}, requestType reflect.Type) (Request, error) {
				return nil, nil
			},
			expectedRequest: nil,
			expectedError:   "Feature AnyFeature not found",
		},
		{
			name:        "parser returns error",
			profileName: "test-profile",
			features: []*TestFeature{
				{
					name:         "TestFeature",
					requestType:  reflect.TypeOf(&TestRequest{}),
					responseType: reflect.TypeOf(&TestResponse{}),
				},
			},
			featureName: "TestFeature",
			rawRequest:  map[string]interface{}{"Value": "test"},
			requestParser: func(raw interface{}, requestType reflect.Type) (Request, error) {
				return nil, fmt.Errorf("parse error")
			},
			expectedRequest: nil,
			expectedError:   "parse error",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			features := make([]Feature, len(tt.features))
			for i, f := range tt.features {
				features[i] = f
			}
			profile := NewProfile(tt.profileName, features...)
			request, err := profile.ParseRequest(tt.featureName, tt.rawRequest, tt.requestParser)

			if tt.expectedError != "" {
				suite.Require().Error(err)
				suite.Assert().Contains(err.Error(), tt.expectedError)
				suite.Assert().Nil(request)
			} else {
				suite.Require().NoError(err)
				suite.Assert().NotNil(request)
				if tt.expectedRequest != nil {
					suite.Assert().Equal(tt.expectedRequest.GetFeatureName(), request.GetFeatureName())
				}
			}
		})
	}
}

func (suite *ProfileTestSuite) TestProfile_ParseResponse() {
	tests := []struct {
		name             string
		profileName      string
		features         []*TestFeature
		featureName      string
		rawResponse      interface{}
		responseParser   func(raw interface{}, responseType reflect.Type) (Response, error)
		expectedResponse Response
		expectedError    string
	}{
		{
			name:        "successful parse",
			profileName: "test-profile",
			features: []*TestFeature{
				{
					name:         "TestFeature",
					requestType:  reflect.TypeOf(&TestRequest{}),
					responseType: reflect.TypeOf(&TestResponse{}),
				},
			},
			featureName: "TestFeature",
			rawResponse: map[string]interface{}{"Result": 42},
			responseParser: func(raw interface{}, responseType reflect.Type) (Response, error) {
				suite.Assert().Equal(reflect.TypeOf(&TestResponse{}), responseType)
				return &TestResponse{Result: 42}, nil
			},
			expectedResponse: &TestResponse{Result: 42},
			expectedError:    "",
		},
		{
			name:        "feature not found",
			profileName: "test-profile",
			features: []*TestFeature{
				{
					name:         "OtherFeature",
					requestType:  reflect.TypeOf(&TestRequest{}),
					responseType: reflect.TypeOf(&TestResponse{}),
				},
			},
			featureName: "NonexistentFeature",
			rawResponse: map[string]interface{}{"Result": 42},
			responseParser: func(raw interface{}, responseType reflect.Type) (Response, error) {
				return nil, nil
			},
			expectedResponse: nil,
			expectedError:    "Feature NonexistentFeature not found",
		},
		{
			name:        "empty profile",
			profileName: "test-profile",
			features:    nil,
			featureName: "AnyFeature",
			rawResponse: map[string]interface{}{"Result": 42},
			responseParser: func(raw interface{}, responseType reflect.Type) (Response, error) {
				return nil, nil
			},
			expectedResponse: nil,
			expectedError:    "Feature AnyFeature not found",
		},
		{
			name:        "parser returns error",
			profileName: "test-profile",
			features: []*TestFeature{
				{
					name:         "TestFeature",
					requestType:  reflect.TypeOf(&TestRequest{}),
					responseType: reflect.TypeOf(&TestResponse{}),
				},
			},
			featureName: "TestFeature",
			rawResponse: map[string]interface{}{"Result": 42},
			responseParser: func(raw interface{}, responseType reflect.Type) (Response, error) {
				return nil, fmt.Errorf("parse error")
			},
			expectedResponse: nil,
			expectedError:    "parse error",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			features := make([]Feature, len(tt.features))
			for i, f := range tt.features {
				features[i] = f
			}
			profile := NewProfile(tt.profileName, features...)
			response, err := profile.ParseResponse(tt.featureName, tt.rawResponse, tt.responseParser)

			if tt.expectedError != "" {
				suite.Require().Error(err)
				suite.Assert().Contains(err.Error(), tt.expectedError)
				suite.Assert().Nil(response)
			} else {
				suite.Require().NoError(err)
				suite.Assert().NotNil(response)
				if tt.expectedResponse != nil {
					suite.Assert().Equal(tt.expectedResponse.GetFeatureName(), response.GetFeatureName())
				}
			}
		})
	}
}
