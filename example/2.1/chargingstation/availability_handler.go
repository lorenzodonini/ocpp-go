package main

import (
	"github.com/xBlaz3kx/ocpp-go/ocpp2.1/availability"
)

func (handler *ChargingStationHandler) OnChangeAvailability(request *availability.ChangeAvailabilityRequest) (response *availability.ChangeAvailabilityResponse, err error) {
	if request.Evse == nil {
		// Changing availability for the entire charging station
		handler.availability = request.OperationalStatus
		// TODO: recursively update the availability for all evse/connectors
		response = availability.NewChangeAvailabilityResponse(availability.ChangeAvailabilityStatusAccepted)
		return
	}
	reqEvse := request.Evse
	if e, ok := handler.evse[reqEvse.ID]; ok {
		// Changing availability for a specific EVSE
		if reqEvse.ConnectorID != nil {
			// Changing availability for a specific connector
			if !e.hasConnector(*reqEvse.ConnectorID) {
				response = availability.NewChangeAvailabilityResponse(availability.ChangeAvailabilityStatusRejected)
			} else {
				connector := e.connectors[*reqEvse.ConnectorID]
				connector.availability = request.OperationalStatus
				e.connectors[*reqEvse.ConnectorID] = connector
				response = availability.NewChangeAvailabilityResponse(availability.ChangeAvailabilityStatusAccepted)
			}
			return
		}
		e.availability = request.OperationalStatus
		response = availability.NewChangeAvailabilityResponse(availability.ChangeAvailabilityStatusAccepted)
		return
	}
	// No EVSE with such ID found
	response = availability.NewChangeAvailabilityResponse(availability.ChangeAvailabilityStatusRejected)
	return
}
