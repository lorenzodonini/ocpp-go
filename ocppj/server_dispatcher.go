package ocppj

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ws"
)

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
	// Sets the maximum timeout to be considered after sending a request.
	// If a response to the request is not received within the specified period, the request
	// is discarded and an error is returned to the caller.
	//
	// One timeout per client runs in the background.
	// The timeout is reset whenever a response comes in, the connection is closed, or the server is stopped.
	//
	// This function must be called before starting the dispatcher, otherwise it may lead to unexpected behavior.
	SetTimeout(timeout time.Duration)
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
	// The callback passes the original client ID, message ID, and request struct of the failed request,
	// along with an error.
	//
	// Calling Stop on the dispatcher will not trigger this callback.
	//
	// If no callback is set, a request will still be removed from the dispatcher when a timeout occurs.
	SetOnRequestCanceled(cb CanceledRequestHandler)
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
	timeout             time.Duration
	timerC              chan string
	running             bool
	stoppedC            chan struct{}
	onRequestCancel     CanceledRequestHandler
	network             ws.WsServer
	mutex               sync.RWMutex
}

// Handler function to be invoked when a request gets canceled (either due to timeout or to other external factors).
type CanceledRequestHandler func(clientID string, requestID string, request ocpp.Request, err *ocpp.Error)

// Utility struct for passing a client context around and cancel pending requests.
type clientTimeoutContext struct {
	ctx    context.Context
	cancel func()
}

func (c clientTimeoutContext) isActive() bool {
	return c.cancel != nil
}

// NewDefaultServerDispatcher creates a new DefaultServerDispatcher struct.
func NewDefaultServerDispatcher(queueMap ServerQueueMap) *DefaultServerDispatcher {
	d := &DefaultServerDispatcher{
		queueMap:         queueMap,
		requestChannel:   nil,
		readyForDispatch: make(chan string, 1),
		timeout:          defaultMessageTimeout,
	}
	d.pendingRequestState = NewServerState(&d.mutex)
	return d
}

func (d *DefaultServerDispatcher) Start() {
	d.requestChannel = make(chan string, 20)
	d.timerC = make(chan string, 10)
	d.stoppedC = make(chan struct{}, 1)
	d.running = true
	go d.messagePump()
}

func (d *DefaultServerDispatcher) IsRunning() bool {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.running
}

func (d *DefaultServerDispatcher) Stop() {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.running = false
	close(d.stoppedC)
}

func (d *DefaultServerDispatcher) SetTimeout(timeout time.Duration) {
	d.timeout = timeout
}

func (d *DefaultServerDispatcher) CreateClient(clientID string) {
	if d.IsRunning() {
		_ = d.queueMap.GetOrCreate(clientID)
	}
}

func (d *DefaultServerDispatcher) DeleteClient(clientID string) {
	d.queueMap.Remove(clientID)
	if d.IsRunning() {
		d.requestChannel <- clientID
	}
}

func (d *DefaultServerDispatcher) SetNetworkServer(server ws.WsServer) {
	d.network = server
}

