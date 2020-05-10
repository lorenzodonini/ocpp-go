package ocpp2

import (
	"fmt"
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

type chargingStation struct {
	client               *ocppj.Client
	messageHandler       ChargingStationHandler
	securityHandler      security.ChargingStationHandler
	provisioningHandler  provisioning.ChargingStationHandler
	authorizationHandler authorization.ChargingStationHandler
	localAuthListHandler localauth.ChargingStationHandler
	transactionsHandler  transactions.ChargingStationHandler
	remoteControlHandler remotecontrol.ChargingStationHandler
	availabilityHandler  availability.ChargingStationHandler
	reservationHandler   reservation.ChargingStationHandler
	tariffCostHandler    tariffcost.ChargingStationHandler
	meterHandler         meter.ChargingStationHandler
	smartChargingHandler smartcharging.ChargingStationHandler
	firmwareHandler      firmware.ChargingStationHandler
	iso15118Handler      iso15118.ChargingStationHandler
	diagnosticsHandler   diagnostics.ChargingStationHandler
	displayHandler       display.ChargingStationHandler
	dataHandler          data.CSMSHandler
	confirmationHandler  chan ocpp.Response
	errorHandler         chan error
}

func (cs *chargingStation) BootNotification(reason provisioning.BootReason, model string, chargePointVendor string, props ...func(request *provisioning.BootNotificationRequest)) (*provisioning.BootNotificationConfirmation, error) {
	request := provisioning.NewBootNotificationRequest(reason, model, chargePointVendor)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cs.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*provisioning.BootNotificationConfirmation), err
	}
}

func (cs *chargingStation) Authorize(idToken string, tokenType types.IdTokenType, props ...func(request *authorization.AuthorizeRequest)) (*authorization.AuthorizeConfirmation, error) {
	request := authorization.NewAuthorizationRequest(idToken, tokenType)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cs.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*authorization.AuthorizeConfirmation), err
	}
}

func (cs *chargingStation) ClearedChargingLimit(chargingLimitSource types.ChargingLimitSourceType, props ...func(request *ClearedChargingLimitRequest)) (*ClearedChargingLimitConfirmation, error) {
	request := NewClearedChargingLimitRequest(chargingLimitSource)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cs.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*ClearedChargingLimitConfirmation), err
	}
}

// Starts a custom data transfer request. Every vendor may implement their own proprietary logic for this message.
func (cs *chargingStation) DataTransfer(vendorId string, props ...func(request *DataTransferRequest)) (*DataTransferConfirmation, error) {
	request := NewDataTransferRequest(vendorId)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cs.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*DataTransferConfirmation), err
	}
}

func (cs *chargingStation) FirmwareStatusNotification(status FirmwareStatus, requestID int, props ...func(request *FirmwareStatusNotificationRequest)) (*FirmwareStatusNotificationConfirmation, error) {
	request := NewFirmwareStatusNotificationRequest(status, requestID)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cs.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*FirmwareStatusNotificationConfirmation), err
	}
}

func (cs *chargingStation) Get15118EVCertificate(schemaVersion string, exiRequest string, props ...func(request *Get15118EVCertificateRequest)) (*Get15118EVCertificateConfirmation, error) {
	request := NewGet15118EVCertificateRequest(schemaVersion, exiRequest)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cs.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*Get15118EVCertificateConfirmation), err
	}
}

func (cs *chargingStation) GetCertificateStatus(ocspRequestData types.OCSPRequestDataType, props ...func(request *GetCertificateStatusRequest)) (*GetCertificateStatusConfirmation, error) {
	request := NewGetCertificateStatusRequest(ocspRequestData)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cs.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*GetCertificateStatusConfirmation), err
	}
}

