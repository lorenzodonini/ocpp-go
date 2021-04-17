package ocppj

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ws"
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
	dispatcher                ServerDispatcher
	RequestState              ServerState
}

// Creates a new Server endpoint.
// Requires a a websocket server. Optionally a structure for queueing/dispatching requests,
// a custom state handler and a list of profiles may be passed.
//
// You may create a simple new server by using these default values:
//	s := ocppj.NewServer(ws.NewServer(), nil, nil)
//
// The dispatcher's associated ClientState will be set during initialization.
func NewServer(wsServer ws.WsServer, dispatcher ServerDispatcher, stateHandler ServerState, profiles ...*ocpp.Profile) *Server {
	if dispatcher == nil {
		dispatcher = NewDefaultServerDispatcher(NewFIFOQueueMap(0))
	}
	if stateHandler == nil {
		d, ok := dispatcher.(*DefaultServerDispatcher)
		if !ok {
			stateHandler = NewServerState(nil)
		} else {
			stateHandler = d.pendingRequestState
		}
	}
	if wsServer == nil {
		wsServer = ws.NewServer()
	}
	dispatcher.SetNetworkServer(wsServer)
	dispatcher.SetPendingRequestState(stateHandler)

	// Create server and add profiles
	s := Server{Endpoint: Endpoint{}, server: wsServer, RequestState: stateHandler, dispatcher: dispatcher}
	for _, profile := range profiles {
		s.AddProfile(profile)
	}
	return &s
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
//
// The function runs indefinitely, until the server is stopped.
// Invoke this function in a separate goroutine, to perform other operations on the main thread.
//
// An error may be returned, if the websocket server couldn't be started.
func (s *Server) Start(listenPort int, listenPath string) {
	// Set internal message handler
	s.server.SetNewClientHandler(func(ws ws.Channel) {
		if s.newClientHandler != nil {
			s.newClientHandler(ws.GetID())
		}
	})
	s.server.SetDisconnectedClientHandler(s.onDisconnected)
	s.server.SetMessageHandler(s.ocppMessageHandler)
	s.dispatcher.Start()
	// Serve & run
	s.server.Start(listenPort, listenPath)
	// TODO: return error?
}

// Stops the server.
// This clears all pending requests and causes the Start function to return.
func (s *Server) Stop() {
	s.server.Stop()
	s.dispatcher.Stop()
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
	if !s.dispatcher.IsRunning() {
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
	if err := s.dispatcher.SendRequest(clientID, RequestBundle{call, jsonMessage}); err != nil {
		log.Errorf("request %v - %v for client %v: %v", call.UniqueId, call.Action, clientID, err)
		return err
	}
	log.Debugf("enqueued request %v - %v for client %v", call.UniqueId, call.Action, clientID)
	return nil
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
	parsedJson, err := ParseRawJsonMessage(data)
	if err != nil {
		log.Error(err)
		return err
	}
	// Get pending requests for client
	pending := s.RequestState.GetClientState(wsChannel.GetID())
	message, err := s.ParseMessage(parsedJson, pending)
	if err != nil {
		ocppErr := err.(*ocpp.Error)
		if ocppErr.MessageId != "" {
			err2 := s.SendError(wsChannel.GetID(), ocppErr.MessageId, ocppErr.Code, ocppErr.Description, nil)
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
		s.dispatcher.CompleteRequest(wsChannel.GetID(), callResult.GetUniqueId())
		if s.responseHandler != nil {
			s.responseHandler(wsChannel.GetID(), callResult.Payload, callResult.UniqueId)
		}
	case CALL_ERROR:
		callError := message.(*CallError)
		s.dispatcher.CompleteRequest(wsChannel.GetID(), callError.GetUniqueId())
		if s.errorHandler != nil {
			s.errorHandler(wsChannel.GetID(), ocpp.NewError(callError.ErrorCode, callError.ErrorDescription, callError.UniqueId), callError.ErrorDetails)
		}
	}
	return nil
}

func (s *Server) onDisconnected(ws ws.Channel) {
	// Clear state for disconnected client
	s.dispatcher.DeleteClient(ws.GetID())
	s.RequestState.ClearClientPendingRequest(ws.GetID())
	// Invoke callback
	if s.disconnectedClientHandler != nil {
		s.disconnectedClientHandler(ws.GetID())
	}
}
