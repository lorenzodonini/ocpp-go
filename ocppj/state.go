package ocppj

import (
	"sync"

	"github.com/lorenzodonini/ocpp-go/ocpp"
)

// Contains the pending request state for messages, associated to a single client-server channel.
// It is used to separate endpoint logic from state management.
type ClientState interface {
	// Sets a Request as pending on the endpoint. Requests are considered pending until a response was received.
	// The function expects a unique message ID and the Request.
	// If an element with the same requestID exists, the new one will be ignored.
	AddPendingRequest(requestID string, req ocpp.Request)
	// Retrieves a pending Request, using the message ID.
	// If no request for the passed message ID is found, a false flag is returned.
	GetPendingRequest(requestID string) (ocpp.Request, bool)
	// Deletes a pending Request from the endpoint, using the message ID.
	// If no such message is currently stored as pending, the call has no effect.
	DeletePendingRequest(requestID string)
	// Clears all currently pending requests. Any confirmation/error,
	// received as a response to a cleared request, will be ignored.
	ClearPendingRequests()
	// Returns true if there currently is at least one pending request, false otherwise.
	HasPendingRequest() bool
}

// ----------------------------
// Request State implementation
// ----------------------------

// Simple implementation of ClientState.
// Supports a single pending request. To add a new pending request, the previous one needs to be deleted.
//
// Uses a mutex internally for concurrent access to the data struct.
type clientState struct {
	requestID      string
	pendingRequest pendingRequest
	mutex          sync.RWMutex
}

// Creates a simple struct implementing ClientState, to be used by client/server dispatchers.
func NewClientState() ClientState {
	return &clientState{}
}

func (s *clientState) AddPendingRequest(requestID string, req ocpp.Request) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if requestID != "" && s.requestID == "" {
		s.requestID = requestID
		s.pendingRequest = pendingRequest{
			request: req,
		}
	}
}

func (s *clientState) GetPendingRequest(requestID string) (ocpp.Request, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.requestID != requestID {
		return nil, false
	}
	return s.pendingRequest.request, true
}

func (s *clientState) DeletePendingRequest(requestID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.requestID != requestID {
		return
	}
	s.requestID = ""
}

func (s *clientState) ClearPendingRequests() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.requestID = ""
}

func (s *clientState) HasPendingRequest() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.requestID != ""
}

// Contains the pending request state for messages associated to all client-server channels.
// It is used to separate endpoint logic from state management.
type ServerState interface {
	// Sets a Request as pending on the endpoint, for a specific client.
	// Requests are considered pending until a response was received.
	// The function expects a client ID, a unique message ID and the Request itself.
	// If an element with the same clientID/requestID exists, the new one will be ignored.
	AddPendingRequest(clientID string, requestID string, req ocpp.Request)
	// Deletes a pending Request from the endpoint, for a specific client, using the message ID.
	// If no such message is currently stored as pending, the call has no effect.
	DeletePendingRequest(clientID string, requestID string)
	// Retrieves a ClientState object, associated to a specific client.
	// If no such state exists, an empty state is returned.
	GetClientState(clientID string) ClientState
	// Returns true if there currently are pending requests for a client, false otherwise.
	HasPendingRequest(clientID string) bool
	// Returns true if there currently is at least one pending request, false otherwise.
	HasPendingRequests() bool
	// Clears currently pending requests for a client. Any confirmation/error,
	// received as a response to a cleared request, will be ignored.
	ClearClientPendingRequest(clientID string)
	// Clears all currently pending requests inside the map. Any confirmation/error,
	// received as a response to a previously sent request, will be ignored.
	//
	// Does not perform a deep deletion; is references to client state objects
	// are stored elsewhere, those will remain unaffected and become invalid.
	ClearAllPendingRequests()
}

// --------------------------------
// Request State Map implementation
// --------------------------------

// Simple implementation of ServerState, using a map.
// Supports any amount of clients and stores the pending requests for each client in a
// clientState struct.
//
// Client data is not deleted automatically; it should be deleted after a client session has ended.
//
// May internally use a mutex for concurrent access to the data struct.
// See NewServerState for more info.
type serverState struct {
	pendingRequestState map[string]ClientState
	mutex               *sync.RWMutex
}

// Creates a simple struct implementing ServerState, to be used by server dispatchers.
//
// If no mutex is passed, then atomic access to the data struct is not guaranteed, and race conditions may arise.
func NewServerState(m *sync.RWMutex) ServerState {
	return &serverState{
		pendingRequestState: map[string]ClientState{},
		mutex:               m,
	}
}

func (d *serverState) AddPendingRequest(clientID string, requestID string, req ocpp.Request) {
	if d.mutex != nil {
		d.mutex.Lock()
		defer d.mutex.Unlock()
	}
	state := d.getOrCreateState(clientID)
	state.AddPendingRequest(requestID, req)
}

func (d *serverState) DeletePendingRequest(clientID string, requestID string) {
	if d.mutex != nil {
		d.mutex.Lock()
		defer d.mutex.Unlock()
	}
	state, exists := d.pendingRequestState[clientID]
	if !exists {
		return
	}
	state.DeletePendingRequest(requestID)
}

func (d *serverState) GetClientState(clientID string) ClientState {
	if d.mutex != nil {
		d.mutex.Lock()
		defer d.mutex.Unlock()
	}
	return d.getOrCreateState(clientID)
}

func (d *serverState) HasPendingRequest(clientID string) bool {
	if d.mutex != nil {
		d.mutex.Lock()
		defer d.mutex.Unlock()
	}
	state, exists := d.pendingRequestState[clientID]
	return exists && state.HasPendingRequest()
}

func (d *serverState) HasPendingRequests() bool {
	if d.mutex != nil {
		d.mutex.Lock()
		defer d.mutex.Unlock()
	}
	for _, s := range d.pendingRequestState {
		if s.HasPendingRequest() {
			return true
		}
	}
	return false
}

func (d *serverState) ClearClientPendingRequest(clientID string) {
	if d.mutex != nil {
		d.mutex.Lock()
		defer d.mutex.Unlock()
	}
	delete(d.pendingRequestState, clientID)
}

func (d *serverState) ClearAllPendingRequests() {
	if d.mutex != nil {
		d.mutex.Lock()
		defer d.mutex.Unlock()
	}
	d.pendingRequestState = map[string]ClientState{}
}

func (d *serverState) getOrCreateState(clientID string) ClientState {
	state, exists := d.pendingRequestState[clientID]
	if !exists {
		state = NewClientState()
		d.pendingRequestState[clientID] = state
	}
	return state
}
