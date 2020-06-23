package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/meter"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"time"
)

// Test
func (suite *OcppV2TestSuite) TestMeterValuesRequestValidation() {
	var requestTable = []GenericTestEntry{
		{meter.MeterValuesRequest{EvseID: 1, MeterValue: []types.MeterValue{ {Timestamp: types.DateTime{Time: time.Now()}, SampledValue: []types.SampledValue{ {Value: 3.14, Context: types.ReadingContextTransactionEnd, Measurand: types.MeasurandPowerActiveExport, Phase: types.PhaseL2, Location: types.LocationBody}}} }}, true},
		{meter.MeterValuesRequest{MeterValue: []types.MeterValue{ {Timestamp: types.DateTime{Time: time.Now()}, SampledValue: []types.SampledValue{ {Value: 3.14, Context: types.ReadingContextTransactionEnd, Measurand: types.MeasurandPowerActiveExport, Phase: types.PhaseL2, Location: types.LocationBody}}} }}, true},
		{meter.MeterValuesRequest{EvseID: 1, MeterValue: []types.MeterValue{}}, false},
		{meter.MeterValuesRequest{EvseID: 1}, false},
		{meter.MeterValuesRequest{EvseID: 1, MeterValue: []types.MeterValue{ {Timestamp: types.DateTime{Time: time.Now()}, SampledValue: []types.SampledValue{ {Value: 3.14, Context: "invalidContext", Measurand: types.MeasurandPowerActiveExport, Phase: types.PhaseL2, Location: types.LocationBody}}} }}, false},
		{meter.MeterValuesRequest{EvseID: -1, MeterValue: []types.MeterValue{ {Timestamp: types.DateTime{Time: time.Now()}, SampledValue: []types.SampledValue{ {Value: 3.14, Context: types.ReadingContextTransactionEnd, Measurand: types.MeasurandPowerActiveExport, Phase: types.PhaseL2, Location: types.LocationBody}}} }}, false},
	}
	ExecuteGenericTestTable(suite.T(), requestTable)
}

func (suite *OcppV2TestSuite) TestMeterValuesConfirmationValidation() {
	var responseTable = []GenericTestEntry{
		{meter.MeterValuesResponse{}, true},
	}
	ExecuteGenericTestTable(suite.T(), responseTable)
}

