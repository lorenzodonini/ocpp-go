// The Meter values functional block contains OCPP 2.0 features for sending meter values to the CSMS.
package meter

import "github.com/lorenzodonini/ocpp-go/ocpp"

// Needs to be implemented by a CSMS for handling messages part of the OCPP 2.0 Meter values profile.
type CSMSHandler interface {
}

// Needs to be implemented by Charging stations for handling messages part of the OCPP 2.0 Meter values profile.
type ChargingStationHandler interface {
}

const ProfileName = "meter"

var Profile = ocpp.NewProfile(
	ProfileName)
