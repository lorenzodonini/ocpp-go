package ocppj

import (
	"fmt"
	"github.com/lorenzodonini/go-ocpp/ws"
	"github.com/pkg/errors"
	"log"
)

type ChargePoint struct {
	Endpoint
	client            ws.WsClient
	Id                string
	callHandler       func(call *Call)
	callResultHandler func(callResult *CallResult)
	callErrorHandler  func(callError *CallError)
	hasPendingRequest bool
}

func NewChargePoint(id string, wsClient ws.WsClient, profiles ...*Profile) *ChargePoint {
	endpoint := Endpoint{pendingRequests: map[string]Request{}}
	for _, profile := range profiles {
		endpoint.AddProfile(profile)
	}
	if wsClient != nil {
		return &ChargePoint{Endpoint: endpoint, client: wsClient, Id: id, hasPendingRequest: false}
	} else {
		return &ChargePoint{Endpoint: endpoint, client: ws.NewClient(), Id: id, hasPendingRequest: false}
	}
}

func (chargePoint *ChargePoint) SetCallHandler(handler func(call *Call)) {
	chargePoint.callHandler = handler
}

func (chargePoint *ChargePoint) SetCallResultHandler(handler func(callResult *CallResult)) {
	chargePoint.callResultHandler = handler
}

func (chargePoint *ChargePoint) SetCallErrorHandler(handler func(callError *CallError)) {
	chargePoint.callErrorHandler = handler
}

// Connects to the given centralSystemUrl and starts running the I/O loop for the underlying connection.
// The write routine runs on a separate goroutine, while the read routine runs on the caller's routine.
// This means, the function is blocking for as long as the ChargePoint is connected to the CentralSystem.
// Whenever the connection is ended, the function returns.
// Call this function in a separate goroutine, to perform other operations on the main thread.
//
// An error may be returned, if the connection failed or if it broke unexpectedly.
func (chargePoint *ChargePoint) Start(centralSystemUrl string) error {
	// Set internal message handler
	chargePoint.client.SetMessageHandler(chargePoint.ocppMessageHandler)
	// Connect & run
	fullUrl := fmt.Sprintf("%v/%v", centralSystemUrl, chargePoint.Id)
	return chargePoint.client.Start(fullUrl)
}

func (chargePoint *ChargePoint) Stop() {
	chargePoint.client.Stop()
	chargePoint.clearPendingRequests()
	chargePoint.hasPendingRequest = false
}

func (chargePoint *ChargePoint) SendRequest(request Request) error {
	err := validate.Struct(request)
	if err != nil {
		return err
	}
	if chargePoint.hasPendingRequest {
		// Cannot send. Protocol is based on response-confirmation
		return errors.Errorf("There already is a pending request. Cannot send a further one before receiving a confirmation first")
	}
	call, err := chargePoint.CreateCall(request.(Request))
	if err != nil {
		return err
	}
	jsonMessage, err := call.MarshalJSON()
	if err != nil {
		return err
	}
	chargePoint.hasPendingRequest = true
	err = chargePoint.client.Write([]byte(jsonMessage))
	if err != nil {
		// Clear pending request
		chargePoint.DeletePendingRequest(call.GetUniqueId())
		chargePoint.hasPendingRequest = false
	}
	return err
}

func (chargePoint *ChargePoint) SendConfirmation(requestId string, confirmation Confirmation) error {
	err := validate.Struct(confirmation)
	if err != nil {
		return err
	}
	callResult, err := chargePoint.CreateCallResult(confirmation, requestId)
	if err != nil {
		return err
	}
	jsonMessage, err := callResult.MarshalJSON()
	if err != nil {
		return err
	}
	return chargePoint.client.Write([]byte(jsonMessage))
}

func (chargePoint *ChargePoint) SendError(requestId string, errorCode ErrorCode, description string, details interface{}) error {
	//TODO: check if error code is valid
	callError := chargePoint.CreateCallError(requestId, errorCode, description, details)
	err := validate.Struct(callError)
	jsonMessage, err := callError.MarshalJSON()
	if err != nil {
		return err
	}
	return chargePoint.client.Write([]byte(jsonMessage))
}

func (chargePoint *ChargePoint) ocppMessageHandler(data []byte) error {
	parsedJson := ParseRawJsonMessage(data)
	message, err := chargePoint.ParseMessage(parsedJson)
	if err != nil {
		if err.MessageId != "" {
			err2 := chargePoint.SendError(err.MessageId, err.ErrorCode, err.Error.Error(), nil)
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
		chargePoint.callHandler(call)
	case CALL_RESULT:
		callResult := message.(*CallResult)
		chargePoint.hasPendingRequest = false
		chargePoint.callResultHandler(callResult)
	case CALL_ERROR:
		callError := message.(*CallError)
		chargePoint.hasPendingRequest = false
		chargePoint.callErrorHandler(callError)
	}
	return nil
}

//func (chargePoint *ChargePoint) SendMessage(message Message) error {
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
//		if chargePoint.hasPendingRequest {
//			// Cannot send. Protocol is based on response-confirmation
//			return errors.Errorf("There already is a pending request. Cannot send a further one before receiving a confirmation first")
//		}
//		chargePoint.pendingRequests[message.GetUniqueId()] = call.Payload
//		chargePoint.hasPendingRequest = true
//	}
//	err = chargePoint.client.Write([]byte(jsonMessage))
//	if err != nil {
//		chargePoint.DeletePendingRequest(message.GetUniqueId())
//		chargePoint.hasPendingRequest = false
//	}
//	return err
//}
