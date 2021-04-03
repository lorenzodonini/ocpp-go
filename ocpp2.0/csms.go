package ocpp2

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/internal/callbackqueue"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/authorization"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/availability"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/data"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/diagnostics"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/display"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/firmware"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/iso15118"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/meter"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/provisioning"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/remotecontrol"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/reservation"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/security"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/smartcharging"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/tariffcost"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/transactions"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
	"github.com/lorenzodonini/ocpp-go/ocppj"
	"github.com/lorenzodonini/ocpp-go/ws"
)

type csms struct {
	server               *ocppj.Server
	securityHandler      security.CSMSHandler
	provisioningHandler  provisioning.CSMSHandler
	authorizationHandler authorization.CSMSHandler
	localAuthListHandler localauth.CSMSHandler
	transactionsHandler  transactions.CSMSHandler
	remoteControlHandler remotecontrol.CSMSHandler
	availabilityHandler  availability.CSMSHandler
	reservationHandler   reservation.CSMSHandler
	tariffCostHandler    tariffcost.CSMSHandler
	meterHandler         meter.CSMSHandler
	smartChargingHandler smartcharging.CSMSHandler
	firmwareHandler      firmware.CSMSHandler
	iso15118Handler      iso15118.CSMSHandler
	diagnosticsHandler   diagnostics.CSMSHandler
	displayHandler       display.CSMSHandler
	dataHandler          data.CSMSHandler
	callbackQueue        callbackqueue.CallbackQueue
	errC                 chan error
}

func newCSMS(server *ocppj.Server) csms {
	if server == nil {
		panic("server must not be nil")
	}
	return csms{
		server:        server,
		callbackQueue: callbackqueue.New(),
	}
}

func (cs *csms) error(err error) {
	if cs.errC != nil {
		cs.errC <- err
	}
}

func (cs *csms) Errors() <-chan error {
	if cs.errC == nil {
		cs.errC = make(chan error, 1)
	}
	return cs.errC
}

