package ocppj

import (
	"github.com/lorenzodonini/go-ocpp/ws"
	"github.com/pkg/errors"
	"log"
)

type CentralSystem struct {
	Endpoint
	server                ws.WsServer
	newChargePointHandler func(chargePointId string)
	requestHandler        func(chargePointId string, request Request, requestId string, action string)
	confirmationHandler   func(chargePointId string, confirmation Confirmation, requestId string)
	errorHandler          func(chargePointId string, errorCode ErrorCode, description string, details interface{}, requestId string)
	clientPendingMessages map[string]string
}

func NewCentralSystem(wsServer ws.WsServer, profiles ...*Profile) *CentralSystem {
	endpoint := Endpoint{pendingRequests: map[string]Request{}}
	for _, profile := range profiles {
		endpoint.AddProfile(profile)
	}
	if wsServer != nil {
		return &CentralSystem{Endpoint: endpoint, server: wsServer, clientPendingMessages: map[string]string{}}
	} else {
		return &CentralSystem{Endpoint: endpoint, server: &ws.Server{}, clientPendingMessages: map[string]string{}}
	}
}

func (centralSystem *CentralSystem) SetRequestHandler(handler func(chargePointId string, request Request, requestId string, action string)) {
	centralSystem.requestHandler = handler
}

func (centralSystem *CentralSystem) SetConfirmationHandler(handler func(chargePointId string, confirmation Confirmation, requestId string)) {
	centralSystem.confirmationHandler = handler
}

func (centralSystem *CentralSystem) SetErrorHandler(handler func(chargePointId string, errorCode ErrorCode, description string, details interface{}, requestId string)) {
	centralSystem.errorHandler = handler
}

func (centralSystem *CentralSystem) SetNewChargePointHandler(handler func(chargePointId string)) {
	centralSystem.newChargePointHandler = handler
}

func (centralSystem *CentralSystem) Start(listenPort int, listenPath string) {
	// Set internal message handler
	centralSystem.server.SetNewClientHandler(func(ws ws.Channel) {
		centralSystem.newChargePointHandler(ws.GetId())
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

func (centralSystem *CentralSystem) SendRequest(chargePointId string, request Request) error {
	err := validate.Struct(request)
	if err != nil {
		return err
	}
	req, ok := centralSystem.clientPendingMessages[chargePointId]
	if ok {
		// Cannot send. Protocol is based on response-confirmation
		return errors.Errorf("There already is a pending request %v for client %v. Cannot send a further one before receiving a confirmation first", req, chargePointId)
	}
	call, err := centralSystem.CreateCall(request.(Request))
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

func (centralSystem *CentralSystem) SendConfirmation(chargePointId string, requestId string, confirmation Confirmation) error {
	err := validate.Struct(confirmation)
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

func (centralSystem *CentralSystem) SendError(chargePointId string, requestId string, errorCode ErrorCode, description string, details interface{}) error {
	callError := centralSystem.CreateCallError(requestId, errorCode, description, details)
	err := validate.Struct(callError)
	if err != nil {
		return err
	}
	jsonMessage, err := callError.MarshalJSON()
	if err != nil {
		return err
	}
	return centralSystem.server.Write(chargePointId, []byte(jsonMessage))
}

//func (centralSystem *CentralSystem) SendMessage(chargePointId string, message Message) error {
//	err := validate.Struct(message)
//	if err != nil {
//		return err
//	}
//	jsonMessage, err := message.MarshalJSON()
//	if err != nil {
//		return err
//	}
//	if message.GetMessageTypeId() == CALL {
//		call := message.(*Call)
//		req, ok := centralSystem.clientPendingMessages[chargePointId]
//		if ok {
//			// Cannot send. Protocol is based on response-confirmation
//			return errors.Errorf("There already is a pending request %v. Cannot send a further one before receiving a confirmation first", req)
//		}
//		centralSystem.AddPendingRequest(message.GetUniqueId(), call.Payload)
//		centralSystem.clientPendingMessages[chargePointId] = call.UniqueId
//	}
//	err = centralSystem.server.Write(chargePointId, []byte(jsonMessage))
//	if err != nil && message.GetMessageTypeId() == CALL {
//		centralSystem.DeletePendingRequest(message.GetUniqueId())
//		delete(centralSystem.clientPendingMessages, chargePointId)
//	}
//	return err
//}

func (centralSystem *CentralSystem) ocppMessageHandler(wsChannel ws.Channel, data []byte) error {
	parsedJson := ParseRawJsonMessage(data)
	message, err := centralSystem.ParseMessage(parsedJson)
	if err != nil {
		if err.MessageId != "" {
			err2 := centralSystem.SendError(wsChannel.GetId(), err.MessageId, err.ErrorCode, err.Error.Error(), nil)
			if err2 != nil {
				return err2
			}
		}
		log.Print(err)
		return err.Error
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
		centralSystem.errorHandler(wsChannel.GetId(), callError.ErrorCode, callError.ErrorDescription, callError.ErrorDetails, callError.UniqueId)
	}
	return nil
}
