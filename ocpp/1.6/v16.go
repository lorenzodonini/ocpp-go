package v16

import (
	"github.com/lorenzodonini/go-ocpp/ocpp"
	"github.com/lorenzodonini/go-ocpp/ws"
)

// v1.6 Charge Point
type ChargePoint interface {
	BootNotification(chargePointModel string, chargePointVendor string) *BootNotificationRequest
	Authorize(idTag string) *AuthorizeRequest
	ChangeAvailability(connectorId int, availabilityType AvailabilityType) *ChangeAvailabilityRequest

	// Logic
	SendRequest(request ocpp.Request) (ocpp.Confirmation, *ocpp.CallError, error)
	SendRequestAsync(request ocpp.Request, callback func(confirmation ocpp.Confirmation, callError *ocpp.CallError)) error
}

type chargePoint struct {
	chargePoint *ocpp.ChargePoint
	confirmationListener chan ocpp.Confirmation
	errorListener chan *ocpp.CallError
}

func (cp chargePoint)BootNotification(chargePointModel string, chargePointVendor string) *BootNotificationRequest {
	return CoreProfile.CreateBootNotification(chargePointModel, chargePointVendor)
}

func (cp chargePoint)Authorize(idTag string) *AuthorizeRequest {
	return CoreProfile.CreateAuthorization(idTag)
}

func (cp chargePoint)ChangeAvailability(connectorId int, availabilityType AvailabilityType) *ChangeAvailabilityRequest {
	return CoreProfile.CreateChangeAvailability(connectorId, availabilityType)
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

func NewChargePoint(id string) ChargePoint {
	cp := chargePoint{chargePoint: ocpp.NewChargePoint(id, ws.NewClient(), CoreProfile.Profile), confirmationListener: make(chan ocpp.Confirmation), errorListener: make(chan *ocpp.CallError)}
	cp.chargePoint.SetCallResultHandler(func(callResult *ocpp.CallResult) {
		cp.confirmationListener <- callResult.Payload
	})
	cp.chargePoint.SetCallErrorHandler(func(callError *ocpp.CallError) {
		cp.errorListener <- callError
	})
	return cp
}

// v1.6 Central System

type CentralSystem interface {
	//TODO: add missing profile methods

	// Logic
	SendRequestAsync(clientId string, request ocpp.Request, callback func(confirmation ocpp.Confirmation, callError *ocpp.CallError)) error
}

type centralSystem struct {
	centralSystem *ocpp.CentralSystem
}

func (cs centralSystem)SendRequestAsync(clientId string, request ocpp.Request, callback func(confirmation ocpp.Confirmation, callError *ocpp.CallError)) error {
	return cs.centralSystem.SendRequest(clientId, request)
}

func NewCentralSystem() CentralSystem {
	cs := centralSystem{centralSystem: ocpp.NewCentralSystem(ws.NewServer(), CoreProfile.Profile)}
	//TODO: handle callbacks per client
	return cs
}
