package ocpp16

import (
	"fmt"
	"github.com/lorenzodonini/go-ocpp/ocppj"
	"github.com/lorenzodonini/go-ocpp/ws"
	"log"
)

// -------------------- v1.6 Charge Point --------------------
type ChargePoint interface {
	// Message
	BootNotification(chargePointModel string, chargePointVendor string, props ...func(request *BootNotificationRequest)) (*BootNotificationConfirmation, *ocppj.CallError, error)
	Authorize(idTag string, props ...func(request *AuthorizeRequest)) (*AuthorizeConfirmation, *ocppj.CallError, error)
	//TODO: add missing profile methods

	// Logic
	SetChargePointCoreListener(listener ChargePointCoreListener)
	SendRequest(request ocppj.Request) (ocppj.Confirmation, *ocppj.CallError, error)
	SendRequestAsync(request ocppj.Request, callback func(confirmation ocppj.Confirmation, callError *ocppj.CallError)) error
	Start(centralSystemUrl string) error
}

type chargePoint struct {
	chargePoint          *ocppj.ChargePoint
	coreListener         ChargePointCoreListener
	confirmationListener chan ocppj.Confirmation
	errorListener        chan *ocppj.CallError
}

func (cp chargePoint) BootNotification(chargePointModel string, chargePointVendor string, props ...func(request *BootNotificationRequest)) (*BootNotificationConfirmation, *ocppj.CallError, error) {
	request := NewBootNotificationRequest(chargePointModel, chargePointVendor)
	for _, fn := range props {
		fn(request)
	}
	confirmation, callError, err := cp.SendRequest(request)
	return confirmation.(*BootNotificationConfirmation), callError, err
}

func (cp chargePoint) Authorize(idTag string, props ...func(request *AuthorizeRequest)) (*AuthorizeConfirmation, *ocppj.CallError, error) {
	request := NewAuthorizationRequest(idTag)
	for _, fn := range props {
		fn(request)
	}
	confirmation, callError, err := cp.SendRequest(request)
	return confirmation.(*AuthorizeConfirmation), callError, err
}

func (cp chargePoint) SetChargePointCoreListener(listener ChargePointCoreListener) {
	cp.coreListener = listener
}

func (cp chargePoint) SendRequest(request ocppj.Request) (ocppj.Confirmation, *ocppj.CallError, error) {
	err := cp.chargePoint.SendRequest(request)
	if err != nil {
		return nil, nil, err
	}
	select {
	case confirmation := <-cp.confirmationListener:
		return confirmation, nil, nil
	case callError := <-cp.errorListener:
		return nil, callError, nil
	}
}

func (cp chargePoint) SendRequestAsync(request ocppj.Request, callback func(confirmation ocppj.Confirmation, callError *ocppj.CallError)) error {
	err := cp.chargePoint.SendRequest(request)
	if err == nil {
		go func() {
			select {
			case confirmation := <-cp.confirmationListener:
				callback(confirmation, nil)
			case callError := <-cp.errorListener:
				callback(nil, callError)
			}
		}()
	}
	return err
}

func (cp chargePoint) sendResponse(call *ocppj.Call, confirmation ocppj.Confirmation, err error) {
	if confirmation != nil {
		callResult := ocppj.CallResult{MessageTypeId: ocppj.CALL_RESULT, UniqueId: call.UniqueId, Payload: confirmation}
		err := cp.chargePoint.SendMessage(&callResult)
		if err != nil {
			log.Printf("Unknown error %v while replying to message %v with CallError", err, call.UniqueId)
			//TODO: handle error somehow
		}
	} else {
		callError := cp.chargePoint.CreateCallError(call.UniqueId, ocppj.InternalError, "Couldn't generate valid confirmation", nil)
		err := cp.chargePoint.SendMessage(callError)
		if err != nil {
			log.Printf("Unknown error %v while replying to message %v with CallError", err, call.UniqueId)
		}
	}
}

func (cp chargePoint) Start(centralSystemUrl string) error {
	// TODO: implement auto-reconnect logic
	return cp.chargePoint.Start(centralSystemUrl)
}

