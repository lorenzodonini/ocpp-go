package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/diagnostics"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"time"
)

// Test
func (suite *OcppV2TestSuite) TestNotifyEventRequestValidation() {
	t := suite.T()
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
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestNotifyEventDataValidation() {
	t := suite.T()
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
	ExecuteGenericTestTable(t, table)
}

func (suite *OcppV2TestSuite) TestNotifyEventResponseValidation() {
	t := suite.T()
	var responseTable = []GenericTestEntry{
		{diagnostics.NotifyEventResponse{}, true},
	}
	ExecuteGenericTestTable(t, responseTable)
}

func (suite *OcppV2TestSuite) TestNotifyEventE2EMocked() {
	t := suite.T()
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

	handler := MockCSMSDiagnosticsHandler{}
	handler.On("OnNotifyEvent", mock.AnythingOfType("string"), mock.Anything).Return(response, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*diagnostics.NotifyEventRequest)
		require.True(t, ok)
		assertDateTimeEquality(t, generatedAt, request.GeneratedAt)
		assert.Equal(t, tbc, request.Tbc)
		assert.Equal(t, seqNo, request.SeqNo)
		require.Len(t, request.EventData, 1)
		assert.Equal(t, eventData.EventID, request.EventData[0].EventID)
		assertDateTimeEquality(t, eventData.Timestamp, request.EventData[0].Timestamp)
		assert.Equal(t, eventData.Trigger, request.EventData[0].Trigger)
		assert.Equal(t, *eventData.Cause, *request.EventData[0].Cause)
		assert.Equal(t, eventData.ActualValue, request.EventData[0].ActualValue)
		assert.Equal(t, eventData.TechCode, request.EventData[0].TechCode)
		assert.Equal(t, eventData.TechInfo, request.EventData[0].TechInfo)
		assert.Equal(t, eventData.Cleared, request.EventData[0].Cleared)
		assert.Equal(t, eventData.TransactionID, request.EventData[0].TransactionID)
		assert.Equal(t, *eventData.VariableMonitoringID, *request.EventData[0].VariableMonitoringID)
		assert.Equal(t, eventData.EventNotificationType, request.EventData[0].EventNotificationType)
		assert.Equal(t, eventData.Component.Name, request.EventData[0].Component.Name)
		assert.Equal(t, eventData.Variable.Name, request.EventData[0].Variable.Name)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	r, err := suite.chargingStation.NotifyEvent(generatedAt, seqNo, []diagnostics.EventData{eventData}, func(request *diagnostics.NotifyEventRequest) {
		request.Tbc = tbc
	})
	assert.Nil(t, err)
	assert.NotNil(t, r)
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
