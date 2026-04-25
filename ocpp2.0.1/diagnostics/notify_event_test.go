package diagnostics

import (
	"reflect"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *diagnosticsTestSuite) TestNotifyEventRequestValidation() {
	t := suite.T()
	eventData := EventData{
		EventID:               1,
		Timestamp:             types.NewDateTime(time.Now()),
		Trigger:               EventTriggerAlerting,
		Cause:                 tests.NewInt(42),
		ActualValue:           "someValue",
		TechCode:              "742",
		TechInfo:              "stacktrace",
		Cleared:               false,
		TransactionID:         "1234",
		VariableMonitoringID:  tests.NewInt(99),
		EventNotificationType: EventPreconfiguredMonitor,
		Component:             types.Component{Name: "component1"},
		Variable:              types.Variable{Name: "variable1"},
	}
	var requestTable = []tests.GenericTestEntry{
		{NotifyEventRequest{GeneratedAt: types.NewDateTime(time.Now()), SeqNo: 0, Tbc: false, EventData: []EventData{eventData, eventData}}, true},
		{NotifyEventRequest{GeneratedAt: types.NewDateTime(time.Now()), SeqNo: 0, Tbc: false, EventData: []EventData{eventData}}, true},
		{NotifyEventRequest{GeneratedAt: types.NewDateTime(time.Now()), SeqNo: 0, EventData: []EventData{eventData}}, true},
		{NotifyEventRequest{GeneratedAt: types.NewDateTime(time.Now()), EventData: []EventData{eventData}}, true},
		{NotifyEventRequest{EventData: []EventData{eventData}}, false},
		{NotifyEventRequest{}, false},
		{NotifyEventRequest{GeneratedAt: types.NewDateTime(time.Now()), SeqNo: -1, Tbc: false, EventData: []EventData{eventData}}, false},
		{NotifyEventRequest{GeneratedAt: types.NewDateTime(time.Now()), SeqNo: 0, Tbc: false, EventData: []EventData{}}, false},
		{NotifyEventRequest{GeneratedAt: types.NewDateTime(time.Now()), SeqNo: 0, Tbc: false}, false},
		{NotifyEventRequest{GeneratedAt: types.NewDateTime(time.Now()), SeqNo: 0, Tbc: false, EventData: []EventData{{Timestamp: types.NewDateTime(time.Now()), Trigger: EventTriggerAlerting, Cause: tests.NewInt(42), ActualValue: "someValue", Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}}}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *diagnosticsTestSuite) TestNotifyEventDataValidation() {
	t := suite.T()
	var table = []tests.GenericTestEntry{
		{EventData{EventID: 1, Timestamp: types.NewDateTime(time.Now()), Trigger: EventTriggerAlerting, Cause: tests.NewInt(42), ActualValue: "someValue", TechCode: "742", TechInfo: "stacktrace", Cleared: false, TransactionID: "1234", VariableMonitoringID: tests.NewInt(99), EventNotificationType: EventPreconfiguredMonitor, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}, true},
		{EventData{EventID: 1, Timestamp: types.NewDateTime(time.Now()), Trigger: EventTriggerAlerting, Cause: tests.NewInt(42), ActualValue: "someValue", TechCode: "742", TechInfo: "stacktrace", Cleared: false, TransactionID: "1234", EventNotificationType: EventPreconfiguredMonitor, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}, true},
		{EventData{EventID: 1, Timestamp: types.NewDateTime(time.Now()), Trigger: EventTriggerAlerting, Cause: tests.NewInt(42), ActualValue: "someValue", TechCode: "742", TechInfo: "stacktrace", Cleared: false, EventNotificationType: EventPreconfiguredMonitor, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}, true},
		{EventData{EventID: 1, Timestamp: types.NewDateTime(time.Now()), Trigger: EventTriggerAlerting, Cause: tests.NewInt(42), ActualValue: "someValue", TechCode: "742", TechInfo: "stacktrace", EventNotificationType: EventPreconfiguredMonitor, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}, true},
		{EventData{EventID: 1, Timestamp: types.NewDateTime(time.Now()), Trigger: EventTriggerAlerting, Cause: tests.NewInt(42), ActualValue: "someValue", TechCode: "742", EventNotificationType: EventPreconfiguredMonitor, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}, true},
		{EventData{EventID: 1, Timestamp: types.NewDateTime(time.Now()), Trigger: EventTriggerAlerting, Cause: tests.NewInt(42), ActualValue: "someValue", EventNotificationType: EventPreconfiguredMonitor, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}, true},
		{EventData{Timestamp: types.NewDateTime(time.Now()), Trigger: EventTriggerAlerting, Cause: tests.NewInt(42), ActualValue: "someValue", EventNotificationType: EventPreconfiguredMonitor, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}, true},
		{EventData{Timestamp: types.NewDateTime(time.Now()), Trigger: EventTriggerAlerting, ActualValue: "someValue", EventNotificationType: EventPreconfiguredMonitor, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}, true},
		{EventData{Trigger: EventTriggerAlerting, ActualValue: "someValue", EventNotificationType: EventPreconfiguredMonitor, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}, false},
		{EventData{Timestamp: types.NewDateTime(time.Now()), ActualValue: "someValue", EventNotificationType: EventPreconfiguredMonitor, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}, false},
		{EventData{Timestamp: types.NewDateTime(time.Now()), Trigger: EventTriggerAlerting, EventNotificationType: EventPreconfiguredMonitor, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}, false},
		{EventData{Timestamp: types.NewDateTime(time.Now()), Trigger: EventTriggerAlerting, Cause: tests.NewInt(42), ActualValue: "someValue", Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}, false},
		{EventData{Timestamp: types.NewDateTime(time.Now()), Trigger: EventTriggerAlerting, Cause: tests.NewInt(42), ActualValue: "someValue", EventNotificationType: EventPreconfiguredMonitor, Component: types.Component{Name: "component1"}}, false},
		{EventData{Timestamp: types.NewDateTime(time.Now()), Trigger: EventTriggerAlerting, Cause: tests.NewInt(42), ActualValue: "someValue", EventNotificationType: EventPreconfiguredMonitor, Variable: types.Variable{Name: "variable1"}}, false},
		{EventData{EventID: -1, Timestamp: types.NewDateTime(time.Now()), Trigger: EventTriggerAlerting, Cause: tests.NewInt(42), ActualValue: "someValue", TechCode: "742", TechInfo: "stacktrace", Cleared: false, TransactionID: "1234", VariableMonitoringID: tests.NewInt(99), EventNotificationType: EventPreconfiguredMonitor, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}, false},
		{EventData{EventID: 1, Timestamp: types.NewDateTime(time.Now()), Trigger: "invalidEventTrigger", Cause: tests.NewInt(42), ActualValue: "someValue", TechCode: "742", TechInfo: "stacktrace", Cleared: false, TransactionID: "1234", VariableMonitoringID: tests.NewInt(99), EventNotificationType: EventPreconfiguredMonitor, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}, false},
		{EventData{EventID: 1, Timestamp: types.NewDateTime(time.Now()), Trigger: EventTriggerAlerting, Cause: tests.NewInt(42), ActualValue: ">2500................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................", TechCode: "742", TechInfo: "stacktrace", Cleared: false, TransactionID: "1234", VariableMonitoringID: tests.NewInt(99), EventNotificationType: EventPreconfiguredMonitor, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}, false},
		{EventData{EventID: 1, Timestamp: types.NewDateTime(time.Now()), Trigger: EventTriggerAlerting, Cause: tests.NewInt(42), ActualValue: "someValue", TechCode: ">50................................................", TechInfo: "stacktrace", Cleared: false, TransactionID: "1234", VariableMonitoringID: tests.NewInt(99), EventNotificationType: EventPreconfiguredMonitor, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}, false},
		{EventData{EventID: 1, Timestamp: types.NewDateTime(time.Now()), Trigger: EventTriggerAlerting, Cause: tests.NewInt(42), ActualValue: "someValue", TechCode: "742", TechInfo: ">500.................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................", Cleared: false, TransactionID: "1234", VariableMonitoringID: tests.NewInt(99), EventNotificationType: EventPreconfiguredMonitor, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}, false},
		{EventData{EventID: 1, Timestamp: types.NewDateTime(time.Now()), Trigger: EventTriggerAlerting, Cause: tests.NewInt(42), ActualValue: "someValue", TechCode: "742", TechInfo: "stacktrace", Cleared: false, TransactionID: ">36..................................", VariableMonitoringID: tests.NewInt(99), EventNotificationType: "invalidEventNotification", Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}, false},
		{EventData{EventID: 1, Timestamp: types.NewDateTime(time.Now()), Trigger: EventTriggerAlerting, Cause: tests.NewInt(42), ActualValue: "someValue", TechCode: "742", TechInfo: "stacktrace", Cleared: false, TransactionID: "1234", VariableMonitoringID: tests.NewInt(99), EventNotificationType: EventPreconfiguredMonitor, Component: types.Component{Name: ">50................................................"}, Variable: types.Variable{Name: "variable1"}}, false},
		{EventData{EventID: 1, Timestamp: types.NewDateTime(time.Now()), Trigger: EventTriggerAlerting, Cause: tests.NewInt(42), ActualValue: "someValue", TechCode: "742", TechInfo: "stacktrace", Cleared: false, TransactionID: "1234", VariableMonitoringID: tests.NewInt(99), EventNotificationType: EventPreconfiguredMonitor, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: ">50................................................"}}, false},
	}
	tests.ExecuteGenericTestTable(t, table)
}

func (suite *diagnosticsTestSuite) TestNotifyEventResponseValidation() {
	t := suite.T()
	var responseTable = []tests.GenericTestEntry{
		{NotifyEventResponse{}, true},
	}
	tests.ExecuteGenericTestTable(t, responseTable)
}

func (suite *diagnosticsTestSuite) TestNotifyEventFeature() {
	feature := NotifyEventFeature{}
	suite.Equal(NotifyEventFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(NotifyEventRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(NotifyEventResponse{}), feature.GetResponseType())
}

func (suite *diagnosticsTestSuite) TestNewNotifyEventRequest() {
	ts := types.NewDateTime(time.Now())
	eventData := []EventData{
		{
			EventID: 1, Timestamp: ts, Trigger: EventTriggerAlerting,
			ActualValue: "val", EventNotificationType: EventPreconfiguredMonitor,
			Component: types.Component{Name: "comp1"}, Variable: types.Variable{Name: "var1"},
		},
	}
	req := NewNotifyEventRequest(ts, 0, eventData)
	suite.NotNil(req)
	suite.Equal(NotifyEventFeatureName, req.GetFeatureName())
	suite.Equal(ts, req.GeneratedAt)
	suite.Equal(0, req.SeqNo)
	suite.Equal(eventData, req.EventData)
}

func (suite *diagnosticsTestSuite) TestNewNotifyEventResponse() {
	resp := NewNotifyEventResponse()
	suite.NotNil(resp)
	suite.Equal(NotifyEventFeatureName, resp.GetFeatureName())
}
