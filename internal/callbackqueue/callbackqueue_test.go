package callbackqueue

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/stretchr/testify/suite"
)

// MockResponse is a simple mock implementation of ocpp.Response for testing
type MockResponse struct {
	Value string
}

func (m *MockResponse) GetFeatureName() string {
	return "Mock"
}

// CallbackQueueTestSuite contains tests for the CallbackQueue
type CallbackQueueTestSuite struct {
	suite.Suite
	queue CallbackQueue
}

func (suite *CallbackQueueTestSuite) SetupTest() {
	suite.queue = New()
}

// TestTryQueueSuccess verifies that TryQueue succeeds when try() returns nil
func (suite *CallbackQueueTestSuite) TestTryQueueSuccess() {
	id := "test-id"
	callbackCalled := atomic.Bool{}
	var receivedResponse ocpp.Response
	var receivedError error

	callback := func(confirmation ocpp.Response, err error) {
		callbackCalled.Store(true)
		receivedResponse = confirmation
		receivedError = err
	}

	try := func() error {
		return nil // Success
	}

	err := suite.queue.TryQueue(id, try, callback)
	suite.Require().NoError(err)
	suite.Assert().False(callbackCalled.Load(), "Callback should not be called yet")

	// Dequeue and invoke callback
	cb, ok := suite.queue.Dequeue(id)
	suite.Require().True(ok)
	suite.Require().NotNil(cb)

	mockResp := &MockResponse{Value: "test"}
	cb(mockResp, nil)

	suite.Assert().True(callbackCalled.Load())
	suite.Assert().Equal(mockResp, receivedResponse)
	suite.Assert().Nil(receivedError)
}

// TestTryQueueFailure verifies that TryQueue removes callback when try() returns error
func (suite *CallbackQueueTestSuite) TestTryQueueFailure() {
	id := "test-id"
	callbackCalled := atomic.Bool{}
	callback := func(confirmation ocpp.Response, err error) {
		callbackCalled.Store(true)
	}

	try := func() error {
		return errors.New("try failed")
	}

	err := suite.queue.TryQueue(id, try, callback)
	suite.Require().Error(err)
	suite.Assert().Equal("try failed", err.Error())
	suite.Assert().False(callbackCalled.Load(), "Callback should not be called")

	// Verify callback was removed from queue
	_, ok := suite.queue.Dequeue(id)
	suite.Assert().False(ok, "Callback should have been removed")
}

// TestDequeueEmpty verifies that Dequeue returns false for empty queue
func (suite *CallbackQueueTestSuite) TestDequeueEmpty() {
	id := "non-existent"

	cb, ok := suite.queue.Dequeue(id)
	suite.Assert().False(ok)
	suite.Assert().Nil(cb)
}

// TestDequeueFIFO verifies that callbacks are dequeued in FIFO order
func (suite *CallbackQueueTestSuite) TestDequeueFIFO() {
	id := "test-id"
	callbacksCalled := make([]int, 0)
	var mu sync.Mutex

	// Queue multiple callbacks
	for i := 0; i < 5; i++ {
		index := i
		callback := func(confirmation ocpp.Response, err error) {
			mu.Lock()
			callbacksCalled = append(callbacksCalled, index)
			mu.Unlock()
		}

		try := func() error {
			return nil
		}

		err := suite.queue.TryQueue(id, try, callback)
		suite.Require().NoError(err)
	}

	// Dequeue and invoke callbacks in order
	for i := 0; i < 5; i++ {
		cb, ok := suite.queue.Dequeue(id)
		suite.Require().True(ok)
		suite.Require().NotNil(cb)
		cb(nil, nil)
	}

	// Verify FIFO order
	mu.Lock()
	suite.Require().Len(callbacksCalled, 5)
	for i := 0; i < 5; i++ {
		suite.Assert().Equal(i, callbacksCalled[i], "Callbacks should be called in FIFO order")
	}
	mu.Unlock()

	// Queue should be empty now
	_, ok := suite.queue.Dequeue(id)
	suite.Assert().False(ok)
}

