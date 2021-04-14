package callbackqueue

import (
	"sync"

	"github.com/lorenzodonini/ocpp-go/ocpp"
)

type CallbackQueue struct {
	callbacksMutex sync.RWMutex
	callbacks      map[string][]func(confirmation ocpp.Response, err error)
}

func New() CallbackQueue {
	return CallbackQueue{
		callbacks: make(map[string][]func(confirmation ocpp.Response, err error)),
	}
}

func (cq *CallbackQueue) Queue(id string, callback func(confirmation ocpp.Response, err error)) {
	cq.callbacksMutex.Lock()
	defer cq.callbacksMutex.Unlock()

	cq.callbacks[id] = append(cq.callbacks[id], callback)
}

func (cq *CallbackQueue) Dequeue(id string) (func(confirmation ocpp.Response, err error), bool) {
	cq.callbacksMutex.Lock()
	defer cq.callbacksMutex.Unlock()

	callbacks, ok := cq.callbacks[id]
	if !ok {
		return nil, false
	}

	if len(callbacks) == 0 {
		panic("Internal CallbackQueue inconsistency")
	}

	callback := callbacks[0]

	if len(callbacks) == 1 {
		delete(cq.callbacks, id)
	} else {
		cq.callbacks[id] = callbacks[1:]
	}

	return callback, ok
}
