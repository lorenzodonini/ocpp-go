package ocppj

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ws"
	log "github.com/sirupsen/logrus"
)

// The endpoint waiting for incoming connections from OCPP clients, in an OCPP-J topology.
// During message exchange, the two roles may be reversed (depending on the message direction), but a server struct remains associated to a central system.
type Server struct {
	Endpoint
	server                    ws.WsServer
	newClientHandler          func(clientID string)
	disconnectedClientHandler func(clientID string)
	requestHandler            func(clientID string, request ocpp.Request, requestId string, action string)
	responseHandler           func(clientID string, response ocpp.Response, requestId string)
	errorHandler              func(clientID string, err *ocpp.Error, details interface{})
	requestQueueMap           ServerQueueMap
	readyForDispatch          chan string
	requestChannel            chan string
	pendingRequests           map[string]string
}

// Creates a new Server endpoint.
// Requires a a websocket server, a structure for queueing requests, and a list of profiles (optional).
//
// You may create a simple new server by using these default values:
//	s := ocppj.NewServer(ws.NewServer(), ocppj.NewFIFOQueueMap(0))
func NewServer(wsServer ws.WsServer, requestMap ServerQueueMap, profiles ...*ocpp.Profile) *Server {
	endpoint := Endpoint{pendingRequests: map[string]ocpp.Request{}}
	for _, profile := range profiles {
		endpoint.AddProfile(profile)
	}
	if wsServer != nil {
		return &Server{Endpoint: endpoint, server: wsServer, requestQueueMap: requestMap, readyForDispatch: make(chan string, 1), pendingRequests: map[string]string{}}
	} else {
		return &Server{Endpoint: endpoint, server: &ws.Server{}, requestQueueMap: requestMap, readyForDispatch: make(chan string, 1), pendingRequests: map[string]string{}}
	}
}

// Registers a handler for incoming requests.
func (s *Server) SetRequestHandler(handler func(clientID string, request ocpp.Request, requestId string, action string)) {
	s.requestHandler = handler
}

// Registers a handler for incoming responses.
func (s *Server) SetResponseHandler(handler func(clientID string, response ocpp.Response, requestId string)) {
	s.responseHandler = handler
}

// Registers a handler for incoming error messages.
func (s *Server) SetErrorHandler(handler func(clientID string, err *ocpp.Error, details interface{})) {
	s.errorHandler = handler
}

// Registers a handler for incoming client connections.
func (s *Server) SetNewClientHandler(handler func(clientID string)) {
	s.newClientHandler = handler
}

// Registers a handler for client disconnections.
func (s *Server) SetDisconnectedClientHandler(handler func(clientID string)) {
	s.disconnectedClientHandler = handler
}

// Starts the underlying Websocket server on a specified listenPort and listenPath.
// The function runs indefinitely, until the server is stopped.
//
// Call this function in a separate goroutine, to perform other operations on the main thread.
//
// An error may be returned, if the websocket server couldn't be started.
func (s *Server) Start(listenPort int, listenPath string) {
	// Set internal message handler
	s.server.SetNewClientHandler(func(ws ws.Channel) {
		if s.newClientHandler != nil {
			s.newClientHandler(ws.GetID())
		}
	})
	s.server.SetDisconnectedClientHandler(func(ws ws.Channel) {
		//TODO: handle reconnection and don't delete request queue
		s.requestQueueMap.Remove(ws.GetID())
		s.requestChannel <- ws.GetID()
		if s.disconnectedClientHandler != nil {
			s.disconnectedClientHandler(ws.GetID())
		}
	})
	s.server.SetMessageHandler(s.ocppMessageHandler)
	s.requestChannel = make(chan string, 1)
	go s.requestPump()
	// Serve & run
	// TODO: return error?
	s.server.Start(listenPort, listenPath)
}

// Stops the server.
// This clears all pending requests and causes the Start function to return.
func (s *Server) Stop() {
	s.server.Stop()
	s.clearPendingRequests()
}

// Sends an OCPP Request to a client, identified by the clientID parameter.
//
// Returns an error in the following cases:
//
// - the server wasn't started
//
// - message validation fails (request is malformed)
//
// - the endpoint doesn't support the feature
//
// - the output queue is full
func (s *Server) SendRequest(clientID string, request ocpp.Request) error {
	if s.requestChannel == nil {
		return fmt.Errorf("ocppj server is not started, couldn't send request")
	}
	err := Validate.Struct(request)
	if err != nil {
		return err
	}
	call, err := s.CreateCall(request.(ocpp.Request))
	if err != nil {
		return err
	}
	jsonMessage, err := call.MarshalJSON()
	if err != nil {
		return err
	}
	// Will not send right away. Queuing message and let it be processed by dedicated requestPump routine
	q := s.requestQueueMap.GetOrCreate(clientID)
	if q.IsFull() {
		return fmt.Errorf("request queue for client %v is full, cannot send new request", clientID)
	}
	err = q.Push(RequestBundle{call, jsonMessage})
	if err != nil {
		return err
	}
	log.Debugf("enqueued request %v - %v for client %v", call.UniqueId, call.Action, clientID)
	// Notify requestPump that a new request for ClientID is pending
	s.requestChannel <- clientID
	return nil
}

