package security

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
)

// -------------------- Security Event Notification Status (CS -> CSMS) --------------------

const SecurityEventNotificationFeatureName = "SecurityEventNotification"

// The field definition of the SecurityEventNotification request payload sent by the Charging Station to the CSMS.
type SecurityEventNotificationRequest struct {
	Type      string          `json:"type" validate:"required,max=50"`                 // Type of the security event. This value should be taken from the Security events list.
	Timestamp *types.DateTime `json:"timestamp" validate:"required"`                   // Date and time at which the event occurred.
	TechInfo  string          `json:"techInfo,omitempty" validate:"omitempty,max=255"` // Additional information about the occurred security event.
}

// This field definition of the SecurityEventNotification response payload, sent by the CSMS to the Charging Station in response to a SecurityEventNotificationRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type SecurityEventNotificationResponse struct {
}

// In case of critical security events, a Charging Station may immediately inform the CSMS of such events,
// via a SecurityEventNotificationRequest.
// The CSMS responds with a SecurityEventNotificationResponse to the Charging Station.
type SecurityEventNotificationFeature struct{}

func (f SecurityEventNotificationFeature) GetFeatureName() string {
	return SecurityEventNotificationFeatureName
}

func (f SecurityEventNotificationFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(SecurityEventNotificationRequest{})
}

func (f SecurityEventNotificationFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(SecurityEventNotificationResponse{})
}

func (r SecurityEventNotificationRequest) GetFeatureName() string {
	return SecurityEventNotificationFeatureName
}

func (c SecurityEventNotificationResponse) GetFeatureName() string {
	return SecurityEventNotificationFeatureName
}

// Creates a new SecurityEventNotificationRequest, containing all required fields. Optional fields may be set afterwards.
func NewSecurityEventNotificationRequest(typ string, timestamp *types.DateTime) *SecurityEventNotificationRequest {
	return &SecurityEventNotificationRequest{Type: typ, Timestamp: timestamp}
}

// Creates a new SecurityEventNotificationResponse, which doesn't contain any required or optional fields.
func NewSecurityEventNotificationResponse() *SecurityEventNotificationResponse {
	return &SecurityEventNotificationResponse{}
}
