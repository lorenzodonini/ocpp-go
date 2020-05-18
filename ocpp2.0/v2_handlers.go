package ocpp2

import (
	"github.com/lorenzodonini/ocpp-go/ocpp"
)

const (
	Get15118EVCertificateFeatureName      = "Get15118EVCertificate"
	GetCertificateStatusFeatureName       = "GetCertificateStatus"
	GetChargingProfilesFeatureName        = "GetChargingProfiles"
	GetCompositeScheduleFeatureName       = "GetCompositeSchedule"
	GetDisplayMessagesFeatureName         = "GetDisplayMessages"
	GetInstalledCertificateIdsFeatureName = "GetInstalledCertificateIds"
	GetLocalListVersionFeatureName        = "GetLocalListVersion"
	GetLogFeatureName                     = "GetLog"
	GetMonitoringReportFeatureName        = "GetMonitoringReport"
	// GetConfigurationFeatureName       = "GetConfiguration"
	// HeartbeatFeatureName              = "Heartbeat"
	// MeterValuesFeatureName            = "MeterValues"
	// RemoteStartTransactionFeatureName = "RemoteStartTransaction"
	// RemoteStopTransactionFeatureName  = "RemoteStopTransaction"
	// ResetFeatureName                  = "Reset"
	// StartTransactionFeatureName       = "StartTransaction"
	// StopTransactionFeatureName        = "StopTransaction"
	// StatusNotificationFeatureName     = "StatusNotification"
	// UnlockConnectorFeatureName        = "UnlockConnector"
	//SetChargingProfileFeatureName   = "SetChargingProfile"
	//GetCompositeScheduleFeatureName = "GetCompositeSchedule"
)

type CSMSHandler interface {
	// OnGet15118EVCertificate is called on the CSMS whenever a Get15118EVCertificateRequest is received from a charging station.
	OnGet15118EVCertificate(chargingStationID string, request *Get15118EVCertificateRequest) (confirmation *Get15118EVCertificateConfirmation, err error)
	// OnGetCertificateStatus is called on the CSMS whenever a GetCertificateStatusRequest is received from a charging station.
	OnGetCertificateStatus(chargingStationID string, request *GetCertificateStatusRequest) (confirmation *GetCertificateStatusConfirmation, err error)
	// OnHeartbeat(chargePointId string, request *HeartbeatRequest) (confirmation *HeartbeatConfirmation, err error)
	// OnMeterValues(chargePointId string, request *MeterValuesRequest) (confirmation *MeterValuesConfirmation, err error)
	// OnStatusNotification(chargePointId string, request *StatusNotificationRequest) (confirmation *StatusNotificationConfirmation, err error)
	// OnStartTransaction(chargePointId string, request *StartTransactionRequest) (confirmation *StartTransactionConfirmation, err error)
	// OnStopTransaction(chargePointId string, request *StopTransactionRequest) (confirmation *StopTransactionConfirmation, err error)
}

type ChargingStationHandler interface {
	// OnGetChargingProfiles is called on a charging station whenever a GetChargingProfilesRequest is received from the CSMS.
	OnGetChargingProfiles(request *GetChargingProfilesRequest) (confirmation *GetChargingProfilesConfirmation, err error)
	// OnGetCompositeSchedule is called on a charging station whenever a GetCompositeScheduleRequest is received from the CSMS.
	OnGetCompositeSchedule(request *GetCompositeScheduleRequest) (confirmation *GetCompositeScheduleConfirmation, err error)
	// OnGetDisplayMessages is called on a charging station whenever a GetDisplayMessagesRequest is received from the CSMS.
	OnGetDisplayMessages(request *GetDisplayMessagesRequest) (confirmation *GetDisplayMessagesConfirmation, err error)
	// OnGetInstalledCertificateIds is called on a charging station whenever a GetInstalledCertificateIdsRequest is received from the CSMS.
	OnGetInstalledCertificateIds(request *GetInstalledCertificateIdsRequest) (confirmation *GetInstalledCertificateIdsConfirmation, err error)
	// OnGetLocalListVersion is called on a charging station whenever a GetLocalListVersionRequest is received from the CSMS.
	OnGetLocalListVersion(request *GetLocalListVersionRequest) (confirmation *GetLocalListVersionConfirmation, err error)
	// OnGetLog is called on a charging station whenever a GetLogRequest is received from the CSMS.
	OnGetLog(request *GetLogRequest) (confirmation *GetLogConfirmation, err error)
	// OnGetMonitoringReport is called on a charging station whenever a GetMonitoringReportRequest is received from the CSMS.
	OnGetMonitoringReport(request *GetMonitoringReportRequest) (confirmation *GetMonitoringReportConfirmation, err error)
	// OnGetConfiguration(request *GetConfigurationRequest) (confirmation *GetConfigurationConfirmation, err error)
	// OnRemoteStartTransaction(request *RemoteStartTransactionRequest) (confirmation *RemoteStartTransactionConfirmation, err error)
	// OnRemoteStopTransaction(request *RemoteStopTransactionRequest) (confirmation *RemoteStopTransactionConfirmation, err error)
	// OnReset(request *ResetRequest) (confirmation *ResetConfirmation, err error)
	// OnUnlockConnector(request *UnlockConnectorRequest) (confirmation *UnlockConnectorConfirmation, err error)
}

var CoreProfileName = "core"

var CoreProfile = ocpp.NewProfile(
	CoreProfileName,
	Get15118EVCertificateFeature{},
	GetCertificateStatusFeature{},
	GetChargingProfilesFeature{},
	GetCompositeScheduleFeature{},
	GetDisplayMessagesFeature{},
	GetInstalledCertificateIdsFeature{},
	GetLocalListVersionFeature{},
	GetLogFeature{},
	GetMonitoringReportFeature{},

//GetConfigurationFeature{},
//HeartbeatFeature{},
//MeterValuesFeature{},
//RemoteStartTransactionFeature{},
//RemoteStopTransactionFeature{},
//StartTransactionFeature{},
//StopTransactionFeature{},
//StatusNotificationFeature{},
//ResetFeature{},
//UnlockConnectorFeature{}
)
