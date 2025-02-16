// Code generated by mockery v2.51.0. DO NOT EDIT.

package mocks

import (
	ws "github.com/lorenzodonini/ocpp-go/ws"
	mock "github.com/stretchr/testify/mock"
)

// MockWsClient is an autogenerated mock type for the WsClient type
type MockWsClient struct {
	mock.Mock
}

type MockWsClient_Expecter struct {
	mock *mock.Mock
}

func (_m *MockWsClient) EXPECT() *MockWsClient_Expecter {
	return &MockWsClient_Expecter{mock: &_m.Mock}
}

// AddOption provides a mock function with given fields: option
func (_m *MockWsClient) AddOption(option interface{}) {
	_m.Called(option)
}

// MockWsClient_AddOption_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddOption'
type MockWsClient_AddOption_Call struct {
	*mock.Call
}

// AddOption is a helper method to define mock.On call
//   - option interface{}
func (_e *MockWsClient_Expecter) AddOption(option interface{}) *MockWsClient_AddOption_Call {
	return &MockWsClient_AddOption_Call{Call: _e.mock.On("AddOption", option)}
}

func (_c *MockWsClient_AddOption_Call) Run(run func(option interface{})) *MockWsClient_AddOption_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}))
	})
	return _c
}

func (_c *MockWsClient_AddOption_Call) Return() *MockWsClient_AddOption_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockWsClient_AddOption_Call) RunAndReturn(run func(interface{})) *MockWsClient_AddOption_Call {
	_c.Run(run)
	return _c
}

// Errors provides a mock function with no fields
func (_m *MockWsClient) Errors() <-chan error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Errors")
	}

	var r0 <-chan error
	if rf, ok := ret.Get(0).(func() <-chan error); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan error)
		}
	}

	return r0
}

// MockWsClient_Errors_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Errors'
type MockWsClient_Errors_Call struct {
	*mock.Call
}

// Errors is a helper method to define mock.On call
func (_e *MockWsClient_Expecter) Errors() *MockWsClient_Errors_Call {
	return &MockWsClient_Errors_Call{Call: _e.mock.On("Errors")}
}

func (_c *MockWsClient_Errors_Call) Run(run func()) *MockWsClient_Errors_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockWsClient_Errors_Call) Return(_a0 <-chan error) *MockWsClient_Errors_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockWsClient_Errors_Call) RunAndReturn(run func() <-chan error) *MockWsClient_Errors_Call {
	_c.Call.Return(run)
	return _c
}

// IsConnected provides a mock function with no fields
func (_m *MockWsClient) IsConnected() bool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for IsConnected")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockWsClient_IsConnected_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsConnected'
type MockWsClient_IsConnected_Call struct {
	*mock.Call
}

// IsConnected is a helper method to define mock.On call
func (_e *MockWsClient_Expecter) IsConnected() *MockWsClient_IsConnected_Call {
	return &MockWsClient_IsConnected_Call{Call: _e.mock.On("IsConnected")}
}

func (_c *MockWsClient_IsConnected_Call) Run(run func()) *MockWsClient_IsConnected_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockWsClient_IsConnected_Call) Return(_a0 bool) *MockWsClient_IsConnected_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockWsClient_IsConnected_Call) RunAndReturn(run func() bool) *MockWsClient_IsConnected_Call {
	_c.Call.Return(run)
	return _c
}

// SetBasicAuth provides a mock function with given fields: username, password
func (_m *MockWsClient) SetBasicAuth(username string, password string) {
	_m.Called(username, password)
}

// MockWsClient_SetBasicAuth_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetBasicAuth'
type MockWsClient_SetBasicAuth_Call struct {
	*mock.Call
}

// SetBasicAuth is a helper method to define mock.On call
//   - username string
//   - password string
func (_e *MockWsClient_Expecter) SetBasicAuth(username interface{}, password interface{}) *MockWsClient_SetBasicAuth_Call {
	return &MockWsClient_SetBasicAuth_Call{Call: _e.mock.On("SetBasicAuth", username, password)}
}

func (_c *MockWsClient_SetBasicAuth_Call) Run(run func(username string, password string)) *MockWsClient_SetBasicAuth_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *MockWsClient_SetBasicAuth_Call) Return() *MockWsClient_SetBasicAuth_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockWsClient_SetBasicAuth_Call) RunAndReturn(run func(string, string)) *MockWsClient_SetBasicAuth_Call {
	_c.Run(run)
	return _c
}

// SetDisconnectedHandler provides a mock function with given fields: handler
func (_m *MockWsClient) SetDisconnectedHandler(handler func(error)) {
	_m.Called(handler)
}

// MockWsClient_SetDisconnectedHandler_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetDisconnectedHandler'
type MockWsClient_SetDisconnectedHandler_Call struct {
	*mock.Call
}

// SetDisconnectedHandler is a helper method to define mock.On call
//   - handler func(error)
func (_e *MockWsClient_Expecter) SetDisconnectedHandler(handler interface{}) *MockWsClient_SetDisconnectedHandler_Call {
	return &MockWsClient_SetDisconnectedHandler_Call{Call: _e.mock.On("SetDisconnectedHandler", handler)}
}

