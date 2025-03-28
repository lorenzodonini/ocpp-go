// Code generated by mockery v2.51.0. DO NOT EDIT.

package mocks

import (
	net "net"
	http "net/http"

	mock "github.com/stretchr/testify/mock"

	websocket "github.com/gorilla/websocket"

	ws "github.com/lorenzodonini/ocpp-go/ws"
)

// MockServer is an autogenerated mock type for the Server type
type MockServer struct {
	mock.Mock
}

type MockServer_Expecter struct {
	mock *mock.Mock
}

func (_m *MockServer) EXPECT() *MockServer_Expecter {
	return &MockServer_Expecter{mock: &_m.Mock}
}

// AddSupportedSubprotocol provides a mock function with given fields: subProto
func (_m *MockServer) AddSupportedSubprotocol(subProto string) {
	_m.Called(subProto)
}

// MockServer_AddSupportedSubprotocol_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddSupportedSubprotocol'
type MockServer_AddSupportedSubprotocol_Call struct {
	*mock.Call
}

// AddSupportedSubprotocol is a helper method to define mock.On call
//   - subProto string
func (_e *MockServer_Expecter) AddSupportedSubprotocol(subProto interface{}) *MockServer_AddSupportedSubprotocol_Call {
	return &MockServer_AddSupportedSubprotocol_Call{Call: _e.mock.On("AddSupportedSubprotocol", subProto)}
}

func (_c *MockServer_AddSupportedSubprotocol_Call) Run(run func(subProto string)) *MockServer_AddSupportedSubprotocol_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockServer_AddSupportedSubprotocol_Call) Return() *MockServer_AddSupportedSubprotocol_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockServer_AddSupportedSubprotocol_Call) RunAndReturn(run func(string)) *MockServer_AddSupportedSubprotocol_Call {
	_c.Run(run)
	return _c
}

// Addr provides a mock function with no fields
func (_m *MockServer) Addr() *net.TCPAddr {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Addr")
	}

	var r0 *net.TCPAddr
	if rf, ok := ret.Get(0).(func() *net.TCPAddr); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*net.TCPAddr)
		}
	}

	return r0
}

// MockServer_Addr_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Addr'
type MockServer_Addr_Call struct {
	*mock.Call
}

// Addr is a helper method to define mock.On call
func (_e *MockServer_Expecter) Addr() *MockServer_Addr_Call {
	return &MockServer_Addr_Call{Call: _e.mock.On("Addr")}
}

func (_c *MockServer_Addr_Call) Run(run func()) *MockServer_Addr_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockServer_Addr_Call) Return(_a0 *net.TCPAddr) *MockServer_Addr_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockServer_Addr_Call) RunAndReturn(run func() *net.TCPAddr) *MockServer_Addr_Call {
	_c.Call.Return(run)
	return _c
}

// Errors provides a mock function with no fields
func (_m *MockServer) Errors() <-chan error {
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

// MockServer_Errors_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Errors'
type MockServer_Errors_Call struct {
	*mock.Call
}

// Errors is a helper method to define mock.On call
func (_e *MockServer_Expecter) Errors() *MockServer_Errors_Call {
	return &MockServer_Errors_Call{Call: _e.mock.On("Errors")}
}

func (_c *MockServer_Errors_Call) Run(run func()) *MockServer_Errors_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockServer_Errors_Call) Return(_a0 <-chan error) *MockServer_Errors_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockServer_Errors_Call) RunAndReturn(run func() <-chan error) *MockServer_Errors_Call {
	_c.Call.Return(run)
	return _c
}

// GetChannel provides a mock function with given fields: websocketId
func (_m *MockServer) GetChannel(websocketId string) (ws.Channel, bool) {
	ret := _m.Called(websocketId)

	if len(ret) == 0 {
		panic("no return value specified for GetChannel")
	}

	var r0 ws.Channel
	var r1 bool
	if rf, ok := ret.Get(0).(func(string) (ws.Channel, bool)); ok {
		return rf(websocketId)
	}
	if rf, ok := ret.Get(0).(func(string) ws.Channel); ok {
		r0 = rf(websocketId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(ws.Channel)
		}
	}

	if rf, ok := ret.Get(1).(func(string) bool); ok {
		r1 = rf(websocketId)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// MockServer_GetChannel_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetChannel'
type MockServer_GetChannel_Call struct {
	*mock.Call
}

// GetChannel is a helper method to define mock.On call
//   - websocketId string
func (_e *MockServer_Expecter) GetChannel(websocketId interface{}) *MockServer_GetChannel_Call {
	return &MockServer_GetChannel_Call{Call: _e.mock.On("GetChannel", websocketId)}
}

func (_c *MockServer_GetChannel_Call) Run(run func(websocketId string)) *MockServer_GetChannel_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockServer_GetChannel_Call) Return(_a0 ws.Channel, _a1 bool) *MockServer_GetChannel_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockServer_GetChannel_Call) RunAndReturn(run func(string) (ws.Channel, bool)) *MockServer_GetChannel_Call {
	_c.Call.Return(run)
	return _c
}

