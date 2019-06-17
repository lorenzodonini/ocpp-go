package ocpp

import (
	"github.com/lorenzodonini/go-ocpp/ws"
	"github.com/pkg/errors"
	"log"
)

type CentralSystem struct {
	Endpoint
	server ws.WsServer
	newChargePointHandler func(chargePointId string)
	callHandler func(chargePointId string, call *Call)
	callResultHandler func(chargePointId string, callResult *CallResult)
	callErrorHandler func(chargePointId string, callError *CallError)
	clientPendingMessages map[string]string
}

func NewCentralSystem(wsServer ws.WsServer, profiles ...*Profile) *CentralSystem {
	endpoint := Endpoint{PendingRequests: map[string]Request{}}
	for _, profile := range profiles {
		endpoint.AddProfile(profile)
	}
	if wsServer != nil {
		return &CentralSystem{Endpoint: endpoint, server: wsServer, clientPendingMessages: map[string]string{}}
	} else {
		return &CentralSystem{Endpoint: endpoint, server: &ws.Server{}, clientPendingMessages: map[string]string{}}
	}
}

func (centralSystem *CentralSystem)SetCallHandler(handler func(chargePointId string, call *Call)) {
	centralSystem.callHandler = handler
}

func (centralSystem *CentralSystem)SetCallResultHandler(handler func(chargePointId string, callResult *CallResult)) {
	centralSystem.callResultHandler = handler
}

func (centralSystem *CentralSystem)SetCallErrorHandler(handler func(chargePointId string, callError *CallError)) {
	centralSystem.callErrorHandler = handler
}

func (centralSystem *CentralSystem)SetNewChargePointHandler(handler func(chargePointId string)) {
	centralSystem.newChargePointHandler = handler
}

func (centralSystem *CentralSystem)Start(listenPort int, listenPath string) {
	// Set internal message handler
	centralSystem.server.SetNewClientHandler(func(ws ws.Channel) {
		centralSystem.newChargePointHandler(ws.GetId())
	})
	centralSystem.server.SetMessageHandler(centralSystem.ocppMessageHandler)
	// Serve & run
	// TODO: return error?
	centralSystem.server.Start(listenPort, listenPath)
}

func (centralSystem *CentralSystem)Stop() {
	centralSystem.server.Stop()
}

func (centralSystem *CentralSystem)SendRequest(chargePointId string, request Request) error {
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
	//TODO: use promise/future for fetching the result
	return centralSystem.server.Write(chargePointId, []byte(jsonMessage))
}

func (centralSystem *CentralSystem)SendMessage(chargePointId string, message Message) error {
	err := validate.Struct(message)
	if err != nil {
		return err
	}
	jsonMessage, err := message.MarshalJSON()
	if err != nil {
		return err
	}
	if message.GetMessageTypeId() == CALL {
		call := message.(*Call)
		req, ok := centralSystem.clientPendingMessages[chargePointId]
		if ok {
			// Cannot send. Protocol is based on response-confirmation
			return errors.Errorf("There already is a pending request %v. Cannot send a further one before receiving a confirmation first", req)
		}
		centralSystem.PendingRequests[message.GetUniqueId()] = call.Payload
		centralSystem.clientPendingMessages[chargePointId] = call.UniqueId
	}
	err = centralSystem.server.Write(chargePointId, []byte(jsonMessage))
	if err != nil {
		return err
	}
	//TODO: use promise/future for fetching the result
	return nil
}

func (centralSystem *CentralSystem)ocppMessageHandler(wsChannel ws.Channel, data []byte) error {
	parsedJson := ParseRawJsonMessage(data)
	message, err := centralSystem.ParseMessage(parsedJson)
	if err != nil {
		// TODO: handle
		log.Printf("Error while parsing message: %v", err)
		return err
	}
	switch message.GetMessageTypeId() {
	case CALL:
		call := message.(*Call)
		centralSystem.callHandler(wsChannel.GetId(), call)
	case CALL_RESULT:
		callResult := message.(*CallResult)
		delete(centralSystem.clientPendingMessages, wsChannel.GetId())
		centralSystem.callResultHandler(wsChannel.GetId(), callResult)
	case CALL_ERROR:
		callError := message.(*CallError)
		delete(centralSystem.clientPendingMessages, wsChannel.GetId())
		centralSystem.callErrorHandler(wsChannel.GetId(), callError)
	}
	return nil
}


