package main

import "github.com/lorenzodonini/ocpp-go/ocpp2.0.1/diagnostics"

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
