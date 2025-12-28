package ocppj_test

import (
	"github.com/lorenzodonini/ocpp-go/ocppj"
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
	empty := suite.queue.IsEmpty()
	suite.Assert().True(empty)
}

func (suite *ClientQueueTestSuite) TestPushElement() {
	req := newMockRequest("somevalue")
	err := suite.queue.Push(req)
	suite.Require().Nil(err)
	suite.Assert().False(suite.queue.IsEmpty())
	suite.Assert().False(suite.queue.IsFull())
	suite.Assert().Equal(1, suite.queue.Size())
}

func (suite *ClientQueueTestSuite) TestQueueSize() {
	for i := 0; i < queueCapacity; i++ {
		req := newMockRequest("somevalue")
		err := suite.queue.Push(req)
		suite.Require().Nil(err)
		suite.Assert().False(suite.queue.IsEmpty())
		suite.Assert().Equal(i+1, suite.queue.Size())
	}
}

func (suite *ClientQueueTestSuite) TestQueueFull() {
	for i := 0; i < queueCapacity+2; i++ {
		req := newMockRequest("somevalue")
		err := suite.queue.Push(req)
		if i < queueCapacity {
			suite.Require().Nil(err)
			if i < queueCapacity-1 {
				suite.Assert().False(suite.queue.IsFull())
			} else {
				suite.Assert().True(suite.queue.IsFull())
			}
		} else {
			suite.Require().NotNil(err)
			suite.Assert().True(suite.queue.IsFull())
		}
	}
}

func (suite *ClientQueueTestSuite) TestPeekElement() {
	req := newMockRequest("somevalue")
	err := suite.queue.Push(req)
	suite.Require().Nil(err)
	el := suite.queue.Peek()
	suite.Require().NotNil(el)
	peeked, ok := el.(*MockRequest)
	suite.Require().True(ok)
	suite.Require().NotNil(peeked)
	suite.Assert().Equal(req.MockValue, peeked.MockValue)
	suite.Assert().False(suite.queue.IsEmpty())
	suite.Assert().False(suite.queue.IsFull())
	suite.Assert().Equal(1, suite.queue.Size())
}

func (suite *ClientQueueTestSuite) TestPopElement() {
	req := newMockRequest("somevalue")
	err := suite.queue.Push(req)
	suite.Require().Nil(err)
	el := suite.queue.Pop()
	suite.Require().NotNil(el)
	popped, ok := el.(*MockRequest)
	suite.Require().True(ok)
	suite.Require().NotNil(popped)
	suite.Assert().Equal(req.MockValue, popped.MockValue)
	suite.Assert().True(suite.queue.IsEmpty())
	suite.Assert().False(suite.queue.IsFull())
}

func (suite *ClientQueueTestSuite) TestQueueNoCapacity() {
	suite.queue = ocppj.NewFIFOClientQueue(0)
	for i := 0; i < 50; i++ {
		req := newMockRequest("somevalue")
		err := suite.queue.Push(req)
		suite.Require().Nil(err)
	}
	suite.Assert().False(suite.queue.IsFull())
}

func (suite *ClientQueueTestSuite) TestQueueClear() {
	for i := 0; i < queueCapacity; i++ {
		req := newMockRequest("somevalue")
		err := suite.queue.Push(req)
		suite.Require().Nil(err)
	}
	suite.Assert().True(suite.queue.IsFull())
	suite.queue.Init()
	suite.Assert().True(suite.queue.IsEmpty())
	suite.Assert().Equal(0, suite.queue.Size())
}

type ServerQueueMapTestSuite struct {
	suite.Suite
	queueMap ocppj.ServerQueueMap
}

func (suite *ServerQueueMapTestSuite) SetupTest() {
	suite.queueMap = ocppj.NewFIFOQueueMap(queueCapacity)
}

func (suite *ServerQueueMapTestSuite) TestAddElement() {
	q := ocppj.NewFIFOClientQueue(0)
	el := "element1"
	_ = q.Push(el)
	id := "test"
	suite.queueMap.Add(id, q)

	retrieved, ok := suite.queueMap.Get(id)
	suite.Require().True(ok)
	suite.Require().NotNil(retrieved)
	suite.Assert().False(retrieved.IsEmpty())
	suite.Assert().Equal(1, retrieved.Size())
	suite.Assert().Equal(el, retrieved.Peek())
}

func (suite *ServerQueueMapTestSuite) TestGetOrCreate() {
	el := "element1"
	id := "test"
	q, ok := suite.queueMap.Get(id)
	suite.Require().False(ok)
	suite.Require().Nil(q)
	q = suite.queueMap.GetOrCreate(id)
	suite.Require().NotNil(q)
	_ = q.Push(el)
	// Verify consistency
	q, ok = suite.queueMap.Get(id)
	suite.Require().True(ok)
	suite.Assert().Equal(1, q.Size())
	suite.Assert().Equal(el, q.Peek())
}

func (suite *ServerQueueMapTestSuite) TestRemove() {
	id := "test"
	q := suite.queueMap.GetOrCreate(id)
	suite.Require().NotNil(q)
	suite.queueMap.Remove(id)
	q, ok := suite.queueMap.Get(id)
	suite.Assert().False(ok)
	suite.Assert().Nil(q)
}
