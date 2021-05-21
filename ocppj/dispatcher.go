package ocppj

import (
	"fmt"
	"sync"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ws"
)

// ClientDispatcher contains the state and logic for handling outgoing messages on a client endpoint.
// This allows the ocpp-j layer to delegate queueing and processing logic to an external entity.
//
// The dispatcher writes outgoing messages directly to the networking layer, using a previously set websocket client.
//
// A ClientState needs to be passed to the dispatcher, before starting it.
// The dispatcher is in charge of managing pending requests while handling the request flow.
type ClientDispatcher interface {
	// Starts the dispatcher. Depending on the implementation, this may
	// start a dedicated goroutine or simply allocate the necessary state.
	Start()
	// Sets the maximum timeout to be considered after sending a request.
	// If a response to the request is not received within the specified period, the request
	// is discarded and an error is returned to the caller.
	//
	// The timeout is reset upon a disconnection/reconnection.
	//
	// This function must be called before starting the dispatcher, otherwise it may lead to unexpected behavior.
	SetTimeout(timeout time.Duration)
	// Returns true, if the dispatcher is currently running, false otherwise.
	// If the dispatcher is paused, the function still returns true.
	IsRunning() bool
	// Returns true, if the dispatcher is currently paused, false otherwise.
	// If the dispatcher is not running at all, the function will still return false.
	IsPaused() bool
	// Dispatches a request. Depending on the implementation, this may first queue a request
	// and process it later, asynchronously, or write it directly to the networking layer.
	//
	// If no network client was set, or the request couldn't be processed, an error is returned.
	SendRequest(req interface{}) error
	// Notifies the dispatcher that a request has been completed (i.e. a response was received).
	// The dispatcher takes care of removing the request marked by the requestID from
	// the pending requests. It will then attempt to process the next queued request.
	CompleteRequest(requestID string)
	// Sets a callback to be invoked when a request gets canceled, due to network timeouts.
	// The callback passes the original message ID, feature name and request struct of the failed request.
	//
	// Calling Stop on the dispatcher will not trigger this callback.
	//
	// If no callback is set, a request will still be removed from the dispatcher when a timeout occurs.
	SetOnRequestCanceled(cb func(string, string, ocpp.Request))
	// Sets the network client, so the dispatcher may send requests using the networking layer directly.
	//
	// This needs to be set before calling the Start method. If not, sending requests will fail.
	SetNetworkClient(client ws.WsClient)
	// Sets the state manager for pending requests in the dispatcher.
	//
	// The state should only be accessed by the dispatcher while running.
	SetPendingRequestState(stateHandler ClientState)
	// Stops a running dispatcher. This will clear all state and empty the internal queues.
	//
	// If an onRequestCanceled callback is set, it won't be triggered by stopping the dispatcher.
	Stop()
	// Notifies that an external event (typically network-related) should pause
	// the dispatcher. Internal timers will be stopped an no further requests
	// will be set to pending. You may keep enqueuing requests.
	// Use the Resume method for re-starting the dispatcher.
	Pause()
	// Undoes a previous pause operation, restarting internal timers and the
	// regular request flow.
	//
	// If there was a pending request before pausing the dispatcher, a response/timeout
	// for this request shall be awaited anew.
	Resume()
}

// pendingRequest is used internally for associating metadata to a pending Request.
type pendingRequest struct {
	request   ocpp.Request
	startTime time.Time
}

// DefaultClientDispatcher is a default implementation of the ClientDispatcher interface.
//
// The dispatcher implements the ClientState as well for simplicity.
// Access to pending requests is thread-safe.
type DefaultClientDispatcher struct {
	requestQueue        RequestQueue
	requestChannel      chan bool
	readyForDispatch    chan bool
	pendingRequestState ClientState
	network             ws.WsClient
	mutex               sync.RWMutex
	onRequestCancel     func(string, string, ocpp.Request)
	timer               *time.Timer
	paused              bool
	timeout             time.Duration
}

const defaultTimeoutTick = 24 * time.Hour
const defaultMessageTimeout = 30 * time.Second

