package ocppj

import (
	"fmt"

	"gopkg.in/go-playground/validator.v9"

	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ws"
)

// The endpoint waiting for incoming connections from OCPP clients, in an OCPP-J topology.
// During message exchange, the two roles may be reversed (depending on the message direction), but a server struct remains associated to a central system.
type Server struct {
	Endpoint
	server                    ws.WsServer
	checkClientHandler        ws.CheckClientHandler
	newClientHandler          ClientHandler
	disconnectedClientHandler ClientHandler
	requestHandler            RequestHandler
	responseHandler           ResponseHandler
	errorHandler              ErrorHandler
	invalidMessageHook        InvalidMessageHook
	dispatcher                ServerDispatcher
	RequestState              ServerState
}

type ClientHandler func(client ws.Channel)
type RequestHandler func(client ws.Channel, request ocpp.Request, requestId string, action string)
type ResponseHandler func(client ws.Channel, response ocpp.Response, requestId string)
type ErrorHandler func(client ws.Channel, err *ocpp.Error, details interface{})
type InvalidMessageHook func(client ws.Channel, err *ocpp.Error, rawJson string, parsedFields []interface{}) *ocpp.Error

// Creates a new Server endpoint.
// Requires a a websocket server. Optionally a structure for queueing/dispatching requests,
// a custom state handler and a list of profiles may be passed.
//
// You may create a simple new server by using these default values:
//
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
func (s *Server) SetRequestHandler(handler RequestHandler) {
	s.requestHandler = handler
}

// Registers a handler for incoming responses.
func (s *Server) SetResponseHandler(handler ResponseHandler) {
	s.responseHandler = handler
}

// Registers a handler for incoming error messages.
func (s *Server) SetErrorHandler(handler ErrorHandler) {
	s.errorHandler = handler
}

// SetInvalidMessageHook registers an optional hook for incoming messages that couldn't be parsed.
// This hook is called when a message is received but cannot be parsed to the target OCPP message struct.
//
// The application is notified synchronously of the error.
// The callback provides the raw JSON string, along with the parsed fields.
// The application MUST return as soon as possible, since the hook is called synchronously and awaits a return value.
//
// The hook does not allow responding to the message directly,
// but the return value will be used to send an OCPP error to the other endpoint.
//
// If no handler is registered (or no error is returned by the hook),
// the internal error message is sent to the client without further processing.
//
// Note: Failing to return from the hook will cause the handler for this client to block indefinitely.
func (s *Server) SetInvalidMessageHook(hook InvalidMessageHook) {
	s.invalidMessageHook = hook
}

// Registers a handler for canceled request messages.
func (s *Server) SetCanceledRequestHandler(handler CanceledRequestHandler) {
	s.dispatcher.SetOnRequestCanceled(handler)
}

// Registers a handler for incoming client connections.
func (s *Server) SetNewClientHandler(handler ClientHandler) {
	s.newClientHandler = handler
}

// Registers a handler for validate incoming client connections.
func (s *Server) SetNewClientValidationHandler(handler ws.CheckClientHandler) {
	s.checkClientHandler = handler
}

// Registers a handler for client disconnections.
func (s *Server) SetDisconnectedClientHandler(handler ClientHandler) {
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
	s.server.SetCheckClientHandler(s.checkClientHandler)
	s.server.SetNewClientHandler(s.onClientConnected)
	s.server.SetDisconnectedClientHandler(s.onClientDisconnected)
	s.server.SetMessageHandler(s.ocppMessageHandler)
	s.dispatcher.Start()
	// Serve & run
	s.server.Start(listenPort, listenPath)
	// TODO: return error?
}

// Stops the server.
// This clears all pending requests and causes the Start function to return.
func (s *Server) Stop() {
	s.dispatcher.Stop()
	s.server.Stop()
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
	call, err := s.CreateCall(request)
	if err != nil {
		return err
	}
	jsonMessage, err := call.MarshalJSON()
	if err != nil {
		return err
	}
	// Will not send right away. Queuing message and let it be processed by dedicated requestPump routine
	if err = s.dispatcher.SendRequest(clientID, RequestBundle{call, jsonMessage}); err != nil {
		log.Errorf("error dispatching request [%s, %s] to %s: %v", call.UniqueId, call.Action, clientID, err)
		return err
	}
	log.Debugf("enqueued CALL [%s, %s] for %s", call.UniqueId, call.Action, clientID)
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
	callResult, err := s.CreateCallResult(response, requestId)
	if err != nil {
		return err
	}
	jsonMessage, err := callResult.MarshalJSON()
	if err != nil {
		return ocpp.NewError(GenericError, err.Error(), requestId)
	}
	if err = s.server.Write(clientID, jsonMessage); err != nil {
		log.Errorf("error sending response [%s] to %s: %v", callResult.GetUniqueId(), clientID, err)
		return ocpp.NewError(GenericError, err.Error(), requestId)
	}
	log.Debugf("sent CALL RESULT [%s] for %s", callResult.GetUniqueId(), clientID)
	log.Debugf("sent JSON message to %s: %s", clientID, string(jsonMessage))
	return nil
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
	callError, err := s.CreateCallError(requestId, errorCode, description, details)
	if err != nil {
		return err
	}
	jsonMessage, err := callError.MarshalJSON()
	if err != nil {
		return ocpp.NewError(GenericError, err.Error(), requestId)
	}
	if err = s.server.Write(clientID, jsonMessage); err != nil {
		log.Errorf("error sending response error [%s] to %s: %v", callError.UniqueId, clientID, err)
		return ocpp.NewError(GenericError, err.Error(), requestId)
	}
	log.Debugf("sent CALL ERROR [%s] for %s", callError.UniqueId, clientID)
	return nil
}