// SetBasicAuthHandler provides a mock function with given fields: handler
func (_m *MockServer) SetBasicAuthHandler(handler func(string, string) bool) {
	_m.Called(handler)
}

// MockServer_SetBasicAuthHandler_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetBasicAuthHandler'
type MockServer_SetBasicAuthHandler_Call struct {
	*mock.Call
}

// SetBasicAuthHandler is a helper method to define mock.On call
//   - handler func(string , string) bool
func (_e *MockServer_Expecter) SetBasicAuthHandler(handler interface{}) *MockServer_SetBasicAuthHandler_Call {
	return &MockServer_SetBasicAuthHandler_Call{Call: _e.mock.On("SetBasicAuthHandler", handler)}
}

func (_c *MockServer_SetBasicAuthHandler_Call) Run(run func(handler func(string, string) bool)) *MockServer_SetBasicAuthHandler_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(func(string, string) bool))
	})
	return _c
}

func (_c *MockServer_SetBasicAuthHandler_Call) Return() *MockServer_SetBasicAuthHandler_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockServer_SetBasicAuthHandler_Call) RunAndReturn(run func(func(string, string) bool)) *MockServer_SetBasicAuthHandler_Call {
	_c.Run(run)
	return _c
}

// SetChargePointIdResolver provides a mock function with given fields: resolver
func (_m *MockServer) SetChargePointIdResolver(resolver func(*http.Request) (string, error)) {
	_m.Called(resolver)
}

// MockServer_SetChargePointIdResolver_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetChargePointIdResolver'
type MockServer_SetChargePointIdResolver_Call struct {
	*mock.Call
}

// SetChargePointIdResolver is a helper method to define mock.On call
//   - resolver func(*http.Request)(string , error)
func (_e *MockServer_Expecter) SetChargePointIdResolver(resolver interface{}) *MockServer_SetChargePointIdResolver_Call {
	return &MockServer_SetChargePointIdResolver_Call{Call: _e.mock.On("SetChargePointIdResolver", resolver)}
}

func (_c *MockServer_SetChargePointIdResolver_Call) Run(run func(resolver func(*http.Request) (string, error))) *MockServer_SetChargePointIdResolver_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(func(*http.Request) (string, error)))
	})
	return _c
}

func (_c *MockServer_SetChargePointIdResolver_Call) Return() *MockServer_SetChargePointIdResolver_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockServer_SetChargePointIdResolver_Call) RunAndReturn(run func(func(*http.Request) (string, error))) *MockServer_SetChargePointIdResolver_Call {
	_c.Run(run)
	return _c
}

// SetCheckClientHandler provides a mock function with given fields: handler
func (_m *MockServer) SetCheckClientHandler(handler ws.CheckClientHandler) {
	_m.Called(handler)
}

// MockServer_SetCheckClientHandler_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetCheckClientHandler'
type MockServer_SetCheckClientHandler_Call struct {
	*mock.Call
}

// SetCheckClientHandler is a helper method to define mock.On call
//   - handler ws.CheckClientHandler
func (_e *MockServer_Expecter) SetCheckClientHandler(handler interface{}) *MockServer_SetCheckClientHandler_Call {
	return &MockServer_SetCheckClientHandler_Call{Call: _e.mock.On("SetCheckClientHandler", handler)}
}

func (_c *MockServer_SetCheckClientHandler_Call) Run(run func(handler ws.CheckClientHandler)) *MockServer_SetCheckClientHandler_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(ws.CheckClientHandler))
	})
	return _c
}

func (_c *MockServer_SetCheckClientHandler_Call) Return() *MockServer_SetCheckClientHandler_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockServer_SetCheckClientHandler_Call) RunAndReturn(run func(ws.CheckClientHandler)) *MockServer_SetCheckClientHandler_Call {
	_c.Run(run)
	return _c
}

// SetCheckOriginHandler provides a mock function with given fields: handler
func (_m *MockServer) SetCheckOriginHandler(handler func(*http.Request) bool) {
	_m.Called(handler)
}

// MockServer_SetCheckOriginHandler_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetCheckOriginHandler'
type MockServer_SetCheckOriginHandler_Call struct {
	*mock.Call
}

// SetCheckOriginHandler is a helper method to define mock.On call
//   - handler func(*http.Request) bool
func (_e *MockServer_Expecter) SetCheckOriginHandler(handler interface{}) *MockServer_SetCheckOriginHandler_Call {
	return &MockServer_SetCheckOriginHandler_Call{Call: _e.mock.On("SetCheckOriginHandler", handler)}
}