func (_c *MockWsClient_SetDisconnectedHandler_Call) Run(run func(handler func(error))) *MockWsClient_SetDisconnectedHandler_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(func(error)))
	})
	return _c
}

func (_c *MockWsClient_SetDisconnectedHandler_Call) Return() *MockWsClient_SetDisconnectedHandler_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockWsClient_SetDisconnectedHandler_Call) RunAndReturn(run func(func(error))) *MockWsClient_SetDisconnectedHandler_Call {
	_c.Run(run)
	return _c
}

// SetHeaderValue provides a mock function with given fields: key, value
func (_m *MockWsClient) SetHeaderValue(key string, value string) {
	_m.Called(key, value)
}

// MockWsClient_SetHeaderValue_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetHeaderValue'
type MockWsClient_SetHeaderValue_Call struct {
	*mock.Call
}

// SetHeaderValue is a helper method to define mock.On call
//   - key string
//   - value string
func (_e *MockWsClient_Expecter) SetHeaderValue(key interface{}, value interface{}) *MockWsClient_SetHeaderValue_Call {
	return &MockWsClient_SetHeaderValue_Call{Call: _e.mock.On("SetHeaderValue", key, value)}
}

func (_c *MockWsClient_SetHeaderValue_Call) Run(run func(key string, value string)) *MockWsClient_SetHeaderValue_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *MockWsClient_SetHeaderValue_Call) Return() *MockWsClient_SetHeaderValue_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockWsClient_SetHeaderValue_Call) RunAndReturn(run func(string, string)) *MockWsClient_SetHeaderValue_Call {
	_c.Run(run)
	return _c
}

// SetMessageHandler provides a mock function with given fields: handler
func (_m *MockWsClient) SetMessageHandler(handler func([]byte) error) {
	_m.Called(handler)
}

// MockWsClient_SetMessageHandler_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetMessageHandler'
type MockWsClient_SetMessageHandler_Call struct {
	*mock.Call
}

// SetMessageHandler is a helper method to define mock.On call
//   - handler func([]byte) error
func (_e *MockWsClient_Expecter) SetMessageHandler(handler interface{}) *MockWsClient_SetMessageHandler_Call {
	return &MockWsClient_SetMessageHandler_Call{Call: _e.mock.On("SetMessageHandler", handler)}
}

func (_c *MockWsClient_SetMessageHandler_Call) Run(run func(handler func([]byte) error)) *MockWsClient_SetMessageHandler_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(func([]byte) error))
	})
	return _c
}

func (_c *MockWsClient_SetMessageHandler_Call) Return() *MockWsClient_SetMessageHandler_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockWsClient_SetMessageHandler_Call) RunAndReturn(run func(func([]byte) error)) *MockWsClient_SetMessageHandler_Call {
	_c.Run(run)
	return _c
}

// SetReconnectedHandler provides a mock function with given fields: handler
func (_m *MockWsClient) SetReconnectedHandler(handler func()) {
	_m.Called(handler)
}

// MockWsClient_SetReconnectedHandler_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetReconnectedHandler'
type MockWsClient_SetReconnectedHandler_Call struct {
	*mock.Call
}

// SetReconnectedHandler is a helper method to define mock.On call
//   - handler func()
func (_e *MockWsClient_Expecter) SetReconnectedHandler(handler interface{}) *MockWsClient_SetReconnectedHandler_Call {
	return &MockWsClient_SetReconnectedHandler_Call{Call: _e.mock.On("SetReconnectedHandler", handler)}
}

func (_c *MockWsClient_SetReconnectedHandler_Call) Run(run func(handler func())) *MockWsClient_SetReconnectedHandler_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(func()))
	})
	return _c
}

func (_c *MockWsClient_SetReconnectedHandler_Call) Return() *MockWsClient_SetReconnectedHandler_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockWsClient_SetReconnectedHandler_Call) RunAndReturn(run func(func())) *MockWsClient_SetReconnectedHandler_Call {
	_c.Run(run)
	return _c
}

// SetRequestedSubProtocol provides a mock function with given fields: subProto
func (_m *MockWsClient) SetRequestedSubProtocol(subProto string) {
	_m.Called(subProto)
}

// MockWsClient_SetRequestedSubProtocol_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetRequestedSubProtocol'
type MockWsClient_SetRequestedSubProtocol_Call struct {
	*mock.Call
}

// SetRequestedSubProtocol is a helper method to define mock.On call
//   - subProto string
func (_e *MockWsClient_Expecter) SetRequestedSubProtocol(subProto interface{}) *MockWsClient_SetRequestedSubProtocol_Call {
	return &MockWsClient_SetRequestedSubProtocol_Call{Call: _e.mock.On("SetRequestedSubProtocol", subProto)}
}

func (_c *MockWsClient_SetRequestedSubProtocol_Call) Run(run func(subProto string)) *MockWsClient_SetRequestedSubProtocol_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockWsClient_SetRequestedSubProtocol_Call) Return() *MockWsClient_SetRequestedSubProtocol_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockWsClient_SetRequestedSubProtocol_Call) RunAndReturn(run func(string)) *MockWsClient_SetRequestedSubProtocol_Call {
	_c.Run(run)
	return _c
}

