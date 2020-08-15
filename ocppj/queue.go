package ocppj

import (
	"container/list"
	"errors"
	"sync"
)

// RequestBundle is a convenience struct for passing a call object struct and the
// raw byte data into the queue containing outgoing requests.
type RequestBundle struct {
	Call *Call
	Data []byte
}

// RequestQueue can be arbitrarily implemented, as long as it conforms to the Queue interface.
//
// A RequestQueue is used by ocppj client and server to manage outgoing requests.
// The underlying data structure must be thread-safe, since different goroutines may access it at the same time.
type RequestQueue interface {
	// Init puts the queue in its initial state. May be used for initial setup or clearing.
	Init()
	// Push appends the given element at the end of the queue.
	// Returns an error if the operation failed (e.g. the queue is full).
	Push(element interface{}) error
	// Peek returns the first element of the queue, without removing it from the data structure.
	Peek() interface{}
	// Pop returns the first element of the queue, removing it from the queue.
	Pop() interface{}
	// Size returns the current size of the queue.
	Size() int
	// IsFull returns true if the queue is currently full, false otherwise.
	IsFull() bool
	// IsEmpty returns true if the queue is currently empty, false otherwise.
	IsEmpty() bool
}

// FIFOClientQueue is a default queue implementation for OCPP-J clients.
type FIFOClientQueue struct {
	requestQueue *list.List
	capacity     int
	mutex        sync.Mutex
}

func (q *FIFOClientQueue) Init() {
	q.requestQueue = q.requestQueue.Init()
}

func (q *FIFOClientQueue) Push(element interface{}) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if q.requestQueue.Len() >= q.capacity {
		return errors.New("request queue is full, cannot push new element")
	}
	q.requestQueue.PushBack(element)
	return nil
}

func (q *FIFOClientQueue) Peek() interface{} {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	el := q.requestQueue.Front()
	if el == nil {
		return nil
	}
	return el.Value
}

func (q *FIFOClientQueue) Pop() interface{} {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	result := q.requestQueue.Front()
	if result != nil {
		return q.requestQueue.Remove(result)
	}
	return nil
}

func (q *FIFOClientQueue) Size() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.requestQueue.Len()
}

func (q *FIFOClientQueue) IsFull() bool {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.requestQueue.Len() >= q.capacity
}

func (q *FIFOClientQueue) IsEmpty() bool {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.requestQueue.Len() == 0
}

// NewFIFOClientQueue creates a new FIFOClientQueue with the given capacity.
//
// A FIFOQueue is backed by a linked list, and the capacity represents the maximum capacity of the queue.
// The capacity cannot change after creation.
func NewFIFOClientQueue(capacity int) *FIFOClientQueue {
	return &FIFOClientQueue{
		requestQueue: list.New(),
		capacity:     capacity,
	}
}
