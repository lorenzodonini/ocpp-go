package ocppj

import (
	"container/list"
	"errors"
)

// RequestBundle is a convenience struct for passing a call object struct and the
// raw byte data into the queue containing outgoing requests.
type RequestBundle struct {
	Call *Call
	Data []byte
}

type RequestQueue interface {
	Init()
	Push(element interface{}) error
	Peek() interface{}
	Pop() interface{}
	Size() int
	IsFull() bool
	IsEmpty() bool
}

// FIFOClientQueue is a default queue implementation for OCPP-J clients.
type FIFOClientQueue struct {
	requestQueue *list.List
	capacity     int
}

func (q *FIFOClientQueue) Init() {
	q.requestQueue = q.requestQueue.Init()
}

func (q *FIFOClientQueue) Push(element interface{}) error {
	if q.requestQueue.Len() >= q.capacity {
		return errors.New("request queue is full, cannot push new element")
	}
	q.requestQueue.PushBack(element)
	return nil
}

func (q *FIFOClientQueue) Peek() interface{} {
	return q.requestQueue.Front().Value
}

func (q *FIFOClientQueue) Pop() interface{} {
	result := q.requestQueue.Front()
	if result != nil {
		return q.requestQueue.Remove(result)
	}
	return nil
}

func (q *FIFOClientQueue) Size() int {
	return q.requestQueue.Len()
}

func (q *FIFOClientQueue) IsFull() bool {
	return q.requestQueue.Len() >= q.capacity
}

func (q *FIFOClientQueue) IsEmpty() bool {
	return q.requestQueue.Len() == 0
}

func NewFIFOClientQueue(capacity int) *FIFOClientQueue {
	return &FIFOClientQueue{
		requestQueue: list.New(),
		capacity:     capacity,
	}
}
