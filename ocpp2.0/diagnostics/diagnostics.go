// The diagnostics functional block contains OCPP 2.0 features than enable remote diagnostics of problems with a charging station.
package diagnostics

import "github.com/lorenzodonini/ocpp-go/ocpp"

// Needs to be implemented by a CSMS for handling messages part of the OCPP 2.0 Diagnostics profile.
type CSMSHandler interface {
	// OnLogStatusNotification is called on the CSMS whenever a LogStatusNotificationRequest is received from a Charging Station.
	OnLogStatusNotification(chargingStationID string, request *LogStatusNotificationRequest) (response *LogStatusNotificationResponse, err error)
	// OnNotifyCustomerInformation is called on the CSMS whenever a NotifyCustomerInformationRequest is received from a Charging Station.
	OnNotifyCustomerInformation(chargingStationID string, request *NotifyCustomerInformationRequest) (response *NotifyCustomerInformationResponse, err error)
	// OnNotifyEvent is called on the CSMS whenever a NotifyEventRequest is received from a Charging Station.
	OnNotifyEvent(chargingStationID string, request *NotifyEventRequest) (response *NotifyEventResponse, err error)
}

// Needs to be implemented by Charging stations for handling messages part of the OCPP 2.0 Diagnostics profile.
type ChargingStationHandler interface {
	// OnClearVariableMonitoring is called on a charging station whenever a ClearVariableMonitoringRequest is received from the CSMS.
	OnClearVariableMonitoring(request *ClearVariableMonitoringRequest) (response *ClearVariableMonitoringResponse, err error)
	// OnCustomerInformation is called on a charging station whenever a CustomerInformationRequest is received from the CSMS.
	OnCustomerInformation(request *CustomerInformationRequest) (response *CustomerInformationResponse, err error)
	// OnGetLog is called on a charging station whenever a GetLogRequest is received from the CSMS.
	OnGetLog(request *GetLogRequest) (response *GetLogResponse, err error)
	// OnGetMonitoringReport is called on a charging station whenever a GetMonitoringReportRequest is received from the CSMS.
	OnGetMonitoringReport(request *GetMonitoringReportRequest) (response *GetMonitoringReportResponse, err error)
}

const ProfileName = "diagnostics"

var Profile = ocpp.NewProfile(
	ProfileName,
	ClearVariableMonitoringFeature{},
	CustomerInformationFeature{},
	GetLogFeature{},
	GetMonitoringReportFeature{},
	LogStatusNotificationFeature{},
	NotifyCustomerInformationFeature{},
	NotifyEventFeature{},
)
