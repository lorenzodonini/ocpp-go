package ocpp2

import (
	"reflect"
)

// -------------------- Firmware Status Notification (CS -> CSMS) --------------------

// Contains an X.509 certificate chain, each first DER encoded into binary, and then base64 encoded.
type CertificateChain struct {
	Certificate      string   `json:"certificate" validate:"required,max=800"`
	ChildCertificate []string `json:"childCertificate,omitempty" validate:"omitempty,max=4,dive,required,max=800"`
}

// The field definition of the Get15118EVCertificate request payload sent by the Charging Station to the CSMS.
type Get15118EVCertificateRequest struct {
	SchemaVersion string `json:"15118SchemaVersion" validate:"required,max=50"`
	ExiRequest    string `json:"exiRequest" validate:"required,max=5500"`
}

// This field definition of the Get15118EVCertificate confirmation payload, sent by the CSMS to the Charging Station in response to a Get15118EVCertificateRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type Get15118EVCertificateConfirmation struct {
	Status                            Certificate15118EVStatus `json:"status" validate:"required,15118EVCertificate"`
	ExiResponse                       string                   `json:"exiResponse" validate:"required,max=5500"`
	ContractSignatureCertificateChain CertificateChain         `json:"contractSignatureCertificateChain" validate:"required"`
	SaProvisioningCertificateChain    CertificateChain         `json:"saProvisioningCertificateChain" validate:"required"`
}

// An EV connected to a Charging Station may request a new certificate.
// The EV initiates installing a new certificate. The Charging Station forwards the request for a new certificate to the CSMS.
// The CSMS responds to Charging Station with a message containing the status and optionally new certificate.
type Get15118EVCertificateFeature struct{}

func (f Get15118EVCertificateFeature) GetFeatureName() string {
	return Get15118EVCertificateFeatureName
}

func (f Get15118EVCertificateFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(Get15118EVCertificateRequest{})
}

func (f Get15118EVCertificateFeature) GetConfirmationType() reflect.Type {
	return reflect.TypeOf(Get15118EVCertificateConfirmation{})
}

func (r Get15118EVCertificateRequest) GetFeatureName() string {
	return Get15118EVCertificateFeatureName
}

func (c Get15118EVCertificateConfirmation) GetFeatureName() string {
	return Get15118EVCertificateFeatureName
}

// Creates a new Get15118EVCertificateRequest, containing all required fields.
func NewGet15118EVCertificateRequest(schemaVersion string, exiRequest string) *Get15118EVCertificateRequest {
	return &Get15118EVCertificateRequest{SchemaVersion: schemaVersion, ExiRequest: exiRequest}
}

// Creates a new Get15118EVCertificateConfirmation, containing all required fields.
func NewGet15118EVCertificateConfirmation(status Certificate15118EVStatus, exiResponse string, contractSignatureCertificateChain CertificateChain, saProvisioningCertificateChain CertificateChain) *Get15118EVCertificateConfirmation {
	return &Get15118EVCertificateConfirmation{Status: status, ExiResponse: exiResponse, ContractSignatureCertificateChain: contractSignatureCertificateChain, SaProvisioningCertificateChain: saProvisioningCertificateChain}
}

func init() {
}
