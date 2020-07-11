package ocpp16_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"time"
)

// Test
func (suite *OcppV16TestSuite) TestMeterValuesRequestValidation() {
	var requestTable = []GenericTestEntry{
		{core.MeterValuesRequest{ConnectorId: 1, TransactionId: newInt(1), MeterValue: []types.MeterValue{{Timestamp: types.NewDateTime(time.Now()), SampledValue: []types.SampledValue{{Value: "value"}}}}}, true},
		{core.MeterValuesRequest{ConnectorId: 1, MeterValue: []types.MeterValue{{Timestamp: types.NewDateTime(time.Now()), SampledValue: []types.SampledValue{{Value: "value"}}}}}, true},
		{core.MeterValuesRequest{MeterValue: []types.MeterValue{{Timestamp: types.NewDateTime(time.Now()), SampledValue: []types.SampledValue{{Value: "value"}}}}}, true},
		{core.MeterValuesRequest{ConnectorId: -1, MeterValue: []types.MeterValue{{Timestamp: types.NewDateTime(time.Now()), SampledValue: []types.SampledValue{{Value: "value"}}}}}, false},
		{core.MeterValuesRequest{ConnectorId: 1, MeterValue: []types.MeterValue{}}, false},
		{core.MeterValuesRequest{ConnectorId: 1}, false},
		{core.MeterValuesRequest{ConnectorId: 1, MeterValue: []types.MeterValue{{Timestamp: types.NewDateTime(time.Now()), SampledValue: []types.SampledValue{}}}}, false},
	}
	ExecuteGenericTestTable(suite.T(), requestTable)
}

func (suite *OcppV16TestSuite) TestMeterValuesConfirmationValidation() {
	var confirmationTable = []GenericTestEntry{
		{core.MeterValuesConfirmation{}, true},
	}
	ExecuteGenericTestTable(suite.T(), confirmationTable)
}

func (suite *OcppV16TestSuite) TestMeterValuesE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	connectorId := 1
	mockValue := "value"
	mockUnit := types.UnitOfMeasureKW
	meterValues := []types.MeterValue{{Timestamp: types.NewDateTime(time.Now()), SampledValue: []types.SampledValue{{Value: mockValue, Unit: mockUnit}}}}
	timestamp := types.DateTime{Time: time.Now()}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"connectorId":%v,"meterValue":[{"timestamp":"%v","sampledValue":[{"value":"%v","unit":"%v"}]}]}]`, messageId, core.MeterValuesFeatureName, connectorId, timestamp.FormatTimestamp(), mockValue, mockUnit)
	responseJson := fmt.Sprintf(`[3,"%v",{}]`, messageId)
	meterValuesConfirmation := core.NewMeterValuesConfirmation()
	channel := NewMockWebSocket(wsId)

	coreListener := MockCentralSystemCoreListener{}
	coreListener.On("OnMeterValues", mock.AnythingOfType("string"), mock.Anything).Return(meterValuesConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*core.MeterValuesRequest)
		require.NotNil(t, request)
		require.True(t, ok)
		assert.Equal(t, connectorId, request.ConnectorId)
		require.Equal(t, 1, len(request.MeterValue))
		mv := request.MeterValue[0]
		assertDateTimeEquality(t, timestamp, *mv.Timestamp)
		require.Equal(t, 1, len(mv.SampledValue))
		sv := mv.SampledValue[0]
		assert.Equal(t, mockValue, sv.Value)
		assert.Equal(t, mockUnit, sv.Unit)
	})
	setupDefaultCentralSystemHandlers(suite, coreListener, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	require.Nil(t, err)
	confirmation, err := suite.chargePoint.MeterValues(connectorId, meterValues)
	require.Nil(t, err)
	require.NotNil(t, confirmation)
}

func (suite *OcppV16TestSuite) TestMeterValuesInvalidEndpoint() {
	messageId := defaultMessageId
	connectorId := 1
	mockValue := "value"
	mockUnit := types.UnitOfMeasureKW
	timestamp := types.DateTime{Time: time.Now()}
	meterValues := []types.MeterValue{{Timestamp: types.NewDateTime(time.Now()), SampledValue: []types.SampledValue{{Value: mockValue, Unit: mockUnit}}}}
	meterValuesRequest := core.NewMeterValuesRequest(connectorId, meterValues)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"connectorId":%v,"meterValue":[{"timestamp":"%v","sampledValue":[{"value":"%v","unit":"%v"}]}]}]`, messageId, core.MeterValuesFeatureName, connectorId, timestamp.FormatTimestamp(), mockValue, mockUnit)
	testUnsupportedRequestFromCentralSystem(suite, meterValuesRequest, requestJson, messageId)
}
