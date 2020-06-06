package ocppj

import (
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ws"
	"github.com/pkg/errors"
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
	clientPendingMessages     map[string]string
}

// Creates a new Server endpoint.
// Requires a a websocket server and a list of profiles (optional).
func NewServer(wsServer ws.WsServer, profiles ...*ocpp.Profile) *Server {
	endpoint := Endpoint{pendingRequests: map[string]ocpp.Request{}}
	for _, profile := range profiles {
		endpoint.AddProfile(profile)
	}
	if wsServer != nil {
		return &Server{Endpoint: endpoint, server: wsServer, clientPendingMessages: map[string]string{}}
	} else {
		return &Server{Endpoint: endpoint, server: &ws.Server{}, clientPendingMessages: map[string]string{}}
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
		delete(s.clientPendingMessages, ws.GetID())
		if s.disconnectedClientHandler != nil {
			s.disconnectedClientHandler(ws.GetID())
		}
	})
	s.server.SetMessageHandler(s.ocppMessageHandler)
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
// - message validation fails (request is malformed)
//
// - another request for that client is already pending
//
// - the endpoint doesn't support the feature
//
// - a network error occurred
func (s *Server) SendRequest(clientID string, request ocpp.Request) error {
	err := Validate.Struct(request)
	if err != nil {
		return err
	}
	req, ok := s.clientPendingMessages[clientID]
	if ok {
		// Cannot send. Protocol is based on request-response
		return errors.Errorf("There already is a pending request %v for client %v. Cannot send a further one before receiving a response first", req, clientID)
	}
	call, err := s.CreateCall(request.(ocpp.Request))
	if err != nil {
		return err
	}
	jsonMessage, err := call.MarshalJSON()
	if err != nil {
		return err
	}
	s.clientPendingMessages[clientID] = call.UniqueId
	err = s.server.Write(clientID, []byte(jsonMessage))
	if err != nil {
		// Clear pending request
		s.DeletePendingRequest(call.GetUniqueId())
		delete(s.clientPendingMessages, clientID)
	}
	return err
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
		delete(s.clientPendingMessages, wsChannel.GetID())
		s.responseHandler(wsChannel.GetID(), callResult.Payload, callResult.UniqueId)
	case CALL_ERROR:
		callError := message.(*CallError)
		delete(s.clientPendingMessages, wsChannel.GetID())
		s.errorHandler(wsChannel.GetID(), ocpp.NewError(callError.ErrorCode, callError.ErrorDescription, callError.UniqueId), callError.ErrorDetails)
	}
	return nil
}
