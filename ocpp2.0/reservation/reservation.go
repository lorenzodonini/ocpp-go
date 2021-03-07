// The reservation functional block contains OCPP 2.0 features that enable EV drivers to make and manage reservations of charging stations.
package reservation

import (
	"github.com/lorenzodonini/ocpp-go/ocpp"
)

// Needs to be implemented by a CSMS for handling messages part of the OCPP 2.0 Reservation profile.
type CSMSHandler interface {
	// OnReservationStatusUpdate is called on the CSMS whenever a ReservationStatusUpdateRequest is received from a charging station.
	OnReservationStatusUpdate(chargingStationID string, request *ReservationStatusUpdateRequest) (resp *ReservationStatusUpdateResponse, err error)
}

// Needs to be implemented by Charging stations for handling messages part of the OCPP 2.0 Reservation profile.
type ChargingStationHandler interface {
	// OnCancelReservation is called on a charging station whenever a CancelReservationRequest is received from the CSMS.
	OnCancelReservation(request *CancelReservationRequest) (resp *CancelReservationResponse, err error)
}

const ProfileName = "reservation"

var Profile = ocpp.NewProfile(
	ProfileName,
	CancelReservationFeature{},
	ReservationStatusUpdateFeature{},
)