func (_c *MockServer_SetCheckOriginHandler_Call) Run(run func(handler func(*http.Request) bool)) *MockServer_SetCheckOriginHandler_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(func(*http.Request) bool))
	})
	return _c
}

func (_c *MockServer_SetCheckOriginHandler_Call) Return() *MockServer_SetCheckOriginHandler_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockServer_SetCheckOriginHandler_Call) RunAndReturn(run func(func(*http.Request) bool)) *MockServer_SetCheckOriginHandler_Call {
	_c.Run(run)
	return _c
}

// SetDisconnectedClientHandler provides a mock function with given fields: handler
func (_m *MockServer) SetDisconnectedClientHandler(handler func(ws.Channel)) {
	_m.Called(handler)
}

// MockServer_SetDisconnectedClientHandler_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetDisconnectedClientHandler'
type MockServer_SetDisconnectedClientHandler_Call struct {
	*mock.Call
}

// SetDisconnectedClientHandler is a helper method to define mock.On call
//   - handler func(ws.Channel)
func (_e *MockServer_Expecter) SetDisconnectedClientHandler(handler interface{}) *MockServer_SetDisconnectedClientHandler_Call {
	return &MockServer_SetDisconnectedClientHandler_Call{Call: _e.mock.On("SetDisconnectedClientHandler", handler)}
}

func (_c *MockServer_SetDisconnectedClientHandler_Call) Run(run func(handler func(ws.Channel))) *MockServer_SetDisconnectedClientHandler_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(func(ws.Channel)))
	})
	return _c
}

func (_c *MockServer_SetDisconnectedClientHandler_Call) Return() *MockServer_SetDisconnectedClientHandler_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockServer_SetDisconnectedClientHandler_Call) RunAndReturn(run func(func(ws.Channel))) *MockServer_SetDisconnectedClientHandler_Call {
	_c.Run(run)
	return _c
}

// SetMessageHandler provides a mock function with given fields: handler
func (_m *MockServer) SetMessageHandler(handler ws.MessageHandler) {
	_m.Called(handler)
}

// MockServer_SetMessageHandler_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetMessageHandler'
type MockServer_SetMessageHandler_Call struct {
	*mock.Call
}

// SetMessageHandler is a helper method to define mock.On call
//   - handler ws.MessageHandler
func (_e *MockServer_Expecter) SetMessageHandler(handler interface{}) *MockServer_SetMessageHandler_Call {
	return &MockServer_SetMessageHandler_Call{Call: _e.mock.On("SetMessageHandler", handler)}
}

func (_c *MockServer_SetMessageHandler_Call) Run(run func(handler ws.MessageHandler)) *MockServer_SetMessageHandler_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(ws.MessageHandler))
	})
	return _c
}

func (_c *MockServer_SetMessageHandler_Call) Return() *MockServer_SetMessageHandler_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockServer_SetMessageHandler_Call) RunAndReturn(run func(ws.MessageHandler)) *MockServer_SetMessageHandler_Call {
	_c.Run(run)
	return _c
}

// SetNewClientHandler provides a mock function with given fields: handler
func (_m *MockServer) SetNewClientHandler(handler ws.ConnectedHandler) {
	_m.Called(handler)
}

// MockServer_SetNewClientHandler_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetNewClientHandler'
type MockServer_SetNewClientHandler_Call struct {
	*mock.Call
}

// SetNewClientHandler is a helper method to define mock.On call
//   - handler ws.ConnectedHandler
func (_e *MockServer_Expecter) SetNewClientHandler(handler interface{}) *MockServer_SetNewClientHandler_Call {
	return &MockServer_SetNewClientHandler_Call{Call: _e.mock.On("SetNewClientHandler", handler)}
}

func (_c *MockServer_SetNewClientHandler_Call) Run(run func(handler ws.ConnectedHandler)) *MockServer_SetNewClientHandler_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(ws.ConnectedHandler))
	})
	return _c
}

func (_c *MockServer_SetNewClientHandler_Call) Return() *MockServer_SetNewClientHandler_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockServer_SetNewClientHandler_Call) RunAndReturn(run func(ws.ConnectedHandler)) *MockServer_SetNewClientHandler_Call {
	_c.Run(run)
	return _c
}

// SetTimeoutConfig provides a mock function with given fields: config
func (_m *MockServer) SetTimeoutConfig(config ws.ServerTimeoutConfig) {
	_m.Called(config)
}

// MockServer_SetTimeoutConfig_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetTimeoutConfig'
type MockServer_SetTimeoutConfig_Call struct {
	*mock.Call
}

