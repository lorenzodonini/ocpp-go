package callbackqueue

import (
	"sync"

	"github.com/lorenzodonini/ocpp-go/ocpp"
)

type RequestType string
type CallbackQueue struct {
	callbacksMutex sync.RWMutex
	callbacks      map[string]map[RequestType][]func(confirmation ocpp.Response, err error)
}

func New() CallbackQueue {
	return CallbackQueue{
		callbacks: make(map[string]map[RequestType][]func(confirmation ocpp.Response, err error)),
	}
}

func (cq *CallbackQueue) TryQueue(id string, requestType RequestType, try func() error, callback func(confirmation ocpp.Response, err error)) error {
	cq.callbacksMutex.Lock()
	defer cq.callbacksMutex.Unlock()

	if _, ok := cq.callbacks[id]; !ok {
		cq.callbacks[id] = make(map[RequestType][]func(confirmation ocpp.Response, err error))
	}
	cq.callbacks[id][requestType] = append(cq.callbacks[id][requestType], callback)

	if err := try(); err != nil {
		// pop off our element
		if callbacks, ok := cq.callbacks[id]; ok {
			delete(callbacks, requestType)
		}
		if len(cq.callbacks[id]) == 0 {
			delete(cq.callbacks, id)
		}
		return err
	}

	return nil
}

func (cq *CallbackQueue) Dequeue(id string, requestType RequestType) (func(confirmation ocpp.Response, err error), bool) {
	cq.callbacksMutex.Lock()
	defer cq.callbacksMutex.Unlock()

	clientCallbacks, ok := cq.callbacks[id]
	if !ok {
		return nil, false
	}

	if len(clientCallbacks) == 0 {
		//panic("Internal CallbackQueue inconsistency")
		return nil, false
	}

	requestTypeCallbacks, ok := clientCallbacks[requestType]
	if !ok {
		if requestType != "" { /* requestType known and not available... */
			return nil, false
		}
		/* requestType any, take first one... */
		for reqType, cb := range clientCallbacks {
			requestType = reqType
			requestTypeCallbacks = append(requestTypeCallbacks, cb...)
			break // only first one
		}
	}

	callback := requestTypeCallbacks[0]

	if len(requestTypeCallbacks) == 1 {
		delete(cq.callbacks[id], requestType)
	} else {
		cq.callbacks[id][requestType] = requestTypeCallbacks[1:]
	}

	return callback, true
}
