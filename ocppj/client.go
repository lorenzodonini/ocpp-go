package ocppj

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ws"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// The endpoint initiating the connection to an OCPP server, in an OCPP-J topology.
// During message exchange, the two roles may be reversed (depending on the message direction), but a client struct remains associated to a charge point/charging station.
type Client struct {
	Endpoint
	client            ws.WsClient
	Id                string
	requestHandler    func(request ocpp.Request, requestId string, action string)
	responseHandler   func(response ocpp.Response, requestId string)
	errorHandler      func(err *ocpp.Error, details interface{})
	hasPendingRequest bool
}

// Creates a new Client endpoint.
// Requires a unique client ID, a websocket client and a list of profiles (optional).
func NewClient(id string, wsClient ws.WsClient, profiles ...*ocpp.Profile) *Client {
	endpoint := Endpoint{pendingRequests: map[string]ocpp.Request{}}
	for _, profile := range profiles {
		endpoint.AddProfile(profile)
	}
	if wsClient != nil {
		return &Client{Endpoint: endpoint, client: wsClient, Id: id, hasPendingRequest: false}
	} else {
		return &Client{Endpoint: endpoint, client: ws.NewClient(), Id: id, hasPendingRequest: false}
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
	return c.client.Start(fullUrl)
}

// Stops the client.
// The underlying I/O loop is stopped and all pending requests are cleared.
func (c *Client) Stop() {
	c.client.Stop()
	c.clearPendingRequests()
	c.hasPendingRequest = false
}

// Sends an OCPP Request to the server.
//
// Returns an error in the following cases:
//
// - message validation fails (request is malformed)
//
// - another request is already pending
//
// - the endpoint doesn't support the feature
//
// - a network error occurred
func (c *Client) SendRequest(request ocpp.Request) error {
	err := Validate.Struct(request)
	if err != nil {
		return err
	}
	if c.hasPendingRequest {
		// Cannot send. Protocol is based on request-response
		return errors.Errorf("There already is a pending request. Cannot send a further one before receiving a response first")
	}
	call, err := c.CreateCall(request.(ocpp.Request))
	if err != nil {
		return err
	}
	jsonMessage, err := call.MarshalJSON()
	if err != nil {
		return err
	}
	c.hasPendingRequest = true
	err = c.client.Write([]byte(jsonMessage))
	if err != nil {
		// Clear pending request
		c.DeletePendingRequest(call.GetUniqueId())
		c.hasPendingRequest = false
	}
	return err
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
		c.hasPendingRequest = false
		c.responseHandler(callResult.Payload, callResult.UniqueId)
	case CALL_ERROR:
		callError := message.(*CallError)
		c.hasPendingRequest = false
		c.errorHandler(ocpp.NewError(callError.ErrorCode, callError.ErrorDescription, callError.UniqueId), callError.ErrorDetails)
	}
	return nil
}
