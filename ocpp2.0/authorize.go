package ocpp2

import (
	"reflect"
)

// -------------------- Authorize (CS -> CSMS) --------------------

// The field definition of the Authorize request payload sent by the Charging Station to the CSMS.
type AuthorizeRequest struct {
	EvseID              []int                 `json:"evseId,omitempty"`
	IdToken             IdToken               `json:"idToken" validate:"required"`
	CertificateHashData []OCSPRequestDataType `json:"15118CertificateHashData,omitempty" validate:"max=4"`
}

// This field definition of the Authorize confirmation payload, sent by the Charging Station to the CSMS in response to an AuthorizeRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type AuthorizeConfirmation struct {
	CertificateStatus CertificateStatus `json:"certificateStatus,omitempty" validate:"omitempty,certificateStatus"`
	EvseID            []int             `json:"evseId,omitempty"`
	IdTokenInfo       IdTokenInfo      `json:"idTokenInfo" validate:"required"`
}

// Before the owner of an electric vehicle can start or stop charging, the Charging Station has to authorize the operation.
// Upon receipt of an AuthorizeRequest, the CSMS SHALL respond with an AuthorizeConfirmation.
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

func (f AuthorizeFeature) GetConfirmationType() reflect.Type {
	return reflect.TypeOf(AuthorizeConfirmation{})
}

func (r AuthorizeRequest) GetFeatureName() string {
	return AuthorizeFeatureName
}

func (c AuthorizeConfirmation) GetFeatureName() string {
	return AuthorizeFeatureName
}

// Creates a new AuthorizeRequest, containing all required fields. There are no optional fields for this message.
func NewAuthorizationRequest(idToken string, tokenType IdTokenType) *AuthorizeRequest {
	return &AuthorizeRequest{IdToken: IdToken{IdToken: idToken, Type: tokenType}}
}

// Creates a new AuthorizeConfirmation. There are no optional fields for this message.
func NewAuthorizationConfirmation(idTokenInfo IdTokenInfo) *AuthorizeConfirmation {
	return &AuthorizeConfirmation{IdTokenInfo: idTokenInfo}
}
