package ocpp16

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocppj"
	"github.com/lorenzodonini/ocpp-go/ws"
	log "github.com/sirupsen/logrus"
)

// -------------------- v1.6 Charge Point --------------------
type ChargePoint interface {
	// Message
	BootNotification(chargePointModel string, chargePointVendor string, props ...func(request *BootNotificationRequest)) (*BootNotificationConfirmation, error)
	Authorize(idTag string, props ...func(request *AuthorizeRequest)) (*AuthorizeConfirmation, error)
	DataTransfer(vendorId string, props ...func(request *DataTransferRequest)) (*DataTransferConfirmation, error)
	Heartbeat(props ...func(request *HeartbeatRequest)) (*HeartbeatConfirmation, error)
	MeterValues(connectorId int, meterValues []MeterValue, props ...func(request *MeterValuesRequest)) (*MeterValuesConfirmation, error)
	StartTransaction(connectorId int, idTag string, meterStart int, timestamp *DateTime, props ...func(request *StartTransactionRequest)) (*StartTransactionConfirmation, error)
	StopTransaction(meterStop int, timestamp *DateTime, transactionId int, props ...func(request *StopTransactionRequest)) (*StopTransactionConfirmation, error)
	StatusNotification(connectorId int, errorCode ChargePointErrorCode, status ChargePointStatus, props ...func(request *StatusNotificationRequest)) (*StatusNotificationConfirmation, error)
	DiagnosticsStatusNotification(status DiagnosticsStatus, props ...func(request *DiagnosticsStatusNotificationRequest)) (*DiagnosticsStatusNotificationConfirmation, error)
	FirmwareStatusNotification(status FirmwareStatus, props ...func(request *FirmwareStatusNotificationRequest)) (*FirmwareStatusNotificationConfirmation, error)
	//TODO: add missing profile methods

	// Logic
	SetChargePointCoreListener(listener ChargePointCoreListener)
	SetLocalAuthListListener(listener ChargePointLocalAuthListListener)
	SetFirmwareManagementListener(listener ChargePointFirmwareManagementListener)
	SendRequest(request ocpp.Request) (ocpp.Confirmation, error)
	SendRequestAsync(request ocpp.Request, callback func(confirmation ocpp.Confirmation, protoError error)) error
	Start(centralSystemUrl string) error
	Stop()
}

type chargePoint struct {
	chargePoint           *ocppj.ChargePoint
	coreListener          ChargePointCoreListener
	localAuthListListener ChargePointLocalAuthListListener
	firmwareListener      ChargePointFirmwareManagementListener
	confirmationListener  chan ocpp.Confirmation
	errorListener         chan error
}

func (cp *chargePoint) BootNotification(chargePointModel string, chargePointVendor string, props ...func(request *BootNotificationRequest)) (*BootNotificationConfirmation, error) {
	request := NewBootNotificationRequest(chargePointModel, chargePointVendor)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cp.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*BootNotificationConfirmation), err
	}
}

func (cp *chargePoint) Authorize(idTag string, props ...func(request *AuthorizeRequest)) (*AuthorizeConfirmation, error) {
	request := NewAuthorizationRequest(idTag)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cp.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*AuthorizeConfirmation), err
	}
}

func (cp *chargePoint) DataTransfer(vendorId string, props ...func(request *DataTransferRequest)) (*DataTransferConfirmation, error) {
	request := NewDataTransferRequest(vendorId)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cp.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*DataTransferConfirmation), err
	}
}

func (cp *chargePoint) Heartbeat(props ...func(request *HeartbeatRequest)) (*HeartbeatConfirmation, error) {
	request := NewHeartbeatRequest()
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cp.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*HeartbeatConfirmation), err
	}
}

func (cp *chargePoint) MeterValues(connectorId int, meterValues []MeterValue, props ...func(request *MeterValuesRequest)) (*MeterValuesConfirmation, error) {
	request := NewMeterValuesRequest(connectorId, meterValues)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cp.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*MeterValuesConfirmation), err
	}
}

func (cp *chargePoint) StartTransaction(connectorId int, idTag string, meterStart int, timestamp *DateTime, props ...func(request *StartTransactionRequest)) (*StartTransactionConfirmation, error) {
	request := NewStartTransactionRequest(connectorId, idTag, meterStart, timestamp)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cp.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*StartTransactionConfirmation), err
	}
}

func (cp *chargePoint) StopTransaction(meterStop int, timestamp *DateTime, transactionId int, props ...func(request *StopTransactionRequest)) (*StopTransactionConfirmation, error) {
	request := NewStopTransactionRequest(meterStop, timestamp, transactionId)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cp.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*StopTransactionConfirmation), err
	}
}