func (cp chargePoint) handleIncomingCall(call *ocppj.Call) {
	if cp.coreListener == nil {
		log.Printf("Cannot handle call %v from central system. Sending CallError instead", call.UniqueId)
		callError := cp.chargePoint.CreateCallError(call.UniqueId, ocppj.NotImplemented, fmt.Sprintf("No handler for action %v implemented", call.Action), nil)
		err := cp.chargePoint.SendMessage(callError)
		if err != nil {
			log.Printf("Unknown error %v while replying to message %v with CallError", err, call.UniqueId)
		}
		return
	}
	var confirmation ocppj.Confirmation = nil
	var err error = nil
	switch call.Action {
	case ChangeAvailabilityFeatureName:
		confirmation, err = cp.coreListener.OnChangeAvailability(call.Payload.(*ChangeAvailabilityRequest))
	default:
		callError := cp.chargePoint.CreateCallError(call.UniqueId, ocppj.NotSupported, fmt.Sprintf("Unsupported action %v on charge point", call.Action), nil)
		err := cp.chargePoint.SendMessage(callError)
		if err != nil {
			log.Printf("Unknown error %v while replying to message %v with CallError", err, call.UniqueId)
		}
		return
	}
	cp.sendResponse(call, confirmation, err)
}

func NewChargePoint(id string) ChargePoint {
	cp := chargePoint{chargePoint: ocppj.NewChargePoint(id, ws.NewClient(), CoreProfile), confirmationListener: make(chan ocppj.Confirmation), errorListener: make(chan *ocppj.CallError)}
	cp.chargePoint.SetCallResultHandler(func(callResult *ocppj.CallResult) {
		cp.confirmationListener <- callResult.Payload
	})
	cp.chargePoint.SetCallErrorHandler(func(callError *ocppj.CallError) {
		cp.errorListener <- callError
	})
	cp.chargePoint.SetCallHandler(cp.handleIncomingCall)
	return cp
}

// -------------------- v1.6 Central System --------------------
type CentralSystem interface {
	// Message
	//TODO: add missing profile methods
	ChangeAvailability(clientId string, callback func(confirmation *ChangeAvailabilityConfirmation, callError *ocppj.CallError), connectorId int, availabilityType AvailabilityType, props ...func(request *ChangeAvailabilityRequest)) error
	// Logic
	SetCentralSystemCoreListener(listener CentralSystemCoreListener)
	SetNewChargePointHandler(handler func(chargePointId string))
	SendRequestAsync(clientId string, request ocppj.Request, callback func(confirmation ocppj.Confirmation, callError *ocppj.CallError)) error
	Start(listenPort int, listenPath string)
}

type centralSystem struct {
	centralSystem *ocppj.CentralSystem
	coreListener  CentralSystemCoreListener
	callbacks     map[string]func(confirmation ocppj.Confirmation, callError *ocppj.CallError)
}

