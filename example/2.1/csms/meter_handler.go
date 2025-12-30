package main

import "github.com/xBlaz3kx/ocpp-go/ocpp2.0.1/meter"

func (c *CSMSHandler) OnMeterValues(chargingStationID string, request *meter.MeterValuesRequest) (response *meter.MeterValuesResponse, err error) {
	logDefault(chargingStationID, request.GetFeatureName()).Infof("received meter values for EVSE %v. Meter values:\n", request.EvseID)
	for _, mv := range request.MeterValue {
		logDefault(chargingStationID, request.GetFeatureName()).Printf("%v", mv)
	}
	response = meter.NewMeterValuesResponse()
	return
}