func (cp *chargePoint) StatusNotification(connectorId int, errorCode ChargePointErrorCode, status ChargePointStatus, props ...func(request *StatusNotificationRequest)) (*StatusNotificationConfirmation, error) {
	request := NewStatusNotificationRequest(connectorId, errorCode, status)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cp.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*StatusNotificationConfirmation), err
	}
}

func (cp *chargePoint) DiagnosticsStatusNotification(status DiagnosticsStatus, props ...func(request *DiagnosticsStatusNotificationRequest)) (*DiagnosticsStatusNotificationConfirmation, error) {
	request := NewDiagnosticsStatusNotificationRequest(status)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cp.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*DiagnosticsStatusNotificationConfirmation), err
	}
}

func (cp *chargePoint) FirmwareStatusNotification(status FirmwareStatus, props ...func(request *FirmwareStatusNotificationRequest)) (*FirmwareStatusNotificationConfirmation, error) {
	request := NewFirmwareStatusNotificationRequest(status)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cp.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*FirmwareStatusNotificationConfirmation), err
	}
}

func (cp *chargePoint) SetChargePointCoreListener(listener ChargePointCoreListener) {
	cp.coreListener = listener
}

func (cp *chargePoint) SetLocalAuthListListener(listener ChargePointLocalAuthListListener) {
	cp.localAuthListListener = listener
}

func (cp *chargePoint) SetFirmwareManagementListener(listener ChargePointFirmwareManagementListener) {
	cp.firmwareListener = listener
}

func (cp *chargePoint) SendRequest(request ocpp.Request) (ocpp.Confirmation, error) {
	// TODO: check for supported feature
	err := cp.chargePoint.SendRequest(request)
	if err != nil {
		return nil, err
	}
	//TODO: timeouts
	select {
	case confirmation := <-cp.confirmationListener:
		return confirmation, nil
	case err = <-cp.errorListener:
		return nil, err
	}
}

func (cp *chargePoint) SendRequestAsync(request ocpp.Request, callback func(confirmation ocpp.Confirmation, err error)) error {
	switch request.GetFeatureName() {
	case AuthorizeFeatureName, BootNotificationFeatureName, DataTransferFeatureName, HeartbeatFeatureName, MeterValuesFeatureName, StartTransactionFeatureName, StopTransactionFeatureName, StatusNotificationFeatureName, DiagnosticsStatusNotificationFeatureName, FirmwareStatusNotificationFeatureName:
		break
	default:
		return fmt.Errorf("unsupported action %v on charge point, cannot send request", request.GetFeatureName())
	}
	err := cp.chargePoint.SendRequest(request)
	if err == nil {
		// Retrieve result asynchronously
		go func() {
			select {
			case confirmation := <-cp.confirmationListener:
				callback(confirmation, nil)
			case protoError := <-cp.errorListener:
				callback(nil, protoError)
			}
		}()
	}
	return err
}

func (cp *chargePoint) sendResponse(confirmation ocpp.Confirmation, err error, requestId string) {
	if confirmation != nil {
		err := cp.chargePoint.SendConfirmation(requestId, confirmation)
		if err != nil {
			log.WithField("request", requestId).Errorf("unknown error %v while replying to message with CallError", err)
			//TODO: handle error somehow
		}
	} else {
		err = cp.chargePoint.SendError(requestId, ocppj.ProtocolError, err.Error(), nil)
		if err != nil {
			log.WithField("request", requestId).Errorf("unknown error %v while replying to message with CallError", err)
		}
	}
}

func (cp *chargePoint) Start(centralSystemUrl string) error {
	// TODO: implement auto-reconnect logic
	return cp.chargePoint.Start(centralSystemUrl)
}

func (cp *chargePoint) Stop() {
	cp.chargePoint.Stop()
}

func (cp *chargePoint) notImplementedError(requestId string, action string) {
	log.WithField("request", requestId).Errorf("cannot handle call from central system. Sending CallError instead")
	err := cp.chargePoint.SendError(requestId, ocppj.NotImplemented, fmt.Sprintf("no handler for action %v implemented", action), nil)
	if err != nil {
		log.WithField("request", requestId).Errorf("unknown error %v while replying to message with CallError", err)
	}
}

func (cp *chargePoint) notSupportedError(requestId string, action string) {
	log.WithField("request", requestId).Errorf("cannot handle call from central system. Sending CallError instead")
	err := cp.chargePoint.SendError(requestId, ocppj.NotSupported, fmt.Sprintf("unsupported action %v on charge point", action), nil)
	if err != nil {
		log.WithField("request", requestId).Errorf("unknown error %v while replying to message with CallError", err)
	}
}

