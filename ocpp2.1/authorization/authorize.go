package authorization

import (
	"reflect"

	"gopkg.in/go-playground/validator.v9"

	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
)

// -------------------- Authorize (CS -> CSMS) --------------------

const AuthorizeFeatureName = "Authorize"

// The Certificate status information.
type AuthorizeCertificateStatus string

const (
	CertificateStatusAccepted               AuthorizeCertificateStatus = "Accepted"
	CertificateStatusSignatureError         AuthorizeCertificateStatus = "SignatureError"
	CertificateStatusCertificateExpired     AuthorizeCertificateStatus = "CertificateExpired"
	CertificateStatusCertificateRevoked     AuthorizeCertificateStatus = "CertificateRevoked"
	CertificateStatusNoCertificateAvailable AuthorizeCertificateStatus = "NoCertificateAvailable"
	CertificateStatusCertChainError         AuthorizeCertificateStatus = "CertChainError"
	CertificateStatusContractCancelled      AuthorizeCertificateStatus = "ContractCancelled"
)

func isValidAuthorizeCertificateStatus(fl validator.FieldLevel) bool {
	status := AuthorizeCertificateStatus(fl.Field().String())
	switch status {
	case CertificateStatusAccepted, CertificateStatusCertChainError, CertificateStatusCertificateExpired, CertificateStatusSignatureError, CertificateStatusNoCertificateAvailable, CertificateStatusCertificateRevoked, CertificateStatusContractCancelled:
		return true
	default:
		return false
	}
}

// The field definition of the Authorize request payload sent by the Charging Station to the CSMS.
type AuthorizeRequest struct {
	Certificate         *string                     `json:"certificate,omitempty" validate:"max=10000"`
	IdToken             types.IdToken               `json:"idToken" validate:"required"`
	CertificateHashData []types.OCSPRequestDataType `json:"iso15118CertificateHashData,omitempty" validate:"max=4,dive"`
}

// This field definition of the Authorize response payload, sent by the Charging Station to the CSMS in response to an AuthorizeRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type AuthorizeResponse struct {
	CertificateStatus     AuthorizeCertificateStatus `json:"certificateStatus,omitempty" validate:"omitempty,authorizeCertificateStatus21"`
	AllowedEnergyTransfer []types.EnergyTransferMode `json:"allowedEnergyTransfer,omitempty" validate:"omitempty,energyTransferMode"`
	IdTokenInfo           types.IdTokenInfo          `json:"idTokenInfo" validate:"required"`
	Tariff                *types.Tariff              `json:"tariff,omitempty" validate:"omitempty,dive"`
}

// Before the owner of an electric vehicle can start or stop charging, the Charging Station has to authorize the operation.
// Upon receipt of an AuthorizeRequest, the CSMS SHALL respond with an AuthorizeResponse.
// This response payload SHALL indicate whether or not the idTag is accepted by the CSMS.
// If the CSMS accepts the idToken then the response payload MUST include an authorization status value indicating acceptance or a reason for rejection.
//
// A Charging Station MAY authorize identifier locally without involving the CSMS, as described in Local Authorization List.
//
// The Charging Station SHALL only supply energy after authorization.
type AuthorizeFeature struct{}

func (f AuthorizeFeature) GetFeatureName() string {
	return AuthorizeFeatureName
}

func (f AuthorizeFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(AuthorizeRequest{})
}

func (f AuthorizeFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(AuthorizeResponse{})
}

func (r AuthorizeRequest) GetFeatureName() string {
	return AuthorizeFeatureName
}

func (c AuthorizeResponse) GetFeatureName() string {
	return AuthorizeFeatureName
}

// Creates a new AuthorizeRequest, containing all required fields. There are no optional fields for this message.
func NewAuthorizationRequest(idToken string, tokenType types.IdTokenType) *AuthorizeRequest {
	return &AuthorizeRequest{IdToken: types.IdToken{IdToken: idToken, Type: tokenType}}
}

// Creates a new AuthorizeResponse. There are no optional fields for this message.
func NewAuthorizationResponse(idTokenInfo types.IdTokenInfo) *AuthorizeResponse {
	return &AuthorizeResponse{IdTokenInfo: idTokenInfo}
}

func init() {
	_ = types.Validate.RegisterValidation("authorizeCertificateStatus21", isValidAuthorizeCertificateStatus)
}