// SetTimeoutConfig is a helper method to define mock.On call
//   - config ws.ServerTimeoutConfig
func (_e *MockServer_Expecter) SetTimeoutConfig(config interface{}) *MockServer_SetTimeoutConfig_Call {
	return &MockServer_SetTimeoutConfig_Call{Call: _e.mock.On("SetTimeoutConfig", config)}
}

func (_c *MockServer_SetTimeoutConfig_Call) Run(run func(config ws.ServerTimeoutConfig)) *MockServer_SetTimeoutConfig_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(ws.ServerTimeoutConfig))
	})
	return _c
}

func (_c *MockServer_SetTimeoutConfig_Call) Return() *MockServer_SetTimeoutConfig_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockServer_SetTimeoutConfig_Call) RunAndReturn(run func(ws.ServerTimeoutConfig)) *MockServer_SetTimeoutConfig_Call {
	_c.Run(run)
	return _c
}

// Start provides a mock function with given fields: port, listenPath
func (_m *MockServer) Start(port int, listenPath string) {
	_m.Called(port, listenPath)
}

// MockServer_Start_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Start'
type MockServer_Start_Call struct {
	*mock.Call
}

// Start is a helper method to define mock.On call
//   - port int
//   - listenPath string
func (_e *MockServer_Expecter) Start(port interface{}, listenPath interface{}) *MockServer_Start_Call {
	return &MockServer_Start_Call{Call: _e.mock.On("Start", port, listenPath)}
}

func (_c *MockServer_Start_Call) Run(run func(port int, listenPath string)) *MockServer_Start_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int), args[1].(string))
	})
	return _c
}

func (_c *MockServer_Start_Call) Return() *MockServer_Start_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockServer_Start_Call) RunAndReturn(run func(int, string)) *MockServer_Start_Call {
	_c.Run(run)
	return _c
}

// Stop provides a mock function with no fields
func (_m *MockServer) Stop() {
	_m.Called()
}

// MockServer_Stop_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Stop'
type MockServer_Stop_Call struct {
	*mock.Call
}

// Stop is a helper method to define mock.On call
func (_e *MockServer_Expecter) Stop() *MockServer_Stop_Call {
	return &MockServer_Stop_Call{Call: _e.mock.On("Stop")}
}

func (_c *MockServer_Stop_Call) Run(run func()) *MockServer_Stop_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockServer_Stop_Call) Return() *MockServer_Stop_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockServer_Stop_Call) RunAndReturn(run func()) *MockServer_Stop_Call {
	_c.Run(run)
	return _c
}

// StopConnection provides a mock function with given fields: id, closeError
func (_m *MockServer) StopConnection(id string, closeError websocket.CloseError) error {
	ret := _m.Called(id, closeError)

	if len(ret) == 0 {
		panic("no return value specified for StopConnection")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, websocket.CloseError) error); ok {
		r0 = rf(id, closeError)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockServer_StopConnection_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'StopConnection'
type MockServer_StopConnection_Call struct {
	*mock.Call
}

// StopConnection is a helper method to define mock.On call
//   - id string
//   - closeError websocket.CloseError
func (_e *MockServer_Expecter) StopConnection(id interface{}, closeError interface{}) *MockServer_StopConnection_Call {
	return &MockServer_StopConnection_Call{Call: _e.mock.On("StopConnection", id, closeError)}
}

func (_c *MockServer_StopConnection_Call) Run(run func(id string, closeError websocket.CloseError)) *MockServer_StopConnection_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(websocket.CloseError))
	})
	return _c
}

func (_c *MockServer_StopConnection_Call) Return(_a0 error) *MockServer_StopConnection_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockServer_StopConnection_Call) RunAndReturn(run func(string, websocket.CloseError) error) *MockServer_StopConnection_Call {
	_c.Call.Return(run)
	return _c
}

// Write provides a mock function with given fields: webSocketId, data
func (_m *MockServer) Write(webSocketId string, data []byte) error {
	ret := _m.Called(webSocketId, data)

	if len(ret) == 0 {
		panic("no return value specified for Write")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, []byte) error); ok {
		r0 = rf(webSocketId, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockServer_Write_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Write'
type MockServer_Write_Call struct {
	*mock.Call
}

// Write is a helper method to define mock.On call
//   - webSocketId string
//   - data []byte
func (_e *MockServer_Expecter) Write(webSocketId interface{}, data interface{}) *MockServer_Write_Call {
	return &MockServer_Write_Call{Call: _e.mock.On("Write", webSocketId, data)}
}

func (_c *MockServer_Write_Call) Run(run func(webSocketId string, data []byte)) *MockServer_Write_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].([]byte))
	})
	return _c
}

func (_c *MockServer_Write_Call) Return(_a0 error) *MockServer_Write_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockServer_Write_Call) RunAndReturn(run func(string, []byte) error) *MockServer_Write_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockServer creates a new instance of MockServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockServer(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockServer {
	mock := &MockServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
