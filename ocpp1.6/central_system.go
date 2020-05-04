package ocpp16

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocppj"
	log "github.com/sirupsen/logrus"
)

type centralSystem struct {
	centralSystem         *ocppj.CentralSystem
	coreListener          CentralSystemCoreHandler
	localAuthListListener CentralSystemLocalAuthListHandler
	firmwareListener      CentralSystemFirmwareManagementHandler
	reservationListener   CentralSystemReservationHandler
	remoteTriggerListener CentralSystemRemoteTriggerHandler
	smartChargingListener CentralSystemSmartChargingHandler
	callbacks             map[string]func(confirmation ocpp.Confirmation, err error)
}

func (cs *centralSystem) ChangeAvailability(clientId string, callback func(confirmation *ChangeAvailabilityConfirmation, err error), connectorId int, availabilityType AvailabilityType, props ...func(request *ChangeAvailabilityRequest)) error {
	request := NewChangeAvailabilityRequest(connectorId, availabilityType)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*ChangeAvailabilityConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *centralSystem) ChangeConfiguration(clientId string, callback func(confirmation *ChangeConfigurationConfirmation, err error), key string, value string, props ...func(request *ChangeConfigurationRequest)) error {
	request := NewChangeConfigurationRequest(key, value)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*ChangeConfigurationConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *centralSystem) ClearCache(clientId string, callback func(confirmation *ClearCacheConfirmation, err error), props ...func(*ClearCacheRequest)) error {
	request := NewClearCacheRequest()
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*ClearCacheConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *centralSystem) DataTransfer(clientId string, callback func(confirmation *DataTransferConfirmation, err error), vendorId string, props ...func(request *DataTransferRequest)) error {
	request := NewDataTransferRequest(vendorId)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*DataTransferConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *centralSystem) GetConfiguration(clientId string, callback func(confirmation *GetConfigurationConfirmation, err error), keys []string, props ...func(request *GetConfigurationRequest)) error {
	request := NewGetConfigurationRequest(keys)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*GetConfigurationConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *centralSystem) RemoteStartTransaction(clientId string, callback func(*RemoteStartTransactionConfirmation, error), idTag string, props ...func(*RemoteStartTransactionRequest)) error {
	request := NewRemoteStartTransactionRequest(idTag)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*RemoteStartTransactionConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *centralSystem) RemoteStopTransaction(clientId string, callback func(*RemoteStopTransactionConfirmation, error), transactionId int, props ...func(request *RemoteStopTransactionRequest)) error {
	request := NewRemoteStopTransactionRequest(transactionId)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*RemoteStopTransactionConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *centralSystem) Reset(clientId string, callback func(*ResetConfirmation, error), resetType ResetType, props ...func(request *ResetRequest)) error {
	request := NewResetRequest(resetType)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*ResetConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *centralSystem) UnlockConnector(clientId string, callback func(*UnlockConnectorConfirmation, error), connectorId int, props ...func(*UnlockConnectorRequest)) error {
	request := NewUnlockConnectorRequest(connectorId)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*UnlockConnectorConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *centralSystem) GetLocalListVersion(clientId string, callback func(*GetLocalListVersionConfirmation, error), props ...func(request *GetLocalListVersionRequest)) error {
	request := NewGetLocalListVersionRequest()
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*GetLocalListVersionConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *centralSystem) SendLocalList(clientId string, callback func(*SendLocalListConfirmation, error), version int, updateType UpdateType, props ...func(request *SendLocalListRequest)) error {
	request := NewSendLocalListRequest(version, updateType)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*SendLocalListConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *centralSystem) GetDiagnostics(clientId string, callback func(*GetDiagnosticsConfirmation, error), location string, props ...func(request *GetDiagnosticsRequest)) error {
	request := NewGetDiagnosticsRequest(location)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*GetDiagnosticsConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *centralSystem) UpdateFirmware(clientId string, callback func(*UpdateFirmwareConfirmation, error), location string, retrieveDate *DateTime, props ...func(request *UpdateFirmwareRequest)) error {
	request := NewUpdateFirmwareRequest(location, retrieveDate)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*UpdateFirmwareConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *centralSystem) ReserveNow(clientId string, callback func(*ReserveNowConfirmation, error), connectorId int, expiryDate *DateTime, idTag string, reservationId int, props ...func(request *ReserveNowRequest)) error {
	request := NewReserveNowRequest(connectorId, expiryDate, idTag, reservationId)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*ReserveNowConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *centralSystem) CancelReservation(clientId string, callback func(*CancelReservationConfirmation, error), reservationId int, props ...func(request *CancelReservationRequest)) error {
	request := NewCancelReservationRequest(reservationId)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*CancelReservationConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *centralSystem) TriggerMessage(clientId string, callback func(*TriggerMessageConfirmation, error), requestedMessage MessageTrigger, props ...func(request *TriggerMessageRequest)) error {
	request := NewTriggerMessageRequest(requestedMessage)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*TriggerMessageConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *centralSystem) SetChargingProfile(clientId string, callback func(*SetChargingProfileConfirmation, error), connectorId int, chargingProfile *ChargingProfile, props ...func(request *SetChargingProfileRequest)) error {
	request := NewSetChargingProfileRequest(connectorId, chargingProfile)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*SetChargingProfileConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *centralSystem) ClearChargingProfile(clientId string, callback func(*ClearChargingProfileConfirmation, error), props ...func(request *ClearChargingProfileRequest)) error {
	request := NewClearChargingProfileRequest()
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*ClearChargingProfileConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *centralSystem) GetCompositeSchedule(clientId string, callback func(*GetCompositeScheduleConfirmation, error), connectorId int, duration int, props ...func(request *GetCompositeScheduleRequest)) error {
	request := NewGetCompositeScheduleRequest(connectorId, duration)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*GetCompositeScheduleConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *centralSystem) SetCentralSystemCoreHandler(listener CentralSystemCoreHandler) {
	cs.coreListener = listener
}

