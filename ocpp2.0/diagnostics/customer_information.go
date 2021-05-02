package diagnostics

import (
	"reflect"

	"gopkg.in/go-playground/validator.v9"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
)

// -------------------- Customer Information (CSMS -> CS) --------------------

const CustomerInformationFeatureName = "CustomerInformation"

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
	RequestID           int                        `json:"requestId" validate:"gte=0"`
	Report              bool                       `json:"report"`
	Clear               bool                       `json:"clear"`
	CustomerIdentifier  string                     `json:"customerIdentifier,omitempty" validate:"max=64"`
	IdToken             *types.IdToken             `json:"idToken,omitempty" validate:"omitempty,dive"`
	CustomerCertificate *types.CertificateHashData `json:"customerCertificate,omitempty" validate:"omitempty,dive"`
}

// This field definition of the CustomerInformation response payload, sent by the Charging Station to the CSMS in response to a CustomerInformationRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type CustomerInformationResponse struct {
	Status     CustomerInformationStatus `json:"status" validate:"required,customerInformationStatus"`
	StatusInfo *types.StatusInfo         `json:"statusInfo,omitempty" validate:"omitempty"`
}

// CSMS can request a Charging Station to clear its Authorization Cache.
// The CSMS SHALL send a CustomerInformationRequest payload for clearing the Charging Station’s Authorization Cache.
// Upon receipt of a CustomerInformationRequest, the Charging Station SHALL respond with a CustomerInformationResponse payload.
// The response payload SHALL indicate whether the Charging Station was able to clear its Authorization Cache.
type CustomerInformationFeature struct{}

func (f CustomerInformationFeature) GetFeatureName() string {
	return CustomerInformationFeatureName
}

func (f CustomerInformationFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(CustomerInformationRequest{})
}

func (f CustomerInformationFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(CustomerInformationResponse{})
}

func (r CustomerInformationRequest) GetFeatureName() string {
	return CustomerInformationFeatureName
}

func (c CustomerInformationResponse) GetFeatureName() string {
	return CustomerInformationFeatureName
}

// Creates a new CustomerInformationRequest, containing all required fields. Additional optional fields may be set afterwards.
func NewCustomerInformationRequest(requestId int, report bool, clear bool) *CustomerInformationRequest {
	return &CustomerInformationRequest{RequestID: requestId, Report: report, Clear: clear}
}

// Creates a new CustomerInformationResponse, containing all required fields. Additional optional fields may be set afterwards.
func NewCustomerInformationResponse(status CustomerInformationStatus) *CustomerInformationResponse {
	return &CustomerInformationResponse{Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("customerInformationStatus", isValidCustomerInformationStatus)
}
