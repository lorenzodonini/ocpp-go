package ocppj_test

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/lorenzodonini/ocpp-go/ws"

	"github.com/stretchr/testify/mock"

	"github.com/lorenzodonini/ocpp-go/ocpp"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocppj"
	"github.com/stretchr/testify/suite"
)

// --------- Custom types ---------

type GetConfigurationRequest struct {
	RequiredKeys []string `json:"requiredKeys,omitempty"`
	OptionalKeys []string `json:"optionalKeys,omitempty"`
}

func (r *GetConfigurationRequest) GetOcppType() reflect.Type {
	return reflect.TypeOf(core.GetConfigurationRequest{})
}

func (r *GetConfigurationRequest) GetFeatureName() string {
	return core.GetConfigurationFeatureName
}

func (r *GetConfigurationRequest) Parse(targetPayload ocpp.Request) error {
	target, ok := targetPayload.(*core.GetConfigurationRequest)
	if !ok {
		return fmt.Errorf("invalid type %T", targetPayload)
	}
	target.Key = make([]string, len(r.RequiredKeys)+len(r.OptionalKeys))
	i := 0
	for _, key := range r.RequiredKeys {
		target.Key[i] = key
		i++
	}
	for _, key := range r.OptionalKeys {
		target.Key[i] = key
		i++
	}
	return nil
}

func (r *GetConfigurationRequest) Serialize(srcPayload ocpp.Request) error {
	src, ok := srcPayload.(*core.GetConfigurationRequest)
	if !ok {
		return fmt.Errorf("invalid type %T", srcPayload)
	}
	r.RequiredKeys = []string{}
	r.OptionalKeys = []string{}
	for _, key := range src.Key {
		if strings.HasPrefix(key, "opt_") {
			r.OptionalKeys = append(r.OptionalKeys, key)
			continue
		} else {
			r.RequiredKeys = append(r.RequiredKeys, key)
		}
	}
	return nil
}

type ConfigurationKey struct {
	Key      string  `json:"key"`
	Readonly string  `json:"readonly"`
	Value    *string `json:"value,omitempty"`
}

type GetConfigurationConfirmation struct {
	ConfigurationKey []ConfigurationKey `json:"configurationKey,omitempty"`
	UnknownKey       []string           `json:"unknownKey,omitempty"`
}

func (c *GetConfigurationConfirmation) GetOcppType() reflect.Type {
	return reflect.TypeOf(core.GetConfigurationConfirmation{})
}

func (c *GetConfigurationConfirmation) GetFeatureName() string {
	return core.GetConfigurationFeatureName
}

func (c *GetConfigurationConfirmation) Parse(targetPayload ocpp.Response) error {
	target, ok := targetPayload.(*core.GetConfigurationConfirmation)
	if !ok {
		return fmt.Errorf("invalid type %T", targetPayload)
	}
	target.UnknownKey = c.UnknownKey
	target.ConfigurationKey = make([]core.ConfigurationKey, len(c.ConfigurationKey))
	for i, key := range c.ConfigurationKey {
		boolVal, err := strconv.ParseBool(key.Readonly)
		if err != nil {
			return err
		}
		target.ConfigurationKey[i] = core.ConfigurationKey{
			Key:      key.Key,
			Readonly: boolVal,
			Value:    key.Value,
		}
	}
	return nil
}

func (c *GetConfigurationConfirmation) Serialize(srcPayload ocpp.Response) error {
	src, ok := srcPayload.(*core.GetConfigurationConfirmation)
	if !ok {
		return fmt.Errorf("invalid type %T", srcPayload)
	}
	c.UnknownKey = src.UnknownKey
	c.ConfigurationKey = make([]ConfigurationKey, len(src.ConfigurationKey))
	for i, key := range src.ConfigurationKey {
		c.ConfigurationKey[i] = ConfigurationKey{
			Key:      key.Key,
			Readonly: fmt.Sprintf("%t", key.Readonly),
			Value:    key.Value,
		}
	}
	return nil
}

// --------- Test suite ---------

type CustomTypeTestSuite struct {
	suite.Suite
	mapper             ocppj.CustomTypeMapper
	chargePoint        *ocppj.Client
	centralSystem      *ocppj.Server
	mockServer         *MockWebsocketServer
	mockClient         *MockWebsocketClient
	clientDispatcher   ocppj.ClientDispatcher
	serverDispatcher   ocppj.ServerDispatcher
	clientRequestQueue ocppj.RequestQueue
	serverRequestMap   ocppj.ServerQueueMap
}

var defaultMessageID = "12345"

