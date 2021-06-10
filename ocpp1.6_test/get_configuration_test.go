package ocpp16_test

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test
func (suite *OcppV16TestSuite) TestGetConfigurationRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{core.GetConfigurationRequest{Key: []string{"key1", "key2"}}, true},
		{core.GetConfigurationRequest{Key: []string{"key1", "key2", "key3", "key4", "key5", "key6"}}, true},
		{core.GetConfigurationRequest{Key: []string{"key1", "key2", "key2"}}, false},
		{core.GetConfigurationRequest{}, true},
		{core.GetConfigurationRequest{Key: []string{}}, true},
		{core.GetConfigurationRequest{Key: []string{">50................................................"}}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV16TestSuite) TestGetConfigurationConfirmationValidation() {
	t := suite.T()
	value1 := "value1"
	value2 := "value2"
	longValue := ">500................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................."
	var confirmationTable = []GenericTestEntry{
		{core.GetConfigurationConfirmation{ConfigurationKey: []core.ConfigurationKey{{Key: "key1", Readonly: true, Value: &value1}}}, true},
		{core.GetConfigurationConfirmation{ConfigurationKey: []core.ConfigurationKey{{Key: "key1", Readonly: true, Value: &value1}, {Key: "key2", Readonly: false, Value: &value2}}}, true},
		{core.GetConfigurationConfirmation{ConfigurationKey: []core.ConfigurationKey{{Key: "key1", Readonly: true, Value: &value1}}, UnknownKey: []string{"keyX"}}, true},
		{core.GetConfigurationConfirmation{ConfigurationKey: []core.ConfigurationKey{{Key: "key1", Readonly: false, Value: &value1}}, UnknownKey: []string{"keyX", "keyY"}}, true},
		{core.GetConfigurationConfirmation{UnknownKey: []string{"keyX"}}, true},
		{core.GetConfigurationConfirmation{UnknownKey: []string{">50................................................"}}, false},
		{core.GetConfigurationConfirmation{ConfigurationKey: []core.ConfigurationKey{{Key: ">50................................................", Readonly: true, Value: &value1}}}, false},
		{core.GetConfigurationConfirmation{ConfigurationKey: []core.ConfigurationKey{{Key: "key1", Readonly: true, Value: &longValue}}}, false},
		//{ocpp16.GetConfigurationConfirmation{ConfigurationKey: []ocpp16.ConfigurationKey{{Key: "key1", Readonly: true, Value: "value1"}, {Key: "key1", Readonly: false, Value: "value2"}}}, false},
	}
	//TODO: additional test cases TBD. See get_configuration.go
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV16TestSuite) TestGetConfigurationE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	key1 := "key1"
	key2 := "key2"
	someValue := "someValue"
	someOtherValue := "someOtherValue"
	requestKeys := []string{key1, key2}
	resultKey1 := core.ConfigurationKey{Key: key1, Readonly: true, Value: &someValue}
	resultKey2 := core.ConfigurationKey{Key: key1, Readonly: true, Value: &someOtherValue}
	resultKeys := []core.ConfigurationKey{resultKey1, resultKey2}
	unknownKeys := []string{"keyX", "keyY"}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"key":["%v","%v"]}]`, messageId, core.GetConfigurationFeatureName, key1, key2)
	responseJson := fmt.Sprintf(`[3,"%v",{"configurationKey":[{"key":"%v","readonly":%v,"value":"%v"},{"key":"%v","readonly":%v,"value":"%v"}],"unknownKey":["%v","%v"]}]`, messageId, resultKey1.Key, resultKey1.Readonly, *resultKey1.Value, resultKey2.Key, resultKey2.Readonly, *resultKey2.Value, unknownKeys[0], unknownKeys[1])
	getConfigurationConfirmation := core.NewGetConfigurationConfirmation(resultKeys)
	getConfigurationConfirmation.UnknownKey = unknownKeys
	channel := NewMockWebSocket(wsId)

	coreListener := MockChargePointCoreListener{}
	coreListener.On("OnGetConfiguration", mock.Anything).Return(getConfigurationConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*core.GetConfigurationRequest)
		require.NotNil(t, request)
		require.True(t, ok)
		require.Len(t, request.Key, len(requestKeys))
		assert.Equal(t, requestKeys, request.Key)
	})
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, coreListener, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	assert.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.centralSystem.GetConfiguration(wsId, func(confirmation *core.GetConfigurationConfirmation, err error) {
		require.Nil(t, err)
		require.NotNil(t, confirmation)
		assert.Equal(t, unknownKeys, confirmation.UnknownKey)
		assert.Equal(t, resultKeys, confirmation.ConfigurationKey)
		resultChannel <- true
	}, requestKeys)
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV16TestSuite) TestGetConfigurationInvalidEndpoint() {
	messageId := defaultMessageId
	key1 := "key1"
	key2 := "key2"
	requestKeys := []string{key1, key2}
	getConfigurationRequest := core.NewGetConfigurationRequest(requestKeys)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"key":["%v","%v"]}]`, messageId, core.GetConfigurationFeatureName, key1, key2)
	testUnsupportedRequestFromChargePoint(suite, getConfigurationRequest, requestJson, messageId)
}
