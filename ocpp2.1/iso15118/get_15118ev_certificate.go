package iso15118

import (
	"reflect"

	"gopkg.in/go-playground/validator.v9"

	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
)

// -------------------- Get 15118EV Certificate (CS -> CSMS) --------------------

const Get15118EVCertificateFeatureName = "Get15118EVCertificate"

// Defines whether certificate needs to be installed or updated.
type CertificateAction string

const (
	CertificateActionInstall CertificateAction = "Install"
	CertificateActionUpdate  CertificateAction = "Update"
)

func isValidCertificateAction(fl validator.FieldLevel) bool {
	status := CertificateAction(fl.Field().String())
	switch status {
	case CertificateActionInstall, CertificateActionUpdate:
		return true
	default:
		return false
	}
}

// The field definition of the Get15118EVCertificate request payload sent by the Charging Station to the CSMS.
type Get15118EVCertificateRequest struct {
	SchemaVersion                    string            `json:"iso15118SchemaVersion" validate:"required,max=50"`
	Action                           CertificateAction `json:"action" validate:"required,certificateAction21"`
	ExiRequest                       string            `json:"exiRequest" validate:"required,max=11000"`
	MaximumContractCertificateChains *int              `json:"maximumContractCertificateChains,omitempty" validate:"omitempty,min=0"`
	PrioritizedEMAIDs                []string          `json:"prioritizedEMAIDs,omitempty" validate:"omitempty,max=8"`
}

// This field definition of the Get15118EVCertificate response payload, sent by the CSMS to the Charging Station in response to a Get15118EVCertificateRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type Get15118EVCertificateResponse struct {
	Status             types.Certificate15118EVStatus `json:"status" validate:"required,15118EVCertificate21"`
	ExiResponse        string                         `json:"exiResponse" validate:"required,max=17000"` // Raw CertificateInstallationRes response for the EV, Base64 encoded.
	RemainingContracts *int                           `json:"remainingContracts,omitempty" validate:"omitempty,gte=0"`
	StatusInfo         *types.StatusInfo              `json:"statusInfo,omitempty" validate:"omitempty"`
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

func (f Get15118EVCertificateFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(Get15118EVCertificateResponse{})
}

func (r Get15118EVCertificateRequest) GetFeatureName() string {
	return Get15118EVCertificateFeatureName
}

func (c Get15118EVCertificateResponse) GetFeatureName() string {
	return Get15118EVCertificateFeatureName
}

// Creates a new Get15118EVCertificateRequest, containing all required fields. There are no optional fields for this message.
func NewGet15118EVCertificateRequest(schemaVersion string, action CertificateAction, exiRequest string) *Get15118EVCertificateRequest {
	return &Get15118EVCertificateRequest{SchemaVersion: schemaVersion, Action: action, ExiRequest: exiRequest}
}

// Creates a new Get15118EVCertificateResponse, containing all required fields.
func NewGet15118EVCertificateResponse(status types.Certificate15118EVStatus, exiResponse string) *Get15118EVCertificateResponse {
	return &Get15118EVCertificateResponse{Status: status, ExiResponse: exiResponse}
}

func init() {
	_ = types.Validate.RegisterValidation("certificateAction21", isValidCertificateAction)
}
