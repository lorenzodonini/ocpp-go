// The firmware functional block contains OCPP 2.0 features that enable firmware updates on a charging station.
package firmware

import "github.com/lorenzodonini/ocpp-go/ocpp"

// Needs to be implemented by a CSMS for handling messages part of the OCPP 2.0 Firmware profile.
type CSMSHandler interface {
	// OnFirmwareStatusNotification is called on the CSMS whenever a FirmwareStatusNotificationRequest is received from a charging station.
	OnFirmwareStatusNotification(chargingStationID string, request *FirmwareStatusNotificationRequest) (response *FirmwareStatusNotificationResponse, err error)
	// OnPublishFirmwareStatusNotification is called on the CSMS whenever a PublishFirmwareStatusNotificationRequest is received from a local controller.
	OnPublishFirmwareStatusNotification(chargingStationID string, request *PublishFirmwareStatusNotificationRequest) (response *PublishFirmwareStatusNotificationResponse, err error)
}

// Needs to be implemented by Charging stations for handling messages part of the OCPP 2.0 Firmware profile.
type ChargingStationHandler interface {
	// OnPublishFirmware is called on a charging station whenever a PublishFirmwareRequest is received from the CSMS.
	OnPublishFirmware(request *PublishFirmwareRequest) (response *PublishFirmwareResponse, err error)
}

const ProfileName = "firmware"

var Profile = ocpp.NewProfile(
	ProfileName,
	FirmwareStatusNotificationFeature{},
	PublishFirmwareFeature{},
	PublishFirmwareStatusNotificationFeature{},
)
