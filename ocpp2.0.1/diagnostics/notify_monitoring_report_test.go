package diagnostics

import (
	"reflect"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *diagnosticsTestSuite) TestNotifyMonitoringReportRequestValidation() {
	t := suite.T()
	validMonitoring := NewVariableMonitoring(1, false, 42.42, MonitorPeriodic, 0)
	invalidMonitoring := NewVariableMonitoring(1, false, 42.42, "invalidMonitorType", 0)
	monitoringData := MonitoringData{
		Component:          types.Component{Name: "component1"},
		Variable:           types.Variable{Name: "variable1"},
		VariableMonitoring: []VariableMonitoring{validMonitoring},
	}
	var requestTable = []tests.GenericTestEntry{
		{NotifyMonitoringReportRequest{RequestID: 42, Tbc: true, SeqNo: 0, GeneratedAt: types.NewDateTime(time.Now()), Monitor: []MonitoringData{monitoringData}}, true},
		{NotifyMonitoringReportRequest{RequestID: 42, Tbc: true, SeqNo: 0, GeneratedAt: types.NewDateTime(time.Now()), Monitor: []MonitoringData{}}, true},
		{NotifyMonitoringReportRequest{RequestID: 42, Tbc: true, SeqNo: 0, GeneratedAt: types.NewDateTime(time.Now())}, true},
		{NotifyMonitoringReportRequest{RequestID: 42, SeqNo: 0, GeneratedAt: types.NewDateTime(time.Now())}, true},
		{NotifyMonitoringReportRequest{RequestID: 42, GeneratedAt: types.NewDateTime(time.Now())}, true},
		{NotifyMonitoringReportRequest{GeneratedAt: types.NewDateTime(time.Now())}, true},
		{NotifyMonitoringReportRequest{}, false},
		{NotifyMonitoringReportRequest{RequestID: -1, Tbc: true, SeqNo: 0, GeneratedAt: types.NewDateTime(time.Now()), Monitor: []MonitoringData{monitoringData}}, false},
		{NotifyMonitoringReportRequest{RequestID: 42, Tbc: true, SeqNo: -1, GeneratedAt: types.NewDateTime(time.Now()), Monitor: []MonitoringData{monitoringData}}, false},
		{NotifyMonitoringReportRequest{RequestID: 42, Tbc: true, SeqNo: 0, GeneratedAt: types.NewDateTime(time.Now()), Monitor: []MonitoringData{{Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}, VariableMonitoring: []VariableMonitoring{invalidMonitoring}}}}, false},
		{NotifyMonitoringReportRequest{RequestID: 42, Tbc: true, SeqNo: 0, GeneratedAt: types.NewDateTime(time.Now()), Monitor: []MonitoringData{{Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}, VariableMonitoring: []VariableMonitoring{}}}}, false},
		{NotifyMonitoringReportRequest{RequestID: 42, Tbc: true, SeqNo: 0, GeneratedAt: types.NewDateTime(time.Now()), Monitor: []MonitoringData{{Component: types.Component{Name: "component1"}, Variable: types.Variable{}, VariableMonitoring: []VariableMonitoring{validMonitoring}}}}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *diagnosticsTestSuite) TestVariableMonitoringValidation() {
	t := suite.T()
	var table = []tests.GenericTestEntry{
		{VariableMonitoring{ID: 1, Transaction: false, Value: 42.42, Type: MonitorPeriodic, Severity: 0}, true},
		{VariableMonitoring{ID: 1, Transaction: false, Value: 42.42, Type: MonitorPeriodic, Severity: 9}, true},
		{VariableMonitoring{ID: 1, Transaction: false, Value: 42.42, Type: MonitorPeriodicClockAligned, Severity: 0}, true},
		{VariableMonitoring{ID: 1, Transaction: false, Value: 42.42, Type: MonitorUpperThreshold, Severity: 0}, true},
		{VariableMonitoring{ID: 1, Transaction: false, Value: 42.42, Type: MonitorLowerThreshold, Severity: 0}, true},
		{VariableMonitoring{ID: 1, Transaction: false, Value: 42.42, Type: MonitorDelta, Severity: 0}, true},
		{VariableMonitoring{ID: 1, Transaction: false, Value: -42.42, Type: MonitorPeriodic, Severity: 0}, true},
		{VariableMonitoring{ID: 1, Transaction: false, Value: 42.42, Type: MonitorPeriodic}, true},
		{VariableMonitoring{ID: 1, Transaction: false, Type: MonitorPeriodic}, true},
		{VariableMonitoring{ID: 1, Type: MonitorPeriodic}, true},
		{VariableMonitoring{Type: MonitorPeriodic}, true},
		{VariableMonitoring{}, false},
		{VariableMonitoring{ID: -1, Transaction: false, Value: 42.42, Type: MonitorPeriodic, Severity: 0}, false},
		{VariableMonitoring{ID: 1, Transaction: false, Value: 42.42, Type: MonitorPeriodic, Severity: -1}, false},
		{VariableMonitoring{ID: 1, Transaction: false, Value: 42.42, Type: MonitorPeriodic, Severity: 10}, false},
		{VariableMonitoring{ID: 1, Transaction: false, Value: 42.42, Type: "invalidMonitorType", Severity: 0}, false},
	}
	tests.ExecuteGenericTestTable(t, table)
}

func (suite *diagnosticsTestSuite) TestNotifyMonitoringReportResponseValidation() {
	t := suite.T()
	var responseTable = []tests.GenericTestEntry{
		{NotifyMonitoringReportResponse{}, true},
	}
	tests.ExecuteGenericTestTable(t, responseTable)
}

func (suite *diagnosticsTestSuite) TestNotifyMonitoringReportFeature() {
	feature := NotifyMonitoringReportFeature{}
	suite.Equal(NotifyMonitoringReportFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(NotifyMonitoringReportRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(NotifyMonitoringReportResponse{}), feature.GetResponseType())
}

func (suite *diagnosticsTestSuite) TestNewNotifyMonitoringReportRequest() {
	ts := types.NewDateTime(time.Now())
	req := NewNotifyMonitoringReportRequest(42, 0, ts, nil)
	suite.NotNil(req)
	suite.Equal(NotifyMonitoringReportFeatureName, req.GetFeatureName())
	suite.Equal(42, req.RequestID)
	suite.Equal(0, req.SeqNo)
	suite.Equal(ts, req.GeneratedAt)
}

func (suite *diagnosticsTestSuite) TestNewNotifyMonitoringReportResponse() {
	resp := NewNotifyMonitoringReportResponse()
	suite.NotNil(resp)
	suite.Equal(NotifyMonitoringReportFeatureName, resp.GetFeatureName())
}