// NewDefaultClientDispatcher creates a new DefaultClientDispatcher struct.
func NewDefaultClientDispatcher(queue RequestQueue) *DefaultClientDispatcher {
	return &DefaultClientDispatcher{
		requestQueue:        queue,
		requestChannel:      nil,
		readyForDispatch:    make(chan bool, 1),
		pendingRequestState: NewClientState(),
		timeout:             defaultMessageTimeout,
	}
}

func (d *DefaultClientDispatcher) SetOnRequestCanceled(cb func(string, string, ocpp.Request)) {
	d.onRequestCancel = cb
}

func (d *DefaultClientDispatcher) SetTimeout(timeout time.Duration) {
	d.timeout = timeout
}

func (d *DefaultClientDispatcher) Start() {
	d.requestChannel = make(chan bool, 1)
	d.timer = time.NewTimer(defaultTimeoutTick) // Default to 24 hours tick
	go d.messagePump()
}

func (d *DefaultClientDispatcher) IsRunning() bool {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	return d.requestChannel != nil
}

func (d *DefaultClientDispatcher) IsPaused() bool {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	return d.paused
}

func (d *DefaultClientDispatcher) Stop() {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	close(d.requestChannel)
	// TODO: clear pending requests?
}

func (d *DefaultClientDispatcher) SetNetworkClient(client ws.WsClient) {
	d.network = client
}

func (d *DefaultClientDispatcher) SetPendingRequestState(state ClientState) {
	d.pendingRequestState = state
}

func (d *DefaultClientDispatcher) SendRequest(req interface{}) error {
	if d.network == nil {
		return fmt.Errorf("cannot SendRequest, no network client was set")
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
				d.requestQueue.Init()
				d.requestChannel = nil
				return
			}
		case _, ok := <-d.timer.C:
			// Timeout elapsed
			if !ok {
				continue
			}
			if d.pendingRequestState.HasPendingRequest() {
				// Current request timed out. Removing request and triggering cancel callback
				el := d.requestQueue.Peek()
				bundle, _ := el.(RequestBundle)
				d.CompleteRequest(bundle.Call.UniqueId)
				if d.onRequestCancel != nil {
					d.onRequestCancel(bundle.Call.UniqueId, bundle.Call.Action, bundle.Call.Payload)
				}
			}
			// No request is currently pending -> set timer to high number
			d.timer.Reset(defaultTimeoutTick)
		case rdy = <-d.readyForDispatch:
			// Ready flag set, keep going
		}
		// Check if dispatcher is paused
		d.mutex.Lock()
		paused := d.paused
		d.mutex.Unlock()
		if paused {
			// Ignore dispatch events as long as dispatcher is paused
			continue
		}
		// Only dispatch request if able to send and request queue isn't empty
		if rdy && !d.requestQueue.IsEmpty() {
			d.dispatchNextRequest()
			rdy = false
			// Set timer
			if !d.timer.Stop() {
				<-d.timer.C
			}
			d.timer.Reset(d.timeout)
		}
	}
}

func (d *DefaultClientDispatcher) dispatchNextRequest() {
	// Get first element in queue
	el := d.requestQueue.Peek()
	bundle, _ := el.(RequestBundle)
	jsonMessage := bundle.Data
	d.pendingRequestState.AddPendingRequest(bundle.Call.UniqueId, bundle.Call.Payload)
	// Attempt to send over network
	err := d.network.Write(jsonMessage)
	if err != nil {
		//TODO: handle retransmission instead of skipping request altogether
		d.CompleteRequest(bundle.Call.GetUniqueId())
		if d.onRequestCancel != nil {
			d.onRequestCancel(bundle.Call.UniqueId, bundle.Call.Action, bundle.Call.Payload)
		}
	}
}

func (d *DefaultClientDispatcher) Pause() {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if !d.timer.Stop() {
		<-d.timer.C
	}
	d.timer.Reset(defaultTimeoutTick)
	d.paused = true
}

