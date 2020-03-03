package ocpp2

import (
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Customer Information (CSMS -> CS) --------------------

// Status returned in response to CustomerInformationRequest.
type CustomerInformationStatus string

const (
	CustomerInformationStatusAccepted CustomerInformationStatus = "Accepted"
	CustomerInformationStatusRejected CustomerInformationStatus = "Rejected"
	CustomerInformationStatusInvalid  CustomerInformationStatus = "Invalid"
)

func isValidCustomerInformationStatus(fl validator.FieldLevel) bool {
	status := CustomerInformationStatus(fl.Field().String())
	switch status {
	case CustomerInformationStatusAccepted, CustomerInformationStatusRejected, CustomerInformationStatusInvalid:
		return true
	default:
		return false
	}
}

// The field definition of the CustomerInformation request payload sent by the CSMS to the Charging Station.
type CustomerInformationRequest struct {
	RequestID           int                  `json:"requestId" validate:"gte=0"`
	Report              bool                 `json:"report"`
	Clear               bool                 `json:"clear"`
	CustomerIdentifier  string               `json:"customerIdentifier,omitempty" validate:"max=64"`
	IdToken             *IdToken             `json:"idToken,omitempty" validate:"omitempty,dive"`
	CustomerCertificate *CertificateHashData `json:"customerCertificate,omitempty" validate:"omitempty,dive"`
}

// This field definition of the CustomerInformation confirmation payload, sent by the Charging Station to the CSMS in response to a CustomerInformationRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type CustomerInformationConfirmation struct {
	Status CustomerInformationStatus `json:"status" validate:"required,customerInformationStatus"`
}

// CSMS can request a Charging Station to clear its Authorization Cache.
// The CSMS SHALL send a CustomerInformationRequest payload for clearing the Charging Station’s Authorization Cache.
// Upon receipt of a CustomerInformationRequest, the Charging Station SHALL respond with a CustomerInformationConfirmation payload.
// The response payload SHALL indicate whether the Charging Station was able to clear its Authorization Cache.
type CustomerInformationFeature struct{}

func (f CustomerInformationFeature) GetFeatureName() string {
	return CustomerInformationFeatureName
}

func (f CustomerInformationFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(CustomerInformationRequest{})
}

func (f CustomerInformationFeature) GetConfirmationType() reflect.Type {
	return reflect.TypeOf(CustomerInformationConfirmation{})
}

func (r CustomerInformationRequest) GetFeatureName() string {
	return CustomerInformationFeatureName
}

func (c CustomerInformationConfirmation) GetFeatureName() string {
	return CustomerInformationFeatureName
}

// Creates a new CustomerInformationRequest, containing all required fields. Additional optional fields may be set afterwards.
func NewCustomerInformationRequest(requestId int, report bool, clear bool) *CustomerInformationRequest {
	return &CustomerInformationRequest{RequestID: requestId, Report: report, Clear: clear}
}

// Creates a new CustomerInformationConfirmation, containing all required fields. There are no optional fields for this message.
func NewCustomerInformationConfirmation(status CustomerInformationStatus) *CustomerInformationConfirmation {
	return &CustomerInformationConfirmation{Status: status}
}

func init() {
	_ = Validate.RegisterValidation("customerInformationStatus", isValidCustomerInformationStatus)
}
