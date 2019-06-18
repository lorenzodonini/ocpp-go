package v16

import (
	"github.com/lorenzodonini/go-ocpp/ocpp"
	"github.com/lorenzodonini/go-ocpp/ws"
	"log"
)

// -------------------- v1.6 Charge Point --------------------
type ChargePoint interface {
	// Message
	BootNotification(chargePointModel string, chargePointVendor string, props... func(request *BootNotificationRequest)) (*BootNotificationConfirmation, *ocpp.CallError, error)
	Authorize(idTag string, props... func(request *AuthorizeRequest)) (*AuthorizeConfirmation, *ocpp.CallError, error)
	//TODO: add missing profile methods

	// Logic
	SetChargePointCoreListener(listener ChargePointCoreListener)
	SendRequest(request ocpp.Request) (ocpp.Confirmation, *ocpp.CallError, error)
	SendRequestAsync(request ocpp.Request, callback func(confirmation ocpp.Confirmation, callError *ocpp.CallError)) error
	Start(centralSystemUrl string) error
}

type chargePoint struct {
	chargePoint *ocpp.ChargePoint
	coreListener ChargePointCoreListener
	confirmationListener chan ocpp.Confirmation
	errorListener chan *ocpp.CallError
}

func (cp chargePoint)BootNotification(chargePointModel string, chargePointVendor string, props... func(request *BootNotificationRequest)) (*BootNotificationConfirmation, *ocpp.CallError, error) {
	request := CoreProfile.CreateBootNotification(chargePointModel, chargePointVendor)
	for _, fn := range props {
		fn(request)
	}
	confirmation, callError, err := cp.SendRequest(request)
	return confirmation.(*BootNotificationConfirmation), callError, err
}


func (cp chargePoint)Authorize(idTag string, props... func(request *AuthorizeRequest)) (*AuthorizeConfirmation, *ocpp.CallError, error) {
	request := CoreProfile.CreateAuthorization(idTag)
	for _, fn := range props {
		fn(request)
	}
	confirmation, callError, err := cp.SendRequest(request)
	return confirmation.(*AuthorizeConfirmation), callError, err
}

func (cp chargePoint)SetChargePointCoreListener(listener ChargePointCoreListener) {
	cp.coreListener = listener
}

func (cp chargePoint)SendRequest(request ocpp.Request) (ocpp.Confirmation, *ocpp.CallError, error) {
	err := cp.chargePoint.SendRequest(request)
	if err != nil {
		return nil, nil, err
	}
	select {
		case confirmation := <- cp.confirmationListener:
			return confirmation, nil, nil
		case callError := <- cp.errorListener:
			return nil, callError, nil
	}
}

func (cp chargePoint)SendRequestAsync(request ocpp.Request, callback func(confirmation ocpp.Confirmation, callError *ocpp.CallError)) error {
	err := cp.chargePoint.SendRequest(request)
	if err == nil {
		go func() {
			select {
			case confirmation := <- cp.confirmationListener:
				callback(confirmation, nil)
			case callError := <- cp.errorListener:
				callback(nil, callError)
			}
		}()
	}
	return err
}

func (cp chargePoint)sendResponse(call *ocpp.Call, confirmation ocpp.Confirmation, err error) {
	if confirmation != nil {
		callResult := ocpp.CallResult{ MessageTypeId: ocpp.CALL_RESULT, UniqueId: call.UniqueId, Payload: confirmation }
		err := cp.chargePoint.SendMessage(&callResult)
		if err != nil {
			//TODO: handle error somehow
			log.Print(err)
		}
	} else {
		//TODO: create call error
		return
	}
}

func (cp chargePoint)Start(centralSystemUrl string) error {
	// TODO: implement auto-reconnect logic
	return cp.chargePoint.Start(centralSystemUrl)
}

func (cp chargePoint)handleIncomingCall(call *ocpp.Call) {
	if cp.coreListener == nil {
		log.Printf("Cannot handle call %v from central system. Sending CallError instead", call.UniqueId)
		//TODO: send call error
		return
	}
	var confirmation ocpp.Confirmation = nil
	var err error = nil
	switch call.Action {
	case ChangeAvailabilityFeatureName:
		confirmation, err = cp.coreListener.OnChangeAvailability(call.Payload.(*ChangeAvailabilityRequest))
	default:
		log.Printf("Unsupported action %v on charge point", call.Action)
		//TODO: send back CallError
	}
	cp.sendResponse(call, confirmation, err)
}

func NewChargePoint(id string) ChargePoint {
	cp := chargePoint{chargePoint: ocpp.NewChargePoint(id, ws.NewClient(), CoreProfile.Profile), confirmationListener: make(chan ocpp.Confirmation), errorListener: make(chan *ocpp.CallError)}
	cp.chargePoint.SetCallResultHandler(func(callResult *ocpp.CallResult) {
		cp.confirmationListener <- callResult.Payload
	})
	cp.chargePoint.SetCallErrorHandler(func(callError *ocpp.CallError) {
		cp.errorListener <- callError
	})
	cp.chargePoint.SetCallHandler(cp.handleIncomingCall)
	return cp
}

