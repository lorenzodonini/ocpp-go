package der

import (
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"reflect"
)

// -------------------- NotifyDERAlarm (CS -> CSMS) --------------------

const NotifyDERAlarm = "NotifyDERAlarm"

// The field definition of the NotifyDERAlarmRequest request payload sent by the CSMS to the Charging Station.
type NotifyDERAlarmRequest struct {
	ControlType    DERControl     `json:"controlType" validate:"required"`
	GridEventFault GridEventFault `json:"gridEventFault,omitempty" validate:"omitempty,gridEventFault"`
	AlarmEnded     *bool          `json:"alarmEnded,omitempty" validate:"omitempty"`
	Timestamp      types.DateTime `json:"timestamp" validate:"required"`
	ExtraInfo      string         `json:"extraInfo,omitempty" validate:"omitempty,max=200"`
}

// This field definition of the NotifyDERAlarmResponse
type NotifyDERAlarmResponse struct {
}

type NotifyDERAlarmFeature struct{}

func (f NotifyDERAlarmFeature) GetFeatureName() string {
	return NotifyDERAlarm
}

func (f NotifyDERAlarmFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(NotifyDERAlarmRequest{})
}

func (f NotifyDERAlarmFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(NotifyDERAlarmResponse{})
}

func (r NotifyDERAlarmRequest) GetFeatureName() string {
	return NotifyDERAlarm
}

func (c NotifyDERAlarmResponse) GetFeatureName() string {
	return NotifyDERAlarm
}

// Creates a new NotifyDERAlarmRequest, containing all required fields. Optional fields may be set afterwards.
func NewNotifyDERAlarmRequest(controlType DERControl, timestamp types.DateTime) *NotifyDERAlarmRequest {
	return &NotifyDERAlarmRequest{
		ControlType: controlType,
		Timestamp:   timestamp,
	}
}

// Creates a new NewAFFRSignalResponse, containing all required fields. Optional fields may be set afterwards.
func NewNotifyDERAlarmResponse() *NotifyDERAlarmResponse {
	return &NotifyDERAlarmResponse{}
}
