package ocpp2_test

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/diagnostics"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

// Test
func (suite *OcppV2TestSuite) TestNotifyEventRequestValidation() {
	eventData := diagnostics.EventData{
		EventID:               1,
		Timestamp:             types.NewDateTime(time.Now()),
		Trigger:               diagnostics.EventTriggerAlerting,
		Cause:                 newInt(42),
		ActualValue:           "someValue",
		TechCode:              "742",
		TechInfo:              "stacktrace",
		Cleared:               false,
		TransactionID:         "1234",
		VariableMonitoringID:  newInt(99),
		EventNotificationType: diagnostics.EventPreconfiguredMonitor,
		Component:             types.Component{Name: "component1"},
		Variable:              types.Variable{Name: "variable1"},
	}
	var requestTable = []GenericTestEntry{
		{diagnostics.NotifyEventRequest{GeneratedAt: types.NewDateTime(time.Now()), SeqNo: 0, Tbc: false, EventData: []diagnostics.EventData{eventData, eventData}}, true},
		{diagnostics.NotifyEventRequest{GeneratedAt: types.NewDateTime(time.Now()), SeqNo: 0, Tbc: false, EventData: []diagnostics.EventData{eventData}}, true},
		{diagnostics.NotifyEventRequest{GeneratedAt: types.NewDateTime(time.Now()), SeqNo: 0, EventData: []diagnostics.EventData{eventData}}, true},
		{diagnostics.NotifyEventRequest{GeneratedAt: types.NewDateTime(time.Now()), EventData: []diagnostics.EventData{eventData}}, true},
		{diagnostics.NotifyEventRequest{EventData: []diagnostics.EventData{eventData}}, false},
		{diagnostics.NotifyEventRequest{}, false},
		{diagnostics.NotifyEventRequest{GeneratedAt: types.NewDateTime(time.Now()), SeqNo: -1, Tbc: false, EventData: []diagnostics.EventData{eventData}}, false},
		{diagnostics.NotifyEventRequest{GeneratedAt: types.NewDateTime(time.Now()), SeqNo: 0, Tbc: false, EventData: []diagnostics.EventData{}}, false},
		{diagnostics.NotifyEventRequest{GeneratedAt: types.NewDateTime(time.Now()), SeqNo: 0, Tbc: false}, false},
		{diagnostics.NotifyEventRequest{GeneratedAt: types.NewDateTime(time.Now()), SeqNo: 0, Tbc: false, EventData: []diagnostics.EventData{{Timestamp: types.NewDateTime(time.Now()), Trigger: diagnostics.EventTriggerAlerting, Cause: newInt(42), ActualValue: "someValue", Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}}}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV2TestSuite) TestNotifyEventDataValidation() {
	var table = []GenericTestEntry{
		{diagnostics.EventData{EventID: 1, Timestamp: types.NewDateTime(time.Now()), Trigger: diagnostics.EventTriggerAlerting, Cause: newInt(42), ActualValue: "someValue", TechCode: "742", TechInfo: "stacktrace", Cleared: false, TransactionID: "1234", VariableMonitoringID: newInt(99), EventNotificationType: diagnostics.EventPreconfiguredMonitor, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}, true},
		{diagnostics.EventData{EventID: 1, Timestamp: types.NewDateTime(time.Now()), Trigger: diagnostics.EventTriggerAlerting, Cause: newInt(42), ActualValue: "someValue", TechCode: "742", TechInfo: "stacktrace", Cleared: false, TransactionID: "1234", EventNotificationType: diagnostics.EventPreconfiguredMonitor, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}, true},
		{diagnostics.EventData{EventID: 1, Timestamp: types.NewDateTime(time.Now()), Trigger: diagnostics.EventTriggerAlerting, Cause: newInt(42), ActualValue: "someValue", TechCode: "742", TechInfo: "stacktrace", Cleared: false, EventNotificationType: diagnostics.EventPreconfiguredMonitor, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}, true},
		{diagnostics.EventData{EventID: 1, Timestamp: types.NewDateTime(time.Now()), Trigger: diagnostics.EventTriggerAlerting, Cause: newInt(42), ActualValue: "someValue", TechCode: "742", TechInfo: "stacktrace", EventNotificationType: diagnostics.EventPreconfiguredMonitor, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}, true},
		{diagnostics.EventData{EventID: 1, Timestamp: types.NewDateTime(time.Now()), Trigger: diagnostics.EventTriggerAlerting, Cause: newInt(42), ActualValue: "someValue", TechCode: "742", EventNotificationType: diagnostics.EventPreconfiguredMonitor, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}, true},
		{diagnostics.EventData{EventID: 1, Timestamp: types.NewDateTime(time.Now()), Trigger: diagnostics.EventTriggerAlerting, Cause: newInt(42), ActualValue: "someValue", EventNotificationType: diagnostics.EventPreconfiguredMonitor, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}, true},
		{diagnostics.EventData{Timestamp: types.NewDateTime(time.Now()), Trigger: diagnostics.EventTriggerAlerting, Cause: newInt(42), ActualValue: "someValue", EventNotificationType: diagnostics.EventPreconfiguredMonitor, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}, true},
		{diagnostics.EventData{Timestamp: types.NewDateTime(time.Now()), Trigger: diagnostics.EventTriggerAlerting, ActualValue: "someValue", EventNotificationType: diagnostics.EventPreconfiguredMonitor, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}, true},
		{diagnostics.EventData{Trigger: diagnostics.EventTriggerAlerting, ActualValue: "someValue", EventNotificationType: diagnostics.EventPreconfiguredMonitor, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}, false},
		{diagnostics.EventData{Timestamp: types.NewDateTime(time.Now()), ActualValue: "someValue", EventNotificationType: diagnostics.EventPreconfiguredMonitor, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}, false},
		{diagnostics.EventData{Timestamp: types.NewDateTime(time.Now()), Trigger: diagnostics.EventTriggerAlerting, EventNotificationType: diagnostics.EventPreconfiguredMonitor, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}, false},
		{diagnostics.EventData{Timestamp: types.NewDateTime(time.Now()), Trigger: diagnostics.EventTriggerAlerting, Cause: newInt(42), ActualValue: "someValue", Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}, false},
		{diagnostics.EventData{Timestamp: types.NewDateTime(time.Now()), Trigger: diagnostics.EventTriggerAlerting, Cause: newInt(42), ActualValue: "someValue", EventNotificationType: diagnostics.EventPreconfiguredMonitor, Component: types.Component{Name: "component1"}}, false},
		{diagnostics.EventData{Timestamp: types.NewDateTime(time.Now()), Trigger: diagnostics.EventTriggerAlerting, Cause: newInt(42), ActualValue: "someValue", EventNotificationType: diagnostics.EventPreconfiguredMonitor, Variable: types.Variable{Name: "variable1"}}, false},
		{diagnostics.EventData{EventID: -1, Timestamp: types.NewDateTime(time.Now()), Trigger: diagnostics.EventTriggerAlerting, Cause: newInt(42), ActualValue: "someValue", TechCode: "742", TechInfo: "stacktrace", Cleared: false, TransactionID: "1234", VariableMonitoringID: newInt(99), EventNotificationType: diagnostics.EventPreconfiguredMonitor, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}, false},
		{diagnostics.EventData{EventID: 1, Timestamp: types.NewDateTime(time.Now()), Trigger: "invalidEventTrigger", Cause: newInt(42), ActualValue: "someValue", TechCode: "742", TechInfo: "stacktrace", Cleared: false, TransactionID: "1234", VariableMonitoringID: newInt(99), EventNotificationType: diagnostics.EventPreconfiguredMonitor, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}, false},
		{diagnostics.EventData{EventID: 1, Timestamp: types.NewDateTime(time.Now()), Trigger: diagnostics.EventTriggerAlerting, Cause: newInt(42), ActualValue: ">2500................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................", TechCode: "742", TechInfo: "stacktrace", Cleared: false, TransactionID: "1234", VariableMonitoringID: newInt(99), EventNotificationType: diagnostics.EventPreconfiguredMonitor, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}, false},
		{diagnostics.EventData{EventID: 1, Timestamp: types.NewDateTime(time.Now()), Trigger: diagnostics.EventTriggerAlerting, Cause: newInt(42), ActualValue: "someValue", TechCode: ">50................................................", TechInfo: "stacktrace", Cleared: false, TransactionID: "1234", VariableMonitoringID: newInt(99), EventNotificationType: diagnostics.EventPreconfiguredMonitor, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}, false},
		{diagnostics.EventData{EventID: 1, Timestamp: types.NewDateTime(time.Now()), Trigger: diagnostics.EventTriggerAlerting, Cause: newInt(42), ActualValue: "someValue", TechCode: "742", TechInfo: ">500.................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................", Cleared: false, TransactionID: "1234", VariableMonitoringID: newInt(99), EventNotificationType: diagnostics.EventPreconfiguredMonitor, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}, false},
		{diagnostics.EventData{EventID: 1, Timestamp: types.NewDateTime(time.Now()), Trigger: diagnostics.EventTriggerAlerting, Cause: newInt(42), ActualValue: "someValue", TechCode: "742", TechInfo: "stacktrace", Cleared: false, TransactionID: ">36..................................", VariableMonitoringID: newInt(99), EventNotificationType: "invalidEventNotification", Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: "variable1"}}, false},
		{diagnostics.EventData{EventID: 1, Timestamp: types.NewDateTime(time.Now()), Trigger: diagnostics.EventTriggerAlerting, Cause: newInt(42), ActualValue: "someValue", TechCode: "742", TechInfo: "stacktrace", Cleared: false, TransactionID: "1234", VariableMonitoringID: newInt(99), EventNotificationType: diagnostics.EventPreconfiguredMonitor, Component: types.Component{Name: ">50................................................"}, Variable: types.Variable{Name: "variable1"}}, false},
		{diagnostics.EventData{EventID: 1, Timestamp: types.NewDateTime(time.Now()), Trigger: diagnostics.EventTriggerAlerting, Cause: newInt(42), ActualValue: "someValue", TechCode: "742", TechInfo: "stacktrace", Cleared: false, TransactionID: "1234", VariableMonitoringID: newInt(99), EventNotificationType: diagnostics.EventPreconfiguredMonitor, Component: types.Component{Name: "component1"}, Variable: types.Variable{Name: ">50................................................"}}, false},
	}
	ExecuteGenericTestTable(suite, table)
}

func (suite *OcppV2TestSuite) TestNotifyEventResponseValidation() {
	var responseTable = []GenericTestEntry{
		{diagnostics.NotifyEventResponse{}, true},
	}
	ExecuteGenericTestTable(suite, responseTable)
}

func (suite *OcppV2TestSuite) TestNotifyEventE2EMocked() {
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	tbc := true
	seqNo := 0
	generatedAt := types.NewDateTime(time.Now())
	eventData := diagnostics.EventData{
		EventID:               1,
		Timestamp:             types.NewDateTime(time.Now()),
		Trigger:               diagnostics.EventTriggerAlerting,
		Cause:                 newInt(42),
		ActualValue:           "someValue",
		TechCode:              "742",
		TechInfo:              "stacktrace",
		Cleared:               true,
		TransactionID:         "1234",
		VariableMonitoringID:  newInt(99),
		EventNotificationType: diagnostics.EventPreconfiguredMonitor,
		Component:             types.Component{Name: "component1"},
		Variable:              types.Variable{Name: "variable1"},
	}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"generatedAt":"%v","seqNo":%v,"tbc":%v,"eventData":[{"eventId":%v,"timestamp":"%v","trigger":"%v","cause":%v,"actualValue":"%v","techCode":"%v","techInfo":"%v","cleared":%v,"transactionId":"%v","variableMonitoringId":%v,"eventNotificationType":"%v","component":{"name":"%v"},"variable":{"name":"%v"}}]}]`,
		messageId, diagnostics.NotifyEventFeatureName, generatedAt.FormatTimestamp(), seqNo, tbc, eventData.EventID, eventData.Timestamp.FormatTimestamp(), eventData.Trigger, *eventData.Cause, eventData.ActualValue, eventData.TechCode, eventData.TechInfo, eventData.Cleared, eventData.TransactionID, *eventData.VariableMonitoringID, eventData.EventNotificationType, eventData.Component.Name, eventData.Variable.Name)
	responseJson := fmt.Sprintf(`[3,"%v",{}]`, messageId)
	response := diagnostics.NewNotifyEventResponse()
	channel := NewMockWebSocket(wsId)

	handler := &MockCSMSDiagnosticsHandler{}
	handler.On("OnNotifyEvent", mock.AnythingOfType("string"), mock.Anything).Return(response, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*diagnostics.NotifyEventRequest)
		suite.Require().True(ok)
		assertDateTimeEquality(suite, generatedAt, request.GeneratedAt)
		suite.Equal(tbc, request.Tbc)
		suite.Equal(seqNo, request.SeqNo)
		suite.Require().Len(request.EventData, 1)
		suite.Equal(eventData.EventID, request.EventData[0].EventID)
		assertDateTimeEquality(suite, eventData.Timestamp, request.EventData[0].Timestamp)
		suite.Equal(eventData.Trigger, request.EventData[0].Trigger)
		suite.Equal(*eventData.Cause, *request.EventData[0].Cause)
		suite.Equal(eventData.ActualValue, request.EventData[0].ActualValue)
		suite.Equal(eventData.TechCode, request.EventData[0].TechCode)
		suite.Equal(eventData.TechInfo, request.EventData[0].TechInfo)
		suite.Equal(eventData.Cleared, request.EventData[0].Cleared)
		suite.Equal(eventData.TransactionID, request.EventData[0].TransactionID)
		suite.Equal(*eventData.VariableMonitoringID, *request.EventData[0].VariableMonitoringID)
		suite.Equal(eventData.EventNotificationType, request.EventData[0].EventNotificationType)
		suite.Equal(eventData.Component.Name, request.EventData[0].Component.Name)
		suite.Equal(eventData.Variable.Name, request.EventData[0].Variable.Name)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	suite.Require().Nil(err)
	r, err := suite.chargingStation.NotifyEvent(generatedAt, seqNo, []diagnostics.EventData{eventData}, func(request *diagnostics.NotifyEventRequest) {
		request.Tbc = tbc
	})
	suite.Nil(err)
	suite.NotNil(r)
}

func (suite *OcppV2TestSuite) TestNotifyEventInvalidEndpoint() {
	messageId := defaultMessageId
	tbc := false
	seqNo := 0
	generatedAt := types.NewDateTime(time.Now())
	eventData := diagnostics.EventData{
		EventID:               1,
		Timestamp:             types.NewDateTime(time.Now()),
		Trigger:               diagnostics.EventTriggerAlerting,
		Cause:                 newInt(42),
		ActualValue:           "someValue",
		TechCode:              "742",
		TechInfo:              "stacktrace",
		Cleared:               true,
		TransactionID:         "1234",
		VariableMonitoringID:  newInt(99),
		EventNotificationType: diagnostics.EventPreconfiguredMonitor,
		Component:             types.Component{Name: "component1"},
		Variable:              types.Variable{Name: "variable1"},
	}
	req := diagnostics.NewNotifyEventRequest(generatedAt, seqNo, []diagnostics.EventData{eventData})
	req.Tbc = tbc
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"generatedAt":"%v","seqNo":%v,"tbc":%v,"eventData":[{"eventId":%v,"timestamp":"%v","trigger":"%v","cause":%v,"actualValue":"%v","techCode":"%v","techInfo":"%v","cleared":%v,"transactionId":"%v","variableMonitoringId":%v,"eventNotificationType":"%v","component":{"name":"%v"},"variable":{"name":"%v"}}]}]`,
		messageId, diagnostics.NotifyEventFeatureName, generatedAt.FormatTimestamp(), seqNo, tbc, eventData.EventID, eventData.Timestamp.FormatTimestamp(), eventData.Trigger, *eventData.Cause, eventData.ActualValue, eventData.TechCode, eventData.TechInfo, eventData.Cleared, eventData.TransactionID, *eventData.VariableMonitoringID, eventData.EventNotificationType, eventData.Component.Name, eventData.Variable.Name)
	testUnsupportedRequestFromCentralSystem(suite, req, requestJson, messageId)
}