//
//// Notifies the central system that the charge point is still online. The central system's response is used for time synchronization purposes. It is recommended to perform this operation once every 24 hours.
//func (cp *chargingStation) Heartbeat(props ...func(request *HeartbeatRequest)) (*HeartbeatConfirmation, error) {
//	request := NewHeartbeatRequest()
//	for _, fn := range props {
//		fn(request)
//	}
//	confirmation, err := cp.SendRequest(request)
//	if err != nil {
//		return nil, err
//	} else {
//		return confirmation.(*HeartbeatConfirmation), err
//	}
//}
//
//// Sends a batch of collected meter values to the central system, for billing and analysis. May be done periodically during ongoing transactions.
//func (cp *chargingStation) MeterValues(connectorId int, meterValues []MeterValue, props ...func(request *MeterValuesRequest)) (*MeterValuesConfirmation, error) {
//	request := NewMeterValuesRequest(connectorId, meterValues)
//	for _, fn := range props {
//		fn(request)
//	}
//	confirmation, err := cp.SendRequest(request)
//	if err != nil {
//		return nil, err
//	} else {
//		return confirmation.(*MeterValuesConfirmation), err
//	}
//}
//
//// Requests to start a transaction for a specific connector. The central system will verify the client's IdTag and either accept or reject the transaction.
//func (cp *chargingStation) StartTransaction(connectorId int, idTag string, meterStart int, timestamp *DateTime, props ...func(request *StartTransactionRequest)) (*StartTransactionConfirmation, error) {
//	request := NewStartTransactionRequest(connectorId, idTag, meterStart, timestamp)
//	for _, fn := range props {
//		fn(request)
//	}
//	confirmation, err := cp.SendRequest(request)
//	if err != nil {
//		return nil, err
//	} else {
//		return confirmation.(*StartTransactionConfirmation), err
//	}
//}
//
//// Stops an ongoing transaction. Typically a batch of meter values is passed along with this message.
//func (cp *chargingStation) StopTransaction(meterStop int, timestamp *DateTime, transactionId int, props ...func(request *StopTransactionRequest)) (*StopTransactionConfirmation, error) {
//	request := NewStopTransactionRequest(meterStop, timestamp, transactionId)
//	for _, fn := range props {
//		fn(request)
//	}
//	confirmation, err := cp.SendRequest(request)
//	if err != nil {
//		return nil, err
//	} else {
//		return confirmation.(*StopTransactionConfirmation), err
//	}
//}
//
//// Notifies the central system of a status update. This may apply to the entire charge point or to a single connector.
//func (cp *chargingStation) StatusNotification(connectorId int, errorCode ChargePointErrorCode, status ChargePointStatus, props ...func(request *StatusNotificationRequest)) (*StatusNotificationConfirmation, error) {
//	request := NewStatusNotificationRequest(connectorId, errorCode, status)
//	for _, fn := range props {
//		fn(request)
//	}
//	confirmation, err := cp.SendRequest(request)
//	if err != nil {
//		return nil, err
//	} else {
//		return confirmation.(*StatusNotificationConfirmation), err
//	}
//}
//
//// Notifies the central system of a status change in the upload of diagnostics data.
//func (cp *chargingStation) DiagnosticsStatusNotification(status DiagnosticsStatus, props ...func(request *DiagnosticsStatusNotificationRequest)) (*DiagnosticsStatusNotificationConfirmation, error) {
//	request := NewDiagnosticsStatusNotificationRequest(status)
//	for _, fn := range props {
//		fn(request)
//	}
//	confirmation, err := cp.SendRequest(request)
//	if err != nil {
//		return nil, err
//	} else {
//		return confirmation.(*DiagnosticsStatusNotificationConfirmation), err
//	}
//}
//
//// Notifies the central system of a status change during the download of a new firmware version.
//func (cp *chargingStation) FirmwareStatusNotification(status FirmwareStatus, props ...func(request *FirmwareStatusNotificationRequest)) (*FirmwareStatusNotificationConfirmation, error) {
//	request := NewFirmwareStatusNotificationRequest(status)
//	for _, fn := range props {
//		fn(request)
//	}
//	confirmation, err := cp.SendRequest(request)
//	if err != nil {
//		return nil, err
//	} else {
//		return confirmation.(*FirmwareStatusNotificationConfirmation), err
//	}
//}

func (cs *chargingStation) SetMessageHandler(handler ChargingStationHandler) {
	cs.messageHandler = handler
}

func (cs *chargingStation) SetSecurityHandler(handler security.ChargingStationHandler) {
	cs.securityHandler = handler
}

func (cs *chargingStation) SetProvisioningHandler(handler provisioning.ChargingStationHandler) {
	cs.provisioningHandler = handler
}

func (cs *chargingStation) SetAuthorizationHandler(handler authorization.ChargingStationHandler) {
	cs.authorizationHandler = handler
}

func (cs *chargingStation) SetLocalAuthListHandler(handler localauth.ChargingStationHandler) {
	cs.localAuthListHandler = handler
}

func (cs *chargingStation) SetTransactionsHandler(handler transactions.ChargingStationHandler) {
	cs.transactionsHandler = handler
}

