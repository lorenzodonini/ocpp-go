// The diagnostics functional block contains OCPP 2.0 features than enable remote diagnostics of problems with a charging station.
package diagnostics

import "github.com/lorenzodonini/ocpp-go/ocpp"

// Needs to be implemented by a CSMS for handling messages part of the OCPP 2.0 Diagnostics profile.
type CSMSHandler interface {
}

// Needs to be implemented by Charging stations for handling messages part of the OCPP 2.0 Diagnostics profile.
type ChargingStationHandler interface {
}

const ProfileName = "diagnostics"

var Profile = ocpp.NewProfile(
	ProfileName)
