// Contains support for reservation of a Charge Point.
package reservation

import "github.com/lorenzodonini/ocpp-go/ocpp"

const (
	CancelReservationFeatureName = "CancelReservation"
	ReserveNowFeatureName        = "ReserveNow"
)

type CentralSystemReservationHandler interface {
}

type ChargePointReservationHandler interface {
	OnReserveNow(request *ReserveNowRequest) (confirmation *ReserveNowConfirmation, err error)
	OnCancelReservation(request *CancelReservationRequest) (confirmation *CancelReservationConfirmation, err error)
}

const ProfileName = "reservation"

var Profile = ocpp.NewProfile(
	ProfileName,
	ReserveNowFeature{},
	CancelReservationFeature{})
