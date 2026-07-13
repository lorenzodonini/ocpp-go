package tariffcost

import "reflect"

// -------------------- NotifyWebPaymentStarted (CS -> CSMS) --------------------

const NotifyWebPaymentStartedFeatureName = "NotifyWebPaymentStarted"

// The field definition of the NotifyWebPaymentStartedRequest request payload sent by the Charging Station to the CSMS.
type NotifyWebPaymentStartedRequest struct {
	EvseId  int `json:"evseId" validate:"required,gte=0"`
	Timeout int `json:"timeout" validate:"required"`
}

// This field definition of the NotifyWebPaymentStartedResponse response payload, sent by the CSMS to the Charging Station.
type NotifyWebPaymentStartedResponse struct {
}

type NotifyWebPaymentStartedFeature struct{}

func (f NotifyWebPaymentStartedFeature) GetFeatureName() string {
	return NotifyWebPaymentStartedFeatureName
}

func (f NotifyWebPaymentStartedFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(NotifyWebPaymentStartedRequest{})
}

func (f NotifyWebPaymentStartedFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(NotifyWebPaymentStartedResponse{})
}

func (r NotifyWebPaymentStartedRequest) GetFeatureName() string {
	return NotifyWebPaymentStartedFeatureName
}

func (c NotifyWebPaymentStartedResponse) GetFeatureName() string {
	return NotifyWebPaymentStartedFeatureName
}

// Creates a new NotifyWebPaymentStartedRequest, containing all required fields.
func NewNotifyWebPaymentStartedRequest(evseId int, timeout int) *NotifyWebPaymentStartedRequest {
	return &NotifyWebPaymentStartedRequest{EvseId: evseId, Timeout: timeout}
}

// Creates a new NotifyWebPaymentStartedResponse, containing all required fields.
func NewNotifyWebPaymentStartedResponse() *NotifyWebPaymentStartedResponse {
	return &NotifyWebPaymentStartedResponse{}
}