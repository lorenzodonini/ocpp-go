// The Smart charging functional block contains OCPP 2.0 features that enable the CSO (or a third party) to influence the charging current/power transferred during a transaction, or set limits to the amount of current/power a Charging Station can draw from the grid.
package smartcharging

import "github.com/lorenzodonini/ocpp-go/ocpp"

// Needs to be implemented by a CSMS for handling messages part of the OCPP 2.0 Smart charging profile.
type CSMSHandler interface {
}

// Needs to be implemented by Charging stations for handling messages part of the OCPP 2.0 Smart charging profile.
type ChargingStationHandler interface {
}

const ProfileName = "smartCharging"

var Profile = ocpp.NewProfile(
	ProfileName)