func (d *DefaultServerDispatcher) SetOnRequestCanceled(cb CanceledRequestHandler) {
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
	var clientCtx clientTimeoutContext
	var clientQueue RequestQueue
	clientContextMap := map[string]clientTimeoutContext{} // Empty at the beginning
	// Dispatcher Loop
	for {
		select {
		case <-d.stoppedC:
			// Server was stopped
			d.queueMap.Init()
			log.Info("stopped processing requests")
			return
		case clientID = <-d.requestChannel:
			// Check whether there is a request queue for the specified client
			clientQueue, ok = d.queueMap.Get(clientID)
			if !ok {
				// No client queue found (client was removed)
				// Deleting and canceling the context
				clientCtx = clientContextMap[clientID]
				delete(clientContextMap, clientID)
				if clientCtx.ctx != nil {
					clientCtx.cancel()
				}
				continue
			}
			// Check whether we can transmit to client
			clientCtx, ok = clientContextMap[clientID]
			if !ok {
				// First request for this client, ready to transmit
				rdy = true
			} else {
				// If there is no active context, the client is ready to transmit
				rdy = !clientCtx.isActive()
			}
		case clientID, ok = <-d.timerC:
			// Timeout elapsed
			if !ok {
				continue
			}
			// Canceling timeout context
			log.Debugf("timeout for client %v, canceling message", clientID)
			clientCtx = clientContextMap[clientID]
			if clientCtx.isActive() {
				clientCtx.cancel()
				clientContextMap[clientID] = clientTimeoutContext{}
			}
			if d.pendingRequestState.HasPendingRequest(clientID) {
				// Current request for client timed out. Removing request and triggering cancel callback
				q, _ := d.queueMap.Get(clientID)
				bundle, _ := q.Peek().(RequestBundle)
				d.CompleteRequest(clientID, bundle.Call.UniqueId)
				log.Infof("request %v for %v timed out", bundle.Call.UniqueId, clientID)
				if d.onRequestCancel != nil {
					d.onRequestCancel(clientID, bundle.Call.UniqueId, bundle.Call.Payload,
						ocpp.NewError(GenericError, "Request timed out", bundle.Call.UniqueId))
				}
			}
		case clientID = <-d.readyForDispatch:
			// Cancel previous timeout (if any)
			clientCtx, ok = clientContextMap[clientID]
			if ok && clientCtx.isActive() {
				clientCtx.cancel()
				clientContextMap[clientID] = clientTimeoutContext{}
			}
			// Client can now transmit again
			clientQueue, ok = d.queueMap.Get(clientID)
			if ok {
				// Ready to transmit
				rdy = true
			}
			log.Debugf("%v ready to transmit again", clientID)
		}
		// Only dispatch request if able to send and request queue isn't empty
		if rdy && clientQueue != nil && !clientQueue.IsEmpty() {
			// Send request & set new context
			clientCtx = d.dispatchNextRequest(clientID)
			clientContextMap[clientID] = clientCtx
			if clientCtx.isActive() {
				go d.waitForTimeout(clientID, clientCtx)
			}
			// Update ready state
			rdy = false
		}
	}
}

func (d *DefaultServerDispatcher) dispatchNextRequest(clientID string) (clientCtx clientTimeoutContext) {
	// Get first element in queue
	q, ok := d.queueMap.Get(clientID)
	if !ok {
		log.Errorf("failed to dispatch next request for %s, no request queue available", clientID)
		return
	}
	el := q.Peek()
	bundle, _ := el.(RequestBundle)
	// Convert message to JSON
	jsonMessage, err := json.Marshal(bundle.Call)
	if err != nil {
		// Cancel request internally
		d.CompleteRequest(clientID, bundle.Call.GetUniqueId())
		if d.onRequestCancel != nil {
			d.onRequestCancel(clientID, bundle.Call.GetUniqueId(), bundle.Call.Payload, ocpp.NewError(InternalError, err.Error(), bundle.Call.GetUniqueId()))
		}
		return
	}
	callID := bundle.Call.GetUniqueId()
	d.pendingRequestState.AddPendingRequest(clientID, callID, bundle.Call.Payload)
	err = d.network.Write(clientID, jsonMessage)
	if err != nil {
		log.Errorf("error while sending message: %v", err)
		//TODO: handle retransmission instead of removing pending request
		d.CompleteRequest(clientID, callID)
		if d.onRequestCancel != nil {
			d.onRequestCancel(clientID, bundle.Call.UniqueId, bundle.Call.Payload,
				ocpp.NewError(InternalError, err.Error(), bundle.Call.UniqueId))
		}
		return
	}
	// Create and return context (only if timeout is set)
	if d.timeout > 0 {
		ctx, cancel := context.WithTimeout(context.TODO(), d.timeout)
		clientCtx = clientTimeoutContext{ctx: ctx, cancel: cancel}
	}
	log.Infof("dispatched request %s for %s", callID, clientID)
	return
}

func (d *DefaultServerDispatcher) waitForTimeout(clientID string, clientCtx clientTimeoutContext) {
	defer clientCtx.cancel()
	log.Debugf("started timeout timer for %s", clientID)
	select {
	case <-clientCtx.ctx.Done():
		err := clientCtx.ctx.Err()
		if err == context.DeadlineExceeded {
			// Timeout triggered, notifying messagePump
			d.mutex.RLock()
			defer d.mutex.RUnlock()
			if d.running {
				d.timerC <- clientID
			}
		} else {
			log.Debugf("timeout canceled for %s", clientID)
		}
	case <-d.stoppedC:
		// Server was stopped, every pending timeout gets canceled
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
		log.Errorf("internal state mismatch: processing response for %v but expected response for %v", requestID, callID)
		return
	}
	q.Pop()
	d.pendingRequestState.DeletePendingRequest(clientID, requestID)
	log.Debugf("completed request %s for %s", callID, clientID)
	// Signal that next message in queue may be sent
	d.readyForDispatch <- clientID
}
