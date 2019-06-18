package v16

import (
	"github.com/lorenzodonini/go-ocpp/ocpp"
)

const (
	BootNotificationFeatureName = "BootNotification"
	AuthorizeFeatureName = "Authorize"
	ChangeAvailabilityFeatureName = "ChangeAvailability"
)

type coreProfile struct {
	*ocpp.Profile
}

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

func (profile* coreProfile)CreateBootNotification(chargePointModel string, chargePointVendor string) *BootNotificationRequest {
	return &BootNotificationRequest{ChargePointModel: chargePointModel, ChargePointVendor: chargePointVendor}
}

func (profile* coreProfile)CreateAuthorization(idTag string) *AuthorizeRequest {
	return &AuthorizeRequest{IdTag: idTag}
}

func (profile* coreProfile)CreateChangeAvailability(connectorId int, availabilityType AvailabilityType) *ChangeAvailabilityRequest {
	return &ChangeAvailabilityRequest{ConnectorId: connectorId, Type: availabilityType}
}

var CoreProfile = coreProfile{
	ocpp.NewProfile("core", BootNotificationFeature{}),
}
