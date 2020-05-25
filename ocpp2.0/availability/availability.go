// The availability functional block contains OCPP 2.0 features for notifying the CSMS of availability and status changes.
// A CSMS can also instruct a charging station to change its availability.
package availability

import "github.com/lorenzodonini/ocpp-go/ocpp"

// Needs to be implemented by a CSMS for handling messages part of the OCPP 2.0 Availability profile.
type CSMSHandler interface {
}

// Needs to be implemented by Charging stations for handling messages part of the OCPP 2.0 Availability profile.
type ChargingStationHandler interface {
	// OnChangeAvailability is called on a charging station whenever a ChangeAvailabilityRequest is received from the CSMS.
	OnChangeAvailability(request *ChangeAvailabilityRequest) (confirmation *ChangeAvailabilityResponse, err error)
}

const ProfileName = "availability"

var Profile = ocpp.NewProfile(
	ProfileName,
	ChangeAvailabilityFeature{},
)
