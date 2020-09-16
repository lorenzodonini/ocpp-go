package ocppj

import (
	"errors"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ws"
	log "github.com/sirupsen/logrus"
	"sync"
)

// Contains the pending request state for messages.
// It is used to separate endpoint logic from state management.
type PendingRequestState interface {
	// Sets a Request as pending on the endpoint. Requests are considered pending until a response was received.
	// The function expects a message unique ID and the Request.
	// If an element with the same requestID exists, it will be overwritten.
	AddPendingRequest(requestID string, req ocpp.Request)
	// Retrieves a pending Request, using the message ID.
	// If no request for the passed message ID is found, a false flag is returned.
	GetPendingRequest(requestID string) (ocpp.Request, bool)
	// Deletes a pending Request from the endpoint, using the message ID.
	DeletePendingRequest(requestID string)
	// Clears all currently pending requests. Any confirmation/error,
	// received as a response to a previously sent request, will be ignored.
	ClearPendingRequests()
}

// ClientDispatcher contains the state and logic for handling outgoing messages on a client endpoint.
// This allows the ocpp-j layer to delegate queueing and processing logic to an external entity.
//
// The dispatcher writes outgoing messages directly to the networking layer, using a previously set websocket client.
//
// A PendingRequestState needs to be passed to the dispatcher, before starting it.
// The dispatcher is in charge of managing pending requests while managing the request flow.
type ClientDispatcher interface {
	Start()
	IsRunning() bool
	SendRequest(req interface{}) error
	CompleteRequest(requestID string)
	SetNetworkClient(client ws.WsClient)
	SetPendingRequestState(stateHandler PendingRequestState)
	Stop()
}

// DefaultClientDispatcher is a default implementation of the ClientDispatcher interface.
//
// The dispatcher implements the PendingRequestState as well for simplicity.
// Access to pending requests is thread-safe.
type DefaultClientDispatcher struct {
	requestQueue     RequestQueue
	requestChannel   chan bool
	readyForDispatch chan bool
	pendingRequests  map[string]ocpp.Request
	network          ws.WsClient
	mutex            sync.Mutex
}

// NewDefaultClientDispatcher creates a new DefaultClientDispatcher struct.
func NewDefaultClientDispatcher(queue RequestQueue) *DefaultClientDispatcher {
	return &DefaultClientDispatcher{
		requestQueue:     queue,
		requestChannel:   nil,
		readyForDispatch: make(chan bool, 1),
		pendingRequests:  map[string]ocpp.Request{},
	}
}

func (d *DefaultClientDispatcher) Start() {
	d.requestChannel = make(chan bool, 1)
	go d.messagePump()
}

func (d *DefaultClientDispatcher) IsRunning() bool {
	return d.requestChannel != nil
}

func (d *DefaultClientDispatcher) Stop() {
	close(d.requestChannel)
	// TODO: clear pending requests?
}

func (d *DefaultClientDispatcher) SetNetworkClient(client ws.WsClient) {
	d.network = client
}

func (d *DefaultClientDispatcher) SetPendingRequestState(_ PendingRequestState) {}

func (d *DefaultClientDispatcher) AddPendingRequest(requestID string, req ocpp.Request) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.pendingRequests[requestID] = req
}

func (d *DefaultClientDispatcher) DeletePendingRequest(requestID string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	delete(d.pendingRequests, requestID)
}

func (d *DefaultClientDispatcher) ClearPendingRequests() {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.pendingRequests = map[string]ocpp.Request{}
}

func (d *DefaultClientDispatcher) GetPendingRequest(requestID string) (ocpp.Request, bool) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	req, ok := d.pendingRequests[requestID]
	return req, ok
}

func (d *DefaultClientDispatcher) SendRequest(req interface{}) error {
	if d.network == nil {
		return errors.New("cannot SendRequest, no network client was set")
	}
	if err := d.requestQueue.Push(req); err != nil {
		return err
	}
	d.requestChannel <- true
	return nil
}