// requestPump processes new outgoing requests for each client and makes sure they are processed sequentially.
// This method is executed by a dedicated coroutine as soon as the server is started and runs indefinitely.
func (s *Server) requestPump() {
	var clientID string
	var ok bool
	var rdy bool
	var clientQueue RequestQueue
	clientReadyMap := map[string]bool{} // Empty at the beginning
	for {
		select {
		case clientID, ok = <-s.requestChannel:
			// Check if channel was closed
			if !ok {
				log.Infof("stopped processing requests")
				s.clearPendingRequests()
				s.requestChannel = nil
				return
			}
			clientQueue, ok = s.requestQueueMap.Get(clientID)
			// Check whether there is a request queue for the specified client
			if !ok {
				// No client queue found, deleting the ready flag
				delete(clientReadyMap, clientID)
				break
			}
			// Check whether can transmit to client
			rdy, ok = clientReadyMap[clientID]
			if !ok {
				// First request sent to client. Setting ready flag to true
				rdy = true
				clientReadyMap[clientID] = rdy
			}
		case clientID = <-s.readyForDispatch:
			// Client can now transmit again
			rdy = true
			clientReadyMap[clientID] = rdy
			clientQueue, _ = s.requestQueueMap.Get(clientID)
		}
		// Only dispatch request if able to send and request queue isn't empty
		if rdy && !clientQueue.IsEmpty() {
			s.dispatchNextRequest(clientID)
			// Update ready state
			rdy = false
			clientReadyMap[clientID] = rdy
		}
	}
}

func (s *Server) dispatchNextRequest(clientID string) {
	// Get first element in queue
	q, _ := s.requestQueueMap.Get(clientID)
	el := q.Peek()
	bundle, _ := el.(RequestBundle)
	jsonMessage := bundle.Data
	callID := bundle.Call.GetUniqueId()
	s.AddPendingRequest(callID, bundle.Call.Payload)
	s.pendingRequests[clientID] = callID
	err := s.server.Write(clientID, jsonMessage)
	if err != nil {
		log.Errorf("error while sending message: %v", err)
		//TODO: handle retransmission instead of removing pending request
		s.completePendingRequest(clientID, callID)
		if s.errorHandler != nil {
			s.errorHandler(clientID, ocpp.NewError(GenericError, err.Error(), callID), err)
		}
	} else {
		// Transmitted correctly
		log.Debugf("sent request %v to client %v: %v", callID, clientID, string(jsonMessage))
	}
}

func (s *Server) completePendingRequest(clientID string, requestID string) {
	q, _ := s.requestQueueMap.Get(clientID)
	el := q.Peek()
	if el == nil {
		log.Errorf("attempting to pop front of queue, but queue is empty")
		return
	}
	bundle, _ := el.(RequestBundle)
	callID := bundle.Call.GetUniqueId()
	if callID != requestID {
		log.Fatalf("internal state mismatch: received response for %v but expected response for %v", requestID, callID)
		return
	}
	q.Pop()
	delete(s.pendingRequests, clientID)
	s.DeletePendingRequest(callID)
	log.Debugf("removed request %v from front of queue", callID)
	// Signal that next message in queue may be sent
	s.readyForDispatch <- clientID
}

// Sends an OCPP Response to a client, identified by the clientID parameter.
// The requestID parameter is required and identifies the previously received request.
//
// Returns an error in the following cases:
//
// - message validation fails (response is malformed)
//
// - the endpoint doesn't support the feature
//
// - a network error occurred
func (s *Server) SendResponse(clientID string, requestId string, response ocpp.Response) error {
	err := Validate.Struct(response)
	if err != nil {
		return err
	}
	callResult, err := s.CreateCallResult(response, requestId)
	if err != nil {
		return err
	}
	jsonMessage, err := callResult.MarshalJSON()
	if err != nil {
		return err
	}
	return s.server.Write(clientID, []byte(jsonMessage))
}

// Sends an OCPP Error to a client, identified by the clientID parameter.
// The requestID parameter is required and identifies the previously received request.
//
// Returns an error in the following cases:
//
// - message validation fails (error is malformed)
//
// - a network error occurred
func (s *Server) SendError(clientID string, requestId string, errorCode ocpp.ErrorCode, description string, details interface{}) error {
	callError := s.CreateCallError(requestId, errorCode, description, details)
	err := Validate.Struct(callError)
	if err != nil {
		return err
	}
	jsonMessage, err := callError.MarshalJSON()
	if err != nil {
		return err
	}
	return s.server.Write(clientID, []byte(jsonMessage))
}

func (s *Server) ocppMessageHandler(wsChannel ws.Channel, data []byte) error {
	parsedJson := ParseRawJsonMessage(data)
	message, err := s.ParseMessage(parsedJson)
	if err != nil {
		if err.MessageId != "" {
			err2 := s.SendError(wsChannel.GetID(), err.MessageId, err.Code, err.Description, nil)
			if err2 != nil {
				return err2
			}
		}
		log.Error(err)
		return err
	}
	switch message.GetMessageTypeId() {
	case CALL:
		call := message.(*Call)
		s.requestHandler(wsChannel.GetID(), call.Payload, call.UniqueId, call.Action)
	case CALL_RESULT:
		callResult := message.(*CallResult)
		s.completePendingRequest(wsChannel.GetID(), callResult.GetUniqueId()) // Remove current request from queue and send next one
		if s.responseHandler != nil {
			s.responseHandler(wsChannel.GetID(), callResult.Payload, callResult.UniqueId)
		}
	case CALL_ERROR:
		callError := message.(*CallError)
		s.completePendingRequest(wsChannel.GetID(), callError.GetUniqueId()) // Remove current request from queue and send next one
		if s.errorHandler != nil {
			s.errorHandler(wsChannel.GetID(), ocpp.NewError(callError.ErrorCode, callError.ErrorDescription, callError.UniqueId), callError.ErrorDetails)
		}
	}
	return nil
}
