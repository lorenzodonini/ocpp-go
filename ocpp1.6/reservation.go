package ocpp16

import "github.com/lorenzodonini/ocpp-go/ocpp"

const (
	CancelReservationFeatureName = "CancelReservation"
	ReserveNowFeatureName        = "ReserveNow"
)

type CentralSystemReservationListener interface {
}

type ChargePointReservationListener interface {
	OnReserveNow(request *ReserveNowRequest) (confirmation *ReserveNowConfirmation, err error)
	OnCancelReservation(request *CancelReservationRequest) (confirmation *CancelReservationConfirmation, err error)
}

const ReservationProfileName = "reservation"

var ReservationProfile = ocpp.NewProfile(
	ReservationProfileName,
	ReserveNowFeature{},
	CancelReservationFeature{})
