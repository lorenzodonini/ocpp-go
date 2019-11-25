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
	//onFirmwareStatusNotification()
}

type ChargePointFirmwareManagementListener interface {
	OnGetDiagnostics(request *GetDiagnosticsRequest) (confirmation *GetDiagnosticsConfirmation, err error)
	//onUpdateFirmware()
}

const FirmwareManagementProfileName = "firmwareManagement"

var FirmwareManagementProfile = ocpp.NewProfile(
	FirmwareManagementProfileName,
	GetDiagnosticsFeature{},
	DiagnosticsStatusNotificationFeature{})
