// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	net "net"

	mock "github.com/stretchr/testify/mock"

	tls "crypto/tls"
)

// MockChannel is an autogenerated mock type for the Channel type
type MockChannel struct {
	mock.Mock
}

type MockChannel_Expecter struct {
	mock *mock.Mock
}

func (_m *MockChannel) EXPECT() *MockChannel_Expecter {
	return &MockChannel_Expecter{mock: &_m.Mock}
}

// ID provides a mock function with given fields:
func (_m *MockChannel) ID() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ID")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockChannel_ID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ID'
type MockChannel_ID_Call struct {
	*mock.Call
}

// ID is a helper method to define mock.On call
func (_e *MockChannel_Expecter) ID() *MockChannel_ID_Call {
	return &MockChannel_ID_Call{Call: _e.mock.On("ID")}
}

func (_c *MockChannel_ID_Call) Run(run func()) *MockChannel_ID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockChannel_ID_Call) Return(_a0 string) *MockChannel_ID_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockChannel_ID_Call) RunAndReturn(run func() string) *MockChannel_ID_Call {
	_c.Call.Return(run)
	return _c
}

// RemoteAddr provides a mock function with given fields:
func (_m *MockChannel) RemoteAddr() net.Addr {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for RemoteAddr")
	}

	var r0 net.Addr
	if rf, ok := ret.Get(0).(func() net.Addr); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(net.Addr)
		}
	}

	return r0
}

// MockChannel_RemoteAddr_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RemoteAddr'
type MockChannel_RemoteAddr_Call struct {
	*mock.Call
}

// RemoteAddr is a helper method to define mock.On call
func (_e *MockChannel_Expecter) RemoteAddr() *MockChannel_RemoteAddr_Call {
	return &MockChannel_RemoteAddr_Call{Call: _e.mock.On("RemoteAddr")}
}

func (_c *MockChannel_RemoteAddr_Call) Run(run func()) *MockChannel_RemoteAddr_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockChannel_RemoteAddr_Call) Return(_a0 net.Addr) *MockChannel_RemoteAddr_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockChannel_RemoteAddr_Call) RunAndReturn(run func() net.Addr) *MockChannel_RemoteAddr_Call {
	_c.Call.Return(run)
	return _c
}

// TLSConnectionState provides a mock function with given fields:
func (_m *MockChannel) TLSConnectionState() *tls.ConnectionState {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for TLSConnectionState")
	}

	var r0 *tls.ConnectionState
	if rf, ok := ret.Get(0).(func() *tls.ConnectionState); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*tls.ConnectionState)
		}
	}

	return r0
}

// MockChannel_TLSConnectionState_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'TLSConnectionState'
type MockChannel_TLSConnectionState_Call struct {
	*mock.Call
}

// TLSConnectionState is a helper method to define mock.On call
func (_e *MockChannel_Expecter) TLSConnectionState() *MockChannel_TLSConnectionState_Call {
	return &MockChannel_TLSConnectionState_Call{Call: _e.mock.On("TLSConnectionState")}
}

func (_c *MockChannel_TLSConnectionState_Call) Run(run func()) *MockChannel_TLSConnectionState_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockChannel_TLSConnectionState_Call) Return(_a0 *tls.ConnectionState) *MockChannel_TLSConnectionState_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockChannel_TLSConnectionState_Call) RunAndReturn(run func() *tls.ConnectionState) *MockChannel_TLSConnectionState_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockChannel creates a new instance of MockChannel. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockChannel(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockChannel {
	mock := &MockChannel{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