// -------------------- v1.6 Central System --------------------
type CentralSystem interface {
	// Message
	//TODO: add missing profile methods
	ChangeAvailability(clientId string, callback func(confirmation *ChangeAvailabilityConfirmation, callError *ocpp.CallError), connectorId int, availabilityType AvailabilityType, props... func(request *ChangeAvailabilityRequest)) error
	// Logic
	SetCentralSystemCoreListener(listener CentralSystemCoreListener)
	SetNewChargePointHandler(handler func(chargePointId string))
	SendRequestAsync(clientId string, request ocpp.Request, callback func(confirmation ocpp.Confirmation, callError *ocpp.CallError)) error
	Start(listenPort int, listenPath string)
}

type centralSystem struct {
	centralSystem *ocpp.CentralSystem
	coreListener CentralSystemCoreListener
	callbacks map[string]func(confirmation ocpp.Confirmation, callError *ocpp.CallError)
}

func (cs centralSystem)ChangeAvailability(clientId string, callback func(confirmation *ChangeAvailabilityConfirmation, callError *ocpp.CallError), connectorId int, availabilityType AvailabilityType, props... func(request *ChangeAvailabilityRequest)) error {
	request := CoreProfile.CreateChangeAvailability(connectorId, availabilityType)
	genericCallback := func(confirmation ocpp.Confirmation, callError *ocpp.CallError) {
		if confirmation != nil {
			callback(confirmation.(*ChangeAvailabilityConfirmation), callError)
		} else {
			callback(nil, callError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs centralSystem)SetCentralSystemCoreListener(listener CentralSystemCoreListener) {
	cs.coreListener = listener
}

func (cs centralSystem)SetNewChargePointHandler(handler func(chargePointId string)) {
	cs.centralSystem.SetNewChargePointHandler(handler)
}

func (cs centralSystem)SendRequestAsync(clientId string, request ocpp.Request, callback func(confirmation ocpp.Confirmation, callError *ocpp.CallError)) error {
	err := cs.centralSystem.SendRequest(clientId, request)
	if err != nil {
		return err
	}
	cs.callbacks[clientId] = callback
	return nil
}

func (cs centralSystem)Start(listenPort int, listenPath string) {
	cs.centralSystem.Start(listenPort, listenPath)
}

func (cs centralSystem)sendResponse(chargePointId string, call *ocpp.Call, confirmation ocpp.Confirmation, err error) {
	if confirmation != nil {
		callResult := ocpp.CallResult{ MessageTypeId: ocpp.CALL_RESULT, UniqueId: call.UniqueId, Payload: confirmation }
		err := cs.centralSystem.SendMessage(chargePointId, &callResult)
		if err != nil {
			//TODO: handle error somehow
			log.Print(err)
		}
	} else {
		//TODO: create call error
		return
	}
}

func (cs centralSystem)handleIncomingCall(chargePointId string, call *ocpp.Call) {
	if cs.coreListener == nil {
		log.Printf("Cannot handle call %v from charge point %v. Sending CallError instead", call.UniqueId, chargePointId)
		//TODO: send call error
		return
	}
	var confirmation ocpp.Confirmation = nil
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
			log.Printf("Unsupported action %v on central system", call.Action)
			//TODO: send back CallError
		}
		cs.sendResponse(chargePointId, call, confirmation, err)
	}()
}

func (cs centralSystem)handleIncomingCallResult(chargePointId string, callResult *ocpp.CallResult) {
	if callback, ok := cs.callbacks[chargePointId]; ok {
		delete(cs.callbacks, chargePointId)
		callback(callResult.Payload, nil)
	} else {
		log.Printf("No handler for Call Result %v from charge point %v", callResult.UniqueId, chargePointId)
	}
}

func (cs centralSystem)handleIncomingCallError(chargePointId string, callResult *ocpp.CallError) {
	if callback, ok := cs.callbacks[chargePointId]; ok {
		delete(cs.callbacks, chargePointId)
		callback(nil, callResult)
	} else {
		log.Printf("No handler for Call Result %v from charge point %v", callResult.UniqueId, chargePointId)
	}
}

func NewCentralSystem() CentralSystem {
	cs := centralSystem{
		centralSystem: ocpp.NewCentralSystem(ws.NewServer(), CoreProfile.Profile),
		callbacks: map[string]func(confirmation ocpp.Confirmation, callError *ocpp.CallError){}}
	cs.centralSystem.SetCallHandler(cs.handleIncomingCall)
	cs.centralSystem.SetCallResultHandler(cs.handleIncomingCallResult)
	cs.centralSystem.SetCallErrorHandler(cs.handleIncomingCallError)
	return cs
}
