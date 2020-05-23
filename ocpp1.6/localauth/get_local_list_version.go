package localauth

import (
	"reflect"
)

// -------------------- Get Local List Version (CS -> CP) --------------------

const GetLocalListVersionFeatureName = "GetLocalListVersion"

// The field definition of the GetLocalListVersion request payload sent by the Central System to the Charge Point.
type GetLocalListVersionRequest struct {
}

// This field definition of the GetLocalListVersion confirmation payload, sent by the Charge Point to the Central System in response to a GetLocalListVersionRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type GetLocalListVersionConfirmation struct {
	ListVersion int `json:"listVersion" validate:"gte=-1"`
}

// Central System can request a Charge Point for the version number of the Local Authorization List.
// The Central System SHALL send a GetLocalListVersionRequest to request this value.
// Upon receipt of a GetLocalListVersionRequest, the Charge Point SHALL respond with a GetLocalListVersionConfirmation.
// The response payload SHALL contain the version number of its Local Authorization List.
// A version number of 0 (zero) SHALL be used to indicate that the local authorization list is empty, and a version number of -1 SHALL be used to indicate that the Charge Point does not support Local Authorization Lists.
type GetLocalListVersionFeature struct{}

func (f GetLocalListVersionFeature) GetFeatureName() string {
	return GetLocalListVersionFeatureName
}

func (f GetLocalListVersionFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(GetLocalListVersionRequest{})
}

func (f GetLocalListVersionFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(GetLocalListVersionConfirmation{})
}

func (r GetLocalListVersionRequest) GetFeatureName() string {
	return GetLocalListVersionFeatureName
}

func (c GetLocalListVersionConfirmation) GetFeatureName() string {
	return GetLocalListVersionFeatureName
}

// Creates a new GetLocalListVersionRequest, which doesn't contain any required or optional fields.
func NewGetLocalListVersionRequest() *GetLocalListVersionRequest {
	return &GetLocalListVersionRequest{}
}

// Creates a new GetLocalListVersionConfirmation, containing all required fields. There are no optional fields for this message.
func NewGetLocalListVersionConfirmation(version int) *GetLocalListVersionConfirmation {
	return &GetLocalListVersionConfirmation{ListVersion: version}
}