func (d *DefaultClientDispatcher) Resume() {
	d.mutex.Lock()
	d.paused = false
	d.mutex.Unlock()
	if d.pendingRequestState.HasPendingRequest() {
		// There is a pending request already. Awaiting response, before dispatching new requests.
		d.timer.Reset(d.timeout)
	} else {
		// Can dispatch a new request. Notifying message pump.
		d.readyForDispatch <- true
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
		log.Errorf("internal state mismatch: received response for %v but expected response for %v", requestId, bundle.Call.UniqueId)
		return
	}
	d.requestQueue.Pop()
	d.pendingRequestState.DeletePendingRequest(requestId)
	log.Debugf("removed request %v from front of queue", bundle.Call.UniqueId)
	// Signal that next message in queue may be sent
	d.readyForDispatch <- true
}

// ServerDispatcher contains the state and logic for handling outgoing messages on a server endpoint.
// This allows the ocpp-j layer to delegate queueing and processing logic to an external entity.
//
// The dispatcher writes outgoing messages directly to the networking layer, using a previously set websocket server.
//
// A ClientState needs to be passed to the dispatcher, before starting it.
// The dispatcher is in charge of managing all pending requests to clients, while handling the request flow.
type ServerDispatcher interface {
	// Starts the dispatcher. Depending on the implementation, this may
	// start a dedicated goroutine or simply allocate the necessary state.
	Start()
	// Returns true, if the dispatcher is currently running, false otherwise.
	// If the dispatcher is paused, the function still returns true.
	IsRunning() bool
	// Dispatches a request for a specific client. Depending on the implementation, this may first queue
	// a request and process it later (asynchronously), or write it directly to the networking layer.
	//
	// If no network server was set, or the request couldn't be processed, an error is returned.
	SendRequest(clientID string, req RequestBundle) error
	// Notifies the dispatcher that a request has been completed (i.e. a response was received),
	// for a specific client.
	// The dispatcher takes care of removing the request marked by the requestID from
	// that client's pending requests. It will then attempt to process the next queued request.
	CompleteRequest(clientID string, requestID string)
	// Sets a callback to be invoked when a request gets canceled, due to network timeouts.
	// The callback passes the original client ID, message ID, feature name and request struct of the failed request.
	//
	// Calling Stop on the dispatcher will not trigger this callback.
	//
	// If no callback is set, a request will still be removed from the dispatcher when a timeout occurs.
	SetOnRequestCanceled(cb func(string, string, string, ocpp.Request))
	// Sets the network server, so the dispatcher may send requests using the networking layer directly.
	//
	// This needs to be set before calling the Start method. If not, sending requests will fail.
	SetNetworkServer(server ws.WsServer)
	// Sets the state manager for pending requests in the dispatcher.
	//
	// The state should only be accessed by the dispatcher while running.
	SetPendingRequestState(stateHandler ServerState)
	// Stops a running dispatcher. This will clear all state and empty the internal queues.
	//
	// If an onRequestCanceled callback is set, it won't be triggered by stopping the dispatcher.
	Stop()
	// Notifies that it is now possible to dispatch requests for a new client.
	//
	// Internal queues are created and requests for the client are now accepted.
	CreateClient(clientID string)
	// Notifies that a client was invalidated (typically caused by a network event).
	//
	// The dispatcher will stop dispatching requests for that specific client.
	// Internal queues for that client are cleared and no further requests will be accepted.
	// Undelivered pending requests are also cleared.
	// The OnRequestCanceled callback will be invoked for each discarded request.
	DeleteClient(clientID string)
}

// DefaultServerDispatcher is a default implementation of the ServerDispatcher interface.
//
// The dispatcher implements the ClientState as well for simplicity.
// Access to pending requests is thread-safe.
type DefaultServerDispatcher struct {
	queueMap            ServerQueueMap
	requestChannel      chan string
	readyForDispatch    chan string
	pendingRequestState ServerState
	onRequestCancel     func(string, string, string, ocpp.Request)
	network             ws.WsServer
	mutex               sync.RWMutex
	stopped             chan struct{}
}