func (cs *centralSystem) SetLocalAuthListHandler(listener CentralSystemLocalAuthListHandler) {
	cs.localAuthListListener = listener
}

func (cs *centralSystem) SetFirmwareManagementHandler(listener CentralSystemFirmwareManagementHandler) {
	cs.firmwareListener = listener
}

func (cs *centralSystem) SetReservationHandler(listener CentralSystemReservationHandler) {
	cs.reservationListener = listener
}

func (cs *centralSystem) SetRemoteTriggerHandler(listener CentralSystemRemoteTriggerHandler) {
	cs.remoteTriggerListener = listener
}

func (cs *centralSystem) SetSmartChargingHandler(listener CentralSystemSmartChargingHandler) {
	cs.smartChargingListener = listener
}

func (cs *centralSystem) SetNewChargePointHandler(handler func(chargePointId string)) {
	cs.centralSystem.SetNewChargePointHandler(handler)
}

func (cs *centralSystem) SetChargePointDisconnectedHandler(handler func(chargePointId string)) {
	cs.centralSystem.SetDisconnectedChargePointHandler(handler)
}

func (cs *centralSystem) SendRequestAsync(clientId string, request ocpp.Request, callback func(confirmation ocpp.Confirmation, err error)) error {
	switch request.GetFeatureName() {
	case ChangeAvailabilityFeatureName, ChangeConfigurationFeatureName, ClearCacheFeatureName, DataTransferFeatureName, GetConfigurationFeatureName, RemoteStartTransactionFeatureName, RemoteStopTransactionFeatureName, ResetFeatureName, UnlockConnectorFeatureName,
		GetLocalListVersionFeatureName, SendLocalListFeatureName,
		GetDiagnosticsFeatureName, UpdateFirmwareFeatureName,
		ReserveNowFeatureName, CancelReservationFeatureName,
		TriggerMessageFeatureName,
		SetChargingProfileFeatureName, ClearChargingProfileFeatureName, GetCompositeScheduleFeatureName:
	default:
		return fmt.Errorf("unsupported action %v on central system, cannot send request", request.GetFeatureName())
	}
	cs.callbacks[clientId] = callback
	err := cs.centralSystem.SendRequest(clientId, request)
	if err != nil {
		delete(cs.callbacks, clientId)
		return err
	}
	return nil
}

func (cs *centralSystem) Start(listenPort int, listenPath string) {
	cs.centralSystem.Start(listenPort, listenPath)
}

func (cs *centralSystem) sendResponse(chargePointId string, confirmation ocpp.Confirmation, err error, requestId string) {
	if confirmation != nil {
		err := cs.centralSystem.SendConfirmation(chargePointId, requestId, confirmation)
		if err != nil {
			//TODO: handle error somehow
			log.Print(err)
		}
	} else {
		err := cs.centralSystem.SendError(chargePointId, requestId, ocppj.ProtocolError, "Couldn't generate valid confirmation", nil)
		if err != nil {
			log.WithFields(log.Fields{
				"client":  chargePointId,
				"request": requestId,
			}).Errorf("unknown error %v while replying to message with CallError", err)
		}
	}
}

func (cs *centralSystem) notImplementedError(chargePointId string, requestId string, action string) {
	log.Warnf("Cannot handle call %v from charge point %v. Sending CallError instead", requestId, chargePointId)
	err := cs.centralSystem.SendError(chargePointId, requestId, ocppj.NotImplemented, fmt.Sprintf("no handler for action %v implemented", action), nil)
	if err != nil {
		log.WithFields(log.Fields{
			"client":  chargePointId,
			"request": requestId,
		}).Errorf("unknown error %v while replying to message with CallError", err)
	}
}

