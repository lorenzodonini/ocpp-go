package ocpp16

import (
	"fmt"
	"github.com/lorenzodonini/go-ocpp/ocpp"
	"github.com/lorenzodonini/go-ocpp/ocppj"
	"github.com/lorenzodonini/go-ocpp/ws"
	"log"
)

// -------------------- v1.6 Charge Point --------------------
type ChargePoint interface {
	// Message
	BootNotification(chargePointModel string, chargePointVendor string, props ...func(request *BootNotificationRequest)) (*BootNotificationConfirmation, error)
	Authorize(idTag string, props ...func(request *AuthorizeRequest)) (*AuthorizeConfirmation, error)
	DataTransfer(vendorId string, props ...func(request *DataTransferRequest)) (*DataTransferConfirmation, error)
	Heartbeat(props ...func(request *HeartbeatRequest)) (*HeartbeatConfirmation, error)
	//TODO: add missing profile methods

	// Logic
	SetChargePointCoreListener(listener ChargePointCoreListener)
	SendRequest(request ocpp.Request) (ocpp.Confirmation, error)
	SendRequestAsync(request ocpp.Request, callback func(confirmation ocpp.Confirmation, protoError error)) error
	Start(centralSystemUrl string) error
}

type chargePoint struct {
	chargePoint          *ocppj.ChargePoint
	coreListener         ChargePointCoreListener
	confirmationListener chan ocpp.Confirmation
	errorListener        chan error
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

func (cp *chargePoint) SetChargePointCoreListener(listener ChargePointCoreListener) {
	cp.coreListener = listener
}

func (cp *chargePoint) SendRequest(request ocpp.Request) (ocpp.Confirmation, error) {
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
	case AuthorizeFeatureName, BootNotificationFeatureName, DataTransferFeatureName, HeartbeatFeatureName:
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
			log.Printf("Unknown error %v while replying to message %v with CallError", err, requestId)
			//TODO: handle error somehow
		}
	} else {
		err = cp.chargePoint.SendError(requestId, ocppj.ProtocolError, err.Error(), nil)
		if err != nil {
			log.Printf("Unknown error %v while replying to message %v with CallError", err, requestId)
		}
	}
}

func (cp *chargePoint) Start(centralSystemUrl string) error {
	// TODO: implement auto-reconnect logic
	return cp.chargePoint.Start(centralSystemUrl)
}

func (cp *chargePoint) handleIncomingRequest(request ocpp.Request, requestId string, action string) {
	if cp.coreListener == nil {
		log.Printf("Cannot handle call %v from central system. Sending CallError instead", requestId)
		err := cp.chargePoint.SendError(requestId, ocppj.NotImplemented, fmt.Sprintf("No handler for action %v implemented", action), nil)
		if err != nil {
			log.Printf("Unknown error %v while replying to message %v with CallError", err, requestId)
		}
		return
	}
	var confirmation ocpp.Confirmation = nil
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
	default:
		err := cp.chargePoint.SendError(requestId, ocppj.NotSupported, fmt.Sprintf("unsupported action %v on charge point", action), nil)
		if err != nil {
			log.Printf("Unknown error %v while replying to message %v with CallError", err, requestId)
		}
		return
	}
	cp.sendResponse(confirmation, err, requestId)
}

func NewChargePoint(id string, dispatcher *ocppj.ChargePoint, client ws.WsClient) ChargePoint {
	if client == nil {
		client = ws.NewClient()
	}
	if dispatcher == nil {
		dispatcher = ocppj.NewChargePoint(id, client, CoreProfile)
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
	// Logic
	SetCentralSystemCoreListener(listener CentralSystemCoreListener)
	SetNewChargePointHandler(handler func(chargePointId string))
	SendRequestAsync(clientId string, request ocpp.Request, callback func(ocpp.Confirmation, error)) error
	Start(listenPort int, listenPath string)
}

type centralSystem struct {
	centralSystem *ocppj.CentralSystem
	coreListener  CentralSystemCoreListener
	callbacks     map[string]func(confirmation ocpp.Confirmation, err error)
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

func (cs *centralSystem) SetCentralSystemCoreListener(listener CentralSystemCoreListener) {
	cs.coreListener = listener
}

func (cs *centralSystem) SetNewChargePointHandler(handler func(chargePointId string)) {
	cs.centralSystem.SetNewChargePointHandler(handler)
}

func (cs *centralSystem) SendRequestAsync(clientId string, request ocpp.Request, callback func(confirmation ocpp.Confirmation, err error)) error {
	switch request.GetFeatureName() {
	case ChangeAvailabilityFeatureName, ChangeConfigurationFeatureName, ClearCacheFeatureName, DataTransferFeatureName, GetConfigurationFeatureName:
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
			log.Printf("Unknown error %v while replying to message %v with CallError", err, requestId)
		}
	}
}

func (cs *centralSystem) handleIncomingRequest(chargePointId string, request ocpp.Request, requestId string, action string) {
	if cs.coreListener == nil {
		log.Printf("Cannot handle call %v from charge point %v. Sending CallError instead", requestId, chargePointId)
		err := cs.centralSystem.SendError(chargePointId, requestId, ocppj.NotImplemented, fmt.Sprintf("No handler for action %v implemented", action), nil)
		if err != nil {
			log.Printf("Unknown error %v while replying to message %v with CallError", err, requestId)
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
		default:
			err := cs.centralSystem.SendError(chargePointId, requestId, ocppj.NotSupported, fmt.Sprintf("unsupported action %v on central system", action), nil)
			if err != nil {
				log.Printf("Unknown error %v while replying to message %v with CallError", err, requestId)
			}
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
		log.Printf("No handler for Call Result %v from charge point %v", requestId, chargePointId)
	}
}

func (cs *centralSystem) handleIncomingError(chargePointId string, err *ocpp.Error, details interface{}) {
	if callback, ok := cs.callbacks[chargePointId]; ok {
		delete(cs.callbacks, chargePointId)
		callback(nil, err)
	} else {
		log.Printf("No handler for Call Result %v from charge point %v", err.MessageId, chargePointId)
	}
}

func NewCentralSystem(dispatcher *ocppj.CentralSystem, server ws.WsServer) CentralSystem {
	if server == nil {
		server = ws.NewServer()
	}
	if dispatcher == nil {
		dispatcher = ocppj.NewCentralSystem(server, CoreProfile)
	}
	cs := centralSystem{
		centralSystem: dispatcher,
		callbacks:     map[string]func(confirmation ocpp.Confirmation, err error){}}
	cs.centralSystem.SetRequestHandler(cs.handleIncomingRequest)
	cs.centralSystem.SetConfirmationHandler(cs.handleIncomingConfirmation)
	cs.centralSystem.SetErrorHandler(cs.handleIncomingError)
	return &cs
}
