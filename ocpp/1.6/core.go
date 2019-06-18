package v16

import (
	"github.com/lorenzodonini/go-ocpp/ocpp"
)

const (
	BootNotificationFeatureName = "BootNotification"
	AuthorizeFeatureName = "Authorize"
	ChangeAvailabilityFeatureName = "ChangeAvailability"
)

type CentralSystemCoreListener interface {
	OnAuthorize(chargePointId string, request *AuthorizeRequest) (confirmation *AuthorizeConfirmation, err error)
	OnBootNotification(chargePointId string, request *BootNotificationRequest) (confirmation *BootNotificationConfirmation, err error)
	//onDataTransfer()
	//onDiagnosticsStatusNotification()
	//onFirmwareStatusNotification()
	//onHeartbeat()
	//onMeterValues()
	//onStatusNotification()
	//onStartTransaction()
	//onStopTransaction()
}

type ChargePointCoreListener interface {
	//onCancelReservation()
	OnChangeAvailability(request *ChangeAvailabilityRequest) (confirmation *ChangeAvailabilityConfirmation, err error)
	//onChangeConfiguration()
	//onClearCache()
	//onClearChargingProfile()
	//onDataTransfer()
	//onGetCompositeSchedule()
	//onGetConfiguration()
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

var CoreProfile = ocpp.NewProfile("core", BootNotificationFeature{}, AuthorizeFeature{})
