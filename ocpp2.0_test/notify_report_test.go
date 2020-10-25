package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/provisioning"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"time"
)

// Tests
func (suite *OcppV2TestSuite) TestNotifyReportRequestValidation() {
	t := suite.T()
	reportData := provisioning.ReportData{
		Component:               types.Component{Name: "component1"},
		Variable:                types.Variable{Name: "variable1"},
		VariableAttribute:       []provisioning.VariableAttribute{provisioning.NewVariableAttribute()},
		VariableCharacteristics: provisioning.NewVariableCharacteristics(provisioning.TypeString, true),
	}
	var requestTable = []GenericTestEntry{
		{provisioning.NewNotifyReportRequest(42, types.NewDateTime(time.Now()), 0), true},
		{provisioning.NotifyReportRequest{RequestID: 42, GeneratedAt: types.NewDateTime(time.Now()), Tbc: true, SeqNo: 0, ReportData: []provisioning.ReportData{reportData}}, true},
		{provisioning.NotifyReportRequest{RequestID: 42, GeneratedAt: types.NewDateTime(time.Now()), Tbc: true, SeqNo: 0, ReportData: []provisioning.ReportData{{Component: types.Component{Name: "comp1"}, Variable: types.Variable{Name: "var1"}, VariableAttribute: []provisioning.VariableAttribute{provisioning.NewVariableAttribute()}}}}, true},
		{provisioning.NotifyReportRequest{RequestID: 42, GeneratedAt: types.NewDateTime(time.Now()), Tbc: true, ReportData: []provisioning.ReportData{reportData}}, true},
		{provisioning.NotifyReportRequest{RequestID: 42, GeneratedAt: types.NewDateTime(time.Now()), ReportData: []provisioning.ReportData{reportData}}, true},
		{provisioning.NotifyReportRequest{RequestID: 42, GeneratedAt: types.NewDateTime(time.Now()), ReportData: []provisioning.ReportData{}}, true},
		{provisioning.NotifyReportRequest{RequestID: 42, GeneratedAt: types.NewDateTime(time.Now())}, true},
		{provisioning.NotifyReportRequest{GeneratedAt: types.NewDateTime(time.Now())}, true},
		{provisioning.NotifyReportRequest{}, false},
		{provisioning.NotifyReportRequest{RequestID: -1, GeneratedAt: types.NewDateTime(time.Now()), Tbc: true, SeqNo: 0, ReportData: []provisioning.ReportData{reportData}}, false},
		{provisioning.NotifyReportRequest{RequestID: 42, GeneratedAt: types.NewDateTime(time.Now()), Tbc: true, SeqNo: -1, ReportData: []provisioning.ReportData{reportData}}, false},
		{provisioning.NotifyReportRequest{RequestID: 42, GeneratedAt: types.NewDateTime(time.Now()), Tbc: true, SeqNo: 0, ReportData: []provisioning.ReportData{{Component: types.Component{}, Variable: types.Variable{Name: "var1"}, VariableAttribute: []provisioning.VariableAttribute{provisioning.NewVariableAttribute()}, VariableCharacteristics: provisioning.NewVariableCharacteristics(provisioning.TypeString, true)}}}, false},
		{provisioning.NotifyReportRequest{RequestID: 42, GeneratedAt: types.NewDateTime(time.Now()), Tbc: true, SeqNo: 0, ReportData: []provisioning.ReportData{{Component: types.Component{Name: "comp1"}, Variable: types.Variable{}, VariableAttribute: []provisioning.VariableAttribute{provisioning.NewVariableAttribute()}, VariableCharacteristics: provisioning.NewVariableCharacteristics(provisioning.TypeString, true)}}}, false},
		{provisioning.NotifyReportRequest{RequestID: 42, GeneratedAt: types.NewDateTime(time.Now()), Tbc: true, SeqNo: 0, ReportData: []provisioning.ReportData{{Component: types.Component{Name: "comp1"}, Variable: types.Variable{Name: "var1"}, VariableAttribute: []provisioning.VariableAttribute{}, VariableCharacteristics: provisioning.NewVariableCharacteristics(provisioning.TypeString, true)}}}, false},
		{provisioning.NotifyReportRequest{RequestID: 42, GeneratedAt: types.NewDateTime(time.Now()), Tbc: true, SeqNo: 0, ReportData: []provisioning.ReportData{{Component: types.Component{Name: "comp1"}, Variable: types.Variable{Name: "var1"}, VariableAttribute: []provisioning.VariableAttribute{provisioning.NewVariableAttribute(), provisioning.NewVariableAttribute(), provisioning.NewVariableAttribute(), provisioning.NewVariableAttribute(), provisioning.NewVariableAttribute()}, VariableCharacteristics: provisioning.NewVariableCharacteristics(provisioning.TypeString, true)}}}, false},
		{provisioning.NotifyReportRequest{RequestID: 42, GeneratedAt: types.NewDateTime(time.Now()), Tbc: true, SeqNo: 0, ReportData: []provisioning.ReportData{{Component: types.Component{Name: "comp1"}, Variable: types.Variable{Name: "var1"}, VariableAttribute: []provisioning.VariableAttribute{provisioning.NewVariableAttribute()}, VariableCharacteristics: provisioning.NewVariableCharacteristics("unknownType", true)}}}, false},
		{provisioning.NotifyReportRequest{RequestID: 42, GeneratedAt: types.NewDateTime(time.Now()), Tbc: true, SeqNo: 0, ReportData: []provisioning.ReportData{{Component: types.Component{Name: "comp1"}, Variable: types.Variable{Name: "var1"}, VariableAttribute: []provisioning.VariableAttribute{{Mutability: "invalidMutability"}}, VariableCharacteristics: provisioning.NewVariableCharacteristics(provisioning.TypeString, true)}}}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestVariableCharacteristicsValidation() {
	t := suite.T()
	var table = []GenericTestEntry{
		{provisioning.NewVariableCharacteristics(provisioning.TypeString, false), true},
		{provisioning.VariableCharacteristics{Unit: "KWh", DataType: provisioning.TypeDecimal, MinLimit: newFloat(1.0), MaxLimit: newFloat(22.0), ValuesList: "7.0", SupportsMonitoring: true}, true},
		{provisioning.VariableCharacteristics{Unit: "KWh", DataType: provisioning.TypeDecimal, MinLimit: newFloat(1.0), MaxLimit: newFloat(22.0), ValuesList: "7.0"}, true},
		{provisioning.VariableCharacteristics{Unit: "KWh", DataType: provisioning.TypeDecimal, MinLimit: newFloat(1.0), MaxLimit: newFloat(22.0)}, true},
		{provisioning.VariableCharacteristics{Unit: "KWh", DataType: provisioning.TypeDecimal, MinLimit: newFloat(-11.0), MaxLimit: newFloat(-2.0)}, true},
		{provisioning.VariableCharacteristics{Unit: "KWh", DataType: provisioning.TypeDecimal, MinLimit: newFloat(-1.0)}, true},
		{provisioning.VariableCharacteristics{Unit: "KWh", DataType: provisioning.TypeDecimal}, true},
		{provisioning.VariableCharacteristics{DataType: provisioning.TypeDecimal}, true},
		{provisioning.VariableCharacteristics{DataType: provisioning.TypeString}, true},
		{provisioning.VariableCharacteristics{DataType: provisioning.TypeInteger}, true},
		{provisioning.VariableCharacteristics{DataType: provisioning.TypeDateTime}, true},
		{provisioning.VariableCharacteristics{DataType: provisioning.TypeBoolean}, true},
		{provisioning.VariableCharacteristics{DataType: provisioning.TypeMemberList}, true},
		{provisioning.VariableCharacteristics{DataType: provisioning.TypeSequenceList}, true},
		{provisioning.VariableCharacteristics{DataType: provisioning.TypeOptionList}, true},
		{provisioning.VariableCharacteristics{}, false},
		{provisioning.VariableCharacteristics{Unit: ">16..............", DataType: provisioning.TypeDecimal, MinLimit: newFloat(1.0), MaxLimit: newFloat(22.0), ValuesList: "7.0", SupportsMonitoring: true}, false},
		{provisioning.VariableCharacteristics{Unit: "KWh", DataType: "invalidDataType", MinLimit: newFloat(1.0), MaxLimit: newFloat(22.0), ValuesList: "7.0", SupportsMonitoring: true}, false},
		{provisioning.VariableCharacteristics{Unit: "KWh", DataType: provisioning.TypeDecimal, MinLimit: newFloat(1.0), MaxLimit: newFloat(22.0), ValuesList: ">1000....................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................", SupportsMonitoring: true}, false},
	}
	ExecuteGenericTestTable(t, table)
}

func (suite *OcppV2TestSuite) TestVariableAttributeValidation() {
	t := suite.T()
	var table = []GenericTestEntry{
		{provisioning.NewVariableAttribute(), true},
		{provisioning.VariableAttribute{Type: types.AttributeActual, Value: "someValue", Mutability: provisioning.MutabilityReadWrite, Persistent: false, Constant: false}, true},
		{provisioning.VariableAttribute{Type: types.AttributeActual, Value: "someValue", Mutability: provisioning.MutabilityWriteOnly, Persistent: false, Constant: false}, true},
		{provisioning.VariableAttribute{Type: types.AttributeActual, Value: "someValue", Mutability: provisioning.MutabilityReadOnly, Persistent: false, Constant: false}, true},
		{provisioning.VariableAttribute{Type: types.AttributeMaxSet, Value: "someValue", Mutability: provisioning.MutabilityReadWrite, Persistent: false, Constant: false}, true},
		{provisioning.VariableAttribute{Type: types.AttributeMinSet, Value: "someValue", Mutability: provisioning.MutabilityReadWrite, Persistent: false, Constant: false}, true},
		{provisioning.VariableAttribute{Type: types.AttributeTarget, Value: "someValue", Mutability: provisioning.MutabilityReadWrite, Persistent: false, Constant: false}, true},
		{provisioning.VariableAttribute{Type: types.AttributeActual, Value: "someValue", Mutability: provisioning.MutabilityReadWrite}, true},
		{provisioning.VariableAttribute{Type: types.AttributeActual, Value: "someValue"}, true},
		{provisioning.VariableAttribute{Value: "someValue"}, true},
		//TODO: enable tests once validation on mutability field is enabled
		//{provisioning.VariableAttribute{Mutability: provisioning.MutabilityWriteOnly}, true},
		//{provisioning.VariableAttribute{}, false},
		//{provisioning.VariableAttribute{Mutability: provisioning.MutabilityReadOnly}, false},
		//{provisioning.VariableAttribute{Mutability: provisioning.MutabilityReadWrite}, false},
		{provisioning.VariableAttribute{Type: "invalidType", Value: "someValue", Mutability: provisioning.MutabilityReadWrite}, false},
		{provisioning.VariableAttribute{Type: types.AttributeActual, Value: "someValue", Mutability: "invalidMutability"}, false},
		{provisioning.VariableAttribute{Type: types.AttributeActual, Value: ">2500................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................", Mutability: provisioning.MutabilityReadWrite}, false},
	}
	ExecuteGenericTestTable(t, table)
}

func (suite *OcppV2TestSuite) TestNotifyReportResponseValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{provisioning.NewNotifyReportResponse(), true},
		{provisioning.NotifyReportResponse{}, true},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

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
	variableCharacteristics := &provisioning.VariableCharacteristics{Unit: "KWh", DataType: provisioning.TypeString, MaxLimit: newFloat(22.0), SupportsMonitoring: true}
	reportData := provisioning.ReportData{
		Component:               types.Component{Name: "component1"},
		Variable:                types.Variable{Name: "variable1"},
		VariableAttribute:       []provisioning.VariableAttribute{variableAttribute},
		VariableCharacteristics: variableCharacteristics,
	}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"requestId":%v,"generatedAt":"%v","tbc":%v,"seqNo":%v,"reportData":[{"component":{"name":"%v"},"variable":{"name":"%v"},"variableAttribute":[{"type":"%v","value":"%v","mutability":"%v"}],"variableCharacteristics":{"unit":"%v","dataType":"%v","maxLimit":%v,"supportsMonitoring":%v}}]}]`,
		messageId, provisioning.NotifyReportFeatureName, requestID, generatedAt.FormatTimestamp(), tbc, seqNo, reportData.Component.Name, reportData.Variable.Name, variableAttribute.Type, variableAttribute.Value, variableAttribute.Mutability, variableCharacteristics.Unit, variableCharacteristics.DataType, *variableCharacteristics.MaxLimit, variableCharacteristics.SupportsMonitoring)
	responseJson := fmt.Sprintf(`[3,"%v",{}]`, messageId)
	fmt.Println(responseJson)
	notifyReportResponse := provisioning.NewNotifyReportResponse()
	channel := NewMockWebSocket(wsId)

	handler := MockCSMSProvisioningHandler{}
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
	variableCharacteristics := &provisioning.VariableCharacteristics{Unit: "KWh", DataType: provisioning.TypeString, MaxLimit: newFloat(22.0), SupportsMonitoring: true}
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
