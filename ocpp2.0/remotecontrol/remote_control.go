// The Remote control functional block contains OCPP 2.0 features for remote-control management from the CSMS.
package remotecontrol

import "github.com/lorenzodonini/ocpp-go/ocpp"

// Needs to be implemented by a CSMS for handling messages part of the OCPP 2.0 Remote control profile.
type CSMSHandler interface {
}

// Needs to be implemented by Charging stations for handling messages part of the OCPP 2.0 Remote control profile.
type ChargingStationHandler interface {
}

const ProfileName = "remoteControl"

var Profile = ocpp.NewProfile(
	ProfileName)