func (cs *chargingStation) SetRemoteControlHandler(handler transactions.ChargingStationHandler) {
	cs.remoteControlHandler = handler
}

func (cs *chargingStation) SetAvailabilityHandler(handler transactions.ChargingStationHandler) {
	cs.availabilityHandler = handler
}

func (cs *chargingStation) SetReservationHandler(handler reservation.ChargingStationHandler) {
	cs.reservationHandler = handler
}

func (cs *chargingStation) SetTariffCostHandler(handler tariffcost.ChargingStationHandler) {
	cs.tariffCostHandler = handler
}

func (cs *chargingStation) SetMeterHandler(handler tariffcost.ChargingStationHandler) {
	cs.meterHandler = handler
}

func (cs *chargingStation) SetSmartChargingHandler(handler smartcharging.ChargingStationHandler) {
	cs.smartChargingHandler = handler
}

func (cs *chargingStation) SetFirmwareHandler(handler firmware.ChargingStationHandler) {
	cs.firmwareHandler = handler
}

func (cs *chargingStation) SetISO15118Handler(handler iso15118.ChargingStationHandler) {
	cs.iso15118Handler = handler
}

func (cs *chargingStation) SetDiagnosticsHandler(handler diagnostics.ChargingStationHandler) {
	cs.diagnosticsHandler = handler
}

func (cs *chargingStation) SetDisplayHandler(handler display.ChargingStationHandler) {
	cs.displayHandler = handler
}

func (cs *chargingStation) SetDataHandler(handler data.ChargingStationHandler) {
	cs.dataHandler = handler
}

//// Registers a handler for incoming local authorization profile messages
//func (cp *chargingStation) SetLocalAuthListHandler(listener ChargePointLocalAuthListListener) {
//	cp.localAuthListListener = listener
//}
//
//// Registers a handler for incoming firmware management profile messages
//func (cp *chargingStation) SetFirmwareManagementHandler(listener ChargePointFirmwareManagementListener) {
//	cp.firmwareListener = listener
//}
//
//// Registers a handler for incoming reservation profile messages
//func (cp *chargingStation) SetReservationHandler(listener ChargePointReservationListener) {
//	cp.reservationListener = listener
//}
//
//// Registers a handler for incoming remote trigger profile messages
//func (cp *chargingStation) SetRemoteTriggerHandler(listener ChargePointRemoteTriggerListener) {
//	cp.remoteTriggerListener = listener
//}
//
//// Registers a handler for incoming smart charging profile messages
//func (cp *chargingStation) SetSmartChargingHandler(listener ChargePointSmartChargingListener) {
//	cp.smartChargingListener = listener
//}

func (cs *chargingStation) SendRequest(request ocpp.Request) (ocpp.Response, error) {
	// TODO: check for supported feature
	err := cs.client.SendRequest(request)
	if err != nil {
		return nil, err
	}
	//TODO: timeouts
	select {
	case confirmation := <-cs.confirmationHandler:
		return confirmation, nil
	case err = <-cs.errorHandler:
		return nil, err
	}
}

func (cs *chargingStation) SendRequestAsync(request ocpp.Request, callback func(confirmation ocpp.Response, err error)) error {
	switch request.GetFeatureName() {
	case authorization.AuthorizeFeatureName, provisioning.BootNotificationFeatureName, ClearedChargingLimitFeatureName, DataTransferFeatureName, FirmwareStatusNotificationFeatureName, Get15118EVCertificateFeatureName, GetCertificateStatusFeatureName:
		break
	default:
		return fmt.Errorf("unsupported action %v on charge point, cannot send request", request.GetFeatureName())
	}
	err := cs.client.SendRequest(request)
	if err == nil {
		// Retrieve result asynchronously
		go func() {
			select {
			case confirmation := <-cs.confirmationHandler:
				callback(confirmation, nil)
			case protoError := <-cs.errorHandler:
				callback(nil, protoError)
			}
		}()
	}
	return err
}

func (cs *chargingStation) sendResponse(confirmation ocpp.Response, err error, requestId string) {
	if confirmation != nil {
		err := cs.client.SendResponse(requestId, confirmation)
		if err != nil {
			log.WithField("request", requestId).Errorf("unknown error %v while replying to message with CallError", err)
			//TODO: handle error somehow
		}
	} else {
		err = cs.client.SendError(requestId, ocppj.ProtocolError, err.Error(), nil)
		if err != nil {
			log.WithField("request", requestId).Errorf("unknown error %v while replying to message with CallError", err)
		}
	}
}

