package ocppj

import (
	"github.com/lorenzodonini/go-ocpp/ocpp"
	"github.com/lorenzodonini/go-ocpp/ws"
	"github.com/pkg/errors"
	"log"
)

type CentralSystem struct {
	Endpoint
	server                         ws.WsServer
	newChargePointHandler          func(chargePointId string)
	disconnectedChargePointHandler func(chargePointId string)
	requestHandler                 func(chargePointId string, request ocpp.Request, requestId string, action string)
	confirmationHandler            func(chargePointId string, confirmation ocpp.Confirmation, requestId string)
	errorHandler                   func(chargePointId string, err *ocpp.Error, details interface{})
	clientPendingMessages          map[string]string
}

func NewCentralSystem(wsServer ws.WsServer, profiles ...*ocpp.Profile) *CentralSystem {
	endpoint := Endpoint{pendingRequests: map[string]ocpp.Request{}}
	for _, profile := range profiles {
		endpoint.AddProfile(profile)
	}
	if wsServer != nil {
		return &CentralSystem{Endpoint: endpoint, server: wsServer, clientPendingMessages: map[string]string{}}
	} else {
		return &CentralSystem{Endpoint: endpoint, server: &ws.Server{}, clientPendingMessages: map[string]string{}}
	}
}

func (centralSystem *CentralSystem) SetRequestHandler(handler func(chargePointId string, request ocpp.Request, requestId string, action string)) {
	centralSystem.requestHandler = handler
}

func (centralSystem *CentralSystem) SetConfirmationHandler(handler func(chargePointId string, confirmation ocpp.Confirmation, requestId string)) {
	centralSystem.confirmationHandler = handler
}

func (centralSystem *CentralSystem) SetErrorHandler(handler func(chargePointId string, err *ocpp.Error, details interface{})) {
	centralSystem.errorHandler = handler
}

func (centralSystem *CentralSystem) SetNewChargePointHandler(handler func(chargePointId string)) {
	centralSystem.newChargePointHandler = handler
}

func (centralSystem *CentralSystem) SetDisconnectedChargePointHandler(handler func(chargePointId string)) {
	centralSystem.disconnectedChargePointHandler = handler
}

func (centralSystem *CentralSystem) Start(listenPort int, listenPath string) {
	// Set internal message handler
	centralSystem.server.SetNewClientHandler(func(ws ws.Channel) {
		if centralSystem.newChargePointHandler != nil {
			centralSystem.newChargePointHandler(ws.GetId())
		}
	})
	centralSystem.server.SetDisconnectedClientHandler(func(ws ws.Channel) {
		if centralSystem.disconnectedChargePointHandler != nil {
			centralSystem.disconnectedChargePointHandler(ws.GetId())
		}
	})
	centralSystem.server.SetMessageHandler(centralSystem.ocppMessageHandler)
	// Serve & run
	// TODO: return error?
	centralSystem.server.Start(listenPort, listenPath)
}

func (centralSystem *CentralSystem) Stop() {
	centralSystem.server.Stop()
	centralSystem.clearPendingRequests()
}

func (centralSystem *CentralSystem) SendRequest(chargePointId string, request ocpp.Request) error {
	err := Validate.Struct(request)
	if err != nil {
		return err
	}
	req, ok := centralSystem.clientPendingMessages[chargePointId]
	if ok {
		// Cannot send. Protocol is based on response-confirmation
		return errors.Errorf("There already is a pending request %v for client %v. Cannot send a further one before receiving a confirmation first", req, chargePointId)
	}
	call, err := centralSystem.CreateCall(request.(ocpp.Request))
	if err != nil {
		return err
	}
	jsonMessage, err := call.MarshalJSON()
	if err != nil {
		return err
	}
	centralSystem.clientPendingMessages[chargePointId] = call.UniqueId
	err = centralSystem.server.Write(chargePointId, []byte(jsonMessage))
	if err != nil {
		// Clear pending request
		centralSystem.DeletePendingRequest(call.GetUniqueId())
		delete(centralSystem.clientPendingMessages, chargePointId)
	}
	return err
}

func (centralSystem *CentralSystem) SendConfirmation(chargePointId string, requestId string, confirmation ocpp.Confirmation) error {
	err := Validate.Struct(confirmation)
	if err != nil {
		return err
	}
	callResult, err := centralSystem.CreateCallResult(confirmation, requestId)
	if err != nil {
		return err
	}
	jsonMessage, err := callResult.MarshalJSON()
	if err != nil {
		return err
	}
	return centralSystem.server.Write(chargePointId, []byte(jsonMessage))
}

func (centralSystem *CentralSystem) SendError(chargePointId string, requestId string, errorCode ocpp.ErrorCode, description string, details interface{}) error {
	callError := centralSystem.CreateCallError(requestId, errorCode, description, details)
	err := Validate.Struct(callError)
	if err != nil {
		return err
	}
	jsonMessage, err := callError.MarshalJSON()
	if err != nil {
		return err
	}
	return centralSystem.server.Write(chargePointId, []byte(jsonMessage))
}

func (centralSystem *CentralSystem) ocppMessageHandler(wsChannel ws.Channel, data []byte) error {
	parsedJson := ParseRawJsonMessage(data)
	message, err := centralSystem.ParseMessage(parsedJson)
	if err != nil {
		if err.MessageId != "" {
			err2 := centralSystem.SendError(wsChannel.GetId(), err.MessageId, err.Code, err.Description, nil)
			if err2 != nil {
				return err2
			}
		}
		log.Print(err)
		return err
	}
	switch message.GetMessageTypeId() {
	case CALL:
		call := message.(*Call)
		centralSystem.requestHandler(wsChannel.GetId(), call.Payload, call.UniqueId, call.Action)
	case CALL_RESULT:
		callResult := message.(*CallResult)
		delete(centralSystem.clientPendingMessages, wsChannel.GetId())
		centralSystem.confirmationHandler(wsChannel.GetId(), callResult.Payload, callResult.UniqueId)
	case CALL_ERROR:
		callError := message.(*CallError)
		delete(centralSystem.clientPendingMessages, wsChannel.GetId())
		centralSystem.errorHandler(wsChannel.GetId(), ocpp.NewError(callError.ErrorCode, callError.ErrorDescription, callError.UniqueId), callError.ErrorDetails)
	}
	return nil
}
