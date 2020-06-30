// The display functional block contains OCPP 2.0 features for managing message that get displayed on a charging station.
package display

import "github.com/lorenzodonini/ocpp-go/ocpp"

// Needs to be implemented by a CSMS for handling messages part of the OCPP 2.0 Display profile.
type CSMSHandler interface {
	// OnNotifyDisplayMessages is called on the CSMS whenever a NotifyDisplayMessagesRequest is received from a Charging Station.
	OnNotifyDisplayMessages(chargingStationID string, request *NotifyDisplayMessagesRequest) (response *NotifyDisplayMessagesResponse, err error)
}

// Needs to be implemented by Charging stations for handling messages part of the OCPP 2.0 Display profile.
type ChargingStationHandler interface {
	// OnClearDisplay is called on a charging station whenever a ClearDisplayRequest is received from the CSMS.
	OnClearDisplay(request *ClearDisplayRequest) (confirmation *ClearDisplayResponse, err error)
	// OnGetDisplayMessages is called on a charging station whenever a GetDisplayMessagesRequest is received from the CSMS.
	OnGetDisplayMessages(request *GetDisplayMessagesRequest) (confirmation *GetDisplayMessagesResponse, err error)
}

const ProfileName = "display"

var Profile = ocpp.NewProfile(
	ProfileName,
	ClearDisplayFeature{},
	GetDisplayMessagesFeature{},
	NotifyDisplayMessagesFeature{},
)
