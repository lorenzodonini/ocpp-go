package smartcharging

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
	"gopkg.in/go-playground/validator.v9"
)

// -------------------- Set Charging Profile (CSMS -> CS) --------------------

const SetChargingProfileFeatureName = "SetChargingProfile"

// Status reported in SetChargingProfileResponse, indicating whether the Charging Station processed
// the message successfully. This does not guarantee the schedule will be followed to the letter.
type ChargingProfileStatus string

const (
	ChargingProfileStatusAccepted ChargingProfileStatus = "Accepted"
	ChargingProfileStatusRejected ChargingProfileStatus = "Rejected"
)

func isValidChargingProfileStatus(fl validator.FieldLevel) bool {
	status := ChargingProfileStatus(fl.Field().String())
	switch status {
	case ChargingProfileStatusAccepted, ChargingProfileStatusRejected:
		return true
	default:
		return false
	}
}

// The field definition of the SetChargingProfile request payload sent by the CSMS to the Charging Station.
type SetChargingProfileRequest struct {
	EvseID          int                    `json:"evseId" validate:"gte=0"`
	ChargingProfile *types.ChargingProfile `json:"chargingProfile" validate:"required"`
}

// This field definition of the SetChargingProfile response payload, sent by the Charging Station to the CSMS in response to a SetChargingProfileRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type SetChargingProfileResponse struct {
	Status     ChargingProfileStatus `json:"status" validate:"required,chargingProfileStatus21"`
	StatusInfo *types.StatusInfo     `json:"statusInfo,omitempty" validate:"omitempty"`
}

// The CSMS may influence the charging power or current drawn from a specific EVSE or
// the entire Charging Station, over a period of time.
// For this purpose, the CSMS calculates a ChargingSchedule to stay within certain limits, then sends a
// SetChargingProfileRequest to the Charging Station. The charging schedule limits may be imposed by any
// external system. The Charging Station responds to this request with a SetChargingProfileResponse.
//
// While charging, the EVSE will continuously adapt the maximum current/power according to the installed
// charging profiles.
type SetChargingProfileFeature struct{}

func (f SetChargingProfileFeature) GetFeatureName() string {
	return SetChargingProfileFeatureName
}

func (f SetChargingProfileFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(SetChargingProfileRequest{})
}

func (f SetChargingProfileFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(SetChargingProfileResponse{})
}

func (r SetChargingProfileRequest) GetFeatureName() string {
	return SetChargingProfileFeatureName
}

func (c SetChargingProfileResponse) GetFeatureName() string {
	return SetChargingProfileFeatureName
}

// Creates a new SetChargingProfileRequest, containing all required fields. There are no optional fields for this message.
func NewSetChargingProfileRequest(evseID int, chargingProfile *types.ChargingProfile) *SetChargingProfileRequest {
	return &SetChargingProfileRequest{
		EvseID:          evseID,
		ChargingProfile: chargingProfile,
	}
}

// Creates a new SetChargingProfileResponse, containing all required fields. Optional fields may be set afterwards.
func NewSetChargingProfileResponse(status ChargingProfileStatus) *SetChargingProfileResponse {
	return &SetChargingProfileResponse{Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("chargingProfileStatus21", isValidChargingProfileStatus)
}