// TestMultipleIDs verifies that different IDs maintain separate queues
func (suite *CallbackQueueTestSuite) TestMultipleIDs() {
	id1 := "id1"
	id2 := "id2"

	callback1Called := atomic.Bool{}
	callback2Called := atomic.Bool{}

	callback1 := func(confirmation ocpp.Response, err error) {
		callback1Called.Store(true)
	}

	callback2 := func(confirmation ocpp.Response, err error) {
		callback2Called.Store(true)
	}

	try := func() error {
		return nil
	}

	// Queue callbacks for different IDs
	err := suite.queue.TryQueue(id1, try, callback1)
	suite.Require().NoError(err)

	err = suite.queue.TryQueue(id2, try, callback2)
	suite.Require().NoError(err)

	// Dequeue from id1
	cb1, ok := suite.queue.Dequeue(id1)
	suite.Require().True(ok)
	cb1(nil, nil)
	suite.Assert().True(callback1Called.Load())
	suite.Assert().False(callback2Called.Load())

	// Dequeue from id2
	cb2, ok := suite.queue.Dequeue(id2)
	suite.Require().True(ok)
	cb2(nil, nil)
	suite.Assert().True(callback2Called.Load())

	// Both queues should be empty
	_, ok = suite.queue.Dequeue(id1)
	suite.Assert().False(ok)
	_, ok = suite.queue.Dequeue(id2)
	suite.Assert().False(ok)
}

// TestDequeueRemovesFromQueue verifies that Dequeue removes the callback from the queue
func (suite *CallbackQueueTestSuite) TestDequeueRemovesFromQueue() {
	id := "test-id"
	callbackCalled := atomic.Int64{}

	callback := func(confirmation ocpp.Response, err error) {
		callbackCalled.Add(1)
	}

	try := func() error {
		return nil
	}

	// Queue callback
	err := suite.queue.TryQueue(id, try, callback)
	suite.Require().NoError(err)

	// Dequeue first time
	cb1, ok := suite.queue.Dequeue(id)
	suite.Require().True(ok)
	cb1(nil, nil)
	suite.Assert().Equal(1, int(callbackCalled.Load()))

	// Try to dequeue again - should fail
	_, ok = suite.queue.Dequeue(id)
	suite.Assert().False(ok)
	suite.Assert().Equal(1, int(callbackCalled.Load()), "Callback should only be called once")
}

// TestTryQueueFailureRemovesLastCallback verifies that on failure, only the last queued callback is removed
func (suite *CallbackQueueTestSuite) TestTryQueueFailureRemovesLastCallback() {
	id := "test-id"
	callbacksCalled := make([]int, 0)
	var mu sync.Mutex

	// Queue two callbacks successfully
	for i := 0; i < 2; i++ {
		index := i
		callback := func(confirmation ocpp.Response, err error) {
			mu.Lock()
			callbacksCalled = append(callbacksCalled, index)
			mu.Unlock()
		}

		try := func() error {
			return nil
		}

		err := suite.queue.TryQueue(id, try, callback)
		suite.Require().NoError(err)
	}

	// Queue a third callback that fails
	callback3 := func(confirmation ocpp.Response, err error) {
		mu.Lock()
		callbacksCalled = append(callbacksCalled, 3)
		mu.Unlock()
	}

	tryFail := func() error {
		return errors.New("failed")
	}

	err := suite.queue.TryQueue(id, tryFail, callback3)
	suite.Require().Error(err)

	// First two callbacks should still be in queue
	cb1, ok := suite.queue.Dequeue(id)
	suite.Require().True(ok)
	cb1(nil, nil)

	cb2, ok := suite.queue.Dequeue(id)
	suite.Require().True(ok)
	cb2(nil, nil)

	// Queue should be empty
	_, ok = suite.queue.Dequeue(id)
	suite.Assert().False(ok)

	// Verify only first two callbacks were called
	mu.Lock()
	suite.Require().Len(callbacksCalled, 2)
	suite.Assert().Equal(0, callbacksCalled[0])
	suite.Assert().Equal(1, callbacksCalled[1])
	mu.Unlock()
}

// TestConcurrentTryQueue verifies thread safety of TryQueue
func (suite *CallbackQueueTestSuite) TestConcurrentTryQueue() {
	id := "test-id"
	numGoroutines := 100
	var wg sync.WaitGroup
	callbacksQueued := make(chan int, numGoroutines)

	// Concurrently queue callbacks
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)

		go func(index int) {
			defer wg.Done()
			callback := func(confirmation ocpp.Response, err error) {
				callbacksQueued <- index
			}

			try := func() error {
				return nil
			}

			err := suite.queue.TryQueue(id, try, callback)
			suite.Require().NoError(err)
		}(i)
	}

	wg.Wait()

	// Dequeue all callbacks
	received := make(map[int]bool)
	for i := 0; i < numGoroutines; i++ {
		cb, ok := suite.queue.Dequeue(id)
		suite.Require().True(ok)
		cb(nil, nil)
	}

	// Collect all received indices
	for i := 0; i < numGoroutines; i++ {
		index := <-callbacksQueued
		received[index] = true
	}

	// Verify all callbacks were queued and dequeued
	suite.Assert().Len(received, numGoroutines)

	// Queue should be empty
	_, ok := suite.queue.Dequeue(id)
	suite.Assert().False(ok)
}

