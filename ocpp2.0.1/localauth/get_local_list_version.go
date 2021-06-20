package localauth

import (
	"reflect"
)

// -------------------- Get Local List Version (CSMS -> CS) --------------------

const GetLocalListVersionFeatureName = "GetLocalListVersion"

// The field definition of the GetLocalListVersion request payload sent by the CSMS to the Charging Station.
type GetLocalListVersionRequest struct {
}

// This field definition of the GetLocalListVersion response payload, sent by the Charging Station to the CSMS in response to a GetLocalListVersionRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type GetLocalListVersionResponse struct {
	VersionNumber int `json:"versionNumber" validate:"gte=0"`
}

// The CSMS can request a Charging Station for the version number of the Local Authorization List by sending a GetLocalListVersionRequest.
// Upon receipt of the GetLocalListVersionRequest Charging Station responds with a GetLocalListVersionResponse containing the version number of its Local Authorization List.
// The Charging Station SHALL use a version number of 0 (zero) to indicate that the Local Authorization List is empty.
type GetLocalListVersionFeature struct{}

func (f GetLocalListVersionFeature) GetFeatureName() string {
	return GetLocalListVersionFeatureName
}

func (f GetLocalListVersionFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(GetLocalListVersionRequest{})
}

func (f GetLocalListVersionFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(GetLocalListVersionResponse{})
}

func (r GetLocalListVersionRequest) GetFeatureName() string {
	return GetLocalListVersionFeatureName
}

func (c GetLocalListVersionResponse) GetFeatureName() string {
	return GetLocalListVersionFeatureName
}

// Creates a new GetLocalListVersionRequest, which doesn't contain any required or optional fields.
func NewGetLocalListVersionRequest() *GetLocalListVersionRequest {
	return &GetLocalListVersionRequest{}
}

// Creates a new GetLocalListVersionResponse, containing all required fields. There are no optional fields for this message.
func NewGetLocalListVersionResponse(version int) *GetLocalListVersionResponse {
	return &GetLocalListVersionResponse{VersionNumber: version}
}
