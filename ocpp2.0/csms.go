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
	log "github.com/sirupsen/logrus"
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

func (cs *csms) SetNewChargingStationHandler(handler func(chargingStationID string)) {
	cs.server.SetNewClientHandler(handler)
}

func (cs *csms) SetChargingStationDisconnectedHandler(handler func(chargingStationID string)) {
	cs.server.SetDisconnectedClientHandler(handler)
}

func (cs *csms) SendRequestAsync(clientId string, request ocpp.Request, callback func(response ocpp.Response, err error)) error {
	featureName := request.GetFeatureName()
	if _, found := cs.server.GetProfileForFeature(featureName); !found {
		return fmt.Errorf("feature %v is unsupported on CSMS (missing profile), cannot send request", featureName)
	}
	switch featureName {
	case reservation.CancelReservationFeatureName, security.CertificateSignedFeatureName, availability.ChangeAvailabilityFeatureName, authorization.ClearCacheFeatureName, smartcharging.ClearChargingProfileFeatureName, display.ClearDisplayFeatureName, diagnostics.ClearVariableMonitoringFeatureName, tariffcost.CostUpdatedFeatureName, diagnostics.CustomerInformationFeatureName, data.DataTransferFeatureName, iso15118.DeleteCertificateFeatureName, provisioning.GetBaseReportFeatureName, smartcharging.GetChargingProfilesFeatureName, smartcharging.GetCompositeScheduleFeatureName, display.GetDisplayMessagesFeatureName, iso15118.GetInstalledCertificateIdsFeatureName, localauth.GetLocalListVersionFeatureName, diagnostics.GetLogFeatureName, diagnostics.GetMonitoringReportFeatureName:
		break
	default:
		return fmt.Errorf("unsupported action %v on CSMS, cannot send request", featureName)
	}

	cs.callbackQueue.Queue(clientId, callback)

	if err := cs.server.SendRequest(clientId, request); err != nil {
		_, _ = cs.callbackQueue.Dequeue(clientId)
		return err
	}
	return nil
}

func (cs *csms) Start(listenPort int, listenPath string) {
	cs.server.Start(listenPort, listenPath)
}

func (cs *csms) sendResponse(chargingStationID string, response ocpp.Response, err error, requestId string) {
	if response != nil {
		err := cs.server.SendResponse(chargingStationID, requestId, response)
		if err != nil {
			//TODO: handle error somehow
			log.Print(err)
		}
	} else {
		err := cs.server.SendError(chargingStationID, requestId, ocppj.ProtocolError, "Couldn't generate valid response", nil)
		if err != nil {
			log.WithFields(log.Fields{
				"client":  chargingStationID,
				"request": requestId,
			}).Errorf("unknown error %v while replying to message with CallError", err)
		}
	}
}

func (cs *csms) notImplementedError(chargingStationID string, requestId string, action string) {
	log.Warnf("Cannot handle call %v from charging station %v. Sending CallError instead", requestId, chargingStationID)
	err := cs.server.SendError(chargingStationID, requestId, ocppj.NotImplemented, fmt.Sprintf("no handler for action %v implemented", action), nil)
	if err != nil {
		log.WithFields(log.Fields{
			"client":  chargingStationID,
			"request": requestId,
		}).Errorf("unknown error %v while replying to message with CallError", err)
	}
}

func (cs *csms) notSupportedError(chargingStationID string, requestId string, action string) {
	log.Warnf("Cannot handle call %v from charging station %v. Sending CallError instead", requestId, chargingStationID)
	err := cs.server.SendError(chargingStationID, requestId, ocppj.NotSupported, fmt.Sprintf("unsupported action %v on CSMS", action), nil)
	if err != nil {
		log.WithFields(log.Fields{
			"client":  chargingStationID,
			"request": requestId,
		}).Errorf("unknown error %v while replying to message with CallError", err)
	}
}

func (cs *csms) handleIncomingRequest(chargingStationID string, request ocpp.Request, requestId string, action string) {
	profile, found := cs.server.GetProfileForFeature(action)
	// Check whether action is supported and a listener for it exists
	if !found {
		cs.notImplementedError(chargingStationID, requestId, action)
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
			cs.notSupportedError(chargingStationID, requestId, action)
			return
		}
	}
	var response ocpp.Response = nil
	var err error = nil
	// Execute in separate goroutine, so the caller goroutine is available
	go func() {
		switch action {
		case provisioning.BootNotificationFeatureName:
			response, err = cs.provisioningHandler.OnBootNotification(chargingStationID, request.(*provisioning.BootNotificationRequest))
		case authorization.AuthorizeFeatureName:
			response, err = cs.authorizationHandler.OnAuthorize(chargingStationID, request.(*authorization.AuthorizeRequest))
		case smartcharging.ClearedChargingLimitFeatureName:
			response, err = cs.smartChargingHandler.OnClearedChargingLimit(chargingStationID, request.(*smartcharging.ClearedChargingLimitRequest))
		case data.DataTransferFeatureName:
			response, err = cs.dataHandler.OnDataTransfer(chargingStationID, request.(*data.DataTransferRequest))
		case firmware.FirmwareStatusNotificationFeatureName:
			response, err = cs.firmwareHandler.OnFirmwareStatusNotification(chargingStationID, request.(*firmware.FirmwareStatusNotificationRequest))
		case iso15118.Get15118EVCertificateFeatureName:
			response, err = cs.iso15118Handler.OnGet15118EVCertificate(chargingStationID, request.(*iso15118.Get15118EVCertificateRequest))
		case iso15118.GetCertificateStatusFeatureName:
			response, err = cs.iso15118Handler.OnGetCertificateStatus(chargingStationID, request.(*iso15118.GetCertificateStatusRequest))
		default:
			cs.notSupportedError(chargingStationID, requestId, action)
			return
		}
		cs.sendResponse(chargingStationID, response, err, requestId)
	}()
}

func (cs *csms) handleIncomingResponse(chargingStationID string, response ocpp.Response, requestId string) {
	if callback, ok := cs.callbackQueue.Dequeue(chargingStationID); ok {
		callback(response, nil)
	} else {
		log.WithFields(log.Fields{
			"client":  chargingStationID,
			"request": requestId,
		}).Errorf("no handler available for Call Result of type %v", response.GetFeatureName())
	}
}

func (cs *csms) handleIncomingError(chargingStationID string, err *ocpp.Error, details interface{}) {
	if callback, ok := cs.callbackQueue.Dequeue(chargingStationID); ok {
		callback(nil, err)
	} else {
		//TODO: print details
		log.WithFields(log.Fields{
			"client":  chargingStationID,
			"request": err.MessageId,
		}).Errorf("no handler available for Call Error %v", err.Code)
	}
}