func (cp *chargePoint) handleIncomingRequest(request ocpp.Request, requestId string, action string) {
	profile, found := cp.chargePoint.GetProfileForFeature(action)
	// Check whether action is supported and a listener for it exists
	if !found {
		cp.notImplementedError(requestId, action)
		return
	} else {
		switch profile.Name {
		case CoreProfileName:
			if cp.coreListener == nil {
				cp.notSupportedError(requestId, action)
				return
			}
		case LocalAuthListProfileName:
			if cp.localAuthListListener == nil {
				cp.notSupportedError(requestId, action)
				return
			}
		case FirmwareManagementProfileName:
			if cp.firmwareListener == nil {
				cp.notSupportedError(requestId, action)
				return
			}
		}
	}
	// Process request
	var confirmation ocpp.Confirmation = nil
	cp.chargePoint.GetProfileForFeature(action)
	var err error = nil
	switch action {
	case ChangeAvailabilityFeatureName:
		confirmation, err = cp.coreListener.OnChangeAvailability(request.(*ChangeAvailabilityRequest))
	case ChangeConfigurationFeatureName:
		confirmation, err = cp.coreListener.OnChangeConfiguration(request.(*ChangeConfigurationRequest))
	case ClearCacheFeatureName:
		confirmation, err = cp.coreListener.OnClearCache(request.(*ClearCacheRequest))
	case DataTransferFeatureName:
		confirmation, err = cp.coreListener.OnDataTransfer(request.(*DataTransferRequest))
	case GetConfigurationFeatureName:
		confirmation, err = cp.coreListener.OnGetConfiguration(request.(*GetConfigurationRequest))
	case RemoteStartTransactionFeatureName:
		confirmation, err = cp.coreListener.OnRemoteStartTransaction(request.(*RemoteStartTransactionRequest))
	case RemoteStopTransactionFeatureName:
		confirmation, err = cp.coreListener.OnRemoteStopTransaction(request.(*RemoteStopTransactionRequest))
	case ResetFeatureName:
		confirmation, err = cp.coreListener.OnReset(request.(*ResetRequest))
	case UnlockConnectorFeatureName:
		confirmation, err = cp.coreListener.OnUnlockConnector(request.(*UnlockConnectorRequest))
	case GetLocalListVersionFeatureName:
		confirmation, err = cp.localAuthListListener.OnGetLocalListVersion(request.(*GetLocalListVersionRequest))
	case SendLocalListFeatureName:
		confirmation, err = cp.localAuthListListener.OnSendLocalList(request.(*SendLocalListRequest))
	case GetDiagnosticsFeatureName:
		confirmation, err = cp.firmwareListener.OnGetDiagnostics(request.(*GetDiagnosticsRequest))
	case UpdateFirmwareFeatureName:
		confirmation, err = cp.firmwareListener.OnUpdateFirmware(request.(*UpdateFirmwareRequest))
	default:
		cp.notSupportedError(requestId, action)
		return
	}
	cp.sendResponse(confirmation, err, requestId)
}

func NewChargePoint(id string, dispatcher *ocppj.ChargePoint, client ws.WsClient) ChargePoint {
	if client == nil {
		client = ws.NewClient()
	}
	if dispatcher == nil {
		dispatcher = ocppj.NewChargePoint(id, client, CoreProfile, LocalAuthListProfile, FirmwareManagementProfile)
	}
	cp := chargePoint{chargePoint: dispatcher, confirmationListener: make(chan ocpp.Confirmation), errorListener: make(chan error)}
	cp.chargePoint.SetConfirmationHandler(func(confirmation ocpp.Confirmation, requestId string) {
		cp.confirmationListener <- confirmation
	})
	cp.chargePoint.SetErrorHandler(func(err *ocpp.Error, details interface{}) {
		cp.errorListener <- err
	})
	cp.chargePoint.SetRequestHandler(cp.handleIncomingRequest)
	return &cp
}