// TestConcurrentDequeue verifies thread safety of Dequeue
func (suite *CallbackQueueTestSuite) TestConcurrentDequeue() {
	id := "test-id"
	numCallbacks := 50
	var wg sync.WaitGroup
	dequeued := make(chan bool, numCallbacks)

	// Queue callbacks
	for i := 0; i < numCallbacks; i++ {
		callback := func(confirmation ocpp.Response, err error) {
			dequeued <- true
		}

		try := func() error {
			return nil
		}

		err := suite.queue.TryQueue(id, try, callback)
		suite.Require().NoError(err)
	}

	// Concurrently dequeue callbacks
	for i := 0; i < numCallbacks; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cb, ok := suite.queue.Dequeue(id)
			if ok && cb != nil {
				cb(nil, nil)
			}
		}()
	}

	wg.Wait()

	// Verify all callbacks were dequeued
	suite.Assert().Len(dequeued, numCallbacks)

	// Queue should be empty
	_, ok := suite.queue.Dequeue(id)
	suite.Assert().False(ok)
}

// TestDequeueWithError verifies that callbacks can be invoked with errors
func (suite *CallbackQueueTestSuite) TestDequeueWithError() {
	id := "test-id"
	var receivedError error
	var receivedResponse ocpp.Response

	callback := func(confirmation ocpp.Response, err error) {
		receivedResponse = confirmation
		receivedError = err
	}

	try := func() error {
		return nil
	}

	err := suite.queue.TryQueue(id, try, callback)
	suite.Require().NoError(err)

	// Dequeue and invoke with error
	ocppErr := ocpp.NewError(ocpp.ErrorCode("GenericError"), "test error", "123")
	cb, ok := suite.queue.Dequeue(id)
	suite.Require().True(ok)
	cb(nil, ocppErr)

	suite.Assert().Nil(receivedResponse)
	suite.Assert().NotNil(receivedError)
	suite.Assert().Equal(ocppErr, receivedError)
}

// TestDequeueWithResponse verifies that callbacks can be invoked with responses
func (suite *CallbackQueueTestSuite) TestDequeueWithResponse() {
	id := "test-id"
	var receivedError error
	var receivedResponse ocpp.Response

	callback := func(confirmation ocpp.Response, err error) {
		receivedResponse = confirmation
		receivedError = err
	}

	try := func() error {
		return nil
	}

	err := suite.queue.TryQueue(id, try, callback)
	suite.Require().NoError(err)

	// Dequeue and invoke with response
	mockResp := &MockResponse{Value: "success"}
	cb, ok := suite.queue.Dequeue(id)
	suite.Require().True(ok)
	cb(mockResp, nil)

	suite.Assert().Equal(mockResp, receivedResponse)
	suite.Assert().Nil(receivedError)
}

// TestMultipleIDsConcurrent verifies concurrent operations on different IDs
func (suite *CallbackQueueTestSuite) TestMultipleIDsConcurrent() {
	numIDs := 10
	callbacksPerID := 5
	var wg sync.WaitGroup

	// Queue callbacks for multiple IDs concurrently
	for idIndex := 0; idIndex < numIDs; idIndex++ {
		id := string(rune('a' + idIndex))
		for cbIndex := 0; cbIndex < callbacksPerID; cbIndex++ {
			wg.Add(1)
			go func(clientID string, callbackIndex int) {
				defer wg.Done()
				callback := func(confirmation ocpp.Response, err error) {
					// Just verify it's called
				}

				try := func() error {
					return nil
				}

				err := suite.queue.TryQueue(clientID, try, callback)
				suite.Require().NoError(err)
			}(id, cbIndex)
		}
	}

	wg.Wait()

	// Verify all callbacks were queued
	for idIndex := 0; idIndex < numIDs; idIndex++ {
		id := string(rune('a' + idIndex))
		dequeuedCount := 0
		for {
			cb, ok := suite.queue.Dequeue(id)
			if !ok {
				break
			}
			cb(nil, nil)
			dequeuedCount++
		}
		suite.Assert().Equal(callbacksPerID, dequeuedCount, "All callbacks should be dequeued for ID %s", id)
	}
}

// TestDequeueEmptyIDAfterFailure verifies that ID is removed when last callback fails
func (suite *CallbackQueueTestSuite) TestDequeueEmptyIDAfterFailure() {
	id := "test-id"

	callback := func(confirmation ocpp.Response, err error) {}

	tryFail := func() error {
		return errors.New("failed")
	}

	// Queue callback that fails
	err := suite.queue.TryQueue(id, tryFail, callback)
	suite.Require().Error(err)

	// ID should be removed from map
	_, ok := suite.queue.Dequeue(id)
	suite.Assert().False(ok)
}

func TestCallbackQueue(t *testing.T) {
	suite.Run(t, new(CallbackQueueTestSuite))
}
