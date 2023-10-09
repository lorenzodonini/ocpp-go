package ocppj

import (
	"fmt"
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

// FIFOClientQueue is a default queue implementation. The queue is thread-safe.
type FIFOClientQueue struct {
	elements []interface{}
	capacity int
	mutex    sync.RWMutex
}

func (q *FIFOClientQueue) Init() {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.elements = make([]interface{}, 0, q.capacity)
}

func (q *FIFOClientQueue) Push(element interface{}) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if len(q.elements) >= q.capacity && q.capacity > 0 {
		return fmt.Errorf("request queue is full, cannot push new element")
	}
	q.elements = append(q.elements, element)
	return nil
}

func (q *FIFOClientQueue) Peek() interface{} {
	q.mutex.RLock()
	defer q.mutex.RUnlock()
	if len(q.elements) == 0 {
		return nil
	}
	return q.elements[0]
}

func (q *FIFOClientQueue) Pop() interface{} {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if len(q.elements) == 0 {
		return nil
	}
	result := q.elements[0]
	q.elements = q.elements[1:]
	return result
}

func (q *FIFOClientQueue) Size() int {
	q.mutex.RLock()
	defer q.mutex.RUnlock()
	return len(q.elements)
}

func (q *FIFOClientQueue) IsFull() bool {
	q.mutex.RLock()
	defer q.mutex.RUnlock()
	return len(q.elements) >= q.capacity && q.capacity > 0
}

func (q *FIFOClientQueue) IsEmpty() bool {
	q.mutex.RLock()
	defer q.mutex.RUnlock()
	return len(q.elements) == 0
}

// NewFIFOClientQueue creates a new FIFOClientQueue with the given capacity.
//
// A FIFOQueue is backed by a slice, and the capacity represents the maximum capacity of the queue.
// Passing capacity = 0 will create a queue without a maximum capacity.
// The capacity cannot change after creation.
func NewFIFOClientQueue(capacity int) *FIFOClientQueue {
	return &FIFOClientQueue{
		elements: make([]interface{}, 0, capacity),
		capacity: capacity,
	}
}

// ServerQueueMap defines the interface for managing client request queues.
//
// An OCPP-J server may serve multiple clients at the same time, so it will need to provide a queue for each client.
type ServerQueueMap interface {
	// Init puts the queue map in its initial state. May be used for initial setup or clearing.
	Init()
	// Get retrieves the queue associated to a specific clientID.
	// If no such element exists, the returned flag will be false.
	Get(clientID string) (RequestQueue, bool)
	// GetOrCreate retrieves the queue associated to a specific clientID.
	// If no such element exists, it is created, added to the map and returned.
	GetOrCreate(clientID string) RequestQueue
	// Remove deletes the queue associated to a specific clientID.
	// If no such element exists, nothing happens.
	Remove(clientID string)
	// Add inserts a new RequestQueue into the map structure.
	// If such element already exists, it will be replaced with the new queue.
	Add(clientID string, queue RequestQueue)
}

// FIFOQueueMap is a default implementation of ServerQueueMap.
// A FIFOQueueMap is backed by a map[string]RequestQueue. The data structure is thread-safe.
//
// When calling the GetOrCreate function, if no entry for a key was found in the map,
// a new RequestQueue with the given capacity will be created.
type FIFOQueueMap struct {
	data          map[string]RequestQueue
	queueCapacity int
	mutex         sync.RWMutex
}

func (f *FIFOQueueMap) Init() {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.data = map[string]RequestQueue{}
}

func (f *FIFOQueueMap) Get(clientID string) (RequestQueue, bool) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	q, ok := f.data[clientID]
	return q, ok
}

func (f *FIFOQueueMap) GetOrCreate(clientID string) RequestQueue {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	q, ok := f.data[clientID]
	if !ok {
		q = NewFIFOClientQueue(f.queueCapacity)
		f.data[clientID] = q
	}
	return q
}

func (f *FIFOQueueMap) Remove(clientID string) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	delete(f.data, clientID)
}

func (f *FIFOQueueMap) Add(clientID string, queue RequestQueue) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.data[clientID] = queue
}

// NewFIFOQueueMap creates a new FIFOQueueMap, which will automatically create queues with the specified capacity.
//
// Passing capacity = 0 will generate queues without a maximum capacity.
// The capacity cannot change after creation.
func NewFIFOQueueMap(clientQueueCapacity int) *FIFOQueueMap {
	return &FIFOQueueMap{data: map[string]RequestQueue{}, queueCapacity: clientQueueCapacity}
}
