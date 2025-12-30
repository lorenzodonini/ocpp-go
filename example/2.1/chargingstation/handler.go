package main

import (
	"fmt"
	"time"

	"github.com/xBlaz3kx/ocpp-go/ocpp2.0.1"
	"github.com/xBlaz3kx/ocpp-go/ocpp2.0.1/availability"
	"github.com/xBlaz3kx/ocpp-go/ocpp2.0.1/localauth"
	"github.com/xBlaz3kx/ocpp-go/ocpp2.0.1/reservation"
	"github.com/xBlaz3kx/ocpp-go/ocpp2.0.1/types"
)

// ConnectorInfo contains some simple state about a single connector.
type ConnectorInfo struct {
	status       availability.ConnectorStatus
	availability availability.OperationalStatus
	typ          reservation.ConnectorType
}

// EVSEInfo contains some simple state
type EVSEInfo struct {
	availability       availability.OperationalStatus
	currentTransaction string
	currentReservation int
	connectors         map[int]ConnectorInfo
	seqNo              int
}

// ChargingStationHandler contains some simple state that a charge point needs to keep.
// In production this will typically be replaced by database/API calls.
type ChargingStationHandler struct {
	availability         availability.OperationalStatus
	evse                 map[int]*EVSEInfo
	configuration        map[string]string
	model                string
	vendor               string
	meterValue           float64
	localAuthList        []localauth.AuthorizationData
	localAuthListVersion int
	monitoringLevel      int
}

var chargingStation ocpp2.ChargingStation

func (evse *EVSEInfo) hasConnector(ID int) bool {
	return ID > 0 && len(evse.connectors) > ID
}

func (evse *EVSEInfo) getConnectorOfType(typ reservation.ConnectorType) int {
	for i, c := range evse.connectors {
		if c.typ == typ {
			return i
		}
	}
	return -1
}

func (evse *EVSEInfo) nextSequence() int {
	seq := evse.seqNo
	evse.seqNo = seq + 1
	return seq
}

var transactionID = 0

func nextTransactionID() string {
	transactionID++
	return fmt.Sprintf("%d", transactionID)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getExpiryDate(info *types.IdTokenInfo) string {
	if info.CacheExpiryDateTime != nil {
		return fmt.Sprintf("authorized until %v", info.CacheExpiryDateTime.String())
	}
	return ""
}

func updateOperationalStatus(stateHandler *ChargingStationHandler, evseID int, status availability.OperationalStatus) {
	if evseID == 0 {
		stateHandler.availability = status
		log.Infof("operational status for charging station updated to: %v", status)
	} else if evse, ok := stateHandler.evse[evseID]; !ok {
		log.Errorf("couldn't update operational status for invalid evse %d", evseID)
		return
	} else {
		evse.availability = status
		log.Infof("operational status for evse %d updated to: %v", evseID, status)
	}
}

func updateConnectorStatus(stateHandler *ChargingStationHandler, evseID int, connector int, status availability.ConnectorStatus) {
	if evse, ok := stateHandler.evse[evseID]; !ok {
		log.Errorf("couldn't update connector status for invalid evse %d", evseID)
		return
	} else if connector < 0 || connector > len(evse.connectors) {
		log.Errorf("couldn't update status for evse %d with invalid connector %d", evseID, connector)
		return
	} else {
		conn := evse.connectors[connector]
		conn.status = status
		evse.connectors[connector] = conn
		// Send asynchronous status update
		response, err := chargingStation.StatusNotification(types.NewDateTime(time.Now()), status, evseID, connector)
		checkError(err)
		logDefault(response.GetFeatureName()).Infof("status for evse %d - connector %d updated to: %v", evseID, connector, status)
	}
}