func (s *CustomTypeTestSuite) SetupTest() {
	mockProfile := ocpp.NewProfile("mock", &MockFeature{})
	mockClient := MockWebsocketClient{}
	mockServer := MockWebsocketServer{}
	s.mockClient = &mockClient
	s.mockServer = &mockServer
	s.clientRequestQueue = ocppj.NewFIFOClientQueue(queueCapacity)
	s.clientDispatcher = ocppj.NewDefaultClientDispatcher(s.clientRequestQueue)
	s.chargePoint = ocppj.NewClient("mock_id", s.mockClient, s.clientDispatcher, nil, mockProfile)
	s.serverRequestMap = ocppj.NewFIFOQueueMap(queueCapacity)
	s.serverDispatcher = ocppj.NewDefaultServerDispatcher(s.serverRequestMap)
	s.centralSystem = ocppj.NewServer(s.mockServer, s.serverDispatcher, nil, mockProfile)
	defaultDialect := ocpp.V16 // set default to version 1.6 format error *for test only
	s.centralSystem.SetDialect(defaultDialect)
	s.chargePoint.SetDialect(defaultDialect)
	s.mapper = ocppj.NewCustomTypeMapper()
	ocppj.SetMessageIdGenerator(func() string {
		return defaultMessageID
	})
}

func (s *CustomTypeTestSuite) TestGetCustomRequest() {
	// Initially, no custom request is registered
	customReq, ok := s.mapper.GetCustomRequest(core.GetConfigurationFeatureName)
	s.False(ok)
	s.Nil(customReq)
	// Register custom request and check getters
	s.mapper.SetCustomRequest(&GetConfigurationRequest{})
	customReq, ok = s.mapper.GetCustomRequest(core.GetConfigurationFeatureName)
	s.True(ok)
	s.NotNil(customReq)
	s.Equal(core.GetConfigurationFeatureName, customReq.GetFeatureName())
}

func (s *CustomTypeTestSuite) TestGetCustomResponse() {
	// Initially, no custom response is registered
	customResp, ok := s.mapper.GetCustomResponse(core.GetConfigurationFeatureName)
	s.False(ok)
	s.Nil(customResp)
	// Register custom response and check getters
	s.mapper.SetCustomResponse(&GetConfigurationConfirmation{})
	customResp, ok = s.mapper.GetCustomResponse(core.GetConfigurationFeatureName)
	s.True(ok)
	s.NotNil(customResp)
	s.Equal(core.GetConfigurationFeatureName, customResp.GetFeatureName())
}

func (s *CustomTypeTestSuite) TestGetConfigurationRequest() {
	// Register custom request
	s.mapper.SetCustomRequest(&GetConfigurationRequest{})
	s.chargePoint.SetCustomTypeMapper(s.mapper)
	s.chargePoint.Profiles = append(s.chargePoint.Profiles, core.Profile)
	s.centralSystem.Profiles = append(s.centralSystem.Profiles, core.Profile)
	// Create mock message
	mockChargePointID := "1234"
	requiredKeys := []string{"key1", "key2"}
	optionalKeys := []string{"opt_key3"}
	expectedJson := fmt.Sprintf(`[%v,"%v","%v",{"requiredKeys":["%v","%v"],"optionalKeys":["%v"]}]`,
		ocppj.CALL, defaultMessageID, core.GetConfigurationFeatureName, requiredKeys[0], requiredKeys[1], optionalKeys[0])
	// Mock lower layer calls to simulate E2E behavior
	msgC := make(chan []byte, 1)
	s.mockServer.On("Write", mock.AnythingOfType("string"), mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		msg, ok := args.Get(1).([]byte)
		s.True(ok)
		msgC <- msg
	})
	s.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	s.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	s.centralSystem.Start(8887, "somePath")
	_ = s.chargePoint.Start("someUrl")
	// Register mock client
	s.serverDispatcher.CreateClient(mockChargePointID)
	s.centralSystem.SetClientCustomTypeMapper(mockChargePointID, s.mapper)
	// Test sending GetConfigurationRequest with custom type
	r := core.NewGetConfigurationRequest(append(requiredKeys, optionalKeys...))
	err := s.centralSystem.SendRequest(mockChargePointID, r)
	s.NoError(err)
	// Await raw message from websocket
	outMsg := <-msgC
	s.Equal([]byte(expectedJson), outMsg)
	// Test parsing on charge point
	s.chargePoint.SetRequestHandler(func(request ocpp.Request, requestId string, action string) {
		s.Equal(defaultMessageID, requestId)
		s.Equal(core.GetConfigurationFeatureName, action)
		s.IsType(&core.GetConfigurationRequest{}, request)
		req, ok := request.(*core.GetConfigurationRequest)
		s.True(ok)
		s.Len(req.Key, len(requiredKeys)+len(optionalKeys))
		i := 0
		// Verify keys are complete and in the correct order
		for _, key := range requiredKeys {
			s.Equal(key, req.Key[i])
			i++
		}
		for _, key := range optionalKeys {
			s.Equal(key, req.Key[i])
			i++
		}
	})
	err = s.mockClient.MessageHandler(outMsg)
	s.NoError(err)
}

