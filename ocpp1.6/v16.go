package ocpp16

import (
	"errors"
	"fmt"
	"github.com/lorenzodonini/go-ocpp/ocppj"
	"github.com/lorenzodonini/go-ocpp/ws"
	"log"
)

// -------------------- v1.6 Charge Point --------------------
type ChargePoint interface {
	// Message
	BootNotification(chargePointModel string, chargePointVendor string, props ...func(request *BootNotificationRequest)) (*BootNotificationConfirmation, *ocppj.ProtoError, error)
	Authorize(idTag string, props ...func(request *AuthorizeRequest)) (*AuthorizeConfirmation, *ocppj.ProtoError, error)
	//TODO: add missing profile methods

	// Logic
	SetChargePointCoreListener(listener ChargePointCoreListener)
	SendRequest(request ocppj.Request) (ocppj.Confirmation, *ocppj.ProtoError, error)
	SendRequestAsync(request ocppj.Request, callback func(confirmation ocppj.Confirmation, callError *ocppj.ProtoError)) error
	Start(centralSystemUrl string) error
}

type chargePoint struct {
	chargePoint          *ocppj.ChargePoint
	coreListener         ChargePointCoreListener
	confirmationListener chan ocppj.Confirmation
	errorListener        chan ocppj.ProtoError
}

func (cp *chargePoint) BootNotification(chargePointModel string, chargePointVendor string, props ...func(request *BootNotificationRequest)) (*BootNotificationConfirmation, *ocppj.ProtoError, error) {
	request := NewBootNotificationRequest(chargePointModel, chargePointVendor)
	for _, fn := range props {
		fn(request)
	}
	confirmation, protoError, err := cp.SendRequest(request)
	return confirmation.(*BootNotificationConfirmation), protoError, err
}

func (cp *chargePoint) Authorize(idTag string, props ...func(request *AuthorizeRequest)) (*AuthorizeConfirmation, *ocppj.ProtoError, error) {
	request := NewAuthorizationRequest(idTag)
	for _, fn := range props {
		fn(request)
	}
	confirmation, protoError, err := cp.SendRequest(request)
	return confirmation.(*AuthorizeConfirmation), protoError, err
}

func (cp *chargePoint) SetChargePointCoreListener(listener ChargePointCoreListener) {
	cp.coreListener = listener
}

func (cp *chargePoint) SendRequest(request ocppj.Request) (ocppj.Confirmation, *ocppj.ProtoError, error) {
	err := cp.chargePoint.SendRequest(request)
	if err != nil {
		return nil, nil, err
	}
	select {
	case confirmation := <-cp.confirmationListener:
		return confirmation, nil, nil
	case protoError := <-cp.errorListener:
		return nil, &protoError, nil
	}
}

func (cp *chargePoint) SendRequestAsync(request ocppj.Request, callback func(confirmation ocppj.Confirmation, protoError *ocppj.ProtoError)) error {
	err := cp.chargePoint.SendRequest(request)
	if err == nil {
		// Retrieve result asynchronously
		go func() {
			select {
			case confirmation := <-cp.confirmationListener:
				callback(confirmation, nil)
			case protoError := <-cp.errorListener:
				callback(nil, &protoError)
			}
		}()
	}
	return err
}

