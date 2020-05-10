// The display functional block contains OCPP 2.0 features for managing message that get displayed on a charging station.
package display

import "github.com/lorenzodonini/ocpp-go/ocpp"

// Needs to be implemented by a CSMS for handling messages part of the OCPP 2.0 Display profile.
type CSMSHandler interface {
}

// Needs to be implemented by Charging stations for handling messages part of the OCPP 2.0 Display profile.
type ChargingStationHandler interface {
}

const ProfileName = "display"

var Profile = ocpp.NewProfile(
	ProfileName)
