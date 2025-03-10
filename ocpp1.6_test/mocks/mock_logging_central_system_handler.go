// Code generated by mockery v2.51.0. DO NOT EDIT.

package mocks

import (
	logging "github.com/lorenzodonini/ocpp-go/ocpp1.6/logging"
	mock "github.com/stretchr/testify/mock"
)

// MockLogCentralSystemHandler is an autogenerated mock type for the CentralSystemHandler type
type MockLogCentralSystemHandler struct {
	mock.Mock
}

type MockLogCentralSystemHandler_Expecter struct {
	mock *mock.Mock
}

func (_m *MockLogCentralSystemHandler) EXPECT() *MockLogCentralSystemHandler_Expecter {
	return &MockLogCentralSystemHandler_Expecter{mock: &_m.Mock}
}

// OnLogStatusNotification provides a mock function with given fields: chargingStationID, request
func (_m *MockLogCentralSystemHandler) OnLogStatusNotification(chargingStationID string, request *logging.LogStatusNotificationRequest) (*logging.LogStatusNotificationResponse, error) {
	ret := _m.Called(chargingStationID, request)

	if len(ret) == 0 {
		panic("no return value specified for OnLogStatusNotification")
	}

	var r0 *logging.LogStatusNotificationResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(string, *logging.LogStatusNotificationRequest) (*logging.LogStatusNotificationResponse, error)); ok {
		return rf(chargingStationID, request)
	}
	if rf, ok := ret.Get(0).(func(string, *logging.LogStatusNotificationRequest) *logging.LogStatusNotificationResponse); ok {
		r0 = rf(chargingStationID, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*logging.LogStatusNotificationResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(string, *logging.LogStatusNotificationRequest) error); ok {
		r1 = rf(chargingStationID, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockLogCentralSystemHandler_OnLogStatusNotification_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'OnLogStatusNotification'
type MockLogCentralSystemHandler_OnLogStatusNotification_Call struct {
	*mock.Call
}

// OnLogStatusNotification is a helper method to define mock.On call
//   - chargingStationID string
//   - request *logging.LogStatusNotificationRequest
func (_e *MockLogCentralSystemHandler_Expecter) OnLogStatusNotification(chargingStationID interface{}, request interface{}) *MockLogCentralSystemHandler_OnLogStatusNotification_Call {
	return &MockLogCentralSystemHandler_OnLogStatusNotification_Call{Call: _e.mock.On("OnLogStatusNotification", chargingStationID, request)}
}

func (_c *MockLogCentralSystemHandler_OnLogStatusNotification_Call) Run(run func(chargingStationID string, request *logging.LogStatusNotificationRequest)) *MockLogCentralSystemHandler_OnLogStatusNotification_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(*logging.LogStatusNotificationRequest))
	})
	return _c
}

func (_c *MockLogCentralSystemHandler_OnLogStatusNotification_Call) Return(response *logging.LogStatusNotificationResponse, err error) *MockLogCentralSystemHandler_OnLogStatusNotification_Call {
	_c.Call.Return(response, err)
	return _c
}

func (_c *MockLogCentralSystemHandler_OnLogStatusNotification_Call) RunAndReturn(run func(string, *logging.LogStatusNotificationRequest) (*logging.LogStatusNotificationResponse, error)) *MockLogCentralSystemHandler_OnLogStatusNotification_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockLogCentralSystemHandler creates a new instance of MockLogCentralSystemHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockLogCentralSystemHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockLogCentralSystemHandler {
	mock := &MockLogCentralSystemHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
