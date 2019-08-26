package ocpp16

import (
	"reflect"
)

// -------------------- Authorize (CP -> CS) --------------------

// The field definition of the Authorize request payload sent by the Charge Point to the Central System.
type AuthorizeRequest struct {
	IdTag string `json:"idTag" validate:"required,max=20"`
}

// This field definition of the Authorize confirmation payload, sent by the Charge Point to the Central System in response to an AuthorizeRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type AuthorizeConfirmation struct {
	IdTagInfo *IdTagInfo `json:"idTagInfo" validate:"required"`
}

// Before the owner of an electric vehicle can start or stop charging, the Charge Point has to authorize the operation.
// The Charge Point SHALL only supply energy after authorization.
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
func NewAuthorizationRequest(idTag string) *AuthorizeRequest {
	return &AuthorizeRequest{IdTag: idTag}
}

// Creates a new AuthorizeConfirmation. There are no optional fields for this message.
func NewAuthorizationConfirmation(idTagInfo *IdTagInfo) *AuthorizeConfirmation {
	return &AuthorizeConfirmation{IdTagInfo: idTagInfo}
}
