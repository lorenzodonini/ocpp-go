package ocppj_test

import (
	"github.com/lorenzodonini/ocpp-go/ocppj"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const clientQueueCapacity = 10

type ClientQueueTestSuite struct {
	suite.Suite
	queue ocppj.RequestQueue
}

func (suite *ClientQueueTestSuite) SetupTest() {
	suite.queue = ocppj.NewFIFOClientQueue(clientQueueCapacity)
}

func (suite *ClientQueueTestSuite) TestQueueEmpty() {
	t := suite.T()
	empty := suite.queue.IsEmpty()
	assert.True(t, empty)
}

func (suite *ClientQueueTestSuite) TestPushElement() {
	t := suite.T()
	req := newMockRequest("somevalue")
	err := suite.queue.Push(req)
	require.Nil(t, err)
	assert.False(t, suite.queue.IsEmpty())
	assert.False(t, suite.queue.IsFull())
	assert.Equal(t, 1, suite.queue.Size())
}

func (suite *ClientQueueTestSuite) TestQueueSize() {
	t := suite.T()
	for i := 0; i < clientQueueCapacity; i++ {
		req := newMockRequest("somevalue")
		err := suite.queue.Push(req)
		require.Nil(t, err)
		assert.False(t, suite.queue.IsEmpty())
		assert.Equal(t, i+1, suite.queue.Size())
	}
}

func (suite *ClientQueueTestSuite) TestQueueFull() {
	t := suite.T()
	for i := 0; i < clientQueueCapacity+2; i++ {
		req := newMockRequest("somevalue")
		err := suite.queue.Push(req)
		if i < clientQueueCapacity {
			require.Nil(t, err)
			if i < clientQueueCapacity-1 {
				assert.False(t, suite.queue.IsFull())
			} else {
				assert.True(t, suite.queue.IsFull())
			}
		} else {
			require.NotNil(t, err)
			assert.True(t, suite.queue.IsFull())
		}
	}
}

func (suite *ClientQueueTestSuite) TestPeekElement() {
	t := suite.T()
	req := newMockRequest("somevalue")
	err := suite.queue.Push(req)
	require.Nil(t, err)
	el := suite.queue.Peek()
	require.NotNil(t, el)
	peeked, ok := el.(*MockRequest)
	require.True(t, ok)
	require.NotNil(t, peeked)
	assert.Equal(t, req.MockValue, peeked.MockValue)
	assert.False(t, suite.queue.IsEmpty())
	assert.False(t, suite.queue.IsFull())
	assert.Equal(t, 1, suite.queue.Size())
}

func (suite *ClientQueueTestSuite) TestPopElement() {
	t := suite.T()
	req := newMockRequest("somevalue")
	err := suite.queue.Push(req)
	require.Nil(t, err)
	el := suite.queue.Pop()
	require.NotNil(t, el)
	popped, ok := el.(*MockRequest)
	require.True(t, ok)
	require.NotNil(t, popped)
	assert.Equal(t, req.MockValue, popped.MockValue)
	assert.True(t, suite.queue.IsEmpty())
	assert.False(t, suite.queue.IsFull())
}
