package ocpp2

import (
	"github.com/lorenzodonini/ocpp-go/ocpp"
)

const (
	BootNotificationFeatureName           = "BootNotification"
	AuthorizeFeatureName                  = "Authorize"
	CancelReservationFeatureName          = "CancelReservation"
	CertificateSignedFeatureName          = "CertificateSigned"
	ChangeAvailabilityFeatureName         = "ChangeAvailability"
	ClearCacheFeatureName                 = "ClearCache"
	ClearDisplayFeatureName               = "ClearDisplay"
	ClearChargingProfileFeatureName       = "ClearChargingProfile"
	ClearedChargingLimitFeatureName       = "ClearedChargingLimit"
	ClearVariableMonitoringFeatureName    = "ClearVariableMonitoring"
	CostUpdatedFeatureName                = "CostUpdated"
	CustomerInformationFeatureName        = "CustomerInformation"
	DataTransferFeatureName               = "DataTransfer"
	DeleteCertificateFeatureName          = "DeleteCertificate"
	FirmwareStatusNotificationFeatureName = "FirmwareStatusNotification"
	Get15118EVCertificateFeatureName      = "Get15118EVCertificate"
	GetBaseReportFeatureName              = "GetBaseReport"
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
	// OnCancelReservation is called on the CSMS whenever an AuthorizeRequest is received from a charging station.
	OnAuthorize(chargingStationID string, request *AuthorizeRequest) (confirmation *AuthorizeConfirmation, err error)
	// OnBootNotification is called on the CSMS whenever a BootNotificationRequest is received from a charging station.
	OnBootNotification(chargingStationID string, request *BootNotificationRequest) (confirmation *BootNotificationConfirmation, err error)
	// OnClearedChargingLimit is called on the CSMS whenever a ClearedChargingLimitRequest is received from a charging station.
	OnClearedChargingLimit(chargingStationID string, request *ClearedChargingLimitRequest) (confirmation *ClearedChargingLimitConfirmation, err error)
	// OnDataTransfer is called on the CSMS whenever a DataTransferRequest is received from a charging station.
	OnDataTransfer(chargingStationID string, request *DataTransferRequest) (confirmation *DataTransferConfirmation, err error)
	// OnFirmwareStatusNotification is called on the CSMS whenever a FirmwareStatusNotificationRequest is received from a charging station.
	OnFirmwareStatusNotification(chargingStationID string, request *FirmwareStatusNotificationRequest) (confirmation *FirmwareStatusNotificationConfirmation, err error)
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
	// OnCancelReservation is called on a charging station whenever a CancelReservationRequest is received from the CSMS.
	OnCancelReservation(request *CancelReservationRequest) (confirmation *CancelReservationConfirmation, err error)
	// OnCertificateSigned is called on a charging station whenever a CertificateSignedRequest is received from the CSMS.
	OnCertificateSigned(request *CertificateSignedRequest) (confirmation *CertificateSignedConfirmation, err error)
	// OnChangeAvailability is called on a charging station whenever a ChangeAvailabilityRequest is received from the CSMS.
	OnChangeAvailability(request *ChangeAvailabilityRequest) (confirmation *ChangeAvailabilityConfirmation, err error)
	// OnClearCache is called on a charging station whenever a ClearCacheRequest is received from the CSMS.
	OnClearCache(request *ClearCacheRequest) (confirmation *ClearCacheConfirmation, err error)
	// OnClearDisplay is called on a charging station whenever a ClearDisplayRequest is received from the CSMS.
	OnClearDisplay(request *ClearDisplayRequest) (confirmation *ClearDisplayConfirmation, err error)
	// OnClearChargingProfile is called on a charging station whenever a ClearChargingProfileRequest is received from the CSMS.
	OnClearChargingProfile(request *ClearChargingProfileRequest) (confirmation *ClearChargingProfileConfirmation, err error)
	// OnClearVariableMonitoring is called on a charging station whenever a ClearVariableMonitoringRequest is received from the CSMS.
	OnClearVariableMonitoring(request *ClearVariableMonitoringRequest) (confirmation *ClearVariableMonitoringConfirmation, err error)
	// OnCostUpdated is called on a charging station whenever a CostUpdatedRequest is received from the CSMS.
	OnCostUpdated(request *CostUpdatedRequest) (confirmation *CostUpdatedConfirmation, err error)
	// OnCustomerInformation is called on a charging station whenever a CustomerInformationRequest is received from the CSMS.
	OnCustomerInformation(request *CustomerInformationRequest) (confirmation *CustomerInformationConfirmation, err error)
	// OnDataTransfer is called on a charging station whenever a DataTransferRequest is received from the CSMS.
	OnDataTransfer(request *DataTransferRequest) (confirmation *DataTransferConfirmation, err error)
	// OnDeleteCertificate is called on a charging station whenever a DeleteCertificateRequest is received from the CSMS.
	OnDeleteCertificate(request *DeleteCertificateRequest) (confirmation *DeleteCertificateConfirmation, err error)
	// OnGetBaseReport is called on a charging station whenever a GetBaseReportRequest is received from the CSMS.
	OnGetBaseReport(request *GetBaseReportRequest) (confirmation *GetBaseReportConfirmation, err error)
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
	BootNotificationFeature{},
	AuthorizeFeature{},
	CancelReservationFeature{},
	CertificateSignedFeature{},
	ChangeAvailabilityFeature{},
	ClearCacheFeature{},
	ClearDisplayFeature{},
	ClearChargingProfileFeature{},
	ClearedChargingLimitFeature{},
	ClearVariableMonitoringFeature{},
	CostUpdatedFeature{},
	CustomerInformationFeature{},
	DataTransferFeature{},
	DeleteCertificateFeature{},
	FirmwareStatusNotificationFeature{},
	Get15118EVCertificateFeature{},
	GetBaseReportFeature{},
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
