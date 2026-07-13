package smartcharging

import "reflect"

// -------------------- NotifyPriorityCharging (CS -> CSMS) --------------------

const NotifyPriorityChargingFeatureName = "NotifyPriorityCharging"

// The field definition of the NotifyPriorityChargingRequest request payload sent by the Charging Station to the CSMS.
type NotifyPriorityChargingRequest struct {
	TransactionId string `json:"transactionId" validate:"required,max=36"`
	Activated     bool   `json:"activated" validate:"required"`
}

// This field definition of the NotifyPriorityChargingResponse response payload, sent by the CSMS to the Charging Station.
type NotifyPriorityChargingResponse struct {
}

type NotifyPriorityChargingFeature struct{}

func (f NotifyPriorityChargingFeature) GetFeatureName() string {
	return NotifyPriorityChargingFeatureName
}

func (f NotifyPriorityChargingFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(NotifyPriorityChargingRequest{})
}

func (f NotifyPriorityChargingFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(NotifyPriorityChargingResponse{})
}

func (r NotifyPriorityChargingRequest) GetFeatureName() string {
	return NotifyPriorityChargingFeatureName
}

func (c NotifyPriorityChargingResponse) GetFeatureName() string {
	return NotifyPriorityChargingFeatureName
}

// Creates a new NotifyPriorityChargingRequest, containing all required fields.
func NewNotifyPriorityChargingRequest(transactionId string, activated bool) *NotifyPriorityChargingRequest {
	return &NotifyPriorityChargingRequest{TransactionId: transactionId, Activated: activated}
}

// Creates a new NotifyPriorityChargingResponse, containing all required fields.
func NewNotifyPriorityChargingResponse() *NotifyPriorityChargingResponse {
	return &NotifyPriorityChargingResponse{}
}