func (suite *OcppV2TestSuite) TestMeterValuesE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	evseId := 1
	signedMeterValue := types.SignedMeterValue{SignedMeterData: "0xdeadbeef", SigningMethod: "ECDSAP256SHA256", EncodingMethod: "DLMS Message", PublicKey: "0xd34dc0de"}
	unitOfMeasure := types.UnitOfMeasure{Unit: "kW", Multiplier: newInt(0)}
	sampledValue := types.SampledValue{Value: 3.14, Context: types.ReadingContextTransactionEnd, Measurand: types.MeasurandPowerActiveExport, Phase: types.PhaseL2, Location: types.LocationBody, SignedMeterValue: &signedMeterValue, UnitOfMeasure: &unitOfMeasure}
	sampledValues := []types.SampledValue{sampledValue}
	meterValue := types.MeterValue{Timestamp: types.DateTime{Time: time.Now()}, SampledValue: sampledValues}
	meterValues := []types.MeterValue{meterValue}
	timestamp := types.DateTime{Time: time.Now()}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"evseId":%v,"meterValue":[{"timestamp":"%v","sampledValue":[{"value":%v,"context":"%v","measurand":"%v","phase":"%v","location":"%v","signedMeterValue":{"signedMeterData":"%v","signingMethod":"%v","encodingMethod":"%v","publicKey":"%v"},"unitOfMeasure":{"unit":"%v","multiplier":%v}}]}]}]`,
		messageId, meter.MeterValuesFeatureName, evseId, timestamp.FormatTimestamp(), sampledValue.Value, sampledValue.Context, sampledValue.Measurand, sampledValue.Phase, sampledValue.Location, signedMeterValue.SignedMeterData, signedMeterValue.SigningMethod, signedMeterValue.EncodingMethod, signedMeterValue.PublicKey, unitOfMeasure.Unit, *unitOfMeasure.Multiplier)
	responseJson := fmt.Sprintf(`[3,"%v",{}]`, messageId)
	response := meter.NewMeterValuesResponse()
	channel := NewMockWebSocket(wsId)

	handler := MockCSMSMeterHandler{}
	handler.On("OnMeterValues", mock.AnythingOfType("string"), mock.Anything).Return(response, nil).Run(func(args mock.Arguments) {
		request := args.Get(1).(*meter.MeterValuesRequest)
		assert.Equal(t, evseId, request.EvseID)
		require.Len(t, request.MeterValue, len(meterValues))
		assertDateTimeEquality(t, &meterValue.Timestamp, &request.MeterValue[0].Timestamp)
		require.Len(t, request.MeterValue[0].SampledValue, len(sampledValues))
		assert.Equal(t, sampledValue.Value, request.MeterValue[0].SampledValue[0].Value)
		assert.Equal(t, sampledValue.Context, request.MeterValue[0].SampledValue[0].Context)
		assert.Equal(t, sampledValue.Measurand, request.MeterValue[0].SampledValue[0].Measurand)
		assert.Equal(t, sampledValue.Phase, request.MeterValue[0].SampledValue[0].Phase)
		assert.Equal(t, sampledValue.Location, request.MeterValue[0].SampledValue[0].Location)
		require.NotNil(t, request.MeterValue[0].SampledValue[0].SignedMeterValue)
		assert.Equal(t, signedMeterValue.SignedMeterData, request.MeterValue[0].SampledValue[0].SignedMeterValue.SignedMeterData)
		assert.Equal(t, signedMeterValue.SigningMethod, request.MeterValue[0].SampledValue[0].SignedMeterValue.SigningMethod)
		assert.Equal(t, signedMeterValue.EncodingMethod, request.MeterValue[0].SampledValue[0].SignedMeterValue.EncodingMethod)
		assert.Equal(t, signedMeterValue.PublicKey, request.MeterValue[0].SampledValue[0].SignedMeterValue.PublicKey)
		require.NotNil(t, request.MeterValue[0].SampledValue[0].UnitOfMeasure)
		assert.Equal(t, unitOfMeasure.Unit, request.MeterValue[0].SampledValue[0].UnitOfMeasure.Unit)
		assert.Equal(t, *unitOfMeasure.Multiplier, *request.MeterValue[0].SampledValue[0].UnitOfMeasure.Multiplier)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	assert.Nil(t, err)
	r, err := suite.chargingStation.MeterValues(evseId, meterValues)
	assert.Nil(t, err)
	assert.NotNil(t, r)
}

func (suite *OcppV2TestSuite) TestMeterValuesInvalidEndpoint() {
	messageId := defaultMessageId
	connectorId := 1
	evseId := 1
	signedMeterValue := types.SignedMeterValue{SignedMeterData: "0xdeadbeef", SigningMethod: "ECDSAP256SHA256", EncodingMethod: "DLMS Message", PublicKey: "0xd34dc0de"}
	unitOfMeasure := types.UnitOfMeasure{Unit: "kW", Multiplier: newInt(0)}
	sampledValue := types.SampledValue{Value: 3.14, Context: types.ReadingContextTransactionEnd, Measurand: types.MeasurandPowerActiveExport, Phase: types.PhaseL2, Location: types.LocationBody, SignedMeterValue: &signedMeterValue, UnitOfMeasure: &unitOfMeasure}
	sampledValues := []types.SampledValue{sampledValue}
	meterValue := types.MeterValue{Timestamp: types.DateTime{Time: time.Now()}, SampledValue: sampledValues}
	meterValues := []types.MeterValue{meterValue}
	timestamp := types.DateTime{Time: time.Now()}
	meterValuesRequest := meter.NewMeterValuesRequest(connectorId, meterValues)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"evseId":%v,"meterValue":[{"timestamp":"%v","sampledValue":[{"value":%v,"context":"%v","measurand":"%v","phase":"%v","location":"%v","signedMeterValue":{"signedMeterData":"%v","signingMethod":"%v","encodingMethod":"%v","publicKey":"%v"},"unitOfMeasure":{"unit":"%v","multiplier":%v}}]}]}]`,
		messageId, meter.MeterValuesFeatureName, evseId, timestamp.FormatTimestamp(), sampledValue.Value, sampledValue.Context, sampledValue.Measurand, sampledValue.Phase, sampledValue.Location, signedMeterValue.SignedMeterData, signedMeterValue.SigningMethod, signedMeterValue.EncodingMethod, signedMeterValue.PublicKey, unitOfMeasure.Unit, *unitOfMeasure.Multiplier)
	testUnsupportedRequestFromCentralSystem(suite, meterValuesRequest, requestJson, messageId)
}
