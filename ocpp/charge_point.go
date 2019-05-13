package ocpp

import (
	"errors"
	"fmt"
	"github.com/lorenzodonini/go-ocpp/ws"
	"log"
)

type ChargePoint struct {
	Endpoint
	client ws.WsClient
	Id string
	callHandler func(call *Call)
	callResultHandler func(callResult *CallResult)
	callErrorHandler func(callError *CallError)
}

func NewChargePoint(id string, wsClient ws.WsClient, profiles ...*Profile) *ChargePoint {
	endpoint := Endpoint{PendingRequests: map[string]Request{}}
	for _, profile := range profiles {
		endpoint.AddProfile(profile)
	}
	if wsClient != nil {
		return &ChargePoint{Endpoint: endpoint, client: wsClient, Id: id}
	} else {
		return &ChargePoint{Endpoint: endpoint, client: &ws.Client{}, Id: id}
	}
}

func (chargePoint *ChargePoint)SeCallHandler(handler func(call *Call)) {
	chargePoint.callHandler = handler
}

func (chargePoint *ChargePoint)SeCallResultHandler(handler func(callResult *CallResult)) {
	chargePoint.callResultHandler = handler
}

func (chargePoint *ChargePoint)SeCalleHandler(handler func(callError *CallError)) {
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

func (chargePoint *ChargePoint)SendMessage(message *Message) error {
	jsonMessage, err := CreateJsonMessage(message)
	if err != nil {
		return err
	}
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
	ocppMessage, ok := message.(Message)
	if !ok {
		return errors.New("couldn't convert parsed data to Message type")
	}
	switch ocppMessage.MessageTypeId {
	case CALL:
		call := message.(Call)
		chargePoint.callHandler(&call)
	case CALL_RESULT:
		callResult := message.(CallResult)
		chargePoint.callResultHandler(&callResult)
	case CALL_ERROR:
		callError := message.(CallError)
		chargePoint.callErrorHandler(&callError)
	}
	return nil
}
