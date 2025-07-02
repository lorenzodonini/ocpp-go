package main

import (
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/diagnostics"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

func (handler *ChargingStationHandler) OnClearVariableMonitoring(request *diagnostics.ClearVariableMonitoringRequest) (response *diagnostics.ClearVariableMonitoringResponse, err error) {
	logDefault(request.GetFeatureName()).Infof("cleared variables %v", request.ID)
	clearMonitoringResult := make([]diagnostics.ClearMonitoringResult, len(request.ID))
	for i, req := range request.ID {
		res := diagnostics.ClearMonitoringResult{
			ID:     req,
			Status: diagnostics.ClearMonitoringStatusAccepted,
		}
		clearMonitoringResult[i] = res
	}
	return diagnostics.NewClearVariableMonitoringResponse(clearMonitoringResult), nil
}

func (handler *ChargingStationHandler) OnCustomerInformation(request *diagnostics.CustomerInformationRequest) (response *diagnostics.CustomerInformationResponse, err error) {
	logDefault(request.GetFeatureName()).Infof("request %d for customer %s (clear %v, report %v)", request.RequestID, request.CustomerIdentifier, request.Clear, request.Report)
	return diagnostics.NewCustomerInformationResponse(diagnostics.CustomerInformationStatusAccepted), nil
}

func (handler *ChargingStationHandler) OnGetLog(request *diagnostics.GetLogRequest) (response *diagnostics.GetLogResponse, err error) {
	logDefault(request.GetFeatureName()).Infof("request %d to upload logs %v to %v accepted", request.RequestID, request.LogType, request.Log.RemoteLocation)
	// TODO: start asynchronous log upload
	return diagnostics.NewGetLogResponse(diagnostics.LogStatusAccepted), nil
}

func (handler *ChargingStationHandler) OnGetMonitoringReport(request *diagnostics.GetMonitoringReportRequest) (response *diagnostics.GetMonitoringReportResponse, err error) {
	logDefault(request.GetFeatureName()).Infof("request %d to upload report with criteria %v, component variables %v", request.RequestID, request.MonitoringCriteria, request.ComponentVariable)
	// TODO: start asynchronous report upload via NotifyMonitoringReportRequest
	return diagnostics.NewGetMonitoringReportResponse(types.GenericDeviceModelStatusAccepted), nil
}

func (handler *ChargingStationHandler) OnSetMonitoringBase(request *diagnostics.SetMonitoringBaseRequest) (response *diagnostics.SetMonitoringBaseResponse, err error) {
	logDefault(request.GetFeatureName()).Infof("monitoring base %s set successfully", request.MonitoringBase)
	return diagnostics.NewSetMonitoringBaseResponse(types.GenericDeviceModelStatusAccepted), nil
}

func (handler *ChargingStationHandler) OnSetMonitoringLevel(request *diagnostics.SetMonitoringLevelRequest) (response *diagnostics.SetMonitoringLevelResponse, err error) {
	handler.monitoringLevel = request.Severity
	logDefault(request.GetFeatureName()).Infof("set monitoring severity level to %d", handler.monitoringLevel)
	return diagnostics.NewSetMonitoringLevelResponse(types.GenericDeviceModelStatusAccepted), nil
}

func (handler *ChargingStationHandler) OnSetVariableMonitoring(request *diagnostics.SetVariableMonitoringRequest) (response *diagnostics.SetVariableMonitoringResponse, err error) {
	setMonitoringResult := make([]diagnostics.SetMonitoringResult, len(request.MonitoringData))
	for i, req := range request.MonitoringData {
		// TODO: configure custom monitoring rules internal for the received SetMonitoringData parameters
		logDefault(request.GetFeatureName()).Infof("set monitoring for component %v, variable %v to type %v = %v, severity %v",
			req.Component.Name, req.Variable.Name, req.Type, req.Value, req.Severity)
		res := diagnostics.SetMonitoringResult{
			ID:        req.ID,
			Status:    diagnostics.SetMonitoringStatusAccepted,
			Type:      req.Type,
			Severity:  req.Severity,
			Component: req.Component,
			Variable:  req.Variable,
		}
		setMonitoringResult[i] = res
	}
	return diagnostics.NewSetVariableMonitoringResponse(setMonitoringResult), nil
}