func (cs *csms) CancelReservation(clientId string, callback func(*reservation.CancelReservationResponse, error), reservationId int, props ...func(request *reservation.CancelReservationRequest)) error {
	request := reservation.NewCancelReservationRequest(reservationId)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(response ocpp.Response, protoError error) {
		if response != nil {
			callback(response.(*reservation.CancelReservationResponse), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) CertificateSigned(clientId string, callback func(*security.CertificateSignedResponse, error), certificate []string, props ...func(*security.CertificateSignedRequest)) error {
	request := security.NewCertificateSignedRequest(certificate)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(response ocpp.Response, protoError error) {
		if response != nil {
			callback(response.(*security.CertificateSignedResponse), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) ChangeAvailability(clientId string, callback func(*availability.ChangeAvailabilityResponse, error), evseID int, operationalStatus availability.OperationalStatus, props ...func(request *availability.ChangeAvailabilityRequest)) error {
	request := availability.NewChangeAvailabilityRequest(evseID, operationalStatus)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(response ocpp.Response, protoError error) {
		if response != nil {
			callback(response.(*availability.ChangeAvailabilityResponse), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) ClearCache(clientId string, callback func(*authorization.ClearCacheResponse, error), props ...func(*authorization.ClearCacheRequest)) error {
	request := authorization.NewClearCacheRequest()
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(response ocpp.Response, protoError error) {
		if response != nil {
			callback(response.(*authorization.ClearCacheResponse), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) ClearChargingProfile(clientId string, callback func(*smartcharging.ClearChargingProfileResponse, error), props ...func(request *smartcharging.ClearChargingProfileRequest)) error {
	request := smartcharging.NewClearChargingProfileRequest()
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(response ocpp.Response, protoError error) {
		if response != nil {
			callback(response.(*smartcharging.ClearChargingProfileResponse), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) ClearDisplay(clientId string, callback func(*display.ClearDisplayResponse, error), id int, props ...func(*display.ClearDisplayRequest)) error {
	request := display.NewClearDisplayRequest(id)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(response ocpp.Response, protoError error) {
		if response != nil {
			callback(response.(*display.ClearDisplayResponse), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) ClearVariableMonitoring(clientId string, callback func(*diagnostics.ClearVariableMonitoringResponse, error), id []int, props ...func(*diagnostics.ClearVariableMonitoringRequest)) error {
	request := diagnostics.NewClearVariableMonitoringRequest(id)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(response ocpp.Response, protoError error) {
		if response != nil {
			callback(response.(*diagnostics.ClearVariableMonitoringResponse), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) CostUpdated(clientId string, callback func(*tariffcost.CostUpdatedResponse, error), totalCost float64, transactionId string, props ...func(*tariffcost.CostUpdatedRequest)) error {
	request := tariffcost.NewCostUpdatedRequest(totalCost, transactionId)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(response ocpp.Response, protoError error) {
		if response != nil {
			callback(response.(*tariffcost.CostUpdatedResponse), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) CustomerInformation(clientId string, callback func(*diagnostics.CustomerInformationResponse, error), requestId int, report bool, clear bool, props ...func(*diagnostics.CustomerInformationRequest)) error {
	request := diagnostics.NewCustomerInformationRequest(requestId, report, clear)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(response ocpp.Response, protoError error) {
		if response != nil {
			callback(response.(*diagnostics.CustomerInformationResponse), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) DataTransfer(clientId string, callback func(*data.DataTransferResponse, error), vendorId string, props ...func(request *data.DataTransferRequest)) error {
	request := data.NewDataTransferRequest(vendorId)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(response ocpp.Response, protoError error) {
		if response != nil {
			callback(response.(*data.DataTransferResponse), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) DeleteCertificate(clientId string, callback func(*iso15118.DeleteCertificateResponse, error), data types.CertificateHashData, props ...func(*iso15118.DeleteCertificateRequest)) error {
	request := iso15118.NewDeleteCertificateRequest(data)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(response ocpp.Response, protoError error) {
		if response != nil {
			callback(response.(*iso15118.DeleteCertificateResponse), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) GetBaseReport(clientId string, callback func(*provisioning.GetBaseReportResponse, error), requestId int, reportBase provisioning.ReportBaseType, props ...func(*provisioning.GetBaseReportRequest)) error {
	request := provisioning.NewGetBaseReportRequest(requestId, reportBase)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(response ocpp.Response, protoError error) {
		if response != nil {
			callback(response.(*provisioning.GetBaseReportResponse), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) GetChargingProfiles(clientId string, callback func(*smartcharging.GetChargingProfilesResponse, error), chargingProfile smartcharging.ChargingProfileCriterion, props ...func(*smartcharging.GetChargingProfilesRequest)) error {
	request := smartcharging.NewGetChargingProfilesRequest(chargingProfile)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(response ocpp.Response, protoError error) {
		if response != nil {
			callback(response.(*smartcharging.GetChargingProfilesResponse), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) GetCompositeSchedule(clientId string, callback func(*smartcharging.GetCompositeScheduleResponse, error), duration int, evseId int, props ...func(*smartcharging.GetCompositeScheduleRequest)) error {
	request := smartcharging.NewGetCompositeScheduleRequest(duration, evseId)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(response ocpp.Response, protoError error) {
		if response != nil {
			callback(response.(*smartcharging.GetCompositeScheduleResponse), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) GetDisplayMessages(clientId string, callback func(*display.GetDisplayMessagesResponse, error), requestId int, props ...func(*display.GetDisplayMessagesRequest)) error {
	request := display.NewGetDisplayMessagesRequest(requestId)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(response ocpp.Response, protoError error) {
		if response != nil {
			callback(response.(*display.GetDisplayMessagesResponse), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) GetInstalledCertificateIds(clientId string, callback func(*iso15118.GetInstalledCertificateIdsResponse, error), typeOfCertificate types.CertificateUse, props ...func(*iso15118.GetInstalledCertificateIdsRequest)) error {
	request := iso15118.NewGetInstalledCertificateIdsRequest(typeOfCertificate)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(response ocpp.Response, protoError error) {
		if response != nil {
			callback(response.(*iso15118.GetInstalledCertificateIdsResponse), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) GetLocalListVersion(clientId string, callback func(*localauth.GetLocalListVersionResponse, error), props ...func(*localauth.GetLocalListVersionRequest)) error {
	request := localauth.NewGetLocalListVersionRequest()
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(response ocpp.Response, protoError error) {
		if response != nil {
			callback(response.(*localauth.GetLocalListVersionResponse), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) GetLog(clientId string, callback func(*diagnostics.GetLogResponse, error), logType diagnostics.LogType, requestID int, logParameters diagnostics.LogParameters, props ...func(*diagnostics.GetLogRequest)) error {
	request := diagnostics.NewGetLogRequest(logType, requestID, logParameters)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(response ocpp.Response, protoError error) {
		if response != nil {
			callback(response.(*diagnostics.GetLogResponse), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) GetMonitoringReport(clientId string, callback func(*diagnostics.GetMonitoringReportResponse, error), props ...func(*diagnostics.GetMonitoringReportRequest)) error {
	request := diagnostics.NewGetMonitoringReportRequest()
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(response ocpp.Response, protoError error) {
		if response != nil {
			callback(response.(*diagnostics.GetMonitoringReportResponse), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) GetReport(clientId string, callback func(*provisioning.GetReportResponse, error), props ...func(*provisioning.GetReportRequest)) error {
	request := provisioning.NewGetReportRequest()
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(response ocpp.Response, protoError error) {
		if response != nil {
			callback(response.(*provisioning.GetReportResponse), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) GetTransactionStatus(clientId string, callback func(*transactions.GetTransactionStatusResponse, error), props ...func(*transactions.GetTransactionStatusRequest)) error {
	request := transactions.NewGetTransactionStatusRequest()
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(response ocpp.Response, protoError error) {
		if response != nil {
			callback(response.(*transactions.GetTransactionStatusResponse), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) GetVariables(clientId string, callback func(*provisioning.GetVariablesResponse, error), variableData []provisioning.GetVariableData, props ...func(*provisioning.GetVariablesRequest)) error {
	request := provisioning.NewGetVariablesRequest(variableData)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(response ocpp.Response, protoError error) {
		if response != nil {
			callback(response.(*provisioning.GetVariablesResponse), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) InstallCertificate(clientId string, callback func(*iso15118.InstallCertificateResponse, error), certificateType types.CertificateUse, certificate string, props ...func(*iso15118.InstallCertificateRequest)) error {
	request := iso15118.NewInstallCertificateRequest(certificateType, certificate)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(response ocpp.Response, protoError error) {
		if response != nil {
			callback(response.(*iso15118.InstallCertificateResponse), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) PublishFirmware(clientId string, callback func(*firmware.PublishFirmwareResponse, error), location string, checksum string, requestID int, props ...func(request *firmware.PublishFirmwareRequest)) error {
	request := firmware.NewPublishFirmwareRequest(location, checksum, requestID)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(response ocpp.Response, protoError error) {
		if response != nil {
			callback(response.(*firmware.PublishFirmwareResponse), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) RequestStartTransaction(clientId string, callback func(*remotecontrol.RequestStartTransactionResponse, error), remoteStartID int, IdToken types.IdTokenType, props ...func(request *remotecontrol.RequestStartTransactionRequest)) error {
	request := remotecontrol.NewRequestStartTransactionRequest(remoteStartID, IdToken)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(response ocpp.Response, protoError error) {
		if response != nil {
			callback(response.(*remotecontrol.RequestStartTransactionResponse), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) RequestStopTransaction(clientId string, callback func(*remotecontrol.RequestStopTransactionResponse, error), transactionID string, props ...func(request *remotecontrol.RequestStopTransactionRequest)) error {
	request := remotecontrol.NewRequestStopTransactionRequest(transactionID)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(response ocpp.Response, protoError error) {
		if response != nil {
			callback(response.(*remotecontrol.RequestStopTransactionResponse), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) ReserveNow(clientId string, callback func(*reservation.ReserveNowResponse, error), id int, expiryDateTime *types.DateTime, idToken types.IdTokenType, props ...func(request *reservation.ReserveNowRequest)) error {
	request := reservation.NewReserveNowRequest(id, expiryDateTime, idToken)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(response ocpp.Response, protoError error) {
		if response != nil {
			callback(response.(*reservation.ReserveNowResponse), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) Reset(clientId string, callback func(*provisioning.ResetResponse, error), t provisioning.ResetType, props ...func(request *provisioning.ResetRequest)) error {
	request := provisioning.NewResetRequest(t)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(response ocpp.Response, protoError error) {
		if response != nil {
			callback(response.(*provisioning.ResetResponse), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) SendLocalList(clientId string, callback func(*localauth.SendLocalListResponse, error), version int, updateType localauth.UpdateType, props ...func(request *localauth.SendLocalListRequest)) error {
	request := localauth.NewSendLocalListRequest(version, updateType)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(response ocpp.Response, protoError error) {
		if response != nil {
			callback(response.(*localauth.SendLocalListResponse), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) SetChargingProfile(clientId string, callback func(*smartcharging.SetChargingProfileResponse, error), evseID int, chargingProfile *types.ChargingProfile, props ...func(request *smartcharging.SetChargingProfileRequest)) error {
	request := smartcharging.NewSetChargingProfileRequest(evseID, chargingProfile)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(response ocpp.Response, protoError error) {
		if response != nil {
			callback(response.(*smartcharging.SetChargingProfileResponse), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) SetDisplayMessage(clientId string, callback func(*display.SetDisplayMessageResponse, error), message display.MessageInfo, props ...func(request *display.SetDisplayMessageRequest)) error {
	request := display.NewSetDisplayMessageRequest(message)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(response ocpp.Response, protoError error) {
		if response != nil {
			callback(response.(*display.SetDisplayMessageResponse), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) SetMonitoringBase(clientId string, callback func(*diagnostics.SetMonitoringBaseResponse, error), monitoringBase diagnostics.MonitoringBase, props ...func(request *diagnostics.SetMonitoringBaseRequest)) error {
	request := diagnostics.NewSetMonitoringBaseRequest(monitoringBase)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(response ocpp.Response, protoError error) {
		if response != nil {
			callback(response.(*diagnostics.SetMonitoringBaseResponse), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) SetMonitoringLevel(clientId string, callback func(*diagnostics.SetMonitoringLevelResponse, error), severity int, props ...func(request *diagnostics.SetMonitoringLevelRequest)) error {
	request := diagnostics.NewSetMonitoringLevelRequest(severity)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(response ocpp.Response, protoError error) {
		if response != nil {
			callback(response.(*diagnostics.SetMonitoringLevelResponse), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) SetNetworkProfile(clientId string, callback func(*provisioning.SetNetworkProfileResponse, error), configurationSlot int, connectionData provisioning.NetworkConnectionProfile, props ...func(request *provisioning.SetNetworkProfileRequest)) error {
	request := provisioning.NewSetNetworkProfileRequest(configurationSlot, connectionData)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(response ocpp.Response, protoError error) {
		if response != nil {
			callback(response.(*provisioning.SetNetworkProfileResponse), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) SetVariableMonitoring(clientId string, callback func(*diagnostics.SetVariableMonitoringResponse, error), data []diagnostics.SetMonitoringData, props ...func(request *diagnostics.SetVariableMonitoringRequest)) error {
	request := diagnostics.NewSetVariableMonitoringRequest(data)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(response ocpp.Response, protoError error) {
		if response != nil {
			callback(response.(*diagnostics.SetVariableMonitoringResponse), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) SetSecurityHandler(handler security.CSMSHandler) {
	cs.securityHandler = handler
}

func (cs *csms) SetProvisioningHandler(handler provisioning.CSMSHandler) {
	cs.provisioningHandler = handler
}

func (cs *csms) SetAuthorizationHandler(handler authorization.CSMSHandler) {
	cs.authorizationHandler = handler
}

func (cs *csms) SetLocalAuthListHandler(handler localauth.CSMSHandler) {
	cs.localAuthListHandler = handler
}

func (cs *csms) SetTransactionsHandler(handler transactions.CSMSHandler) {
	cs.transactionsHandler = handler
}

func (cs *csms) SetRemoteControlHandler(handler remotecontrol.CSMSHandler) {
	cs.remoteControlHandler = handler
}

func (cs *csms) SetAvailabilityHandler(handler availability.CSMSHandler) {
	cs.availabilityHandler = handler
}

func (cs *csms) SetReservationHandler(handler reservation.CSMSHandler) {
	cs.reservationHandler = handler
}

func (cs *csms) SetTariffCostHandler(handler tariffcost.CSMSHandler) {
	cs.tariffCostHandler = handler
}

func (cs *csms) SetMeterHandler(handler meter.CSMSHandler) {
	cs.meterHandler = handler
}

func (cs *csms) SetSmartChargingHandler(handler smartcharging.CSMSHandler) {
	cs.smartChargingHandler = handler
}

func (cs *csms) SetFirmwareHandler(handler firmware.CSMSHandler) {
	cs.firmwareHandler = handler
}

func (cs *csms) SetISO15118Handler(handler iso15118.CSMSHandler) {
	cs.iso15118Handler = handler
}

func (cs *csms) SetDiagnosticsHandler(handler diagnostics.CSMSHandler) {
	cs.diagnosticsHandler = handler
}

func (cs *csms) SetDisplayHandler(handler display.CSMSHandler) {
	cs.displayHandler = handler
}

func (cs *csms) SetDataHandler(handler data.CSMSHandler) {
	cs.dataHandler = handler
}

func (cs *csms) SetNewChargingStationHandler(handler ChargingStationConnectionHandler) {
	cs.server.SetNewClientHandler(func(chargingStation ws.Channel) {
		handler(chargingStation)
	})
}

func (cs *csms) SetChargingStationDisconnectedHandler(handler ChargingStationConnectionHandler) {
	cs.server.SetDisconnectedClientHandler(func(chargingStation ws.Channel) {
		handler(chargingStation)
	})
}

func (cs *csms) SendRequestAsync(clientId string, request ocpp.Request, callback func(response ocpp.Response, err error)) error {
	featureName := request.GetFeatureName()
	if _, found := cs.server.GetProfileForFeature(featureName); !found {
		return fmt.Errorf("feature %v is unsupported on CSMS (missing profile), cannot send request", featureName)
	}
	switch featureName {
	case reservation.CancelReservationFeatureName,
		security.CertificateSignedFeatureName,
		availability.ChangeAvailabilityFeatureName,
		authorization.ClearCacheFeatureName,
		smartcharging.ClearChargingProfileFeatureName,
		display.ClearDisplayFeatureName,
		diagnostics.ClearVariableMonitoringFeatureName,
		tariffcost.CostUpdatedFeatureName,
		diagnostics.CustomerInformationFeatureName,
		data.DataTransferFeatureName,
		iso15118.DeleteCertificateFeatureName,
		provisioning.GetBaseReportFeatureName,
		smartcharging.GetChargingProfilesFeatureName,
		smartcharging.GetCompositeScheduleFeatureName,
		display.GetDisplayMessagesFeatureName,
		iso15118.GetInstalledCertificateIdsFeatureName,
		localauth.GetLocalListVersionFeatureName,
		diagnostics.GetLogFeatureName,
		diagnostics.GetMonitoringReportFeatureName,
		provisioning.GetReportFeatureName,
		transactions.GetTransactionStatusFeatureName,
		provisioning.GetVariablesFeatureName,
		iso15118.InstallCertificateFeatureName,
		firmware.PublishFirmwareFeatureName,
		remotecontrol.RequestStartTransactionFeatureName,
		remotecontrol.RequestStopTransactionFeatureName,
		reservation.ReserveNowFeatureName,
		provisioning.ResetFeatureName,
		localauth.SendLocalListFeatureName,
		smartcharging.SetChargingProfileFeatureName,
		display.SetDisplayMessageFeatureName,
		diagnostics.SetMonitoringBaseFeatureName,
		diagnostics.SetMonitoringLevelFeatureName,
		provisioning.SetNetworkProfileFeatureName,
		diagnostics.SetVariableMonitoringFeatureName:
		break
	default:
		return fmt.Errorf("unsupported action %v on CSMS, cannot send request", featureName)
	}

	send := func() error {
		return cs.server.SendRequest(clientId, request)
	}
	return cs.callbackQueue.TryQueue(clientId, send, callback)
}

func (cs *csms) Start(listenPort int, listenPath string) {
	cs.server.Start(listenPort, listenPath)
}

func (cs *csms) sendResponse(chargingStationID string, response ocpp.Response, err error, requestId string) {
	if err != nil {
		err := cs.server.SendError(chargingStationID, requestId, ocppj.ProtocolError, "Couldn't generate valid confirmation", nil)
		if err != nil {
			err = fmt.Errorf("replying cs %s to request %s with 'protocol error': %w", chargingStationID, requestId, err)
			cs.error(err)
		}
		return
	}
	if response == nil {
		err = fmt.Errorf("empty response to %s for request %s", chargingStationID, requestId)
		cs.error(err)
		return
	}
	// send response
	err = cs.server.SendResponse(chargingStationID, requestId, response)
	if err != nil {
		err = fmt.Errorf("replying cs %s to request %s: %w", chargingStationID, requestId, err)
		cs.error(err)
	}
}

func (cs *csms) notImplementedError(chargingStationID string, requestId string, action string) {
	err := cs.server.SendError(chargingStationID, requestId, ocppj.NotImplemented, fmt.Sprintf("no handler for action %v implemented", action), nil)
	if err != nil {
		err = fmt.Errorf("replying cs %s to request %s with 'not implemented': %w", chargingStationID, requestId, err)
		cs.error(err)
	}
}

func (cs *csms) notSupportedError(chargingStationID string, requestId string, action string) {
	err := cs.server.SendError(chargingStationID, requestId, ocppj.NotSupported, fmt.Sprintf("unsupported action %v on CSMS", action), nil)
	if err != nil {
		err = fmt.Errorf("replying cs %s to request %s with 'not supported': %w", chargingStationID, requestId, err)
		cs.error(err)
	}
}

func (cs *csms) handleIncomingRequest(chargingStation ChargingStationConnection, request ocpp.Request, requestId string, action string) {
	profile, found := cs.server.GetProfileForFeature(action)
	// Check whether action is supported and a listener for it exists
	if !found {
		cs.notImplementedError(chargingStation.ID(), requestId, action)
		return
	} else {
		supported := true
		switch profile.Name {
		case authorization.ProfileName:
			if cs.authorizationHandler == nil {
				supported = false
			}
		case availability.ProfileName:
			if cs.availabilityHandler == nil {
				supported = false
			}
		case data.ProfileName:
			if cs.dataHandler == nil {
				supported = false
			}
		case diagnostics.ProfileName:
			if cs.diagnosticsHandler == nil {
				supported = false
			}
		case display.ProfileName:
			if cs.displayHandler == nil {
				supported = false
			}
		case firmware.ProfileName:
			if cs.firmwareHandler == nil {
				supported = false
			}
		case iso15118.ProfileName:
			if cs.iso15118Handler == nil {
				supported = false
			}
		case localauth.ProfileName:
			if cs.localAuthListHandler == nil {
				supported = false
			}
		case meter.ProfileName:
			if cs.meterHandler == nil {
				supported = false
			}
		case provisioning.ProfileName:
			if cs.provisioningHandler == nil {
				supported = false
			}
		case remotecontrol.ProfileName:
			if cs.remoteControlHandler == nil {
				supported = false
			}
		case reservation.ProfileName:
			if cs.reservationHandler == nil {
				supported = false
			}
		case security.ProfileName:
			if cs.securityHandler == nil {
				supported = false
			}
		case smartcharging.ProfileName:
			if cs.smartChargingHandler == nil {
				supported = false
			}
		case tariffcost.ProfileName:
			if cs.tariffCostHandler == nil {
				supported = false
			}
		case transactions.ProfileName:
			if cs.transactionsHandler == nil {
				supported = false
			}
		}
		if !supported {
			cs.notSupportedError(chargingStation.ID(), requestId, action)
			return
		}
	}
	var response ocpp.Response = nil
	var err error = nil
	// Execute in separate goroutine, so the caller goroutine is available
	go func() {
		switch action {
		case provisioning.BootNotificationFeatureName:
			response, err = cs.provisioningHandler.OnBootNotification(chargingStation.ID(), request.(*provisioning.BootNotificationRequest))
		case authorization.AuthorizeFeatureName:
			response, err = cs.authorizationHandler.OnAuthorize(chargingStation.ID(), request.(*authorization.AuthorizeRequest))
		case smartcharging.ClearedChargingLimitFeatureName:
			response, err = cs.smartChargingHandler.OnClearedChargingLimit(chargingStation.ID(), request.(*smartcharging.ClearedChargingLimitRequest))
		case data.DataTransferFeatureName:
			response, err = cs.dataHandler.OnDataTransfer(chargingStation.ID(), request.(*data.DataTransferRequest))
		case firmware.FirmwareStatusNotificationFeatureName:
			response, err = cs.firmwareHandler.OnFirmwareStatusNotification(chargingStation.ID(), request.(*firmware.FirmwareStatusNotificationRequest))
		case iso15118.Get15118EVCertificateFeatureName:
			response, err = cs.iso15118Handler.OnGet15118EVCertificate(chargingStation.ID(), request.(*iso15118.Get15118EVCertificateRequest))
		case iso15118.GetCertificateStatusFeatureName:
			response, err = cs.iso15118Handler.OnGetCertificateStatus(chargingStation.ID(), request.(*iso15118.GetCertificateStatusRequest))
		case availability.HeartbeatFeatureName:
			response, err = cs.availabilityHandler.OnHeartbeat(chargingStation.ID(), request.(*availability.HeartbeatRequest))
		case diagnostics.LogStatusNotificationFeatureName:
			response, err = cs.diagnosticsHandler.OnLogStatusNotification(chargingStation.ID(), request.(*diagnostics.LogStatusNotificationRequest))
		case meter.MeterValuesFeatureName:
			response, err = cs.meterHandler.OnMeterValues(chargingStation.ID(), request.(*meter.MeterValuesRequest))
		case smartcharging.NotifyChargingLimitFeatureName:
			response, err = cs.smartChargingHandler.OnNotifyChargingLimit(chargingStation.ID(), request.(*smartcharging.NotifyChargingLimitRequest))
		case diagnostics.NotifyCustomerInformationFeatureName:
			response, err = cs.diagnosticsHandler.OnNotifyCustomerInformation(chargingStation.ID(), request.(*diagnostics.NotifyCustomerInformationRequest))
		case display.NotifyDisplayMessagesFeatureName:
			response, err = cs.displayHandler.OnNotifyDisplayMessages(chargingStation.ID(), request.(*display.NotifyDisplayMessagesRequest))
		case smartcharging.NotifyEVChargingNeedsFeatureName:
			response, err = cs.smartChargingHandler.OnNotifyEVChargingNeeds(chargingStation.ID(), request.(*smartcharging.NotifyEVChargingNeedsRequest))
		case smartcharging.NotifyEVChargingScheduleFeatureName:
			response, err = cs.smartChargingHandler.OnNotifyEVChargingSchedule(chargingStation.ID(), request.(*smartcharging.NotifyEVChargingScheduleRequest))
		case diagnostics.NotifyEventFeatureName:
			response, err = cs.diagnosticsHandler.OnNotifyEvent(chargingStation.ID(), request.(*diagnostics.NotifyEventRequest))
		case diagnostics.NotifyMonitoringReportFeatureName:
			response, err = cs.diagnosticsHandler.OnNotifyMonitoringReport(chargingStation.ID(), request.(*diagnostics.NotifyMonitoringReportRequest))
		case provisioning.NotifyReportFeatureName:
			response, err = cs.provisioningHandler.OnNotifyReport(chargingStation.ID(), request.(*provisioning.NotifyReportRequest))
		case firmware.PublishFirmwareStatusNotificationFeatureName:
			response, err = cs.firmwareHandler.OnPublishFirmwareStatusNotification(chargingStation.ID(), request.(*firmware.PublishFirmwareStatusNotificationRequest))
		case smartcharging.ReportChargingProfilesFeatureName:
			response, err = cs.smartChargingHandler.OnReportChargingProfiles(chargingStation.ID(), request.(*smartcharging.ReportChargingProfilesRequest))
		case reservation.ReservationStatusUpdateFeatureName:
			response, err = cs.reservationHandler.OnReservationStatusUpdate(chargingStation.ID(), request.(*reservation.ReservationStatusUpdateRequest))
		case security.SecurityEventNotificationFeatureName:
			response, err = cs.securityHandler.OnSecurityEventNotification(chargingStation.ID(), request.(*security.SecurityEventNotificationRequest))
		default:
			cs.notSupportedError(chargingStation.ID(), requestId, action)
			return
		}
		cs.sendResponse(chargingStation.ID(), response, err, requestId)
	}()
}

func (cs *csms) handleIncomingResponse(chargingStation ChargingStationConnection, response ocpp.Response, requestId string) {
	if callback, ok := cs.callbackQueue.Dequeue(chargingStation.ID()); ok {
		callback(response, nil)
	} else {
		err := fmt.Errorf("no handler available for call of type %v from client %s for request %s", response.GetFeatureName(), chargingStation.ID(), requestId)
		cs.error(err)
	}
}

func (cs *csms) handleIncomingError(chargingStation ChargingStationConnection, err *ocpp.Error, details interface{}) {
	if callback, ok := cs.callbackQueue.Dequeue(chargingStation.ID()); ok {
		callback(nil, err)
	} else {
		cs.error(fmt.Errorf("no handler available for call error %w from client %s", err, chargingStation.ID()))
	}
}