func (cs *centralSystem) notSupportedError(chargePointId string, requestId string, action string) {
	log.Warnf("Cannot handle call %v from charge point %v. Sending CallError instead", requestId, chargePointId)
	err := cs.centralSystem.SendError(chargePointId, requestId, ocppj.NotSupported, fmt.Sprintf("unsupported action %v on central system", action), nil)
	if err != nil {
		log.WithFields(log.Fields{
			"client":  chargePointId,
			"request": requestId,
		}).Errorf("unknown error %v while replying to message with CallError", err)
	}
}

func (cs *centralSystem) handleIncomingRequest(chargePointId string, request ocpp.Request, requestId string, action string) {
	profile, found := cs.centralSystem.GetProfileForFeature(action)
	// Check whether action is supported and a listener for it exists
	if !found {
		cs.notImplementedError(chargePointId, requestId, action)
		return
	} else {
		switch profile.Name {
		case CoreProfileName:
			if cs.coreListener == nil {
				cs.notSupportedError(chargePointId, requestId, action)
				return
			}
		case LocalAuthListProfileName:
			if cs.localAuthListListener == nil {
				cs.notSupportedError(chargePointId, requestId, action)
				return
			}
		case FirmwareManagementProfileName:
			if cs.firmwareListener == nil {
				cs.notSupportedError(chargePointId, requestId, action)
				return
			}
		case ReservationProfileName:
			if cs.reservationListener == nil {
				cs.notSupportedError(chargePointId, requestId, action)
				return
			}
		case RemoteTriggerProfileName:
			if cs.remoteTriggerListener == nil {
				cs.notSupportedError(chargePointId, requestId, action)
				return
			}
		case SmartChargingProfileName:
			if cs.smartChargingListener == nil {
				cs.notSupportedError(chargePointId, requestId, action)
				return
			}
		}
	}
	var confirmation ocpp.Confirmation = nil
	var err error = nil
	// Execute in separate goroutine, so the caller goroutine is available
	go func() {
		switch action {
		case BootNotificationFeatureName:
			confirmation, err = cs.coreListener.OnBootNotification(chargePointId, request.(*BootNotificationRequest))
		case AuthorizeFeatureName:
			confirmation, err = cs.coreListener.OnAuthorize(chargePointId, request.(*AuthorizeRequest))
		case DataTransferFeatureName:
			confirmation, err = cs.coreListener.OnDataTransfer(chargePointId, request.(*DataTransferRequest))
		case HeartbeatFeatureName:
			confirmation, err = cs.coreListener.OnHeartbeat(chargePointId, request.(*HeartbeatRequest))
		case MeterValuesFeatureName:
			confirmation, err = cs.coreListener.OnMeterValues(chargePointId, request.(*MeterValuesRequest))
		case StartTransactionFeatureName:
			confirmation, err = cs.coreListener.OnStartTransaction(chargePointId, request.(*StartTransactionRequest))
		case StopTransactionFeatureName:
			confirmation, err = cs.coreListener.OnStopTransaction(chargePointId, request.(*StopTransactionRequest))
		case StatusNotificationFeatureName:
			confirmation, err = cs.coreListener.OnStatusNotification(chargePointId, request.(*StatusNotificationRequest))
		case DiagnosticsStatusNotificationFeatureName:
			confirmation, err = cs.firmwareListener.OnDiagnosticsStatusNotification(chargePointId, request.(*DiagnosticsStatusNotificationRequest))
		case FirmwareStatusNotificationFeatureName:
			confirmation, err = cs.firmwareListener.OnFirmwareStatusNotification(chargePointId, request.(*FirmwareStatusNotificationRequest))
		default:
			cs.notSupportedError(chargePointId, requestId, action)
			return
		}
		cs.sendResponse(chargePointId, confirmation, err, requestId)
	}()
}

func (cs *centralSystem) handleIncomingConfirmation(chargePointId string, confirmation ocpp.Confirmation, requestId string) {
	if callback, ok := cs.callbacks[chargePointId]; ok {
		delete(cs.callbacks, chargePointId)
		callback(confirmation, nil)
	} else {
		log.WithFields(log.Fields{
			"client":  chargePointId,
			"request": requestId,
		}).Errorf("no handler available for Call Result of type %v", confirmation.GetFeatureName())
	}
}

func (cs *centralSystem) handleIncomingError(chargePointId string, err *ocpp.Error, details interface{}) {
	if callback, ok := cs.callbacks[chargePointId]; ok {
		delete(cs.callbacks, chargePointId)
		callback(nil, err)
	} else {
		log.WithFields(log.Fields{
			"client":  chargePointId,
			"request": err.MessageId,
		}).Errorf("no handler available for Call Error %v", err.Code)
	}
}
