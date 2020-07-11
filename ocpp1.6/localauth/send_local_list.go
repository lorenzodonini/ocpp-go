package localauth

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Send Local List (CS -> CP) --------------------

const SendLocalListFeatureName = "SendLocalList"

type UpdateType string
type UpdateStatus string

const (
	UpdateTypeDifferential      UpdateType   = "Differential"
	UpdateTypeFull              UpdateType   = "Full"
	UpdateStatusAccepted        UpdateStatus = "Accepted"
	UpdateStatusFailed          UpdateStatus = "Failed"
	UpdateStatusNotSupported    UpdateStatus = "NotSupported"
	UpdateStatusVersionMismatch UpdateStatus = "VersionMismatch"
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

func isValidUpdateStatus(fl validator.FieldLevel) bool {
	status := UpdateStatus(fl.Field().String())
	switch status {
	case UpdateStatusAccepted, UpdateStatusFailed, UpdateStatusNotSupported, UpdateStatusVersionMismatch:
		return true
	default:
		return false
	}
}

type AuthorizationData struct {
	IdTag     string           `json:"idTag" validate:"required,max=20"`
	IdTagInfo *types.IdTagInfo `json:"idTagInfo,omitempty"` //TODO: validate required if update type is Full
}

// The field definition of the SendLocalList request payload sent by the Central System to the Charge Point.
// If no (empty) localAuthorizationList is given and the updateType is Full, all identifications are removed from the list.
//
// Requesting a Differential update without (empty) localAuthorizationList will have no effect on the list.
// All idTags in the localAuthorizationList MUST be unique, no duplicate values are allowed.
type SendLocalListRequest struct {
	ListVersion            int                 `json:"listVersion" validate:"gte=0"`
	LocalAuthorizationList []AuthorizationData `json:"localAuthorizationList,omitempty" validate:"omitempty,dive"`
	UpdateType             UpdateType          `json:"updateType" validate:"required,updateType"`
}

// This field definition of the SendLocalList confirmation payload, sent by the Charge Point to the Central System in response to a SendLocalListRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type SendLocalListConfirmation struct {
	Status UpdateStatus `json:"status" validate:"required,updateStatus"`
}

// Central System can send a Local Authorization List that a Charge Point can use for authorization of idTags.
// The list MAY be either a full list to replace the current list in the Charge Point or it MAY be a differential list
// with updates to be applied to the current list in the Charge Point.
// Upon receipt of a SendLocalListRequest the Charge Point SHALL respond with a SendLocalListConfirmation.
// The Central System SHALL send a SendLocalListRequest to send the list to a Charge Point.
// The request payload SHALL contain the type of update (full or differential) and the version number that the Charge Point MUST associate with the local authorization list after it has been updated.
// The response payload SHALL indicate whether the Charge Point has accepted the update of the local authorization list.
// If the status is Failed or VersionMismatch and the updateType was Differential, then Central System SHOULD retry sending the full local authorization list with updateType Full.
type SendLocalListFeature struct{}

func (f SendLocalListFeature) GetFeatureName() string {
	return SendLocalListFeatureName
}

func (f SendLocalListFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(SendLocalListRequest{})
}

func (f SendLocalListFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(SendLocalListConfirmation{})
}

func (r SendLocalListRequest) GetFeatureName() string {
	return SendLocalListFeatureName
}

func (c SendLocalListConfirmation) GetFeatureName() string {
	return SendLocalListFeatureName
}

// Creates a new SendLocalListRequest, containing all required field. Optional fields may be set afterwards.
func NewSendLocalListRequest(version int, updateType UpdateType) *SendLocalListRequest {
	return &SendLocalListRequest{ListVersion: version, UpdateType: updateType}
}

// Creates a new SendLocalListConfirmation, containing all required fields. There are no optional fields for this message.
func NewSendLocalListConfirmation(status UpdateStatus) *SendLocalListConfirmation {
	return &SendLocalListConfirmation{Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("updateStatus", isValidUpdateStatus)
	_ = types.Validate.RegisterValidation("updateType", isValidUpdateType)
	//TODO: validation for SendLocalListMaxLength
}
