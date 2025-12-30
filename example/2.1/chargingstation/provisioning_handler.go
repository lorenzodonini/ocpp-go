package main

import (
	"github.com/xBlaz3kx/ocpp-go/ocpp"
	"github.com/xBlaz3kx/ocpp-go/ocpp2.0.1/provisioning"
	"github.com/xBlaz3kx/ocpp-go/ocppj"
)

func (handler *ChargingStationHandler) OnGetBaseReport(request *provisioning.GetBaseReportRequest) (response *provisioning.GetBaseReportResponse, err error) {
	logDefault(request.GetFeatureName()).Warnf("Unsupported feature")
	return nil, ocpp.NewHandlerError(ocppj.NotSupported, "Not supported")
}

func (handler *ChargingStationHandler) OnGetReport(request *provisioning.GetReportRequest) (response *provisioning.GetReportResponse, err error) {
	logDefault(request.GetFeatureName()).Warnf("Unsupported feature")
	return nil, ocpp.NewHandlerError(ocppj.NotSupported, "Not supported")
}

func (handler *ChargingStationHandler) OnGetVariables(request *provisioning.GetVariablesRequest) (response *provisioning.GetVariablesResponse, err error) {
	logDefault(request.GetFeatureName()).Warnf("Unsupported feature")
	return nil, ocpp.NewHandlerError(ocppj.NotSupported, "Not supported")
}

func (handler *ChargingStationHandler) OnReset(request *provisioning.ResetRequest) (response *provisioning.ResetResponse, err error) {
	logDefault(request.GetFeatureName()).Info("reset handled")
	response = provisioning.NewResetResponse(provisioning.ResetStatusAccepted)
	return
}

func (handler *ChargingStationHandler) OnSetNetworkProfile(request *provisioning.SetNetworkProfileRequest) (response *provisioning.SetNetworkProfileResponse, err error) {
	logDefault(request.GetFeatureName()).Warnf("Unsupported feature")
	return nil, ocpp.NewHandlerError(ocppj.NotSupported, "Not supported")
}

func (handler *ChargingStationHandler) OnSetVariables(request *provisioning.SetVariablesRequest) (response *provisioning.SetVariablesResponse, err error) {
	logDefault(request.GetFeatureName()).Warnf("Unsupported feature")
	return nil, ocpp.NewHandlerError(ocppj.NotSupported, "Not supported")
}

//func (handler *ChargingStationHandler) OnChangeConfiguration(request *core.ChangeConfigurationRequest) (confirmation *core.ChangeConfigurationConfirmation, err error) {
//	configKey, ok := handler.configuration[request.Key]
//	if !ok {
//		logDefault(request.GetFeatureName()).Errorf("couldn't change configuration for unsupported parameter %v", configKey.Key)
//		return core.NewChangeConfigurationConfirmation(core.ConfigurationStatusNotSupported), nil
//	} else if configKey.Readonly {
//		logDefault(request.GetFeatureName()).Errorf("couldn't change configuration for readonly parameter %v", configKey.Key)
//		return core.NewChangeConfigurationConfirmation(core.ConfigurationStatusRejected), nil
//	}
//	configKey.Value = &request.Value
//	handler.configuration[request.Key] = configKey
//	logDefault(request.GetFeatureName()).Infof("changed configuration for parameter %v to %v", configKey.Key, configKey.Value)
//	return core.NewChangeConfigurationConfirmation(core.ConfigurationStatusAccepted), nil
//}
//
//func (handler *ChargingStationHandler) OnGetConfiguration(request *core.GetConfigurationRequest) (confirmation *core.GetConfigurationConfirmation, err error) {
//	var resultKeys []core.ConfigurationKey
//	var unknownKeys []string
//	for _, key := range request.Key {
//		configKey, ok := handler.configuration[key]
//		if !ok {
//			unknownKeys = append(unknownKeys, *configKey.Value)
//		} else {
//			resultKeys = append(resultKeys, configKey)
//		}
//	}
//	if len(request.Key) == 0 {
//		// Return config for all keys∂
//		for _, v := range handler.configuration {
//			resultKeys = append(resultKeys, v)
//		}
//	}
//	logDefault(request.GetFeatureName()).Infof("returning configuration for requested keys: %v", request.Key)
//	conf := core.NewGetConfigurationConfirmation(resultKeys)
//	conf.UnknownKey = unknownKeys
//	return conf, nil
//}
