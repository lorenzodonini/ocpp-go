package ocpp

import (
	"container/list"
	"errors"
	"github.com/lorenzodonini/go-ocpp/ws"
	"log"
)

type CentralSystem struct {
	Endpoint
	server ws.WsServer
	callHandler func(chargePointId string, call *Call)
	callResultHandler func(chargePointId string, callResult *CallResult)
	callErrorHandler func(chargePointId string, callError *CallError)
	clientQueues map[string]*list.List
}

func NewCentralSystem(wsServer ws.WsServer, profiles ...*Profile) *CentralSystem {
	endpoint := Endpoint{PendingRequests: map[string]Request{}}
	for _, profile := range profiles {
		endpoint.AddProfile(profile)
	}
	if wsServer != nil {
		return &CentralSystem{Endpoint: endpoint, server: wsServer, clientQueues: map[string]*list.List{}}
	} else {
		return &CentralSystem{Endpoint: endpoint, server: &ws.Server{}, clientQueues: map[string]*list.List{}}
	}
}

func (centralSystem *CentralSystem)SeCallHandler(handler func(chargePointId string, call *Call)) {
	centralSystem.callHandler = handler
}

func (centralSystem *CentralSystem)SeCallResultHandler(handler func(chargePointId string, callResult *CallResult)) {
	centralSystem.callResultHandler = handler
}

func (centralSystem *CentralSystem)SeCalleHandler(handler func(chargePointId string, callError *CallError)) {
	centralSystem.callErrorHandler = handler
}

func (centralSystem *CentralSystem)Start(listenPort int, listenPath string) {
	// Set internal message handler
	centralSystem.server.SetMessageHandler(centralSystem.ocppMessageHandler)
	// Serve & run
	// TODO: return error?
	centralSystem.server.Start(listenPort, listenPath)
}

func (centralSystem *CentralSystem)Stop() {
	centralSystem.server.Stop()
}

func (centralSystem *CentralSystem)SendRequest(chargePointId string, request Request) error {
	centralSystem.clientQueues[chargePointId].PushBack(request)
	if centralSystem.clientQueues[chargePointId].Len() > 1 {
		// Cannot send right away
		return nil
	}
	err := centralSystem.processCallQueue(chargePointId)
	if err != nil {
		return err
	}
	//TODO: use promise/future for fetching the result
	return nil
}

func (centralSystem *CentralSystem)SendMessage(chargePointId string, ocppMessage interface{}) error {
	message, ok := ocppMessage.(*Message)
	if !ok {
		return errors.New("invalid ocpp message. Couldn't send")
	}
	jsonMessage, err := message.ToJson()
	if err != nil {
		return err
	}
	if message.MessageTypeId == CALL {
		call := ocppMessage.(*Call)
		centralSystem.PendingRequests[message.UniqueId] = call.Payload
	}
	err = centralSystem.server.Write(chargePointId, []byte(jsonMessage))
	if err != nil {
		return err
	}
	//TODO: use promise/future for fetching the result
	return nil
}

func (centralSystem *CentralSystem)processCallQueue(chargePointId string) error {
	if centralSystem.clientQueues[chargePointId].Len() == 0 {
		return nil
	}
	element := centralSystem.clientQueues[chargePointId].Front()
	request := element.Value
	call, err := centralSystem.CreateCall(request.(Request))
	jsonMessage, err := call.ToJson()
	if err != nil {
		return err
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
	ocppMessage, ok := message.(Message)
	if !ok {
		return errors.New("couldn't convert parsed data to Message type")
	}
	switch ocppMessage.MessageTypeId {
	case CALL:
		call := message.(Call)
		centralSystem.callHandler(wsChannel.GetId(), &call)
	case CALL_RESULT:
		callResult := message.(CallResult)
		centralSystem.callResultHandler(wsChannel.GetId(), &callResult)
		err := centralSystem.processCallQueue(wsChannel.GetId())
		if err != nil {
			return err
		}
	case CALL_ERROR:
		callError := message.(CallError)
		centralSystem.callErrorHandler(wsChannel.GetId(), &callError)
		err := centralSystem.processCallQueue(wsChannel.GetId())
		if err != nil {
			return err
		}
	}
	return nil
}


