package ocppj

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ws"
)

// The endpoint initiating the connection to an OCPP server, in an OCPP-J topology.
// During message exchange, the two roles may be reversed (depending on the message direction), but a client struct remains associated to a charge point/charging station.
type Client struct {
	Endpoint
	client                ws.WsClient
	Id                    string
	requestHandler        func(request ocpp.Request, requestId string, action string)
	responseHandler       func(response ocpp.Response, requestId string)
	errorHandler          func(err *ocpp.Error, details interface{})
	onDisconnectedHandler func(err error)
	onReconnectedHandler  func()
	dispatcher            ClientDispatcher
	RequestState          ClientState
}

// Creates a new Client endpoint.
// Requires a unique client ID, a websocket client, a struct for queueing/dispatching requests,
// a state handler and a list of supported profiles (optional).
//
// You may create a simple new server by using these default values:
//	s := ocppj.NewClient(ws.NewClient(), nil, nil)
//
// The wsClient parameter cannot be nil. Refer to the ws package for information on how to create and
// customize a websocket client.
func NewClient(id string, wsClient ws.WsClient, dispatcher ClientDispatcher, stateHandler ClientState, profiles ...*ocpp.Profile) *Client {
	endpoint := Endpoint{}
	if wsClient == nil {
		panic("wsClient parameter cannot be nil")
	}
	for _, profile := range profiles {
		endpoint.AddProfile(profile)
	}
	if dispatcher == nil {
		dispatcher = NewDefaultClientDispatcher(NewFIFOClientQueue(10))
	}
	if stateHandler == nil {
		stateHandler = NewClientState()
	}
	dispatcher.SetNetworkClient(wsClient)
	dispatcher.SetPendingRequestState(stateHandler)
	return &Client{Endpoint: endpoint, client: wsClient, Id: id, dispatcher: dispatcher, RequestState: stateHandler}
}

// Registers a handler for incoming requests.
func (c *Client) SetRequestHandler(handler func(request ocpp.Request, requestId string, action string)) {
	c.requestHandler = handler
}

// Registers a handler for incoming responses.
func (c *Client) SetResponseHandler(handler func(response ocpp.Response, requestId string)) {
	c.responseHandler = handler
}

// Registers a handler for incoming error messages.
func (c *Client) SetErrorHandler(handler func(err *ocpp.Error, details interface{})) {
	c.errorHandler = handler
}

func (c *Client) SetOnDisconnectedHandler(handler func(err error)) {
	c.onDisconnectedHandler = handler
}

func (c *Client) SetOnReconnectedHandler(handler func()) {
	c.onReconnectedHandler = handler
}

// Registers the handler to be called on timeout.
func (c *Client) SetOnRequestCanceled(handler func(requestId string, request ocpp.Request, err *ocpp.Error)) {
	c.dispatcher.SetOnRequestCanceled(handler)
}

// Connects to the given serverURL and starts running the I/O loop for the underlying connection.
//
// If the connection is established successfully, the function returns control to the caller immediately.
// The read/write routines are run on dedicated goroutines, so the main thread can perform other operations.
//
// In case of disconnection, the client handles re-connection automatically.
// The client will attempt to re-connect to the server forever, until it is stopped by invoking the Stop method.
//
// An error may be returned, if establishing the connection failed.
func (c *Client) Start(serverURL string) error {
	// Set internal message handler
	c.client.SetMessageHandler(c.ocppMessageHandler)
	c.client.SetDisconnectedHandler(c.onDisconnected)
	c.client.SetReconnectedHandler(c.onReconnected)
	// Connect & run
	fullUrl := fmt.Sprintf("%v/%v", serverURL, c.Id)
	err := c.client.Start(fullUrl)
	if err == nil {
		c.dispatcher.Start()
	}
	return err
}

// Stops the client.
// The underlying I/O loop is stopped and all pending requests are cleared.
func (c *Client) Stop() {
	// Overwrite handler to intercept disconnected signal
	cleanupC := make(chan struct{}, 1)
	c.client.SetDisconnectedHandler(func(err error) {
		cleanupC <- struct{}{}
	})
	c.client.Stop()
	c.dispatcher.Stop()
	// Wait for websocket to be cleaned up
	<-cleanupC
}

func (c *Client) IsConnected() bool {
	return c.client.IsConnected()
}

