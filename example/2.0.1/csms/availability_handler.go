package main

import (
	"fmt"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/availability"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

func (c *CSMSHandler) OnHeartbeat(chargingStationID string, request *availability.HeartbeatRequest) (response *availability.HeartbeatResponse, err error) {
	logDefault(chargingStationID, request.GetFeatureName()).Infof("heartbeat handled")
	response = availability.NewHeartbeatResponse(types.DateTime{Time: time.Now()})
	return
}

func (c *CSMSHandler) OnStatusNotification(chargingStationID string, request *availability.StatusNotificationRequest) (response *availability.StatusNotificationResponse, err error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	info, ok := c.chargingStations[chargingStationID]
	if !ok {
		return nil, fmt.Errorf("unknown charging station %v", chargingStationID)
	}

	if request.ConnectorID > 0 {
		connectorInfo := info.getConnector(request.ConnectorID)
		connectorInfo.status = request.ConnectorStatus
		logDefault(chargingStationID, request.GetFeatureName()).Infof("connector %v updated status to %v", request.ConnectorID, request.ConnectorStatus)
	} else {
		logDefault(chargingStationID, request.GetFeatureName()).Infof("couldn't update status for invalid connector %v", request.ConnectorID)
	}
	response = availability.NewStatusNotificationResponse()
	return
}