func (s *CustomTypeTestSuite) TestGetConfigurationResponse() {
	// Register custom request
	s.mapper.SetCustomResponse(&GetConfigurationConfirmation{})
	s.chargePoint.SetCustomTypeMapper(s.mapper)
	s.chargePoint.Profiles = append(s.chargePoint.Profiles, core.Profile)
	s.centralSystem.Profiles = append(s.centralSystem.Profiles, core.Profile)
	// Create mock response
	mockChargePointID := "1234"
	// Mock lower layer calls to simulate E2E behavior
	msgC := make(chan []byte, 1)
	s.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil)
	s.mockClient.On("Write", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		msg, ok := args.Get(0).([]byte)
		s.True(ok)
		msgC <- msg
	})
	s.mockServer.On("Start", mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil)
	s.centralSystem.Start(8887, "somePath")
	_ = s.chargePoint.Start("someUrl")
	// Register mock client
	s.serverDispatcher.CreateClient(mockChargePointID)
	s.centralSystem.SetClientCustomTypeMapper(mockChargePointID, s.mapper)
	// Register pending request
	rq := core.NewGetConfigurationRequest([]string{"key1", "key2", "opt_key3"})
	s.centralSystem.RequestState.AddPendingRequest(mockChargePointID, defaultMessageID, rq)
	// Test sending GetConfigurationConfirmation with custom type
	dummyVal := "someValue"
	rs := core.NewGetConfigurationConfirmation([]core.ConfigurationKey{
		{Key: "key1", Readonly: false, Value: &dummyVal},
		{Key: "opt_key3", Readonly: true, Value: &dummyVal},
	})
	rs.UnknownKey = []string{"key2"}
	expectedJson := fmt.Sprintf(`[%v,"%v",{"configurationKey":[{"key":"%v","readonly":"%v","value":"%v"},{"key":"%v","readonly":"%v","value":"%v"}],"unknownKey":["%v"]}]`,
		ocppj.CALL_RESULT, defaultMessageID,
		rs.ConfigurationKey[0].Key, rs.ConfigurationKey[0].Readonly, *rs.ConfigurationKey[0].Value,
		rs.ConfigurationKey[1].Key, rs.ConfigurationKey[1].Readonly, *rs.ConfigurationKey[1].Value,
		rs.UnknownKey[0])
	err := s.chargePoint.SendResponse(defaultMessageID, rs)
	s.NoError(err)
	// Await raw message from websocket
	outMsg := <-msgC
	s.Equal([]byte(expectedJson), outMsg)
	// Test parsing on central system
	s.centralSystem.SetResponseHandler(func(client ws.Channel, response ocpp.Response, requestId string) {
		s.Equal(defaultMessageID, requestId)
		s.Equal(core.GetConfigurationFeatureName, response.GetFeatureName())
		s.IsType(&core.GetConfigurationConfirmation{}, response)
		resp, ok := response.(*core.GetConfigurationConfirmation)
		s.True(ok)
		// Verify contents
		s.Len(resp.ConfigurationKey, 2)
		s.Equal(rs.ConfigurationKey[0].Key, resp.ConfigurationKey[0].Key)
		s.Equal(rs.ConfigurationKey[0].Readonly, resp.ConfigurationKey[0].Readonly)
		s.Equal(*rs.ConfigurationKey[0].Value, *resp.ConfigurationKey[0].Value)
		s.Equal(rs.ConfigurationKey[1].Key, resp.ConfigurationKey[1].Key)
		s.Equal(rs.ConfigurationKey[1].Readonly, resp.ConfigurationKey[1].Readonly)
		s.Equal(*rs.ConfigurationKey[1].Value, *resp.ConfigurationKey[1].Value)
		s.Len(resp.UnknownKey, 1)
		s.Equal(rs.UnknownKey[0], resp.UnknownKey[0])
	})
	channel := NewMockWebSocket(mockChargePointID)
	err = s.mockServer.MessageHandler(channel, outMsg)
	s.NoError(err)
}
