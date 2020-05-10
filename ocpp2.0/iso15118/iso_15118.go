// The ISO 15118 functional block contains OCPP 2.0 features that allow:
//
// - communication between EV and an EVSE
//
// - support for certificate-based authentication and authorization at the charging station, i.e. plug and charge
package iso15118

import "github.com/lorenzodonini/ocpp-go/ocpp"

// Needs to be implemented by a CSMS for handling messages part of the OCPP 2.0 ISO 15118 profile.
type CSMSHandler interface {
}

// Needs to be implemented by Charging stations for handling messages part of the OCPP 2.0 ISO 15118 profile.
type ChargingStationHandler interface {
}

const ProfileName = "iso15118"

var Profile = ocpp.NewProfile(
	ProfileName)
