// The reservation functional block contains OCPP 2.0 features that enable EV drivers to make and manage reservations of charging stations.
package reservation

import (
	"github.com/lorenzodonini/ocpp-go/ocpp"
)

// Needs to be implemented by a CSMS for handling messages part of the OCPP 2.0 Reservation profile.
type CSMSHandler interface {
}

// Needs to be implemented by Charging stations for handling messages part of the OCPP 2.0 Reservation profile.
type ChargingStationHandler interface {
	// OnCancelReservation is called on a charging station whenever a CancelReservationRequest is received from the CSMS.
	OnCancelReservation(request *CancelReservationRequest) (confirmation *CancelReservationResponse, err error)
}

const ProfileName = "reservation"

var Profile = ocpp.NewProfile(
	ProfileName,
	CancelReservationFeature{},
	)