// SetTimeoutConfig provides a mock function with given fields: config
func (_m *MockWsClient) SetTimeoutConfig(config ws.ClientTimeoutConfig) {
	_m.Called(config)
}

// MockWsClient_SetTimeoutConfig_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetTimeoutConfig'
type MockWsClient_SetTimeoutConfig_Call struct {
	*mock.Call
}

// SetTimeoutConfig is a helper method to define mock.On call
//   - config ws.ClientTimeoutConfig
func (_e *MockWsClient_Expecter) SetTimeoutConfig(config interface{}) *MockWsClient_SetTimeoutConfig_Call {
	return &MockWsClient_SetTimeoutConfig_Call{Call: _e.mock.On("SetTimeoutConfig", config)}
}

func (_c *MockWsClient_SetTimeoutConfig_Call) Run(run func(config ws.ClientTimeoutConfig)) *MockWsClient_SetTimeoutConfig_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(ws.ClientTimeoutConfig))
	})
	return _c
}

func (_c *MockWsClient_SetTimeoutConfig_Call) Return() *MockWsClient_SetTimeoutConfig_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockWsClient_SetTimeoutConfig_Call) RunAndReturn(run func(ws.ClientTimeoutConfig)) *MockWsClient_SetTimeoutConfig_Call {
	_c.Run(run)
	return _c
}

// Start provides a mock function with given fields: url
func (_m *MockWsClient) Start(url string) error {
	ret := _m.Called(url)

	if len(ret) == 0 {
		panic("no return value specified for Start")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(url)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockWsClient_Start_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Start'
type MockWsClient_Start_Call struct {
	*mock.Call
}

// Start is a helper method to define mock.On call
//   - url string
func (_e *MockWsClient_Expecter) Start(url interface{}) *MockWsClient_Start_Call {
	return &MockWsClient_Start_Call{Call: _e.mock.On("Start", url)}
}

func (_c *MockWsClient_Start_Call) Run(run func(url string)) *MockWsClient_Start_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockWsClient_Start_Call) Return(_a0 error) *MockWsClient_Start_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockWsClient_Start_Call) RunAndReturn(run func(string) error) *MockWsClient_Start_Call {
	_c.Call.Return(run)
	return _c
}

// StartWithRetries provides a mock function with given fields: url
func (_m *MockWsClient) StartWithRetries(url string) {
	_m.Called(url)
}

// MockWsClient_StartWithRetries_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'StartWithRetries'
type MockWsClient_StartWithRetries_Call struct {
	*mock.Call
}

// StartWithRetries is a helper method to define mock.On call
//   - url string
func (_e *MockWsClient_Expecter) StartWithRetries(url interface{}) *MockWsClient_StartWithRetries_Call {
	return &MockWsClient_StartWithRetries_Call{Call: _e.mock.On("StartWithRetries", url)}
}

func (_c *MockWsClient_StartWithRetries_Call) Run(run func(url string)) *MockWsClient_StartWithRetries_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockWsClient_StartWithRetries_Call) Return() *MockWsClient_StartWithRetries_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockWsClient_StartWithRetries_Call) RunAndReturn(run func(string)) *MockWsClient_StartWithRetries_Call {
	_c.Run(run)
	return _c
}

// Stop provides a mock function with no fields
func (_m *MockWsClient) Stop() {
	_m.Called()
}

// MockWsClient_Stop_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Stop'
type MockWsClient_Stop_Call struct {
	*mock.Call
}

// Stop is a helper method to define mock.On call
func (_e *MockWsClient_Expecter) Stop() *MockWsClient_Stop_Call {
	return &MockWsClient_Stop_Call{Call: _e.mock.On("Stop")}
}

func (_c *MockWsClient_Stop_Call) Run(run func()) *MockWsClient_Stop_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockWsClient_Stop_Call) Return() *MockWsClient_Stop_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockWsClient_Stop_Call) RunAndReturn(run func()) *MockWsClient_Stop_Call {
	_c.Run(run)
	return _c
}

// Write provides a mock function with given fields: data
func (_m *MockWsClient) Write(data []byte) error {
	ret := _m.Called(data)

	if len(ret) == 0 {
		panic("no return value specified for Write")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func([]byte) error); ok {
		r0 = rf(data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockWsClient_Write_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Write'
type MockWsClient_Write_Call struct {
	*mock.Call
}

// Write is a helper method to define mock.On call
//   - data []byte
func (_e *MockWsClient_Expecter) Write(data interface{}) *MockWsClient_Write_Call {
	return &MockWsClient_Write_Call{Call: _e.mock.On("Write", data)}
}

func (_c *MockWsClient_Write_Call) Run(run func(data []byte)) *MockWsClient_Write_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]byte))
	})
	return _c
}

func (_c *MockWsClient_Write_Call) Return(_a0 error) *MockWsClient_Write_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockWsClient_Write_Call) RunAndReturn(run func([]byte) error) *MockWsClient_Write_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockWsClient creates a new instance of MockWsClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockWsClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockWsClient {
	mock := &MockWsClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
