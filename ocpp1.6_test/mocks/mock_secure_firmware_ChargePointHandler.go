// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	securefirmware "github.com/lorenzodonini/ocpp-go/ocpp1.6/securefirmware"
	mock "github.com/stretchr/testify/mock"
)

// MockSecureFirmwareChargePointHandler is an autogenerated mock type for the ChargePointHandler type
type MockSecureFirmwareChargePointHandler struct {
	mock.Mock
}

type MockSecureFirmwareChargePointHandler_Expecter struct {
	mock *mock.Mock
}

func (_m *MockSecureFirmwareChargePointHandler) EXPECT() *MockSecureFirmwareChargePointHandler_Expecter {
	return &MockSecureFirmwareChargePointHandler_Expecter{mock: &_m.Mock}
}

// OnSignedUpdateFirmware provides a mock function with given fields: request
func (_m *MockSecureFirmwareChargePointHandler) OnSignedUpdateFirmware(request *securefirmware.SignedUpdateFirmwareRequest) (*securefirmware.SignedUpdateFirmwareResponse, error) {
	ret := _m.Called(request)

	if len(ret) == 0 {
		panic("no return value specified for OnSignedUpdateFirmware")
	}

	var r0 *securefirmware.SignedUpdateFirmwareResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(*securefirmware.SignedUpdateFirmwareRequest) (*securefirmware.SignedUpdateFirmwareResponse, error)); ok {
		return rf(request)
	}
	if rf, ok := ret.Get(0).(func(*securefirmware.SignedUpdateFirmwareRequest) *securefirmware.SignedUpdateFirmwareResponse); ok {
		r0 = rf(request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*securefirmware.SignedUpdateFirmwareResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(*securefirmware.SignedUpdateFirmwareRequest) error); ok {
		r1 = rf(request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockSecureFirmwareChargePointHandler_OnSignedUpdateFirmware_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'OnSignedUpdateFirmware'
type MockSecureFirmwareChargePointHandler_OnSignedUpdateFirmware_Call struct {
	*mock.Call
}

// OnSignedUpdateFirmware is a helper method to define mock.On call
//   - request *securefirmware.SignedUpdateFirmwareRequest
func (_e *MockSecureFirmwareChargePointHandler_Expecter) OnSignedUpdateFirmware(request interface{}) *MockSecureFirmwareChargePointHandler_OnSignedUpdateFirmware_Call {
	return &MockSecureFirmwareChargePointHandler_OnSignedUpdateFirmware_Call{Call: _e.mock.On("OnSignedUpdateFirmware", request)}
}

func (_c *MockSecureFirmwareChargePointHandler_OnSignedUpdateFirmware_Call) Run(run func(request *securefirmware.SignedUpdateFirmwareRequest)) *MockSecureFirmwareChargePointHandler_OnSignedUpdateFirmware_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*securefirmware.SignedUpdateFirmwareRequest))
	})
	return _c
}

func (_c *MockSecureFirmwareChargePointHandler_OnSignedUpdateFirmware_Call) Return(response *securefirmware.SignedUpdateFirmwareResponse, err error) *MockSecureFirmwareChargePointHandler_OnSignedUpdateFirmware_Call {
	_c.Call.Return(response, err)
	return _c
}

func (_c *MockSecureFirmwareChargePointHandler_OnSignedUpdateFirmware_Call) RunAndReturn(run func(*securefirmware.SignedUpdateFirmwareRequest) (*securefirmware.SignedUpdateFirmwareResponse, error)) *MockSecureFirmwareChargePointHandler_OnSignedUpdateFirmware_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockSecureFirmwareChargePointHandler creates a new instance of MockSecureFirmwareChargePointHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockSecureFirmwareChargePointHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockSecureFirmwareChargePointHandler {
	mock := &MockSecureFirmwareChargePointHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
