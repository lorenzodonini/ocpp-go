package provisioning

import (
	"reflect"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *provisioningTestSuite) TestNotifyReportRequestValidation() {
	t := suite.T()
	reportData := ReportData{
		Component:               types.Component{Name: "component1"},
		Variable:                types.Variable{Name: "variable1"},
		VariableAttribute:       []VariableAttribute{NewVariableAttribute()},
		VariableCharacteristics: NewVariableCharacteristics(TypeString, true),
	}
	var requestTable = []tests.GenericTestEntry{
		{NewNotifyReportRequest(42, types.NewDateTime(time.Now()), 0), true},
		{NotifyReportRequest{RequestID: 42, GeneratedAt: types.NewDateTime(time.Now()), Tbc: true, SeqNo: 0, ReportData: []ReportData{reportData}}, true},
		{NotifyReportRequest{RequestID: 42, GeneratedAt: types.NewDateTime(time.Now()), Tbc: true, SeqNo: 0, ReportData: []ReportData{{Component: types.Component{Name: "comp1"}, Variable: types.Variable{Name: "var1"}, VariableAttribute: []VariableAttribute{NewVariableAttribute()}}}}, true},
		{NotifyReportRequest{RequestID: 42, GeneratedAt: types.NewDateTime(time.Now()), Tbc: true, ReportData: []ReportData{reportData}}, true},
		{NotifyReportRequest{RequestID: 42, GeneratedAt: types.NewDateTime(time.Now()), ReportData: []ReportData{reportData}}, true},
		{NotifyReportRequest{RequestID: 42, GeneratedAt: types.NewDateTime(time.Now()), ReportData: []ReportData{}}, true},
		{NotifyReportRequest{RequestID: 42, GeneratedAt: types.NewDateTime(time.Now())}, true},
		{NotifyReportRequest{GeneratedAt: types.NewDateTime(time.Now())}, true},
		{NotifyReportRequest{}, false},
		{NotifyReportRequest{RequestID: -1, GeneratedAt: types.NewDateTime(time.Now()), Tbc: true, SeqNo: 0, ReportData: []ReportData{reportData}}, false},
		{NotifyReportRequest{RequestID: 42, GeneratedAt: types.NewDateTime(time.Now()), Tbc: true, SeqNo: -1, ReportData: []ReportData{reportData}}, false},
		{NotifyReportRequest{RequestID: 42, GeneratedAt: types.NewDateTime(time.Now()), Tbc: true, SeqNo: 0, ReportData: []ReportData{{Component: types.Component{}, Variable: types.Variable{Name: "var1"}, VariableAttribute: []VariableAttribute{NewVariableAttribute()}, VariableCharacteristics: NewVariableCharacteristics(TypeString, true)}}}, false},
		{NotifyReportRequest{RequestID: 42, GeneratedAt: types.NewDateTime(time.Now()), Tbc: true, SeqNo: 0, ReportData: []ReportData{{Component: types.Component{Name: "comp1"}, Variable: types.Variable{}, VariableAttribute: []VariableAttribute{NewVariableAttribute()}, VariableCharacteristics: NewVariableCharacteristics(TypeString, true)}}}, false},
		{NotifyReportRequest{RequestID: 42, GeneratedAt: types.NewDateTime(time.Now()), Tbc: true, SeqNo: 0, ReportData: []ReportData{{Component: types.Component{Name: "comp1"}, Variable: types.Variable{Name: "var1"}, VariableAttribute: []VariableAttribute{}, VariableCharacteristics: NewVariableCharacteristics(TypeString, true)}}}, false},
		{NotifyReportRequest{RequestID: 42, GeneratedAt: types.NewDateTime(time.Now()), Tbc: true, SeqNo: 0, ReportData: []ReportData{{Component: types.Component{Name: "comp1"}, Variable: types.Variable{Name: "var1"}, VariableAttribute: []VariableAttribute{NewVariableAttribute(), NewVariableAttribute(), NewVariableAttribute(), NewVariableAttribute(), NewVariableAttribute()}, VariableCharacteristics: NewVariableCharacteristics(TypeString, true)}}}, false},
		{NotifyReportRequest{RequestID: 42, GeneratedAt: types.NewDateTime(time.Now()), Tbc: true, SeqNo: 0, ReportData: []ReportData{{Component: types.Component{Name: "comp1"}, Variable: types.Variable{Name: "var1"}, VariableAttribute: []VariableAttribute{NewVariableAttribute()}, VariableCharacteristics: NewVariableCharacteristics("unknownType", true)}}}, false},
		{NotifyReportRequest{RequestID: 42, GeneratedAt: types.NewDateTime(time.Now()), Tbc: true, SeqNo: 0, ReportData: []ReportData{{Component: types.Component{Name: "comp1"}, Variable: types.Variable{Name: "var1"}, VariableAttribute: []VariableAttribute{{Mutability: "invalidMutability"}}, VariableCharacteristics: NewVariableCharacteristics(TypeString, true)}}}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *provisioningTestSuite) TestVariableCharacteristicsValidation() {
	t := suite.T()
	var table = []tests.GenericTestEntry{
		{NewVariableCharacteristics(TypeString, false), true},
		{VariableCharacteristics{Unit: "KWh", DataType: TypeDecimal, MinLimit: tests.NewFloat(1.0), MaxLimit: tests.NewFloat(22.0), ValuesList: "7.0", SupportsMonitoring: true}, true},
		{VariableCharacteristics{Unit: "KWh", DataType: TypeDecimal, MinLimit: tests.NewFloat(1.0), MaxLimit: tests.NewFloat(22.0), ValuesList: "7.0"}, true},
		{VariableCharacteristics{Unit: "KWh", DataType: TypeDecimal, MinLimit: tests.NewFloat(1.0), MaxLimit: tests.NewFloat(22.0)}, true},
		{VariableCharacteristics{Unit: "KWh", DataType: TypeDecimal, MinLimit: tests.NewFloat(-11.0), MaxLimit: tests.NewFloat(-2.0)}, true},
		{VariableCharacteristics{Unit: "KWh", DataType: TypeDecimal, MinLimit: tests.NewFloat(-1.0)}, true},
		{VariableCharacteristics{Unit: "KWh", DataType: TypeDecimal}, true},
		{VariableCharacteristics{DataType: TypeDecimal}, true},
		{VariableCharacteristics{DataType: TypeString}, true},
		{VariableCharacteristics{DataType: TypeInteger}, true},
		{VariableCharacteristics{DataType: TypeDateTime}, true},
		{VariableCharacteristics{DataType: TypeBoolean}, true},
		{VariableCharacteristics{DataType: TypeMemberList}, true},
		{VariableCharacteristics{DataType: TypeSequenceList}, true},
		{VariableCharacteristics{DataType: TypeOptionList}, true},
		{VariableCharacteristics{}, false},
		{VariableCharacteristics{Unit: ">16..............", DataType: TypeDecimal, MinLimit: tests.NewFloat(1.0), MaxLimit: tests.NewFloat(22.0), ValuesList: "7.0", SupportsMonitoring: true}, false},
		{VariableCharacteristics{Unit: "KWh", DataType: "invalidDataType", MinLimit: tests.NewFloat(1.0), MaxLimit: tests.NewFloat(22.0), ValuesList: "7.0", SupportsMonitoring: true}, false},
		{VariableCharacteristics{Unit: "KWh", DataType: TypeDecimal, MinLimit: tests.NewFloat(1.0), MaxLimit: tests.NewFloat(22.0), ValuesList: ">1000....................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................", SupportsMonitoring: true}, false},
	}
	tests.ExecuteGenericTestTable(t, table)
}

func (suite *provisioningTestSuite) TestVariableAttributeValidation() {
	t := suite.T()
	var table = []tests.GenericTestEntry{
		{NewVariableAttribute(), true},
		{VariableAttribute{Type: types.AttributeActual, Value: "someValue", Mutability: MutabilityReadWrite, Persistent: false, Constant: false}, true},
		{VariableAttribute{Type: types.AttributeActual, Value: "someValue", Mutability: MutabilityWriteOnly, Persistent: false, Constant: false}, true},
		{VariableAttribute{Type: types.AttributeActual, Value: "someValue", Mutability: MutabilityReadOnly, Persistent: false, Constant: false}, true},
		{VariableAttribute{Type: types.AttributeMaxSet, Value: "someValue", Mutability: MutabilityReadWrite, Persistent: false, Constant: false}, true},
		{VariableAttribute{Type: types.AttributeMinSet, Value: "someValue", Mutability: MutabilityReadWrite, Persistent: false, Constant: false}, true},
		{VariableAttribute{Type: types.AttributeTarget, Value: "someValue", Mutability: MutabilityReadWrite, Persistent: false, Constant: false}, true},
		{VariableAttribute{Type: types.AttributeActual, Value: "someValue", Mutability: MutabilityReadWrite}, true},
		{VariableAttribute{Type: types.AttributeActual, Value: "someValue"}, true},
		{VariableAttribute{Value: "someValue"}, true},
		//TODO: enable tests once validation on mutability field is enabled
		//{VariableAttribute{Mutability: MutabilityWriteOnly}, true},
		//{VariableAttribute{}, false},
		//{VariableAttribute{Mutability: MutabilityReadOnly}, false},
		//{VariableAttribute{Mutability: MutabilityReadWrite}, false},
		{VariableAttribute{Type: "invalidType", Value: "someValue", Mutability: MutabilityReadWrite}, false},
		{VariableAttribute{Type: types.AttributeActual, Value: "someValue", Mutability: "invalidMutability"}, false},
		{VariableAttribute{Type: types.AttributeActual, Value: ">2500................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................", Mutability: MutabilityReadWrite}, false},
	}
	tests.ExecuteGenericTestTable(t, table)
}

func (suite *provisioningTestSuite) TestNotifyReportResponseValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{NewNotifyReportResponse(), true},
		{NotifyReportResponse{}, true},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *provisioningTestSuite) TestNotifyReportFeature() {
	feature := NotifyReportFeature{}
	suite.Equal(NotifyReportFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(NotifyReportRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(NotifyReportResponse{}), feature.GetResponseType())
}

func (suite *provisioningTestSuite) TestNewNotifyReportRequest() {
	ts := types.NewDateTime(time.Now())
	req := NewNotifyReportRequest(42, ts, 0)
	suite.NotNil(req)
	suite.Equal(NotifyReportFeatureName, req.GetFeatureName())
	suite.Equal(42, req.RequestID)
	suite.Equal(ts, req.GeneratedAt)
	suite.Equal(0, req.SeqNo)
}

func (suite *provisioningTestSuite) TestNewNotifyReportResponse() {
	resp := NewNotifyReportResponse()
	suite.NotNil(resp)
	suite.Equal(NotifyReportFeatureName, resp.GetFeatureName())
}
