// The authorization functional block contains OCPP 2.0 features that show tariff and costs to an EV driver, when supported by the charging station.
package tariffcost

import "github.com/lorenzodonini/ocpp-go/ocpp"

// Needs to be implemented by a CSMS for handling messages part of the OCPP 2.0 Tariff and cost profile.
type CSMSHandler interface {
}

// Needs to be implemented by Charging stations for handling messages part of the OCPP 2.0 Tariff and cost profile.
type ChargingStationHandler interface {
}

const ProfileName = "tariffCost"

var Profile = ocpp.NewProfile(
	ProfileName)
