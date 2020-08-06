package smartcharging

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"gopkg.in/go-playground/validator.v9"
)

// -------------------- Set Charging Profile (CS -> CP) --------------------

const SetChargingProfileFeatureName = "SetChargingProfile"

// Status reported in SetChargingProfileConfirmation.
type ChargingProfileStatus string

const (
	ChargingProfileStatusAccepted       ChargingProfileStatus = "Accepted"
	ChargingProfileStatusRejected       ChargingProfileStatus = "Rejected"
	ChargingProfileStatusNotImplemented ChargingProfileStatus = "NotImplemented"
)

func isValidChargingProfileStatus(fl validator.FieldLevel) bool {
	status := ChargingProfileStatus(fl.Field().String())
	switch status {
	case ChargingProfileStatusAccepted, ChargingProfileStatusRejected, ChargingProfileStatusNotImplemented:
		return true
	default:
		return false
	}
}

// The field definition of the SetChargingProfile request payload sent by the Central System to the Charge Point.
type SetChargingProfileRequest struct {
	ConnectorId     int                    `json:"connectorId" validate:"gte=0"`
	ChargingProfile *types.ChargingProfile `json:"csChargingProfiles" validate:"required"`
}

// This field definition of the SetChargingProfile confirmation payload, sent by the Charge Point to the Central System in response to a SetChargingProfileRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type SetChargingProfileConfirmation struct {
	Status ChargingProfileStatus `json:"status" validate:"required,chargingProfileStatus"`
}

// The Central System MAY send charging profiles to a Charge Point that are to be used as default charging profiles.
// Such charging profiles MAY be sent at any time. If a charging profile with the same chargingProfileId, or the same combination
// of stackLevel / ChargingProfilePurpose, exists on the Charge Point, the new charging profile SHALL replace the existing charging profile,
// otherwise it SHALL be added. The Charge Point SHALL then re-evaluate its collection of charge profiles to determine which charging profile will become active.
type SetChargingProfileFeature struct{}

func (f SetChargingProfileFeature) GetFeatureName() string {
	return SetChargingProfileFeatureName
}

func (f SetChargingProfileFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(SetChargingProfileRequest{})
}

func (f SetChargingProfileFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(SetChargingProfileConfirmation{})
}

func (r SetChargingProfileRequest) GetFeatureName() string {
	return SetChargingProfileFeatureName
}

func (c SetChargingProfileConfirmation) GetFeatureName() string {
	return SetChargingProfileFeatureName
}

// Creates a new SetChargingProfileRequest, containing all required fields. There are no optional fields for this message.
func NewSetChargingProfileRequest(connectorId int, chargingProfile *types.ChargingProfile) *SetChargingProfileRequest {
	return &SetChargingProfileRequest{ConnectorId: connectorId, ChargingProfile: chargingProfile}
}

// Creates a new SetChargingProfileConfirmation, containing all required fields. There are no optional fields for this message.
func NewSetChargingProfileConfirmation(status ChargingProfileStatus) *SetChargingProfileConfirmation {
	return &SetChargingProfileConfirmation{Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("chargingProfileStatus", isValidChargingProfileStatus)
}
