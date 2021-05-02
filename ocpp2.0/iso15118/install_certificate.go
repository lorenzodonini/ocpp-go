package iso15118

import (
	"reflect"

	"gopkg.in/go-playground/validator.v9"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
)

// -------------------- Clear Display (CSMS -> CS) --------------------

const InstallCertificateFeatureName = "InstallCertificate"

// Charging Station indicates if installation was successful.
type InstallCertificateStatus string

const (
	CertificateStatusAccepted InstallCertificateStatus = "Accepted"
	CertificateStatusRejected InstallCertificateStatus = "Rejected"
	CertificateStatusFailed   InstallCertificateStatus = "Failed"
)

func isValidInstallCertificateStatus(fl validator.FieldLevel) bool {
	status := InstallCertificateStatus(fl.Field().String())
	switch status {
	case CertificateStatusAccepted, CertificateStatusRejected, CertificateStatusFailed:
		return true
	default:
		return false
	}
}

// The field definition of the InstallCertificate request payload sent by the CSMS to the Charging Station.
type InstallCertificateRequest struct {
	CertificateType types.CertificateUse `json:"certificateType" validate:"required,certificateUse"` // Indicates the certificate type that is sent.
	Certificate     string               `json:"certificate" validate:"required,max=5500"`           // A PEM encoded X.509 certificate.
}

// This field definition of the InstallCertificate response payload, sent by the Charging Station to the CSMS in response to a InstallCertificateRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type InstallCertificateResponse struct {
	Status     InstallCertificateStatus `json:"status" validate:"required,installCertificateStatus"`
	StatusInfo *types.StatusInfo        `json:"statusInfo,omitempty" validate:"omitempty"`
}

// The CSMS requests the Charging Station to install a new certificate by sending an InstallCertificateRequest.
// The certificate may be a root CA certificate, a Sub-CA certificate for an eMobility Operator, Charging Station operator, or a V2G root certificate.
//
// The Charging Station responds with an InstallCertificateResponse.
type InstallCertificateFeature struct{}

func (f InstallCertificateFeature) GetFeatureName() string {
	return InstallCertificateFeatureName
}

func (f InstallCertificateFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(InstallCertificateRequest{})
}

func (f InstallCertificateFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(InstallCertificateResponse{})
}

func (r InstallCertificateRequest) GetFeatureName() string {
	return InstallCertificateFeatureName
}

func (c InstallCertificateResponse) GetFeatureName() string {
	return InstallCertificateFeatureName
}

// Creates a new InstallCertificateRequest, containing all required fields. There are no optional fields for this message.
func NewInstallCertificateRequest(certificateType types.CertificateUse, certificate string) *InstallCertificateRequest {
	return &InstallCertificateRequest{CertificateType: certificateType, Certificate: certificate}
}

// Creates a new InstallCertificateResponse, containing all required fields. There are no optional fields for this message.
func NewInstallCertificateResponse(status InstallCertificateStatus) *InstallCertificateResponse {
	return &InstallCertificateResponse{Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("installCertificateStatus", isValidInstallCertificateStatus)
}
