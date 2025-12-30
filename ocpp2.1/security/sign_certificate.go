package security

import (
	"reflect"

	"github.com/xBlaz3kx/ocpp-go/ocpp2.1/types"
)

// -------------------- Sign Certificate (CS -> CSMS) --------------------

const SignCertificateFeatureName = "SignCertificate"

// The field definition of the SignCertificate request payload sent by the Charging Station to the CSMS.
type SignCertificateRequest struct {
	CSR                 string                      `json:"csr" validate:"required,max=5500"`                                         // The Charging Station SHALL send the public key in form of a Certificate Signing Request (CSR) as described in RFC 2986 and then PEM encoded.
	CertificateType     types.CertificateSigningUse `json:"certificateType,omitempty" validate:"omitempty,certificateSigningUse21"`   // Indicates the type of certificate that is to be signed.
	RequestId           *int                        `json:"requestId,omitempty"`                                                      // RequestId to match this message with the CertificateSignedRequest
	HashRootCertificate types.CertificateHashData   `json:"hashRootCertificate,omitempty" validate:"omitempty,certificateHashData21"` // The hash of the root certificate to identify	the PKI to use.
}

// This field definition of the SignCertificate response payload, sent by the CSMS to the Charging Station in response to a SignCertificateRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type SignCertificateResponse struct {
	Status     types.GenericStatus `json:"status" validate:"required,genericStatus21"` // Specifies whether the CSMS can process the request.
	StatusInfo *types.StatusInfo   `json:"statusInfo,omitempty" validate:"omitempty"`  // Detailed status information.
}

// If a Charging Station detected, that its certificate is due to expire, it will generate a new public/private key pair,
// then send a SignCertificateRequest to the CSMS containing a valid Certificate Signing Request.
//
// The CSMS responds with a SignCertificateResponse and will then forward the CSR to a CA server.
// Once the CA has issues a valid certificate, the CSMS will send a CertificateSignedRequest to the
// charging station (asynchronously).
type SignCertificateFeature struct{}

func (f SignCertificateFeature) GetFeatureName() string {
	return SignCertificateFeatureName
}

func (f SignCertificateFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(SignCertificateRequest{})
}

func (f SignCertificateFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(SignCertificateResponse{})
}

func (r SignCertificateRequest) GetFeatureName() string {
	return SignCertificateFeatureName
}

func (c SignCertificateResponse) GetFeatureName() string {
	return SignCertificateFeatureName
}

// Creates a new SignCertificateRequest, containing all required fields. Optional fields may be set afterwards.
func NewSignCertificateRequest(csr string) *SignCertificateRequest {
	return &SignCertificateRequest{CSR: csr}
}

// Creates a new SignCertificateResponse, containing all required fields. Optional fields may be set afterwards.
func NewSignCertificateResponse(status types.GenericStatus) *SignCertificateResponse {
	return &SignCertificateResponse{Status: status}
}
