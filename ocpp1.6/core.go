package ocpp16

import (
	"github.com/lorenzodonini/go-ocpp/ocppj"
)

const (
	BootNotificationFeatureName    = "BootNotification"
	AuthorizeFeatureName           = "Authorize"
	ChangeAvailabilityFeatureName  = "ChangeAvailability"
	ChangeConfigurationFeatureName = "ChangeConfiguration"
	DataTransferFeatureName        = "DataTransfer"
	GetConfigurationFeatureName    = "GetConfiguration"
	ClearCacheFeatureName          = "ClearCache"
	HeartbeatFeatureName           = "Heartbeat"
)

type CentralSystemCoreListener interface {
	OnAuthorize(chargePointId string, request *AuthorizeRequest) (confirmation *AuthorizeConfirmation, err error)
	OnBootNotification(chargePointId string, request *BootNotificationRequest) (confirmation *BootNotificationConfirmation, err error)
	OnDataTransfer(chargePointId string, request *DataTransferRequest) (confirmation *DataTransferConfirmation, err error)
	//onDiagnosticsStatusNotification()
	//onFirmwareStatusNotification()
	OnHeartbeat(chargePointId string, request *HeartbeatRequest) (confirmation *HeartbeatConfirmation, err error)
	//onMeterValues()
	//onStatusNotification()
	//onStartTransaction()
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
	//onRemoteStartTransaction()
	//onRemoteStopTransaction()
	//onReserveNow()
	//onReset()
	//onSendLocalList()
	//onSetChargingProfile()
	//onTriggerMessage()
	//onUnlockConnector()
	//onUpdateFirmware()
}

var CoreProfile = ocppj.NewProfile("core", BootNotificationFeature{}, AuthorizeFeature{}, ChangeAvailabilityFeature{}, ChangeConfigurationFeature{}, ClearCacheFeature{}, DataTransferFeature{}, GetConfigurationFeature{}, HeartbeatFeature{})
