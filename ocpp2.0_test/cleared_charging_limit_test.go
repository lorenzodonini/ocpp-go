package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Tests
func (suite *OcppV2TestSuite) TestClearedChargingLimitRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{ocpp2.ClearedChargingLimitRequest{ChargingLimitSource: ocpp2.ChargingLimitSourceEMS, EvseID: newInt(0)}, true},
		{ocpp2.ClearedChargingLimitRequest{ChargingLimitSource: ocpp2.ChargingLimitSourceEMS}, true},
		{ocpp2.ClearedChargingLimitRequest{}, false},
		{ocpp2.ClearedChargingLimitRequest{ChargingLimitSource: ocpp2.ChargingLimitSourceEMS, EvseID: newInt(-1)}, false},
		{ocpp2.ClearedChargingLimitRequest{ChargingLimitSource: "invalidChargingLimitSource"}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestClearedChargingLimitConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{ocpp2.ClearedChargingLimitConfirmation{}, true},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestClearedChargingLimitE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := "1234"
	wsUrl := "someUrl"
	chargingLimitSource := ocpp2.ChargingLimitSourceEMS
	evseID := 42
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"chargingLimitSource":"%v","evseId":%v}]`, messageId, ocpp2.ClearedChargingLimitFeatureName, chargingLimitSource, evseID)
	responseJson := fmt.Sprintf(`[3,"%v",{}]`, messageId)
	clearedChargingLimitConfirmation := ocpp2.NewClearedChargingLimitConfirmation()
	channel := NewMockWebSocket(wsId)

	coreListener := MockCentralSystemCoreListener{}
	coreListener.On("OnClearedChargingLimit", mock.AnythingOfType("string"), mock.Anything).Return(clearedChargingLimitConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*ocpp2.ClearedChargingLimitRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		assert.Equal(t, chargingLimitSource, request.ChargingLimitSource)
		assert.Equal(t, evseID, *request.EvseID)
	})
	setupDefaultCentralSystemHandlers(suite, coreListener, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run test
	suite.csms.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	require.Nil(t, err)
	confirmation, err := suite.chargePoint.ClearChargingLimit(chargingLimitSource, func(request *ocpp2.ClearedChargingLimitRequest) {
		request.EvseID = newInt(evseID)
	})
	require.Nil(t, err)
	require.NotNil(t, confirmation)
}

func (suite *OcppV2TestSuite) TestClearedChargingLimitInvalidEndpoint() {
	messageId := defaultMessageId
	chargingLimitSource := ocpp2.ChargingLimitSourceEMS
	evseID := 42
	clearedChargingLimitRequest := ocpp2.NewClearedChargingLimitRequest(chargingLimitSource)
	clearedChargingLimitRequest.EvseID = newInt(evseID)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"chargingLimitSource":"%v","evseId":%v}]`, messageId, ocpp2.ClearedChargingLimitFeatureName, chargingLimitSource, evseID)
	testUnsupportedRequestFromCentralSystem(suite, clearedChargingLimitRequest, requestJson, messageId)
}
