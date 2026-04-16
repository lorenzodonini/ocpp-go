package authorization

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
)

// -------------------- Authorize (CS -> CSMS) --------------------

const AuthorizeFeatureName = "Authorize"

// AuthorizeRequest is the payload for an Authorize request from CS to CSMS.
type AuthorizeRequest struct {
	IdToken                     types.IdToken               `json:"idToken" validate:"required"`
	Certificate                 string                      `json:"certificate,omitempty" validate:"omitempty,max=10000"`
	ISO15118CertificateHashData []types.OCSPRequestDataType `json:"iso15118CertificateHashData,omitempty" validate:"omitempty,max=4,dive"`
	CustomData                  *types.CustomDataType       `json:"customData,omitempty" validate:"omitempty"`
}

// AuthorizeResponse is the payload for an Authorize response from CSMS to CS.
type AuthorizeResponse struct {
	CertificateStatus     types.AuthorizeCertificateStatus `json:"certificateStatus,omitempty" validate:"omitempty,authorizeCertificateStatus21"`
	IdTokenInfo           types.IdTokenInfo                `json:"idTokenInfo" validate:"required"`
	AllowedEnergyTransfer []types.EnergyTransferMode       `json:"allowedEnergyTransfer,omitempty" validate:"omitempty,dive,energyTransferMode21"`
	CustomData            *types.CustomDataType            `json:"customData,omitempty" validate:"omitempty"`
}

// AuthorizeFeature represents the Authorize feature.
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

// NewAuthorizeRequest creates a new AuthorizeRequest.
func NewAuthorizeRequest(idToken string, tokenType types.IdTokenType) *AuthorizeRequest {
	return &AuthorizeRequest{IdToken: types.IdToken{IdToken: idToken, Type: tokenType}}
}

// NewAuthorizeResponse creates a new AuthorizeResponse.
func NewAuthorizeResponse(idTokenInfo types.IdTokenInfo) *AuthorizeResponse {
	return &AuthorizeResponse{IdTokenInfo: idTokenInfo}
}
