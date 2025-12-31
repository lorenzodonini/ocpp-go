package ocppj_test

import (
	"github.com/lorenzodonini/ocpp-go/ocppj"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const queueCapacity = 10

type ClientQueueTestSuite struct {
	suite.Suite
	queue ocppj.RequestQueue
}

func (suite *ClientQueueTestSuite) SetupTest() {
	suite.queue = ocppj.NewFIFOClientQueue(queueCapacity)
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
	for i := 0; i < queueCapacity; i++ {
		req := newMockRequest("somevalue")
		err := suite.queue.Push(req)
		require.Nil(t, err)
		assert.False(t, suite.queue.IsEmpty())
		assert.Equal(t, i+1, suite.queue.Size())
	}
}

func (suite *ClientQueueTestSuite) TestQueueFull() {
	t := suite.T()
	for i := 0; i < queueCapacity+2; i++ {
		req := newMockRequest("somevalue")
		err := suite.queue.Push(req)
		if i < queueCapacity {
			require.Nil(t, err)
			if i < queueCapacity-1 {
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

func (suite *ClientQueueTestSuite) TestQueueNoCapacity() {
	t := suite.T()
	suite.queue = ocppj.NewFIFOClientQueue(0)
	for i := 0; i < 50; i++ {
		req := newMockRequest("somevalue")
		err := suite.queue.Push(req)
		require.Nil(t, err)
	}
	assert.False(t, suite.queue.IsFull())
}

func (suite *ClientQueueTestSuite) TestQueueClear() {
	t := suite.T()
	for i := 0; i < queueCapacity; i++ {
		req := newMockRequest("somevalue")
		err := suite.queue.Push(req)
		require.Nil(t, err)
	}
	assert.True(t, suite.queue.IsFull())
	suite.queue.Init()
	assert.True(t, suite.queue.IsEmpty())
	assert.Equal(t, 0, suite.queue.Size())
}

type ServerQueueMapTestSuite struct {
	suite.Suite
	queueMap ocppj.ServerQueueMap
}

func (suite *ServerQueueMapTestSuite) SetupTest() {
	suite.queueMap = ocppj.NewFIFOQueueMap(queueCapacity)
}

func (suite *ServerQueueMapTestSuite) TestAddElement() {
	t := suite.T()
	q := ocppj.NewFIFOClientQueue(0)
	el := "element1"
	_ = q.Push(el)
	id := "test"
	suite.queueMap.Add(id, q)

	retrieved, ok := suite.queueMap.Get(id)
	require.True(t, ok)
	require.NotNil(t, retrieved)
	assert.False(t, retrieved.IsEmpty())
	assert.Equal(t, 1, retrieved.Size())
	assert.Equal(t, el, retrieved.Peek())
}

func (suite *ServerQueueMapTestSuite) TestGetOrCreate() {
	t := suite.T()
	el := "element1"
	id := "test"
	q, ok := suite.queueMap.Get(id)
	require.False(t, ok)
	require.Nil(t, q)
	q = suite.queueMap.GetOrCreate(id)
	require.NotNil(t, q)
	_ = q.Push(el)
	// Verify consistency
	q, ok = suite.queueMap.Get(id)
	require.True(t, ok)
	assert.Equal(t, 1, q.Size())
	assert.Equal(t, el, q.Peek())
}

func (suite *ServerQueueMapTestSuite) TestRemove() {
	t := suite.T()
	id := "test"
	q := suite.queueMap.GetOrCreate(id)
	require.NotNil(t, q)
	suite.queueMap.Remove(id)
	q, ok := suite.queueMap.Get(id)
	assert.False(t, ok)
	assert.Nil(t, q)
}

func (suite *ServerQueueMapTestSuite) TestSize() {
	suite.Equal(0, suite.queueMap.Size())
	id := "test"
	q := suite.queueMap.GetOrCreate(id)
	suite.NotNil(q)
	err := q.Push("something")
	suite.NoError(err)
	suite.Equal(1, suite.queueMap.Size())
	suite.queueMap.Remove(id)
	suite.Equal(0, suite.queueMap.Size())
}

func (suite *ServerQueueMapTestSuite) TestSizePerClient() {
	for i := 0; i < queueCapacity; i++ {
		id := "client" + string(rune(i%3+'A')) // Three different clients: A, B, C
		q := suite.queueMap.GetOrCreate(id)
		suite.NotNil(q)
		req := newMockRequest("somevalue")
		err := q.Push(req)
		suite.Nil(err)
	}

	expectedMap := map[string]int{
		"clientA": 4,
		"clientB": 3,
		"clientC": 3,
	}

	suite.Equal(expectedMap, suite.queueMap.SizePerClient())
}
