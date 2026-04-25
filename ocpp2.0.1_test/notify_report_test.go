package ocpp2_test

import (
	"fmt"
	"time"

	"github.com/lorenzodonini/ocpp-go/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/provisioning"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

// Tests
func (suite *OcppV2TestSuite) TestNotifyReportE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := "1234"
	wsUrl := "someUrl"
	generatedAt := types.NewDateTime(time.Now())
	seqNo := 0
	requestID := 42
	tbc := true
	variableAttribute := provisioning.VariableAttribute{Type: types.AttributeTarget, Value: "someValue", Mutability: provisioning.MutabilityReadWrite}
	variableCharacteristics := &provisioning.VariableCharacteristics{Unit: "KWh", DataType: provisioning.TypeString, MaxLimit: tests.NewFloat(22.0), SupportsMonitoring: true}
	reportData := provisioning.ReportData{
		Component:               types.Component{Name: "component1"},
		Variable:                types.Variable{Name: "variable1"},
		VariableAttribute:       []provisioning.VariableAttribute{variableAttribute},
		VariableCharacteristics: variableCharacteristics,
	}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"requestId":%v,"generatedAt":"%v","tbc":%v,"seqNo":%v,"reportData":[{"component":{"name":"%v"},"variable":{"name":"%v"},"variableAttribute":[{"type":"%v","value":"%v","mutability":"%v"}],"variableCharacteristics":{"unit":"%v","dataType":"%v","maxLimit":%v,"supportsMonitoring":%v}}]}]`,
		messageId, provisioning.NotifyReportFeatureName, requestID, generatedAt.FormatTimestamp(), tbc, seqNo, reportData.Component.Name, reportData.Variable.Name, variableAttribute.Type, variableAttribute.Value, variableAttribute.Mutability, variableCharacteristics.Unit, variableCharacteristics.DataType, *variableCharacteristics.MaxLimit, variableCharacteristics.SupportsMonitoring)
	responseJson := fmt.Sprintf(`[3,"%v",{}]`, messageId)
	notifyReportResponse := provisioning.NewNotifyReportResponse()
	channel := NewMockWebSocket(wsId)

	handler := &MockCSMSProvisioningHandler{}
	handler.On("OnNotifyReport", mock.AnythingOfType("string"), mock.Anything).Return(notifyReportResponse, nil).Run(func(args mock.Arguments) {
		request := args.Get(1).(*provisioning.NotifyReportRequest)
		assert.Equal(t, requestID, request.RequestID)
		assertDateTimeEquality(t, generatedAt, request.GeneratedAt)
		assert.Equal(t, seqNo, request.SeqNo)
		assert.Equal(t, tbc, request.Tbc)
		require.Len(t, request.ReportData, 1)
		assert.Equal(t, reportData.Component.Name, request.ReportData[0].Component.Name)
		assert.Equal(t, reportData.Variable.Name, request.ReportData[0].Variable.Name)
		require.Len(t, request.ReportData[0].VariableAttribute, len(reportData.VariableAttribute))
		assert.Equal(t, variableAttribute.Mutability, request.ReportData[0].VariableAttribute[0].Mutability)
		assert.Equal(t, variableAttribute.Value, request.ReportData[0].VariableAttribute[0].Value)
		assert.Equal(t, variableAttribute.Type, request.ReportData[0].VariableAttribute[0].Type)
		assert.Equal(t, variableAttribute.Constant, request.ReportData[0].VariableAttribute[0].Constant)
		assert.Equal(t, variableAttribute.Persistent, request.ReportData[0].VariableAttribute[0].Persistent)
		require.NotNil(t, request.ReportData[0].VariableCharacteristics)
		assert.Equal(t, variableCharacteristics.Unit, request.ReportData[0].VariableCharacteristics.Unit)
		assert.Equal(t, variableCharacteristics.DataType, request.ReportData[0].VariableCharacteristics.DataType)
		assert.Equal(t, *variableCharacteristics.MaxLimit, *request.ReportData[0].VariableCharacteristics.MaxLimit)
		assert.Equal(t, variableCharacteristics.SupportsMonitoring, request.ReportData[0].VariableCharacteristics.SupportsMonitoring)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	response, err := suite.chargingStation.NotifyReport(requestID, generatedAt, seqNo, func(request *provisioning.NotifyReportRequest) {
		request.ReportData = []provisioning.ReportData{reportData}
		request.Tbc = tbc
	})
	require.Nil(t, err)
	require.NotNil(t, response)
}

func (suite *OcppV2TestSuite) TestNotifyReportInvalidEndpoint() {
	messageId := defaultMessageId
	generatedAt := types.NewDateTime(time.Now())
	seqNo := 0
	requestID := 42
	tbc := true
	variableAttribute := provisioning.VariableAttribute{Type: types.AttributeTarget, Value: "someValue", Mutability: provisioning.MutabilityReadWrite}
	variableCharacteristics := &provisioning.VariableCharacteristics{Unit: "KWh", DataType: provisioning.TypeString, MaxLimit: tests.NewFloat(22.0), SupportsMonitoring: true}
	reportData := provisioning.ReportData{
		Component:               types.Component{Name: "component1"},
		Variable:                types.Variable{Name: "variable1"},
		VariableAttribute:       []provisioning.VariableAttribute{variableAttribute},
		VariableCharacteristics: variableCharacteristics,
	}
	request := provisioning.NewNotifyReportRequest(requestID, generatedAt, seqNo)
	request.ReportData = []provisioning.ReportData{reportData}
	request.Tbc = tbc
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"requestId":%v,"generatedAt":"%v","tbc":%v,"seqNo":%v,"reportData":[{"component":{"name":"%v"},"variable":{"name":"%v"},"variableAttribute":[{"type":"%v","value":"%v","mutability":"%v"}],"variableCharacteristics":{"unit":"%v","dataType":"%v","maxLimit":%v,"supportsMonitoring":%v}}]}]`,
		messageId, provisioning.NotifyReportFeatureName, requestID, generatedAt.FormatTimestamp(), tbc, seqNo, reportData.Component.Name, reportData.Variable.Name, variableAttribute.Type, variableAttribute.Value, variableAttribute.Mutability, variableCharacteristics.Unit, variableCharacteristics.DataType, *variableCharacteristics.MaxLimit, variableCharacteristics.SupportsMonitoring)
	testUnsupportedRequestFromCentralSystem(suite, request, requestJson, messageId)
}
