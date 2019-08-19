package ocpp16_test

import (
	"fmt"
	"github.com/lorenzodonini/go-ocpp/ocpp1.6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"time"
)

// Test
func (suite *OcppV16TestSuite) TestMeterValuesRequestValidation() {
	var requestTable = []RequestTestEntry{
		{ocpp16.MeterValuesRequest{ConnectorId: 1, TransactionId: 1, MeterValue: []ocpp16.MeterValue{{Timestamp: ocpp16.NewDateTime(time.Now()), SampledValue: []ocpp16.SampledValue{{Value: "value"}}}}}, true},
		{ocpp16.MeterValuesRequest{ConnectorId: 1, MeterValue: []ocpp16.MeterValue{{Timestamp: ocpp16.NewDateTime(time.Now()), SampledValue: []ocpp16.SampledValue{{Value: "value"}}}}}, true},
		{ocpp16.MeterValuesRequest{MeterValue: []ocpp16.MeterValue{{Timestamp: ocpp16.NewDateTime(time.Now()), SampledValue: []ocpp16.SampledValue{{Value: "value"}}}}}, true},
		{ocpp16.MeterValuesRequest{ConnectorId: -1, MeterValue: []ocpp16.MeterValue{{Timestamp: ocpp16.NewDateTime(time.Now()), SampledValue: []ocpp16.SampledValue{{Value: "value"}}}}}, false},
		{ocpp16.MeterValuesRequest{ConnectorId: 1, MeterValue: []ocpp16.MeterValue{}}, false},
		{ocpp16.MeterValuesRequest{ConnectorId: 1}, false},
		{ocpp16.MeterValuesRequest{ConnectorId: 1, MeterValue: []ocpp16.MeterValue{{Timestamp: ocpp16.NewDateTime(time.Now()), SampledValue: []ocpp16.SampledValue{}}}}, false},
	}
	ExecuteRequestTestTable(suite.T(), requestTable)
}

func (suite *OcppV16TestSuite) TestMeterValuesConfirmationValidation() {
	var confirmationTable = []ConfirmationTestEntry{
		{ocpp16.MeterValuesConfirmation{}, true},
	}
	ExecuteConfirmationTestTable(suite.T(), confirmationTable)
}

func (suite *OcppV16TestSuite) TestMeterValuesE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	connectorId := 1
	mockValue := "value"
	mockUnit := ocpp16.UnitOfMeasureKW
	meterValues := []ocpp16.MeterValue{{Timestamp: ocpp16.NewDateTime(time.Now()), SampledValue: []ocpp16.SampledValue{{Value: mockValue, Unit: mockUnit}}}}
	timestamp := ocpp16.DateTime{Time: time.Now()}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"connectorId":%v,"meterValue":[{"timestamp":"%v","sampledValue":[{"value":"%v","unit":"%v"}]}]}]`, messageId, ocpp16.MeterValuesFeatureName, connectorId, timestamp.Format(ocpp16.ISO8601), mockValue, mockUnit)
	responseJson := fmt.Sprintf(`[3,"%v",{}]`, messageId)
	meterValuesConfirmation := ocpp16.NewMeterValuesConfirmation()
	channel := NewMockWebSocket(wsId)

	coreListener := MockCentralSystemCoreListener{}
	coreListener.On("OnMeterValues", mock.AnythingOfType("string"), mock.Anything).Return(meterValuesConfirmation, nil).Run(func(args mock.Arguments) {
		request := args.Get(1).(*ocpp16.MeterValuesRequest)
		assert.Equal(t, connectorId, request.ConnectorId)
		assert.Equal(t, 1, len(request.MeterValue))
		mv := request.MeterValue[0]
		assertDateTimeEquality(t, timestamp, *mv.Timestamp)
		assert.Equal(t, 1, len(mv.SampledValue))
		sv := mv.SampledValue[0]
		assert.Equal(t, mockValue, sv.Value)
		assert.Equal(t, mockUnit, sv.Unit)
	})
	setupDefaultCentralSystemHandlers(suite, coreListener, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	assert.Nil(t, err)
	confirmation, err := suite.chargePoint.MeterValues(connectorId, meterValues)
	assert.Nil(t, err)
	assert.NotNil(t, confirmation)
}

func (suite *OcppV16TestSuite) TestMeterValuesInvalidEndpoint() {
	messageId := defaultMessageId
	connectorId := 1
	mockValue := "value"
	mockUnit := ocpp16.UnitOfMeasureKW
	timestamp := ocpp16.DateTime{Time: time.Now()}
	meterValues := []ocpp16.MeterValue{{Timestamp: ocpp16.NewDateTime(time.Now()), SampledValue: []ocpp16.SampledValue{{Value: mockValue, Unit: mockUnit}}}}
	meterValuesRequest := ocpp16.NewMeterValuesRequest(connectorId, meterValues)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"connectorId":%v,"meterValue":[{"timestamp":"%v","sampledValue":[{"value":"%v","unit":"%v"}]}]}]`, messageId, ocpp16.MeterValuesFeatureName, connectorId, timestamp.Format(ocpp16.ISO8601), mockValue, mockUnit)
	testUnsupportedRequestFromCentralSystem(suite, meterValuesRequest, requestJson, messageId)
}
