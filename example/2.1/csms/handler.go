package main

import (
	"github.com/sirupsen/logrus"

	"github.com/xBlaz3kx/ocpp-go/ocpp2.1/availability"
	"github.com/xBlaz3kx/ocpp-go/ocpp2.1/firmware"
	"github.com/xBlaz3kx/ocpp-go/ocpp2.1/types"
)

// TransactionInfo contains info about a transaction
type TransactionInfo struct {
	id          int
	startTime   *types.DateTime
	endTime     *types.DateTime
	startMeter  int
	endMeter    int
	connectorID int
	idTag       string
}

func (ti *TransactionInfo) hasTransactionEnded() bool {
	return ti.endTime != nil && !ti.endTime.IsZero()
}

// ConnectorInfo contains status and ongoing transaction ID for a connector
type ConnectorInfo struct {
	status             availability.ConnectorStatus
	currentTransaction int
}

func (ci *ConnectorInfo) hasTransactionInProgress() bool {
	return ci.currentTransaction >= 0
}

// ChargingStationState contains some simple state for a connected charging station
type ChargingStationState struct {
	status         availability.ChangeAvailabilityStatus
	firmwareStatus firmware.FirmwareStatus
	connectors     map[int]*ConnectorInfo
	transactions   map[int]*TransactionInfo
}

func (s *ChargingStationState) getConnector(id int) *ConnectorInfo {
	ci, ok := s.connectors[id]
	if !ok {
		ci = &ConnectorInfo{currentTransaction: -1}
		s.connectors[id] = ci
	}
	return ci
}

// CSMSHandler contains some simple state that a CSMS may want to keep.
// In production this will typically be replaced by database/API calls.
type CSMSHandler struct {
	chargingStations map[string]*ChargingStationState
}

// Utility functions
func logDefault(chargingStationID string, feature string) *logrus.Entry {
	return log.WithFields(logrus.Fields{"client": chargingStationID, "message": feature})
}
