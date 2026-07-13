// The display functional block contains OCPP 2.1 features for managing message that get displayed on a charging station.
package display

import "github.com/lorenzodonini/ocpp-go/ocpp"

// Needs to be implemented by a CSMS for handling messages part of the OCPP 2.1 Display profile.
type CSMSHandler interface {
	// OnNotifyDisplayMessages is called on the CSMS whenever a NotifyDisplayMessagesRequest is received from a Charging Station.
	OnNotifyDisplayMessages(chargingStationID string, request *NotifyDisplayMessagesRequest) (response *NotifyDisplayMessagesResponse, err error)
}

// Needs to be implemented by Charging stations for handling messages part of the OCPP 2.1 Display profile.
type ChargingStationHandler interface {
	// OnClearDisplay is called on a charging station whenever a ClearDisplayMessageRequest is received from the CSMS.
	OnClearDisplay(request *ClearDisplayMessageRequest) (confirmation *ClearDisplayResponse, err error)
	// OnGetDisplayMessages is called on a charging station whenever a GetDisplayMessagesRequest is received from the CSMS.
	OnGetDisplayMessages(request *GetDisplayMessagesRequest) (confirmation *GetDisplayMessagesResponse, err error)
	// OnSetDisplayMessage is called on a charging station whenever a SetDisplayMessageRequest is received from the CSMS.
	OnSetDisplayMessage(request *SetDisplayMessageRequest) (response *SetDisplayMessageResponse, err error)
}

const ProfileName = "Display"

var Profile = ocpp.NewProfile(
	ProfileName,
	ClearDisplayFeature{},
	GetDisplayMessagesFeature{},
	NotifyDisplayMessagesFeature{},
	SetDisplayMessageFeature{},
)
