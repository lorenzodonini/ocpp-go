// Code generated by mockery v2.51.0. DO NOT EDIT.

package mocks

import (
	ocpp "github.com/lorenzodonini/ocpp-go/ocpp"
	mock "github.com/stretchr/testify/mock"

	ws "github.com/lorenzodonini/ocpp-go/ws"
)

// MockInvalidMessageHook is an autogenerated mock type for the InvalidMessageHook type
type MockInvalidMessageHook struct {
	mock.Mock
}

type MockInvalidMessageHook_Expecter struct {
	mock *mock.Mock
}

func (_m *MockInvalidMessageHook) EXPECT() *MockInvalidMessageHook_Expecter {
	return &MockInvalidMessageHook_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: client, err, rawJson, parsedFields
func (_m *MockInvalidMessageHook) Execute(client ws.Channel, err *ocpp.Error, rawJson string, parsedFields []interface{}) *ocpp.Error {
	ret := _m.Called(client, err, rawJson, parsedFields)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 *ocpp.Error
	if rf, ok := ret.Get(0).(func(ws.Channel, *ocpp.Error, string, []interface{}) *ocpp.Error); ok {
		r0 = rf(client, err, rawJson, parsedFields)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ocpp.Error)
		}
	}

	return r0
}

// MockInvalidMessageHook_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type MockInvalidMessageHook_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - client ws.Channel
//   - err *ocpp.Error
//   - rawJson string
//   - parsedFields []interface{}
func (_e *MockInvalidMessageHook_Expecter) Execute(client interface{}, err interface{}, rawJson interface{}, parsedFields interface{}) *MockInvalidMessageHook_Execute_Call {
	return &MockInvalidMessageHook_Execute_Call{Call: _e.mock.On("Execute", client, err, rawJson, parsedFields)}
}

func (_c *MockInvalidMessageHook_Execute_Call) Run(run func(client ws.Channel, err *ocpp.Error, rawJson string, parsedFields []interface{})) *MockInvalidMessageHook_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(ws.Channel), args[1].(*ocpp.Error), args[2].(string), args[3].([]interface{}))
	})
	return _c
}

func (_c *MockInvalidMessageHook_Execute_Call) Return(_a0 *ocpp.Error) *MockInvalidMessageHook_Execute_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockInvalidMessageHook_Execute_Call) RunAndReturn(run func(ws.Channel, *ocpp.Error, string, []interface{}) *ocpp.Error) *MockInvalidMessageHook_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockInvalidMessageHook creates a new instance of MockInvalidMessageHook. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockInvalidMessageHook(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockInvalidMessageHook {
	mock := &MockInvalidMessageHook{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
