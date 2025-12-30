package localauth

import (
	"reflect"

	"github.com/xBlaz3kx/ocpp-go/ocpp2.1/types"
	"gopkg.in/go-playground/validator.v9"
)

// -------------------- Send Local List (CSMS -> CS) --------------------

const SendLocalListFeatureName = "SendLocalList"

// Indicates the type of update (full or differential) for a SendLocalListRequest.
type UpdateType string

const (
	UpdateTypeDifferential UpdateType = "Differential" // Indicates that the current Local Authorization List must be updated with the values in this message.
	UpdateTypeFull         UpdateType = "Full"         // Indicates that the current Local Authorization List must be replaced by the values in this message.
)

func isValidUpdateType(fl validator.FieldLevel) bool {
	status := UpdateType(fl.Field().String())
	switch status {
	case UpdateTypeDifferential, UpdateTypeFull:
		return true
	default:
		return false
	}
}

// Indicates whether the Charging Station has successfully received and applied the update of the Local Authorization List.
type SendLocalListStatus string

const (
	SendLocalListStatusAccepted        SendLocalListStatus = "Accepted"        // Local Authorization List successfully updated.
	SendLocalListStatusFailed          SendLocalListStatus = "Failed"          // Failed to update the Local Authorization List.
	SendLocalListStatusVersionMismatch SendLocalListStatus = "VersionMismatch" // Version number in the request for a differential update is less or equal then version number of current list.
)

func isValidSendLocalListStatus(fl validator.FieldLevel) bool {
	status := SendLocalListStatus(fl.Field().String())
	switch status {
	case SendLocalListStatusAccepted, SendLocalListStatusFailed, SendLocalListStatusVersionMismatch:
		return true
	default:
		return false
	}
}

// Contains the identifier to use for authorization.
type AuthorizationData struct {
	IdTokenInfo *types.IdTokenInfo `json:"idTokenInfo,omitempty" validate:"omitempty"`
	IdToken     types.IdToken      `json:"idToken"`
}

// The field definition of the SendLocalList request payload sent by the CSMS to the Charging Station.
type SendLocalListRequest struct {
	VersionNumber          int                 `json:"versionNumber" validate:"gte=0"`
	UpdateType             UpdateType          `json:"updateType" validate:"required,updateType21"`
	LocalAuthorizationList []AuthorizationData `json:"localAuthorizationList,omitempty" validate:"omitempty,dive"`
}

// This field definition of the SendLocalList response payload, sent by the Charging Station to the CSMS in response to a SendLocalListRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type SendLocalListResponse struct {
	Status     SendLocalListStatus `json:"status" validate:"required,sendLocalListStatus21"`
	StatusInfo *types.StatusInfo   `json:"statusInfo,omitempty" validate:"omitempty"`
}

// Enables the CSMS to send a Local Authorization List which a Charging Station can use for the
// authorization of idTokens.
// The list MAY be either a full list to replace the current list in the Charging Station or it MAY
// be a differential list with updates to be applied to the current list in the Charging Station.
//
// To install or update a local authorization list, the CSMS sends a SendLocalListRequest to a
// Charging Station, which responds with a SendLocalListResponse, containing the status of the operation.
//
// If LocalAuthListEnabled is configured to false on a charging station, this operation will have no effect.
type SendLocalListFeature struct{}

func (f SendLocalListFeature) GetFeatureName() string {
	return SendLocalListFeatureName
}

func (f SendLocalListFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(SendLocalListRequest{})
}

func (f SendLocalListFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(SendLocalListResponse{})
}

func (r SendLocalListRequest) GetFeatureName() string {
	return SendLocalListFeatureName
}

func (c SendLocalListResponse) GetFeatureName() string {
	return SendLocalListFeatureName
}

// Creates a new SendLocalListRequest, which doesn't contain any required or optional fields.
func NewSendLocalListRequest(versionNumber int, updateType UpdateType) *SendLocalListRequest {
	return &SendLocalListRequest{VersionNumber: versionNumber, UpdateType: updateType}
}

// Creates a new SendLocalListResponse, containing all required fields. There are no optional fields for this message.
func NewSendLocalListResponse(status SendLocalListStatus) *SendLocalListResponse {
	return &SendLocalListResponse{Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("updateType21", isValidUpdateType)
	_ = types.Validate.RegisterValidation("sendLocalListStatus21", isValidSendLocalListStatus)
}
