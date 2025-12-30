package main

import (
	"fmt"
	"github.com/xBlaz3kx/ocpp-go/ocpp2.0.1/firmware"
)

func (c *CSMSHandler) OnFirmwareStatusNotification(chargingStationID string, request *firmware.FirmwareStatusNotificationRequest) (response *firmware.FirmwareStatusNotificationResponse, err error) {
	info, ok := c.chargingStations[chargingStationID]
	if !ok {
		err = fmt.Errorf("unknown charging station %v", chargingStationID)
		return
	}
	info.firmwareStatus = request.Status
	logDefault(chargingStationID, request.GetFeatureName()).Infof("updated firmware status to %v", request.Status)
	response = firmware.NewFirmwareStatusNotificationResponse()
	return
}

func (c *CSMSHandler) OnPublishFirmwareStatusNotification(chargingStationID string, request *firmware.PublishFirmwareStatusNotificationRequest) (response *firmware.PublishFirmwareStatusNotificationResponse, err error) {
	if len(request.Location) > 0 {
		logDefault(chargingStationID, request.GetFeatureName()).Infof("firmware download status on local controller: %v, download locations: %v", request.Status, request.Location)
	} else {
		logDefault(chargingStationID, request.GetFeatureName()).Infof("firmware download status on local controller: %v", request.Status)
	}
	response = firmware.NewPublishFirmwareStatusNotificationResponse()
	return
}
