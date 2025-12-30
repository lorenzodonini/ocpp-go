package main

import (
	"github.com/xBlaz3kx/ocpp-go/ocpp"
	"github.com/xBlaz3kx/ocpp-go/ocpp2.0.1/smartcharging"
	"github.com/xBlaz3kx/ocpp-go/ocppj"
)

func (handler *ChargingStationHandler) OnClearChargingProfile(request *smartcharging.ClearChargingProfileRequest) (response *smartcharging.ClearChargingProfileResponse, err error) {
	logDefault(request.GetFeatureName()).Warnf("Unsupported feature")
	return nil, ocpp.NewHandlerError(ocppj.NotSupported, "Not supported")
}

func (handler *ChargingStationHandler) OnGetChargingProfiles(request *smartcharging.GetChargingProfilesRequest) (response *smartcharging.GetChargingProfilesResponse, err error) {
	logDefault(request.GetFeatureName()).Warnf("Unsupported feature")
	return nil, ocpp.NewHandlerError(ocppj.NotSupported, "Not supported")
}

func (handler *ChargingStationHandler) OnGetCompositeSchedule(request *smartcharging.GetCompositeScheduleRequest) (response *smartcharging.GetCompositeScheduleResponse, err error) {
	logDefault(request.GetFeatureName()).Warnf("Unsupported feature")
	return nil, ocpp.NewHandlerError(ocppj.NotSupported, "Not supported")
}

func (handler *ChargingStationHandler) OnSetChargingProfile(request *smartcharging.SetChargingProfileRequest) (response *smartcharging.SetChargingProfileResponse, err error) {
	logDefault(request.GetFeatureName()).Warnf("Unsupported feature")
	return nil, ocpp.NewHandlerError(ocppj.NotSupported, "Not supported")
}