func (d *DefaultClientDispatcher) messagePump() {
	rdy := true // Ready to transmit at the beginning
	for {
		select {
		case _, ok := <-d.requestChannel:
			// New request was posted
			if !ok {
				log.Infof("stopped processing requests")
				d.requestQueue.Init()
				d.requestChannel = nil
				return
			}
		case rdy = <-d.readyForDispatch:
		}
		// Only dispatch request if able to send and request queue isn't empty
		if rdy && !d.requestQueue.IsEmpty() {
			d.dispatchNextRequest()
			rdy = false
		}
	}
}

func (d *DefaultClientDispatcher) dispatchNextRequest() {
	// Get first element in queue
	el := d.requestQueue.Peek()
	bundle, _ := el.(RequestBundle)
	jsonMessage := bundle.Data
	d.AddPendingRequest(bundle.Call.UniqueId, bundle.Call.Payload)

	err := d.network.Write(jsonMessage)
	if err != nil {
		log.Errorf("error while sending message: %v", err)
		//TODO: handle retransmission instead of removing pending request
		d.DeletePendingRequest(bundle.Call.GetUniqueId())
		d.CompleteRequest(bundle.Call.GetUniqueId())
		// TODO: throw error?
	} else {
		// Transmitted correctly
		log.Debugf("sent request %v: %v", bundle.Call.UniqueId, string(jsonMessage))
	}
}

func (d *DefaultClientDispatcher) CompleteRequest(requestId string) {
	el := d.requestQueue.Peek()
	if el == nil {
		log.Errorf("attempting to pop front of queue, but queue is empty")
		return
	}
	bundle, _ := el.(RequestBundle)
	if bundle.Call.UniqueId != requestId {
		log.Fatalf("internal state mismatch: received response for %v but expected response for %v", requestId, bundle.Call.UniqueId)
		return
	}
	d.requestQueue.Pop()
	d.DeletePendingRequest(requestId)
	log.Debugf("removed request %v from front of queue", bundle.Call.UniqueId)
	// Signal that next message in queue may be sent
	d.readyForDispatch <- true
}

// ServerDispatcher contains the state and logic for handling outgoing messages on a server endpoint.
// This allows the ocpp-j layer to delegate queueing and processing logic to an external entity.
//
// The dispatcher writes outgoing messages directly to the networking layer, using a previously set websocket server.
//
// A PendingRequestState needs to be passed to the dispatcher, before starting it.
// The dispatcher is in charge of managing all pending requests to clients, while managing the request flow.
type ServerDispatcher interface {
	Start()
	IsRunning() bool
	SendRequest(clientID string, req RequestBundle) error
	CompleteRequest(clientID string, requestID string)
	SetNetworkServer(server ws.WsServer)
	SetPendingRequestState(stateHandler PendingRequestState)
	Stop()
}

// DefaultServerDispatcher is a default implementation of the ServerDispatcher interface.
//
// The dispatcher implements the PendingRequestState as well for simplicity.
// Access to pending requests is thread-safe.
type DefaultServerDispatcher struct {
	queueMap         ServerQueueMap
	requestChannel   chan string
	readyForDispatch chan string
	pendingRequests  map[string]ocpp.Request
	network          ws.WsServer
	mutex            sync.Mutex
}

// NewDefaultServerDispatcher creates a new DefaultServerDispatcher struct.
func NewDefaultServerDispatcher(queueMap ServerQueueMap) *DefaultServerDispatcher {
	return &DefaultServerDispatcher{
		queueMap:         queueMap,
		requestChannel:   nil,
		readyForDispatch: make(chan string, 1),
		pendingRequests:  map[string]ocpp.Request{},
	}
}

func (d *DefaultServerDispatcher) Start() {
	d.requestChannel = make(chan string, 1)
	go d.messagePump()
}

func (d *DefaultServerDispatcher) IsRunning() bool {
	return d.requestChannel != nil
}