func (cp *chargePoint) sendResponse(confirmation ocppj.Confirmation, err error, requestId string) {
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

func (cp *chargePoint) handleIncomingRequest(request ocppj.Request, requestId string, action string) {
	if cp.coreListener == nil {
		log.Printf("Cannot handle call %v from central system. Sending CallError instead", requestId)
		err := cp.chargePoint.SendError(requestId, ocppj.NotImplemented, fmt.Sprintf("No handler for action %v implemented", action), nil)
		if err != nil {
			log.Printf("Unknown error %v while replying to message %v with CallError", err, requestId)
		}
		return
	}
	var confirmation ocppj.Confirmation = nil
	var err error = nil
	switch action {
	case ChangeAvailabilityFeatureName:
		confirmation, err = cp.coreListener.OnChangeAvailability(request.(*ChangeAvailabilityRequest))
	default:
		err := cp.chargePoint.SendError(requestId, ocppj.NotSupported, fmt.Sprintf("Unsupported action %v on charge point", action), nil)
		if err != nil {
			log.Printf("Unknown error %v while replying to message %v with CallError", err, requestId)
		}
		return
	}
	cp.sendResponse(confirmation, err, requestId)
}

func NewChargePoint(id string, client ws.WsClient) ChargePoint {
	if client == nil {
		client = ws.NewClient()
	}
	cp := chargePoint{chargePoint: ocppj.NewChargePoint(id, client, CoreProfile), confirmationListener: make(chan ocppj.Confirmation), errorListener: make(chan ocppj.ProtoError)}
	cp.chargePoint.SetConfirmationHandler(func(confirmation ocppj.Confirmation, requestId string) {
		cp.confirmationListener <- confirmation
	})
	cp.chargePoint.SetErrorHandler(func(errorCode ocppj.ErrorCode, description string, details interface{}, requestId string) {
		protoError := ocppj.ProtoError{Error: errors.New(description), ErrorCode: errorCode, MessageId: requestId}
		cp.errorListener <- protoError
	})
	cp.chargePoint.SetRequestHandler(cp.handleIncomingRequest)
	return &cp
}

// -------------------- v1.6 Central System --------------------
type CentralSystem interface {
	// Message
	//TODO: add missing profile methods
	ChangeAvailability(clientId string, callback func(confirmation *ChangeAvailabilityConfirmation, callError *ocppj.ProtoError), connectorId int, availabilityType AvailabilityType, props ...func(request *ChangeAvailabilityRequest)) error
	// Logic
	SetCentralSystemCoreListener(listener CentralSystemCoreListener)
	SetNewChargePointHandler(handler func(chargePointId string))
	SendRequestAsync(clientId string, request ocppj.Request, callback func(confirmation ocppj.Confirmation, callError *ocppj.ProtoError)) error
	Start(listenPort int, listenPath string)
}

type centralSystem struct {
	centralSystem *ocppj.CentralSystem
	coreListener  CentralSystemCoreListener
	callbacks     map[string]func(confirmation ocppj.Confirmation, callError *ocppj.ProtoError)
}

func (cs *centralSystem) ChangeAvailability(clientId string, callback func(confirmation *ChangeAvailabilityConfirmation, protoError *ocppj.ProtoError), connectorId int, availabilityType AvailabilityType, props ...func(request *ChangeAvailabilityRequest)) error {
	request := NewChangeAvailabilityRequest(connectorId, availabilityType)
	genericCallback := func(confirmation ocppj.Confirmation, protoError *ocppj.ProtoError) {
		if confirmation != nil {
			callback(confirmation.(*ChangeAvailabilityConfirmation), protoError)
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

func (cs *centralSystem) SendRequestAsync(clientId string, request ocppj.Request, callback func(confirmation ocppj.Confirmation, protoError *ocppj.ProtoError)) error {
	err := cs.centralSystem.SendRequest(clientId, request)
	if err != nil {
		return err
	}
	cs.callbacks[clientId] = callback
	return nil
}

func (cs *centralSystem) Start(listenPort int, listenPath string) {
	cs.centralSystem.Start(listenPort, listenPath)
}

func (cs *centralSystem) sendResponse(chargePointId string, confirmation ocppj.Confirmation, err error, requestId string) {
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

func (cs *centralSystem) handleIncomingRequest(chargePointId string, request ocppj.Request, requestId string, action string) {
	if cs.coreListener == nil {
		log.Printf("Cannot handle call %v from charge point %v. Sending CallError instead", requestId, chargePointId)
		err := cs.centralSystem.SendError(chargePointId, requestId, ocppj.NotImplemented, fmt.Sprintf("No handler for action %v implemented", action), nil)
		if err != nil {
			log.Printf("Unknown error %v while replying to message %v with CallError", err, requestId)
		}
	}
	var confirmation ocppj.Confirmation = nil
	var err error = nil
	// Execute in separate goroutine, so the caller goroutine is available
	go func() {
		switch action {
		case BootNotificationFeatureName:
			confirmation, err = cs.coreListener.OnBootNotification(chargePointId, request.(*BootNotificationRequest))
			break
		case AuthorizeFeatureName:
			confirmation, err = cs.coreListener.OnAuthorize(chargePointId, request.(*AuthorizeRequest))
			break
		default:
			err := cs.centralSystem.SendError(chargePointId, requestId, ocppj.NotSupported, fmt.Sprintf("Unsupported action %v on central system", action), nil)
			if err != nil {
				log.Printf("Unknown error %v while replying to message %v with CallError", err, requestId)
			}
			return
		}
		cs.sendResponse(chargePointId, confirmation, err, requestId)
	}()
}

func (cs *centralSystem) handleIncomingConfirmation(chargePointId string, confirmation ocppj.Confirmation, requestId string) {
	if callback, ok := cs.callbacks[chargePointId]; ok {
		delete(cs.callbacks, chargePointId)
		callback(confirmation, nil)
	} else {
		log.Printf("No handler for Call Result %v from charge point %v", requestId, chargePointId)
	}
}

func (cs *centralSystem) handleIncomingError(chargePointId string, errorCode ocppj.ErrorCode, description string, details interface{}, requestId string) {
	if callback, ok := cs.callbacks[chargePointId]; ok {
		delete(cs.callbacks, chargePointId)
		protoError := ocppj.ProtoError{Error: errors.New(description), ErrorCode: errorCode, MessageId: requestId}
		callback(nil, &protoError)
	} else {
		log.Printf("No handler for Call Result %v from charge point %v", requestId, chargePointId)
	}
}

func NewCentralSystem(server ws.WsServer) CentralSystem {
	if server == nil {
		server = ws.NewServer()
	}
	cs := centralSystem{
		centralSystem: ocppj.NewCentralSystem(server, CoreProfile),
		callbacks:     map[string]func(confirmation ocppj.Confirmation, callError *ocppj.ProtoError){}}
	cs.centralSystem.SetRequestHandler(cs.handleIncomingRequest)
	cs.centralSystem.SetConfirmationHandler(cs.handleIncomingConfirmation)
	cs.centralSystem.SetErrorHandler(cs.handleIncomingError)
	return &cs
}