// NewDefaultServerDispatcher creates a new DefaultServerDispatcher struct.
func NewDefaultServerDispatcher(queueMap ServerQueueMap) *DefaultServerDispatcher {
	d := &DefaultServerDispatcher{
		queueMap:         queueMap,
		requestChannel:   nil,
		readyForDispatch: make(chan string, 1),
	}
	d.pendingRequestState = NewServerState(&d.mutex)
	return d
}

func (d *DefaultServerDispatcher) Start() {
	d.requestChannel = make(chan string, 1)
	d.stopped = make(chan struct{})
	go d.messagePump()
}

func (d *DefaultServerDispatcher) IsRunning() bool {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	return d.requestChannel != nil
}

func (d *DefaultServerDispatcher) Stop() {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	close(d.requestChannel)
	// TODO: clear pending requests?
}

func (d *DefaultServerDispatcher) CreateClient(clientID string) {
	_ = d.queueMap.GetOrCreate(clientID)
}

func (d *DefaultServerDispatcher) DeleteClient(clientID string) {
	d.queueMap.Remove(clientID)
	select {
	case <-d.stopped:
	case d.requestChannel <- clientID:
	}
}

func (d *DefaultServerDispatcher) SetNetworkServer(server ws.WsServer) {
	d.network = server
}

func (d *DefaultServerDispatcher) SetOnRequestCanceled(cb func(string, string, string, ocpp.Request)) {
	d.onRequestCancel = cb
}

func (d *DefaultServerDispatcher) SetPendingRequestState(state ServerState) {
	d.pendingRequestState = state
}

func (d *DefaultServerDispatcher) SendRequest(clientID string, req RequestBundle) error {
	if d.network == nil {
		return fmt.Errorf("cannot send request %v, no network server was set", req.Call.UniqueId)
	}
	q, ok := d.queueMap.Get(clientID)
	if !ok {
		return fmt.Errorf("cannot send request %s, no client %s exists", req.Call.UniqueId, clientID)
	}
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
				d.queueMap.Init()
				d.requestChannel = nil
				close(d.stopped)
				log.Info("stopped processing requests")
				return
			}
			clientQueue, ok = d.queueMap.Get(clientID)
			// Check whether there is a request queue for the specified client
			if !ok {
				// No client queue found, deleting the ready flag
				delete(clientReadyMap, clientID)
				rdy = false
				break
			}
			// Check whether can transmit to client
			rdy, ok = clientReadyMap[clientID]
			if !ok {
				// First request for this client. Setting ready flag to true
				rdy = true
				clientReadyMap[clientID] = rdy
			}
			//TODO: check for response timeout
		case clientID = <-d.readyForDispatch:
			// Client can now transmit again
			clientQueue, rdy = d.queueMap.Get(clientID)
			if rdy {
				clientReadyMap[clientID] = rdy
			}
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
	q, ok := d.queueMap.Get(clientID)
	if !ok {
		log.Errorf("failed to dispatch next request for client %s, no request queue available", clientID)
		return
	}
	el := q.Peek()
	bundle, _ := el.(RequestBundle)
	jsonMessage := bundle.Data
	callID := bundle.Call.GetUniqueId()
	d.pendingRequestState.AddPendingRequest(clientID, callID, bundle.Call.Payload)
	err := d.network.Write(clientID, jsonMessage)
	if err != nil {
		log.Errorf("error while sending message: %v", err)
		//TODO: handle retransmission instead of removing pending request
		d.CompleteRequest(clientID, callID)
		if d.onRequestCancel != nil {
			d.onRequestCancel(clientID, callID, bundle.Call.Action, bundle.Call.Payload)
		}
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
		log.Errorf("internal state mismatch: received response for %v but expected response for %v", requestID, callID)
		return
	}
	q.Pop()
	d.pendingRequestState.DeletePendingRequest(clientID, requestID)
	log.Debugf("removed request %v from front of queue", callID)
	// Signal that next message in queue may be sent
	d.readyForDispatch <- clientID
}
