// The diagnostics functional block contains OCPP 1.6J extension features than enable remote firmware updates on charging stations.
package securefirmware

import "github.com/lorenzodonini/ocpp-go/ocpp"

type CentralSystemHandler interface {
	OnSignedFirmwareStatusNotification(chargingStationID string, request *SignedFirmwareStatusNotificationRequest) (response *SignedFirmwareStatusNotificationResponse, err error)
}

// Needs to be implemented by Charging stations for handling messages part of the OCPP 1.6j security extension.
type ChargePointHandler interface {
	OnSignedUpdateFirmware(request *SignedUpdateFirmwareRequest) (response *SignedUpdateFirmwareResponse, err error)
}

const ProfileName = "SecureFirmwareUpdate"

var Profile = ocpp.NewProfile(
	ProfileName,
	SignedFirmwareStatusNotificationFeature{},
	SignedUpdateFirmwareFeature{},
)