// -------------------- v1.6 Central System --------------------
type CentralSystem interface {
	// Message
	//TODO: add missing profile methods
	ChangeAvailability(clientId string, callback func(*ChangeAvailabilityConfirmation, error), connectorId int, availabilityType AvailabilityType, props ...func(*ChangeAvailabilityRequest)) error
	ChangeConfiguration(clientId string, callback func(*ChangeConfigurationConfirmation, error), key string, value string, props ...func(*ChangeConfigurationRequest)) error
	ClearCache(clientId string, callback func(*ClearCacheConfirmation, error), props ...func(*ClearCacheRequest)) error
	DataTransfer(clientId string, callback func(*DataTransferConfirmation, error), vendorId string, props ...func(*DataTransferRequest)) error
	GetConfiguration(clientId string, callback func(*GetConfigurationConfirmation, error), keys []string, props ...func(*GetConfigurationRequest)) error
	RemoteStartTransaction(clientId string, callback func(*RemoteStartTransactionConfirmation, error), idTag string, props ...func(*RemoteStartTransactionRequest)) error
	RemoteStopTransaction(clientId string, callback func(*RemoteStopTransactionConfirmation, error), transactionId int, props ...func(request *RemoteStopTransactionRequest)) error
	Reset(clientId string, callback func(*ResetConfirmation, error), resetType ResetType, props ...func(*ResetRequest)) error
	UnlockConnector(clientId string, callback func(*UnlockConnectorConfirmation, error), connectorId int, props ...func(*UnlockConnectorRequest)) error
	GetLocalListVersion(clientId string, callback func(*GetLocalListVersionConfirmation, error), props ...func(request *GetLocalListVersionRequest)) error
	SendLocalList(clientId string, callback func(*SendLocalListConfirmation, error), version int, updateType UpdateType, props ...func(request *SendLocalListRequest)) error
	GetDiagnostics(clientId string, callback func(*GetDiagnosticsConfirmation, error), location string, props ...func(request *GetDiagnosticsRequest)) error
	UpdateFirmware(clientId string, callback func(*UpdateFirmwareConfirmation, error), location string, retrieveDate *DateTime, props ...func(request *UpdateFirmwareRequest)) error
	// Logic
	SetCentralSystemCoreListener(listener CentralSystemCoreListener)
	SetLocalAuthListListener(listener CentralSystemLocalAuthListListener)
	SetFirmwareManagementListener(listener CentralSystemFirmwareManagementListener)
	SetNewChargePointHandler(handler func(chargePointId string))
	SetChargePointDisconnectedHandler(handler func(chargePointId string))
	SendRequestAsync(clientId string, request ocpp.Request, callback func(ocpp.Confirmation, error)) error
	Start(listenPort int, listenPath string)
}

type centralSystem struct {
	centralSystem         *ocppj.CentralSystem
	coreListener          CentralSystemCoreListener
	localAuthListListener CentralSystemLocalAuthListListener
	firmwareListener      CentralSystemFirmwareManagementListener
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

func (cs *centralSystem) SetCentralSystemCoreListener(listener CentralSystemCoreListener) {
	cs.coreListener = listener
}

func (cs *centralSystem) SetLocalAuthListListener(listener CentralSystemLocalAuthListListener) {
	cs.localAuthListListener = listener
}

func (cs *centralSystem) SetFirmwareManagementListener(listener CentralSystemFirmwareManagementListener) {
	cs.firmwareListener = listener
}

func (cs *centralSystem) SetNewChargePointHandler(handler func(chargePointId string)) {
	cs.centralSystem.SetNewChargePointHandler(handler)
}

func (cs *centralSystem) SetChargePointDisconnectedHandler(handler func(chargePointId string)) {
	cs.centralSystem.SetDisconnectedChargePointHandler(handler)
}

func (cs *centralSystem) SendRequestAsync(clientId string, request ocpp.Request, callback func(confirmation ocpp.Confirmation, err error)) error {
	switch request.GetFeatureName() {
	case ChangeAvailabilityFeatureName, ChangeConfigurationFeatureName, ClearCacheFeatureName, DataTransferFeatureName, GetConfigurationFeatureName, RemoteStartTransactionFeatureName, RemoteStopTransactionFeatureName, ResetFeatureName, UnlockConnectorFeatureName, GetLocalListVersionFeatureName, SendLocalListFeatureName, GetDiagnosticsFeatureName, UpdateFirmwareFeatureName:
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

func NewCentralSystem(dispatcher *ocppj.CentralSystem, server ws.WsServer) CentralSystem {
	if server == nil {
		server = ws.NewServer()
	}
	if dispatcher == nil {
		dispatcher = ocppj.NewCentralSystem(server, CoreProfile, LocalAuthListProfile, FirmwareManagementProfile)
	}
	cs := centralSystem{
		centralSystem: dispatcher,
		callbacks:     map[string]func(confirmation ocpp.Confirmation, err error){}}
	cs.centralSystem.SetRequestHandler(cs.handleIncomingRequest)
	cs.centralSystem.SetConfirmationHandler(cs.handleIncomingConfirmation)
	cs.centralSystem.SetErrorHandler(cs.handleIncomingError)
	return &cs
}

func init() {
	log.New()
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetLevel(log.InfoLevel)
}
