// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	localauth "github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	mock "github.com/stretchr/testify/mock"
)

// MockLocalAuthListChargePointHandler is an autogenerated mock type for the ChargePointHandler type
type MockLocalAuthListChargePointHandler struct {
	mock.Mock
}

type MockLocalAuthListChargePointHandler_Expecter struct {
	mock *mock.Mock
}

func (_m *MockLocalAuthListChargePointHandler) EXPECT() *MockLocalAuthListChargePointHandler_Expecter {
	return &MockLocalAuthListChargePointHandler_Expecter{mock: &_m.Mock}
}

// OnGetLocalListVersion provides a mock function with given fields: request
func (_m *MockLocalAuthListChargePointHandler) OnGetLocalListVersion(request *localauth.GetLocalListVersionRequest) (*localauth.GetLocalListVersionConfirmation, error) {
	ret := _m.Called(request)

	if len(ret) == 0 {
		panic("no return value specified for OnGetLocalListVersion")
	}

	var r0 *localauth.GetLocalListVersionConfirmation
	var r1 error
	if rf, ok := ret.Get(0).(func(*localauth.GetLocalListVersionRequest) (*localauth.GetLocalListVersionConfirmation, error)); ok {
		return rf(request)
	}
	if rf, ok := ret.Get(0).(func(*localauth.GetLocalListVersionRequest) *localauth.GetLocalListVersionConfirmation); ok {
		r0 = rf(request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*localauth.GetLocalListVersionConfirmation)
		}
	}

	if rf, ok := ret.Get(1).(func(*localauth.GetLocalListVersionRequest) error); ok {
		r1 = rf(request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockLocalAuthListChargePointHandler_OnGetLocalListVersion_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'OnGetLocalListVersion'
type MockLocalAuthListChargePointHandler_OnGetLocalListVersion_Call struct {
	*mock.Call
}

// OnGetLocalListVersion is a helper method to define mock.On call
//   - request *localauth.GetLocalListVersionRequest
func (_e *MockLocalAuthListChargePointHandler_Expecter) OnGetLocalListVersion(request interface{}) *MockLocalAuthListChargePointHandler_OnGetLocalListVersion_Call {
	return &MockLocalAuthListChargePointHandler_OnGetLocalListVersion_Call{Call: _e.mock.On("OnGetLocalListVersion", request)}
}

func (_c *MockLocalAuthListChargePointHandler_OnGetLocalListVersion_Call) Run(run func(request *localauth.GetLocalListVersionRequest)) *MockLocalAuthListChargePointHandler_OnGetLocalListVersion_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*localauth.GetLocalListVersionRequest))
	})
	return _c
}

func (_c *MockLocalAuthListChargePointHandler_OnGetLocalListVersion_Call) Return(confirmation *localauth.GetLocalListVersionConfirmation, err error) *MockLocalAuthListChargePointHandler_OnGetLocalListVersion_Call {
	_c.Call.Return(confirmation, err)
	return _c
}

func (_c *MockLocalAuthListChargePointHandler_OnGetLocalListVersion_Call) RunAndReturn(run func(*localauth.GetLocalListVersionRequest) (*localauth.GetLocalListVersionConfirmation, error)) *MockLocalAuthListChargePointHandler_OnGetLocalListVersion_Call {
	_c.Call.Return(run)
	return _c
}

// OnSendLocalList provides a mock function with given fields: request
func (_m *MockLocalAuthListChargePointHandler) OnSendLocalList(request *localauth.SendLocalListRequest) (*localauth.SendLocalListConfirmation, error) {
	ret := _m.Called(request)

	if len(ret) == 0 {
		panic("no return value specified for OnSendLocalList")
	}

	var r0 *localauth.SendLocalListConfirmation
	var r1 error
	if rf, ok := ret.Get(0).(func(*localauth.SendLocalListRequest) (*localauth.SendLocalListConfirmation, error)); ok {
		return rf(request)
	}
	if rf, ok := ret.Get(0).(func(*localauth.SendLocalListRequest) *localauth.SendLocalListConfirmation); ok {
		r0 = rf(request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*localauth.SendLocalListConfirmation)
		}
	}

	if rf, ok := ret.Get(1).(func(*localauth.SendLocalListRequest) error); ok {
		r1 = rf(request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockLocalAuthListChargePointHandler_OnSendLocalList_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'OnSendLocalList'
type MockLocalAuthListChargePointHandler_OnSendLocalList_Call struct {
	*mock.Call
}

// OnSendLocalList is a helper method to define mock.On call
//   - request *localauth.SendLocalListRequest
func (_e *MockLocalAuthListChargePointHandler_Expecter) OnSendLocalList(request interface{}) *MockLocalAuthListChargePointHandler_OnSendLocalList_Call {
	return &MockLocalAuthListChargePointHandler_OnSendLocalList_Call{Call: _e.mock.On("OnSendLocalList", request)}
}

func (_c *MockLocalAuthListChargePointHandler_OnSendLocalList_Call) Run(run func(request *localauth.SendLocalListRequest)) *MockLocalAuthListChargePointHandler_OnSendLocalList_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*localauth.SendLocalListRequest))
	})
	return _c
}

func (_c *MockLocalAuthListChargePointHandler_OnSendLocalList_Call) Return(confirmation *localauth.SendLocalListConfirmation, err error) *MockLocalAuthListChargePointHandler_OnSendLocalList_Call {
	_c.Call.Return(confirmation, err)
	return _c
}

func (_c *MockLocalAuthListChargePointHandler_OnSendLocalList_Call) RunAndReturn(run func(*localauth.SendLocalListRequest) (*localauth.SendLocalListConfirmation, error)) *MockLocalAuthListChargePointHandler_OnSendLocalList_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockLocalAuthListChargePointHandler creates a new instance of MockLocalAuthListChargePointHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockLocalAuthListChargePointHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockLocalAuthListChargePointHandler {
	mock := &MockLocalAuthListChargePointHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}