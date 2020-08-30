package ocppj

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ws"
	log "github.com/sirupsen/logrus"
)

// The endpoint initiating the connection to an OCPP server, in an OCPP-J topology.
// During message exchange, the two roles may be reversed (depending on the message direction), but a client struct remains associated to a charge point/charging station.
type Client struct {
	Endpoint
	client           ws.WsClient
	Id               string
	requestHandler   func(request ocpp.Request, requestId string, action string)
	responseHandler  func(response ocpp.Response, requestId string)
	errorHandler     func(err *ocpp.Error, details interface{})
	requestQueue     RequestQueue
	requestChannel   chan bool
	readyForDispatch chan bool
}

// Creates a new Client endpoint.
// Requires a unique client ID, a websocket client, a struct for queueing requests and a list of supported profiles (optional).
//
// You may create a simple new server by using these default values:
//	s := ocppj.NewClient(ws.NewClient(), ocppj.NewFIFOClientQueue())
func NewClient(id string, wsClient ws.WsClient, requestQueue RequestQueue, profiles ...*ocpp.Profile) *Client {
	endpoint := Endpoint{pendingRequests: map[string]ocpp.Request{}}
	for _, profile := range profiles {
		endpoint.AddProfile(profile)
	}
	if wsClient != nil {
		return &Client{Endpoint: endpoint, client: wsClient, Id: id, requestQueue: requestQueue, readyForDispatch: make(chan bool)}
	} else {
		return &Client{Endpoint: endpoint, client: ws.NewClient(), Id: id, requestQueue: requestQueue, readyForDispatch: make(chan bool)}
	}
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

// Connects to the given serverURL and starts running the I/O loop for the underlying connection.
// The write routine runs on a separate goroutine, while the read routine runs on the caller's routine.
// This means, the function is blocking for as long as the Client is connected to the Server.
//
// Whenever the connection is ended, the function returns.
//
// Call this function in a separate goroutine, to perform other operations on the main thread.
//
// An error may be returned, if the connection failed or if it broke unexpectedly.
func (c *Client) Start(serverURL string) error {
	// Set internal message handler
	c.client.SetMessageHandler(c.ocppMessageHandler)
	// Connect & run
	fullUrl := fmt.Sprintf("%v/%v", serverURL, c.Id)
	err := c.client.Start(fullUrl)
	if err == nil {
		c.requestChannel = make(chan bool, 1)
		go c.requestPump()
	}
	return err
}

// Stops the client.
// The underlying I/O loop is stopped and all pending requests are cleared.
func (c *Client) Stop() {
	c.client.Stop()
	c.clearPendingRequests()
	close(c.requestChannel)
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
	if c.requestChannel == nil {
		return fmt.Errorf("ocppj client is not started, couldn't send request")
	}
	err := Validate.Struct(request)
	if err != nil {
		return err
	}
	call, err := c.CreateCall(request)
	if err != nil {
		return err
	}
	jsonMessage, err := call.MarshalJSON()
	if err != nil {
		return err
	}
	// Will not send right away. Queuing message and let it be processed by dedicated requestPump routine
	if err := c.requestQueue.Push(RequestBundle{Call: call, Data: jsonMessage}); err != nil {
		log.Errorf("request %v - %v: %v", call.UniqueId, call.Action, err)
		return err
	}
	log.Debugf("enqueued request %v - %v", call.UniqueId, call.Action)
	c.requestChannel <- true
	return nil
}

// requestPump processes new outgoing requests and makes sure they are enqueued correctly.
// This method is executed by a dedicated coroutine as soon as the client is started and runs indefinitely.
func (c *Client) requestPump() {
	rdy := true // Ready to transmit at the beginning
	for {
		select {
		case _, ok := <-c.requestChannel:
			// Enqueue new request
			if !ok {
				log.Infof("stopped processing requests")
				c.requestQueue.Init()
				c.requestChannel = nil
				return
			}
		case rdy = <-c.readyForDispatch:
		}
		// Only dispatch request if able to send and request queue isn't empty
		if rdy && !c.requestQueue.IsEmpty() {
			c.dispatchNextRequest()
			rdy = false
		}
	}
}

func (c *Client) dispatchNextRequest() {
	// Get first element in queue
	el := c.requestQueue.Peek()
	bundle, _ := el.(RequestBundle)
	jsonMessage := bundle.Data
	c.AddPendingRequest(bundle.Call.UniqueId, bundle.Call.Payload)

	err := c.client.Write(jsonMessage)
	if err != nil {
		log.Errorf("error while sending message: %v", err)
		//TODO: handle retransmission instead of removing pending request
		c.DeletePendingRequest(bundle.Call.GetUniqueId())
		c.completePendingRequest(bundle.Call.GetUniqueId())
		if c.errorHandler != nil {
			c.errorHandler(ocpp.NewError(GenericError, err.Error(), bundle.Call.GetUniqueId()), err)
		}
	} else {
		// Transmitted correctly
		log.Debugf("sent request %v: %v", bundle.Call.UniqueId, string(jsonMessage))
	}
}

func (c *Client) completePendingRequest(requestId string) {
	el := c.requestQueue.Peek()
	if el == nil {
		log.Errorf("attempting to pop front of queue, but queue is empty")
		return
	}
	bundle, _ := el.(RequestBundle)
	if bundle.Call.UniqueId != requestId {
		log.Fatalf("internal state mismatch: received response for %v but expected response for %v", requestId, bundle.Call.UniqueId)
		return
	}
	c.requestQueue.Pop()
	log.Debugf("removed request %v from front of queue", bundle.Call.UniqueId)
	// Signal that next message in queue may be sent
	c.readyForDispatch <- true
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
	err := Validate.Struct(response)
	if err != nil {
		return err
	}
	callResult, err := c.CreateCallResult(response, requestId)
	if err != nil {
		return err
	}
	jsonMessage, err := callResult.MarshalJSON()
	if err != nil {
		return err
	}
	return c.client.Write([]byte(jsonMessage))
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
	callError := c.CreateCallError(requestId, errorCode, description, details)
	err := Validate.Struct(callError)
	if err != nil {
		return err
	}
	jsonMessage, err := callError.MarshalJSON()
	if err != nil {
		return err
	}
	return c.client.Write([]byte(jsonMessage))
}

func (c *Client) ocppMessageHandler(data []byte) error {
	parsedJson := ParseRawJsonMessage(data)
	message, err := c.ParseMessage(parsedJson)
	if err != nil {
		if err.MessageId != "" {
			err2 := c.SendError(err.MessageId, err.Code, err.Description, nil)
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
		c.requestHandler(call.Payload, call.UniqueId, call.Action)
	case CALL_RESULT:
		callResult := message.(*CallResult)
		c.DeletePendingRequest(callResult.GetUniqueId())
		c.completePendingRequest(callResult.UniqueId) // Remove current request from queue and send next one
		if c.responseHandler != nil {
			c.responseHandler(callResult.Payload, callResult.UniqueId)
		}
	case CALL_ERROR:
		callError := message.(*CallError) // Remove current request from queue and send next one
		c.DeletePendingRequest(callError.GetUniqueId())
		c.completePendingRequest(callError.UniqueId)
		if c.errorHandler != nil {
			c.errorHandler(ocpp.NewError(callError.ErrorCode, callError.ErrorDescription, callError.UniqueId), callError.ErrorDetails)
		}
	}
	return nil
}