func (cs *chargingStation) Start(csmsUrl string) error {
	// TODO: implement auto-reconnect logic
	return cs.client.Start(csmsUrl)
}

func (cs *chargingStation) Stop() {
	cs.client.Stop()
}

func (cs *chargingStation) notImplementedError(requestId string, action string) {
	log.WithField("request", requestId).Errorf("cannot handle call from central system. Sending CallError instead")
	err := cs.client.SendError(requestId, ocppj.NotImplemented, fmt.Sprintf("no handler for action %v implemented", action), nil)
	if err != nil {
		log.WithField("request", requestId).Errorf("unknown error %v while replying to message with CallError", err)
	}
}

func (cs *chargingStation) notSupportedError(requestId string, action string) {
	log.WithField("request", requestId).Errorf("cannot handle call from central system. Sending CallError instead")
	err := cs.client.SendError(requestId, ocppj.NotSupported, fmt.Sprintf("unsupported action %v on charge point", action), nil)
	if err != nil {
		log.WithField("request", requestId).Errorf("unknown error %v while replying to message with CallError", err)
	}
}

func (cs *chargingStation) handleIncomingRequest(request ocpp.Request, requestId string, action string) {
	profile, found := cs.client.GetProfileForFeature(action)
	// Check whether action is supported and a listener for it exists
	if !found {
		cs.notImplementedError(requestId, action)
		return
	} else {
		switch profile.Name {
		case CoreProfileName:
			if cs.messageHandler == nil {
				cs.notSupportedError(requestId, action)
				return
			}
		case security.ProfileName:
			if cs.securityHandler == nil {
				cs.notSupportedError(requestId, action)
				return
			}
		case provisioning.ProfileName:
			if cs.provisioningHandler == nil {
				cs.notSupportedError(requestId, action)
				return
			}
		case authorization.ProfileName:
			if cs.authorizationHandler == nil {
				cs.notSupportedError(requestId, action)
				return
			}
			//case LocalAuthListProfileName:
			//	if cp.localAuthListListener == nil {
			//		cp.notSupportedError(requestId, action)
			//		return
			//	}
			//case FirmwareManagementProfileName:
			//	if cp.firmwareListener == nil {
			//		cp.notSupportedError(requestId, action)
			//		return
			//	}
			//case ReservationProfileName:
			//	if cp.reservationListener == nil {
			//		cp.notSupportedError(requestId, action)
			//		return
			//	}
			//case RemoteTriggerProfileName:
			//	if cp.remoteTriggerListener == nil {
			//		cp.notSupportedError(requestId, action)
			//		return
			//	}
			//case SmartChargingProfileName:
			//	if cp.smartChargingListener == nil {
			//		cp.notSupportedError(requestId, action)
			//		return
			//	}
		}
	}
	// Process request
	var confirmation ocpp.Response = nil
	cs.client.GetProfileForFeature(action)
	var err error = nil
	switch action {
	case CancelReservationFeatureName:
		confirmation, err = cs.messageHandler.OnCancelReservation(request.(*CancelReservationRequest))
	case security.CertificateSignedFeatureName:
		confirmation, err = cs.securityHandler.OnCertificateSigned(request.(*security.CertificateSignedRequest))
	case ChangeAvailabilityFeatureName:
		confirmation, err = cs.messageHandler.OnChangeAvailability(request.(*ChangeAvailabilityRequest))
	//case ChangeConfigurationFeatureName:
	//	confirmation, err = cp.messageHandler.OnChangeConfiguration(request.(*ChangeConfigurationRequest))
	case ClearCacheFeatureName:
		confirmation, err = cs.messageHandler.OnClearCache(request.(*ClearCacheRequest))
	case ClearChargingProfileFeatureName:
		confirmation, err = cs.messageHandler.OnClearChargingProfile(request.(*ClearChargingProfileRequest))
	case ClearDisplayFeatureName:
		confirmation, err = cs.messageHandler.OnClearDisplay(request.(*ClearDisplayRequest))
	case ClearVariableMonitoringFeatureName:
		confirmation, err = cs.messageHandler.OnClearVariableMonitoring(request.(*ClearVariableMonitoringRequest))
	case CostUpdatedFeatureName:
		confirmation, err = cs.messageHandler.OnCostUpdated(request.(*CostUpdatedRequest))
	case CustomerInformationFeatureName:
		confirmation, err = cs.messageHandler.OnCustomerInformation(request.(*CustomerInformationRequest))
	case DataTransferFeatureName:
		confirmation, err = cs.messageHandler.OnDataTransfer(request.(*DataTransferRequest))
	case DeleteCertificateFeatureName:
		confirmation, err = cs.messageHandler.OnDeleteCertificate(request.(*DeleteCertificateRequest))
	case provisioning.GetBaseReportFeatureName:
		confirmation, err = cs.provisioningHandler.OnGetBaseReport(request.(*provisioning.GetBaseReportRequest))
	case GetChargingProfilesFeatureName:
		confirmation, err = cs.messageHandler.OnGetChargingProfiles(request.(*GetChargingProfilesRequest))
	case GetCompositeScheduleFeatureName:
		confirmation, err = cs.messageHandler.OnGetCompositeSchedule(request.(*GetCompositeScheduleRequest))
	case GetDisplayMessagesFeatureName:
		confirmation, err = cs.messageHandler.OnGetDisplayMessages(request.(*GetDisplayMessagesRequest))
	case GetInstalledCertificateIdsFeatureName:
		confirmation, err = cs.messageHandler.OnGetInstalledCertificateIds(request.(*GetInstalledCertificateIdsRequest))
	case GetLocalListVersionFeatureName:
		confirmation, err = cs.messageHandler.OnGetLocalListVersion(request.(*GetLocalListVersionRequest))
	case GetLogFeatureName:
		confirmation, err = cs.messageHandler.OnGetLog(request.(*GetLogRequest))
	case GetMonitoringReportFeatureName:
		confirmation, err = cs.messageHandler.OnGetMonitoringReport(request.(*GetMonitoringReportRequest))
	//case GetConfigurationFeatureName:
	//	confirmation, err = cp.messageHandler.OnGetConfiguration(request.(*GetConfigurationRequest))
	//case RemoteStartTransactionFeatureName:
	//	confirmation, err = cp.messageHandler.OnRemoteStartTransaction(request.(*RemoteStartTransactionRequest))
	//case RemoteStopTransactionFeatureName:
	//	confirmation, err = cp.messageHandler.OnRemoteStopTransaction(request.(*RemoteStopTransactionRequest))
	//case ResetFeatureName:
	//	confirmation, err = cp.messageHandler.OnReset(request.(*ResetRequest))
	//case UnlockConnectorFeatureName:
	//	confirmation, err = cp.messageHandler.OnUnlockConnector(request.(*UnlockConnectorRequest))
	//case GetLocalListVersionFeatureName:
	//	confirmation, err = cp.localAuthListListener.OnGetLocalListVersion(request.(*GetLocalListVersionRequest))
	//case SendLocalListFeatureName:
	//	confirmation, err = cp.localAuthListListener.OnSendLocalList(request.(*SendLocalListRequest))
	//case GetDiagnosticsFeatureName:
	//	confirmation, err = cp.firmwareListener.OnGetDiagnostics(request.(*GetDiagnosticsRequest))
	//case UpdateFirmwareFeatureName:
	//	confirmation, err = cp.firmwareListener.OnUpdateFirmware(request.(*UpdateFirmwareRequest))
	//case ReserveNowFeatureName:
	//	confirmation, err = cp.reservationListener.OnReserveNow(request.(*ReserveNowRequest))
	//case CancelReservationFeatureName:
	//	confirmation, err = cp.reservationListener.OnCancelReservation(request.(*CancelReservationRequest))
	//case TriggerMessageFeatureName:
	//	confirmation, err = cp.remoteTriggerListener.OnTriggerMessage(request.(*TriggerMessageRequest))
	//case SetChargingProfileFeatureName:
	//	confirmation, err = cp.smartChargingListener.OnSetChargingProfile(request.(*SetChargingProfileRequest))
	//case ClearChargingProfileFeatureName:
	//	confirmation, err = cp.smartChargingListener.OnClearChargingProfile(request.(*ClearChargingProfileRequest))
	//case GetCompositeScheduleFeatureName:
	//	confirmation, err = cp.smartChargingListener.OnGetCompositeSchedule(request.(*GetCompositeScheduleRequest))
	default:
		cs.notSupportedError(requestId, action)
		return
	}
	cs.sendResponse(confirmation, err, requestId)
}
