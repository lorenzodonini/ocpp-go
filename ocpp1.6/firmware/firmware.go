// Contains support for firmware update management and diagnostic log file download.
package firmware

import (
	"github.com/lorenzodonini/ocpp-go/ocpp"
)

// Needs to be implemented by Central systems for handling messages part of the OCPP 1.6 FirmwareManagement profile.
type CentralSystemHandler interface {
	OnDiagnosticsStatusNotification(chargePointId string, request *DiagnosticsStatusNotificationRequest) (confirmation *DiagnosticsStatusNotificationConfirmation, err error)
	OnFirmwareStatusNotification(chargePointId string, request *FirmwareStatusNotificationRequest) (confirmation *FirmwareStatusNotificationConfirmation, err error)
}

// Needs to be implemented by Charge points for handling messages part of the OCPP 1.6 FirmwareManagement profile.
type ChargePointHandler interface {
	OnGetDiagnostics(request *GetDiagnosticsRequest) (confirmation *GetDiagnosticsConfirmation, err error)
	OnUpdateFirmware(request *UpdateFirmwareRequest) (confirmation *UpdateFirmwareConfirmation, err error)
}

// The profile name
const ProfileName = "firmwareManagement"

// Provides support for firmware update management and diagnostic log file download.
var Profile = ocpp.NewProfile(
	ProfileName,
	GetDiagnosticsFeature{},
	DiagnosticsStatusNotificationFeature{},
	FirmwareStatusNotificationFeature{},
	UpdateFirmwareFeature{})
