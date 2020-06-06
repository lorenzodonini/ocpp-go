// Contains support for reservation of a Charge Point.
package reservation

import "github.com/lorenzodonini/ocpp-go/ocpp"

// Needs to be implemented by Central systems for handling messages part of the OCPP 1.6 Reservation profile.
type CentralSystemHandler interface {
}

// Needs to be implemented by Charge points for handling messages part of the OCPP 1.6 Reservation profile.
type ChargePointHandler interface {
	OnReserveNow(request *ReserveNowRequest) (confirmation *ReserveNowConfirmation, err error)
	OnCancelReservation(request *CancelReservationRequest) (confirmation *CancelReservationConfirmation, err error)
}

// The profile name
const ProfileName = "reservation"

// Provides support for for reservation of a Charge Point.
var Profile = ocpp.NewProfile(
	ProfileName,
	ReserveNowFeature{},
	CancelReservationFeature{})