func (cs centralSystem) ChangeAvailability(clientId string, callback func(confirmation *ChangeAvailabilityConfirmation, callError *ocppj.CallError), connectorId int, availabilityType AvailabilityType, props ...func(request *ChangeAvailabilityRequest)) error {
	request := NewChangeAvailabilityRequest(connectorId, availabilityType)
	genericCallback := func(confirmation ocppj.Confirmation, callError *ocppj.CallError) {
		if confirmation != nil {
			callback(confirmation.(*ChangeAvailabilityConfirmation), callError)
		} else {
			callback(nil, callError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs centralSystem) SetCentralSystemCoreListener(listener CentralSystemCoreListener) {
	cs.coreListener = listener
}

func (cs centralSystem) SetNewChargePointHandler(handler func(chargePointId string)) {
	cs.centralSystem.SetNewChargePointHandler(handler)
}

func (cs centralSystem) SendRequestAsync(clientId string, request ocppj.Request, callback func(confirmation ocppj.Confirmation, callError *ocppj.CallError)) error {
	err := cs.centralSystem.SendRequest(clientId, request)
	if err != nil {
		return err
	}
	cs.callbacks[clientId] = callback
	return nil
}

func (cs centralSystem) Start(listenPort int, listenPath string) {
	cs.centralSystem.Start(listenPort, listenPath)
}

func (cs centralSystem) sendResponse(chargePointId string, call *ocppj.Call, confirmation ocppj.Confirmation, err error) {
	if confirmation != nil {
		callResult := ocppj.CallResult{MessageTypeId: ocppj.CALL_RESULT, UniqueId: call.UniqueId, Payload: confirmation}
		err := cs.centralSystem.SendMessage(chargePointId, &callResult)
		if err != nil {
			//TODO: handle error somehow
			log.Print(err)
		}
	} else {
		callError := cs.centralSystem.CreateCallError(call.UniqueId, ocppj.InternalError, "Couldn't generate valid confirmation", nil)
		err := cs.centralSystem.SendMessage(chargePointId, callError)
		if err != nil {
			log.Printf("Unknown error %v while replying to message %v with CallError", err, call.UniqueId)
		}
	}
}

func (cs centralSystem) handleIncomingCall(chargePointId string, call *ocppj.Call) {
	if cs.coreListener == nil {
		log.Printf("Cannot handle call %v from charge point %v. Sending CallError instead", call.UniqueId, chargePointId)
		callError := cs.centralSystem.CreateCallError(call.UniqueId, ocppj.NotImplemented, fmt.Sprintf("No handler for action %v implemented", call.Action), nil)
		err := cs.centralSystem.SendMessage(chargePointId, callError)
		if err != nil {
			log.Printf("Unknown error %v while replying to message %v with CallError", err, call.UniqueId)
		}
	}
	var confirmation ocppj.Confirmation = nil
	var err error = nil
	// Execute in separate goroutine, so the caller goroutine is available
	go func() {
		switch call.Action {
		case BootNotificationFeatureName:
			confirmation, err = cs.coreListener.OnBootNotification(chargePointId, call.Payload.(*BootNotificationRequest))
			break
		case AuthorizeFeatureName:
			confirmation, err = cs.coreListener.OnAuthorize(chargePointId, call.Payload.(*AuthorizeRequest))
			break
		default:
			callError := cs.centralSystem.CreateCallError(call.UniqueId, ocppj.NotSupported, fmt.Sprintf("Unsupported action %v on central system", call.Action), nil)
			err := cs.centralSystem.SendMessage(chargePointId, callError)
			if err != nil {
				log.Printf("Unknown error %v while replying to message %v with CallError", err, call.UniqueId)
			}
			return
		}
		cs.sendResponse(chargePointId, call, confirmation, err)
	}()
}

func (cs centralSystem) handleIncomingCallResult(chargePointId string, callResult *ocppj.CallResult) {
	if callback, ok := cs.callbacks[chargePointId]; ok {
		delete(cs.callbacks, chargePointId)
		callback(callResult.Payload, nil)
	} else {
		log.Printf("No handler for Call Result %v from charge point %v", callResult.UniqueId, chargePointId)
	}
}

func (cs centralSystem) handleIncomingCallError(chargePointId string, callResult *ocppj.CallError) {
	if callback, ok := cs.callbacks[chargePointId]; ok {
		delete(cs.callbacks, chargePointId)
		callback(nil, callResult)
	} else {
		log.Printf("No handler for Call Result %v from charge point %v", callResult.UniqueId, chargePointId)
	}
}

func NewCentralSystem() CentralSystem {
	cs := centralSystem{
		centralSystem: ocppj.NewCentralSystem(ws.NewServer(), CoreProfile),
		callbacks:     map[string]func(confirmation ocppj.Confirmation, callError *ocppj.CallError){}}
	cs.centralSystem.SetCallHandler(cs.handleIncomingCall)
	cs.centralSystem.SetCallResultHandler(cs.handleIncomingCallResult)
	cs.centralSystem.SetCallErrorHandler(cs.handleIncomingCallError)
	return cs
}
