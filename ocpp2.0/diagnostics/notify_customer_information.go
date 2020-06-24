package diagnostics

import (
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
	"reflect"
)

// -------------------- Notify Customer Information (CS -> CSMS) --------------------

const NotifyCustomerInformationFeatureName = "NotifyCustomerInformation"

// The field definition of the NotifyCustomerInformation request payload sent by a Charging Station to the CSMS.
type NotifyCustomerInformationRequest struct {
	Data        string         `json:"data" validate:"required,max=512"`   // (Part of) the requested data. No format specified in which the data is returned. Should be human readable.
	Tbc         bool           `json:"tbc,omitempty" validate:"omitempty"` // “to be continued” indicator. Indicates whether another part of the monitoringData follows in an upcoming notifyMonitoringReportRequest message. Default value when omitted is false.
	SeqNo       int            `json:"seqNo" validate:"gte=0"`             // Sequence number of this message. First message starts at 0.
	GeneratedAt types.DateTime `json:"generatedAt" validate:"required"`    // Timestamp of the moment this message was generated at the Charging Station.
	RequestID   int            `json:"requestId" validate:"gte=0"`         // The Id of the request.
}

// This field definition of the NotifyCustomerInformation response payload, sent by the CSMS to the Charging Station in response to a NotifyCustomerInformationRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type NotifyCustomerInformationResponse struct {
}

// The CSMS may send a message to the Charging Station to retrieve raw customer information, for example to be compliant with local privacy laws.
// The Charging Station notifies the CSMS by sending one or more reports.
// For each report, the Charging station shall send a NotifyCustomerInformationRequest to the CSMS.
//
// The CSMS responds with a NotifyCustomerInformationResponse message to the Charging Station for each received NotifyCustomerInformationRequest.
type NotifyCustomerInformationFeature struct{}

func (f NotifyCustomerInformationFeature) GetFeatureName() string {
	return NotifyCustomerInformationFeatureName
}

func (f NotifyCustomerInformationFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(NotifyCustomerInformationRequest{})
}

func (f NotifyCustomerInformationFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(NotifyCustomerInformationResponse{})
}

func (r NotifyCustomerInformationRequest) GetFeatureName() string {
	return NotifyCustomerInformationFeatureName
}

func (c NotifyCustomerInformationResponse) GetFeatureName() string {
	return NotifyCustomerInformationFeatureName
}

// Creates a new NotifyCustomerInformationRequest, containing all required fields. Optional fields may be set afterwards.
func NewNotifyCustomerInformationRequest(Data string, seqNo int, generatedAt types.DateTime, requestID int) *NotifyCustomerInformationRequest {
	return &NotifyCustomerInformationRequest{Data: Data, SeqNo: seqNo, GeneratedAt: generatedAt, RequestID: requestID}
}

// Creates a new NotifyCustomerInformationResponse, which doesn't contain any required or optional fields.
func NewNotifyCustomerInformationResponse() *NotifyCustomerInformationResponse {
	return &NotifyCustomerInformationResponse{}
}