func (d *DefaultServerDispatcher) Stop() {
	close(d.requestChannel)
	// TODO: clear pending requests?
}

func (d *DefaultServerDispatcher) SetNetworkServer(server ws.WsServer) {
	d.network = server
}

func (d *DefaultServerDispatcher) SetPendingRequestState(_ PendingRequestState) {}

func (d *DefaultServerDispatcher) AddPendingRequest(requestID string, req ocpp.Request) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.pendingRequests[requestID] = req
}

func (d *DefaultServerDispatcher) DeletePendingRequest(requestID string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	delete(d.pendingRequests, requestID)
}

func (d *DefaultServerDispatcher) ClearPendingRequests() {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.pendingRequests = map[string]ocpp.Request{}
}

func (d *DefaultServerDispatcher) GetPendingRequest(requestID string) (ocpp.Request, bool) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	req, ok := d.pendingRequests[requestID]
	return req, ok
}

func (d *DefaultServerDispatcher) SendRequest(clientID string, req RequestBundle) error {
	if d.network == nil {
		return errors.New("cannot SendRequest, no network server was set")
	}
	q := d.queueMap.GetOrCreate(clientID)
	if err := q.Push(req); err != nil {
		return err
	}
	d.requestChannel <- clientID
	return nil
}

// requestPump processes new outgoing requests for each client and makes sure they are processed sequentially.
// This method is executed by a dedicated coroutine as soon as the server is started and runs indefinitely.
func (d *DefaultServerDispatcher) messagePump() {
	var clientID string
	var ok bool
	var rdy bool
	var clientQueue RequestQueue
	clientReadyMap := map[string]bool{} // Empty at the beginning
	for {
		select {
		case clientID, ok = <-d.requestChannel:
			// Check if channel was closed
			if !ok {
				log.Infof("stopped processing requests")
				// TODO: clear queue map?
				d.requestChannel = nil
				return
			}
			clientQueue, ok = d.queueMap.Get(clientID)
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
		case clientID = <-d.readyForDispatch:
			// Client can now transmit again
			rdy = true
			clientReadyMap[clientID] = rdy
			clientQueue, _ = d.queueMap.Get(clientID)
		}
		// Only dispatch request if able to send and request queue isn't empty
		if rdy && !clientQueue.IsEmpty() {
			d.dispatchNextRequest(clientID)
			// Update ready state
			rdy = false
			clientReadyMap[clientID] = rdy
		}
	}
}

func (d *DefaultServerDispatcher) dispatchNextRequest(clientID string) {
	// Get first element in queue
	q, _ := d.queueMap.Get(clientID)
	el := q.Peek()
	bundle, _ := el.(RequestBundle)
	jsonMessage := bundle.Data
	callID := bundle.Call.GetUniqueId()
	d.AddPendingRequest(callID, bundle.Call.Payload)
	err := d.network.Write(clientID, jsonMessage)
	if err != nil {
		log.Errorf("error while sending message: %v", err)
		//TODO: handle retransmission instead of removing pending request
		d.DeletePendingRequest(bundle.Call.GetUniqueId())
		d.CompleteRequest(clientID, bundle.Call.GetUniqueId())
		// TODO: throw error?
		//if s.errorHandler != nil {
		//	s.errorHandler(clientID, ocpp.NewError(GenericError, err.Error(), callID), err)
		//}
	} else {
		// Transmitted correctly
		log.Debugf("sent request %v to client %v: %v", callID, clientID, string(jsonMessage))
	}
}

func (d *DefaultServerDispatcher) CompleteRequest(clientID string, requestID string) {
	q, ok := d.queueMap.Get(clientID)
	if !ok {
		log.Errorf("attempting to complete request for client %v, but no matching queue found", clientID)
		return
	}
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
	d.DeletePendingRequest(requestID)
	log.Debugf("removed request %v from front of queue", callID)
	// Signal that next message in queue may be sent
	d.readyForDispatch <- clientID
}
