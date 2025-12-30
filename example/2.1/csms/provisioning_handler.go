package main

import (
	"github.com/xBlaz3kx/ocpp-go/ocpp2.0.1/provisioning"
	"github.com/xBlaz3kx/ocpp-go/ocpp2.0.1/types"
	"time"
)

func (c *CSMSHandler) OnBootNotification(chargingStationID string, request *provisioning.BootNotificationRequest) (response *provisioning.BootNotificationResponse, err error) {
	logDefault(chargingStationID, request.GetFeatureName()).Infof("boot confirmed for %v %v, serial: %v, firmare version: %v, reason: %v",
		request.ChargingStation.VendorName, request.ChargingStation.Model, request.ChargingStation.SerialNumber, request.ChargingStation.FirmwareVersion, request.Reason)
	response = provisioning.NewBootNotificationResponse(types.NewDateTime(time.Now()), defaultHeartbeatInterval, provisioning.RegistrationStatusAccepted)
	return
}

func (c *CSMSHandler) OnNotifyReport(chargingStationID string, request *provisioning.NotifyReportRequest) (response *provisioning.NotifyReportResponse, err error) {
	logDefault(chargingStationID, request.GetFeatureName()).Infof("data report %v, seq. %v:\n", request.RequestID, request.SeqNo)
	for _, d := range request.ReportData {
		logDefault(chargingStationID, request.GetFeatureName()).Printf("%v", d)
	}
	response = provisioning.NewNotifyReportResponse()
	return
}
