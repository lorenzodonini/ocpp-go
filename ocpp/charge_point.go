package ocpp

import (
	"fmt"
	"github.com/lorenzodonini/go-ocpp/ws"
	"github.com/pkg/errors"
	"log"
)

type ChargePoint struct {
	Endpoint
	client ws.WsClient
	Id string
	callHandler func(call *Call)
	callResultHandler func(callResult *CallResult)
	callErrorHandler func(callError *CallError)
	pendingRequest string
}

func NewChargePoint(id string, wsClient ws.WsClient, profiles ...*Profile) *ChargePoint {
	endpoint := Endpoint{PendingRequests: map[string]Request{}}
	for _, profile := range profiles {
		endpoint.AddProfile(profile)
	}
	if wsClient != nil {
		return &ChargePoint{Endpoint: endpoint, client: wsClient, Id: id, pendingRequest: ""}
	} else {
		return &ChargePoint{Endpoint: endpoint, client: &ws.Client{}, Id: id, pendingRequest: ""}
	}
}

func (chargePoint *ChargePoint)SetCallHandler(handler func(call *Call)) {
	chargePoint.callHandler = handler
}

func (chargePoint *ChargePoint)SetCallResultHandler(handler func(callResult *CallResult)) {
	chargePoint.callResultHandler = handler
}

func (chargePoint *ChargePoint)SetCallErrorHandler(handler func(callError *CallError)) {
	chargePoint.callErrorHandler = handler
}

// Connects to the given centralSystemUrl and starts running the I/O loop for the underlying connection.
// The write routine runs on a separate goroutine, while the read routine runs on the caller's routine.
// This means, the function is blocking for as long as the ChargePoint is connected to the CentralSystem.
// Whenever the connection breaks, the function returns.
// Call this function in a separate goroutine, to perform other operations on the main thread.
//
// An error may be returned, if the connection failed or if it broke unexpectedly.
func (chargePoint *ChargePoint)Start(centralSystemUrl string) error {
	// Set internal message handler
	chargePoint.client.SetMessageHandler(chargePoint.ocppMessageHandler)
	// Connect & run
	fullUrl := fmt.Sprintf("%v/%v", centralSystemUrl, chargePoint.Id)
	return chargePoint.client.Start(fullUrl)
}

func (chargePoint *ChargePoint)Stop() {
	chargePoint.client.Stop()
}

func (chargePoint *ChargePoint)SendRequest(request Request) error {
	err := validate.Struct(request)
	if err != nil {
		return err
	}
	if chargePoint.pendingRequest != "" {
		// Cannot send. Protocol is based on response-confirmation
		return errors.Errorf("There already is a pending request %v. Cannot send a further one before receiving a confirmation first", chargePoint.pendingRequest)
	}
	call, err := chargePoint.CreateCall(request.(Request))
	jsonMessage, err := call.MarshalJSON()
	if err != nil {
		return err
	}
	chargePoint.pendingRequest = call.UniqueId
	chargePoint.client.Write([]byte(jsonMessage))
	//TODO: use promise/future for fetching the result
	return nil
}

func (chargePoint *ChargePoint)ocppMessageHandler(data []byte) error {
	parsedJson := ParseRawJsonMessage(data)
	message, err := chargePoint.ParseMessage(parsedJson)
	if err != nil {
		// TODO: handle
		log.Printf("Error while parsing message: %v", err)
		return err
	}
	switch message.GetMessageTypeId() {
	case CALL:
		call := message.(*Call)
		chargePoint.callHandler(call)
	case CALL_RESULT:
		callResult := message.(*CallResult)
		chargePoint.pendingRequest = ""
		chargePoint.callResultHandler(callResult)
	case CALL_ERROR:
		callError := message.(*CallError)
		chargePoint.pendingRequest = ""
		chargePoint.callErrorHandler(callError)
	}
	return nil
}

func (chargePoint *ChargePoint)SendMessage(message Message) error {
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
		if chargePoint.pendingRequest != "" {
			// Cannot send. Protocol is based on response-confirmation
			return errors.Errorf("There already is a pending request %v. Cannot send a further one before receiving a confirmation first", chargePoint.pendingRequest)
		}
		chargePoint.PendingRequests[message.GetUniqueId()] = call.Payload
		chargePoint.pendingRequest = call.UniqueId
	}
	chargePoint.client.Write([]byte(jsonMessage))
	//TODO: use promise/future for fetching the result
	return nil
}
