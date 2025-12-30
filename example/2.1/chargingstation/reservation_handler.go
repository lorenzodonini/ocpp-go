package main

import (
	"fmt"
	"github.com/xBlaz3kx/ocpp-go/ocpp2.1/availability"
	"github.com/xBlaz3kx/ocpp-go/ocpp2.1/reservation"
	"github.com/xBlaz3kx/ocpp-go/ocpp2.1/types"
)

func (handler *ChargingStationHandler) OnCancelReservation(request *reservation.CancelReservationRequest) (resp *reservation.CancelReservationResponse, err error) {
	for i, e := range handler.evse {
		if e.currentReservation == request.ReservationID {
			// Found reservation -> cancel
			e.currentReservation = -1
			for j := range e.connectors {
				if e.connectors[j].status == availability.ConnectorStatusReserved {
					go updateConnectorStatus(handler, i, j, availability.ConnectorStatusAvailable)
					break
				}
			}
			logDefault(request.GetFeatureName()).Infof("reservation %v for evse %d canceled", request.ReservationID, i)
			resp = reservation.NewCancelReservationResponse(reservation.CancelReservationStatusAccepted)
			return
		}
	}
	// Didn't find reservation -> reject
	logDefault(request.GetFeatureName()).Infof("couldn't cancel reservation %v: reservation not found!", request.ReservationID)
	resp = reservation.NewCancelReservationResponse(reservation.CancelReservationStatusRejected)
	return
}

func (handler *ChargingStationHandler) OnReserveNow(request *reservation.ReserveNowRequest) (resp *reservation.ReserveNowResponse, err error) {
	var reservedEvse int
	var reservedConnector int
	var status reservation.ReserveNowStatus

	status, reservedEvse, reservedConnector, err = handler.findConnector(request.EvseID, request.ConnectorType)
	if err != nil {
		logDefault(request.GetFeatureName()).Error(err)
	}
	resp = reservation.NewReserveNowResponse(status)
	if resp.Status != reservation.ReserveNowStatusAccepted {
		resp.StatusInfo = types.NewStatusInfo("code", err.Error())
		return
	}
	// Complete reservation
	evse := handler.evse[reservedEvse]
	evse.currentReservation = request.ID
	logDefault(request.GetFeatureName()).Infof("reservation %v accepted for evse %v, connector %v",
		request.ID, reservedEvse, reservedConnector)
	go updateConnectorStatus(handler, reservedEvse, reservedConnector, availability.ConnectorStatusReserved)

	// TODO: the logic above is incomplete. Advanced support for reservation management is missing.
	// TODO: automatically remove reservation after expiryDate
	return
}

func (handler *ChargingStationHandler) findConnector(requestedEVSE *int, connectorType reservation.ConnectorType) (status reservation.ReserveNowStatus, evseID int, connectorID int, err error) {
	status = reservation.ReserveNowStatusAccepted
	if requestedEVSE != nil {
		evseID = *requestedEVSE
		evse, ok := handler.evse[evseID]
		if !ok {
			status = reservation.ReserveNowStatusRejected
			err = fmt.Errorf("couldn't reserve a connector for invalid evse %d", evseID)
			return
		} else if evse.currentReservation != 0 {
			status = reservation.ReserveNowStatusOccupied
			err = fmt.Errorf("evse %v already has a pending reservation", evseID)
			return
		} else if evse.availability == availability.OperationalStatusInoperative {
			status = reservation.ReserveNowStatusUnavailable
			err = fmt.Errorf("evse %v is currently not operative", evseID)
			return
		}
		for i, c := range evse.connectors {
			if connectorType != "" && c.typ != connectorType {
				continue
			}
			switch c.status {
			case availability.ConnectorStatusReserved, availability.ConnectorStatusOccupied:
				status = reservation.ReserveNowStatusOccupied
			case availability.ConnectorStatusUnavailable:
				status = reservation.ReserveNowStatusUnavailable
			case availability.ConnectorStatusFaulted:
				status = reservation.ReserveNowStatusUnavailable
			case availability.ConnectorStatusAvailable:
				// Found an available connector
				status = reservation.ReserveNowStatusAccepted
				connectorID = i
				return
			}
		}
	} else {
		// Find suitable evse + connector
		for j, e := range handler.evse {
			evseID = j
			if e.currentReservation != 0 {
				status = reservation.ReserveNowStatusOccupied
				err = fmt.Errorf("evse %v already has a pending reservation", evseID)
				return
			} else if e.availability == availability.OperationalStatusInoperative {
				status = reservation.ReserveNowStatusUnavailable
				err = fmt.Errorf("evse %v is currently not operative", evseID)
				return
			}
			for i, c := range e.connectors {
				if connectorType != "" && c.typ != connectorType {
					continue
				}
				switch c.status {
				case availability.ConnectorStatusReserved, availability.ConnectorStatusOccupied:
					status = reservation.ReserveNowStatusOccupied
				case availability.ConnectorStatusUnavailable:
					status = reservation.ReserveNowStatusUnavailable
				case availability.ConnectorStatusFaulted:
					status = reservation.ReserveNowStatusUnavailable
				case availability.ConnectorStatusAvailable:
					// Found an available connector
					status = reservation.ReserveNowStatusAccepted
					connectorID = i
					return
				}
			}
		}
		if status == "" {
			status = reservation.ReserveNowStatusRejected
			err = fmt.Errorf("no available evse found")
		}
	}
	return
}