// Sends an OCPP Request to the server.
// The protocol is based on request-response and cannot send multiple messages concurrently.
// To guarantee this, outgoing messages are added to a queue and processed sequentially.
//
// Returns an error in the following cases:
//
// - the client wasn't started
//
// - message validation fails (request is malformed)
//
// - the endpoint doesn't support the feature
//
// - the output queue is full
func (c *Client) SendRequest(request ocpp.Request) error {
	if !c.dispatcher.IsRunning() {
		return fmt.Errorf("ocppj client is not started, couldn't send request")
	}
	call, err := c.CreateCall(request)
	if err != nil {
		return err
	}
	jsonMessage, err := call.MarshalJSON()
	if err != nil {
		return err
	}
	// Message will be processed by dispatcher. A dedicated mechanism allows to delegate the message queue handling.
	if err = c.dispatcher.SendRequest(RequestBundle{Call: call, Data: jsonMessage}); err != nil {
		log.Errorf("error dispatching request [%s, %s]: %v", call.UniqueId, call.Action, err)
		return err
	}
	log.Debugf("enqueued CALL [%s, %s]", call.UniqueId, call.Action)
	return nil
}

// Sends an OCPP Response to the server.
// The requestID parameter is required and identifies the previously received request.
//
// Returns an error in the following cases:
//
// - message validation fails (response is malformed)
//
// - the endpoint doesn't support the feature
//
// - a network error occurred
func (c *Client) SendResponse(requestId string, response ocpp.Response) error {
	callResult, err := c.CreateCallResult(response, requestId)
	if err != nil {
		return err
	}
	jsonMessage, err := callResult.MarshalJSON()
	if err != nil {
		return err
	}
	if err = c.client.Write(jsonMessage); err != nil {
		log.Errorf("error sending response [%s]: %v", callResult.UniqueId, err)
		return err
	}
	log.Debugf("sent CALL RESULT [%s]", callResult.UniqueId)
	return nil
}

// Sends an OCPP Error to the server.
// The requestID parameter is required and identifies the previously received request.
//
// Returns an error in the following cases:
//
// - message validation fails (error is malformed)
//
// - a network error occurred
func (c *Client) SendError(requestId string, errorCode ocpp.ErrorCode, description string, details interface{}) error {
	callError, err := c.CreateCallError(requestId, errorCode, description, details)
	if err != nil {
		return err
	}
	jsonMessage, err := callError.MarshalJSON()
	if err != nil {
		return err
	}
	if err = c.client.Write(jsonMessage); err != nil {
		log.Errorf("error sending response error [%s]: %v", callError.UniqueId, err)
		return err
	}
	log.Debugf("sent CALL ERROR [%s]", callError.UniqueId)
	return nil
}

func (c *Client) ocppMessageHandler(data []byte) error {
	parsedJson, err := ParseRawJsonMessage(data)
	if err != nil {
		log.Error(err)
		return err
	}
	message, err := c.ParseMessage(parsedJson, c.RequestState)
	if err != nil {
		ocppErr := err.(*ocpp.Error)
		if ocppErr.MessageId != "" {
			err2 := c.SendError(ocppErr.MessageId, ocppErr.Code, ocppErr.Description, nil)
			if err2 != nil {
				return err2
			}
		}
		log.Error(err)
		return err
	}
	if message != nil {
		switch message.GetMessageTypeId() {
		case CALL:
			call := message.(*Call)
			log.Debugf("handling incoming CALL [%s, %s]", call.UniqueId, call.Action)
			c.requestHandler(call.Payload, call.UniqueId, call.Action)
		case CALL_RESULT:
			callResult := message.(*CallResult)
			log.Debugf("handling incoming CALL RESULT [%s]", callResult.UniqueId)
			c.dispatcher.CompleteRequest(callResult.GetUniqueId()) // Remove current request from queue and send next one
			if c.responseHandler != nil {
				c.responseHandler(callResult.Payload, callResult.UniqueId)
			}
		case CALL_ERROR:
			callError := message.(*CallError)
			log.Debugf("handling incoming CALL ERROR [%s]", callError.UniqueId)
			c.dispatcher.CompleteRequest(callError.GetUniqueId()) // Remove current request from queue and send next one
			if c.errorHandler != nil {
				c.errorHandler(ocpp.NewError(callError.ErrorCode, callError.ErrorDescription, callError.UniqueId), callError.ErrorDetails)
			}
		}
	}
	return nil
}

func (c *Client) onDisconnected(err error) {
	log.Error("disconnected from server", err)
	c.dispatcher.Pause()
	if c.onDisconnectedHandler != nil {
		c.onDisconnectedHandler(err)
	}
}

func (c *Client) onReconnected() {
	if c.onReconnectedHandler != nil {
		c.onReconnectedHandler()
	}
	c.dispatcher.Resume()
}
