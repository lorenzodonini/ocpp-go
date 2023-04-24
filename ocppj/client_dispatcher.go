package ocppj

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp"
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
	SendRequest(req RequestBundle) error
	// Notifies the dispatcher that a request has been completed (i.e. a response was received).
	// The dispatcher takes care of removing the request marked by the requestID from
	// the pending requests. It will then attempt to process the next queued request.
	//
	// The original RequestBundle is returned after the request was marked as complete.
	// The caller may access the metadata and invoke the respective callback function.
	// If the original request is no longer in the queue, the return struct is empty.
	CompleteRequest(requestID string) RequestBundle
	// Sets a callback to be invoked when a request gets canceled, due to network timeouts or internal errors.
	// The callback passes the original message ID and request struct of the failed request, along with an error.
	//
	// Calling Stop on the dispatcher will not trigger this callback.
	//
	// If no callback is set, a request will still be removed from the dispatcher when a timeout occurs.
	SetOnRequestCanceled(func(request RequestBundle, err error))
	// Sets the function that allows the dispatcher to send requests without accessing the network layer directly.
	//
	// This needs to be set before calling the Start method. If not, sending requests will fail.
	SetNetworkSendHandler(func(data []byte) error)
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

// Used internally for sending control signals to the main message pump.
type controlSignal int

const (
	signalPause  controlSignal = 1
	signalResume controlSignal = 2
	signalStop   controlSignal = 3
)

// pendingRequest is used internally for associating metadata to a pending Request.
type pendingRequest struct {
	request ocpp.Request
}

// DefaultClientDispatcher is a default implementation of the ClientDispatcher interface.
//
// The dispatcher implements the ClientState as well for simplicity.
// Access to pending requests is thread-safe.
type DefaultClientDispatcher struct {
	requestQueue        RequestQueue
	requestChannel      chan struct{}
	controlChannel      chan controlSignal
	pendingRequestState ClientState
	networkSendHandler  func(data []byte) error
	mutex               sync.RWMutex
	onRequestCancel     func(bundle RequestBundle, err error)
	timeoutContext      context.Context
	timeoutCancelFn     context.CancelFunc
	paused              bool
	timeout             time.Duration
}

const defaultMessageTimeout = 30 * time.Second

// NewDefaultClientDispatcher creates a new DefaultClientDispatcher struct.
func NewDefaultClientDispatcher(queue RequestQueue) *DefaultClientDispatcher {
	return &DefaultClientDispatcher{
		requestQueue:        queue,
		requestChannel:      nil,
		pendingRequestState: NewClientState(),
		timeout:             defaultMessageTimeout,
	}
}

func (d *DefaultClientDispatcher) SetOnRequestCanceled(cb func(request RequestBundle, err error)) {
	d.onRequestCancel = cb
}

func (d *DefaultClientDispatcher) SetTimeout(timeout time.Duration) {
	d.timeout = timeout
}

func (d *DefaultClientDispatcher) Start() {
	d.requestChannel = make(chan struct{}, 1)
	d.controlChannel = make(chan controlSignal, 1)
	d.timeoutContext = context.TODO()
	d.timeoutCancelFn = func() {}
	go d.messagePump()
}

func (d *DefaultClientDispatcher) IsRunning() bool {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.requestChannel != nil
}

func (d *DefaultClientDispatcher) IsPaused() bool {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.paused
}

func (d *DefaultClientDispatcher) Stop() {
	d.controlChannel <- signalStop
	// TODO: clear pending requests?
}

func (d *DefaultClientDispatcher) SetNetworkSendHandler(handler func(data []byte) error) {
	d.networkSendHandler = handler
}

func (d *DefaultClientDispatcher) SetPendingRequestState(state ClientState) {
	d.pendingRequestState = state
}

func (d *DefaultClientDispatcher) SendRequest(req RequestBundle) error {
	if d.networkSendHandler == nil {
		return fmt.Errorf("cannot SendRequest, no network sending function was set")
	}
	if err := d.requestQueue.Push(req); err != nil {
		return err
	}
	d.requestChannel <- struct{}{}
	return nil
}

