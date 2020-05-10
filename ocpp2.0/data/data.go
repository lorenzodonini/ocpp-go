// The data transfer functional block enables parties to add custom commands and extensions to OCPP 2.0.
package data

import "github.com/lorenzodonini/ocpp-go/ocpp"

// Needs to be implemented by a CSMS for handling messages part of the OCPP 2.0 Data transfer profile.
type CSMSHandler interface {
}

// Needs to be implemented by Charging stations for handling messages part of the OCPP 2.0 Data transfer profile.
type ChargingStationHandler interface {
}

const ProfileName = "data"

var Profile = ocpp.NewProfile(
	ProfileName)
