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

type CentralSystemCoreListener interface {
	OnAuthorize(chargePointId string, request *AuthorizeRequest) (confirmation *AuthorizeConfirmation, err error)
	OnBootNotification(chargePointId string, request *BootNotificationRequest) (confirmation *BootNotificationConfirmation, err error)
	OnClearedChargingLimit(chargePointId string, request *ClearedChargingLimitRequest) (confirmation *ClearedChargingLimitConfirmation, err error)
	OnDataTransfer(chargePointId string, request *DataTransferRequest) (confirmation *DataTransferConfirmation, err error)
	OnFirmwareStatusNotification(chargePointId string, request *FirmwareStatusNotificationRequest) (confirmation *FirmwareStatusNotificationConfirmation, err error)
	OnGet15118EVCertificate(chargePointId string, request *Get15118EVCertificateRequest) (confirmation *Get15118EVCertificateConfirmation, err error)
	OnGetCertificateStatus(chargePointId string, request *GetCertificateStatusRequest) (confirmation *GetCertificateStatusConfirmation, err error)
	// OnHeartbeat(chargePointId string, request *HeartbeatRequest) (confirmation *HeartbeatConfirmation, err error)
	// OnMeterValues(chargePointId string, request *MeterValuesRequest) (confirmation *MeterValuesConfirmation, err error)
	// OnStatusNotification(chargePointId string, request *StatusNotificationRequest) (confirmation *StatusNotificationConfirmation, err error)
	// OnStartTransaction(chargePointId string, request *StartTransactionRequest) (confirmation *StartTransactionConfirmation, err error)
	// OnStopTransaction(chargePointId string, request *StopTransactionRequest) (confirmation *StopTransactionConfirmation, err error)
}

type ChargePointCoreListener interface {
	OnCancelReservation(request *CancelReservationRequest) (confirmation *CancelReservationConfirmation, err error)
	OnCertificateSigned(request *CertificateSignedRequest) (confirmation *CertificateSignedConfirmation, err error)
	OnChangeAvailability(request *ChangeAvailabilityRequest) (confirmation *ChangeAvailabilityConfirmation, err error)
	// OnChangeConfiguration(request *ChangeConfigurationRequest) (confirmation *ChangeConfigurationConfirmation, err error)
	OnClearCache(request *ClearCacheRequest) (confirmation *ClearCacheConfirmation, err error)
	OnClearDisplay(request *ClearDisplayRequest) (confirmation *ClearDisplayConfirmation, err error)
	OnClearChargingProfile(request *ClearChargingProfileRequest) (confirmation *ClearChargingProfileConfirmation, err error)
	OnClearVariableMonitoring(request *ClearVariableMonitoringRequest) (confirmation *ClearVariableMonitoringConfirmation, err error)
	OnCostUpdated(request *CostUpdatedRequest) (confirmation *CostUpdatedConfirmation, err error)
	OnCustomerInformation(request *CustomerInformationRequest) (confirmation *CustomerInformationConfirmation, err error)
	OnDataTransfer(request *DataTransferRequest) (confirmation *DataTransferConfirmation, err error)
	OnDeleteCertificate(request *DeleteCertificateRequest) (confirmation *DeleteCertificateConfirmation, err error)
	OnGetBaseReport(request *GetBaseReportRequest) (confirmation *GetBaseReportConfirmation, err error)
	OnGetChargingProfiles(request *GetChargingProfilesRequest) (confirmation *GetChargingProfilesConfirmation, err error)
	OnGetCompositeSchedule(request *GetCompositeScheduleRequest) (confirmation *GetCompositeScheduleConfirmation, err error)
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
	//ChangeConfigurationFeature{},
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
