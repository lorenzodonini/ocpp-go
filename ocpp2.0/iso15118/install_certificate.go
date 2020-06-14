package iso15118

import (
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
	"reflect"
)

// -------------------- Clear Display (CSMS -> CS) --------------------

const InstallCertificateFeatureName = "InstallCertificate"

// The field definition of the InstallCertificate request payload sent by the CSMS to the Charging Station.
type InstallCertificateRequest struct {
	CertificateType types.CertificateUse `json:"certificateType" validate:"required,certificateUse"`
	Certificate     string               `json:"certificate" validate:"required,max=800"`
}

// This field definition of the InstallCertificate response payload, sent by the Charging Station to the CSMS in response to a InstallCertificateRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type InstallCertificateResponse struct {
	Status types.CertificateStatus `json:"status" validate:"required,certificateStatus"`
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
func NewInstallCertificateResponse(status types.CertificateStatus) *InstallCertificateResponse {
	return &InstallCertificateResponse{Status: status}
}
