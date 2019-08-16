package ocpp16

import (
	"github.com/lorenzodonini/go-ocpp/ocpp"
)

const (
	BootNotificationFeatureName       = "BootNotification"
	AuthorizeFeatureName              = "Authorize"
	ChangeAvailabilityFeatureName     = "ChangeAvailability"
	ChangeConfigurationFeatureName    = "ChangeConfiguration"
	DataTransferFeatureName           = "DataTransfer"
	GetConfigurationFeatureName       = "GetConfiguration"
	ClearCacheFeatureName             = "ClearCache"
	HeartbeatFeatureName              = "Heartbeat"
	RemoteStartTransactionFeatureName = "RemoteStartTransaction"
	ResetFeatureName                  = "Reset"
	StartTransactionFeatureName       = "StartTransaction"
	StatusNotificationFeatureName     = "StatusNotification"
	UnlockConnectorFeatureName        = "UnlockConnector"
)

type CentralSystemCoreListener interface {
	OnAuthorize(chargePointId string, request *AuthorizeRequest) (confirmation *AuthorizeConfirmation, err error)
	OnBootNotification(chargePointId string, request *BootNotificationRequest) (confirmation *BootNotificationConfirmation, err error)
	OnDataTransfer(chargePointId string, request *DataTransferRequest) (confirmation *DataTransferConfirmation, err error)
	//onDiagnosticsStatusNotification()
	//onFirmwareStatusNotification()
	OnHeartbeat(chargePointId string, request *HeartbeatRequest) (confirmation *HeartbeatConfirmation, err error)
	//onMeterValues()
	OnStatusNotification(chargePointId string, request *StatusNotificationRequest) (confirmation *StatusNotificationConfirmation, err error)
	OnStartTransaction(chargePointId string, request *StartTransactionRequest) (confirmation *StartTransactionConfirmation, err error)
	//onStopTransaction()
}

type ChargePointCoreListener interface {
	//onCancelReservation()
	OnChangeAvailability(request *ChangeAvailabilityRequest) (confirmation *ChangeAvailabilityConfirmation, err error)
	OnChangeConfiguration(request *ChangeConfigurationRequest) (confirmation *ChangeConfigurationConfirmation, err error)
	OnClearCache(request *ClearCacheRequest) (confirmation *ClearCacheConfirmation, err error)
	//onClearChargingProfile()
	OnDataTransfer(request *DataTransferRequest) (confirmation *DataTransferConfirmation, err error)
	//onGetCompositeSchedule()
	OnGetConfiguration(request *GetConfigurationRequest) (confirmation *GetConfigurationConfirmation, err error)
	//onGetDiagnostics()
	//onGetLocalListVersion()
	OnRemoteStartTransaction(request *RemoteStartTransactionRequest) (confirmation *RemoteStartTransactionConfirmation, err error)
	//onRemoteStopTransaction()
	//onReserveNow()
	OnReset(request *ResetRequest) (confirmation *ResetConfirmation, err error)
	//onSendLocalList()
	//onSetChargingProfile()
	//onTriggerMessage()
	OnUnlockConnector(request *UnlockConnectorRequest) (confirmation *UnlockConnectorConfirmation, err error)
	//onUpdateFirmware()
}

var CoreProfile = ocpp.NewProfile(
	"core",
	BootNotificationFeature{},
	AuthorizeFeature{},
	ChangeAvailabilityFeature{},
	ChangeConfigurationFeature{},
	ClearCacheFeature{},
	DataTransferFeature{},
	GetConfigurationFeature{},
	HeartbeatFeature{},
	RemoteStartTransactionFeature{},
	StartTransactionFeature{},
	StatusNotificationFeature{},
	ResetFeature{},
	UnlockConnectorFeature{})