func (s *Server) ocppMessageHandler(wsChannel ws.Channel, data []byte) error {
	parsedJson, err := ParseRawJsonMessage(data)
	if err != nil {
		log.Error(err)
		return err
	}
	log.Debugf("received JSON message from %s: %s", wsChannel.ID(), string(data))
	// Get pending requests for client
	pending := s.RequestState.GetClientState(wsChannel.ID())
	message, err := s.ParseMessage(parsedJson, pending)
	if err != nil {
		ocppErr := err.(*ocpp.Error)
		messageID := ocppErr.MessageId
		// Support ad-hoc callback for invalid message handling
		if s.invalidMessageHook != nil {
			err2 := s.invalidMessageHook(wsChannel, ocppErr, string(data), parsedJson)
			// If the hook returns an error, use it as output error. If not, use the original error.
			if err2 != nil {
				ocppErr = err2
				ocppErr.MessageId = messageID
			}
		}
		err = ocppErr
		// Send error to other endpoint if a message ID is available
		if ocppErr.MessageId != "" {
			err2 := s.SendError(wsChannel.ID(), ocppErr.MessageId, ocppErr.Code, ocppErr.Description, nil)
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
			log.Debugf("handling incoming CALL [%s, %s] from %s", call.UniqueId, call.Action, wsChannel.ID())
			if s.requestHandler != nil {
				s.requestHandler(wsChannel, call.Payload, call.UniqueId, call.Action)
			}
		case CALL_RESULT:
			callResult := message.(*CallResult)
			log.Debugf("handling incoming CALL RESULT [%s] from %s", callResult.UniqueId, wsChannel.ID())
			s.dispatcher.CompleteRequest(wsChannel.ID(), callResult.GetUniqueId())
			if s.responseHandler != nil {
				s.responseHandler(wsChannel, callResult.Payload, callResult.UniqueId)
			}
		case CALL_ERROR:
			callError := message.(*CallError)
			log.Debugf("handling incoming CALL RESULT [%s] from %s", callError.UniqueId, wsChannel.ID())
			s.dispatcher.CompleteRequest(wsChannel.ID(), callError.GetUniqueId())
			if s.errorHandler != nil {
				s.errorHandler(wsChannel, ocpp.NewError(callError.ErrorCode, callError.ErrorDescription, callError.UniqueId), callError.ErrorDetails)
			}
		}
	}
	return nil
}

// HandleFailedResponseError allows to handle failures while sending responses (either CALL_RESULT or CALL_ERROR).
// It internally analyzes and creates an ocpp.Error based on the given error.
// It will the attempt to send it to the client.
//
// The function helps to prevent starvation on the other endpoint, which is caused by a response never reaching it.
// The method will, however, only attempt to send a default error once.
// If this operation fails, the other endpoint may still starve.
func (s *Server) HandleFailedResponseError(clientID string, requestID string, err error, featureName string) {
	log.Debugf("handling error for failed response [%s]", requestID)
	var responseErr *ocpp.Error
	// There's several possible errors: invalid profile, invalid payload or send error
	switch err.(type) {
	case validator.ValidationErrors:
		// Validation error
		validationErr := err.(validator.ValidationErrors)
		responseErr = errorFromValidation(validationErr, requestID, featureName)
	case *ocpp.Error:
		// Internal OCPP error
		responseErr = err.(*ocpp.Error)
	case error:
		// Unknown error
		responseErr = ocpp.NewError(GenericError, err.Error(), requestID)
	}
	// Send an OCPP error to the target, since no regular response could be sent
	_ = s.SendError(clientID, requestID, responseErr.Code, responseErr.Description, nil)
}

func (s *Server) onClientConnected(ws ws.Channel) {
	// Create state for connected client
	s.dispatcher.CreateClient(ws.ID())
	// Invoke callback
	if s.newClientHandler != nil {
		s.newClientHandler(ws)
	}
}

func (s *Server) onClientDisconnected(ws ws.Channel) {
	// Clear state for disconnected client
	s.dispatcher.DeleteClient(ws.ID())
	s.RequestState.ClearClientPendingRequest(ws.ID())
	// Invoke callback
	if s.disconnectedClientHandler != nil {
		s.disconnectedClientHandler(ws)
	}
}
