package main

import (
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/smartcharging"
	"github.com/lorenzodonini/ocpp-go/ocppj"
)

func (c *CSMSHandler) OnClearedChargingLimit(chargingStationID string, request *smartcharging.ClearedChargingLimitRequest) (response *smartcharging.ClearedChargingLimitResponse, err error) {
	logDefault(chargingStationID, request.GetFeatureName()).Warnf("Unsupported feature")
	return nil, ocpp.NewHandlerError(ocppj.NotSupported, "Not supported")
}

func (c *CSMSHandler) OnNotifyChargingLimit(chargingStationID string, request *smartcharging.NotifyChargingLimitRequest) (response *smartcharging.NotifyChargingLimitResponse, err error) {
	logDefault(chargingStationID, request.GetFeatureName()).Warnf("Unsupported feature")
	return nil, ocpp.NewHandlerError(ocppj.NotSupported, "Not supported")
}

func (c *CSMSHandler) OnNotifyEVChargingNeeds(chargingStationID string, request *smartcharging.NotifyEVChargingNeedsRequest) (response *smartcharging.NotifyEVChargingNeedsResponse, err error) {
	logDefault(chargingStationID, request.GetFeatureName()).Warnf("Unsupported feature")
	return nil, ocpp.NewHandlerError(ocppj.NotSupported, "Not supported")
}

func (c *CSMSHandler) OnNotifyEVChargingSchedule(chargingStationID string, request *smartcharging.NotifyEVChargingScheduleRequest) (response *smartcharging.NotifyEVChargingScheduleResponse, err error) {
	logDefault(chargingStationID, request.GetFeatureName()).Warnf("Unsupported feature")
	return nil, ocpp.NewHandlerError(ocppj.NotSupported, "Not supported")
}

func (c *CSMSHandler) OnReportChargingProfiles(chargingStationID string, request *smartcharging.ReportChargingProfilesRequest) (response *smartcharging.ReportChargingProfilesResponse, err error) {
	logDefault(chargingStationID, request.GetFeatureName()).Warnf("Unsupported feature")
	return nil, ocpp.NewHandlerError(ocppj.NotSupported, "Not supported")
}
