package ocpp2

import (
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Certificate Signed (CSMS -> CS) --------------------

// Status returned in response to CertificateSignedRequest, that indicates whether certificate signing has been accepted or rejected.
type CertificateSignedStatus string

const (
	CertificateSignedStatusAccepted CertificateSignedStatus = "Accepted"
	CertificateSignedStatusRejected CertificateSignedStatus = "Rejected"
)

func isValidCertificateSignedStatus(fl validator.FieldLevel) bool {
	status := CertificateSignedStatus(fl.Field().String())
	switch status {
	case CertificateSignedStatusAccepted, CertificateSignedStatusRejected:
		return true
	default:
		return false
	}
}

// The field definition of the CertificateSignedRequest PDU sent by the CSMS to the Charging Station.
type CertificateSignedRequest struct {
	Cert              []string              `json:"cert" validate:"required,min=1,dive,max=800"`
	TypeOfCertificate CertificateSigningUse `json:"typeOfCertificate,omitempty" validate:"omitempty,certificateSigningUse"`
}

// The field definition of the CertificateSignedResponse payload sent by the Charging Station to the CSMS in response to a CertificateSignedRequest.
type CertificateSignedConfirmation struct {
	Status CertificateSignedStatus `json:"status" validate:"required,certificateSignedStatus"`
}

// During the a certificate update procedure, the CSMS sends a new certificate, signed by a CA, to the Charging Station with a CertificateSignedRequest.
// The Charging Station verifies the signed certificate, installs it locally and responds with a CertificateSignedResponse to the the CSMS with the status Accepted or Rejected.
type CertificateSignedFeature struct{}

func (f CertificateSignedFeature) GetFeatureName() string {
	return CertificateSignedFeatureName
}

func (f CertificateSignedFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(CertificateSignedRequest{})
}

func (f CertificateSignedFeature) GetConfirmationType() reflect.Type {
	return reflect.TypeOf(CertificateSignedConfirmation{})
}

func (r CertificateSignedRequest) GetFeatureName() string {
	return CertificateSignedFeatureName
}

func (c CertificateSignedConfirmation) GetFeatureName() string {
	return CertificateSignedFeatureName
}

// Creates a new CertificateSignedRequest, containing all required fields. Additional optional fields may be set afterwards.
func NewCertificateSignedRequest(certificate []string) *CertificateSignedRequest {
	return &CertificateSignedRequest{Cert: certificate}
}

// Creates a new ChangeAvailabilityConfirmation, containing all required fields. There are no optional fields for this message.
func NewCertificateSignedConfirmation(status CertificateSignedStatus) *CertificateSignedConfirmation {
	return &CertificateSignedConfirmation{Status: status}
}

func init() {
	_ = Validate.RegisterValidation("certificateSignedStatus", isValidCertificateSignedStatus)
}
