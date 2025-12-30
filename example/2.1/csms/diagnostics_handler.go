package main

import (
	"github.com/xBlaz3kx/ocpp-go/ocpp2.1/diagnostics"
	"github.com/xBlaz3kx/ocpp-go/ocpp2.1/types"
)

func (c *CSMSHandler) OnLogStatusNotification(chargingStationID string, request *diagnostics.LogStatusNotificationRequest) (response *diagnostics.LogStatusNotificationResponse, err error) {
	logDefault(chargingStationID, request.GetFeatureName()).Infof("log upload status: %v", request.Status)
	response = diagnostics.NewLogStatusNotificationResponse()
	return
}

func (c *CSMSHandler) OnNotifyCustomerInformation(chargingStationID string, request *diagnostics.NotifyCustomerInformationRequest) (response *diagnostics.NotifyCustomerInformationResponse, err error) {
	logDefault(chargingStationID, request.GetFeatureName()).Infof("data report for request %v: %v", request.RequestID, request.Data)
	response = diagnostics.NewNotifyCustomerInformationResponse()
	return
}

func (c *CSMSHandler) OnNotifyEvent(chargingStationID string, request *diagnostics.NotifyEventRequest) (response *diagnostics.NotifyEventResponse, err error) {
	logDefault(chargingStationID, request.GetFeatureName()).Infof("report part %v for events:\n", request.SeqNo)
	for _, ed := range request.EventData {
		logDefault(chargingStationID, request.GetFeatureName()).Infof("%v", ed)
	}
	response = diagnostics.NewNotifyEventResponse()
	return
}

func (c *CSMSHandler) OnNotifyMonitoringReport(chargingStationID string, request *diagnostics.NotifyMonitoringReportRequest) (response *diagnostics.NotifyMonitoringReportResponse, err error) {
	logDefault(chargingStationID, request.GetFeatureName()).Infof("report part %v for monitored variables:\n", request.SeqNo)
	for _, md := range request.Monitor {
		logDefault(chargingStationID, request.GetFeatureName()).Infof("%v", md)
	}
	response = diagnostics.NewNotifyMonitoringReportResponse()
	return
}

func (c *CSMSHandler) OnOpenPeriodicEventStream(chargingStationID string, request *diagnostics.OpenPeriodicEventStreamRequest) (response *diagnostics.OpenPeriodicEventStreamResponse, err error) {
	logDefault(chargingStationID, request.GetFeatureName()).Infof("open periodic event stream: id=%d, variableMonitoringId=%d",
		request.ConstantStreamData.Id, request.ConstantStreamData.VariableMonitoringId)
	return diagnostics.NewOpenPeriodicEventStreamResponse(types.GenericStatusAccepted), nil
}

func (c *CSMSHandler) OnClosePeriodicEventStream(chargingStationID string, request *diagnostics.ClosePeriodicEventStreamRequest) (response *diagnostics.ClosePeriodicEventStreamResponse, err error) {
	logDefault(chargingStationID, request.GetFeatureName()).Infof("close periodic event stream: id=%d", request.Id)
	return diagnostics.NewClosePeriodicEventStreamResponse(), nil
}

func (c *CSMSHandler) OnNotifyPeriodicEventStream(chargingStationID string, request *diagnostics.NotifyPeriodicEventStream) {
	logDefault(chargingStationID, request.GetFeatureName()).Infof("periodic event stream notification: id=%d, pending=%d, baseTime=%v, dataCount=%d",
		request.ID, request.Pending, request.BaseTime, len(request.Data))
	for i, data := range request.Data {
		logDefault(chargingStationID, request.GetFeatureName()).Infof("  data[%d]: t=%v, v=%s", i, data.T, data.V)
	}
}
