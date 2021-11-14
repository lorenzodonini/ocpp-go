package main

import "github.com/lorenzodonini/ocpp-go/ocpp2.0.1/reservation"

func (c *CSMSHandler) OnReservationStatusUpdate(chargingStationID string, request *reservation.ReservationStatusUpdateRequest) (response *reservation.ReservationStatusUpdateResponse, err error) {
	logDefault(chargingStationID, request.GetFeatureName()).Infof("updated status of reservation %v to: %v", request.ReservationID, request.Status)
	response = reservation.NewReservationStatusUpdateResponse()
	return
}