func (d *DefaultClientDispatcher) messagePump() {
	for {
		select {
		case signal, _ := <-d.controlChannel:
			if !d.onControlSignal(signal) {
				continue
			}
		case _, ok := <-d.requestChannel:
			if !ok {
				// Dispatcher was stopped. Clearing request queue and resetting state.
				d.requestQueue.Init()
				d.requestChannel = nil
				// TODO: clear pending requests?
				return
			} else {
				// New request was posted, continue
			}
		case <-d.timeoutContext.Done():
			// Request interrupted
			if d.pendingRequestState.HasPendingRequest() {
				// Get current request
				el := d.requestQueue.Peek()
				bundle, _ := el.(RequestBundle)
				// Mark request as completed
				d.CompleteRequest(bundle.Call.UniqueId)
				// Check which errors occurred
				ctxErr := d.timeoutContext.Err()
				if ctxErr == context.DeadlineExceeded {
					// Current request timed out. Notifying upper layer.
					if d.onRequestCancel != nil {
						d.onRequestCancel(bundle, ocpp.NewError(GenericError, "Request timed out", bundle.Call.GetUniqueId()))
					}
				} else if ctxErr == context.Canceled {
					// Current request canceled by user. Notifying upper layer.
					if d.onRequestCancel != nil {
						d.onRequestCancel(bundle, ocpp.NewError(GenericError, "Request canceled by user", bundle.Call.GetUniqueId()))
					}
				}
			}
			// Context has outlived its purpose, reset
			d.timeoutCancelFn()
			d.timeoutContext = context.TODO()
		}

		if d.IsPaused() {
			// Ignore dispatch events as long as dispatcher is paused
			continue
		}
		// Only dispatch a request if the request queue isn't empty and there is no pending request already
		if !d.requestQueue.IsEmpty() && !d.pendingRequestState.HasPendingRequest() {
			err := d.dispatchNextRequest()
			if err != nil {
				// Network error while sending request. Request will be retried later.
				log.Error(err)
			}
			if err == nil {
				d.startRequestTimeout()
			}
		}
	}
}

func (d *DefaultClientDispatcher) onControlSignal(signal controlSignal) (ready bool) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	switch signal {
	case signalPause:
		d.paused = true
		// Reset timeout
		d.timeoutCancelFn()
		d.timeoutContext = context.TODO()
	case signalResume:
		d.paused = false
		if d.pendingRequestState.HasPendingRequest() {
			// There is a pending request already. Awaiting response, before dispatching new requests.
			// The original timer is reset to its original duration for this purpose.
			d.startRequestTimeout()
		}
		// Message pump will attempt to dispatch a new request, if there are no pending requests
		return true
	case signalStop:
		close(d.requestChannel)
		// TODO: clear pending requests?
	}
	return false
}

func (d *DefaultClientDispatcher) dispatchNextRequest() error {
	// Get first element in queue
	el := d.requestQueue.Peek()
	bundle, _ := el.(RequestBundle)
	// Convert message to JSON
	jsonMessage, err := json.Marshal(bundle.Call)
	if err != nil {
		// Cancel request internally
		d.CompleteRequest(bundle.Call.GetUniqueId())
		if d.onRequestCancel != nil {
			d.onRequestCancel(bundle, ocpp.NewError(GenericError, err.Error(), bundle.Call.GetUniqueId()))
		}
		return err
	}
	d.pendingRequestState.AddPendingRequest(bundle.Call.UniqueId, bundle.Call.Payload)
	// Attempt to send over network
	err = d.networkSendHandler(jsonMessage)
	if err != nil {
		//TODO: handle retransmission instead of skipping request altogether upon network failure
		d.CompleteRequest(bundle.Call.GetUniqueId())
		if d.onRequestCancel != nil {
			d.onRequestCancel(bundle, ocpp.NewError(GenericError, err.Error(), bundle.Call.UniqueId))
		}
	} else {
		log.Infof("dispatched request %s to server", bundle.Call.UniqueId)
		log.Debugf("sent JSON message to server: %s", string(jsonMessage))
	}
	return err
}

func (d *DefaultClientDispatcher) startRequestTimeout() {
	// Get current queue element
	el := d.requestQueue.Peek()
	bundle, _ := el.(RequestBundle)
	if d.timeout == 0 {
		// Don't start a timer, but set base context and dummy cancel function.
		d.timeoutContext = bundle.Context
		d.timeoutCancelFn = func() {}
	} else {
		// Start timeout
		d.timeoutContext, d.timeoutCancelFn = context.WithTimeout(bundle.Context, d.timeout)
	}
}

func (d *DefaultClientDispatcher) Pause() {
	d.controlChannel <- signalPause
}

func (d *DefaultClientDispatcher) Resume() {
	d.controlChannel <- signalResume
}

func (d *DefaultClientDispatcher) CompleteRequest(requestId string) (bundle RequestBundle) {
	// Critical section, lock prevents concurrent invocations of CompleteRequest,
	// which may lead to race conditions in the request queue.
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	el := d.requestQueue.Peek()
	if el == nil {
		log.Errorf("attempting to pop front of queue, but queue is empty")
		return
	}
	bundle, _ = el.(RequestBundle)
	if bundle.Call.UniqueId != requestId {
		log.Errorf("internal state mismatch: received response for %v but expected response for %v", requestId, bundle.Call.UniqueId)
		return
	}
	d.requestQueue.Pop()
	d.pendingRequestState.DeletePendingRequest(requestId)
	log.Debugf("removed request %v from front of queue", bundle.Call.UniqueId)
	// Signal that next message in queue may be sent
	d.requestChannel <- struct{}{}
	return
}
