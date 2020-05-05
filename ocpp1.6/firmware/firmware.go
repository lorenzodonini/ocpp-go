// Contains support for firmware update management and diagnostic log file download.
package firmware

import (
	"github.com/lorenzodonini/ocpp-go/ocpp"
)

const (
	GetDiagnosticsFeatureName                = "GetDiagnostics"
	DiagnosticsStatusNotificationFeatureName = "DiagnosticsStatusNotification"
	FirmwareStatusNotificationFeatureName    = "FirmwareStatusNotification"
	UpdateFirmwareFeatureName                = "UpdateFirmware"
)

type CentralSystemFirmwareManagementHandler interface {
	OnDiagnosticsStatusNotification(chargePointId string, request *DiagnosticsStatusNotificationRequest) (confirmation *DiagnosticsStatusNotificationConfirmation, err error)
	OnFirmwareStatusNotification(chargePointId string, request *FirmwareStatusNotificationRequest) (confirmation *FirmwareStatusNotificationConfirmation, err error)
}

type ChargePointFirmwareManagementHandler interface {
	OnGetDiagnostics(request *GetDiagnosticsRequest) (confirmation *GetDiagnosticsConfirmation, err error)
	OnUpdateFirmware(request *UpdateFirmwareRequest) (confirmation *UpdateFirmwareConfirmation, err error)
}

const ProfileName = "firmwareManagement"

// Provides support for firmware update management and diagnostic log file download.
var Profile = ocpp.NewProfile(
	ProfileName,
	GetDiagnosticsFeature{},
	DiagnosticsStatusNotificationFeature{},
	FirmwareStatusNotificationFeature{},
	UpdateFirmwareFeature{})
