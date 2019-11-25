package ocpp16

import "github.com/lorenzodonini/ocpp-go/ocpp"

const (
	GetDiagnosticsFeatureName = "GetDiagnostics"
	DiagnosticsStatusNotificationFeatureName = "DiagnosticsStatusNotification"
	FirmwareStatusNotificationFeatureName = "FirmwareStatusNotification"
	UpdateFirmwareFeatureName = "UpdateFirmware"
)

type CentralSystemFirmwareManagementListener interface {
	OnDiagnosticsStatusNotification(chargePointId string, request *DiagnosticsStatusNotificationRequest) (confirmation *DiagnosticsStatusNotificationConfirmation, err error)
	OnFirmwareStatusNotification(chargePointId string, request *FirmwareStatusNotificationRequest) (confirmation *FirmwareStatusNotificationConfirmation, err error)
}

type ChargePointFirmwareManagementListener interface {
	OnGetDiagnostics(request *GetDiagnosticsRequest) (confirmation *GetDiagnosticsConfirmation, err error)
	OnUpdateFirmware(request *UpdateFirmwareRequest) (confirmation *UpdateFirmwareConfirmation, err error)
}

const FirmwareManagementProfileName = "firmwareManagement"

var FirmwareManagementProfile = ocpp.NewProfile(
	FirmwareManagementProfileName,
	GetDiagnosticsFeature{},
	DiagnosticsStatusNotificationFeature{},
	FirmwareStatusNotificationFeature{},
	UpdateFirmwareFeature{})
